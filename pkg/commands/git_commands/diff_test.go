package git_commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestTextDifferRenderedDiff(t *testing.T) {
	var args []string
	var dir string
	var beforeContent, afterContent string

	// The two texts are diffed as files in a temporary directory that is removed
	// again afterwards, so read them here, while they are still there.
	runner := oscommands.NewFakeRunner(t).ExpectFunc("text diff",
		func(cmdObj *oscommands.CmdObj) bool {
			args = cmdObj.GetCmd().Args
			dir = cmdObj.GetCmd().Dir
			before, _ := os.ReadFile(filepath.Join(dir, "old message"))
			after, _ := os.ReadFile(filepath.Join(dir, "new message"))
			beforeContent, afterContent = string(before), string(after)
			return true
		},
		"\x1b[1mdiff --git old message new message\x1b[m\n"+
			"\x1b[1mindex 9ebe4a6..ec3885a 100644\x1b[m\n"+
			"\x1b[1m--- old message\x1b[m\n"+
			"\x1b[1m+++ new message\x1b[m\n"+
			"\x1b[36m@@ -1,3 +1,3 @@\x1b[m\n"+
			"-Fix the widget\n"+
			"+Fix the widget on startup\n",
		// --no-index implies --exit-code
		errors.New("exit status 1"))

	instance := buildDiffCommands(commonDeps{runner: runner})

	differ, err := instance.NewTextDiffer(config.DiffRendererValues{Width: 80, DiffContext: 3}, 24)
	assert.NoError(t, err)
	output, err := differ.RenderedDiff(
		NamedText{Name: "old message", Content: "Fix the widget\n"},
		NamedText{Name: "new message", Content: "Fix the widget on startup\n"})

	assert.NoError(t, err)
	assert.Equal(t, "-Fix the widget\n+Fix the widget on startup\n", output)

	assert.Equal(t, []string{
		"git", "-c", "diff.default.xfuncname=x^", "diff", "--no-ext-diff", "--unified=3", "--find-renames=50%",
		"--no-index", "--no-prefix", "--color=always", "--", "old message", "new message",
	}, args)
	assert.Equal(t, "Fix the widget\n", beforeContent)
	assert.Equal(t, "Fix the widget on startup\n", afterContent)

	assert.NoDirExists(t, dir)
}

func TestTextDifferPassesOnTheRenderersGitArgs(t *testing.T) {
	userConfig := config.GetDefaultConfig()
	userConfig.Git.DiffRenderers = []config.DiffRendererConfig{
		{Type: "rawGit", Args: []string{"--color-words"}},
	}

	var args []string
	runner := oscommands.NewFakeRunner(t).ExpectFunc("text diff",
		func(cmdObj *oscommands.CmdObj) bool {
			args = cmdObj.GetCmd().Args
			return true
		}, "", nil)

	instance := buildDiffCommands(commonDeps{runner: runner, userConfig: userConfig})

	differ, err := instance.NewTextDiffer(config.DiffRendererValues{Width: 80, DiffContext: 3}, 24)
	assert.NoError(t, err)
	_, err = differ.RenderedDiff(
		NamedText{Name: "old message", Content: "one"},
		NamedText{Name: "new message", Content: "two"})

	assert.NoError(t, err)
	assert.Contains(t, args, "--color-words")
}

func TestStripDiffHeaders(t *testing.T) {
	scenarios := []struct {
		name           string
		diff           string
		expectedOutput string
	}{
		{
			name: "colored diff",
			diff: "\x1b[1mdiff --git old message new message\x1b[m\n" +
				"\x1b[1mindex 9ebe4a6..ec3885a 100644\x1b[m\n" +
				"\x1b[1m--- old message\x1b[m\n" +
				"\x1b[1m+++ new message\x1b[m\n" +
				"\x1b[36m@@ -1,3 +1,3 @@\x1b[m\n" +
				"one \x1b[32mtwo\x1b[m\n",
			expectedOutput: "one \x1b[32mtwo\x1b[m\n",
		},
		{
			name: "uncolored diff with several hunks",
			diff: "diff --git old message new message\n" +
				"index 9ebe4a6..ec3885a 100644\n" +
				"--- old message\n" +
				"+++ new message\n" +
				"@@ -1,3 +1,3 @@\n" +
				"one two\n" +
				"@@ -20,3 +20,3 @@\n" +
				"three four\n",
			expectedOutput: "one two\n" +
				"@@ -20,3 +20,3 @@\n" +
				"three four\n",
		},
		{
			name:           "empty diff",
			diff:           "",
			expectedOutput: "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expectedOutput, stripDiffHeaders(s.diff))
		})
	}
}
