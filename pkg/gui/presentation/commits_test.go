package presentation

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/gookit/color"
	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func init() {
	color.ForceSetColorLevel(terminfo.ColorLevelNone)
}

func formatExpected(expected string) string {
	return strings.TrimSpace(strings.ReplaceAll(expected, "\t", ""))
}

func TestGetCommitListDisplayStrings(t *testing.T) {
	scenarios := []struct {
		testName                 string
		commits                  []*models.Commit
		branches                 []*models.Branch
		currentBranchName        string
		hasUpdateRefConfig       bool
		fullDescription          bool
		cherryPickedCommitShaSet *set.Set[string]
		markedBaseCommit         string
		diffName                 string
		timeFormat               string
		shortTimeFormat          string
		now                      time.Time
		parseEmoji               bool
		selectedCommitSha        string
		startIdx                 int
		endIdx                   int
		showGraph                bool
		bisectInfo               *git_commands.BisectInfo
		showYouAreHereLabel      bool
		expected                 string
		focus                    bool
	}{
		{
			testName:                 "no commits",
			commits:                  []*models.Commit{},
			startIdx:                 0,
			endIdx:                   1,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:                 "",
		},
		{
			testName: "some commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1"},
				{Name: "commit2", Sha: "sha2"},
			},
			startIdx:                 0,
			endIdx:                   2,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 commit1
		sha2 commit2
						`),
		},
		{
			testName: "commit with tags",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Tags: []string{"tag1", "tag2"}},
				{Name: "commit2", Sha: "sha2"},
			},
			startIdx:                 0,
			endIdx:                   2,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 tag1 tag2 commit1
		sha2 commit2
						`),
		},
		{
			testName: "show local branch head, except the current branch, main branches, or merged branches",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1"},
				{Name: "commit2", Sha: "sha2"},
				{Name: "commit3", Sha: "sha3"},
				{Name: "commit4", Sha: "sha4", Status: models.StatusMerged},
			},
			branches: []*models.Branch{
				{Name: "current-branch", CommitHash: "sha1", Head: true},
				{Name: "other-branch", CommitHash: "sha2", Head: false},
				{Name: "master", CommitHash: "sha3", Head: false},
				{Name: "old-branch", CommitHash: "sha4", Head: false},
			},
			currentBranchName:        "current-branch",
			hasUpdateRefConfig:       true,
			startIdx:                 0,
			endIdx:                   4,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 commit1
		sha2 * commit2
		sha3 commit3
		sha4 commit4
						`),
		},
		{
			testName: "show local branch head for head commit if updateRefs is on",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1"},
				{Name: "commit2", Sha: "sha2"},
			},
			branches: []*models.Branch{
				{Name: "current-branch", CommitHash: "sha1", Head: true},
				{Name: "other-branch", CommitHash: "sha1", Head: false},
			},
			currentBranchName:        "current-branch",
			hasUpdateRefConfig:       true,
			startIdx:                 0,
			endIdx:                   2,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 * commit1
		sha2 commit2
						`),
		},
		{
			testName: "don't show local branch head for head commit if updateRefs is off",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1"},
				{Name: "commit2", Sha: "sha2"},
			},
			branches: []*models.Branch{
				{Name: "current-branch", CommitHash: "sha1", Head: true},
				{Name: "other-branch", CommitHash: "sha1", Head: false},
			},
			currentBranchName:        "current-branch",
			hasUpdateRefConfig:       false,
			startIdx:                 0,
			endIdx:                   2,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 commit1
		sha2 commit2
						`),
		},
		{
			testName: "show local branch head and tag if both exist",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1"},
				{Name: "commit2", Sha: "sha2", Tags: []string{"some-tag"}},
				{Name: "commit3", Sha: "sha3"},
			},
			branches: []*models.Branch{
				{Name: "some-branch", CommitHash: "sha2"},
			},
			startIdx:                 0,
			endIdx:                   3,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 commit1
		sha2 * some-tag commit2
		sha3 commit3
						`),
		},
		{
			testName: "showing graph",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:                 0,
			endIdx:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 ⏣─╮ commit1
		sha2 ◯ │ commit2
		sha3 ◯─╯ commit3
		sha4 ◯ commit4
		sha5 ◯ commit5
						`),
		},
		{
			testName: "showing graph, including rebase commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: todo.Pick},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: todo.Pick},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:                 0,
			endIdx:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 pick  commit1
		sha2 pick  commit2
		sha3       ◯ <-- YOU ARE HERE --- commit3
		sha4       ◯ commit4
		sha5       ◯ commit5
				`),
		},
		{
			testName: "showing graph, including rebase commits, with offset",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: todo.Pick},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: todo.Pick},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:                 1,
			endIdx:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha2 pick  commit2
		sha3       ◯ <-- YOU ARE HERE --- commit3
		sha4       ◯ commit4
		sha5       ◯ commit5
				`),
		},
		{
			testName: "startIdx is past TODO commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: todo.Pick},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: todo.Pick},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:                 3,
			endIdx:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha4 ◯ commit4
		sha5 ◯ commit5
				`),
		},
		{
			testName: "only showing TODO commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: todo.Pick},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: todo.Pick},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:                 0,
			endIdx:                   2,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 pick  commit1
		sha2 pick  commit2
				`),
		},
		{
			testName: "no TODO commits, towards bottom",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:                 4,
			endIdx:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
			sha5 ◯ commit5
				`),
		},
		{
			testName: "only TODO commits except last",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: todo.Pick},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: todo.Pick},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}, Action: todo.Pick},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}, Action: todo.Pick},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:                 0,
			endIdx:                   2,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
			sha1 pick  commit1
			sha2 pick  commit2
				`),
		},
		{
			testName: "don't show YOU ARE HERE label when not asked for (e.g. in branches panel)",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2"}, Action: todo.Pick},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
			},
			startIdx:                 0,
			endIdx:                   3,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      false,
			now:                      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: formatExpected(`
		sha1 pick  commit1
		sha2       ◯ commit2
		sha3       ◯ commit3
				`),
		},
		{
			testName: "custom time format",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", UnixTimestamp: 1577844184, AuthorName: "Jesse Duffield"},
				{Name: "commit2", Sha: "sha2", UnixTimestamp: 1576844184, AuthorName: "Jesse Duffield"},
			},
			fullDescription:          true,
			timeFormat:               "2006-01-02",
			shortTimeFormat:          "3:04PM",
			startIdx:                 0,
			endIdx:                   2,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			now:                      time.Date(2020, 1, 1, 5, 3, 4, 0, time.UTC),
			expected: formatExpected(`
		sha1 2:03AM     Jesse Duffield    commit1
		sha2 2019-12-20 Jesse Duffield    commit2
						`),
		},
	}

	os.Setenv("TZ", "UTC")

	focusing := false
	for _, scenario := range scenarios {
		if scenario.focus {
			focusing = true
		}
	}

	common := utils.NewDummyCommon()

	for _, s := range scenarios {
		s := s
		if !focusing || s.focus {
			t.Run(s.testName, func(t *testing.T) {
				result := GetCommitListDisplayStrings(
					common,
					s.commits,
					s.branches,
					s.currentBranchName,
					s.hasUpdateRefConfig,
					s.fullDescription,
					s.cherryPickedCommitShaSet,
					s.diffName,
					s.markedBaseCommit,
					s.timeFormat,
					s.shortTimeFormat,
					s.now,
					s.parseEmoji,
					s.selectedCommitSha,
					s.startIdx,
					s.endIdx,
					s.showGraph,
					s.bisectInfo,
					s.showYouAreHereLabel,
				)

				renderedLines, _ := utils.RenderDisplayStrings(result, nil)
				renderedResult := strings.Join(renderedLines, "\n")
				t.Logf("\n%s", renderedResult)

				assert.EqualValues(t, s.expected, renderedResult)
			})
		}
	}
}
