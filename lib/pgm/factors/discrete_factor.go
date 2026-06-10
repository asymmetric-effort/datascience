package factors

import (
	"fmt"
	"math"
	"strings"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// DiscreteFactor represents a discrete factor (potential) over a set of
// variables, each with a finite cardinality. Values are stored in a flat
// NDArray in row-major order corresponding to the variable order.
type DiscreteFactor struct {
	variables   []string
	cardinality []int
	values      *numgo.NDArray
}

// NewDiscreteFactor creates a new DiscreteFactor.
// variables and cardinality must have the same length.
// values length must equal the product of all cardinalities.
func NewDiscreteFactor(variables []string, cardinality []int, values []float64) (*DiscreteFactor, error) {
	if len(variables) != len(cardinality) {
		return nil, fmt.Errorf("factors: variables length %d != cardinality length %d", len(variables), len(cardinality))
	}
	expectedSize := 1
	for _, c := range cardinality {
		if c <= 0 {
			return nil, fmt.Errorf("factors: cardinality must be positive, got %d", c)
		}
		expectedSize *= c
	}
	if len(values) != expectedSize {
		return nil, fmt.Errorf("factors: values length %d != expected %d", len(values), expectedSize)
	}
	// Check for duplicate variable names.
	seen := make(map[string]bool, len(variables))
	for _, v := range variables {
		if seen[v] {
			return nil, fmt.Errorf("factors: duplicate variable %q", v)
		}
		seen[v] = true
	}

	vars := make([]string, len(variables))
	copy(vars, variables)
	card := make([]int, len(cardinality))
	copy(card, cardinality)

	nd := numgo.NewNDArray(card, values)
	return &DiscreteFactor{
		variables:   vars,
		cardinality: card,
		values:      nd,
	}, nil
}

// Variables returns a copy of the variable names.
func (f *DiscreteFactor) Variables() []string {
	v := make([]string, len(f.variables))
	copy(v, f.variables)
	return v
}

// Cardinality returns a copy of the cardinality slice.
func (f *DiscreteFactor) Cardinality() []int {
	c := make([]int, len(f.cardinality))
	copy(c, f.cardinality)
	return c
}

// Values returns the underlying NDArray.
func (f *DiscreteFactor) Values() *numgo.NDArray {
	return f.values
}

// varIndex returns the axis index for a variable, or -1 if not found.
func (f *DiscreteFactor) varIndex(name string) int {
	for i, v := range f.variables {
		if v == name {
			return i
		}
	}
	return -1
}

// assignmentToIndices converts a variable-name assignment map to an ordered
// slice of indices matching the factor's variable order.
func (f *DiscreteFactor) assignmentToIndices(assignment map[string]int) []int {
	indices := make([]int, len(f.variables))
	for i, v := range f.variables {
		indices[i] = assignment[v]
	}
	return indices
}

// GetValue returns the factor value for a complete assignment.
func (f *DiscreteFactor) GetValue(assignment map[string]int) float64 {
	return f.values.Get(f.assignmentToIndices(assignment)...)
}

// SetValue sets the factor value for a complete assignment.
func (f *DiscreteFactor) SetValue(assignment map[string]int, value float64) {
	f.values.Set(value, f.assignmentToIndices(assignment)...)
}

// flatToAssignment decomposes a flat index into a map of variable assignments
// based on the factor's variable ordering and cardinalities.
func (f *DiscreteFactor) flatToAssignment(flat int) map[string]int {
	assignment := make(map[string]int, len(f.variables))
	rem := flat
	// Row-major: leftmost dimension varies slowest.
	for i := len(f.variables) - 1; i >= 0; i-- {
		assignment[f.variables[i]] = rem % f.cardinality[i]
		rem /= f.cardinality[i]
	}
	return assignment
}

// totalSize returns the product of all cardinalities.
func (f *DiscreteFactor) totalSize() int {
	size := 1
	for _, c := range f.cardinality {
		size *= c
	}
	return size
}

// Marginalize sums out the given variables and returns a new factor over
// the remaining variables.
func (f *DiscreteFactor) Marginalize(variables []string) (*DiscreteFactor, error) {
	if len(variables) == 0 {
		return f.Copy(), nil
	}
	// Validate that all variables to marginalize exist.
	margSet := make(map[string]bool, len(variables))
	for _, v := range variables {
		if f.varIndex(v) == -1 {
			return nil, fmt.Errorf("factors: variable %q not in factor", v)
		}
		margSet[v] = true
	}
	if len(margSet) == len(f.variables) {
		return nil, fmt.Errorf("factors: cannot marginalize all variables")
	}

	// Build remaining variables and cardinalities.
	var newVars []string
	var newCard []int
	for i, v := range f.variables {
		if !margSet[v] {
			newVars = append(newVars, v)
			newCard = append(newCard, f.cardinality[i])
		}
	}

	newSize := 1
	for _, c := range newCard {
		newSize *= c
	}
	newValues := make([]float64, newSize)

	data := f.values.Data()
	totalSize := f.totalSize()

	// Iterate over all assignments in the original factor.
	for flat := 0; flat < totalSize; flat++ {
		assignment := f.flatToAssignment(flat)
		// Compute flat index in the new factor.
		newFlat := 0
		stride := 1
		for i := len(newVars) - 1; i >= 0; i-- {
			newFlat += assignment[newVars[i]] * stride
			stride *= newCard[i]
		}
		newValues[newFlat] += data[flat]
	}

	return NewDiscreteFactor(newVars, newCard, newValues)
}

// Reduce fixes the given variables to specified values and returns a new
// factor over the remaining (non-evidence) variables.
func (f *DiscreteFactor) Reduce(evidence map[string]int) (*DiscreteFactor, error) {
	if len(evidence) == 0 {
		return f.Copy(), nil
	}
	// Validate evidence variables.
	for v, val := range evidence {
		idx := f.varIndex(v)
		if idx == -1 {
			return nil, fmt.Errorf("factors: evidence variable %q not in factor", v)
		}
		if val < 0 || val >= f.cardinality[idx] {
			return nil, fmt.Errorf("factors: evidence value %d out of range for variable %q (card %d)", val, v, f.cardinality[idx])
		}
	}

	// Build remaining variables.
	var newVars []string
	var newCard []int
	for i, v := range f.variables {
		if _, ok := evidence[v]; !ok {
			newVars = append(newVars, v)
			newCard = append(newCard, f.cardinality[i])
		}
	}

	if len(newVars) == 0 {
		// All variables are evidence; return a scalar factor.
		// Use a single dummy dimension.
		val := f.GetValue(evidence)
		return NewDiscreteFactor(nil, nil, []float64{val})
	}

	newSize := 1
	for _, c := range newCard {
		newSize *= c
	}
	newValues := make([]float64, newSize)

	data := f.values.Data()
	totalSize := f.totalSize()

	for flat := 0; flat < totalSize; flat++ {
		assignment := f.flatToAssignment(flat)
		// Check if this assignment matches the evidence.
		match := true
		for v, val := range evidence {
			if assignment[v] != val {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		// Compute new flat index.
		newFlat := 0
		stride := 1
		for i := len(newVars) - 1; i >= 0; i-- {
			newFlat += assignment[newVars[i]] * stride
			stride *= newCard[i]
		}
		newValues[newFlat] = data[flat]
	}

	return NewDiscreteFactor(newVars, newCard, newValues)
}

// Normalize normalizes the factor values in-place so they sum to 1.
func (f *DiscreteFactor) Normalize() {
	data := f.values.Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	if sum == 0 || math.IsNaN(sum) || math.IsInf(sum, 0) {
		return
	}
	totalSize := f.totalSize()
	indices := make([]int, len(f.variables))
	for flat := 0; flat < totalSize; flat++ {
		// Decompose flat index.
		rem := flat
		for i := len(f.variables) - 1; i >= 0; i-- {
			indices[i] = rem % f.cardinality[i]
			rem /= f.cardinality[i]
		}
		f.values.Set(data[flat]/sum, indices...)
	}
}

// Copy returns a deep copy of the factor.
func (f *DiscreteFactor) Copy() *DiscreteFactor {
	vars := make([]string, len(f.variables))
	copy(vars, f.variables)
	card := make([]int, len(f.cardinality))
	copy(card, f.cardinality)
	return &DiscreteFactor{
		variables:   vars,
		cardinality: card,
		values:      f.values.Copy(),
	}
}

// String returns a human-readable representation of the factor.
func (f *DiscreteFactor) String() string {
	var b strings.Builder
	b.WriteString("DiscreteFactor(")
	b.WriteString("variables=")
	b.WriteString(fmt.Sprintf("%v", f.variables))
	b.WriteString(", cardinality=")
	b.WriteString(fmt.Sprintf("%v", f.cardinality))
	b.WriteString(", values=")
	b.WriteString(f.values.String())
	b.WriteString(")")
	return b.String()
}
