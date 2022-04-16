package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

const HORIZONTAL_SCROLL_FACTOR = 3

type ViewTrait struct {
	view *gocui.View
}

var _ types.IViewTrait = &ViewTrait{}

func NewViewTrait(view *gocui.View) *ViewTrait {
	return &ViewTrait{view: view}
}

func (self *ViewTrait) FocusPoint(yIdx int) {
	self.view.FocusPoint(self.view.OriginX(), yIdx)
}

func (self *ViewTrait) SetViewPortContent(content string) {
	_, y := self.view.Origin()
	self.view.OverwriteLines(y, content)
}

func (self *ViewTrait) SetContent(content string) {
	self.view.SetContent(content)
}

func (self *ViewTrait) SetHighlight(highlight bool) {
	self.view.Highlight = highlight
}

func (self *ViewTrait) SetFooter(value string) {
	self.view.Footer = value
}

func (self *ViewTrait) SetOriginX(value int) {
	_ = self.view.SetOriginX(value)
}

// tells us the start of line indexes shown in the view currently as well as the capacity of lines shown in the viewport.
func (self *ViewTrait) ViewPortYBounds() (int, int) {
	_, start := self.view.Origin()
	length := self.view.InnerHeight() + 1
	return start, length
}

func (self *ViewTrait) ScrollLeft() {
	self.view.ScrollLeft(self.horizontalScrollAmount())
}

func (self *ViewTrait) ScrollRight() {
	self.view.ScrollRight(self.horizontalScrollAmount())
}

func (self *ViewTrait) horizontalScrollAmount() int {
	return self.view.InnerWidth() / HORIZONTAL_SCROLL_FACTOR
}

func (self *ViewTrait) ScrollUp() {
	self.view.ScrollUp(1)
}

func (self *ViewTrait) ScrollDown() {
	self.view.ScrollDown(1)
}

// this returns the amount we'll scroll if we want to scroll by a page.
func (self *ViewTrait) PageDelta() int {
	_, height := self.view.Size()

	delta := height - 1
	if delta == 0 {
		return 1
	}

	return delta
}

func (self *ViewTrait) SelectedLineIdx() int {
	return self.view.SelectedLineIdx()
}
