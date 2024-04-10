package icons

import (
	"strings"

	"github.com/lobes/lazytask/pkg/commands/models"
)

var (
	BRANCH_ICON                  = "\U000f062c" // 󰘬
	DETACHED_HEAD_ICON           = "\ue729"     // 
	TAG_ICON                     = "\uf02b"     // 
	COMMIT_ICON                  = "\U000f0718" // 󰜘
	MERGE_COMMIT_ICON            = "\U000f062d" // 󰘭
	DEFAULT_REMOTE_ICON          = "\uf02a2"    // 󰊢
	STASH_ICON                   = "\uf01c"     // 
	LINKED_WORKTREE_ICON         = "\U000f0339" // 󰌹
	MISSING_LINKED_WORKTREE_ICON = "\U000f033a" // 󰌺
)

var remoteIcons = map[string]string{
	"github.com":    "\ue709",     // 
	"bitbucket.org": "\ue703",     // 
	"gitlab.com":    "\uf296",     // 
	"dev.azure.com": "\U000f0805", // 󰠅
}

func patchGitIconsForNerdFontsV2() {
	BRANCH_ICON = "\ufb2b"                  // שׂ
	COMMIT_ICON = "\ufc16"                  // ﰖ
	MERGE_COMMIT_ICON = "\ufb2c"            // שּׁ
	DEFAULT_REMOTE_ICON = "\uf7a1"          // 
	LINKED_WORKTREE_ICON = "\uf838"         // 
	MISSING_LINKED_WORKTREE_ICON = "\uf839" // 

	remoteIcons["dev.azure.com"] = "\ufd03" // ﴃ
}

func IconForBranch(branch *models.Branch) string {
	if branch.DetachedHead {
		return DETACHED_HEAD_ICON
	}
	return BRANCH_ICON
}

func IconForRemoteBranch(branch *models.RemoteBranch) string {
	return BRANCH_ICON
}

func IconForTag(tag *models.Tag) string {
	return TAG_ICON
}

func IconForCommit(commit *models.Commit) string {
	if len(commit.Parents) > 1 {
		return MERGE_COMMIT_ICON
	}
	return COMMIT_ICON
}

func IconForRemote(remote *models.Remote) string {
	for domain, icon := range remoteIcons {
		for _, url := range remote.Urls {
			if strings.Contains(url, domain) {
				return icon
			}
		}
	}
	return DEFAULT_REMOTE_ICON
}

func IconForStash(stash *models.StashEntry) string {
	return STASH_ICON
}

func IconForWorktree(missing bool) string {
	if missing {
		return MISSING_LINKED_WORKTREE_ICON
	}
	return LINKED_WORKTREE_ICON
}
