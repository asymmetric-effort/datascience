package models

import (
	"fmt"
	"io"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// ---------------------------------------------------------------------------
// Dependency-injection interfaces and internal testable functions.
//
// These exist to make defensive error guards reachable in tests. The public
// API signatures are unchanged; each public method delegates to an internal
// function that accepts interfaces so that failing mocks can be injected.
// ---------------------------------------------------------------------------

// factorizer abstracts factor extraction to enable testing of defensive
// error paths in Predict, PredictProbability, and GetStateProbability.
type factorizer interface {
	ToMarkovFactors() ([]*factors.DiscreteFactor, error)
}

// veQuerier abstracts variable-elimination query operations to enable
// testing of defensive error paths in inference-dependent methods.
type veQuerier interface {
	query(factorList []*factors.DiscreteFactor, queryVars []string, evidence map[string]int) (*factors.DiscreteFactor, error)
	mapQuery(factorList []*factors.DiscreteFactor, queryVars []string, evidence map[string]int) (map[string]int, error)
}

// bifWriter abstracts io.Writer for writeBIF to enable testing of
// defensive write-error paths in Save.
type bifWriter interface {
	io.Writer
}

// topologicalSorter abstracts topological ordering to enable testing
// of defensive error paths in Simulate and sampleBN.
type topologicalSorter interface {
	TopologicalOrder() ([]string, error)
}

// cpdCreator abstracts CPD creation to enable testing of defensive
// error paths in FitUpdate, NaiveBayes.Fit, and Do.
type cpdCreator func(variable string, variableCard int, values [][]float64, evidence []string, evidenceCard []int) (*factors.TabularCPD, error)

// defaultCPDCreator is the production CPD creator that delegates to
// factors.NewTabularCPD.
func defaultCPDCreator(variable string, variableCard int, values [][]float64, evidence []string, evidenceCard []int) (*factors.TabularCPD, error) {
	return factors.NewTabularCPD(variable, variableCard, values, evidence, evidenceCard)
}

// defaultVEQuerier is the production variable-elimination querier.
type defaultVEQuerier struct{}

func (defaultVEQuerier) query(factorList []*factors.DiscreteFactor, queryVars []string, evidence map[string]int) (*factors.DiscreteFactor, error) {
	return veQuery(factorList, queryVars, evidence)
}

func (defaultVEQuerier) mapQuery(factorList []*factors.DiscreteFactor, queryVars []string, evidence map[string]int) (map[string]int, error) {
	return veMAP(factorList, queryVars, evidence)
}

// ---------------------------------------------------------------------------
// Testable implementations that accept injected dependencies.
// ---------------------------------------------------------------------------

// predictImpl is the testable implementation of Predict. It accepts a
// factorizer interface to allow injection of failing mocks in tests for
// defensive error path coverage.
func predictImpl(f factorizer, veq veQuerier, data *tabgo.DataFrame, nodes []string, colVals map[string][]any) (*tabgo.DataFrame, error) {
	markovFactors, err := f.ToMarkovFactors()
	if err != nil {
		return nil, err // defensive: tested via mock factorizer
	}

	nRows := data.Len()
	cols := data.Columns()

	resultCols := make(map[string][]any, len(cols))
	for _, c := range cols {
		resultCols[c] = make([]any, nRows)
	}

	for row := 0; row < nRows; row++ {
		evidence := make(map[string]int)
		var queryVars []string

		for _, c := range cols {
			val := colVals[c][row]
			if val == nil {
				queryVars = append(queryVars, c)
			} else {
				evidence[c] = toInt(val)
			}
		}

		if len(queryVars) == 0 {
			for _, c := range cols {
				resultCols[c][row] = colVals[c][row]
			}
			continue
		}

		mapAssignment, err := veq.mapQuery(markovFactors, queryVars, evidence)
		if err != nil {
			return nil, fmt.Errorf("models: Predict row %d: %w", row, err)
		}

		for _, c := range cols {
			if colVals[c][row] == nil {
				resultCols[c][row] = mapAssignment[c]
			} else {
				resultCols[c][row] = colVals[c][row]
			}
		}
	}

	seriesMap := make(map[string]*tabgo.Series, len(cols))
	for _, c := range cols {
		seriesMap[c] = tabgo.NewSeries(c, resultCols[c])
	}
	return tabgo.NewDataFrame(seriesMap), nil
}

// predictProbabilityImpl is the testable implementation of PredictProbability.
// It accepts a factorizer interface to allow injection of failing mocks in
// tests for defensive error path coverage.
func predictProbabilityImpl(f factorizer, veq veQuerier, data *tabgo.DataFrame, colVals map[string][]any) (map[string][]float64, error) {
	markovFactors, err := f.ToMarkovFactors()
	if err != nil {
		return nil, err // defensive: tested via mock factorizer
	}

	cols := data.Columns()
	nRows := data.Len()
	result := make(map[string][]float64)

	for row := 0; row < nRows; row++ {
		evidence := make(map[string]int)
		var queryVars []string

		for _, c := range cols {
			val := colVals[c][row]
			if val == nil {
				queryVars = append(queryVars, c)
			} else {
				evidence[c] = toInt(val)
			}
		}

		if len(queryVars) == 0 {
			continue
		}

		resultFactor, err := veq.query(markovFactors, queryVars, evidence)
		if err != nil {
			return nil, fmt.Errorf("models: PredictProbability row %d: %w", row, err)
		}

		for _, qv := range queryVars {
			otherVars := make([]string, 0)
			for _, v := range resultFactor.Variables() {
				if v != qv {
					otherVars = append(otherVars, v)
				}
			}
			var marginal *factors.DiscreteFactor
			if len(otherVars) > 0 {
				marginal, err = resultFactor.Marginalize(otherVars)
				if err != nil {
					return nil, fmt.Errorf("models: PredictProbability marginalize %q: %w", qv, err)
				}
			} else {
				marginal = resultFactor.Copy()
			}
			probs := marginal.Values().Data()
			result[qv] = append(result[qv], probs...)
		}
	}

	return result, nil
}

// getStateProbabilityImpl is the testable implementation of GetStateProbability.
// It accepts a factorizer interface to allow injection of failing mocks in
// tests for defensive error path coverage.
func getStateProbabilityImpl(f factorizer, veq veQuerier, states map[string]int, allNodes []string) (float64, error) {
	markovFactors, err := f.ToMarkovFactors()
	if err != nil {
		return 0, err // defensive: tested via mock factorizer
	}

	specified := make(map[string]bool, len(states))
	for v := range states {
		specified[v] = true
	}

	allSpecified := true
	for _, n := range allNodes {
		if !specified[n] {
			allSpecified = false
			break
		}
	}

	if allSpecified {
		product, err := factors.FactorProduct(markovFactors...)
		if err != nil {
			return 0, fmt.Errorf("models: GetStateProbability: %w", err)
		}
		return product.GetValue(states), nil
	}

	specifiedVars := make([]string, 0, len(states))
	for _, n := range allNodes {
		if specified[n] {
			specifiedVars = append(specifiedVars, n)
		}
	}

	resultFactor, err := veq.query(markovFactors, specifiedVars, nil)
	if err != nil {
		return 0, fmt.Errorf("models: GetStateProbability: %w", err)
	}

	return resultFactor.GetValue(states), nil
}

// writeBIFImpl is the testable implementation of writeBIF. It accepts a
// bifWriter interface to allow injection of failing writers in tests for
// defensive write-error path coverage.
func writeBIFImpl(w bifWriter, bn *BayesianNetwork) error {
	return bn.writeBIF(w)
}

// fitUpdateImpl is the testable implementation of FitUpdate. It accepts a
// cpdCreator function to allow injection of failing CPD creators in tests
// for defensive error path coverage.
func fitUpdateImpl(bn *BayesianNetwork, data *tabgo.DataFrame, nPrevSamples int, createCPD cpdCreator) error {
	if nPrevSamples < 0 {
		return fmt.Errorf("models: nPrevSamples must be non-negative, got %d", nPrevSamples)
	}

	nodes := bn.Nodes()
	nRows := data.Len()
	if nRows == 0 {
		return nil
	}

	colVals := make(map[string][]any, len(nodes))
	for _, node := range nodes {
		colVals[node] = data.Column(node).Values()
	}

	for _, node := range nodes {
		cpd := bn.cpds[node]
		if cpd == nil {
			return fmt.Errorf("models: node %q has no CPD for FitUpdate", node)
		}

		parents := bn.Parents(node)
		childCard := cpd.VariableCard()
		evidenceCard := cpd.EvidenceCard()

		numParentConfigs := 1
		for _, ec := range evidenceCard {
			numParentConfigs *= ec
		}

		counts := make([][]float64, childCard)
		for i := range counts {
			counts[i] = make([]float64, numParentConfigs)
		}
		parentConfigCounts := make([]float64, numParentConfigs)

		for row := 0; row < nRows; row++ {
			childVal := toInt(colVals[node][row])
			if childVal < 0 || childVal >= childCard {
				continue
			}

			pc := 0
			stride := 1
			valid := true
			for pi := len(parents) - 1; pi >= 0; pi-- {
				pVal := toInt(colVals[parents[pi]][row])
				if pVal < 0 || pVal >= evidenceCard[pi] {
					valid = false
					break
				}
				pc += pVal * stride
				stride *= evidenceCard[pi]
			}
			if !valid {
				continue
			}

			counts[childVal][pc]++
			parentConfigCounts[pc]++
		}

		oldFactor := cpd.ToFactor()
		oldData := oldFactor.Values().Data()

		newValues := make([][]float64, childCard)
		for cs := 0; cs < childCard; cs++ {
			newValues[cs] = make([]float64, numParentConfigs)
			for pc := 0; pc < numParentConfigs; pc++ {
				oldProb := oldData[cs*numParentConfigs+pc]
				newCount := counts[cs][pc]
				total := float64(nPrevSamples) + parentConfigCounts[pc]
				if total > 0 {
					newValues[cs][pc] = (float64(nPrevSamples)*oldProb + newCount) / total
				} else {
					newValues[cs][pc] = oldProb
				}
			}
		}

		newCPD, err := createCPD(node, childCard, newValues, parents, evidenceCard)
		if err != nil {
			return fmt.Errorf("models: FitUpdate CPD for %q: %w", node, err)
		}
		bn.cpds[node] = newCPD
	}

	return nil
}

// doImpl is the testable implementation of Do. It accepts a cpdCreator
// function to allow injection of failing CPD creators in tests for
// defensive error path coverage.
func doImpl(bn *BayesianNetwork, nodes map[string]int, createCPD cpdCreator) (*BayesianNetwork, error) {
	if len(nodes) == 0 {
		return bn.Copy(), nil
	}

	mutilated := bn.Copy()

	for node, state := range nodes {
		if !mutilated.dag.HasNode(node) {
			return nil, fmt.Errorf("models: do-intervention on unknown node %q", node)
		}

		cpd := mutilated.cpds[node]
		if cpd == nil {
			return nil, fmt.Errorf("models: node %q has no CPD for do-intervention", node)
		}

		card := cpd.VariableCard()
		if state < 0 || state >= card {
			return nil, fmt.Errorf("models: do-intervention state %d out of range for %q (card %d)", state, node, card)
		}

		parents := mutilated.Parents(node)
		for _, p := range parents {
			_ = mutilated.dag.RemoveEdge(p, node)
		}

		vals := make([][]float64, card)
		for i := 0; i < card; i++ {
			if i == state {
				vals[i] = []float64{1.0}
			} else {
				vals[i] = []float64{0.0}
			}
		}
		deltaCPD, err := createCPD(node, card, vals, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("models: Do failed to create delta CPD for %q: %w", node, err)
		}
		mutilated.cpds[node] = deltaCPD
	}

	return mutilated, nil
}
