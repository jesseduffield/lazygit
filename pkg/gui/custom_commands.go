package gui

import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
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
		CheckedOutBranch:       gui.currentBranch(),
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
				return gui.surfaceError(err)
			}

			if customCommand.Subprocess {
				return gui.runSubprocessWithSuspenseAndRefresh(gui.OSCommand.PrepareShellSubProcess(cmdStr))
			}

			loadingText := customCommand.LoadingText
			if loadingText == "" {
				loadingText = gui.Tr.LcRunningCustomCommandStatus
			}
			return gui.WithWaitingStatus(loadingText, func() error {
				if err := gui.OSCommand.WithSpan(gui.Tr.Spans.CustomCommand).RunShellCommand(cmdStr); err != nil {
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

			// need to do this because f's value will change with each iteration
			wrappedF := f

			switch prompt.Type {
			case "input":
				f = func() error {
					title, err := gui.resolveTemplate(prompt.Title, promptResponses)
					if err != nil {
						return gui.surfaceError(err)
					}

					initialValue, err := gui.resolveTemplate(prompt.InitialValue, promptResponses)
					if err != nil {
						return gui.surfaceError(err)
					}

					return gui.prompt(promptOpts{
						title:          title,
						initialContent: initialValue,
						handleConfirm: func(str string) error {
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
							return gui.surfaceError(err)
						}

						description, err := gui.resolveTemplate(option.Description, promptResponses)
						if err != nil {
							return gui.surfaceError(err)
						}

						value, err := gui.resolveTemplate(option.Value, promptResponses)
						if err != nil {
							return gui.surfaceError(err)
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
						return gui.surfaceError(err)
					}

					return gui.createMenu(title, menuItems, createMenuOptions{showCancel: true})
				}
			case "menuFromCommand":
				f = func() error {
					// Collect cmd to run from config
					cmdStr, err := gui.resolveTemplate(prompt.Command, promptResponses)
					if err != nil {
						return gui.surfaceError(err)
					}

					// Collect Filter regexp
					filter, err := gui.resolveTemplate(prompt.Filter, promptResponses)
					if err != nil {
						return gui.surfaceError(err)
					}
					reg, err := regexp.Compile(filter)
					if err != nil {
						return gui.surfaceError(errors.New("unable to parse filter regex, error: " + err.Error()))
					}

					// Run and save output
					message, err := gui.GitCommand.RunCommandWithOutput(cmdStr)
					if err != nil {
						return gui.surfaceError(err)
					}

					// Need to make a menu out of what the cmd has displayed
					candidates := []string{}
					buff := bytes.NewBuffer(nil)
					temp, err := template.New("format").Parse(prompt.Format)
					if err != nil {
						return gui.surfaceError(errors.New("unable to parse format, error: " + err.Error()))
					}
					for _, str := range strings.Split(string(message), "\n") {
						if str == "" {
							continue
						}
						tmplData := map[string]string{}
						out := reg.FindAllStringSubmatch(str, -1)
						if len(out) > 0 {
							for groupIdx, group := range reg.SubexpNames() {
								// Record matched group with group ids
								matchName := "group_" + strconv.Itoa(groupIdx)
								tmplData[matchName] = group
								// Record last named group non-empty matches as group matches
								if group != "" {
									tmplData[group] = out[0][idx]
								}
							}
						}
						err = temp.Execute(buff, tmplData)
						if err != nil {
							return gui.surfaceError(err)
						}

						candidates = append(candidates, strings.TrimSpace(buff.String()))
						buff.Reset()
					}

					menuItems := make([]*menuItem, len(candidates))
					for i := range candidates {
						menuItems[i] = &menuItem{
							displayStrings: []string{candidates[i]},
							onPress: func() error {
								promptResponses[idx] = candidates[i]
								return wrappedF()
							},
						}
					}

					title, err := gui.resolveTemplate(prompt.Title, promptResponses)
					if err != nil {
						return gui.surfaceError(err)
					}

					return gui.createMenu(title, menuItems, createMenuOptions{showCancel: true})
				}
			default:
				return gui.createErrorPanel("custom command prompt must have a type of 'input', 'menu' or 'menuFromCommand'")
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
			allContextKeyStrings := make([]string, len(allContextKeys))
			for i := range allContextKeys {
				allContextKeyStrings[i] = string(allContextKeys[i])
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
