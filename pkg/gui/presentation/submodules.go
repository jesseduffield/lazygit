package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
)

func GetSubmoduleListDisplayStrings(submodules []*models.SubmoduleConfig) [][]string {
	lines := make([][]string, len(submodules))

	for i := range submodules {
		lines[i] = getSubmoduleDisplayStrings(submodules[i])
	}

	return lines
}

func getSubmoduleDisplayStrings(s *models.SubmoduleConfig) []string {
	return []string{theme.DefaultTextColor.Sprint(s.Name)}
}
