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
	CheckedOutBranch     *commands.Branch
	PromptResponses      []string
}

func (gui *Gui) resolveTemplate(templateStr string, promptResponses []string) (string, error) {
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
		CheckedOutBranch:     gui.currentBranch(),
		PromptResponses:      promptResponses,
	}

	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, objects); err != nil {
		return "", err
	}

	cmdStr := buf.String()

	return cmdStr, nil
}

func (gui *Gui) handleCustomCommandKeybinding(customCommand CustomCommand) func() error {
	return func() error {
		promptResponses := make([]string, len(customCommand.Prompts))

		f := func() error {
			cmdStr, err := gui.resolveTemplate(customCommand.Command, promptResponses)
			if err != nil {
				return gui.surfaceError(err)
			}

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

		// if we have prompts we'll recursively wrap our confirm handlers with more prompts
		// until we reach the actual command
		for reverseIdx := range customCommand.Prompts {
			idx := len(customCommand.Prompts) - 1 - reverseIdx

			// going backwards so the outermost prompt is the first one
			prompt := customCommand.Prompts[idx]

			gui.Log.Warn(prompt.Title)

			wrappedF := f // need to do this because f's value will change with each iteration
			f = func() error {
				return gui.prompt(
					prompt.Title,
					prompt.InitialValue,
					func(str string) error {
						promptResponses[idx] = str

						return wrappedF()
					},
				)
			}
		}

		return f()
	}
}

type CustomCommandPrompt struct {
	Title        string `yaml:"title"`
	InitialValue string `yaml:"initialValue"`
}

type CustomCommand struct {
	Key        string                `yaml:"key"`
	Context    string                `yaml:"context"`
	Command    string                `yaml:"command"`
	Subprocess bool                  `yaml:"subprocess"`
	Prompts    []CustomCommandPrompt `yaml:"prompts"`
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
