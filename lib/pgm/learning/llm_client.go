package learning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LLMClient is a generic interface for large language model interactions,
// used by ExpertInLoop causal discovery to query an LLM for causal reasoning.
type LLMClient interface {
	// Complete sends a single prompt and returns the model's completion.
	Complete(prompt string, opts ...LLMOption) (string, error)

	// ChatComplete sends a sequence of chat messages and returns the model's reply.
	ChatComplete(messages []Message, opts ...LLMOption) (string, error)
}

// Message represents a single message in a chat conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMConfig holds configuration options for an LLM request.
type LLMConfig struct {
	Temperature float64
	MaxTokens   int
	Model       string
}

// LLMOption is a functional option for configuring LLM requests.
type LLMOption func(*LLMConfig)

// WithTemperature sets the sampling temperature for the LLM request.
func WithTemperature(t float64) LLMOption {
	return func(c *LLMConfig) {
		c.Temperature = t
	}
}

// WithMaxTokens sets the maximum number of tokens in the LLM response.
func WithMaxTokens(n int) LLMOption {
	return func(c *LLMConfig) {
		c.MaxTokens = n
	}
}

// WithModel sets the model to use for the LLM request.
func WithModel(m string) LLMOption {
	return func(c *LLMConfig) {
		c.Model = m
	}
}

// tokenBucket implements a simple token-bucket rate limiter.
type tokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

// newTokenBucket creates a token bucket that allows maxPerMinute requests per minute.
func newTokenBucket(maxPerMinute int) *tokenBucket {
	rate := float64(maxPerMinute) / 60.0
	return &tokenBucket{
		tokens:     float64(maxPerMinute),
		maxTokens:  float64(maxPerMinute),
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

// wait blocks until a token is available, then consumes one.
func (tb *tokenBucket) wait() {
	for {
		tb.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(tb.lastRefill).Seconds()
		tb.tokens += elapsed * tb.refillRate
		if tb.tokens > tb.maxTokens {
			tb.tokens = tb.maxTokens
		}
		tb.lastRefill = now

		if tb.tokens >= 1.0 {
			tb.tokens--
			tb.mu.Unlock()
			return
		}
		// Calculate how long until we have a token.
		waitDur := time.Duration((1.0 - tb.tokens) / tb.refillRate * float64(time.Second))
		tb.mu.Unlock()
		time.Sleep(waitDur)
	}
}

// HTTPLLMClient is an OpenAI-compatible HTTP client that implements LLMClient.
// It supports rate limiting and retry with exponential backoff.
type HTTPLLMClient struct {
	endpoint     string
	apiKey       string
	defaultModel string
	httpClient   *http.Client
	rateLimiter  *tokenBucket
	maxRetries   int
}

// NewHTTPLLMClient creates a new HTTP-based LLM client compatible with the
// OpenAI chat completions API. The endpoint should be a base URL (e.g.,
// "https://api.openai.com/v1"). Requests are rate-limited to 60 per minute.
func NewHTTPLLMClient(endpoint, apiKey string, defaultModel string) *HTTPLLMClient {
	return &HTTPLLMClient{
		endpoint:     strings.TrimRight(endpoint, "/"),
		apiKey:       apiKey,
		defaultModel: defaultModel,
		httpClient:   &http.Client{Timeout: 60 * time.Second},
		rateLimiter:  newTokenBucket(60),
		maxRetries:   3,
	}
}

// chatRequest is the JSON body sent to the chat completions endpoint.
type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// chatResponse is the JSON response from the chat completions endpoint.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Complete sends a single prompt to the LLM and returns the response text.
func (c *HTTPLLMClient) Complete(prompt string, opts ...LLMOption) (string, error) {
	messages := []Message{
		{Role: "user", Content: prompt},
	}
	return c.ChatComplete(messages, opts...)
}

// ChatComplete sends a sequence of chat messages and returns the model's reply.
func (c *HTTPLLMClient) ChatComplete(messages []Message, opts ...LLMOption) (string, error) {
	cfg := LLMConfig{
		Model: c.defaultModel,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	reqBody := chatRequest{
		Model:       cfg.Model,
		Messages:    messages,
		Temperature: cfg.Temperature,
		MaxTokens:   cfg.MaxTokens,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("llm: failed to marshal request: %w", err)
	}

	url := c.endpoint + "/chat/completions"

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			time.Sleep(backoff)
		}

		c.rateLimiter.wait()

		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			return "", fmt.Errorf("llm: failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if c.apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+c.apiKey)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("llm: request failed: %w", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("llm: failed to read response: %w", err)
			continue
		}

		// Retry on 429 (rate limited) or 5xx (server error).
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("llm: server returned status %d: %s", resp.StatusCode, string(respBody))
			continue
		}

		if resp.StatusCode != 200 {
			return "", fmt.Errorf("llm: server returned status %d: %s", resp.StatusCode, string(respBody))
		}

		var chatResp chatResponse
		if err := json.Unmarshal(respBody, &chatResp); err != nil {
			return "", fmt.Errorf("llm: failed to parse response: %w", err)
		}

		if chatResp.Error != nil {
			return "", fmt.Errorf("llm: API error: %s", chatResp.Error.Message)
		}

		if len(chatResp.Choices) == 0 {
			return "", fmt.Errorf("llm: no choices in response")
		}

		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("llm: all %d retries exhausted: %w", c.maxRetries, lastErr)
}

// isAlpha returns true if b is an ASCII letter.
func isAlpha(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

// CausalPromptTemplate provides templates for formatting causal reasoning
// queries to an LLM and parsing the responses.
type CausalPromptTemplate struct{}

// FormatCausalDirectionQuery formats a prompt asking the LLM whether var1
// causes var2 within a given context.
func (t CausalPromptTemplate) FormatCausalDirectionQuery(var1, var2, context string) string {
	return fmt.Sprintf(
		"Does %s cause %s? Context: %s. Answer YES, NO, or UNKNOWN.",
		var1, var2, context,
	)
}

// FormatIndependenceQuery formats a prompt asking the LLM whether var1 and
// var2 are conditionally independent given a set of conditioning variables.
func (t CausalPromptTemplate) FormatIndependenceQuery(var1, var2 string, given []string) string {
	if len(given) == 0 {
		return fmt.Sprintf(
			"Are %s and %s statistically independent? Answer YES, NO, or UNKNOWN with a confidence between 0 and 1.",
			var1, var2,
		)
	}
	return fmt.Sprintf(
		"Are %s and %s conditionally independent given {%s}? Answer YES, NO, or UNKNOWN with a confidence between 0 and 1.",
		var1, var2, strings.Join(given, ", "),
	)
}

// ParseCausalResponse parses an LLM response to extract a causal direction
// (YES, NO, or UNKNOWN) and an optional confidence value. The confidence
// defaults to 0.5 if not present in the response.
func (t CausalPromptTemplate) ParseCausalResponse(response string) (direction string, confidence float64, err error) {
	upper := strings.ToUpper(response)

	// Extract direction. Check UNKNOWN before NO because "UNKNOWN" contains "NO".
	// Use word-boundary-aware matching to avoid false positives (e.g., "not" matching "NO").
	containsWord := func(s, word string) bool {
		idx := 0
		for {
			pos := strings.Index(s[idx:], word)
			if pos < 0 {
				return false
			}
			pos += idx
			// Check left boundary.
			leftOK := pos == 0 || !isAlpha(s[pos-1])
			// Check right boundary.
			rightOK := pos+len(word) >= len(s) || !isAlpha(s[pos+len(word)])
			if leftOK && rightOK {
				return true
			}
			idx = pos + 1
		}
	}

	switch {
	case containsWord(upper, "UNKNOWN"):
		direction = "UNKNOWN"
	case containsWord(upper, "YES"):
		direction = "YES"
	case containsWord(upper, "NO"):
		direction = "NO"
	default:
		return "", 0, fmt.Errorf("llm: could not parse direction from response: %q", response)
	}

	// Try to extract a confidence value.
	confidence = 0.5 // default
	// Look for patterns like "confidence: 0.8" or "confidence 0.8" or "(0.8)".
	for _, prefix := range []string{"confidence:", "confidence", "conf:"} {
		idx := strings.Index(strings.ToLower(response), prefix)
		if idx >= 0 {
			after := strings.TrimSpace(response[idx+len(prefix):])
			var val float64
			if _, scanErr := fmt.Sscanf(after, "%f", &val); scanErr == nil && val >= 0 && val <= 1 {
				confidence = val
				return direction, confidence, nil
			}
		}
	}

	// Try to find a decimal in parentheses like (0.85).
	if idx := strings.Index(response, "("); idx >= 0 {
		if end := strings.Index(response[idx:], ")"); end >= 0 {
			inner := strings.TrimSpace(response[idx+1 : idx+end])
			var val float64
			if _, scanErr := fmt.Sscanf(inner, "%f", &val); scanErr == nil && val >= 0 && val <= 1 {
				confidence = val
			}
		}
	}

	return direction, confidence, nil
}
