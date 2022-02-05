package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

func (self *ViewTrait) SetFooter(value string) {
	self.view.Footer = value
}

func (self *ViewTrait) SetOriginX(value int) {
	_ = self.view.SetOriginX(value)
}

// tells us the bounds of line indexes shown in the view currently
func (self *ViewTrait) ViewPortYBounds() (int, int) {
	_, min := self.view.Origin()
	max := self.view.InnerHeight() + 1
	return min, max
}

func (self *ViewTrait) ScrollLeft() {
	newOriginX := utils.Max(self.view.OriginX()-self.view.InnerWidth()/HORIZONTAL_SCROLL_FACTOR, 0)
	_ = self.view.SetOriginX(newOriginX)
}

func (self *ViewTrait) ScrollRight() {
	_ = self.view.SetOriginX(self.view.OriginX() + self.view.InnerWidth()/HORIZONTAL_SCROLL_FACTOR)
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
