package factors

import (
	"fmt"
	"math"
)

// TabularCPD represents a conditional probability distribution in tabular form.
// It wraps a DiscreteFactor where the first variable is the child and the
// remaining variables are evidence (parents).
//
// The values parameter to NewTabularCPD is organized as a 2D table:
// rows correspond to child variable states, columns correspond to parent
// configurations (in row-major order of evidence variables).
type TabularCPD struct {
	variable     string
	variableCard int
	evidence     []string
	evidenceCard []int
	factor       *DiscreteFactor
}

// NewTabularCPD creates a new TabularCPD.
//
// values is a 2D slice of shape [variableCard][numParentConfigs] where
// numParentConfigs = product of evidenceCard. Each column should sum to 1.
//
// The resulting factor has variables ordered as [variable] + evidence,
// with cardinality [variableCard] + evidenceCard.
func NewTabularCPD(variable string, variableCard int, values [][]float64,
	evidence []string, evidenceCard []int) (*TabularCPD, error) {

	if variableCard <= 0 {
		return nil, fmt.Errorf("factors: variableCard must be positive, got %d", variableCard)
	}
	if len(evidence) != len(evidenceCard) {
		return nil, fmt.Errorf("factors: evidence length %d != evidenceCard length %d",
			len(evidence), len(evidenceCard))
	}
	for _, ec := range evidenceCard {
		if ec <= 0 {
			return nil, fmt.Errorf("factors: evidence cardinality must be positive, got %d", ec)
		}
	}

	numParentConfigs := 1
	for _, ec := range evidenceCard {
		numParentConfigs *= ec
	}

	if len(values) != variableCard {
		return nil, fmt.Errorf("factors: values must have %d rows (variableCard), got %d",
			variableCard, len(values))
	}
	for i, row := range values {
		if len(row) != numParentConfigs {
			return nil, fmt.Errorf("factors: values row %d has length %d, expected %d",
				i, len(row), numParentConfigs)
		}
	}

	// Build flat values in row-major order for [variable, evidence...].
	// The variable is the first axis, evidence follow.
	// Flat index = variable_state * numParentConfigs + parent_config.
	allVars := make([]string, 0, 1+len(evidence))
	allVars = append(allVars, variable)
	allVars = append(allVars, evidence...)

	allCard := make([]int, 0, 1+len(evidenceCard))
	allCard = append(allCard, variableCard)
	allCard = append(allCard, evidenceCard...)

	totalSize := variableCard * numParentConfigs
	flat := make([]float64, totalSize)
	for childState := 0; childState < variableCard; childState++ {
		for parentConfig := 0; parentConfig < numParentConfigs; parentConfig++ {
			flat[childState*numParentConfigs+parentConfig] = values[childState][parentConfig]
		}
	}

	factor, err := NewDiscreteFactor(allVars, allCard, flat)
	if err != nil {
		return nil, fmt.Errorf("factors: failed to create underlying factor: %w", err)
	}

	ev := make([]string, len(evidence))
	copy(ev, evidence)
	ec := make([]int, len(evidenceCard))
	copy(ec, evidenceCard)

	return &TabularCPD{
		variable:     variable,
		variableCard: variableCard,
		evidence:     ev,
		evidenceCard: ec,
		factor:       factor,
	}, nil
}

// Variable returns the child variable name.
func (cpd *TabularCPD) Variable() string {
	return cpd.variable
}

// VariableCard returns the cardinality of the child variable.
func (cpd *TabularCPD) VariableCard() int {
	return cpd.variableCard
}

// Evidence returns a copy of the evidence variable names.
func (cpd *TabularCPD) Evidence() []string {
	ev := make([]string, len(cpd.evidence))
	copy(ev, cpd.evidence)
	return ev
}

// EvidenceCard returns a copy of the evidence cardinalities.
func (cpd *TabularCPD) EvidenceCard() []int {
	ec := make([]int, len(cpd.evidenceCard))
	copy(ec, cpd.evidenceCard)
	return ec
}

// Validate checks that each column of the CPD sums to 1 (within tolerance).
func (cpd *TabularCPD) Validate() error {
	const tol = 1e-6
	numParentConfigs := 1
	for _, ec := range cpd.evidenceCard {
		numParentConfigs *= ec
	}

	data := cpd.factor.values.Data()

	for parentConfig := 0; parentConfig < numParentConfigs; parentConfig++ {
		sum := 0.0
		for childState := 0; childState < cpd.variableCard; childState++ {
			sum += data[childState*numParentConfigs+parentConfig]
		}
		if math.Abs(sum-1.0) > tol {
			return fmt.Errorf("factors: CPD column %d sums to %f, expected 1.0", parentConfig, sum)
		}
	}
	return nil
}

// ToFactor returns the underlying DiscreteFactor (a copy).
func (cpd *TabularCPD) ToFactor() *DiscreteFactor {
	return cpd.factor.Copy()
}

// Copy returns a deep copy of the TabularCPD.
func (cpd *TabularCPD) Copy() *TabularCPD {
	ev := make([]string, len(cpd.evidence))
	copy(ev, cpd.evidence)
	ec := make([]int, len(cpd.evidenceCard))
	copy(ec, cpd.evidenceCard)
	return &TabularCPD{
		variable:     cpd.variable,
		variableCard: cpd.variableCard,
		evidence:     ev,
		evidenceCard: ec,
		factor:       cpd.factor.Copy(),
	}
}
