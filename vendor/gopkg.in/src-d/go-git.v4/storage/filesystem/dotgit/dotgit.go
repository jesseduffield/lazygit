// https://github.com/git/git/blob/master/Documentation/gitrepository-layout.txt
package dotgit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	stdioutil "io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/utils/ioutil"

	"gopkg.in/src-d/go-billy.v4"
)

const (
	suffix         = ".git"
	packedRefsPath = "packed-refs"
	configPath     = "config"
	indexPath      = "index"
	shallowPath    = "shallow"
	modulePath     = "modules"
	objectsPath    = "objects"
	packPath       = "pack"
	refsPath       = "refs"

	tmpPackedRefsPrefix = "._packed-refs"

	packExt = ".pack"
	idxExt  = ".idx"
)

var (
	// ErrNotFound is returned by New when the path is not found.
	ErrNotFound = errors.New("path not found")
	// ErrIdxNotFound is returned by Idxfile when the idx file is not found
	ErrIdxNotFound = errors.New("idx file not found")
	// ErrPackfileNotFound is returned by Packfile when the packfile is not found
	ErrPackfileNotFound = errors.New("packfile not found")
	// ErrConfigNotFound is returned by Config when the config is not found
	ErrConfigNotFound = errors.New("config file not found")
	// ErrPackedRefsDuplicatedRef is returned when a duplicated reference is
	// found in the packed-ref file. This is usually the case for corrupted git
	// repositories.
	ErrPackedRefsDuplicatedRef = errors.New("duplicated ref found in packed-ref file")
	// ErrPackedRefsBadFormat is returned when the packed-ref file corrupt.
	ErrPackedRefsBadFormat = errors.New("malformed packed-ref")
	// ErrSymRefTargetNotFound is returned when a symbolic reference is
	// targeting a non-existing object. This usually means the repository
	// is corrupt.
	ErrSymRefTargetNotFound = errors.New("symbolic reference target not found")
)

// The DotGit type represents a local git repository on disk. This
// type is not zero-value-safe, use the New function to initialize it.
type DotGit struct {
	fs billy.Filesystem
}

// New returns a DotGit value ready to be used. The path argument must
// be the absolute path of a git repository directory (e.g.
// "/foo/bar/.git").
func New(fs billy.Filesystem) *DotGit {
	return &DotGit{fs: fs}
}

// Initialize creates all the folder scaffolding.
func (d *DotGit) Initialize() error {
	mustExists := []string{
		d.fs.Join("objects", "info"),
		d.fs.Join("objects", "pack"),
		d.fs.Join("refs", "heads"),
		d.fs.Join("refs", "tags"),
	}

	for _, path := range mustExists {
		_, err := d.fs.Stat(path)
		if err == nil {
			continue
		}

		if !os.IsNotExist(err) {
			return err
		}

		if err := d.fs.MkdirAll(path, os.ModeDir|os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// ConfigWriter returns a file pointer for write to the config file
func (d *DotGit) ConfigWriter() (billy.File, error) {
	return d.fs.Create(configPath)
}

// Config returns a file pointer for read to the config file
func (d *DotGit) Config() (billy.File, error) {
	return d.fs.Open(configPath)
}

// IndexWriter returns a file pointer for write to the index file
func (d *DotGit) IndexWriter() (billy.File, error) {
	return d.fs.Create(indexPath)
}

// Index returns a file pointer for read to the index file
func (d *DotGit) Index() (billy.File, error) {
	return d.fs.Open(indexPath)
}

// ShallowWriter returns a file pointer for write to the shallow file
func (d *DotGit) ShallowWriter() (billy.File, error) {
	return d.fs.Create(shallowPath)
}

// Shallow returns a file pointer for read to the shallow file
func (d *DotGit) Shallow() (billy.File, error) {
	f, err := d.fs.Open(shallowPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	return f, nil
}

// NewObjectPack return a writer for a new packfile, it saves the packfile to
// disk and also generates and save the index for the given packfile.
func (d *DotGit) NewObjectPack() (*PackWriter, error) {
	return newPackWrite(d.fs)
}

// ObjectPacks returns the list of availables packfiles
func (d *DotGit) ObjectPacks() ([]plumbing.Hash, error) {
	packDir := d.fs.Join(objectsPath, packPath)
	files, err := d.fs.ReadDir(packDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var packs []plumbing.Hash
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), packExt) {
			continue
		}

		n := f.Name()
		h := plumbing.NewHash(n[5 : len(n)-5]) //pack-(hash).pack
		if h.IsZero() {
			// Ignore files with badly-formatted names.
			continue
		}
		packs = append(packs, h)
	}

	return packs, nil
}

func (d *DotGit) objectPackPath(hash plumbing.Hash, extension string) string {
	return d.fs.Join(objectsPath, packPath, fmt.Sprintf("pack-%s.%s", hash.String(), extension))
}

func (d *DotGit) objectPackOpen(hash plumbing.Hash, extension string) (billy.File, error) {
	pack, err := d.fs.Open(d.objectPackPath(hash, extension))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrPackfileNotFound
		}

		return nil, err
	}

	return pack, nil
}

// ObjectPack returns a fs.File of the given packfile
func (d *DotGit) ObjectPack(hash plumbing.Hash) (billy.File, error) {
	return d.objectPackOpen(hash, `pack`)
}

// ObjectPackIdx returns a fs.File of the index file for a given packfile
func (d *DotGit) ObjectPackIdx(hash plumbing.Hash) (billy.File, error) {
	return d.objectPackOpen(hash, `idx`)
}

func (d *DotGit) DeleteOldObjectPackAndIndex(hash plumbing.Hash, t time.Time) error {
	path := d.objectPackPath(hash, `pack`)
	if !t.IsZero() {
		fi, err := d.fs.Stat(path)
		if err != nil {
			return err
		}
		// too new, skip deletion.
		if !fi.ModTime().Before(t) {
			return nil
		}
	}
	err := d.fs.Remove(path)
	if err != nil {
		return err
	}
	return d.fs.Remove(d.objectPackPath(hash, `idx`))
}

// NewObject return a writer for a new object file.
func (d *DotGit) NewObject() (*ObjectWriter, error) {
	return newObjectWriter(d.fs)
}

// Objects returns a slice with the hashes of objects found under the
// .git/objects/ directory.
func (d *DotGit) Objects() ([]plumbing.Hash, error) {
	var objects []plumbing.Hash
	err := d.ForEachObjectHash(func(hash plumbing.Hash) error {
		objects = append(objects, hash)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return objects, nil
}

// Objects returns a slice with the hashes of objects found under the
// .git/objects/ directory.
func (d *DotGit) ForEachObjectHash(fun func(plumbing.Hash) error) error {
	files, err := d.fs.ReadDir(objectsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, f := range files {
		if f.IsDir() && len(f.Name()) == 2 && isHex(f.Name()) {
			base := f.Name()
			d, err := d.fs.ReadDir(d.fs.Join(objectsPath, base))
			if err != nil {
				return err
			}

			for _, o := range d {
				h := plumbing.NewHash(base + o.Name())
				if h.IsZero() {
					// Ignore files with badly-formatted names.
					continue
				}
				err = fun(h)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (d *DotGit) objectPath(h plumbing.Hash) string {
	hash := h.String()
	return d.fs.Join(objectsPath, hash[0:2], hash[2:40])
}

// Object returns a fs.File pointing the object file, if exists
func (d *DotGit) Object(h plumbing.Hash) (billy.File, error) {
	return d.fs.Open(d.objectPath(h))
}

// ObjectStat returns a os.FileInfo pointing the object file, if exists
func (d *DotGit) ObjectStat(h plumbing.Hash) (os.FileInfo, error) {
	return d.fs.Stat(d.objectPath(h))
}

// ObjectDelete removes the object file, if exists
func (d *DotGit) ObjectDelete(h plumbing.Hash) error {
	return d.fs.Remove(d.objectPath(h))
}

func (d *DotGit) readReferenceFrom(rd io.Reader, name string) (ref *plumbing.Reference, err error) {
	b, err := stdioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	line := strings.TrimSpace(string(b))
	return plumbing.NewReferenceFromStrings(name, line), nil
}

func (d *DotGit) checkReferenceAndTruncate(f billy.File, old *plumbing.Reference) error {
	if old == nil {
		return nil
	}
	ref, err := d.readReferenceFrom(f, old.Name().String())
	if err != nil {
		return err
	}
	if ref.Hash() != old.Hash() {
		return fmt.Errorf("reference has changed concurrently")
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return f.Truncate(0)
}

func (d *DotGit) SetRef(r, old *plumbing.Reference) error {
	var content string
	switch r.Type() {
	case plumbing.SymbolicReference:
		content = fmt.Sprintf("ref: %s\n", r.Target())
	case plumbing.HashReference:
		content = fmt.Sprintln(r.Hash().String())
	}

	fileName := r.Name().String()

	return d.setRef(fileName, content, old)
}

// Refs scans the git directory collecting references, which it returns.
// Symbolic references are resolved and included in the output.
func (d *DotGit) Refs() ([]*plumbing.Reference, error) {
	var refs []*plumbing.Reference
	var seen = make(map[plumbing.ReferenceName]bool)
	if err := d.addRefsFromRefDir(&refs, seen); err != nil {
		return nil, err
	}

	if err := d.addRefsFromPackedRefs(&refs, seen); err != nil {
		return nil, err
	}

	if err := d.addRefFromHEAD(&refs); err != nil {
		return nil, err
	}

	return refs, nil
}

// Ref returns the reference for a given reference name.
func (d *DotGit) Ref(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	ref, err := d.readReferenceFile(".", name.String())
	if err == nil {
		return ref, nil
	}

	return d.packedRef(name)
}

func (d *DotGit) findPackedRefsInFile(f billy.File) ([]*plumbing.Reference, error) {
	s := bufio.NewScanner(f)
	var refs []*plumbing.Reference
	for s.Scan() {
		ref, err := d.processLine(s.Text())
		if err != nil {
			return nil, err
		}

		if ref != nil {
			refs = append(refs, ref)
		}
	}

	return refs, s.Err()
}

func (d *DotGit) findPackedRefs() (r []*plumbing.Reference, err error) {
	f, err := d.fs.Open(packedRefsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	defer ioutil.CheckClose(f, &err)
	return d.findPackedRefsInFile(f)
}

func (d *DotGit) packedRef(name plumbing.ReferenceName) (*plumbing.Reference, error) {
	refs, err := d.findPackedRefs()
	if err != nil {
		return nil, err
	}

	for _, ref := range refs {
		if ref.Name() == name {
			return ref, nil
		}
	}

	return nil, plumbing.ErrReferenceNotFound
}

// RemoveRef removes a reference by name.
func (d *DotGit) RemoveRef(name plumbing.ReferenceName) error {
	path := d.fs.Join(".", name.String())
	_, err := d.fs.Stat(path)
	if err == nil {
		err = d.fs.Remove(path)
		// Drop down to remove it from the packed refs file, too.
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return d.rewritePackedRefsWithoutRef(name)
}

func (d *DotGit) addRefsFromPackedRefs(refs *[]*plumbing.Reference, seen map[plumbing.ReferenceName]bool) (err error) {
	packedRefs, err := d.findPackedRefs()
	if err != nil {
		return err
	}

	for _, ref := range packedRefs {
		if !seen[ref.Name()] {
			*refs = append(*refs, ref)
			seen[ref.Name()] = true
		}
	}
	return nil
}

func (d *DotGit) addRefsFromPackedRefsFile(refs *[]*plumbing.Reference, f billy.File, seen map[plumbing.ReferenceName]bool) (err error) {
	packedRefs, err := d.findPackedRefsInFile(f)
	if err != nil {
		return err
	}

	for _, ref := range packedRefs {
		if !seen[ref.Name()] {
			*refs = append(*refs, ref)
			seen[ref.Name()] = true
		}
	}
	return nil
}

func (d *DotGit) openAndLockPackedRefs(doCreate bool) (
	pr billy.File, err error) {
	var f billy.File
	defer func() {
		if err != nil && f != nil {
			ioutil.CheckClose(f, &err)
		}
	}()

	// File mode is retrieved from a constant defined in the target specific
	// files (dotgit_rewrite_packed_refs_*). Some modes are not available
	// in all filesystems.
	openFlags := d.openAndLockPackedRefsMode()
	if doCreate {
		openFlags |= os.O_CREATE
	}

	// Keep trying to open and lock the file until we're sure the file
	// didn't change between the open and the lock.
	for {
		f, err = d.fs.OpenFile(packedRefsPath, openFlags, 0600)
		if err != nil {
			if os.IsNotExist(err) && !doCreate {
				return nil, nil
			}

			return nil, err
		}
		fi, err := d.fs.Stat(packedRefsPath)
		if err != nil {
			return nil, err
		}
		mtime := fi.ModTime()

		err = f.Lock()
		if err != nil {
			return nil, err
		}

		fi, err = d.fs.Stat(packedRefsPath)
		if err != nil {
			return nil, err
		}
		if mtime.Equal(fi.ModTime()) {
			break
		}
		// The file has changed since we opened it.  Close and retry.
		err = f.Close()
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

func (d *DotGit) rewritePackedRefsWithoutRef(name plumbing.ReferenceName) (err error) {
	pr, err := d.openAndLockPackedRefs(false)
	if err != nil {
		return err
	}
	if pr == nil {
		return nil
	}
	defer ioutil.CheckClose(pr, &err)

	// Creating the temp file in the same directory as the target file
	// improves our chances for rename operation to be atomic.
	tmp, err := d.fs.TempFile("", tmpPackedRefsPrefix)
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		ioutil.CheckClose(tmp, &err)
		_ = d.fs.Remove(tmpName) // don't check err, we might have renamed it
	}()

	s := bufio.NewScanner(pr)
	found := false
	for s.Scan() {
		line := s.Text()
		ref, err := d.processLine(line)
		if err != nil {
			return err
		}

		if ref != nil && ref.Name() == name {
			found = true
			continue
		}

		if _, err := fmt.Fprintln(tmp, line); err != nil {
			return err
		}
	}

	if err := s.Err(); err != nil {
		return err
	}

	if !found {
		return nil
	}

	return d.rewritePackedRefsWhileLocked(tmp, pr)
}

// process lines from a packed-refs file
func (d *DotGit) processLine(line string) (*plumbing.Reference, error) {
	if len(line) == 0 {
		return nil, nil
	}

	switch line[0] {
	case '#': // comment - ignore
		return nil, nil
	case '^': // annotated tag commit of the previous line - ignore
		return nil, nil
	default:
		ws := strings.Split(line, " ") // hash then ref
		if len(ws) != 2 {
			return nil, ErrPackedRefsBadFormat
		}

		return plumbing.NewReferenceFromStrings(ws[1], ws[0]), nil
	}
}

func (d *DotGit) addRefsFromRefDir(refs *[]*plumbing.Reference, seen map[plumbing.ReferenceName]bool) error {
	return d.walkReferencesTree(refs, []string{refsPath}, seen)
}

func (d *DotGit) walkReferencesTree(refs *[]*plumbing.Reference, relPath []string, seen map[plumbing.ReferenceName]bool) error {
	files, err := d.fs.ReadDir(d.fs.Join(relPath...))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, f := range files {
		newRelPath := append(append([]string(nil), relPath...), f.Name())
		if f.IsDir() {
			if err = d.walkReferencesTree(refs, newRelPath, seen); err != nil {
				return err
			}

			continue
		}

		ref, err := d.readReferenceFile(".", strings.Join(newRelPath, "/"))
		if err != nil {
			return err
		}

		if ref != nil && !seen[ref.Name()] {
			*refs = append(*refs, ref)
			seen[ref.Name()] = true
		}
	}

	return nil
}

func (d *DotGit) addRefFromHEAD(refs *[]*plumbing.Reference) error {
	ref, err := d.readReferenceFile(".", "HEAD")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	*refs = append(*refs, ref)
	return nil
}

func (d *DotGit) readReferenceFile(path, name string) (ref *plumbing.Reference, err error) {
	path = d.fs.Join(path, d.fs.Join(strings.Split(name, "/")...))
	f, err := d.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer ioutil.CheckClose(f, &err)

	return d.readReferenceFrom(f, name)
}

func (d *DotGit) CountLooseRefs() (int, error) {
	var refs []*plumbing.Reference
	var seen = make(map[plumbing.ReferenceName]bool)
	if err := d.addRefsFromRefDir(&refs, seen); err != nil {
		return 0, err
	}

	return len(refs), nil
}

// PackRefs packs all loose refs into the packed-refs file.
//
// This implementation only works under the assumption that the view
// of the file system won't be updated during this operation.  This
// strategy would not work on a general file system though, without
// locking each loose reference and checking it again before deleting
// the file, because otherwise an updated reference could sneak in and
// then be deleted by the packed-refs process.  Alternatively, every
// ref update could also lock packed-refs, so only one lock is
// required during ref-packing.  But that would worsen performance in
// the common case.
//
// TODO: add an "all" boolean like the `git pack-refs --all` flag.
// When `all` is false, it would only pack refs that have already been
// packed, plus all tags.
func (d *DotGit) PackRefs() (err error) {
	// Lock packed-refs, and create it if it doesn't exist yet.
	f, err := d.openAndLockPackedRefs(true)
	if err != nil {
		return err
	}
	defer ioutil.CheckClose(f, &err)

	// Gather all refs using addRefsFromRefDir and addRefsFromPackedRefs.
	var refs []*plumbing.Reference
	seen := make(map[plumbing.ReferenceName]bool)
	if err = d.addRefsFromRefDir(&refs, seen); err != nil {
		return err
	}
	if len(refs) == 0 {
		// Nothing to do!
		return nil
	}
	numLooseRefs := len(refs)
	if err = d.addRefsFromPackedRefsFile(&refs, f, seen); err != nil {
		return err
	}

	// Write them all to a new temp packed-refs file.
	tmp, err := d.fs.TempFile("", tmpPackedRefsPrefix)
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		ioutil.CheckClose(tmp, &err)
		_ = d.fs.Remove(tmpName) // don't check err, we might have renamed it
	}()

	w := bufio.NewWriter(tmp)
	for _, ref := range refs {
		_, err = w.WriteString(ref.String() + "\n")
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	// Rename the temp packed-refs file.
	err = d.rewritePackedRefsWhileLocked(tmp, f)
	if err != nil {
		return err
	}

	// Delete all the loose refs, while still holding the packed-refs
	// lock.
	for _, ref := range refs[:numLooseRefs] {
		path := d.fs.Join(".", ref.Name().String())
		err = d.fs.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// Module return a billy.Filesystem pointing to the module folder
func (d *DotGit) Module(name string) (billy.Filesystem, error) {
	return d.fs.Chroot(d.fs.Join(modulePath, name))
}

// Alternates returns DotGit(s) based off paths in objects/info/alternates if
// available. This can be used to checks if it's a shared repository.
func (d *DotGit) Alternates() ([]*DotGit, error) {
	altpath := d.fs.Join("objects", "info", "alternates")
	f, err := d.fs.Open(altpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var alternates []*DotGit

	// Read alternate paths line-by-line and create DotGit objects.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		path := scanner.Text()
		if !filepath.IsAbs(path) {
			// For relative paths, we can perform an internal conversion to
			// slash so that they work cross-platform.
			slashPath := filepath.ToSlash(path)
			// If the path is not absolute, it must be relative to object
			// database (.git/objects/info).
			// https://www.kernel.org/pub/software/scm/git/docs/gitrepository-layout.html
			// Hence, derive a path relative to DotGit's root.
			// "../../../reponame/.git/" -> "../../reponame/.git"
			// Remove the first ../
			relpath := filepath.Join(strings.Split(slashPath, "/")[1:]...)
			normalPath := filepath.FromSlash(relpath)
			path = filepath.Join(d.fs.Root(), normalPath)
		}
		fs := osfs.New(filepath.Dir(path))
		alternates = append(alternates, New(fs))
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return alternates, nil
}

func isHex(s string) bool {
	for _, b := range []byte(s) {
		if isNum(b) {
			continue
		}
		if isHexAlpha(b) {
			continue
		}

		return false
	}

	return true
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9'
}

func isHexAlpha(b byte) bool {
	return b >= 'a' && b <= 'f' || b >= 'A' && b <= 'F'
}
