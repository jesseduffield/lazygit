package commands

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (c *GitCommand) GetRemotes() ([]*models.Remote, error) {
	// get remote branches
	unescaped := "git branch -r"
	remoteBranchesStr, err := c.OSCommand.RunCommandWithOutput(unescaped)
	if err != nil {
		return nil, err
	}

	goGitRemotes, err := c.Repo.Remotes()
	if err != nil {
		return nil, err
	}

	// first step is to get our remotes from go-git
	remotes := make([]*models.Remote, len(goGitRemotes))
	for i, goGitRemote := range goGitRemotes {
		remoteName := goGitRemote.Config().Name

		re := regexp.MustCompile(fmt.Sprintf(`%s\/([\S]+)`, remoteName))
		matches := re.FindAllStringSubmatch(remoteBranchesStr, -1)
		branches := make([]*models.RemoteBranch, len(matches))

		for j, match := range matches {
			branch := &models.RemoteBranch{
				Name:       match[1],
				RemoteName: remoteName,
			}

			unescaped := "git --no-pager log -1 " + branch.FullName() + " --format=%at"
			remoteBranchDateStr, err := c.OSCommand.RunCommandWithOutput(unescaped)
			if err != nil {
				return nil, err
			}

			i, err := strconv.ParseInt(strings.TrimSuffix(remoteBranchDateStr, "\n"), 10, 64)

			if err != nil {
				return nil, err
			}

			tm := time.Unix(i, 0)

			if err != nil {
				return nil, err
			}

			branch.LastCommitUnixTime = tm
			branches[j] = branch
		}

		// sort branches by commit date
		sort.Slice(branches, func(i, j int) bool {
			return (branches[j].LastCommitUnixTime.Before(branches[i].LastCommitUnixTime))
		})

		remotes[i] = &models.Remote{
			Name:     goGitRemote.Config().Name,
			Urls:     goGitRemote.Config().URLs,
			Branches: branches,
		}
	}

	// now lets sort our remotes by name alphabetically
	sort.Slice(remotes, func(i, j int) bool {
		// we want origin at the top because we'll be most likely to want it
		if remotes[i].Name == "origin" {
			return true
		}
		if remotes[j].Name == "origin" {
			return false
		}
		return strings.ToLower(remotes[i].Name) < strings.ToLower(remotes[j].Name)
	})

	return remotes, nil
}
