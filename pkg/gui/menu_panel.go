package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// handleMenuPress is called when a user presses the menu key.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleMenuPress(g *gocui.Gui, v *gocui.View) error {

	lineNumber := gui.getItemPosition(v)
	if gui.State.Keys[lineNumber].Key == nil {
		return nil
	}

	if len(gui.State.Keys) > lineNumber {

		err := gui.handleMenuClose(gui.g, v)
		if err != nil {
			gui.Log.Errorf("Failed to handleMenuClose at handleMenuPress: %s\n", err)
			return err
		}

		vv, err := gui.g.View(gui.State.Keys[lineNumber].ViewName)
		if err != nil {
			gui.Log.Errorf("Failed to get view at handleMenuPress: %s\n", err)
			return err
		}

		err = gui.State.Keys[lineNumber].Handler(gui.g, vv)
		if err != nil {
			gui.Log.Errorf("Failed to call handler for %s\n at handleMenuPress: %s\n", gui.State.Keys[lineNumber].Description, err)
			return err
		}

		return nil
	}

	return nil
}

// handleMenuSelect doesn't do anything at the moment
// but it is needed for switch in newLineFocused
// TODO find use
func (gui *Gui) handleMenuSelect(v *gocui.View) error {
	return nil
}

// renderMenuOptions renders the menu options.
// returns an error if something goes wrong.
func (gui *Gui) renderMenuOptions() error {

	optionsMap := map[string]string{
		"esc/q": gui.Tr.SLocalize("close"),
		"↑ ↓":   gui.Tr.SLocalize("navigate"),
		"space": gui.Tr.SLocalize("execute"),
	}

	err := gui.renderOptionsMap(gui.g, optionsMap)
	if err != nil {
		gui.Log.Errorf("Failed to renderOptionsMap at renderMenuOptions: %s\n", err)
		return err
	}

	return err
}

// handleMenuClose is called when the user closes the menu
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleMenuClose(g *gocui.Gui, v *gocui.View) error {

	err := gui.g.DeleteView("menu")
	if err != nil {
		return err
	}

	err = gui.returnFocus(gui.g, v)
	if err != nil {
		gui.Log.Errorf("Failed to return focus at handleMenuClose: %s\n", err)
		return err
	}

	return nil
}

// GetKey returns the string representation of the key.
// binding: what to get the key from.
// returns the key
func (gui *Gui) GetKey(binding Binding) string {
	r, ok := binding.Key.(rune)
	key := ""

	if ok {
		key = string(r)
	} else if binding.KeyReadable != "" {
		key = binding.KeyReadable
	}

	return key
}

// GetMaxKeyLength returns the size of a key from an array of bindings
// bindings: what to get the keys from.
// returns the size
func (gui *Gui) GetMaxKeyLength(bindings []Binding) int {

	max := 0

	for _, binding := range bindings {
		keyLength := len(gui.GetKey(binding))
		if keyLength > max {
			max = keyLength
		}
	}

	return max
}

// handleMenu is called when a user opens the menu.
// g and v are passed by the gocui library.
// returns an error if something goes wrong.
func (gui *Gui) handleMenu(g *gocui.Gui, v *gocui.View) error {

	var (
		contentGlobal, contentPanel   []string
		bindingsGlobal, bindingsPanel []Binding
	)

	// clear keys slice, so we don't have ghost elements
	gui.State.Keys = gui.State.Keys[:0]
	bindings := gui.GetKeybindings()
	padWidth := gui.GetMaxKeyLength(bindings)

	for _, binding := range bindings {
		key := gui.GetKey(binding)
		if key != "" && binding.Description != "" {
			content := fmt.Sprintf("%s  %s", utils.WithPadding(key, padWidth), binding.Description)
			switch binding.ViewName {
			case "":
				contentGlobal = append(contentGlobal, content)
				bindingsGlobal = append(bindingsGlobal, binding)
			case v.Name():
				contentPanel = append(contentPanel, content)
				bindingsPanel = append(bindingsPanel, binding)
			}
		}
	}

	// append dummy element to have a separator between
	// panel and global keybindings
	contentPanel = append(contentPanel, "")
	bindingsPanel = append(bindingsPanel, Binding{})

	content := append(contentPanel, contentGlobal...)
	gui.State.Keys = append(bindingsPanel, bindingsGlobal...)
	// append newline at the end so the last line would be selectable
	contentJoined := strings.Join(content, "\n") + "\n"

	// y1-1 so there will not be an extra space at the end of panel
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(contentJoined)
	menuView, _ := gui.g.SetView("menu", x0, y0, x1, y1-1, 0)
	menuView.Title = strings.Title(gui.Tr.SLocalize("menu"))
	menuView.FgColor = gocui.ColorWhite

	err := gui.renderMenuOptions()
	if err != nil {
		gui.Log.Errorf("Failed to render menu options at handleMenu: %s\n", err)
		return err
	}

	// TODO change this to renderString
	fmt.Fprint(menuView, contentJoined)

	gui.g.Update(func(g *gocui.Gui) error {

		_, err := gui.g.SetViewOnTop("menu")
		if err != nil {
			gui.Log.Errorf("Failed to setViewOnTop at handleMenu: %s\n", err)
			return err
		}

		err = gui.switchFocus(gui.g, v, menuView)
		if err != nil {
			gui.Log.Errorf("Failed to switchFocus at handleMenu")
			return err
		}

		return nil
	})

	return nil
}
