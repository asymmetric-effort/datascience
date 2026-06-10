// Package schedule provides learning rate schedule functions,
// analogous to tf.keras.optimizers.schedules.
package schedule

import "math"

// Schedule is the interface for learning rate schedules.
type Schedule interface {
	// LearningRate returns the learning rate for the given step.
	LearningRate(step int) float64
	// Name returns the schedule name.
	Name() string
}

// ExponentialDecay implements exponential learning rate decay.
// lr = initialLR * decayRate ^ (step / decaySteps).
// Analogous to tf.keras.optimizers.schedules.ExponentialDecay.
type ExponentialDecay struct {
	InitialLR  float64
	DecaySteps int
	DecayRate  float64
	Staircase  bool
}

// NewExponentialDecay creates a new ExponentialDecay schedule.
func NewExponentialDecay(initialLR float64, decaySteps int, decayRate float64, staircase bool) *ExponentialDecay {
	return &ExponentialDecay{
		InitialLR:  initialLR,
		DecaySteps: decaySteps,
		DecayRate:  decayRate,
		Staircase:  staircase,
	}
}

// LearningRate returns the learning rate for the given step.
func (e *ExponentialDecay) LearningRate(step int) float64 {
	exponent := float64(step) / float64(e.DecaySteps)
	if e.Staircase {
		exponent = math.Floor(exponent)
	}
	return e.InitialLR * math.Pow(e.DecayRate, exponent)
}

// Name returns "exponential_decay".
func (e *ExponentialDecay) Name() string {
	return "exponential_decay"
}

// CosineDecay implements cosine annealing learning rate decay.
// lr = initialLR * 0.5 * (1 + cos(pi * step / totalSteps)).
// Analogous to tf.keras.optimizers.schedules.CosineDecay.
type CosineDecay struct {
	InitialLR  float64
	TotalSteps int
	Alpha      float64 // minimum learning rate as fraction of initialLR
}

// NewCosineDecay creates a new CosineDecay schedule.
func NewCosineDecay(initialLR float64, totalSteps int, alpha float64) *CosineDecay {
	return &CosineDecay{
		InitialLR:  initialLR,
		TotalSteps: totalSteps,
		Alpha:      alpha,
	}
}

// LearningRate returns the learning rate for the given step.
func (c *CosineDecay) LearningRate(step int) float64 {
	if step >= c.TotalSteps {
		return c.InitialLR * c.Alpha
	}
	progress := float64(step) / float64(c.TotalSteps)
	cosine := 0.5 * (1 + math.Cos(math.Pi*progress))
	return c.InitialLR * (c.Alpha + (1-c.Alpha)*cosine)
}

// Name returns "cosine_decay".
func (c *CosineDecay) Name() string {
	return "cosine_decay"
}

// PiecewiseConstant implements a piecewise constant learning rate.
// Analogous to tf.keras.optimizers.schedules.PiecewiseConstantDecay.
type PiecewiseConstant struct {
	Boundaries []int
	Values     []float64
}

// NewPiecewiseConstant creates a new PiecewiseConstant schedule.
// len(values) must equal len(boundaries) + 1.
func NewPiecewiseConstant(boundaries []int, values []float64) *PiecewiseConstant {
	return &PiecewiseConstant{
		Boundaries: boundaries,
		Values:     values,
	}
}

// LearningRate returns the learning rate for the given step.
func (p *PiecewiseConstant) LearningRate(step int) float64 {
	for i, b := range p.Boundaries {
		if step < b {
			return p.Values[i]
		}
	}
	return p.Values[len(p.Values)-1]
}

// Name returns "piecewise_constant".
func (p *PiecewiseConstant) Name() string {
	return "piecewise_constant"
}

// WarmupSchedule implements linear warmup followed by another schedule.
type WarmupSchedule struct {
	WarmupSteps int
	PeakLR      float64
	PostWarmup  Schedule
}

// NewWarmupSchedule creates a warmup schedule.
func NewWarmupSchedule(warmupSteps int, peakLR float64, postWarmup Schedule) *WarmupSchedule {
	return &WarmupSchedule{
		WarmupSteps: warmupSteps,
		PeakLR:      peakLR,
		PostWarmup:  postWarmup,
	}
}

// LearningRate returns the learning rate for the given step.
func (w *WarmupSchedule) LearningRate(step int) float64 {
	if step < w.WarmupSteps {
		return w.PeakLR * float64(step+1) / float64(w.WarmupSteps)
	}
	return w.PostWarmup.LearningRate(step - w.WarmupSteps)
}

// Name returns "warmup".
func (w *WarmupSchedule) Name() string {
	return "warmup"
}
