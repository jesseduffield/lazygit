package gui

import (
	"github.com/jesseduffield/lazygit/pkg/gui/boxlayout"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (gui *Gui) mainSectionChildren() []*boxlayout.Box {
	currentViewName := gui.currentViewName()

	// if we're not in split mode we can just show the one main panel. Likewise if
	// the main panel is focused and we're in full-screen mode
	if !gui.State.SplitMainPanel || (gui.State.ScreenMode == SCREEN_FULL && currentViewName == "main") {
		return []*boxlayout.Box{
			{
				ViewName: "main",
				Weight:   1,
			},
		}
	}

	main := "main"
	secondary := "secondary"
	if gui.secondaryViewFocused() {
		// when you think you've focused the secondary view, we've actually just swapped them around in the layout
		main, secondary = secondary, main
	}

	return []*boxlayout.Box{
		{
			ViewName: main,
			Weight:   1,
		},
		{
			ViewName: secondary,
			Weight:   1,
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

func (gui *Gui) infoSectionChildren(informationStr string, appStatus string) []*boxlayout.Box {
	if gui.State.Searching.isSearching {
		return []*boxlayout.Box{
			{
				ViewName: "searchPrefix",
				Size:     len(SEARCH_PREFIX),
			},
			{
				ViewName: "search",
				Weight:   1,
			},
		}
	}

	result := []*boxlayout.Box{}

	if len(appStatus) > 0 {
		result = append(result,
			&boxlayout.Box{
				ViewName: "appStatus",
				Size:     len(appStatus) + len(INFO_SECTION_PADDING),
			},
		)
	}

	result = append(result,
		[]*boxlayout.Box{
			{
				ViewName: "options",
				Weight:   1,
			},
			{
				ViewName: "information",
				// unlike appStatus, informationStr has various colors so we need to decolorise before taking the length
				Size: len(INFO_SECTION_PADDING) + len(utils.Decolorise(informationStr)),
			},
		}...,
	)

	return result
}

func (gui *Gui) getViewDimensions(informationStr string, appStatus string) map[string]boxlayout.Dimensions {
	width, height := gui.g.Size()

	sideSectionWeight, mainSectionWeight := gui.getMidSectionWeights()

	sidePanelsDirection := boxlayout.COLUMN
	portraitMode := width <= 84 && height > 45
	if portraitMode {
		sidePanelsDirection = boxlayout.ROW
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
						ConditionalChildren: gui.sidePanelChildren,
					},
					{
						ConditionalDirection: func(width int, height int) int {
							mainPanelSplitMode := gui.Config.GetUserConfig().GetString("gui.mainPanelSplitMode")

							switch mainPanelSplitMode {
							case "vertical":
								return boxlayout.ROW
							case "horizontal":
								return boxlayout.COLUMN
							default:
								if width < 160 && height > 30 { // 2 80 character width panels
									return boxlayout.ROW
								} else {
									return boxlayout.COLUMN
								}
							}
						},
						Direction: boxlayout.COLUMN,
						Weight:    mainSectionWeight,
						Children:  gui.mainSectionChildren(),
					},
				},
			},
			{
				Direction: boxlayout.COLUMN,
				Size:      1,
				Children:  gui.infoSectionChildren(informationStr, appStatus),
			},
		},
	}

	return boxlayout.ArrangeViews(root, 0, 0, width, height)
}

func (gui *Gui) sidePanelChildren(width int, height int) []*boxlayout.Box {
	currentCyclableViewName := gui.currentCyclableViewName()

	if gui.State.ScreenMode == SCREEN_FULL || gui.State.ScreenMode == SCREEN_HALF {
		fullHeightBox := func(viewName string) *boxlayout.Box {
			if viewName == currentCyclableViewName {
				return &boxlayout.Box{
					ViewName: viewName,
					Weight:   1,
				}
			} else {
				return &boxlayout.Box{
					ViewName: viewName,
					Size:     0,
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
		accordianMode := gui.Config.GetUserConfig().GetBool("gui.expandFocusedSidePanel")
		accordianBox := func(defaultBox *boxlayout.Box) *boxlayout.Box {
			if accordianMode && defaultBox.ViewName == currentCyclableViewName {
				return &boxlayout.Box{
					ViewName: defaultBox.ViewName,
					Weight:   2,
				}
			}

			return defaultBox
		}

		return []*boxlayout.Box{
			{
				ViewName: "status",
				Size:     3,
			},
			accordianBox(&boxlayout.Box{ViewName: "files", Weight: 1}),
			accordianBox(&boxlayout.Box{ViewName: "branches", Weight: 1}),
			accordianBox(&boxlayout.Box{ViewName: "commits", Weight: 1}),
			accordianBox(&boxlayout.Box{ViewName: "stash", Size: 3}),
		}
	} else {
		squashedHeight := 1
		if height >= 21 {
			squashedHeight = 3
		}

		squashedSidePanelBox := func(viewName string) *boxlayout.Box {
			if viewName == currentCyclableViewName {
				return &boxlayout.Box{
					ViewName: viewName,
					Weight:   1,
				}
			} else {
				return &boxlayout.Box{
					ViewName: viewName,
					Size:     squashedHeight,
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

func (gui *Gui) currentCyclableViewName() string {
	// there is always a cyclable context in the context stack. We'll look from top to bottom
	for idx := range gui.State.ContextStack {
		reversedIdx := len(gui.State.ContextStack) - 1 - idx
		context := gui.State.ContextStack[reversedIdx]

		if context.GetKind() == SIDE_CONTEXT {
			viewName := context.GetViewName()

			// unfortunate result of the fact that these are separate views, have to map explicitly
			if viewName == "commitFiles" {
				return "commits"
			}

			return viewName
		}
	}

	return "files" // default
}
