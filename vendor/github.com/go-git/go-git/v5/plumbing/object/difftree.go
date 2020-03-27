package object

import (
	"bytes"
	"context"

	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/go-git/go-git/v5/utils/merkletrie/noder"
)

// DiffTree compares the content and mode of the blobs found via two
// tree objects.
func DiffTree(a, b *Tree) (Changes, error) {
	return DiffTreeContext(context.Background(), a, b)
}

// DiffTree compares the content and mode of the blobs found via two
// tree objects. Provided context must be non-nil.
// An error will be return if context expires
func DiffTreeContext(ctx context.Context, a, b *Tree) (Changes, error) {
	from := NewTreeRootNode(a)
	to := NewTreeRootNode(b)

	hashEqual := func(a, b noder.Hasher) bool {
		return bytes.Equal(a.Hash(), b.Hash())
	}

	merkletrieChanges, err := merkletrie.DiffTreeContext(ctx, from, to, hashEqual)
	if err != nil {
		if err == merkletrie.ErrCanceled {
			return nil, ErrCanceled
		}
		return nil, err
	}

	return newChanges(merkletrieChanges)
}
