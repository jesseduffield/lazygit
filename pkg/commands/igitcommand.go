// auto-generated file
// created via: ifacemaker --file="pkg/commands/*.go" --struct=GitCommand --iface=IGitCommand --pkg=commands

package commands

import (
	"os/exec"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
)

type IGitCommand interface {
	// NewBranch create new branch
	NewBranch(name string, base string) error
	// CurrentBranchName get the current branch name and displayname.
	// the first returned string is the name and the second is the displayname
	// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
	CurrentBranchName() (string, string, error)
	// DeleteBranch delete branch
	DeleteBranch(branch string, force bool) error
	Checkout(branch string, options CheckoutOptions) error
	// GetBranchGraph gets the color-formatted graph of the log for the given branch
	// Currently it limits the result to 100 commits, but when we get async stuff
	// working we can do lazy loading
	GetBranchGraph(branchName string) (string, error)
	GetUpstreamForBranch(branchName string) (string, error)
	GetBranchGraphCmdStr(branchName string) string
	SetUpstreamBranch(upstream string) error
	SetBranchUpstream(remoteName string, remoteBranchName string, branchName string) error
	GetCurrentBranchUpstreamDifferenceCount() (string, string)
	GetBranchUpstreamDifferenceCount(branchName string) (string, string)
	// GetCommitDifferences checks how many pushables/pullables there are for the
	// current branch
	GetCommitDifferences(from, to string) (string, string)
	// Merge merge
	Merge(branchName string, opts MergeOpts) error
	// AbortMerge abort merge
	AbortMerge() error
	IsHeadDetached() bool
	// ResetHardHead runs `git reset --hard`
	ResetHard(ref string) error
	// ResetSoft runs `git reset --soft HEAD`
	ResetSoft(ref string) error
	ResetMixed(ref string) error
	RenameBranch(oldName string, newName string) error
	// RenameCommit renames the topmost commit with the given name
	RenameCommit(name string) error
	// ResetToCommit reset to commit
	ResetToCommit(sha string, strength string, options oscommands.RunCommandOptions) error
	CommitCmdStr(message string, flags string) string
	// Get the subject of the HEAD commit
	GetHeadCommitMessage() (string, error)
	GetCommitMessage(commitSha string) (string, error)
	GetCommitMessageFirstLine(sha string) (string, error)
	// AmendHead amends HEAD with whatever is staged in your working tree
	AmendHead() error
	AmendHeadCmdStr() string
	ShowCmdStr(sha string, filterPath string) string
	// Revert reverts the selected commit by sha
	Revert(sha string) error
	RevertMerge(sha string, parentNumber int) error
	// CherryPickCommits begins an interactive rebase with the given shas being cherry picked onto HEAD
	CherryPickCommits(commits []*models.Commit) error
	// CreateFixupCommit creates a commit that fixes up a previous commit
	CreateFixupCommit(sha string) error
	ConfiguredPager() string
	GetPager(width int) string
	GetConfigValue(key string) string
	// UsingGpg tells us whether the user has gpg enabled so that we can know
	// whether we need to run a subprocess to allow them to enter their password
	UsingGpg() bool
	// CatFile obtains the content of a file
	CatFile(fileName string) (string, error)
	OpenMergeToolCmd() string
	OpenMergeTool() error
	// StageFile stages a file
	StageFile(fileName string) error
	// StageAll stages all files
	StageAll() error
	// UnstageAll unstages all files
	UnstageAll() error
	// UnStageFile unstages a file
	// we accept an array of filenames for the cases where a file has been renamed i.e.
	// we accept the current name and the previous name
	UnStageFile(fileNames []string, reset bool) error
	BeforeAndAfterFileForRename(file *models.File) (*models.File, *models.File, error)
	// DiscardAllFileChanges directly
	DiscardAllFileChanges(file *models.File) error
	DiscardAllDirChanges(node *filetree.FileNode) error
	DiscardUnstagedDirChanges(node *filetree.FileNode) error
	RemoveUntrackedDirFiles(node *filetree.FileNode) error
	// DiscardUnstagedFileChanges directly
	DiscardUnstagedFileChanges(file *models.File) error
	// Ignore adds a file to the gitignore for the repo
	Ignore(filename string) error
	// WorktreeFileDiff returns the diff of a file
	WorktreeFileDiff(file *models.File, plain bool, cached bool) string
	WorktreeFileDiffCmdStr(node models.IFile, plain bool, cached bool) string
	ApplyPatch(patch string, flags ...string) error
	// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
	// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
	ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error)
	ShowFileDiffCmdStr(from string, to string, reverse bool, fileName string, plain bool) string
	// CheckoutFile checks out the file for the given commit
	CheckoutFile(commitSha, fileName string) error
	// DiscardOldFileChanges discards changes to a file from an old commit
	DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error
	// DiscardAnyUnstagedFileChanges discards any unstages file changes via `git checkout -- .`
	DiscardAnyUnstagedFileChanges() error
	// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
	RemoveTrackedFiles(name string) error
	// RemoveUntrackedFiles runs `git clean -fd`
	RemoveUntrackedFiles() error
	// ResetAndClean removes all unstaged changes and removes all untracked files
	ResetAndClean() error
	EditFileCmdStr(filename string) (string, error)
	WithSpan(span string) IGitCommand
	RunCommand(formatString string, formatArgs ...interface{}) error
	RunCommandWithOutput(formatString string, formatArgs ...interface{}) (string, error)
	// GetFilesInDiff get the specified commit files
	GetFilesInDiff(from string, to string, reverse bool) ([]*models.CommitFile, error)
	GetStatusFiles(opts GetStatusFileOptions) []*models.File
	GitStatus(opts GitStatusOptions) (string, error)
	// GetReflogCommits only returns the new reflog commits since the given lastReflogCommit
	// if none is passed (i.e. it's value is nil) then we get all the reflog commits
	GetReflogCommits(lastReflogCommit *models.Commit, filterPath string) ([]*models.Commit, bool, error)
	GetRemotes() ([]*models.Remote, error)
	// GetStashEntries stash entries
	GetStashEntries(filterPath string) []*models.StashEntry
	GetTags() ([]*models.Tag, error)
	// DeletePatchesFromCommit applies a patch in reverse for a commit
	DeletePatchesFromCommit(commits []*models.Commit, commitIndex int, p *patch.PatchManager) error
	MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int, p *patch.PatchManager) error
	MovePatchIntoIndex(commits []*models.Commit, commitIdx int, p *patch.PatchManager, stash bool) error
	PullPatchIntoNewCommit(commits []*models.Commit, commitIdx int, p *patch.PatchManager) error
	RewordCommit(commits []*models.Commit, index int) (*exec.Cmd, error)
	MoveCommitDown(commits []*models.Commit, index int) error
	InteractiveRebase(commits []*models.Commit, index int, action string) error
	// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
	// we tell git to run lazygit to edit the todo list, and we pass the client
	// lazygit a todo string to write to the todo file
	PrepareInteractiveRebaseCommand(baseSha string, todo string, overrideEditor bool) (*exec.Cmd, error)
	GenerateGenericRebaseTodo(commits []*models.Commit, actionIndex int, action string) (string, string, error)
	// AmendTo amends the given commit with whatever files are staged
	AmendTo(sha string) error
	// EditRebaseTodo sets the action at a given index in the git-rebase-todo file
	EditRebaseTodo(index int, action string) error
	// MoveTodoDown moves a rebase todo item down by one position
	MoveTodoDown(index int) error
	// SquashAllAboveFixupCommits squashes all fixup! commits above the given one
	SquashAllAboveFixupCommits(sha string) error
	// BeginInteractiveRebaseForCommit starts an interactive rebase to edit the current
	// commit and pick all others. After this you'll want to call `c.GenericMergeOrRebaseAction("rebase", "continue")`
	BeginInteractiveRebaseForCommit(commits []*models.Commit, commitIndex int) error
	// RebaseBranch interactive rebases onto a branch
	RebaseBranch(branchName string) error
	// GenericMerge takes a commandType of "merge" or "rebase" and a command of "abort", "skip" or "continue"
	// By default we skip the editor in the case where a commit will be made
	GenericMergeOrRebaseAction(commandType string, command string) error
	AddRemote(name string, url string) error
	RemoveRemote(name string) error
	RenameRemote(oldRemoteName string, newRemoteName string) error
	UpdateRemoteUrl(remoteName string, updatedUrl string) error
	DeleteRemoteBranch(remoteName string, branchName string, promptUserForCredential func(string) string) error
	// CheckRemoteBranchExists Returns remote branch
	CheckRemoteBranchExists(branch *models.Branch) bool
	// GetRemoteURL returns current repo remote url
	GetRemoteURL() string
	// StashDo modify stash
	StashDo(index int, method string) error
	// StashSave save stash
	// TODO: before calling this, check if there is anything to save
	StashSave(message string) error
	// GetStashEntryDiff stash diff
	ShowStashEntryCmdStr(index int) string
	// StashSaveStagedChanges stashes only the currently staged changes. This takes a few steps
	// shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
	StashSaveStagedChanges(message string) error
	// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
	// and "interactive" for interactive rebase
	RebaseMode() (string, error)
	WorkingTreeState() string
	// IsInMergeState states whether we are still mid-merge
	IsInMergeState() (bool, error)
	IsBareRepo() bool
	GetSubmoduleConfigs() ([]*models.SubmoduleConfig, error)
	SubmoduleStash(submodule *models.SubmoduleConfig) error
	SubmoduleReset(submodule *models.SubmoduleConfig) error
	SubmoduleUpdateAll() error
	SubmoduleDelete(submodule *models.SubmoduleConfig) error
	SubmoduleAdd(name string, path string, url string) error
	SubmoduleUpdateUrl(name string, path string, newUrl string) error
	SubmoduleInit(path string) error
	SubmoduleUpdate(path string) error
	SubmoduleBulkInitCmdStr() string
	SubmoduleBulkUpdateCmdStr() string
	SubmoduleForceBulkUpdateCmdStr() string
	SubmoduleBulkDeinitCmdStr() string
	ResetSubmodules(submodules []*models.SubmoduleConfig) error
	// Push pushes to a branch
	Push(branchName string, force bool, upstream string, args string, promptUserForCredential func(string) string) error
	// Fetch fetch git repo
	Fetch(opts FetchOptions) error
	FastForward(branchName string, remoteName string, remoteBranchName string, promptUserForCredential func(string) string) error
	FetchRemote(remoteName string, promptUserForCredential func(string) string) error
	CreateLightweightTag(tagName string, commitSha string) error
	DeleteTag(tagName string) error
	PushTag(remoteName string, tagName string, promptUserForCredential func(string) string) error
	GetPushToCurrent() bool
	FindRemoteForBranchInConfig(string) (string, error)
	GetOSCommand() *oscommands.OSCommand
}
