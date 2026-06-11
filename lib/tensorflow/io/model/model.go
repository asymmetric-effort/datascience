// Package model provides model serialization (save/load) for go-tensorflow,
// analogous to tf.keras.models.save_model / load_model.
package model

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/asymmetric-effort/datascience/internal/safepath"
)

// maxWeightFileSize limits the maximum weight file size (1GB).
const maxWeightFileSize = 1 << 30

// LayerConfig describes a single layer's configuration for serialization.
type LayerConfig struct {
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

// ModelConfig describes the full model architecture for serialization.
type ModelConfig struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Layers       []LayerConfig `json:"layers"`
	LearningRate float64       `json:"learning_rate,omitempty"`
	LossFunction string        `json:"loss_function,omitempty"`
}

// SaveConfig writes a model configuration to a JSON file.
func SaveConfig(config *ModelConfig, path string) error {
	path, err := safepath.Validate(path)
	if err != nil {
		return fmt.Errorf("model: %w", err)
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadConfig reads a model configuration from a JSON file.
func LoadConfig(path string) (*ModelConfig, error) {
	path, err := safepath.Validate(path)
	if err != nil {
		return nil, fmt.Errorf("model: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var config ModelConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &config, nil
}

// SaveWeights writes model weights to a binary file.
// Format: [uint64 count][float64 values...]
func SaveWeights(weights []float64, path string) error {
	path, err := safepath.Validate(path)
	if err != nil {
		return fmt.Errorf("model: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create weights file: %w", err)
	}
	defer f.Close()

	// Write count.
	if err := binary.Write(f, binary.LittleEndian, uint64(len(weights))); err != nil {
		return fmt.Errorf("write count: %w", err)
	}
	// Write weights.
	for _, w := range weights {
		if err := binary.Write(f, binary.LittleEndian, w); err != nil {
			return fmt.Errorf("write weight: %w", err)
		}
	}
	return nil
}

// LoadWeights reads model weights from a binary file.
func LoadWeights(path string) ([]float64, error) {
	path, err := safepath.Validate(path)
	if err != nil {
		return nil, fmt.Errorf("model: %w", err)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open weights file: %w", err)
	}
	defer f.Close()

	// Check file size.
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat weights file: %w", err)
	}
	if info.Size() > maxWeightFileSize {
		return nil, fmt.Errorf("weights file too large: %d bytes (max %d)", info.Size(), maxWeightFileSize)
	}

	var count uint64
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("read count: %w", err)
	}

	// Validate count against actual file size to prevent OOM from crafted files.
	// File layout: 8 bytes (count) + count*8 bytes (float64 values).
	maxCount := uint64(info.Size()-8) / 8
	if count > maxCount {
		return nil, fmt.Errorf("weight count %d exceeds file capacity (%d bytes)", count, info.Size())
	}

	weights := make([]float64, count)
	for i := range weights {
		if err := binary.Read(f, binary.LittleEndian, &weights[i]); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("unexpected EOF at weight %d of %d", i, count)
			}
			return nil, fmt.Errorf("read weight %d: %w", i, err)
		}
	}
	return weights, nil
}

// SaveModel saves both config and weights to a directory.
func SaveModel(config *ModelConfig, weights []float64, dir string) error {
	dir, err := safepath.Validate(dir)
	if err != nil {
		return fmt.Errorf("model: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	if err := SaveConfig(config, filepath.Join(dir, "config.json")); err != nil {
		return err
	}
	if err := SaveWeights(weights, filepath.Join(dir, "weights.bin")); err != nil {
		return err
	}
	return nil
}

// LoadModel loads both config and weights from a directory.
func LoadModel(dir string) (*ModelConfig, []float64, error) {
	dir, err := safepath.Validate(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("model: %w", err)
	}
	config, err := LoadConfig(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil, nil, err
	}
	weights, err := LoadWeights(filepath.Join(dir, "weights.bin"))
	if err != nil {
		return nil, nil, err
	}
	return config, weights, nil
}

// Checkpoint saves a training checkpoint with epoch and loss info.
type Checkpoint struct {
	Epoch   int     `json:"epoch"`
	Loss    float64 `json:"loss"`
	Weights string  `json:"weights_file"`
}

// SaveCheckpoint saves a checkpoint to disk.
func SaveCheckpoint(ckpt *Checkpoint, weights []float64, dir string) error {
	dir, err := safepath.Validate(dir)
	if err != nil {
		return fmt.Errorf("model: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	ckpt.Weights = fmt.Sprintf("weights_epoch_%d.bin", ckpt.Epoch)
	weightsPath := filepath.Join(dir, ckpt.Weights)
	if err := SaveWeights(weights, weightsPath); err != nil {
		return err
	}

	data, err := json.MarshalIndent(ckpt, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal checkpoint: %w", err)
	}
	ckptPath := filepath.Join(dir, fmt.Sprintf("checkpoint_epoch_%d.json", ckpt.Epoch))
	return os.WriteFile(ckptPath, data, 0o644)
}

// LoadCheckpoint loads a checkpoint from disk.
func LoadCheckpoint(path string) (*Checkpoint, error) {
	path, err := safepath.Validate(path)
	if err != nil {
		return nil, fmt.Errorf("model: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read checkpoint: %w", err)
	}
	var ckpt Checkpoint
	if err := json.Unmarshal(data, &ckpt); err != nil {
		return nil, fmt.Errorf("unmarshal checkpoint: %w", err)
	}
	return &ckpt, nil
}

// Float64ToBytes converts a float64 to its 8-byte binary representation.
func Float64ToBytes(f float64) [8]byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf
}

// BytesToFloat64 converts 8 bytes to a float64.
func BytesToFloat64(b [8]byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(b[:]))
}
