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
func (m *mockBackend) Close() error { m.closed = true; return nil }

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
