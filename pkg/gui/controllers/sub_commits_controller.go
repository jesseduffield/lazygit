package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SubCommitsController struct {
	baseController
	*ListControllerTrait[*models.Commit]
	c *ControllerCommon
}

var _ types.IController = &SubCommitsController{}

func NewSubCommitsController(
	c *ControllerCommon,
) *SubCommitsController {
	return &SubCommitsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().SubCommits,
			c.Contexts().SubCommits.GetSelected,
			c.Contexts().SubCommits.GetSelectedItems,
		),
		c: c,
	}
}

func (self *SubCommitsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Keys:              opts.GetKeys(opts.Config.Branches.ViewBranchDiff),
			Handler:           self.viewBranchDiff,
			GetDisabledReason: self.require(self.commitsOfABranch),
			Description:       self.c.Tr.ViewBranchDiff,
			Tooltip:           self.c.Tr.ViewBranchDiffTooltip,
		},
	}
}

// viewBranchDiff lists the files of the diff of the branch whose commits this panel
// shows.
func (self *SubCommitsController) viewBranchDiff() error {
	return (&BranchDiffAction{c: self.c}).View(self.context().GetRef(), self.context())
}

// commitsOfABranch disables the branch diff where this panel was entered from
// something other than a local or remote branch, such as a tag or the reflog.
func (self *SubCommitsController) commitsOfABranch() *types.DisabledReason {
	switch self.context().GetRef().(type) {
	case *models.Branch, *models.RemoteBranch:
		return nil
	default:
		return &types.DisabledReason{Text: self.c.Tr.NotCommitsOfABranch}
	}
}

func (self *SubCommitsController) Context() types.Context {
	return self.context()
}

func (self *SubCommitsController) context() *context.SubCommitsContext {
	return self.c.Contexts().SubCommits
}

func (self *SubCommitsController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			commit := self.context().GetSelected()
			var task types.UpdateTask
			if commit == nil {
				task = types.NewRenderStringTask("No commits")
			} else {
				refRange := self.context().GetSelectedRefRangeForDiffFiles()
				task = self.c.Helpers().Diff.GetUpdateTaskForRenderingCommitsDiff(commit, refRange)
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title:    "Commit",
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Task:     task,
				},
				Secondary: secondaryPatchPanelUpdateOpts(self.c),
			})
		})
	}
}

func (self *SubCommitsController) GetOnFocus() func(types.OnFocusOpts) {
	return func(types.OnFocusOpts) {
		loadAllCommitsIfSelectionIsNearEnd(self.c, self.context(), types.SUB_COMMITS)
	}
}
