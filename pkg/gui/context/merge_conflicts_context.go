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

func (self *MergeConflictsContext) RenderAndFocus() error {
	self.setContent()
	self.FocusSelection()

	self.c.Render()

	return nil
}

func (self *MergeConflictsContext) Render() error {
	self.setContent()

	self.c.Render()

	return nil
}

func (self *MergeConflictsContext) GetContentToRender() string {
	if self.GetState() == nil {
		return ""
	}

	return mergeconflicts.ColoredConflictFile(self.GetState())
}

func (self *MergeConflictsContext) setContent() {
	self.GetView().SetContent(self.GetContentToRender())
}

func (self *MergeConflictsContext) FocusSelection() {
	if !self.IsUserScrolling() {
		_ = self.GetView().SetOriginY(self.GetOriginY())
	}

	self.SetSelectedLineRange()
}

func (self *MergeConflictsContext) SetSelectedLineRange() {
	startIdx, endIdx := self.GetState().GetSelectedRange()
	view := self.GetView()
	originY := view.OriginY()
	// As far as the view is concerned, we are always selecting a range
	view.SetRangeSelectStart(startIdx)
	view.SetCursorY(endIdx - originY)
}

func (self *MergeConflictsContext) GetOriginY() int {
	view := self.GetView()
	conflictMiddle := self.GetState().GetConflictMiddle()
	return int(math.Max(0, float64(conflictMiddle-(view.Height()/2))))
}
