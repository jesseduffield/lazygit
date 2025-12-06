package git_commands

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

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

func (self *WorkingTreeCommands) OpenMergeToolCmdObj() *oscommands.CmdObj {
	return self.cmd.New(NewGitCmd("mergetool").ToArgv())
}

// StageFile stages a file
func (self *WorkingTreeCommands) StageFile(path string) error {
	return self.StageFiles([]string{path}, nil)
}

func (self *WorkingTreeCommands) StageFiles(paths []string, extraArgs []string) error {
	cmdArgs := NewGitCmd("add").
		Arg(extraArgs...).
		Arg("--").
		Arg(paths...).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// StageAll stages all files
func (self *WorkingTreeCommands) StageAll(onlyTrackedFiles bool) error {
	cmdArgs := NewGitCmd("add").
		ArgIfElse(onlyTrackedFiles, "-u", "-A").
		ToArgv()

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
	}
	return self.UnstageUntrackedFiles(paths)
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
		if f.Path == file.PreviousPath {
			beforeFile = f
		}

		if f.Path == file.Path {
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

// DiscardAllFilesChanges discards changes for multiple files in batch
func (self *WorkingTreeCommands) DiscardAllFilesChanges(files []*models.File) error {
	// Group files by their discard strategy
	var (
		aaStatusFiles      []*models.File
		duStatusFiles      []*models.File
		filesToReset       []*models.File
		addedFilesToRemove []*models.File
		filesToCheckout    []*models.File
	)

	// Helper function to categorize a file into the appropriate group
	categorizeFile := func(file *models.File) {
		if file.ShortStatus == "AA" {
			aaStatusFiles = append(aaStatusFiles, file)
			return
		}

		if file.ShortStatus == "DU" {
			duStatusFiles = append(duStatusFiles, file)
			return
		}

		// Track which files need to be reset first
		needsReset := file.HasStagedChanges || file.HasMergeConflicts
		if needsReset {
			filesToReset = append(filesToReset, file)
		}

		if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
		} else if file.Added {
			addedFilesToRemove = append(addedFilesToRemove, file)
		} else {
			filesToCheckout = append(filesToCheckout, file)
		}
	}

	for _, file := range files {
		if file.IsRename() {
			// Get the before and after files for the rename and add them to the appropriate groups
			beforeFile, afterFile, err := self.BeforeAndAfterFileForRename(file)
			if err != nil {
				return err
			}
			categorizeFile(beforeFile)
			categorizeFile(afterFile)
			continue
		}

		categorizeFile(file)
	}

	// Batch reset files that need resetting
	if len(filesToReset) > 0 {
		paths := make([]string, len(filesToReset))
		for i, file := range filesToReset {
			paths[i] = file.Path
		}
		if err := self.cmd.New(
			NewGitCmd("reset").Arg("--").Arg(paths...).ToArgv(),
		).Run(); err != nil {
			return err
		}
	}

	// Batch remove DU status files
	if len(duStatusFiles) > 0 {
		paths := make([]string, len(duStatusFiles))
		for i, file := range duStatusFiles {
			paths[i] = file.Path
		}
		if err := self.cmd.New(
			NewGitCmd("rm").Arg("--").Arg(paths...).ToArgv(),
		).Run(); err != nil {
			return err
		}
	}

	// Batch checkout --ours for AA status files
	if len(aaStatusFiles) > 0 {
		paths := make([]string, len(aaStatusFiles))
		for i, file := range aaStatusFiles {
			paths[i] = file.Path
		}
		if err := self.cmd.New(
			NewGitCmd("checkout").Arg("--ours", "--").Arg(paths...).ToArgv(),
		).Run(); err != nil {
			return err
		}
		// Stage them after checkout
		if err := self.cmd.New(
			NewGitCmd("add").Arg("--").Arg(paths...).ToArgv(),
		).Run(); err != nil {
			return err
		}
	}

	// Remove added files from filesystem
	for _, file := range addedFilesToRemove {
		if err := self.os.RemoveFile(file.Path); err != nil {
			return err
		}
	}

	// Batch checkout other files
	if len(filesToCheckout) > 0 {
		paths := make([]string, len(filesToCheckout))
		for i, file := range filesToCheckout {
			paths[i] = file.Path
		}
		if err := self.cmd.New(
			NewGitCmd("checkout").Arg("--").Arg(paths...).ToArgv(),
		).Run(); err != nil {
			return err
		}
	}

	return nil
}

type IFileNode interface {
	ForEachFile(cb func(*models.File) error) error
	GetFilePathsMatching(test func(*models.File) bool) []string
	GetPath() string
	// Returns file if the node is not a directory, otherwise returns nil
	GetFile() *models.File
}

// DiscardUnstagedFilesChanges discards unstaged changes for multiple files in batch
func (self *WorkingTreeCommands) DiscardUnstagedFilesChanges(files []*models.File) error {
	var (
		addedFilesToRemove     []*models.File
		trackedFilesToCheckout []*models.File
	)

	for _, file := range files {
		// Only remove files that are added but not staged
		if file.Added && !file.HasStagedChanges {
			addedFilesToRemove = append(addedFilesToRemove, file)
		} else {
			// Checkout tracked files to discard unstaged changes
			trackedFilesToCheckout = append(trackedFilesToCheckout, file)
		}
	}

	// Remove added files from filesystem
	for _, file := range addedFilesToRemove {
		if err := self.os.RemoveFile(file.Path); err != nil {
			return err
		}
	}

	// Batch checkout tracked files
	if len(trackedFilesToCheckout) > 0 {
		paths := make([]string, len(trackedFilesToCheckout))
		for i, file := range trackedFilesToCheckout {
			paths[i] = file.Path
		}
		cmdArgs := NewGitCmd("checkout").Arg("--").Arg(paths...).ToArgv()
		if err := self.cmd.New(cmdArgs).Run(); err != nil {
			return err
		}
	}

	return nil
}

// Escapes special characters in a filename for gitignore and exclude files
func escapeFilename(filename string) string {
	re := regexp.MustCompile(`^[!#]|[\[\]*]`)
	return re.ReplaceAllString(filename, `\${0}`)
}

// Ignore adds a file to the gitignore for the repo
func (self *WorkingTreeCommands) Ignore(filename string) error {
	return self.os.AppendLineToFile(".gitignore", escapeFilename(filename))
}

// Exclude adds a file to the .git/info/exclude for the repo
func (self *WorkingTreeCommands) Exclude(filename string) error {
	excludeFile := filepath.Join(self.repoPaths.repoGitDirPath, "info", "exclude")
	return self.os.AppendLineToFile(excludeFile, escapeFilename(filename))
}

// WorktreeFileDiff returns the diff of a file
func (self *WorkingTreeCommands) WorktreeFileDiff(file *models.File, plain bool, cached bool) string {
	// for now we assume an error means the file was deleted
	s, _ := self.WorktreeFileDiffCmdObj(file, plain, cached).RunWithOutput()
	return s
}

func (self *WorkingTreeCommands) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool) *oscommands.CmdObj {
	colorArg := self.pagerConfig.GetColorArg()
	if plain {
		colorArg = "never"
	}

	contextSize := self.UserConfig().Git.DiffContextSize
	prevPath := node.GetPreviousPath()
	noIndex := !node.GetIsTracked() && !node.GetHasStagedChanges() && !cached && node.GetIsFile()
	extDiffCmd := self.pagerConfig.GetExternalDiffCommand()
	useExtDiff := extDiffCmd != "" && !plain
	useExtDiffGitConfig := self.pagerConfig.GetUseExternalDiffGitConfig() && !plain

	cmdArgs := NewGitCmd("diff").
		ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
		ArgIfElse(useExtDiff || useExtDiffGitConfig, "--ext-diff", "--no-ext-diff").
		Arg("--submodule").
		Arg(fmt.Sprintf("--unified=%d", contextSize)).
		Arg(fmt.Sprintf("--color=%s", colorArg)).
		ArgIf(!plain && self.UserConfig().Git.IgnoreWhitespaceInDiffView, "--ignore-all-space").
		Arg(fmt.Sprintf("--find-renames=%d%%", self.UserConfig().Git.RenameSimilarityThreshold)).
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

func (self *WorkingTreeCommands) ShowFileDiffCmdObj(from string, to string, reverse bool, fileName string, plain bool) *oscommands.CmdObj {
	contextSize := self.UserConfig().Git.DiffContextSize

	colorArg := self.pagerConfig.GetColorArg()
	if plain {
		colorArg = "never"
	}

	extDiffCmd := self.pagerConfig.GetExternalDiffCommand()
	useExtDiff := extDiffCmd != "" && !plain
	useExtDiffGitConfig := self.pagerConfig.GetUseExternalDiffGitConfig() && !plain

	cmdArgs := NewGitCmd("diff").
		Config("diff.noprefix=false").
		ConfigIf(useExtDiff, "diff.external="+extDiffCmd).
		ArgIfElse(useExtDiff || useExtDiffGitConfig, "--ext-diff", "--no-ext-diff").
		Arg("--submodule").
		Arg(fmt.Sprintf("--unified=%d", contextSize)).
		Arg("--no-renames").
		Arg(fmt.Sprintf("--color=%s", colorArg)).
		Arg(from).
		Arg(to).
		ArgIf(reverse, "-R").
		ArgIf(!plain && self.UserConfig().Git.IgnoreWhitespaceInDiffView, "--ignore-all-space").
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

func (self *WorkingTreeCommands) RemoveConflictedFile(name string) error {
	cmdArgs := NewGitCmd("rm").Arg("--", name).
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

func (self *WorkingTreeCommands) ShowFileAtStage(path string, stage int) (string, error) {
	cmdArgs := NewGitCmd("show").
		Arg(fmt.Sprintf(":%d:%s", stage, path)).
		ToArgv()

	return self.cmd.New(cmdArgs).RunWithOutput()
}

func (self *WorkingTreeCommands) ObjectIDAtStage(path string, stage int) (string, error) {
	cmdArgs := NewGitCmd("rev-parse").
		Arg(fmt.Sprintf(":%d:%s", stage, path)).
		ToArgv()

	output, err := self.cmd.New(cmdArgs).RunWithOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

func (self *WorkingTreeCommands) MergeFileForFiles(strategy string, oursFilepath string, baseFilepath string, theirsFilepath string) (string, error) {
	cmdArgs := NewGitCmd("merge-file").
		Arg(strategy).
		Arg("--stdout").
		Arg(oursFilepath, baseFilepath, theirsFilepath).
		ToArgv()

	return self.cmd.New(cmdArgs).RunWithOutput()
}

// OIDs mode (Git 2.43+)
func (self *WorkingTreeCommands) MergeFileForObjectIDs(strategy string, oursID string, baseID string, theirsID string) (string, error) {
	cmdArgs := NewGitCmd("merge-file").
		Arg(strategy).
		Arg("--stdout").
		Arg("--object-id").
		Arg(oursID, baseID, theirsID).
		ToArgv()

	return self.cmd.New(cmdArgs).RunWithOutput()
}

// Returns all tracked files in the repo (not in the working tree). The returned entries are
// relative paths to the repo root, using '/' as the path separator on all platforms.
// Does not really belong in WorkingTreeCommands, but it's close enough, and we don't seem to have a
// better place for it right now.
func (self *WorkingTreeCommands) AllRepoFiles() ([]string, error) {
	cmdArgs := NewGitCmd("ls-files").Arg("-z").ToArgv()
	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(strings.TrimRight(output, "\x00"), "\x00"), nil
}
