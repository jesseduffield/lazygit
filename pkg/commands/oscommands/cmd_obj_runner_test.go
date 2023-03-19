package oscommands

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

func getRunner() *cmdObjRunner {
	log := utils.NewDummyLog()
	return &cmdObjRunner{
		log:   log,
		guiIO: NewNullGuiIO(log),
	}
}

func TestProcessOutput(t *testing.T) {
	defaultPromptUserForCredential := func(ct CredentialType) string {
		switch ct {
		case Password:
			return "password"
		case Username:
			return "username"
		case Passphrase:
			return "passphrase"
		case PIN:
			return "pin"
		default:
			panic("unexpected credential type")
		}
	}

	scenarios := []struct {
		name                    string
		promptUserForCredential func(CredentialType) string
		output                  string
		expectedToWrite         string
	}{
		{
			name:                    "no output",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "",
			expectedToWrite:         "",
		},
		{
			name:                    "password prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Password:",
			expectedToWrite:         "password",
		},
		{
			name:                    "password prompt 2",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Bill's password:",
			expectedToWrite:         "password",
		},
		{
			name:                    "password prompt 3",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Password for 'Bill':",
			expectedToWrite:         "password",
		},
		{
			name:                    "username prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Username for 'Bill':",
			expectedToWrite:         "username",
		},
		{
			name:                    "passphrase prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Enter passphrase for key '123':",
			expectedToWrite:         "passphrase",
		},
		{
			name:                    "pin prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Enter PIN for key '123':",
			expectedToWrite:         "pin",
		},
		{
			name:                    "username and password prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Password:\nUsername for 'Alice':\n",
			expectedToWrite:         "passwordusername",
		},
		{
			name:                    "user submits empty credential",
			promptUserForCredential: func(ct CredentialType) string { return "" },
			output:                  "Password:\n",
			expectedToWrite:         "",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			runner := getRunner()
			reader := strings.NewReader(scenario.output)
			writer := &strings.Builder{}

			runner.processOutput(reader, writer, scenario.promptUserForCredential)

			if writer.String() != scenario.expectedToWrite {
				t.Errorf("expected to write '%s' but got '%s'", scenario.expectedToWrite, writer.String())
			}
		})
	}
}
