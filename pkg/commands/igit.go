// hmm auto-generated for testing purposes. To re-generate, do: <ifacemaker --file="pkg/commands/*.go" --struct=Git --iface=IGit --pkg=commands -o pkg/commands/igit.go --doc false --comment="$(cat pkg/commands/auto-generation-message.txt)"> from the root directory of the repo and fix up any missing imports

package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/sirupsen/logrus"
)

//counterfeiter:generate . IGit
type IGit interface {
	Branches() IBranchesMgr
	Commits() ICommitsMgr
	Worktree() IWorktreeMgr
	Submodules() ISubmodulesMgr
	Status() IStatusMgr
	Stash() IStashMgr
	Tags() ITagsMgr
	Remotes() IRemotesMgr

	// config
	IGitConfigMgr

	FindRemoteForBranchInConfig(branchName string) (string, error)

	// diffing
	WorktreeFileDiff(file *models.File, plain bool, cached bool) string
	WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool) ICmdObj
	ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error)
	ShowFileDiffCmdObj(from string, to string, reverse bool, path string, plain bool, showRenames bool) ICmdObj
	DiffEndArgs(from string, to string, reverse bool, path string) string

	// commands
	ICommander

	SetCredentialHandlers(promptUserForCredential func(CredentialKind) string, handleCredentialError func(error))

	// common
	GetLog() *logrus.Entry
	WithSpan(span string) IGit
	GetOS() oscommands.IOS

	// flow
	FlowStart(branchType string, name string) ICmdObj
	FlowFinish(branchType string, name string) ICmdObj
	GetGitFlowRegexpConfig() (string, error)

	// loaders
	GetFilesInDiff(from string, to string, reverse bool) ([]*models.CommitFile, error)
	GetReflogCommits(lastReflogCommit *models.Commit, filterPath string) ([]*models.Commit, bool, error)

	// patch
	NewPatchManager() *patch.PatchManager
	DeletePatchesFromCommit(commits []*models.Commit, commitIndex int, p *patch.PatchManager) error
	MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int, p *patch.PatchManager) error
	MovePatchIntoIndex(commits []*models.Commit, commitIdx int, p *patch.PatchManager, stash bool) error
	PullPatchIntoNewCommit(commits []*models.Commit, commitIdx int, p *patch.PatchManager) error

	// rebasing
	DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error
	GenericAbortCmdObj() ICmdObj
	GenericContinueCmdObj() ICmdObj
	GenericMergeOrRebaseCmdObj(action string) ICmdObj
	AbortRebase() error
	ContinueRebase() error
	MergeOrRebase() string
	GetRewordCommitCmdObj(commits []*models.Commit, index int) (ICmdObj, error)
	MoveCommitDown(commits []*models.Commit, index int) error
	InteractiveRebase(commits []*models.Commit, index int, action string) error
	InteractiveRebaseCmdObj(baseSha string, todo string, overrideEditor bool) ICmdObj
	GenerateGenericRebaseTodo(commits []*models.Commit, actionIndex int, action string) (string, string, error)
	AmendTo(sha string) error
	EditRebaseTodo(index int, action string) error
	MoveTodoDown(index int) error
	SquashAllAboveFixupCommits(sha string) error
	BeginInteractiveRebaseForCommit(commits []*models.Commit, commitIndex int) error
	RebaseBranch(branchName string) error
	GenericMergeOrRebaseAction(commandType string, command string) error
	CherryPickCommits(commits []*models.Commit) error

	// sync
	Push(opts PushOpts) (bool, error)
	Fetch(opts FetchOptions) error
	FetchInBackground(opts FetchOptions) error
	FastForward(branchName string, remoteName string, remoteBranchName string) error
	FetchRemote(remoteName string) error
	PushRef(remoteName string, refName string) error
	DeleteRemoteRef(remoteName string, ref string) error
}
