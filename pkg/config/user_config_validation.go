package config

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/constants"
)

func (config *UserConfig) Validate() error {
	if err := validateEnum("gui.statusPanelView", config.Gui.StatusPanelView,
		[]string{"dashboard", "allBranchesLog"}); err != nil {
		return err
	}
	if err := validateEnum("gui.showDivergenceFromBaseBranch", config.Gui.ShowDivergenceFromBaseBranch,
		[]string{"none", "onlyArrow", "arrowAndNumber"}); err != nil {
		return err
	}
	if err := validateEnum("git.autoForwardBranches", config.Git.AutoForwardBranches,
		[]string{"none", "onlyMainBranches", "allBranches"}); err != nil {
		return err
	}
	if err := validateEnum("git.localBranchSortOrder", config.Git.LocalBranchSortOrder,
		[]string{"date", "recency", "alphabetical"}); err != nil {
		return err
	}
	if err := validateEnum("git.remoteBranchSortOrder", config.Git.RemoteBranchSortOrder,
		[]string{"date", "alphabetical"}); err != nil {
		return err
	}
	if err := validateEnum("git.log.order", config.Git.Log.Order,
		[]string{"date-order", "author-date-order", "topo-order", "default"}); err != nil {
		return err
	}
	if err := validateEnum("git.log.showGraph", config.Git.Log.ShowGraph,
		[]string{"always", "never", "when-maximised"}); err != nil {
		return err
	}
	if err := validateKeybindings(config.Keybinding); err != nil {
		return err
	}
	if err := validateCustomCommands(config.CustomCommands); err != nil {
		return err
	}
	return nil
}

func validateEnum(name string, value string, allowedValues []string) error {
	if slices.Contains(allowedValues, value) {
		return nil
	}
	allowedValuesStr := strings.Join(allowedValues, ", ")
	return fmt.Errorf("Unexpected value '%s' for '%s'. Allowed values: %s", value, name, allowedValuesStr)
}

func validateKeybindingsRecurse(path string, node any) error {
	value := reflect.ValueOf(node)
	if value.Kind() == reflect.Struct {
		for _, field := range reflect.VisibleFields(reflect.TypeOf(node)) {
			var newPath string
			if len(path) == 0 {
				newPath = field.Name
			} else {
				newPath = fmt.Sprintf("%s.%s", path, field.Name)
			}
			if err := validateKeybindingsRecurse(newPath,
				value.FieldByName(field.Name).Interface()); err != nil {
				return err
			}
		}
	} else if value.Kind() == reflect.Slice {
		for i := range value.Len() {
			if err := validateKeybindingsRecurse(
				fmt.Sprintf("%s[%d]", path, i), value.Index(i).Interface()); err != nil {
				return err
			}
		}
	} else if value.Kind() == reflect.String {
		key := node.(string)
		if !isValidKeybindingKey(key) {
			return fmt.Errorf("Unrecognized key '%s' for keybinding '%s'. For permitted values see %s",
				key, path, constants.Links.Docs.CustomKeybindings)
		}
	} else {
		log.Fatalf("Unexpected type for property '%s': %s", path, value.Kind())
	}
	return nil
}

func validateKeybindings(keybindingConfig KeybindingConfig) error {
	if err := validateKeybindingsRecurse("", keybindingConfig); err != nil {
		return err
	}

	if len(keybindingConfig.Universal.JumpToBlock) != 5 {
		return fmt.Errorf("keybinding.universal.jumpToBlock must have 5 elements; found %d.",
			len(keybindingConfig.Universal.JumpToBlock))
	}

	return nil
}

func validateCustomCommandKey(key string) error {
	if !isValidKeybindingKey(key) {
		return fmt.Errorf("Unrecognized key '%s' for custom command. For permitted values see %s",
			key, constants.Links.Docs.CustomKeybindings)
	}
	return nil
}

func validateCustomCommands(customCommands []CustomCommand) error {
	for _, customCommand := range customCommands {
		if err := validateCustomCommandKey(customCommand.Key); err != nil {
			return err
		}

		if len(customCommand.CommandMenu) > 0 {
			if len(customCommand.Context) > 0 ||
				len(customCommand.Command) > 0 ||
				len(customCommand.Prompts) > 0 ||
				len(customCommand.LoadingText) > 0 ||
				len(customCommand.Output) > 0 ||
				len(customCommand.OutputTitle) > 0 ||
				customCommand.After != nil {
				commandRef := ""
				if len(customCommand.Key) > 0 {
					commandRef = fmt.Sprintf(" with key '%s'", customCommand.Key)
				}
				return fmt.Errorf("Error with custom command%s: it is not allowed to use both commandMenu and any of the other fields except key and description.", commandRef)
			}

			if err := validateCustomCommands(customCommand.CommandMenu); err != nil {
				return err
			}
		} else {
			if err := validateEnum("customCommand.output", customCommand.Output,
				[]string{"", "none", "terminal", "log", "logWithPty", "popup"}); err != nil {
				return err
			}
		}
	}
	return nil
}
