package data

import (
	"fmt"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func makeArray(t *testing.T, dims []int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray(dims, data)
}

func TestFromSlices(t *testing.T) {
	arr := makeArray(t, []int{5, 3}, make([]float64, 15))
	ds, err := FromSlices(arr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ds.Len() != 5 {
		t.Errorf("Len() = %d, want 5", ds.Len())
	}
}

func TestFromSlicesScalar(t *testing.T) {
	arr := numgo.Zeros() // scalar: shape=[], size=1
	_, err := FromSlices(arr)
	if err == nil {
		t.Error("expected error for scalar array")
	}
}

func TestFromNDArraySlice(t *testing.T) {
	slices := []*numgo.NDArray{
		makeArray(t, []int{3}, []float64{1, 2, 3}),
		makeArray(t, []int{3}, []float64{4, 5, 6}),
	}
	ds := FromNDArraySlice(slices)
	if ds.Len() != 2 {
		t.Errorf("Len() = %d, want 2", ds.Len())
	}
}

func TestDatasetBatch(t *testing.T) {
	arr := makeArray(t, []int{6}, []float64{1, 2, 3, 4, 5, 6})
	ds, _ := FromSlices(arr)
	batched := ds.Batch(2)
	it := batched.Iterator()

	batch1, err := it.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batch1) != 2 {
		t.Errorf("batch1 size = %d, want 2", len(batch1))
	}

	batch2, _ := it.Next()
	if len(batch2) != 2 {
		t.Errorf("batch2 size = %d, want 2", len(batch2))
	}

	batch3, _ := it.Next()
	if len(batch3) != 2 {
		t.Errorf("batch3 size = %d, want 2", len(batch3))
	}

	batch4, _ := it.Next()
	if batch4 != nil {
		t.Error("expected nil after exhausting dataset")
	}
}

func TestDatasetBatchPartial(t *testing.T) {
	arr := makeArray(t, []int{5}, []float64{1, 2, 3, 4, 5})
	ds, _ := FromSlices(arr)
	it := ds.Batch(3).Iterator()

	batch1, _ := it.Next()
	if len(batch1) != 3 {
		t.Errorf("batch1 size = %d, want 3", len(batch1))
	}
	batch2, _ := it.Next()
	if len(batch2) != 2 {
		t.Errorf("batch2 size = %d, want 2 (partial)", len(batch2))
	}
}

func TestDatasetShuffle(t *testing.T) {
	arr := makeArray(t, []int{10}, []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	ds, _ := FromSlices(arr)
	it := ds.Shuffle(42).Batch(10).Iterator()
	batch, _ := it.Next()

	// Check that not all elements are in original order.
	inOrder := true
	for i, b := range batch {
		if b.Data()[0] != float64(i) {
			inOrder = false
			break
		}
	}
	if inOrder {
		t.Error("shuffle did not change order")
	}

	// Check all elements present.
	sum := 0.0
	for _, b := range batch {
		sum += b.Data()[0]
	}
	if sum != 45 {
		t.Errorf("sum = %f, want 45", sum)
	}
}

func TestDatasetRepeat(t *testing.T) {
	arr := makeArray(t, []int{2}, []float64{1, 2})
	ds, _ := FromSlices(arr)
	it := ds.Repeat(1).Batch(2).Iterator() // original + 1 repeat = 2 passes

	batch1, _ := it.Next()
	if len(batch1) != 2 {
		t.Errorf("batch1 size = %d, want 2", len(batch1))
	}
	batch2, _ := it.Next()
	if len(batch2) != 2 {
		t.Errorf("batch2 size = %d, want 2", len(batch2))
	}
	batch3, _ := it.Next()
	if batch3 != nil {
		t.Error("expected nil after repeat exhausted")
	}
}

func TestDatasetMap(t *testing.T) {
	arr := makeArray(t, []int{3}, []float64{1, 2, 3})
	ds, _ := FromSlices(arr)
	doubled := ds.Map(func(a *numgo.NDArray) (*numgo.NDArray, error) {
		data := a.Data()
		for i := range data {
			data[i] *= 2
		}
		return numgo.NewNDArray(a.Shape(), data), nil
	})
	it := doubled.Batch(3).Iterator()
	batch, _ := it.Next()
	want := []float64{2, 4, 6}
	for i, b := range batch {
		if b.Data()[0] != want[i] {
			t.Errorf("batch[%d] = %f, want %f", i, b.Data()[0], want[i])
		}
	}
}

func TestDatasetMapError(t *testing.T) {
	arr := makeArray(t, []int{2}, []float64{1, 2})
	ds, _ := FromSlices(arr)
	errDs := ds.Map(func(a *numgo.NDArray) (*numgo.NDArray, error) {
		return nil, fmt.Errorf("map error")
	})
	it := errDs.Iterator()
	_, err := it.Next()
	if err == nil {
		t.Error("expected error from map function")
	}
}

func TestDatasetTake(t *testing.T) {
	arr := makeArray(t, []int{10}, make([]float64, 10))
	ds, _ := FromSlices(arr)
	it := ds.Take(3).Iterator()
	batch1, _ := it.Next()
	if len(batch1) != 1 {
		t.Errorf("batch1 size = %d, want 1", len(batch1))
	}
	batch2, _ := it.Next()
	if len(batch2) != 1 {
		t.Errorf("batch2 size = %d, want 1", len(batch2))
	}
	batch3, _ := it.Next()
	if len(batch3) != 1 {
		t.Errorf("batch3 size = %d, want 1", len(batch3))
	}
	batch4, _ := it.Next()
	if batch4 != nil {
		t.Error("expected nil after take limit")
	}
}

func TestDatasetChained(t *testing.T) {
	arr := makeArray(t, []int{6}, []float64{10, 20, 30, 40, 50, 60})
	ds, _ := FromSlices(arr)
	pipeline := ds.Shuffle(99).Batch(2).Take(4)
	it := pipeline.Iterator()

	count := 0
	for {
		batch, err := it.Next()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if batch == nil {
			break
		}
		count += len(batch)
	}
	if count != 4 {
		t.Errorf("got %d elements, want 4", count)
	}
}
