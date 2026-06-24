package patch

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// newTestPatchBuilder returns a PatchBuilder whose files all resolve to the given
// diff, started for a dummy commit.
func newTestPatchBuilder(diff string) *PatchBuilder {
	pb := NewPatchBuilder(logrus.New().WithField("test", "test"),
		func(from, to string, reverse bool, filename string, plain bool) (string, error) {
			return diff, nil
		},
		nil)
	pb.Start("from", "to", false, true)
	return pb
}

// In simpleDiff the deletion "-orange" is patch line index 6 (old-file line 2) and the
// addition "+grape" is index 7 (new-file line 2).
func TestPatchLineIndicesForLines(t *testing.T) {
	pb := newTestPatchBuilder(simpleDiff)

	indices, err := pb.PatchLineIndicesForLines("filename", []LineIdentity{
		{LineNumber: 2, IsDeletion: true},  // -orange
		{LineNumber: 2, IsDeletion: false}, // +grape
		{LineNumber: 1, IsDeletion: false}, // " apple" — a context line, no change index
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{6, 7}, indices, "context-line identity is skipped; change lines map to their indices")
}

func TestIncludedLineIdentities(t *testing.T) {
	pb := newTestPatchBuilder(simpleDiff)

	// Nothing included yet.
	assert.Empty(t, pb.IncludedLineIdentities("filename"))

	// Include only the deletion: it comes back as its identity, the addition does not.
	assert.NoError(t, pb.AddFileLineRange("filename", []int{6}))
	assert.Equal(t,
		[]LineIdentity{{LineNumber: 2, IsDeletion: true}},
		pb.IncludedLineIdentities("filename"))

	// Including the addition too yields both identities (order-independent).
	assert.NoError(t, pb.AddFileLineRange("filename", []int{7}))
	assert.ElementsMatch(t,
		[]LineIdentity{{LineNumber: 2, IsDeletion: true}, {LineNumber: 2, IsDeletion: false}},
		pb.IncludedLineIdentities("filename"))
}
