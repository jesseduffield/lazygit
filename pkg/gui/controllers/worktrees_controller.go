package controllers

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type WorktreesController struct {
	baseController
	*ListControllerTrait[*models.Worktree]
	c *ControllerCommon
}

var _ types.IController = &WorktreesController{}

func NewWorktreesController(
	c *ControllerCommon,
) *WorktreesController {
	return &WorktreesController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.Worktree](
			c,
			c.Contexts().Worktrees,
			c.Contexts().Worktrees.GetSelected,
			c.Contexts().Worktrees.GetSelectedItems,
		),
		c: c,
	}
}

func (self *WorktreesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:             opts.GetKey(opts.Config.Universal.New),
			Handler:         self.add,
			Description:     self.c.Tr.NewWorktree,
			DisplayOnScreen: true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Switch,
			Tooltip:           self.c.Tr.SwitchToWorktreeTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Confirm),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:           self.withItem(self.open),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenInEditor,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.remove),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Remove,
			Tooltip:           self.c.Tr.RemoveWorktreeTooltip,
			DisplayOnScreen:   true,
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
			if worktree.IsMain {
				main = style.FgDefault.Sprintf(" %s", self.c.Tr.MainWorktree)
			}

			missing := ""
			if worktree.IsPathMissing {
				missing = style.FgRed.Sprintf(" %s", self.c.Tr.MissingWorktree)
			}

			var builder strings.Builder
			w := tabwriter.NewWriter(&builder, 0, 0, 2, ' ', 0)
			_, _ = fmt.Fprintf(w, "%s:\t%s%s\n", self.c.Tr.Name, style.FgGreen.Sprint(worktree.Name), main)
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
	if worktree.IsMain {
		return errors.New(self.c.Tr.CantDeleteMainWorktree)
	}

	if worktree.IsCurrent {
		return errors.New(self.c.Tr.CantDeleteCurrentWorktree)
	}

	return self.c.Helpers().Worktree.Remove(worktree, false)
}

func (self *WorktreesController) GetOnClick() func() error {
	return self.withItemGraceful(self.enter)
}

func (self *WorktreesController) enter(worktree *models.Worktree) error {
	return self.c.Helpers().Worktree.Switch(worktree, context.WORKTREES_CONTEXT_KEY)
}

func (self *WorktreesController) open(worktree *models.Worktree) error {
	return self.c.Helpers().Files.OpenDirInEditor(worktree.Path)
}

func (self *WorktreesController) context() *context.WorktreesContext {
	return self.c.Contexts().Worktrees
}
