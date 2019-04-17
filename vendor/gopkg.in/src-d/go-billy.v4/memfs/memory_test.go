package memfs

import (
	"testing"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/test"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MemorySuite struct {
	test.FilesystemSuite
	path string
}

var _ = Suite(&MemorySuite{})

func (s *MemorySuite) SetUpTest(c *C) {
	s.FilesystemSuite = test.NewFilesystemSuite(New())
}

func (s *MemorySuite) TestCapabilities(c *C) {
	_, ok := s.FS.(billy.Capable)
	c.Assert(ok, Equals, true)

	caps := billy.Capabilities(s.FS)
	c.Assert(caps, Equals, billy.DefaultCapabilities&^billy.LockCapability)
}
