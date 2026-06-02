package readwrite

import (
	"fmt"
	"io"
	"strings"

	"github.com/asymmetric-effort/pgmgo/src/models"
)

// WriteBIF serializes a BayesianNetwork to BIF (Bayesian Interchange Format).
// Variable blocks are written in the order returned by bn.Nodes(), followed by
// probability blocks in the same order.
func WriteBIF(w io.Writer, bn *models.BayesianNetwork) error {
	// Network header.
	if _, err := fmt.Fprintf(w, "network unknown {\n}\n"); err != nil {
		return fmt.Errorf("readwrite: write error: %w", err)
	}

	nodes := bn.Nodes()

	// Write variable blocks.
	for _, node := range nodes {
		states := bn.GetStates(node)
		if len(states) == 0 {
			return fmt.Errorf("readwrite: variable %q has no state names", node)
		}
		if _, err := fmt.Fprintf(w, "\nvariable %s {\n", node); err != nil {
			return fmt.Errorf("readwrite: write error: %w", err)
		}
		if _, err := fmt.Fprintf(w, "  type discrete [ %d ] { %s };\n",
			len(states), strings.Join(states, ", ")); err != nil {
			return fmt.Errorf("readwrite: write error: %w", err)
		}
		if _, err := fmt.Fprintf(w, "}\n"); err != nil {
			return fmt.Errorf("readwrite: write error: %w", err)
		}
	}

	// Write probability blocks.
	for _, node := range nodes {
		cpd := bn.GetCPD(node)
		if cpd == nil {
			return fmt.Errorf("readwrite: variable %q has no CPD", node)
		}

		evidence := cpd.Evidence()
		states := bn.GetStates(node)

		// Header.
		if len(evidence) == 0 {
			if _, err := fmt.Fprintf(w, "\nprobability ( %s ) {\n", node); err != nil {
				return fmt.Errorf("readwrite: write error: %w", err)
			}
		} else {
			if _, err := fmt.Fprintf(w, "\nprobability ( %s | %s ) {\n",
				node, strings.Join(evidence, ", ")); err != nil {
				return fmt.Errorf("readwrite: write error: %w", err)
			}
		}

		// Get flat data from the underlying factor.
		data := cpd.ToFactor().Values().Data()
		childCard := cpd.VariableCard()
		evidenceCard := cpd.EvidenceCard()

		numParentConfigs := 1
		for _, ec := range evidenceCard {
			numParentConfigs *= ec
		}

		if len(evidence) == 0 {
			// Unconditional: table val1, val2;
			var parts []string
			for cs := 0; cs < childCard; cs++ {
				parts = append(parts, bifFormatFloat(data[cs]))
			}
			if _, err := fmt.Fprintf(w, "  table %s;\n", strings.Join(parts, ", ")); err != nil {
				return fmt.Errorf("readwrite: write error: %w", err)
			}
		} else {
			// Conditional: (State1, State2) val1, val2;
			// Data is stored as: flat[childState * numParentConfigs + parentConfig]
			for pc := 0; pc < numParentConfigs; pc++ {
				// Decompose parent config into state indices.
				parentStateNames := bifDecomposeParentConfig(pc, evidence, evidenceCard, bn)

				var valParts []string
				for cs := 0; cs < childCard; cs++ {
					valParts = append(valParts, bifFormatFloat(data[cs*numParentConfigs+pc]))
				}

				if _, err := fmt.Fprintf(w, "  (%s) %s;\n",
					strings.Join(parentStateNames, ", "),
					strings.Join(valParts, ", ")); err != nil {
					return fmt.Errorf("readwrite: write error: %w", err)
				}
			}
		}

		_ = states // used implicitly via bn.GetStates in decompose

		if _, err := fmt.Fprintf(w, "}\n"); err != nil {
			return fmt.Errorf("readwrite: write error: %w", err)
		}
	}

	return nil
}

// bifDecomposeParentConfig converts a flat parent configuration index to state names.
func bifDecomposeParentConfig(pc int, evidence []string, evidenceCard []int, bn *models.BayesianNetwork) []string {
	// Row-major decomposition: leftmost parent varies slowest.
	indices := make([]int, len(evidence))
	rem := pc
	for i := len(evidence) - 1; i >= 0; i-- {
		indices[i] = rem % evidenceCard[i]
		rem /= evidenceCard[i]
	}

	names := make([]string, len(evidence))
	for i, ev := range evidence {
		states := bn.GetStates(ev)
		if indices[i] < len(states) {
			names[i] = states[indices[i]]
		} else {
			names[i] = fmt.Sprintf("state%d", indices[i])
		}
	}
	return names
}

// bifFormatFloat formats a float64 for BIF output, using minimal precision.
func bifFormatFloat(v float64) string {
	s := fmt.Sprintf("%.10g", v)
	return s
}
