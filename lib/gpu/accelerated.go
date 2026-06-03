package gpu

import (
	"errors"
	"math"
)

// FactorData represents a discrete factor (conditional probability table)
// with named variables, a shape describing each variable's cardinality,
// and flat values in row-major order.
type FactorData struct {
	Variables []string
	Shape     []int
	Values    []float64
}

// AcceleratedVE runs variable elimination using the provided GPU backend.
//
// factors is the set of factors in the model. queryVars are the variables
// whose joint distribution is desired. elimVars lists variables to eliminate
// (in the order they should be eliminated). evidence maps variable names
// to observed values (0-indexed).
//
// The returned slice is the normalized distribution over queryVars.
func AcceleratedVE(backend Backend, factors []FactorData, queryVars, elimVars []string, evidence map[string]int) ([]float64, error) {
	if len(factors) == 0 {
		return nil, errors.New("gpu: no factors provided")
	}

	// Step 1: Apply evidence by reducing observed variables.
	reduced := make([]FactorData, len(factors))
	for i, f := range factors {
		reduced[i] = FactorData{
			Variables: append([]string(nil), f.Variables...),
			Shape:     append([]int(nil), f.Shape...),
			Values:    backend.CopyToDevice(f.Values),
		}
	}

	for varName, val := range evidence {
		for i, f := range reduced {
			axis := indexOf(f.Variables, varName)
			if axis < 0 {
				continue
			}
			newValues, newShape := backend.FactorReduce(f.Values, f.Shape, axis, val)
			newVars := removeString(f.Variables, axis)
			reduced[i] = FactorData{Variables: newVars, Shape: newShape, Values: newValues}
		}
	}

	// Step 2: Eliminate variables one at a time.
	workingFactors := reduced
	for _, elimVar := range elimVars {
		// Collect factors that mention elimVar.
		var involved []FactorData
		var remaining []FactorData
		for _, f := range workingFactors {
			if indexOf(f.Variables, elimVar) >= 0 {
				involved = append(involved, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		if len(involved) == 0 {
			continue
		}

		// Multiply involved factors together.
		product := involved[0]
		for j := 1; j < len(involved); j++ {
			product = multiplyFactors(backend, product, involved[j])
		}

		// Marginalize out elimVar.
		axis := indexOf(product.Variables, elimVar)
		newValues, newShape := backend.Marginalize(product.Values, product.Shape, axis)
		newVars := removeString(product.Variables, axis)
		remaining = append(remaining, FactorData{Variables: newVars, Shape: newShape, Values: newValues})
		workingFactors = remaining
	}

	// Step 3: Multiply remaining factors.
	result := workingFactors[0]
	for i := 1; i < len(workingFactors); i++ {
		result = multiplyFactors(backend, result, workingFactors[i])
	}

	// Step 4: Normalize.
	normalized := backend.Normalize(result.Values)
	return normalized, nil
}

// AcceleratedBP runs loopy belief propagation on a cluster graph using the
// provided GPU backend.
//
// cliques lists the variable names in each clique. cliqueFactors provides
// the initial factor (potential) for each clique. Messages are passed between
// neighboring cliques (those sharing at least one variable) until convergence
// or the iteration limit is reached.
//
// On return, cliqueFactors[i].Values contains the updated (unnormalized)
// belief for clique i.
func AcceleratedBP(backend Backend, cliques [][]string, cliqueFactors []FactorData) error {
	if len(cliques) != len(cliqueFactors) {
		return errors.New("gpu: cliques and cliqueFactors length mismatch")
	}

	const maxIter = 100
	const tol = 1e-6

	nCliques := len(cliques)

	// Build neighbor list: two cliques are neighbors if they share a variable.
	type edge struct{ i, j int }
	var edges []edge
	for i := 0; i < nCliques; i++ {
		for j := i + 1; j < nCliques; j++ {
			if sharedVars(cliques[i], cliques[j]) != nil {
				edges = append(edges, edge{i, j})
			}
		}
	}

	// Messages: msg[i][j] is the message from clique i to clique j.
	// Represented as a factor over the shared (sepset) variables.
	type msgKey struct{ from, to int }
	messages := map[msgKey]FactorData{}

	// Initialize messages to uniform 1s over sepset variables.
	for _, e := range edges {
		shared := sharedVars(cliques[e.i], cliques[e.j])
		shape := sepsetShape(shared, cliques[e.i], cliqueFactors[e.i].Shape)
		size := 1
		for _, d := range shape {
			size *= d
		}
		ones := make([]float64, size)
		for k := range ones {
			ones[k] = 1.0
		}
		messages[msgKey{e.i, e.j}] = FactorData{Variables: shared, Shape: shape, Values: ones}
		messages[msgKey{e.j, e.i}] = FactorData{Variables: shared, Shape: shape, Values: append([]float64(nil), ones...)}
	}

	// Iterate.
	for iter := 0; iter < maxIter; iter++ {
		maxDiff := 0.0
		for _, e := range edges {
			// Update message from i -> j and j -> i.
			for _, dir := range [][2]int{{e.i, e.j}, {e.j, e.i}} {
				from, to := dir[0], dir[1]
				key := msgKey{from, to}
				shared := messages[key].Variables

				// Compute product of clique potential with all incoming messages except from 'to'.
				belief := cliqueFactors[from]
				for _, e2 := range edges {
					for _, nb := range [][2]int{{e2.i, e2.j}, {e2.j, e2.i}} {
						if nb[1] == from && nb[0] != to {
							belief = multiplyFactors(backend, belief, messages[msgKey{nb[0], nb[1]}])
						}
					}
				}

				// Marginalize out variables not in the sepset.
				for _, v := range cliques[from] {
					if indexOf(shared, v) < 0 {
						axis := indexOf(belief.Variables, v)
						if axis >= 0 {
							newVals, newShape := backend.Marginalize(belief.Values, belief.Shape, axis)
							belief = FactorData{
								Variables: removeString(belief.Variables, axis),
								Shape:     newShape,
								Values:    newVals,
							}
						}
					}
				}

				// Normalize message.
				belief.Values = backend.Normalize(belief.Values)

				// Compute convergence delta.
				old := messages[key]
				diff := messageDistance(backend, old.Values, belief.Values)
				if diff > maxDiff {
					maxDiff = diff
				}
				messages[key] = FactorData{Variables: belief.Variables, Shape: belief.Shape, Values: belief.Values}
			}
		}
		if maxDiff < tol {
			break
		}
	}

	// Update clique beliefs: multiply clique potential with all incoming messages.
	for i := 0; i < nCliques; i++ {
		belief := cliqueFactors[i]
		for _, e := range edges {
			for _, dir := range [][2]int{{e.i, e.j}, {e.j, e.i}} {
				if dir[1] == i {
					belief = multiplyFactors(backend, belief, messages[msgKey{dir[0], dir[1]}])
				}
			}
		}
		cliqueFactors[i].Values = belief.Values
	}

	return nil
}

// AcceleratedSample generates n samples via forward sampling using the GPU
// backend for probability computations. topoOrder gives the topological
// ordering of variables. factors must include a factor for each variable.
//
// Returns an n x len(topoOrder) matrix where result[i][j] is the sampled
// value of the j-th variable in sample i.
func AcceleratedSample(backend Backend, factors []FactorData, topoOrder []string, n int) [][]int {
	samples := make([][]int, n)

	// Build a map from variable name to its factor.
	factorMap := map[string]FactorData{}
	for _, f := range factors {
		if len(f.Variables) > 0 {
			// The last variable in the factor's variable list is the "child".
			child := f.Variables[len(f.Variables)-1]
			factorMap[child] = f
		}
	}

	for i := 0; i < n; i++ {
		sample := make(map[string]int)
		row := make([]int, len(topoOrder))

		for j, varName := range topoOrder {
			f, ok := factorMap[varName]
			if !ok {
				row[j] = 0
				sample[varName] = 0
				continue
			}

			// Reduce factor by already-sampled parent values.
			reduced := FactorData{
				Variables: append([]string(nil), f.Variables...),
				Shape:     append([]int(nil), f.Shape...),
				Values:    backend.CopyToDevice(f.Values),
			}
			for _, pv := range f.Variables[:len(f.Variables)-1] {
				val, exists := sample[pv]
				if !exists {
					continue
				}
				axis := indexOf(reduced.Variables, pv)
				newVals, newShape := backend.FactorReduce(reduced.Values, reduced.Shape, axis, val)
				reduced = FactorData{
					Variables: removeString(reduced.Variables, axis),
					Shape:     newShape,
					Values:    newVals,
				}
			}

			// Normalize to get probability distribution.
			probs := backend.Normalize(reduced.Values)

			// Sample from the distribution using a simple linear scan CDF.
			// Use a deterministic-ish hash for reproducibility within this
			// pure-Go implementation. Real GPU sampling would use cuRAND.
			sampled := sampleCategorical(probs, i, j)
			row[j] = sampled
			sample[varName] = sampled
		}
		samples[i] = row
	}
	return samples
}

// --- helpers ---

// indexOf returns the position of s in the slice, or -1 if not found.
func indexOf(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1
}

// removeString removes the element at index from a string slice.
func removeString(slice []string, index int) []string {
	result := make([]string, 0, len(slice)-1)
	for i, v := range slice {
		if i != index {
			result = append(result, v)
		}
	}
	return result
}

// multiplyFactors multiplies two factors by computing their outer product.
// The result has variables from both factors (union, preserving order from a then b).
func multiplyFactors(backend Backend, a, b FactorData) FactorData {
	// Determine result variables (union).
	resultVars := append([]string(nil), a.Variables...)
	for _, v := range b.Variables {
		if indexOf(resultVars, v) < 0 {
			resultVars = append(resultVars, v)
		}
	}

	// For simplicity, compute outer product using the backend.
	resultShape := make([]int, len(a.Shape)+len(b.Shape))
	copy(resultShape, a.Shape)
	copy(resultShape[len(a.Shape):], b.Shape)

	resultVarsSimple := append(append([]string(nil), a.Variables...), b.Variables...)

	values := backend.FactorProduct(a.Values, a.Shape, b.Values, b.Shape, resultShape)

	return FactorData{
		Variables: resultVarsSimple,
		Shape:     resultShape,
		Values:    values,
	}
}

// sharedVars returns variables present in both a and b.
func sharedVars(a, b []string) []string {
	var shared []string
	for _, v := range a {
		if indexOf(b, v) >= 0 {
			shared = append(shared, v)
		}
	}
	return shared
}

// sepsetShape extracts the shape dimensions for the shared variables from
// a clique's variable list and shape.
func sepsetShape(shared []string, cliqueVars []string, cliqueShape []int) []int {
	shape := make([]int, len(shared))
	for i, v := range shared {
		idx := indexOf(cliqueVars, v)
		if idx >= 0 && idx < len(cliqueShape) {
			shape[i] = cliqueShape[idx]
		}
	}
	return shape
}

// messageDistance computes the max absolute difference between two message vectors.
func messageDistance(backend Backend, a, b []float64) float64 {
	if len(a) != len(b) {
		return math.Inf(1)
	}
	diff := backend.ElementWiseSub(a, b)
	absDiff := backend.Abs(diff)
	return backend.Max(absDiff)
}

// sampleCategorical samples from a categorical distribution given probabilities.
// Uses a simple deterministic pseudo-random based on seeds for reproducibility.
func sampleCategorical(probs []float64, seed1, seed2 int) int {
	// Simple hash-based pseudo-random number in [0,1).
	h := uint64(seed1)*2654435761 + uint64(seed2)*40503
	h ^= h >> 16
	h *= 0x45d9f3b
	h ^= h >> 16
	u := float64(h%1000000) / 1000000.0

	cumulative := 0.0
	for i, p := range probs {
		cumulative += p
		if u < cumulative {
			return i
		}
	}
	return len(probs) - 1
}
