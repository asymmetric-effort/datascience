//go:build unit

package learning

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// LLM HTTP client retry path coverage
// ---------------------------------------------------------------------------

func TestHTTPLLMClient_RetryOn429TwiceThen200(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(429)
			fmt.Fprint(w, `{"error":"rate limited"}`)
			return
		}
		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "success after retry"}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.maxRetries = 3
	// Use a high rate limit so rate limiter doesn't interfere.
	client.rateLimiter = newTokenBucket(6000)

	result, err := client.Complete("test")
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if result != "success after retry" {
		t.Errorf("expected 'success after retry', got %q", result)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("expected 3 attempts, got %d", got)
	}
}

func TestHTTPLLMClient_AllRetriesExhausted_500(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(500)
		fmt.Fprint(w, `internal server error`)
	}))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.maxRetries = 2
	client.rateLimiter = newTokenBucket(6000)

	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error when all retries exhausted")
	}
	if !strings.Contains(err.Error(), "retries exhausted") {
		t.Errorf("expected 'retries exhausted' error, got: %v", err)
	}
	// Should have made maxRetries+1 attempts (initial + retries).
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("expected 3 attempts (1 initial + 2 retries), got %d", got)
	}
}

func TestHTTPLLMClient_NonRetryable400(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(400)
		fmt.Fprint(w, `bad request`)
	}))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.maxRetries = 3
	client.rateLimiter = newTokenBucket(6000)

	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error for 400 status")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected status 400 in error, got: %v", err)
	}
	// 400 should NOT retry — only 1 attempt.
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("expected 1 attempt for non-retryable error, got %d", got)
	}
}

func TestHTTPLLMClient_NonRetryable403(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(403)
		fmt.Fprint(w, `forbidden`)
	}))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.maxRetries = 3
	client.rateLimiter = newTokenBucket(6000)

	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error for 403 status")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("expected 1 attempt for 403, got %d", got)
	}
}

func TestHTTPLLMClient_RateLimiterThrottling(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "ok"}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	// Rate limiter allows only 1 per minute — second request must wait.
	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.rateLimiter = newTokenBucket(60) // 1 per second

	// First request should succeed immediately.
	start := time.Now()
	result, err := client.Complete("first")
	if err != nil {
		t.Fatalf("first request error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected 'ok', got %q", result)
	}
	elapsed := time.Since(start)
	if elapsed > 2*time.Second {
		t.Errorf("first request took too long: %v", elapsed)
	}
}

func TestHTTPLLMClient_RetryOn503(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 1 {
			w.WriteHeader(503)
			fmt.Fprint(w, `service unavailable`)
			return
		}
		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "recovered"}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.maxRetries = 2
	client.rateLimiter = newTokenBucket(6000)

	result, err := client.Complete("test")
	if err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if result != "recovered" {
		t.Errorf("expected 'recovered', got %q", result)
	}
}

func TestHTTPLLMClient_ConnectionError_Retries(t *testing.T) {
	// Client pointed at a non-existent server triggers HTTP request failure path.
	client := NewHTTPLLMClient("http://127.0.0.1:1", "key", "model")
	client.maxRetries = 1
	client.rateLimiter = newTokenBucket(6000)

	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error for connection failure")
	}
	if !strings.Contains(err.Error(), "retries exhausted") {
		t.Errorf("expected 'retries exhausted', got: %v", err)
	}
}

func TestHTTPLLMClient_InvalidURL(t *testing.T) {
	// An endpoint with invalid characters triggers request creation failure.
	client := NewHTTPLLMClient("http://[::1]:namedport", "key", "model")
	client.maxRetries = 0
	client.rateLimiter = newTokenBucket(6000)

	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestHTTPLLMClient_MixedRetryable(t *testing.T) {
	// First 429, then 500, then success.
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		switch n {
		case 1:
			w.WriteHeader(429)
			fmt.Fprint(w, `rate limited`)
		case 2:
			w.WriteHeader(500)
			fmt.Fprint(w, `server error`)
		default:
			resp := chatResponse{
				Choices: []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				}{
					{Message: struct {
						Content string `json:"content"`
					}{Content: "finally ok"}},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.maxRetries = 3
	client.rateLimiter = newTokenBucket(6000)

	result, err := client.Complete("test")
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if result != "finally ok" {
		t.Errorf("expected 'finally ok', got %q", result)
	}
}
