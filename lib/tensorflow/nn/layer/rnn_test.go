package layer

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// --- Embedding Tests ---

func TestEmbedding(t *testing.T) {
	emb, err := NewEmbedding(10, 4, deterministicRNG())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if emb.VocabSize() != 10 {
		t.Errorf("VocabSize() = %d, want 10", emb.VocabSize())
	}
	if emb.EmbeddingDim() != 4 {
		t.Errorf("EmbeddingDim() = %d, want 4", emb.EmbeddingDim())
	}

	// Input: batch=2, seqLen=3, indices [0,1,2], [3,4,5]
	input := makeTestTensor(t, []int{2, 3}, []float64{0, 1, 2, 3, 4, 5})
	result, err := emb.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{2, 3, 4}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestEmbeddingOutOfRange(t *testing.T) {
	emb, _ := NewEmbedding(5, 3, deterministicRNG())
	input := makeTestTensor(t, []int{1, 2}, []float64{0, 10}) // 10 > vocabSize
	_, err := emb.Forward(input)
	if err == nil {
		t.Error("expected error for out of range index")
	}
}

func TestEmbeddingWrongRank(t *testing.T) {
	emb, _ := NewEmbedding(5, 3, deterministicRNG())
	input := makeTestTensor(t, []int{6}, []float64{0, 1, 2, 3, 4, 0})
	_, err := emb.Forward(input)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestEmbeddingInvalidParams(t *testing.T) {
	_, err := NewEmbedding(0, 4, deterministicRNG())
	if err == nil {
		t.Error("expected error for zero vocab size")
	}
}

func TestEmbeddingWeights(t *testing.T) {
	emb, _ := NewEmbedding(5, 3, deterministicRNG())
	w := emb.Weights()
	expectedShape := []int{5, 3}
	if !shapeEq(w.Shape(), expectedShape) {
		t.Errorf("weights shape = %v, want %v", w.Shape(), expectedShape)
	}
}

// --- LSTM Tests ---

func TestLSTM(t *testing.T) {
	lstm, err := NewLSTM(4, 8, false, deterministicRNG())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lstm.Units() != 8 {
		t.Errorf("Units() = %d, want 8", lstm.Units())
	}
	if lstm.ReturnSequences() {
		t.Error("expected ReturnSequences=false")
	}

	// Input: batch=2, timeSteps=3, features=4
	inData := make([]float64, 24)
	for i := range 24 {
		inData[i] = float64(i) * 0.1
	}
	input := numgo.NewNDArray([]int{2, 3, 4}, inData)
	result, err := lstm.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{2, 8}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
	// Values should be bounded by tanh.
	for i, v := range result.Data() {
		if math.Abs(v) > 1.0 {
			t.Errorf("data[%d] = %f, outside [-1, 1]", i, v)
		}
	}
}

func TestLSTMReturnSequences(t *testing.T) {
	lstm, _ := NewLSTM(4, 8, true, deterministicRNG())
	if !lstm.ReturnSequences() {
		t.Error("expected ReturnSequences=true")
	}
	input := makeTestTensor(t, []int{2, 3, 4}, make([]float64, 24))
	result, err := lstm.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{2, 3, 8}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestLSTMInvalidParams(t *testing.T) {
	_, err := NewLSTM(0, 8, false, deterministicRNG())
	if err == nil {
		t.Error("expected error for zero inputSize")
	}
}

func TestLSTMWrongInput(t *testing.T) {
	lstm, _ := NewLSTM(4, 8, false, deterministicRNG())
	input := makeTestTensor(t, []int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	_, err := lstm.Forward(input)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestLSTMWrongFeatures(t *testing.T) {
	lstm, _ := NewLSTM(4, 8, false, deterministicRNG())
	input := makeTestTensor(t, []int{1, 2, 5}, make([]float64, 10)) // features=5, not 4
	_, err := lstm.Forward(input)
	if err == nil {
		t.Error("expected error for wrong feature count")
	}
}

// --- GRU Tests ---

func TestGRU(t *testing.T) {
	gru, err := NewGRU(4, 8, false, deterministicRNG())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gru.Units() != 8 {
		t.Errorf("Units() = %d, want 8", gru.Units())
	}
	if gru.ReturnSequences() {
		t.Error("expected ReturnSequences=false")
	}

	inData := make([]float64, 24)
	for i := range 24 {
		inData[i] = float64(i) * 0.1
	}
	input := numgo.NewNDArray([]int{2, 3, 4}, inData)
	result, err := gru.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{2, 8}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestGRUReturnSequences(t *testing.T) {
	gru, _ := NewGRU(4, 8, true, deterministicRNG())
	if !gru.ReturnSequences() {
		t.Error("expected ReturnSequences=true")
	}
	input := makeTestTensor(t, []int{2, 3, 4}, make([]float64, 24))
	result, err := gru.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{2, 3, 8}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestGRUInvalidParams(t *testing.T) {
	_, err := NewGRU(0, 8, false, deterministicRNG())
	if err == nil {
		t.Error("expected error for zero inputSize")
	}
}

func TestGRUWrongInput(t *testing.T) {
	gru, _ := NewGRU(4, 8, false, deterministicRNG())
	input := makeTestTensor(t, []int{2, 4}, []float64{1, 2, 3, 4, 5, 6, 7, 8})
	_, err := gru.Forward(input)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

// --- MultiHeadAttention Tests ---

func TestMultiHeadAttention(t *testing.T) {
	mha, err := NewMultiHeadAttention(8, 2, deterministicRNG())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mha.NumHeads() != 2 {
		t.Errorf("NumHeads() = %d, want 2", mha.NumHeads())
	}
	if mha.DModel() != 8 {
		t.Errorf("DModel() = %d, want 8", mha.DModel())
	}

	// Self-attention: batch=1, seqLen=3, dModel=8
	inData := make([]float64, 24)
	for i := range 24 {
		inData[i] = float64(i) * 0.01
	}
	input := numgo.NewNDArray([]int{1, 3, 8}, inData)
	result, err := mha.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{1, 3, 8}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestMultiHeadAttentionCrossAttention(t *testing.T) {
	mha, _ := NewMultiHeadAttention(4, 2, deterministicRNG())
	query := makeTestTensor(t, []int{1, 2, 4}, make([]float64, 8))
	kv := makeTestTensor(t, []int{1, 5, 4}, make([]float64, 20))
	result, err := mha.ForwardQKV(query, kv, kv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{1, 2, 4}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestMultiHeadAttentionInvalidParams(t *testing.T) {
	_, err := NewMultiHeadAttention(0, 2, deterministicRNG())
	if err == nil {
		t.Error("expected error for zero dModel")
	}
	_, err = NewMultiHeadAttention(5, 2, deterministicRNG())
	if err == nil {
		t.Error("expected error for dModel not divisible by numHeads")
	}
}

func TestMultiHeadAttentionWrongRank(t *testing.T) {
	mha, _ := NewMultiHeadAttention(4, 2, deterministicRNG())
	input := makeTestTensor(t, []int{4, 4}, make([]float64, 16))
	_, err := mha.Forward(input)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestMultiHeadAttentionWrongDim(t *testing.T) {
	mha, _ := NewMultiHeadAttention(4, 2, deterministicRNG())
	input := makeTestTensor(t, []int{1, 3, 6}, make([]float64, 18))
	_, err := mha.Forward(input)
	if err == nil {
		t.Error("expected error for wrong feature dim")
	}
}

func TestMultiHeadAttentionBatchMismatch(t *testing.T) {
	mha, _ := NewMultiHeadAttention(4, 2, deterministicRNG())
	q := makeTestTensor(t, []int{1, 2, 4}, make([]float64, 8))
	k := makeTestTensor(t, []int{2, 2, 4}, make([]float64, 16))
	v := makeTestTensor(t, []int{1, 2, 4}, make([]float64, 8))
	_, err := mha.ForwardQKV(q, k, v)
	if err == nil {
		t.Error("expected error for batch mismatch")
	}
}

func TestMultiHeadAttentionKVLenMismatch(t *testing.T) {
	mha, _ := NewMultiHeadAttention(4, 2, deterministicRNG())
	q := makeTestTensor(t, []int{1, 2, 4}, make([]float64, 8))
	k := makeTestTensor(t, []int{1, 3, 4}, make([]float64, 12))
	v := makeTestTensor(t, []int{1, 5, 4}, make([]float64, 20))
	_, err := mha.ForwardQKV(q, k, v)
	if err == nil {
		t.Error("expected error for K/V sequence length mismatch")
	}
}

func TestLSTMNonZeroOutput(t *testing.T) {
	// Ensure LSTM produces non-zero output with non-zero input.
	lstm, _ := NewLSTM(2, 3, false, deterministicRNG())
	input := numgo.NewNDArray([]int{1, 4, 2}, []float64{1, 2, 3, 4, 5, 6, 7, 8})
	result, err := lstm.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hasNonZero := false
	for _, v := range result.Data() {
		if v != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Error("LSTM output is all zeros for non-zero input")
	}
}
