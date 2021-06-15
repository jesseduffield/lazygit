// hmm auto-generated for testing purposes. To re-generate, do: <ifacemaker --file="pkg/commands/*.go" --struct=Git --iface=IGit --pkg=commands -o pkg/commands/igit.go --doc false --comment="$(cat pkg/commands/auto-generation-message.txt)"> from the root directory of the repo and fix up any missing imports

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

package commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/loaders"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/sirupsen/logrus"
)

// IGit ...
type IGit interface {
	// branches
	Branches() IBranchesMgr

	CurrentBranchName() (string, string, error)
	DeleteBranch(branch string, force bool) error
	Checkout(branch string, options CheckoutOptions) error
	GetBranchGraph(branchName string) (string, error)
	GetUpstreamForBranch(branchName string) (string, error)
	GetBranchGraphCmdObj(branchName string) ICmdObj
	SetUpstreamBranch(upstream string) error
	SetBranchUpstream(remoteName string, remoteBranchName string, branchName string) error
	GetCurrentBranchUpstreamDifferenceCount() (string, string)
	GetBranchUpstreamDifferenceCount(branchName string) (string, string)
	RenameBranch(oldName string, newName string) error
	FindRemoteForBranchInConfig(branchName string) (string, error)
	AllBranchesCmdObj() ICmdObj
	Merge(branchName string, opts MergeOpts) error
	AbortMerge() error

	// resetting
	ResetHard(ref string) error
	ResetSoft(ref string) error
	ResetMixed(ref string) error
	ResetToRef(ref string, strength string, options ResetToCommitOptions) error

	// commits
	Commits() ICommitsMgr

	// config
	GetPager(width int) string
	GetConfigValue(key string) string
	UsingGpg() bool
	GetPushToCurrent() bool

	// diffing
	GetCommitDifferences(from, to string) (string, string)
	WorktreeFileDiff(file *models.File, plain bool, cached bool) string
	WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool) ICmdObj
	ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error)
	ShowFileDiffCmdObj(from string, to string, reverse bool, path string, plain bool, showRenames bool) ICmdObj
	DiffEndArgs(from string, to string, reverse bool, path string) string

	CatFile(fileName string) (string, error)
	EditFileCmdObj(filename string) (ICmdObj, error)

	// worktree
	OpenMergeToolCmdObj() ICmdObj
	StageFile(fileName string) error
	StageAll() error
	UnstageAll() error
	UnStageFile(fileNames []string, reset bool) error
	BeforeAndAfterFileForRename(file *models.File) (*models.File, *models.File, error)
	DiscardAllFileChanges(file *models.File) error
	DiscardAllDirChanges(node *filetree.FileNode) error
	DiscardUnstagedDirChanges(node *filetree.FileNode) error
	RemoveUntrackedDirFiles(node *filetree.FileNode) error
	DiscardUnstagedFileChanges(file *models.File) error
	Ignore(filename string) error
	CheckoutFile(commitSha, fileName string) error
	DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error
	DiscardAnyUnstagedFileChanges() error
	RemoveTrackedFiles(name string) error
	RemoveUntrackedFiles() error
	ResetAndClean() error
	GetStatusFiles(opts loaders.LoadStatusFilesOpts) []*models.File

	// commands
	ICommander

	SetCredentialHandlers(promptUserForCredential func(CredentialKind) string, handleCredentialError func(error))

	// common
	GetLog() *logrus.Entry
	WithSpan(span string) IGit
	GetOS() *oscommands.OS

	// flow
	FlowStart(branchType string, name string) ICmdObj
	FlowFinish(branchType string, name string) ICmdObj
	GetGitFlowRegexpConfig() (string, error)

	// loaders
	GetFilesInDiff(from string, to string, reverse bool) ([]*models.CommitFile, error)
	GetReflogCommits(lastReflogCommit *models.Commit, filterPath string) ([]*models.Commit, bool, error)
	GetRemotes() ([]*models.Remote, error)
	GetStashEntries(filterPath string) []*models.StashEntry
	GetTags() ([]*models.Tag, error)

	// patch
	ApplyPatch(patch string, flags ...string) error
	NewPatchManager() *patch.PatchManager
	DeletePatchesFromCommit(commits []*models.Commit, commitIndex int, p *patch.PatchManager) error
	MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int, p *patch.PatchManager) error
	MovePatchIntoIndex(commits []*models.Commit, commitIdx int, p *patch.PatchManager, stash bool) error
	PullPatchIntoNewCommit(commits []*models.Commit, commitIdx int, p *patch.PatchManager) error

	// rebasing
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

	// remotes
	AddRemote(name string, url string) error
	RemoveRemote(name string) error
	RenameRemote(oldRemoteName string, newRemoteName string) error
	UpdateRemoteUrl(remoteName string, updatedUrl string) error
	DeleteRemoteBranch(remoteName string, branchName string) error
	CheckRemoteBranchExists(branch *models.Branch) bool
	GetRemoteURL() string

	// stash
	StashDo(index int, method string) error
	StashSave(message string) error
	ShowStashEntryCmdObj(index int) ICmdObj
	StashSaveStagedChanges(message string) error

	// state/info
	RebaseMode() (WorkingTreeState, error)
	WorkingTreeState() WorkingTreeState
	IsInMergeState() (bool, error)
	IsBareRepo() bool
	IsHeadDetached() bool

	// submodules
	GetSubmoduleConfigs() ([]*models.SubmoduleConfig, error)
	SubmoduleStash(submodule *models.SubmoduleConfig) error
	SubmoduleReset(submodule *models.SubmoduleConfig) error
	SubmoduleDelete(submodule *models.SubmoduleConfig) error
	SubmoduleAdd(name string, path string, url string) error
	SubmoduleUpdateUrl(name string, path string, newUrl string) error
	SubmoduleInit(path string) error
	SubmoduleUpdate(path string) error
	SubmoduleBulkInitCmdObj() ICmdObj
	SubmoduleBulkUpdateCmdObj() ICmdObj
	SubmoduleForceBulkUpdateCmdObj() ICmdObj
	SubmoduleBulkDeinitCmdObj() ICmdObj
	ResetSubmodules(submodules []*models.SubmoduleConfig) error

	// sync
	Push(opts PushOpts) (bool, error)
	Fetch(opts FetchOptions) error
	FetchInBackground(opts FetchOptions) error
	FastForward(branchName string, remoteName string, remoteBranchName string) error
	FetchRemote(remoteName string) error

	// tags
	DeleteTag(tagName string) error
	PushTag(remoteName string, tagName string) error
	CreateLightweightTag(tagName string, commitSha string) error
}
