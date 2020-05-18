package gui

import (
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) mainSectionChildren() []*box {
	currentViewName := gui.currentViewName()

	// if we're not in split mode we can just show the one main panel. Likewise if
	// the main panel is focused and we're in full-screen mode
	if !gui.State.SplitMainPanel || (gui.State.ScreenMode == SCREEN_FULL && currentViewName == "main") {
		return []*box{
			{
				viewName: "main",
				weight:   1,
			},
		}
	}

	main := "main"
	secondary := "secondary"
	if gui.secondaryViewFocused() {
		// when you think you've focused the secondary view, we've actually just swapped them around in the layout
		main, secondary = secondary, main
	}

	return []*box{
		{
			viewName: main,
			weight:   1,
		},
		{
			viewName: secondary,
			weight:   1,
		},
	}
}

func (gui *Gui) getMidSectionWeights() (int, int) {
	currentViewName := gui.currentViewName()

	// we originally specified this as a ratio i.e. .20 would correspond to a weight of 1 against 4
	sidePanelWidthRatio := gui.Config.GetUserConfig().GetFloat64("gui.sidePanelWidth")
	// we could make this better by creating ratios like 2:3 rather than always 1:something
	mainSectionWeight := int(1/sidePanelWidthRatio) - 1
	sideSectionWeight := 1

	if gui.State.SplitMainPanel {
		mainSectionWeight = 5 // need to shrink side panel to make way for main panels if side-by-side
	}

	if currentViewName == "main" {
		if gui.State.ScreenMode == SCREEN_HALF || gui.State.ScreenMode == SCREEN_FULL {
			sideSectionWeight = 0
		}
	} else {
		if gui.State.ScreenMode == SCREEN_HALF {
			mainSectionWeight = 1
		} else if gui.State.ScreenMode == SCREEN_FULL {
			mainSectionWeight = 0
		}
	}

	return sideSectionWeight, mainSectionWeight
}

func (gui *Gui) infoSectionChildren(informationStr string, appStatus string) []*box {
	if gui.State.Searching.isSearching {
		return []*box{
			{
				viewName: "searchPrefix",
				size:     len(SEARCH_PREFIX),
			},
			{
				viewName: "search",
				weight:   1,
			},
		}
	}

	result := []*box{}

	if len(appStatus) > 0 {
		result = append(result,
			&box{
				viewName: "appStatus",
				size:     len(appStatus) + len(INFO_SECTION_PADDING),
			},
		)
	}

	result = append(result,
		[]*box{
			{
				viewName: "options",
				weight:   1,
			},
			{
				viewName: "information",
				// unlike appStatus, informationStr has various colors so we need to decolorise before taking the length
				size: len(INFO_SECTION_PADDING) + len(utils.Decolorise(informationStr)),
			},
		}...,
	)

	return result
}

func (gui *Gui) getViewDimensions(informationStr string, appStatus string) map[string]dimensions {
	width, height := gui.g.Size()

	sideSectionWeight, mainSectionWeight := gui.getMidSectionWeights()

	sidePanelsDirection := COLUMN
	portraitMode := width <= 84 && height > 50
	if portraitMode {
		sidePanelsDirection = ROW
	}

	root := &box{
		direction: ROW,
		children: []*box{
			{
				direction: sidePanelsDirection,
				weight:    1,
				children: []*box{
					{
						direction:           ROW,
						weight:              sideSectionWeight,
						conditionalChildren: gui.sidePanelChildren,
					},
					{
						conditionalDirection: func(width int, height int) int {
							if width < 160 && height > 30 { // 2 80 character width panels
								return ROW
							} else {
								return COLUMN
							}
						},
						direction: COLUMN,
						weight:    mainSectionWeight,
						children:  gui.mainSectionChildren(),
					},
				},
			},
			{
				direction: COLUMN,
				size:      1,
				children:  gui.infoSectionChildren(informationStr, appStatus),
			},
		},
	}

	return gui.arrangeViews(root, 0, 0, width, height)
}

func (gui *Gui) sidePanelChildren(width int, height int) []*box {
	currentCyclableViewName := gui.currentCyclableViewName()

	if gui.State.ScreenMode == SCREEN_FULL || gui.State.ScreenMode == SCREEN_HALF {
		fullHeightBox := func(viewName string) *box {
			if viewName == currentCyclableViewName {
				return &box{
					viewName: viewName,
					weight:   1,
				}
			} else {
				return &box{
					viewName: viewName,
					size:     0,
				}
			}
		}

		return []*box{
			fullHeightBox("status"),
			fullHeightBox("files"),
			fullHeightBox("branches"),
			fullHeightBox("commits"),
			fullHeightBox("stash"),
		}
	} else if height >= 28 {
		return []*box{
			{
				viewName: "status",
				size:     3,
			},
			{
				viewName: "files",
				weight:   1,
			},
			{
				viewName: "branches",
				weight:   1,
			},
			{
				viewName: "commits",
				weight:   1,
			},
			{
				viewName: "stash",
				size:     3,
			},
		}
	} else {
		squashedHeight := 1
		if height >= 21 {
			squashedHeight = 3
		}

		squashedSidePanelBox := func(viewName string) *box {
			if viewName == currentCyclableViewName {
				return &box{
					viewName: viewName,
					weight:   1,
				}
			} else {
				return &box{
					viewName: viewName,
					size:     squashedHeight,
				}
			}
		}

		return []*box{
			squashedSidePanelBox("status"),
			squashedSidePanelBox("files"),
			squashedSidePanelBox("branches"),
			squashedSidePanelBox("commits"),
			squashedSidePanelBox("stash"),
		}
	}
}

func (gui *Gui) currentCyclableViewName() string {
	currView := gui.g.CurrentView()
	currentCyclebleView := gui.State.PreviousView
	if currView != nil {
		viewName := currView.Name()
		usePreviousView := true
		for _, view := range gui.getCyclableViews() {
			if view == viewName {
				currentCyclebleView = viewName
				usePreviousView = false
				break
			}
		}
		if usePreviousView {
			currentCyclebleView = gui.State.PreviousView
		}
	}

	// unfortunate result of the fact that these are separate views, have to map explicitly
	if currentCyclebleView == "commitFiles" {
		return "commits"
	}

	return currentCyclebleView
}
