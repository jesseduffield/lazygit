package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/sirupsen/logrus"
)

//counterfeiter:generate . ICommitsMgr
type ICommitsMgr interface {
	RewordHead(name string) error
	CommitCmdObj(message string, flags string) ICmdObj
	GetHeadMessage() (string, error)
	GetMessage(commitSha string) (string, error)
	GetMessageFirstLine(sha string) (string, error)
	AmendHead() error
	AmendHeadCmdObj() ICmdObj
	ShowCmdObj(sha string, filterPath string) ICmdObj
	Revert(sha string) error
	RevertMerge(sha string, parentNumber int) error
	CreateFixupCommit(sha string) error
	Load(LoadCommitsOptions) ([]*models.Commit, error)
	MergeRebasingCommits(commits []*models.Commit) ([]*models.Commit, error)
}

type CommitsMgr struct {
	ICommander

	commitsLoader *CommitsLoader
	config        IGitConfigMgr
}

func NewCommitsMgr(
	commander ICommander,
	config IGitConfigMgr,
	branchesMgr IBranchesMgr,
	statusMgr IStatusMgr,
	log *logrus.Entry,
	oS *oscommands.OS,
	tr *i18n.TranslationSet,
	dotGitDir string,
) *CommitsMgr {
	commitsLoader := NewCommitsLoader(
		log, branchesMgr, statusMgr, oS, tr, dotGitDir, commander,
	)

	return &CommitsMgr{
		commitsLoader: commitsLoader,
		ICommander:    commander,
		config:        config,
	}
}

func (c *CommitsMgr) Load(opts LoadCommitsOptions) ([]*models.Commit, error) {
	return c.commitsLoader.Load(opts)
}

func (c *CommitsMgr) MergeRebasingCommits(commits []*models.Commit) ([]*models.Commit, error) {
	return c.commitsLoader.MergeRebasingCommits(commits)
}

// RenameCommit renames the topmost commit with the given name
func (c *CommitsMgr) RewordHead(name string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("commit --allow-empty --amend --only -m %s", c.Quote(name)))
}

func (c *CommitsMgr) CommitCmdObj(message string, flags string) ICmdObj {
	splitMessage := strings.Split(message, "\n")
	lineArgs := ""
	for _, line := range splitMessage {
		lineArgs += fmt.Sprintf(" -m %s", c.Quote(line))
	}

	flagsStr := ""
	if flags != "" {
		flagsStr = fmt.Sprintf(" %s", flags)
	}

	cmdStr := fmt.Sprintf("commit%s%s", flagsStr, lineArgs)

	return c.BuildGitCmdObjFromStr(cmdStr)
}

// Get the subject of the HEAD commit
func (c *CommitsMgr) GetHeadMessage() (string, error) {
	cmdObj := c.BuildGitCmdObjFromStr("log -1 --pretty=%s")
	message, err := c.RunWithOutput(cmdObj)
	return strings.TrimSpace(message), err
}

func (c *CommitsMgr) GetMessage(commitSha string) (string, error) {
	messageWithHeader, err := c.RunWithOutput(
		c.BuildGitCmdObjFromStr("rev-list --format=%B --max-count=1 " + commitSha),
	)
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "\n")
	return strings.TrimSpace(message), err
}

func (c *CommitsMgr) GetMessageFirstLine(sha string) (string, error) {
	return c.RunWithOutput(
		c.BuildGitCmdObjFromStr(fmt.Sprintf("show --no-patch --pretty=format:%%s %s", sha)),
	)
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (c *CommitsMgr) AmendHead() error {
	return c.Run(c.AmendHeadCmdObj())
}

func (c *CommitsMgr) AmendHeadCmdObj() ICmdObj {
	return c.BuildGitCmdObjFromStr("commit --amend --no-edit --allow-empty")
}

func (c *CommitsMgr) ShowCmdObj(sha string, filterPath string) ICmdObj {
	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" -- %s", c.Quote(filterPath))
	}
	return c.BuildGitCmdObjFromStr(
		fmt.Sprintf("show --submodule --color=%s --no-renames --stat -p %s%s", c.config.ColorArg(), sha, filterPathArg),
	)
}

// Revert reverts the selected commit by sha
func (c *CommitsMgr) Revert(sha string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("revert %s", sha))
}

func (c *CommitsMgr) RevertMerge(sha string, parentNumber int) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("revert %s -m %d", sha, parentNumber))
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (c *CommitsMgr) CreateFixupCommit(sha string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("commit --fixup=%s", sha))
}
