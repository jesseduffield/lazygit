package commands

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func (c *GitCommand) GetRemotes() ([]*Remote, error) {
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
	remotes := make([]*Remote, len(goGitRemotes))
	for i, goGitRemote := range goGitRemotes {
		remoteName := goGitRemote.Config().Name

		re := regexp.MustCompile(fmt.Sprintf(`%s\/([\S]+)`, remoteName))
		matches := re.FindAllStringSubmatch(remoteBranchesStr, -1)
		branches := make([]*RemoteBranch, len(matches))
		for j, match := range matches {
			branches[j] = &RemoteBranch{
				Name:       match[1],
				RemoteName: remoteName,
			}
		}

		remotes[i] = &Remote{
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
