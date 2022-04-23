package icons

import "github.com/jesseduffield/lazygit/pkg/commands/models"

const BRANCH_ICON = "\ufb2b"       // שׂ
const COMMIT_ICON = "\ufc16"       // ﰖ
const MERGE_COMMIT_ICON = "\ufb2c" // שּׁ

func IconForCommit(commit *models.Commit) string {
	if len(commit.Parents) > 1 {
		return MERGE_COMMIT_ICON
	}
	return COMMIT_ICON
}
