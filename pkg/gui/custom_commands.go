package gui

import (
	"bytes"
	"log"
	"text/template"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
)

type CustomCommandObjects struct {
	SelectedLocalCommit  *commands.Commit
	SelectedReflogCommit *commands.Commit
	SelectedSubCommit    *commands.Commit
	SelectedFile         *commands.File
	SelectedLocalBranch  *commands.Branch
	SelectedRemoteBranch *commands.RemoteBranch
	SelectedRemote       *commands.Remote
	SelectedTag          *commands.Tag
	SelectedStashEntry   *commands.StashEntry
	SelectedCommitFile   *commands.CommitFile
	CurrentBranch        *commands.Branch
}

func (gui *Gui) handleCustomCommandKeybinding(customCommand CustomCommand) func() error {
	return func() error {
		objects := CustomCommandObjects{
			SelectedFile:         gui.getSelectedFile(),
			SelectedLocalCommit:  gui.getSelectedLocalCommit(),
			SelectedReflogCommit: gui.getSelectedReflogCommit(),
			SelectedLocalBranch:  gui.getSelectedBranch(),
			SelectedRemoteBranch: gui.getSelectedRemoteBranch(),
			SelectedRemote:       gui.getSelectedRemote(),
			SelectedTag:          gui.getSelectedTag(),
			SelectedStashEntry:   gui.getSelectedStashEntry(),
			SelectedCommitFile:   gui.getSelectedCommitFile(),
			SelectedSubCommit:    gui.getSelectedSubCommit(),
			CurrentBranch:        gui.currentBranch(),
		}

		tmpl, err := template.New("custom command template").Parse(customCommand.Command)
		if err != nil {
			return gui.surfaceError(err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, objects); err != nil {
			return gui.surfaceError(err)
		}

		cmdStr := buf.String()

		if customCommand.Subprocess {
			gui.PrepareSubProcess(cmdStr)
			return nil
		}

		return gui.WithWaitingStatus(gui.Tr.SLocalize("runningCustomCommandStatus"), func() error {
			gui.OSCommand.PrepareSubProcess(cmdStr)

			if err := gui.OSCommand.RunCommand(cmdStr); err != nil {
				return gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{})
		})
	}
}

type CustomCommand struct {
	Key        string `yaml:"key"`
	Context    string `yaml:"context"`
	Command    string `yaml:"command"`
	Subprocess bool   `yaml:"subprocess"`
}

func (gui *Gui) GetCustomCommandKeybindings() []*Binding {
	bindings := []*Binding{}

	var customCommands []CustomCommand

	if err := gui.Config.GetUserConfig().UnmarshalKey("customCommands", &customCommands); err != nil {
		log.Fatalf("Error parsing custom command keybindings: %v", err)
	}

	for _, customCommand := range customCommands {
		var viewName string
		if customCommand.Context == "global" || customCommand.Context == "" {
			viewName = ""
		} else {
			context := gui.contextForContextKey(customCommand.Context)
			if context == nil {
				log.Fatalf("Error when setting custom command keybindings: unknown context: %s", customCommand.Context)
			}
			// here we assume that a given context will always belong to the same view.
			// Currently this is a safe bet but it's by no means guaranteed in the long term
			// and we might need to make some changes in the future to support it.
			viewName = context.GetViewName()
		}

		bindings = append(bindings, &Binding{
			ViewName:    viewName,
			Contexts:    []string{customCommand.Context},
			Key:         gui.getKey(customCommand.Key),
			Modifier:    gocui.ModNone,
			Handler:     gui.wrappedHandler(gui.handleCustomCommandKeybinding(customCommand)),
			Description: customCommand.Command,
		})
	}

	return bindings
}
