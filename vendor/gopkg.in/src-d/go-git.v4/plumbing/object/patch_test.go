package object

import (
	. "gopkg.in/check.v1"
	fixtures "gopkg.in/src-d/go-git-fixtures.v3"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type PatchSuite struct {
	BaseObjectsSuite
}

var _ = Suite(&PatchSuite{})

func (s *PatchSuite) TestStatsWithSubmodules(c *C) {
	storer, err := filesystem.NewStorage(
		fixtures.ByURL("https://github.com/git-fixtures/submodule.git").One().DotGit())

	commit, err := GetCommit(storer, plumbing.NewHash("b685400c1f9316f350965a5993d350bc746b0bf4"))

	tree, err := commit.Tree()
	c.Assert(err, IsNil)

	e, err := tree.entry("basic")
	c.Assert(err, IsNil)

	ch := &Change{
		From: ChangeEntry{
			Name:      "basic",
			Tree:      tree,
			TreeEntry: *e,
		},
		To: ChangeEntry{
			Name:      "basic",
			Tree:      tree,
			TreeEntry: *e,
		},
	}

	p, err := getPatch("", ch)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
}
