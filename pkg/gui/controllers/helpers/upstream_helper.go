package helpers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type UpstreamHelper struct {
	c *HelperCommon

	suggestions *SuggestionsHelper
}

type IUpstreamHelper interface {
	PromptForUpstream(suggestedBranch string, branchSelectionTitle func(string) string, onConfirm func(Upstream) error) error
}

var _ IUpstreamHelper = &UpstreamHelper{}

func NewUpstreamHelper(
	c *HelperCommon,
	suggestions *SuggestionsHelper,
) *UpstreamHelper {
	return &UpstreamHelper{
		c:           c,
		suggestions: suggestions,
	}
}

type Upstream struct {
	Remote string
	Branch string
}

func (self *UpstreamHelper) promptForUpstreamBranch(chosenRemote string, initialBranch string, branchSelectionTitle func(string) string, onConfirm func(Upstream) error) error {
	remoteDoesNotExist := lo.NoneBy(self.c.Model().Remotes, func(remote *models.Remote) bool {
		return remote.Name == chosenRemote
	})
	if remoteDoesNotExist {
		return fmt.Errorf(self.c.Tr.NoValidRemoteName, chosenRemote)
	}

	self.c.Prompt(types.PromptOpts{
		Title:               branchSelectionTitle(chosenRemote),
		InitialContent:      initialBranch,
		FindSuggestionsFunc: self.suggestions.GetRemoteBranchesForRemoteSuggestionsFunc(chosenRemote),
		HandleConfirm: func(chosenBranch string) error {
			self.c.Log.Debugf("User selected branch '%s' on remote '%s'", chosenRemote, chosenBranch)
			return onConfirm(Upstream{chosenRemote, chosenBranch})
		},
	})
	return nil
}

// Creates a series of two prompts that will gather an upstream remote and branch from the user.
// The first prompt gathers the remote name, with a pre-filled default of origin, or the first remote in the list.
// After selecting a remote, only branches present on that remote will be displayed.
// If there is only one remote in the current repository, the remote prompt will be skipped.
//
// suggestedBranch pre-fills the second prompt, but can be an empty string if no suggestion would make sense. Often it is the current local branch name.
// branchSelectionTitle allows customization of the second prompt to remind the user what their action entails.
// onConfirm is called once the user has fully specified what their desired upstream is.
func (self *UpstreamHelper) PromptForUpstream(suggestedBranch string, branchSelectionTitle func(string) string, onConfirm func(Upstream) error) error {
	if len(self.c.Model().Remotes) == 1 {
		remote := self.c.Model().Remotes[0].Name
		self.c.Log.Debugf("Defaulting to only remote %s", remote)
		return self.promptForUpstreamBranch(remote, suggestedBranch, branchSelectionTitle, onConfirm)
	} else {
		suggestedRemote := getSuggestedRemote(self.c.Model().Remotes)
		self.c.Prompt(types.PromptOpts{
			Title:               self.c.Tr.SelectTargetRemote,
			InitialContent:      suggestedRemote,
			FindSuggestionsFunc: self.suggestions.GetRemoteSuggestionsFunc(),
			HandleConfirm: func(chosenRemote string) error {
				return self.promptForUpstreamBranch(chosenRemote, suggestedBranch, branchSelectionTitle, onConfirm)
			},
		})
	}

	return nil
}

func getSuggestedRemote(remotes []*models.Remote) string {
	if len(remotes) == 0 {
		return "origin"
	}

	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote.Name
		}
	}

	return remotes[0].Name
}
