package config

import (
	"fmt"
	"slices"
	"strings"
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
	return nil
}

func validateEnum(name string, value string, allowedValues []string) error {
	if slices.Contains(allowedValues, value) {
		return nil
	}
	allowedValuesStr := strings.Join(allowedValues, ", ")
	return fmt.Errorf("Unexpected value '%s' for '%s'. Allowed values: %s", value, name, allowedValuesStr)
}
