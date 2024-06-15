package config

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

func (config *UserConfig) Validate() error {
	if strings.HasPrefix(string(config.Gui.CommitAuthorFormat), "truncateTo:") {
		regex := regexp.MustCompile(`truncateTo:\d+(,\d+)?$`)
		if !regex.MatchString(string(config.Gui.CommitAuthorFormat)) {
			return fmt.Errorf("Invalid value for 'gui.commitAuthorFormat'. Expected format: 'truncateTo:<normal_length>[,<extended_length>]'")
		}
	} else if err := validateEnum("gui.commitAuthorFormat", string(config.Gui.CommitAuthorFormat),
		[]string{"auto", "short", "full"}); err != nil {
		return err
	}

	if err := validateEnum("gui.statusPanelView", config.Gui.StatusPanelView,
		[]string{"dashboard", "allBranchesLog"}); err != nil {
		return err
	}
	if err := validateEnum("gui.showDivergenceFromBaseBranch", config.Gui.ShowDivergenceFromBaseBranch,
		[]string{"none", "onlyArrow", "arrowAndNumber"}); err != nil {
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
