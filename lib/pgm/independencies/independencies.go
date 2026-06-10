package independencies

import (
	"fmt"
	"sort"
	"strings"
)

// Independencies represents a collection of IndependenceAssertion values.
type Independencies struct {
	assertions []*IndependenceAssertion
}

// NewIndependencies creates a new empty Independencies collection.
func NewIndependencies() *Independencies {
	return &Independencies{
		assertions: make([]*IndependenceAssertion, 0),
	}
}

// Add appends one or more assertions to the collection, skipping duplicates.
func (ind *Independencies) Add(assertions ...*IndependenceAssertion) {
	for _, a := range assertions {
		if a == nil {
			continue
		}
		if !ind.Contains(a) {
			ind.assertions = append(ind.assertions, a)
		}
	}
}

// Remove removes an assertion from the collection (by equality).
func (ind *Independencies) Remove(assertion *IndependenceAssertion) {
	if assertion == nil {
		return
	}
	for i, a := range ind.assertions {
		if a.Equals(assertion) {
			ind.assertions = append(ind.assertions[:i], ind.assertions[i+1:]...)
			return
		}
	}
}

// Contains returns true if the collection contains an assertion equal to the given one.
func (ind *Independencies) Contains(assertion *IndependenceAssertion) bool {
	if assertion == nil {
		return false
	}
	for _, a := range ind.assertions {
		if a.Equals(assertion) {
			return true
		}
	}
	return false
}

// GetAssertions returns a copy of the assertion slice.
func (ind *Independencies) GetAssertions() []*IndependenceAssertion {
	result := make([]*IndependenceAssertion, len(ind.assertions))
	copy(result, ind.assertions)
	return result
}

// Len returns the number of assertions in the collection.
func (ind *Independencies) Len() int {
	return len(ind.assertions)
}

// IsEquivalent returns true if this collection contains the same set of assertions
// as other (order independent).
func (ind *Independencies) IsEquivalent(other *Independencies) bool {
	if other == nil {
		return false
	}
	if ind.Len() != other.Len() {
		return false
	}
	for _, a := range ind.assertions {
		if !other.Contains(a) {
			return false
		}
	}
	return true
}

// GetAllVariables returns the union of all variables appearing in any assertion
// in the collection, sorted.
func (ind *Independencies) GetAllVariables() []string {
	varSet := make(map[string]bool)
	for _, a := range ind.assertions {
		for _, v := range a.event1 {
			varSet[v] = true
		}
		for _, v := range a.event2 {
			varSet[v] = true
		}
		for _, v := range a.given {
			varSet[v] = true
		}
	}
	result := make([]string, 0, len(varSet))
	for v := range varSet {
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}

// Closure computes the closure of the independence assertions under the
// graphoid axioms: symmetry, decomposition, weak union, and contraction.
// Returns a new Independencies containing the original assertions plus
// all derivable ones.
func (ind *Independencies) Closure() *Independencies {
	result := NewIndependencies()
	for _, a := range ind.assertions {
		result.Add(a)
	}

	changed := true
	for changed {
		changed = false
		current := result.GetAssertions()

		for _, a := range current {
			// Symmetry: X _|_ Y | Z => Y _|_ X | Z
			sym := NewIndependenceAssertion(a.Event2(), a.Event1(), a.Given())
			if !result.Contains(sym) {
				result.Add(sym)
				changed = true
			}

			// Decomposition: X _|_ {Y, W} | Z => X _|_ Y | Z and X _|_ W | Z
			e2 := a.Event2()
			if len(e2) > 1 {
				for _, v := range e2 {
					decomp := NewIndependenceAssertion(a.Event1(), []string{v}, a.Given())
					if !result.Contains(decomp) {
						result.Add(decomp)
						changed = true
					}
				}
			}

			// Weak union: X _|_ {Y, W} | Z => X _|_ Y | Z ∪ {W}
			if len(e2) > 1 {
				for i, y := range e2 {
					rest := make([]string, 0, len(e2)-1)
					for j, w := range e2 {
						if j != i {
							rest = append(rest, w)
						}
					}
					newGiven := append(a.Given(), rest...)
					wu := NewIndependenceAssertion(a.Event1(), []string{y}, newGiven)
					if !result.Contains(wu) {
						result.Add(wu)
						changed = true
					}
				}
			}
		}
	}

	return result
}

// Entails checks if the assertions in this collection imply the given assertion.
// Uses a simple containment check.
func (ind *Independencies) Entails(assertion *IndependenceAssertion) bool {
	if assertion == nil {
		return true
	}
	for _, a := range ind.assertions {
		if a.Contains(assertion) {
			return true
		}
	}
	return false
}

// Reduce removes redundant assertions from the collection. An assertion is
// redundant if it is entailed by another assertion via containment.
// Returns a new reduced Independencies.
func (ind *Independencies) Reduce() *Independencies {
	result := NewIndependencies()
	assertions := ind.GetAssertions()

	for i, a := range assertions {
		entailed := false
		for j, b := range assertions {
			if i != j && b.Contains(a) {
				entailed = true
				break
			}
		}
		if !entailed {
			result.Add(a)
		}
	}

	return result
}

// LatexString returns a LaTeX representation of all assertions.
func (ind *Independencies) LatexString() string {
	if len(ind.assertions) == 0 {
		return "\\emptyset"
	}
	parts := make([]string, len(ind.assertions))
	for i, a := range ind.assertions {
		parts[i] = a.LatexString()
	}
	return strings.Join(parts, ", \\quad ")
}

// GetFactorizedProduct returns a string representation of the factorized
// product implied by the independence assertions.
func (ind *Independencies) GetFactorizedProduct() string {
	if len(ind.assertions) == 0 {
		return "P(V)"
	}
	var parts []string
	for _, a := range ind.assertions {
		x := strings.Join(a.Event1(), ", ")
		given := a.Given()
		if len(given) > 0 {
			parts = append(parts, fmt.Sprintf("P(%s | %s)", x, strings.Join(given, ", ")))
		} else {
			parts = append(parts, fmt.Sprintf("P(%s)", x))
		}
	}
	return strings.Join(parts, " * ")
}

// String returns a human-readable representation of all assertions in the collection.
func (ind *Independencies) String() string {
	if len(ind.assertions) == 0 {
		return "{}"
	}
	parts := make([]string, len(ind.assertions))
	for i, a := range ind.assertions {
		parts[i] = a.String()
	}
	return "{\n  " + strings.Join(parts, ",\n  ") + "\n}"
}
