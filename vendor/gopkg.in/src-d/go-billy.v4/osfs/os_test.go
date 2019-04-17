package osfs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/test"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type OSSuite struct {
	test.FilesystemSuite
	path string
}

var _ = Suite(&OSSuite{})

func (s *OSSuite) SetUpTest(c *C) {
	s.path, _ = ioutil.TempDir(os.TempDir(), "go-billy-osfs-test")
	s.FilesystemSuite = test.NewFilesystemSuite(New(s.path))
}

func (s *OSSuite) TearDownTest(c *C) {
	err := os.RemoveAll(s.path)
	c.Assert(err, IsNil)
}

func (s *OSSuite) TestOpenDoesNotCreateDir(c *C) {
	_, err := s.FS.Open("dir/non-existent")
	c.Assert(err, NotNil)

	_, err = os.Stat(filepath.Join(s.path, "dir"))
	c.Assert(os.IsNotExist(err), Equals, true)
}

func (s *OSSuite) TestCapabilities(c *C) {
	_, ok := s.FS.(billy.Capable)
	c.Assert(ok, Equals, true)

	caps := billy.Capabilities(s.FS)
	c.Assert(caps, Equals, billy.AllCapabilities)
}
