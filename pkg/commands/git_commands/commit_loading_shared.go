package git_commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
)

func loadCommits(
	cmd *oscommands.CmdObj,
	filterPath string,
	parseLogLine func(string) (*models.Commit, bool),
) ([]*models.Commit, error) {
	commits := []*models.Commit{}

	var commit *models.Commit
	var filterPaths []string
	// A string pool that stores interned strings to reduce memory usage
	pool := make(map[string]string)

	finishLastCommit := func() {
		if commit != nil {
			// Only set the filter paths if we have one that is not contained in the original
			// filter path. When filtering on a directory, all file paths will start with that
			// directory, so we needn't bother storing the individual paths. Likewise, if we
			// filter on a file and the file path hasn't changed, we needn't store it either.
			// Only if a file has been moved or renamed do we need to store the paths, but then
			// we need them all so that we can properly render a diff for the rename.
			if lo.SomeBy(filterPaths, func(path string) bool {
				return !strings.HasPrefix(path, filterPath)
			}) {
				commit.FilterPaths = lo.Map(filterPaths, func(path string, _ int) string {
					if v, ok := pool[path]; ok {
						return v
					}
					pool[path] = path
					return path
				})
			}
			commits = append(commits, commit)
			commit = nil
			filterPaths = nil
		}
	}
	err := cmd.RunAndProcessLines(func(line string) (bool, error) {
		if line == "" {
			return false, nil
		}

		if line[0] == '+' {
			finishLastCommit()
			var stop bool
			commit, stop = parseLogLine(line[1:])
			if stop {
				commit = nil
				return true, nil
			}
		} else if commit != nil && filterPath != "" {
			// We are filtering by path, and this line is the output of the --name-status flag
			fields := strings.Split(line, "\t")
			// We don't bother looking at the first field (it will be 'A', 'M', 'R072' or a bunch of others).
			// All we care about is the path(s), and there will be one for 'M' and 'A', and two for 'R' or 'C',
			// in which case we want them both so that we can show the diff between the two.
			if len(fields) > 1 {
				filterPaths = append(filterPaths, fields[1:]...)
			}
		}
		return false, nil
	})
	finishLastCommit()
	return commits, err
}
