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
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// CatFile obtains the content of a file
func (c *Git) CatFile(fileName string) (string, error) {
	return c.GetOS().CatFile(fileName)
}

func (c *Git) OpenMergeToolCmdObj() ICmdObj {
	return BuildGitCmdObjFromStr("mergetool")
}

// StageFile stages a file
func (c *Git) StageFile(fileName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("add -- %s", c.GetOS().Quote(fileName)))
}

// StageAll stages all files
func (c *Git) StageAll() error {
	return c.RunGitCmdFromStr("add -A")
}

// UnstageAll unstages all files
func (c *Git) UnstageAll() error {
	return c.RunGitCmdFromStr("reset")
}

// UnStageFile unstages a file
// we accept an array of filenames for the cases where a file has been renamed i.e.
// we accept the current name and the previous name
func (c *Git) UnStageFile(fileNames []string, reset bool) error {
	cmdFormat := "rm --cached --force -- %s"
	if reset {
		cmdFormat = "reset HEAD -- %s"
	}

	for _, name := range fileNames {
		if err := c.RunGitCmdFromStr(fmt.Sprintf(cmdFormat, c.GetOS().Quote(name))); err != nil {
			return err
		}
	}
	return nil
}

func (c *Git) BeforeAndAfterFileForRename(file *models.File) (*models.File, *models.File, error) {
	if !file.IsRename() {
		return nil, nil, errors.New("Expected renamed file")
	}

	// we've got a file that represents a rename from one file to another. Here we will refetch
	// all files, passing the --no-renames flag and then recursively call the function
	// again for the before file and after file.

	filesWithoutRenames := c.GetStatusFiles(loaders.LoadStatusFilesOpts{NoRenames: true})
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
func (c *Git) DiscardAllFileChanges(file *models.File) error {
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

	quotedFileName := c.GetOS().Quote(file.Name)

	if file.ShortStatus == "AA" {
		if err := c.RunGitCmdFromStr(fmt.Sprintf("checkout --ours --  %s", quotedFileName)); err != nil {
			return err
		}
		if err := c.RunGitCmdFromStr(fmt.Sprintf("add %s", quotedFileName)); err != nil {
			return err
		}
		return nil
	}

	if file.ShortStatus == "DU" {
		return c.RunGitCmdFromStr(fmt.Sprintf("rm %s", quotedFileName))
	}

	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges || file.HasMergeConflicts {
		if err := c.RunGitCmdFromStr(fmt.Sprintf("reset -- %s", quotedFileName)); err != nil {
			return err
		}
	}

	if file.ShortStatus == "DD" || file.ShortStatus == "AU" {
		return nil
	}

	if file.Added {
		return c.GetOS().RemoveFile(file.Name)
	}
	return c.DiscardUnstagedFileChanges(file)
}

func (c *Git) DiscardAllDirChanges(node *filetree.FileNode) error {
	// this could be more efficient but we would need to handle all the edge cases
	return node.ForEachFile(c.DiscardAllFileChanges)
}

func (c *Git) DiscardUnstagedDirChanges(node *filetree.FileNode) error {
	if err := c.RemoveUntrackedDirFiles(node); err != nil {
		return err
	}

	quotedPath := c.GetOS().Quote(node.GetPath())
	if err := c.RunGitCmdFromStr(fmt.Sprintf("checkout -- %s", quotedPath)); err != nil {
		return err
	}

	return nil
}

func (c *Git) RemoveUntrackedDirFiles(node *filetree.FileNode) error {
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
func (c *Git) DiscardUnstagedFileChanges(file *models.File) error {
	quotedFileName := c.GetOS().Quote(file.Name)
	return c.RunGitCmdFromStr(fmt.Sprintf("checkout -- %s", quotedFileName))
}

// Ignore adds a file to the gitignore for the repo
func (c *Git) Ignore(filename string) error {
	return c.GetOS().AppendLineToFile(".gitignore", filename)
}

func (c *Git) ApplyPatch(patch string, flags ...string) error {
	filepath := filepath.Join(c.config.GetUserConfigDir(), utils.GetCurrentRepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	c.log.Infof("saving temporary patch to %s", filepath)
	if err := c.GetOS().CreateFileWithContent(filepath, patch); err != nil {
		return err
	}

	flagStr := ""
	for _, flag := range flags {
		flagStr += " --" + flag
	}

	return c.RunGitCmdFromStr(fmt.Sprintf("apply %s %s", flagStr, c.GetOS().Quote(filepath)))
}

// CheckoutFile checks out the file for the given commit
func (c *Git) CheckoutFile(commitSha, fileName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("checkout %s %s", commitSha, fileName))
}

// DiscardOldFileChanges discards changes to a file from an old commit
func (c *Git) DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// check if file exists in previous commit (this command returns an error if the file doesn't exist)
	if err := c.RunGitCmdFromStr(fmt.Sprintf("cat-file -e HEAD^:%s", fileName)); err != nil {
		if err := c.GetOS().Remove(fileName); err != nil {
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

	return c.ContinueRebase()
}

// DiscardAnyUnstagedFileChanges discards any unstages file changes via `git checkout -- .`
func (c *Git) DiscardAnyUnstagedFileChanges() error {
	return c.RunGitCmdFromStr("checkout -- .")
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (c *Git) RemoveTrackedFiles(name string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("rm -r --cached %s", name))
}

// RemoveUntrackedFiles runs `git clean -fd`
func (c *Git) RemoveUntrackedFiles() error {
	return c.RunGitCmdFromStr("clean -fd")
}

// ResetAndClean removes all unstaged changes and removes all untracked files
func (c *Git) ResetAndClean() error {
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

func (c *Git) EditFileCmdObj(filename string) (ICmdObj, error) {
	editor := c.config.GetUserConfig().OS.EditCommand

	if editor == "" {
		editor = c.GetConfigValue("core.editor")
	}

	if editor == "" {
		editor = c.GetOS().Getenv("GIT_EDITOR")
	}
	if editor == "" {
		editor = c.GetOS().Getenv("VISUAL")
	}
	if editor == "" {
		editor = c.GetOS().Getenv("EDITOR")
	}
	if editor == "" {
		if err := c.Run(oscommands.NewCmdObjFromStr("which vi")); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return nil, errors.New("No editor defined in config file, $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
	}

	cmdObj := c.BuildShellCmdObj(fmt.Sprintf("%s %s", editor, c.GetOS().Quote(filename)))

	return cmdObj, nil
}
