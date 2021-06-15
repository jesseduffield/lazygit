package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

//counterfeiter:generate . IBranchesMgr
type IBranchesMgr interface {
	NewBranch(name string, base string) error
	CurrentBranchName() (string, string, error)
	AllBranchesCmdObj() ICmdObj
	GetBranchGraphCmdObj(branchName string) ICmdObj
	Delete(branch string, force bool) error
	Merge(branchName string, opts MergeOpts) error
	Checkout(branch string, options CheckoutOpts) error
	GetUpstream(branchName string) (string, error)
	SetUpstream(upstream string, branchName string) error
	// RenameBranch(oldName string, newName string) error
	// FindRemoteForBranchInConfig(branchName string) (string, error)
	// AbortMerge() error
}

type BranchesMgr struct {
	commander ICommander
	config    IGitConfig
}

func NewBranchesMgr(commander ICommander, config IGitConfig) *BranchesMgr {
	return &BranchesMgr{
		commander: commander,
		config:    config,
	}
}

// NewBranch create new branch
func (c *BranchesMgr) NewBranch(name string, base string) error {
	return c.commander.RunGitCmdFromStr(fmt.Sprintf("checkout -b %s %s", name, base))
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *BranchesMgr) CurrentBranchName() (string, string, error) {
	branchName, err := c.commander.RunWithOutput(
		c.commander.BuildGitCmdObjFromStr("symbolic-ref --short HEAD"),
	)

	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}

	output, err := c.commander.RunWithOutput(
		c.commander.BuildGitCmdObjFromStr("branch --contains"),
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

func (c *BranchesMgr) AllBranchesCmdObj() ICmdObj {
	cmdStr := stripGitPrefixFromCmdStr(c.config.GetUserConfig().Git.AllBranchesLogCmd)

	return BuildGitCmdObjFromStr(cmdStr)
}

func (c *BranchesMgr) GetBranchGraphCmdObj(branchName string) ICmdObj {
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

func (c *BranchesMgr) Delete(branch string, force bool) error {
	forceFlag := "-d"
	if force {
		forceFlag = "-D"
	}

	return c.commander.RunGitCmdFromStr(fmt.Sprintf("branch %s %s", forceFlag, branch))
}

type MergeOpts struct {
	FastForwardOnly bool
}

// Merge merge
func (c *BranchesMgr) Merge(branchName string, opts MergeOpts) error {
	mergeArgs := c.config.GetUserConfig().Git.Merging.Args

	cmdStr := "merge --no-edit"
	if opts.FastForwardOnly {
		cmdStr += " --ff-only"
	}

	if mergeArgs != "" {
		cmdStr += " " + mergeArgs
	}

	cmdStr += " " + branchName

	return c.commander.RunGitCmdFromStr(cmdStr)
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOpts struct {
	Force   bool
	EnvVars []string
}

func (c *BranchesMgr) Checkout(branch string, options CheckoutOpts) error {
	forceArg := ""
	if options.Force {
		forceArg = " --force"
	}

	cmdObj := c.commander.BuildGitCmdObjFromStr(fmt.Sprintf("checkout%s %s", forceArg, branch))
	cmdObj.AddEnvVars(options.EnvVars...)

	return c.commander.Run(cmdObj)
}

func (c *BranchesMgr) GetUpstream(branchName string) (string, error) {
	output, err := c.commander.RunWithOutput(
		BuildGitCmdObjFromStr(fmt.Sprintf("rev-parse --abbrev-ref --symbolic-full-name %s@{u}", branchName)),
	)
	return strings.TrimSpace(output), err
}

// upstream is of the form remote/branchname
func (c *BranchesMgr) SetUpstream(upstream string, branchName string) error {
	return c.commander.RunGitCmdFromStr(fmt.Sprintf("branch --set-upstream-to=%s %s", upstream, branchName))
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
