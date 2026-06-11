package model

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	config := &ModelConfig{
		Name:    "test_model",
		Version: "1.0.0",
		Layers: []LayerConfig{
			{Type: "Dense", Params: map[string]any{"units": float64(64)}},
			{Type: "Dense", Params: map[string]any{"units": float64(10)}},
		},
		LearningRate: 0.001,
		LossFunction: "mse",
	}

	err := SaveConfig(config, path)
	if err != nil {
		t.Fatalf("SaveConfig error: %v", err)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if loaded.Name != "test_model" {
		t.Errorf("Name = %q, want %q", loaded.Name, "test_model")
	}
	if len(loaded.Layers) != 2 {
		t.Errorf("Layers count = %d, want 2", len(loaded.Layers))
	}
	if loaded.LearningRate != 0.001 {
		t.Errorf("LR = %f, want 0.001", loaded.LearningRate)
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSaveLoadWeights(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "weights.bin")

	weights := []float64{1.5, -2.3, 0.0, math.Pi, math.E}
	err := SaveWeights(weights, path)
	if err != nil {
		t.Fatalf("SaveWeights error: %v", err)
	}

	loaded, err := LoadWeights(path)
	if err != nil {
		t.Fatalf("LoadWeights error: %v", err)
	}
	if len(loaded) != len(weights) {
		t.Fatalf("loaded %d weights, want %d", len(loaded), len(weights))
	}
	for i, v := range loaded {
		if v != weights[i] {
			t.Errorf("weights[%d] = %f, want %f", i, v, weights[i])
		}
	}
}

func TestLoadWeightsNotFound(t *testing.T) {
	_, err := LoadWeights("/nonexistent/path.bin")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadWeightsTruncated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.bin")
	// Write count=10 but only 1 weight.
	f, _ := os.Create(path)
	data := make([]byte, 8+8) // count + 1 weight
	data[0] = 10              // count = 10
	f.Write(data)
	f.Close()

	_, err := LoadWeights(path)
	if err == nil {
		t.Error("expected error for truncated file")
	}
}

func TestSaveLoadModel(t *testing.T) {
	dir := t.TempDir()
	modelDir := filepath.Join(dir, "my_model")

	config := &ModelConfig{
		Name:    "test",
		Version: "0.1.0",
		Layers:  []LayerConfig{{Type: "Dense", Params: map[string]any{"units": float64(3)}}},
	}
	weights := []float64{1, 2, 3, 4, 5}

	err := SaveModel(config, weights, modelDir)
	if err != nil {
		t.Fatalf("SaveModel error: %v", err)
	}

	loadedConfig, loadedWeights, err := LoadModel(modelDir)
	if err != nil {
		t.Fatalf("LoadModel error: %v", err)
	}
	if loadedConfig.Name != "test" {
		t.Errorf("Name = %q, want %q", loadedConfig.Name, "test")
	}
	if len(loadedWeights) != 5 {
		t.Errorf("loaded %d weights, want 5", len(loadedWeights))
	}
}

func TestLoadModelMissing(t *testing.T) {
	_, _, err := LoadModel("/nonexistent/dir")
	if err == nil {
		t.Error("expected error for missing dir")
	}
}

func TestSaveLoadCheckpoint(t *testing.T) {
	dir := t.TempDir()
	ckpt := &Checkpoint{Epoch: 5, Loss: 0.123}
	weights := []float64{1, 2, 3}

	err := SaveCheckpoint(ckpt, weights, dir)
	if err != nil {
		t.Fatalf("SaveCheckpoint error: %v", err)
	}

	loaded, err := LoadCheckpoint(filepath.Join(dir, "checkpoint_epoch_5.json"))
	if err != nil {
		t.Fatalf("LoadCheckpoint error: %v", err)
	}
	if loaded.Epoch != 5 {
		t.Errorf("Epoch = %d, want 5", loaded.Epoch)
	}
	if loaded.Loss != 0.123 {
		t.Errorf("Loss = %f, want 0.123", loaded.Loss)
	}
	if loaded.Weights != "weights_epoch_5.bin" {
		t.Errorf("Weights = %q", loaded.Weights)
	}

	// Verify weights file.
	w, err := LoadWeights(filepath.Join(dir, loaded.Weights))
	if err != nil {
		t.Fatalf("LoadWeights error: %v", err)
	}
	if len(w) != 3 || w[0] != 1 {
		t.Errorf("loaded weights = %v, want [1 2 3]", w)
	}
}

func TestLoadCheckpointNotFound(t *testing.T) {
	_, err := LoadCheckpoint("/nonexistent/ckpt.json")
	if err == nil {
		t.Error("expected error for missing checkpoint")
	}
}

func TestLoadCheckpointInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadCheckpoint(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestFloat64Conversion(t *testing.T) {
	values := []float64{0, 1, -1, math.Pi, math.Inf(1), math.Inf(-1)}
	for _, v := range values {
		b := Float64ToBytes(v)
		got := BytesToFloat64(b)
		if v != got && !math.IsNaN(v) {
			t.Errorf("roundtrip %f failed: got %f", v, got)
		}
	}
}

func TestSaveConfigBadPath(t *testing.T) {
	config := &ModelConfig{Name: "test"}
	err := SaveConfig(config, "/nonexistent/dir/config.json")
	if err == nil {
		t.Error("expected error for bad path")
	}
}

func TestSaveWeightsBadPath(t *testing.T) {
	err := SaveWeights([]float64{1, 2}, "/nonexistent/dir/w.bin")
	if err == nil {
		t.Error("expected error for bad path")
	}
}

func TestLoadWeightsEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.bin")
	os.WriteFile(path, []byte{}, 0o644)
	_, err := LoadWeights(path)
	if err == nil {
		t.Error("expected error for empty file")
	}
}

func TestSaveModelBadDir(t *testing.T) {
	config := &ModelConfig{Name: "test"}
	// Use a path that exists as a file, not a dir.
	dir := t.TempDir()
	filePath := filepath.Join(dir, "afile")
	os.WriteFile(filePath, []byte("x"), 0o644)
	err := SaveModel(config, []float64{1}, filepath.Join(filePath, "subdir"))
	if err == nil {
		t.Error("expected error for bad dir")
	}
}

func TestLoadModelMissingWeights(t *testing.T) {
	dir := t.TempDir()
	// Create config but no weights.
	config := &ModelConfig{Name: "test"}
	SaveConfig(config, filepath.Join(dir, "config.json"))
	_, _, err := LoadModel(dir)
	if err == nil {
		t.Error("expected error for missing weights")
	}
}

func TestSaveCheckpointBadDir(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "afile")
	os.WriteFile(filePath, []byte("x"), 0o644)
	ckpt := &Checkpoint{Epoch: 1, Loss: 0.5}
	err := SaveCheckpoint(ckpt, []float64{1}, filepath.Join(filePath, "subdir"))
	if err == nil {
		t.Error("expected error for bad dir")
	}
}

func TestSaveWeightsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.bin")
	err := SaveWeights([]float64{}, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := LoadWeights(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("expected empty weights, got %d", len(loaded))
	}
}

func TestLoadWeightsCraftedCountRejectsOversizedCount(t *testing.T) {
	// Craft a weights file where the count header claims far more elements
	// than the file actually contains. This must return an error, not OOM.
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.bin")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	// Write a count of 1 billion (would need ~8GB) but file is only 16 bytes.
	var count uint64 = 1_000_000_000
	binary.Write(f, binary.LittleEndian, count)
	binary.Write(f, binary.LittleEndian, float64(1.0)) // one real weight
	f.Close()

	_, err = LoadWeights(path)
	if err == nil {
		t.Fatal("LoadWeights should reject file where count exceeds file capacity")
	}
}

func TestLoadWeightsValidCountAtBoundary(t *testing.T) {
	// A file with count=2 and exactly 2 float64s should load fine.
	dir := t.TempDir()
	path := filepath.Join(dir, "ok.bin")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	var count uint64 = 2
	binary.Write(f, binary.LittleEndian, count)
	binary.Write(f, binary.LittleEndian, float64(3.14))
	binary.Write(f, binary.LittleEndian, float64(2.72))
	f.Close()

	weights, err := LoadWeights(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(weights) != 2 {
		t.Fatalf("expected 2 weights, got %d", len(weights))
	}
	if weights[0] != 3.14 || weights[1] != 2.72 {
		t.Errorf("weights = %v, want [3.14, 2.72]", weights)
	}
}

func TestPathTraversalRejected(t *testing.T) {
	cases := []struct {
		name string
		fn   func() error
	}{
		{"SaveConfig", func() error { return SaveConfig(&ModelConfig{}, "../../evil.json") }},
		{"LoadConfig", func() error { _, err := LoadConfig("../../evil.json"); return err }},
		{"SaveWeights", func() error { return SaveWeights(nil, "../../evil.bin") }},
		{"LoadWeights", func() error { _, err := LoadWeights("../../evil.bin"); return err }},
		{"SaveModel", func() error { return SaveModel(&ModelConfig{}, nil, "../../evil") }},
		{"LoadModel", func() error { _, _, err := LoadModel("../../evil"); return err }},
		{"SaveCheckpoint", func() error {
			return SaveCheckpoint(&Checkpoint{Epoch: 1}, nil, "../../evil")
		}},
		{"LoadCheckpoint", func() error { _, err := LoadCheckpoint("../../evil.json"); return err }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn()
			if err == nil {
				t.Fatal("expected error for traversal path")
			}
			if !strings.Contains(err.Error(), "directory traversal") {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidPathsAccepted(t *testing.T) {
	dir := t.TempDir()

	// SaveModel with a valid path should succeed.
	config := &ModelConfig{Name: "test", Layers: []LayerConfig{{Type: "dense"}}}
	err := SaveModel(config, []float64{1.0, 2.0}, dir)
	if err != nil {
		t.Fatalf("SaveModel to valid path failed: %v", err)
	}

	// LoadModel should succeed.
	loaded, weights, err := LoadModel(dir)
	if err != nil {
		t.Fatalf("LoadModel from valid path failed: %v", err)
	}
	if loaded.Name != "test" {
		t.Errorf("loaded config name = %q, want %q", loaded.Name, "test")
	}
	if len(weights) != 2 {
		t.Errorf("loaded %d weights, want 2", len(weights))
	}
}
