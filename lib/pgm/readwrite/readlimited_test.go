//go:build unit

package readwrite

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadLimitedAcceptsSmallInput(t *testing.T) {
	input := []byte("<small>data</small>")
	data, err := readLimited(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, input) {
		t.Errorf("data mismatch")
	}
}

func TestReadLimitedAcceptsExactlyMaxSize(t *testing.T) {
	input := make([]byte, MaxInputSize)
	for i := range input {
		input[i] = 'x'
	}
	data, err := readLimited(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error for exactly MaxInputSize: %v", err)
	}
	if len(data) != MaxInputSize {
		t.Errorf("expected %d bytes, got %d", MaxInputSize, len(data))
	}
}

func TestReadLimitedRejectsOversizedInput(t *testing.T) {
	input := strings.NewReader(strings.Repeat("x", MaxInputSize+1))
	_, err := readLimited(input)
	if err == nil {
		t.Fatal("expected error for input exceeding MaxInputSize")
	}
	if !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestReadLimitedRejectsLargeStream(t *testing.T) {
	_, err := readLimited(&infiniteReader{})
	if err == nil {
		t.Fatal("expected error for infinite stream")
	}
}

// infiniteReader returns 'A' bytes forever.
type infiniteReader struct{}

func (r *infiniteReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'A'
	}
	return len(p), nil
}

func TestXMLBIFRejectsOversizedInput(t *testing.T) {
	input := strings.NewReader(strings.Repeat("<", MaxInputSize+1))
	_, err := ReadXMLBIF(input)
	if err == nil {
		t.Fatal("ReadXMLBIF should reject oversized input")
	}
}

func TestReadJSONRejectsOversizedInput(t *testing.T) {
	input := strings.NewReader(strings.Repeat("{", MaxInputSize+1))
	_, err := ReadJSON(input)
	if err == nil {
		t.Fatal("ReadJSON should reject oversized input")
	}
}

// --- readLimitedN tests ---

func TestReadLimitedNCustomSize(t *testing.T) {
	input := []byte("hello world")
	data, err := readLimitedN(bytes.NewReader(input), 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, input) {
		t.Errorf("data mismatch")
	}
}

func TestReadLimitedNRejectsAtCustomLimit(t *testing.T) {
	_, err := readLimitedN(strings.NewReader("abcdef"), 5)
	if err == nil {
		t.Fatal("expected error for input exceeding custom limit")
	}
}

func TestReadLimitedNAcceptsExactCustomLimit(t *testing.T) {
	data, err := readLimitedN(strings.NewReader("abcde"), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "abcde" {
		t.Errorf("got %q, want %q", string(data), "abcde")
	}
}

func TestReadLimitedNRejectsNonPositiveLimit(t *testing.T) {
	_, err := readLimitedN(strings.NewReader("x"), 0)
	if err == nil {
		t.Fatal("expected error for maxBytes=0")
	}
	_, err = readLimitedN(strings.NewReader("x"), -1)
	if err == nil {
		t.Fatal("expected error for maxBytes=-1")
	}
}

// --- WithLimit variant tests ---

func TestReadXMLBIFWithLimitAcceptsLargerInput(t *testing.T) {
	// Input larger than MaxInputSize but within custom limit.
	// Not valid XML, so it will fail on parse — but it should NOT fail on size.
	bigInput := strings.Repeat("<x>", MaxInputSize/3+1) // > MaxInputSize
	_, err := ReadXMLBIFWithLimit(strings.NewReader(bigInput), len(bigInput)+1)
	if err == nil {
		t.Skip("unexpectedly valid XML")
	}
	// Should fail on XML parse, NOT on size limit.
	if strings.Contains(err.Error(), "exceeds maximum size") {
		t.Errorf("WithLimit should have allowed input within custom limit, got: %v", err)
	}
}

func TestReadXMLBIFWithLimitRejectsOverCustomLimit(t *testing.T) {
	input := strings.Repeat("x", 1000)
	_, err := ReadXMLBIFWithLimit(strings.NewReader(input), 500)
	if err == nil {
		t.Fatal("expected rejection for input exceeding custom limit")
	}
}

func TestReadJSONWithLimitAcceptsLargerInput(t *testing.T) {
	bigInput := strings.Repeat("{", MaxInputSize+100)
	_, err := ReadJSONWithLimit(strings.NewReader(bigInput), MaxInputSize+200)
	if err == nil {
		t.Skip("unexpectedly valid JSON")
	}
	if strings.Contains(err.Error(), "exceeds maximum size") {
		t.Errorf("WithLimit should have allowed input within custom limit, got: %v", err)
	}
}
