package icons

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

var (
	BRANCH_ICON                  = IconProperties{"\U000f062c", 239} // 󰘬
	DETACHED_HEAD_ICON           = IconProperties{"\ue729", 239}     // 
	TAG_ICON                     = IconProperties{"\uf02b", 239}     // 
	COMMIT_ICON                  = IconProperties{"\U000f0718", 239} // 󰜘
	MERGE_COMMIT_ICON            = IconProperties{"\U000f062d", 239} // 󰘭
	DEFAULT_REMOTE_ICON          = IconProperties{"\uf02a2", 239}    // 󰊢
	STASH_ICON                   = IconProperties{"\uf01c", 239}     // 
	LINKED_WORKTREE_ICON         = IconProperties{"\U000f0339", 239} // 󰌹
	MISSING_LINKED_WORKTREE_ICON = IconProperties{"\U000f033a", 239} // 󰌺
)

var remoteIcons = map[string]IconProperties{
	"github.com":    {"\ue709", 239},     // 
	"bitbucket.org": {"\ue703", 239},     // 
	"gitlab.com":    {"\uf296", 239},     // 
	"dev.azure.com": {"\U000f0805", 239}, // 󰠅
}

func patchGitIconsForNerdFontsV2() {
	BRANCH_ICON = IconProperties{"\ufb2b", 239}                  // שׂ
	COMMIT_ICON = IconProperties{"\ufc16", 239}                  // ﰖ
	MERGE_COMMIT_ICON = IconProperties{"\ufb2c", 239}            // שּׁ
	DEFAULT_REMOTE_ICON = IconProperties{"\uf7a1", 239}          // 
	LINKED_WORKTREE_ICON = IconProperties{"\uf838", 239}         // 
	MISSING_LINKED_WORKTREE_ICON = IconProperties{"\uf839", 239} // 

	remoteIcons["dev.azure.com"] = IconProperties{"\ufd03", 239} // ﴃ
}

func IconForBranch(branch *models.Branch) IconProperties {
	if branch.DetachedHead {
		return DETACHED_HEAD_ICON
	}
	return BRANCH_ICON
}

func IconForRemoteBranch(branch *models.RemoteBranch) IconProperties {
	return BRANCH_ICON
}

func IconForTag(tag *models.Tag) IconProperties {
	return TAG_ICON
}

func IconForCommit(commit *models.Commit) IconProperties {
	if len(commit.Parents) > 1 {
		return MERGE_COMMIT_ICON
	}
	return COMMIT_ICON
}

func IconForRemote(remote *models.Remote) IconProperties {
	for domain, icon := range remoteIcons {
		for _, url := range remote.Urls {
			if strings.Contains(url, domain) {
				return icon
			}
		}
	}
	return DEFAULT_REMOTE_ICON
}

func IconForStash(stash *models.StashEntry) IconProperties {
	return STASH_ICON
}

func IconForWorktree(missing bool) IconProperties {
	if missing {
		return MISSING_LINKED_WORKTREE_ICON
	}
	return LINKED_WORKTREE_ICON
}
