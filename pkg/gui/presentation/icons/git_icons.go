package icons

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

const (
	BRANCH_ICON         = "\U000f062c" // 󰘬
	DETACHED_HEAD_ICON  = "\ue729" // 
	TAG_ICON            = "\uf02b" // 
	COMMIT_ICON         = "\U000f0718" // 󰜘
	MERGE_COMMIT_ICON   = "\U000f062d" // 󰘭
	DEFAULT_REMOTE_ICON = "\U000f02a2" // 󰊢
	STASH_ICON          = "\uf01c" // 
)

type remoteIcon struct {
	domain string
	icon   string
}

var remoteIcons = []remoteIcon{
	{domain: "github.com", icon: "\ue709"},    // 
	{domain: "bitbucket.org", icon: "\ue703"}, // 
	{domain: "gitlab.com", icon: "\uf296"},    // 
	{domain: "dev.azure.com", icon: "\U000f0805"}, // 󰠅
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
	for _, r := range remoteIcons {
		for _, url := range remote.Urls {
			if strings.Contains(url, r.domain) {
				return r.icon
			}
		}
	}
	return DEFAULT_REMOTE_ICON
}

func IconForStash(stash *models.StashEntry) string {
	return STASH_ICON
}
