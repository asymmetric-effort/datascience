// Package gpu provides a compute backend abstraction for accelerated
// factor operations. The default implementation uses pure Go on CPU.
// Future CGO backends (CUDA, OpenCL) can be plugged in via the Backend interface.
package gpu
