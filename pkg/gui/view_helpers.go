package gui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spkg/bom"
)

// TODO implement struct idea and add a boolean cyclable or something
var cyclableViews = []string{"status", "files", "branches", "commits", "stash"}

// nextView is called when the user presses the nextView keybinding.
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
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

	focusedView, err := gui.g.View(focusedViewName)
	if err != nil {
		gui.Log.Errorf("Failed to get the focusedView at nextView: %s\n", err)
	}

	err = gui.switchFocus(v, focusedView)
	if err != nil {
		gui.Log.Errorf("Failed to switchFocus at nextView: %s\n", err)
		return err
	}

	return nil
}

// previousView is called when the user presses the previouseView keybinding.
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
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

	focusedView, err := gui.g.View(focusedViewName)
	if err != nil {
		gui.Log.Errorf("Failed to get focusedViewName at previousView: %s\n", err)
		return err
	}

	err = gui.switchFocus(v, focusedView)
	if err != nil {
		gui.Log.Errorf("Failed to switchFocus at previousView: %s\n", err)
		return err
	}

	return nil
}

// newLineFocused acts as the general switch when a user presses the up or down
// key.
// v: what view this should affect.
// returns an error when something goes wrong.
func (gui *Gui) newLineFocused(v *gocui.View) error {
	mainView, _ := gui.g.View("main")
	mainView.SetOrigin(0, 0)

	switch v.Name() {
	case "menu":
		return gui.handleMenuSelect()
	case "status":
		return gui.handleStatusSelect()
	case "files":
		return gui.handleFileSelect()
	case "branches":
		return gui.handleBranchSelect(v)
	case "confirmation":
		return nil
	case "commitMessage":
		return gui.handleCommitFocused()
	case "main":
		// TODO: pull this out into a 'view focused' function
		err := gui.refreshMergePanel()
		if err != nil {
			gui.Log.Errorf("Failed to refreshMergePanel at newLineFocused: %s\n", err)
			return err
		}
		v.Highlight = false
		return nil
	case "commits":
		return gui.handleCommitSelect()
	case "stash":
		return gui.handleStashEntrySelect()
	default:
		panic(gui.Tr.SLocalize("NoViewMachingNewLineFocusedSwitchStatement"))
	}
}

// returnFocus TODO
// v: the view to be the new "previousView".
// returns an error when something goes wrong.
func (gui *Gui) returnFocus(v *gocui.View) error {
	previousView, err := gui.g.View(gui.State.PreviousView)
	if err != nil {
		// always fall back to files view if there's no 'previous' view stored
		previousView, err = gui.g.View("files")
		if err != nil {
			gui.Log.Error(err)
		}
	}

	err = gui.switchFocus(v, previousView)
	if err != nil {
		gui.Log.Errorf("Failed to switchFocus at returnFocus: %s\n", err)
		return err
	}

	return nil
}

// switchFocus gets called when the user wants to switch the focus to a different
// view.
// oldView: what the user may go back to.
// newView: what to go to now.
// returns an error when something goes wrong.
func (gui *Gui) switchFocus(oldView, newView *gocui.View) error {
	// TODO call gui.State.PreviousView
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
		gui.State.PreviousView = oldView.Name()
	}
	newView.Highlight = true

	message := gui.Tr.TemplateLocalize(
		"newFocusedViewIs",
		Teml{
			"newFocusedView": newView.Name(),
		},
	)
	gui.Log.Info(message)

	_, err := gui.g.SetCurrentView(newView.Name())
	if err != nil {
		gui.Log.Errorf("Failed to SetCurrentView at switchFocus: %s\n", err)
		return err
	}

	gui.g.Cursor = newView.Editable

	err = gui.newLineFocused(newView)
	if err != nil {
		gui.Log.Errorf("Failed to newLineFocuses at switchFocus: %s\n", err)
		return err
	}

	return nil
}

// getItemPosition is called when a function needs to know the current position
// of the cursor.
// v: the view to check.
// returns an integer representing the position.
func (gui *Gui) getItemPosition(v *gocui.View) int {
	err := gui.correctCursor(v)
	if err != nil {
		gui.Log.Errorf("Failed to correctCursor at getItemPosition: %s\n", err)
	}

	_, cy := v.Cursor()
	_, oy := v.Origin()

	return oy + cy
}

// cursorUp is called when the user presses the cursorUp keybind,
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
func (gui *Gui) cursorUp(g *gocui.Gui, v *gocui.View) error {
	// swallowing cursor movements in main
	// TODO: pull this out
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

	err := gui.newLineFocused(v)
	if err != nil {
		gui.Log.Errorf("Failed to newLineFocused at cursorUp: %s\n", err)
		return err
	}

	return nil
}

// cursorDown is called when the user presses the cursorUp keybind,
// g and v are passed by the gocui library.
// returns an error when something goes wrong.
func (gui *Gui) cursorDown(g *gocui.Gui, v *gocui.View) error {
	// swallowing cursor movements in main
	// TODO: pull this out
	if v == nil || v.Name() == "main" {
		return nil
	}

	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	if cy+oy >= len(v.BufferLines())-2 {
		return nil
	}

	err := v.SetCursor(cx, cy+1)
	if err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	err = gui.newLineFocused(v)
	if err != nil {
		gui.Log.Errorf("Failed to newLineFocused at cursorDown: %s\n")
		return err
	}

	return nil
}

// resetOrigin resets the cursor to 0,0.
// v: what view to affect.
// returns an error when something goes wrong.
func (gui *Gui) resetOrigin(v *gocui.View) error {
	err := v.SetCursor(0, 0)
	if err != nil {
		gui.Log.Errorf("Failed to setCursor at resetOrigin: %s\n", err)
		return err
	}

	err = v.SetOrigin(0, 0)
	if err != nil {
		gui.Log.Errorf("Failed to setOrigin at resetOrigin: %s\n")
		return err
	}

	return nil
}

// correctCursor corrects the cursor if the cursor down past the last item.
// It does so by move it to the last line.
// v: the view to affect.
// returns an error when something goes wrong.
func (gui *Gui) correctCursor(v *gocui.View) error {
	cx, cy := v.Cursor()
	_, oy := v.Origin()

	lineCount := len(v.BufferLines()) - 2

	if lineCount < 0 {
		lineCount = 0
	}

	if cy >= lineCount-oy {
		err := v.SetCursor(cx, lineCount-oy)
		if err != nil {
			gui.Log.Errorf("Failed to set cursor at correctCursor: %s\n", err)
			return err
		}

		return nil
	}

	return nil
}

// renderString renders a string on the given view.
// viewName: the view to populate.
// s: the string to write.
// returns an error if something goes wrong.
func (gui *Gui) renderString(viewName, s string) error {
	gui.g.Update(func(*gocui.Gui) error {

		v, err := gui.g.View(viewName)
		if err != nil {
			gui.Log.Errorf("Failed to get view %s at renderString: %s", viewName, err)
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

// optionsMapToString parses an options map to a string.
// optionsMap: what to parse.
// returns a string representing the map.
func (gui *Gui) optionsMapToString(optionsMap map[string]string) string {
	optionsArray := make([]string, 0)
	for key, description := range optionsMap {
		optionsArray = append(optionsArray, key+": "+description)
	}

	sort.Strings(optionsArray)

	return strings.Join(optionsArray, ", ")
}

// renderOptionsMap renders the options to the options view.
// optionsMap: what to render.
// returns an error when something goes wrong.
func (gui *Gui) renderOptionsMap(optionsMap map[string]string) error {
	return gui.renderString("options", gui.optionsMapToString(optionsMap))
}

// trimmedContent trims the content.
// v: the view of which the content should be trimmed.
// returns an string representing the trimmed content.
func (gui *Gui) trimmedContent(v *gocui.View) string {
	return strings.TrimSpace(v.Buffer())
}

// resizeCurrentPopupPanel resizes the current popup panel.
// returns an error if something went wrong.
func (gui *Gui) resizeCurrentPopupPanel() error {
	v := gui.g.CurrentView()
	if v.Name() == "commitMessage" || v.Name() == "confirmation" {
		return gui.resizePopupPanel(v)
	}

	return nil
}

// resizePopupPanel resizes the popup panel.
// v: the view to resize.
// returns an error when something goes wrong.
func (gui *Gui) resizePopupPanel(v *gocui.View) error {
	content := v.Buffer()
	// If the confirmation panel is already displayed, just resize the width,
	// otherwise continue
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(content)
	vx0, vy0, vx1, vy1 := v.Dimensions()
	if vx0 == x0 && vy0 == y0 && vx1 == x1 && vy1 == y1 {
		return nil
	}

	gui.Log.Info(gui.Tr.SLocalize("resizingPopupPanel"))

	_, err := gui.g.SetView(v.Name(), x0, y0, x1, y1, 0)
	if err != nil {
		gui.Log.Errorf("Failed to get view %s at resizePopupPanel: %s\n", v.Name(), err)
		return err
	}

	return nil
}
