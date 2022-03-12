package loaders

// "*|feat/detect-purge|origin/feat/detect-purge|[ahead 1]"
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObtainBranchTrimHeads(t *testing.T) {
	// Given
	split := []string{"", "heads/a_branch", "", ""}

	// When
	branch := obtainBranch(split)

	// Then
	assert.EqualValues(t, "a_branch", branch.Name)
}

func TestObtainBranchNoUpstream(t *testing.T) {
	// Given
	split := []string{"", "a_branch", "", ""}

	// When
	branch := obtainBranch(split)

	// Then
	// We get the default values for pullables and pushables i.e. "?"
	assert.EqualValues(t, "a_branch", branch.Name)
	assert.EqualValues(t, false, branch.Head)
	assert.EqualValues(t, "?", branch.Pushables)
	assert.EqualValues(t, "?", branch.Pullables)
}

func TestObtainBranchIsHead(t *testing.T) {
	// Given
	split := []string{"*", "", "", ""}

	// When
	branch := obtainBranch(split)

	// Then
	assert.EqualValues(t, true, branch.Head)
}

func TestObtainBranchIsBehindAndAhead(t *testing.T) {
	// Given
	split := []string{"", "a_branch", "a_remote/a_branch", "[behind 2, ahead 3]"}

	// When
	branch := obtainBranch(split)

	// Then
	assert.EqualValues(t, "a_branch", branch.Name)
	assert.EqualValues(t, false, branch.Head)
	assert.EqualValues(t, "2", branch.Pullables)
	assert.EqualValues(t, "3", branch.Pushables)
}

func TestObtainBranchDeletedInRemote(t *testing.T) {
	// Given
	split := []string{"", "a_branch", "a_remote/a_branch", "[gone]"}

	// When
	branch := obtainBranch(split)

	// Then
	assert.EqualValues(t, "a_branch", branch.Name)
	assert.EqualValues(t, false, branch.Head)
	assert.EqualValues(t, "d", branch.Pullables)
	assert.EqualValues(t, "?", branch.Pushables)
}
