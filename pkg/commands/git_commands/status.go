package git_commands

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/spf13/afero"
)

type StatusCommands struct {
	*GitCommon
}

func NewStatusCommands(
	gitCommon *GitCommon,
) *StatusCommands {
	return &StatusCommands{
		GitCommon: gitCommon,
	}
}

func (self *StatusCommands) WorkingTreeState() models.WorkingTreeState {
	result := models.WorkingTreeState{}
	result.Rebasing, _ = self.IsInRebase()
	result.Merging, _ = self.IsInMergeState()
	result.CherryPicking, _ = self.IsInCherryPick()
	result.Reverting, _ = self.IsInRevert()
	return result
}

func (self *StatusCommands) IsBareRepo() bool {
	return self.repoPaths.isBareRepo
}

func (self *StatusCommands) IsInRebase() (bool, error) {
	exists, err := self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge"))
	if err == nil && exists {
		return true, nil
	}
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-apply"))
}

// IsInMergeState states whether we are still mid-merge
func (self *StatusCommands) IsInMergeState() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "MERGE_HEAD"))
}

func (self *StatusCommands) IsInCherryPick() (bool, error) {
	exists, err := self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "CHERRY_PICK_HEAD"))
	if err != nil || !exists {
		return exists, err
	}
	// Sometimes, CHERRY_PICK_HEAD is present during rebases even if no
	// cherry-pick is in progress. I suppose this is because rebase used to be
	// implemented as a series of cherry-picks, so this could be remnants of
	// code that is shared between cherry-pick and rebase, or something. The way
	// to tell if this is the case is to check for the presence of the
	// stopped-sha file, which records the sha of the last pick that was
	// executed before the rebase stopped, and seeing if the sha in that file is
	// the same as the one in CHERRY_PICK_HEAD.
	cherryPickHead, err := os.ReadFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "CHERRY_PICK_HEAD"))
	if err != nil {
		return false, err
	}
	stoppedSha, err := os.ReadFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge", "stopped-sha"))
	if err != nil {
		// If we get an error we assume the file doesn't exist
		return true, nil
	}
	cherryPickHeadStr := strings.TrimSpace(string(cherryPickHead))
	stoppedShaStr := strings.TrimSpace(string(stoppedSha))
	// Need to use HasPrefix here because the cherry-pick HEAD is a full sha1,
	// but stopped-sha is an abbreviated sha1
	if strings.HasPrefix(cherryPickHeadStr, stoppedShaStr) {
		return false, nil
	}
	return true, nil
}

func (self *StatusCommands) IsInRevert() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "REVERT_HEAD"))
}

// RefsSnapshot returns a string fingerprint of the current state of local
// branches and HEAD. Comparing two snapshots byte-for-byte tells us whether
// any local ref or HEAD has moved since the last snapshot.
func (self *StatusCommands) RefsSnapshot() (string, error) {
	t := time.Now()
	defer func() { self.Log.Infof("RefsSnapshot took %s", time.Since(t)) }()

	refsArgs := NewGitCmd("for-each-ref").
		Arg("--format=%(objectname) %(refname)").
		Arg("refs/heads").
		ToArgv()
	refs, err := self.cmd.New(refsArgs).DontLog().RunWithOutput()
	if err != nil {
		return "", err
	}

	head, err := self.headSnapshot()
	if err != nil {
		return "", err
	}

	return refs + head, nil
}

// headSnapshot returns a fingerprint of HEAD that distinguishes "detached at
// commit X" from "on a branch that points at X". The commit hash alone can't
// tell those apart, which matters at the end of a rebase: HEAD reattaches to
// the branch without the hash changing, and we'd otherwise miss that refresh.
//
// We read .git/HEAD directly rather than shelling out: it's faster (no child
// process) and its content is exactly the symref-or-hash distinction we want
// ("ref: refs/heads/foo" when attached, the raw hash when detached). The
// reftable backend, however, doesn't keep a real .git/HEAD — it writes a fixed
// stub ("ref: refs/heads/.invalid") that never reflects the actual HEAD. When
// we see that stub (or the file is missing/unreadable) we fall back to
// porcelain commands, which are backend-agnostic.
func (self *StatusCommands) headSnapshot() (string, error) {
	headPath := filepath.Join(self.repoPaths.WorktreeGitDirPath(), "HEAD")
	if content, err := afero.ReadFile(self.Fs, headPath); err == nil {
		head := strings.TrimSpace(string(content))
		if head != "" && head != "ref: refs/heads/.invalid" {
			return head, nil
		}
	}

	// symbolic-ref gives the branch when HEAD is attached and fails when it's
	// detached, in which case rev-parse gives the commit HEAD points at.
	symbolicRefArgs := NewGitCmd("symbolic-ref").Arg("HEAD").ToArgv()
	if symref, err := self.cmd.New(symbolicRefArgs).DontLog().RunWithOutput(); err == nil {
		return strings.TrimSpace(symref), nil
	}

	revParseArgs := NewGitCmd("rev-parse").Arg("HEAD").ToArgv()
	head, err := self.cmd.New(revParseArgs).DontLog().RunWithOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(head), nil
}

// Full ref (e.g. "refs/heads/mybranch") of the branch that is currently
// being rebased, or empty string when we're not in a rebase
func (self *StatusCommands) BranchBeingRebased() string {
	for _, dir := range []string{"rebase-merge", "rebase-apply"} {
		if bytesContent, err := os.ReadFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), dir, "head-name")); err == nil {
			return strings.TrimSpace(string(bytesContent))
		}
	}
	return ""
}
