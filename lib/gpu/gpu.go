package gpu

// defaultBackend is the active compute backend. It starts as a CPUBackend.
var defaultBackend Backend = NewCPUBackend()

// SetBackend sets the active compute backend.
func SetBackend(b Backend) {
	defaultBackend = b
}

// GetBackend returns the active compute backend.
func GetBackend() Backend {
	return defaultBackend
}

// ResetBackend restores the default CPU backend.
func ResetBackend() {
	defaultBackend = NewCPUBackend()
}
