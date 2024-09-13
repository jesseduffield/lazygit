package custom_commands

import (
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// takes a custom command and returns a function that will be called when the corresponding user-defined keybinding is pressed
type HandlerCreator struct {
	c                    *helpers.HelperCommon
	sessionStateLoader   *SessionStateLoader
	resolver             *Resolver
	menuGenerator        *MenuGenerator
	suggestionsHelper    *helpers.SuggestionsHelper
	mergeAndRebaseHelper *helpers.MergeAndRebaseHelper
}

func NewHandlerCreator(
	c *helpers.HelperCommon,
	sessionStateLoader *SessionStateLoader,
	suggestionsHelper *helpers.SuggestionsHelper,
	mergeAndRebaseHelper *helpers.MergeAndRebaseHelper,
) *HandlerCreator {
	resolver := NewResolver(c.Common)
	menuGenerator := NewMenuGenerator(c.Common)

	return &HandlerCreator{
		c:                    c,
		sessionStateLoader:   sessionStateLoader,
		resolver:             resolver,
		menuGenerator:        menuGenerator,
		suggestionsHelper:    suggestionsHelper,
		mergeAndRebaseHelper: mergeAndRebaseHelper,
	}
}

func (self *HandlerCreator) call(customCommand config.CustomCommand) func() error {
	return func() error {
		sessionState := self.sessionStateLoader.call()
		promptResponses := make([]string, len(customCommand.Prompts))
		form := make(map[string]string)

		f := func() error { return self.finalHandler(customCommand, sessionState, promptResponses, form) }

		// if we have prompts we'll recursively wrap our confirm handlers with more prompts
		// until we reach the actual command
		for reverseIdx := range customCommand.Prompts {
			// reassigning so that we don't end up with an infinite recursion
			g := f
			idx := len(customCommand.Prompts) - 1 - reverseIdx

			// going backwards so the outermost prompt is the first one
			prompt := customCommand.Prompts[idx]

			wrappedF := func(response string) error {
				promptResponses[idx] = response
				form[prompt.Key] = response
				return g()
			}

			resolveTemplate := self.getResolveTemplateFn(form, promptResponses, sessionState)

			switch prompt.Type {
			case "input":
				f = func() error {
					resolvedPrompt, err := self.resolver.resolvePrompt(&prompt, resolveTemplate)
					if err != nil {
						return err
					}
					return self.inputPrompt(resolvedPrompt, wrappedF)
				}
			case "menu":
				f = func() error {
					resolvedPrompt, err := self.resolver.resolvePrompt(&prompt, resolveTemplate)
					if err != nil {
						return err
					}
					return self.menuPrompt(resolvedPrompt, wrappedF)
				}
			case "menuFromCommand":
				f = func() error {
					resolvedPrompt, err := self.resolver.resolvePrompt(&prompt, resolveTemplate)
					if err != nil {
						return err
					}
					return self.menuPromptFromCommand(resolvedPrompt, wrappedF)
				}
			case "confirm":
				f = func() error {
					resolvedPrompt, err := self.resolver.resolvePrompt(&prompt, resolveTemplate)
					if err != nil {
						return err
					}
					return self.confirmPrompt(resolvedPrompt, g)
				}
			default:
				return errors.New("custom command prompt must have a type of 'input', 'menu', 'menuFromCommand', or 'confirm'")
			}
		}

		return f()
	}
}

func (self *HandlerCreator) inputPrompt(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	findSuggestionsFn, err := self.generateFindSuggestionsFunc(prompt)
	if err != nil {
		return err
	}

	self.c.Prompt(types.PromptOpts{
		Title:               prompt.Title,
		InitialContent:      prompt.InitialValue,
		FindSuggestionsFunc: findSuggestionsFn,
		HandleConfirm: func(str string) error {
			return wrappedF(str)
		},
	})

	return nil
}

func (self *HandlerCreator) generateFindSuggestionsFunc(prompt *config.CustomCommandPrompt) (func(string) []*types.Suggestion, error) {
	if prompt.Suggestions.Preset != "" && prompt.Suggestions.Command != "" {
		return nil, fmt.Errorf(
			"Custom command prompt cannot have both a preset and a command for suggestions. Preset: '%s', Command: '%s'",
			prompt.Suggestions.Preset,
			prompt.Suggestions.Command,
		)
	} else if prompt.Suggestions.Preset != "" {
		return self.getPresetSuggestionsFn(prompt.Suggestions.Preset)
	} else if prompt.Suggestions.Command != "" {
		return self.getCommandSuggestionsFn(prompt.Suggestions.Command)
	}

	return nil, nil
}

func (self *HandlerCreator) getCommandSuggestionsFn(command string) (func(string) []*types.Suggestion, error) {
	lines := []*types.Suggestion{}
	err := self.c.OS().Cmd.NewShell(command).RunAndProcessLines(func(line string) (bool, error) {
		lines = append(lines, &types.Suggestion{Value: line, Label: line})
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return func(currentWord string) []*types.Suggestion {
		return lo.Filter(lines, func(suggestion *types.Suggestion, _ int) bool {
			return strings.Contains(strings.ToLower(suggestion.Value), strings.ToLower(currentWord))
		})
	}, nil
}

func (self *HandlerCreator) getPresetSuggestionsFn(preset string) (func(string) []*types.Suggestion, error) {
	switch preset {
	case "authors":
		return self.suggestionsHelper.GetAuthorsSuggestionsFunc(), nil
	case "branches":
		return self.suggestionsHelper.GetBranchNameSuggestionsFunc(), nil
	case "files":
		return self.suggestionsHelper.GetFilePathSuggestionsFunc(), nil
	case "refs":
		return self.suggestionsHelper.GetRefsSuggestionsFunc(), nil
	case "remotes":
		return self.suggestionsHelper.GetRemoteSuggestionsFunc(), nil
	case "remoteBranches":
		return self.suggestionsHelper.GetRemoteBranchesSuggestionsFunc("/"), nil
	case "tags":
		return self.suggestionsHelper.GetTagsSuggestionsFunc(), nil
	default:
		return nil, fmt.Errorf("Unknown value for suggestionsPreset in custom command: %s. Valid values: files, branches, remotes, remoteBranches, refs", preset)
	}
}

func (self *HandlerCreator) confirmPrompt(prompt *config.CustomCommandPrompt, handleConfirm func() error) error {
	self.c.Confirm(types.ConfirmOpts{
		Title:         prompt.Title,
		Prompt:        prompt.Body,
		HandleConfirm: handleConfirm,
	})

	return nil
}

func (self *HandlerCreator) menuPrompt(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	menuItems := lo.Map(prompt.Options, func(option config.CustomCommandMenuOption, _ int) *types.MenuItem {
		return &types.MenuItem{
			LabelColumns: []string{option.Name, style.FgYellow.Sprint(option.Description)},
			OnPress: func() error {
				return wrappedF(option.Value)
			},
		}
	})

	return self.c.Menu(types.CreateMenuOptions{Title: prompt.Title, Items: menuItems})
}

func (self *HandlerCreator) menuPromptFromCommand(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	// Run and save output
	message, err := self.c.Git().Custom.RunWithOutput(prompt.Command)
	if err != nil {
		return err
	}

	// Need to make a menu out of what the cmd has displayed
	candidates, err := self.menuGenerator.call(message, prompt.Filter, prompt.ValueFormat, prompt.LabelFormat)
	if err != nil {
		return err
	}

	menuItems := lo.Map(candidates, func(candidate *commandMenuItem, _ int) *types.MenuItem {
		return &types.MenuItem{
			LabelColumns: []string{candidate.label},
			OnPress: func() error {
				return wrappedF(candidate.value)
			},
		}
	})

	return self.c.Menu(types.CreateMenuOptions{Title: prompt.Title, Items: menuItems})
}

type CustomCommandObjects struct {
	*SessionState
	PromptResponses []string
	Form            map[string]string
}

func (self *HandlerCreator) getResolveTemplateFn(form map[string]string, promptResponses []string, sessionState *SessionState) func(string) (string, error) {
	objects := CustomCommandObjects{
		SessionState:    sessionState,
		PromptResponses: promptResponses,
		Form:            form,
	}

	funcs := template.FuncMap{
		"quote": self.c.OS().Quote,
	}

	return func(templateStr string) (string, error) { return utils.ResolveTemplate(templateStr, objects, funcs) }
}

func (self *HandlerCreator) finalHandler(customCommand config.CustomCommand, sessionState *SessionState, promptResponses []string, form map[string]string) error {
	resolveTemplate := self.getResolveTemplateFn(form, promptResponses, sessionState)
	cmdStr, err := resolveTemplate(customCommand.Command)
	if err != nil {
		return err
	}

	cmdObj := self.c.OS().Cmd.NewShell(cmdStr)

	if customCommand.Subprocess {
		return self.c.RunSubprocessAndRefresh(cmdObj)
	}

	loadingText := customCommand.LoadingText
	if loadingText == "" {
		loadingText = self.c.Tr.RunningCustomCommandStatus
	}

	return self.c.WithWaitingStatus(loadingText, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.CustomCommand)

		if customCommand.Stream {
			cmdObj.StreamOutput()
		}
		output, err := cmdObj.RunWithOutput()

		if refreshErr := self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC}); err != nil {
			self.c.Log.Error(refreshErr)
		}

		if err != nil {
			if customCommand.After.CheckForConflicts {
				return self.mergeAndRebaseHelper.CheckForConflicts(err)
			}

			return err
		}

		if customCommand.ShowOutput {
			if strings.TrimSpace(output) == "" {
				output = self.c.Tr.EmptyOutput
			}

			title := cmdStr
			if customCommand.OutputTitle != "" {
				title, err = resolveTemplate(customCommand.OutputTitle)
				if err != nil {
					return err
				}
			}
			self.c.Alert(title, output)
		}

		return nil
	})
}
