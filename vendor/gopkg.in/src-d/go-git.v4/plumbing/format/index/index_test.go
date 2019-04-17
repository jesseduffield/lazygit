package index

import (
	"path/filepath"

	. "gopkg.in/check.v1"
)

func (s *IndexSuite) TestIndexAdd(c *C) {
	idx := &Index{}
	e := idx.Add("foo")
	e.Size = 42

	e, err := idx.Entry("foo")
	c.Assert(err, IsNil)
	c.Assert(e.Name, Equals, "foo")
	c.Assert(e.Size, Equals, uint32(42))
}

func (s *IndexSuite) TestIndexEntry(c *C) {
	idx := &Index{
		Entries: []*Entry{
			{Name: "foo", Size: 42},
			{Name: "bar", Size: 82},
		},
	}

	e, err := idx.Entry("foo")
	c.Assert(err, IsNil)
	c.Assert(e.Name, Equals, "foo")

	e, err = idx.Entry("missing")
	c.Assert(e, IsNil)
	c.Assert(err, Equals, ErrEntryNotFound)
}

func (s *IndexSuite) TestIndexRemove(c *C) {
	idx := &Index{
		Entries: []*Entry{
			{Name: "foo", Size: 42},
			{Name: "bar", Size: 82},
		},
	}

	e, err := idx.Remove("foo")
	c.Assert(err, IsNil)
	c.Assert(e.Name, Equals, "foo")

	e, err = idx.Remove("foo")
	c.Assert(e, IsNil)
	c.Assert(err, Equals, ErrEntryNotFound)
}

func (s *IndexSuite) TestIndexGlob(c *C) {
	idx := &Index{
		Entries: []*Entry{
			{Name: "foo/bar/bar", Size: 42},
			{Name: "foo/baz/qux", Size: 42},
			{Name: "fux", Size: 82},
		},
	}

	m, err := idx.Glob(filepath.Join("foo", "b*"))
	c.Assert(err, IsNil)
	c.Assert(m, HasLen, 2)
	c.Assert(m[0].Name, Equals, "foo/bar/bar")
	c.Assert(m[1].Name, Equals, "foo/baz/qux")

	m, err = idx.Glob("f*")
	c.Assert(err, IsNil)
	c.Assert(m, HasLen, 3)

	m, err = idx.Glob("f*/baz/q*")
	c.Assert(err, IsNil)
	c.Assert(m, HasLen, 1)
}
