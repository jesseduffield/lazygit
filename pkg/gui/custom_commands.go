package gui

import (
	"bytes"
	"text/template"

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

func (gui *Gui) handleCustomCommandKeybinding(templateStr string) func() error {
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

		tmpl, err := template.New("custom command template").Parse(templateStr)
		if err != nil {
			return gui.surfaceError(err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, objects); err != nil {
			return gui.surfaceError(err)
		}

		cmdStr := buf.String()

		return gui.WithWaitingStatus(gui.Tr.SLocalize("runningCustomCommandStatus"), func() error {
			if err := gui.OSCommand.RunCommand(cmdStr); err != nil {
				return gui.surfaceError(err)
			}
			return gui.refreshSidePanels(refreshOptions{})
		})
	}
}
