package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	gogit "gopkg.in/src-d/go-git.v4"
)

type fileInfoMock struct {
	name        string
	size        int64
	fileMode    os.FileMode
	fileModTime time.Time
	isDir       bool
	sys         interface{}
}

func (f fileInfoMock) Name() string {
	return f.name
}

func (f fileInfoMock) Size() int64 {
	return f.size
}

func (f fileInfoMock) Mode() os.FileMode {
	return f.fileMode
}

func (f fileInfoMock) ModTime() time.Time {
	return f.fileModTime
}

func (f fileInfoMock) IsDir() bool {
	return f.isDir
}

func (f fileInfoMock) Sys() interface{} {
	return f.sys
}

func newDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
}

func newDummyGitCommand() *GitCommand {
	return &GitCommand{
		Log:       newDummyLog(),
		OSCommand: newDummyOSCommand(),
		Tr:        i18n.NewLocalizer(newDummyLog()),
	}
}

func TestVerifyInGitRepo(t *testing.T) {
	type scenario struct {
		runCmdWithOutput func(string) (string, error)
		test             func(error)
	}

	scenarios := []scenario{
		{
			func(string) (string, error) {
				return "", nil
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			func(string) (string, error) {
				return "", ErrGitRepositoryInvalid
			},
			func(err error) {
				assert.Error(t, err)
				assert.Equal(t, ErrGitRepositoryInvalid, err)
			},
		},
	}

	for _, s := range scenarios {
		s.test(verifyInGitRepo(s.runCmdWithOutput))
	}
}

func TestNavigateToRepoRootDirectory(t *testing.T) {
	type scenario struct {
		stat  func(string) (os.FileInfo, error)
		chdir func(string) error
		test  func(error)
	}

	scenarios := []scenario{
		{
			func(string) (os.FileInfo, error) {
				return fileInfoMock{isDir: true}, nil
			},
			func(string) error {
				return nil
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			func(string) (os.FileInfo, error) {
				return nil, fmt.Errorf("An error occurred")
			},
			func(string) error {
				return nil
			},
			func(err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "An error occurred")
			},
		},
		{
			func(string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			func(string) error {
				return fmt.Errorf("An error occurred")
			},
			func(err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "An error occurred")
			},
		},
		{
			func(string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			func(string) error {
				return fmt.Errorf("An error occurred")
			},
			func(err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "An error occurred")
			},
		},
	}

	for _, s := range scenarios {
		s.test(navigateToRepoRootDirectory(s.stat, s.chdir))
	}
}

func TestSetupRepositoryAndWorktree(t *testing.T) {
	type scenario struct {
		openGitRepository func(string) (*gogit.Repository, error)
		sLocalize         func(string) string
		test              func(*gogit.Repository, *gogit.Worktree, error)
	}

	scenarios := []scenario{
		{
			func(string) (*gogit.Repository, error) {
				return nil, fmt.Errorf(`unquoted '\' must be followed by new line`)
			},
			func(string) string {
				return "error translated"
			},
			func(r *gogit.Repository, w *gogit.Worktree, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "error translated")
			},
		},
		{
			func(string) (*gogit.Repository, error) {
				return nil, fmt.Errorf("Error from inside gogit")
			},
			func(string) string { return "" },
			func(r *gogit.Repository, w *gogit.Worktree, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "Error from inside gogit")
			},
		},
		{
			func(string) (*gogit.Repository, error) {
				return &gogit.Repository{}, nil
			},
			func(string) string { return "" },
			func(r *gogit.Repository, w *gogit.Worktree, err error) {
				assert.Error(t, err)
				assert.Equal(t, gogit.ErrIsBareRepository, err)
			},
		},
		{
			func(string) (*gogit.Repository, error) {
				assert.NoError(t, os.RemoveAll("/tmp/lazygit-test"))
				r, err := gogit.PlainInit("/tmp/lazygit-test", false)
				assert.NoError(t, err)
				return r, nil
			},
			func(string) string { return "" },
			func(r *gogit.Repository, w *gogit.Worktree, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s.test(setupRepositoryAndWorktree(s.openGitRepository, s.sLocalize))
	}
}

func TestGitCommandGetStashEntries(t *testing.T) {
	type scenario struct {
		command func(string, ...string) *exec.Cmd
		test    func([]StashEntry)
	}

	scenarios := []scenario{
		{
			func(string, ...string) *exec.Cmd {
				return exec.Command("echo")
			},
			func(entries []StashEntry) {
				assert.Len(t, entries, 0)
			},
		},
		{
			func(string, ...string) *exec.Cmd {
				return exec.Command("echo", "WIP on add-pkg-commands-test: 55c6af2 increase parallel build\nWIP on master: bb86a3f update github template")
			},
			func(entries []StashEntry) {
				expected := []StashEntry{
					{
						0,
						"WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
						"WIP on add-pkg-commands-test: 55c6af2 increase parallel build",
					},
					{
						1,
						"WIP on master: bb86a3f update github template",
						"WIP on master: bb86a3f update github template",
					},
				}

				assert.Len(t, entries, 2)
				assert.EqualValues(t, expected, entries)
			},
		},
	}

	for _, s := range scenarios {
		gitCmd := newDummyGitCommand()
		gitCmd.OSCommand.command = s.command

		s.test(gitCmd.GetStashEntries())
	}
}

func TestGetStashEntryDiff(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "show", "-p", "--color", "stash@{1}"}, args)

		return exec.Command("echo")
	}

	_, err := gitCmd.GetStashEntryDiff(1)

	assert.NoError(t, err)
}

func TestGetStatusFiles(t *testing.T) {
	type scenario struct {
		command func(string, ...string) *exec.Cmd
		test    func([]File)
	}

	scenarios := []scenario{
		{
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command("echo")
			},
			func(files []File) {
				assert.Len(t, files, 0)
			},
		},
		{
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command(
					"echo",
					"MM file1.txt\nA  file3.txt\nAM file2.txt\n?? file4.txt",
				)
			},
			func(files []File) {
				assert.Len(t, files, 4)

				expected := []File{
					{
						Name:               "file1.txt",
						HasStagedChanges:   true,
						HasUnstagedChanges: true,
						Tracked:            true,
						Deleted:            false,
						HasMergeConflicts:  false,
						DisplayString:      "MM file1.txt",
						Type:               "other",
					},
					{
						Name:               "file3.txt",
						HasStagedChanges:   true,
						HasUnstagedChanges: false,
						Tracked:            false,
						Deleted:            false,
						HasMergeConflicts:  false,
						DisplayString:      "A  file3.txt",
						Type:               "other",
					},
					{
						Name:               "file2.txt",
						HasStagedChanges:   true,
						HasUnstagedChanges: true,
						Tracked:            false,
						Deleted:            false,
						HasMergeConflicts:  false,
						DisplayString:      "AM file2.txt",
						Type:               "other",
					},
					{
						Name:               "file4.txt",
						HasStagedChanges:   false,
						HasUnstagedChanges: true,
						Tracked:            false,
						Deleted:            false,
						HasMergeConflicts:  false,
						DisplayString:      "?? file4.txt",
						Type:               "other",
					},
				}

				assert.EqualValues(t, expected, files)
			},
		},
	}

	for _, s := range scenarios {
		gitCmd := newDummyGitCommand()
		gitCmd.OSCommand.command = s.command

		s.test(gitCmd.GetStatusFiles())
	}
}

func TestGitCommandStashDo(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "drop", "stash@{1}"}, args)

		return exec.Command("echo")
	}

	assert.NoError(t, gitCmd.StashDo(1, "drop"))
}

func TestGitCommandStashSave(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "save", "A stash message"}, args)

		return exec.Command("echo")
	}

	assert.NoError(t, gitCmd.StashSave("A stash message"))
}

func TestGitCommandCommitAmend(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"commit", "--amend", "--allow-empty"}, args)

		return exec.Command("echo")
	}

	_, err := gitCmd.PrepareCommitAmendSubProcess().CombinedOutput()
	assert.NoError(t, err)
}

func TestGitCommandMergeStatusFiles(t *testing.T) {
	type scenario struct {
		oldFiles []File
		newFiles []File
		test     func([]File)
	}

	scenarios := []scenario{
		{
			[]File{},
			[]File{
				{
					Name: "new_file.txt",
				},
			},
			func(files []File) {
				expected := []File{
					{
						Name: "new_file.txt",
					},
				}

				assert.Len(t, files, 1)
				assert.EqualValues(t, expected, files)
			},
		},
		{
			[]File{
				{
					Name: "new_file1.txt",
				},
				{
					Name: "new_file2.txt",
				},
				{
					Name: "new_file3.txt",
				},
			},
			[]File{
				{
					Name: "new_file4.txt",
				},
				{
					Name: "new_file5.txt",
				},
				{
					Name: "new_file1.txt",
				},
			},
			func(files []File) {
				expected := []File{
					{
						Name: "new_file1.txt",
					},
					{
						Name: "new_file4.txt",
					},
					{
						Name: "new_file5.txt",
					},
				}

				assert.Len(t, files, 3)
				assert.EqualValues(t, expected, files)
			},
		},
	}

	for _, s := range scenarios {
		gitCmd := newDummyGitCommand()

		s.test(gitCmd.MergeStatusFiles(s.oldFiles, s.newFiles))
	}
}

func TestGitCommandDiff(t *testing.T) {
	gitCommand := newDummyGitCommand()
	assert.NoError(t, test.GenerateRepo("lots_of_diffs.sh"))

	files := []File{
		{
			Name:               "deleted_staged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            true,
			Deleted:            true,
			HasMergeConflicts:  false,
			DisplayString:      " D deleted_staged",
		},
		{
			Name:               "file with space staged",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "A  \"file with space staged\"",
		},
		{
			Name:               "file with space unstaged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "?? file with space unstaged",
		},
		{
			Name:               "modified_unstaged",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            true,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "M  modified_unstaged",
		},
		{
			Name:               "modified_staged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            true,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      " M modified_staged",
		},
		{
			Name:               "renamed_before -> renamed_after",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            true,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "R  renamed_before -> renamed_after",
		},
		{
			Name:               "untracked_unstaged",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "?? untracked_unstaged",
		},
		{
			Name:               "untracked_staged",
			HasStagedChanges:   true,
			HasUnstagedChanges: false,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "A  untracked_staged",
		},
		{
			Name:               "master",
			HasStagedChanges:   false,
			HasUnstagedChanges: true,
			Tracked:            false,
			Deleted:            false,
			HasMergeConflicts:  false,
			DisplayString:      "?? master",
		},
	}

	for _, file := range files {
		assert.NotContains(t, gitCommand.Diff(file), "error")
	}
}
