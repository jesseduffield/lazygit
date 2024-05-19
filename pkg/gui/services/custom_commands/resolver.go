package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
)

// takes a prompt that is defined in terms of template strings and resolves the templates to contain actual values
type Resolver struct {
	c *common.Common
}

func NewResolver(c *common.Common) *Resolver {
	return &Resolver{c: c}
}

func (self *Resolver) resolvePrompt(
	prompt *config.CustomCommandPrompt,
	resolveTemplate func(string) (string, error),
) (*config.CustomCommandPrompt, error) {
	var err error
	result := &config.CustomCommandPrompt{
		ValueFormat: prompt.ValueFormat,
		LabelFormat: prompt.LabelFormat,
	}

	result.Title, err = resolveTemplate(prompt.Title)
	if err != nil {
		return nil, err
	}

	result.InitialValue, err = resolveTemplate(prompt.InitialValue)
	if err != nil {
		return nil, err
	}

	result.Suggestions.Preset, err = resolveTemplate(prompt.Suggestions.Preset)
	if err != nil {
		return nil, err
	}

	result.Suggestions.Command, err = resolveTemplate(prompt.Suggestions.Command)
	if err != nil {
		return nil, err
	}

	result.Body, err = resolveTemplate(prompt.Body)
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

	if prompt.Type == "menu" {
		result.Options, err = self.resolveMenuOptions(prompt, resolveTemplate)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (self *Resolver) resolveMenuOptions(prompt *config.CustomCommandPrompt, resolveTemplate func(string) (string, error)) ([]config.CustomCommandMenuOption, error) {
	newOptions := make([]config.CustomCommandMenuOption, 0, len(prompt.Options))
	for _, option := range prompt.Options {
		newOption, err := self.resolveMenuOption(&option, resolveTemplate)
		if err != nil {
			return nil, err
		}
		newOptions = append(newOptions, *newOption)
	}

	return newOptions, nil
}

func (self *Resolver) resolveMenuOption(option *config.CustomCommandMenuOption, resolveTemplate func(string) (string, error)) (*config.CustomCommandMenuOption, error) {
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

type CustomCommandObject struct {
	// deprecated. Use Responses instead
	PromptResponses []string
	Form            map[string]string
}
