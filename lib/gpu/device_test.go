//go:build unit

package gpu

import (
	"testing"
)

func TestDeviceTypeString(t *testing.T) {
	tests := []struct {
		dt   DeviceType
		want string
	}{
		{CPU, "CPU"},
		{CUDA, "CUDA"},
		{OpenCL, "OpenCL"},
		{MPS, "MPS"},
		{DeviceType(99), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.dt.String(); got != tt.want {
			t.Errorf("DeviceType(%d).String() = %q, want %q", tt.dt, got, tt.want)
		}
	}
}

func TestDetectDevices(t *testing.T) {
	devices := DetectDevices()
	if len(devices) == 0 {
		t.Fatal("DetectDevices should return at least one device (CPU)")
	}
	found := false
	for _, d := range devices {
		if d.Type == CPU && d.Available {
			found = true
			break
		}
	}
	if !found {
		t.Error("DetectDevices should include an available CPU device")
	}
}

func TestDetectDevicesCPUFields(t *testing.T) {
	devices := DetectDevices()
	cpu := devices[0]
	if cpu.Type != CPU {
		t.Errorf("expected CPU type, got %v", cpu.Type)
	}
	if cpu.Name != "cpu" {
		t.Errorf("expected name 'cpu', got %q", cpu.Name)
	}
	if cpu.Index != 0 {
		t.Errorf("expected index 0, got %d", cpu.Index)
	}
	if !cpu.Available {
		t.Error("CPU device should be available")
	}
}

func TestSelectBackend(t *testing.T) {
	b := SelectBackend()
	if b == nil {
		t.Fatal("SelectBackend should not return nil")
	}
	if b.Name() != "cpu" {
		t.Errorf("SelectBackend should return CPU for now, got %q", b.Name())
	}
}
