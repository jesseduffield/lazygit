package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type WorkingTreeCommands struct {
	*common.Common

	cmd        oscommands.ICmdObjBuilder
	os         WorkingTreeOSCommand
	submodule  *SubmoduleCommands
	fileLoader *loaders.FileLoader
}

type WorkingTreeOSCommand interface {
	RemoveFile(string) error
	CreateFileWithContent(string, string) error
	AppendLineToFile(string, string) error
}

func NewWorkingTreeCommands(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
	submoduleCommands *SubmoduleCommands,
	osCommand WorkingTreeOSCommand,
	fileLoader *loaders.FileLoader,
) *WorkingTreeCommands {
	return &WorkingTreeCommands{
		Common:     common,
		cmd:        cmd,
		os:         osCommand,
		submodule:  submoduleCommands,
		fileLoader: fileLoader,
	}
}

func (self *WorkingTreeCommands) OpenMergeToolCmdObj() oscommands.ICmdObj {
	return self.cmd.New("git mergetool")
}

func (self *WorkingTreeCommands) OpenMergeTool() error {
	return self.OpenMergeToolCmdObj().Run()
}

// StageFile stages a file
func (self *WorkingTreeCommands) StageFile(fileName string) error {
	return self.cmd.New("git add -- " + self.cmd.Quote(fileName)).Run()
}

// StageAll stages all files
func (self *WorkingTreeCommands) StageAll() error {
	return self.cmd.New("git add -A").Run()
}

// UnstageAll unstages all files
func (self *WorkingTreeCommands) UnstageAll() error {
	return self.cmd.New("git reset").Run()
}

// UnStageFile unstages a file
// we accept an array of filenames for the cases where a file has been renamed i.e.
// we accept the current name and the previous name
func (self *WorkingTreeCommands) UnStageFile(fileNames []string, reset bool) error {
	command := "git rm --cached --force -- %s"
	if reset {
		command = "git reset HEAD -- %s"
	}

	for _, name := range fileNames {
		err := self.cmd.New(fmt.Sprintf(command, self.cmd.Quote(name))).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *WorkingTreeCommands) BeforeAndAfterFileForRename(file *models.File) (*models.File, *models.File, error) {
	if !file.IsRename() {
		return nil, nil, errors.New("Expected renamed file")
	}

	// we've got a file that represents a rename from one file to another. Here we will refetch
	// all files, passing the --no-renames flag and then recursively call the function
	// again for the before file and after file.

	filesWithoutRenames := c.fileLoader.GetStatusFiles(loaders.GetStatusFileOptions{NoRenames: true})

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
func (c *WorkingTreeCommands) DiscardAllFileChanges(file *models.File) error {
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

	quotedFileName := c.cmd.Quote(file.Name)

	if file.ShortStatus == "AA" {
		if err := c.cmd.New("git checkout --ours --  " + quotedFileName).Run(); err != nil {
			return err
		}
		if err := c.cmd.New("git add -- " + quotedFileName).Run(); err != nil {
			return err
		}
		return nil
	}

	if file.ShortStatus == "DU" {
		return c.cmd.New("git rm -- " + quotedFileName).Run()
	}

	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := c.cmd.New("git reset -- " + quotedFileName).Run(); err != nil {
			return err
		}
	}

	if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
		return nil
	}

	if file.Added {
		return c.os.RemoveFile(file.Name)
	}
	return c.DiscardUnstagedFileChanges(file)
}

func (c *WorkingTreeCommands) DiscardAllDirChanges(node *filetree.FileNode) error {
	// this could be more efficient but we would need to handle all the edge cases
	return node.ForEachFile(c.DiscardAllFileChanges)
}

func (c *WorkingTreeCommands) DiscardUnstagedDirChanges(node *filetree.FileNode) error {
	if err := c.RemoveUntrackedDirFiles(node); err != nil {
		return err
	}

	quotedPath := c.cmd.Quote(node.GetPath())
	if err := c.cmd.New("git checkout -- " + quotedPath).Run(); err != nil {
		return err
	}

	return nil
}

func (c *WorkingTreeCommands) RemoveUntrackedDirFiles(node *filetree.FileNode) error {
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
func (c *WorkingTreeCommands) DiscardUnstagedFileChanges(file *models.File) error {
	quotedFileName := c.cmd.Quote(file.Name)
	return c.cmd.New("git checkout -- " + quotedFileName).Run()
}

// Ignore adds a file to the gitignore for the repo
func (c *WorkingTreeCommands) Ignore(filename string) error {
	return c.os.AppendLineToFile(".gitignore", filename)
}

// WorktreeFileDiff returns the diff of a file
func (c *WorkingTreeCommands) WorktreeFileDiff(file *models.File, plain bool, cached bool, ignoreWhitespace bool) string {
	// for now we assume an error means the file was deleted
	s, _ := c.WorktreeFileDiffCmdObj(file, plain, cached, ignoreWhitespace).RunWithOutput()
	return s
}

func (c *WorkingTreeCommands) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool, ignoreWhitespace bool) oscommands.ICmdObj {
	cachedArg := ""
	trackedArg := "--"
	colorArg := c.UserConfig.Git.Paging.ColorArg
	quotedPath := c.cmd.Quote(node.GetPath())
	ignoreWhitespaceArg := ""
	contextSize := c.UserConfig.Git.DiffContextSize
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

	return c.cmd.New(cmdStr).DontLog()
}

func (c *WorkingTreeCommands) ApplyPatch(patch string, flags ...string) error {
	filepath := filepath.Join(oscommands.GetTempDir(), utils.GetCurrentRepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	c.Log.Infof("saving temporary patch to %s", filepath)
	if err := c.os.CreateFileWithContent(filepath, patch); err != nil {
		return err
	}

	flagStr := ""
	for _, flag := range flags {
		flagStr += " --" + flag
	}

	return c.cmd.New(fmt.Sprintf("git apply%s %s", flagStr, c.cmd.Quote(filepath))).Run()
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (c *WorkingTreeCommands) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	return c.ShowFileDiffCmdObj(from, to, reverse, fileName, plain).RunWithOutput()
}

func (c *WorkingTreeCommands) ShowFileDiffCmdObj(from string, to string, reverse bool, fileName string, plain bool) oscommands.ICmdObj {
	colorArg := c.UserConfig.Git.Paging.ColorArg
	contextSize := c.UserConfig.Git.DiffContextSize
	if plain {
		colorArg = "never"
	}

	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	return c.cmd.
		New(
			fmt.Sprintf(
				"git diff --submodule --no-ext-diff --unified=%d --no-renames --color=%s %s %s %s -- %s",
				contextSize, colorArg, from, to, reverseFlag, c.cmd.Quote(fileName)),
		).
		DontLog()
}

// CheckoutFile checks out the file for the given commit
func (c *WorkingTreeCommands) CheckoutFile(commitSha, fileName string) error {
	return c.cmd.New(fmt.Sprintf("git checkout %s -- %s", commitSha, c.cmd.Quote(fileName))).Run()
}

// DiscardAnyUnstagedFileChanges discards any unstages file changes via `git checkout -- .`
func (c *WorkingTreeCommands) DiscardAnyUnstagedFileChanges() error {
	return c.cmd.New("git checkout -- .").Run()
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (c *WorkingTreeCommands) RemoveTrackedFiles(name string) error {
	return c.cmd.New("git rm -r --cached -- " + c.cmd.Quote(name)).Run()
}

// RemoveUntrackedFiles runs `git clean -fd`
func (c *WorkingTreeCommands) RemoveUntrackedFiles() error {
	return c.cmd.New("git clean -fd").Run()
}

// ResetAndClean removes all unstaged changes and removes all untracked files
func (c *WorkingTreeCommands) ResetAndClean() error {
	submoduleConfigs, err := c.submodule.GetConfigs()
	if err != nil {
		return err
	}

	if len(submoduleConfigs) > 0 {
		if err := c.submodule.ResetSubmodules(submoduleConfigs); err != nil {
			return err
		}
	}

	if err := c.ResetHard("HEAD"); err != nil {
		return err
	}

	return c.RemoveUntrackedFiles()
}

// ResetHardHead runs `git reset --hard`
func (self *WorkingTreeCommands) ResetHard(ref string) error {
	return self.cmd.New("git reset --hard " + self.cmd.Quote(ref)).Run()
}

// ResetSoft runs `git reset --soft HEAD`
func (self *WorkingTreeCommands) ResetSoft(ref string) error {
	return self.cmd.New("git reset --soft " + self.cmd.Quote(ref)).Run()
}

func (self *WorkingTreeCommands) ResetMixed(ref string) error {
	return self.cmd.New("git reset --mixed " + self.cmd.Quote(ref)).Run()
}
