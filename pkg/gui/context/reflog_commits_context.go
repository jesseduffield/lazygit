package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ReflogCommitsContext struct {
	*FilteredListViewModel[*models.Commit]
	*ListContextTrait
}

var _ types.IListContext = (*ReflogCommitsContext)(nil)

func NewReflogCommitsContext(
	getItems func() []*models.Commit,
	guiContextState GuiContextState,
	view *gocui.View,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.HelperCommon,
) *ReflogCommitsContext {
	viewModel := NewFilteredListViewModel(getItems, guiContextState.Needle, commitToString)

	getDisplayStrings := func(startIdx int, length int) [][]string {
		return presentation.GetReflogCommitListDisplayStrings(
			viewModel.getModel(),
			guiContextState.ScreenMode() != types.SCREEN_NORMAL,
			cherryPickedCommitShaSet(guiContextState),
			guiContextState.Modes().Diffing.Ref,
			c.UserConfig.Git.ParseEmoji,
		)
	}

	return &ReflogCommitsContext{
		FilteredListViewModel: viewModel,
		ListContextTrait: &ListContextTrait{
			Context: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
				ViewName:   "commits",
				WindowName: "commits",
				Key:        REFLOG_COMMITS_CONTEXT_KEY,
				Kind:       types.SIDE_CONTEXT,
				Focusable:  true,
			}), ContextCallbackOpts{
				OnFocus:        onFocus,
				OnFocusLost:    onFocusLost,
				OnRenderToMain: onRenderToMain,
			}),
			list:              viewModel,
			viewTrait:         NewViewTrait(view),
			getDisplayStrings: getDisplayStrings,
			c:                 c,
		},
	}
}

func (self *ReflogCommitsContext) GetSelectedItemId() string {
	item := self.GetSelected()
	if item == nil {
		return ""
	}

	return item.ID()
}

func (self *ReflogCommitsContext) CanRebase() bool {
	return false
}

func (self *ReflogCommitsContext) GetSelectedRef() types.Ref {
	commit := self.GetSelected()
	if commit == nil {
		return nil
	}
	return commit
}

func (self *ReflogCommitsContext) GetCommits() []*models.Commit {
	return self.getModel()
}
