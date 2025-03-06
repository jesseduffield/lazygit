// Package revlist provides support to access the ancestors of commits, in a
// similar way as the git-rev-list command.
package revlist

import (
	"fmt"
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Objects applies a complementary set. It gets all the hashes from all
// the reachable objects from the given objects. Ignore param are object hashes
// that we want to ignore on the result. All that objects must be accessible
// from the object storer.
func Objects(
	s storer.EncodedObjectStorer,
	objs,
	ignore []plumbing.Hash,
) ([]plumbing.Hash, error) {
	return ObjectsWithStorageForIgnores(s, s, objs, ignore)
}

// ObjectsWithStorageForIgnores is the same as Objects, but a
// secondary storage layer can be provided, to be used to finding the
// full set of objects to be ignored while finding the reachable
// objects.  This is useful when the main `s` storage layer is slow
// and/or remote, while the ignore list is available somewhere local.
func ObjectsWithStorageForIgnores(
	s, ignoreStore storer.EncodedObjectStorer,
	objs,
	ignore []plumbing.Hash,
) ([]plumbing.Hash, error) {
	ignore, err := objects(ignoreStore, ignore, nil, true)
	if err != nil {
		return nil, err
	}

	return objects(s, objs, ignore, false)
}

func objects(
	s storer.EncodedObjectStorer,
	objects,
	ignore []plumbing.Hash,
	allowMissingObjects bool,
) ([]plumbing.Hash, error) {
	seen := hashListToSet(ignore)
	result := make(map[plumbing.Hash]bool)
	visited := make(map[plumbing.Hash]bool)

	walkerFunc := func(h plumbing.Hash) {
		if !seen[h] {
			result[h] = true
			seen[h] = true
		}
	}

	for _, h := range objects {
		if err := processObject(s, h, seen, visited, ignore, walkerFunc); err != nil {
			if allowMissingObjects && err == plumbing.ErrObjectNotFound {
				continue
			}

			return nil, err
		}
	}

	return hashSetToList(result), nil
}

// processObject obtains the object using the hash an process it depending of its type
func processObject(
	s storer.EncodedObjectStorer,
	h plumbing.Hash,
	seen map[plumbing.Hash]bool,
	visited map[plumbing.Hash]bool,
	ignore []plumbing.Hash,
	walkerFunc func(h plumbing.Hash),
) error {
	if seen[h] {
		return nil
	}

	o, err := s.EncodedObject(plumbing.AnyObject, h)
	if err != nil {
		return err
	}

	do, err := object.DecodeObject(s, o)
	if err != nil {
		return err
	}

	switch do := do.(type) {
	case *object.Commit:
		return reachableObjects(do, seen, visited, ignore, walkerFunc)
	case *object.Tree:
		return iterateCommitTrees(seen, do, walkerFunc)
	case *object.Tag:
		walkerFunc(do.Hash)
		return processObject(s, do.Target, seen, visited, ignore, walkerFunc)
	case *object.Blob:
		walkerFunc(do.Hash)
	default:
		return fmt.Errorf("object type not valid: %s. "+
			"Object reference: %s", o.Type(), o.Hash())
	}

	return nil
}

// reachableObjects returns, using the callback function, all the reachable
// objects from the specified commit. To avoid to iterate over seen commits,
// if a commit hash is into the 'seen' set, we will not iterate all his trees
// and blobs objects.
func reachableObjects(
	commit *object.Commit,
	seen map[plumbing.Hash]bool,
	visited map[plumbing.Hash]bool,
	ignore []plumbing.Hash,
	cb func(h plumbing.Hash),
) error {
	i := object.NewCommitPreorderIter(commit, seen, ignore)
	pending := make(map[plumbing.Hash]bool)
	addPendingParents(pending, visited, commit)
	for {
		commit, err := i.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if pending[commit.Hash] {
			delete(pending, commit.Hash)
		}

		addPendingParents(pending, visited, commit)

		if visited[commit.Hash] && len(pending) == 0 {
			break
		}

		if seen[commit.Hash] {
			continue
		}

		cb(commit.Hash)

		tree, err := commit.Tree()
		if err != nil {
			return err
		}

		if err := iterateCommitTrees(seen, tree, cb); err != nil {
			return err
		}
	}

	return nil
}

func addPendingParents(pending, visited map[plumbing.Hash]bool, commit *object.Commit) {
	for _, p := range commit.ParentHashes {
		if !visited[p] {
			pending[p] = true
		}
	}
}

// iterateCommitTrees iterate all reachable trees from the given commit
func iterateCommitTrees(
	seen map[plumbing.Hash]bool,
	tree *object.Tree,
	cb func(h plumbing.Hash),
) error {
	if seen[tree.Hash] {
		return nil
	}

	cb(tree.Hash)

	treeWalker := object.NewTreeWalker(tree, true, seen)

	for {
		_, e, err := treeWalker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if e.Mode == filemode.Submodule {
			continue
		}

		if seen[e.Hash] {
			continue
		}

		cb(e.Hash)
	}

	return nil
}

func hashSetToList(hashes map[plumbing.Hash]bool) []plumbing.Hash {
	var result []plumbing.Hash
	for key := range hashes {
		result = append(result, key)
	}

	return result
}

func hashListToSet(hashes []plumbing.Hash) map[plumbing.Hash]bool {
	result := make(map[plumbing.Hash]bool)
	for _, h := range hashes {
		result[h] = true
	}

	return result
}
