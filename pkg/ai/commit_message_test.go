package ai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCommitMessage(t *testing.T) {
	var gotAuth, gotPath string
	var gotReq chatRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotReq)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"<think>hmm</think>\n` + "```\\nAdd foo\\n\\nBecause bar.\\n```" + `"}}]}`))
	}))
	defer server.Close()

	msg, err := GenerateCommitMessage(context.Background(), Options{
		BaseURL: server.URL + "/openai/v1/",
		APIKey:  "secret",
		Model:   "test-model",
	}, Input{
		Diff:                 "diff --git a/foo b/foo",
		RecentCommitSubjects: []string{"Fix thing"},
	})

	assert.NoError(t, err)
	assert.Equal(t, "Add foo\n\nBecause bar.", msg)
	assert.Equal(t, "Bearer secret", gotAuth)
	assert.Equal(t, "/openai/v1/chat/completions", gotPath)
	assert.Equal(t, "test-model", gotReq.Model)
	assert.Len(t, gotReq.Messages, 2)
	assert.Equal(t, DefaultSystemPrompt, gotReq.Messages[0].Content)
	assert.Contains(t, gotReq.Messages[1].Content, "- Fix thing")
	assert.Contains(t, gotReq.Messages[1].Content, "diff --git a/foo b/foo")
}

func TestGenerateCommitMessageAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"Rate limit reached"}}`))
	}))
	defer server.Close()

	_, err := GenerateCommitMessage(context.Background(), Options{BaseURL: server.URL, Model: "m"}, Input{Diff: "x"})
	assert.EqualError(t, err, "AI request failed (HTTP 429): Rate limit reached")
}

func TestGenerateCommitMessageNoChanges(t *testing.T) {
	_, err := GenerateCommitMessage(context.Background(), Options{}, Input{Diff: "  \n"})
	assert.ErrorIs(t, err, ErrNoChanges)
}

func TestBuildUserPromptTruncates(t *testing.T) {
	prompt := BuildUserPrompt(Input{Diff: strings.Repeat("a", 100), DiffStat: " foo | 1 +"}, 10)
	assert.Contains(t, prompt, "Summary of staged changes:\nfoo | 1 +")
	assert.Contains(t, prompt, "Staged diff:\naaaaaaaaaa\n[... diff truncated ...]")
	assert.NotContains(t, prompt, strings.Repeat("a", 11))
}

func TestCleanMessage(t *testing.T) {
	scenarios := []struct {
		in, out string
	}{
		{"  Add foo  \n", "Add foo"},
		{"\"Add foo\"", "Add foo"},
		{"```text\nAdd foo\n\nBody\n```", "Add foo\n\nBody"},
		{"<think>\nreasoning\n</think>\n\nFix bar", "Fix bar"},
	}
	for _, s := range scenarios {
		assert.Equal(t, s.out, CleanMessage(s.in))
	}
}
