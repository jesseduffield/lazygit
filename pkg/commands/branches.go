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
	return c.Cmd.New(fmt.Sprintf("git checkout -b %s %s", c.OSCommand.Quote(name), c.OSCommand.Quote(base))).Run()
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *GitCommand) CurrentBranchName() (string, string, error) {
	branchName, err := c.Cmd.New("git symbolic-ref --short HEAD").RunWithOutput()
	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}
	output, err := c.Cmd.New("git branch --contains").RunWithOutput()
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

	return c.Cmd.New(fmt.Sprintf("%s %s", command, c.OSCommand.Quote(branch))).Run()
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

	return c.Cmd.New(fmt.Sprintf("git checkout%s %s", forceArg, c.OSCommand.Quote(branch))).
		// prevents git from prompting us for input which would freeze the program
		// TODO: see if this is actually needed here
		AddEnvVars("GIT_TERMINAL_PROMPT=0").
		AddEnvVars(options.EnvVars...).
		Run()
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *GitCommand) GetBranchGraph(branchName string) (string, error) {
	return c.GetBranchGraphCmdObj(branchName).RunWithOutput()
}

func (c *GitCommand) GetUpstreamForBranch(branchName string) (string, error) {
	output, err := c.Cmd.New(fmt.Sprintf("git rev-parse --abbrev-ref --symbolic-full-name %s@{u}", c.OSCommand.Quote(branchName))).RunWithOutput()
	return strings.TrimSpace(output), err
}

func (c *GitCommand) GetBranchGraphCmdObj(branchName string) oscommands.ICmdObj {
	branchLogCmdTemplate := c.UserConfig.Git.BranchLogCmd
	templateValues := map[string]string{
		"branchName": c.OSCommand.Quote(branchName),
	}
	return c.Cmd.New(utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues))
}

func (c *GitCommand) SetUpstreamBranch(upstream string) error {
	return c.Cmd.New("git branch -u " + c.OSCommand.Quote(upstream)).Run()
}

func (c *GitCommand) SetBranchUpstream(remoteName string, remoteBranchName string, branchName string) error {
	return c.Cmd.New(fmt.Sprintf("git branch --set-upstream-to=%s/%s %s", c.OSCommand.Quote(remoteName), c.OSCommand.Quote(remoteBranchName), c.OSCommand.Quote(branchName))).Run()
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
	pushableCount, err := c.Cmd.New(fmt.Sprintf(command, to, from)).RunWithOutput()
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.Cmd.New(fmt.Sprintf(command, from, to)).RunWithOutput()
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
	mergeArgs := c.UserConfig.Git.Merging.Args

	command := fmt.Sprintf("git merge --no-edit %s %s", mergeArgs, c.OSCommand.Quote(branchName))
	if opts.FastForwardOnly {
		command = fmt.Sprintf("%s --ff-only", command)
	}

	return c.OSCommand.Cmd.New(command).Run()
}

// AbortMerge abort merge
func (c *GitCommand) AbortMerge() error {
	return c.Cmd.New("git merge --abort").Run()
}

func (c *GitCommand) IsHeadDetached() bool {
	err := c.Cmd.New("git symbolic-ref -q HEAD").Run()
	return err != nil
}

// ResetHardHead runs `git reset --hard`
func (c *GitCommand) ResetHard(ref string) error {
	return c.Cmd.New("git reset --hard " + c.OSCommand.Quote(ref)).Run()
}

// ResetSoft runs `git reset --soft HEAD`
func (c *GitCommand) ResetSoft(ref string) error {
	return c.Cmd.New("git reset --soft " + c.OSCommand.Quote(ref)).Run()
}

func (c *GitCommand) ResetMixed(ref string) error {
	return c.Cmd.New("git reset --mixed " + c.OSCommand.Quote(ref)).Run()
}

func (c *GitCommand) RenameBranch(oldName string, newName string) error {
	return c.Cmd.New(fmt.Sprintf("git branch --move %s %s", c.OSCommand.Quote(oldName), c.OSCommand.Quote(newName))).Run()
}

func (c *GitCommand) GetRawBranches() (string, error) {
	return c.Cmd.New(`git for-each-ref --sort=-committerdate --format="%(HEAD)|%(refname:short)|%(upstream:short)|%(upstream:track)" refs/heads`).RunWithOutput()
}
