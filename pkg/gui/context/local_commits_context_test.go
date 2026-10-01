package context

import (
	"testing"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestAddCommitDropIndicator(t *testing.T) {
	pendingHeader := &NonModelItem{Index: 0, Content: "pending"}
	commitsHeader := &NonModelItem{Index: 3, Content: "commits"}
	indicator := &commitDropIndicator{insertionIndex: 3}
	spinnerConfig := config.SpinnerConfig{Frames: []string{"one", "two"}, Rate: 100}

	items := addCommitDropIndicator(
		[]*NonModelItem{pendingHeader}, indicator, "drop here", "moving commits here", spinnerConfig, time.UnixMilli(0),
	)
	items = append(items, commitsHeader)

	assert.Equal(t, []*NonModelItem{
		pendingHeader,
		{
			Index:   3,
			Content: style.FgCyan.SetBold().Sprint("━━━━━━ drop here ━━━━━━"),
			Column:  6,
		},
		commitsHeader,
	}, items)
	assert.Equal(t, 6, modelIndexToViewIndex(4, items, 3))
	assert.Equal(t, 3, viewIndexToModelIndex(4, items, 4))
}

func TestAddMovingCommitsIndicator(t *testing.T) {
	items := addCommitDropIndicator(
		nil,
		&commitDropIndicator{insertionIndex: 2, moving: true},
		"drop here",
		"moving commits here",
		config.SpinnerConfig{Frames: []string{"one", "two"}, Rate: 100},
		time.UnixMilli(100),
	)

	assert.Equal(t, []*NonModelItem{
		{
			Index:   2,
			Content: style.FgCyan.SetBold().Sprint("━━━━━━ moving commits here two ━━━━━━"),
			Column:  6,
		},
	}, items)
}

func TestCommitsShownInDiff(t *testing.T) {
	hashPool := &utils.StringPool{}
	newer := models.NewCommit(hashPool, models.NewCommitOpts{Hash: "newer"})
	older := models.NewCommit(hashPool, models.NewCommitOpts{Hash: "older"})
	selected := []*models.Commit{newer, older}

	scenarios := []struct {
		name           string
		selectedCommit *models.Commit
		refRange       *types.RefRange
		expected       []*models.Commit
	}{
		{
			name:           "a range is diffed as a whole",
			selectedCommit: newer,
			refRange:       &types.RefRange{From: older, To: newer},
			expected:       selected,
		},
		{
			name:           "without a range to diff, only the commit at the cursor is",
			selectedCommit: newer,
			expected:       []*models.Commit{newer},
		},
		{
			name: "nothing is diffed while nothing is selected",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, commitsShownInDiff(selected, s.selectedCommit, s.refRange))
		})
	}
}

func TestPullRequestBaseForCommits(t *testing.T) {
	hashPool := &utils.StringPool{}
	commit := func(hash string, parent string, status models.CommitStatus) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{
			Hash:    hash,
			Parents: lo.Ternary(parent == "", []string{}, []string{parent}),
			Status:  status,
		})
	}

	// A branch of three pushed commits whose tip was amended, on top of a commit that
	// is in a main branch already, as the panel lists them: newest first.
	amended := commit("amended", "third", models.StatusUnpushed)
	third := commit("third", "second", models.StatusPushed)
	second := commit("second", "first", models.StatusPushed)
	first := commit("first", "merged", models.StatusPushed)
	merged := commit("merged", "ancient", models.StatusMerged)
	ancient := commit("ancient", "", models.StatusMerged)
	allCommits := []*models.Commit{amended, third, second, first, merged, ancient}

	scenarios := []struct {
		name     string
		commits  []*models.Commit
		expected string
	}{
		{
			name:     "a single commit starts after its parent",
			commits:  []*models.Commit{second},
			expected: "first",
		},
		{
			name:     "a range starts after the parent of its oldest commit",
			commits:  []*models.Commit{third, second},
			expected: "first",
		},
		{
			name:    "the pull request's first commit starts where the pull request does",
			commits: []*models.Commit{first},
		},
		{
			name:    "so does a range reaching down to it",
			commits: []*models.Commit{third, second, first},
		},
		{
			name:    "and so does the first commit of the repository",
			commits: []*models.Commit{ancient},
		},
		{
			name:    "a parent the panel doesn't list is none of the pull request's",
			commits: []*models.Commit{commit("elsewhere", "unlisted", models.StatusPushed)},
		},
		{
			name: "nothing is shown",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, pullRequestBaseForCommits(allCommits, s.commits))
		})
	}
}
