// Package callbacks provides training callbacks,
// analogous to tf.keras.callbacks.
package callbacks

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
)

// Callback is the interface for training callbacks.
type Callback interface {
	// OnEpochEnd is called at the end of each epoch.
	// Returns true if training should stop.
	OnEpochEnd(epoch int, logs map[string]float64) bool
	// Name returns the callback name.
	Name() string
}

// EarlyStopping stops training when a monitored metric stops improving.
// Analogous to tf.keras.callbacks.EarlyStopping.
type EarlyStopping struct {
	monitor  string
	patience int
	minDelta float64
	mode     string // "min" or "max"
	best     float64
	wait     int
	stopped  bool
}

// NewEarlyStopping creates a new EarlyStopping callback.
// monitor: metric name to monitor. patience: number of epochs with no improvement.
// minDelta: minimum change to qualify as improvement.
// mode: "min" (lower is better) or "max" (higher is better).
func NewEarlyStopping(monitor string, patience int, minDelta float64, mode string) *EarlyStopping {
	best := math.Inf(1)
	if mode == "max" {
		best = math.Inf(-1)
	}
	return &EarlyStopping{
		monitor:  monitor,
		patience: patience,
		minDelta: minDelta,
		mode:     mode,
		best:     best,
	}
}

// OnEpochEnd checks if the monitored metric improved.
func (e *EarlyStopping) OnEpochEnd(epoch int, logs map[string]float64) bool {
	val, ok := logs[e.monitor]
	if !ok {
		return false
	}

	improved := false
	if e.mode == "min" {
		improved = val < e.best-e.minDelta
	} else {
		improved = val > e.best+e.minDelta
	}

	if improved {
		e.best = val
		e.wait = 0
	} else {
		e.wait++
		if e.wait >= e.patience {
			e.stopped = true
			return true
		}
	}
	return false
}

// Name returns "early_stopping".
func (e *EarlyStopping) Name() string {
	return "early_stopping"
}

// Stopped returns whether training was stopped early.
func (e *EarlyStopping) Stopped() bool {
	return e.stopped
}

// BestValue returns the best observed value of the monitored metric.
func (e *EarlyStopping) BestValue() float64 {
	return e.best
}

// LearningRateScheduler adjusts the learning rate at each epoch.
// Analogous to tf.keras.callbacks.LearningRateScheduler.
type LearningRateScheduler struct {
	schedule  func(epoch int, currentLR float64) float64
	currentLR float64
}

// NewLearningRateScheduler creates a new LR scheduler callback.
// schedule is a function that takes (epoch, currentLR) and returns the new LR.
func NewLearningRateScheduler(initialLR float64, schedule func(epoch int, currentLR float64) float64) *LearningRateScheduler {
	return &LearningRateScheduler{
		schedule:  schedule,
		currentLR: initialLR,
	}
}

// OnEpochEnd updates the learning rate.
func (l *LearningRateScheduler) OnEpochEnd(epoch int, logs map[string]float64) bool {
	l.currentLR = l.schedule(epoch, l.currentLR)
	logs["lr"] = l.currentLR
	return false
}

// Name returns "lr_scheduler".
func (l *LearningRateScheduler) Name() string {
	return "lr_scheduler"
}

// CurrentLR returns the current learning rate.
func (l *LearningRateScheduler) CurrentLR() float64 {
	return l.currentLR
}

// CSVLogger logs training metrics to a CSV file.
// Analogous to tf.keras.callbacks.CSVLogger.
type CSVLogger struct {
	filepath string
	writer   *csv.Writer
	file     *os.File
	headers  []string
	started  bool
}

// NewCSVLogger creates a new CSVLogger callback.
func NewCSVLogger(filepath string) *CSVLogger {
	return &CSVLogger{filepath: filepath}
}

// OnEpochEnd writes metrics to the CSV file.
func (c *CSVLogger) OnEpochEnd(epoch int, logs map[string]float64) bool {
	if !c.started {
		var err error
		c.file, err = os.Create(c.filepath)
		if err != nil {
			return false
		}
		c.writer = csv.NewWriter(c.file)

		c.headers = []string{"epoch"}
		for k := range logs {
			c.headers = append(c.headers, k)
		}
		_ = c.writer.Write(c.headers)
		c.started = true
	}

	row := []string{strconv.Itoa(epoch)}
	for _, h := range c.headers[1:] {
		row = append(row, fmt.Sprintf("%.6f", logs[h]))
	}
	_ = c.writer.Write(row)
	c.writer.Flush()
	return false
}

// Name returns "csv_logger".
func (c *CSVLogger) Name() string {
	return "csv_logger"
}

// Close closes the CSV file.
func (c *CSVLogger) Close() error {
	if c.file != nil {
		c.writer.Flush()
		return c.file.Close()
	}
	return nil
}

// ReduceLROnPlateau reduces learning rate when a metric has stopped improving.
// Analogous to tf.keras.callbacks.ReduceLROnPlateau.
type ReduceLROnPlateau struct {
	monitor   string
	factor    float64
	patience  int
	minLR     float64
	mode      string
	best      float64
	wait      int
	currentLR float64
}

// NewReduceLROnPlateau creates a new ReduceLROnPlateau callback.
func NewReduceLROnPlateau(monitor string, factor float64, patience int, minLR float64, mode string, initialLR float64) *ReduceLROnPlateau {
	best := math.Inf(1)
	if mode == "max" {
		best = math.Inf(-1)
	}
	return &ReduceLROnPlateau{
		monitor:   monitor,
		factor:    factor,
		patience:  patience,
		minLR:     minLR,
		mode:      mode,
		best:      best,
		currentLR: initialLR,
	}
}

// OnEpochEnd checks if LR should be reduced.
func (r *ReduceLROnPlateau) OnEpochEnd(epoch int, logs map[string]float64) bool {
	val, ok := logs[r.monitor]
	if !ok {
		return false
	}

	improved := false
	if r.mode == "min" {
		improved = val < r.best
	} else {
		improved = val > r.best
	}

	if improved {
		r.best = val
		r.wait = 0
	} else {
		r.wait++
		if r.wait >= r.patience {
			newLR := r.currentLR * r.factor
			if newLR < r.minLR {
				newLR = r.minLR
			}
			r.currentLR = newLR
			logs["lr"] = r.currentLR
			r.wait = 0
		}
	}
	return false
}

// Name returns "reduce_lr_on_plateau".
func (r *ReduceLROnPlateau) Name() string {
	return "reduce_lr_on_plateau"
}

// CurrentLR returns the current learning rate.
func (r *ReduceLROnPlateau) CurrentLR() float64 {
	return r.currentLR
}
