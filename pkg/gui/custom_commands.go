package gui

import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

type commandMenuEntry struct {
	label string
	value string
}

func (gui *Gui) getResolveTemplateFn(promptResponses []string) func(string) (string, error) {
	objects := CustomCommandObjects{
		SelectedFile:           gui.getSelectedFile(),
		SelectedPath:           gui.getSelectedPath(),
		SelectedLocalCommit:    gui.State.Contexts.LocalCommits.GetSelected(),
		SelectedReflogCommit:   gui.State.Contexts.ReflogCommits.GetSelected(),
		SelectedLocalBranch:    gui.State.Contexts.Branches.GetSelected(),
		SelectedRemoteBranch:   gui.State.Contexts.RemoteBranches.GetSelected(),
		SelectedRemote:         gui.State.Contexts.Remotes.GetSelected(),
		SelectedTag:            gui.State.Contexts.Tags.GetSelected(),
		SelectedStashEntry:     gui.State.Contexts.Stash.GetSelected(),
		SelectedCommitFile:     gui.getSelectedCommitFile(),
		SelectedCommitFilePath: gui.getSelectedCommitFilePath(),
		SelectedSubCommit:      gui.State.Contexts.SubCommits.GetSelected(),
		CheckedOutBranch:       gui.helpers.Refs.GetCheckedOutRef(),
		PromptResponses:        promptResponses,
	}

	return func(templateStr string) (string, error) { return utils.ResolveTemplate(templateStr, objects) }
}

func resolveCustomCommandPrompt(prompt *config.CustomCommandPrompt, resolveTemplate func(string) (string, error)) (*config.CustomCommandPrompt, error) {
	var err error
	result := &config.CustomCommandPrompt{}

	result.Title, err = resolveTemplate(prompt.Title)
	if err != nil {
		return nil, err
	}

	result.InitialValue, err = resolveTemplate(prompt.InitialValue)
	if err != nil {
		return nil, err
	}

	result.Command, err = resolveTemplate(prompt.Command)
	if err != nil {
		return nil, err
	}

	result.Filter, err = resolveTemplate(prompt.Filter)
	if err != nil {
		return nil, err
	}

	if len(prompt.Options) > 0 {
		newOptions := make([]config.CustomCommandMenuOption, len(prompt.Options))
		for _, option := range prompt.Options {
			option := option
			newOption, err := resolveMenuOption(&option, resolveTemplate)
			if err != nil {
				return nil, err
			}
			newOptions = append(newOptions, *newOption)
		}
		prompt.Options = newOptions
	}

	return result, nil
}

func resolveMenuOption(option *config.CustomCommandMenuOption, resolveTemplate func(string) (string, error)) (*config.CustomCommandMenuOption, error) {
	nameTemplate := option.Name
	if nameTemplate == "" {
		// this allows you to only pass values rather than bother with names/descriptions
		nameTemplate = option.Value
	}

	name, err := resolveTemplate(nameTemplate)
	if err != nil {
		return nil, err
	}

	description, err := resolveTemplate(option.Description)
	if err != nil {
		return nil, err
	}

	value, err := resolveTemplate(option.Value)
	if err != nil {
		return nil, err
	}

	return &config.CustomCommandMenuOption{
		Name:        name,
		Description: description,
		Value:       value,
	}, nil
}

func (gui *Gui) inputPrompt(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	return gui.c.Prompt(types.PromptOpts{
		Title:          prompt.Title,
		InitialContent: prompt.InitialValue,
		HandleConfirm: func(str string) error {
			return wrappedF(str)
		},
	})
}

func (gui *Gui) menuPrompt(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	menuItems := make([]*types.MenuItem, len(prompt.Options))
	for i, option := range prompt.Options {
		option := option
		menuItems[i] = &types.MenuItem{
			DisplayStrings: []string{option.Name, style.FgYellow.Sprint(option.Description)},
			OnPress: func() error {
				return wrappedF(option.Value)
			},
		}
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: prompt.Title, Items: menuItems})
}

func (gui *Gui) GenerateMenuCandidates(commandOutput, filter, valueFormat, labelFormat string) ([]commandMenuEntry, error) {
	reg, err := regexp.Compile(filter)
	if err != nil {
		return nil, gui.c.Error(errors.New("unable to parse filter regex, error: " + err.Error()))
	}

	buff := bytes.NewBuffer(nil)

	valueTemp, err := template.New("format").Parse(valueFormat)
	if err != nil {
		return nil, gui.c.Error(errors.New("unable to parse value format, error: " + err.Error()))
	}

	colorFuncMap := style.TemplateFuncMapAddColors(template.FuncMap{})

	descTemp, err := template.New("format").Funcs(colorFuncMap).Parse(labelFormat)
	if err != nil {
		return nil, gui.c.Error(errors.New("unable to parse label format, error: " + err.Error()))
	}

	candidates := []commandMenuEntry{}
	for _, str := range strings.Split(commandOutput, "\n") {
		if str == "" {
			continue
		}

		tmplData := map[string]string{}
		out := reg.FindAllStringSubmatch(str, -1)
		if len(out) > 0 {
			for groupIdx, group := range reg.SubexpNames() {
				// Record matched group with group ids
				matchName := "group_" + strconv.Itoa(groupIdx)
				tmplData[matchName] = out[0][groupIdx]
				// Record last named group non-empty matches as group matches
				if group != "" {
					tmplData[group] = out[0][groupIdx]
				}
			}
		}

		err = valueTemp.Execute(buff, tmplData)
		if err != nil {
			return candidates, gui.c.Error(err)
		}
		entry := commandMenuEntry{
			value: strings.TrimSpace(buff.String()),
		}

		if labelFormat != "" {
			buff.Reset()
			err = descTemp.Execute(buff, tmplData)
			if err != nil {
				return candidates, gui.c.Error(err)
			}
			entry.label = strings.TrimSpace(buff.String())
		} else {
			entry.label = entry.value
		}

		candidates = append(candidates, entry)

		buff.Reset()
	}
	return candidates, err
}

func (gui *Gui) menuPromptFromCommand(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	// Run and save output
	message, err := gui.git.Custom.RunWithOutput(prompt.Command)
	if err != nil {
		return gui.c.Error(err)
	}

	// Need to make a menu out of what the cmd has displayed
	candidates, err := gui.GenerateMenuCandidates(message, prompt.Filter, prompt.ValueFormat, prompt.LabelFormat)
	if err != nil {
		return gui.c.Error(err)
	}

	menuItems := make([]*types.MenuItem, len(candidates))
	for i := range candidates {
		i := i
		menuItems[i] = &types.MenuItem{
			DisplayStrings: []string{candidates[i].label},
			OnPress: func() error {
				return wrappedF(candidates[i].value)
			},
		}
	}

	return gui.c.Menu(types.CreateMenuOptions{Title: prompt.Title, Items: menuItems})
}

func (gui *Gui) handleCustomCommandKeybinding(customCommand config.CustomCommand) func() error {
	return func() error {
		promptResponses := make([]string, len(customCommand.Prompts))

		f := func() error {
			resolveTemplate := gui.getResolveTemplateFn(promptResponses)
			cmdStr, err := resolveTemplate(customCommand.Command)
			if err != nil {
				return gui.c.Error(err)
			}

			if customCommand.Subprocess {
				return gui.runSubprocessWithSuspenseAndRefresh(gui.os.Cmd.NewShell(cmdStr))
			}

			loadingText := customCommand.LoadingText
			if loadingText == "" {
				loadingText = gui.c.Tr.LcRunningCustomCommandStatus
			}

			return gui.c.WithWaitingStatus(loadingText, func() error {
				gui.c.LogAction(gui.c.Tr.Actions.CustomCommand)
				cmdObj := gui.os.Cmd.NewShell(cmdStr)
				if customCommand.Stream {
					cmdObj.StreamOutput()
				}
				err := cmdObj.Run()
				if err != nil {
					return gui.c.Error(err)
				}
				return gui.c.Refresh(types.RefreshOptions{})
			})
		}

		// if we have prompts we'll recursively wrap our confirm handlers with more prompts
		// until we reach the actual command
		for reverseIdx := range customCommand.Prompts {
			idx := len(customCommand.Prompts) - 1 - reverseIdx

			// going backwards so the outermost prompt is the first one
			prompt := customCommand.Prompts[idx]

			wrappedF := func(response string) error {
				promptResponses[idx] = response
				return f()
			}

			resolveTemplate := gui.getResolveTemplateFn(promptResponses)
			resolvedPrompt, err := resolveCustomCommandPrompt(&prompt, resolveTemplate)
			if err != nil {
				return gui.c.Error(err)
			}

			switch prompt.Type {
			case "input":
				f = func() error {
					return gui.inputPrompt(resolvedPrompt, wrappedF)
				}
			case "menu":
				f = func() error {
					return gui.menuPrompt(resolvedPrompt, wrappedF)
				}
			case "menuFromCommand":
				f = func() error {
					return gui.menuPromptFromCommand(resolvedPrompt, wrappedF)
				}
			default:
				return gui.c.ErrorMsg("custom command prompt must have a type of 'input', 'menu' or 'menuFromCommand'")
			}
		}

		return f()
	}
}

func (gui *Gui) GetCustomCommandKeybindings() []*types.Binding {
	bindings := []*types.Binding{}
	customCommands := gui.c.UserConfig.CustomCommands

	for _, customCommand := range customCommands {
		var viewName string
		var contexts []string
		switch customCommand.Context {
		case "global":
			viewName = ""
		case "":
			log.Fatalf("Error parsing custom command keybindings: context not provided (use context: 'global' for the global context). Key: %s, Command: %s", customCommand.Key, customCommand.Command)
		default:
			ctx, ok := gui.contextForContextKey(types.ContextKey(customCommand.Context))
			// stupid golang making me build an array of strings for this.
			allContextKeyStrings := make([]string, len(context.AllContextKeys))
			for i := range context.AllContextKeys {
				allContextKeyStrings[i] = string(context.AllContextKeys[i])
			}
			if !ok {
				log.Fatalf("Error when setting custom command keybindings: unknown context: %s. Key: %s, Command: %s.\nPermitted contexts: %s", customCommand.Context, customCommand.Key, customCommand.Command, strings.Join(allContextKeyStrings, ", "))
			}
			// here we assume that a given context will always belong to the same view.
			// Currently this is a safe bet but it's by no means guaranteed in the long term
			// and we might need to make some changes in the future to support it.
			viewName = ctx.GetViewName()
			contexts = []string{customCommand.Context}
		}

		description := customCommand.Description
		if description == "" {
			description = customCommand.Command
		}

		bindings = append(bindings, &types.Binding{
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
