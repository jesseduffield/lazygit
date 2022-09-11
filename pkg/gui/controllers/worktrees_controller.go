package controllers

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.delete),
			Description: self.c.Tr.DeleteWorktree,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.create,
			Description: self.c.Tr.CreateWorktree,
		},
	}

	return bindings
}

func (self *WorktreesController) GetOnRenderToMain() func() error {
	return func() error {
		var task types.UpdateTask
		worktree := self.context().GetSelected()
		if worktree == nil {
			task = types.NewRenderStringTask(self.c.Tr.NoWorktreesThisRepo)
		} else {
			main := ""
			if worktree.Main() {
				main = style.FgDefault.Sprintf(" %s", self.c.Tr.MainWorktree)
			}

			missing := ""
			if self.c.Git().Worktree.IsWorktreePathMissing(worktree) {
				missing = style.FgRed.Sprintf(" %s", self.c.Tr.MissingWorktree)
			}

			var builder strings.Builder
			w := tabwriter.NewWriter(&builder, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintf(w, "%s:\t%s%s\n", self.c.Tr.Name, style.FgGreen.Sprint(worktree.Name()), main)
			_, _ = fmt.Fprintf(w, "%s:\t%s\n", self.c.Tr.Branch, style.FgYellow.Sprint(worktree.Branch))
			_, _ = fmt.Fprintf(w, "%s:\t%s%s\n", self.c.Tr.Path, style.FgCyan.Sprint(worktree.Path), missing)
			_ = w.Flush()

			task = types.NewRenderStringTask(builder.String())
		}

		return self.c.RenderToMainViews(types.RefreshMainOpts{
			Pair: self.c.MainViewPairs().Normal,
			Main: &types.ViewUpdateOpts{
				Title: self.c.Tr.WorktreeTitle,
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

func (self *WorktreesController) create() error {
	return self.c.Helpers().Worktree.NewWorktree()
}

func (self *WorktreesController) delete(worktree *models.Worktree) error {
	if worktree.Main() {
		return self.c.ErrorMsg(self.c.Tr.CantDeleteMainWorktree)
	}

	if self.c.Git().Worktree.IsCurrentWorktree(worktree) {
		return self.c.ErrorMsg(self.c.Tr.CantDeleteCurrentWorktree)
	}

	return self.deleteWithForce(worktree, false)
}

func (self *WorktreesController) deleteWithForce(worktree *models.Worktree, force bool) error {
	title := self.c.Tr.DeleteWorktreeTitle
	var templateStr string
	if force {
		templateStr = self.c.Tr.ForceDeleteWorktreePrompt
	} else {
		templateStr = self.c.Tr.DeleteWorktreePrompt
	}
	message := utils.ResolvePlaceholderString(
		templateStr,
		map[string]string{
			"worktreeName": worktree.Name(),
		},
	)

	return self.c.Confirm(types.ConfirmOpts{
		Title:  title,
		Prompt: message,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.DeleteWorktree)
			if err := self.c.Git().Worktree.Delete(worktree.Path, force); err != nil {
				errMessage := err.Error()
				if !force {
					return self.deleteWithForce(worktree, true)
				}
				return self.c.ErrorMsg(errMessage)
			}
			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.WORKTREES}})
		},
	})
}

func (self *WorktreesController) GetOnClick() func() error {
	return self.checkSelected(self.enter)
}

func (self *WorktreesController) enter(worktree *models.Worktree) error {
	// if we were in a submodule, we want to forget about that stack of repos
	// so that hitting escape in the new repo does nothing
	self.c.State().GetRepoPathStack().Clear()

	return self.c.Helpers().Repos.DispatchSwitchTo(worktree.Path, true, self.c.Tr.ErrWorktreeMovedOrDeleted)
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
