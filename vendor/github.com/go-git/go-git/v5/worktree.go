package git

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/utils/ioutil"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/go-git/go-git/v5/utils/sync"
)

var (
	ErrWorktreeNotClean                = errors.New("worktree is not clean")
	ErrSubmoduleNotFound               = errors.New("submodule not found")
	ErrUnstagedChanges                 = errors.New("worktree contains unstaged changes")
	ErrGitModulesSymlink               = errors.New(gitmodulesFile + " is a symlink")
	ErrNonFastForwardUpdate            = errors.New("non-fast-forward update")
	ErrRestoreWorktreeOnlyNotSupported = errors.New("worktree only is not supported")
)

// Worktree represents a git worktree.
type Worktree struct {
	// Filesystem underlying filesystem.
	Filesystem billy.Filesystem
	// External excludes not found in the repository .gitignore
	Excludes []gitignore.Pattern

	r *Repository
}

// Pull incorporates changes from a remote repository into the current branch.
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
//
// Pull only supports merges where the can be resolved as a fast-forward.
func (w *Worktree) Pull(o *PullOptions) error {
	return w.PullContext(context.Background(), o)
}

// PullContext incorporates changes from a remote repository into the current
// branch. Returns nil if the operation is successful, NoErrAlreadyUpToDate if
// there are no changes to be fetched, or an error.
//
// Pull only supports merges where the can be resolved as a fast-forward.
//
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects the
// transport operations.
func (w *Worktree) PullContext(ctx context.Context, o *PullOptions) error {
	if err := o.Validate(); err != nil {
		return err
	}

	remote, err := w.r.Remote(o.RemoteName)
	if err != nil {
		return err
	}

	fetchHead, err := remote.fetch(ctx, &FetchOptions{
		RemoteName:      o.RemoteName,
		RemoteURL:       o.RemoteURL,
		Depth:           o.Depth,
		Auth:            o.Auth,
		Progress:        o.Progress,
		Force:           o.Force,
		InsecureSkipTLS: o.InsecureSkipTLS,
		CABundle:        o.CABundle,
		ProxyOptions:    o.ProxyOptions,
	})

	updated := true
	if err == NoErrAlreadyUpToDate {
		updated = false
	} else if err != nil {
		return err
	}

	ref, err := storer.ResolveReference(fetchHead, o.ReferenceName)
	if err != nil {
		return err
	}

	head, err := w.r.Head()
	if err == nil {
		// if we don't have a shallows list, just ignore it
		shallowList, _ := w.r.Storer.Shallow()

		var earliestShallow *plumbing.Hash
		if len(shallowList) > 0 {
			earliestShallow = &shallowList[0]
		}

		headAheadOfRef, err := isFastForward(w.r.Storer, ref.Hash(), head.Hash(), earliestShallow)
		if err != nil {
			return err
		}

		if !updated && headAheadOfRef {
			return NoErrAlreadyUpToDate
		}

		ff, err := isFastForward(w.r.Storer, head.Hash(), ref.Hash(), earliestShallow)
		if err != nil {
			return err
		}

		if !ff {
			return ErrNonFastForwardUpdate
		}
	}

	if err != nil && err != plumbing.ErrReferenceNotFound {
		return err
	}

	if err := w.updateHEAD(ref.Hash()); err != nil {
		return err
	}

	if err := w.Reset(&ResetOptions{
		Mode:   MergeReset,
		Commit: ref.Hash(),
	}); err != nil {
		return err
	}

	if o.RecurseSubmodules != NoRecurseSubmodules {
		return w.updateSubmodules(ctx, &SubmoduleUpdateOptions{
			RecurseSubmodules: o.RecurseSubmodules,
			Auth:              o.Auth,
		})
	}

	return nil
}

func (w *Worktree) updateSubmodules(ctx context.Context, o *SubmoduleUpdateOptions) error {
	s, err := w.Submodules()
	if err != nil {
		return err
	}
	o.Init = true
	return s.UpdateContext(ctx, o)
}

// Checkout switch branches or restore working tree files.
func (w *Worktree) Checkout(opts *CheckoutOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	if opts.Create {
		if err := w.createBranch(opts); err != nil {
			return err
		}
	}

	c, err := w.getCommitFromCheckoutOptions(opts)
	if err != nil {
		return err
	}

	ro := &ResetOptions{Commit: c, Mode: MergeReset}
	if opts.Force {
		ro.Mode = HardReset
	} else if opts.Keep {
		ro.Mode = SoftReset
	}

	if !opts.Hash.IsZero() && !opts.Create {
		err = w.setHEADToCommit(opts.Hash)
	} else {
		err = w.setHEADToBranch(opts.Branch, c)
	}

	if err != nil {
		return err
	}

	if len(opts.SparseCheckoutDirectories) > 0 {
		return w.ResetSparsely(ro, opts.SparseCheckoutDirectories)
	}

	return w.Reset(ro)
}

func (w *Worktree) createBranch(opts *CheckoutOptions) error {
	if err := opts.Branch.Validate(); err != nil {
		return err
	}

	_, err := w.r.Storer.Reference(opts.Branch)
	if err == nil {
		return fmt.Errorf("a branch named %q already exists", opts.Branch)
	}

	if err != plumbing.ErrReferenceNotFound {
		return err
	}

	if opts.Hash.IsZero() {
		ref, err := w.r.Head()
		if err != nil {
			return err
		}

		opts.Hash = ref.Hash()
	}

	return w.r.Storer.SetReference(
		plumbing.NewHashReference(opts.Branch, opts.Hash),
	)
}

func (w *Worktree) getCommitFromCheckoutOptions(opts *CheckoutOptions) (plumbing.Hash, error) {
	hash := opts.Hash
	if hash.IsZero() {
		b, err := w.r.Reference(opts.Branch, true)
		if err != nil {
			return plumbing.ZeroHash, err
		}

		hash = b.Hash()
	}

	o, err := w.r.Object(plumbing.AnyObject, hash)
	if err != nil {
		return plumbing.ZeroHash, err
	}

	switch o := o.(type) {
	case *object.Tag:
		if o.TargetType != plumbing.CommitObject {
			return plumbing.ZeroHash, fmt.Errorf("%w: tag target %q", object.ErrUnsupportedObject, o.TargetType)
		}

		return o.Target, nil
	case *object.Commit:
		return o.Hash, nil
	}

	return plumbing.ZeroHash, fmt.Errorf("%w: %q", object.ErrUnsupportedObject, o.Type())
}

func (w *Worktree) setHEADToCommit(commit plumbing.Hash) error {
	head := plumbing.NewHashReference(plumbing.HEAD, commit)
	return w.r.Storer.SetReference(head)
}

func (w *Worktree) setHEADToBranch(branch plumbing.ReferenceName, commit plumbing.Hash) error {
	target, err := w.r.Storer.Reference(branch)
	if err != nil {
		return err
	}

	var head *plumbing.Reference
	if target.Name().IsBranch() {
		head = plumbing.NewSymbolicReference(plumbing.HEAD, target.Name())
	} else {
		head = plumbing.NewHashReference(plumbing.HEAD, commit)
	}

	return w.r.Storer.SetReference(head)
}

func (w *Worktree) ResetSparsely(opts *ResetOptions, dirs []string) error {
	if err := opts.Validate(w.r); err != nil {
		return err
	}

	if opts.Mode == MergeReset {
		unstaged, err := w.containsUnstagedChanges()
		if err != nil {
			return err
		}

		if unstaged {
			return ErrUnstagedChanges
		}
	}

	if err := w.setHEADCommit(opts.Commit); err != nil {
		return err
	}

	if opts.Mode == SoftReset {
		return nil
	}

	t, err := w.r.getTreeFromCommitHash(opts.Commit)
	if err != nil {
		return err
	}

	if opts.Mode == MixedReset || opts.Mode == MergeReset || opts.Mode == HardReset {
		if err := w.resetIndex(t, dirs, opts.Files); err != nil {
			return err
		}
	}

	if opts.Mode == MergeReset || opts.Mode == HardReset {
		if err := w.resetWorktree(t, opts.Files); err != nil {
			return err
		}
	}

	return nil
}

// Restore restores specified files in the working tree or stage with contents from
// a restore source. If a path is tracked but does not exist in the restore,
// source, it will be removed to match the source.
//
// If Staged and Worktree are true, then the restore source will be the index.
// If only Staged is true, then the restore source will be HEAD.
// If only Worktree is true or neither Staged nor Worktree are true, will
// result in ErrRestoreWorktreeOnlyNotSupported because restoring the working
// tree while leaving the stage untouched is not currently supported.
//
// Restore with no files specified will return ErrNoRestorePaths.
func (w *Worktree) Restore(o *RestoreOptions) error {
	if err := o.Validate(); err != nil {
		return err
	}

	if o.Staged {
		opts := &ResetOptions{
			Files: o.Files,
		}

		if o.Worktree {
			// If we are doing both Worktree and Staging then it is a hard reset
			opts.Mode = HardReset
		} else {
			// If we are doing just staging then it is a mixed reset
			opts.Mode = MixedReset
		}

		return w.Reset(opts)
	}

	return ErrRestoreWorktreeOnlyNotSupported
}

// Reset the worktree to a specified state.
func (w *Worktree) Reset(opts *ResetOptions) error {
	return w.ResetSparsely(opts, nil)
}

func (w *Worktree) resetIndex(t *object.Tree, dirs []string, files []string) error {
	idx, err := w.r.Storer.Index()
	if err != nil {
		return err
	}

	b := newIndexBuilder(idx)

	changes, err := w.diffTreeWithStaging(t, true)
	if err != nil {
		return err
	}

	for _, ch := range changes {
		a, err := ch.Action()
		if err != nil {
			return err
		}

		var name string
		var e *object.TreeEntry

		switch a {
		case merkletrie.Modify, merkletrie.Insert:
			name = ch.To.String()
			e, err = t.FindEntry(name)
			if err != nil {
				return err
			}
		case merkletrie.Delete:
			name = ch.From.String()
		}

		if len(files) > 0 {
			contains := inFiles(files, name)
			if !contains {
				continue
			}
		}

		b.Remove(name)
		if e == nil {
			continue
		}

		b.Add(&index.Entry{
			Name: name,
			Hash: e.Hash,
			Mode: e.Mode,
		})

	}

	b.Write(idx)

	if len(dirs) > 0 {
		idx.SkipUnless(dirs)
	}

	return w.r.Storer.SetIndex(idx)
}

func inFiles(files []string, v string) bool {
	v = filepath.Clean(v)
	for _, s := range files {
		if filepath.Clean(s) == v {
			return true
		}
	}

	return false
}

func (w *Worktree) resetWorktree(t *object.Tree, files []string) error {
	changes, err := w.diffStagingWithWorktree(true, false)
	if err != nil {
		return err
	}

	idx, err := w.r.Storer.Index()
	if err != nil {
		return err
	}
	b := newIndexBuilder(idx)

	for _, ch := range changes {
		if err := w.validChange(ch); err != nil {
			return err
		}

		if len(files) > 0 {
			file := ""
			if ch.From != nil {
				file = ch.From.String()
			} else if ch.To != nil {
				file = ch.To.String()
			}

			if file == "" {
				continue
			}

			contains := inFiles(files, file)
			if !contains {
				continue
			}
		}

		if err := w.checkoutChange(ch, t, b); err != nil {
			return err
		}
	}

	b.Write(idx)
	return w.r.Storer.SetIndex(idx)
}

// worktreeDeny is a list of paths that are not allowed
// to be used when resetting the worktree.
var worktreeDeny = map[string]struct{}{
	// .git
	GitDirName: {},

	// For other historical reasons, file names that do not conform to the 8.3
	// format (up to eight characters for the basename, three for the file
	// extension, certain characters not allowed such as `+`, etc) are associated
	// with a so-called "short name", at least on the `C:` drive by default.
	// Which means that `git~1/` is a valid way to refer to `.git/`.
	"git~1": {},
}

// validPath checks whether paths are valid.
// The rules around invalid paths could differ from upstream based on how
// filesystems are managed within go-git, but they are largely the same.
//
// For upstream rules:
// https://github.com/git/git/blob/564d0252ca632e0264ed670534a51d18a689ef5d/read-cache.c#L946
// https://github.com/git/git/blob/564d0252ca632e0264ed670534a51d18a689ef5d/path.c#L1383
func validPath(paths ...string) error {
	for _, p := range paths {
		parts := strings.FieldsFunc(p, func(r rune) bool { return (r == '\\' || r == '/') })
		if len(parts) == 0 {
			return fmt.Errorf("invalid path: %q", p)
		}

		if _, denied := worktreeDeny[strings.ToLower(parts[0])]; denied {
			return fmt.Errorf("invalid path prefix: %q", p)
		}

		if runtime.GOOS == "windows" {
			// Volume names are not supported, in both formats: \\ and <DRIVE_LETTER>:.
			if vol := filepath.VolumeName(p); vol != "" {
				return fmt.Errorf("invalid path: %q", p)
			}

			if !windowsValidPath(parts[0]) {
				return fmt.Errorf("invalid path: %q", p)
			}
		}

		for _, part := range parts {
			if part == ".." {
				return fmt.Errorf("invalid path %q: cannot use '..'", p)
			}
		}
	}
	return nil
}

// windowsPathReplacer defines the chars that need to be replaced
// as part of windowsValidPath.
var windowsPathReplacer *strings.Replacer

func init() {
	windowsPathReplacer = strings.NewReplacer(" ", "", ".", "")
}

func windowsValidPath(part string) bool {
	if len(part) > 3 && strings.EqualFold(part[:4], GitDirName) {
		// For historical reasons, file names that end in spaces or periods are
		// automatically trimmed. Therefore, `.git . . ./` is a valid way to refer
		// to `.git/`.
		if windowsPathReplacer.Replace(part[4:]) == "" {
			return false
		}

		// For yet other historical reasons, NTFS supports so-called "Alternate Data
		// Streams", i.e. metadata associated with a given file, referred to via
		// `<filename>:<stream-name>:<stream-type>`. There exists a default stream
		// type for directories, allowing `.git/` to be accessed via
		// `.git::$INDEX_ALLOCATION/`.
		//
		// For performance reasons, _all_ Alternate Data Streams of `.git/` are
		// forbidden, not just `::$INDEX_ALLOCATION`.
		if len(part) > 4 && part[4:5] == ":" {
			return false
		}
	}
	return true
}

func (w *Worktree) validChange(ch merkletrie.Change) error {
	action, err := ch.Action()
	if err != nil {
		return nil
	}

	switch action {
	case merkletrie.Delete:
		return validPath(ch.From.String())
	case merkletrie.Insert:
		return validPath(ch.To.String())
	case merkletrie.Modify:
		return validPath(ch.From.String(), ch.To.String())
	}

	return nil
}

func (w *Worktree) checkoutChange(ch merkletrie.Change, t *object.Tree, idx *indexBuilder) error {
	a, err := ch.Action()
	if err != nil {
		return err
	}

	var e *object.TreeEntry
	var name string
	var isSubmodule bool

	switch a {
	case merkletrie.Modify, merkletrie.Insert:
		name = ch.To.String()
		e, err = t.FindEntry(name)
		if err != nil {
			return err
		}

		isSubmodule = e.Mode == filemode.Submodule
	case merkletrie.Delete:
		return rmFileAndDirsIfEmpty(w.Filesystem, ch.From.String())
	}

	if isSubmodule {
		return w.checkoutChangeSubmodule(name, a, e, idx)
	}

	return w.checkoutChangeRegularFile(name, a, t, e, idx)
}

func (w *Worktree) containsUnstagedChanges() (bool, error) {
	ch, err := w.diffStagingWithWorktree(false, true)
	if err != nil {
		return false, err
	}

	for _, c := range ch {
		a, err := c.Action()
		if err != nil {
			return false, err
		}

		if a == merkletrie.Insert {
			continue
		}

		return true, nil
	}

	return false, nil
}

func (w *Worktree) setHEADCommit(commit plumbing.Hash) error {
	head, err := w.r.Reference(plumbing.HEAD, false)
	if err != nil {
		return err
	}

	if head.Type() == plumbing.HashReference {
		head = plumbing.NewHashReference(plumbing.HEAD, commit)
		return w.r.Storer.SetReference(head)
	}

	branch, err := w.r.Reference(head.Target(), false)
	if err != nil {
		return err
	}

	if !branch.Name().IsBranch() {
		return fmt.Errorf("invalid HEAD target should be a branch, found %s", branch.Type())
	}

	branch = plumbing.NewHashReference(branch.Name(), commit)
	return w.r.Storer.SetReference(branch)
}

func (w *Worktree) checkoutChangeSubmodule(name string,
	a merkletrie.Action,
	e *object.TreeEntry,
	idx *indexBuilder,
) error {
	switch a {
	case merkletrie.Modify:
		sub, err := w.Submodule(name)
		if err != nil {
			return err
		}

		if !sub.initialized {
			return nil
		}

		return w.addIndexFromTreeEntry(name, e, idx)
	case merkletrie.Insert:
		mode, err := e.Mode.ToOSFileMode()
		if err != nil {
			return err
		}

		if err := w.Filesystem.MkdirAll(name, mode); err != nil {
			return err
		}

		return w.addIndexFromTreeEntry(name, e, idx)
	}

	return nil
}

func (w *Worktree) checkoutChangeRegularFile(name string,
	a merkletrie.Action,
	t *object.Tree,
	e *object.TreeEntry,
	idx *indexBuilder,
) error {
	switch a {
	case merkletrie.Modify:
		idx.Remove(name)

		// to apply perm changes the file is deleted, billy doesn't implement
		// chmod
		if err := w.Filesystem.Remove(name); err != nil {
			return err
		}

		fallthrough
	case merkletrie.Insert:
		f, err := t.File(name)
		if err != nil {
			return err
		}

		if err := w.checkoutFile(f); err != nil {
			return err
		}

		return w.addIndexFromFile(name, e.Hash, f.Mode, idx)
	}

	return nil
}

func (w *Worktree) checkoutFile(f *object.File) (err error) {
	mode, err := f.Mode.ToOSFileMode()
	if err != nil {
		return
	}

	if mode&os.ModeSymlink != 0 {
		return w.checkoutFileSymlink(f)
	}

	from, err := f.Reader()
	if err != nil {
		return
	}

	defer ioutil.CheckClose(from, &err)

	to, err := w.Filesystem.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
	if err != nil {
		return
	}

	defer ioutil.CheckClose(to, &err)
	buf := sync.GetByteSlice()
	_, err = io.CopyBuffer(to, from, *buf)
	sync.PutByteSlice(buf)
	return
}

func (w *Worktree) checkoutFileSymlink(f *object.File) (err error) {
	// https://github.com/git/git/commit/10ecfa76491e4923988337b2e2243b05376b40de
	if strings.EqualFold(f.Name, gitmodulesFile) {
		return ErrGitModulesSymlink
	}

	from, err := f.Reader()
	if err != nil {
		return
	}

	defer ioutil.CheckClose(from, &err)

	bytes, err := io.ReadAll(from)
	if err != nil {
		return
	}

	err = w.Filesystem.Symlink(string(bytes), f.Name)

	// On windows, this might fail.
	// Follow Git on Windows behavior by writing the link as it is.
	if err != nil && isSymlinkWindowsNonAdmin(err) {
		mode, _ := f.Mode.ToOSFileMode()

		to, err := w.Filesystem.OpenFile(f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode.Perm())
		if err != nil {
			return err
		}

		defer ioutil.CheckClose(to, &err)

		_, err = to.Write(bytes)
		return err
	}
	return
}

func (w *Worktree) addIndexFromTreeEntry(name string, f *object.TreeEntry, idx *indexBuilder) error {
	idx.Remove(name)
	idx.Add(&index.Entry{
		Hash: f.Hash,
		Name: name,
		Mode: filemode.Submodule,
	})
	return nil
}

func (w *Worktree) addIndexFromFile(name string, h plumbing.Hash, mode filemode.FileMode, idx *indexBuilder) error {
	idx.Remove(name)
	fi, err := w.Filesystem.Lstat(name)
	if err != nil {
		return err
	}

	e := &index.Entry{
		Hash:       h,
		Name:       name,
		Mode:       mode,
		ModifiedAt: fi.ModTime(),
		Size:       uint32(fi.Size()),
	}

	// if the FileInfo.Sys() comes from os the ctime, dev, inode, uid and gid
	// can be retrieved, otherwise this doesn't apply
	if fillSystemInfo != nil {
		fillSystemInfo(e, fi.Sys())
	}
	idx.Add(e)
	return nil
}

func (r *Repository) getTreeFromCommitHash(commit plumbing.Hash) (*object.Tree, error) {
	c, err := r.CommitObject(commit)
	if err != nil {
		return nil, err
	}

	return c.Tree()
}

var fillSystemInfo func(e *index.Entry, sys interface{})

const gitmodulesFile = ".gitmodules"

// Submodule returns the submodule with the given name
func (w *Worktree) Submodule(name string) (*Submodule, error) {
	l, err := w.Submodules()
	if err != nil {
		return nil, err
	}

	for _, m := range l {
		if m.Config().Name == name {
			return m, nil
		}
	}

	return nil, ErrSubmoduleNotFound
}

// Submodules returns all the available submodules
func (w *Worktree) Submodules() (Submodules, error) {
	l := make(Submodules, 0)
	m, err := w.readGitmodulesFile()
	if err != nil || m == nil {
		return l, err
	}

	c, err := w.r.Config()
	if err != nil {
		return nil, err
	}

	for _, s := range m.Submodules {
		l = append(l, w.newSubmodule(s, c.Submodules[s.Name]))
	}

	return l, nil
}

func (w *Worktree) newSubmodule(fromModules, fromConfig *config.Submodule) *Submodule {
	m := &Submodule{w: w}
	m.initialized = fromConfig != nil

	if !m.initialized {
		m.c = fromModules
		return m
	}

	m.c = fromConfig
	m.c.Path = fromModules.Path
	return m
}

func (w *Worktree) isSymlink(path string) bool {
	if s, err := w.Filesystem.Lstat(path); err == nil {
		return s.Mode()&os.ModeSymlink != 0
	}
	return false
}

func (w *Worktree) readGitmodulesFile() (*config.Modules, error) {
	if w.isSymlink(gitmodulesFile) {
		return nil, ErrGitModulesSymlink
	}

	f, err := w.Filesystem.Open(gitmodulesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	defer f.Close()
	input, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	m := config.NewModules()
	if err := m.Unmarshal(input); err != nil {
		return m, err
	}

	return m, nil
}

// Clean the worktree by removing untracked files.
// An empty dir could be removed - this is what  `git clean -f -d .` does.
func (w *Worktree) Clean(opts *CleanOptions) error {
	s, err := w.Status()
	if err != nil {
		return err
	}

	root := ""
	files, err := w.Filesystem.ReadDir(root)
	if err != nil {
		return err
	}
	return w.doClean(s, opts, root, files)
}

func (w *Worktree) doClean(status Status, opts *CleanOptions, dir string, files []os.FileInfo) error {
	for _, fi := range files {
		if fi.Name() == GitDirName {
			continue
		}

		// relative path under the root
		path := filepath.Join(dir, fi.Name())
		if fi.IsDir() {
			if !opts.Dir {
				continue
			}

			subfiles, err := w.Filesystem.ReadDir(path)
			if err != nil {
				return err
			}
			err = w.doClean(status, opts, path, subfiles)
			if err != nil {
				return err
			}
		} else {
			if status.IsUntracked(path) {
				if err := w.Filesystem.Remove(path); err != nil {
					return err
				}
			}
		}
	}

	if opts.Dir && dir != "" {
		_, err := removeDirIfEmpty(w.Filesystem, dir)
		return err
	}

	return nil
}

// GrepResult is structure of a grep result.
type GrepResult struct {
	// FileName is the name of file which contains match.
	FileName string
	// LineNumber is the line number of a file at which a match was found.
	LineNumber int
	// Content is the content of the file at the matching line.
	Content string
	// TreeName is the name of the tree (reference name/commit hash) at
	// which the match was performed.
	TreeName string
}

func (gr GrepResult) String() string {
	return fmt.Sprintf("%s:%s:%d:%s", gr.TreeName, gr.FileName, gr.LineNumber, gr.Content)
}

// Grep performs grep on a repository.
func (r *Repository) Grep(opts *GrepOptions) ([]GrepResult, error) {
	if err := opts.validate(r); err != nil {
		return nil, err
	}

	// Obtain commit hash from options (CommitHash or ReferenceName).
	var commitHash plumbing.Hash
	// treeName contains the value of TreeName in GrepResult.
	var treeName string

	if opts.ReferenceName != "" {
		ref, err := r.Reference(opts.ReferenceName, true)
		if err != nil {
			return nil, err
		}
		commitHash = ref.Hash()
		treeName = opts.ReferenceName.String()
	} else if !opts.CommitHash.IsZero() {
		commitHash = opts.CommitHash
		treeName = opts.CommitHash.String()
	}

	// Obtain a tree from the commit hash and get a tracked files iterator from
	// the tree.
	tree, err := r.getTreeFromCommitHash(commitHash)
	if err != nil {
		return nil, err
	}
	fileiter := tree.Files()

	return findMatchInFiles(fileiter, treeName, opts)
}

// Grep performs grep on a worktree.
func (w *Worktree) Grep(opts *GrepOptions) ([]GrepResult, error) {
	return w.r.Grep(opts)
}

// findMatchInFiles takes a FileIter, worktree name and GrepOptions, and
// returns a slice of GrepResult containing the result of regex pattern matching
// in content of all the files.
func findMatchInFiles(fileiter *object.FileIter, treeName string, opts *GrepOptions) ([]GrepResult, error) {
	var results []GrepResult

	err := fileiter.ForEach(func(file *object.File) error {
		var fileInPathSpec bool

		// When no pathspecs are provided, search all the files.
		if len(opts.PathSpecs) == 0 {
			fileInPathSpec = true
		}

		// Check if the file name matches with the pathspec. Break out of the
		// loop once a match is found.
		for _, pathSpec := range opts.PathSpecs {
			if pathSpec != nil && pathSpec.MatchString(file.Name) {
				fileInPathSpec = true
				break
			}
		}

		// If the file does not match with any of the pathspec, skip it.
		if !fileInPathSpec {
			return nil
		}

		grepResults, err := findMatchInFile(file, treeName, opts)
		if err != nil {
			return err
		}
		results = append(results, grepResults...)

		return nil
	})

	return results, err
}

// findMatchInFile takes a single File, worktree name and GrepOptions,
// and returns a slice of GrepResult containing the result of regex pattern
// matching in the given file.
func findMatchInFile(file *object.File, treeName string, opts *GrepOptions) ([]GrepResult, error) {
	var grepResults []GrepResult

	content, err := file.Contents()
	if err != nil {
		return grepResults, err
	}

	// Split the file content and parse line-by-line.
	contentByLine := strings.Split(content, "\n")
	for lineNum, cnt := range contentByLine {
		addToResult := false

		// Match the patterns and content. Break out of the loop once a
		// match is found.
		for _, pattern := range opts.Patterns {
			if pattern != nil && pattern.MatchString(cnt) {
				// Add to result only if invert match is not enabled.
				if !opts.InvertMatch {
					addToResult = true
					break
				}
			} else if opts.InvertMatch {
				// If matching fails, and invert match is enabled, add to
				// results.
				addToResult = true
				break
			}
		}

		if addToResult {
			grepResults = append(grepResults, GrepResult{
				FileName:   file.Name,
				LineNumber: lineNum + 1,
				Content:    cnt,
				TreeName:   treeName,
			})
		}
	}

	return grepResults, nil
}

// will walk up the directory tree removing all encountered empty
// directories, not just the one containing this file
func rmFileAndDirsIfEmpty(fs billy.Filesystem, name string) error {
	if err := util.RemoveAll(fs, name); err != nil {
		return err
	}

	dir := filepath.Dir(name)
	for {
		removed, err := removeDirIfEmpty(fs, dir)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if !removed {
			// directory was not empty and not removed,
			// stop checking parents
			break
		}

		// move to parent directory
		dir = filepath.Dir(dir)
	}

	return nil
}

// removeDirIfEmpty will remove the supplied directory `dir` if
// `dir` is empty
// returns true if the directory was removed
func removeDirIfEmpty(fs billy.Filesystem, dir string) (bool, error) {
	files, err := fs.ReadDir(dir)
	if err != nil {
		return false, err
	}

	if len(files) > 0 {
		return false, nil
	}

	err = fs.Remove(dir)
	if err != nil {
		return false, err
	}

	return true, nil
}

type indexBuilder struct {
	entries map[string]*index.Entry
}

func newIndexBuilder(idx *index.Index) *indexBuilder {
	entries := make(map[string]*index.Entry, len(idx.Entries))
	for _, e := range idx.Entries {
		entries[e.Name] = e
	}
	return &indexBuilder{
		entries: entries,
	}
}

func (b *indexBuilder) Write(idx *index.Index) {
	idx.Entries = idx.Entries[:0]
	for _, e := range b.entries {
		idx.Entries = append(idx.Entries, e)
	}
}

func (b *indexBuilder) Add(e *index.Entry) {
	b.entries[e.Name] = e
}

func (b *indexBuilder) Remove(name string) {
	delete(b.entries, filepath.ToSlash(name))
}
