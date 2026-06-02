package inference

import (
	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// MinNeighborsOrder returns an elimination ordering for the given variables
// using the min-neighbors (min-degree) heuristic. At each step the variable
// that appears in the fewest remaining factors is chosen for elimination.
//
// This is a greedy heuristic that tends to produce smaller intermediate
// factors, which keeps variable elimination tractable for many practical
// networks.
func MinNeighborsOrder(factorList []*factors.DiscreteFactor, eliminateVars []string) []string {
	if len(eliminateVars) == 0 {
		return nil
	}

	// Build a mutable copy of the variable set to eliminate.
	remaining := make(map[string]bool, len(eliminateVars))
	for _, v := range eliminateVars {
		remaining[v] = true
	}

	// Build a mutable list of variable sets for each factor.
	factorVarSets := make([]map[string]bool, len(factorList))
	for i, f := range factorList {
		factorVarSets[i] = make(map[string]bool)
		for _, v := range f.Variables() {
			factorVarSets[i][v] = true
		}
	}

	order := make([]string, 0, len(eliminateVars))

	for len(remaining) > 0 {
		// Pick the variable appearing in the fewest factor sets.
		bestVar := ""
		bestCount := int(^uint(0) >> 1) // max int

		for v := range remaining {
			count := 0
			for _, vs := range factorVarSets {
				if vs[v] {
					count++
				}
			}
			if count < bestCount {
				bestCount = count
				bestVar = v
			}
		}

		order = append(order, bestVar)
		delete(remaining, bestVar)

		// Simulate elimination: merge all factor sets containing bestVar
		// into one combined set, and remove bestVar from it.
		var merged map[string]bool
		var keepIndices []int

		for i, vs := range factorVarSets {
			if vs[bestVar] {
				if merged == nil {
					merged = make(map[string]bool)
				}
				for v := range vs {
					merged[v] = true
				}
			} else {
				keepIndices = append(keepIndices, i)
			}
		}

		if merged != nil {
			delete(merged, bestVar)
			newFactorVarSets := make([]map[string]bool, 0, len(keepIndices)+1)
			for _, idx := range keepIndices {
				newFactorVarSets = append(newFactorVarSets, factorVarSets[idx])
			}
			newFactorVarSets = append(newFactorVarSets, merged)
			factorVarSets = newFactorVarSets
		}
	}

	return order
}
