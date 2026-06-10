package callbacks

import (
	"math"
	"os"
	"testing"
)

// --- EarlyStopping ---

func TestEarlyStoppingMin(t *testing.T) {
	es := NewEarlyStopping("loss", 3, 0.001, "min")
	if es.Name() != "early_stopping" {
		t.Errorf("Name() = %q", es.Name())
	}
	// Improving loss.
	es.OnEpochEnd(0, map[string]float64{"loss": 1.0})
	es.OnEpochEnd(1, map[string]float64{"loss": 0.8})
	es.OnEpochEnd(2, map[string]float64{"loss": 0.6})
	if es.Stopped() {
		t.Error("should not have stopped while improving")
	}
	// Stagnant loss for 3 epochs.
	es.OnEpochEnd(3, map[string]float64{"loss": 0.6})
	es.OnEpochEnd(4, map[string]float64{"loss": 0.6})
	stop := es.OnEpochEnd(5, map[string]float64{"loss": 0.6})
	if !stop || !es.Stopped() {
		t.Error("should have stopped after patience exhausted")
	}
	if math.Abs(es.BestValue()-0.6) > 1e-10 {
		t.Errorf("BestValue() = %f, want 0.6", es.BestValue())
	}
}

func TestEarlyStoppingMax(t *testing.T) {
	es := NewEarlyStopping("accuracy", 2, 0, "max")
	es.OnEpochEnd(0, map[string]float64{"accuracy": 0.5})
	es.OnEpochEnd(1, map[string]float64{"accuracy": 0.7})
	es.OnEpochEnd(2, map[string]float64{"accuracy": 0.7})
	stop := es.OnEpochEnd(3, map[string]float64{"accuracy": 0.7})
	if !stop {
		t.Error("should have stopped")
	}
}

func TestEarlyStoppingMissingMetric(t *testing.T) {
	es := NewEarlyStopping("loss", 2, 0, "min")
	stop := es.OnEpochEnd(0, map[string]float64{"accuracy": 0.5})
	if stop {
		t.Error("should not stop when metric is missing")
	}
}

// --- LearningRateScheduler ---

func TestLearningRateScheduler(t *testing.T) {
	schedule := func(epoch int, lr float64) float64 {
		return lr * 0.9
	}
	lrs := NewLearningRateScheduler(0.1, schedule)
	if lrs.Name() != "lr_scheduler" {
		t.Errorf("Name() = %q", lrs.Name())
	}
	logs := map[string]float64{}
	lrs.OnEpochEnd(0, logs)
	if math.Abs(lrs.CurrentLR()-0.09) > 1e-10 {
		t.Errorf("CurrentLR() = %f, want 0.09", lrs.CurrentLR())
	}
	if math.Abs(logs["lr"]-0.09) > 1e-10 {
		t.Errorf("logs[lr] = %f, want 0.09", logs["lr"])
	}
	lrs.OnEpochEnd(1, logs)
	if math.Abs(lrs.CurrentLR()-0.081) > 1e-10 {
		t.Errorf("CurrentLR() = %f, want 0.081", lrs.CurrentLR())
	}
}

// --- CSVLogger ---

func TestCSVLogger(t *testing.T) {
	path := t.TempDir() + "/training.csv"
	logger := NewCSVLogger(path)
	if logger.Name() != "csv_logger" {
		t.Errorf("Name() = %q", logger.Name())
	}
	logger.OnEpochEnd(0, map[string]float64{"loss": 1.5, "accuracy": 0.3})
	logger.OnEpochEnd(1, map[string]float64{"loss": 1.0, "accuracy": 0.5})
	err := logger.Close()
	if err != nil {
		t.Fatalf("Close error: %v", err)
	}

	// Verify file exists and has content.
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(data) == 0 {
		t.Error("CSV file is empty")
	}
}

func TestCSVLoggerCloseWithoutOpen(t *testing.T) {
	logger := NewCSVLogger("/nonexistent/path.csv")
	err := logger.Close()
	if err != nil {
		t.Errorf("Close should not error when file not opened: %v", err)
	}
}

// --- ReduceLROnPlateau ---

func TestReduceLROnPlateau(t *testing.T) {
	r := NewReduceLROnPlateau("loss", 0.5, 2, 0.001, "min", 0.1)
	if r.Name() != "reduce_lr_on_plateau" {
		t.Errorf("Name() = %q", r.Name())
	}
	// Improving.
	r.OnEpochEnd(0, map[string]float64{"loss": 1.0})
	r.OnEpochEnd(1, map[string]float64{"loss": 0.5})
	if math.Abs(r.CurrentLR()-0.1) > 1e-10 {
		t.Errorf("LR should still be 0.1, got %f", r.CurrentLR())
	}
	// Plateau for 2 epochs.
	r.OnEpochEnd(2, map[string]float64{"loss": 0.5})
	logs := map[string]float64{"loss": 0.5}
	r.OnEpochEnd(3, logs)
	if math.Abs(r.CurrentLR()-0.05) > 1e-10 {
		t.Errorf("LR should be 0.05, got %f", r.CurrentLR())
	}
}

func TestReduceLROnPlateauMax(t *testing.T) {
	r := NewReduceLROnPlateau("acc", 0.5, 1, 0.001, "max", 0.1)
	r.OnEpochEnd(0, map[string]float64{"acc": 0.9})
	logs := map[string]float64{"acc": 0.9}
	r.OnEpochEnd(1, logs)
	if math.Abs(r.CurrentLR()-0.05) > 1e-10 {
		t.Errorf("LR should be 0.05, got %f", r.CurrentLR())
	}
}

func TestReduceLROnPlateauMinLR(t *testing.T) {
	r := NewReduceLROnPlateau("loss", 0.1, 1, 0.05, "min", 0.1)
	r.OnEpochEnd(0, map[string]float64{"loss": 1.0})
	r.OnEpochEnd(1, map[string]float64{"loss": 1.0})
	// 0.1 * 0.1 = 0.01 < minLR 0.05, so should clamp to 0.05.
	if r.CurrentLR() < 0.05 {
		t.Errorf("LR should not go below minLR, got %f", r.CurrentLR())
	}
}

func TestReduceLROnPlateauMissing(t *testing.T) {
	r := NewReduceLROnPlateau("loss", 0.5, 1, 0.001, "min", 0.1)
	r.OnEpochEnd(0, map[string]float64{"accuracy": 0.5})
	if math.Abs(r.CurrentLR()-0.1) > 1e-10 {
		t.Errorf("LR should be unchanged when metric missing, got %f", r.CurrentLR())
	}
}

// --- Interface Compliance ---

func TestCallbackInterface(t *testing.T) {
	var _ Callback = &EarlyStopping{}
	var _ Callback = &LearningRateScheduler{}
	var _ Callback = &CSVLogger{}
	var _ Callback = &ReduceLROnPlateau{}
}
