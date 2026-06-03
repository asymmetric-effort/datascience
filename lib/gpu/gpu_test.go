//go:build unit

package gpu

import (
	"testing"
)

func TestGetBackendDefault(t *testing.T) {
	ResetBackend()
	b := GetBackend()
	if b == nil {
		t.Fatal("expected non-nil default backend")
	}
	if b.Name() != "cpu" {
		t.Errorf("expected default backend name 'cpu', got %q", b.Name())
	}
}

// mockBackend is a minimal Backend for testing the registry.
type mockBackend struct {
	name   string
	closed bool
}

func (m *mockBackend) Name() string                                  { return m.name }
func (m *mockBackend) IsAvailable() bool                             { return true }
func (m *mockBackend) MatMul(a, b []float64, mm, k, n int) []float64 { return nil }
func (m *mockBackend) ElementWiseMul(a, b []float64) []float64       { return nil }
func (m *mockBackend) Sum(a []float64) float64                       { return 0 }
func (m *mockBackend) Normalize(a []float64) []float64               { return nil }
func (m *mockBackend) FactorProduct(aV []float64, aS []int, bV []float64, bS []int, rS []int) []float64 {
	return nil
}
func (m *mockBackend) Marginalize(v []float64, s []int, axis int) ([]float64, []int) {
	return nil, nil
}
func (m *mockBackend) Close() error                               { m.closed = true; return nil }
func (m *mockBackend) ElementWiseAdd(a, b []float64) []float64    { return nil }
func (m *mockBackend) ElementWiseSub(a, b []float64) []float64    { return nil }
func (m *mockBackend) ElementWiseDiv(a, b []float64) []float64    { return nil }
func (m *mockBackend) ScalarMul(a []float64, s float64) []float64 { return nil }
func (m *mockBackend) ScalarAdd(a []float64, s float64) []float64 { return nil }
func (m *mockBackend) Exp(a []float64) []float64                  { return nil }
func (m *mockBackend) Log(a []float64) []float64                  { return nil }
func (m *mockBackend) Sqrt(a []float64) []float64                 { return nil }
func (m *mockBackend) Abs(a []float64) []float64                  { return nil }
func (m *mockBackend) Max(a []float64) float64                    { return 0 }
func (m *mockBackend) Min(a []float64) float64                    { return 0 }
func (m *mockBackend) ArgMax(a []float64) int                     { return 0 }
func (m *mockBackend) ArgMin(a []float64) int                     { return 0 }
func (m *mockBackend) Dot(a, b []float64) float64                 { return 0 }
func (m *mockBackend) FactorReduce(v []float64, s []int, axis, index int) ([]float64, []int) {
	return nil, nil
}
func (m *mockBackend) FactorMaximize(v []float64, s []int, axis int) ([]float64, []int) {
	return nil, nil
}
func (m *mockBackend) LogSumExp(a []float64) float64                          { return 0 }
func (m *mockBackend) Softmax(a []float64) []float64                          { return nil }
func (m *mockBackend) BatchMatMul(a, b []float64, bs, mm, k, n int) []float64 { return nil }
func (m *mockBackend) BatchNormalize(a []float64, bs, n int) []float64        { return nil }
func (m *mockBackend) Alloc(size int) []float64                               { return make([]float64, size) }
func (m *mockBackend) Free(data []float64)                                    {}
func (m *mockBackend) CopyToDevice(data []float64) []float64                  { return data }
func (m *mockBackend) CopyFromDevice(data []float64) []float64                { return data }
func (m *mockBackend) DeviceCount() int                                       { return 1 }
func (m *mockBackend) DeviceName(index int) string                            { return m.name }
func (m *mockBackend) MemoryUsed() int64                                      { return 0 }
func (m *mockBackend) MemoryTotal() int64                                     { return 0 }

func TestSetBackend(t *testing.T) {
	ResetBackend()
	mock := &mockBackend{name: "mock-gpu"}
	SetBackend(mock)
	b := GetBackend()
	if b.Name() != "mock-gpu" {
		t.Errorf("expected 'mock-gpu', got %q", b.Name())
	}
	// Restore default for other tests.
	ResetBackend()
}

func TestResetBackend(t *testing.T) {
	mock := &mockBackend{name: "temp"}
	SetBackend(mock)
	ResetBackend()
	b := GetBackend()
	if b.Name() != "cpu" {
		t.Errorf("expected 'cpu' after reset, got %q", b.Name())
	}
}

func TestSetBackendOverwrite(t *testing.T) {
	ResetBackend()
	m1 := &mockBackend{name: "first"}
	m2 := &mockBackend{name: "second"}
	SetBackend(m1)
	SetBackend(m2)
	b := GetBackend()
	if b.Name() != "second" {
		t.Errorf("expected 'second', got %q", b.Name())
	}
	ResetBackend()
}
