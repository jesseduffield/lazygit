package gui

import (
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type CustomCommandObjects struct {
	SelectedLocalCommit    *models.Commit
	SelectedReflogCommit   *models.Commit
	SelectedSubCommit      *models.Commit
	SelectedFile           *models.File
	SelectedPath           string
	SelectedLocalBranch    *models.Branch
	SelectedRemoteBranch   *models.RemoteBranch
	SelectedRemote         *models.Remote
	SelectedTag            *models.Tag
	SelectedStashEntry     *models.StashEntry
	SelectedCommitFile     *models.CommitFile
	SelectedCommitFilePath string
	CheckedOutBranch       *models.Branch
	PromptResponses        []string
}

func (gui *Gui) resolveTemplate(templateStr string, promptResponses []string) (string, error) {
	objects := CustomCommandObjects{
		SelectedFile:           gui.getSelectedFile(),
		SelectedPath:           gui.getSelectedPath(),
		SelectedLocalCommit:    gui.getSelectedLocalCommit(),
		SelectedReflogCommit:   gui.getSelectedReflogCommit(),
		SelectedLocalBranch:    gui.getSelectedBranch(),
		SelectedRemoteBranch:   gui.getSelectedRemoteBranch(),
		SelectedRemote:         gui.getSelectedRemote(),
		SelectedTag:            gui.getSelectedTag(),
		SelectedStashEntry:     gui.getSelectedStashEntry(),
		SelectedCommitFile:     gui.getSelectedCommitFile(),
		SelectedCommitFilePath: gui.getSelectedCommitFilePath(),
		SelectedSubCommit:      gui.getSelectedSubCommit(),
		CheckedOutBranch:       gui.CurrentBranch(),
		PromptResponses:        promptResponses,
	}

	return utils.ResolveTemplate(templateStr, objects)
}

func (gui *Gui) handleCustomCommandKeybinding(customCommand config.CustomCommand) func() error {
	return func() error {
		promptResponses := make([]string, len(customCommand.Prompts))

		f := func() error {
			cmdStr, err := gui.resolveTemplate(customCommand.Command, promptResponses)
			if err != nil {
				return gui.SurfaceError(err)
			}

			cmdObj := gui.GitCommand.BuildShellCmdObj(cmdStr)

			if customCommand.Subprocess {
				return gui.runSubprocessWithSuspenseAndRefresh(cmdObj)
			}

			loadingText := customCommand.LoadingText
			if loadingText == "" {
				loadingText = gui.Tr.LcRunningCustomCommandStatus
			}
			return gui.WithWaitingStatus(loadingText, func() error {
				if err := gui.OSCommand.WithSpan(gui.Tr.Spans.CustomCommand).Run(cmdObj); err != nil {
					return gui.SurfaceError(err)
				}
				return gui.RefreshSidePanels(RefreshOptions{})
			})
		}

		// if we have prompts we'll recursively wrap our confirm handlers with more prompts
		// until we reach the actual command
		for reverseIdx := range customCommand.Prompts {
			idx := len(customCommand.Prompts) - 1 - reverseIdx

			// going backwards so the outermost prompt is the first one
			prompt := customCommand.Prompts[idx]

			// need to do this because f's value will change with each iteration
			wrappedF := f

			switch prompt.Type {
			case "input":
				f = func() error {
					title, err := gui.resolveTemplate(prompt.Title, promptResponses)
					if err != nil {
						return gui.SurfaceError(err)
					}

					initialValue, err := gui.resolveTemplate(prompt.InitialValue, promptResponses)
					if err != nil {
						return gui.SurfaceError(err)
					}

					return gui.Prompt(PromptOpts{
						Title:          title,
						InitialContent: initialValue,
						HandleConfirm: func(str string) error {
							promptResponses[idx] = str

							return wrappedF()
						},
					})
				}
			case "menu":
				f = func() error {
					// need to make a menu here some how
					menuItems := make([]*menuItem, len(prompt.Options))
					for i, option := range prompt.Options {
						option := option

						nameTemplate := option.Name
						if nameTemplate == "" {
							// this allows you to only pass values rather than bother with names/descriptions
							nameTemplate = option.Value
						}
						name, err := gui.resolveTemplate(nameTemplate, promptResponses)
						if err != nil {
							return gui.SurfaceError(err)
						}

						description, err := gui.resolveTemplate(option.Description, promptResponses)
						if err != nil {
							return gui.SurfaceError(err)
						}

						value, err := gui.resolveTemplate(option.Value, promptResponses)
						if err != nil {
							return gui.SurfaceError(err)
						}

						menuItems[i] = &menuItem{
							displayStrings: []string{name, utils.ColoredString(description, color.FgYellow)},
							onPress: func() error {
								promptResponses[idx] = value

								return wrappedF()
							},
						}
					}

					title, err := gui.resolveTemplate(prompt.Title, promptResponses)
					if err != nil {
						return gui.SurfaceError(err)
					}

					return gui.createMenu(title, menuItems, createMenuOptions{showCancel: true})
				}
			default:
				return gui.CreateErrorPanel("custom command prompt must have a type of 'input' or 'menu'")
			}

		}

		return f()
	}
}

func (gui *Gui) GetCustomCommandKeybindings() []*Binding {
	bindings := []*Binding{}
	customCommands := gui.Config.GetUserConfig().CustomCommands

	for _, customCommand := range customCommands {
		var viewName string
		var contexts []string
		switch customCommand.Context {
		case "global":
			viewName = ""
		case "":
			log.Fatalf("Error parsing custom command keybindings: context not provided (use context: 'global' for the global context). Key: %s, Command: %s", customCommand.Key, customCommand.Command)
		default:
			context, ok := gui.contextForContextKey(ContextKey(customCommand.Context))
			// stupid golang making me build an array of strings for this.
			allContextKeyStrings := make([]string, len(AllContextKeys))
			for i := range AllContextKeys {
				allContextKeyStrings[i] = string(AllContextKeys[i])
			}
			if !ok {
				log.Fatalf("Error when setting custom command keybindings: unknown context: %s. Key: %s, Command: %s.\nPermitted contexts: %s", customCommand.Context, customCommand.Key, customCommand.Command, strings.Join(allContextKeyStrings, ", "))
			}
			// here we assume that a given context will always belong to the same view.
			// Currently this is a safe bet but it's by no means guaranteed in the long term
			// and we might need to make some changes in the future to support it.
			viewName = context.GetViewName()
			contexts = []string{customCommand.Context}
		}

		description := customCommand.Description
		if description == "" {
			description = customCommand.Command
		}

		bindings = append(bindings, &Binding{
			ViewName:    viewName,
			Contexts:    contexts,
			Key:         gui.getKey(customCommand.Key),
			Modifier:    gocui.ModNone,
			Handler:     gui.handleCustomCommandKeybinding(customCommand),
			Description: description,
		})
	}

	return bindings
}
