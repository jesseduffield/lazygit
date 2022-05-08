package graph

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation/authors"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/xo/terminfo"
)

func init() {
	// on CI we've got no color capability so we're forcing it here
	color.ForceSetColorLevel(terminfo.ColorLevelMillions)
}

func TestRenderCommitGraph(t *testing.T) {
	tests := []struct {
		name           string
		commits        []*models.Commit
		expectedOutput string
	}{
		{
			name: "with some merges",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3"}},
				{Sha: "3", Parents: []string{"4"}},
				{Sha: "4", Parents: []string{"5", "7"}},
				{Sha: "7", Parents: []string{"5"}},
				{Sha: "5", Parents: []string{"8"}},
				{Sha: "8", Parents: []string{"9"}},
				{Sha: "9", Parents: []string{"A", "B"}},
				{Sha: "B", Parents: []string{"D"}},
				{Sha: "D", Parents: []string{"D"}},
				{Sha: "A", Parents: []string{"E"}},
				{Sha: "E", Parents: []string{"F"}},
				{Sha: "F", Parents: []string{"D"}},
				{Sha: "D", Parents: []string{"G"}},
			},
			expectedOutput: `
			1 ◯
			2 ◯
			3 ◯
			4 ⏣─╮
			7 │ ◯
			5 ◯─╯
			8 ◯
			9 ⏣─╮
			B │ ◯
			D │ ◯
			A ◯ │
			E ◯ │
			F ◯ │
			D ◯─╯`,
		},
		{
			name: "with a path that has room to move to the left",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "4", Parents: []string{"3", "5"}},
				{Sha: "3", Parents: []string{"5"}},
				{Sha: "5", Parents: []string{"6"}},
				{Sha: "6", Parents: []string{"7"}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			4 │ ⏣─╮
			3 ◯─╯ │
			5 ◯───╯
			6 ◯`,
		},
		{
			name: "with a new commit",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "4", Parents: []string{"3", "5"}},
				{Sha: "Z", Parents: []string{"Z"}},
				{Sha: "3", Parents: []string{"5"}},
				{Sha: "5", Parents: []string{"6"}},
				{Sha: "6", Parents: []string{"7"}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			4 │ ⏣─╮
			Z │ │ │ ◯
			3 ◯─╯ │ │
			5 ◯───╯ │
			6 ◯ ╭───╯`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "3", Parents: []string{"5", "4"}},
				{Sha: "5", Parents: []string{"7", "8"}},
				{Sha: "4", Parents: []string{"7"}},
				{Sha: "7", Parents: []string{"11"}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			4 │ ◯─╯ │
			7 ◯─╯ ╭─╯`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "3", Parents: []string{"5", "4"}},
				{Sha: "5", Parents: []string{"7", "8"}},
				{Sha: "7", Parents: []string{"4", "A"}},
				{Sha: "4", Parents: []string{"B"}},
				{Sha: "B", Parents: []string{"C"}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			7 ⏣─│─│─│─╮
			4 ◯─┴─╯ │ │
			B ◯ ╭───╯ │`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2", "3"}},
				{Sha: "3", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"4", "5"}},
				{Sha: "4", Parents: []string{"6", "7"}},
				{Sha: "6", Parents: []string{"8"}},
			},
			expectedOutput: `
			1 ⏣─╮
			3 │ ◯
			2 ⏣─│
			4 ⏣─│─╮
			6 ◯ │ │`,
		},
		{
			name: "new merge path fills gap before continuing path on right",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2", "3", "4", "5"}},
				{Sha: "4", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"A"}},
				{Sha: "A", Parents: []string{"6", "B"}},
				{Sha: "B", Parents: []string{"C"}},
			},
			expectedOutput: `
			1 ⏣─┬─┬─╮
			4 │ │ ◯ │
			2 ◯─│─╯ │
			A ⏣─│─╮ │
			B │ │ ◯ │`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "3", Parents: []string{"5", "4"}},
				{Sha: "5", Parents: []string{"7", "8"}},
				{Sha: "7", Parents: []string{"4", "A"}},
				{Sha: "4", Parents: []string{"B"}},
				{Sha: "B", Parents: []string{"C"}},
				{Sha: "C", Parents: []string{"D"}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			7 ⏣─│─│─│─╮
			4 ◯─┴─╯ │ │
			B ◯ ╭───╯ │
			C ◯ │ ╭───╯`,
		},
		{
			name: "with a path that has room to move to the left and continues",
			commits: []*models.Commit{
				{Sha: "1", Parents: []string{"2"}},
				{Sha: "2", Parents: []string{"3", "4"}},
				{Sha: "3", Parents: []string{"5", "4"}},
				{Sha: "5", Parents: []string{"7", "G"}},
				{Sha: "7", Parents: []string{"8", "A"}},
				{Sha: "8", Parents: []string{"4", "E"}},
				{Sha: "4", Parents: []string{"B"}},
				{Sha: "B", Parents: []string{"C"}},
				{Sha: "C", Parents: []string{"D"}},
				{Sha: "D", Parents: []string{"F"}},
			},
			expectedOutput: `
			1 ◯
			2 ⏣─╮
			3 ⏣─│─╮
			5 ⏣─│─│─╮
			7 ⏣─│─│─│─╮
			8 ⏣─│─│─│─│─╮
			4 ◯─┴─╯ │ │ │
			B ◯ ╭───╯ │ │
			C ◯ │ ╭───╯ │
			D ◯ │ │ ╭───╯`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			getStyle := func(c *models.Commit) style.TextStyle { return style.FgDefault }
			lines := RenderCommitGraph(test.commits, "blah", getStyle)

			trimmedExpectedOutput := ""
			for _, line := range strings.Split(strings.TrimPrefix(test.expectedOutput, "\n"), "\n") {
				trimmedExpectedOutput += strings.TrimSpace(line) + "\n"
			}

			t.Log("\nexpected: \n" + trimmedExpectedOutput)

			output := ""
			for i, line := range lines {
				description := test.commits[i].Sha
				output += strings.TrimSpace(description+" "+utils.Decolorise(line)) + "\n"
			}
			t.Log("\nactual: \n" + output)

			assert.Equal(t,
				trimmedExpectedOutput,
				output)
		})
	}
}

func TestRenderPipeSet(t *testing.T) {
	cyan := style.FgCyan
	red := style.FgRed
	green := style.FgGreen
	// blue := style.FgBlue
	yellow := style.FgYellow
	magenta := style.FgMagenta
	nothing := style.Nothing

	tests := []struct {
		name           string
		pipes          []*Pipe
		commit         *models.Commit
		prevCommit     *models.Commit
		expectedStr    string
		expectedStyles []style.TextStyle
	}{
		{
			name: "single cell",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: TERMINATES, style: cyan},
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "c", kind: STARTS, style: green},
			},
			prevCommit:     &models.Commit{Sha: "a"},
			expectedStr:    "◯",
			expectedStyles: []style.TextStyle{green},
		},
		{
			name: "single cell, selected",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "selected", kind: TERMINATES, style: cyan},
				{fromPos: 0, toPos: 0, fromSha: "selected", toSha: "c", kind: STARTS, style: green},
			},
			prevCommit:     &models.Commit{Sha: "a"},
			expectedStr:    "◯",
			expectedStyles: []style.TextStyle{highlightStyle},
		},
		{
			name: "terminating hook and starting hook, selected",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "selected", kind: TERMINATES, style: cyan},
				{fromPos: 1, toPos: 0, fromSha: "c", toSha: "selected", kind: TERMINATES, style: yellow},
				{fromPos: 0, toPos: 0, fromSha: "selected", toSha: "d", kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromSha: "selected", toSha: "e", kind: STARTS, style: green},
			},
			prevCommit:  &models.Commit{Sha: "a"},
			expectedStr: "⏣─╮",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle,
			},
		},
		{
			name: "terminating hook and starting hook, prioritise the terminating one",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 0, fromSha: "c", toSha: "b", kind: TERMINATES, style: magenta},
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "d", kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromSha: "b", toSha: "e", kind: STARTS, style: green},
			},
			prevCommit:  &models.Commit{Sha: "a"},
			expectedStr: "⏣─│",
			expectedStyles: []style.TextStyle{
				green, green, magenta,
			},
		},
		{
			name: "starting and terminating pipe sharing some space",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a1", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "a2", toSha: "a3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromSha: "b1", toSha: "b2", kind: CONTINUES, style: magenta},
				{fromPos: 3, toPos: 0, fromSha: "e1", toSha: "a2", kind: TERMINATES, style: green},
				{fromPos: 0, toPos: 2, fromSha: "a2", toSha: "c3", kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Sha: "a1"},
			expectedStr: "⏣─│─┬─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, magenta, yellow, yellow, green, green,
			},
		},
		{
			name: "starting and terminating pipe sharing some space, with selection",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a1", toSha: "selected", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "selected", toSha: "a3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromSha: "b1", toSha: "b2", kind: CONTINUES, style: magenta},
				{fromPos: 3, toPos: 0, fromSha: "e1", toSha: "selected", kind: TERMINATES, style: green},
				{fromPos: 0, toPos: 2, fromSha: "selected", toSha: "c3", kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Sha: "a1"},
			expectedStr: "⏣───╮ ╯",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle, highlightStyle, highlightStyle, nothing, green,
			},
		},
		{
			name: "many terminating pipes",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a1", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "a2", toSha: "a3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 0, fromSha: "b1", toSha: "a2", kind: TERMINATES, style: magenta},
				{fromPos: 2, toPos: 0, fromSha: "c1", toSha: "a2", kind: TERMINATES, style: green},
			},
			prevCommit:  &models.Commit{Sha: "a1"},
			expectedStr: "◯─┴─╯",
			expectedStyles: []style.TextStyle{
				yellow, magenta, magenta, green, green,
			},
		},
		{
			name: "starting pipe passing through",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a1", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "a2", toSha: "a3", kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 3, fromSha: "a2", toSha: "d3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromSha: "b1", toSha: "b3", kind: CONTINUES, style: magenta},
				{fromPos: 2, toPos: 2, fromSha: "c1", toSha: "c3", kind: CONTINUES, style: green},
			},
			prevCommit:  &models.Commit{Sha: "a1"},
			expectedStr: "⏣─│─│─╮",
			expectedStyles: []style.TextStyle{
				yellow, yellow, magenta, yellow, green, yellow, yellow,
			},
		},
		{
			name: "starting and terminating path crossing continuing path",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a1", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "a2", toSha: "a3", kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 1, fromSha: "a2", toSha: "b3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromSha: "b1", toSha: "a2", kind: CONTINUES, style: green},
				{fromPos: 2, toPos: 0, fromSha: "c1", toSha: "a2", kind: TERMINATES, style: magenta},
			},
			prevCommit:  &models.Commit{Sha: "a1"},
			expectedStr: "⏣─│─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, green, magenta, magenta,
			},
		},
		{
			name: "another clash of starting and terminating paths",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a1", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "a2", toSha: "a3", kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 1, fromSha: "a2", toSha: "b3", kind: STARTS, style: yellow},
				{fromPos: 2, toPos: 2, fromSha: "c1", toSha: "c3", kind: CONTINUES, style: green},
				{fromPos: 3, toPos: 0, fromSha: "d1", toSha: "a2", kind: TERMINATES, style: magenta},
			},
			prevCommit:  &models.Commit{Sha: "a1"},
			expectedStr: "⏣─┬─│─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, yellow, magenta, green, magenta, magenta,
			},
		},
		{
			name: "commit whose previous commit is selected",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "selected", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "a2", toSha: "a3", kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Sha: "selected"},
			expectedStr: "◯",
			expectedStyles: []style.TextStyle{
				yellow,
			},
		},
		{
			name: "commit whose previous commit is selected and is a merge commit",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "selected", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 1, fromSha: "selected", toSha: "b3", kind: CONTINUES, style: red},
			},
			prevCommit:  &models.Commit{Sha: "selected"},
			expectedStr: "◯ │",
			expectedStyles: []style.TextStyle{
				highlightStyle, nothing, highlightStyle,
			},
		},
		{
			name: "commit whose previous commit is selected and is a merge commit, with continuing pipe inbetween",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "selected", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 1, fromSha: "z1", toSha: "z3", kind: CONTINUES, style: green},
				{fromPos: 2, toPos: 2, fromSha: "selected", toSha: "b3", kind: CONTINUES, style: red},
			},
			prevCommit:  &models.Commit{Sha: "selected"},
			expectedStr: "◯ │ │",
			expectedStyles: []style.TextStyle{
				highlightStyle, nothing, green, nothing, highlightStyle,
			},
		},
		{
			name: "when previous commit is selected, not a merge commit, and spawns a continuing pipe",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a1", toSha: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromSha: "a2", toSha: "a3", kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromSha: "a2", toSha: "b3", kind: STARTS, style: green},
				{fromPos: 1, toPos: 0, fromSha: "selected", toSha: "a2", kind: TERMINATES, style: yellow},
			},
			prevCommit:  &models.Commit{Sha: "selected"},
			expectedStr: "⏣─╯",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actualStr := renderPipeSet(test.pipes, "selected", test.prevCommit)
			t.Log("actual cells:")
			t.Log(actualStr)
			expectedStr := ""
			if len([]rune(test.expectedStr)) != len(test.expectedStyles) {
				t.Fatalf("Error in test setup: you have %d characters in the expected output (%s) but have specified %d styles", len([]rune(test.expectedStr)), test.expectedStr, len(test.expectedStyles))
			}
			for i, char := range []rune(test.expectedStr) {
				expectedStr += test.expectedStyles[i].Sprint(string(char))
			}
			expectedStr += " "
			t.Log("expected cells:")
			t.Log(expectedStr)

			assert.Equal(t, expectedStr, actualStr)
		})
	}
}

func TestGetNextPipes(t *testing.T) {
	tests := []struct {
		prevPipes []*Pipe
		commit    *models.Commit
		expected  []*Pipe
	}{
		{
			prevPipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Sha:     "b",
				Parents: []string{"c"},
			},
			expected: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: TERMINATES, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "c", kind: STARTS, style: style.FgDefault},
			},
		},
		{
			prevPipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "a", toSha: "b", kind: TERMINATES, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "c", kind: STARTS, style: style.FgDefault},
				{fromPos: 0, toPos: 1, fromSha: "b", toSha: "d", kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Sha:     "d",
				Parents: []string{"e"},
			},
			expected: []*Pipe{
				{fromPos: 0, toPos: 0, fromSha: "b", toSha: "c", kind: CONTINUES, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromSha: "b", toSha: "d", kind: TERMINATES, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromSha: "d", toSha: "e", kind: STARTS, style: style.FgDefault},
			},
		},
	}

	for _, test := range tests {
		getStyle := func(c *models.Commit) style.TextStyle { return style.FgDefault }
		pipes := getNextPipes(test.prevPipes, test.commit, getStyle)
		// rendering cells so that it's easier to see what went wrong
		actualStr := renderPipeSet(pipes, "selected", nil)
		expectedStr := renderPipeSet(test.expected, "selected", nil)
		t.Log("expected cells:")
		t.Log(expectedStr)
		t.Log("actual cells:")
		t.Log(actualStr)
		assert.EqualValues(t, test.expected, pipes)
	}
}

func BenchmarkRenderCommitGraph(b *testing.B) {
	commits := generateCommits(50)
	getStyle := func(commit *models.Commit) style.TextStyle {
		return authors.AuthorStyle(commit.AuthorName)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RenderCommitGraph(commits, "selected", getStyle)
	}
}

func generateCommits(count int) []*models.Commit {
	rand.Seed(1234)
	pool := []*models.Commit{{Sha: "a", AuthorName: "A"}}
	commits := make([]*models.Commit, 0, count)
	authorPool := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	for len(commits) < count {
		currentCommitIdx := rand.Intn(len(pool))
		currentCommit := pool[currentCommitIdx]
		pool = append(pool[0:currentCommitIdx], pool[currentCommitIdx+1:]...)
		// I need to pick a random number of parents to add
		parentCount := rand.Intn(2) + 1

		for j := 0; j < parentCount; j++ {
			reuseParent := rand.Intn(6) != 1 && j <= len(pool)-1 && j != 0
			var newParent *models.Commit
			if reuseParent {
				newParent = pool[j]
			} else {
				newParent = &models.Commit{
					Sha:        fmt.Sprintf("%s%d", currentCommit.Sha, j),
					AuthorName: authorPool[rand.Intn(len(authorPool))],
				}
				pool = append(pool, newParent)
			}
			currentCommit.Parents = append(currentCommit.Parents, newParent.Sha)
		}

		commits = append(commits, currentCommit)
	}

	return commits
}
