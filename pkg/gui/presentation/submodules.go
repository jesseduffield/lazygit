package presentation

import (
	"fmt"

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
	// Pad right with some spaces to the end of the HEAD so that it (hopefully) aligns well.
	// Put the HEAD first because those are more likely to be similar lengths than the repo name.
	name := fmt.Sprintf("%-20s %s", s.Head,
		s.Name,
	)
	if s.ParentModule != nil {
		indentation := ""
		for p := s.ParentModule; p != nil; p = p.ParentModule {
			indentation += "  "
		}

		name = indentation + "- " + s.Name
	}

	if s.NumStagedFiles != 0 {
		name = fmt.Sprintf(
			"%s +%d",
			name,
			s.NumStagedFiles,
		)
	}

	if s.NumUnstagedChanges != 0 {
		name = fmt.Sprintf(
			"%s !%d",
			name,
			s.NumUnstagedChanges,
		)
	}

	if s.NumUntrackedChanges != 0 {
		name = fmt.Sprintf(
			"%s ?%d ",
			name,
			s.NumUntrackedChanges,
		)
	}

	return []string{theme.DefaultTextColor.Sprint(name)}
}
