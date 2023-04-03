package presentation

import (
	"os"
	"strings"
	"testing"

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
		fullDescription          bool
		cherryPickedCommitShaSet *set.Set[string]
		diffName                 string
		timeFormat               string
		parseEmoji               bool
		selectedCommitSha        string
		startIdx                 int
		length                   int
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
			length:                   1,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			expected:                 "",
		},
		{
			testName: "some commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1"},
				{Name: "commit2", Sha: "sha2"},
			},
			startIdx:                 0,
			length:                   2,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			expected: formatExpected(`
		sha1 commit1
		sha2 commit2
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
			length:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
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
			length:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
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
			length:                   10,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
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
			length:                   2,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
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
			length:                   2,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
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
			length:                   2,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
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
			length:                   2,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      true,
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
			length:                   5,
			showGraph:                true,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			showYouAreHereLabel:      false,
			expected: formatExpected(`
		sha1 pick  commit1
		sha2       ◯ commit2
		sha3       ◯ commit3
				`),
		},
		{
			testName: "custom time format",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", UnixTimestamp: 1652443200, AuthorName: "Jesse Duffield"},
				{Name: "commit2", Sha: "sha2", UnixTimestamp: 1652529600, AuthorName: "Jesse Duffield"},
			},
			fullDescription:          true,
			timeFormat:               "2006-01-02 15:04:05",
			startIdx:                 0,
			length:                   2,
			showGraph:                false,
			bisectInfo:               git_commands.NewNullBisectInfo(),
			cherryPickedCommitShaSet: set.New[string](),
			expected: formatExpected(`
		sha1 2022-05-13 12:00:00 Jesse Duffield    commit1
		sha2 2022-05-14 12:00:00 Jesse Duffield    commit2
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
					s.fullDescription,
					s.cherryPickedCommitShaSet,
					s.diffName,
					s.timeFormat,
					s.parseEmoji,
					s.selectedCommitSha,
					s.startIdx,
					s.length,
					s.showGraph,
					s.bisectInfo,
					s.showYouAreHereLabel,
				)

				renderedResult := utils.RenderDisplayStrings(result)
				t.Logf("\n%s", renderedResult)

				assert.EqualValues(t, s.expected, renderedResult)
			})
		}
	}
}
