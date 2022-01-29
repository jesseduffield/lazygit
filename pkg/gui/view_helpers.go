package gui

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) getCyclableWindows() []string {
	return []string{"status", "files", "branches", "commits", "stash"}
}

func getScopeNames(scopes []types.RefreshableView) []string {
	scopeNameMap := map[types.RefreshableView]string{
		types.COMMITS:     "commits",
		types.BRANCHES:    "branches",
		types.FILES:       "files",
		types.SUBMODULES:  "submodules",
		types.STASH:       "stash",
		types.REFLOG:      "reflog",
		types.TAGS:        "tags",
		types.REMOTES:     "remotes",
		types.STATUS:      "status",
		types.BISECT_INFO: "bisect",
	}

	scopeNames := make([]string, len(scopes))
	for i, scope := range scopes {
		scopeNames[i] = scopeNameMap[scope]
	}

	return scopeNames
}

func getModeName(mode types.RefreshMode) string {
	switch mode {
	case types.SYNC:
		return "sync"
	case types.ASYNC:
		return "async"
	case types.BLOCK_UI:
		return "block-ui"
	default:
		return "unknown mode"
	}
}

func arrToMap(arr []types.RefreshableView) map[types.RefreshableView]bool {
	output := map[types.RefreshableView]bool{}
	for _, el := range arr {
		output[el] = true
	}
	return output
}

func (gui *Gui) Refresh(options types.RefreshOptions) error {
	if options.Scope == nil {
		gui.c.Log.Infof(
			"refreshing all scopes in %s mode",
			getModeName(options.Mode),
		)
	} else {
		gui.c.Log.Infof(
			"refreshing the following scopes in %s mode: %s",
			getModeName(options.Mode),
			strings.Join(getScopeNames(options.Scope), ","),
		)
	}

	wg := sync.WaitGroup{}

	f := func() {
		var scopeMap map[types.RefreshableView]bool
		if len(options.Scope) == 0 {
			scopeMap = arrToMap([]types.RefreshableView{
				types.COMMITS,
				types.BRANCHES,
				types.FILES,
				types.STASH,
				types.REFLOG,
				types.TAGS,
				types.REMOTES,
				types.STATUS,
				types.BISECT_INFO,
			})
		} else {
			scopeMap = arrToMap(options.Scope)
		}

		refresh := func(f func()) {
			wg.Add(1)
			func() {
				if options.Mode == types.ASYNC {
					go utils.Safe(f)
				} else {
					f()
				}
				wg.Done()
			}()
		}

		if scopeMap[types.COMMITS] || scopeMap[types.BRANCHES] || scopeMap[types.REFLOG] || scopeMap[types.BISECT_INFO] {
			refresh(gui.refreshCommits)
		} else if scopeMap[types.REBASE_COMMITS] {
			// the above block handles rebase commits so we only need to call this one
			// if we've asked specifically for rebase commits and not those other things
			refresh(func() { _ = gui.refreshRebaseCommits() })
		}

		if scopeMap[types.FILES] || scopeMap[types.SUBMODULES] {
			refresh(func() { _ = gui.refreshFilesAndSubmodules() })
		}

		if scopeMap[types.STASH] {
			refresh(func() { _ = gui.refreshStashEntries() })
		}

		if scopeMap[types.TAGS] {
			refresh(func() { _ = gui.refreshTags() })
		}

		if scopeMap[types.REMOTES] {
			refresh(func() { _ = gui.refreshRemotes() })
		}

		wg.Wait()

		gui.refreshStatus()

		if options.Then != nil {
			options.Then()
		}
	}

	if options.Mode == types.BLOCK_UI {
		gui.OnUIThread(func() error {
			f()
			return nil
		})
	} else {
		f()
	}

	return nil
}

func (gui *Gui) resetOrigin(v *gocui.View) error {
	_ = v.SetCursor(0, 0)
	return v.SetOrigin(0, 0)
}

func (gui *Gui) cleanString(s string) string {
	output := string(bom.Clean([]byte(s)))
	return utils.NormalizeLinefeeds(output)
}

func (gui *Gui) setViewContent(v *gocui.View, s string) {
	v.SetContent(gui.cleanString(s))
}

// renderString resets the origin of a view and sets its content
func (gui *Gui) renderString(view *gocui.View, s string) error {
	if err := view.SetOrigin(0, 0); err != nil {
		return err
	}
	if err := view.SetCursor(0, 0); err != nil {
		return err
	}
	gui.setViewContent(view, s)
	return nil
}

func (gui *Gui) optionsMapToString(optionsMap map[string]string) string {
	optionsArray := make([]string, 0)
	for key, description := range optionsMap {
		optionsArray = append(optionsArray, key+": "+description)
	}
	sort.Strings(optionsArray)
	return strings.Join(optionsArray, ", ")
}

func (gui *Gui) renderOptionsMap(optionsMap map[string]string) {
	_ = gui.renderString(gui.Views.Options, gui.optionsMapToString(optionsMap))
}

func (gui *Gui) currentViewName() string {
	currentView := gui.g.CurrentView()
	if currentView == nil {
		return ""
	}
	return currentView.Name()
}

func (gui *Gui) resizeCurrentPopupPanel() error {
	v := gui.g.CurrentView()
	if v == nil {
		return nil
	}
	if gui.isPopupPanel(v.Name()) {
		return gui.resizePopupPanel(v, v.Buffer())
	}
	return nil
}

func (gui *Gui) resizePopupPanel(v *gocui.View, content string) error {
	// If the confirmation panel is already displayed, just resize the width,
	// otherwise continue
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(v.Wrap, content)
	vx0, vy0, vx1, vy1 := v.Dimensions()
	if vx0 == x0 && vy0 == y0 && vx1 == x1 && vy1 == y1 {
		return nil
	}
	_, err := gui.g.SetView(v.Name(), x0, y0, x1, y1, 0)
	return err
}

func (gui *Gui) changeSelectedLine(panelState types.IListPanelState, total int, change int) {
	// TODO: find out why we're doing this
	line := panelState.GetSelectedLineIdx()

	if line == -1 {
		return
	}
	var newLine int
	if line+change < 0 {
		newLine = 0
	} else if line+change >= total {
		newLine = total - 1
	} else {
		newLine = line + change
	}

	panelState.SetSelectedLineIdx(newLine)
}

func (gui *Gui) refreshSelectedLine(panelState types.IListPanelState, total int) {
	line := panelState.GetSelectedLineIdx()

	if line == -1 && total > 0 {
		panelState.SetSelectedLineIdx(0)
	} else if total-1 < line {
		panelState.SetSelectedLineIdx(total - 1)
	}
}

func (gui *Gui) renderDisplayStrings(v *gocui.View, displayStrings [][]string) {
	list := utils.RenderDisplayStrings(displayStrings)
	v.SetContent(list)
}

func (gui *Gui) renderDisplayStringsInViewPort(v *gocui.View, displayStrings [][]string) {
	list := utils.RenderDisplayStrings(displayStrings)
	_, y := v.Origin()
	v.OverwriteLines(y, list)
}

func (gui *Gui) globalOptionsMap() map[string]string {
	keybindingConfig := gui.c.UserConfig.Keybinding

	return map[string]string{
		fmt.Sprintf("%s/%s", gui.getKeyDisplay(keybindingConfig.Universal.ScrollUpMain), gui.getKeyDisplay(keybindingConfig.Universal.ScrollDownMain)):                                                                                                               gui.c.Tr.LcScroll,
		fmt.Sprintf("%s %s %s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevBlock), gui.getKeyDisplay(keybindingConfig.Universal.NextBlock), gui.getKeyDisplay(keybindingConfig.Universal.PrevItem), gui.getKeyDisplay(keybindingConfig.Universal.NextItem)): gui.c.Tr.LcNavigate,
		gui.getKeyDisplay(keybindingConfig.Universal.Return):     gui.c.Tr.LcCancel,
		gui.getKeyDisplay(keybindingConfig.Universal.Quit):       gui.c.Tr.LcQuit,
		gui.getKeyDisplay(keybindingConfig.Universal.OptionMenu): gui.c.Tr.LcMenu,
		fmt.Sprintf("%s-%s", gui.getKeyDisplay(keybindingConfig.Universal.JumpToBlock[0]), gui.getKeyDisplay(keybindingConfig.Universal.JumpToBlock[len(keybindingConfig.Universal.JumpToBlock)-1])): gui.c.Tr.LcJump,
		fmt.Sprintf("%s/%s", gui.getKeyDisplay(keybindingConfig.Universal.ScrollLeft), gui.getKeyDisplay(keybindingConfig.Universal.ScrollRight)):                                                    gui.c.Tr.LcScrollLeftRight,
	}
}

func (gui *Gui) isPopupPanel(viewName string) bool {
	return viewName == "commitMessage" || viewName == "credentials" || viewName == "confirmation" || viewName == "menu"
}

func (gui *Gui) popupPanelFocused() bool {
	return gui.isPopupPanel(gui.currentViewName())
}

// secondaryViewFocused tells us whether it appears that the secondary view is focused. The view is actually never focused for real: we just swap the main and secondary views and then you're still focused on the main view so that we can give you access to all its keybindings for free. I will probably regret this design decision soon enough.
func (gui *Gui) secondaryViewFocused() bool {
	state := gui.State.Panels.LineByLine
	return state != nil && state.SecondaryFocused
}

func (gui *Gui) onViewTabClick(viewName string, tabIndex int) error {
	context := gui.State.ViewTabContextMap[viewName][tabIndex].Contexts[0]

	return gui.c.PushContext(context)
}

func (gui *Gui) handleNextTab() error {
	v := getTabbedView(gui)
	if v == nil {
		return nil
	}

	return gui.onViewTabClick(
		v.Name(),
		utils.ModuloWithWrap(v.TabIndex+1, len(v.Tabs)),
	)
}

func (gui *Gui) handlePrevTab() error {
	v := getTabbedView(gui)
	if v == nil {
		return nil
	}

	return gui.onViewTabClick(
		v.Name(),
		utils.ModuloWithWrap(v.TabIndex-1, len(v.Tabs)),
	)
}

// this is the distance we will move the cursor when paging up or down in a view
func (gui *Gui) pageDelta(view *gocui.View) int {
	_, height := view.Size()

	delta := height - 1
	if delta == 0 {
		return 1
	}

	return delta
}

func getTabbedView(gui *Gui) *gocui.View {
	// It safe assumption that only static contexts have tabs
	context := gui.currentStaticContext()
	view, _ := gui.g.View(context.GetViewName())
	return view
}

func (gui *Gui) render() {
	gui.OnUIThread(func() error { return nil })
}
