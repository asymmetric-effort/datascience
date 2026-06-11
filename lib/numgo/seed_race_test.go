//go:build unit

package numgo

import (
	"sync"
	"testing"
)

func TestSeedConcurrentNoRace(t *testing.T) {
	// Under -race, this test verifies that concurrent Seed() calls
	// do not trigger the race detector.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			Seed(seed)
		}(int64(i))
	}
	wg.Wait()
}

func TestSeedThenLoad(t *testing.T) {
	Seed(42)
	rng := defaultRNG.Load()
	if rng == nil {
		t.Fatal("defaultRNG.Load() returned nil after Seed()")
	}
}
