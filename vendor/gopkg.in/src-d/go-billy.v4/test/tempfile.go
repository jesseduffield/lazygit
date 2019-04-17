package test

import (
	"strings"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/util"
)

// TempFileSuite is a convenient test suite to validate any implementation of
// billy.TempFile
type TempFileSuite struct {
	FS interface {
		billy.Basic
		billy.TempFile
	}
}

func (s *TempFileSuite) TestTempFile(c *C) {
	f, err := s.FS.TempFile("", "bar")
	c.Assert(err, IsNil)
	c.Assert(f.Close(), IsNil)

	c.Assert(strings.Index(f.Name(), "bar"), Not(Equals), -1)
}

func (s *TempFileSuite) TestTempFileWithPath(c *C) {
	f, err := s.FS.TempFile("foo", "bar")
	c.Assert(err, IsNil)
	c.Assert(f.Close(), IsNil)

	c.Assert(strings.HasPrefix(f.Name(), s.FS.Join("foo", "bar")), Equals, true)
}

func (s *TempFileSuite) TestTempFileFullWithPath(c *C) {
	f, err := s.FS.TempFile("/foo", "bar")
	c.Assert(err, IsNil)
	c.Assert(f.Close(), IsNil)

	c.Assert(strings.Index(f.Name(), s.FS.Join("foo", "bar")), Not(Equals), -1)
}

func (s *TempFileSuite) TestRemoveTempFile(c *C) {
	f, err := s.FS.TempFile("test-dir", "test-prefix")
	c.Assert(err, IsNil)

	fn := f.Name()
	c.Assert(f.Close(), IsNil)
	c.Assert(s.FS.Remove(fn), IsNil)
}

func (s *TempFileSuite) TestRenameTempFile(c *C) {
	f, err := s.FS.TempFile("test-dir", "test-prefix")
	c.Assert(err, IsNil)

	fn := f.Name()
	c.Assert(f.Close(), IsNil)
	c.Assert(s.FS.Rename(fn, "other-path"), IsNil)
}

func (s *TempFileSuite) TestTempFileMany(c *C) {
	for i := 0; i < 1024; i++ {
		var fs []billy.File

		for j := 0; j < 100; j++ {
			f, err := s.FS.TempFile("test-dir", "test-prefix")
			c.Assert(err, IsNil)
			fs = append(fs, f)
		}

		for _, f := range fs {
			c.Assert(f.Close(), IsNil)
			c.Assert(s.FS.Remove(f.Name()), IsNil)
		}
	}
}

func (s *TempFileSuite) TestTempFileManyWithUtil(c *C) {
	for i := 0; i < 1024; i++ {
		var fs []billy.File

		for j := 0; j < 100; j++ {
			f, err := util.TempFile(s.FS, "test-dir", "test-prefix")
			c.Assert(err, IsNil)
			fs = append(fs, f)
		}

		for _, f := range fs {
			c.Assert(f.Close(), IsNil)
			c.Assert(s.FS.Remove(f.Name()), IsNil)
		}
	}
}
