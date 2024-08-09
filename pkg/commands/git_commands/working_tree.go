package git_commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
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
	return self.cmd.New(NewGitCmd("mergetool").ToArgv())
}

// StageFile stages a file
func (self *WorkingTreeCommands) StageFile(path string) error {
	return self.StageFiles([]string{path})
}

func (self *WorkingTreeCommands) StageFiles(paths []string) error {
	cmdArgs := NewGitCmd("add").Arg("--").Arg(paths...).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// StageAll stages all files
func (self *WorkingTreeCommands) StageAll() error {
	cmdArgs := NewGitCmd("add").Arg("-A").ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// UnstageAll unstages all files
func (self *WorkingTreeCommands) UnstageAll() error {
	return self.cmd.New(NewGitCmd("reset").ToArgv()).Run()
}

// UnStageFile unstages a file
// we accept an array of filenames for the cases where a file has been renamed i.e.
// we accept the current name and the previous name
func (self *WorkingTreeCommands) UnStageFile(paths []string, tracked bool) error {
	if tracked {
		return self.UnstageTrackedFiles(paths)
	} else {
		return self.UnstageUntrackedFiles(paths)
	}
}

func (self *WorkingTreeCommands) UnstageTrackedFiles(paths []string) error {
	return self.cmd.New(NewGitCmd("reset").Arg("HEAD", "--").Arg(paths...).ToArgv()).Run()
}

func (self *WorkingTreeCommands) UnstageUntrackedFiles(paths []string) error {
	return self.cmd.New(NewGitCmd("rm").Arg("--cached", "--force", "--").Arg(paths...).ToArgv()).Run()
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

	if file.ShortStatus == "AA" {
		if err := self.cmd.New(
			NewGitCmd("checkout").Arg("--ours", "--", file.Name).ToArgv(),
		).Run(); err != nil {
			return err
		}

		if err := self.cmd.New(
			NewGitCmd("add").Arg("--", file.Name).ToArgv(),
		).Run(); err != nil {
			return err
		}
		return nil
	}

	if file.ShortStatus == "DU" {
		return self.cmd.New(
			NewGitCmd("rm").Arg("--", file.Name).ToArgv(),
		).Run()
	}

	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := self.cmd.New(
			NewGitCmd("reset").Arg("--", file.Name).ToArgv(),
		).Run(); err != nil {
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
	// Returns file if the node is not a directory, otherwise returns nil
	GetFile() *models.File
}

func (self *WorkingTreeCommands) DiscardAllDirChanges(node IFileNode) error {
	// this could be more efficient but we would need to handle all the edge cases
	return node.ForEachFile(self.DiscardAllFileChanges)
}

func (self *WorkingTreeCommands) DiscardUnstagedDirChanges(node IFileNode) error {
	file := node.GetFile()
	if file == nil {
		if err := self.RemoveUntrackedDirFiles(node); err != nil {
			return err
		}

		cmdArgs := NewGitCmd("checkout").Arg("--", node.GetPath()).ToArgv()
		if err := self.cmd.New(cmdArgs).Run(); err != nil {
			return err
		}
	} else {
		if file.Added && !file.HasStagedChanges {
			return self.os.RemoveFile(file.Name)
		}

		if err := self.DiscardUnstagedFileChanges(file); err != nil {
			return err
		}
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

func (self *WorkingTreeCommands) DiscardUnstagedFileChanges(file *models.File) error {
	cmdArgs := NewGitCmd("checkout").Arg("--", file.Name).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// Ignore adds a file to the gitignore for the repo
func (self *WorkingTreeCommands) Ignore(filename string) error {
	return self.os.AppendLineToFile(".gitignore", filename)
}

// Exclude adds a file to the .git/info/exclude for the repo
func (self *WorkingTreeCommands) Exclude(filename string) error {
	excludeFile := filepath.Join(self.repoPaths.repoGitDirPath, "info", "exclude")
	return self.os.AppendLineToFile(excludeFile, filename)
}

// WorktreeFileDiff returns the diff of a file
func (self *WorkingTreeCommands) WorktreeFileDiff(file *models.File, plain bool, cached bool) string {
	// for now we assume an error means the file was deleted
	s, _ := self.WorktreeFileDiffCmdObj(file, plain, cached).RunWithOutput()
	return s
}

func (self *WorkingTreeCommands) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool) oscommands.ICmdObj {
	colorArg := self.UserConfig().Git.Paging.ColorArg
	if plain {
		colorArg = "never"
	}

	contextSize := self.AppState.DiffContextSize
	prevPath := node.GetPreviousPath()
	noIndex := !node.GetIsTracked() && !node.GetHasStagedChanges() && !cached && node.GetIsFile()
	extDiffCmd := self.UserConfig().Git.Paging.ExternalDiffCommand
	useExtDiff := extDiffCmd != "" && !plain

	cmdArgs := NewGitCmd("diff").
		ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
		ArgIfElse(useExtDiff, "--ext-diff", "--no-ext-diff").
		Arg("--submodule").
		Arg(fmt.Sprintf("--unified=%d", contextSize)).
		Arg(fmt.Sprintf("--color=%s", colorArg)).
		ArgIf(!plain && self.AppState.IgnoreWhitespaceInDiffView, "--ignore-all-space").
		Arg(fmt.Sprintf("--find-renames=%d%%", self.AppState.RenameSimilarityThreshold)).
		ArgIf(cached, "--cached").
		ArgIf(noIndex, "--no-index").
		Arg("--").
		ArgIf(noIndex, "/dev/null").
		Arg(node.GetPath()).
		ArgIf(prevPath != "", prevPath).
		Dir(self.repoPaths.worktreePath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (self *WorkingTreeCommands) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	return self.ShowFileDiffCmdObj(from, to, reverse, fileName, plain).RunWithOutput()
}

func (self *WorkingTreeCommands) ShowFileDiffCmdObj(from string, to string, reverse bool, fileName string, plain bool) oscommands.ICmdObj {
	contextSize := self.AppState.DiffContextSize

	colorArg := self.UserConfig().Git.Paging.ColorArg
	if plain {
		colorArg = "never"
	}

	extDiffCmd := self.UserConfig().Git.Paging.ExternalDiffCommand
	useExtDiff := extDiffCmd != "" && !plain

	cmdArgs := NewGitCmd("diff").
		Config("diff.noprefix=false").
		ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
		ArgIfElse(useExtDiff, "--ext-diff", "--no-ext-diff").
		Arg("--submodule").
		Arg(fmt.Sprintf("--unified=%d", contextSize)).
		Arg("--no-renames").
		Arg(fmt.Sprintf("--color=%s", colorArg)).
		Arg(from).
		Arg(to).
		ArgIf(reverse, "-R").
		ArgIf(!plain && self.AppState.IgnoreWhitespaceInDiffView, "--ignore-all-space").
		Arg("--").
		Arg(fileName).
		Dir(self.repoPaths.worktreePath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

// CheckoutFile checks out the file for the given commit
func (self *WorkingTreeCommands) CheckoutFile(commitHash, fileName string) error {
	cmdArgs := NewGitCmd("checkout").Arg(commitHash, "--", fileName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// DiscardAnyUnstagedFileChanges discards any unstaged file changes via `git checkout -- .`
func (self *WorkingTreeCommands) DiscardAnyUnstagedFileChanges() error {
	cmdArgs := NewGitCmd("checkout").Arg("--", ".").
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (self *WorkingTreeCommands) RemoveTrackedFiles(name string) error {
	cmdArgs := NewGitCmd("rm").Arg("-r", "--cached", "--", name).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// RemoveUntrackedFiles runs `git clean -fd`
func (self *WorkingTreeCommands) RemoveUntrackedFiles() error {
	cmdArgs := NewGitCmd("clean").Arg("-fd").ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// ResetAndClean removes all unstaged changes and removes all untracked files
func (self *WorkingTreeCommands) ResetAndClean() error {
	submoduleConfigs, err := self.submodule.GetConfigs(nil)
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

// ResetHard runs `git reset --hard`
func (self *WorkingTreeCommands) ResetHard(ref string) error {
	cmdArgs := NewGitCmd("reset").Arg("--hard", ref).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// ResetSoft runs `git reset --soft HEAD`
func (self *WorkingTreeCommands) ResetSoft(ref string) error {
	cmdArgs := NewGitCmd("reset").Arg("--soft", ref).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *WorkingTreeCommands) ResetMixed(ref string) error {
	cmdArgs := NewGitCmd("reset").Arg("--mixed", ref).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}
