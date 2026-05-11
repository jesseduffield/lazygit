package git_commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
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
			NewGitCmd("checkout").Arg("--ours", "--", file.Path).ToArgv(),
		).Run(); err != nil {
			return err
		}

		if err := self.cmd.New(
			NewGitCmd("add").Arg("--", file.Path).ToArgv(),
		).Run(); err != nil {
			return err
		}
		return nil
	}

	if file.ShortStatus == "DU" {
		return self.cmd.New(
			NewGitCmd("rm").Arg("--", file.Path).ToArgv(),
		).Run()
	}

	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := self.cmd.New(
			NewGitCmd("reset").Arg("--", file.Path).ToArgv(),
		).Run(); err != nil {
			return err
		}
	}

	if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
		return nil
	}

	if file.Added {
		return self.os.RemoveFile(file.Path)
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

func (self *WorkingTreeCommands) DiscardAllDirChanges(nodes []IFileNode) error {
	// Collect files into buckets so we can batch git calls where possible.
	var specialFiles []*models.File // renames, AA, DU — handled individually
	var filesToReset []string       // need `git reset` first (staged or conflicted)
	var filesToCheckout []string    // need `git checkout` (after optional reset)
	var filesToRemove []string      // added files to delete from disk

	for _, node := range nodes {
		_ = node.ForEachFile(func(file *models.File) error {
			// Renames and certain merge-conflict statuses need per-file logic.
			if file.IsRename() || file.ShortStatus == "AA" || file.ShortStatus == "DU" {
				specialFiles = append(specialFiles, file)
				return nil
			}

			if file.HasStagedChanges || file.HasMergeConflicts {
				filesToReset = append(filesToReset, file.Path)
				// DD and AU are done after the reset; no checkout or remove needed.
				if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
					return nil
				}
				if file.Added {
					filesToRemove = append(filesToRemove, file.Path)
				} else {
					filesToCheckout = append(filesToCheckout, file.Path)
				}
				return nil
			}

			// No staged changes below this point.
			if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
				return nil
			}

			if file.Added {
				filesToRemove = append(filesToRemove, file.Path)
				return nil
			}

			filesToCheckout = append(filesToCheckout, file.Path)
			return nil
		})
	}

	for _, file := range specialFiles {
		if err := self.DiscardAllFileChanges(file); err != nil {
			return err
		}
	}

	if err := runGitCmdOnPaths("reset", filesToReset, self.cmd); err != nil {
		return err
	}

	if err := self.removeFiles(filesToRemove, nodes); err != nil {
		return err
	}

	return runGitCmdOnPaths("checkout", filesToCheckout, self.cmd)
}

func (self *WorkingTreeCommands) DiscardUnstagedDirChanges(nodes []IFileNode) error {
	// Collect files into buckets so we can batch git calls where possible.
	// Use specific file paths rather than directory paths, so that an active
	// filter (e.g. from pressing `/`) only discards visible files.
	var filesToRemove []string   // purely untracked: remove from disk
	var filesToCheckout []string // tracked or staged: restore via checkout

	for _, node := range nodes {
		_ = node.ForEachFile(func(file *models.File) error {
			if !file.Tracked && !file.HasStagedChanges {
				filesToRemove = append(filesToRemove, file.Path)
			} else {
				// Include staged files: a file that is staged but also has
				// additional unstaged changes (AM status) needs checkout to
				// discard those changes.
				filesToCheckout = append(filesToCheckout, file.Path)
			}
			return nil
		})
	}

	if err := self.removeFiles(filesToRemove, nodes); err != nil {
		return err
	}

	return runGitCmdOnPaths("checkout", filesToCheckout, self.cmd)
}

// Removes the given files from disk, and also removes any directories that have become empty
// because of this.
func (self *WorkingTreeCommands) removeFiles(paths []string, selectedNodes []IFileNode) error {
	for _, path := range paths {
		if err := self.os.RemoveFile(path); err != nil {
			return err
		}
	}

	return self.removeEmptyDirs(paths, selectedDirPaths(selectedNodes))
}

// Removes empty directories left behind after deleting files, but only for directories that
// are at or below a selected directory node. It works bottom-up so that nested empty directories
// are also cleaned up. Directories that still have contents are skipped.
func (self *WorkingTreeCommands) removeEmptyDirs(removedFilePaths []string, selectedDirs []string) error {
	candidates := set.NewFromSlice(
		lo.FilterMap(removedFilePaths, func(filePath string, _ int) (string, bool) {
			dir := path.Dir(filePath)
			return dir, dir != "." && isUnderSelectedDir(dir, selectedDirs)
		}))

	for {
		var removed []string
		for _, dir := range candidates.ToSlice() {
			empty, err := self.os.IsDirEmpty(dir)
			if err != nil {
				return err
			}
			if empty {
				if err := self.os.RemoveDir(dir); err != nil {
					return err
				}
				removed = append(removed, dir)
			}
		}
		if len(removed) == 0 {
			break
		}
		for _, dir := range removed {
			candidates.Remove(dir)
			if parent := path.Dir(dir); parent != "." && isUnderSelectedDir(parent, selectedDirs) {
				candidates.Add(parent)
			}
		}
	}
	return nil
}

func isUnderSelectedDir(path string, selectedDirs []string) bool {
	isSubdir := func(parent, child string) bool {
		rel, err := filepath.Rel(parent, child)
		return err == nil && !strings.HasPrefix(rel, "..")
	}

	return lo.SomeBy(selectedDirs, func(selectedDir string) bool {
		return isSubdir(selectedDir, path)
	})
}

func selectedDirPaths(nodes []IFileNode) []string {
	return lo.FilterMap(nodes, func(node IFileNode, _ int) (string, bool) {
		return node.GetPath(), node.GetFile() == nil
	})
}

func (self *WorkingTreeCommands) RemoveUntrackedDirFiles(node IFileNode) error {
	untrackedFilePaths := node.GetFilePathsMatching(
		func(file *models.File) bool { return !file.GetIsTracked() && !file.GetHasStagedChanges() },
	)

	for _, path := range untrackedFilePaths {
		if err := self.os.RemoveFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (self *WorkingTreeCommands) DiscardUnstagedFileChanges(file *models.File) error {
	cmdArgs := NewGitCmd("checkout").Arg("--", file.Path).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// Escapes special characters in a filename for gitignore and exclude files, and prepends `/`
func escapeFilename(filename string) string {
	re := regexp.MustCompile(`^[!#]|[\[\]*]`)
	return "/" + re.ReplaceAllString(filename, `\${0}`)
}

// Ignore adds a file to the gitignore for the repo
func (self *WorkingTreeCommands) Ignore(filename string) error {
	return self.os.AppendLineToFile(".gitignore", escapeFilename(filename))
}

// Exclude adds a file to the .git/info/exclude for the repo
func (self *WorkingTreeCommands) Exclude(filename string) error {
	infoDir := filepath.Join(self.repoPaths.repoGitDirPath, "info")
	if err := os.MkdirAll(infoDir, 0o755); err != nil {
		return err
	}
	excludeFile := filepath.Join(infoDir, "exclude")
	return self.os.AppendLineToFile(excludeFile, escapeFilename(filename))
}

// WorktreeFileDiff returns the diff of a file
func (self *WorkingTreeCommands) WorktreeFileDiff(file *models.File, plain bool, cached bool) string {
	// for now we assume an error means the file was deleted
	s, _ := self.WorktreeFileDiffCmdObj(file, plain, cached, nil).RunWithOutput()
	return s
}

// WorktreeFileDiffCmdObj returns a command object for diffing a file or directory
// in the working tree. When pathOverrides is non-empty, those paths are used instead of
// the node's path (used to diff only filtered/visible files within a directory).
func (self *WorkingTreeCommands) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool, pathOverrides []string) *oscommands.CmdObj {
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

	paths := pathOverrides
	if len(paths) == 0 {
		paths = []string{node.GetPath()}
	}

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
		Arg(paths...).
		ArgIf(prevPath != "", prevPath).
		Dir(self.repoPaths.worktreePath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (self *WorkingTreeCommands) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	return self.ShowFileDiffCmdObj(from, to, reverse, []string{fileName}, plain).RunWithOutput()
}

func (self *WorkingTreeCommands) ShowFileDiffCmdObj(from string, to string, reverse bool, fileNames []string, plain bool) *oscommands.CmdObj {
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
		Arg(fileNames...).
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
