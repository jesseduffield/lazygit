package gui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const gitOutputHeader = "Git output:"

func englishCopyToClipboardLogLineMatcher() func(string) bool {
	return func(line string) bool {
		return logLineMatchesTemplate(line, "Copying '{{.str}}' to clipboard", "{{.str}}")
	}
}

func TestGitOutputBlocksFromCommandLogLines(t *testing.T) {
	t.Parallel()

	lines := []string{
		"Push",
		"  git push",
		"",
		gitOutputHeader,
		"line1",
		"line2",
	}

	assert.Equal(t, []string{"Push\n  git push\n\nGit output:\nline1\nline2"}, gitOutputBlocksFromCommandLogLines(lines, gitOutputHeader, englishCopyToClipboardLogLineMatcher()))
}

func TestGitOutputBlocksIncludeIndentedStderr(t *testing.T) {
	t.Parallel()

	lines := []string{
		"Push",
		"  git push",
		gitOutputHeader,
		"  at foo.go:10",
		"  at bar.go:20",
		"hook failed",
	}

	assert.Equal(t, []string{"Push\n  git push\n\nGit output:\n  at foo.go:10\n  at bar.go:20\nhook failed"}, gitOutputBlocksFromCommandLogLines(lines, gitOutputHeader, englishCopyToClipboardLogLineMatcher()))
}

func TestGitOutputBlocksIncludeToolErrorWithIndentedContext(t *testing.T) {
	t.Parallel()

	lines := []string{
		"Push",
		"  git push",
		gitOutputHeader,
		"Error: validation failed",
		"  line 42: syntax error",
		"more output",
		"Stage file",
		"  git add foo",
		gitOutputHeader,
		"second command output",
	}

	assert.Equal(t, []string{
		"Push\n  git push\n\nGit output:\nError: validation failed\n  line 42: syntax error\nmore output",
		"Stage file\n  git add foo\n\nGit output:\nsecond command output",
	}, gitOutputBlocksFromCommandLogLines(lines, gitOutputHeader, englishCopyToClipboardLogLineMatcher()))
}

func TestGitOutputBlocksSkipCopyNotifications(t *testing.T) {
	t.Parallel()

	lines := []string{
		"Push",
		"  git push",
		gitOutputHeader,
		"hook line",
		"  Copying 'hook line' to clipboard",
		"hook line 2",
	}

	assert.Equal(t, []string{"Push\n  git push\n\nGit output:\nhook line\nhook line 2"}, gitOutputBlocksFromCommandLogLines(lines, gitOutputHeader, englishCopyToClipboardLogLineMatcher()))
}

func TestGitOutputBlocksEndAtNextCommandLogEntry(t *testing.T) {
	t.Parallel()

	lines := []string{
		"Push",
		"  git push",
		gitOutputHeader,
		"first command output",
		"Stage file",
		"  git add foo",
		"",
		gitOutputHeader,
		"second command output",
	}

	assert.Equal(t, []string{
		"Push\n  git push\n\nGit output:\nfirst command output",
		"Stage file\n  git add foo\n\nGit output:\nsecond command output",
	}, gitOutputBlocksFromCommandLogLines(lines, gitOutputHeader, englishCopyToClipboardLogLineMatcher()))
}

func TestGitOutputBlocksMultipleBlocksJoined(t *testing.T) {
	t.Parallel()

	lines := []string{
		"Push",
		"  git push",
		gitOutputHeader,
		"first command",
		"Pull",
		"  git pull",
		gitOutputHeader,
		"second command",
	}

	blocks := gitOutputBlocksFromCommandLogLines(lines, gitOutputHeader, englishCopyToClipboardLogLineMatcher())
	assert.Equal(t, "Push\n  git push\n\nGit output:\nfirst command\n\nPull\n  git pull\n\nGit output:\nsecond command", joinGitOutputBlocks(blocks))
}

func TestLogLineMatchesCopyToClipboardTemplate(t *testing.T) {
	t.Parallel()

	matcher := englishCopyToClipboardLogLineMatcher()
	assert.True(t, matcher("  Copying 'hook line' to clipboard"))
	assert.False(t, matcher("Push"))
}

func joinGitOutputBlocks(blocks []string) string {
	result := ""
	for i, block := range blocks {
		if i > 0 {
			result += "\n\n"
		}
		result += block
	}
	return result
}
