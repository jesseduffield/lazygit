package icons

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

const (
	GIT_BRANCH_ICON         = "\ufb2b" // שׂ
	GIT_COMMIT_ICON         = "\ufc16" // ﰖ
	GIT_DEFAULT_REMOTE_ICON = "\uf7a1" // 
	GIT_DETACHED_HEAD_ICON  = "\ue729" // 
	GIT_DIFF_ADDED_ICON     = "\uf457" // 
	GIT_DIFF_IGNORED_ICON   = "\uf474" // 
	GIT_DIFF_MODIFIED_ICON  = "\uf459" // 
	GIT_DIFF_REMOVED_ICON   = "\uf458" // 
	GIT_DIFF_RENAMED_ICON   = "\uf45a" // 
	GIT_PULL_ICON           = "\uf0ed" // 
	GIT_PUSH_ICON           = "\uf0ee" // 
	GIT_MERGE_COMMIT_ICON   = "\ufb2c" // שּׁ
	GIT_TAG_ICON            = "\uf02b" // 
)

const (
	GIT_MERGE_COMMIT_SYMBOL  = "⏣"
	GIT_COMMIT_SYMBOL = "◯"
)

type remoteIcon struct {
	domain string
	icon   string
}

var remoteIcons = []remoteIcon{
	{domain: "github.com",    icon: "\ue708"}, // 
	{domain: "bitbucket.org", icon: "\ue703"}, // 
	{domain: "gitlab.com",    icon: "\uf296"}, // 
	{domain: "dev.azure.com", icon: "\ufd03"}, // ﴃ
}

func IconForMergeCommitCell() string {
	if IsIconEnabled() {
		return GIT_MERGE_COMMIT_ICON
	} else {
		return GIT_MERGE_COMMIT_SYMBOL
	}
}

func IconForCommitCell() string {
	if IsIconEnabled() {
		return GIT_COMMIT_ICON
	} else {
		return GIT_COMMIT_SYMBOL
	}
}

func IconForBranch(branch *models.Branch) string {
	if branch.DisplayName != "" {
		return GIT_DETACHED_HEAD_ICON
	}
	return GIT_BRANCH_ICON
}

func IconForPull() string {
	if IsIconEnabled() {
		return GIT_PULL_ICON + " "
	} else {
		return "↓"
	}
}

func IconForPush() string {
	if IsIconEnabled() {
		return GIT_PUSH_ICON + " "
	} else {
		return "↑"
	}
}

func IconForRemoteBranch(branch *models.RemoteBranch) string {
	return GIT_BRANCH_ICON
}

func IconForTag(tag *models.Tag) string {
	return GIT_TAG_ICON
}

func IconForCommit(commit *models.Commit) string {
	if len(commit.Parents) > 1 {
		return GIT_MERGE_COMMIT_ICON
	}
	return GIT_COMMIT_ICON
}

func IconForRemote(remote *models.Remote) string {
	for _, r := range remoteIcons {
		for _, url := range remote.Urls {
			if strings.Contains(url, r.domain) {
				return r.icon
			}
		}
	}
	return GIT_DEFAULT_REMOTE_ICON
}


func IconForChangeStatus(changeStatus string) string {
	switch changeStatus {
	case "?", "A":
		return GIT_DIFF_ADDED_ICON
	case "M", "C":
		return GIT_DIFF_MODIFIED_ICON
	case "R", "T":
		return GIT_DIFF_RENAMED_ICON
	case "D":
		return GIT_DIFF_REMOVED_ICON
	default:
		return " "
	}
}
