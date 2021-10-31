package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

func (c *GitCommand) AddRemote(name string, url string) error {
	return c.RunCommand("git remote add %s %s", c.OSCommand.Quote(name), c.OSCommand.Quote(url))
}

func (c *GitCommand) RemoveRemote(name string) error {
	return c.RunCommand("git remote remove %s", c.OSCommand.Quote(name))
}

func (c *GitCommand) RenameRemote(oldRemoteName string, newRemoteName string) error {
	return c.RunCommand("git remote rename %s %s", c.OSCommand.Quote(oldRemoteName), c.OSCommand.Quote(newRemoteName))
}

func (c *GitCommand) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	return c.RunCommand("git remote set-url %s %s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(updatedUrl))
}

func (c *GitCommand) DeleteRemoteBranch(remoteName string, branchName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git push %s --delete %s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(branchName))
	cmdObj := c.NewCmdObjFromStr(command)
	return c.DetectUnamePass(cmdObj, promptUserForCredential)
}

func (c *GitCommand) DetectUnamePass(cmdObj oscommands.ICmdObj, promptUserForCredential func(string) string) error {
	return c.OSCommand.DetectUnamePass(cmdObj, c.GetCmdWriter(), promptUserForCredential)
}

// CheckRemoteBranchExists Returns remote branch
func (c *GitCommand) CheckRemoteBranchExists(branchName string) bool {
	_, err := c.OSCommand.RunCommandWithOutput(
		"git show-ref --verify -- refs/remotes/origin/%s",
		c.OSCommand.Quote(branchName),
	)

	return err == nil
}

// GetRemoteURL returns current repo remote url
func (c *GitCommand) GetRemoteURL() string {
	return c.GitConfig.Get("remote.origin.url")
}

func GetRepoInfoFromURL(url string) *RepoInformation {
	isHTTP := strings.HasPrefix(url, "http")

	if isHTTP {
		splits := strings.Split(url, "/")
		owner := strings.Join(splits[3:len(splits)-1], "/")
		repo := strings.TrimSuffix(splits[len(splits)-1], ".git")

		return &RepoInformation{
			Owner:      owner,
			Repository: repo,
		}
	}

	tmpSplit := strings.Split(url, ":")
	splits := strings.Split(tmpSplit[1], "/")
	owner := strings.Join(splits[0:len(splits)-1], "/")
	repo := strings.TrimSuffix(splits[len(splits)-1], ".git")

	return &RepoInformation{
		Owner:      owner,
		Repository: repo,
	}
}

func (c *GitCommand) GetRemotesToOwnersMap() (map[string]string, error) {
	remotes, err := c.GetRemotes()
	if err != nil {
		return nil, err
	}

	res := map[string]string{}
	for _, remote := range remotes {
		res[remote.Name] = GetRepoInfoFromURL(remote.Urls[0]).Owner
	}
	return res, nil
}
