package commands

import (
	"fmt"
	"strings"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type PushOpts struct {
	Force                   bool
	PromptUserForCredential func(CredentialKind) string
	SetUpstream             bool
	DestinationRemote       string
	DestinationBranch       string
}

func (c *GitCommand) Push(opts PushOpts) error {
	cmd := buildGitCmd("push", []string{opts.DestinationRemote, opts.DestinationBranch},
		map[string]bool{
			"--follow-tags":      c.GetConfigValue("push.followTags") != "false",
			"--force-with-lease": opts.Force,
			"--set-upstream":     opts.SetUpstream,
		})

	return c.GetOSCommand().DetectUnamePass(cmd, opts.PromptUserForCredential)
}

func buildGitCmd(command string, positionalArgs []string, kwArgs map[string]bool) string {
	parts := []string{"git", command}

	args := make([]string, 0, len(kwArgs))
	for arg, include := range kwArgs {
		if include {
			args = append(args, arg)
		}
	}
	utils.SortAlphabeticalInPlace(args)

	presentPosArgs := utils.ExcludeEmpty(positionalArgs)

	parts = append(parts, presentPosArgs...)
	parts = append(parts, args...)

	return strings.Join(parts, " ")
}

type FetchOptions struct {
	PromptUserForCredential func(CredentialKind) string
	RemoteName              string
	BranchName              string
}

// Fetch fetch git repo
func (c *GitCommand) Fetch(opts FetchOptions) error {
	command := "git fetch"

	if opts.RemoteName != "" {
		command = fmt.Sprintf("%s %s", command, opts.RemoteName)
	}
	if opts.BranchName != "" {
		command = fmt.Sprintf("%s %s", command, opts.BranchName)
	}

	return c.GetOSCommand().DetectUnamePass(command, func(question CredentialKind) string {
		if opts.PromptUserForCredential != nil {
			return opts.PromptUserForCredential(question)
		}
		return "\n"
	})
}

func (c *GitCommand) FastForward(branchName string, remoteName string, remoteBranchName string, promptUserForCredential func(CredentialKind) string) error {
	command := fmt.Sprintf("git fetch %s %s:%s", remoteName, remoteBranchName, branchName)
	return c.GetOSCommand().DetectUnamePass(command, promptUserForCredential)
}

func (c *GitCommand) FetchRemote(remoteName string, promptUserForCredential func(CredentialKind) string) error {
	command := fmt.Sprintf("git fetch %s", remoteName)
	return c.GetOSCommand().DetectUnamePass(command, promptUserForCredential)
}
