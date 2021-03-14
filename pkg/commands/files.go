package commands

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mgutz/str"
)

// CatFile obtains the content of a file
func (c *GitCommand) CatFile(fileName string) (string, error) {
	return c.OSCommand.RunCommandWithOutput("%s %s", c.OSCommand.Platform.CatCmd, c.OSCommand.Quote(fileName))
}

// StageFile stages a file
func (c *GitCommand) StageFile(fileName string) error {
	// renamed files look like "file1 -> file2"
	fileNames := strings.Split(fileName, models.RENAME_SEPARATOR)
	return c.OSCommand.RunCommand("git add -- %s", c.OSCommand.Quote(fileNames[len(fileNames)-1]))
}

// StageAll stages all files
func (c *GitCommand) StageAll() error {
	return c.OSCommand.RunCommand("git add -A")
}

// UnstageAll stages all files
func (c *GitCommand) UnstageAll() error {
	return c.OSCommand.RunCommand("git reset")
}

// UnStageFile unstages a file
func (c *GitCommand) UnStageFile(fileName string, tracked bool) error {
	command := "git rm --cached --force -- %s"
	if tracked {
		command = "git reset HEAD -- %s"
	}

	// renamed files look like "file1 -> file2"
	fileNames := strings.Split(fileName, models.RENAME_SEPARATOR)
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

	// we've got a file that represents a rename from one file to another. Unfortunately
	// our File abstraction fails to consider this case, so here we will refetch
	// all files, passing the --no-renames flag and then recursively call the function
	// again for the before file and after file. At some point we should fix the abstraction itself

	split := strings.Split(file.Name, models.RENAME_SEPARATOR)
	filesWithoutRenames := c.GetStatusFiles(GetStatusFileOptions{NoRenames: true})
	var beforeFile *models.File
	var afterFile *models.File
	for _, f := range filesWithoutRenames {
		if f.Name == split[0] {
			beforeFile = f
		}
		if f.Name == split[1] {
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

	// if the file isn't tracked, we assume you want to delete it
	quotedFileName := c.OSCommand.Quote(file.Name)
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := c.OSCommand.RunCommand("git reset -- %s", quotedFileName); err != nil {
			return err
		}
	}

	if !file.Tracked {
		return c.removeFile(file.Name)
	}
	return c.DiscardUnstagedFileChanges(file)
}

// DiscardUnstagedFileChanges directly
func (c *GitCommand) DiscardUnstagedFileChanges(file *models.File) error {
	quotedFileName := c.OSCommand.Quote(file.Name)
	return c.OSCommand.RunCommand("git checkout -- %s", quotedFileName)
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

func (c *GitCommand) WorktreeFileDiffCmdStr(file *models.File, plain bool, cached bool) string {
	cachedArg := ""
	trackedArg := "--"
	colorArg := c.colorArg()
	split := strings.Split(file.Name, models.RENAME_SEPARATOR) // in case of a renamed file we get the new filename
	fileName := c.OSCommand.Quote(split[len(split)-1])
	if cached {
		cachedArg = "--cached"
	}
	if !file.Tracked && !file.HasStagedChanges && !cached {
		trackedArg = "--no-index -- /dev/null"
	}
	if plain {
		colorArg = "never"
	}

	return fmt.Sprintf("git diff --submodule --no-ext-diff --color=%s %s %s %s", colorArg, cachedArg, trackedArg, fileName)
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

	return c.OSCommand.RunCommand("git apply %s %s", flagStr, c.OSCommand.Quote(filepath))
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
	return c.OSCommand.RunCommand("git checkout %s %s", commitSha, fileName)
}

// DiscardOldFileChanges discards changes to a file from an old commit
func (c *GitCommand) DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// check if file exists in previous commit (this command returns an error if the file doesn't exist)
	if err := c.OSCommand.RunCommand("git cat-file -e HEAD^:%s", fileName); err != nil {
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
	cmd, err := c.AmendHead()
	if cmd != nil {
		return errors.New("received unexpected pointer to cmd")
	}
	if err != nil {
		return err
	}

	// continue
	return c.GenericMergeOrRebaseAction("rebase", "continue")
}

// DiscardAnyUnstagedFileChanges discards any unstages file changes via `git checkout -- .`
func (c *GitCommand) DiscardAnyUnstagedFileChanges() error {
	return c.OSCommand.RunCommand("git checkout -- .")
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (c *GitCommand) RemoveTrackedFiles(name string) error {
	return c.OSCommand.RunCommand("git rm -r --cached %s", name)
}

// RemoveUntrackedFiles runs `git clean -fd`
func (c *GitCommand) RemoveUntrackedFiles() error {
	return c.OSCommand.RunCommand("git clean -fd")
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

// EditFile opens a file in a subprocess using whatever editor is available,
// falling back to core.editor, VISUAL, EDITOR, then vi
func (c *GitCommand) EditFile(filename string) (*exec.Cmd, error) {
	editor := c.GetConfigValue("core.editor")

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
		return nil, errors.New("No editor defined in $VISUAL, $EDITOR, or git config")
	}

	splitCmd := str.ToArgv(fmt.Sprintf("%s %s", editor, c.OSCommand.Quote(filename)))

	return c.OSCommand.PrepareSubProcess(splitCmd[0], splitCmd[1:]...), nil
}
