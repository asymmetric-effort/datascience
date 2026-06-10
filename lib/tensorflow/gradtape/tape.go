// Package gradtape implements reverse-mode automatic differentiation,
// analogous to tf.GradientTape.
package gradtape

import "fmt"

// maxTapeOperations limits the size of the recorded operations buffer.
const maxTapeOperations = 1 << 20 // ~1M operations

// OpType identifies the type of recorded operation.
type OpType uint8

const (
	// OpAdd represents addition.
	OpAdd OpType = iota
	// OpSub represents subtraction.
	OpSub
	// OpMul represents multiplication.
	OpMul
	// OpDiv represents division.
	OpDiv
	// OpNeg represents negation.
	OpNeg
	// OpReLU represents the ReLU activation.
	OpReLU
	// OpSigmoid represents the sigmoid activation.
	OpSigmoid
	// OpExp represents exponentiation.
	OpExp
	// OpLog represents natural logarithm.
	OpLog
	// OpSum represents summation reduction.
	OpSum
	// OpMatMul represents matrix multiplication.
	OpMatMul
)

// operation records a single operation on the tape.
type operation struct {
	op     OpType
	inputs []int // indices into the values slice
	output int   // index into the values slice
	cached []float64
}

// Variable represents a tracked value in the computation graph.
type Variable struct {
	tape  *Tape
	index int
}

// Value returns the current value of the variable.
func (v *Variable) Value() []float64 {
	out := make([]float64, len(v.tape.values[v.index]))
	copy(out, v.tape.values[v.index])
	return out
}

// Index returns the internal index of this variable on the tape.
func (v *Variable) Index() int {
	return v.index
}

// Tape records operations for reverse-mode automatic differentiation.
type Tape struct {
	values [][]float64
	ops    []operation
}

// New creates a new gradient tape.
func New() *Tape {
	return &Tape{}
}

// Variable creates a new tracked variable with the given initial values.
func (t *Tape) Variable(data []float64) *Variable {
	buf := make([]float64, len(data))
	copy(buf, data)
	idx := len(t.values)
	t.values = append(t.values, buf)
	return &Variable{tape: t, index: idx}
}

// Add records element-wise addition: out = a + b.
func (t *Tape) Add(a, b *Variable) (*Variable, error) {
	if len(a.tape.values[a.index]) != len(b.tape.values[b.index]) {
		return nil, fmt.Errorf("add: length mismatch %d vs %d", len(a.tape.values[a.index]), len(b.tape.values[b.index]))
	}
	if err := t.checkCapacity(); err != nil {
		return nil, err
	}
	dataA := t.values[a.index]
	dataB := t.values[b.index]
	result := make([]float64, len(dataA))
	for i := range dataA {
		result[i] = dataA[i] + dataB[i]
	}
	idx := len(t.values)
	t.values = append(t.values, result)
	t.ops = append(t.ops, operation{op: OpAdd, inputs: []int{a.index, b.index}, output: idx})
	return &Variable{tape: t, index: idx}, nil
}

// Sub records element-wise subtraction: out = a - b.
func (t *Tape) Sub(a, b *Variable) (*Variable, error) {
	if len(a.tape.values[a.index]) != len(b.tape.values[b.index]) {
		return nil, fmt.Errorf("sub: length mismatch %d vs %d", len(a.tape.values[a.index]), len(b.tape.values[b.index]))
	}
	if err := t.checkCapacity(); err != nil {
		return nil, err
	}
	dataA := t.values[a.index]
	dataB := t.values[b.index]
	result := make([]float64, len(dataA))
	for i := range dataA {
		result[i] = dataA[i] - dataB[i]
	}
	idx := len(t.values)
	t.values = append(t.values, result)
	t.ops = append(t.ops, operation{op: OpSub, inputs: []int{a.index, b.index}, output: idx})
	return &Variable{tape: t, index: idx}, nil
}

// Mul records element-wise multiplication: out = a * b.
func (t *Tape) Mul(a, b *Variable) (*Variable, error) {
	if len(a.tape.values[a.index]) != len(b.tape.values[b.index]) {
		return nil, fmt.Errorf("mul: length mismatch %d vs %d", len(a.tape.values[a.index]), len(b.tape.values[b.index]))
	}
	if err := t.checkCapacity(); err != nil {
		return nil, err
	}
	dataA := t.values[a.index]
	dataB := t.values[b.index]
	result := make([]float64, len(dataA))
	for i := range dataA {
		result[i] = dataA[i] * dataB[i]
	}
	idx := len(t.values)
	t.values = append(t.values, result)
	t.ops = append(t.ops, operation{
		op:     OpMul,
		inputs: []int{a.index, b.index},
		output: idx,
		cached: append(append([]float64{}, dataA...), dataB...),
	})
	return &Variable{tape: t, index: idx}, nil
}

// Neg records element-wise negation: out = -a.
func (t *Tape) Neg(a *Variable) (*Variable, error) {
	if err := t.checkCapacity(); err != nil {
		return nil, err
	}
	dataA := t.values[a.index]
	result := make([]float64, len(dataA))
	for i := range dataA {
		result[i] = -dataA[i]
	}
	idx := len(t.values)
	t.values = append(t.values, result)
	t.ops = append(t.ops, operation{op: OpNeg, inputs: []int{a.index}, output: idx})
	return &Variable{tape: t, index: idx}, nil
}

// Sum records a summation reduction: out = sum(a). Returns a single-element variable.
func (t *Tape) Sum(a *Variable) (*Variable, error) {
	if err := t.checkCapacity(); err != nil {
		return nil, err
	}
	dataA := t.values[a.index]
	total := 0.0
	for _, v := range dataA {
		total += v
	}
	idx := len(t.values)
	t.values = append(t.values, []float64{total})
	t.ops = append(t.ops, operation{
		op:     OpSum,
		inputs: []int{a.index},
		output: idx,
		cached: []float64{float64(len(dataA))},
	})
	return &Variable{tape: t, index: idx}, nil
}

// Gradient computes the gradient of the output variable with respect to all
// tracked variables using reverse-mode autodiff. Returns a map from variable
// index to gradient values.
func (t *Tape) Gradient(output *Variable) map[int][]float64 {
	grads := make(map[int][]float64)

	// Seed the output gradient with ones.
	outGrad := make([]float64, len(t.values[output.index]))
	for i := range outGrad {
		outGrad[i] = 1.0
	}
	grads[output.index] = outGrad

	// Walk operations in reverse order (reverse-mode).
	for i := len(t.ops) - 1; i >= 0; i-- {
		op := t.ops[i]
		outG := grads[op.output]
		if outG == nil {
			continue
		}

		switch op.op {
		case OpAdd:
			accumGrad(grads, op.inputs[0], outG)
			accumGrad(grads, op.inputs[1], outG)

		case OpSub:
			accumGrad(grads, op.inputs[0], outG)
			negG := make([]float64, len(outG))
			for j := range outG {
				negG[j] = -outG[j]
			}
			accumGrad(grads, op.inputs[1], negG)

		case OpMul:
			n := len(outG)
			cachedA := op.cached[:n]
			cachedB := op.cached[n:]
			gradA := make([]float64, n)
			gradB := make([]float64, n)
			for j := range outG {
				gradA[j] = outG[j] * cachedB[j]
				gradB[j] = outG[j] * cachedA[j]
			}
			accumGrad(grads, op.inputs[0], gradA)
			accumGrad(grads, op.inputs[1], gradB)

		case OpNeg:
			negG := make([]float64, len(outG))
			for j := range outG {
				negG[j] = -outG[j]
			}
			accumGrad(grads, op.inputs[0], negG)

		case OpSum:
			inputLen := int(op.cached[0])
			expanded := make([]float64, inputLen)
			for j := range expanded {
				expanded[j] = outG[0]
			}
			accumGrad(grads, op.inputs[0], expanded)
		}
	}

	return grads
}

func accumGrad(grads map[int][]float64, idx int, grad []float64) {
	existing, ok := grads[idx]
	if !ok {
		buf := make([]float64, len(grad))
		copy(buf, grad)
		grads[idx] = buf
		return
	}
	for i := range existing {
		existing[i] += grad[i]
	}
}

func (t *Tape) checkCapacity() error {
	if len(t.ops) >= maxTapeOperations {
		return fmt.Errorf("tape capacity exceeded: maximum %d operations", maxTapeOperations)
	}
	return nil
}
