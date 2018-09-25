package gui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

var cyclableViews = []string{"status", "files", "branches", "commits", "stash"}

func (gui *Gui) refreshSidePanels(g *gocui.Gui) error {
	gui.refreshBranches(g)
	gui.refreshFiles(g)
	gui.refreshCommits(g)
	return nil
}

func (gui *Gui) nextView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == cyclableViews[len(cyclableViews)-1] {
		focusedViewName = cyclableViews[0]
	} else {
		for i := range cyclableViews {
			if v.Name() == cyclableViews[i] {
				focusedViewName = cyclableViews[i+1]
				break
			}
			if i == len(cyclableViews)-1 {
				message := gui.Tr.TemplateLocalize(
					"IssntListOfViews",
					Teml{
						"name": v.Name(),
					},
				)
				gui.Log.Info(message)
				return nil
			}
		}
	}
	focusedView, err := g.View(focusedViewName)
	if err != nil {
		panic(err)
	}
	return gui.switchFocus(g, v, focusedView)
}

func (gui *Gui) previousView(g *gocui.Gui, v *gocui.View) error {
	var focusedViewName string
	if v == nil || v.Name() == cyclableViews[0] {
		focusedViewName = cyclableViews[len(cyclableViews)-1]
	} else {
		for i := range cyclableViews {
			if v.Name() == cyclableViews[i] {
				focusedViewName = cyclableViews[i-1] // TODO: make this work properly
				break
			}
			if i == len(cyclableViews)-1 {
				message := gui.Tr.TemplateLocalize(
					"IssntListOfViews",
					Teml{
						"name": v.Name(),
					},
				)
				gui.Log.Info(message)
				return nil
			}
		}
	}
	focusedView, err := g.View(focusedViewName)
	if err != nil {
		panic(err)
	}
	return gui.switchFocus(g, v, focusedView)
}

func (gui *Gui) newLineFocused(g *gocui.Gui, v *gocui.View) error {
	mainView, _ := g.View("main")
	mainView.SetOrigin(0, 0)

	switch v.Name() {
	case "menu":
		return gui.handleMenuSelect(g, v)
	case "status":
		return gui.handleStatusSelect(g, v)
	case "files":
		return gui.handleFileSelect(g, v)
	case "branches":
		return gui.handleBranchSelect(g, v)
	case "confirmation":
		return nil
	case "commitMessage":
		return gui.handleCommitFocused(g, v)
	case "main":
		// TODO: pull this out into a 'view focused' function
		gui.refreshMergePanel(g)
		v.Highlight = false
		return nil
	case "commits":
		return gui.handleCommitSelect(g, v)
	case "stash":
		return gui.handleStashEntrySelect(g, v)
	default:
		panic(gui.Tr.SLocalize("NoViewMachingNewLineFocusedSwitchStatement"))
	}
}

func (gui *Gui) returnFocus(g *gocui.Gui, v *gocui.View) error {
	previousView, err := g.View(gui.State.PreviousView)
	if err != nil {
		// always fall back to files view if there's no 'previous' view stored
		previousView, err = g.View("files")
		if err != nil {
			gui.Log.Error(err)
		}
	}
	return gui.switchFocus(g, v, previousView)
}

// pass in oldView = nil if you don't want to be able to return to your old view
func (gui *Gui) switchFocus(g *gocui.Gui, oldView, newView *gocui.View) error {
	// we assume we'll never want to return focus to a confirmation panel i.e.
	// we should never stack confirmation panels
	if oldView != nil && oldView.Name() != "confirmation" {
		oldView.Highlight = false
		message := gui.Tr.TemplateLocalize(
			"settingPreviewsViewTo",
			Teml{
				"oldViewName": oldView.Name(),
			},
		)
		gui.Log.Info(message)

		// second class panels should never have focus restored to them because
		// once they lose focus they are effectively 'destroyed'
		secondClassPanels := []string{"confirmation", "menu"}
		if !utils.IncludesString(secondClassPanels, oldView.Name()) {
			gui.State.PreviousView = oldView.Name()
		}
	}

	newView.Highlight = true
	message := gui.Tr.TemplateLocalize(
		"newFocusedViewIs",
		Teml{
			"newFocusedView": newView.Name(),
		},
	)
	gui.Log.Info(message)
	if _, err := g.SetCurrentView(newView.Name()); err != nil {
		return err
	}
	g.Cursor = newView.Editable

	return gui.newLineFocused(g, newView)
}

func (gui *Gui) getItemPosition(v *gocui.View) int {
	gui.correctCursor(v)
	_, cy := v.Cursor()
	_, oy := v.Origin()
	return oy + cy
}

func (gui *Gui) cursorUp(g *gocui.Gui, v *gocui.View) error {
	// swallowing cursor movements in main
	if v == nil || v.Name() == "main" {
		return nil
	}

	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	gui.newLineFocused(g, v)
	return nil
}

func (gui *Gui) cursorDown(g *gocui.Gui, v *gocui.View) error {
	// swallowing cursor movements in main
	if v == nil || v.Name() == "main" {
		return nil
	}
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	ly := v.LinesHeight() - 1
	_, height := v.Size()
	maxY := height - 1

	// if we are at the end we just return
	if cy+oy == ly {
		return nil
	}

	var err error
	if cy < maxY {
		err = v.SetCursor(cx, cy+1)
	} else {
		err = v.SetOrigin(ox, oy+1)
	}
	if err != nil {
		return err
	}

	gui.newLineFocused(g, v)
	return nil
}

func (gui *Gui) resetOrigin(v *gocui.View) error {
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}
	return v.SetOrigin(0, 0)
}

// if the cursor down past the last item, move it to the last line
func (gui *Gui) correctCursor(v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	_, height := v.Size()
	maxY := height - 1
	ly := v.LinesHeight() - 1
	if oy+cy <= ly {
		return nil
	}
	newCy := utils.Min(ly, maxY)
	if err := v.SetCursor(cx, newCy); err != nil {
		return err
	}
	if err := v.SetOrigin(ox, ly-newCy); err != nil {
		return err
	}
	return nil
}

func (gui *Gui) renderString(g *gocui.Gui, viewName, s string) error {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(viewName)
		// just in case the view disappeared as this function was called, we'll
		// silently return if it's not found
		if err != nil {
			return nil
		}
		v.Clear()
		output := string(bom.Clean([]byte(s)))
		output = utils.NormalizeLinefeeds(output)
		fmt.Fprint(v, output)
		v.Wrap = true
		return nil
	})
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

func (gui *Gui) renderOptionsMap(g *gocui.Gui, optionsMap map[string]string) error {
	return gui.renderString(g, "options", gui.optionsMapToString(optionsMap))
}

// TODO: refactor properly
// i'm so sorry but had to add this getBranchesView
func (gui *Gui) getFilesView(g *gocui.Gui) *gocui.View {
	v, _ := g.View("files")
	return v
}

func (gui *Gui) getCommitsView(g *gocui.Gui) *gocui.View {
	v, _ := g.View("commits")
	return v
}

func (gui *Gui) getCommitMessageView(g *gocui.Gui) *gocui.View {
	v, _ := g.View("commitMessage")
	return v
}

func (gui *Gui) getBranchesView(g *gocui.Gui) *gocui.View {
	v, _ := g.View("branches")
	return v
}

func (gui *Gui) trimmedContent(v *gocui.View) string {
	return strings.TrimSpace(v.Buffer())
}

func (gui *Gui) currentViewName(g *gocui.Gui) string {
	currentView := g.CurrentView()
	return currentView.Name()
}

func (gui *Gui) resizeCurrentPopupPanel(g *gocui.Gui) error {
	v := g.CurrentView()
	if v.Name() == "commitMessage" || v.Name() == "confirmation" {
		return gui.resizePopupPanel(g, v)
	}
	return nil
}

func (gui *Gui) resizePopupPanel(g *gocui.Gui, v *gocui.View) error {
	// If the confirmation panel is already displayed, just resize the width,
	// otherwise continue
	content := v.Buffer()
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(g, content)
	vx0, vy0, vx1, vy1 := v.Dimensions()
	if vx0 == x0 && vy0 == y0 && vx1 == x1 && vy1 == y1 {
		return nil
	}
	gui.Log.Info(gui.Tr.SLocalize("resizingPopupPanel"))
	_, err := g.SetView(v.Name(), x0, y0, x1, y1, 0)
	return err
}
