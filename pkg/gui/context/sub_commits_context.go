package context

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type SubCommitsContext struct {
	*SubCommitsViewModel
	*ViewportListContextTrait
	*DynamicTitleBuilder
}

var (
	_ types.IListContext    = (*SubCommitsContext)(nil)
	_ types.DiffableContext = (*SubCommitsContext)(nil)
)

func NewSubCommitsContext(
	c *ContextCommon,
) *SubCommitsContext {
	viewModel := &SubCommitsViewModel{
		BasicViewModel: NewBasicViewModel(
			func() []*models.Commit { return c.Model().SubCommits },
		),
		ref:          nil,
		limitCommits: true,
	}

	getDisplayStrings := func(startIdx int, length int) [][]string {
		selectedCommitSha := ""
		if c.CurrentContext().GetKey() == SUB_COMMITS_CONTEXT_KEY {
			selectedCommit := viewModel.GetSelected()
			if selectedCommit != nil {
				selectedCommitSha = selectedCommit.Sha
			}
		}
		return presentation.GetCommitListDisplayStrings(
			c.Common,
			c.Model().SubCommits,
			c.State().GetRepoState().GetScreenMode() != types.SCREEN_NORMAL,
			c.Modes().CherryPicking.SelectedShaSet(),
			c.Modes().Diffing.Ref,
			c.UserConfig.Gui.TimeFormat,
			c.UserConfig.Git.ParseEmoji,
			selectedCommitSha,
			startIdx,
			length,
			shouldShowGraph(c),
			git_commands.NewNullBisectInfo(),
			false,
		)
	}

	return &SubCommitsContext{
		SubCommitsViewModel: viewModel,
		DynamicTitleBuilder: NewDynamicTitleBuilder(c.Tr.SubCommitsDynamicTitle),
		ViewportListContextTrait: &ViewportListContextTrait{
			ListContextTrait: &ListContextTrait{
				Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
					View:       c.Views().SubCommits,
					WindowName: "branches",
					Key:        SUB_COMMITS_CONTEXT_KEY,
					Kind:       types.SIDE_CONTEXT,
					Focusable:  true,
					Transient:  true,
				})),
				list:              viewModel,
				getDisplayStrings: getDisplayStrings,
				c:                 c,
			},
		},
	}
}

type SubCommitsViewModel struct {
	// name of the ref that the sub-commits are shown for
	ref types.Ref
	*BasicViewModel[*models.Commit]

	limitCommits bool
}

func (self *SubCommitsViewModel) SetRef(ref types.Ref) {
	self.ref = ref
}

func (self *SubCommitsViewModel) GetRef() types.Ref {
	return self.ref
}

func (self *SubCommitsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *SubCommitsContext) CanRebase() bool {
	return false
}

func (self *SubCommitsContext) GetSelectedRef() types.Ref {
	commit := self.GetSelected()
	if commit == nil {
		return nil
	}
	return commit
}

func (self *SubCommitsContext) GetCommits() []*models.Commit {
	return self.getModel()
}

func (self *SubCommitsContext) Title() string {
	return fmt.Sprintf(self.c.Tr.SubCommitsDynamicTitle, utils.TruncateWithEllipsis(self.ref.RefName(), 50))
}

func (self *SubCommitsContext) SetLimitCommits(value bool) {
	self.limitCommits = value
}

func (self *SubCommitsContext) GetLimitCommits() bool {
	return self.limitCommits
}

func (self *SubCommitsContext) GetDiffTerminals() []string {
	itemId := self.GetSelectedItemId()

	return []string{itemId}
}
