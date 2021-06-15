package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// NewBranch create new branch
func (c *Git) NewBranch(name string, base string) error {
	return c.Run(
		BuildGitCmdObj("checkout", []string{name, base}, map[string]bool{"-b": true}),
	)
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *Git) CurrentBranchName() (string, string, error) {
	branchName, err := c.RunWithOutput(
		BuildGitCmdObjFromStr("symbolic-ref --short HEAD"),
	)

	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}
	output, err := c.RunWithOutput(
		BuildGitCmdObjFromStr("branch --contains"),
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
func (c *Git) DeleteBranch(branch string, force bool) error {
	return c.Run(
		BuildGitCmdObj("branch", []string{branch}, map[string]bool{"-d": !force, "-D": force}),
	)
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOptions struct {
	Force   bool
	EnvVars []string
}

func (c *Git) Checkout(branch string, options CheckoutOptions) error {
	cmdObj := BuildGitCmdObj("checkout", []string{branch}, map[string]bool{"--force": options.Force})
	cmdObj.AddEnvVars(options.EnvVars...)

	return c.Run(cmdObj)
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *Git) GetBranchGraph(branchName string) (string, error) {
	return c.RunWithOutput(c.GetBranchGraphCmdObj(branchName))
}

func (c *Git) GetUpstreamForBranch(branchName string) (string, error) {
	output, err := c.RunWithOutput(BuildGitCmdObjFromStr(fmt.Sprintf("rev-parse --abbrev-ref --symbolic-full-name %s@{u}", branchName)))
	return strings.TrimSpace(output), err
}

func (c *Git) GetBranchGraphCmdObj(branchName string) ICmdObj {
	branchLogCmdTemplate := c.config.GetUserConfig().Git.BranchLogCmd
	templateValues := map[string]string{
		"branchName": branchName,
	}
	cmdObj := oscommands.NewCmdObjFromStr(
		utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues),
	)
	SetDefaultEnvVars(cmdObj)

	return cmdObj
}

func (c *Git) SetUpstreamBranch(upstream string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("branch -u %s", upstream))
}

func (c *Git) SetBranchUpstream(remoteName string, remoteBranchName string, branchName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("branch --set-upstream-to=%s/%s %s", remoteName, remoteBranchName, branchName))
}

func (c *Git) GetCurrentBranchUpstreamDifferenceCount() (string, string) {
	return c.GetBranchUpstreamDifferenceCount("HEAD")
}

func (c *Git) GetBranchUpstreamDifferenceCount(branchName string) (string, string) {
	return c.GetCommitDifferences(branchName, branchName+"@{u}")
}

// GetCommitDifferences checks how many pushables/pullables there are for the
// current branch
func (c *Git) GetCommitDifferences(from, to string) (string, string) {
	pushableCount, err := c.GetCommitDifference(to, from)
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.GetCommitDifference(from, to)
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

func (c *Git) GetCommitDifference(from string, to string) (string, error) {
	return c.RunWithOutput(BuildGitCmdObjFromStr(fmt.Sprintf("rev-list %s..%s --count", from, to)))
}

type MergeOpts struct {
	FastForwardOnly bool
}

// Merge merge
func (c *Git) Merge(branchName string, opts MergeOpts) error {
	mergeArgs := c.config.GetUserConfig().Git.Merging.Args

	cmdStr := fmt.Sprintf("merge --no-edit %s %s", mergeArgs, branchName)
	if opts.FastForwardOnly {
		cmdStr = fmt.Sprintf("%s --ff-only", cmdStr)
	}

	return c.RunGitCmdFromStr(cmdStr)
}

// AbortMerge abort merge
func (c *Git) AbortMerge() error {
	return c.RunGitCmdFromStr("merge --abort")
}

func (c *Git) IsHeadDetached() bool {
	err := c.RunGitCmdFromStr("symbolic-ref -q HEAD")
	return err != nil
}

// ResetHardHead runs `git reset --hard`
func (c *Git) ResetHard(ref string) error {
	return c.RunGitCmdFromStr("reset --hard " + ref)
}

// ResetSoft runs `git reset --soft HEAD`
func (c *Git) ResetSoft(ref string) error {
	return c.RunGitCmdFromStr("reset --soft " + ref)
}

func (c *Git) ResetMixed(ref string) error {
	return c.RunGitCmdFromStr("reset --mixed " + ref)
}

func (c *Git) RenameBranch(oldName string, newName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("branch --move %s %s", oldName, newName))
}

// ResetToCommit reset to commit
func (c *Git) ResetToRef(ref string, strength string, options ResetToCommitOptions) error {
	cmdObj := BuildGitCmdObjFromStr(fmt.Sprintf("reset --%s %s", strength, ref))
	cmdObj.AddEnvVars(options.EnvVars...)

	return c.Run(cmdObj)
}
