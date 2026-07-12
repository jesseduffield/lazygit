package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestParseGitFlowPrefixMap(t *testing.T) {
	type scenario struct {
		testName     string
		legacyOutput string
		nextOutput   string
		expected     map[string]string
	}
	scenarios := []scenario{
		{
			testName:     "empty inputs",
			legacyOutput: "",
			nextOutput:   "",
			expected:     map[string]string{},
		},
		{
			testName:     "legacy only",
			legacyOutput: "gitflow.prefix.feature feature/\ngitflow.prefix.hotfix hotfix/",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName:     "next only",
			legacyOutput: "",
			nextOutput:   "gitflow.branch.feature.prefix feature/\ngitflow.branch.release.prefix release/",
			expected:     map[string]string{"feature/": "feature", "release/": "release"},
		},
		{
			testName:     "legacy wins on duplicate prefix",
			legacyOutput: "gitflow.prefix.foo feature/",
			nextOutput:   "gitflow.branch.bar.prefix feature/",
			expected:     map[string]string{"feature/": "foo"},
		},
		{
			testName:     "prefix normalized with trailing slash from legacy",
			legacyOutput: "gitflow.prefix.feature feature",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature"},
		},
		{
			testName:     "malformed legacy lines skipped",
			legacyOutput: "gitflow.prefix.feature feature/\nnot-a-valid-line\ngitflow.prefix.hotfix hotfix/",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName:     "blank lines and whitespace ignored",
			legacyOutput: "  \n  gitflow.prefix.feature feature/  \n  \n  ",
			nextOutput:   "",
			expected:     map[string]string{"feature/": "feature"},
		},
	}
	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			got := parseGitFlowPrefixMap(s.legacyOutput, s.nextOutput)
			assert.Equal(t, s.expected, got)
		})
	}
}

func TestGetGitFlowPrefixMap(t *testing.T) {
	type scenario struct {
		testName               string
		gitConfigMockResponses map[string]string
		expected               map[string]string
	}
	scenarios := []scenario{
		{
			testName:               "empty when both queries empty",
			gitConfigMockResponses: nil,
			expected:               map[string]string{},
		},
		{
			testName: "correct map from legacy-only output",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature/\ngitflow.prefix.hotfix hotfix/",
			},
			expected: map[string]string{"feature/": "feature", "hotfix/": "hotfix"},
		},
		{
			testName: "correct map from git-flow-next-only output",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow\\.branch\\..*\\.prefix": "gitflow.branch.feature.prefix feature/\ngitflow.branch.release.prefix release/",
			},
			expected: map[string]string{"feature/": "feature", "release/": "release"},
		},
		{
			testName: "merged map with legacy winning when both have same prefix",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix":                 "gitflow.prefix.foo feature/",
				"--local --get-regexp gitflow\\.branch\\..*\\.prefix": "gitflow.branch.bar.prefix feature/",
			},
			expected: map[string]string{"feature/": "foo"},
		},
		{
			testName: "prefix normalized with trailing slash",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature",
			},
			expected: map[string]string{"feature/": "feature"},
		},
		{
			testName: "malformed lines skipped",
			gitConfigMockResponses: map[string]string{
				"--local --get-regexp gitflow.prefix": "gitflow.prefix.feature feature/\nnot-a-valid-line\n",
			},
			expected: map[string]string{"feature/": "feature"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			config := NewConfigCommands(common.NewDummyCommon(), git_config.NewFakeGitConfig(s.gitConfigMockResponses))
			got := config.GetGitFlowPrefixMap()
			assert.Equal(t, s.expected, got)
		})
	}
}

func TestIsTerminalPinentryProgram(t *testing.T) {
	scenarios := []struct {
		testName string
		program  string
		expected bool
	}{
		{testName: "empty", program: "", expected: false},
		{testName: "gtk", program: "/usr/bin/pinentry-gtk-2", expected: false},
		{testName: "qt", program: "/usr/bin/pinentry-qt", expected: false},
		{testName: "mac", program: "/usr/local/bin/pinentry-mac", expected: false},
		{testName: "plain pinentry (curses fallback binary)", program: "/usr/bin/pinentry", expected: false},
		{testName: "tty", program: "/usr/bin/pinentry-tty", expected: true},
		{testName: "curses", program: "/usr/bin/pinentry-curses", expected: true},
		{testName: "case insensitive", program: "/usr/bin/Pinentry-TTY", expected: true},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.Equal(t, s.expected, isTerminalPinentryProgram(s.program))
		})
	}
}

func TestParseGpgConfPinentryProgram(t *testing.T) {
	scenarios := []struct {
		testName string
		output   string
		expected string
	}{
		{
			testName: "empty output",
			output:   "",
			expected: "",
		},
		{
			testName: "no pinentry-program line",
			output:   "gpg-agent-name:0:0:0:1:0:0:0:0:\n",
			expected: "",
		},
		{
			testName: "value unset",
			output:   "pinentry-program:0:24:Pinentry to use for password entry:1:path:0:0:0:\n",
			expected: "",
		},
		{
			testName: "plain path",
			output:   "pinentry-program:0:24:Pinentry to use for password entry:1:path:0:0:0:/usr/bin/pinentry-curses\n",
			expected: "/usr/bin/pinentry-curses",
		},
		{
			testName: "percent-encoded path with spaces and colons",
			output:   `pinentry-program:0:24:Pinentry to use for password entry:1:path:0:0:0:/usr/local/my%20agent/pinentry-tty%3av2`,
			expected: "/usr/local/my agent/pinentry-tty:v2",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.Equal(t, s.expected, parseGpgConfPinentryProgram(s.output))
		})
	}
}

func TestParseGpgAgentConfPinentryProgram(t *testing.T) {
	scenarios := []struct {
		testName string
		content  string
		expected string
	}{
		{
			testName: "empty file",
			content:  "",
			expected: "",
		},
		{
			testName: "no pinentry-program setting",
			content:  "# a comment\ndefault-cache-ttl 600\n",
			expected: "",
		},
		{
			testName: "simple setting",
			content:  "pinentry-program /usr/bin/pinentry-curses\n",
			expected: "/usr/bin/pinentry-curses",
		},
		{
			testName: "setting among other lines, with leading whitespace",
			content:  "default-cache-ttl 600\n  pinentry-program /usr/bin/pinentry-tty\nmax-cache-ttl 7200\n",
			expected: "/usr/bin/pinentry-tty",
		},
		{
			testName: "commented out setting is ignored",
			content:  "# pinentry-program /usr/bin/pinentry-tty\n",
			expected: "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.Equal(t, s.expected, parseGpgAgentConfPinentryProgram(s.content))
		})
	}
}
