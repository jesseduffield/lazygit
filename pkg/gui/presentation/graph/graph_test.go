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

func TestRenderCommitGraph(t *testing.T) {
	tests := []struct {
		name           string
		commits        []*models.Commit
		expectedOutput string
	}{
		{
			name: "with some merges",
			commits: []*models.Commit{
				{Hash: "1", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"3"}},
				{Hash: "3", Parents: []string{"4"}},
				{Hash: "4", Parents: []string{"5", "7"}},
				{Hash: "7", Parents: []string{"5"}},
				{Hash: "5", Parents: []string{"8"}},
				{Hash: "8", Parents: []string{"9"}},
				{Hash: "9", Parents: []string{"A", "B"}},
				{Hash: "B", Parents: []string{"D"}},
				{Hash: "D", Parents: []string{"D"}},
				{Hash: "A", Parents: []string{"E"}},
				{Hash: "E", Parents: []string{"F"}},
				{Hash: "F", Parents: []string{"D"}},
				{Hash: "D", Parents: []string{"G"}},
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
				{Hash: "1", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"3", "4"}},
				{Hash: "4", Parents: []string{"3", "5"}},
				{Hash: "3", Parents: []string{"5"}},
				{Hash: "5", Parents: []string{"6"}},
				{Hash: "6", Parents: []string{"7"}},
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
				{Hash: "1", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"3", "4"}},
				{Hash: "4", Parents: []string{"3", "5"}},
				{Hash: "Z", Parents: []string{"Z"}},
				{Hash: "3", Parents: []string{"5"}},
				{Hash: "5", Parents: []string{"6"}},
				{Hash: "6", Parents: []string{"7"}},
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
				{Hash: "1", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"3", "4"}},
				{Hash: "3", Parents: []string{"5", "4"}},
				{Hash: "5", Parents: []string{"7", "8"}},
				{Hash: "4", Parents: []string{"7"}},
				{Hash: "7", Parents: []string{"11"}},
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
				{Hash: "1", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"3", "4"}},
				{Hash: "3", Parents: []string{"5", "4"}},
				{Hash: "5", Parents: []string{"7", "8"}},
				{Hash: "7", Parents: []string{"4", "A"}},
				{Hash: "4", Parents: []string{"B"}},
				{Hash: "B", Parents: []string{"C"}},
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
				{Hash: "1", Parents: []string{"2", "3"}},
				{Hash: "3", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"4", "5"}},
				{Hash: "4", Parents: []string{"6", "7"}},
				{Hash: "6", Parents: []string{"8"}},
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
				{Hash: "1", Parents: []string{"2", "3", "4", "5"}},
				{Hash: "4", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"A"}},
				{Hash: "A", Parents: []string{"6", "B"}},
				{Hash: "B", Parents: []string{"C"}},
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
				{Hash: "1", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"3", "4"}},
				{Hash: "3", Parents: []string{"5", "4"}},
				{Hash: "5", Parents: []string{"7", "8"}},
				{Hash: "7", Parents: []string{"4", "A"}},
				{Hash: "4", Parents: []string{"B"}},
				{Hash: "B", Parents: []string{"C"}},
				{Hash: "C", Parents: []string{"D"}},
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
				{Hash: "1", Parents: []string{"2"}},
				{Hash: "2", Parents: []string{"3", "4"}},
				{Hash: "3", Parents: []string{"5", "4"}},
				{Hash: "5", Parents: []string{"7", "G"}},
				{Hash: "7", Parents: []string{"8", "A"}},
				{Hash: "8", Parents: []string{"4", "E"}},
				{Hash: "4", Parents: []string{"B"}},
				{Hash: "B", Parents: []string{"C"}},
				{Hash: "C", Parents: []string{"D"}},
				{Hash: "D", Parents: []string{"F"}},
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

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, test := range tests {
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
				description := test.commits[i].Hash
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
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "b", kind: TERMINATES, style: cyan},
				{fromPos: 0, toPos: 0, fromHash: "b", toHash: "c", kind: STARTS, style: green},
			},
			prevCommit:     &models.Commit{Hash: "a"},
			expectedStr:    "◯",
			expectedStyles: []style.TextStyle{green},
		},
		{
			name: "single cell, selected",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "selected", kind: TERMINATES, style: cyan},
				{fromPos: 0, toPos: 0, fromHash: "selected", toHash: "c", kind: STARTS, style: green},
			},
			prevCommit:     &models.Commit{Hash: "a"},
			expectedStr:    "◯",
			expectedStyles: []style.TextStyle{highlightStyle},
		},
		{
			name: "terminating hook and starting hook, selected",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "selected", kind: TERMINATES, style: cyan},
				{fromPos: 1, toPos: 0, fromHash: "c", toHash: "selected", kind: TERMINATES, style: yellow},
				{fromPos: 0, toPos: 0, fromHash: "selected", toHash: "d", kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromHash: "selected", toHash: "e", kind: STARTS, style: green},
			},
			prevCommit:  &models.Commit{Hash: "a"},
			expectedStr: "⏣─╮",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle,
			},
		},
		{
			name: "terminating hook and starting hook, prioritise the terminating one",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "b", kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 0, fromHash: "c", toHash: "b", kind: TERMINATES, style: magenta},
				{fromPos: 0, toPos: 0, fromHash: "b", toHash: "d", kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromHash: "b", toHash: "e", kind: STARTS, style: green},
			},
			prevCommit:  &models.Commit{Hash: "a"},
			expectedStr: "⏣─│",
			expectedStyles: []style.TextStyle{
				green, green, magenta,
			},
		},
		{
			name: "starting and terminating pipe sharing some space",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a1", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "a2", toHash: "a3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: "b1", toHash: "b2", kind: CONTINUES, style: magenta},
				{fromPos: 3, toPos: 0, fromHash: "e1", toHash: "a2", kind: TERMINATES, style: green},
				{fromPos: 0, toPos: 2, fromHash: "a2", toHash: "c3", kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: "a1"},
			expectedStr: "⏣─│─┬─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, magenta, yellow, yellow, green, green,
			},
		},
		{
			name: "starting and terminating pipe sharing some space, with selection",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a1", toHash: "selected", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "selected", toHash: "a3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: "b1", toHash: "b2", kind: CONTINUES, style: magenta},
				{fromPos: 3, toPos: 0, fromHash: "e1", toHash: "selected", kind: TERMINATES, style: green},
				{fromPos: 0, toPos: 2, fromHash: "selected", toHash: "c3", kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: "a1"},
			expectedStr: "⏣───╮ ╯",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle, highlightStyle, highlightStyle, nothing, green,
			},
		},
		{
			name: "many terminating pipes",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a1", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "a2", toHash: "a3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 0, fromHash: "b1", toHash: "a2", kind: TERMINATES, style: magenta},
				{fromPos: 2, toPos: 0, fromHash: "c1", toHash: "a2", kind: TERMINATES, style: green},
			},
			prevCommit:  &models.Commit{Hash: "a1"},
			expectedStr: "◯─┴─╯",
			expectedStyles: []style.TextStyle{
				yellow, magenta, magenta, green, green,
			},
		},
		{
			name: "starting pipe passing through",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a1", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "a2", toHash: "a3", kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 3, fromHash: "a2", toHash: "d3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: "b1", toHash: "b3", kind: CONTINUES, style: magenta},
				{fromPos: 2, toPos: 2, fromHash: "c1", toHash: "c3", kind: CONTINUES, style: green},
			},
			prevCommit:  &models.Commit{Hash: "a1"},
			expectedStr: "⏣─│─│─╮",
			expectedStyles: []style.TextStyle{
				yellow, yellow, magenta, yellow, green, yellow, yellow,
			},
		},
		{
			name: "starting and terminating path crossing continuing path",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a1", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "a2", toHash: "a3", kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 1, fromHash: "a2", toHash: "b3", kind: STARTS, style: yellow},
				{fromPos: 1, toPos: 1, fromHash: "b1", toHash: "a2", kind: CONTINUES, style: green},
				{fromPos: 2, toPos: 0, fromHash: "c1", toHash: "a2", kind: TERMINATES, style: magenta},
			},
			prevCommit:  &models.Commit{Hash: "a1"},
			expectedStr: "⏣─│─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, green, magenta, magenta,
			},
		},
		{
			name: "another clash of starting and terminating paths",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a1", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "a2", toHash: "a3", kind: STARTS, style: yellow},
				{fromPos: 0, toPos: 1, fromHash: "a2", toHash: "b3", kind: STARTS, style: yellow},
				{fromPos: 2, toPos: 2, fromHash: "c1", toHash: "c3", kind: CONTINUES, style: green},
				{fromPos: 3, toPos: 0, fromHash: "d1", toHash: "a2", kind: TERMINATES, style: magenta},
			},
			prevCommit:  &models.Commit{Hash: "a1"},
			expectedStr: "⏣─┬─│─╯",
			expectedStyles: []style.TextStyle{
				yellow, yellow, yellow, magenta, green, magenta, magenta,
			},
		},
		{
			name: "commit whose previous commit is selected",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "selected", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "a2", toHash: "a3", kind: STARTS, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: "selected"},
			expectedStr: "◯",
			expectedStyles: []style.TextStyle{
				yellow,
			},
		},
		{
			name: "commit whose previous commit is selected and is a merge commit",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "selected", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 1, fromHash: "selected", toHash: "b3", kind: CONTINUES, style: red},
			},
			prevCommit:  &models.Commit{Hash: "selected"},
			expectedStr: "◯ │",
			expectedStyles: []style.TextStyle{
				highlightStyle, nothing, highlightStyle,
			},
		},
		{
			name: "commit whose previous commit is selected and is a merge commit, with continuing pipe inbetween",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "selected", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 1, toPos: 1, fromHash: "z1", toHash: "z3", kind: CONTINUES, style: green},
				{fromPos: 2, toPos: 2, fromHash: "selected", toHash: "b3", kind: CONTINUES, style: red},
			},
			prevCommit:  &models.Commit{Hash: "selected"},
			expectedStr: "◯ │ │",
			expectedStyles: []style.TextStyle{
				highlightStyle, nothing, green, nothing, highlightStyle,
			},
		},
		{
			name: "when previous commit is selected, not a merge commit, and spawns a continuing pipe",
			pipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a1", toHash: "a2", kind: TERMINATES, style: red},
				{fromPos: 0, toPos: 0, fromHash: "a2", toHash: "a3", kind: STARTS, style: green},
				{fromPos: 0, toPos: 1, fromHash: "a2", toHash: "b3", kind: STARTS, style: green},
				{fromPos: 1, toPos: 0, fromHash: "selected", toHash: "a2", kind: TERMINATES, style: yellow},
			},
			prevCommit:  &models.Commit{Hash: "selected"},
			expectedStr: "⏣─╯",
			expectedStyles: []style.TextStyle{
				highlightStyle, highlightStyle, highlightStyle,
			},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

	for _, test := range tests {
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
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "b", kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Hash:    "b",
				Parents: []string{"c"},
			},
			expected: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "b", kind: TERMINATES, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromHash: "b", toHash: "c", kind: STARTS, style: style.FgDefault},
			},
		},
		{
			prevPipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "b", kind: TERMINATES, style: style.FgDefault},
				{fromPos: 0, toPos: 0, fromHash: "b", toHash: "c", kind: STARTS, style: style.FgDefault},
				{fromPos: 0, toPos: 1, fromHash: "b", toHash: "d", kind: STARTS, style: style.FgDefault},
			},
			commit: &models.Commit{
				Hash:    "d",
				Parents: []string{"e"},
			},
			expected: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "b", toHash: "c", kind: CONTINUES, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromHash: "b", toHash: "d", kind: TERMINATES, style: style.FgDefault},
				{fromPos: 1, toPos: 1, fromHash: "d", toHash: "e", kind: STARTS, style: style.FgDefault},
			},
		},
		{
			prevPipes: []*Pipe{
				{fromPos: 0, toPos: 0, fromHash: "a", toHash: "root", kind: TERMINATES, style: style.FgDefault},
			},
			commit: &models.Commit{
				Hash:    "root",
				Parents: []string{},
			},
			expected: []*Pipe{
				{fromPos: 1, toPos: 1, fromHash: "root", toHash: models.EmptyTreeCommitHash, kind: STARTS, style: style.FgDefault},
			},
		},
	}

	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

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
	oldColorLevel := color.ForceSetColorLevel(terminfo.ColorLevelMillions)
	defer color.ForceSetColorLevel(oldColorLevel)

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
	rnd := rand.New(rand.NewSource(1234))
	pool := []*models.Commit{{Hash: "a", AuthorName: "A"}}
	commits := make([]*models.Commit, 0, count)
	authorPool := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	for len(commits) < count {
		currentCommitIdx := rnd.Intn(len(pool))
		currentCommit := pool[currentCommitIdx]
		pool = append(pool[0:currentCommitIdx], pool[currentCommitIdx+1:]...)
		// I need to pick a random number of parents to add
		parentCount := rnd.Intn(2) + 1

		for j := 0; j < parentCount; j++ {
			reuseParent := rnd.Intn(6) != 1 && j <= len(pool)-1 && j != 0
			var newParent *models.Commit
			if reuseParent {
				newParent = pool[j]
			} else {
				newParent = &models.Commit{
					Hash:       fmt.Sprintf("%s%d", currentCommit.Hash, j),
					AuthorName: authorPool[rnd.Intn(len(authorPool))],
				}
				pool = append(pool, newParent)
			}
			currentCommit.Parents = append(currentCommit.Parents, newParent.Hash)
		}

		commits = append(commits, currentCommit)
	}

	return commits
}
