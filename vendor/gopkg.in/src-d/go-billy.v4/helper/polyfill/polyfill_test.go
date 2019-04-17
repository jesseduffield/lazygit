package polyfill

import (
	"path/filepath"
	"testing"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/test"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&PolyfillSuite{})

type PolyfillSuite struct {
	Helper     billy.Filesystem
	Underlying billy.Filesystem
}

func (s *PolyfillSuite) SetUpTest(c *C) {
	s.Helper = New(&test.BasicMock{})
}

func (s *PolyfillSuite) TestTempFile(c *C) {
	_, err := s.Helper.TempFile("", "")
	c.Assert(err, Equals, billy.ErrNotSupported)
}

func (s *PolyfillSuite) TestReadDir(c *C) {
	_, err := s.Helper.ReadDir("")
	c.Assert(err, Equals, billy.ErrNotSupported)
}

func (s *PolyfillSuite) TestMkdirAll(c *C) {
	err := s.Helper.MkdirAll("", 0)
	c.Assert(err, Equals, billy.ErrNotSupported)
}

func (s *PolyfillSuite) TestSymlink(c *C) {
	err := s.Helper.Symlink("", "")
	c.Assert(err, Equals, billy.ErrNotSupported)
}

func (s *PolyfillSuite) TestReadlink(c *C) {
	_, err := s.Helper.Readlink("")
	c.Assert(err, Equals, billy.ErrNotSupported)
}

func (s *PolyfillSuite) TestLstat(c *C) {
	_, err := s.Helper.Lstat("")
	c.Assert(err, Equals, billy.ErrNotSupported)
}

func (s *PolyfillSuite) TestChroot(c *C) {
	_, err := s.Helper.Chroot("")
	c.Assert(err, Equals, billy.ErrNotSupported)
}

func (s *PolyfillSuite) TestRoot(c *C) {
	c.Assert(s.Helper.Root(), Equals, string(filepath.Separator))
}

func (s *PolyfillSuite) TestCapabilities(c *C) {
	testCapabilities(c, new(test.BasicMock))
	testCapabilities(c, new(test.OnlyReadCapFs))
	testCapabilities(c, new(test.NoLockCapFs))
}

func testCapabilities(c *C, basic billy.Basic) {
	baseCapabilities := billy.Capabilities(basic)

	fs := New(basic)
	capabilities := billy.Capabilities(fs)

	c.Assert(capabilities, Equals, baseCapabilities)
}
