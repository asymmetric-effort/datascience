package inference

import (
	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

// ---------------------------------------------------------------------------
// Dependency-injection interfaces for inference package.
//
// These interfaces exist to make defensive error guards in computeMessage,
// Query, and elimination steps reachable in tests. The public API signatures
// are unchanged; internal functions accept interfaces for mock injection.
// ---------------------------------------------------------------------------

// factorMultiplier abstracts factor product and marginalize operations to
// enable testing of defensive error paths in message computation and
// variable elimination.
type factorMultiplier interface {
	Product(factors ...*factors.DiscreteFactor) (*factors.DiscreteFactor, error)
	Marginalize(f *factors.DiscreteFactor, vars []string) (*factors.DiscreteFactor, error)
}

// factorReducer abstracts evidence reduction to enable testing of
// defensive error paths in Query methods.
type factorReducer interface {
	Reduce(f *factors.DiscreteFactor, evidence map[string]int) (*factors.DiscreteFactor, error)
}

// defaultFactorMultiplier is the production implementation that delegates
// to the factors package.
type defaultFactorMultiplier struct{}

func (defaultFactorMultiplier) Product(fs ...*factors.DiscreteFactor) (*factors.DiscreteFactor, error) {
	return factors.FactorProduct(fs...)
}

func (defaultFactorMultiplier) Marginalize(f *factors.DiscreteFactor, vars []string) (*factors.DiscreteFactor, error) {
	return f.Marginalize(vars)
}

// defaultFactorReducer is the production implementation.
type defaultFactorReducer struct{}

func (defaultFactorReducer) Reduce(f *factors.DiscreteFactor, evidence map[string]int) (*factors.DiscreteFactor, error) {
	return f.Reduce(evidence)
}

// ---------------------------------------------------------------------------
// Testable internal functions for BeliefPropagation.
// ---------------------------------------------------------------------------

// computeMessageImpl is the testable implementation of computeMessage.
// It accepts a factorMultiplier interface to allow injection of failing
// mocks in tests for defensive error path coverage.
func computeMessageImpl(
	bp *BeliefPropagation,
	src, dst int,
	fm factorMultiplier,
) (*factors.DiscreteFactor, error) {
	current := bp.potentials[src]

	for _, nb := range bp.neighbors[src] {
		if nb == dst {
			continue
		}
		key := msgKey(nb, src)
		if msg, ok := bp.messages[key]; ok {
			prod, err := fm.Product(current, msg)
			if err != nil {
				return nil, err
			}
			current = prod
		}
	}

	sepKey := edgeKey(src, dst)
	sepVars := bp.separators[sepKey]

	sepSet := make(map[string]bool, len(sepVars))
	for _, v := range sepVars {
		sepSet[v] = true
	}

	currentVars := current.Variables()
	var margVars []string
	for _, v := range currentVars {
		if !sepSet[v] {
			margVars = append(margVars, v)
		}
	}

	if len(margVars) == 0 {
		return current.Copy(), nil
	}

	msg, err := fm.Marginalize(current, margVars)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// eliminateVariableImpl is the testable implementation of eliminateVariable.
// It accepts a factorMultiplier interface to allow injection of failing
// mocks in tests for defensive error path coverage.
func eliminateVariableImpl(
	factorList []*factors.DiscreteFactor,
	variable string,
	fm factorMultiplier,
) ([]*factors.DiscreteFactor, error) {
	var containing []*factors.DiscreteFactor
	var remaining []*factors.DiscreteFactor

	for _, f := range factorList {
		if varSet(f)[variable] {
			containing = append(containing, f)
		} else {
			remaining = append(remaining, f)
		}
	}

	if len(containing) == 0 {
		return factorList, nil
	}

	product, err := fm.Product(containing...)
	if err != nil {
		return nil, err
	}

	prodVars := product.Variables()
	if len(prodVars) == 1 && prodVars[0] == variable {
		return remaining, nil
	}

	marginalized, err := fm.Marginalize(product, []string{variable})
	if err != nil {
		return nil, err
	}

	return append(remaining, marginalized), nil
}

// identificationChecker abstracts identification criterion checks
// to enable testing all dispatch branches in IdentificationMethod.
type identificationChecker interface {
	canBackdoor(treatment, outcome string) bool
	canFrontdoor(treatment, outcome string) bool
	canIV(treatment, outcome string) bool
}

// defaultIdentificationChecker delegates to CausalInference methods.
type defaultIdentificationChecker struct {
	ci *CausalInference
}

func (d defaultIdentificationChecker) canBackdoor(treatment, outcome string) bool {
	return d.ci.canIdentifyByBackdoor(treatment, outcome)
}

func (d defaultIdentificationChecker) canFrontdoor(treatment, outcome string) bool {
	return d.ci.canIdentifyByFrontdoor(treatment, outcome)
}

func (d defaultIdentificationChecker) canIV(treatment, outcome string) bool {
	return d.ci.canIdentifyByIV(treatment, outcome)
}

// identificationMethodImpl is the testable implementation of IdentificationMethod.
func identificationMethodImpl(treatment, outcome string, ic identificationChecker) string {
	if ic.canBackdoor(treatment, outcome) {
		return "backdoor"
	}
	if ic.canFrontdoor(treatment, outcome) {
		return "frontdoor"
	}
	if ic.canIV(treatment, outcome) {
		return "iv"
	}
	return "none"
}

// maxEliminateVariableImpl is the testable implementation of maxEliminateVariable.
// It accepts a factorMultiplier interface to allow injection of failing
// mocks in tests for defensive error path coverage.
func maxEliminateVariableImpl(
	factorList []*factors.DiscreteFactor,
	variable string,
	fm factorMultiplier,
) ([]*factors.DiscreteFactor, error) {
	var containing []*factors.DiscreteFactor
	var remaining []*factors.DiscreteFactor

	for _, f := range factorList {
		if varSet(f)[variable] {
			containing = append(containing, f)
		} else {
			remaining = append(remaining, f)
		}
	}

	if len(containing) == 0 {
		return factorList, nil
	}

	product, err := fm.Product(containing...)
	if err != nil {
		return nil, err
	}

	prodVars := product.Variables()
	if len(prodVars) == 1 && prodVars[0] == variable {
		return remaining, nil
	}

	maximized, err := maxMarginalize(product, variable)
	if err != nil {
		return nil, err
	}

	return append(remaining, maximized), nil
}

// reduceAllImpl is the testable implementation of reduceAll.
// It accepts a factorReducer interface to allow injection of failing
// mocks in tests for defensive error path coverage.
func reduceAllImpl(
	factorList []*factors.DiscreteFactor,
	evidence map[string]int,
	fr factorReducer,
) ([]*factors.DiscreteFactor, error) {
	result := make([]*factors.DiscreteFactor, 0, len(factorList))
	for _, f := range factorList {
		fVars := varSet(f)
		applicable := make(map[string]int)
		for v, val := range evidence {
			if fVars[v] {
				applicable[v] = val
			}
		}
		reduced, err := fr.Reduce(f, applicable)
		if err != nil {
			return nil, err
		}
		result = append(result, reduced)
	}
	return result, nil
}
