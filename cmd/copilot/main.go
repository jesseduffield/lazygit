package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sanity-io/litter"
	"github.com/zalando/go-keyring"
)

const (
	KEYRING_SERVICE = "lazygit"
	KEYRING_USER    = "github-copilot"
)

// ICopilotChat defines the interface for chat operations
type ICopilotChat interface {
	Authenticate() error
	IsAuthenticated() bool
	Chat(request Request) (string, error)
}

var _ ICopilotChat = &CopilotChat{}

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type Model string

const (
	Gpt4o           Model = "gpt-4o-2024-05-13"
	Gpt4            Model = "gpt-4"
	Gpt3_5Turbo     Model = "gpt-3.5-turbo"
	O1Preview       Model = "o1-preview-2024-09-12"
	O1Mini          Model = "o1-mini-2024-09-12"
	Claude3_5Sonnet Model = "claude-3.5-sonnet"
)

const (
	COPILOT_CHAT_COMPLETION_URL = "https://api.githubcopilot.com/chat/completions"
	COPILOT_CHAT_AUTH_URL       = "https://api.github.com/copilot_internal/v2/token"
	EDITOR_VERSION              = "Lazygit/0.44.0"
	COPILOT_INTEGRATION_ID      = "vscode-chat"
)
const (
	CACHE_FILE_NAME = ".copilot_auth.json"
)
const (
	CHECK_INTERVAL = 30 * time.Second
	MAX_AUTH_TIME  = 5 * time.Minute
)
const (
	GITHUB_CLIENT_ID = "Iv1.b507a08c87ecfe98"
)

type ChatMessage struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Intent      bool          `json:"intent"`
	N           int           `json:"n"`
	Stream      bool          `json:"stream"`
	Temperature float32       `json:"temperature"`
	Model       Model         `json:"model"`
	Messages    []ChatMessage `json:"messages"`
}

type ContentFilterResult struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity"`
}

type ContentFilterResults struct {
	Hate     ContentFilterResult `json:"hate"`
	SelfHarm ContentFilterResult `json:"self_harm"`
	Sexual   ContentFilterResult `json:"sexual"`
	Violence ContentFilterResult `json:"violence"`
}

type ChatResponse struct {
	Choices             []ResponseChoice     `json:"choices"`
	Created             int64                `json:"created"`
	ID                  string               `json:"id"`
	Model               string               `json:"model"`
	SystemFingerprint   string               `json:"system_fingerprint"`
	PromptFilterResults []PromptFilterResult `json:"prompt_filter_results"`
	Usage               Usage                `json:"usage"`
}

type ResponseChoice struct {
	ContentFilterResults ContentFilterResults `json:"content_filter_results"`
	FinishReason         string               `json:"finish_reason"`
	Index                int                  `json:"index"`
	Message              ChatMessage          `json:"message"`
}

type PromptFilterResult struct {
	ContentFilterResults ContentFilterResults `json:"content_filter_results"`
	PromptIndex          int                  `json:"prompt_index"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ApiTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type ApiToken struct {
	ApiKey    string
	ExpiresAt time.Time
}
type CacheData struct {
	OAuthToken string    `json:"oauth_token"`
	ApiKey     string    `json:"api_key"`
	ExpiresAt  time.Time `json:"expires_at"`
}
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type DeviceTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error,omitempty"`
}

type CopilotChat struct {
	OAuthToken string
	ApiToken   *ApiToken
	Client     *http.Client
	mu         sync.Mutex
}

// TODO: import a library to count the number of tokens in a string
func (m Model) MaxTokenCount() int {
	switch m {
	case Gpt4o:
		return 64000
	case Gpt4:
		return 32768
	case Gpt3_5Turbo:
		return 12288
	case O1Mini:
		return 20000
	case O1Preview:
		return 20000
	case Claude3_5Sonnet:
		return 200000
	default:
		return 0
	}
}

func NewCopilotChat(client *http.Client) *CopilotChat {
	if client == nil {
		client = &http.Client{}
	}

	chat := &CopilotChat{
		Client: client,
	}

	if err := chat.loadFromKeyring(); err != nil {
		log.Printf("Warning: Failed to load from keyring: %v", err)
	}

	return chat
}

func (self *CopilotChat) saveToKeyring() error {
	data := CacheData{
		OAuthToken: self.OAuthToken,
		ApiKey:     self.ApiToken.ApiKey,
		ExpiresAt:  self.ApiToken.ExpiresAt,
	}

	fileData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal keyring data: %v", err)
	}

	if err := keyring.Set(KEYRING_SERVICE, KEYRING_USER, string(fileData)); err != nil {
		return fmt.Errorf("failed to save to keyring: %v", err)
	}

	return nil
}

func (self *CopilotChat) loadFromKeyring() error {
	jsonData, err := keyring.Get(KEYRING_SERVICE, KEYRING_USER)
	if err != nil {
		if err == keyring.ErrNotFound {
			return nil // No credentials stored yet
		}
		return fmt.Errorf("failed to get credentials from keyring: %v", err)
	}

	var data CacheData
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return fmt.Errorf("failed to unmarshal Keyring data: %v", err)
	}

	// Always load OAuth token if it exists
	if data.OAuthToken != "" {
		self.OAuthToken = data.OAuthToken
	}

	// If we have a valid API key, use it
	if data.ApiKey != "" && data.ExpiresAt.After(time.Now()) {
		self.ApiToken = &ApiToken{
			ApiKey:    data.ApiKey,
			ExpiresAt: data.ExpiresAt,
		}
		fmt.Println("Loaded valid API key from keyring")
		return nil
	}

	// If we have OAuth token but no valid API key, fetch a new one
	if self.OAuthToken != "" {
		fmt.Println("OAuth token found, fetching new API key...")
		if err := self.fetchNewApiToken(); err != nil {
			return fmt.Errorf("failed to fetch new API token: %v", err)
		}
		fmt.Println("Successfully fetched new API key")
		return nil
	}

	return nil
}

func (self *CopilotChat) fetchNewApiToken() error {
	apiTokenReq, err := http.NewRequest(http.MethodGet, COPILOT_CHAT_AUTH_URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create API token request: %v", err)
	}

	apiTokenReq.Header.Set("Authorization", fmt.Sprintf("token %s", self.OAuthToken))
	setHeaders(apiTokenReq, "")

	apiTokenResp, err := self.Client.Do(apiTokenReq)
	if err != nil {
		return fmt.Errorf("failed to get API token: %v", err)
	}
	defer apiTokenResp.Body.Close()

	var apiTokenResponse ApiTokenResponse
	if err := json.NewDecoder(apiTokenResp.Body).Decode(&apiTokenResponse); err != nil {
		return fmt.Errorf("failed to decode API token response: %v", err)
	}

	self.ApiToken = &ApiToken{
		ApiKey:    apiTokenResponse.Token,
		ExpiresAt: time.Unix(apiTokenResponse.ExpiresAt, 0),
	}

	return self.saveToKeyring()
}

func setHeaders(req *http.Request, contentType string) {
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Editor-Version", EDITOR_VERSION)
	req.Header.Set("Copilot-Integration-Id", COPILOT_INTEGRATION_ID)
}

func (self *CopilotChat) Authenticate() error {
	// Try to load from keyring first
	if err := self.loadFromKeyring(); err == nil && self.IsAuthenticated() {
		return nil
	}

	self.mu.Lock()
	defer self.mu.Unlock()

	// Step 1: Request device and user codes
	deviceCodeReq, err := http.NewRequest(
		http.MethodPost,
		"https://github.com/login/device/code",
		strings.NewReader(fmt.Sprintf(
			"client_id=%s&scope=copilot",
			GITHUB_CLIENT_ID,
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to create device code request: %v", err)
	}
	deviceCodeReq.Header.Set("Accept", "application/json")
	deviceCodeReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := self.Client.Do(deviceCodeReq)
	if err != nil {
		return fmt.Errorf("failed to get device code: %v", err)
	}
	defer resp.Body.Close()

	var deviceCode DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceCode); err != nil {
		return fmt.Errorf("failed to decode device code response: %v", err)
	}

	// Step 2: Display user code and verification URL
	fmt.Printf("\nPlease visit: %s\n", deviceCode.VerificationUri)
	fmt.Printf("And enter code: %s\n\n", deviceCode.UserCode)

	// Step 3: Poll for the access token with timeout
	startTime := time.Now()
	attempts := 0

	// FIXME: There is probably a better way to do this
	for {
		if time.Since(startTime) >= MAX_AUTH_TIME {
			return fmt.Errorf("authentication timed out after 5 minutes")
		}

		time.Sleep(CHECK_INTERVAL)
		attempts++
		fmt.Printf("Checking for authentication... attempt %d\n", attempts)

		tokenReq, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token",
			strings.NewReader(fmt.Sprintf(
				"client_id=%s&device_code=%s&grant_type=urn:ietf:params:oauth:grant-type:device_code", GITHUB_CLIENT_ID,
				deviceCode.DeviceCode)))
		if err != nil {
			return fmt.Errorf("failed to create token request: %v", err)
		}
		tokenReq.Header.Set("Accept", "application/json")
		tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		tokenResp, err := self.Client.Do(tokenReq)
		if err != nil {
			return fmt.Errorf("failed to get access token: %v", err)
		}

		var tokenResponse DeviceTokenResponse
		if err := json.NewDecoder(tokenResp.Body).Decode(&tokenResponse); err != nil {
			tokenResp.Body.Close()
			return fmt.Errorf("failed to decode token response: %v", err)
		}
		tokenResp.Body.Close()

		if tokenResponse.Error == "authorization_pending" {
			fmt.Println("Login not detected. Please visit the URL and enter the code.")
			continue
		}
		if tokenResponse.Error != "" {
			if time.Since(startTime) >= MAX_AUTH_TIME {
				return fmt.Errorf("authentication timed out after 5 minutes")
			}
			continue
		}

		// Successfully got the access token
		self.OAuthToken = tokenResponse.AccessToken

		// Now get the Copilot API token using fetchNewApiToken
		if err := self.fetchNewApiToken(); err != nil {
			return fmt.Errorf("failed to fetch API token: %v", err)
		}

		fmt.Println("Successfully authenticated!")
		// Save the new credentials to cache
		if err := self.saveToKeyring(); err != nil {
			log.Printf("Warning: Failed to save credentials to keyring: %v", err)
		}
		return nil
	}
}

func (self *CopilotChat) IsAuthenticated() bool {
	if self.ApiToken == nil {
		return false
	}
	return self.ApiToken.ExpiresAt.After(time.Now())
}

func (self *CopilotChat) Chat(request Request) (string, error) {
	fmt.Println("Chatting with Copilot...")

	if !self.IsAuthenticated() {
		fmt.Println("Not authenticated with Copilot. Authenticating...")
		if err := self.Authenticate(); err != nil {
			return "", fmt.Errorf("authentication failed: %v", err)
		}
	}

	apiKey := self.ApiToken.ApiKey
	fmt.Println("Authenticated with Copilot!")
	fmt.Println("API Key: ", apiKey)

	litter.Dump(self)

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	fmt.Println("Mounting request body: ", string(requestBody))

	self.mu.Lock()
	defer self.mu.Unlock()

	req, err := http.NewRequest(http.MethodPost, COPILOT_CHAT_COMPLETION_URL, strings.NewReader(string(requestBody)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	setHeaders(req, "")

	response, err := self.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("failed to get completion: %s", string(body))
	}

	var chatResponse ChatResponse
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&chatResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResponse.Choices[0].Message.Content, nil
}

func main() {
	client := &http.Client{}
	fmt.Println("Starting...")
	copilotChat := NewCopilotChat(client)

	fmt.Println("Chatting...")

	err := copilotChat.Authenticate()
	if err != nil {
		if strings.Contains(err.Error(), "timed out") {
			log.Fatalf("Authentication process timed out. Please try again later.")
		}
		log.Fatalf("Error during authentication: %v", err)
	}

	fmt.Println("Authenticated!")

	messages := []ChatMessage{
		{
			Role:    RoleUser,
			Content: "Describe what is Lazygit in one sentence",
		},
	}

	request := Request{
		Intent:      true,
		N:           1,
		Stream:      false,
		Temperature: 0.1,
		Model:       Gpt4o,
		Messages:    messages,
	}

	response, err := copilotChat.Chat(request)
	if err != nil {
		log.Fatalf("Error during chat: %v", err)
	}

	fmt.Println(response)
}
