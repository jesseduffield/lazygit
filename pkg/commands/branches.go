package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewBranch create new branch
func (c *GitCommand) NewBranch(name string, base string) error {
	return c.RunExecutable(
		BuildGitCmdObj("checkout", []string{name, base}, map[string]bool{"-b": true}),
	)
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *GitCommand) CurrentBranchName() (string, string, error) {
	branchName, err := c.RunExecutableWithOutput(
		BuildGitCmdObj("symbolic-ref", []string{"HEAD"}, map[string]bool{"--short": true}),
	)
	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}
	output, err := c.RunExecutableWithOutput(
		BuildGitCmdObj("branch", nil, map[string]bool{"--contains": true}),
	)
	if err != nil {
		return "", "", err
	}
	for _, line := range utils.SplitLines(output) {
		re := regexp.MustCompile(CurrentBranchNameRegex)
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			branchName = match[1]
			displayBranchName := match[0][2:]
			return branchName, displayBranchName, nil
		}
	}
	return "HEAD", "HEAD", nil
}

// DeleteBranch delete branch
func (c *GitCommand) DeleteBranch(branch string, force bool) error {
	return c.GetOSCommand().RunExecutable(
		BuildGitCmdObj("branch", []string{branch}, map[string]bool{"-d": !force, "-D": force}),
	)
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOptions struct {
	Force   bool
	EnvVars []string
}

func (c *GitCommand) Checkout(branch string, options CheckoutOptions) error {
	cmdObj := BuildGitCmdObj("checkout", []string{branch}, map[string]bool{"--force": options.Force})
	cmdObj.AddEnvVars(options.EnvVars...)

	return c.GetOSCommand().RunCommandWithOptions(cmdObj)
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *GitCommand) GetBranchGraph(branchName string) (string, error) {
	return c.GetOSCommand().RunExecutableWithOutput(c.GetBranchGraphCmdObj(branchName))
}

func (c *GitCommand) GetUpstreamForBranch(branchName string) (string, error) {
	output, err := c.RunCommandWithOutput("git rev-parse --abbrev-ref --symbolic-full-name %s@{u}", branchName)
	return strings.TrimSpace(output), err
}

func (c *GitCommand) GetBranchGraphCmdObj(branchName string) *oscommands.CmdObj {
	branchLogCmdTemplate := c.config.GetUserConfig().Git.BranchLogCmd
	templateValues := map[string]string{
		"branchName": branchName,
	}
	str := utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues)
	cmdObj := &oscommands.CmdObj{CmdStr: str}
	DisableOptionalLocks(cmdObj)

	return cmdObj
}

func (c *GitCommand) SetUpstreamBranch(upstream string) error {
	return c.RunCommand("git branch -u %s", upstream)
}

func (c *GitCommand) SetBranchUpstream(remoteName string, remoteBranchName string, branchName string) error {
	return c.RunCommand("git branch --set-upstream-to=%s/%s %s", remoteName, remoteBranchName, branchName)
}

func (c *GitCommand) GetCurrentBranchUpstreamDifferenceCount() (string, string) {
	return c.GetCommitDifferences("HEAD", "HEAD@{u}")
}

func (c *GitCommand) GetBranchUpstreamDifferenceCount(branchName string) (string, string) {
	return c.GetCommitDifferences(branchName, branchName+"@{u}")
}

// GetCommitDifferences checks how many pushables/pullables there are for the
// current branch
func (c *GitCommand) GetCommitDifferences(from, to string) (string, string) {
	command := "git rev-list %s..%s --count"
	pushableCount, err := c.GetOSCommand().RunCommandWithOutput(command, to, from)
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.GetOSCommand().RunCommandWithOutput(command, from, to)
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

type MergeOpts struct {
	FastForwardOnly bool
}

// Merge merge
func (c *GitCommand) Merge(branchName string, opts MergeOpts) error {
	mergeArgs := c.config.GetUserConfig().Git.Merging.Args

	command := fmt.Sprintf("git merge --no-edit %s %s", mergeArgs, branchName)
	if opts.FastForwardOnly {
		command = fmt.Sprintf("%s --ff-only", command)
	}

	return c.GetOSCommand().RunCommand(command)
}

// AbortMerge abort merge
func (c *GitCommand) AbortMerge() error {
	return c.RunCommand("git merge --abort")
}

func (c *GitCommand) IsHeadDetached() bool {
	err := c.RunCommand("git symbolic-ref -q HEAD")
	return err != nil
}

// ResetHardHead runs `git reset --hard`
func (c *GitCommand) ResetHard(ref string) error {
	return c.RunCommand("git reset --hard " + ref)
}

// ResetSoft runs `git reset --soft HEAD`
func (c *GitCommand) ResetSoft(ref string) error {
	return c.RunCommand("git reset --soft " + ref)
}

func (c *GitCommand) ResetMixed(ref string) error {
	return c.RunCommand("git reset --mixed " + ref)
}

func (c *GitCommand) RenameBranch(oldName string, newName string) error {
	return c.RunCommand("git branch --move %s %s", oldName, newName)
}
