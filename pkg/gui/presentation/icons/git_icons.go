package icons

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

const (
	BRANCH_ICON         = "\ufb2b" // שׂ
	DETACHED_HEAD_ICON  = "\ue729" // 
	TAG_ICON            = "\uf02b" // 
	COMMIT_ICON         = "\ufc16" // ﰖ
	MERGE_COMMIT_ICON   = "\ufb2c" // שּׁ
	DEFAULT_REMOTE_ICON = "\uf7a1" // 
	STAGED_FILE_ICON    = "\uf833 " // 
)

type remoteIcon struct {
	domain string
	icon   string
}

var remoteIcons = []remoteIcon{
	{domain: "github.com", icon: "\ue709"},    // 
	{domain: "bitbucket.org", icon: "\ue703"}, // 
	{domain: "gitlab.com", icon: "\uf296"},    // 
	{domain: "dev.azure.com", icon: "\ufd03"}, // ﴃ
}

func IconForBranch(branch *models.Branch) string {
	if branch.DisplayName != "" {
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


func IconForChangeStatus(changeStatus string) string {
	nfOctDiffAdded    := "\uf457 " // 
	nfOctDiffModified := "\uf459 " // 
	nfOctDiffRemoved  := "\uf458 " // 
	nfOctDiffRenamed  := "\uf45a " // 

	switch changeStatus {
	case "?", "A":
		return nfOctDiffAdded
	case "M", "C":
		return nfOctDiffModified
	case "R", "T":
		return nfOctDiffRenamed
	case "D":
		return nfOctDiffRemoved
	default:
		return "  "
	}
}
