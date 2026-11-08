package context

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
)

// Returns the index of the commit that the "go to bottom" command selects in a
// list of commits
func indexForGotoBottomInCommits(commits []*models.Commit, selectedIdx int) int {
	if selectedIdx >= 0 && selectedIdx < len(commits)-1 {
		if commits[selectedIdx+1].Status != models.StatusMerged {
			_, idx, found := lo.FindIndexOf(commits, func(c *models.Commit) bool {
				return c.Status == models.StatusMerged
			})
			if found {
				return idx - 1
			}
		}
	}

	return len(commits) - 1
}
