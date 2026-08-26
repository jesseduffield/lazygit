package helpers

import (
	"errors"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type UpstreamHelper struct {
	c *HelperCommon

	getRemoteBranchesSuggestionsFunc func(string) func(string) []*types.Suggestion
}

func NewUpstreamHelper(
	c *HelperCommon,
	getRemoteBranchesSuggestionsFunc func(string) func(string) []*types.Suggestion,
) *UpstreamHelper {
	return &UpstreamHelper{
		c:                                c,
		getRemoteBranchesSuggestionsFunc: getRemoteBranchesSuggestionsFunc,
	}
}

func (self *UpstreamHelper) ParseUpstream(upstream string) (string, string, error) {
	var upstreamBranch, upstreamRemote string
	split := strings.Split(upstream, " ")
	if len(split) != 2 {
		return "", "", errors.New(self.c.Tr.InvalidUpstream)
	}

	upstreamRemote = split[0]
	upstreamBranch = split[1]

	return upstreamRemote, upstreamBranch, nil
}

func (self *UpstreamHelper) promptForUpstream(initialContent string, isPush bool, onConfirm func(string) error) error {
	title := self.c.Tr.EnterUpstream
	if isPush {
		title = self.c.Tr.EnterUpstreamNewBranch
	}

	self.c.Prompt(types.PromptOpts{
		Title:               title,
		InitialContent:      initialContent,
		FindSuggestionsFunc: self.getRemoteBranchesSuggestionsFunc(" "),
		HandleConfirm:       onConfirm,
	})

	return nil
}

func (self *UpstreamHelper) PromptForUpstreamWithInitialContent(currentBranch *models.Branch, isPush bool, onConfirm func(string) error) error {
	initialContent := currentBranch.UpstreamRemote + " " + currentBranch.UpstreamBranch
	if !currentBranch.IsTrackingRemote() {
		initialContent = self.GetSuggestedRemote() + " " + currentBranch.Name
	}

	return self.promptForUpstream(initialContent, isPush, onConfirm)
}

func (self *UpstreamHelper) PromptForUpstreamWithoutInitialContent(_ *models.Branch, onConfirm func(string) error) error {
	return self.promptForUpstream("", false, onConfirm)
}

func (self *UpstreamHelper) GetSuggestedRemote() string {
	return getSuggestedRemote(self.c.Model().Remotes)
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
