package context

import (
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PatchExplorerContext struct {
	*SimpleContext

	state                  *patch_exploring.State
	viewTrait              *ViewTrait
	getIncludedLineIndices func() []int
	c                      *types.HelperCommon
	mutex                  *sync.Mutex
}

var _ types.IPatchExplorerContext = (*PatchExplorerContext)(nil)

func NewPatchExplorerContext(
	view *gocui.View,
	windowName string,
	key types.ContextKey,

	onFocus func(types.OnFocusOpts) error,
	onFocusLost func(opts types.OnFocusLostOpts) error,
	getIncludedLineIndices func() []int,

	c *types.HelperCommon,
) *PatchExplorerContext {
	return &PatchExplorerContext{
		state:                  nil,
		viewTrait:              NewViewTrait(view),
		c:                      c,
		mutex:                  &sync.Mutex{},
		getIncludedLineIndices: getIncludedLineIndices,
		SimpleContext: NewSimpleContext(NewBaseContext(NewBaseContextOpts{
			View:       view,
			WindowName: windowName,
			Key:        key,
			Kind:       types.MAIN_CONTEXT,
			Focusable:  true,
		}), ContextCallbackOpts{
			OnFocus:     onFocus,
			OnFocusLost: onFocusLost,
		}),
	}
}

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
	self.GetView().SetContent(self.GetContentToRender(isFocused))

	if err := self.focusSelection(); err != nil {
		return err
	}

	self.c.Render()

	return nil
}

func (self *PatchExplorerContext) Render(isFocused bool) error {
	self.GetView().SetContent(self.GetContentToRender(isFocused))

	self.c.Render()

	return nil
}

func (self *PatchExplorerContext) Focus() error {
	if err := self.focusSelection(); err != nil {
		return err
	}

	self.c.Render()

	return nil
}

func (self *PatchExplorerContext) focusSelection() error {
	view := self.GetView()
	state := self.GetState()
	_, viewHeight := view.Size()
	bufferHeight := viewHeight - 1
	_, origin := view.Origin()

	selectedLineIdx := state.GetSelectedLineIdx()

	newOrigin := state.CalculateOrigin(origin, bufferHeight)

	if err := view.SetOriginY(newOrigin); err != nil {
		return err
	}

	return view.SetCursor(0, selectedLineIdx-newOrigin)
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

func (self *PatchExplorerContext) GetMutex() *sync.Mutex {
	return self.mutex
}
