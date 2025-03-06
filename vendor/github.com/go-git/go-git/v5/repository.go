package git

import (
	"bytes"
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/internal/path_util"
	"github.com/go-git/go-git/v5/internal/revision"
	"github.com/go-git/go-git/v5/internal/url"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	formatcfg "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/plumbing/hash"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/filesystem/dotgit"
	"github.com/go-git/go-git/v5/utils/ioutil"
)

// GitDirName this is a special folder where all the git stuff is.
const GitDirName = ".git"

var (
	// ErrBranchExists an error stating the specified branch already exists
	ErrBranchExists = errors.New("branch already exists")
	// ErrBranchNotFound an error stating the specified branch does not exist
	ErrBranchNotFound = errors.New("branch not found")
	// ErrTagExists an error stating the specified tag already exists
	ErrTagExists = errors.New("tag already exists")
	// ErrTagNotFound an error stating the specified tag does not exist
	ErrTagNotFound = errors.New("tag not found")
	// ErrFetching is returned when the packfile could not be downloaded
	ErrFetching = errors.New("unable to fetch packfile")

	ErrInvalidReference            = errors.New("invalid reference, should be a tag or a branch")
	ErrRepositoryNotExists         = errors.New("repository does not exist")
	ErrRepositoryIncomplete        = errors.New("repository's commondir path does not exist")
	ErrRepositoryAlreadyExists     = errors.New("repository already exists")
	ErrRemoteNotFound              = errors.New("remote not found")
	ErrRemoteExists                = errors.New("remote already exists")
	ErrAnonymousRemoteName         = errors.New("anonymous remote name must be 'anonymous'")
	ErrWorktreeNotProvided         = errors.New("worktree should be provided")
	ErrIsBareRepository            = errors.New("worktree not available in a bare repository")
	ErrUnableToResolveCommit       = errors.New("unable to resolve commit")
	ErrPackedObjectsNotSupported   = errors.New("packed objects not supported")
	ErrSHA256NotSupported          = errors.New("go-git was not compiled with SHA256 support")
	ErrAlternatePathNotSupported   = errors.New("alternate path must use the file scheme")
	ErrUnsupportedMergeStrategy    = errors.New("unsupported merge strategy")
	ErrFastForwardMergeNotPossible = errors.New("not possible to fast-forward merge changes")
)

// Repository represents a git repository
type Repository struct {
	Storer storage.Storer

	r  map[string]*Remote
	wt billy.Filesystem
}

type InitOptions struct {
	// The default branch (e.g. "refs/heads/master")
	DefaultBranch plumbing.ReferenceName
}

// Init creates an empty git repository, based on the given Storer and worktree.
// The worktree Filesystem is optional, if nil a bare repository is created. If
// the given storer is not empty ErrRepositoryAlreadyExists is returned
func Init(s storage.Storer, worktree billy.Filesystem) (*Repository, error) {
	options := InitOptions{
		DefaultBranch: plumbing.Master,
	}
	return InitWithOptions(s, worktree, options)
}

func InitWithOptions(s storage.Storer, worktree billy.Filesystem, options InitOptions) (*Repository, error) {
	if err := initStorer(s); err != nil {
		return nil, err
	}

	if options.DefaultBranch == "" {
		options.DefaultBranch = plumbing.Master
	}

	if err := options.DefaultBranch.Validate(); err != nil {
		return nil, err
	}

	r := newRepository(s, worktree)
	_, err := r.Reference(plumbing.HEAD, false)
	switch err {
	case plumbing.ErrReferenceNotFound:
	case nil:
		return nil, ErrRepositoryAlreadyExists
	default:
		return nil, err
	}

	h := plumbing.NewSymbolicReference(plumbing.HEAD, options.DefaultBranch)
	if err := s.SetReference(h); err != nil {
		return nil, err
	}

	if worktree == nil {
		_ = r.setIsBare(true)
		return r, nil
	}

	return r, setWorktreeAndStoragePaths(r, worktree)
}

func initStorer(s storer.Storer) error {
	i, ok := s.(storer.Initializer)
	if !ok {
		return nil
	}

	return i.Init()
}

func setWorktreeAndStoragePaths(r *Repository, worktree billy.Filesystem) error {
	type fsBased interface {
		Filesystem() billy.Filesystem
	}

	// .git file is only created if the storage is file based and the file
	// system is osfs.OS
	fs, isFSBased := r.Storer.(fsBased)
	if !isFSBased {
		return nil
	}

	if err := createDotGitFile(worktree, fs.Filesystem()); err != nil {
		return err
	}

	return setConfigWorktree(r, worktree, fs.Filesystem())
}

func createDotGitFile(worktree, storage billy.Filesystem) error {
	path, err := filepath.Rel(worktree.Root(), storage.Root())
	if err != nil {
		path = storage.Root()
	}

	if path == GitDirName {
		// not needed, since the folder is the default place
		return nil
	}

	f, err := worktree.Create(GitDirName)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = fmt.Fprintf(f, "gitdir: %s\n", path)
	return err
}

func setConfigWorktree(r *Repository, worktree, storage billy.Filesystem) error {
	path, err := filepath.Rel(storage.Root(), worktree.Root())
	if err != nil {
		path = worktree.Root()
	}

	if path == ".." {
		// not needed, since the folder is the default place
		return nil
	}

	cfg, err := r.Config()
	if err != nil {
		return err
	}

	cfg.Core.Worktree = path
	return r.Storer.SetConfig(cfg)
}

// Open opens a git repository using the given Storer and worktree filesystem,
// if the given storer is complete empty ErrRepositoryNotExists is returned.
// The worktree can be nil when the repository being opened is bare, if the
// repository is a normal one (not bare) and worktree is nil the err
// ErrWorktreeNotProvided is returned
func Open(s storage.Storer, worktree billy.Filesystem) (*Repository, error) {
	_, err := s.Reference(plumbing.HEAD)
	if err == plumbing.ErrReferenceNotFound {
		return nil, ErrRepositoryNotExists
	}

	if err != nil {
		return nil, err
	}

	return newRepository(s, worktree), nil
}

// Clone a repository into the given Storer and worktree Filesystem with the
// given options, if worktree is nil a bare repository is created. If the given
// storer is not empty ErrRepositoryAlreadyExists is returned.
func Clone(s storage.Storer, worktree billy.Filesystem, o *CloneOptions) (*Repository, error) {
	return CloneContext(context.Background(), s, worktree, o)
}

// CloneContext a repository into the given Storer and worktree Filesystem with
// the given options, if worktree is nil a bare repository is created. If the
// given storer is not empty ErrRepositoryAlreadyExists is returned.
//
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects the
// transport operations.
func CloneContext(
	ctx context.Context, s storage.Storer, worktree billy.Filesystem, o *CloneOptions,
) (*Repository, error) {
	r, err := Init(s, worktree)
	if err != nil {
		return nil, err
	}

	return r, r.clone(ctx, o)
}

// PlainInit create an empty git repository at the given path. isBare defines
// if the repository will have worktree (non-bare) or not (bare), if the path
// is not empty ErrRepositoryAlreadyExists is returned.
func PlainInit(path string, isBare bool) (*Repository, error) {
	return PlainInitWithOptions(path, &PlainInitOptions{
		Bare: isBare,
	})
}

func PlainInitWithOptions(path string, opts *PlainInitOptions) (*Repository, error) {
	if opts == nil {
		opts = &PlainInitOptions{}
	}

	var wt, dot billy.Filesystem

	if opts.Bare {
		dot = osfs.New(path)
	} else {
		wt = osfs.New(path)
		dot, _ = wt.Chroot(GitDirName)
	}

	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	r, err := InitWithOptions(s, wt, opts.InitOptions)
	if err != nil {
		return nil, err
	}

	cfg, err := r.Config()
	if err != nil {
		return nil, err
	}

	if opts.ObjectFormat != "" {
		if opts.ObjectFormat == formatcfg.SHA256 && hash.CryptoType != crypto.SHA256 {
			return nil, ErrSHA256NotSupported
		}

		cfg.Core.RepositoryFormatVersion = formatcfg.Version_1
		cfg.Extensions.ObjectFormat = opts.ObjectFormat
	}

	err = r.Storer.SetConfig(cfg)
	if err != nil {
		return nil, err
	}

	return r, err
}

// PlainOpen opens a git repository from the given path. It detects if the
// repository is bare or a normal one. If the path doesn't contain a valid
// repository ErrRepositoryNotExists is returned
func PlainOpen(path string) (*Repository, error) {
	return PlainOpenWithOptions(path, &PlainOpenOptions{})
}

// PlainOpenWithOptions opens a git repository from the given path with specific
// options. See PlainOpen for more info.
func PlainOpenWithOptions(path string, o *PlainOpenOptions) (*Repository, error) {
	dot, wt, err := dotGitToOSFilesystems(path, o.DetectDotGit)
	if err != nil {
		return nil, err
	}

	if _, err := dot.Stat(""); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrRepositoryNotExists
		}

		return nil, err
	}

	var repositoryFs billy.Filesystem

	if o.EnableDotGitCommonDir {
		dotGitCommon, err := dotGitCommonDirectory(dot)
		if err != nil {
			return nil, err
		}
		repositoryFs = dotgit.NewRepositoryFilesystem(dot, dotGitCommon)
	} else {
		repositoryFs = dot
	}

	s := filesystem.NewStorage(repositoryFs, cache.NewObjectLRUDefault())

	return Open(s, wt)
}

func dotGitToOSFilesystems(path string, detect bool) (dot, wt billy.Filesystem, err error) {
	path, err = path_util.ReplaceTildeWithHome(path)
	if err != nil {
		return nil, nil, err
	}

	if path, err = filepath.Abs(path); err != nil {
		return nil, nil, err
	}

	var fs billy.Filesystem
	var fi os.FileInfo
	for {
		fs = osfs.New(path)

		pathinfo, err := fs.Stat("/")
		if !os.IsNotExist(err) {
			if pathinfo == nil {
				return nil, nil, err
			}
			if !pathinfo.IsDir() && detect {
				fs = osfs.New(filepath.Dir(path))
			}
		}

		fi, err = fs.Stat(GitDirName)
		if err == nil {
			// no error; stop
			break
		}
		if !os.IsNotExist(err) {
			// unknown error; stop
			return nil, nil, err
		}
		if detect {
			// try its parent as long as we haven't reached
			// the root dir
			if dir := filepath.Dir(path); dir != path {
				path = dir
				continue
			}
		}
		// not detecting via parent dirs and the dir does not exist;
		// stop
		return fs, nil, nil
	}

	if fi.IsDir() {
		dot, err = fs.Chroot(GitDirName)
		return dot, fs, err
	}

	dot, err = dotGitFileToOSFilesystem(path, fs)
	if err != nil {
		return nil, nil, err
	}

	return dot, fs, nil
}

func dotGitFileToOSFilesystem(path string, fs billy.Filesystem) (bfs billy.Filesystem, err error) {
	f, err := fs.Open(GitDirName)
	if err != nil {
		return nil, err
	}
	defer ioutil.CheckClose(f, &err)

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	line := string(b)
	const prefix = "gitdir: "
	if !strings.HasPrefix(line, prefix) {
		return nil, fmt.Errorf(".git file has no %s prefix", prefix)
	}

	gitdir := strings.Split(line[len(prefix):], "\n")[0]
	gitdir = strings.TrimSpace(gitdir)
	if filepath.IsAbs(gitdir) {
		return osfs.New(gitdir), nil
	}

	return osfs.New(fs.Join(path, gitdir)), nil
}

func dotGitCommonDirectory(fs billy.Filesystem) (commonDir billy.Filesystem, err error) {
	f, err := fs.Open("commondir")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(b) > 0 {
		path := strings.TrimSpace(string(b))
		if filepath.IsAbs(path) {
			commonDir = osfs.New(path)
		} else {
			commonDir = osfs.New(filepath.Join(fs.Root(), path))
		}
		if _, err := commonDir.Stat(""); err != nil {
			if os.IsNotExist(err) {
				return nil, ErrRepositoryIncomplete
			}

			return nil, err
		}
	}

	return commonDir, nil
}

// PlainClone a repository into the path with the given options, isBare defines
// if the new repository will be bare or normal. If the path is not empty
// ErrRepositoryAlreadyExists is returned.
//
// TODO(mcuadros): move isBare to CloneOptions in v5
func PlainClone(path string, isBare bool, o *CloneOptions) (*Repository, error) {
	return PlainCloneContext(context.Background(), path, isBare, o)
}

// PlainCloneContext a repository into the path with the given options, isBare
// defines if the new repository will be bare or normal. If the path is not empty
// ErrRepositoryAlreadyExists is returned.
//
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects the
// transport operations.
//
// TODO(mcuadros): move isBare to CloneOptions in v5
// TODO(smola): refuse upfront to clone on a non-empty directory in v5, see #1027
func PlainCloneContext(ctx context.Context, path string, isBare bool, o *CloneOptions) (*Repository, error) {
	cleanup, cleanupParent, err := checkIfCleanupIsNeeded(path)
	if err != nil {
		return nil, err
	}

	if o.Mirror {
		isBare = true
	}
	r, err := PlainInit(path, isBare)
	if err != nil {
		return nil, err
	}

	err = r.clone(ctx, o)
	if err != nil && err != ErrRepositoryAlreadyExists {
		if cleanup {
			_ = cleanUpDir(path, cleanupParent)
		}
	}

	return r, err
}

func newRepository(s storage.Storer, worktree billy.Filesystem) *Repository {
	return &Repository{
		Storer: s,
		wt:     worktree,
		r:      make(map[string]*Remote),
	}
}

func checkIfCleanupIsNeeded(path string) (cleanup bool, cleanParent bool, err error) {
	fi, err := osfs.Default.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, true, nil
		}

		return false, false, err
	}

	if !fi.IsDir() {
		return false, false, fmt.Errorf("path is not a directory: %s", path)
	}

	files, err := osfs.Default.ReadDir(path)
	if err != nil {
		return false, false, err
	}

	if len(files) == 0 {
		return true, false, nil
	}

	return false, false, nil
}

func cleanUpDir(path string, all bool) error {
	if all {
		return util.RemoveAll(osfs.Default, path)
	}

	files, err := osfs.Default.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range files {
		if err := util.RemoveAll(osfs.Default, osfs.Default.Join(path, fi.Name())); err != nil {
			return err
		}
	}

	return err
}

// Config return the repository config. In a filesystem backed repository this
// means read the `.git/config`.
func (r *Repository) Config() (*config.Config, error) {
	return r.Storer.Config()
}

// SetConfig marshall and writes the repository config. In a filesystem backed
// repository this means write the `.git/config`. This function should be called
// with the result of `Repository.Config` and never with the output of
// `Repository.ConfigScoped`.
func (r *Repository) SetConfig(cfg *config.Config) error {
	return r.Storer.SetConfig(cfg)
}

// ConfigScoped returns the repository config, merged with requested scope and
// lower. For example if, config.GlobalScope is given the local and global config
// are returned merged in one config value.
func (r *Repository) ConfigScoped(scope config.Scope) (*config.Config, error) {
	// TODO(mcuadros): v6, add this as ConfigOptions.Scoped

	var err error
	system := config.NewConfig()
	if scope >= config.SystemScope {
		system, err = config.LoadConfig(config.SystemScope)
		if err != nil {
			return nil, err
		}
	}

	global := config.NewConfig()
	if scope >= config.GlobalScope {
		global, err = config.LoadConfig(config.GlobalScope)
		if err != nil {
			return nil, err
		}
	}

	local, err := r.Storer.Config()
	if err != nil {
		return nil, err
	}

	_ = mergo.Merge(global, system)
	_ = mergo.Merge(local, global)
	return local, nil
}

// Remote return a remote if exists
func (r *Repository) Remote(name string) (*Remote, error) {
	cfg, err := r.Config()
	if err != nil {
		return nil, err
	}

	c, ok := cfg.Remotes[name]
	if !ok {
		return nil, ErrRemoteNotFound
	}

	return NewRemote(r.Storer, c), nil
}

// Remotes returns a list with all the remotes
func (r *Repository) Remotes() ([]*Remote, error) {
	cfg, err := r.Config()
	if err != nil {
		return nil, err
	}

	remotes := make([]*Remote, len(cfg.Remotes))

	var i int
	for _, c := range cfg.Remotes {
		remotes[i] = NewRemote(r.Storer, c)
		i++
	}

	return remotes, nil
}

// CreateRemote creates a new remote
func (r *Repository) CreateRemote(c *config.RemoteConfig) (*Remote, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	remote := NewRemote(r.Storer, c)

	cfg, err := r.Config()
	if err != nil {
		return nil, err
	}

	if _, ok := cfg.Remotes[c.Name]; ok {
		return nil, ErrRemoteExists
	}

	cfg.Remotes[c.Name] = c
	return remote, r.Storer.SetConfig(cfg)
}

// CreateRemoteAnonymous creates a new anonymous remote. c.Name must be "anonymous".
// It's used like 'git fetch git@github.com:src-d/go-git.git master:master'.
func (r *Repository) CreateRemoteAnonymous(c *config.RemoteConfig) (*Remote, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	if c.Name != "anonymous" {
		return nil, ErrAnonymousRemoteName
	}

	remote := NewRemote(r.Storer, c)

	return remote, nil
}

// DeleteRemote delete a remote from the repository and delete the config
func (r *Repository) DeleteRemote(name string) error {
	cfg, err := r.Config()
	if err != nil {
		return err
	}

	if _, ok := cfg.Remotes[name]; !ok {
		return ErrRemoteNotFound
	}

	delete(cfg.Remotes, name)
	return r.Storer.SetConfig(cfg)
}

// Branch return a Branch if exists
func (r *Repository) Branch(name string) (*config.Branch, error) {
	cfg, err := r.Config()
	if err != nil {
		return nil, err
	}

	b, ok := cfg.Branches[name]
	if !ok {
		return nil, ErrBranchNotFound
	}

	return b, nil
}

// CreateBranch creates a new Branch
func (r *Repository) CreateBranch(c *config.Branch) error {
	if err := c.Validate(); err != nil {
		return err
	}

	cfg, err := r.Config()
	if err != nil {
		return err
	}

	if _, ok := cfg.Branches[c.Name]; ok {
		return ErrBranchExists
	}

	cfg.Branches[c.Name] = c
	return r.Storer.SetConfig(cfg)
}

// DeleteBranch delete a Branch from the repository and delete the config
func (r *Repository) DeleteBranch(name string) error {
	cfg, err := r.Config()
	if err != nil {
		return err
	}

	if _, ok := cfg.Branches[name]; !ok {
		return ErrBranchNotFound
	}

	delete(cfg.Branches, name)
	return r.Storer.SetConfig(cfg)
}

// CreateTag creates a tag. If opts is included, the tag is an annotated tag,
// otherwise a lightweight tag is created.
func (r *Repository) CreateTag(name string, hash plumbing.Hash, opts *CreateTagOptions) (*plumbing.Reference, error) {
	rname := plumbing.NewTagReferenceName(name)
	if err := rname.Validate(); err != nil {
		return nil, err
	}

	_, err := r.Storer.Reference(rname)
	switch err {
	case nil:
		// Tag exists, this is an error
		return nil, ErrTagExists
	case plumbing.ErrReferenceNotFound:
		// Tag missing, available for creation, pass this
	default:
		// Some other error
		return nil, err
	}

	var target plumbing.Hash
	if opts != nil {
		target, err = r.createTagObject(name, hash, opts)
		if err != nil {
			return nil, err
		}
	} else {
		target = hash
	}

	ref := plumbing.NewHashReference(rname, target)
	if err = r.Storer.SetReference(ref); err != nil {
		return nil, err
	}

	return ref, nil
}

func (r *Repository) createTagObject(name string, hash plumbing.Hash, opts *CreateTagOptions) (plumbing.Hash, error) {
	if err := opts.Validate(r, hash); err != nil {
		return plumbing.ZeroHash, err
	}

	rawobj, err := object.GetObject(r.Storer, hash)
	if err != nil {
		return plumbing.ZeroHash, err
	}

	tag := &object.Tag{
		Name:       name,
		Tagger:     *opts.Tagger,
		Message:    opts.Message,
		TargetType: rawobj.Type(),
		Target:     hash,
	}

	if opts.SignKey != nil {
		sig, err := r.buildTagSignature(tag, opts.SignKey)
		if err != nil {
			return plumbing.ZeroHash, err
		}

		tag.PGPSignature = sig
	}

	obj := r.Storer.NewEncodedObject()
	if err := tag.Encode(obj); err != nil {
		return plumbing.ZeroHash, err
	}

	return r.Storer.SetEncodedObject(obj)
}

func (r *Repository) buildTagSignature(tag *object.Tag, signKey *openpgp.Entity) (string, error) {
	encoded := &plumbing.MemoryObject{}
	if err := tag.Encode(encoded); err != nil {
		return "", err
	}

	rdr, err := encoded.Reader()
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := openpgp.ArmoredDetachSign(&b, signKey, rdr, nil); err != nil {
		return "", err
	}

	return b.String(), nil
}

// Tag returns a tag from the repository.
//
// If you want to check to see if the tag is an annotated tag, you can call
// TagObject on the hash of the reference in ForEach:
//
//	ref, err := r.Tag("v0.1.0")
//	if err != nil {
//	  // Handle error
//	}
//
//	obj, err := r.TagObject(ref.Hash())
//	switch err {
//	case nil:
//	  // Tag object present
//	case plumbing.ErrObjectNotFound:
//	  // Not a tag object
//	default:
//	  // Some other error
//	}
func (r *Repository) Tag(name string) (*plumbing.Reference, error) {
	ref, err := r.Reference(plumbing.ReferenceName(path.Join("refs", "tags", name)), false)
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			// Return a friendly error for this one, versus just ReferenceNotFound.
			return nil, ErrTagNotFound
		}

		return nil, err
	}

	return ref, nil
}

// DeleteTag deletes a tag from the repository.
func (r *Repository) DeleteTag(name string) error {
	_, err := r.Tag(name)
	if err != nil {
		return err
	}

	return r.Storer.RemoveReference(plumbing.ReferenceName(path.Join("refs", "tags", name)))
}

func (r *Repository) resolveToCommitHash(h plumbing.Hash) (plumbing.Hash, error) {
	obj, err := r.Storer.EncodedObject(plumbing.AnyObject, h)
	if err != nil {
		return plumbing.ZeroHash, err
	}
	switch obj.Type() {
	case plumbing.TagObject:
		t, err := object.DecodeTag(r.Storer, obj)
		if err != nil {
			return plumbing.ZeroHash, err
		}
		return r.resolveToCommitHash(t.Target)
	case plumbing.CommitObject:
		return h, nil
	default:
		return plumbing.ZeroHash, ErrUnableToResolveCommit
	}
}

// Clone clones a remote repository
func (r *Repository) clone(ctx context.Context, o *CloneOptions) error {
	if err := o.Validate(); err != nil {
		return err
	}

	c := &config.RemoteConfig{
		Name:   o.RemoteName,
		URLs:   []string{o.URL},
		Fetch:  r.cloneRefSpec(o),
		Mirror: o.Mirror,
	}

	if _, err := r.CreateRemote(c); err != nil {
		return err
	}

	// When the repository to clone is on the local machine,
	// instead of using hard links, automatically setup .git/objects/info/alternates
	// to share the objects with the source repository
	if o.Shared {
		if !url.IsLocalEndpoint(o.URL) {
			return ErrAlternatePathNotSupported
		}
		altpath := o.URL
		remoteRepo, err := PlainOpen(o.URL)
		if err != nil {
			return fmt.Errorf("failed to open remote repository: %w", err)
		}
		conf, err := remoteRepo.Config()
		if err != nil {
			return fmt.Errorf("failed to read remote repository configuration: %w", err)
		}
		if !conf.Core.IsBare {
			altpath = path.Join(altpath, GitDirName)
		}
		if err := r.Storer.AddAlternate(altpath); err != nil {
			return fmt.Errorf("failed to add alternate file to git objects dir: %w", err)
		}
	}

	ref, err := r.fetchAndUpdateReferences(ctx, &FetchOptions{
		RefSpecs:        c.Fetch,
		Depth:           o.Depth,
		Auth:            o.Auth,
		Progress:        o.Progress,
		Tags:            o.Tags,
		RemoteName:      o.RemoteName,
		InsecureSkipTLS: o.InsecureSkipTLS,
		CABundle:        o.CABundle,
		ProxyOptions:    o.ProxyOptions,
	}, o.ReferenceName)
	if err != nil {
		return err
	}

	if r.wt != nil && !o.NoCheckout {
		w, err := r.Worktree()
		if err != nil {
			return err
		}

		head, err := r.Head()
		if err != nil {
			return err
		}

		if err := w.Reset(&ResetOptions{
			Mode:   MergeReset,
			Commit: head.Hash(),
		}); err != nil {
			return err
		}

		if o.RecurseSubmodules != NoRecurseSubmodules {
			if err := w.updateSubmodules(ctx, &SubmoduleUpdateOptions{
				RecurseSubmodules: o.RecurseSubmodules,
				Depth: func() int {
					if o.ShallowSubmodules {
						return 1
					}
					return 0
				}(),
				Auth: o.Auth,
			}); err != nil {
				return err
			}
		}
	}

	if err := r.updateRemoteConfigIfNeeded(o, c, ref); err != nil {
		return err
	}

	if !o.Mirror && ref.Name().IsBranch() {
		branchRef := ref.Name()
		branchName := strings.Split(string(branchRef), "refs/heads/")[1]

		b := &config.Branch{
			Name:  branchName,
			Merge: branchRef,
		}

		if o.RemoteName == "" {
			b.Remote = "origin"
		} else {
			b.Remote = o.RemoteName
		}

		if err := r.CreateBranch(b); err != nil {
			return err
		}
	}

	return nil
}

const (
	refspecTag              = "+refs/tags/%s:refs/tags/%[1]s"
	refspecSingleBranch     = "+refs/heads/%s:refs/remotes/%s/%[1]s"
	refspecSingleBranchHEAD = "+HEAD:refs/remotes/%s/HEAD"
)

func (r *Repository) cloneRefSpec(o *CloneOptions) []config.RefSpec {
	switch {
	case o.Mirror:
		return []config.RefSpec{"+refs/*:refs/*"}
	case o.ReferenceName.IsTag():
		return []config.RefSpec{
			config.RefSpec(fmt.Sprintf(refspecTag, o.ReferenceName.Short())),
		}
	case o.SingleBranch && o.ReferenceName == plumbing.HEAD:
		return []config.RefSpec{
			config.RefSpec(fmt.Sprintf(refspecSingleBranchHEAD, o.RemoteName)),
		}
	case o.SingleBranch:
		return []config.RefSpec{
			config.RefSpec(fmt.Sprintf(refspecSingleBranch, o.ReferenceName.Short(), o.RemoteName)),
		}
	default:
		return []config.RefSpec{
			config.RefSpec(fmt.Sprintf(config.DefaultFetchRefSpec, o.RemoteName)),
		}
	}
}

func (r *Repository) setIsBare(isBare bool) error {
	cfg, err := r.Config()
	if err != nil {
		return err
	}

	cfg.Core.IsBare = isBare
	return r.Storer.SetConfig(cfg)
}

func (r *Repository) updateRemoteConfigIfNeeded(o *CloneOptions, c *config.RemoteConfig, _ *plumbing.Reference) error {
	if !o.SingleBranch {
		return nil
	}

	c.Fetch = r.cloneRefSpec(o)

	cfg, err := r.Config()
	if err != nil {
		return err
	}

	cfg.Remotes[c.Name] = c
	return r.Storer.SetConfig(cfg)
}

func (r *Repository) fetchAndUpdateReferences(
	ctx context.Context, o *FetchOptions, ref plumbing.ReferenceName,
) (*plumbing.Reference, error) {

	if err := o.Validate(); err != nil {
		return nil, err
	}

	remote, err := r.Remote(o.RemoteName)
	if err != nil {
		return nil, err
	}

	objsUpdated := true
	remoteRefs, err := remote.fetch(ctx, o)
	if err == NoErrAlreadyUpToDate {
		objsUpdated = false
	} else if err == packfile.ErrEmptyPackfile {
		return nil, ErrFetching
	} else if err != nil {
		return nil, err
	}

	resolvedRef, err := expand_ref(remoteRefs, ref)
	if err != nil {
		return nil, err
	}

	refsUpdated, err := r.updateReferences(remote.c.Fetch, resolvedRef)
	if err != nil {
		return nil, err
	}

	if !objsUpdated && !refsUpdated {
		return nil, NoErrAlreadyUpToDate
	}

	return resolvedRef, nil
}

func (r *Repository) updateReferences(spec []config.RefSpec,
	resolvedRef *plumbing.Reference) (updated bool, err error) {

	if !resolvedRef.Name().IsBranch() {
		// Detached HEAD mode
		h, err := r.resolveToCommitHash(resolvedRef.Hash())
		if err != nil {
			return false, err
		}
		head := plumbing.NewHashReference(plumbing.HEAD, h)
		return updateReferenceStorerIfNeeded(r.Storer, head)
	}

	refs := []*plumbing.Reference{
		// Create local reference for the resolved ref
		resolvedRef,
		// Create local symbolic HEAD
		plumbing.NewSymbolicReference(plumbing.HEAD, resolvedRef.Name()),
	}

	refs = append(refs, r.calculateRemoteHeadReference(spec, resolvedRef)...)

	for _, ref := range refs {
		u, err := updateReferenceStorerIfNeeded(r.Storer, ref)
		if err != nil {
			return updated, err
		}

		if u {
			updated = true
		}
	}

	return
}

func (r *Repository) calculateRemoteHeadReference(spec []config.RefSpec,
	resolvedHead *plumbing.Reference) []*plumbing.Reference {

	var refs []*plumbing.Reference

	// Create resolved HEAD reference with remote prefix if it does not
	// exist. This is needed when using single branch and HEAD.
	for _, rs := range spec {
		name := resolvedHead.Name()
		if !rs.Match(name) {
			continue
		}

		name = rs.Dst(name)
		_, err := r.Storer.Reference(name)
		if err == plumbing.ErrReferenceNotFound {
			refs = append(refs, plumbing.NewHashReference(name, resolvedHead.Hash()))
		}
	}

	return refs
}

func checkAndUpdateReferenceStorerIfNeeded(
	s storer.ReferenceStorer, r, old *plumbing.Reference) (
	updated bool, err error) {
	p, err := s.Reference(r.Name())
	if err != nil && err != plumbing.ErrReferenceNotFound {
		return false, err
	}

	// we use the string method to compare references, is the easiest way
	if err == plumbing.ErrReferenceNotFound || r.String() != p.String() {
		if err := s.CheckAndSetReference(r, old); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func updateReferenceStorerIfNeeded(
	s storer.ReferenceStorer, r *plumbing.Reference) (updated bool, err error) {
	return checkAndUpdateReferenceStorerIfNeeded(s, r, nil)
}

// Fetch fetches references along with the objects necessary to complete
// their histories, from the remote named as FetchOptions.RemoteName.
//
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
func (r *Repository) Fetch(o *FetchOptions) error {
	return r.FetchContext(context.Background(), o)
}

// FetchContext fetches references along with the objects necessary to complete
// their histories, from the remote named as FetchOptions.RemoteName.
//
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
//
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects the
// transport operations.
func (r *Repository) FetchContext(ctx context.Context, o *FetchOptions) error {
	if err := o.Validate(); err != nil {
		return err
	}

	remote, err := r.Remote(o.RemoteName)
	if err != nil {
		return err
	}

	return remote.FetchContext(ctx, o)
}

// Push performs a push to the remote. Returns NoErrAlreadyUpToDate if
// the remote was already up-to-date, from the remote named as
// FetchOptions.RemoteName.
func (r *Repository) Push(o *PushOptions) error {
	return r.PushContext(context.Background(), o)
}

// PushContext performs a push to the remote. Returns NoErrAlreadyUpToDate if
// the remote was already up-to-date, from the remote named as
// FetchOptions.RemoteName.
//
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects the
// transport operations.
func (r *Repository) PushContext(ctx context.Context, o *PushOptions) error {
	if err := o.Validate(); err != nil {
		return err
	}

	remote, err := r.Remote(o.RemoteName)
	if err != nil {
		return err
	}

	return remote.PushContext(ctx, o)
}

// Log returns the commit history from the given LogOptions.
func (r *Repository) Log(o *LogOptions) (object.CommitIter, error) {
	fn := commitIterFunc(o.Order)
	if fn == nil {
		return nil, fmt.Errorf("invalid Order=%v", o.Order)
	}

	var (
		it  object.CommitIter
		err error
	)
	if o.All {
		it, err = r.logAll(fn)
	} else {
		it, err = r.log(o.From, fn)
	}

	if err != nil {
		return nil, err
	}

	if o.FileName != nil {
		// for `git log --all` also check parent (if the next commit comes from the real parent)
		it = r.logWithFile(*o.FileName, it, o.All)
	}
	if o.PathFilter != nil {
		it = r.logWithPathFilter(o.PathFilter, it, o.All)
	}

	if o.Since != nil || o.Until != nil {
		limitOptions := object.LogLimitOptions{Since: o.Since, Until: o.Until}
		it = r.logWithLimit(it, limitOptions)
	}

	return it, nil
}

func (r *Repository) log(from plumbing.Hash, commitIterFunc func(*object.Commit) object.CommitIter) (object.CommitIter, error) {
	h := from
	if from == plumbing.ZeroHash {
		head, err := r.Head()
		if err != nil {
			return nil, err
		}

		h = head.Hash()
	}

	commit, err := r.CommitObject(h)
	if err != nil {
		return nil, err
	}
	return commitIterFunc(commit), nil
}

func (r *Repository) logAll(commitIterFunc func(*object.Commit) object.CommitIter) (object.CommitIter, error) {
	return object.NewCommitAllIter(r.Storer, commitIterFunc)
}

func (*Repository) logWithFile(fileName string, commitIter object.CommitIter, checkParent bool) object.CommitIter {
	return object.NewCommitPathIterFromIter(
		func(path string) bool {
			return path == fileName
		},
		commitIter,
		checkParent,
	)
}

func (*Repository) logWithPathFilter(pathFilter func(string) bool, commitIter object.CommitIter, checkParent bool) object.CommitIter {
	return object.NewCommitPathIterFromIter(
		pathFilter,
		commitIter,
		checkParent,
	)
}

func (*Repository) logWithLimit(commitIter object.CommitIter, limitOptions object.LogLimitOptions) object.CommitIter {
	return object.NewCommitLimitIterFromIter(commitIter, limitOptions)
}

func commitIterFunc(order LogOrder) func(c *object.Commit) object.CommitIter {
	switch order {
	case LogOrderDefault:
		return func(c *object.Commit) object.CommitIter {
			return object.NewCommitPreorderIter(c, nil, nil)
		}
	case LogOrderDFS:
		return func(c *object.Commit) object.CommitIter {
			return object.NewCommitPreorderIter(c, nil, nil)
		}
	case LogOrderDFSPost:
		return func(c *object.Commit) object.CommitIter {
			return object.NewCommitPostorderIter(c, nil)
		}
	case LogOrderBSF:
		return func(c *object.Commit) object.CommitIter {
			return object.NewCommitIterBSF(c, nil, nil)
		}
	case LogOrderCommitterTime:
		return func(c *object.Commit) object.CommitIter {
			return object.NewCommitIterCTime(c, nil, nil)
		}
	}
	return nil
}

// Tags returns all the tag References in a repository.
//
// If you want to check to see if the tag is an annotated tag, you can call
// TagObject on the hash Reference passed in through ForEach:
//
//	iter, err := r.Tags()
//	if err != nil {
//	  // Handle error
//	}
//
//	if err := iter.ForEach(func (ref *plumbing.Reference) error {
//	  obj, err := r.TagObject(ref.Hash())
//	  switch err {
//	  case nil:
//	    // Tag object present
//	  case plumbing.ErrObjectNotFound:
//	    // Not a tag object
//	  default:
//	    // Some other error
//	    return err
//	  }
//	}); err != nil {
//	  // Handle outer iterator error
//	}
func (r *Repository) Tags() (storer.ReferenceIter, error) {
	refIter, err := r.Storer.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsTag()
		}, refIter), nil
}

// Branches returns all the References that are Branches.
func (r *Repository) Branches() (storer.ReferenceIter, error) {
	refIter, err := r.Storer.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsBranch()
		}, refIter), nil
}

// Notes returns all the References that are notes. For more information:
// https://git-scm.com/docs/git-notes
func (r *Repository) Notes() (storer.ReferenceIter, error) {
	refIter, err := r.Storer.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(
		func(r *plumbing.Reference) bool {
			return r.Name().IsNote()
		}, refIter), nil
}

// TreeObject return a Tree with the given hash. If not found
// plumbing.ErrObjectNotFound is returned
func (r *Repository) TreeObject(h plumbing.Hash) (*object.Tree, error) {
	return object.GetTree(r.Storer, h)
}

// TreeObjects returns an unsorted TreeIter with all the trees in the repository
func (r *Repository) TreeObjects() (*object.TreeIter, error) {
	iter, err := r.Storer.IterEncodedObjects(plumbing.TreeObject)
	if err != nil {
		return nil, err
	}

	return object.NewTreeIter(r.Storer, iter), nil
}

// CommitObject return a Commit with the given hash. If not found
// plumbing.ErrObjectNotFound is returned.
func (r *Repository) CommitObject(h plumbing.Hash) (*object.Commit, error) {
	return object.GetCommit(r.Storer, h)
}

// CommitObjects returns an unsorted CommitIter with all the commits in the repository.
func (r *Repository) CommitObjects() (object.CommitIter, error) {
	iter, err := r.Storer.IterEncodedObjects(plumbing.CommitObject)
	if err != nil {
		return nil, err
	}

	return object.NewCommitIter(r.Storer, iter), nil
}

// BlobObject returns a Blob with the given hash. If not found
// plumbing.ErrObjectNotFound is returned.
func (r *Repository) BlobObject(h plumbing.Hash) (*object.Blob, error) {
	return object.GetBlob(r.Storer, h)
}

// BlobObjects returns an unsorted BlobIter with all the blobs in the repository.
func (r *Repository) BlobObjects() (*object.BlobIter, error) {
	iter, err := r.Storer.IterEncodedObjects(plumbing.BlobObject)
	if err != nil {
		return nil, err
	}

	return object.NewBlobIter(r.Storer, iter), nil
}

// TagObject returns a Tag with the given hash. If not found
// plumbing.ErrObjectNotFound is returned. This method only returns
// annotated Tags, no lightweight Tags.
func (r *Repository) TagObject(h plumbing.Hash) (*object.Tag, error) {
	return object.GetTag(r.Storer, h)
}

// TagObjects returns a unsorted TagIter that can step through all of the annotated
// tags in the repository.
func (r *Repository) TagObjects() (*object.TagIter, error) {
	iter, err := r.Storer.IterEncodedObjects(plumbing.TagObject)
	if err != nil {
		return nil, err
	}

	return object.NewTagIter(r.Storer, iter), nil
}

// Object returns an Object with the given hash. If not found
// plumbing.ErrObjectNotFound is returned.
func (r *Repository) Object(t plumbing.ObjectType, h plumbing.Hash) (object.Object, error) {
	obj, err := r.Storer.EncodedObject(t, h)
	if err != nil {
		return nil, err
	}

	return object.DecodeObject(r.Storer, obj)
}

// Objects returns an unsorted ObjectIter with all the objects in the repository.
func (r *Repository) Objects() (*object.ObjectIter, error) {
	iter, err := r.Storer.IterEncodedObjects(plumbing.AnyObject)
	if err != nil {
		return nil, err
	}

	return object.NewObjectIter(r.Storer, iter), nil
}

// Head returns the reference where HEAD is pointing to.
func (r *Repository) Head() (*plumbing.Reference, error) {
	return storer.ResolveReference(r.Storer, plumbing.HEAD)
}

// Reference returns the reference for a given reference name. If resolved is
// true, any symbolic reference will be resolved.
func (r *Repository) Reference(name plumbing.ReferenceName, resolved bool) (
	*plumbing.Reference, error) {

	if resolved {
		return storer.ResolveReference(r.Storer, name)
	}

	return r.Storer.Reference(name)
}

// References returns an unsorted ReferenceIter for all references.
func (r *Repository) References() (storer.ReferenceIter, error) {
	return r.Storer.IterReferences()
}

// Worktree returns a worktree based on the given fs, if nil the default
// worktree will be used.
func (r *Repository) Worktree() (*Worktree, error) {
	if r.wt == nil {
		return nil, ErrIsBareRepository
	}

	return &Worktree{r: r, Filesystem: r.wt}, nil
}

func expand_ref(s storer.ReferenceStorer, ref plumbing.ReferenceName) (*plumbing.Reference, error) {
	// For improving troubleshooting, this preserves the error for the provided `ref`,
	// and returns the error for that specific ref in case all parse rules fails.
	var ret error
	for _, rule := range plumbing.RefRevParseRules {
		resolvedRef, err := storer.ResolveReference(s, plumbing.ReferenceName(fmt.Sprintf(rule, ref)))

		if err == nil {
			return resolvedRef, nil
		} else if ret == nil {
			ret = err
		}
	}

	return nil, ret
}

// ResolveRevision resolves revision to corresponding hash. It will always
// resolve to a commit hash, not a tree or annotated tag.
//
// Implemented resolvers : HEAD, branch, tag, heads/branch, refs/heads/branch,
// refs/tags/tag, refs/remotes/origin/branch, refs/remotes/origin/HEAD, tilde and caret (HEAD~1, master~^, tag~2, ref/heads/master~1, ...), selection by text (HEAD^{/fix nasty bug}), hash (prefix and full)
func (r *Repository) ResolveRevision(in plumbing.Revision) (*plumbing.Hash, error) {
	rev := in.String()
	if rev == "" {
		return &plumbing.ZeroHash, plumbing.ErrReferenceNotFound
	}

	p := revision.NewParserFromString(rev)
	items, err := p.Parse()

	if err != nil {
		return nil, err
	}

	var commit *object.Commit

	for _, item := range items {
		switch item := item.(type) {
		case revision.Ref:
			revisionRef := item

			var tryHashes []plumbing.Hash

			tryHashes = append(tryHashes, r.resolveHashPrefix(string(revisionRef))...)

			ref, err := expand_ref(r.Storer, plumbing.ReferenceName(revisionRef))
			if err == nil {
				tryHashes = append(tryHashes, ref.Hash())
			}

			// in ambiguous cases, `git rev-parse` will emit a warning, but
			// will always return the oid in preference to a ref; we don't have
			// the ability to emit a warning here, so (for speed purposes)
			// don't bother to detect the ambiguity either, just return in the
			// priority that git would.
			gotOne := false
			for _, hash := range tryHashes {
				commitObj, err := r.CommitObject(hash)
				if err == nil {
					commit = commitObj
					gotOne = true
					break
				}

				tagObj, err := r.TagObject(hash)
				if err == nil {
					// If the tag target lookup fails here, this most likely
					// represents some sort of repo corruption, so let the
					// error bubble up.
					tagCommit, err := tagObj.Commit()
					if err != nil {
						return &plumbing.ZeroHash, err
					}
					commit = tagCommit
					gotOne = true
					break
				}
			}

			if !gotOne {
				return &plumbing.ZeroHash, plumbing.ErrReferenceNotFound
			}

		case revision.CaretPath:
			depth := item.Depth

			if depth == 0 {
				break
			}

			iter := commit.Parents()

			c, err := iter.Next()

			if err != nil {
				return &plumbing.ZeroHash, err
			}

			if depth == 1 {
				commit = c

				break
			}

			c, err = iter.Next()

			if err != nil {
				return &plumbing.ZeroHash, err
			}

			commit = c
		case revision.TildePath:
			for i := 0; i < item.Depth; i++ {
				c, err := commit.Parents().Next()

				if err != nil {
					return &plumbing.ZeroHash, err
				}

				commit = c
			}
		case revision.CaretReg:
			history := object.NewCommitPreorderIter(commit, nil, nil)

			re := item.Regexp
			negate := item.Negate

			var c *object.Commit

			err := history.ForEach(func(hc *object.Commit) error {
				if !negate && re.MatchString(hc.Message) {
					c = hc
					return storer.ErrStop
				}

				if negate && !re.MatchString(hc.Message) {
					c = hc
					return storer.ErrStop
				}

				return nil
			})
			if err != nil {
				return &plumbing.ZeroHash, err
			}

			if c == nil {
				return &plumbing.ZeroHash, fmt.Errorf("no commit message match regexp: %q", re.String())
			}

			commit = c
		}
	}

	if commit == nil {
		return &plumbing.ZeroHash, plumbing.ErrReferenceNotFound
	}

	return &commit.Hash, nil
}

// resolveHashPrefix returns a list of potential hashes that the given string
// is a prefix of. It quietly swallows errors, returning nil.
func (r *Repository) resolveHashPrefix(hashStr string) []plumbing.Hash {
	// Handle complete and partial hashes.
	// plumbing.NewHash forces args into a full 20 byte hash, which isn't suitable
	// for partial hashes since they will become zero-filled.

	if hashStr == "" {
		return nil
	}
	if len(hashStr) == len(plumbing.ZeroHash)*2 {
		// Only a full hash is possible.
		hexb, err := hex.DecodeString(hashStr)
		if err != nil {
			return nil
		}
		var h plumbing.Hash
		copy(h[:], hexb)
		return []plumbing.Hash{h}
	}

	// Partial hash.
	// hex.DecodeString only decodes to complete bytes, so only works with pairs of hex digits.
	evenHex := hashStr[:len(hashStr)&^1]
	hexb, err := hex.DecodeString(evenHex)
	if err != nil {
		return nil
	}
	candidates := expandPartialHash(r.Storer, hexb)
	if len(evenHex) == len(hashStr) {
		// The prefix was an exact number of bytes.
		return candidates
	}
	// Do another prefix check to ensure the dangling nybble is correct.
	var hashes []plumbing.Hash
	for _, h := range candidates {
		if strings.HasPrefix(h.String(), hashStr) {
			hashes = append(hashes, h)
		}
	}
	return hashes
}

type RepackConfig struct {
	// UseRefDeltas configures whether packfile encoder will use reference deltas.
	// By default OFSDeltaObject is used.
	UseRefDeltas bool
	// OnlyDeletePacksOlderThan if set to non-zero value
	// selects only objects older than the time provided.
	OnlyDeletePacksOlderThan time.Time
}

func (r *Repository) RepackObjects(cfg *RepackConfig) (err error) {
	pos, ok := r.Storer.(storer.PackedObjectStorer)
	if !ok {
		return ErrPackedObjectsNotSupported
	}

	// Get the existing object packs.
	hs, err := pos.ObjectPacks()
	if err != nil {
		return err
	}

	// Create a new pack.
	nh, err := r.createNewObjectPack(cfg)
	if err != nil {
		return err
	}

	// Delete old packs.
	for _, h := range hs {
		// Skip if new hash is the same as an old one.
		if h == nh {
			continue
		}
		err = pos.DeleteOldObjectPackAndIndex(h, cfg.OnlyDeletePacksOlderThan)
		if err != nil {
			return err
		}
	}

	return nil
}

// Merge merges the reference branch into the current branch.
//
// If the merge is not possible (or supported) returns an error without changing
// the HEAD for the current branch. Possible errors include:
//   - The merge strategy is not supported.
//   - The specific strategy cannot be used (e.g. using FastForwardMerge when one is not possible).
func (r *Repository) Merge(ref plumbing.Reference, opts MergeOptions) error {
	if opts.Strategy != FastForwardMerge {
		return ErrUnsupportedMergeStrategy
	}

	// Ignore error as not having a shallow list is optional here.
	shallowList, _ := r.Storer.Shallow()
	var earliestShallow *plumbing.Hash
	if len(shallowList) > 0 {
		earliestShallow = &shallowList[0]
	}

	head, err := r.Head()
	if err != nil {
		return err
	}

	ff, err := isFastForward(r.Storer, head.Hash(), ref.Hash(), earliestShallow)
	if err != nil {
		return err
	}

	if !ff {
		return ErrFastForwardMergeNotPossible
	}

	return r.Storer.SetReference(plumbing.NewHashReference(head.Name(), ref.Hash()))
}

// createNewObjectPack is a helper for RepackObjects taking care
// of creating a new pack. It is used so the PackfileWriter
// deferred close has the right scope.
func (r *Repository) createNewObjectPack(cfg *RepackConfig) (h plumbing.Hash, err error) {
	ow := newObjectWalker(r.Storer)
	err = ow.walkAllRefs()
	if err != nil {
		return h, err
	}
	objs := make([]plumbing.Hash, 0, len(ow.seen))
	for h := range ow.seen {
		objs = append(objs, h)
	}
	pfw, ok := r.Storer.(storer.PackfileWriter)
	if !ok {
		return h, fmt.Errorf("Repository storer is not a storer.PackfileWriter")
	}
	wc, err := pfw.PackfileWriter()
	if err != nil {
		return h, err
	}
	defer ioutil.CheckClose(wc, &err)
	scfg, err := r.Config()
	if err != nil {
		return h, err
	}
	enc := packfile.NewEncoder(wc, r.Storer, cfg.UseRefDeltas)
	h, err = enc.Encode(objs, scfg.Pack.Window)
	if err != nil {
		return h, err
	}

	// Delete the packed, loose objects.
	if los, ok := r.Storer.(storer.LooseObjectStorer); ok {
		err = los.ForEachObjectHash(func(hash plumbing.Hash) error {
			if ow.isSeen(hash) {
				err = los.DeleteLooseObject(hash)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return h, err
		}
	}

	return h, err
}

func expandPartialHash(st storer.EncodedObjectStorer, prefix []byte) (hashes []plumbing.Hash) {
	// The fast version is implemented by storage/filesystem.ObjectStorage.
	type fastIter interface {
		HashesWithPrefix(prefix []byte) ([]plumbing.Hash, error)
	}
	if fi, ok := st.(fastIter); ok {
		h, err := fi.HashesWithPrefix(prefix)
		if err != nil {
			return nil
		}
		return h
	}

	// Slow path.
	iter, err := st.IterEncodedObjects(plumbing.AnyObject)
	if err != nil {
		return nil
	}
	iter.ForEach(func(obj plumbing.EncodedObject) error {
		h := obj.Hash()
		if bytes.HasPrefix(h[:], prefix) {
			hashes = append(hashes, h)
		}
		return nil
	})
	return
}
