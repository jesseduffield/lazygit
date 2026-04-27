package controllers

import (
	"log"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

// CollapsibleSideWindows are the side windows that can be collapsed.
// Status (1) and stash (5) are excluded.
var CollapsibleSideWindows = []string{"files", "branches", "commits"}

type CollapseSideWindowController struct {
	baseController
	c *ControllerCommon
}

func NewCollapseSideWindowController(
	c *ControllerCommon,
) *CollapseSideWindowController {
	return &CollapseSideWindowController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *CollapseSideWindowController) Context() types.Context {
	return nil
}

func (self *CollapseSideWindowController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	windows := self.c.Helpers().Window.SideWindows()

	if len(opts.Config.Universal.JumpToBlock) != len(windows) {
		log.Fatal("Jump to block keybindings cannot be set. Exactly 5 keybindings must be supplied.")
	}

	bindings := make([]*types.Binding, 0, len(CollapsibleSideWindows))
	for index, window := range windows {
		if !slices.Contains(CollapsibleSideWindows, window) {
			continue
		}
		bindings = append(bindings, &types.Binding{
			ViewName:    "",
			Description: self.c.Tr.CollapseSidePanel,
			Tooltip:     "Toggle collapsing the side panel to a minimal size. The collapsed state is persisted across restarts.",
			Key:         opts.GetKey(opts.Config.Universal.JumpToBlock[index]),
			Modifier:    gocui.ModAlt,
			Handler:     opts.Guards.NoPopupPanel(self.toggleCollapse(window)),
		})
	}
	return bindings
}

func (self *CollapseSideWindowController) toggleCollapse(window string) func() error {
	return func() error {
		appState := self.c.GetAppState()
		isCollapsed := slices.Contains(appState.CollapsedSideWindows, window)

		if isCollapsed {
			// Uncollapse the window and focus it
			appState.CollapsedSideWindows = lo.Filter(appState.CollapsedSideWindows, func(w string, _ int) bool {
				return w != window
			})
			self.c.SaveAppStateAndLogError()

			context := self.c.Helpers().Window.GetContextForWindow(window)
			self.c.Context().Push(context, types.OnFocusOpts{})
		} else {
			// Collapse the window
			appState.CollapsedSideWindows = append(appState.CollapsedSideWindows, window)
			self.c.SaveAppStateAndLogError()

			// If the collapsed window was focused, move focus to the nearest non-collapsed side window
			if self.c.Helpers().Window.CurrentWindow() == window {
				allSideWindows := self.c.Helpers().Window.SideWindows()
				for _, candidate := range allSideWindows {
					if candidate != window && !slices.Contains(appState.CollapsedSideWindows, candidate) {
						context := self.c.Helpers().Window.GetContextForWindow(candidate)
						self.c.Context().Push(context, types.OnFocusOpts{})
						break
					}
				}
			}
		}
		return nil
	}
}
