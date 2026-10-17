package presentation

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/samber/lo"
)

func GetRepoListDisplayStrings(repos []*models.Repo, currentPath string) [][]string {
	return lo.Map(repos, func(repo *models.Repo, _ int) []string {
		current := ""
		if IsCurrentRepo(repo, currentPath) {
			current = "  *"
		}
		branch := style.FgCyan.Sprint(repo.Branch)
		if repo.Dirty {
			branch += style.FgYellow.Sprint("*")
		}
		return []string{
			style.FgGreen.Sprint(current),
			theme.DefaultTextColor.Sprint(repo.Name),
			branch,
		}
	})
}

// IsCurrentRepo reports whether repo is the repo at currentPath. currentPath
// comes from git, which writes forward slashes on Windows too.
func IsCurrentRepo(repo *models.Repo, currentPath string) bool {
	return repo.Path == filepath.Clean(filepath.FromSlash(currentPath))
}
