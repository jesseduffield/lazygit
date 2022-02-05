package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TagsController struct {
	baseController

	c          *types.ControllerCommon
	context    *context.TagsContext
	git        *commands.GitCommand
	contexts   *context.ContextTree
	tagsHelper *TagsHelper

	refsHelper        IRefsHelper
	suggestionsHelper ISuggestionsHelper

	switchToSubCommitsContext func(string) error
}

var _ types.IController = &TagsController{}

func NewTagsController(
	c *types.ControllerCommon,
	context *context.TagsContext,
	git *commands.GitCommand,
	contexts *context.ContextTree,
	tagsHelper *TagsHelper,
	refsHelper IRefsHelper,
	suggestionsHelper ISuggestionsHelper,

	switchToSubCommitsContext func(string) error,
) *TagsController {
	return &TagsController{
		baseController:    baseController{},
		c:                 c,
		context:           context,
		git:               git,
		contexts:          contexts,
		tagsHelper:        tagsHelper,
		refsHelper:        refsHelper,
		suggestionsHelper: suggestionsHelper,

		switchToSubCommitsContext: switchToSubCommitsContext,
	}
}

func (self *TagsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.withSelectedTag(self.checkout),
			Description: self.c.Tr.LcCheckout,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.withSelectedTag(self.delete),
			Description: self.c.Tr.LcDeleteTag,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.PushTag),
			Handler:     self.withSelectedTag(self.push),
			Description: self.c.Tr.LcPushTag,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.create,
			Description: self.c.Tr.LcCreateTag,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.withSelectedTag(self.createResetMenu),
			Description: self.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.withSelectedTag(self.enter),
			Description: self.c.Tr.LcViewCommits,
		},
	}

	return bindings
}

func (self *TagsController) checkout(tag *models.Tag) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutTag)
	if err := self.refsHelper.CheckoutRef(tag.Name, types.CheckoutRefOptions{}); err != nil {
		return err
	}
	return self.c.PushContext(self.contexts.Branches)
}

func (self *TagsController) enter(tag *models.Tag) error {
	return self.switchToSubCommitsContext(tag.Name)
}

func (self *TagsController) delete(tag *models.Tag) error {
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteTagPrompt,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.DeleteTagTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.DeleteTag)
			if err := self.git.Tag.Delete(tag.Name); err != nil {
				return self.c.Error(err)
			}
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
		},
	})
}

func (self *TagsController) push(tag *models.Tag) error {
	title := utils.ResolvePlaceholderString(
		self.c.Tr.PushTagTitle,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	return self.c.Prompt(types.PromptOpts{
		Title:               title,
		InitialContent:      "origin",
		FindSuggestionsFunc: self.suggestionsHelper.GetRemoteSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			return self.c.WithWaitingStatus(self.c.Tr.PushingTagStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.PushTag)
				err := self.git.Tag.Push(response, tag.Name)
				if err != nil {
					_ = self.c.Error(err)
				}

				return nil
			})
		},
	})
}

func (self *TagsController) createResetMenu(tag *models.Tag) error {
	return self.refsHelper.CreateGitResetMenu(tag.Name)
}

func (self *TagsController) create() error {
	// leaving commit SHA blank so that we're just creating the tag for the current commit
	return self.tagsHelper.CreateTagMenu("", func() { self.context.GetPanelState().SetSelectedLineIdx(0) })
}

func (self *TagsController) withSelectedTag(f func(tag *models.Tag) error) func() error {
	return func() error {
		tag := self.context.GetSelected()
		if tag == nil {
			return nil
		}

		return f(tag)
	}
}

func (self *TagsController) Context() types.Context {
	return self.context
}
