package commands

import (
	"fmt"
	"regexp"
)

func (c *GitCommand) GetRemotes() ([]*Remote, error) {
	// get remote branches
	remoteBranchesStr, err := c.OSCommand.RunCommandWithOutput("git for-each-ref --format='%(refname:strip=2)' refs/remotes")
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
		name := goGitRemote.Config().Name

		re := regexp.MustCompile(fmt.Sprintf("%s\\/(.*)", name))
		matches := re.FindAllStringSubmatch(remoteBranchesStr, -1)
		branches := make([]*Branch, len(matches))
		for j, match := range matches {
			branches[j] = &Branch{
				Name: match[1],
			}
		}

		remotes[i] = &Remote{
			Name:     goGitRemote.Config().Name,
			Urls:     goGitRemote.Config().URLs,
			Branches: branches,
		}
	}

	return remotes, nil
}
