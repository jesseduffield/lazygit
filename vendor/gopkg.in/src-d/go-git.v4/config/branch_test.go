package config

import (
	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type BranchSuite struct{}

var _ = Suite(&BranchSuite{})

func (b *BranchSuite) TestValidateName(c *C) {
	goodBranch := Branch{
		Name:   "master",
		Remote: "some_remote",
		Merge:  "refs/heads/master",
	}
	badBranch := Branch{
		Remote: "some_remote",
		Merge:  "refs/heads/master",
	}
	c.Assert(goodBranch.Validate(), IsNil)
	c.Assert(badBranch.Validate(), NotNil)
}

func (b *BranchSuite) TestValidateMerge(c *C) {
	goodBranch := Branch{
		Name:   "master",
		Remote: "some_remote",
		Merge:  "refs/heads/master",
	}
	badBranch := Branch{
		Name:   "master",
		Remote: "some_remote",
		Merge:  "blah",
	}
	c.Assert(goodBranch.Validate(), IsNil)
	c.Assert(badBranch.Validate(), NotNil)
}

func (b *BranchSuite) TestMarshall(c *C) {
	expected := []byte(`[core]
	bare = false
[branch "branch-tracking-on-clone"]
	remote = fork
	merge = refs/heads/branch-tracking-on-clone
`)

	cfg := NewConfig()
	cfg.Branches["branch-tracking-on-clone"] = &Branch{
		Name:   "branch-tracking-on-clone",
		Remote: "fork",
		Merge:  plumbing.ReferenceName("refs/heads/branch-tracking-on-clone"),
	}

	actual, err := cfg.Marshal()
	c.Assert(err, IsNil)
	c.Assert(string(actual), Equals, string(expected))
}

func (b *BranchSuite) TestUnmarshall(c *C) {
	input := []byte(`[core]
	bare = false
[branch "branch-tracking-on-clone"]
	remote = fork
	merge = refs/heads/branch-tracking-on-clone
`)

	cfg := NewConfig()
	err := cfg.Unmarshal(input)
	c.Assert(err, IsNil)
	branch := cfg.Branches["branch-tracking-on-clone"]
	c.Assert(branch.Name, Equals, "branch-tracking-on-clone")
	c.Assert(branch.Remote, Equals, "fork")
	c.Assert(branch.Merge, Equals, plumbing.ReferenceName("refs/heads/branch-tracking-on-clone"))
}
