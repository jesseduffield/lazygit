package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// takes a custom command and returns a function that will be called when the corresponding user-defined keybinding is pressed
type HandlerCreator struct {
	c                  *types.HelperCommon
	os                 *oscommands.OSCommand
	git                *commands.GitCommand
	sessionStateLoader *SessionStateLoader
	resolver           *Resolver
	menuGenerator      *MenuGenerator
}

func NewHandlerCreator(
	c *types.HelperCommon,
	os *oscommands.OSCommand,
	git *commands.GitCommand,
	sessionStateLoader *SessionStateLoader,
) *HandlerCreator {
	resolver := NewResolver(c.Common)
	menuGenerator := NewMenuGenerator(c.Common)

	return &HandlerCreator{
		c:                  c,
		os:                 os,
		git:                git,
		sessionStateLoader: sessionStateLoader,
		resolver:           resolver,
		menuGenerator:      menuGenerator,
	}
}

func (self *HandlerCreator) call(customCommand config.CustomCommand) func() error {
	return func() error {
		sessionState := self.sessionStateLoader.call()
		promptResponses := make([]string, len(customCommand.Prompts))

		f := func() error { return self.finalHandler(customCommand, sessionState, promptResponses) }

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
				return g()
			}

			resolveTemplate := self.getResolveTemplateFn(promptResponses, sessionState)
			resolvedPrompt, err := self.resolver.resolvePrompt(&prompt, resolveTemplate)
			if err != nil {
				return self.c.Error(err)
			}

			switch prompt.Type {
			case "input":
				f = func() error {
					return self.inputPrompt(resolvedPrompt, wrappedF)
				}
			case "menu":
				f = func() error {
					return self.menuPrompt(resolvedPrompt, wrappedF)
				}
			case "menuFromCommand":
				f = func() error {
					return self.menuPromptFromCommand(resolvedPrompt, wrappedF)
				}
			default:
				return self.c.ErrorMsg("custom command prompt must have a type of 'input', 'menu' or 'menuFromCommand'")
			}
		}

		return f()
	}
}

func (self *HandlerCreator) inputPrompt(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	return self.c.Prompt(types.PromptOpts{
		Title:          prompt.Title,
		InitialContent: prompt.InitialValue,
		HandleConfirm: func(str string) error {
			return wrappedF(str)
		},
	})
}

func (self *HandlerCreator) menuPrompt(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
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

	return self.c.Menu(types.CreateMenuOptions{Title: prompt.Title, Items: menuItems})
}

func (self *HandlerCreator) menuPromptFromCommand(prompt *config.CustomCommandPrompt, wrappedF func(string) error) error {
	// Run and save output
	message, err := self.git.Custom.RunWithOutput(prompt.Command)
	if err != nil {
		return self.c.Error(err)
	}

	// Need to make a menu out of what the cmd has displayed
	candidates, err := self.menuGenerator.call(message, prompt.Filter, prompt.ValueFormat, prompt.LabelFormat)
	if err != nil {
		return self.c.Error(err)
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

	return self.c.Menu(types.CreateMenuOptions{Title: prompt.Title, Items: menuItems})
}

type CustomCommandObjects struct {
	*SessionState
	PromptResponses []string
}

func (self *HandlerCreator) getResolveTemplateFn(promptResponses []string, sessionState *SessionState) func(string) (string, error) {
	objects := CustomCommandObjects{
		SessionState:    sessionState,
		PromptResponses: promptResponses,
	}

	return func(templateStr string) (string, error) { return utils.ResolveTemplate(templateStr, objects) }
}

func (self *HandlerCreator) finalHandler(customCommand config.CustomCommand, sessionState *SessionState, promptResponses []string) error {
	resolveTemplate := self.getResolveTemplateFn(promptResponses, sessionState)
	cmdStr, err := resolveTemplate(customCommand.Command)
	if err != nil {
		return self.c.Error(err)
	}

	cmdObj := self.os.Cmd.NewShell(cmdStr)

	if customCommand.Subprocess {
		return self.c.RunSubprocessAndRefresh(cmdObj)
	}

	loadingText := customCommand.LoadingText
	if loadingText == "" {
		loadingText = self.c.Tr.LcRunningCustomCommandStatus
	}

	return self.c.WithWaitingStatus(loadingText, func() error {
		self.c.LogAction(self.c.Tr.Actions.CustomCommand)

		if customCommand.Stream {
			cmdObj.StreamOutput()
		}
		err := cmdObj.Run()
		if err != nil {
			return self.c.Error(err)
		}
		return self.c.Refresh(types.RefreshOptions{})
	})
}
