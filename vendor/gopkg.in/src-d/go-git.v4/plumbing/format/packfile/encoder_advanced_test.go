package packfile_test

import (
	"bytes"
	"math/rand"
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing"
	. "gopkg.in/src-d/go-git.v4/plumbing/format/packfile"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type EncoderAdvancedSuite struct {
	fixtures.Suite
}

var _ = Suite(&EncoderAdvancedSuite{})

func (s *EncoderAdvancedSuite) TestEncodeDecode(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	fixs := fixtures.Basic().ByTag("packfile").ByTag(".git")
	fixs = append(fixs, fixtures.ByURL("https://github.com/src-d/go-git.git").
		ByTag("packfile").ByTag(".git").One())
	fixs.Test(c, func(f *fixtures.Fixture) {
		storage, err := filesystem.NewStorage(f.DotGit())
		c.Assert(err, IsNil)
		s.testEncodeDecode(c, storage, 10)
	})

}

func (s *EncoderAdvancedSuite) TestEncodeDecodeNoDeltaCompression(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	fixs := fixtures.Basic().ByTag("packfile").ByTag(".git")
	fixs = append(fixs, fixtures.ByURL("https://github.com/src-d/go-git.git").
		ByTag("packfile").ByTag(".git").One())
	fixs.Test(c, func(f *fixtures.Fixture) {
		storage, err := filesystem.NewStorage(f.DotGit())
		c.Assert(err, IsNil)
		s.testEncodeDecode(c, storage, 0)
	})
}

func (s *EncoderAdvancedSuite) testEncodeDecode(c *C, storage storer.Storer, packWindow uint) {

	objIter, err := storage.IterEncodedObjects(plumbing.AnyObject)
	c.Assert(err, IsNil)

	expectedObjects := map[plumbing.Hash]bool{}
	var hashes []plumbing.Hash
	err = objIter.ForEach(func(o plumbing.EncodedObject) error {
		expectedObjects[o.Hash()] = true
		hashes = append(hashes, o.Hash())
		return err

	})
	c.Assert(err, IsNil)

	// Shuffle hashes to avoid delta selector getting order right just because
	// the initial order is correct.
	auxHashes := make([]plumbing.Hash, len(hashes))
	for i, j := range rand.Perm(len(hashes)) {
		auxHashes[j] = hashes[i]
	}
	hashes = auxHashes

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, storage, false)
	encodeHash, err := enc.Encode(hashes, packWindow)
	c.Assert(err, IsNil)

	scanner := NewScanner(buf)
	storage = memory.NewStorage()
	d, err := NewDecoder(scanner, storage)
	c.Assert(err, IsNil)
	decodeHash, err := d.Decode()
	c.Assert(err, IsNil)

	c.Assert(encodeHash, Equals, decodeHash)

	objIter, err = storage.IterEncodedObjects(plumbing.AnyObject)
	c.Assert(err, IsNil)
	obtainedObjects := map[plumbing.Hash]bool{}
	err = objIter.ForEach(func(o plumbing.EncodedObject) error {
		obtainedObjects[o.Hash()] = true
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(obtainedObjects, DeepEquals, expectedObjects)

	for h := range obtainedObjects {
		if !expectedObjects[h] {
			c.Errorf("obtained unexpected object: %s", h)
		}
	}

	for h := range expectedObjects {
		if !obtainedObjects[h] {
			c.Errorf("missing object: %s", h)
		}
	}
}
