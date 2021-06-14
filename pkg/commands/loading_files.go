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

	statusOutput, err := c.GitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg})
	if err != nil {
		c.log.Error(err)
	}
	statusStrings := utils.SplitLines(statusOutput)
	files := []*models.File{}

	for _, statusString := range statusStrings {
		if strings.HasPrefix(statusString, "warning") {
			c.log.Warningf("warning when calling git status: %s", statusString)
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
			Type:                    c.GetOSCommand().FileType(name),
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

	statusLines, err := c.RunWithOutput(BuildGitCmdObjFromStr(fmt.Sprintf("status %s --porcelain -z %s", opts.UntrackedFilesArg, noRenamesFlag)))
	if err != nil {
		return "", err
	}

	splitLines := strings.Split(statusLines, "\x00")
	// if a line starts with 'R' then the next line is the original file.
	for i := 0; i < len(splitLines)-1; i++ {
		original := splitLines[i]
		if strings.HasPrefix(original, "R  ") {
			next := splitLines[i+1]
			updated := "R  " + next + RENAME_SEPARATOR + strings.TrimPrefix(original, "R  ")
			splitLines[i] = updated
			splitLines = append(splitLines[0:i+1], splitLines[i+2:]...)
		}
	}

	return strings.Join(splitLines, "\n"), nil
}
