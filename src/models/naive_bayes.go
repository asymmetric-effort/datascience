package models

import (
	"fmt"
	"math"
	"sort"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/independencies"
)

// NaiveBayes represents a Naive Bayes classifier — a Bayesian network with a
// star topology where the class variable is the root and all features are
// conditionally independent given the class.
type NaiveBayes struct {
	*BayesianNetwork
	classVariable string
	features      []string
}

// NewNaiveBayes creates a new NaiveBayes model with the given class variable
// and feature variables. It constructs a star topology DAG where classVariable
// is a parent of each feature.
func NewNaiveBayes(classVariable string, features []string) (*NaiveBayes, error) {
	if classVariable == "" {
		return nil, fmt.Errorf("models: classVariable must not be empty")
	}
	if len(features) == 0 {
		return nil, fmt.Errorf("models: features must not be empty")
	}

	// Check for duplicates.
	seen := map[string]bool{classVariable: true}
	for _, f := range features {
		if f == classVariable {
			return nil, fmt.Errorf("models: feature %q is the same as the class variable", f)
		}
		if seen[f] {
			return nil, fmt.Errorf("models: duplicate feature %q", f)
		}
		seen[f] = true
	}

	bn := NewBayesianNetwork()
	if err := bn.AddNode(classVariable); err != nil {
		return nil, err
	}
	for _, f := range features {
		if err := bn.AddNode(f); err != nil {
			return nil, err
		}
		if err := bn.AddEdge(classVariable, f); err != nil {
			return nil, err
		}
	}

	feats := make([]string, len(features))
	copy(feats, features)

	return &NaiveBayes{
		BayesianNetwork: bn,
		classVariable:   classVariable,
		features:        feats,
	}, nil
}

// ClassVariable returns the name of the class variable.
func (nb *NaiveBayes) ClassVariable() string {
	return nb.classVariable
}

// Features returns a copy of the feature variable names.
func (nb *NaiveBayes) Features() []string {
	f := make([]string, len(nb.features))
	copy(f, nb.features)
	return f
}

// Fit estimates parameters from data using maximum likelihood estimation (MLE).
// The DataFrame must contain columns for the class variable and all features.
// All values must be non-negative integers representing discrete state indices.
// Fit delegates to nbFitImpl with the default CPD creator, preserving the
// public API signature.
func (nb *NaiveBayes) Fit(data *tabgo.DataFrame) error {
	return nbFitImpl(nb, data, defaultCPDCreator)
}

// PredictProbability returns the posterior probability of each class for each
// row in data. The result is a slice of length data.Len(), where each element
// is a slice of class probabilities.
func (nb *NaiveBayes) PredictProbability(data *tabgo.DataFrame) ([][]float64, error) {
	if data == nil {
		return nil, fmt.Errorf("models: data must not be nil")
	}
	if err := nb.BayesianNetwork.CheckModel(); err != nil {
		return nil, fmt.Errorf("models: model is not valid: %w", err)
	}

	classCPD := nb.GetCPD(nb.classVariable)
	if classCPD == nil {
		return nil, fmt.Errorf("models: no CPD for class variable %q", nb.classVariable)
	}
	classCard := classCPD.VariableCard()
	classPrior := classCPD.ToFactor().Values().Data()

	nRows := data.Len()
	result := make([][]float64, nRows)

	// Pre-fetch feature values.
	featVals := make([][]int, len(nb.features))
	featCPDs := make([]*factors.TabularCPD, len(nb.features))
	for i, feat := range nb.features {
		featVals[i] = data.Column(feat).Int()
		featCPDs[i] = nb.GetCPD(feat)
		if featCPDs[i] == nil {
			return nil, fmt.Errorf("models: no CPD for feature %q", feat)
		}
	}

	for row := 0; row < nRows; row++ {
		posterior := make([]float64, classCard)

		for c := 0; c < classCard; c++ {
			logProb := math.Log(classPrior[c])

			for i := range nb.features {
				cpd := featCPDs[i]
				featCard := cpd.VariableCard()
				fv := featVals[i][row]
				if fv < 0 || fv >= featCard {
					return nil, fmt.Errorf("models: feature value %d out of range for %q (card %d)",
						fv, nb.features[i], featCard)
				}
				// CPD data layout: data[featState * numParentConfigs + parentConfig]
				// Parent config for single parent (class) with state c is just c.
				val := cpd.ToFactor().Values().Data()[fv*classCard+c]
				if val <= 0 {
					logProb = math.Inf(-1)
					break
				}
				logProb += math.Log(val)
			}
			posterior[c] = logProb
		}

		// Convert from log-space, using log-sum-exp for numerical stability.
		maxLog := math.Inf(-1)
		for _, lp := range posterior {
			if lp > maxLog {
				maxLog = lp
			}
		}

		sum := 0.0
		for i := range posterior {
			if math.IsInf(maxLog, -1) {
				posterior[i] = 0
			} else {
				posterior[i] = math.Exp(posterior[i] - maxLog)
			}
			sum += posterior[i]
		}
		if sum > 0 {
			for i := range posterior {
				posterior[i] /= sum
			}
		}

		result[row] = posterior
	}

	return result, nil
}

// Predict returns the predicted class (index of highest posterior probability)
// for each row in data.
func (nb *NaiveBayes) Predict(data *tabgo.DataFrame) ([]int, error) {
	probs, err := nb.PredictProbability(data)
	if err != nil {
		return nil, err
	}

	predictions := make([]int, len(probs))
	for i, p := range probs {
		bestClass := 0
		bestProb := p[0]
		for c := 1; c < len(p); c++ {
			if p[c] > bestProb {
				bestProb = p[c]
				bestClass = c
			}
		}
		predictions[i] = bestClass
	}

	return predictions, nil
}

// AddEdge overrides the embedded BayesianNetwork's AddEdge to enforce the
// star topology: only edges from the class variable to a feature are allowed.
func (nb *NaiveBayes) AddEdge(from, to string) error {
	if from != nb.classVariable {
		return fmt.Errorf("models: NaiveBayes only allows edges from the class variable %q, got from=%q", nb.classVariable, from)
	}
	// Check that 'to' is a known feature.
	found := false
	for _, f := range nb.features {
		if f == to {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("models: NaiveBayes only allows edges to feature variables, %q is not a feature", to)
	}
	return nb.BayesianNetwork.AddEdge(from, to)
}

// AddEdgesFrom adds multiple edges from the given parent to the given children.
// Each edge must satisfy the star topology constraint.
func (nb *NaiveBayes) AddEdgesFrom(from string, toList []string) error {
	for _, to := range toList {
		if err := nb.AddEdge(from, to); err != nil {
			return err
		}
	}
	return nil
}

// ActiveTrailNodes returns the set of nodes reachable from the given variable
// via active trails given the observed variables. In a naive Bayes structure
// (class -> features), the d-separation rules simplify:
//   - From the class variable: all features not in observed are reachable, plus
//     any observed feature (since it's a direct child). All features are reachable.
//   - From a feature: the class is reachable unless it is observed and there is
//     no other path. With evidence on a feature (v-structure at class), the class
//     and other features become reachable. Without evidence on the feature, the
//     class is reachable and through it all non-observed features.
//
// This implements a BFS on the augmented DAG following Bayes-ball rules.
func (nb *NaiveBayes) ActiveTrailNodes(variable string, observed map[string]bool) ([]string, error) {
	if !nb.BayesianNetwork.dag.HasNode(variable) {
		return nil, fmt.Errorf("models: variable %q not in the network", variable)
	}
	if observed == nil {
		observed = make(map[string]bool)
	}

	// Bayes-ball algorithm for the star topology.
	// Nodes are visited as (node, direction) where direction is "up" (toward parents)
	// or "down" (toward children).
	type visit struct {
		node string
		up   bool // true = going up (toward parent), false = going down (toward child)
	}

	reachable := make(map[string]bool)
	visited := make(map[visit]bool)
	queue := []visit{}

	// Start: schedule both directions from the source.
	queue = append(queue, visit{variable, true}, visit{variable, false})

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if visited[cur] {
			continue
		}
		visited[cur] = true

		node := cur.node
		isObserved := observed[node]

		if cur.up { // ball going up (from child to parent)
			if !isObserved {
				reachable[node] = true
				// Pass ball up to parents.
				for _, p := range nb.BayesianNetwork.Parents(node) {
					queue = append(queue, visit{p, true})
				}
				// Pass ball down to children.
				for _, c := range nb.BayesianNetwork.Children(node) {
					queue = append(queue, visit{c, false})
				}
			}
		} else { // ball going down (from parent to child)
			if !isObserved {
				reachable[node] = true
				// Continue down to children.
				for _, c := range nb.BayesianNetwork.Children(node) {
					queue = append(queue, visit{c, false})
				}
			} else {
				// Observed node blocks downward, but we can bounce up to parents
				// (v-structure activation).
				for _, p := range nb.BayesianNetwork.Parents(node) {
					queue = append(queue, visit{p, true})
				}
			}
		}
	}

	// Remove the source itself from the result.
	delete(reachable, variable)

	result := make([]string, 0, len(reachable))
	for n := range reachable {
		result = append(result, n)
	}
	sort.Strings(result)
	return result, nil
}

// LocalIndependencies returns the set of independence assertions implied by
// the naive Bayes structure. In a naive Bayes model, every pair of features
// is conditionally independent given the class variable:
//
//	(f_i _|_ f_j | classVariable) for all i != j
func (nb *NaiveBayes) LocalIndependencies() *independencies.Independencies {
	ind := independencies.NewIndependencies()
	given := []string{nb.classVariable}

	for i := 0; i < len(nb.features); i++ {
		// Each feature is independent of all other features given the class.
		var others []string
		for j := 0; j < len(nb.features); j++ {
			if i != j {
				others = append(others, nb.features[j])
			}
		}
		if len(others) > 0 {
			ind.Add(independencies.NewIndependenceAssertion(
				[]string{nb.features[i]},
				others,
				given,
			))
		}
	}

	return ind
}
