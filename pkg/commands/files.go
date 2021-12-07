package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// CatFile obtains the content of a file
func (c *GitCommand) CatFile(fileName string) (string, error) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", nil
	}
	return string(buf), nil
}

func (c *GitCommand) OpenMergeToolCmdObj() oscommands.ICmdObj {
	return c.NewCmdObj("git mergetool")
}

func (c *GitCommand) OpenMergeTool() error {
	return c.Run(c.OpenMergeToolCmdObj())
}

// StageFile stages a file
func (c *GitCommand) StageFile(fileName string) error {
	return c.Run(c.NewCmdObj("git add -- " + c.OSCommand.Quote(fileName)))
}

// StageAll stages all files
func (c *GitCommand) StageAll() error {
	return c.Run(c.NewCmdObj("git add -A"))
}

// UnstageAll unstages all files
func (c *GitCommand) UnstageAll() error {
	return c.Run(c.NewCmdObj("git reset"))
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
		cmdObj := c.NewCmdObj(fmt.Sprintf(command, c.OSCommand.Quote(name)))
		if err := c.Run(cmdObj); err != nil {
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
		if err := c.Run(c.NewCmdObj("git checkout --ours --  " + quotedFileName)); err != nil {
			return err
		}
		if err := c.Run(c.NewCmdObj("git add -- " + quotedFileName)); err != nil {
			return err
		}
		return nil
	}

	if file.ShortStatus == "DU" {
		return c.Run(c.NewCmdObj("git rm -- " + quotedFileName))
	}

	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := c.Run(c.NewCmdObj("git reset -- " + quotedFileName)); err != nil {
			return err
		}
	}

	if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
		return nil
	}

	if file.Added {
		return c.OSCommand.RemoveFile(file.Name)
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
	if err := c.Run(c.NewCmdObj("git checkout -- " + quotedPath)); err != nil {
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
	return c.Run(c.NewCmdObj("git checkout -- " + quotedFileName))
}

// Ignore adds a file to the gitignore for the repo
func (c *GitCommand) Ignore(filename string) error {
	return c.OSCommand.AppendLineToFile(".gitignore", filename)
}

// WorktreeFileDiff returns the diff of a file
func (c *GitCommand) WorktreeFileDiff(file *models.File, plain bool, cached bool, ignoreWhitespace bool) string {
	// for now we assume an error means the file was deleted
	s, _ := c.OSCommand.RunWithOutput(c.WorktreeFileDiffCmdObj(file, plain, cached, ignoreWhitespace))
	return s
}

func (c *GitCommand) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool, ignoreWhitespace bool) oscommands.ICmdObj {
	cachedArg := ""
	trackedArg := "--"
	colorArg := c.colorArg()
	quotedPath := c.OSCommand.Quote(node.GetPath())
	ignoreWhitespaceArg := ""
	contextSize := c.Config.GetUserConfig().Git.DiffContextSize
	if cached {
		cachedArg = "--cached"
	}
	if !node.GetIsTracked() && !node.GetHasStagedChanges() && !cached {
		trackedArg = "--no-index -- /dev/null"
	}
	if plain {
		colorArg = "never"
	}
	if ignoreWhitespace {
		ignoreWhitespaceArg = "--ignore-all-space"
	}

	cmdStr := fmt.Sprintf("git diff --submodule --no-ext-diff --unified=%d --color=%s %s %s %s %s", contextSize, colorArg, ignoreWhitespaceArg, cachedArg, trackedArg, quotedPath)

	return c.NewCmdObj(cmdStr)
}

func (c *GitCommand) ApplyPatch(patch string, flags ...string) error {
	filepath := filepath.Join(c.Config.GetTempDir(), utils.GetCurrentRepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	c.Log.Infof("saving temporary patch to %s", filepath)
	if err := c.OSCommand.CreateFileWithContent(filepath, patch); err != nil {
		return err
	}

	flagStr := ""
	for _, flag := range flags {
		flagStr += " --" + flag
	}

	return c.Run(c.NewCmdObj(fmt.Sprintf("git apply %s %s", flagStr, c.OSCommand.Quote(filepath))))
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (c *GitCommand) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	cmdObj := c.ShowFileDiffCmdObj(from, to, reverse, fileName, plain)
	return c.RunWithOutput(cmdObj)
}

func (c *GitCommand) ShowFileDiffCmdObj(from string, to string, reverse bool, fileName string, plain bool) oscommands.ICmdObj {
	colorArg := c.colorArg()
	contextSize := c.Config.GetUserConfig().Git.DiffContextSize
	if plain {
		colorArg = "never"
	}

	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	return c.NewCmdObj(fmt.Sprintf("git diff --submodule --no-ext-diff --unified=%d --no-renames --color=%s %s %s %s -- %s", contextSize, colorArg, from, to, reverseFlag, c.OSCommand.Quote(fileName)))
}

// CheckoutFile checks out the file for the given commit
func (c *GitCommand) CheckoutFile(commitSha, fileName string) error {
	return c.Run(c.NewCmdObj(fmt.Sprintf("git checkout %s -- %s", commitSha, c.OSCommand.Quote(fileName))))
}

// DiscardOldFileChanges discards changes to a file from an old commit
func (c *GitCommand) DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// check if file exists in previous commit (this command returns an error if the file doesn't exist)
	if err := c.Run(c.NewCmdObj("git cat-file -e HEAD^:" + c.OSCommand.Quote(fileName))); err != nil {
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
	return c.Run(c.NewCmdObj("git checkout -- ."))
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (c *GitCommand) RemoveTrackedFiles(name string) error {
	return c.Run(c.NewCmdObj("git rm -r --cached -- " + c.OSCommand.Quote(name)))
}

// RemoveUntrackedFiles runs `git clean -fd`
func (c *GitCommand) RemoveUntrackedFiles() error {
	return c.Run(c.NewCmdObj("git clean -fd"))
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

func (c *GitCommand) EditFileCmdStr(filename string, lineNumber int) (string, error) {
	editor := c.Config.GetUserConfig().OS.EditCommand

	if editor == "" {
		editor = c.GitConfig.Get("core.editor")
	}

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
		if err := c.OSCommand.Run(c.NewCmdObj("which vi")); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return "", errors.New("No editor defined in config file, $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
	}

	templateValues := map[string]string{
		"editor":   editor,
		"filename": c.OSCommand.Quote(filename),
		"line":     strconv.Itoa(lineNumber),
	}

	editCmdTemplate := c.Config.GetUserConfig().OS.EditCommandTemplate
	return utils.ResolvePlaceholderString(editCmdTemplate, templateValues), nil
}
