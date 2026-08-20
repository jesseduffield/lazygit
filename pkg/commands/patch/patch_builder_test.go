package patch

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// newTestPatchBuilder returns a patch builder started for a dummy commit, in which
// every file's diff is the given one.
func newTestPatchBuilder(diff string) *PatchBuilder {
	patchBuilder := NewPatchBuilder(logrus.New().WithField("test", "test"),
		func(from string, to string, reverse bool, filename string, previousPath string) (string, error) {
			return diff, nil
		},
		// Nothing here renders the patch, so it needs no directory to be
		// materialized into.
		nil)
	patchBuilder.Start("from", "to", false, true)
	return patchBuilder
}

// In simpleDiff the deletion "-orange" is line index 6 of the parsed diff (line 2 of
// the old file) and the addition "+grape" is index 7 (line 2 of the new file).
func TestPatchLineIndicesForLines(t *testing.T) {
	patchBuilder := newTestPatchBuilder(simpleDiff)

	indices, everyChange, err := patchBuilder.PatchLineIndicesForLines("filename", "", []LineIdentity{
		{LineNumber: 2, IsDeletion: true},  // -orange
		{LineNumber: 2, IsDeletion: false}, // +grape
		{LineNumber: 1, IsDeletion: false}, // " apple", a context line
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{6, 7}, indices, "the context line names no change line")
	assert.True(t, everyChange, "the two changes are all the diff has")
}

// A renamed file's rename header makes its change lines sit further down the diff, and
// its old-file line numbers are of the file under its previous name.
func TestPatchLineIndicesForLinesOfARenamedFile(t *testing.T) {
	patchBuilder := newTestPatchBuilder(renameWithModificationDiff)

	indices, everyChange, err := patchBuilder.PatchLineIndicesForLines("newname", "oldname", []LineIdentity{
		{LineNumber: 2, IsDeletion: true},  // -orange
		{LineNumber: 2, IsDeletion: false}, // +grape
	})
	assert.NoError(t, err)
	assert.Equal(t, []int{9, 10}, indices)
	assert.True(t, everyChange)
}

// everyChange is about the diff alone: whether anything the file changes was left out
// of the selection. What that then means for the patch is the caller's question.
func TestPatchLineIndicesForLinesEveryChange(t *testing.T) {
	patchBuilder := newTestPatchBuilder(newFile)

	_, everyChange, err := patchBuilder.PatchLineIndicesForLines("newfile", "", []LineIdentity{
		{LineNumber: 1},
		{LineNumber: 2},
	})
	assert.NoError(t, err)
	assert.False(t, everyChange, "the file's third added line is left out")

	_, everyChange, err = patchBuilder.PatchLineIndicesForLines("newfile", "", []LineIdentity{
		{LineNumber: 1},
		{LineNumber: 2},
		{LineNumber: 3},
	})
	assert.NoError(t, err)
	assert.True(t, everyChange)

	// A line the diff doesn't have doesn't stand in for one it does.
	patchBuilder = newTestPatchBuilder(deletedFile)
	_, everyChange, err = patchBuilder.PatchLineIndicesForLines("newfile", "", []LineIdentity{
		{LineNumber: 1, IsDeletion: true},
		{LineNumber: 2, IsDeletion: true},
		{LineNumber: 4, IsDeletion: true},
	})
	assert.NoError(t, err)
	assert.False(t, everyChange)
}

func TestIncludedLineIdentities(t *testing.T) {
	patchBuilder := newTestPatchBuilder(simpleDiff)

	// A file no part of the patch has nothing included.
	assert.Empty(t, patchBuilder.IncludedLineIdentities("filename"))

	// With only the deletion in, only its identity comes back.
	assert.NoError(t, patchBuilder.AddFileLineRange("filename", "", []int{6}))
	assert.Equal(t,
		[]LineIdentity{{LineNumber: 2, IsDeletion: true}},
		patchBuilder.IncludedLineIdentities("filename"))

	// With the addition in as well, both do.
	assert.NoError(t, patchBuilder.AddFileLineRange("filename", "", []int{7}))
	assert.ElementsMatch(t,
		[]LineIdentity{{LineNumber: 2, IsDeletion: true}, {LineNumber: 2, IsDeletion: false}},
		patchBuilder.IncludedLineIdentities("filename"))
}

func TestFilesInPatch(t *testing.T) {
	patchBuilder := newTestPatchBuilder(simpleDiff)

	// A file no part of the patch is no part of what the patch is materialized from.
	assert.Empty(t, patchBuilder.FilesInPatch())

	assert.NoError(t, patchBuilder.AddFileLineRange("filename", "", []int{6}))
	assert.Equal(t,
		[]PatchFile{{Path: "filename", ContentPath: "filename"}},
		patchBuilder.FilesInPatch())
}

// A renamed file's content is under the name it had before whatever the patch calls the
// file, and the patch calls it by the name it had before only where it carries the
// rename — a partial selection has the rename stripped and names the file by the new one.
func TestFilesInPatchOfARenamedFile(t *testing.T) {
	patchBuilder := newTestPatchBuilder(renameWithModificationDiff)

	assert.NoError(t, patchBuilder.AddFileLineRange("newname", "oldname", []int{9}))
	assert.Equal(t,
		[]PatchFile{{Path: "newname", ContentPath: "oldname"}},
		patchBuilder.FilesInPatch())

	assert.NoError(t, patchBuilder.AddFileWhole("newname", "oldname"))
	assert.Equal(t,
		[]PatchFile{{Path: "oldname", ContentPath: "oldname"}},
		patchBuilder.FilesInPatch())
}

// A file taken into the patch whole has every one of its change lines in it.
func TestIncludedLineIdentitiesOfAWholeFile(t *testing.T) {
	patchBuilder := newTestPatchBuilder(simpleDiff)

	assert.NoError(t, patchBuilder.AddFileWhole("filename", ""))
	assert.ElementsMatch(t,
		[]LineIdentity{{LineNumber: 2, IsDeletion: true}, {LineNumber: 2, IsDeletion: false}},
		patchBuilder.IncludedLineIdentities("filename"))
}
