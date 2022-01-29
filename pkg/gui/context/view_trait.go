package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

const HORIZONTAL_SCROLL_FACTOR = 3

type ViewTrait struct {
	getView func() *gocui.View
}

func NewViewTrait(getView func() *gocui.View) *ViewTrait {
	return &ViewTrait{getView: getView}
}

func (self *ViewTrait) FocusPoint(yIdx int) {
	view := self.getView()
	view.FocusPoint(view.OriginX(), yIdx)
}

func (self *ViewTrait) SetViewPortContent(content string) {
	view := self.getView()

	_, y := view.Origin()
	view.OverwriteLines(y, content)
}

func (self *ViewTrait) SetContent(content string) {
	self.getView().SetContent(content)
}

func (self *ViewTrait) SetFooter(value string) {
	self.getView().Footer = value
}

func (self *ViewTrait) SetOriginX(value int) {
	self.getView().SetOriginX(value)
}

// tells us the bounds of line indexes shown in the view currently
func (self *ViewTrait) ViewPortYBounds() (int, int) {
	view := self.getView()

	_, min := view.Origin()
	max := view.InnerHeight() + 1
	return min, max
}

func (self *ViewTrait) ScrollLeft() {
	view := self.getView()

	newOriginX := utils.Max(view.OriginX()-view.InnerWidth()/HORIZONTAL_SCROLL_FACTOR, 0)
	_ = view.SetOriginX(newOriginX)
}

func (self *ViewTrait) ScrollRight() {
	view := self.getView()

	_ = view.SetOriginX(view.OriginX() + view.InnerWidth()/HORIZONTAL_SCROLL_FACTOR)
}

// this returns the amount we'll scroll if we want to scroll by a page.
func (self *ViewTrait) PageDelta() int {
	view := self.getView()

	_, height := view.Size()

	delta := height - 1
	if delta == 0 {
		return 1
	}

	return delta
}

func (self *ViewTrait) SelectedLineIdx() int {
	return self.getView().SelectedLineIdx()
}
