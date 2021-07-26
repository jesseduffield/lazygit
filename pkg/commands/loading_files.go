package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

const RENAME_SEPARATOR = " -> "

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

	statusStrings, err := c.GitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg})
	if err != nil {
		c.Log.Error(err)
	}
	files := []*models.File{}

	for _, statusString := range statusStrings {
		if strings.HasPrefix(statusString, "warning") {
			c.Log.Warningf("warning when calling git status: %s", statusString)
			continue
		}
		change := statusString[0:2]
		stagedChange := change[0:1]
		unstagedChange := statusString[1:2]
		name := statusString[3:]
		untracked := utils.IncludesString([]string{"??", "A ", "AM"}, change)
		hasNoStagedChanges := utils.IncludesString([]string{" ", "U", "?"}, stagedChange)
		hasMergeConflicts := utils.IncludesString([]string{"DD", "AA", "UU", "AU", "UA", "UD", "DU"}, change)
		hasInlineMergeConflicts := utils.IncludesString([]string{"UU", "AA"}, change)
		previousName := ""
		if strings.Contains(name, RENAME_SEPARATOR) {
			split := strings.Split(name, RENAME_SEPARATOR)
			name = split[1]
			previousName = split[0]
		}

		file := &models.File{
			Name:                    name,
			PreviousName:            previousName,
			DisplayString:           statusString,
			HasStagedChanges:        !hasNoStagedChanges,
			HasUnstagedChanges:      unstagedChange != " ",
			Tracked:                 !untracked,
			Deleted:                 unstagedChange == "D" || stagedChange == "D",
			Added:                   unstagedChange == "A" || untracked,
			HasMergeConflicts:       hasMergeConflicts,
			HasInlineMergeConflicts: hasInlineMergeConflicts,
			Type:                    c.OSCommand.FileType(name),
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

func (c *GitCommand) GitStatus(opts GitStatusOptions) ([]string, error) {
	noRenamesFlag := ""
	if opts.NoRenames {
		noRenamesFlag = "--no-renames"
	}

	statusLines, err := c.RunCommandWithOutput("git status %s --porcelain -z %s", opts.UntrackedFilesArg, noRenamesFlag)
	if err != nil {
		return []string{}, err
	}

	splitLines := strings.Split(statusLines, "\x00")
	response := []string{}

	for i := 0; i < len(splitLines); i++ {
		original := splitLines[i]
		if len(original) < 2 {
			continue
		} else if strings.HasPrefix(original, "R  ") {
			// if a line starts with 'R' then the next line is the original file.
			next := strings.TrimSpace(splitLines[i+1])
			original = "R  " + next + RENAME_SEPARATOR + strings.TrimPrefix(original, "R  ")
			i++
		}
		response = append(response, original)
	}

	return response, nil
}
