package packfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/utils/ioutil"
	"github.com/go-git/go-git/v5/utils/sync"
)

var (
	// ErrReferenceDeltaNotFound is returned when the reference delta is not
	// found.
	ErrReferenceDeltaNotFound = errors.New("reference delta not found")

	// ErrNotSeekableSource is returned when the source for the parser is not
	// seekable and a storage was not provided, so it can't be parsed.
	ErrNotSeekableSource = errors.New("parser source is not seekable and storage was not provided")

	// ErrDeltaNotCached is returned when the delta could not be found in cache.
	ErrDeltaNotCached = errors.New("delta could not be found in cache")
)

// Observer interface is implemented by index encoders.
type Observer interface {
	// OnHeader is called when a new packfile is opened.
	OnHeader(count uint32) error
	// OnInflatedObjectHeader is called for each object header read.
	OnInflatedObjectHeader(t plumbing.ObjectType, objSize int64, pos int64) error
	// OnInflatedObjectContent is called for each decoded object.
	OnInflatedObjectContent(h plumbing.Hash, pos int64, crc uint32, content []byte) error
	// OnFooter is called when decoding is done.
	OnFooter(h plumbing.Hash) error
}

// Parser decodes a packfile and calls any observer associated to it. Is used
// to generate indexes.
type Parser struct {
	storage    storer.EncodedObjectStorer
	scanner    *Scanner
	count      uint32
	oi         []*objectInfo
	oiByHash   map[plumbing.Hash]*objectInfo
	oiByOffset map[int64]*objectInfo
	checksum   plumbing.Hash

	cache *cache.BufferLRU
	// delta content by offset, only used if source is not seekable
	deltas map[int64][]byte

	ob []Observer
}

// NewParser creates a new Parser. The Scanner source must be seekable.
// If it's not, NewParserWithStorage should be used instead.
func NewParser(scanner *Scanner, ob ...Observer) (*Parser, error) {
	return NewParserWithStorage(scanner, nil, ob...)
}

// NewParserWithStorage creates a new Parser. The scanner source must either
// be seekable or a storage must be provided.
func NewParserWithStorage(
	scanner *Scanner,
	storage storer.EncodedObjectStorer,
	ob ...Observer,
) (*Parser, error) {
	if !scanner.IsSeekable && storage == nil {
		return nil, ErrNotSeekableSource
	}

	var deltas map[int64][]byte
	if !scanner.IsSeekable {
		deltas = make(map[int64][]byte)
	}

	return &Parser{
		storage: storage,
		scanner: scanner,
		ob:      ob,
		count:   0,
		cache:   cache.NewBufferLRUDefault(),
		deltas:  deltas,
	}, nil
}

func (p *Parser) forEachObserver(f func(o Observer) error) error {
	for _, o := range p.ob {
		if err := f(o); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) onHeader(count uint32) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnHeader(count)
	})
}

func (p *Parser) onInflatedObjectHeader(
	t plumbing.ObjectType,
	objSize int64,
	pos int64,
) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnInflatedObjectHeader(t, objSize, pos)
	})
}

func (p *Parser) onInflatedObjectContent(
	h plumbing.Hash,
	pos int64,
	crc uint32,
	content []byte,
) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnInflatedObjectContent(h, pos, crc, content)
	})
}

func (p *Parser) onFooter(h plumbing.Hash) error {
	return p.forEachObserver(func(o Observer) error {
		return o.OnFooter(h)
	})
}

// Parse start decoding phase of the packfile.
func (p *Parser) Parse() (plumbing.Hash, error) {
	if err := p.init(); err != nil {
		return plumbing.ZeroHash, err
	}

	if err := p.indexObjects(); err != nil {
		return plumbing.ZeroHash, err
	}

	var err error
	p.checksum, err = p.scanner.Checksum()
	if err != nil && err != io.EOF {
		return plumbing.ZeroHash, err
	}

	if err := p.resolveDeltas(); err != nil {
		return plumbing.ZeroHash, err
	}

	if err := p.onFooter(p.checksum); err != nil {
		return plumbing.ZeroHash, err
	}

	return p.checksum, nil
}

func (p *Parser) init() error {
	_, c, err := p.scanner.Header()
	if err != nil {
		return err
	}

	if err := p.onHeader(c); err != nil {
		return err
	}

	p.count = c
	p.oiByHash = make(map[plumbing.Hash]*objectInfo, p.count)
	p.oiByOffset = make(map[int64]*objectInfo, p.count)
	p.oi = make([]*objectInfo, p.count)

	return nil
}

type objectHeaderWriter func(typ plumbing.ObjectType, sz int64) error

type lazyObjectWriter interface {
	// LazyWriter enables an object to be lazily written.
	// It returns:
	// - w: a writer to receive the object's content.
	// - lwh: a func to write the object header.
	// - err: any error from the initial writer creation process.
	//
	// Note that if the object header is not written BEFORE the writer
	// is used, this will result in an invalid object.
	LazyWriter() (w io.WriteCloser, lwh objectHeaderWriter, err error)
}

func (p *Parser) indexObjects() error {
	buf := sync.GetBytesBuffer()
	defer sync.PutBytesBuffer(buf)

	for i := uint32(0); i < p.count; i++ {
		oh, err := p.scanner.NextObjectHeader()
		if err != nil {
			return err
		}

		delta := false
		var ota *objectInfo
		switch t := oh.Type; t {
		case plumbing.OFSDeltaObject:
			delta = true

			parent, ok := p.oiByOffset[oh.OffsetReference]
			if !ok {
				return plumbing.ErrObjectNotFound
			}

			ota = newDeltaObject(oh.Offset, oh.Length, t, parent)
			parent.Children = append(parent.Children, ota)
		case plumbing.REFDeltaObject:
			delta = true
			parent, ok := p.oiByHash[oh.Reference]
			if !ok {
				// can't find referenced object in this pack file
				// this must be a "thin" pack.
				parent = &objectInfo{ //Placeholder parent
					SHA1:        oh.Reference,
					ExternalRef: true, // mark as an external reference that must be resolved
					Type:        plumbing.AnyObject,
					DiskType:    plumbing.AnyObject,
				}
				p.oiByHash[oh.Reference] = parent
			}
			ota = newDeltaObject(oh.Offset, oh.Length, t, parent)
			parent.Children = append(parent.Children, ota)

		default:
			ota = newBaseObject(oh.Offset, oh.Length, t)
		}

		hasher := plumbing.NewHasher(oh.Type, oh.Length)
		writers := []io.Writer{hasher}
		var obj *plumbing.MemoryObject

		// Lazy writing is only available for non-delta objects.
		if p.storage != nil && !delta {
			// When a storage is set and supports lazy writing,
			// use that instead of creating a memory object.
			if low, ok := p.storage.(lazyObjectWriter); ok {
				ow, lwh, err := low.LazyWriter()
				if err != nil {
					return err
				}

				if err = lwh(oh.Type, oh.Length); err != nil {
					return err
				}

				defer ow.Close()
				writers = append(writers, ow)
			} else {
				obj = new(plumbing.MemoryObject)
				obj.SetSize(oh.Length)
				obj.SetType(oh.Type)

				writers = append(writers, obj)
			}
		}
		if delta && !p.scanner.IsSeekable {
			buf.Reset()
			buf.Grow(int(oh.Length))
			writers = append(writers, buf)
		}

		mw := io.MultiWriter(writers...)

		_, crc, err := p.scanner.NextObject(mw)
		if err != nil {
			return err
		}

		// Non delta objects needs to be added into the storage. This
		// is only required when lazy writing is not supported.
		if obj != nil {
			if _, err := p.storage.SetEncodedObject(obj); err != nil {
				return err
			}
		}

		ota.Crc32 = crc
		ota.Length = oh.Length

		if !delta {
			sha1 := hasher.Sum()

			// Move children of placeholder parent into actual parent, in case this
			// was a non-external delta reference.
			if placeholder, ok := p.oiByHash[sha1]; ok {
				ota.Children = placeholder.Children
				for _, c := range ota.Children {
					c.Parent = ota
				}
			}

			ota.SHA1 = sha1
			p.oiByHash[ota.SHA1] = ota
		}

		if delta && !p.scanner.IsSeekable {
			data := buf.Bytes()
			p.deltas[oh.Offset] = make([]byte, len(data))
			copy(p.deltas[oh.Offset], data)
		}

		p.oiByOffset[oh.Offset] = ota
		p.oi[i] = ota
	}

	return nil
}

func (p *Parser) resolveDeltas() error {
	buf := sync.GetBytesBuffer()
	defer sync.PutBytesBuffer(buf)

	for _, obj := range p.oi {
		buf.Reset()
		buf.Grow(int(obj.Length))
		err := p.get(obj, buf)
		if err != nil {
			return err
		}

		if err := p.onInflatedObjectHeader(obj.Type, obj.Length, obj.Offset); err != nil {
			return err
		}

		if err := p.onInflatedObjectContent(obj.SHA1, obj.Offset, obj.Crc32, nil); err != nil {
			return err
		}

		if !obj.IsDelta() && len(obj.Children) > 0 {
			// Dealing with an io.ReaderAt object, means we can
			// create it once and reuse across all children.
			r := bytes.NewReader(buf.Bytes())
			for _, child := range obj.Children {
				// Even though we are discarding the output, we still need to read it to
				// so that the scanner can advance to the next object, and the SHA1 can be
				// calculated.
				if err := p.resolveObject(io.Discard, child, r); err != nil {
					return err
				}
				p.resolveExternalRef(child)
			}

			// Remove the delta from the cache.
			if obj.DiskType.IsDelta() && !p.scanner.IsSeekable {
				delete(p.deltas, obj.Offset)
			}
		}
	}

	return nil
}

func (p *Parser) resolveExternalRef(o *objectInfo) {
	if ref, ok := p.oiByHash[o.SHA1]; ok && ref.ExternalRef {
		p.oiByHash[o.SHA1] = o
		o.Children = ref.Children
		for _, c := range o.Children {
			c.Parent = o
		}
	}
}

func (p *Parser) get(o *objectInfo, buf *bytes.Buffer) (err error) {
	if !o.ExternalRef { // skip cache check for placeholder parents
		b, ok := p.cache.Get(o.Offset)
		if ok {
			_, err := buf.Write(b)
			return err
		}
	}

	// If it's not on the cache and is not a delta we can try to find it in the
	// storage, if there's one. External refs must enter here.
	if p.storage != nil && !o.Type.IsDelta() {
		var e plumbing.EncodedObject
		e, err = p.storage.EncodedObject(plumbing.AnyObject, o.SHA1)
		if err != nil {
			return err
		}
		o.Type = e.Type()

		var r io.ReadCloser
		r, err = e.Reader()
		if err != nil {
			return err
		}

		defer ioutil.CheckClose(r, &err)

		_, err = buf.ReadFrom(io.LimitReader(r, e.Size()))
		return err
	}

	if o.ExternalRef {
		// we were not able to resolve a ref in a thin pack
		return ErrReferenceDeltaNotFound
	}

	if o.DiskType.IsDelta() {
		b := sync.GetBytesBuffer()
		defer sync.PutBytesBuffer(b)
		buf.Grow(int(o.Length))
		err := p.get(o.Parent, b)
		if err != nil {
			return err
		}

		err = p.resolveObject(buf, o, bytes.NewReader(b.Bytes()))
		if err != nil {
			return err
		}
	} else {
		err := p.readData(buf, o)
		if err != nil {
			return err
		}
	}

	// If the scanner is seekable, caching this data into
	// memory by offset seems wasteful.
	// There is a trade-off to be considered here in terms
	// of execution time vs memory consumption.
	//
	// TODO: improve seekable execution time, so that we can
	// skip this cache.
	if len(o.Children) > 0 {
		data := make([]byte, buf.Len())
		copy(data, buf.Bytes())
		p.cache.Put(o.Offset, data)
	}
	return nil
}

// resolveObject resolves an object from base, using information
// provided by o.
//
// This call has the side-effect of changing field values
// from the object info o:
//   - Type: OFSDeltaObject may become the target type (e.g. Blob).
//   - Size: The size may be update with the target size.
//   - Hash: Zero hashes will be calculated as part of the object
//     resolution. Hence why this process can't be avoided even when w
//     is an io.Discard.
//
// base must be an io.ReaderAt, which is a requirement from
// patchDeltaStream. The main reason being that reversing an
// delta object may lead to going backs and forths within base,
// which is not supported by io.Reader.
func (p *Parser) resolveObject(
	w io.Writer,
	o *objectInfo,
	base io.ReaderAt,
) error {
	if !o.DiskType.IsDelta() {
		return nil
	}
	buf := sync.GetBytesBuffer()
	defer sync.PutBytesBuffer(buf)
	err := p.readData(buf, o)
	if err != nil {
		return err
	}

	writers := []io.Writer{w}
	var obj *plumbing.MemoryObject
	var lwh objectHeaderWriter

	if p.storage != nil {
		if low, ok := p.storage.(lazyObjectWriter); ok {
			ow, wh, err := low.LazyWriter()
			if err != nil {
				return err
			}
			lwh = wh

			defer ow.Close()
			writers = append(writers, ow)
		} else {
			obj = new(plumbing.MemoryObject)
			ow, err := obj.Writer()
			if err != nil {
				return err
			}

			writers = append(writers, ow)
		}
	}

	mw := io.MultiWriter(writers...)

	err = applyPatchBase(o, base, buf, mw, lwh)
	if err != nil {
		return err
	}

	if obj != nil {
		obj.SetType(o.Type)
		obj.SetSize(o.Size()) // Size here is correct as it was populated by applyPatchBase.
		if _, err := p.storage.SetEncodedObject(obj); err != nil {
			return err
		}
	}
	return err
}

func (p *Parser) readData(w io.Writer, o *objectInfo) error {
	if !p.scanner.IsSeekable && o.DiskType.IsDelta() {
		data, ok := p.deltas[o.Offset]
		if !ok {
			return ErrDeltaNotCached
		}
		_, err := w.Write(data)
		return err
	}

	if _, err := p.scanner.SeekObjectHeader(o.Offset); err != nil {
		return err
	}

	if _, _, err := p.scanner.NextObject(w); err != nil {
		return err
	}
	return nil
}

// applyPatchBase applies the patch to target.
//
// Note that ota will be updated based on the description in resolveObject.
func applyPatchBase(ota *objectInfo, base io.ReaderAt, delta io.Reader, target io.Writer, wh objectHeaderWriter) error {
	if target == nil {
		return fmt.Errorf("cannot apply patch against nil target")
	}

	typ := ota.Type
	if ota.SHA1 == plumbing.ZeroHash {
		typ = ota.Parent.Type
	}

	sz, h, err := patchDeltaWriter(target, base, delta, typ, wh)
	if err != nil {
		return err
	}

	if ota.SHA1 == plumbing.ZeroHash {
		ota.Type = typ
		ota.Length = int64(sz)
		ota.SHA1 = h
	}

	return nil
}

func getSHA1(t plumbing.ObjectType, data []byte) (plumbing.Hash, error) {
	hasher := plumbing.NewHasher(t, int64(len(data)))
	if _, err := hasher.Write(data); err != nil {
		return plumbing.ZeroHash, err
	}

	return hasher.Sum(), nil
}

type objectInfo struct {
	Offset      int64
	Length      int64
	Type        plumbing.ObjectType
	DiskType    plumbing.ObjectType
	ExternalRef bool // indicates this is an external reference in a thin pack file

	Crc32 uint32

	Parent   *objectInfo
	Children []*objectInfo
	SHA1     plumbing.Hash
}

func newBaseObject(offset, length int64, t plumbing.ObjectType) *objectInfo {
	return newDeltaObject(offset, length, t, nil)
}

func newDeltaObject(
	offset, length int64,
	t plumbing.ObjectType,
	parent *objectInfo,
) *objectInfo {
	obj := &objectInfo{
		Offset:   offset,
		Length:   length,
		Type:     t,
		DiskType: t,
		Crc32:    0,
		Parent:   parent,
	}

	return obj
}

func (o *objectInfo) IsDelta() bool {
	return o.Type.IsDelta()
}

func (o *objectInfo) Size() int64 {
	return o.Length
}
