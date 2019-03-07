package git

import (
	"io/ioutil"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/stretchr/testify/assert"
)

// NewDummyPatchParser constructs a new dummy patch parser for testing
func NewDummyPatchParser() *PatchParser {
	return &PatchParser{
		Log: commands.NewDummyLog(),
	}
}

func TestParsePatch(t *testing.T) {
	type scenario struct {
		testName               string
		patchFilename          string
		shouldError            bool
		expectedStageableLines []int
		expectedHunkStarts     []int
	}

	scenarios := []scenario{
		{
			"Diff with one hunk",
			"testdata/testPatchBefore.diff",
			false,
			[]int{8, 9, 10, 11},
			[]int{4},
		},
		{
			"Diff with two hunks",
			"testdata/testPatchBefore2.diff",
			false,
			[]int{8, 9, 10, 11, 12, 13, 20, 21, 22, 23, 24, 25, 26, 27, 28, 33, 34, 35, 36, 37, 45, 46, 47, 48, 49, 50, 51, 52, 53},
			[]int{4, 41},
		},
		{
			"Unstaged file",
			"testdata/addedFile.diff",
			false,
			[]int{6},
			[]int{5},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			p := NewDummyPatchParser()
			beforePatch, err := ioutil.ReadFile(s.patchFilename)
			if err != nil {
				panic("Cannot open file at " + s.patchFilename)
			}
			hunkStarts, stageableLines, err := p.ParsePatch(string(beforePatch))
			if s.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, s.expectedStageableLines, stageableLines)
				assert.Equal(t, s.expectedHunkStarts, hunkStarts)
			}
		})
	}
}
