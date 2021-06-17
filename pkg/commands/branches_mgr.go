package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// this takes something like:
// * (HEAD detached at 264fc6f5)
//	remotes
// and returns '264fc6f5' as the second match
const CurrentBranchNameRegex = `(?m)^\*.*?([^ ]*?)\)?$`

//counterfeiter:generate . IBranchesMgr
type IBranchesMgr interface {
	NewBranch(name string, base string) error
	AllBranchesCmdObj() ICmdObj
	GetBranchGraphCmdObj(branchName string) ICmdObj
	Delete(branch string, force bool) error
	Merge(branchName string, opts MergeOpts) error
	AbortMerge() error
	Checkout(branch string, options CheckoutOpts) error
	GetUpstream(branchName string) (string, error)
	SetUpstream(upstream string, branchName string) error
	RenameBranch(oldName string, newName string) error
	ResetToRef(ref string, strength ResetStrength, opts ResetToRefOpts) error
	LoadBranches(reflogCommits []*models.Commit) []*models.Branch
}

type BranchesMgr struct {
	*MgrCtx

	statusMgr      IStatusMgr
	branchesLoader *BranchesLoader
}

func NewBranchesMgr(mgrCtx *MgrCtx, statusMgr IStatusMgr) *BranchesMgr {
	mgr := &BranchesMgr{MgrCtx: mgrCtx}

	mgr.branchesLoader = NewBranchesLoader(mgrCtx, statusMgr)

	return mgr
}

func (c *BranchesMgr) LoadBranches(reflogCommits []*models.Commit) []*models.Branch {
	return c.branchesLoader.Load(reflogCommits)
}

// NewBranch create new branch
func (c *BranchesMgr) NewBranch(name string, base string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("checkout -b %s %s", name, base))
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

	return c.RunGitCmdFromStr(fmt.Sprintf("branch %s %s", forceFlag, branch))
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

	return c.RunGitCmdFromStr(cmdStr)
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

	cmdObj := c.BuildGitCmdObjFromStr(fmt.Sprintf("checkout%s %s", forceArg, branch))
	cmdObj.AddEnvVars(options.EnvVars...)

	return c.Run(cmdObj)
}

func (c *BranchesMgr) GetUpstream(branchName string) (string, error) {
	output, err := c.RunWithOutput(
		BuildGitCmdObjFromStr(fmt.Sprintf("rev-parse --abbrev-ref --symbolic-full-name %s@{u}", branchName)),
	)
	return strings.TrimSpace(output), err
}

// upstream is of the form remote/branchname
func (c *BranchesMgr) SetUpstream(upstream string, branchName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("branch --set-upstream-to=%s %s", upstream, branchName))
}

func (c *BranchesMgr) RenameBranch(oldName string, newName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("branch --move %s %s", oldName, newName))
}

func (c *BranchesMgr) AbortMerge() error {
	return c.RunGitCmdFromStr("merge --abort")
}

type ResetStrength string

const (
	SOFT  ResetStrength = "soft"
	MIXED               = "mixed"
	HARD                = "hard"
)

type ResetToRefOpts struct {
	EnvVars []string
}

// ResetToCommit reset to commit
func (c *BranchesMgr) ResetToRef(ref string, strength ResetStrength, options ResetToRefOpts) error {
	cmdObj := BuildGitCmdObjFromStr(fmt.Sprintf("reset --%s %s", string(strength), ref))
	cmdObj.AddEnvVars(options.EnvVars...)

	return c.Run(cmdObj)
}
