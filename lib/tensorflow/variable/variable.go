// Package variable provides a trainable variable type,
// analogous to tf.Variable.
package variable

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

var nextID atomic.Int64

// Variable represents a mutable NDArray that persists across calls.
// Analogous to tf.Variable.
type Variable struct {
	mu        sync.RWMutex
	id        int64
	name      string
	value     *numgo.NDArray
	trainable bool
}

// New creates a new Variable with the given initial value and name.
func New(initialValue *numgo.NDArray, name string, trainable bool) *Variable {
	return &Variable{
		id:        nextID.Add(1),
		name:      name,
		value:     initialValue.Copy(),
		trainable: trainable,
	}
}

// ID returns the unique identifier of this variable.
func (v *Variable) ID() int64 {
	return v.id
}

// Name returns the name of this variable.
func (v *Variable) Name() string {
	return v.name
}

// Value returns a copy of the current NDArray value.
func (v *Variable) Value() *numgo.NDArray {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.value.Copy()
}

// shapesEqual returns true if two shape slices are identical.
func shapesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Assign replaces the variable's value. The new value must have the same shape.
func (v *Variable) Assign(newValue *numgo.NDArray) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !shapesEqual(v.value.Shape(), newValue.Shape()) {
		return fmt.Errorf("assign: shape mismatch %v vs %v", v.value.Shape(), newValue.Shape())
	}
	v.value = newValue.Copy()
	return nil
}

// AssignAdd adds an NDArray to the variable's current value in place.
func (v *Variable) AssignAdd(delta *numgo.NDArray) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !shapesEqual(v.value.Shape(), delta.Shape()) {
		return fmt.Errorf("assign_add: shape mismatch %v vs %v", v.value.Shape(), delta.Shape())
	}
	v.value = numgo.Add(v.value, delta)
	return nil
}

// AssignSub subtracts an NDArray from the variable's current value in place.
func (v *Variable) AssignSub(delta *numgo.NDArray) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !shapesEqual(v.value.Shape(), delta.Shape()) {
		return fmt.Errorf("assign_sub: shape mismatch %v vs %v", v.value.Shape(), delta.Shape())
	}
	v.value = numgo.Sub(v.value, delta)
	return nil
}

// Shape returns the shape of the variable's NDArray.
func (v *Variable) Shape() []int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.value.Shape()
}

// Trainable returns whether this variable participates in gradient computation.
func (v *Variable) Trainable() bool {
	return v.trainable
}

// NumElements returns the total number of elements in the variable.
func (v *Variable) NumElements() int {
	return v.value.Size()
}
