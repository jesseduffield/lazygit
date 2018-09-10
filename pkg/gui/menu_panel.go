package gui

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// I need to store the handler function in state and it will take an interface and do something with it
// I need to have another function describing how to display one of the structs
// perhaps this calls for an interface where the struct is Binding and the interface has the methods Display and Execute
// but this means that for the one struct I can only have one possible display/execute function, but I want to use whatever I want.
// Would I ever need to use different handlers for different things? Maybe I should assume not given that I can cross that bridge when I come to it

func (gui *Gui) handleMenuPress(g *gocui.Gui, v *gocui.View) error {
	lineNumber := gui.getItemPosition(v)
	if gui.State.Keys[lineNumber].Key == nil {
		return nil
	}
	if len(gui.State.Keys) > lineNumber {
		err := gui.handleMenuClose(g, v)
		if err != nil {
			return err
		}
		return gui.State.Keys[lineNumber].Handler(g, v)
	}
	return nil
}

func (gui *Gui) handleMenuSelect(g *gocui.Gui, v *gocui.View) error {
	// doing nothing for now
	// but it is needed for switch in newLineFocused
	return nil
}

func (gui *Gui) renderMenuOptions(g *gocui.Gui) error {
	optionsMap := map[string]string{
		"esc/q": gui.Tr.SLocalize("close"),
		"↑ ↓":   gui.Tr.SLocalize("navigate"),
		"space": gui.Tr.SLocalize("execute"),
	}
	return gui.renderOptionsMap(g, optionsMap)
}

func (gui *Gui) handleMenuClose(g *gocui.Gui, v *gocui.View) error {
	// better to delete because for example after closing update confirmation panel,
	// the focus isn't set back to any of panels and one is unable to even quit
	//_, err := g.SetViewOnBottom(v.Name())
	err := g.DeleteView("menu")
	if err != nil {
		return err
	}
	return gui.returnFocus(g, v)
}

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
	contentJoined := strings.Join(content, "\n")

	// y1-1 so there will not be an extra space at the end of panel
	x0, y0, x1, y1 := gui.getConfirmationPanelDimensions(g, contentJoined)
	menuView, _ := g.SetView("menu", x0, y0, x1, y1, 0)
	menuView.Title = strings.Title(gui.Tr.SLocalize("menu"))
	menuView.FgColor = gocui.ColorWhite

	if err := gui.renderMenuOptions(g); err != nil {
		return err
	}

	fmt.Fprint(menuView, contentJoined)

	g.Update(func(g *gocui.Gui) error {
		_, err := g.SetViewOnTop("menu")
		if err != nil {
			return err
		}
		return gui.switchFocus(g, v, menuView)
	})
	return nil
}
