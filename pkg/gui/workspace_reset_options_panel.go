package gui

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

type workspaceResetOption struct {
	handler     func() error
	description string
	command     string
}

// GetDisplayStrings is a function.
func (r *workspaceResetOption) GetDisplayStrings(isFocused bool) []string {
	return []string{r.description, color.New(color.FgRed).Sprint(r.command)}
}

func (gui *Gui) handleCreateResetMenu(g *gocui.Gui, v *gocui.View) error {
	options := []*workspaceResetOption{
		{
			description: gui.Tr.SLocalize("discardAllChangesToAllFiles"),
			command:     "reset --hard HEAD && git clean -fd",
			handler: func() error {
				if err := gui.GitCommand.ResetAndClean(); err != nil {
					return gui.createErrorPanel(gui.g, err.Error())
				}

				return gui.refreshFiles()
			},
		},
		{
			description: gui.Tr.SLocalize("discardAnyUnstagedChanges"),
			command:     "git checkout -- .",
			handler: func() error {
				if err := gui.GitCommand.DiscardAnyUnstagedFileChanges(); err != nil {
					return gui.createErrorPanel(gui.g, err.Error())
				}

				return gui.refreshFiles()
			},
		},
		{
			description: gui.Tr.SLocalize("discardUntrackedFiles"),
			command:     "git clean -fd",
			handler: func() error {
				if err := gui.GitCommand.RemoveUntrackedFiles(); err != nil {
					return gui.createErrorPanel(gui.g, err.Error())
				}

				return gui.refreshFiles()
			},
		},
		{
			description: gui.Tr.SLocalize("softReset"),
			command:     "git reset --soft HEAD",
			handler: func() error {
				if err := gui.GitCommand.ResetSoft("HEAD"); err != nil {
					return gui.createErrorPanel(gui.g, err.Error())
				}

				return gui.refreshFiles()
			},
		},
		{
			description: gui.Tr.SLocalize("hardReset"),
			command:     "git reset --hard HEAD",
			handler: func() error {
				if err := gui.GitCommand.ResetHard("HEAD"); err != nil {
					return gui.createErrorPanel(gui.g, err.Error())
				}

				return gui.refreshFiles()
			},
		},
		{
			description: gui.Tr.SLocalize("hardResetUpstream"),
			command:     "git reset --hard @{upstream}",
			handler: func() error {
				if err := gui.GitCommand.ResetHard("@{upstream}"); err != nil {
					return gui.createErrorPanel(gui.g, err.Error())
				}

				return gui.refreshSidePanels(gui.g)
			},
		},
		{
			description: gui.Tr.SLocalize("cancel"),
			handler: func() error {
				return nil
			},
		},
	}

	handleMenuPress := func(index int) error {
		return options[index].handler()
	}

	return gui.createMenu("", options, len(options), handleMenuPress)
}
