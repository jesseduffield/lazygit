package git

import (
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type PruneSuite struct {
	BaseSuite
}

var _ = Suite(&PruneSuite{})

func (s *PruneSuite) testPrune(c *C, deleteTime time.Time) {
	srcFs := fixtures.ByTag("unpacked").One().DotGit()
	var sto storage.Storer
	var err error
	sto, err = filesystem.NewStorage(srcFs)
	c.Assert(err, IsNil)

	los := sto.(storer.LooseObjectStorer)
	c.Assert(los, NotNil)

	count := 0
	err = los.ForEachObjectHash(func(_ plumbing.Hash) error {
		count++
		return nil
	})
	c.Assert(err, IsNil)

	r, err := Open(sto, srcFs)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	// Remove a branch so we can prune some objects.
	err = sto.RemoveReference(plumbing.ReferenceName("refs/heads/v4"))
	c.Assert(err, IsNil)
	err = sto.RemoveReference(plumbing.ReferenceName("refs/remotes/origin/v4"))
	c.Assert(err, IsNil)

	err = r.Prune(PruneOptions{
		OnlyObjectsOlderThan: deleteTime,
		Handler:              r.DeleteObject,
	})
	c.Assert(err, IsNil)

	newCount := 0
	err = los.ForEachObjectHash(func(_ plumbing.Hash) error {
		newCount++
		return nil
	})
	if deleteTime.IsZero() {
		c.Assert(newCount < count, Equals, true)
	} else {
		// Assume a delete time older than any of the objects was passed in.
		c.Assert(newCount, Equals, count)
	}
}

func (s *PruneSuite) TestPrune(c *C) {
	s.testPrune(c, time.Time{})
}

func (s *PruneSuite) TestPruneWithNoDelete(c *C) {
	s.testPrune(c, time.Unix(0, 1))
}
