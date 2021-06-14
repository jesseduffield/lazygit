package commands

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func (c *GitCommand) GetRemotes() ([]*models.Remote, error) {
	// get remote branches
	remoteBranchesStr, err := c.RunWithOutput(
		BuildGitCmdObjFromStr("branch -r"),
	)
	if err != nil {
		return nil, err
	}

	goGitRemotes, err := c.repo.Remotes()
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
			branches[j] = &models.RemoteBranch{
				Name:       match[1],
				RemoteName: remoteName,
			}
		}

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
