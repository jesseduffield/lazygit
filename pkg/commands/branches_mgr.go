package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

//counterfeiter:generate . IBranchesMgr
type IBranchesMgr interface {
	NewBranch(name string, base string) error
	CurrentBranchName() (string, string, error)
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
	ICommander

	branchesLoader *BranchesLoader
	config         IGitConfigMgr
	log            *logrus.Entry
}

func NewBranchesMgr(commander ICommander, config IGitConfigMgr, log *logrus.Entry) *BranchesMgr {
	mgr := &BranchesMgr{
		ICommander: commander,
		config:     config,
	}

	mgr.branchesLoader = NewBranchesLoader(commander, mgr.CurrentBranchName, log)

	return mgr
}

func (c *BranchesMgr) LoadBranches(reflogCommits []*models.Commit) []*models.Branch {
	return c.branchesLoader.Load(reflogCommits)
}

// NewBranch create new branch
func (c *BranchesMgr) NewBranch(name string, base string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("checkout -b %s %s", name, base))
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *BranchesMgr) CurrentBranchName() (string, string, error) {
	branchName, err := c.RunWithOutput(
		c.BuildGitCmdObjFromStr("symbolic-ref --short HEAD"),
	)

	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}

	output, err := c.RunWithOutput(
		c.BuildGitCmdObjFromStr("branch --contains"),
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
