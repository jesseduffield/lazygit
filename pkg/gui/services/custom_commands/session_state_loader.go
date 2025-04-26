package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/samber/lo"
)

// loads the session state at the time that a custom command is invoked, for use
// in the custom command's template strings
type SessionStateLoader struct {
	c          *helpers.HelperCommon
	refsHelper *helpers.RefsHelper
}

func NewSessionStateLoader(c *helpers.HelperCommon, refsHelper *helpers.RefsHelper) *SessionStateLoader {
	return &SessionStateLoader{
		c:          c,
		refsHelper: refsHelper,
	}
}

func commitShimFromModelCommit(commit *models.Commit) *Commit {
	if commit == nil {
		return nil
	}

	return &Commit{
		Hash:          commit.Hash(),
		Sha:           commit.Hash(),
		Name:          commit.Name,
		Status:        commit.Status,
		Action:        commit.Action,
		Tags:          commit.Tags,
		ExtraInfo:     commit.ExtraInfo,
		AuthorName:    commit.AuthorName,
		AuthorEmail:   commit.AuthorEmail,
		UnixTimestamp: commit.UnixTimestamp,
		Divergence:    commit.Divergence,
		Parents:       commit.Parents(),
	}
}

func fileShimFromModelFile(file *models.File) *File {
	if file == nil {
		return nil
	}

	return &File{
		Name:                    file.Path,
		PreviousName:            file.PreviousPath,
		HasStagedChanges:        file.HasStagedChanges,
		HasUnstagedChanges:      file.HasUnstagedChanges,
		Tracked:                 file.Tracked,
		Added:                   file.Added,
		Deleted:                 file.Deleted,
		HasMergeConflicts:       file.HasMergeConflicts,
		HasInlineMergeConflicts: file.HasInlineMergeConflicts,
		DisplayString:           file.DisplayString,
		ShortStatus:             file.ShortStatus,
		IsWorktree:              file.IsWorktree,
	}
}

func branchShimFromModelBranch(branch *models.Branch) *Branch {
	if branch == nil {
		return nil
	}

	return &Branch{
		Name:           branch.Name,
		DisplayName:    branch.DisplayName,
		Recency:        branch.Recency,
		Pushables:      branch.AheadForPull,
		Pullables:      branch.BehindForPull,
		AheadForPull:   branch.AheadForPull,
		BehindForPull:  branch.BehindForPull,
		AheadForPush:   branch.AheadForPush,
		BehindForPush:  branch.BehindForPush,
		UpstreamGone:   branch.UpstreamGone,
		Head:           branch.Head,
		DetachedHead:   branch.DetachedHead,
		UpstreamRemote: branch.UpstreamRemote,
		UpstreamBranch: branch.UpstreamBranch,
		Subject:        branch.Subject,
		CommitHash:     branch.CommitHash,
	}
}

func remoteBranchShimFromModelRemoteBranch(remoteBranch *models.RemoteBranch) *RemoteBranch {
	if remoteBranch == nil {
		return nil
	}

	return &RemoteBranch{
		Name:       remoteBranch.Name,
		RemoteName: remoteBranch.RemoteName,
	}
}

func remoteShimFromModelRemote(remote *models.Remote) *Remote {
	if remote == nil {
		return nil
	}

	return &Remote{
		Name: remote.Name,
		Urls: remote.Urls,
		Branches: lo.Map(remote.Branches, func(branch *models.RemoteBranch, _ int) *RemoteBranch {
			return remoteBranchShimFromModelRemoteBranch(branch)
		}),
	}
}

func tagShimFromModelRemote(tag *models.Tag) *Tag {
	if tag == nil {
		return nil
	}

	return &Tag{
		Name:    tag.Name,
		Message: tag.Message,
	}
}

func stashEntryShimFromModelRemote(stashEntry *models.StashEntry) *StashEntry {
	if stashEntry == nil {
		return nil
	}

	return &StashEntry{
		Index:   stashEntry.Index,
		Recency: stashEntry.Recency,
		Name:    stashEntry.Name,
	}
}

func commitFileShimFromModelRemote(commitFile *models.CommitFile) *CommitFile {
	if commitFile == nil {
		return nil
	}

	return &CommitFile{
		Name:         commitFile.Path,
		ChangeStatus: commitFile.ChangeStatus,
	}
}

func worktreeShimFromModelRemote(worktree *models.Worktree) *Worktree {
	if worktree == nil {
		return nil
	}

	return &Worktree{
		IsMain:        worktree.IsMain,
		IsCurrent:     worktree.IsCurrent,
		Path:          worktree.Path,
		IsPathMissing: worktree.IsPathMissing,
		GitDir:        worktree.GitDir,
		Branch:        worktree.Branch,
		Name:          worktree.Name,
	}
}

type CommitRange struct {
	From string
	To   string
}

func makeCommitRange(commits []*models.Commit, _ int, _ int) *CommitRange {
	if len(commits) == 0 {
		return nil
	}

	return &CommitRange{
		From: commits[len(commits)-1].Hash(),
		To:   commits[0].Hash(),
	}
}

// SessionState captures the current state of the application for use in custom commands
type SessionState struct {
	SelectedLocalCommit    *Commit // deprecated, use SelectedCommit
	SelectedReflogCommit   *Commit // deprecated, use SelectedCommit
	SelectedSubCommit      *Commit // deprecated, use SelectedCommit
	SelectedCommit         *Commit
	SelectedCommitRange    *CommitRange
	SelectedFile           *File
	SelectedPath           string
	SelectedLocalBranch    *Branch
	SelectedRemoteBranch   *RemoteBranch
	SelectedRemote         *Remote
	SelectedTag            *Tag
	SelectedStashEntry     *StashEntry
	SelectedCommitFile     *CommitFile
	SelectedCommitFilePath string
	SelectedWorktree       *Worktree
	CheckedOutBranch       *Branch
}

func (self *SessionStateLoader) call() *SessionState {
	selectedLocalCommit := commitShimFromModelCommit(self.c.Contexts().LocalCommits.GetSelected())
	selectedLocalCommitRange := makeCommitRange(self.c.Contexts().LocalCommits.GetSelectedItems())
	selectedReflogCommit := commitShimFromModelCommit(self.c.Contexts().ReflogCommits.GetSelected())
	selectedReflogCommitRange := makeCommitRange(self.c.Contexts().ReflogCommits.GetSelectedItems())
	selectedSubCommit := commitShimFromModelCommit(self.c.Contexts().SubCommits.GetSelected())
	selectedSubCommitRange := makeCommitRange(self.c.Contexts().SubCommits.GetSelectedItems())

	selectedCommit := selectedLocalCommit
	selectedCommitRange := selectedLocalCommitRange
	if self.c.Context().IsCurrentOrParent(self.c.Contexts().ReflogCommits) {
		selectedCommit = selectedReflogCommit
		selectedCommitRange = selectedReflogCommitRange
	} else if self.c.Context().IsCurrentOrParent(self.c.Contexts().SubCommits) {
		selectedCommit = selectedSubCommit
		selectedCommitRange = selectedSubCommitRange
	}

	selectedPath := self.c.Contexts().Files.GetSelectedPath()
	selectedCommitFilePath := self.c.Contexts().CommitFiles.GetSelectedPath()

	if self.c.Context().IsCurrent(self.c.Contexts().CommitFiles) {
		selectedPath = selectedCommitFilePath
	}

	return &SessionState{
		SelectedFile:           fileShimFromModelFile(self.c.Contexts().Files.GetSelectedFile()),
		SelectedPath:           selectedPath,
		SelectedLocalCommit:    selectedLocalCommit,
		SelectedReflogCommit:   selectedReflogCommit,
		SelectedSubCommit:      selectedSubCommit,
		SelectedCommit:         selectedCommit,
		SelectedCommitRange:    selectedCommitRange,
		SelectedLocalBranch:    branchShimFromModelBranch(self.c.Contexts().Branches.GetSelected()),
		SelectedRemoteBranch:   remoteBranchShimFromModelRemoteBranch(self.c.Contexts().RemoteBranches.GetSelected()),
		SelectedRemote:         remoteShimFromModelRemote(self.c.Contexts().Remotes.GetSelected()),
		SelectedTag:            tagShimFromModelRemote(self.c.Contexts().Tags.GetSelected()),
		SelectedStashEntry:     stashEntryShimFromModelRemote(self.c.Contexts().Stash.GetSelected()),
		SelectedCommitFile:     commitFileShimFromModelRemote(self.c.Contexts().CommitFiles.GetSelectedFile()),
		SelectedCommitFilePath: selectedCommitFilePath,
		SelectedWorktree:       worktreeShimFromModelRemote(self.c.Contexts().Worktrees.GetSelected()),
		CheckedOutBranch:       branchShimFromModelBranch(self.refsHelper.GetCheckedOutRef()),
	}
}
