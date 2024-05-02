package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazycore/pkg/boxlayout"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/mattn/go-runewidth"
	"golang.org/x/exp/slices"
)

// In this file we use the boxlayout package, along with knowledge about the app's state,
// to arrange the windows (i.e. panels) on the screen.

type WindowArrangementHelper struct {
	c               *HelperCommon
	windowHelper    *WindowHelper
	modeHelper      *ModeHelper
	appStatusHelper *AppStatusHelper
}

func NewWindowArrangementHelper(
	c *HelperCommon,
	windowHelper *WindowHelper,
	modeHelper *ModeHelper,
	appStatusHelper *AppStatusHelper,
) *WindowArrangementHelper {
	return &WindowArrangementHelper{
		c:               c,
		windowHelper:    windowHelper,
		modeHelper:      modeHelper,
		appStatusHelper: appStatusHelper,
	}
}

type WindowArrangementArgs struct {
	// Width of the screen (in characters)
	Width int
	// Height of the screen (in characters)
	Height int
	// User config
	UserConfig *config.UserConfig
	// Name of the currently focused window
	CurrentWindow string
	// Name of the current static window (meaning popups are ignored)
	CurrentStaticWindow string
	// Name of the current side window (i.e. the current window in the left
	// section of the UI)
	CurrentSideWindow string
	// Whether the main panel is split (as is the case e.g. when a file has both
	// staged and unstaged changes)
	SplitMainPanel bool
	// The current screen mode (normal, half, full)
	ScreenMode types.WindowMaximisation
	// The content shown on the bottom left of the screen when showing a loader
	// or toast e.g. 'Rebasing /'
	AppStatus string
	// The content shown on the bottom right of the screen (e.g. the 'donate',
	// 'ask question' links or a message about the current mode e.g. rebase mode)
	InformationStr string
	// Whether to show the extras window which contains the command log context
	ShowExtrasWindow bool
	// Whether we are in a demo (which is used for generating demo gifs for the
	// repo's readme)
	InDemo bool
	// Whether any mode is active (e.g. rebasing, cherry picking, etc)
	IsAnyModeActive bool
	// Whether the search prompt is shown in the bottom left
	InSearchPrompt bool
	// One of '' (not searching), 'Search: ', and 'Filter: '
	SearchPrefix string
}

func (self *WindowArrangementHelper) GetWindowDimensions(informationStr string, appStatus string) map[string]boxlayout.Dimensions {
	width, height := self.c.GocuiGui().Size()
	repoState := self.c.State().GetRepoState()

	var searchPrefix string
	if repoState.GetSearchState().SearchType() == types.SearchTypeSearch {
		searchPrefix = self.c.Tr.SearchPrefix
	} else {
		searchPrefix = self.c.Tr.FilterPrefix
	}

	args := WindowArrangementArgs{
		Width:               width,
		Height:              height,
		UserConfig:          self.c.UserConfig,
		CurrentWindow:       self.windowHelper.CurrentWindow(),
		CurrentSideWindow:   self.c.CurrentSideContext().GetWindowName(),
		CurrentStaticWindow: self.c.CurrentStaticContext().GetWindowName(),
		SplitMainPanel:      repoState.GetSplitMainPanel(),
		ScreenMode:          repoState.GetScreenMode(),
		AppStatus:           appStatus,
		InformationStr:      informationStr,
		ShowExtrasWindow:    self.c.State().GetShowExtrasWindow(),
		InDemo:              self.c.InDemo(),
		IsAnyModeActive:     self.modeHelper.IsAnyModeActive(),
		InSearchPrompt:      repoState.InSearchPrompt(),
		SearchPrefix:        searchPrefix,
	}

	return GetWindowDimensions(args)
}

func shouldUsePortraitMode(args WindowArrangementArgs) bool {
	if args.ScreenMode == types.SCREEN_HALF {
		return args.UserConfig.Gui.EnlargedSideViewLocation == "top"
	}

	switch args.UserConfig.Gui.PortraitMode {
	case "never":
		return false
	case "always":
		return true
	default: // "auto" or any garbage values in PortraitMode value
		return args.Width <= 84 && args.Height > 45
	}
}

func GetWindowDimensions(args WindowArrangementArgs) map[string]boxlayout.Dimensions {
	sideSectionWeight, mainSectionWeight := getMidSectionWeights(args)

	sidePanelsDirection := boxlayout.COLUMN
	if shouldUsePortraitMode(args) {
		sidePanelsDirection = boxlayout.ROW
	}

	showInfoSection := args.UserConfig.Gui.ShowBottomLine ||
		args.InSearchPrompt ||
		args.IsAnyModeActive ||
		args.AppStatus != ""
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
						ConditionalChildren: sidePanelChildren(args),
					},
					{
						Direction: boxlayout.ROW,
						Weight:    mainSectionWeight,
						Children:  mainPanelChildren(args),
					},
				},
			},
			{
				Direction: boxlayout.COLUMN,
				Size:      infoSectionSize,
				Children:  infoSectionChildren(args),
			},
		},
	}

	layerOneWindows := boxlayout.ArrangeWindows(root, 0, 0, args.Width, args.Height)
	limitWindows := boxlayout.ArrangeWindows(&boxlayout.Box{Window: "limit"}, 0, 0, args.Width, args.Height)

	return MergeMaps(layerOneWindows, limitWindows)
}

func mainPanelChildren(args WindowArrangementArgs) []*boxlayout.Box {
	mainPanelsDirection := boxlayout.ROW
	if splitMainPanelSideBySide(args) {
		mainPanelsDirection = boxlayout.COLUMN
	}

	result := []*boxlayout.Box{
		{
			Direction: mainPanelsDirection,
			Children:  mainSectionChildren(args),
			Weight:    1,
		},
	}
	if args.ShowExtrasWindow {
		result = append(result, &boxlayout.Box{
			Window: "extras",
			Size:   getExtrasWindowSize(args),
		})
	}
	return result
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

func mainSectionChildren(args WindowArrangementArgs) []*boxlayout.Box {
	// if we're not in split mode we can just show the one main panel. Likewise if
	// the main panel is focused and we're in full-screen mode
	if !args.SplitMainPanel || (args.ScreenMode == types.SCREEN_FULL && args.CurrentWindow == "main") {
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

func getMidSectionWeights(args WindowArrangementArgs) (int, int) {
	// we originally specified this as a ratio i.e. .20 would correspond to a weight of 1 against 4
	sidePanelWidthRatio := args.UserConfig.Gui.SidePanelWidth
	// we could make this better by creating ratios like 2:3 rather than always 1:something
	mainSectionWeight := int(1/sidePanelWidthRatio) - 1
	sideSectionWeight := 1

	if splitMainPanelSideBySide(args) {
		mainSectionWeight = 5 // need to shrink side panel to make way for main panels if side-by-side
	}

	if args.CurrentWindow == "main" {
		if args.ScreenMode == types.SCREEN_HALF || args.ScreenMode == types.SCREEN_FULL {
			sideSectionWeight = 0
		}
	} else {
		if args.ScreenMode == types.SCREEN_HALF {
			if args.UserConfig.Gui.EnlargedSideViewLocation == "top" {
				mainSectionWeight = 2
			} else {
				mainSectionWeight = 1
			}
		} else if args.ScreenMode == types.SCREEN_FULL {
			mainSectionWeight = 0
		}
	}

	return sideSectionWeight, mainSectionWeight
}

func infoSectionChildren(args WindowArrangementArgs) []*boxlayout.Box {
	if args.InSearchPrompt {
		return []*boxlayout.Box{
			{
				Window: "searchPrefix",
				Size:   runewidth.StringWidth(args.SearchPrefix),
			},
			{
				Window: "search",
				Weight: 1,
			},
		}
	}

	statusSpacerPrefix := "statusSpacer"
	spacerBoxIndex := 0
	maxSpacerBoxIndex := 2 // See pkg/gui/types/views.go
	// Returns a box with size 1 to be used as padding between views
	spacerBox := func() *boxlayout.Box {
		spacerBoxIndex++

		if spacerBoxIndex > maxSpacerBoxIndex {
			panic("Too many spacer boxes")
		}

		return &boxlayout.Box{Window: fmt.Sprintf("%s%d", statusSpacerPrefix, spacerBoxIndex), Size: 1}
	}

	// Returns a box with weight 1 to be used as flexible padding between views
	flexibleSpacerBox := func() *boxlayout.Box {
		spacerBoxIndex++

		if spacerBoxIndex > maxSpacerBoxIndex {
			panic("Too many spacer boxes")
		}

		return &boxlayout.Box{Window: fmt.Sprintf("%s%d", statusSpacerPrefix, spacerBoxIndex), Weight: 1}
	}

	// Adds spacer boxes inbetween given boxes
	insertSpacerBoxes := func(boxes []*boxlayout.Box) []*boxlayout.Box {
		for i := len(boxes) - 1; i >= 1; i-- {
			// ignore existing spacer boxes
			if !strings.HasPrefix(boxes[i].Window, statusSpacerPrefix) {
				boxes = slices.Insert(boxes, i, spacerBox())
			}
		}
		return boxes
	}

	// First collect the real views that we want to show, we'll add spacers in
	// between at the end
	var result []*boxlayout.Box

	if !args.InDemo {
		// app status appears very briefly in demos and dislodges the caption,
		// so better not to show it at all
		if args.AppStatus != "" {
			result = append(result, &boxlayout.Box{Window: "appStatus", Size: runewidth.StringWidth(args.AppStatus)})
		}
	}

	if args.UserConfig.Gui.ShowBottomLine {
		result = append(result, &boxlayout.Box{Window: "options", Weight: 1})
	}

	if (!args.InDemo && args.UserConfig.Gui.ShowBottomLine) || args.IsAnyModeActive {
		result = append(result,
			&boxlayout.Box{
				Window: "information",
				// unlike appStatus, informationStr has various colors so we need to decolorise before taking the length
				Size: runewidth.StringWidth(utils.Decolorise(args.InformationStr)),
			})
	}

	if len(result) == 2 && result[0].Window == "appStatus" {
		// Only status and information are showing; need to insert a flexible
		// spacer between the two, so that information is right-aligned. Note
		// that the call to insertSpacerBoxes below will still insert a 1-char
		// spacer in addition (right after the flexible one); this is needed for
		// the case that there's not enough room, to ensure there's always at
		// least one space.
		result = slices.Insert(result, 1, flexibleSpacerBox())
	} else if len(result) == 1 {
		if result[0].Window == "information" {
			// Only information is showing; need to add a flexible spacer so
			// that information is right-aligned
			result = slices.Insert(result, 0, flexibleSpacerBox())
		} else {
			// Only status is showing; need to make it flexible so that it
			// extends over the whole width
			result[0].Size = 0
			result[0].Weight = 1
		}
	}

	if len(result) > 0 {
		// If we have at least one view, insert 1-char wide spacer boxes between them.
		result = insertSpacerBoxes(result)
	}

	return result
}

func splitMainPanelSideBySide(args WindowArrangementArgs) bool {
	if !args.SplitMainPanel {
		return false
	}

	mainPanelSplitMode := args.UserConfig.Gui.MainPanelSplitMode
	switch mainPanelSplitMode {
	case "vertical":
		return false
	case "horizontal":
		return true
	default:
		if args.Width < 200 && args.Height > 30 { // 2 80 character width panels + 40 width for side panel
			return false
		} else {
			return true
		}
	}
}

func getExtrasWindowSize(args WindowArrangementArgs) int {
	var baseSize int
	// The 'extras' window contains the command log context
	if args.CurrentStaticWindow == "extras" {
		baseSize = 1000 // my way of saying 'fill the available space'
	} else if args.Height < 40 {
		baseSize = 1
	} else {
		baseSize = args.UserConfig.Gui.CommandLogSize
	}

	frameSize := 2
	return baseSize + frameSize
}

// The stash window by default only contains one line so that it's not hogging
// too much space, but if you access it it should take up some space. This is
// the default behaviour when accordion mode is NOT in effect. If it is in effect
// then when it's accessed it will have weight 2, not 1.
func getDefaultStashWindowBox(args WindowArrangementArgs) *boxlayout.Box {
	box := &boxlayout.Box{Window: "stash"}
	// if the stash window is anywhere in our stack we should enlargen it
	if args.CurrentSideWindow == "stash" {
		box.Weight = 1
	} else {
		box.Size = 3
	}

	return box
}

func sidePanelChildren(args WindowArrangementArgs) func(width int, height int) []*boxlayout.Box {
	return func(width int, height int) []*boxlayout.Box {
		if args.ScreenMode == types.SCREEN_FULL || args.ScreenMode == types.SCREEN_HALF {
			fullHeightBox := func(window string) *boxlayout.Box {
				if window == args.CurrentSideWindow {
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
			accordionMode := args.UserConfig.Gui.ExpandFocusedSidePanel
			accordionBox := func(defaultBox *boxlayout.Box) *boxlayout.Box {
				if accordionMode && defaultBox.Window == args.CurrentSideWindow {
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
				accordionBox(getDefaultStashWindowBox(args)),
			}
		} else {
			squashedHeight := 1
			if height >= 21 {
				squashedHeight = 3
			}

			squashedSidePanelBox := func(window string) *boxlayout.Box {
				if window == args.CurrentSideWindow {
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
}
