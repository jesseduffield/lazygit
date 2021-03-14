package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// GetStatusFiles git status files
type GetStatusFileOptions struct {
	NoRenames bool
}

func (c *GitCommand) GetStatusFiles(opts GetStatusFileOptions) []*models.File {
	// check if config wants us ignoring untracked files
	untrackedFilesSetting := c.GetConfigValue("status.showUntrackedFiles")

	if untrackedFilesSetting == "" {
		untrackedFilesSetting = "all"
	}
	untrackedFilesArg := fmt.Sprintf("--untracked-files=%s", untrackedFilesSetting)

	statusOutput, err := c.GitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg})
	if err != nil {
		c.Log.Error(err)
	}
	statusStrings := utils.SplitLines(statusOutput)
	files := []*models.File{}

	for _, statusString := range statusStrings {
		if strings.HasPrefix(statusString, "warning") {
			c.Log.Warningf("warning when calling git status: %s", statusString)
			continue
		}
		change := statusString[0:2]
		stagedChange := change[0:1]
		unstagedChange := statusString[1:2]
		filename := statusString[3:]
		untracked := utils.IncludesString([]string{"??", "A ", "AM"}, change)
		hasNoStagedChanges := utils.IncludesString([]string{" ", "U", "?"}, stagedChange)
		hasMergeConflicts := utils.IncludesString([]string{"DD", "AA", "UU", "AU", "UA", "UD", "DU"}, change)
		hasInlineMergeConflicts := utils.IncludesString([]string{"UU", "AA"}, change)

		file := &models.File{
			Name:                    filename,
			DisplayString:           statusString,
			HasStagedChanges:        !hasNoStagedChanges,
			HasUnstagedChanges:      unstagedChange != " ",
			Tracked:                 !untracked,
			Deleted:                 unstagedChange == "D" || stagedChange == "D",
			HasMergeConflicts:       hasMergeConflicts,
			HasInlineMergeConflicts: hasInlineMergeConflicts,
			Type:                    c.OSCommand.FileType(filename),
			ShortStatus:             change,
		}
		files = append(files, file)
	}
	return files
}

// GitStatus returns the plaintext short status of the repo
type GitStatusOptions struct {
	NoRenames         bool
	UntrackedFilesArg string
}

func (c *GitCommand) GitStatus(opts GitStatusOptions) (string, error) {
	noRenamesFlag := ""
	if opts.NoRenames {
		noRenamesFlag = "--no-renames"
	}

	statusLines, err := c.OSCommand.RunCommandWithOutput("git status %s --porcelain -z %s", opts.UntrackedFilesArg, noRenamesFlag)
	if err != nil {
		return "", err
	}

	splitLines := strings.Split(statusLines, "\x00")
	// if a line starts with 'R' then the next line is the original file.
	for i := 0; i < len(splitLines)-1; i++ {
		original := splitLines[i]
		if strings.HasPrefix(original, "R  ") {
			next := splitLines[i+1]
			updated := "R  " + next + models.RENAME_SEPARATOR + strings.TrimPrefix(original, "R  ")
			splitLines[i] = updated
			splitLines = append(splitLines[0:i+1], splitLines[i+2:]...)
		}
	}

	return strings.Join(splitLines, "\n"), nil
}

// MergeStatusFiles merge status files
func (c *GitCommand) MergeStatusFiles(oldFiles, newFiles []*models.File, selectedFile *models.File) []*models.File {
	if len(oldFiles) == 0 {
		return newFiles
	}

	appendedIndexes := []int{}

	// retain position of files we already could see
	result := []*models.File{}
	for _, oldFile := range oldFiles {
		for newIndex, newFile := range newFiles {
			if utils.IncludesInt(appendedIndexes, newIndex) {
				continue
			}
			// if we just staged B and in doing so created 'A -> B' and we are currently have oldFile: A and newFile: 'A -> B', we want to wait until we come across B so the our cursor isn't jumping anywhere
			waitForMatchingFile := selectedFile != nil && newFile.IsRename() && !selectedFile.IsRename() && newFile.Matches(selectedFile) && !oldFile.Matches(selectedFile)

			if oldFile.Matches(newFile) && !waitForMatchingFile {
				result = append(result, newFile)
				appendedIndexes = append(appendedIndexes, newIndex)
			}
		}
	}

	// append any new files to the end
	for index, newFile := range newFiles {
		if !utils.IncludesInt(appendedIndexes, index) {
			result = append(result, newFile)
		}
	}

	return result
}
