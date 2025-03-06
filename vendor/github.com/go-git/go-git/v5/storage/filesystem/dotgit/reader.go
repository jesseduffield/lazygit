package dotgit

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/objfile"
	"github.com/go-git/go-git/v5/utils/ioutil"
)

var _ (plumbing.EncodedObject) = &EncodedObject{}

type EncodedObject struct {
	dir *DotGit
	h   plumbing.Hash
	t   plumbing.ObjectType
	sz  int64
}

func (e *EncodedObject) Hash() plumbing.Hash {
	return e.h
}

func (e *EncodedObject) Reader() (io.ReadCloser, error) {
	f, err := e.dir.Object(e.h)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, plumbing.ErrObjectNotFound
		}

		return nil, err
	}
	r, err := objfile.NewReader(f)
	if err != nil {
		return nil, err
	}

	t, size, err := r.Header()
	if err != nil {
		_ = r.Close()
		return nil, err
	}
	if t != e.t {
		_ = r.Close()
		return nil, objfile.ErrHeader
	}
	if size != e.sz {
		_ = r.Close()
		return nil, objfile.ErrHeader
	}
	return ioutil.NewReadCloserWithCloser(r, f.Close), nil
}

func (e *EncodedObject) SetType(plumbing.ObjectType) {}

func (e *EncodedObject) Type() plumbing.ObjectType {
	return e.t
}

func (e *EncodedObject) Size() int64 {
	return e.sz
}

func (e *EncodedObject) SetSize(int64) {}

func (e *EncodedObject) Writer() (io.WriteCloser, error) {
	return nil, fmt.Errorf("not supported")
}

func NewEncodedObject(dir *DotGit, h plumbing.Hash, t plumbing.ObjectType, size int64) *EncodedObject {
	return &EncodedObject{
		dir: dir,
		h:   h,
		t:   t,
		sz:  size,
	}
}
