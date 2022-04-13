package helpers

import (
	"errors"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type UpstreamHelper struct {
	c     *types.HelperCommon
	model *types.Model

	getRemoteBranchesSuggestionsFunc func(string) func(string) []*types.Suggestion
}

type IUpstreamHelper interface {
	ParseUpstream(string) (string, string, error)
	PromptForUpstreamWithInitialContent(*models.Branch, func(string) error) error
	PromptForUpstreamWithoutInitialContent(*models.Branch, func(string) error) error
	GetSuggestedRemote() string
}

var _ IUpstreamHelper = &UpstreamHelper{}

func NewUpstreamHelper(
	c *types.HelperCommon,
	model *types.Model,
	getRemoteBranchesSuggestionsFunc func(string) func(string) []*types.Suggestion,
) *UpstreamHelper {
	return &UpstreamHelper{
		c:                                c,
		model:                            model,
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

func (self *UpstreamHelper) promptForUpstream(currentBranch *models.Branch, initialContent string, onConfirm func(string) error) error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.EnterUpstream,
		InitialContent:      initialContent,
		FindSuggestionsFunc: self.getRemoteBranchesSuggestionsFunc(" "),
		HandleConfirm:       onConfirm,
	})
}

func (self *UpstreamHelper) PromptForUpstreamWithInitialContent(currentBranch *models.Branch, onConfirm func(string) error) error {
	suggestedRemote := self.GetSuggestedRemote()
	initialContent := suggestedRemote + " " + currentBranch.Name

	return self.promptForUpstream(currentBranch, initialContent, onConfirm)
}

func (self *UpstreamHelper) PromptForUpstreamWithoutInitialContent(currentBranch *models.Branch, onConfirm func(string) error) error {
	return self.promptForUpstream(currentBranch, "", onConfirm)
}

func (self *UpstreamHelper) GetSuggestedRemote() string {
	return getSuggestedRemote(self.model.Remotes)
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
