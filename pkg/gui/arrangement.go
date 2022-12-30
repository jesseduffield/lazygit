package gui

import (
	"github.com/jesseduffield/lazycore/pkg/boxlayout"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mattn/go-runewidth"
)

// In this file we use the boxlayout package, along with knowledge about the app's state,
// to arrange the windows (i.e. panels) on the screen.

const INFO_SECTION_PADDING = " "

type WindowArranger struct {
	gui *Gui
}

func (self *WindowArranger) getWindowDimensions(informationStr string, appStatus string) map[string]boxlayout.Dimensions {
	width, height := self.gui.g.Size()

	sideSectionWeight, mainSectionWeight := self.getMidSectionWeights()

	sidePanelsDirection := boxlayout.COLUMN
	portraitMode := width <= 84 && height > 45
	if portraitMode {
		sidePanelsDirection = boxlayout.ROW
	}

	mainPanelsDirection := boxlayout.ROW
	if self.splitMainPanelSideBySide() {
		mainPanelsDirection = boxlayout.COLUMN
	}

	extrasWindowSize := self.getExtrasWindowSize(height)

	showInfoSection := self.gui.c.UserConfig.Gui.ShowBottomLine || self.gui.State.Searching.isSearching || self.gui.isAnyModeActive() || self.gui.statusManager.showStatus()
	infoSectionSize := 0
	if showInfoSection {
		infoSectionSize = 1
	}

	root := &boxlayout.Box{
		Direction: boxlayout.ROW,
		Children: []*boxlayout.Box{
			{
				Direction: sidePanelsDirection,
				Weight:    1,
				Children: []*boxlayout.Box{
					{
						Direction:           boxlayout.ROW,
						Weight:              sideSectionWeight,
						ConditionalChildren: self.sidePanelChildren,
					},
					{
						Direction: boxlayout.ROW,
						Weight:    mainSectionWeight,
						Children: []*boxlayout.Box{
							{
								Direction: mainPanelsDirection,
								Children:  self.mainSectionChildren(),
								Weight:    1,
							},
							{
								Window: "extras",
								Size:   extrasWindowSize,
							},
						},
					},
				},
			},
			{
				Direction: boxlayout.COLUMN,
				Size:      infoSectionSize,
				Children:  self.infoSectionChildren(informationStr, appStatus),
			},
		},
	}

	layerOneWindows := boxlayout.ArrangeWindows(root, 0, 0, width, height)
	limitWindows := boxlayout.ArrangeWindows(&boxlayout.Box{Window: "limit"}, 0, 0, width, height)

	return MergeMaps(layerOneWindows, limitWindows)
}

func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := map[K]V{}
	for _, currMap := range maps {
		for key, value := range currMap {
			result[key] = value
		}
	}

	return result
}

func (self *WindowArranger) mainSectionChildren() []*boxlayout.Box {
	currentWindow := self.gui.helpers.Window.CurrentWindow()

	// if we're not in split mode we can just show the one main panel. Likewise if
	// the main panel is focused and we're in full-screen mode
	if !self.gui.isMainPanelSplit() || (self.gui.State.ScreenMode == SCREEN_FULL && currentWindow == "main") {
		return []*boxlayout.Box{
			{
				Window: "main",
				Weight: 1,
			},
		}
	}

	return []*boxlayout.Box{
		{
			Window: "main",
			Weight: 1,
		},
		{
			Window: "secondary",
			Weight: 1,
		},
	}
}

func (self *WindowArranger) getMidSectionWeights() (int, int) {
	currentWindow := self.gui.helpers.Window.CurrentWindow()

	// we originally specified this as a ratio i.e. .20 would correspond to a weight of 1 against 4
	sidePanelWidthRatio := self.gui.c.UserConfig.Gui.SidePanelWidth
	// we could make this better by creating ratios like 2:3 rather than always 1:something
	mainSectionWeight := int(1/sidePanelWidthRatio) - 1
	sideSectionWeight := 1

	if self.splitMainPanelSideBySide() {
		mainSectionWeight = 5 // need to shrink side panel to make way for main panels if side-by-side
	}

	if currentWindow == "main" {
		if self.gui.State.ScreenMode == SCREEN_HALF || self.gui.State.ScreenMode == SCREEN_FULL {
			sideSectionWeight = 0
		}
	} else {
		if self.gui.State.ScreenMode == SCREEN_HALF {
			mainSectionWeight = 1
		} else if self.gui.State.ScreenMode == SCREEN_FULL {
			mainSectionWeight = 0
		}
	}

	return sideSectionWeight, mainSectionWeight
}

func (self *WindowArranger) infoSectionChildren(informationStr string, appStatus string) []*boxlayout.Box {
	if self.gui.State.Searching.isSearching {
		return []*boxlayout.Box{
			{
				Window: "searchPrefix",
				Size:   runewidth.StringWidth(SEARCH_PREFIX),
			},
			{
				Window: "search",
				Weight: 1,
			},
		}
	}

	appStatusBox := &boxlayout.Box{Window: "appStatus"}
	optionsBox := &boxlayout.Box{Window: "options"}

	if !self.gui.c.UserConfig.Gui.ShowBottomLine {
		optionsBox.Weight = 0
		appStatusBox.Weight = 1
	} else {
		optionsBox.Weight = 1
		appStatusBox.Size = runewidth.StringWidth(INFO_SECTION_PADDING) + runewidth.StringWidth(appStatus)
	}

	result := []*boxlayout.Box{appStatusBox, optionsBox}

	if self.gui.c.UserConfig.Gui.ShowBottomLine || self.gui.isAnyModeActive() {
		result = append(result, &boxlayout.Box{
			Window: "information",
			// unlike appStatus, informationStr has various colors so we need to decolorise before taking the length
			Size: runewidth.StringWidth(INFO_SECTION_PADDING) + runewidth.StringWidth(utils.Decolorise(informationStr)),
		})
	}

	return result
}

func (self *WindowArranger) splitMainPanelSideBySide() bool {
	if !self.gui.isMainPanelSplit() {
		return false
	}

	mainPanelSplitMode := self.gui.c.UserConfig.Gui.MainPanelSplitMode
	width, height := self.gui.g.Size()

	switch mainPanelSplitMode {
	case "vertical":
		return false
	case "horizontal":
		return true
	default:
		if width < 200 && height > 30 { // 2 80 character width panels + 40 width for side panel
			return false
		} else {
			return true
		}
	}
}

func (self *WindowArranger) getExtrasWindowSize(screenHeight int) int {
	if !self.gui.ShowExtrasWindow {
		return 0
	}

	var baseSize int
	if self.gui.c.CurrentStaticContext().GetKey() == context.COMMAND_LOG_CONTEXT_KEY {
		baseSize = 1000 // my way of saying 'fill the available space'
	} else if screenHeight < 40 {
		baseSize = 1
	} else {
		baseSize = self.gui.c.UserConfig.Gui.CommandLogSize
	}

	frameSize := 2
	return baseSize + frameSize
}

// The stash window by default only contains one line so that it's not hogging
// too much space, but if you access it it should take up some space. This is
// the default behaviour when accordion mode is NOT in effect. If it is in effect
// then when it's accessed it will have weight 2, not 1.
func (self *WindowArranger) getDefaultStashWindowBox() *boxlayout.Box {
	self.gui.State.ContextMgr.RLock()
	defer self.gui.State.ContextMgr.RUnlock()

	box := &boxlayout.Box{Window: "stash"}
	stashWindowAccessed := false
	for _, context := range self.gui.State.ContextMgr.ContextStack {
		if context.GetWindowName() == "stash" {
			stashWindowAccessed = true
		}
	}
	// if the stash window is anywhere in our stack we should enlargen it
	if stashWindowAccessed {
		box.Weight = 1
	} else {
		box.Size = 3
	}

	return box
}

func (self *WindowArranger) sidePanelChildren(width int, height int) []*boxlayout.Box {
	currentWindow := self.currentSideWindowName()

	if self.gui.State.ScreenMode == SCREEN_FULL || self.gui.State.ScreenMode == SCREEN_HALF {
		fullHeightBox := func(window string) *boxlayout.Box {
			if window == currentWindow {
				return &boxlayout.Box{
					Window: window,
					Weight: 1,
				}
			} else {
				return &boxlayout.Box{
					Window: window,
					Size:   0,
				}
			}
		}

		return []*boxlayout.Box{
			fullHeightBox("status"),
			fullHeightBox("files"),
			fullHeightBox("branches"),
			fullHeightBox("commits"),
			fullHeightBox("stash"),
		}
	} else if height >= 28 {
		accordionMode := self.gui.c.UserConfig.Gui.ExpandFocusedSidePanel
		accordionBox := func(defaultBox *boxlayout.Box) *boxlayout.Box {
			if accordionMode && defaultBox.Window == currentWindow {
				return &boxlayout.Box{
					Window: defaultBox.Window,
					Weight: 2,
				}
			}

			return defaultBox
		}

		return []*boxlayout.Box{
			{
				Window: "status",
				Size:   3,
			},
			accordionBox(&boxlayout.Box{Window: "files", Weight: 1}),
			accordionBox(&boxlayout.Box{Window: "branches", Weight: 1}),
			accordionBox(&boxlayout.Box{Window: "commits", Weight: 1}),
			accordionBox(self.getDefaultStashWindowBox()),
		}
	} else {
		squashedHeight := 1
		if height >= 21 {
			squashedHeight = 3
		}

		squashedSidePanelBox := func(window string) *boxlayout.Box {
			if window == currentWindow {
				return &boxlayout.Box{
					Window: window,
					Weight: 1,
				}
			} else {
				return &boxlayout.Box{
					Window: window,
					Size:   squashedHeight,
				}
			}
		}

		return []*boxlayout.Box{
			squashedSidePanelBox("status"),
			squashedSidePanelBox("files"),
			squashedSidePanelBox("branches"),
			squashedSidePanelBox("commits"),
			squashedSidePanelBox("stash"),
		}
	}
}

func (self *WindowArranger) currentSideWindowName() string {
	// there is always one and only one cyclable context in the context stack. We'll look from top to bottom
	self.gui.State.ContextMgr.RLock()
	defer self.gui.State.ContextMgr.RUnlock()

	for idx := range self.gui.State.ContextMgr.ContextStack {
		reversedIdx := len(self.gui.State.ContextMgr.ContextStack) - 1 - idx
		context := self.gui.State.ContextMgr.ContextStack[reversedIdx]

		if context.GetKind() == types.SIDE_CONTEXT {
			return context.GetWindowName()
		}
	}

	return "files" // default
}
