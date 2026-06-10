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
)

// ---------------------------------------------------------------------------
// LLMConfig and options
// ---------------------------------------------------------------------------

func TestLLMConfig_Defaults(t *testing.T) {
	cfg := LLMConfig{}
	if cfg.Temperature != 0 {
		t.Errorf("expected default temperature 0, got %f", cfg.Temperature)
	}
	if cfg.MaxTokens != 0 {
		t.Errorf("expected default max tokens 0, got %d", cfg.MaxTokens)
	}
	if cfg.Model != "" {
		t.Errorf("expected default model empty, got %q", cfg.Model)
	}
}

func TestWithTemperature(t *testing.T) {
	cfg := LLMConfig{}
	WithTemperature(0.7)(&cfg)
	if cfg.Temperature != 0.7 {
		t.Errorf("expected temperature 0.7, got %f", cfg.Temperature)
	}
}

func TestWithMaxTokens(t *testing.T) {
	cfg := LLMConfig{}
	WithMaxTokens(512)(&cfg)
	if cfg.MaxTokens != 512 {
		t.Errorf("expected max tokens 512, got %d", cfg.MaxTokens)
	}
}

func TestWithModel(t *testing.T) {
	cfg := LLMConfig{}
	WithModel("gpt-4")(&cfg)
	if cfg.Model != "gpt-4" {
		t.Errorf("expected model gpt-4, got %q", cfg.Model)
	}
}

func TestLLMOptions_Chained(t *testing.T) {
	cfg := LLMConfig{}
	opts := []LLMOption{
		WithTemperature(0.3),
		WithMaxTokens(100),
		WithModel("test-model"),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Temperature != 0.3 {
		t.Errorf("expected temperature 0.3, got %f", cfg.Temperature)
	}
	if cfg.MaxTokens != 100 {
		t.Errorf("expected max tokens 100, got %d", cfg.MaxTokens)
	}
	if cfg.Model != "test-model" {
		t.Errorf("expected model test-model, got %q", cfg.Model)
	}
}

// ---------------------------------------------------------------------------
// Message struct
// ---------------------------------------------------------------------------

func TestMessage_Fields(t *testing.T) {
	msg := Message{Role: "user", Content: "hello"}
	if msg.Role != "user" {
		t.Errorf("expected role user, got %q", msg.Role)
	}
	if msg.Content != "hello" {
		t.Errorf("expected content hello, got %q", msg.Content)
	}
}

func TestMessage_JSON(t *testing.T) {
	msg := Message{Role: "assistant", Content: "world"}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if decoded.Role != "assistant" || decoded.Content != "world" {
		t.Errorf("round-trip failed: got %+v", decoded)
	}
}

// ---------------------------------------------------------------------------
// LLMClient interface compliance
// ---------------------------------------------------------------------------

func TestHTTPLLMClient_ImplementsInterface(t *testing.T) {
	// Compile-time check that HTTPLLMClient satisfies LLMClient.
	var _ LLMClient = (*HTTPLLMClient)(nil)
}

// ---------------------------------------------------------------------------
// HTTPLLMClient with mock server
// ---------------------------------------------------------------------------

func newMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func successHandler(content string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: content}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func TestHTTPLLMClient_Complete(t *testing.T) {
	srv := newMockServer(successHandler("test response"))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "test-key", "test-model")
	result, err := client.Complete("hello")
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if result != "test response" {
		t.Errorf("expected 'test response', got %q", result)
	}
}

func TestHTTPLLMClient_ChatComplete(t *testing.T) {
	srv := newMockServer(successHandler("chat reply"))
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "test-key", "test-model")
	messages := []Message{
		{Role: "system", Content: "You are helpful."},
		{Role: "user", Content: "Hi"},
	}
	result, err := client.ChatComplete(messages)
	if err != nil {
		t.Fatalf("ChatComplete() error: %v", err)
	}
	if result != "chat reply" {
		t.Errorf("expected 'chat reply', got %q", result)
	}
}

func TestHTTPLLMClient_SendsCorrectRequest(t *testing.T) {
	var receivedBody chatRequest
	var receivedAuth string
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected path /chat/completions, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		json.NewDecoder(r.Body).Decode(&receivedBody)
		successHandler("ok")(w, r)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "my-api-key", "default-model")
	_, err := client.Complete("test prompt", WithTemperature(0.5), WithMaxTokens(200), WithModel("custom-model"))
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}

	if receivedAuth != "Bearer my-api-key" {
		t.Errorf("expected auth 'Bearer my-api-key', got %q", receivedAuth)
	}
	if receivedBody.Model != "custom-model" {
		t.Errorf("expected model custom-model, got %q", receivedBody.Model)
	}
	if receivedBody.Temperature != 0.5 {
		t.Errorf("expected temperature 0.5, got %f", receivedBody.Temperature)
	}
	if receivedBody.MaxTokens != 200 {
		t.Errorf("expected max tokens 200, got %d", receivedBody.MaxTokens)
	}
	if len(receivedBody.Messages) != 1 || receivedBody.Messages[0].Content != "test prompt" {
		t.Errorf("unexpected messages: %+v", receivedBody.Messages)
	}
}

func TestHTTPLLMClient_UsesDefaultModel(t *testing.T) {
	var receivedBody chatRequest
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		successHandler("ok")(w, r)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "my-default")
	_, err := client.Complete("prompt")
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if receivedBody.Model != "my-default" {
		t.Errorf("expected default model my-default, got %q", receivedBody.Model)
	}
}

func TestHTTPLLMClient_NoAPIKey(t *testing.T) {
	var receivedAuth string
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		successHandler("ok")(w, r)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "", "model")
	_, err := client.Complete("test")
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if receivedAuth != "" {
		t.Errorf("expected no auth header, got %q", receivedAuth)
	}
}

func TestHTTPLLMClient_EmptyChoices(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"choices":[]}`)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
	if !strings.Contains(err.Error(), "no choices") {
		t.Errorf("expected 'no choices' error, got: %v", err)
	}
}

func TestHTTPLLMClient_APIError(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"error":{"message":"invalid model"}}`)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error for API error response")
	}
	if !strings.Contains(err.Error(), "invalid model") {
		t.Errorf("expected 'invalid model' error, got: %v", err)
	}
}

func TestHTTPLLMClient_NonRetryableHTTPError(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprint(w, `bad request`)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error for 400 status")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected status 400 in error, got: %v", err)
	}
}

func TestHTTPLLMClient_RetryOn429(t *testing.T) {
	var attempts int32
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 2 {
			w.WriteHeader(429)
			fmt.Fprint(w, `rate limited`)
			return
		}
		successHandler("recovered")(w, r)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	// Override max retries to speed up test (uses short backoff).
	client.maxRetries = 3

	result, err := client.Complete("test")
	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if result != "recovered" {
		t.Errorf("expected 'recovered', got %q", result)
	}
	if atomic.LoadInt32(&attempts) < 3 {
		t.Errorf("expected at least 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

func TestHTTPLLMClient_RetryOn500(t *testing.T) {
	var attempts int32
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n <= 1 {
			w.WriteHeader(500)
			fmt.Fprint(w, `internal error`)
			return
		}
		successHandler("ok")(w, r)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	result, err := client.Complete("test")
	if err != nil {
		t.Fatalf("expected success after retry, got error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected 'ok', got %q", result)
	}
}

func TestHTTPLLMClient_AllRetriesExhausted(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprint(w, `server error`)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	client.maxRetries = 1 // Reduce for speed.

	_, err := client.Complete("test")
	if err == nil {
		t.Fatal("expected error when all retries exhausted")
	}
	if !strings.Contains(err.Error(), "retries exhausted") {
		t.Errorf("expected 'retries exhausted' error, got: %v", err)
	}
}

func TestHTTPLLMClient_EndpointTrailingSlash(t *testing.T) {
	var receivedPath string
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		successHandler("ok")(w, r)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL+"/", "key", "model")
	_, err := client.Complete("test")
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}
	if receivedPath != "/chat/completions" {
		t.Errorf("expected /chat/completions, got %s", receivedPath)
	}
}

// ---------------------------------------------------------------------------
// CausalPromptTemplate
// ---------------------------------------------------------------------------

func TestCausalPromptTemplate_FormatCausalDirectionQuery(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	query := tmpl.FormatCausalDirectionQuery("Smoking", "Cancer", "epidemiological study")

	if !strings.Contains(query, "Smoking") {
		t.Error("expected query to contain var1")
	}
	if !strings.Contains(query, "Cancer") {
		t.Error("expected query to contain var2")
	}
	if !strings.Contains(query, "epidemiological study") {
		t.Error("expected query to contain context")
	}
	if !strings.Contains(query, "YES, NO, or UNKNOWN") {
		t.Error("expected query to contain answer format")
	}
}

func TestCausalPromptTemplate_FormatIndependenceQuery_NoGiven(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	query := tmpl.FormatIndependenceQuery("X", "Y", nil)

	if !strings.Contains(query, "X") || !strings.Contains(query, "Y") {
		t.Error("expected query to contain variables")
	}
	if !strings.Contains(query, "independent") {
		t.Error("expected query to mention independence")
	}
	if strings.Contains(query, "given") {
		t.Error("expected no conditioning set in query")
	}
}

func TestCausalPromptTemplate_FormatIndependenceQuery_EmptyGiven(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	query := tmpl.FormatIndependenceQuery("X", "Y", []string{})

	if strings.Contains(query, "given") {
		t.Error("expected no conditioning set for empty given")
	}
}

func TestCausalPromptTemplate_FormatIndependenceQuery_WithGiven(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	query := tmpl.FormatIndependenceQuery("X", "Y", []string{"Z", "W"})

	if !strings.Contains(query, "conditionally independent") {
		t.Error("expected query to mention conditional independence")
	}
	if !strings.Contains(query, "Z, W") {
		t.Error("expected query to contain conditioning variables")
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_YES(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	dir, conf, err := tmpl.ParseCausalResponse("YES, I believe so. Confidence: 0.9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != "YES" {
		t.Errorf("expected YES, got %q", dir)
	}
	if conf != 0.9 {
		t.Errorf("expected confidence 0.9, got %f", conf)
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_NO(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	dir, conf, err := tmpl.ParseCausalResponse("NO. confidence: 0.8")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != "NO" {
		t.Errorf("expected NO, got %q", dir)
	}
	if conf != 0.8 {
		t.Errorf("expected confidence 0.8, got %f", conf)
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_UNKNOWN(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	dir, conf, err := tmpl.ParseCausalResponse("UNKNOWN - not enough data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != "UNKNOWN" {
		t.Errorf("expected UNKNOWN, got %q", dir)
	}
	if conf != 0.5 {
		t.Errorf("expected default confidence 0.5, got %f", conf)
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_Lowercase(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	dir, _, err := tmpl.ParseCausalResponse("yes, that is correct")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != "YES" {
		t.Errorf("expected YES, got %q", dir)
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_WithParenConfidence(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	dir, conf, err := tmpl.ParseCausalResponse("YES (0.85)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != "YES" {
		t.Errorf("expected YES, got %q", dir)
	}
	if conf != 0.85 {
		t.Errorf("expected confidence 0.85, got %f", conf)
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_NoDirection(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	_, _, err := tmpl.ParseCausalResponse("I'm not sure what you mean")
	if err == nil {
		t.Fatal("expected error for unparseable response")
	}
	if !strings.Contains(err.Error(), "could not parse direction") {
		t.Errorf("expected 'could not parse direction' error, got: %v", err)
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_DefaultConfidence(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	_, conf, err := tmpl.ParseCausalResponse("YES")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conf != 0.5 {
		t.Errorf("expected default confidence 0.5, got %f", conf)
	}
}

func TestCausalPromptTemplate_ParseCausalResponse_ConfPrefix(t *testing.T) {
	tmpl := CausalPromptTemplate{}
	_, conf, err := tmpl.ParseCausalResponse("NO. conf: 0.75")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conf != 0.75 {
		t.Errorf("expected confidence 0.75, got %f", conf)
	}
}

// ---------------------------------------------------------------------------
// Integration: HTTPLLMClient + CausalPromptTemplate with mock server
// ---------------------------------------------------------------------------

func TestHTTPLLMClient_CausalWorkflow(t *testing.T) {
	srv := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		json.NewDecoder(r.Body).Decode(&req)
		// Echo back a causal response.
		content := req.Messages[len(req.Messages)-1].Content
		var reply string
		if strings.Contains(content, "cause") {
			reply = "YES, based on the evidence. Confidence: 0.92"
		} else {
			reply = "NO. Confidence: 0.7"
		}
		resp := chatResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: reply}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	client := NewHTTPLLMClient(srv.URL, "key", "model")
	tmpl := CausalPromptTemplate{}

	// Test causal direction query.
	prompt := tmpl.FormatCausalDirectionQuery("Smoking", "Cancer", "medical research")
	response, err := client.Complete(prompt, WithTemperature(0.1))
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}

	dir, conf, err := tmpl.ParseCausalResponse(response)
	if err != nil {
		t.Fatalf("ParseCausalResponse() error: %v", err)
	}
	if dir != "YES" {
		t.Errorf("expected YES, got %q", dir)
	}
	if conf != 0.92 {
		t.Errorf("expected confidence 0.92, got %f", conf)
	}

	// Test independence query.
	prompt = tmpl.FormatIndependenceQuery("X", "Y", []string{"Z"})
	response, err = client.Complete(prompt)
	if err != nil {
		t.Fatalf("Complete() error: %v", err)
	}

	dir, conf, err = tmpl.ParseCausalResponse(response)
	if err != nil {
		t.Fatalf("ParseCausalResponse() error: %v", err)
	}
	if dir != "NO" {
		t.Errorf("expected NO, got %q", dir)
	}
	if conf != 0.7 {
		t.Errorf("expected confidence 0.7, got %f", conf)
	}
}
