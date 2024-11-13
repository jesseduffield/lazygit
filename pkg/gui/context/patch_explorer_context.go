package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	deadlock "github.com/sasha-s/go-deadlock"
)

type PatchExplorerContext struct {
	*SimpleContext
	*SearchTrait

	state                  *patch_exploring.State
	viewTrait              *ViewTrait
	getIncludedLineIndices func() []int
	c                      *ContextCommon
	mutex                  *deadlock.Mutex
}

var (
	_ types.IPatchExplorerContext = (*PatchExplorerContext)(nil)
	_ types.ISearchableContext    = (*PatchExplorerContext)(nil)
)

func NewPatchExplorerContext(
	view *gocui.View,
	windowName string,
	key types.ContextKey,

	getIncludedLineIndices func() []int,

	c *ContextCommon,
) *PatchExplorerContext {
	ctx := &PatchExplorerContext{
		state:                  nil,
		viewTrait:              NewViewTrait(view),
		c:                      c,
		mutex:                  &deadlock.Mutex{},
		getIncludedLineIndices: getIncludedLineIndices,
		SimpleContext: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
			View:                       view,
			WindowName:                 windowName,
			Key:                        key,
			Kind:                       types.MAIN_CONTEXT,
			Focusable:                  true,
			HighlightOnFocus:           true,
			NeedsRerenderOnWidthChange: types.NEEDS_RERENDER_ON_WIDTH_CHANGE_WHEN_WIDTH_CHANGES,
		})),
		SearchTrait: NewSearchTrait(c),
	}

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(
		func(selectedLineIdx int) error {
			ctx.GetMutex().Lock()
			defer ctx.GetMutex().Unlock()
			ctx.NavigateTo(selectedLineIdx)
			return nil
		}),
	)

	ctx.SetHandleRenderFunc(ctx.OnViewWidthChanged)

	return ctx
}

func (self *PatchExplorerContext) IsPatchExplorerContext() {}

func (self *PatchExplorerContext) GetState() *patch_exploring.State {
	return self.state
}

func (self *PatchExplorerContext) SetState(state *patch_exploring.State) {
	self.state = state
}

func (self *PatchExplorerContext) GetViewTrait() types.IViewTrait {
	return self.viewTrait
}

func (self *PatchExplorerContext) GetIncludedLineIndices() []int {
	return self.getIncludedLineIndices()
}

func (self *PatchExplorerContext) RenderAndFocus() {
	self.setContent()

	self.FocusSelection()
	self.c.Render()
}

func (self *PatchExplorerContext) Render() {
	self.setContent()

	self.c.Render()
}

func (self *PatchExplorerContext) Focus() {
	self.FocusSelection()
	self.c.Render()
}

func (self *PatchExplorerContext) setContent() {
	self.GetView().SetContent(self.GetContentToRender())
}

func (self *PatchExplorerContext) FocusSelection() {
	view := self.GetView()
	state := self.GetState()
	bufferHeight := view.InnerHeight()
	_, origin := view.Origin()
	numLines := view.ViewLinesHeight()

	newOriginY := state.CalculateOrigin(origin, bufferHeight, numLines)

	view.SetOriginY(newOriginY)

	startIdx, endIdx := state.SelectedViewRange()
	// As far as the view is concerned, we are always selecting a range
	view.SetRangeSelectStart(startIdx)
	view.SetCursorY(endIdx - newOriginY)
}

func (self *PatchExplorerContext) GetContentToRender() string {
	if self.GetState() == nil {
		return ""
	}

	return self.GetState().RenderForLineIndices(self.GetIncludedLineIndices())
}

func (self *PatchExplorerContext) NavigateTo(selectedLineIdx int) {
	self.GetState().SetLineSelectMode()
	self.GetState().SelectLine(selectedLineIdx)

	self.RenderAndFocus()
}

func (self *PatchExplorerContext) GetMutex() *deadlock.Mutex {
	return self.mutex
}

func (self *PatchExplorerContext) ModelSearchResults(searchStr string, caseSensitive bool) []gocui.SearchPosition {
	return nil
}

func (self *PatchExplorerContext) OnViewWidthChanged() {
	if state := self.GetState(); state != nil {
		state.OnViewWidthChanged(self.GetView())
		self.setContent()
		self.RenderAndFocus()
	}
}
