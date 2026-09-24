package presentation

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestGetRemoteBranchListDisplayStrings_NilBranch(t *testing.T) {
	branch := &models.RemoteBranch{Name: "main", RemoteName: "origin"}
	branches := []*models.RemoteBranch{nil, branch, nil}

	result := GetRemoteBranchListDisplayStrings(branches, "")

	// nil entries must be skipped; only the valid branch produces a row
	assert.Len(t, result, 1)
}
