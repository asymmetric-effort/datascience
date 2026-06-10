// Package data provides a data pipeline for feeding arrays to models,
// analogous to tf.data.Dataset.
package data

import (
	"fmt"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// maxBufferSize limits internal buffer sizes in the pipeline.
const maxBufferSize = 1 << 24 // 16M elements

// Dataset represents a pipeline of array batches.
// Analogous to tf.data.Dataset.
type Dataset struct {
	data        []*numgo.NDArray
	batchSize   int
	repeatCount int // -1 = infinite, 0 = no repeat
	shuffleSeed int64
	shuffled    bool
	mapFn       func(*numgo.NDArray) (*numgo.NDArray, error)
	takeCount   int // -1 = take all
}

// FromSlices creates a Dataset from an NDArray by slicing along axis 0.
// Each element in the dataset is a single sample.
// Analogous to tf.data.Dataset.from_tensor_slices.
func FromSlices(a *numgo.NDArray) (*Dataset, error) {
	shape := a.Shape()
	if len(shape) < 1 {
		return nil, fmt.Errorf("from_slices: array must have ndim >= 1")
	}
	numSamples := shape[0]
	elemShape := make([]int, len(shape)-1)
	for i := 1; i < len(shape); i++ {
		elemShape[i-1] = shape[i]
	}

	elemSize := 1
	for _, d := range elemShape {
		elemSize *= d
	}

	data := a.Data()
	slices := make([]*numgo.NDArray, numSamples)
	for i := range numSamples {
		elemData := make([]float64, elemSize)
		copy(elemData, data[i*elemSize:(i+1)*elemSize])
		slices[i] = numgo.NewNDArray(elemShape, elemData)
	}

	return &Dataset{
		data:      slices,
		batchSize: 1,
		takeCount: -1,
	}, nil
}

// FromNDArraySlice creates a Dataset from a slice of NDArrays directly.
func FromNDArraySlice(arrays []*numgo.NDArray) *Dataset {
	return &Dataset{
		data:      arrays,
		batchSize: 1,
		takeCount: -1,
	}
}

// Batch sets the batch size for the dataset.
// Analogous to tf.data.Dataset.batch.
func (d *Dataset) Batch(batchSize int) *Dataset {
	return &Dataset{
		data:        d.data,
		batchSize:   batchSize,
		repeatCount: d.repeatCount,
		shuffleSeed: d.shuffleSeed,
		shuffled:    d.shuffled,
		mapFn:       d.mapFn,
		takeCount:   d.takeCount,
	}
}

// Shuffle marks the dataset for shuffling with the given seed.
// Analogous to tf.data.Dataset.shuffle.
func (d *Dataset) Shuffle(seed int64) *Dataset {
	return &Dataset{
		data:        d.data,
		batchSize:   d.batchSize,
		repeatCount: d.repeatCount,
		shuffleSeed: seed,
		shuffled:    true,
		mapFn:       d.mapFn,
		takeCount:   d.takeCount,
	}
}

// Repeat sets the dataset to repeat the given number of times.
// count=-1 means repeat indefinitely (use Take to limit).
// Analogous to tf.data.Dataset.repeat.
func (d *Dataset) Repeat(count int) *Dataset {
	return &Dataset{
		data:        d.data,
		batchSize:   d.batchSize,
		repeatCount: count,
		shuffleSeed: d.shuffleSeed,
		shuffled:    d.shuffled,
		mapFn:       d.mapFn,
		takeCount:   d.takeCount,
	}
}

// Map applies a transformation function to each element.
// Analogous to tf.data.Dataset.map.
func (d *Dataset) Map(fn func(*numgo.NDArray) (*numgo.NDArray, error)) *Dataset {
	return &Dataset{
		data:        d.data,
		batchSize:   d.batchSize,
		repeatCount: d.repeatCount,
		shuffleSeed: d.shuffleSeed,
		shuffled:    d.shuffled,
		mapFn:       fn,
		takeCount:   d.takeCount,
	}
}

// Take limits the dataset to the first count elements.
// Analogous to tf.data.Dataset.take.
func (d *Dataset) Take(count int) *Dataset {
	return &Dataset{
		data:        d.data,
		batchSize:   d.batchSize,
		repeatCount: d.repeatCount,
		shuffleSeed: d.shuffleSeed,
		shuffled:    d.shuffled,
		mapFn:       d.mapFn,
		takeCount:   count,
	}
}

// Len returns the number of raw elements in the dataset.
func (d *Dataset) Len() int {
	return len(d.data)
}

// Iterator creates an iterator over the dataset.
func (d *Dataset) Iterator() *Iterator {
	elements := make([]*numgo.NDArray, len(d.data))
	copy(elements, d.data)

	if d.shuffled {
		// Fisher-Yates shuffle with xorshift.
		state := uint64(d.shuffleSeed)
		if state == 0 {
			state = 0xdeadbeef
		}
		for i := len(elements) - 1; i > 0; i-- {
			state ^= state << 13
			state ^= state >> 7
			state ^= state << 17
			j := int(state % uint64(i+1))
			elements[i], elements[j] = elements[j], elements[i]
		}
	}

	return &Iterator{
		data:        elements,
		batchSize:   d.batchSize,
		repeatCount: d.repeatCount,
		mapFn:       d.mapFn,
		takeCount:   d.takeCount,
		pos:         0,
		repeatsLeft: d.repeatCount,
		taken:       0,
	}
}

// Iterator iterates over dataset batches.
type Iterator struct {
	data        []*numgo.NDArray
	batchSize   int
	repeatCount int
	mapFn       func(*numgo.NDArray) (*numgo.NDArray, error)
	takeCount   int
	pos         int
	repeatsLeft int
	taken       int
}

// Next returns the next batch of arrays.
// Returns nil, nil when the iterator is exhausted.
func (it *Iterator) Next() ([]*numgo.NDArray, error) {
	if it.takeCount >= 0 && it.taken >= it.takeCount {
		return nil, nil
	}

	if it.pos >= len(it.data) {
		if it.repeatsLeft == 0 {
			return nil, nil
		}
		if it.repeatsLeft > 0 {
			it.repeatsLeft--
		}
		it.pos = 0
	}

	end := it.pos + it.batchSize
	if end > len(it.data) {
		end = len(it.data)
	}

	batch := make([]*numgo.NDArray, end-it.pos)
	for i, arr := range it.data[it.pos:end] {
		if it.mapFn != nil {
			mapped, err := it.mapFn(arr)
			if err != nil {
				return nil, err
			}
			batch[i] = mapped
		} else {
			batch[i] = arr
		}
	}

	it.pos = end
	it.taken += len(batch)
	return batch, nil
}
