package packfile

import (
	"bytes"
	"io"

	"gopkg.in/src-d/go-git.v4/plumbing"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type ScannerSuite struct {
	fixtures.Suite
}

var _ = Suite(&ScannerSuite{})

func (s *ScannerSuite) TestHeader(c *C) {
	r := fixtures.Basic().One().Packfile()
	p := NewScanner(r)

	version, objects, err := p.Header()
	c.Assert(err, IsNil)
	c.Assert(version, Equals, VersionSupported)
	c.Assert(objects, Equals, uint32(31))
}

func (s *ScannerSuite) TestNextObjectHeaderWithoutHeader(c *C) {
	r := fixtures.Basic().One().Packfile()
	p := NewScanner(r)

	h, err := p.NextObjectHeader()
	c.Assert(err, IsNil)
	c.Assert(h, DeepEquals, &expectedHeadersOFS[0])

	version, objects, err := p.Header()
	c.Assert(err, IsNil)
	c.Assert(version, Equals, VersionSupported)
	c.Assert(objects, Equals, uint32(31))
}

func (s *ScannerSuite) TestNextObjectHeaderREFDelta(c *C) {
	s.testNextObjectHeader(c, "ref-delta", expectedHeadersREF, expectedCRCREF)
}

func (s *ScannerSuite) TestNextObjectHeaderOFSDelta(c *C) {
	s.testNextObjectHeader(c, "ofs-delta", expectedHeadersOFS, expectedCRCOFS)
}

func (s *ScannerSuite) testNextObjectHeader(c *C, tag string,
	expected []ObjectHeader, expectedCRC []uint32) {

	r := fixtures.Basic().ByTag(tag).One().Packfile()
	p := NewScanner(r)

	_, objects, err := p.Header()
	c.Assert(err, IsNil)

	for i := 0; i < int(objects); i++ {
		h, err := p.NextObjectHeader()
		c.Assert(err, IsNil)
		c.Assert(*h, DeepEquals, expected[i])

		buf := bytes.NewBuffer(nil)
		n, crcFromScanner, err := p.NextObject(buf)
		c.Assert(err, IsNil)
		c.Assert(n, Equals, h.Length)
		c.Assert(crcFromScanner, Equals, expectedCRC[i])
	}

	n, err := p.Checksum()
	c.Assert(err, IsNil)
	c.Assert(n, HasLen, 20)
}

func (s *ScannerSuite) TestNextObjectHeaderWithOutReadObject(c *C) {
	f := fixtures.Basic().ByTag("ref-delta").One()
	r := f.Packfile()
	p := NewScanner(r)

	_, objects, err := p.Header()
	c.Assert(err, IsNil)

	for i := 0; i < int(objects); i++ {
		h, _ := p.NextObjectHeader()
		c.Assert(err, IsNil)
		c.Assert(*h, DeepEquals, expectedHeadersREF[i])
	}

	err = p.discardObjectIfNeeded()
	c.Assert(err, IsNil)

	n, err := p.Checksum()
	c.Assert(err, IsNil)
	c.Assert(n, Equals, f.PackfileHash)
}

func (s *ScannerSuite) TestNextObjectHeaderWithOutReadObjectNonSeekable(c *C) {
	f := fixtures.Basic().ByTag("ref-delta").One()
	r := io.MultiReader(f.Packfile())
	p := NewScanner(r)

	_, objects, err := p.Header()
	c.Assert(err, IsNil)

	for i := 0; i < int(objects); i++ {
		h, _ := p.NextObjectHeader()
		c.Assert(err, IsNil)
		c.Assert(*h, DeepEquals, expectedHeadersREF[i])
	}

	err = p.discardObjectIfNeeded()
	c.Assert(err, IsNil)

	n, err := p.Checksum()
	c.Assert(err, IsNil)
	c.Assert(n, Equals, f.PackfileHash)
}

var expectedHeadersOFS = []ObjectHeader{
	{Type: plumbing.CommitObject, Offset: 12, Length: 254},
	{Type: plumbing.OFSDeltaObject, Offset: 186, Length: 93, OffsetReference: 12},
	{Type: plumbing.CommitObject, Offset: 286, Length: 242},
	{Type: plumbing.CommitObject, Offset: 449, Length: 242},
	{Type: plumbing.CommitObject, Offset: 615, Length: 333},
	{Type: plumbing.CommitObject, Offset: 838, Length: 332},
	{Type: plumbing.CommitObject, Offset: 1063, Length: 244},
	{Type: plumbing.CommitObject, Offset: 1230, Length: 243},
	{Type: plumbing.CommitObject, Offset: 1392, Length: 187},
	{Type: plumbing.BlobObject, Offset: 1524, Length: 189},
	{Type: plumbing.BlobObject, Offset: 1685, Length: 18},
	{Type: plumbing.BlobObject, Offset: 1713, Length: 1072},
	{Type: plumbing.BlobObject, Offset: 2351, Length: 76110},
	{Type: plumbing.BlobObject, Offset: 78050, Length: 2780},
	{Type: plumbing.BlobObject, Offset: 78882, Length: 217848},
	{Type: plumbing.BlobObject, Offset: 80725, Length: 706},
	{Type: plumbing.BlobObject, Offset: 80998, Length: 11488},
	{Type: plumbing.BlobObject, Offset: 84032, Length: 78},
	{Type: plumbing.TreeObject, Offset: 84115, Length: 272},
	{Type: plumbing.OFSDeltaObject, Offset: 84375, Length: 43, OffsetReference: 84115},
	{Type: plumbing.TreeObject, Offset: 84430, Length: 38},
	{Type: plumbing.TreeObject, Offset: 84479, Length: 75},
	{Type: plumbing.TreeObject, Offset: 84559, Length: 38},
	{Type: plumbing.TreeObject, Offset: 84608, Length: 34},
	{Type: plumbing.BlobObject, Offset: 84653, Length: 9},
	{Type: plumbing.OFSDeltaObject, Offset: 84671, Length: 6, OffsetReference: 84375},
	{Type: plumbing.OFSDeltaObject, Offset: 84688, Length: 9, OffsetReference: 84375},
	{Type: plumbing.OFSDeltaObject, Offset: 84708, Length: 6, OffsetReference: 84375},
	{Type: plumbing.OFSDeltaObject, Offset: 84725, Length: 5, OffsetReference: 84115},
	{Type: plumbing.OFSDeltaObject, Offset: 84741, Length: 8, OffsetReference: 84375},
	{Type: plumbing.OFSDeltaObject, Offset: 84760, Length: 4, OffsetReference: 84741},
}

var expectedCRCOFS = []uint32{
	0xaa07ba4b,
	0xf706df58,
	0x12438846,
	0x2905a38c,
	0xd9429436,
	0xbecfde4e,
	0x780e4b3e,
	0xdc18344f,
	0xcf4e4280,
	0x1f08118a,
	0xafded7b8,
	0xcc1428ed,
	0x1631d22f,
	0xbfff5850,
	0xd108e1d8,
	0x8e97ba25,
	0x7316ff70,
	0xdb4fce56,
	0x901cce2c,
	0xec4552b0,
	0x847905bf,
	0x3689459a,
	0xe67af94a,
	0xc2314a2e,
	0xcd987848,
	0x8a853a6d,
	0x70c6518,
	0x4f4108e2,
	0xd6fe09e9,
	0xf07a2804,
	0x1d75d6be,
}

var expectedHeadersREF = []ObjectHeader{
	{Type: plumbing.CommitObject, Offset: 12, Length: 254},
	{Type: plumbing.REFDeltaObject, Offset: 186, Length: 93,
		Reference: plumbing.NewHash("e8d3ffab552895c19b9fcf7aa264d277cde33881")},
	{Type: plumbing.CommitObject, Offset: 304, Length: 242},
	{Type: plumbing.CommitObject, Offset: 467, Length: 242},
	{Type: plumbing.CommitObject, Offset: 633, Length: 333},
	{Type: plumbing.CommitObject, Offset: 856, Length: 332},
	{Type: plumbing.CommitObject, Offset: 1081, Length: 243},
	{Type: plumbing.CommitObject, Offset: 1243, Length: 244},
	{Type: plumbing.CommitObject, Offset: 1410, Length: 187},
	{Type: plumbing.BlobObject, Offset: 1542, Length: 189},
	{Type: plumbing.BlobObject, Offset: 1703, Length: 18},
	{Type: plumbing.BlobObject, Offset: 1731, Length: 1072},
	{Type: plumbing.BlobObject, Offset: 2369, Length: 76110},
	{Type: plumbing.TreeObject, Offset: 78068, Length: 38},
	{Type: plumbing.BlobObject, Offset: 78117, Length: 2780},
	{Type: plumbing.TreeObject, Offset: 79049, Length: 75},
	{Type: plumbing.BlobObject, Offset: 79129, Length: 217848},
	{Type: plumbing.BlobObject, Offset: 80972, Length: 706},
	{Type: plumbing.TreeObject, Offset: 81265, Length: 38},
	{Type: plumbing.BlobObject, Offset: 81314, Length: 11488},
	{Type: plumbing.TreeObject, Offset: 84752, Length: 34},
	{Type: plumbing.BlobObject, Offset: 84797, Length: 78},
	{Type: plumbing.TreeObject, Offset: 84880, Length: 271},
	{Type: plumbing.REFDeltaObject, Offset: 85141, Length: 6,
		Reference: plumbing.NewHash("a8d315b2b1c615d43042c3a62402b8a54288cf5c")},
	{Type: plumbing.REFDeltaObject, Offset: 85176, Length: 37,
		Reference: plumbing.NewHash("fb72698cab7617ac416264415f13224dfd7a165e")},
	{Type: plumbing.BlobObject, Offset: 85244, Length: 9},
	{Type: plumbing.REFDeltaObject, Offset: 85262, Length: 9,
		Reference: plumbing.NewHash("fb72698cab7617ac416264415f13224dfd7a165e")},
	{Type: plumbing.REFDeltaObject, Offset: 85300, Length: 6,
		Reference: plumbing.NewHash("fb72698cab7617ac416264415f13224dfd7a165e")},
	{Type: plumbing.TreeObject, Offset: 85335, Length: 110},
	{Type: plumbing.REFDeltaObject, Offset: 85448, Length: 8,
		Reference: plumbing.NewHash("eba74343e2f15d62adedfd8c883ee0262b5c8021")},
	{Type: plumbing.TreeObject, Offset: 85485, Length: 73},
}

var expectedCRCREF = []uint32{
	0xaa07ba4b,
	0xfb4725a4,
	0x12438846,
	0x2905a38c,
	0xd9429436,
	0xbecfde4e,
	0xdc18344f,
	0x780e4b3e,
	0xcf4e4280,
	0x1f08118a,
	0xafded7b8,
	0xcc1428ed,
	0x1631d22f,
	0x847905bf,
	0x3e20f31d,
	0x3689459a,
	0xd108e1d8,
	0x71143d4a,
	0xe67af94a,
	0x739fb89f,
	0xc2314a2e,
	0x87864926,
	0x415d752f,
	0xf72fb182,
	0x3ffa37d4,
	0xcd987848,
	0x2f20ac8f,
	0xf2f0575,
	0x7d8726e1,
	0x740bf39,
	0x26af4735,
}
