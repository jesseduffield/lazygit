package filter_by_author

import (
	"fmt"
	"strings"

	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

type AuthorInfo struct {
	name            string
	numberOfCommits int
}

func commonSetup(shell *Shell) {
	authors := []AuthorInfo{{"Yang Wen-li", 3}, {"Siegfried Kircheis", 1}, {"Paul Oberstein", 8}}
	totalCommits := 0
	repoStartDaysAgo := 100

	for _, authorInfo := range authors {
		for i := 0; i < authorInfo.numberOfCommits; i++ {
			authorEmail := strings.ToLower(strings.ReplaceAll(authorInfo.name, " ", ".")) + "@email.com"
			commitMessage := fmt.Sprintf("commit %d", i)

			shell.SetAuthor(authorInfo.name, authorEmail)
			shell.EmptyCommitDaysAgo(commitMessage, repoStartDaysAgo-totalCommits)
			totalCommits++
		}
	}
}
