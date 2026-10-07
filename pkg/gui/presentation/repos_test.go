package presentation

import (
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetRepoListDisplayStrings(t *testing.T) {
	repos := []*models.Repo{
		{Path: "/w/alpha", Name: "alpha", Branch: "feature"},
		{Path: "/w/beta", Name: "beta", Branch: "master", Dirty: true},
	}

	rows := GetRepoListDisplayStrings(repos, "/w/alpha")

	decolorised := lo.Map(rows, func(row []string, _ int) []string {
		return lo.Map(row, func(cell string, _ int) string { return utils.Decolorise(cell) })
	})
	assert.Equal(t, [][]string{
		{"  *", "alpha", "feature"},
		{"", "beta", "master*"},
	}, decolorised)
}

func TestGetRepoListDisplayStringsGitStylePath(t *testing.T) {
	repos := []*models.Repo{
		{Path: filepath.FromSlash("/w/alpha"), Name: "alpha", Branch: "feature"},
		{Path: filepath.FromSlash("/w/beta"), Name: "beta", Branch: "master"},
	}

	rows := GetRepoListDisplayStrings(repos, "/w/beta/")

	assert.Equal(t, []string{"", "  *"}, lo.Map(rows, func(row []string, _ int) string {
		return utils.Decolorise(row[0])
	}))
}
