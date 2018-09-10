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
		testName string
		runCmd   func(string) error
		test     func(error)
	}

	scenarios := []scenario{
		{
			"Valid git repository",
			func(string) error {
				return nil
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Not a valid git repository",
			func(string) error {
				return fmt.Errorf("fatal: Not a git repository (or any of the parent directories): .git")
			},
			func(err error) {
				assert.Error(t, err)
				assert.Regexp(t, "fatal: .ot a git repository \\(or any of the parent directories\\): \\.git", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			s.test(verifyInGitRepo(s.runCmd))
		})
	}
}

func TestNavigateToRepoRootDirectory(t *testing.T) {
	type scenario struct {
		testName string
		stat     func(string) (os.FileInfo, error)
		chdir    func(string) error
		test     func(error)
	}

	scenarios := []scenario{
		{
			"Navigate to git repository",
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
			"An error occurred when getting path informations",
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
			"An error occurred when trying to move one path backward",
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
		t.Run(s.testName, func(t *testing.T) {
			s.test(navigateToRepoRootDirectory(s.stat, s.chdir))
		})
	}
}

func TestSetupRepositoryAndWorktree(t *testing.T) {
	type scenario struct {
		testName          string
		openGitRepository func(string) (*gogit.Repository, error)
		sLocalize         func(string) string
		test              func(*gogit.Repository, *gogit.Worktree, error)
	}

	scenarios := []scenario{
		{
			"A gitconfig parsing error occurred",
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
			"A gogit error occurred",
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
			"An error occurred cause git repository is a bare repository",
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
			"Setup done properly",
			func(string) (*gogit.Repository, error) {
				assert.NoError(t, os.RemoveAll("/tmp/lazygit-test"))
				r, err := gogit.PlainInit("/tmp/lazygit-test", false)
				assert.NoError(t, err)
				return r, nil
			},
			func(string) string { return "" },
			func(r *gogit.Repository, w *gogit.Worktree, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, w)
				assert.NotNil(t, r)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			s.test(setupRepositoryAndWorktree(s.openGitRepository, s.sLocalize))
		})
	}
}

func TestNewGitCommand(t *testing.T) {
	actual, err := os.Getwd()
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.Chdir(actual))
	}()

	type scenario struct {
		testName string
		setup    func()
		test     func(*GitCommand, error)
	}

	scenarios := []scenario{
		{
			"An error occurred, folder doesn't contains a git repository",
			func() {
				assert.NoError(t, os.Chdir("/tmp"))
			},
			func(gitCmd *GitCommand, err error) {
				assert.Error(t, err)
				assert.Regexp(t, "fatal: .ot a git repository \\(or any of the parent directories\\): \\.git", err.Error())
			},
		},
		{
			"New GitCommand object created",
			func() {
				assert.NoError(t, os.RemoveAll("/tmp/lazygit-test"))
				_, err := gogit.PlainInit("/tmp/lazygit-test", false)
				assert.NoError(t, err)
				assert.NoError(t, os.Chdir("/tmp/lazygit-test"))
			},
			func(gitCmd *GitCommand, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			s.setup()
			s.test(NewGitCommand(newDummyLog(), newDummyOSCommand(), i18n.NewLocalizer(newDummyLog())))
		})
	}
}

func TestGitCommandGetStashEntries(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func([]StashEntry)
	}

	scenarios := []scenario{
		{
			"No stash entries found",
			func(string, ...string) *exec.Cmd {
				return exec.Command("echo")
			},
			func(entries []StashEntry) {
				assert.Len(t, entries, 0)
			},
		},
		{
			"Several stash entries found",
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := newDummyGitCommand()
			gitCmd.OSCommand.command = s.command

			s.test(gitCmd.GetStashEntries())
		})
	}
}

func TestGitCommandGetStashEntryDiff(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"stash", "show", "-p", "--color", "stash@{1}"}, args)

		return exec.Command("echo")
	}

	_, err := gitCmd.GetStashEntryDiff(1)

	assert.NoError(t, err)
}

func TestGitCommandGetStatusFiles(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func([]File)
	}

	scenarios := []scenario{
		{
			"No files found",
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command("echo")
			},
			func(files []File) {
				assert.Len(t, files, 0)
			},
		},
		{
			"Several files found",
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := newDummyGitCommand()
			gitCmd.OSCommand.command = s.command

			s.test(gitCmd.GetStatusFiles())
		})
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
		testName string
		oldFiles []File
		newFiles []File
		test     func([]File)
	}

	scenarios := []scenario{
		{
			"Old file and new file are the same",
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
			"Several files to merge, with some identical",
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
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := newDummyGitCommand()

			s.test(gitCmd.MergeStatusFiles(s.oldFiles, s.newFiles))
		})
	}
}

func TestGitCommandUpstreamDifferentCount(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(string, string)
	}

	scenarios := []scenario{
		{
			"Can't retrieve pushable count",
			func(string, ...string) *exec.Cmd {
				return exec.Command("exit", "1")
			},
			func(pushableCount string, pullableCount string) {
				assert.EqualValues(t, "?", pushableCount)
				assert.EqualValues(t, "?", pullableCount)
			},
		},
		{
			"Can't retrieve pullable count",
			func(cmd string, args ...string) *exec.Cmd {
				if args[1] == "head..@{u}" {
					return exec.Command("exit", "1")
				}

				return exec.Command("echo")
			},
			func(pushableCount string, pullableCount string) {
				assert.EqualValues(t, "?", pushableCount)
				assert.EqualValues(t, "?", pullableCount)
			},
		},
		{
			"Retrieve pullable and pushable count",
			func(cmd string, args ...string) *exec.Cmd {
				if args[1] == "head..@{u}" {
					return exec.Command("echo", "10")
				}

				return exec.Command("echo", "11")
			},
			func(pushableCount string, pullableCount string) {
				assert.EqualValues(t, "11", pushableCount)
				assert.EqualValues(t, "10", pullableCount)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := newDummyGitCommand()
			gitCmd.OSCommand.command = s.command
			s.test(gitCmd.UpstreamDifferenceCount())
		})
	}
}

func TestGitCommandGetCommitsToPush(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func([]string)
	}

	scenarios := []scenario{
		{
			"Can't retrieve pushable commits",
			func(string, ...string) *exec.Cmd {
				return exec.Command("exit", "1")
			},
			func(pushables []string) {
				assert.EqualValues(t, []string{}, pushables)
			},
		},
		{
			"Retrieve pushable commits",
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command("echo", "8a2bb0e\n78976bc")
			},
			func(pushables []string) {
				assert.Len(t, pushables, 2)
				assert.EqualValues(t, []string{"8a2bb0e", "78976bc"}, pushables)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := newDummyGitCommand()
			gitCmd.OSCommand.command = s.command
			s.test(gitCmd.GetCommitsToPush())
		})
	}
}

func TestGitCommandRenameCommit(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"commit", "--allow-empty", "--amend", "-m", "test"}, args)

		return exec.Command("echo")
	}

	assert.NoError(t, gitCmd.RenameCommit("test"))
}

func TestGitCommandResetToCommit(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"reset", "78976bc"}, args)

		return exec.Command("echo")
	}

	assert.NoError(t, gitCmd.ResetToCommit("78976bc"))
}

func TestGitCommandNewBranch(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"checkout", "-b", "test"}, args)

		return exec.Command("echo")
	}

	assert.NoError(t, gitCmd.NewBranch("test"))
}

func TestGitCommandDeleteBranch(t *testing.T) {
	type scenario struct {
		testName string
		branch   string
		force    bool
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"Delete a branch",
			"test",
			false,
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"branch", "-d", "test"}, args)

				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"Force delete a branch",
			"test",
			true,
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"branch", "-D", "test"}, args)

				return exec.Command("echo")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := newDummyGitCommand()
			gitCmd.OSCommand.command = s.command
			s.test(gitCmd.DeleteBranch(s.branch, s.force))
		})
	}
}

func TestGitCommandMerge(t *testing.T) {
	gitCmd := newDummyGitCommand()
	gitCmd.OSCommand.command = func(cmd string, args ...string) *exec.Cmd {
		assert.EqualValues(t, "git", cmd)
		assert.EqualValues(t, []string{"merge", "--no-edit", "test"}, args)

		return exec.Command("echo")
	}

	assert.NoError(t, gitCmd.Merge("test"))
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
		t.Run(file.Name, func(t *testing.T) {
			assert.NotContains(t, gitCommand.Diff(file), "error")
		})
	}
}
