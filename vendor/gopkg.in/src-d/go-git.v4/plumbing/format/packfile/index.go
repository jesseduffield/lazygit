package packfile

import (
	"sort"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/idxfile"
)

// Index is an in-memory representation of a packfile index.
// This uses idxfile.Idxfile under the hood to obtain indexes from .idx files
// or to store them.
type Index struct {
	byHash   map[plumbing.Hash]*idxfile.Entry
	byOffset []*idxfile.Entry // sorted by their offset
}

// NewIndex creates a new empty index with the given size. Size is a hint and
// can be 0. It is recommended to set it to the number of objects to be indexed
// if it is known beforehand (e.g. reading from a packfile).
func NewIndex(size int) *Index {
	return &Index{
		byHash:   make(map[plumbing.Hash]*idxfile.Entry, size),
		byOffset: make([]*idxfile.Entry, 0, size),
	}
}

// NewIndexFromIdxFile creates a new Index from an idxfile.IdxFile.
func NewIndexFromIdxFile(idxf *idxfile.Idxfile) *Index {
	idx := &Index{
		byHash:   make(map[plumbing.Hash]*idxfile.Entry, idxf.ObjectCount),
		byOffset: make([]*idxfile.Entry, 0, idxf.ObjectCount),
	}
	sorted := true
	for i, e := range idxf.Entries {
		idx.addUnsorted(e)
		if i > 0 && idx.byOffset[i-1].Offset >= e.Offset {
			sorted = false
		}
	}

	// If the idxfile was loaded from a regular packfile index
	// then it will already be in offset order, in which case we
	// can avoid doing a relatively expensive idempotent sort.
	if !sorted {
		sort.Sort(orderByOffset(idx.byOffset))
	}

	return idx
}

// orderByOffset is a sort.Interface adapter that arranges
// a slice of entries by their offset.
type orderByOffset []*idxfile.Entry

func (o orderByOffset) Len() int           { return len(o) }
func (o orderByOffset) Less(i, j int) bool { return o[i].Offset < o[j].Offset }
func (o orderByOffset) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

// Add adds a new Entry with the given values to the index.
func (idx *Index) Add(h plumbing.Hash, offset uint64, crc32 uint32) {
	e := &idxfile.Entry{
		Hash:   h,
		Offset: offset,
		CRC32:  crc32,
	}
	idx.byHash[e.Hash] = e

	// Find the right position in byOffset.
	// Look for the first position whose offset is *greater* than e.Offset.
	i := sort.Search(len(idx.byOffset), func(i int) bool {
		return idx.byOffset[i].Offset > offset
	})
	if i == len(idx.byOffset) {
		// Simple case: add it to the end.
		idx.byOffset = append(idx.byOffset, e)
		return
	}
	// Harder case: shift existing entries down by one to make room.
	// Append a nil entry first so we can use existing capacity in case
	// the index was carefully preallocated.
	idx.byOffset = append(idx.byOffset, nil)
	copy(idx.byOffset[i+1:], idx.byOffset[i:len(idx.byOffset)-1])
	idx.byOffset[i] = e
}

func (idx *Index) addUnsorted(e *idxfile.Entry) {
	idx.byHash[e.Hash] = e
	idx.byOffset = append(idx.byOffset, e)
}

// LookupHash looks an entry up by its hash. An idxfile.Entry is returned and
// a bool, which is true if it was found or false if it wasn't.
func (idx *Index) LookupHash(h plumbing.Hash) (*idxfile.Entry, bool) {
	e, ok := idx.byHash[h]
	return e, ok
}

// LookupHash looks an entry up by its offset in the packfile. An idxfile.Entry
// is returned and a bool, which is true if it was found or false if it wasn't.
func (idx *Index) LookupOffset(offset uint64) (*idxfile.Entry, bool) {
	i := sort.Search(len(idx.byOffset), func(i int) bool {
		return idx.byOffset[i].Offset >= offset
	})
	if i >= len(idx.byOffset) || idx.byOffset[i].Offset != offset {
		return nil, false // not present
	}
	return idx.byOffset[i], true
}

// Size returns the number of entries in the index.
func (idx *Index) Size() int {
	return len(idx.byHash)
}

// ToIdxFile converts the index to an idxfile.Idxfile, which can then be used
// to serialize.
func (idx *Index) ToIdxFile() *idxfile.Idxfile {
	idxf := idxfile.NewIdxfile()
	for _, e := range idx.byHash {
		idxf.Entries = append(idxf.Entries, e)
	}

	return idxf
}
