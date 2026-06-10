package learning

import (
	"math"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// BICScore returns a ScoreFunc that computes the Bayesian Information Criterion
// (BIC) local score for a variable given its parents and data. The BIC score
// penalizes model complexity: BIC = LL - (k/2) * ln(n), where LL is the
// log-likelihood, k is the number of free parameters, and n is the sample size.
//
// This provides a built-in scoring function for use with HillClimbSearch, GES,
// ExhaustiveSearch, and MMHC without requiring the structure_score package.
func BICScore() ScoreFunc {
	return func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		n := data.Len()
		if n == 0 {
			return 0
		}

		counts, parentCounts, card, _ := localCountTable(variable, parents, data)

		// Compute log-likelihood.
		ll := 0.0
		for pc, childCounts := range counts {
			total := parentCounts[pc]
			if total == 0 {
				continue
			}
			for _, count := range childCounts {
				if count > 0 {
					ll += float64(count) * math.Log(float64(count)/float64(total))
				}
			}
		}

		// Number of parent configurations.
		numPC := len(parentCounts)
		if numPC == 0 {
			numPC = 1
		}

		// Number of free parameters: (card - 1) * numParentConfigs.
		k := float64((card - 1) * numPC)

		// BIC = LL - (k/2) * ln(n).
		return ll - (k/2)*math.Log(float64(n))
	}
}

// K2Score returns a ScoreFunc that computes the K2 (Cooper-Herskovits) score
// for a variable given its parents and data. The K2 score is based on the
// marginal likelihood with a uniform Dirichlet prior (pseudo-count = 1).
func K2Score() ScoreFunc {
	return func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		counts, parentCounts, card, _ := localCountTable(variable, parents, data)

		score := 0.0
		for pc, childCounts := range counts {
			nj := parentCounts[pc]
			// log(Gamma(card) / Gamma(nj + card))
			score += lgammaScore(float64(card)) - lgammaScore(float64(nj+card))
			for _, nijk := range childCounts {
				// log(Gamma(nijk + 1))
				score += lgammaScore(float64(nijk + 1))
			}
			// Account for states with zero counts.
			zeroStates := card - len(childCounts)
			if zeroStates > 0 {
				score += float64(zeroStates) * lgammaScore(1.0)
			}
		}

		return score
	}
}

// BDeuScore returns a ScoreFunc that computes the BDeu (Bayesian Dirichlet
// equivalent uniform) score with the given equivalent sample size.
func BDeuScore(equivalentSampleSize float64) ScoreFunc {
	return func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		counts, parentCounts, card, _ := localCountTable(variable, parents, data)

		numPC := len(parentCounts)
		if numPC == 0 {
			numPC = 1
		}

		alpha := equivalentSampleSize / float64(numPC)
		alphaK := alpha / float64(card)

		score := 0.0
		for pc, childCounts := range counts {
			nj := parentCounts[pc]
			score += lgammaScore(alpha) - lgammaScore(float64(nj)+alpha)
			for _, nijk := range childCounts {
				score += lgammaScore(float64(nijk)+alphaK) - lgammaScore(alphaK)
			}
			zeroStates := card - len(childCounts)
			if zeroStates > 0 {
				// lgammaScore(0 + alphaK) - lgammaScore(alphaK) = 0
				// so zero-count states contribute nothing extra.
			}
		}

		return score
	}
}

// AICScore returns a ScoreFunc that computes the Akaike Information Criterion
// local score. AIC = LL - k, where k is the number of free parameters.
func AICScore() ScoreFunc {
	return func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		n := data.Len()
		if n == 0 {
			return 0
		}

		counts, parentCounts, card, _ := localCountTable(variable, parents, data)

		ll := 0.0
		for pc, childCounts := range counts {
			total := parentCounts[pc]
			if total == 0 {
				continue
			}
			for _, count := range childCounts {
				if count > 0 {
					ll += float64(count) * math.Log(float64(count)/float64(total))
				}
			}
		}

		numPC := len(parentCounts)
		if numPC == 0 {
			numPC = 1
		}
		k := float64((card - 1) * numPC)

		return ll - k
	}
}

// localCountTable builds a count table for a variable given its parents.
// Returns counts[parentConfigKey][variableValue] and parentCounts[parentConfigKey].
func localCountTable(variable string, parents []string, data *tabgo.DataFrame) (
	counts map[string]map[int]int,
	parentCounts map[string]int,
	card int,
	n int,
) {
	n = data.Len()
	varVals := data.Column(variable).Int()

	// Determine cardinality.
	card = 0
	for _, v := range varVals {
		if v+1 > card {
			card = v + 1
		}
	}
	if card < 1 {
		card = 1
	}

	// Get parent column data.
	parentData := make([][]int, len(parents))
	parentCards := make([]int, len(parents))
	for i, p := range parents {
		parentData[i] = data.Column(p).Int()
		pc := 0
		for _, v := range parentData[i] {
			if v+1 > pc {
				pc = v + 1
			}
		}
		if pc < 1 {
			pc = 1
		}
		parentCards[i] = pc
	}

	counts = make(map[string]map[int]int)
	parentCounts = make(map[string]int)

	for row := 0; row < n; row++ {
		// Build parent config key as integer string.
		pcIdx := 0
		for i := range parents {
			pcIdx = pcIdx*parentCards[i] + parentData[i][row]
		}
		key := intToString(pcIdx)

		if counts[key] == nil {
			counts[key] = make(map[int]int)
		}
		counts[key][varVals[row]]++
		parentCounts[key]++
	}

	return counts, parentCounts, card, n
}

// intToString converts an int to a string without importing strconv.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	buf := make([]byte, 0, 12)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// Reverse.
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

// lgammaScore computes log(Gamma(x)) using the standard library.
func lgammaScore(x float64) float64 {
	v, _ := math.Lgamma(x)
	return v
}
