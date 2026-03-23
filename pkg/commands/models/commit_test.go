package models

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCommitParentRefNameRootUsesEmptyTreeParent(t *testing.T) {
	pool := &utils.StringPool{}
	oid := pool.Add("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	c := NewCommit(pool, NewCommitOpts{Hash: "abc", Parents: nil, EmptyTreeParent: oid})
	assert.Equal(t, *oid, c.ParentRefName())
}

func TestCommitParentRefNameRootFallsBackToSHA1EmptyTree(t *testing.T) {
	pool := &utils.StringPool{}
	c := NewCommit(pool, NewCommitOpts{Hash: "abc", Parents: nil})
	assert.Equal(t, EmptyTreeCommitHash, c.ParentRefName())
}
