package controllers

import "fmt"

type PortraitModeActions struct {
	c *ControllerCommon
}

var portraitModes = []string{"auto", "always", "never"}

func (self *PortraitModeActions) Next() error {
	return self.cycle(1)
}

func (self *PortraitModeActions) Prev() error {
	return self.cycle(-1)
}

func (self *PortraitModeActions) cycle(direction int) error {
	current := self.c.UserConfig().Gui.PortraitMode
	currentIdx := 0
	for i, mode := range portraitModes {
		if mode == current {
			currentIdx = i
			break
		}
	}
	newIdx := (currentIdx + direction + len(portraitModes)) % len(portraitModes)
	self.c.UserConfig().Gui.PortraitMode = portraitModes[newIdx]

	self.rerenderViewsWithScreenModeDependentContent()

	self.c.Toast(fmt.Sprintf("Portrait mode: %s", portraitModes[newIdx]))
	return nil
}

func (self *PortraitModeActions) rerenderViewsWithScreenModeDependentContent() {
	(&ScreenModeActions{c: self.c}).rerenderViewsWithScreenModeDependentContent()
}
