package commands

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// GetReflogCommits only returns the new reflog commits since the given lastReflogCommit
// if none is passed (i.e. it's value is nil) then we get all the reflog commits
func (c *GitCommand) GetReflogCommits(lastReflogCommit *models.Commit, filterPath string) ([]*models.Commit, bool, error) {
	commits := make([]*models.Commit, 0)
	re := regexp.MustCompile(`(\w+).*HEAD@\{([^\}]+)\}: (.*)`)

	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" --follow -- %s", c.GetOSCommand().Quote(filterPath))
	}

	cmdObj := BuildGitCmdObjFromStr(fmt.Sprintf("reflog --abbrev=20 --date=unix %s", filterPathArg))
	onlyObtainedNewReflogCommits := false
	err := oscommands.RunLineOutputCmd(cmdObj, func(line string) (bool, error) {
		match := re.FindStringSubmatch(line)
		if len(match) <= 1 {
			return false, nil
		}

		unixTimestamp, _ := strconv.Atoi(match[2])

		commit := &models.Commit{
			Sha:           match[1],
			Name:          match[3],
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
