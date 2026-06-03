package gpu

// DeviceType represents the type of compute device.
type DeviceType int

const (
	// CPU represents a CPU compute device.
	CPU DeviceType = iota
	// CUDA represents an NVIDIA CUDA GPU device.
	CUDA
	// OpenCL represents an OpenCL compute device.
	OpenCL
	// MPS represents an Apple Metal Performance Shaders device.
	MPS
)

// String returns the human-readable name of the device type.
func (d DeviceType) String() string {
	switch d {
	case CPU:
		return "CPU"
	case CUDA:
		return "CUDA"
	case OpenCL:
		return "OpenCL"
	case MPS:
		return "MPS"
	default:
		return "Unknown"
	}
}

// DeviceInfo describes a single compute device.
type DeviceInfo struct {
	// Type is the kind of device (CPU, CUDA, OpenCL, MPS).
	Type DeviceType
	// Name is the human-readable device name.
	Name string
	// Index is the device index within its type.
	Index int
	// Memory is the total device memory in bytes (0 if unknown).
	Memory int64
	// Available indicates whether the device can be used.
	Available bool
}

// DetectDevices returns a list of available compute devices.
// Currently only the CPU device is detected. Future CGO backends
// will add CUDA, OpenCL, and MPS detection here.
func DetectDevices() []DeviceInfo {
	devices := []DeviceInfo{
		{
			Type:      CPU,
			Name:      "cpu",
			Index:     0,
			Memory:    0,
			Available: true,
		},
	}
	// Future: probe for CUDA devices via CGO.
	// Future: probe for OpenCL devices via CGO.
	// Future: probe for MPS devices via CGO (macOS only).
	return devices
}

// SelectBackend returns a Backend for the best available device.
// It prefers CUDA > OpenCL > MPS > CPU. Currently only CPU is implemented.
func SelectBackend() Backend {
	// Future: check for CUDA, OpenCL, MPS availability and return
	// the appropriate backend.
	return NewCPUBackend()
}
