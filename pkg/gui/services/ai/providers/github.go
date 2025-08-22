package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/config"
)

const (
	copilotAuthDeviceCodeURL  = "https://github.com/login/device/code"
	copilotAuthTokenURL       = "https://github.com/login/oauth/access_token"
	copilotChatAuthURL        = "https://api.github.com/copilot_internal/v2/token"
	copilotChatCompletionsURL = "https://api.githubcopilot.com/chat/completions"
	copilotEditorVersion      = "vscode/1.95.3"
	copilotUserAgent          = "curl/7.81.0"
	copilotClientID           = "Iv1.b507a08c87ecfe98"
)

// GitHubProvider implements the Provider interface for GitHub Copilot API
type GitHubProvider struct {
	config      config.AIConfig
	httpClient  *http.Client
	accessToken *AccessToken
}

// AccessToken response from GitHub Copilot's token endpoint
type AccessToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	Endpoints struct {
		API           string `json:"api"`
		OriginTracker string `json:"origin-tracker"`
		Proxy         string `json:"proxy"`
		Telemetry     string `json:"telemetry"`
	} `json:"endpoints"`
	ErrorDetails *struct {
		URL            string `json:"url,omitempty"`
		Message        string `json:"message,omitempty"`
		Title          string `json:"title,omitempty"`
		NotificationID string `json:"notification_id,omitempty"`
	} `json:"error_details,omitempty"`
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type DeviceTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error,omitempty"`
}

type FailedRequestResponse struct {
	DocumentationURL string `json:"documentation_url"`
	Message          string `json:"message"`
}

type OAuthTokenWrapper struct {
	User        string `json:"user"`
	OAuthToken  string `json:"oauth_token"`
	GithubAppID string `json:"githubAppId"`
}

type OAuthToken struct {
	GithubWrapper OAuthTokenWrapper `json:"github.com:Iv1.b507a08c87ecfe98"`
}

// GitHub Copilot API request/response structures
type copilotRequest struct {
	Messages    []copilotMessage `json:"messages"`
	Model       string           `json:"model"`
	Temperature float64          `json:"temperature"`
	TopP        float64          `json:"top_p"`
	N           int              `json:"n"`
	Stream      bool             `json:"stream"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
}

type copilotMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type copilotResponse struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Choices []copilotChoice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *copilotError `json:"error,omitempty"`
}

type copilotChoice struct {
	Index        int            `json:"index"`
	Message      copilotMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

type copilotError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// NewGitHubProvider creates a new GitHub Copilot provider
func NewGitHubProvider(config config.AIConfig) *GitHubProvider {
	return &GitHubProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateMessage generates a commit message using GitHub Copilot API
func (p *GitHubProvider) GenerateMessage(prompt string) (string, error) {
	// Authenticate and get access token
	if err := p.authenticate(); err != nil {
		return "", fmt.Errorf("failed to authenticate with GitHub Copilot: %w", err)
	}

	// Build request
	request := p.buildRequest(prompt)

	// Make API call
	ctx := context.Background()
	response, err := p.makeAPICall(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to call GitHub Copilot API: %w", err)
	}

	// Extract message
	message, err := p.extractMessage(response)
	if err != nil {
		return "", fmt.Errorf("failed to extract message from response: %w", err)
	}

	return message, nil
}

// Name returns the provider name
func (p *GitHubProvider) Name() string {
	return "github"
}

// ValidateConfig validates the GitHub Copilot configuration
func (p *GitHubProvider) ValidateConfig() error {
	// GitHub Copilot uses OAuth authentication, so we don't require an API key
	// We'll validate during authentication flow instead
	return nil
}

// authenticate handles GitHub Copilot authentication
func (p *GitHubProvider) authenticate() error {
	// Check if we have a valid access token
	if p.accessToken != nil && p.accessToken.ExpiresAt > time.Now().Unix() {
		return nil
	}

	// Get OAuth token from config files or login flow
	oauthToken, err := p.getOAuthToken()
	if err != nil {
		return fmt.Errorf("failed to get OAuth token: %w", err)
	}

	// Get access token for Copilot API
	accessToken, err := p.getCopilotAccessToken(oauthToken)
	if err != nil {
		return fmt.Errorf("failed to get Copilot access token: %w", err)
	}

	p.accessToken = &accessToken
	return nil
}

// getOAuthToken gets the OAuth token from config files or initiates login flow
func (p *GitHubProvider) getOAuthToken() (string, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".config/github-copilot")
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("LOCALAPPDATA"), "github-copilot")
	}

	// Support both legacy and current config file locations
	legacyConfigPath := filepath.Join(configPath, "hosts.json")
	currentConfigPath := filepath.Join(configPath, "apps.json")

	// Try to get token from config files
	configFiles := []string{legacyConfigPath, currentConfigPath}
	for _, path := range configFiles {
		token, err := p.extractTokenFromFile(path)
		if err == nil && token != "" {
			return token, nil
		}
	}

	// If no token found, initiate login flow
	return p.loginFlow(currentConfigPath)
}

// extractTokenFromFile extracts OAuth token from config file
func (p *GitHubProvider) extractTokenFromFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read Copilot configuration file at %s: %w", path, err)
	}

	var config map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &config); err != nil {
		return "", fmt.Errorf("failed to parse Copilot configuration file at %s: %w", path, err)
	}

	for key, value := range config {
		if key == "github.com" || strings.HasPrefix(key, "github.com:") {
			var tokenData map[string]string
			if err := json.Unmarshal(value, &tokenData); err != nil {
				continue
			}
			if token, exists := tokenData["oauth_token"]; exists {
				return token, nil
			}
		}
	}

	return "", fmt.Errorf("no token found in %s", path)
}

// loginFlow initiates the GitHub OAuth device flow
func (p *GitHubProvider) loginFlow(configPath string) (string, error) {
	data := strings.NewReader(fmt.Sprintf("client_id=%s&scope=copilot", copilotClientID))
	req, err := http.NewRequest("POST", copilotAuthDeviceCodeURL, data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get device code: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to decode device code response: %w", err)
	}

	parsedData, err := url.ParseQuery(string(responseBody))
	if err != nil {
		return "", fmt.Errorf("failed to parse device code response: %w", err)
	}

	deviceCodeResp := DeviceCodeResponse{
		UserCode:        parsedData.Get("user_code"),
		DeviceCode:      parsedData.Get("device_code"),
		VerificationURI: parsedData.Get("verification_uri"),
	}

	deviceCodeResp.ExpiresIn, _ = strconv.Atoi(parsedData.Get("expires_in"))
	deviceCodeResp.Interval, _ = strconv.Atoi(parsedData.Get("interval"))

	fmt.Printf("Please go to %s and enter the code %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)

	oAuthToken, err := p.fetchRefreshToken(deviceCodeResp.DeviceCode, deviceCodeResp.Interval, deviceCodeResp.ExpiresIn)
	if err != nil {
		return "", err
	}

	err = p.saveOAuthToken(OAuthToken{
		GithubWrapper: OAuthTokenWrapper{
			User:        "",
			OAuthToken:  oAuthToken.AccessToken,
			GithubAppID: copilotClientID,
		},
	}, configPath)
	if err != nil {
		return "", err
	}

	return oAuthToken.AccessToken, nil
}

// fetchRefreshToken polls for the OAuth token
func (p *GitHubProvider) fetchRefreshToken(deviceCode string, interval int, expiresIn int) (DeviceTokenResponse, error) {
	var accessTokenResp DeviceTokenResponse
	var errResp FailedRequestResponse

	time.Sleep(30 * time.Second) // Give user time to open browser

	endTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if time.Now().After(endTime) {
			return DeviceTokenResponse{}, fmt.Errorf("authorization polling timeout")
		}

		fmt.Println("Trying to fetch token...")
		data := strings.NewReader(fmt.Sprintf(
			"client_id=%s&device_code=%s&grant_type=urn:ietf:params:oauth:grant-type:device_code",
			copilotClientID, deviceCode,
		))

		req, err := http.NewRequest("POST", copilotAuthTokenURL, data)
		if err != nil {
			return DeviceTokenResponse{}, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return DeviceTokenResponse{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
				return DeviceTokenResponse{}, err
			}
			return DeviceTokenResponse{}, fmt.Errorf("failed to check refresh token: %s", errResp.Message)
		}

		if err := json.NewDecoder(resp.Body).Decode(&accessTokenResp); err != nil {
			return DeviceTokenResponse{}, err
		}

		if accessTokenResp.AccessToken != "" {
			return accessTokenResp, nil
		}

		if accessTokenResp.Error != "" && accessTokenResp.Error != "authorization_pending" {
			return DeviceTokenResponse{}, fmt.Errorf("token error: %s", accessTokenResp.Error)
		}
	}

	return DeviceTokenResponse{}, fmt.Errorf("authorization polling failed or timed out")
}

// saveOAuthToken saves the OAuth token to the config file
func (p *GitHubProvider) saveOAuthToken(oAuthToken OAuthToken, configPath string) error {
	fileContent, err := json.Marshal(oAuthToken)
	if err != nil {
		return fmt.Errorf("error marshaling oAuthToken: %w", err)
	}

	configDir := filepath.Dir(configPath)
	if err = os.MkdirAll(configDir, 0o700); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	err = os.WriteFile(configPath, fileContent, 0o700)
	if err != nil {
		return fmt.Errorf("error writing oAuthToken to %s: %w", configPath, err)
	}

	return nil
}

// getCopilotAccessToken exchanges OAuth token for Copilot access token
func (p *GitHubProvider) getCopilotAccessToken(oauthToken string) (AccessToken, error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, copilotChatAuthURL, nil)
	if err != nil {
		return AccessToken{}, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Authorization", "token "+oauthToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Editor-Version", copilotEditorVersion)
	req.Header.Set("User-Agent", copilotUserAgent)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return AccessToken{}, fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResponse AccessToken
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return AccessToken{}, fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResponse.ErrorDetails != nil {
		return AccessToken{}, fmt.Errorf("token error: %s", tokenResponse.ErrorDetails.Message)
	}

	return tokenResponse, nil
}

// buildRequest creates the GitHub Copilot API request payload
func (p *GitHubProvider) buildRequest(prompt string) *copilotRequest {
	model := p.config.Model
	if model == "" {
		model = "gpt-4"
	}

	temperature := p.config.Temperature
	if temperature == 0 {
		temperature = 0.7
	}

	maxTokens := p.config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 500
	}

	return &copilotRequest{
		Messages: []copilotMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that generates concise, descriptive git commit messages. Focus on what was changed and why, following conventional commit format when appropriate.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Generate a git commit message for the following changes:\n\n%s", prompt),
			},
		},
		Model:       model,
		Temperature: temperature,
		TopP:        1.0,
		N:           1,
		Stream:      false,
		MaxTokens:   maxTokens,
	}
}

// makeAPICall makes the HTTP request to GitHub Copilot API
func (p *GitHubProvider) makeAPICall(ctx context.Context, request *copilotRequest) (*copilotResponse, error) {
	// Serialize request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := copilotChatCompletionsURL
	if p.config.BaseURL != "" {
		url = p.config.BaseURL
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.accessToken.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Editor-Version", copilotEditorVersion)
	req.Header.Set("User-Agent", copilotUserAgent)

	// Make request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse response
	var apiResponse copilotResponse
	if err := json.Unmarshal(responseBody, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResponse.Error != nil {
		return nil, fmt.Errorf("GitHub Copilot API error: %s", apiResponse.Error.Message)
	}

	return &apiResponse, nil
}

// extractMessage extracts the generated message from GitHub Copilot response
func (p *GitHubProvider) extractMessage(response *copilotResponse) (string, error) {
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	choice := response.Choices[0]
	message := choice.Message.Content

	if message == "" {
		return "", fmt.Errorf("empty message in response")
	}

	// Clean up the message - remove quotes and extra whitespace
	message = strings.Trim(message, "\"'\n\r\t ")

	return message, nil
}
