package commit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GenerateCommitMessageWithAi = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Generate a commit message with AI from the staged changes, overwriting an existing message",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// A fake OpenAI-compatible server. It only returns a commit message if
		// the request contains the staged diff and the recent commit subjects.
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Messages []struct {
					Content string `json:"content"`
				} `json:"messages"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)
			if r.URL.Path != "/v1/chat/completions" || len(req.Messages) != 2 ||
				!strings.Contains(req.Messages[1].Content, "+staged content") ||
				!strings.Contains(req.Messages[1].Content, "- initial commit") {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error":{"message":"unexpected request"}}`))
				return
			}
			_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Add staged file\n\nThis adds the staged file."}}]}`))
		}))

		aiConfig := &config.GetUserConfig().Git.Commit.AI
		aiConfig.BaseURL = server.URL + "/v1"
		aiConfig.APIKeyEnv = ""
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file", "staged content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("existing message").
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Generate commit message with AI")).
			Confirm()

		t.ExpectPopup().Alert().Title(Equals("Generate commit message with AI")).
			Content(Contains("overwrite the current commit message")).
			Confirm()

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("Add staged file")).
			SwitchToDescription().
			Content(Equals("This adds the staged file.")).
			SwitchToSummary().
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("Add staged file"),
				Contains("initial commit"),
			)
	},
})
