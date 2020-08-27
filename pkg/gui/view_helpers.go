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
const (
	COMMITS = iota
	BRANCHES
	FILES
	STASH
	REFLOG
	TAGS
	REMOTES
	STATUS
)

const (
	SYNC     = iota // wait until everything is done before returning
	ASYNC           // return immediately, allowing each independent thing to update itself
	BLOCK_UI        // wrap code in an update call to ensure UI updates all at once and keybindings aren't executed till complete
)

type refreshOptions struct {
	then  func()
	scope []int // e.g. []int{COMMITS, BRANCHES}. Leave empty to refresh everything
	mode  int   // one of SYNC (default), ASYNC, and BLOCK_UI
}

func intArrToMap(arr []int) map[int]bool {
	output := map[int]bool{}
	for _, el := range arr {
		output[el] = true
	}
	return output
}

func (gui *Gui) refreshSidePanels(options refreshOptions) error {
	wg := sync.WaitGroup{}

	f := func() {
		var scopeMap map[int]bool
		if len(options.scope) == 0 {
			scopeMap = intArrToMap([]int{COMMITS, BRANCHES, FILES, STASH, REFLOG, TAGS, REMOTES, STATUS})
		} else {
			scopeMap = intArrToMap(options.scope)
		}

		if scopeMap[COMMITS] || scopeMap[BRANCHES] || scopeMap[REFLOG] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go gui.refreshCommits()
				} else {
					gui.refreshCommits()
				}
				wg.Done()
			}()
		}

		if scopeMap[FILES] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go gui.refreshFiles()
				} else {
					gui.refreshFiles()
				}
				wg.Done()
			}()
		}

		if scopeMap[STASH] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go gui.refreshStashEntries()
				} else {
					gui.refreshStashEntries()
				}
				wg.Done()
			}()
		}

		if scopeMap[TAGS] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go gui.refreshTags()
				} else {
					gui.refreshTags()
				}
				wg.Done()
			}()
		}

		if scopeMap[REMOTES] {
			wg.Add(1)
			func() {
				if options.mode == ASYNC {
					go gui.refreshRemotes()
				} else {
					gui.refreshRemotes()
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
func (gui *Gui) renderString(viewName, s string) {
	gui.g.Update(func(*gocui.Gui) error {
		return gui.renderStringSync(viewName, s)
	})
}

func (gui *Gui) renderStringSync(viewName, s string) error {
	v, err := gui.g.View(viewName)
	if err != nil {
		return nil // return gracefully if view has been deleted
	}
	if err := v.SetOrigin(0, 0); err != nil {
		return err
	}
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}
	gui.setViewContent(v, s)
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
	gui.renderString("options", gui.optionsMapToString(optionsMap))
}

// TODO: refactor properly
// i'm so sorry but had to add this getBranchesView
func (gui *Gui) getFilesView() *gocui.View {
	v, _ := gui.g.View("files")
	return v
}

func (gui *Gui) getCommitsView() *gocui.View {
	v, _ := gui.g.View("commits")
	return v
}

func (gui *Gui) getCommitMessageView() *gocui.View {
	v, _ := gui.g.View("commitMessage")
	return v
}

func (gui *Gui) getBranchesView() *gocui.View {
	v, _ := gui.g.View("branches")
	return v
}

func (gui *Gui) getMainView() *gocui.View {
	v, _ := gui.g.View("main")
	return v
}

func (gui *Gui) getSecondaryView() *gocui.View {
	v, _ := gui.g.View("secondary")
	return v
}

func (gui *Gui) getStashView() *gocui.View {
	v, _ := gui.g.View("stash")
	return v
}

func (gui *Gui) getCommitFilesView() *gocui.View {
	v, _ := gui.g.View("commitFiles")
	return v
}

func (gui *Gui) getMenuView() *gocui.View {
	v, _ := gui.g.View("menu")
	return v
}

func (gui *Gui) getSearchView() *gocui.View {
	v, _ := gui.g.View("search")
	return v
}

func (gui *Gui) getStatusView() *gocui.View {
	v, _ := gui.g.View("status")
	return v
}

func (gui *Gui) getConfirmationView() *gocui.View {
	v, _ := gui.g.View("confirmation")
	return v
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
	gui.Log.Info(gui.Tr.SLocalize("resizingPopupPanel"))
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
	return map[string]string{
		fmt.Sprintf("%s/%s", gui.getKeyDisplay("universal.scrollUpMain"), gui.getKeyDisplay("universal.scrollDownMain")):                                                                                 gui.Tr.SLocalize("scroll"),
		fmt.Sprintf("%s %s %s %s", gui.getKeyDisplay("universal.prevBlock"), gui.getKeyDisplay("universal.nextBlock"), gui.getKeyDisplay("universal.prevItem"), gui.getKeyDisplay("universal.nextItem")): gui.Tr.SLocalize("navigate"),
		gui.getKeyDisplay("universal.return"):     gui.Tr.SLocalize("cancel"),
		gui.getKeyDisplay("universal.quit"):       gui.Tr.SLocalize("quit"),
		gui.getKeyDisplay("universal.optionMenu"): gui.Tr.SLocalize("menu"),
		"1-5": gui.Tr.SLocalize("jump"),
	}
}

func (gui *Gui) isPopupPanel(viewName string) bool {
	return viewName == "commitMessage" || viewName == "credentials" || viewName == "confirmation" || viewName == "menu"
}

func (gui *Gui) popupPanelFocused() bool {
	return gui.isPopupPanel(gui.currentViewName())
}

// often gocui wants functions in the form `func(g *gocui.Gui, v *gocui.View) error`
// but sometimes we just have a function that returns an error, so this is a
// convenience wrapper to give gocui what it wants.
func (gui *Gui) wrappedHandler(f func() error) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f()
	}
}

// secondaryViewFocused tells us whether it appears that the secondary view is focused. The view is actually never focused for real: we just swap the main and secondary views and then you're still focused on the main view so that we can give you access to all its keybindings for free. I will probably regret this design decision soon enough.
func (gui *Gui) secondaryViewFocused() bool {
	return gui.State.Panels.LineByLine != nil && gui.State.Panels.LineByLine.SecondaryFocused
}

func (gui *Gui) clearEditorView(v *gocui.View) {
	v.Clear()
	_ = v.SetCursor(0, 0)
	_ = v.SetOrigin(0, 0)
}

func (gui *Gui) onViewTabClick(viewName string, tabIndex int) error {
	context := gui.ViewTabContextMap[viewName][tabIndex].contexts[0]

	return gui.switchContext(context)
}

func (gui *Gui) handleNextTab(g *gocui.Gui, v *gocui.View) error {
	return gui.onViewTabClick(
		v.Name(),
		utils.ModuloWithWrap(v.TabIndex+1, len(v.Tabs)),
	)
}

func (gui *Gui) handlePrevTab(g *gocui.Gui, v *gocui.View) error {
	return gui.onViewTabClick(
		v.Name(),
		utils.ModuloWithWrap(v.TabIndex-1, len(v.Tabs)),
	)
}
