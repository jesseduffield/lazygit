package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func GetSubmoduleListDisplayStrings(submodules []*models.SubmoduleConfig) [][]string {
	lines := make([][]string, len(submodules))

	for i := range submodules {
		lines[i] = getSubmoduleDisplayStrings(submodules[i])
	}

	return lines
}

func getSubmoduleDisplayStrings(s *models.SubmoduleConfig) []string {
	return []string{utils.ColoredString(s.Name, theme.DefaultTextColor)}
}
