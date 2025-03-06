package packfile

import (
	"io"

	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/format/idxfile"
	"github.com/go-git/go-git/v5/utils/ioutil"
)

// FSObject is an object from the packfile on the filesystem.
type FSObject struct {
	hash                 plumbing.Hash
	offset               int64
	size                 int64
	typ                  plumbing.ObjectType
	index                idxfile.Index
	fs                   billy.Filesystem
	path                 string
	cache                cache.Object
	largeObjectThreshold int64
}

// NewFSObject creates a new filesystem object.
func NewFSObject(
	hash plumbing.Hash,
	finalType plumbing.ObjectType,
	offset int64,
	contentSize int64,
	index idxfile.Index,
	fs billy.Filesystem,
	path string,
	cache cache.Object,
	largeObjectThreshold int64,
) *FSObject {
	return &FSObject{
		hash:                 hash,
		offset:               offset,
		size:                 contentSize,
		typ:                  finalType,
		index:                index,
		fs:                   fs,
		path:                 path,
		cache:                cache,
		largeObjectThreshold: largeObjectThreshold,
	}
}

// Reader implements the plumbing.EncodedObject interface.
func (o *FSObject) Reader() (io.ReadCloser, error) {
	obj, ok := o.cache.Get(o.hash)
	if ok && obj != o {
		reader, err := obj.Reader()
		if err != nil {
			return nil, err
		}

		return reader, nil
	}

	f, err := o.fs.Open(o.path)
	if err != nil {
		return nil, err
	}

	p := NewPackfileWithCache(o.index, nil, f, o.cache, o.largeObjectThreshold)
	if o.largeObjectThreshold > 0 && o.size > o.largeObjectThreshold {
		// We have a big object
		h, err := p.objectHeaderAtOffset(o.offset)
		if err != nil {
			return nil, err
		}

		r, err := p.getReaderDirect(h)
		if err != nil {
			_ = f.Close()
			return nil, err
		}
		return ioutil.NewReadCloserWithCloser(r, f.Close), nil
	}
	r, err := p.getObjectContent(o.offset)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return r, nil
}

// SetSize implements the plumbing.EncodedObject interface. This method
// is a noop.
func (o *FSObject) SetSize(int64) {}

// SetType implements the plumbing.EncodedObject interface. This method is
// a noop.
func (o *FSObject) SetType(plumbing.ObjectType) {}

// Hash implements the plumbing.EncodedObject interface.
func (o *FSObject) Hash() plumbing.Hash { return o.hash }

// Size implements the plumbing.EncodedObject interface.
func (o *FSObject) Size() int64 { return o.size }

// Type implements the plumbing.EncodedObject interface.
func (o *FSObject) Type() plumbing.ObjectType {
	return o.typ
}

// Writer implements the plumbing.EncodedObject interface. This method always
// returns a nil writer.
func (o *FSObject) Writer() (io.WriteCloser, error) {
	return nil, nil
}
