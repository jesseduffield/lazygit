package git_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

func loadCommits(
	cmd *oscommands.CmdObj,
	parseLogLine func(string) (*models.Commit, bool),
) ([]*models.Commit, error) {
	commits := []*models.Commit{}

	err := cmd.RunAndProcessLines(func(line string) (bool, error) {
		commit, stop := parseLogLine(line)
		if stop {
			return true, nil
		}
		commits = append(commits, commit)
		return false, nil
	})
	return commits, err
}
