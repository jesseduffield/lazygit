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

var _ types.IPatchExplorerContext = (*PatchExplorerContext)(nil)

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
			View:             view,
			WindowName:       windowName,
			Key:              key,
			Kind:             types.MAIN_CONTEXT,
			Focusable:        true,
			HighlightOnFocus: true,
		})),
		SearchTrait: NewSearchTrait(c),
	}

	ctx.GetView().SetOnSelectItem(ctx.SearchTrait.onSelectItemWrapper(
		func(selectedLineIdx int) error {
			ctx.GetMutex().Lock()
			defer ctx.GetMutex().Unlock()
			return ctx.NavigateTo(ctx.c.IsCurrentContext(ctx), selectedLineIdx)
		}),
	)

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

func (self *PatchExplorerContext) RenderAndFocus(isFocused bool) error {
	self.setContent(isFocused)

	self.FocusSelection()
	self.c.Render()

	return nil
}

func (self *PatchExplorerContext) Render(isFocused bool) error {
	self.setContent(isFocused)

	self.c.Render()

	return nil
}

func (self *PatchExplorerContext) Focus() error {
	self.FocusSelection()
	self.c.Render()

	return nil
}

func (self *PatchExplorerContext) setContent(isFocused bool) {
	self.GetView().SetContent(self.GetContentToRender(isFocused))
}

func (self *PatchExplorerContext) FocusSelection() {
	view := self.GetView()
	state := self.GetState()
	_, viewHeight := view.Size()
	bufferHeight := viewHeight - 1
	_, origin := view.Origin()

	newOriginY := state.CalculateOrigin(origin, bufferHeight)

	_ = view.SetOriginY(newOriginY)

	view.SetCursorY(state.GetSelectedLineIdx() - newOriginY)
}

func (self *PatchExplorerContext) GetContentToRender(isFocused bool) string {
	if self.GetState() == nil {
		return ""
	}

	return self.GetState().RenderForLineIndices(isFocused, self.GetIncludedLineIndices())
}

func (self *PatchExplorerContext) NavigateTo(isFocused bool, selectedLineIdx int) error {
	self.GetState().SetLineSelectMode()
	self.GetState().SelectLine(selectedLineIdx)

	return self.RenderAndFocus(isFocused)
}

func (self *PatchExplorerContext) GetMutex() *deadlock.Mutex {
	return self.mutex
}
