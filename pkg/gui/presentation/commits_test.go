package presentation

import (
	"strings"
	"testing"

	"github.com/gookit/color"
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
		cherryPickedCommitShaMap map[string]bool
		diffName                 string
		parseEmoji               bool
		selectedCommitSha        string
		startIdx                 int
		length                   int
		showGraph                bool
		bisectInfo               *git_commands.BisectInfo
		expected                 string
		focus                    bool
	}{
		{
			testName:   "no commits",
			commits:    []*models.Commit{},
			startIdx:   0,
			length:     1,
			showGraph:  false,
			bisectInfo: git_commands.NewNullBisectInfo(),
			expected:   "",
		},
		{
			testName: "some commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1"},
				{Name: "commit2", Sha: "sha2"},
			},
			startIdx:   0,
			length:     2,
			showGraph:  false,
			bisectInfo: git_commands.NewNullBisectInfo(),
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
			startIdx:   0,
			length:     5,
			showGraph:  true,
			bisectInfo: git_commands.NewNullBisectInfo(),
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
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: "pick"},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: "pick"},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:   0,
			length:     5,
			showGraph:  true,
			bisectInfo: git_commands.NewNullBisectInfo(),
			expected: formatExpected(`
		sha1 pick  commit1
		sha2 pick  commit2
		sha3       ◯ commit3
		sha4       ◯ commit4
		sha5       ◯ commit5
				`),
		},
		{
			testName: "showing graph, including rebase commits, with offset",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: "pick"},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: "pick"},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:   1,
			length:     10,
			showGraph:  true,
			bisectInfo: git_commands.NewNullBisectInfo(),
			expected: formatExpected(`
		sha2 pick  commit2
		sha3       ◯ commit3
		sha4       ◯ commit4
		sha5       ◯ commit5
				`),
		},
		{
			testName: "startIdx is passed TODO commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: "pick"},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: "pick"},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:   3,
			length:     2,
			showGraph:  true,
			bisectInfo: git_commands.NewNullBisectInfo(),
			expected: formatExpected(`
		sha4 ◯ commit4
		sha5 ◯ commit5
				`),
		},
		{
			testName: "only showing TODO commits",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: "pick"},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: "pick"},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:   0,
			length:     2,
			showGraph:  true,
			bisectInfo: git_commands.NewNullBisectInfo(),
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
			startIdx:   4,
			length:     2,
			showGraph:  true,
			bisectInfo: git_commands.NewNullBisectInfo(),
			expected: formatExpected(`
			sha5 ◯ commit5
				`),
		},
		{
			testName: "only TODO commits except last",
			commits: []*models.Commit{
				{Name: "commit1", Sha: "sha1", Parents: []string{"sha2", "sha3"}, Action: "pick"},
				{Name: "commit2", Sha: "sha2", Parents: []string{"sha3"}, Action: "pick"},
				{Name: "commit3", Sha: "sha3", Parents: []string{"sha4"}, Action: "pick"},
				{Name: "commit4", Sha: "sha4", Parents: []string{"sha5"}, Action: "pick"},
				{Name: "commit5", Sha: "sha5", Parents: []string{"sha7"}},
			},
			startIdx:   0,
			length:     2,
			showGraph:  true,
			bisectInfo: git_commands.NewNullBisectInfo(),
			expected: formatExpected(`
			sha1 pick  commit1
			sha2 pick  commit2
				`),
		},
	}

	focusing := false
	for _, scenario := range scenarios {
		if scenario.focus {
			focusing = true
		}
	}

	for _, s := range scenarios {
		s := s
		if !focusing || s.focus {
			t.Run(s.testName, func(t *testing.T) {
				result := GetCommitListDisplayStrings(
					s.commits,
					s.fullDescription,
					s.cherryPickedCommitShaMap,
					s.diffName,
					s.parseEmoji,
					s.selectedCommitSha,
					s.startIdx,
					s.length,
					s.showGraph,
					s.bisectInfo,
				)

				renderedResult := utils.RenderDisplayStrings(result)
				t.Logf("\n%s", renderedResult)

				assert.EqualValues(t, s.expected, renderedResult)
			})
		}
	}
}
