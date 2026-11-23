package helpers

import (
	"errors"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type UpstreamHelper struct {
	c *HelperCommon

	getRemoteBranchesSuggestionsFunc func(string, ...string) func(string) []*types.Suggestion
}

func NewUpstreamHelper(
	c *HelperCommon,
	getRemoteBranchesSuggestionsFunc func(string, ...string) func(string) []*types.Suggestion,
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

func (self *UpstreamHelper) promptForUpstream(initialContent string, additionalSuggestions []string, onConfirm func(string) error) error {
	self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.EnterUpstream,
		InitialContent:      initialContent,
		FindSuggestionsFunc: self.getRemoteBranchesSuggestionsFunc(" ", additionalSuggestions...),
		HandleConfirm:       onConfirm,
	})

	return nil
}

func (self *UpstreamHelper) PromptForUpstreamWithInitialContent(currentBranch *models.Branch, onConfirm func(string) error) error {
	suggestedRemote := "origin"
	remoteNames := lo.Map(self.c.Model().Remotes, func(r *models.Remote, _ int) string { return r.Name })
	if len(remoteNames) > 0 {
		// remove 'origin' or the first one from the list
		i := max(0, slices.Index(remoteNames, "origin"))
		suggestedRemote = remoteNames[i]
		remoteNames = slices.Delete(remoteNames, i, i+1)
	}

	initialContent := suggestedRemote + " " + currentBranch.Name
	additionalSuggestions := lo.Map(remoteNames, func(remote string, _ int) string { return remote + " " + currentBranch.Name })

	return self.promptForUpstream(initialContent, additionalSuggestions, onConfirm)
}

func (self *UpstreamHelper) PromptForUpstreamWithoutInitialContent(_ *models.Branch, onConfirm func(string) error) error {
	return self.promptForUpstream("", nil, onConfirm)
}
