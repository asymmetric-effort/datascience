// Package optimizer provides optimization algorithms for training neural networks,
// analogous to tf.keras.optimizers.
package optimizer

import "math"

// maxMomentBufferSize limits the number of parameters an optimizer can track.
const maxMomentBufferSize = 1 << 24 // ~16M parameters

// SGD implements stochastic gradient descent with optional momentum.
// Analogous to tf.keras.optimizers.SGD.
type SGD struct {
	lr       float64
	momentum float64
	velocity map[int][]float64
}

// NewSGD creates a new SGD optimizer with the given learning rate and momentum.
func NewSGD(lr, momentum float64) *SGD {
	return &SGD{
		lr:       lr,
		momentum: momentum,
		velocity: make(map[int][]float64),
	}
}

// Step updates parameters using SGD with momentum.
// paramID uniquely identifies the parameter for velocity tracking.
// params and grads are modified in place: params -= lr * update.
func (s *SGD) Step(paramID int, params, grads []float64) {
	v, ok := s.velocity[paramID]
	if !ok || len(v) != len(params) {
		v = make([]float64, len(params))
		s.velocity[paramID] = v
	}
	for i := range params {
		v[i] = s.momentum*v[i] + grads[i]
		params[i] -= s.lr * v[i]
	}
}

// LearningRate returns the current learning rate.
func (s *SGD) LearningRate() float64 {
	return s.lr
}

// SetLearningRate updates the learning rate.
func (s *SGD) SetLearningRate(lr float64) {
	s.lr = lr
}

// Adam implements the Adam optimizer.
// Analogous to tf.keras.optimizers.Adam.
type Adam struct {
	lr      float64
	beta1   float64
	beta2   float64
	epsilon float64
	step    int
	m       map[int][]float64 // first moment estimates
	v       map[int][]float64 // second moment estimates
}

// NewAdam creates a new Adam optimizer with default hyperparameters.
func NewAdam(lr float64) *Adam {
	return &Adam{
		lr:      lr,
		beta1:   0.9,
		beta2:   0.999,
		epsilon: 1e-8,
		m:       make(map[int][]float64),
		v:       make(map[int][]float64),
	}
}

// NewAdamWithParams creates a new Adam optimizer with custom hyperparameters.
func NewAdamWithParams(lr, beta1, beta2, epsilon float64) *Adam {
	return &Adam{
		lr:      lr,
		beta1:   beta1,
		beta2:   beta2,
		epsilon: epsilon,
		m:       make(map[int][]float64),
		v:       make(map[int][]float64),
	}
}

// Step updates parameters using the Adam algorithm.
// paramID uniquely identifies the parameter for moment tracking.
// params and grads are modified in place.
func (a *Adam) Step(paramID int, params, grads []float64) {
	a.step++

	m, ok := a.m[paramID]
	if !ok || len(m) != len(params) {
		m = make([]float64, len(params))
		a.m[paramID] = m
	}
	v, ok := a.v[paramID]
	if !ok || len(v) != len(params) {
		v = make([]float64, len(params))
		a.v[paramID] = v
	}

	beta1Power := math.Pow(a.beta1, float64(a.step))
	beta2Power := math.Pow(a.beta2, float64(a.step))

	for i := range params {
		m[i] = a.beta1*m[i] + (1-a.beta1)*grads[i]
		v[i] = a.beta2*v[i] + (1-a.beta2)*grads[i]*grads[i]

		mHat := m[i] / (1 - beta1Power)
		vHat := v[i] / (1 - beta2Power)

		params[i] -= a.lr * mHat / (math.Sqrt(vHat) + a.epsilon)
	}
}

// LearningRate returns the current learning rate.
func (a *Adam) LearningRate() float64 {
	return a.lr
}

// SetLearningRate updates the learning rate.
func (a *Adam) SetLearningRate(lr float64) {
	a.lr = lr
}

// RMSProp implements the RMSProp optimizer.
// Analogous to tf.keras.optimizers.RMSprop.
type RMSProp struct {
	lr      float64
	rho     float64
	epsilon float64
	cache   map[int][]float64
}

// NewRMSProp creates a new RMSProp optimizer.
func NewRMSProp(lr float64) *RMSProp {
	return &RMSProp{
		lr:      lr,
		rho:     0.9,
		epsilon: 1e-8,
		cache:   make(map[int][]float64),
	}
}

// Step updates parameters using the RMSProp algorithm.
func (r *RMSProp) Step(paramID int, params, grads []float64) {
	c, ok := r.cache[paramID]
	if !ok || len(c) != len(params) {
		c = make([]float64, len(params))
		r.cache[paramID] = c
	}
	for i := range params {
		c[i] = r.rho*c[i] + (1-r.rho)*grads[i]*grads[i]
		params[i] -= r.lr * grads[i] / (math.Sqrt(c[i]) + r.epsilon)
	}
}

// LearningRate returns the current learning rate.
func (r *RMSProp) LearningRate() float64 {
	return r.lr
}

// SetLearningRate updates the learning rate.
func (r *RMSProp) SetLearningRate(lr float64) {
	r.lr = lr
}
