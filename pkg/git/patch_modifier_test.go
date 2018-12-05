package git

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func newDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
}

func newDummyPatchModifier() *PatchModifier {
	return &PatchModifier{
		Log: newDummyLog(),
	}
}
func TestModifyPatchForLine(t *testing.T) {
	type scenario struct {
		testName              string
		patchFilename         string
		lineNumber            int
		shouldError           bool
		expectedPatchFilename string
	}

	scenarios := []scenario{
		{
			"Removing one line",
			"testdata/testPatchBefore.diff",
			8,
			false,
			"testdata/testPatchAfter1.diff",
		},
		{
			"Adding one line",
			"testdata/testPatchBefore.diff",
			10,
			false,
			"testdata/testPatchAfter2.diff",
		},
		{
			"Adding one line in top hunk in diff with multiple hunks",
			"testdata/testPatchBefore2.diff",
			20,
			false,
			"testdata/testPatchAfter3.diff",
		},
		{
			"Adding one line in top hunk in diff with multiple hunks",
			"testdata/testPatchBefore2.diff",
			53,
			false,
			"testdata/testPatchAfter4.diff",
		},
		{
			"adding unstaged file with a single line",
			"testdata/addedFile.diff",
			6,
			false,
			"testdata/addedFile.diff",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			p := newDummyPatchModifier()
			beforePatch, err := ioutil.ReadFile(s.patchFilename)
			if err != nil {
				panic("Cannot open file at " + s.patchFilename)
			}
			afterPatch, err := p.ModifyPatchForLine(string(beforePatch), s.lineNumber)
			if s.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				expected, err := ioutil.ReadFile(s.expectedPatchFilename)
				if err != nil {
					panic("Cannot open file at " + s.expectedPatchFilename)
				}
				assert.Equal(t, string(expected), afterPatch)
			}
		})
	}
}
