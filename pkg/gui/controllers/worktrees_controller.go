package controllers

import (
	"fmt"
	"os"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type WorktreesController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &WorktreesController{}

func NewWorktreesController(
	common *ControllerCommon,
) *WorktreesController {
	return &WorktreesController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *WorktreesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.EnterWorktree,
		},
		//{
		//	Key:         opts.GetKey(opts.Config.Universal.Remove),
		//	Handler:     self.withSelectedTag(self.delete),
		//	Description: self.c.Tr.LcDeleteTag,
		//},
		//{
		//	Key:         opts.GetKey(opts.Config.Branches.PushTag),
		//	Handler:     self.withSelectedTag(self.push),
		//	Description: self.c.Tr.LcPushTag,
		//},
		//{
		//	Key:         opts.GetKey(opts.Config.Universal.New),
		//	Handler:     self.create,
		//	Description: self.c.Tr.LcCreateTag,
		//},
		//{
		//	Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
		//	Handler:     self.withSelectedTag(self.createResetMenu),
		//	Description: self.c.Tr.LcViewResetOptions,
		//	OpensMenu:   true,
		//},
	}

	return bindings
}

func (self *WorktreesController) GetOnRenderToMain() func() error {
	return func() error {
		var task types.UpdateTask
		worktree := self.context().GetSelected()
		if worktree == nil {
			task = types.NewRenderStringTask("No worktrees")
		} else {
			task = types.NewRenderStringTask(fmt.Sprintf("%s\nPath: %s", style.FgGreen.Sprint(worktree.Name), worktree.Path))
		}

		return self.c.RenderToMainViews(types.RefreshMainOpts{
			Pair: self.c.MainViewPairs().Normal,
			Main: &types.ViewUpdateOpts{
				Title: "Worktree",
				Task:  task,
			},
		})
	}
}

//func (self *WorktreesController) switchToWorktree(worktree *models.Worktree) error {
//	//self.c.LogAction(self.c.Tr.Actions.CheckoutTag)
//	//if err := self.helpers.Refs.CheckoutRef(tag.Name, types.CheckoutRefOptions{}); err != nil {
//	//	return err
//	//}
//	//return self.c.PushContext(self.contexts.Branches)
//
//	wd, err := os.Getwd()
//	if err != nil {
//		return err
//	}
//	gui.RepoPathStack.Push(wd)
//
//	return gui.dispatchSwitchToRepo(submodule.Path, true)
//}

//	func (self *WorktreesController) delete(tag *models.Tag) error {
//		prompt := utils.ResolvePlaceholderString(
//			self.c.Tr.DeleteTagPrompt,
//			map[string]string{
//				"tagName": tag.Name,
//			},
//		)
//
//		return self.c.Confirm(types.ConfirmOpts{
//			Title:  self.c.Tr.DeleteTagTitle,
//			Prompt: prompt,
//			HandleConfirm: func() error {
//				self.c.LogAction(self.c.Tr.Actions.DeleteTag)
//				if err := self.git.Tag.Delete(tag.Name); err != nil {
//					return self.c.Error(err)
//				}
//				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.COMMITS, types.TAGS}})
//			},
//		})
//	}
//
//	func (self *WorktreesController) push(tag *models.Tag) error {
//		title := utils.ResolvePlaceholderString(
//			self.c.Tr.PushTagTitle,
//			map[string]string{
//				"tagName": tag.Name,
//			},
//		)
//
//		return self.c.Prompt(types.PromptOpts{
//			Title:               title,
//			InitialContent:      "origin",
//			FindSuggestionsFunc: self.helpers.Suggestions.GetRemoteSuggestionsFunc(),
//			HandleConfirm: func(response string) error {
//				return self.c.WithWaitingStatus(self.c.Tr.PushingTagStatus, func() error {
//					self.c.LogAction(self.c.Tr.Actions.PushTag)
//					err := self.git.Tag.Push(response, tag.Name)
//					if err != nil {
//						_ = self.c.Error(err)
//					}
//
//					return nil
//				})
//			},
//		})
//	}
//
//	func (self *WorktreesController) createResetMenu(tag *models.Tag) error {
//		return self.helpers.Refs.CreateGitResetMenu(tag.Name)
//	}
//
//	func (self *WorktreesController) create() error {
//		// leaving commit SHA blank so that we're just creating the tag for the current commit
//		return self.helpers.Tags.CreateTagMenu("", func() { self.context().SetSelectedLineIdx(0) })
//	}

func (self *WorktreesController) GetOnClick() func() error {
	return self.checkSelected(self.enter)
}

func (self *WorktreesController) enter(worktree *models.Worktree) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	self.c.State().GetRepoPathStack().Push(wd)

	return self.c.Helpers().Repos.DispatchSwitchToRepo(worktree.Path, true)
}

func (self *WorktreesController) checkSelected(callback func(worktree *models.Worktree) error) func() error {
	return func() error {
		worktree := self.context().GetSelected()
		if worktree == nil {
			return nil
		}

		return callback(worktree)
	}
}

func (self *WorktreesController) Context() types.Context {
	return self.context()
}

func (self *WorktreesController) context() *context.WorktreesContext {
	return self.c.Contexts().Worktrees
}
