package git_commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// this takes something like:
// * (HEAD detached at 264fc6f5)
// remotes
// and returns '264fc6f5' as the second match
const CurrentBranchNameRegex = `(?m)^\*.*?([^ ]*?)\)?$`

type BranchCommands struct {
	*GitCommon
}

func NewBranchCommands(gitCommon *GitCommon) *BranchCommands {
	return &BranchCommands{
		GitCommon: gitCommon,
	}
}

// New creates a new branch
func (self *BranchCommands) New(name string, base string) error {
	return self.cmd.New(fmt.Sprintf("git checkout -b %s %s", self.cmd.Quote(name), self.cmd.Quote(base))).Run()
}

// CurrentBranchInfo get the current branch information.
func (self *BranchCommands) CurrentBranchInfo() (BranchInfo, error) {
	branchName, err := self.cmd.New("git symbolic-ref --short HEAD").DontLog().RunWithOutput()
	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return BranchInfo{
			RefName:      trimmedBranchName,
			DisplayName:  trimmedBranchName,
			DetachedHead: false,
		}, nil
	}
	output, err := self.cmd.New("git branch --contains").DontLog().RunWithOutput()
	if err != nil {
		return BranchInfo{}, err
	}
	for _, line := range utils.SplitLines(output) {
		re := regexp.MustCompile(CurrentBranchNameRegex)
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			branchName = match[1]
			displayBranchName := match[0][2:]
			return BranchInfo{
				RefName:      branchName,
				DisplayName:  displayBranchName,
				DetachedHead: true,
			}, nil
		}
	}
	return BranchInfo{
		RefName:      "HEAD",
		DisplayName:  "HEAD",
		DetachedHead: true,
	}, nil
}

// Delete delete branch
func (self *BranchCommands) Delete(branch string, force bool) error {
	command := "git branch -d"

	if force {
		command = "git branch -D"
	}

	return self.cmd.New(fmt.Sprintf("%s %s", command, self.cmd.Quote(branch))).Run()
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOptions struct {
	Force   bool
	EnvVars []string
}

func (self *BranchCommands) Checkout(branch string, options CheckoutOptions) error {
	forceArg := ""
	if options.Force {
		forceArg = " --force"
	}

	return self.cmd.New(fmt.Sprintf("git checkout%s %s", forceArg, self.cmd.Quote(branch))).
		// prevents git from prompting us for input which would freeze the program
		// TODO: see if this is actually needed here
		AddEnvVars("GIT_TERMINAL_PROMPT=0").
		AddEnvVars(options.EnvVars...).
		Run()
}

// GetGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (self *BranchCommands) GetGraph(branchName string) (string, error) {
	return self.GetGraphCmdObj(branchName).DontLog().RunWithOutput()
}

func (self *BranchCommands) GetGraphCmdObj(branchName string) oscommands.ICmdObj {
	branchLogCmdTemplate := self.UserConfig.Git.BranchLogCmd
	templateValues := map[string]string{
		"branchName": self.cmd.Quote(branchName),
	}
	return self.cmd.New(utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues)).DontLog()
}

func (self *BranchCommands) SetCurrentBranchUpstream(remoteName string, remoteBranchName string) error {
	return self.cmd.New(fmt.Sprintf("git branch --set-upstream-to=%s/%s", self.cmd.Quote(remoteName), self.cmd.Quote(remoteBranchName))).Run()
}

func (self *BranchCommands) SetUpstream(remoteName string, remoteBranchName string, branchName string) error {
	return self.cmd.New(fmt.Sprintf("git branch --set-upstream-to=%s/%s %s", self.cmd.Quote(remoteName), self.cmd.Quote(remoteBranchName), self.cmd.Quote(branchName))).Run()
}

func (self *BranchCommands) UnsetUpstream(branchName string) error {
	return self.cmd.New(fmt.Sprintf("git branch --unset-upstream %s", self.cmd.Quote(branchName))).Run()
}

func (self *BranchCommands) GetCurrentBranchUpstreamDifferenceCount() (string, string) {
	return self.GetCommitDifferences("HEAD", "HEAD@{u}")
}

func (self *BranchCommands) GetUpstreamDifferenceCount(branchName string) (string, string) {
	return self.GetCommitDifferences(branchName, branchName+"@{u}")
}

// GetCommitDifferences checks how many pushables/pullables there are for the
// current branch
func (self *BranchCommands) GetCommitDifferences(from, to string) (string, string) {
	command := "git rev-list %s..%s --count"
	pushableCount, err := self.cmd.New(fmt.Sprintf(command, to, from)).DontLog().RunWithOutput()
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := self.cmd.New(fmt.Sprintf(command, from, to)).DontLog().RunWithOutput()
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

func (self *BranchCommands) IsHeadDetached() bool {
	err := self.cmd.New("git symbolic-ref -q HEAD").DontLog().Run()
	return err != nil
}

func (self *BranchCommands) Rename(oldName string, newName string) error {
	return self.cmd.New(fmt.Sprintf("git branch --move %s %s", self.cmd.Quote(oldName), self.cmd.Quote(newName))).Run()
}

func (self *BranchCommands) GetRawBranches() (string, error) {
	return self.cmd.New(`git for-each-ref --sort=-committerdate --format="%(HEAD)%00%(refname:short)%00%(upstream:short)%00%(upstream:track)" refs/heads`).DontLog().RunWithOutput()
}

type MergeOpts struct {
	FastForwardOnly bool
}

func (self *BranchCommands) Merge(branchName string, opts MergeOpts) error {
	mergeArg := ""
	if self.UserConfig.Git.Merging.Args != "" {
		mergeArg = " " + self.UserConfig.Git.Merging.Args
	}

	command := fmt.Sprintf("git merge --no-edit%s %s", mergeArg, self.cmd.Quote(branchName))
	if opts.FastForwardOnly {
		command = fmt.Sprintf("%s --ff-only", command)
	}

	return self.cmd.New(command).Run()
}

func (self *BranchCommands) AllBranchesLogCmdObj() oscommands.ICmdObj {
	return self.cmd.New(self.UserConfig.Git.AllBranchesLogCmd).DontLog()
}
