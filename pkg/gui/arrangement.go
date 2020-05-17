package gui

func (gui *Gui) getViewDimensions() map[string]dimensions {
	// things to consider:
	// current cyclable view
	// three modes: normal, squashed, portrait
	// three fullscreen modes: regular, half, full
	// half/full ignored squashed but not portrait where the orientation is just swapped
	// height (for squashing)
	// width (for portrait mode)
	// let's start by saying squashing and portrait mode are mutuall exclusive. If you're in portrait mode and you end up in a squashed state you're pretty much fucked either way.
	// having said that, half and fullscreen mode can be combined with the other two. Fullscreen is gonna be the same across all of the options. Half mode will just be split vertically rather than horizontally for

	// need

	// we'll start assuming normal mode
	// in normal mode pick the split between the main panel and the side panels, then give every cyclablable view roughly equal height.

	// the options panel has 1 height always

	// so our views are:
	// "status"
	// "files"
	// "branches"
	// "commits"
	// "stash"
	// "main"
	// "commitFiles"
	// "secondary"

	// would be good to have a way of describing this programmatically, like with html:

	// <box>

	width, height := gui.g.Size()

	portraitMode := width <= 84 && height > 50

	main := "main"
	secondary := "secondary"
	if gui.State.Panels.LineByLine != nil && gui.State.Panels.LineByLine.SecondaryFocused {
		main, secondary = secondary, main
	}

	mainSectionChildren := []*box{
		{
			viewName: main,
			weight:   1,
		},
	}

	if gui.State.SplitMainPanel {
		mainSectionChildren = append(mainSectionChildren, &box{
			viewName: secondary,
			weight:   1,
		})
	}

	currentCyclableViewName := gui.currentCyclableViewName()

	var sideSectionChildren []*box
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

		sideSectionChildren = []*box{
			fullHeightBox("status"),
			fullHeightBox("files"),
			fullHeightBox("branches"),
			fullHeightBox("commits"),
			fullHeightBox("stash"),
		}
	} else if height >= 28 && !portraitMode {
		sideSectionChildren = []*box{
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

		sideSectionChildren = []*box{
			squashedSidePanelBox("status"),
			squashedSidePanelBox("files"),
			squashedSidePanelBox("branches"),
			squashedSidePanelBox("commits"),
			squashedSidePanelBox("stash"),
		}
	}

	// we originally specified this as a ratio i.e. .20 would correspond to a weight of 1 against 4
	sidePanelWidthRatio := gui.Config.GetUserConfig().GetFloat64("gui.sidePanelWidth")
	// we could make this better by creating ratios like 2:3 rather than always 1:something
	mainSectionWeight := int(1/sidePanelWidthRatio) - 1
	sideSectionWeight := 1

	if gui.State.SplitMainPanel {
		mainSectionWeight = 5 // need to shrink side panel to make way for main panels if side-by-side
	}
	currentViewName := gui.currentViewName()
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

	sidePanelsDirection := COLUMN
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
						direction: ROW,
						weight:    sideSectionWeight,
						children:  sideSectionChildren,
					},
					{
						conditionalDirection: func(width int, _height int) int {
							if width < 160 && height > 30 { // 2 80 character width panels
								return ROW
							} else {
								return COLUMN
							}
						},
						direction: COLUMN,
						weight:    mainSectionWeight,
						children:  mainSectionChildren,
					},
				},
			},
			// TODO: actually handle options here. Currently we're just hard-coding it to be set on the bottom row in our layout function given that we need some custom logic to have it share space with other views on that row.
			{
				viewName: "options",
				size:     1,
			},
		},
	}

	return gui.layoutViews(root, 0, 0, width, height)
}
