package context

import (
	"math"

	"github.com/jesseduffield/lazygit/pkg/gui/mergeconflicts"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/sasha-s/go-deadlock"
)

type MergeConflictsContext struct {
	types.Context
	viewModel *ConflictsViewModel
	c         *ContextCommon
	mutex     *deadlock.Mutex
}

type ConflictsViewModel struct {
	state *mergeconflicts.State

	// userVerticalScrolling tells us if the user has started scrolling through the file themselves
	// in which case we won't auto-scroll to a conflict.
	userVerticalScrolling bool
}

func NewMergeConflictsContext(
	c *ContextCommon,
) *MergeConflictsContext {
	viewModel := &ConflictsViewModel{
		state:                 mergeconflicts.NewState(),
		userVerticalScrolling: false,
	}

	return &MergeConflictsContext{
		viewModel: viewModel,
		mutex:     &deadlock.Mutex{},
		Context: NewSimpleContext(
			NewBaseContext(NewBaseContextOpts{
				Kind:             types.MAIN_CONTEXT,
				View:             c.Views().MergeConflicts,
				WindowName:       "main",
				Key:              MERGE_CONFLICTS_CONTEXT_KEY,
				Focusable:        true,
				HighlightOnFocus: true,
			}),
		),
		c: c,
	}
}

func (self *MergeConflictsContext) GetState() *mergeconflicts.State {
	return self.viewModel.state
}

func (self *MergeConflictsContext) SetState(state *mergeconflicts.State) {
	self.viewModel.state = state
}

func (self *MergeConflictsContext) GetMutex() *deadlock.Mutex {
	return self.mutex
}

func (self *MergeConflictsContext) SetUserScrolling(isScrolling bool) {
	self.viewModel.userVerticalScrolling = isScrolling
}

func (self *MergeConflictsContext) IsUserScrolling() bool {
	return self.viewModel.userVerticalScrolling
}

func (self *MergeConflictsContext) RenderAndFocus(isFocused bool) error {
	self.setContent(isFocused)
	self.FocusSelection()

	self.c.Render()

	return nil
}

func (self *MergeConflictsContext) Render(isFocused bool) error {
	self.setContent(isFocused)

	self.c.Render()

	return nil
}

func (self *MergeConflictsContext) GetContentToRender(isFocused bool) string {
	if self.GetState() == nil {
		return ""
	}

	return mergeconflicts.ColoredConflictFile(self.GetState(), isFocused)
}

func (self *MergeConflictsContext) setContent(isFocused bool) {
	self.GetView().SetContent(self.GetContentToRender(isFocused))
}

func (self *MergeConflictsContext) FocusSelection() {
	if !self.IsUserScrolling() {
		_ = self.GetView().SetOriginY(self.GetOriginY())
	}
}

func (self *MergeConflictsContext) GetOriginY() int {
	view := self.GetView()
	conflictMiddle := self.GetState().GetConflictMiddle()
	return int(math.Max(0, float64(conflictMiddle-(view.Height()/2))))
}
