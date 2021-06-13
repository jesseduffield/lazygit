package commands

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (c *GitCommand) getUnfilteredStashEntries() []*models.StashEntry {
	rawString, _ := c.RunCommandWithOutput(
		BuildGitCmdObjFromStr("stash list --pretty='%gs'"),
	)
	stashEntries := []*models.StashEntry{}
	for i, line := range utils.SplitLines(rawString) {
		stashEntries = append(stashEntries, stashEntryFromLine(line, i))
	}
	return stashEntries
}

// GetStashEntries stash entries
func (c *GitCommand) GetStashEntries(filterPath string) []*models.StashEntry {
	if filterPath == "" {
		return c.getUnfilteredStashEntries()
	}

	rawString, err := c.RunCommandWithOutput(
		BuildGitCmdObjFromStr("stash list --name-only"),
	)
	if err != nil {
		return c.getUnfilteredStashEntries()
	}
	stashEntries := []*models.StashEntry{}
	var currentStashEntry *models.StashEntry
	lines := utils.SplitLines(rawString)
	isAStash := func(line string) bool { return strings.HasPrefix(line, "stash@{") }
	re := regexp.MustCompile(`stash@\{(\d+)\}`)

outer:
	for i := 0; i < len(lines); i++ {
		if !isAStash(lines[i]) {
			continue
		}
		match := re.FindStringSubmatch(lines[i])
		idx, err := strconv.Atoi(match[1])
		if err != nil {
			return c.getUnfilteredStashEntries()
		}
		currentStashEntry = stashEntryFromLine(lines[i], idx)
		for i+1 < len(lines) && !isAStash(lines[i+1]) {
			i++
			if lines[i] == filterPath {
				stashEntries = append(stashEntries, currentStashEntry)
				continue outer
			}
		}
	}
	return stashEntries
}

func stashEntryFromLine(line string, index int) *models.StashEntry {
	return &models.StashEntry{
		Name:  line,
		Index: index,
	}
}
