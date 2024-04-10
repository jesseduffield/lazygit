package presentation

import (
	"github.com/lobes/lazytask/pkg/commands/models"
	"github.com/lobes/lazytask/pkg/theme"
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
