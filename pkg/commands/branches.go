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
	return c.RunCommand("git checkout -b %s %s", name, base)
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *GitCommand) CurrentBranchName() (string, string, error) {
	branchName, err := c.RunCommandWithOutput("git symbolic-ref --short HEAD")
	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}
	output, err := c.RunCommandWithOutput("git branch --contains")
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
	command := "git branch -d"

	if force {
		command = "git branch -D"
	}

	return c.OSCommand.RunCommand("%s %s", command, branch)
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOptions struct {
	Force   bool
	EnvVars []string
}

func (c *GitCommand) Checkout(branch string, options CheckoutOptions) error {
	forceArg := ""
	if options.Force {
		forceArg = " --force"
	}
	return c.OSCommand.RunCommandWithOptions(fmt.Sprintf("git checkout%s %s", forceArg, branch), oscommands.RunCommandOptions{EnvVars: options.EnvVars})
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *GitCommand) GetBranchGraph(branchName string) (string, error) {
	cmdStr := c.GetBranchGraphCmdStr(branchName)
	return c.OSCommand.RunCommandWithOutput(cmdStr)
}

func (c *GitCommand) GetUpstreamForBranch(branchName string) (string, error) {
	output, err := c.RunCommandWithOutput("git rev-parse --abbrev-ref --symbolic-full-name %s@{u}", branchName)
	return strings.TrimSpace(output), err
}

func (c *GitCommand) GetBranchGraphCmdStr(branchName string) string {
	branchLogCmdTemplate := c.Config.GetUserConfig().Git.BranchLogCmd
	templateValues := map[string]string{
		"branchName": branchName,
	}
	return utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues)
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

// TryGetMainBranch tries to get the main branch
// For old repo's this is usually master and for newer once it is main
// Note that this function might return a branch name that does not exist
func (c *GitCommand) TryGetMainBranch() string {
	// Firstly try to get the main branch from the origins
	cmdStr := `git branch --remotes --list '*/HEAD' --format '%(symref:short)'`
	output, _ := c.OSCommand.RunCommandWithOutput(cmdStr)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "/", 2)
		if len(parts) == 2 {
			return parts[len(parts)-1]
		}
	}

	// If we didn't find a main branch from the remotes we check if there is a main or master branch
	cmdStr = `git branch --format "%(refname:short)" --list "master" --list "main"`
	output, _ = c.OSCommand.RunCommandWithOutput(cmdStr)
	lines = append(lines[:0], strings.Split(strings.TrimSpace(output), "\n")...)

	// Some repo's use "master" while having a "main" branch
	// so we favor "master" over "main"
	foundMainBranch := false
	for _, line := range lines {
		if line == "master" {
			return line
		}
		if line == "main" {
			foundMainBranch = true
		}
	}
	if foundMainBranch {
		return "main"
	}

	return "master"
}

func (c *GitCommand) HasDevelopmentBranch() bool {
	cmdStr := `git branch --format "%(refname:short)" --list "develop"`
	output, _ := c.OSCommand.RunCommandWithOutput(cmdStr)
	return strings.TrimSpace(output) == "develop"
}

// GetCommitDifferences checks how many pushables/pullables there are for the
// current branch
func (c *GitCommand) GetCommitDifferences(from, to string) (string, string) {
	command := "git rev-list %s..%s --count"
	pushableCount, err := c.OSCommand.RunCommandWithOutput(command, to, from)
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.OSCommand.RunCommandWithOutput(command, from, to)
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
	mergeArgs := c.Config.GetUserConfig().Git.Merging.Args

	command := fmt.Sprintf("git merge --no-edit %s %s", mergeArgs, branchName)
	if opts.FastForwardOnly {
		command = fmt.Sprintf("%s --ff-only", command)
	}

	return c.OSCommand.RunCommand(command)
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
