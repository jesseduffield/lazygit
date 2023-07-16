package controllers

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/jesseduffield/gocui"
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
			Handler:     self.checkSelected(self.remove),
			Description: self.c.Tr.RemoveWorktree,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.add,
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

func (self *WorktreesController) add() error {
	return self.c.Helpers().Worktree.NewWorktree()
}

func (self *WorktreesController) remove(worktree *models.Worktree) error {
	if worktree.Main() {
		return self.c.ErrorMsg(self.c.Tr.CantDeleteMainWorktree)
	}

	if self.c.Git().Worktree.IsCurrentWorktree(worktree) {
		return self.c.ErrorMsg(self.c.Tr.CantDeleteCurrentWorktree)
	}

	return self.removeWithForce(worktree, false)
}

func (self *WorktreesController) removeWithForce(worktree *models.Worktree, force bool) error {
	title := self.c.Tr.RemoveWorktreeTitle
	var templateStr string
	if force {
		templateStr = self.c.Tr.ForceRemoveWorktreePrompt
	} else {
		templateStr = self.c.Tr.RemoveWorktreePrompt
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
			return self.c.WithWaitingStatus(self.c.Tr.RemovingWorktree, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.RemovingWorktree)
				if err := self.c.Git().Worktree.Delete(worktree.Path, force); err != nil {
					errMessage := err.Error()
					if !force {
						return self.removeWithForce(worktree, true)
					}
					return self.c.ErrorMsg(errMessage)
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.WORKTREES}})
			})
		},
	})
}

func (self *WorktreesController) GetOnClick() func() error {
	return self.checkSelected(self.enter)
}

func (self *WorktreesController) enter(worktree *models.Worktree) error {
	if self.c.Git().Worktree.IsCurrentWorktree(worktree) {
		return self.c.ErrorMsg(self.c.Tr.AlreadyInWorktree)
	}

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
