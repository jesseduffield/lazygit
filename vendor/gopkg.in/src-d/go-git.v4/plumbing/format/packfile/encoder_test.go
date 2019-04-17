package packfile

import (
	"bytes"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type EncoderSuite struct {
	fixtures.Suite
	buf   *bytes.Buffer
	store *memory.Storage
	enc   *Encoder
}

var _ = Suite(&EncoderSuite{})

func (s *EncoderSuite) SetUpTest(c *C) {
	s.buf = bytes.NewBuffer(nil)
	s.store = memory.NewStorage()
	s.enc = NewEncoder(s.buf, s.store, false)
}

func (s *EncoderSuite) TestCorrectPackHeader(c *C) {
	hash, err := s.enc.Encode([]plumbing.Hash{}, 10)
	c.Assert(err, IsNil)

	hb := [20]byte(hash)

	// PACK + VERSION + OBJECTS + HASH
	expectedResult := []byte{'P', 'A', 'C', 'K', 0, 0, 0, 2, 0, 0, 0, 0}
	expectedResult = append(expectedResult, hb[:]...)

	result := s.buf.Bytes()

	c.Assert(result, DeepEquals, expectedResult)
}

func (s *EncoderSuite) TestCorrectPackWithOneEmptyObject(c *C) {
	o := &plumbing.MemoryObject{}
	o.SetType(plumbing.CommitObject)
	o.SetSize(0)
	_, err := s.store.SetEncodedObject(o)
	c.Assert(err, IsNil)

	hash, err := s.enc.Encode([]plumbing.Hash{o.Hash()}, 10)
	c.Assert(err, IsNil)

	// PACK + VERSION(2) + OBJECT NUMBER(1)
	expectedResult := []byte{'P', 'A', 'C', 'K', 0, 0, 0, 2, 0, 0, 0, 1}
	// OBJECT HEADER(TYPE + SIZE)= 0001 0000
	expectedResult = append(expectedResult, []byte{16}...)

	// Zlib header
	expectedResult = append(expectedResult,
		[]byte{120, 156, 1, 0, 0, 255, 255, 0, 0, 0, 1}...)

	// + HASH
	hb := [20]byte(hash)
	expectedResult = append(expectedResult, hb[:]...)

	result := s.buf.Bytes()

	c.Assert(result, DeepEquals, expectedResult)
}

func (s *EncoderSuite) TestMaxObjectSize(c *C) {
	o := s.store.NewEncodedObject()
	o.SetSize(9223372036854775807)
	o.SetType(plumbing.CommitObject)
	_, err := s.store.SetEncodedObject(o)
	c.Assert(err, IsNil)
	hash, err := s.enc.Encode([]plumbing.Hash{o.Hash()}, 10)
	c.Assert(err, IsNil)
	c.Assert(hash.IsZero(), Not(Equals), true)
}

func (s *EncoderSuite) TestHashNotFound(c *C) {
	h, err := s.enc.Encode([]plumbing.Hash{plumbing.NewHash("BAD")}, 10)
	c.Assert(h, Equals, plumbing.ZeroHash)
	c.Assert(err, NotNil)
	c.Assert(err, Equals, plumbing.ErrObjectNotFound)
}

func (s *EncoderSuite) TestDecodeEncodeWithDeltaDecodeREF(c *C) {
	s.enc = NewEncoder(s.buf, s.store, true)
	s.simpleDeltaTest(c)
}

func (s *EncoderSuite) TestDecodeEncodeWithDeltaDecodeOFS(c *C) {
	s.enc = NewEncoder(s.buf, s.store, false)
	s.simpleDeltaTest(c)
}

func (s *EncoderSuite) TestDecodeEncodeWithDeltasDecodeREF(c *C) {
	s.enc = NewEncoder(s.buf, s.store, true)
	s.deltaOverDeltaTest(c)
}

func (s *EncoderSuite) TestDecodeEncodeWithDeltasDecodeOFS(c *C) {
	s.enc = NewEncoder(s.buf, s.store, false)
	s.deltaOverDeltaTest(c)
}

func (s *EncoderSuite) TestDecodeEncodeWithCycleREF(c *C) {
	s.enc = NewEncoder(s.buf, s.store, true)
	s.deltaOverDeltaCyclicTest(c)
}

func (s *EncoderSuite) TestDecodeEncodeWithCycleOFS(c *C) {
	s.enc = NewEncoder(s.buf, s.store, false)
	s.deltaOverDeltaCyclicTest(c)
}

func (s *EncoderSuite) simpleDeltaTest(c *C) {
	srcObject := newObject(plumbing.BlobObject, []byte("0"))
	targetObject := newObject(plumbing.BlobObject, []byte("01"))

	deltaObject, err := GetDelta(srcObject, targetObject)
	c.Assert(err, IsNil)

	srcToPack := newObjectToPack(srcObject)
	encHash, err := s.enc.encode([]*ObjectToPack{
		srcToPack,
		newDeltaObjectToPack(srcToPack, targetObject, deltaObject),
	})
	c.Assert(err, IsNil)

	scanner := NewScanner(s.buf)

	storage := memory.NewStorage()
	d, err := NewDecoder(scanner, storage)
	c.Assert(err, IsNil)

	decHash, err := d.Decode()
	c.Assert(err, IsNil)

	c.Assert(encHash, Equals, decHash)

	decSrc, err := storage.EncodedObject(srcObject.Type(), srcObject.Hash())
	c.Assert(err, IsNil)
	c.Assert(decSrc, DeepEquals, srcObject)

	decTarget, err := storage.EncodedObject(targetObject.Type(), targetObject.Hash())
	c.Assert(err, IsNil)
	c.Assert(decTarget, DeepEquals, targetObject)
}

func (s *EncoderSuite) deltaOverDeltaTest(c *C) {
	srcObject := newObject(plumbing.BlobObject, []byte("0"))
	targetObject := newObject(plumbing.BlobObject, []byte("01"))
	otherTargetObject := newObject(plumbing.BlobObject, []byte("011111"))

	deltaObject, err := GetDelta(srcObject, targetObject)
	c.Assert(err, IsNil)
	c.Assert(deltaObject.Hash(), Not(Equals), plumbing.ZeroHash)

	otherDeltaObject, err := GetDelta(targetObject, otherTargetObject)
	c.Assert(err, IsNil)
	c.Assert(otherDeltaObject.Hash(), Not(Equals), plumbing.ZeroHash)

	srcToPack := newObjectToPack(srcObject)
	targetToPack := newObjectToPack(targetObject)
	encHash, err := s.enc.encode([]*ObjectToPack{
		targetToPack,
		srcToPack,
		newDeltaObjectToPack(srcToPack, targetObject, deltaObject),
		newDeltaObjectToPack(targetToPack, otherTargetObject, otherDeltaObject),
	})
	c.Assert(err, IsNil)

	scanner := NewScanner(s.buf)
	storage := memory.NewStorage()
	d, err := NewDecoder(scanner, storage)
	c.Assert(err, IsNil)

	decHash, err := d.Decode()
	c.Assert(err, IsNil)

	c.Assert(encHash, Equals, decHash)

	decSrc, err := storage.EncodedObject(srcObject.Type(), srcObject.Hash())
	c.Assert(err, IsNil)
	c.Assert(decSrc, DeepEquals, srcObject)

	decTarget, err := storage.EncodedObject(targetObject.Type(), targetObject.Hash())
	c.Assert(err, IsNil)
	c.Assert(decTarget, DeepEquals, targetObject)

	decOtherTarget, err := storage.EncodedObject(otherTargetObject.Type(), otherTargetObject.Hash())
	c.Assert(err, IsNil)
	c.Assert(decOtherTarget, DeepEquals, otherTargetObject)
}

func (s *EncoderSuite) deltaOverDeltaCyclicTest(c *C) {
	o1 := newObject(plumbing.BlobObject, []byte("0"))
	o2 := newObject(plumbing.BlobObject, []byte("01"))
	o3 := newObject(plumbing.BlobObject, []byte("011111"))
	o4 := newObject(plumbing.BlobObject, []byte("01111100000"))

	_, err := s.store.SetEncodedObject(o1)
	c.Assert(err, IsNil)
	_, err = s.store.SetEncodedObject(o2)
	c.Assert(err, IsNil)
	_, err = s.store.SetEncodedObject(o3)
	c.Assert(err, IsNil)
	_, err = s.store.SetEncodedObject(o4)
	c.Assert(err, IsNil)

	d2, err := GetDelta(o1, o2)
	c.Assert(err, IsNil)

	d3, err := GetDelta(o4, o3)
	c.Assert(err, IsNil)

	d4, err := GetDelta(o3, o4)
	c.Assert(err, IsNil)

	po1 := newObjectToPack(o1)
	pd2 := newDeltaObjectToPack(po1, o2, d2)
	pd3 := newObjectToPack(o3)
	pd4 := newObjectToPack(o4)

	pd3.SetDelta(pd4, d3)
	pd4.SetDelta(pd3, d4)

	// SetOriginal is used by delta selector when generating ObjectToPack.
	// It also fills type, hash and size values to be used when Original
	// is nil.
	po1.SetOriginal(po1.Original)
	pd2.SetOriginal(pd2.Original)
	pd2.CleanOriginal()

	pd3.SetOriginal(pd3.Original)
	pd3.CleanOriginal()

	pd4.SetOriginal(pd4.Original)

	encHash, err := s.enc.encode([]*ObjectToPack{
		po1,
		pd2,
		pd3,
		pd4,
	})
	c.Assert(err, IsNil)

	scanner := NewScanner(s.buf)
	storage := memory.NewStorage()
	d, err := NewDecoder(scanner, storage)
	c.Assert(err, IsNil)

	decHash, err := d.Decode()
	c.Assert(err, IsNil)

	c.Assert(encHash, Equals, decHash)

	decSrc, err := storage.EncodedObject(o1.Type(), o1.Hash())
	c.Assert(err, IsNil)
	c.Assert(decSrc, DeepEquals, o1)

	decTarget, err := storage.EncodedObject(o2.Type(), o2.Hash())
	c.Assert(err, IsNil)
	c.Assert(decTarget, DeepEquals, o2)

	decOtherTarget, err := storage.EncodedObject(o3.Type(), o3.Hash())
	c.Assert(err, IsNil)
	c.Assert(decOtherTarget, DeepEquals, o3)

	decAnotherTarget, err := storage.EncodedObject(o4.Type(), o4.Hash())
	c.Assert(err, IsNil)
	c.Assert(decAnotherTarget, DeepEquals, o4)
}
