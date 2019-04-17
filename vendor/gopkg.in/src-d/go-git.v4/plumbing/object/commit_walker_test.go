package object

import (
	"gopkg.in/src-d/go-git.v4/plumbing"

	. "gopkg.in/check.v1"
)

type CommitWalkerSuite struct {
	BaseObjectsSuite
}

var _ = Suite(&CommitWalkerSuite{})

func (s *CommitWalkerSuite) TestCommitPreIterator(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitPreorderIter(commit, nil, nil).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 8)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
		"af2d6a6954d532f8ffb47615169c8fdf9d383a1a",
		"1669dce138d9b841a518c64b10914d88f5e488ea",
		"35e85108805c84807bc66a02d91535e1e24b38b9",
		"b029517f6300c2da0f4b651b8642506cd6aaf45d",
		"a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69",
		"b8e471f58bcbca63b07bda20e428190409c2db47",
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitPreIteratorWithIgnore(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitPreorderIter(commit, nil, []plumbing.Hash{
		plumbing.NewHash("af2d6a6954d532f8ffb47615169c8fdf9d383a1a"),
	}).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 2)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitPreIteratorWithSeenExternal(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	seenExternal := map[plumbing.Hash]bool{
		plumbing.NewHash("af2d6a6954d532f8ffb47615169c8fdf9d383a1a"): true,
	}
	NewCommitPreorderIter(commit, seenExternal, nil).
		ForEach(func(c *Commit) error {
			commits = append(commits, c)
			return nil
		})

	c.Assert(commits, HasLen, 2)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitPostIterator(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitPostorderIter(commit, nil).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 8)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
		"af2d6a6954d532f8ffb47615169c8fdf9d383a1a",
		"1669dce138d9b841a518c64b10914d88f5e488ea",
		"a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69",
		"b8e471f58bcbca63b07bda20e428190409c2db47",
		"b029517f6300c2da0f4b651b8642506cd6aaf45d",
		"35e85108805c84807bc66a02d91535e1e24b38b9",
	}

	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitPostIteratorWithIgnore(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitPostorderIter(commit, []plumbing.Hash{
		plumbing.NewHash("af2d6a6954d532f8ffb47615169c8fdf9d383a1a"),
	}).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 2)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitCTimeIterator(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitIterCTime(commit, nil, nil).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 8)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5", // 2015-04-05T23:30:47+02:00
		"918c48b83bd081e863dbe1b80f8998f058cd8294", // 2015-03-31T13:56:18+02:00
		"af2d6a6954d532f8ffb47615169c8fdf9d383a1a", // 2015-03-31T13:51:51+02:00
		"1669dce138d9b841a518c64b10914d88f5e488ea", // 2015-03-31T13:48:14+02:00
		"a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69", // 2015-03-31T13:47:14+02:00
		"35e85108805c84807bc66a02d91535e1e24b38b9", // 2015-03-31T13:46:24+02:00
		"b8e471f58bcbca63b07bda20e428190409c2db47", // 2015-03-31T13:44:52+02:00
		"b029517f6300c2da0f4b651b8642506cd6aaf45d", // 2015-03-31T13:42:21+02:00
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitCTimeIteratorWithIgnore(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitIterCTime(commit, nil, []plumbing.Hash{
		plumbing.NewHash("af2d6a6954d532f8ffb47615169c8fdf9d383a1a"),
	}).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 2)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitBSFIterator(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitIterBSF(commit, nil, nil).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 8)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
		"af2d6a6954d532f8ffb47615169c8fdf9d383a1a",
		"1669dce138d9b841a518c64b10914d88f5e488ea",
		"35e85108805c84807bc66a02d91535e1e24b38b9",
		"a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69",
		"b029517f6300c2da0f4b651b8642506cd6aaf45d",
		"b8e471f58bcbca63b07bda20e428190409c2db47",
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}

func (s *CommitWalkerSuite) TestCommitBSFIteratorWithIgnore(c *C) {
	commit := s.commit(c, s.Fixture.Head)

	var commits []*Commit
	NewCommitIterBSF(commit, nil, []plumbing.Hash{
		plumbing.NewHash("af2d6a6954d532f8ffb47615169c8fdf9d383a1a"),
	}).ForEach(func(c *Commit) error {
		commits = append(commits, c)
		return nil
	})

	c.Assert(commits, HasLen, 2)

	expected := []string{
		"6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"918c48b83bd081e863dbe1b80f8998f058cd8294",
	}
	for i, commit := range commits {
		c.Assert(commit.Hash.String(), Equals, expected[i])
	}
}
