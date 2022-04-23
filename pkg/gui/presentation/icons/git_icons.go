package icons

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

const BRANCH_ICON = "\ufb2b"        // שׂ
const DETACHED_HEAD_ICON = "\ue729" // 
const TAG_ICON = "\uf02b"           // 
const COMMIT_ICON = "\ufc16"        // ﰖ
const MERGE_COMMIT_ICON = "\ufb2c"  // שּׁ

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
