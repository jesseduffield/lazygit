package gui

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

func (gui *Gui) getCyclableWindows() []string {
	return []string{"status", "files", "branches", "commits", "stash"}
}

// models/views that we can refresh
type RefreshableView int

const (
	COMMITS RefreshableView = iota
	BRANCHES
	FILES
	STASH
	REFLOG
	TAGS
	REMOTES
	STATUS
	SUBMODULES
)

func getScopeNames(scopes []RefreshableView) []string {
	scopeNameMap := map[RefreshableView]string{
		COMMITS:    "commits",
		BRANCHES:   "branches",
		FILES:      "files",
		SUBMODULES: "submodules",
		STASH:      "stash",
		REFLOG:     "reflog",
		TAGS:       "tags",
		REMOTES:    "remotes",
		STATUS:     "status",
	}

	scopeNames := make([]string, len(scopes))
	for i, scope := range scopes {
		scopeNames[i] = scopeNameMap[scope]
	}

	return scopeNames
}

func getModeName(mode RefreshMode) string {
	switch mode {
	case SYNC:
		return "sync"
	case ASYNC:
		return "async"
	case BLOCK_UI:
		return "block-ui"
	default:
		return "unknown mode"
	}
}

type RefreshMode int

const (
	SYNC     RefreshMode = iota // wait until everything is done before returning
	ASYNC                       // return immediately, allowing each independent thing to update itself
	BLOCK_UI                    // wrap code in an update call to ensure UI updates all at once and keybindings aren't executed till complete
)

type refreshOptions struct {
	then  func()
	scope []RefreshableView // e.g. []int{COMMITS, BRANCHES}. Leave empty to refresh everything
	mode  RefreshMode       // one of SYNC (default), ASYNC, and BLOCK_UI
}

func arrToMap(arr []RefreshableView) map[RefreshableView]bool {
	output := map[RefreshableView]bool{}
	for _, el := range arr {
		output[el] = true
	}
	return output
}

func (gui *Gui) refreshSidePanels(options refreshOptions) error {
	if options.scope == nil {
		gui.Log.Infof(
			"refreshing all scopes in %s mode",
			getModeName(options.mode),
		)
	} else {
		gui.Log.Infof(
			"refreshing the following scopes in %s mode: %s",
			getModeName(options.mode),
			strings.Join(getScopeNames(options.scope), ","),
		)
	}

	wg := sync.WaitGroup{}

	f := func() {
		var scopeMap map[RefreshableView]bool
		if len(options.scope) == 0 {
			scopeMap = arrToMap([]RefreshableView{COMMITS, BRANCHES, FILES, STASH, REFLOG, TAGS, REMOTES, STATUS})
		} else {
			scopeMap = arrToMap(options.scope)
		}

		if scopeMap[COMMITS] || scopeMap[BRANCHES] || scopeMap[REFLOG] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go utils.Safe(func() { _ = gui.refreshCommits() })
				} else {
					_ = gui.refreshCommits()
				}
				wg.Done()
			}()
		}

		if scopeMap[FILES] || scopeMap[SUBMODULES] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go utils.Safe(func() { _ = gui.refreshFilesAndSubmodules() })
				} else {
					_ = gui.refreshFilesAndSubmodules()
				}
				wg.Done()
			}()
		}

		if scopeMap[STASH] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go utils.Safe(func() { _ = gui.refreshStashEntries() })
				} else {
					_ = gui.refreshStashEntries()
				}
				wg.Done()
			}()
		}

		if scopeMap[TAGS] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go utils.Safe(func() { _ = gui.refreshTags() })
				} else {
					_ = gui.refreshTags()
				}
				wg.Done()
			}()
		}

		if scopeMap[REMOTES] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go utils.Safe(func() { _ = gui.refreshRemotes() })
				} else {
					_ = gui.refreshRemotes()
				}
				wg.Done()
			}()
		}

		wg.Wait()

		gui.refreshStatus()

		if options.then != nil {
			options.then()
		}
	}

	if options.mode == BLOCK_UI {
		gui.g.Update(func(g *gocui.Gui) error {
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
	v.Clear()
	fmt.Fprint(v, gui.cleanString(s))
}

// renderString resets the origin of a view and sets its content
func (gui *Gui) renderString(view *gocui.View, s string) {
	gui.g.Update(func(*gocui.Gui) error {
		return gui.renderStringSync(view, s)
	})
}

func (gui *Gui) renderStringSync(view *gocui.View, s string) error {
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
	gui.renderString(gui.Views.Options, gui.optionsMapToString(optionsMap))
}

func (gui *Gui) trimmedContent(v *gocui.View) string {
	return strings.TrimSpace(v.Buffer())
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
		return gui.resizePopupPanel(v)
	}
	return nil
}

func (gui *Gui) resizePopupPanel(v *gocui.View) error {
	// If the confirmation panel is already displayed, just resize the width,
	// otherwise continue
	content := v.Buffer()
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(v.Wrap, content)
	vx0, vy0, vx1, vy1 := v.Dimensions()
	if vx0 == x0 && vy0 == y0 && vx1 == x1 && vy1 == y1 {
		return nil
	}
	_, err := gui.g.SetView(v.Name(), x0, y0, x1, y1, 0)
	return err
}

func (gui *Gui) changeSelectedLine(panelState IListPanelState, total int, change int) {
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

func (gui *Gui) refreshSelectedLine(panelState IListPanelState, total int) {
	line := panelState.GetSelectedLineIdx()

	if line == -1 && total > 0 {
		panelState.SetSelectedLineIdx(0)
	} else if total-1 < line {
		panelState.SetSelectedLineIdx(total - 1)
	}
}

func (gui *Gui) renderDisplayStrings(v *gocui.View, displayStrings [][]string) {
	gui.g.Update(func(g *gocui.Gui) error {
		list := utils.RenderDisplayStrings(displayStrings)
		v.Clear()
		fmt.Fprint(v, list)
		return nil
	})
}

func (gui *Gui) globalOptionsMap() map[string]string {
	keybindingConfig := gui.Config.GetUserConfig().Keybinding

	return map[string]string{
		fmt.Sprintf("%s/%s", gui.getKeyDisplay(keybindingConfig.Universal.ScrollUpMain), gui.getKeyDisplay(keybindingConfig.Universal.ScrollDownMain)):                                                                                                               gui.Tr.LcScroll,
		fmt.Sprintf("%s %s %s %s", gui.getKeyDisplay(keybindingConfig.Universal.PrevBlock), gui.getKeyDisplay(keybindingConfig.Universal.NextBlock), gui.getKeyDisplay(keybindingConfig.Universal.PrevItem), gui.getKeyDisplay(keybindingConfig.Universal.NextItem)): gui.Tr.LcNavigate,
		gui.getKeyDisplay(keybindingConfig.Universal.Return):     gui.Tr.LcCancel,
		gui.getKeyDisplay(keybindingConfig.Universal.Quit):       gui.Tr.LcQuit,
		gui.getKeyDisplay(keybindingConfig.Universal.OptionMenu): gui.Tr.LcMenu,
		"1-5": gui.Tr.LcJump,
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

func (gui *Gui) clearEditorView(v *gocui.View) {
	v.Clear()
	_ = v.SetCursor(0, 0)
	_ = v.SetOrigin(0, 0)
}

func (gui *Gui) onViewTabClick(viewName string, tabIndex int) error {
	context := gui.State.ViewTabContextMap[viewName][tabIndex].contexts[0]

	return gui.pushContext(context)
}

func (gui *Gui) handleNextTab() error {
	v := gui.g.CurrentView()
	if v == nil {
		return nil
	}

	return gui.onViewTabClick(
		v.Name(),
		utils.ModuloWithWrap(v.TabIndex+1, len(v.Tabs)),
	)
}

func (gui *Gui) handlePrevTab() error {
	v := gui.g.CurrentView()
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
