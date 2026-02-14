package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type OllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

type commitMessageJSON struct {
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "qwen2.5-coder:14b"
	}

	return &OllamaProvider{
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *OllamaProvider) Name() ProviderType {
	return ProviderOllama
}

func (p *OllamaProvider) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	prompt := p.buildPrompt(req)

	ollamaReq := ollamaRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if ollamaResp.Error != "" {
		return nil, fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	result := p.parseResponse(ollamaResp.Response, req.IncludeDescription)
	result.Prompt = prompt
	result.RawResponse = ollamaResp.Response

	return result, nil
}

func (p *OllamaProvider) buildPrompt(req GenerateRequest) string {
	var sb strings.Builder

	sb.WriteString("Generate a commit message based on the following git diff.\n\n")

	sb.WriteString("Current branch: ")
	sb.WriteString(req.BranchName)
	sb.WriteString("\n\n")

	sb.WriteString("Requirements:\n")
	sb.WriteString("- Be concise and descriptive\n")
	sb.WriteString("- Focus on WHAT changed and WHY, not HOW\n")
	sb.WriteString("- Use imperative mood (\"Add feature\" not \"Added feature\")\n")

	ticketPattern := regexp.MustCompile(`[A-Z]+-\d+`)
	ticketID := ticketPattern.FindString(req.BranchName)

	if req.UseConventionalCommits {
		sb.WriteString("- Use conventional commit format: <type>(<scope>): <description>\n")
		sb.WriteString("- Types: feat, fix, docs, style, refactor, test, chore\n")
		if ticketID != "" {
			sb.WriteString(fmt.Sprintf("- REQUIRED: Use ticket ID %s as the scope\n", ticketID))
			sb.WriteString(fmt.Sprintf("- Example: feat(%s): add new feature\n", ticketID))
		} else {
			sb.WriteString("- Example: feat(auth): add OAuth login support\n")
		}
	} else if ticketID != "" {
		sb.WriteString(fmt.Sprintf("- Include ticket ID %s\n", ticketID))
	}

	sb.WriteString("\nGit diff:\n")
	sb.WriteString(req.Diff)
	sb.WriteString("\n\n")

	sb.WriteString("Create a JSON object with these fields:\n")
	sb.WriteString("- message: the commit message\n")
	if req.IncludeDescription {
		sb.WriteString("- description: detailed explanation of the changes\n")
	}

	return sb.String()
}

func (p *OllamaProvider) parseResponse(response string, includeDescription bool) *GenerateResponse {
	result := &GenerateResponse{}

	var commitMsg commitMessageJSON
	if err := json.Unmarshal([]byte(response), &commitMsg); err != nil {
		result.Message = strings.TrimSpace(response)
		return result
	}

	result.Message = commitMsg.Message
	if includeDescription {
		result.Description = commitMsg.Description
	}

	return result
}
