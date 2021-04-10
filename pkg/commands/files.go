package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/jesseduffield/lazygit/pkg/test"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// CatFile obtains the content of a file
func (c *GitCommand) CatFile(fileName string) (string, error) {
	return c.OSCommand.CatFile(fileName)
}

// StageFile stages a file
func (c *GitCommand) StageFile(fileName string) error {
	return c.RunCommand("git add -- %s", c.OSCommand.Quote(fileName))
}

// StageAll stages all files
func (c *GitCommand) StageAll() error {
	return c.RunCommand("git add -A")
}

// UnstageAll unstages all files
func (c *GitCommand) UnstageAll() error {
	return c.RunCommand("git reset")
}

// UnStageFile unstages a file
// we accept an array of filenames for the cases where a file has been renamed i.e.
// we accept the current name and the previous name
func (c *GitCommand) UnStageFile(fileNames []string, reset bool) error {
	command := "git rm --cached --force -- %s"
	if reset {
		command = "git reset HEAD -- %s"
	}

	for _, name := range fileNames {
		if err := c.OSCommand.RunCommand(command, c.OSCommand.Quote(name)); err != nil {
			return err
		}
	}
	return nil
}

func (c *GitCommand) BeforeAndAfterFileForRename(file *models.File) (*models.File, *models.File, error) {

	if !file.IsRename() {
		return nil, nil, errors.New("Expected renamed file")
	}

	// we've got a file that represents a rename from one file to another. Here we will refetch
	// all files, passing the --no-renames flag and then recursively call the function
	// again for the before file and after file.

	filesWithoutRenames := c.GetStatusFiles(GetStatusFileOptions{NoRenames: true})
	var beforeFile *models.File
	var afterFile *models.File
	for _, f := range filesWithoutRenames {
		if f.Name == file.PreviousName {
			beforeFile = f
		}

		if f.Name == file.Name {
			afterFile = f
		}
	}

	if beforeFile == nil || afterFile == nil {
		return nil, nil, errors.New("Could not find deleted file or new file for file rename")
	}

	if beforeFile.IsRename() || afterFile.IsRename() {
		// probably won't happen but we want to ensure we don't get an infinite loop
		return nil, nil, errors.New("Nested rename found")
	}

	return beforeFile, afterFile, nil
}

// DiscardAllFileChanges directly
func (c *GitCommand) DiscardAllFileChanges(file *models.File) error {
	if file.IsRename() {
		beforeFile, afterFile, err := c.BeforeAndAfterFileForRename(file)
		if err != nil {
			return err
		}

		if err := c.DiscardAllFileChanges(beforeFile); err != nil {
			return err
		}

		if err := c.DiscardAllFileChanges(afterFile); err != nil {
			return err
		}

		return nil
	}

	quotedFileName := c.OSCommand.Quote(file.Name)

	if file.ShortStatus == "AA" {
		if err := c.RunCommand("git checkout --ours --  %s", quotedFileName); err != nil {
			return err
		}
		if err := c.RunCommand("git add %s", quotedFileName); err != nil {
			return err
		}
		return nil
	}

	if file.ShortStatus == "DU" {
		return c.RunCommand("git rm %s", quotedFileName)
	}

	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := c.RunCommand("git reset -- %s", quotedFileName); err != nil {
			return err
		}
	}

	if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
		return nil
	}

	if file.Added {
		return c.removeFile(file.Name)
	}
	return c.DiscardUnstagedFileChanges(file)
}

func (c *GitCommand) DiscardAllDirChanges(node *filetree.FileNode) error {
	// this could be more efficient but we would need to handle all the edge cases
	return node.ForEachFile(c.DiscardAllFileChanges)
}

func (c *GitCommand) DiscardUnstagedDirChanges(node *filetree.FileNode) error {
	if err := c.RemoveUntrackedDirFiles(node); err != nil {
		return err
	}

	quotedPath := c.OSCommand.Quote(node.GetPath())
	if err := c.RunCommand("git checkout -- %s", quotedPath); err != nil {
		return err
	}

	return nil
}

func (c *GitCommand) RemoveUntrackedDirFiles(node *filetree.FileNode) error {
	untrackedFilePaths := node.GetPathsMatching(
		func(n *filetree.FileNode) bool { return n.File != nil && !n.File.GetIsTracked() },
	)

	for _, path := range untrackedFilePaths {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// DiscardUnstagedFileChanges directly
func (c *GitCommand) DiscardUnstagedFileChanges(file *models.File) error {
	quotedFileName := c.OSCommand.Quote(file.Name)
	return c.RunCommand("git checkout -- %s", quotedFileName)
}

// Ignore adds a file to the gitignore for the repo
func (c *GitCommand) Ignore(filename string) error {
	return c.OSCommand.AppendLineToFile(".gitignore", filename)
}

// WorktreeFileDiff returns the diff of a file
func (c *GitCommand) WorktreeFileDiff(file *models.File, plain bool, cached bool) string {
	// for now we assume an error means the file was deleted
	s, _ := c.OSCommand.RunCommandWithOutput(c.WorktreeFileDiffCmdStr(file, plain, cached))
	return s
}

func (c *GitCommand) WorktreeFileDiffCmdStr(node models.IFile, plain bool, cached bool) string {
	cachedArg := ""
	trackedArg := "--"
	colorArg := c.colorArg()
	path := c.OSCommand.Quote(node.GetPath())
	if cached {
		cachedArg = "--cached"
	}
	if !node.GetIsTracked() && !node.GetHasStagedChanges() && !cached {
		trackedArg = "--no-index -- /dev/null"
	}
	if plain {
		colorArg = "never"
	}

	return fmt.Sprintf("git diff --submodule --no-ext-diff --color=%s %s %s %s", colorArg, cachedArg, trackedArg, path)
}

func (c *GitCommand) ApplyPatch(patch string, flags ...string) error {
	filepath := filepath.Join(c.Config.GetUserConfigDir(), utils.GetCurrentRepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	c.Log.Infof("saving temporary patch to %s", filepath)
	if err := c.OSCommand.CreateFileWithContent(filepath, patch); err != nil {
		return err
	}

	flagStr := ""
	for _, flag := range flags {
		flagStr += " --" + flag
	}

	return c.RunCommand("git apply %s %s", flagStr, c.OSCommand.Quote(filepath))
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (c *GitCommand) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	cmdStr := c.ShowFileDiffCmdStr(from, to, reverse, fileName, plain)
	return c.OSCommand.RunCommandWithOutput(cmdStr)
}

func (c *GitCommand) ShowFileDiffCmdStr(from string, to string, reverse bool, fileName string, plain bool) string {
	colorArg := c.colorArg()
	if plain {
		colorArg = "never"
	}

	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	return fmt.Sprintf("git diff --submodule --no-ext-diff --no-renames --color=%s %s %s %s -- %s", colorArg, from, to, reverseFlag, fileName)
}

// CheckoutFile checks out the file for the given commit
func (c *GitCommand) CheckoutFile(commitSha, fileName string) error {
	return c.RunCommand("git checkout %s %s", commitSha, fileName)
}

// DiscardOldFileChanges discards changes to a file from an old commit
func (c *GitCommand) DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// check if file exists in previous commit (this command returns an error if the file doesn't exist)
	if err := c.RunCommand("git cat-file -e HEAD^:%s", fileName); err != nil {
		if err := c.OSCommand.Remove(fileName); err != nil {
			return err
		}
		if err := c.StageFile(fileName); err != nil {
			return err
		}
	} else if err := c.CheckoutFile("HEAD^", fileName); err != nil {
		return err
	}

	// amend the commit
	err := c.AmendHead()
	if err != nil {
		return err
	}

	// continue
	return c.GenericMergeOrRebaseAction("rebase", "continue")
}

// DiscardAnyUnstagedFileChanges discards any unstages file changes via `git checkout -- .`
func (c *GitCommand) DiscardAnyUnstagedFileChanges() error {
	return c.RunCommand("git checkout -- .")
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (c *GitCommand) RemoveTrackedFiles(name string) error {
	return c.RunCommand("git rm -r --cached %s", name)
}

// RemoveUntrackedFiles runs `git clean -fd`
func (c *GitCommand) RemoveUntrackedFiles() error {
	return c.RunCommand("git clean -fd")
}

// ResetAndClean removes all unstaged changes and removes all untracked files
func (c *GitCommand) ResetAndClean() error {
	submoduleConfigs, err := c.GetSubmoduleConfigs()
	if err != nil {
		return err
	}

	if len(submoduleConfigs) > 0 {
		if err := c.ResetSubmodules(submoduleConfigs); err != nil {
			return err
		}
	}

	if err := c.ResetHard("HEAD"); err != nil {
		return err
	}

	return c.RemoveUntrackedFiles()
}

func (c *GitCommand) EditFileCmdStr(filename string) (string, error) {
	editor := c.GetConfigValue("core.editor")

	if editor == "" {
		editor = c.OSCommand.Getenv("GIT_EDITOR")
	}
	if editor == "" {
		editor = c.OSCommand.Getenv("VISUAL")
	}
	if editor == "" {
		editor = c.OSCommand.Getenv("EDITOR")
	}
	if editor == "" {
		if err := c.OSCommand.RunCommand("which vi"); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return "", errors.New("No editor defined in $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
	}

	return fmt.Sprintf("%s %s", editor, c.OSCommand.Quote(filename)), nil
}

func TestGitCommandApplyPatch(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			func(cmd string, args ...string) *exec.Cmd {
				assert.Equal(t, "git", cmd)
				assert.EqualValues(t, []string{"apply", "--cached"}, args[0:2])
				filename := args[2]
				content, err := ioutil.ReadFile(filename)
				assert.NoError(t, err)

				assert.Equal(t, "test", string(content))

				return secureexec.Command("echo", "done")
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"command returns error",
			func(cmd string, args ...string) *exec.Cmd {
				assert.Equal(t, "git", cmd)
				assert.EqualValues(t, []string{"apply", "--cached"}, args[0:2])
				filename := args[2]
				// TODO: Ideally we want to mock out OSCommand here so that we're not
				// double handling testing it's CreateTempFile functionality,
				// but it is going to take a bit of work to make a proper mock for it
				// so I'm leaving it for another PR
				content, err := ioutil.ReadFile(filename)
				assert.NoError(t, err)

				assert.Equal(t, "test", string(content))

				return secureexec.Command("test")
			},
			func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.ApplyPatch("test", "cached"))
		})
	}
}

// TestGitCommandDiscardOldFileChanges is a function.
func TestGitCommandDiscardOldFileChanges(t *testing.T) {
	type scenario struct {
		testName          string
		getGitConfigValue func(string) (string, error)
		commits           []*models.Commit
		commitIndex       int
		fileName          string
		command           func(string, ...string) *exec.Cmd
		test              func(error)
	}

	scenarios := []scenario{
		{
			"returns error when index outside of range of commits",
			func(string) (string, error) {
				return "", nil
			},
			[]*models.Commit{},
			0,
			"test999.txt",
			nil,
			func(err error) {
				assert.Error(t, err)
			},
		},
		{
			"returns error when using gpg",
			func(string) (string, error) {
				return "true", nil
			},
			[]*models.Commit{{Name: "commit", Sha: "123456"}},
			0,
			"test999.txt",
			nil,
			func(err error) {
				assert.Error(t, err)
			},
		},
		{
			"checks out file if it already existed",
			func(string) (string, error) {
				return "", nil
			},
			[]*models.Commit{
				{Name: "commit", Sha: "123456"},
				{Name: "commit2", Sha: "abcdef"},
			},
			0,
			"test999.txt",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  "git rebase --interactive --autostash --keep-empty abcdef",
					Replace: "echo",
				},
				{
					Expect:  "git cat-file -e HEAD^:test999.txt",
					Replace: "echo",
				},
				{
					Expect:  "git checkout HEAD^ test999.txt",
					Replace: "echo",
				},
				{
					Expect:  "git commit --amend --no-edit --allow-empty",
					Replace: "echo",
				},
				{
					Expect:  "git rebase --continue",
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
		// test for when the file was created within the commit requires a refactor to support proper mocks
		// currently we'd need to mock out the os.Remove function and that's gonna introduce tech debt
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			gitCmd.getGitConfigValue = s.getGitConfigValue
			s.test(gitCmd.DiscardOldFileChanges(s.commits, s.commitIndex, s.fileName))
		})
	}
}

// TestGitCommandDiscardUnstagedFileChanges is a function.
func TestGitCommandDiscardUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		file     *models.File
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			&models.File{Name: "test.txt"},
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  `git checkout -- "test.txt"`,
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.DiscardUnstagedFileChanges(s.file))
		})
	}
}

// TestGitCommandDiscardAnyUnstagedFileChanges is a function.
func TestGitCommandDiscardAnyUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  `git checkout -- .`,
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.DiscardAnyUnstagedFileChanges())
		})
	}
}

// TestGitCommandRemoveUntrackedFiles is a function.
func TestGitCommandRemoveUntrackedFiles(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			test.CreateMockCommand(t, []*test.CommandSwapper{
				{
					Expect:  `git clean -fd`,
					Replace: "echo",
				},
			}),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	gitCmd := NewDummyGitCommand()

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd.OSCommand.Command = s.command
			s.test(gitCmd.RemoveUntrackedFiles())
		})
	}
}

// TestEditFileCmdStr is a function.
func TestEditFileCmdStr(t *testing.T) {
	type scenario struct {
		filename          string
		command           func(string, ...string) *exec.Cmd
		getenv            func(string) string
		getGitConfigValue func(string) (string, error)
		test              func(string, error)
	}

	scenarios := []scenario{
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmdStr string, err error) {
				assert.EqualError(t, err, "No editor defined in $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				return ""
			},
			func(cf string) (string, error) {
				return "nano", nil
			},
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "nano \"test\"", cmdStr)
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				if env == "VISUAL" {
					return "nano"
				}

				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmdStr string, err error) {
				assert.NoError(t, err)
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("exit", "1")
			},
			func(env string) string {
				if env == "EDITOR" {
					return "emacs"
				}

				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "emacs \"test\"", cmdStr)
			},
		},
		{
			"test",
			func(name string, arg ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("echo")
			},
			func(env string) string {
				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "vi \"test\"", cmdStr)
			},
		},
		{
			"file/with space",
			func(name string, args ...string) *exec.Cmd {
				assert.Equal(t, "which", name)
				return secureexec.Command("echo")
			},
			func(env string) string {
				return ""
			},
			func(cf string) (string, error) {
				return "", nil
			},
			func(cmdStr string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "vi \"file/with space\"", cmdStr)
			},
		},
	}

	for _, s := range scenarios {
		gitCmd := NewDummyGitCommand()
		gitCmd.OSCommand.Command = s.command
		gitCmd.OSCommand.Getenv = s.getenv
		gitCmd.getGitConfigValue = s.getGitConfigValue
		s.test(gitCmd.EditFileCmdStr(s.filename))
	}
}
