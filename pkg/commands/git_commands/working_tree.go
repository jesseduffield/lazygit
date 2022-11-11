package git_commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type WorkingTreeCommands struct {
	*GitCommon
	submodule  *SubmoduleCommands
	fileLoader *FileLoader
}

func NewWorkingTreeCommands(
	gitCommon *GitCommon,
	submodule *SubmoduleCommands,
	fileLoader *FileLoader,
) *WorkingTreeCommands {
	return &WorkingTreeCommands{
		GitCommon:  gitCommon,
		submodule:  submodule,
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
func (self *WorkingTreeCommands) StageFile(path string) error {
	return self.StageFiles([]string{path})
}

func (self *WorkingTreeCommands) StageFiles(paths []string) error {
	quotedPaths := slices.Map(paths, func(path string) string {
		return self.cmd.Quote(path)
	})
	return self.cmd.New(fmt.Sprintf("git add -- %s", strings.Join(quotedPaths, " "))).Run()
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

func (self *WorkingTreeCommands) BeforeAndAfterFileForRename(file *models.File) (*models.File, *models.File, error) {
	if !file.IsRename() {
		return nil, nil, errors.New("Expected renamed file")
	}

	// we've got a file that represents a rename from one file to another. Here we will refetch
	// all files, passing the --no-renames flag and then recursively call the function
	// again for the before file and after file.

	filesWithoutRenames := self.fileLoader.GetStatusFiles(GetStatusFileOptions{NoRenames: true})

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
func (self *WorkingTreeCommands) DiscardAllFileChanges(file *models.File) error {
	if file.IsRename() {
		beforeFile, afterFile, err := self.BeforeAndAfterFileForRename(file)
		if err != nil {
			return err
		}

		if err := self.DiscardAllFileChanges(beforeFile); err != nil {
			return err
		}

		if err := self.DiscardAllFileChanges(afterFile); err != nil {
			return err
		}

		return nil
	}

	quotedFileName := self.cmd.Quote(file.Name)

	if file.ShortStatus == "AA" {
		if err := self.cmd.New("git checkout --ours --  " + quotedFileName).Run(); err != nil {
			return err
		}
		if err := self.cmd.New("git add -- " + quotedFileName).Run(); err != nil {
			return err
		}
		return nil
	}

	if file.ShortStatus == "DU" {
		return self.cmd.New("git rm -- " + quotedFileName).Run()
	}

	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := self.cmd.New("git reset -- " + quotedFileName).Run(); err != nil {
			return err
		}
	}

	if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
		return nil
	}

	if file.Added {
		return self.os.RemoveFile(file.Name)
	}
	return self.DiscardUnstagedFileChanges(file)
}

type IFileNode interface {
	ForEachFile(cb func(*models.File) error) error
	GetFilePathsMatching(test func(*models.File) bool) []string
	GetPath() string
}

func (self *WorkingTreeCommands) DiscardAllDirChanges(node IFileNode) error {
	// this could be more efficient but we would need to handle all the edge cases
	return node.ForEachFile(self.DiscardAllFileChanges)
}

func (self *WorkingTreeCommands) DiscardUnstagedDirChanges(node IFileNode) error {
	if err := self.RemoveUntrackedDirFiles(node); err != nil {
		return err
	}

	quotedPath := self.cmd.Quote(node.GetPath())
	if err := self.cmd.New("git checkout -- " + quotedPath).Run(); err != nil {
		return err
	}

	return nil
}

func (self *WorkingTreeCommands) RemoveUntrackedDirFiles(node IFileNode) error {
	untrackedFilePaths := node.GetFilePathsMatching(
		func(file *models.File) bool { return !file.GetIsTracked() },
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
func (self *WorkingTreeCommands) DiscardUnstagedFileChanges(file *models.File) error {
	quotedFileName := self.cmd.Quote(file.Name)
	return self.cmd.New("git checkout -- " + quotedFileName).Run()
}

// Ignore adds a file to the gitignore for the repo
func (self *WorkingTreeCommands) Ignore(filename string) error {
	return self.os.AppendLineToFile(".gitignore", filename)
}

// Exclude adds a file to the .git/info/exclude for the repo
func (self *WorkingTreeCommands) Exclude(filename string) error {
	return self.os.AppendLineToFile(".git/info/exclude", filename)
}

// WorktreeFileDiff returns the diff of a file
func (self *WorkingTreeCommands) WorktreeFileDiff(file *models.File, plain bool, cached bool, ignoreWhitespace bool) string {
	// for now we assume an error means the file was deleted
	s, _ := self.WorktreeFileDiffCmdObj(file, plain, cached, ignoreWhitespace).RunWithOutput()
	return s
}

func (self *WorkingTreeCommands) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool, ignoreWhitespace bool) oscommands.ICmdObj {
	cachedArg := ""
	trackedArg := "--"
	colorArg := self.UserConfig.Git.Paging.ColorArg
	quotedPath := self.cmd.Quote(node.GetPath())
	quotedPrevPath := ""
	ignoreWhitespaceArg := ""
	contextSize := self.UserConfig.Git.DiffContextSize
	if cached {
		cachedArg = " --cached"
	}
	if !node.GetIsTracked() && !node.GetHasStagedChanges() && !cached && node.GetIsFile() {
		trackedArg = "--no-index -- /dev/null"
	}
	if plain {
		colorArg = "never"
	}
	if ignoreWhitespace {
		ignoreWhitespaceArg = " --ignore-all-space"
	}
	if prevPath := node.GetPreviousPath(); prevPath != "" {
		quotedPrevPath = " " + self.cmd.Quote(prevPath)
	}

	cmdStr := fmt.Sprintf("git diff --submodule --no-ext-diff --unified=%d --color=%s%s%s %s %s%s", contextSize, colorArg, ignoreWhitespaceArg, cachedArg, trackedArg, quotedPath, quotedPrevPath)

	return self.cmd.New(cmdStr).DontLog()
}

func (self *WorkingTreeCommands) ApplyPatch(patch string, flags ...string) error {
	filepath, err := self.SaveTemporaryPatch(patch)
	if err != nil {
		return err
	}

	return self.ApplyPatchFile(filepath, flags...)
}

func (self *WorkingTreeCommands) ApplyPatchFile(filepath string, flags ...string) error {
	flagStr := ""
	for _, flag := range flags {
		flagStr += " --" + flag
	}

	return self.cmd.New(fmt.Sprintf("git apply%s %s", flagStr, self.cmd.Quote(filepath))).Run()
}

func (self *WorkingTreeCommands) SaveTemporaryPatch(patch string) (string, error) {
	filepath := filepath.Join(self.os.GetTempDir(), utils.GetCurrentRepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	self.Log.Infof("saving temporary patch to %s", filepath)
	if err := self.os.CreateFileWithContent(filepath, patch); err != nil {
		return "", err
	}
	return filepath, nil
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (self *WorkingTreeCommands) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	return self.ShowFileDiffCmdObj(from, to, reverse, fileName, plain).RunWithOutput()
}

func (self *WorkingTreeCommands) ShowFileDiffCmdObj(from string, to string, reverse bool, fileName string, plain bool) oscommands.ICmdObj {
	colorArg := self.UserConfig.Git.Paging.ColorArg
	contextSize := self.UserConfig.Git.DiffContextSize
	if plain {
		colorArg = "never"
	}

	reverseFlag := ""
	if reverse {
		reverseFlag = " -R"
	}

	return self.cmd.
		New(
			fmt.Sprintf(
				"git diff --submodule --no-ext-diff --unified=%d --no-renames --color=%s%s%s%s -- %s",
				contextSize, colorArg, pad(from), pad(to), reverseFlag, self.cmd.Quote(fileName)),
		).
		DontLog()
}

// CheckoutFile checks out the file for the given commit
func (self *WorkingTreeCommands) CheckoutFile(commitSha, fileName string) error {
	return self.cmd.New(fmt.Sprintf("git checkout %s -- %s", commitSha, self.cmd.Quote(fileName))).Run()
}

// DiscardAnyUnstagedFileChanges discards any unstages file changes via `git checkout -- .`
func (self *WorkingTreeCommands) DiscardAnyUnstagedFileChanges() error {
	return self.cmd.New("git checkout -- .").Run()
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (self *WorkingTreeCommands) RemoveTrackedFiles(name string) error {
	return self.cmd.New("git rm -r --cached -- " + self.cmd.Quote(name)).Run()
}

// RemoveUntrackedFiles runs `git clean -fd`
func (self *WorkingTreeCommands) RemoveUntrackedFiles() error {
	return self.cmd.New("git clean -fd").Run()
}

// ResetAndClean removes all unstaged changes and removes all untracked files
func (self *WorkingTreeCommands) ResetAndClean() error {
	submoduleConfigs, err := self.submodule.GetConfigs()
	if err != nil {
		return err
	}

	if len(submoduleConfigs) > 0 {
		if err := self.submodule.ResetSubmodules(submoduleConfigs); err != nil {
			return err
		}
	}

	if err := self.ResetHard("HEAD"); err != nil {
		return err
	}

	return self.RemoveUntrackedFiles()
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

// so that we don't have unnecessary space in our commands we use this helper function to prepend spaces to args so that in the format string we can go '%s%s%s' and if any args are missing we won't have gaps.
func pad(str string) string {
	if str == "" {
		return ""
	}

	return " " + str
}
