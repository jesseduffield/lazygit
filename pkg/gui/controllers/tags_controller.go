package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TagsController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &TagsController{}

func NewTagsController(
	common *ControllerCommon,
) *TagsController {
	return &TagsController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *TagsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.withSelectedTag(self.checkout),
			Description: self.c.Tr.Checkout,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.withSelectedTag(self.delete),
			Description: self.c.Tr.DeleteTag,
		},
		{
			Key:         opts.GetKey(opts.Config.Branches.PushTag),
			Handler:     self.withSelectedTag(self.push),
			Description: self.c.Tr.PushTag,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.create,
			Description: self.c.Tr.CreateTag,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.withSelectedTag(self.createResetMenu),
			Description: self.c.Tr.ViewResetOptions,
			OpensMenu:   true,
		},
	}

	return bindings
}

func (self *TagsController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			var task types.UpdateTask
			tag := self.context().GetSelected()
			if tag == nil {
				task = types.NewRenderStringTask("No tags")
			} else {
				cmdObj := self.c.Git().Branch.GetGraphCmdObj(tag.FullRefName())
				task = types.NewRunCommandTask(cmdObj.GetCmd())
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title: "Tag",
					Task:  task,
				},
			})
		})
	}
}

func (self *TagsController) checkout(tag *models.Tag) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutTag)
	if err := self.c.Helpers().Refs.CheckoutRef(tag.Name, types.CheckoutRefOptions{}); err != nil {
		return err
	}
	return self.c.PushContext(self.c.Contexts().Branches)
}

func (self *TagsController) delete(tag *models.Tag) error {
	prompt := utils.ResolvePlaceholderString(
		self.c.Tr.DeleteTagPrompt,
		map[string]string{
			"tagName": tag.Name,
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DeleteTagTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.DeleteTag)
			if err := self.c.Git().Tag.Delete(tag.Name); err != nil {
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
		FindSuggestionsFunc: self.c.Helpers().Suggestions.GetRemoteSuggestionsFunc(),
		HandleConfirm: func(response string) error {
			return self.c.WithWaitingStatus(self.c.Tr.PushingTagStatus, func(task gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.PushTag)
				err := self.c.Git().Tag.Push(task, response, tag.Name)
				if err != nil {
					_ = self.c.Error(err)
				}

				return nil
			})
		},
	})
}

func (self *TagsController) createResetMenu(tag *models.Tag) error {
	return self.c.Helpers().Refs.CreateGitResetMenu(tag.Name)
}

func (self *TagsController) create() error {
	// leaving commit SHA blank so that we're just creating the tag for the current commit
	return self.c.Helpers().Tags.CreateTagMenu("", func() { self.context().SetSelectedLineIdx(0) })
}

func (self *TagsController) withSelectedTag(f func(tag *models.Tag) error) func() error {
	return func() error {
		tag := self.context().GetSelected()
		if tag == nil {
			return nil
		}

		return f(tag)
	}
}

func (self *TagsController) Context() types.Context {
	return self.context()
}

func (self *TagsController) context() *context.TagsContext {
	return self.c.Contexts().Tags
}
