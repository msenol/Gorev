// Package providers implements AI provider integrations
// This file contains a shared HTTP client for OpenAI-compatible APIs (DRY)
package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// OpenAICompatibleClient is a shared HTTP client for OpenAI-compatible APIs
// Used by both OpenRouter and Anannas providers (DRY principle)
type OpenAICompatibleClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	userAgent  string
}

// NewOpenAICompatibleClient creates a new OpenAI-compatible HTTP client
func NewOpenAICompatibleClient(baseURL, apiKey string) *OpenAICompatibleClient {
	return &OpenAICompatibleClient{
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: "Gorev/0.18.0",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetTimeout sets the HTTP client timeout
func (c *OpenAICompatibleClient) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// Do executes an HTTP request with proper headers and error handling
func (c *OpenAICompatibleClient) Do(ctx context.Context, method, endpoint string, body, result interface{}) error {
	// Build full URL
	url := c.baseURL + endpoint

	// Marshal request body if provided
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)

	// Add HTTP-Referer header for OpenRouter (required)
	if reqBody != nil {
		req.Header.Set("HTTP-Referer", "https://github.com/msenol/Gorev")
		req.Header.Set("X-Title", "Gorev Task Manager")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	// Check for non-OK status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Unmarshal response into result if provided
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w (body: %s)", err, string(respBody))
		}
	}

	return nil
}

// Get executes a GET request
func (c *OpenAICompatibleClient) Get(ctx context.Context, endpoint string, result interface{}) error {
	return c.Do(ctx, http.MethodGet, endpoint, nil, result)
}

// Post executes a POST request
func (c *OpenAICompatibleClient) Post(ctx context.Context, endpoint string, body, result interface{}) error {
	return c.Do(ctx, http.MethodPost, endpoint, body, result)
}

// ModelsResponse represents the response from listing models
type ModelsResponse struct {
	Data []ModelInfo `json:"data"`
}

// ModelInfo represents information about an AI model
type ModelInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	ContextWindow int     `json:"context_window"`
	Pricing       Pricing `json:"pricing"`
}

// Pricing represents model pricing information
type Pricing struct {
	Input  PriceValue `json:"prompt"`
	Output PriceValue `json:"completion"`
}

// PriceValue is a flexible type that can be unmarshaled from both strings and floats
type PriceValue float64

// UnmarshalJSON implements custom JSON unmarshaling for PriceValue
func (p *PriceValue) UnmarshalJSON(data []byte) error {
	// Try unmarshaling as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*p = PriceValue(f)
		return nil
	}
	// Try unmarshaling as float
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*p = PriceValue(f)
	return nil
}

// ChatRequest represents an OpenAI-compatible chat completion request
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// ChatResponse represents an OpenAI-compatible chat completion response
type ChatResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []ChatChoice     `json:"choices"`
	Usage   ChatUsage        `json:"usage"`
	Error   *ChatErrorDetail `json:"error,omitempty"`
}

// ChatChoice represents a choice in the chat response
type ChatChoice struct {
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	FinishReason string       `json:"finish_reason"`
}

// ChatUsage represents token usage information
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatErrorDetail represents error details in a chat response
type ChatErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// ListModels retrieves available models from the provider
func (c *OpenAICompatibleClient) ListModels(ctx context.Context) ([]ModelInfo, error) {
	var resp ModelsResponse
	if err := c.Get(ctx, "/models", &resp); err != nil {
		return nil, fmt.Errorf("list models: %w", err)
	}
	return resp.Data, nil
}

// ChatCompletion sends a chat completion request
func (c *OpenAICompatibleClient) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	var resp ChatResponse
	if err := c.Post(ctx, "/chat/completions", req, &resp); err != nil {
		return nil, fmt.Errorf("chat completion: %w", err)
	}

	// Check for API-level errors in response
	if resp.Error != nil {
		return nil, fmt.Errorf("API error: %s (type: %s, code: %s)", resp.Error.Message, resp.Error.Type, resp.Error.Code)
	}

	return &resp, nil
}

// BaseURL returns the base URL of the client
func (c *OpenAICompatibleClient) BaseURL() string {
	return c.baseURL
}
