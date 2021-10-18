package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// GetReflogCommits only returns the new reflog commits since the given lastReflogCommit
// if none is passed (i.e. it's value is nil) then we get all the reflog commits
func (c *GitCommand) GetReflogCommits(lastReflogCommit *models.Commit, filterPath string) ([]*models.Commit, bool, error) {
	commits := make([]*models.Commit, 0)

	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" --follow -- %s", c.OSCommand.Quote(filterPath))
	}

	cmd := c.OSCommand.ExecutableFromString(fmt.Sprintf(`git log -g --abbrev=20 --format="%%h %%ct %%gs" %s`, filterPathArg))
	onlyObtainedNewReflogCommits := false
	err := oscommands.RunLineOutputCmd(cmd, func(line string) (bool, error) {
		fields := strings.SplitN(line, " ", 3)
		if len(fields) <= 2 {
			return false, nil
		}

		unixTimestamp, _ := strconv.Atoi(fields[1])

		commit := &models.Commit{
			Sha:           fields[0],
			Name:          fields[2],
			UnixTimestamp: int64(unixTimestamp),
			Status:        "reflog",
		}

		if lastReflogCommit != nil && commit.Sha == lastReflogCommit.Sha && commit.UnixTimestamp == lastReflogCommit.UnixTimestamp {
			onlyObtainedNewReflogCommits = true
			// after this point we already have these reflogs loaded so we'll simply return the new ones
			return true, nil
		}

		commits = append(commits, commit)
		return false, nil
	})
	if err != nil {
		return nil, false, err
	}

	return commits, onlyObtainedNewReflogCommits, nil
}
