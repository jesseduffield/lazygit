package presentation

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gookit/color"
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func formatExpected(expected string) string {
	return strings.TrimSpace(strings.ReplaceAll(expected, "\t", ""))
}

func TestGetCommitListDisplayStrings(t *testing.T) {
	scenarios := []struct {
		testName                  string
		commitOpts                []models.NewCommitOpts
		branches                  []*models.Branch
		currentBranchName         string
		hasUpdateRefConfig        bool
		fullDescription           bool
		cherryPickedCommitHashSet *set.Set[string]
		markedBaseCommit          string
		diffName                  string
		timeFormat                string
		shortTimeFormat           string
		now                       time.Time
		parseEmoji                bool
		selectedCommitHashPtr     *string
		startIdx                  int
		endIdx                    int
		showGraph                 bool
		bisectInfo                *git_commands.BisectInfo
		expected                  string
		focus                     bool
		setupConfig               func(c *common.Common)
	}{
		{
			testName:                  "no commits",
			commitOpts:                []models.NewCommitOpts{},
			startIdx:                  0,
			endIdx:                    1,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:                  "",
		},
		{
			testName: "some commits",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1"},
				{Name: "commit2", Hash: "hash2"},
			},
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 commit1
		hash2 commit2
						`),
		},
		{
			testName: "some commits formatted with email",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", AuthorEmail: "email1@test.com"},
				{Name: "commit2", Hash: "hash2", AuthorEmail: "email2@test.com"},
			},
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			setupConfig: func(c *common.Common) {
				c.UserConfig().Git.Log.CustomPaneLogFormat = "%ae %s"
			},
			expected: formatExpected(`
		email1@test.com commit1
		email2@test.com commit2
						`),
		},
		{
			testName: "commit with tags",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Tags: []string{"tag1", "tag2"}},
				{Name: "commit2", Hash: "hash2"},
			},
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 tag1 tag2 commit1
		hash2 commit2
						`),
		},
		{
			testName: "show local branch head, except the current branch, main branches, or merged branches",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1"},
				{Name: "commit2", Hash: "hash2"},
				{Name: "commit3", Hash: "hash3"},
				{Name: "commit4", Hash: "hash4", Status: models.StatusMerged},
			},
			branches: []*models.Branch{
				{Name: "current-branch", CommitHash: "hash1", Head: true},
				{Name: "other-branch", CommitHash: "hash2", Head: false},
				{Name: "master", CommitHash: "hash3", Head: false},
				{Name: "old-branch", CommitHash: "hash4", Head: false},
			},
			currentBranchName:         "current-branch",
			hasUpdateRefConfig:        true,
			startIdx:                  0,
			endIdx:                    4,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 commit1
		hash2 * commit2
		hash3 commit3
		hash4 commit4
						`),
		},
		{
			testName: "show local branch head for head commit if updateRefs is on",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1"},
				{Name: "commit2", Hash: "hash2"},
			},
			branches: []*models.Branch{
				{Name: "current-branch", CommitHash: "hash1", Head: true},
				{Name: "other-branch", CommitHash: "hash1", Head: false},
			},
			currentBranchName:         "current-branch",
			hasUpdateRefConfig:        true,
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 * commit1
		hash2 commit2
						`),
		},
		{
			testName: "don't show local branch head for head commit if updateRefs is off",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1"},
				{Name: "commit2", Hash: "hash2"},
			},
			branches: []*models.Branch{
				{Name: "current-branch", CommitHash: "hash1", Head: true},
				{Name: "other-branch", CommitHash: "hash1", Head: false},
			},
			currentBranchName:         "current-branch",
			hasUpdateRefConfig:        false,
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 commit1
		hash2 commit2
						`),
		},
		{
			testName: "show local branch head and tag if both exist",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1"},
				{Name: "commit2", Hash: "hash2", Tags: []string{"some-tag"}},
				{Name: "commit3", Hash: "hash3"},
			},
			branches: []*models.Branch{
				{Name: "some-branch", CommitHash: "hash2"},
			},
			startIdx:                  0,
			endIdx:                    3,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 commit1
		hash2 * some-tag commit2
		hash3 commit3
						`),
		},
		{
			testName: "showing graph",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  0,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 ⏣─╮ commit1
		hash2 ◯ │ commit2
		hash3 ◯─╯ commit3
		hash4 ◯ commit4
		hash5 ◯ commit5
						`),
		},
		{
			testName: "showing graph with custom formatting",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  0,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			setupConfig: func(c *common.Common) {
				c.UserConfig().Git.Log.CustomPaneLogFormat = "%h %g"
			},
			expected: formatExpected(`
		hash1 ⏣─╮
		hash2 ◯ │
		hash3 ◯─╯
		hash4 ◯
		hash5 ◯
						`),
		},
		{
			testName: "showing graph, including rebase commits",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}, Action: todo.Pick},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}, Action: todo.Pick},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  0,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 pick commit1
		hash2 pick commit2
		hash3      ◯ commit3
		hash4      ◯ commit4
		hash5      ◯ commit5
				`),
		},
		{
			testName: "showing graph, including rebase commits, with offset",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}, Action: todo.Pick},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}, Action: todo.Pick},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  1,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash2 pick commit2
		hash3      ◯ commit3
		hash4      ◯ commit4
		hash5      ◯ commit5
				`),
		},
		{
			testName: "startIdx is past TODO commits",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}, Action: todo.Pick},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}, Action: todo.Pick},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  3,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash4 ◯ commit4
		hash5 ◯ commit5
				`),
		},
		{
			testName: "only showing TODO commits",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}, Action: todo.Pick},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}, Action: todo.Pick},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		hash1 pick commit1
		hash2 pick commit2
				`),
		},
		{
			testName: "no TODO commits, towards bottom",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  4,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
			hash5 ◯ commit5
				`),
		},
		{
			testName: "only TODO commits except last",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", Parents: []string{"hash2", "hash3"}, Action: todo.Pick},
				{Name: "commit2", Hash: "hash2", Parents: []string{"hash3"}, Action: todo.Pick},
				{Name: "commit3", Hash: "hash3", Parents: []string{"hash4"}, Action: todo.Pick},
				{Name: "commit4", Hash: "hash4", Parents: []string{"hash5"}, Action: todo.Pick},
				{Name: "commit5", Hash: "hash5", Parents: []string{"hash7"}},
			},
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
			hash1 pick commit1
			hash2 pick commit2
				`),
		},
		{
			testName: "graph in divergence view - all commits visible",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1r", Parents: []string{"hash2r"}, Divergence: models.DivergenceRight},
				{Name: "commit2", Hash: "hash2r", Parents: []string{"hash3r", "hash5r"}, Divergence: models.DivergenceRight},
				{Name: "commit3", Hash: "hash3r", Parents: []string{"hash4r"}, Divergence: models.DivergenceRight},
				{Name: "commit1", Hash: "hash1l", Parents: []string{"hash2l"}, Divergence: models.DivergenceLeft},
				{Name: "commit2", Hash: "hash2l", Parents: []string{"hash3l", "hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit3", Hash: "hash3l", Parents: []string{"hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit4", Hash: "hash4l", Parents: []string{"hash5l"}, Divergence: models.DivergenceLeft},
				{Name: "commit5", Hash: "hash5l", Parents: []string{"hash6l"}, Divergence: models.DivergenceLeft},
			},
			startIdx:                  0,
			endIdx:                    8,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		↓ hash1r ◯ commit1
		↓ hash2r ⏣─╮ commit2
		↓ hash3r ◯ │ commit3
		↑ hash1l ◯ commit1
		↑ hash2l ⏣─╮ commit2
		↑ hash3l ◯ │ commit3
		↑ hash4l ◯─╯ commit4
		↑ hash5l ◯ commit5
				`),
		},
		{
			testName: "graph in divergence view - not all remote commits visible",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1r", Parents: []string{"hash2r"}, Divergence: models.DivergenceRight},
				{Name: "commit2", Hash: "hash2r", Parents: []string{"hash3r", "hash5r"}, Divergence: models.DivergenceRight},
				{Name: "commit3", Hash: "hash3r", Parents: []string{"hash4r"}, Divergence: models.DivergenceRight},
				{Name: "commit1", Hash: "hash1l", Parents: []string{"hash2l"}, Divergence: models.DivergenceLeft},
				{Name: "commit2", Hash: "hash2l", Parents: []string{"hash3l", "hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit3", Hash: "hash3l", Parents: []string{"hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit4", Hash: "hash4l", Parents: []string{"hash5l"}, Divergence: models.DivergenceLeft},
				{Name: "commit5", Hash: "hash5l", Parents: []string{"hash6l"}, Divergence: models.DivergenceLeft},
			},
			startIdx:                  2,
			endIdx:                    8,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		↓ hash3r ◯ │ commit3
		↑ hash1l ◯ commit1
		↑ hash2l ⏣─╮ commit2
		↑ hash3l ◯ │ commit3
		↑ hash4l ◯─╯ commit4
		↑ hash5l ◯ commit5
				`),
		},
		{
			testName: "graph in divergence view - not all local commits",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1r", Parents: []string{"hash2r"}, Divergence: models.DivergenceRight},
				{Name: "commit2", Hash: "hash2r", Parents: []string{"hash3r", "hash5r"}, Divergence: models.DivergenceRight},
				{Name: "commit3", Hash: "hash3r", Parents: []string{"hash4r"}, Divergence: models.DivergenceRight},
				{Name: "commit1", Hash: "hash1l", Parents: []string{"hash2l"}, Divergence: models.DivergenceLeft},
				{Name: "commit2", Hash: "hash2l", Parents: []string{"hash3l", "hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit3", Hash: "hash3l", Parents: []string{"hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit4", Hash: "hash4l", Parents: []string{"hash5l"}, Divergence: models.DivergenceLeft},
				{Name: "commit5", Hash: "hash5l", Parents: []string{"hash6l"}, Divergence: models.DivergenceLeft},
			},
			startIdx:                  0,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		↓ hash1r ◯ commit1
		↓ hash2r ⏣─╮ commit2
		↓ hash3r ◯ │ commit3
		↑ hash1l ◯ commit1
		↑ hash2l ⏣─╮ commit2
				`),
		},
		{
			testName: "graph in divergence view - no remote commits visible",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1r", Parents: []string{"hash2r"}, Divergence: models.DivergenceRight},
				{Name: "commit2", Hash: "hash2r", Parents: []string{"hash3r", "hash5r"}, Divergence: models.DivergenceRight},
				{Name: "commit3", Hash: "hash3r", Parents: []string{"hash4r"}, Divergence: models.DivergenceRight},
				{Name: "commit1", Hash: "hash1l", Parents: []string{"hash2l"}, Divergence: models.DivergenceLeft},
				{Name: "commit2", Hash: "hash2l", Parents: []string{"hash3l", "hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit3", Hash: "hash3l", Parents: []string{"hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit4", Hash: "hash4l", Parents: []string{"hash5l"}, Divergence: models.DivergenceLeft},
				{Name: "commit5", Hash: "hash5l", Parents: []string{"hash6l"}, Divergence: models.DivergenceLeft},
			},
			startIdx:                  4,
			endIdx:                    8,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		↑ hash2l ⏣─╮ commit2
		↑ hash3l ◯ │ commit3
		↑ hash4l ◯─╯ commit4
		↑ hash5l ◯ commit5
				`),
		},
		{
			testName: "graph in divergence view - no local commits visible",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1r", Parents: []string{"hash2r"}, Divergence: models.DivergenceRight},
				{Name: "commit2", Hash: "hash2r", Parents: []string{"hash3r", "hash5r"}, Divergence: models.DivergenceRight},
				{Name: "commit3", Hash: "hash3r", Parents: []string{"hash4r"}, Divergence: models.DivergenceRight},
				{Name: "commit1", Hash: "hash1l", Parents: []string{"hash2l"}, Divergence: models.DivergenceLeft},
				{Name: "commit2", Hash: "hash2l", Parents: []string{"hash3l", "hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit3", Hash: "hash3l", Parents: []string{"hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit4", Hash: "hash4l", Parents: []string{"hash5l"}, Divergence: models.DivergenceLeft},
				{Name: "commit5", Hash: "hash5l", Parents: []string{"hash6l"}, Divergence: models.DivergenceLeft},
			},
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		↓ hash1r ◯ commit1
		↓ hash2r ⏣─╮ commit2
				`),
		},
		{
			testName: "graph in divergence view - no remote commits present",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1l", Parents: []string{"hash2l"}, Divergence: models.DivergenceLeft},
				{Name: "commit2", Hash: "hash2l", Parents: []string{"hash3l", "hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit3", Hash: "hash3l", Parents: []string{"hash4l"}, Divergence: models.DivergenceLeft},
				{Name: "commit4", Hash: "hash4l", Parents: []string{"hash5l"}, Divergence: models.DivergenceLeft},
				{Name: "commit5", Hash: "hash5l", Parents: []string{"hash6l"}, Divergence: models.DivergenceLeft},
			},
			startIdx:                  0,
			endIdx:                    5,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		↑ hash1l ◯ commit1
		↑ hash2l ⏣─╮ commit2
		↑ hash3l ◯ │ commit3
		↑ hash4l ◯─╯ commit4
		↑ hash5l ◯ commit5
				`),
		},
		{
			testName: "graph in divergence view - no local commits present",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1r", Parents: []string{"hash2r"}, Divergence: models.DivergenceRight},
				{Name: "commit2", Hash: "hash2r", Parents: []string{"hash3r", "hash5r"}, Divergence: models.DivergenceRight},
				{Name: "commit3", Hash: "hash3r", Parents: []string{"hash4r"}, Divergence: models.DivergenceRight},
			},
			startIdx:                  0,
			endIdx:                    3,
			showGraph:                 true,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		↓ hash1r ◯ commit1
		↓ hash2r ⏣─╮ commit2
		↓ hash3r ◯ │ commit3
				`),
		},
		{
			testName: "custom time format",
			commitOpts: []models.NewCommitOpts{
				{Name: "commit1", Hash: "hash1", UnixTimestamp: 1577844184, AuthorName: "Jesse Duffield"},
				{Name: "commit2", Hash: "hash2", UnixTimestamp: 1576844184, AuthorName: "Jesse Duffield"},
			},
			fullDescription:           true,
			timeFormat:                "2006-01-02",
			shortTimeFormat:           "3:04PM",
			startIdx:                  0,
			endIdx:                    2,
			showGraph:                 false,
			bisectInfo:                git_commands.NewNullBisectInfo(),
			cherryPickedCommitHashSet: set.New[string](),
			now:                       time.Date(2020, 1, 1, 5, 3, 4, 0, time.UTC),
			expected: formatExpected(`
		hash1 2:03AM     Jesse Duffield    commit1
		hash2 2019-12-20 Jesse Duffield    commit2
						`),
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelNone)
	defer color.ForceSetColorLevel(oldColorLevel)

	os.Setenv("TZ", "UTC")

	focusing := false
	for _, scenario := range scenarios {
		if scenario.focus {
			focusing = true
		}
	}

	for _, s := range scenarios {
		if !focusing || s.focus {
			t.Run(s.testName, func(t *testing.T) {
				common := common.NewDummyCommon()

				hashPool := &utils.StringPool{}

				commits := lo.Map(s.commitOpts,
					func(opts models.NewCommitOpts, _ int) *models.Commit { return models.NewCommit(hashPool, opts) })

				if s.setupConfig != nil {
					s.setupConfig(common)
				}

				result := GetCommitListDisplayStrings(
					common,
					commits,
					s.branches,
					s.currentBranchName,
					s.hasUpdateRefConfig,
					s.fullDescription,
					s.cherryPickedCommitHashSet,
					s.diffName,
					s.markedBaseCommit,
					s.timeFormat,
					s.shortTimeFormat,
					s.now,
					s.parseEmoji,
					s.selectedCommitHashPtr,
					s.startIdx,
					s.endIdx,
					s.showGraph,
					s.bisectInfo,
				)

				renderedLines, _ := utils.RenderDisplayStrings(result, nil)
				renderedResult := strings.Join(renderedLines, "\n")
				t.Logf("\n%s", renderedResult)

				assert.EqualValues(t, s.expected, renderedResult)
			})
		}
	}
}
