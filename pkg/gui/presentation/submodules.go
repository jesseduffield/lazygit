package presentation

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

func GetSubmoduleListDisplayStrings(submodules []*models.SubmoduleConfig) [][]string {
	return lo.Map(submodules, func(submodule *models.SubmoduleConfig, _ int) []string {
		return getSubmoduleDisplayStrings(submodule)
	})
}

func getSubmoduleDisplayStrings(s *models.SubmoduleConfig) []string {
	name := s.Name
	if s.ParentModule != nil {
		indentation := ""
		for p := s.ParentModule; p != nil; p = p.ParentModule {
			indentation += "  "
		}

		name = indentation + "- " + s.Name
	}

	return []string{theme.DefaultTextColor.Sprint(name)}
}
