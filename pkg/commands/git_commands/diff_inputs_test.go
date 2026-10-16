package git_commands

import (
	"errors"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetDiffInputs(t *testing.T) {
	lsFilesArgs := []string{"-C", "/repo", "ls-files", "-s", "-z", "--", "file1", "conflicted", "untracked"}
	lsFilesOutput := "100644 1111111111111111111111111111111111111111 0\tfile1\x00" +
		"100644 2222222222222222222222222222222222222222 1\tconflicted\x00" +
		"100644 3333333333333333333333333333333333333333 2\tconflicted\x00"
	revParseArgs := []string{"rev-parse", "--verify", "--quiet", "HEAD"}

	fs := afero.NewMemMapFs()
	assert.NoError(t, afero.WriteFile(fs, "/repo/file1", []byte("one\n"), 0o644))
	assert.NoError(t, afero.WriteFile(fs, "/repo/untracked", []byte("two\n"), 0o644))
	modTime := time.Unix(1700000000, 123456789)
	assert.NoError(t, fs.Chtimes("/repo/file1", modTime, modTime))
	assert.NoError(t, fs.Chtimes("/repo/untracked", modTime, modTime))

	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(lsFilesArgs, lsFilesOutput, nil).
		ExpectGitArgs(revParseArgs, "abcdef\n", nil)
	instance := buildWorkingTreeCommands(commonDeps{runner: runner, fs: fs, repoPaths: MockRepoPaths("/repo")})

	inputs, err := instance.GetDiffInputs([]string{"file1", "conflicted", "untracked"}, true, true)
	assert.NoError(t, err)
	runner.CheckForMissingCalls()

	assert.Equal(t, &DiffInputs{
		head: "abcdef",
		indexEntries: map[string]string{
			"file1": "100644 1111111111111111111111111111111111111111 0\tfile1\n",
			"conflicted": "100644 2222222222222222222222222222222222222222 1\tconflicted\n" +
				"100644 3333333333333333333333333333333333333333 2\tconflicted\n",
			"untracked": "",
		},
		files: map[string]string{
			"file1":      "-rw-r--r-- 4 1700000000123456789",
			"conflicted": "",
			"untracked":  "-rw-r--r-- 4 1700000000123456789",
		},
	}, inputs)
	assert.ElementsMatch(t, []string{"file1", "conflicted", "untracked"}, inputs.Paths())
}

func TestGetDiffInputsOneSide(t *testing.T) {
	lsFilesArgs := []string{"-C", "/repo", "ls-files", "-s", "-z", "--", "file1"}
	revParseArgs := []string{"rev-parse", "--verify", "--quiet", "HEAD"}

	// The unstaged side doesn't need HEAD.
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs(lsFilesArgs, "", nil)
	instance := buildWorkingTreeCommands(commonDeps{runner: runner, fs: afero.NewMemMapFs(), repoPaths: MockRepoPaths("/repo")})
	inputs, err := instance.GetDiffInputs([]string{"file1"}, false, true)
	assert.NoError(t, err)
	runner.CheckForMissingCalls()
	assert.Equal(t, &DiffInputs{indexEntries: map[string]string{"file1": ""}, files: map[string]string{"file1": ""}}, inputs)

	// The staged side doesn't need the working tree. In a repo without commits, there is
	// no HEAD.
	runner = oscommands.NewFakeRunner(t).
		ExpectGitArgs(lsFilesArgs, "", nil).
		ExpectGitArgs(revParseArgs, "", errors.New("exit status 1"))
	instance = buildWorkingTreeCommands(commonDeps{runner: runner, fs: afero.NewMemMapFs(), repoPaths: MockRepoPaths("/repo")})
	inputs, err = instance.GetDiffInputs([]string{"file1"}, true, false)
	assert.NoError(t, err)
	runner.CheckForMissingCalls()
	assert.Equal(t, &DiffInputs{indexEntries: map[string]string{"file1": ""}}, inputs)

	// A failure to read the index is one.
	runner = oscommands.NewFakeRunner(t).
		ExpectGitArgs(lsFilesArgs, "", errors.New("fatal"))
	instance = buildWorkingTreeCommands(commonDeps{runner: runner, fs: afero.NewMemMapFs(), repoPaths: MockRepoPaths("/repo")})
	_, err = instance.GetDiffInputs([]string{"file1"}, true, true)
	assert.Error(t, err)
}

func TestDiffInputsSameFor(t *testing.T) {
	inputs := func(head string, entry string, file string) *DiffInputs {
		return &DiffInputs{
			head:         head,
			indexEntries: map[string]string{"file1": entry, "file2": "100644 2222 0\tfile2\n"},
			files:        map[string]string{"file1": file, "file2": "-rw-r--r-- 2 2"},
		}
	}
	base := inputs("abc", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 1")
	paths := []string{"file1"}

	scenarios := []struct {
		testName   string
		other      *DiffInputs
		stagedSide bool
		expected   bool
	}{
		{"same, unstaged side", inputs("abc", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 1"), false, true},
		{"same, staged side", inputs("abc", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 1"), true, true},
		{"index entry changed, unstaged side", inputs("abc", "100644 9999 0\tfile1\n", "-rw-r--r-- 1 1"), false, false},
		{"index entry changed, staged side", inputs("abc", "100644 9999 0\tfile1\n", "-rw-r--r-- 1 1"), true, false},
		{"file changed, unstaged side", inputs("abc", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 2"), false, false},
		{"file changed, staged side", inputs("abc", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 2"), true, true},
		{"HEAD moved, unstaged side", inputs("def", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 1"), false, true},
		{"HEAD moved, staged side", inputs("def", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 1"), true, false},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.Equal(t, s.expected, base.SameFor(s.other, paths, s.stagedSide))
		})
	}

	// Only the files asked about count.
	other := inputs("abc", "100644 1111 0\tfile1\n", "-rw-r--r-- 1 1")
	other.indexEntries["file2"] = "100644 9999 0\tfile2\n"
	assert.True(t, base.SameFor(other, paths, false))
}
