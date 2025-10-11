package ai

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/jesseduffield/lazygit/pkg/config"
)

// ChatMessage and ChatRequest/Response mirror the OpenAI-compatible schema
// for the /chat/completions endpoint.
type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatRequest struct {
    Model       string        `json:"model"`
    Messages    []ChatMessage `json:"messages"`
    Temperature float64       `json:"temperature,omitempty"`
    MaxTokens   int           `json:"max_tokens,omitempty"`
}

type ChatChoice struct {
    Index        int         `json:"index"`
    FinishReason string      `json:"finish_reason"`
    Message      ChatMessage `json:"message"`
}

type ChatResponse struct {
    ID      string       `json:"id"`
    Object  string       `json:"object"`
    Created int64        `json:"created"`
    Choices []ChatChoice `json:"choices"`
}

type Client struct {
    http     *http.Client
    baseURL  string
    apiKey   string
    model    string
    temp     float64
    maxTokens int
}

func defaultBaseURL(provider string) string {
    switch strings.ToLower(provider) {
    case "openrouter":
        return "https://openrouter.ai/api/v1"
    case "openai", "":
        return "https://api.openai.com/v1"
    default:
        return provider // if somebody passed a URL by mistake into provider
    }
}

func defaultAPIKeyEnv(provider string) string {
    switch strings.ToLower(provider) {
    case "openrouter":
        return "OPENROUTER_API_KEY"
    default:
        return "OPENAI_API_KEY"
    }
}

func NewClientFromConfig(cfg config.AIConfig) (*Client, error) {
    // Validate configuration first
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("invalid AI configuration: %w", err)
    }

    apiKeyEnv := cfg.APIKeyEnv
    if apiKeyEnv == "" {
        apiKeyEnv = defaultAPIKeyEnv(cfg.Provider)
    }
    key := os.Getenv(apiKeyEnv)
    if key == "" {
        return nil, fmt.Errorf("missing API key; set %s", apiKeyEnv)
    }

    baseURL := cfg.BaseURL
    if baseURL == "" {
        baseURL = defaultBaseURL(cfg.Provider)
    }

    return &Client{
        http: &http.Client{Timeout: 45 * time.Second},
        baseURL: strings.TrimRight(baseURL, "/"),
        apiKey: key,
        model: cfg.Model,
        temp: cfg.Temperature,
        maxTokens: cfg.MaxTokens,
    }, nil
}

func (c *Client) chat(ctx context.Context, messages []ChatMessage) (string, error) {
    reqBody := ChatRequest{
        Model:       c.model,
        Messages:    messages,
        Temperature: c.temp,
        MaxTokens:   c.maxTokens,
    }
    b, err := json.Marshal(reqBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }
    url := c.baseURL + "/chat/completions"

    req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

    resp, err := c.http.Do(req)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("AI request failed (%s): %s", resp.Status, string(body))
    }

    var cr ChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }
    if len(cr.Choices) == 0 {
        return "", errors.New("no choices in AI response")
    }
    return cr.Choices[0].Message.Content, nil
}

// GenerateCommitMessage produces a subject and body from a diff and preferences.
func (c *Client) GenerateCommitMessage(ctx context.Context, diff string, style string, wrap int) (subject string, body string, err error) {
    sys := ChatMessage{Role: "system", Content: systemPrompt(style, wrap)}
    user := ChatMessage{Role: "user", Content: buildUserPrompt(diff)}
    out, err := c.chatWithRetry(ctx, []ChatMessage{sys, user})
    if err != nil {
        return "", "", err
    }

    // Normalize newlines and split into subject/body
    s := strings.TrimSpace(out)
    lines := strings.Split(s, "\n")
    if len(lines) == 0 {
        return "", "", errors.New("empty AI response")
    }
    subject = strings.TrimSpace(lines[0])
    body = strings.TrimSpace(strings.Join(lines[1:], "\n"))
    // Trim subject to 72 chars if it’s excessively long
    if len([]rune(subject)) > 72 {
        r := []rune(subject)
        subject = string(r[:72])
    }
    return subject, body, nil
}

// chatWithRetry implements retry logic with exponential backoff for transient failures
func (c *Client) chatWithRetry(ctx context.Context, messages []ChatMessage) (string, error) {
    var lastErr error
    maxRetries := 3
    baseDelay := time.Second

    for attempt := 0; attempt < maxRetries; attempt++ {
        result, err := c.chat(ctx, messages)
        if err == nil {
            return result, nil
        }

        lastErr = err

        // Don't retry on context cancellation or certain error types
        if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
            return "", err
        }

        // Don't retry on the last attempt
        if attempt == maxRetries-1 {
            break
        }

        // Calculate delay with exponential backoff
        delay := baseDelay * time.Duration(1<<uint(attempt))

        // Add some jitter to avoid thundering herd
        jitter := time.Duration(float64(delay) * 0.1 * (2*time.Now().UnixNano()%2 - 1))
        delay += jitter

        select {
        case <-ctx.Done():
            return "", ctx.Err()
        case <-time.After(delay):
            // Continue to next retry
        }
    }

    return "", fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

