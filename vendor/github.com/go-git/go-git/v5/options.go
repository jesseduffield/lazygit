package git

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	formatcfg "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// SubmoduleRescursivity defines how depth will affect any submodule recursive
// operation.
type SubmoduleRescursivity uint

const (
	// DefaultRemoteName name of the default Remote, just like git command.
	DefaultRemoteName = "origin"

	// NoRecurseSubmodules disables the recursion for a submodule operation.
	NoRecurseSubmodules SubmoduleRescursivity = 0
	// DefaultSubmoduleRecursionDepth allow recursion in a submodule operation.
	DefaultSubmoduleRecursionDepth SubmoduleRescursivity = 10
)

var (
	ErrMissingURL = errors.New("URL field is required")
)

// CloneOptions describes how a clone should be performed.
type CloneOptions struct {
	// The (possibly remote) repository URL to clone from.
	URL string
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// Name of the remote to be added, by default `origin`.
	RemoteName string
	// Remote branch to clone.
	ReferenceName plumbing.ReferenceName
	// Fetch only ReferenceName if true.
	SingleBranch bool
	// Mirror clones the repository as a mirror.
	//
	// Compared to a bare clone, mirror not only maps local branches of the
	// source to local branches of the target, it maps all refs (including
	// remote-tracking branches, notes etc.) and sets up a refspec configuration
	// such that all these refs are overwritten by a git remote update in the
	// target repository.
	Mirror bool
	// No checkout of HEAD after clone if true.
	NoCheckout bool
	// Limit fetching to the specified number of commits.
	Depth int
	// RecurseSubmodules after the clone is created, initialize all submodules
	// within, using their default settings. This option is ignored if the
	// cloned repository does not have a worktree.
	RecurseSubmodules SubmoduleRescursivity
	// ShallowSubmodules limit cloning submodules to the 1 level of depth.
	// It matches the git command --shallow-submodules.
	ShallowSubmodules bool
	// Progress is where the human readable information sent by the server is
	// stored, if nil nothing is stored and the capability (if supported)
	// no-progress, is sent to the server to avoid send this information.
	Progress sideband.Progress
	// Tags describe how the tags will be fetched from the remote repository,
	// by default is AllTags.
	Tags TagMode
	// InsecureSkipTLS skips ssl verify if protocol is https
	InsecureSkipTLS bool
	// CABundle specify additional ca bundle with system cert pool
	CABundle []byte
	// ProxyOptions provides info required for connecting to a proxy.
	ProxyOptions transport.ProxyOptions
	// When the repository to clone is on the local machine, instead of
	// using hard links, automatically setup .git/objects/info/alternates
	// to share the objects with the source repository.
	// The resulting repository starts out without any object of its own.
	// NOTE: this is a possibly dangerous operation; do not use it unless
	// you understand what it does.
	//
	// [Reference]: https://git-scm.com/docs/git-clone#Documentation/git-clone.txt---shared
	Shared bool
}

// MergeOptions describes how a merge should be performed.
type MergeOptions struct {
	// Strategy defines the merge strategy to be used.
	Strategy MergeStrategy
}

// MergeStrategy represents the different types of merge strategies.
type MergeStrategy int8

const (
	// FastForwardMerge represents a Git merge strategy where the current
	// branch can be simply updated to point to the HEAD of the branch being
	// merged. This is only possible if the history of the branch being merged
	// is a linear descendant of the current branch, with no conflicting commits.
	//
	// This is the default option.
	FastForwardMerge MergeStrategy = iota
)

// Validate validates the fields and sets the default values.
func (o *CloneOptions) Validate() error {
	if o.URL == "" {
		return ErrMissingURL
	}

	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	if o.ReferenceName == "" {
		o.ReferenceName = plumbing.HEAD
	}

	if o.Tags == InvalidTagMode {
		o.Tags = AllTags
	}

	return nil
}

// PullOptions describes how a pull should be performed.
type PullOptions struct {
	// Name of the remote to be pulled. If empty, uses the default.
	RemoteName string
	// RemoteURL overrides the remote repo address with a custom URL
	RemoteURL string
	// Remote branch to clone. If empty, uses HEAD.
	ReferenceName plumbing.ReferenceName
	// Fetch only ReferenceName if true.
	SingleBranch bool
	// Limit fetching to the specified number of commits.
	Depth int
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// RecurseSubmodules controls if new commits of all populated submodules
	// should be fetched too.
	RecurseSubmodules SubmoduleRescursivity
	// Progress is where the human readable information sent by the server is
	// stored, if nil nothing is stored and the capability (if supported)
	// no-progress, is sent to the server to avoid send this information.
	Progress sideband.Progress
	// Force allows the pull to update a local branch even when the remote
	// branch does not descend from it.
	Force bool
	// InsecureSkipTLS skips ssl verify if protocol is https
	InsecureSkipTLS bool
	// CABundle specify additional ca bundle with system cert pool
	CABundle []byte
	// ProxyOptions provides info required for connecting to a proxy.
	ProxyOptions transport.ProxyOptions
}

// Validate validates the fields and sets the default values.
func (o *PullOptions) Validate() error {
	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	if o.ReferenceName == "" {
		o.ReferenceName = plumbing.HEAD
	}

	return nil
}

type TagMode int

const (
	InvalidTagMode TagMode = iota
	// TagFollowing any tag that points into the histories being fetched is also
	// fetched. TagFollowing requires a server with `include-tag` capability
	// in order to fetch the annotated tags objects.
	TagFollowing
	// AllTags fetch all tags from the remote (i.e., fetch remote tags
	// refs/tags/* into local tags with the same name)
	AllTags
	// NoTags fetch no tags from the remote at all
	NoTags
)

// FetchOptions describes how a fetch should be performed
type FetchOptions struct {
	// Name of the remote to fetch from. Defaults to origin.
	RemoteName string
	// RemoteURL overrides the remote repo address with a custom URL
	RemoteURL string
	RefSpecs  []config.RefSpec
	// Depth limit fetching to the specified number of commits from the tip of
	// each remote branch history.
	Depth int
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// Progress is where the human readable information sent by the server is
	// stored, if nil nothing is stored and the capability (if supported)
	// no-progress, is sent to the server to avoid send this information.
	Progress sideband.Progress
	// Tags describe how the tags will be fetched from the remote repository,
	// by default is TagFollowing.
	Tags TagMode
	// Force allows the fetch to update a local branch even when the remote
	// branch does not descend from it.
	Force bool
	// InsecureSkipTLS skips ssl verify if protocol is https
	InsecureSkipTLS bool
	// CABundle specify additional ca bundle with system cert pool
	CABundle []byte
	// ProxyOptions provides info required for connecting to a proxy.
	ProxyOptions transport.ProxyOptions
	// Prune specify that local refs that match given RefSpecs and that do
	// not exist remotely will be removed.
	Prune bool
}

// Validate validates the fields and sets the default values.
func (o *FetchOptions) Validate() error {
	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	if o.Tags == InvalidTagMode {
		o.Tags = TagFollowing
	}

	for _, r := range o.RefSpecs {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// PushOptions describes how a push should be performed.
type PushOptions struct {
	// RemoteName is the name of the remote to be pushed to.
	RemoteName string
	// RemoteURL overrides the remote repo address with a custom URL
	RemoteURL string
	// RefSpecs specify what destination ref to update with what source object.
	//
	// The format of a <refspec> parameter is an optional plus +, followed by
	//  the source object <src>, followed by a colon :, followed by the destination ref <dst>.
	// The <src> is often the name of the branch you would want to push, but it can be a SHA-1.
	// The <dst> tells which ref on the remote side is updated with this push.
	//
	// A refspec with empty src can be used to delete a reference.
	RefSpecs []config.RefSpec
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// Progress is where the human readable information sent by the server is
	// stored, if nil nothing is stored.
	Progress sideband.Progress
	// Prune specify that remote refs that match given RefSpecs and that do
	// not exist locally will be removed.
	Prune bool
	// Force allows the push to update a remote branch even when the local
	// branch does not descend from it.
	Force bool
	// InsecureSkipTLS skips ssl verify if protocol is https
	InsecureSkipTLS bool
	// CABundle specify additional ca bundle with system cert pool
	CABundle []byte
	// RequireRemoteRefs only allows a remote ref to be updated if its current
	// value is the one specified here.
	RequireRemoteRefs []config.RefSpec
	// FollowTags will send any annotated tags with a commit target reachable from
	// the refs already being pushed
	FollowTags bool
	// ForceWithLease allows a force push as long as the remote ref adheres to a "lease"
	ForceWithLease *ForceWithLease
	// PushOptions sets options to be transferred to the server during push.
	Options map[string]string
	// Atomic sets option to be an atomic push
	Atomic bool
	// ProxyOptions provides info required for connecting to a proxy.
	ProxyOptions transport.ProxyOptions
}

// ForceWithLease sets fields on the lease
// If neither RefName nor Hash are set, ForceWithLease protects
// all refs in the refspec by ensuring the ref of the remote in the local repsitory
// matches the one in the ref advertisement.
type ForceWithLease struct {
	// RefName, when set will protect the ref by ensuring it matches the
	// hash in the ref advertisement.
	RefName plumbing.ReferenceName
	// Hash is the expected object id of RefName. The push will be rejected unless this
	// matches the corresponding object id of RefName in the refs advertisement.
	Hash plumbing.Hash
}

// Validate validates the fields and sets the default values.
func (o *PushOptions) Validate() error {
	if o.RemoteName == "" {
		o.RemoteName = DefaultRemoteName
	}

	if len(o.RefSpecs) == 0 {
		o.RefSpecs = []config.RefSpec{
			config.RefSpec(config.DefaultPushRefSpec),
		}
	}

	for _, r := range o.RefSpecs {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// SubmoduleUpdateOptions describes how a submodule update should be performed.
type SubmoduleUpdateOptions struct {
	// Init, if true initializes the submodules recorded in the index.
	Init bool
	// NoFetch tell to the update command to not fetch new objects from the
	// remote site.
	NoFetch bool
	// RecurseSubmodules the update is performed not only in the submodules of
	// the current repository but also in any nested submodules inside those
	// submodules (and so on). Until the SubmoduleRescursivity is reached.
	RecurseSubmodules SubmoduleRescursivity
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// Depth limit fetching to the specified number of commits from the tip of
	// each remote branch history.
	Depth int
}

var (
	ErrBranchHashExclusive  = errors.New("Branch and Hash are mutually exclusive")
	ErrCreateRequiresBranch = errors.New("Branch is mandatory when Create is used")
)

// CheckoutOptions describes how a checkout operation should be performed.
type CheckoutOptions struct {
	// Hash is the hash of a commit or tag to be checked out. If used, HEAD
	// will be in detached mode. If Create is not used, Branch and Hash are
	// mutually exclusive.
	Hash plumbing.Hash
	// Branch to be checked out, if Branch and Hash are empty is set to `master`.
	Branch plumbing.ReferenceName
	// Create a new branch named Branch and start it at Hash.
	Create bool
	// Force, if true when switching branches, proceed even if the index or the
	// working tree differs from HEAD. This is used to throw away local changes
	Force bool
	// Keep, if true when switching branches, local changes (the index or the
	// working tree changes) will be kept so that they can be committed to the
	// target branch. Force and Keep are mutually exclusive, should not be both
	// set to true.
	Keep bool
	// SparseCheckoutDirectories
	SparseCheckoutDirectories []string
}

// Validate validates the fields and sets the default values.
func (o *CheckoutOptions) Validate() error {
	if !o.Create && !o.Hash.IsZero() && o.Branch != "" {
		return ErrBranchHashExclusive
	}

	if o.Create && o.Branch == "" {
		return ErrCreateRequiresBranch
	}

	if o.Branch == "" {
		o.Branch = plumbing.Master
	}

	return nil
}

// ResetMode defines the mode of a reset operation.
type ResetMode int8

const (
	// MixedReset resets the index but not the working tree (i.e., the changed
	// files are preserved but not marked for commit) and reports what has not
	// been updated. This is the default action.
	MixedReset ResetMode = iota
	// HardReset resets the index and working tree. Any changes to tracked files
	// in the working tree are discarded.
	HardReset
	// MergeReset resets the index and updates the files in the working tree
	// that are different between Commit and HEAD, but keeps those which are
	// different between the index and working tree (i.e. which have changes
	// which have not been added).
	//
	// If a file that is different between Commit and the index has unstaged
	// changes, reset is aborted.
	MergeReset
	// SoftReset does not touch the index file or the working tree at all (but
	// resets the head to <commit>, just like all modes do). This leaves all
	// your changed files "Changes to be committed", as git status would put it.
	SoftReset
)

// ResetOptions describes how a reset operation should be performed.
type ResetOptions struct {
	// Commit, if commit is present set the current branch head (HEAD) to it.
	Commit plumbing.Hash
	// Mode, form resets the current branch head to Commit and possibly updates
	// the index (resetting it to the tree of Commit) and the working tree
	// depending on Mode. If empty MixedReset is used.
	Mode ResetMode
	// Files, if not empty will constrain the reseting the index to only files
	// specified in this list.
	Files []string
}

// Validate validates the fields and sets the default values.
func (o *ResetOptions) Validate(r *Repository) error {
	if o.Commit == plumbing.ZeroHash {
		ref, err := r.Head()
		if err != nil {
			return err
		}

		o.Commit = ref.Hash()
	} else {
		_, err := r.CommitObject(o.Commit)
		if err != nil {
			return fmt.Errorf("invalid reset option: %w", err)
		}
	}

	return nil
}

type LogOrder int8

const (
	LogOrderDefault LogOrder = iota
	LogOrderDFS
	LogOrderDFSPost
	LogOrderBSF
	LogOrderCommitterTime
)

// LogOptions describes how a log action should be performed.
type LogOptions struct {
	// When the From option is set the log will only contain commits
	// reachable from it. If this option is not set, HEAD will be used as
	// the default From.
	From plumbing.Hash

	// The default traversal algorithm is Depth-first search
	// set Order=LogOrderCommitterTime for ordering by committer time (more compatible with `git log`)
	// set Order=LogOrderBSF for Breadth-first search
	Order LogOrder

	// Show only those commits in which the specified file was inserted/updated.
	// It is equivalent to running `git log -- <file-name>`.
	// this field is kept for compatibility, it can be replaced with PathFilter
	FileName *string

	// Filter commits based on the path of files that are updated
	// takes file path as argument and should return true if the file is desired
	// It can be used to implement `git log -- <path>`
	// either <path> is a file path, or directory path, or a regexp of file/directory path
	PathFilter func(string) bool

	// Pretend as if all the refs in refs/, along with HEAD, are listed on the command line as <commit>.
	// It is equivalent to running `git log --all`.
	// If set on true, the From option will be ignored.
	All bool

	// Show commits more recent than a specific date.
	// It is equivalent to running `git log --since <date>` or `git log --after <date>`.
	Since *time.Time

	// Show commits older than a specific date.
	// It is equivalent to running `git log --until <date>` or `git log --before <date>`.
	Until *time.Time
}

var (
	ErrMissingAuthor = errors.New("author field is required")
)

// AddOptions describes how an `add` operation should be performed
type AddOptions struct {
	// All equivalent to `git add -A`, update the index not only where the
	// working tree has a file matching `Path` but also where the index already
	// has an entry. This adds, modifies, and removes index entries to match the
	// working tree.  If no `Path` nor `Glob` is given when `All` option is
	// used, all files in the entire working tree are updated.
	All bool
	// Path is the exact filepath to the file or directory to be added.
	Path string
	// Glob adds all paths, matching pattern, to the index. If pattern matches a
	// directory path, all directory contents are added to the index recursively.
	Glob string
	// SkipStatus adds the path with no status check. This option is relevant only
	// when the `Path` option is specified and does not apply when the `All` option is used.
	// Notice that when passing an ignored path it will be added anyway.
	// When true it can speed up adding files to the worktree in very large repositories.
	SkipStatus bool
}

// Validate validates the fields and sets the default values.
func (o *AddOptions) Validate(r *Repository) error {
	if o.Path != "" && o.Glob != "" {
		return fmt.Errorf("fields Path and Glob are mutual exclusive")
	}

	return nil
}

// CommitOptions describes how a commit operation should be performed.
type CommitOptions struct {
	// All automatically stage files that have been modified and deleted, but
	// new files you have not told Git about are not affected.
	All bool
	// AllowEmptyCommits enable empty commits to be created. An empty commit
	// is when no changes to the tree were made, but a new commit message is
	// provided. The default behavior is false, which results in ErrEmptyCommit.
	AllowEmptyCommits bool
	// Author is the author's signature of the commit. If Author is empty the
	// Name and Email is read from the config, and time.Now it's used as When.
	Author *object.Signature
	// Committer is the committer's signature of the commit. If Committer is
	// nil the Author signature is used.
	Committer *object.Signature
	// Parents are the parents commits for the new commit, by default when
	// len(Parents) is zero, the hash of HEAD reference is used.
	Parents []plumbing.Hash
	// SignKey denotes a key to sign the commit with. A nil value here means the
	// commit will not be signed. The private key must be present and already
	// decrypted.
	SignKey *openpgp.Entity
	// Signer denotes a cryptographic signer to sign the commit with.
	// A nil value here means the commit will not be signed.
	// Takes precedence over SignKey.
	Signer Signer
	// Amend will create a new commit object and replace the commit that HEAD currently
	// points to. Cannot be used with All nor Parents.
	Amend bool
}

// Validate validates the fields and sets the default values.
func (o *CommitOptions) Validate(r *Repository) error {
	if o.All && o.Amend {
		return errors.New("all and amend cannot be used together")
	}

	if o.Amend && len(o.Parents) > 0 {
		return errors.New("parents cannot be used with amend")
	}

	if o.Author == nil {
		if err := o.loadConfigAuthorAndCommitter(r); err != nil {
			return err
		}
	}

	if o.Committer == nil {
		o.Committer = o.Author
	}

	if len(o.Parents) == 0 {
		head, err := r.Head()
		if err != nil && err != plumbing.ErrReferenceNotFound {
			return err
		}

		if head != nil {
			o.Parents = []plumbing.Hash{head.Hash()}
		}
	}

	return nil
}

func (o *CommitOptions) loadConfigAuthorAndCommitter(r *Repository) error {
	cfg, err := r.ConfigScoped(config.SystemScope)
	if err != nil {
		return err
	}

	if o.Author == nil && cfg.Author.Email != "" && cfg.Author.Name != "" {
		o.Author = &object.Signature{
			Name:  cfg.Author.Name,
			Email: cfg.Author.Email,
			When:  time.Now(),
		}
	}

	if o.Committer == nil && cfg.Committer.Email != "" && cfg.Committer.Name != "" {
		o.Committer = &object.Signature{
			Name:  cfg.Committer.Name,
			Email: cfg.Committer.Email,
			When:  time.Now(),
		}
	}

	if o.Author == nil && cfg.User.Email != "" && cfg.User.Name != "" {
		o.Author = &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		}
	}

	if o.Author == nil {
		return ErrMissingAuthor
	}

	return nil
}

var (
	ErrMissingName    = errors.New("name field is required")
	ErrMissingTagger  = errors.New("tagger field is required")
	ErrMissingMessage = errors.New("message field is required")
)

// CreateTagOptions describes how a tag object should be created.
type CreateTagOptions struct {
	// Tagger defines the signature of the tag creator. If Tagger is empty the
	// Name and Email is read from the config, and time.Now it's used as When.
	Tagger *object.Signature
	// Message defines the annotation of the tag. It is canonicalized during
	// validation into the format expected by git - no leading whitespace and
	// ending in a newline.
	Message string
	// SignKey denotes a key to sign the tag with. A nil value here means the tag
	// will not be signed. The private key must be present and already decrypted.
	SignKey *openpgp.Entity
}

// Validate validates the fields and sets the default values.
func (o *CreateTagOptions) Validate(r *Repository, hash plumbing.Hash) error {
	if o.Tagger == nil {
		if err := o.loadConfigTagger(r); err != nil {
			return err
		}
	}

	if o.Message == "" {
		return ErrMissingMessage
	}

	// Canonicalize the message into the expected message format.
	o.Message = strings.TrimSpace(o.Message) + "\n"

	return nil
}

func (o *CreateTagOptions) loadConfigTagger(r *Repository) error {
	cfg, err := r.ConfigScoped(config.SystemScope)
	if err != nil {
		return err
	}

	if o.Tagger == nil && cfg.Author.Email != "" && cfg.Author.Name != "" {
		o.Tagger = &object.Signature{
			Name:  cfg.Author.Name,
			Email: cfg.Author.Email,
			When:  time.Now(),
		}
	}

	if o.Tagger == nil && cfg.User.Email != "" && cfg.User.Name != "" {
		o.Tagger = &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		}
	}

	if o.Tagger == nil {
		return ErrMissingTagger
	}

	return nil
}

// ListOptions describes how a remote list should be performed.
type ListOptions struct {
	// Auth credentials, if required, to use with the remote repository.
	Auth transport.AuthMethod
	// InsecureSkipTLS skips ssl verify if protocol is https
	InsecureSkipTLS bool
	// CABundle specify additional ca bundle with system cert pool
	CABundle []byte
	// PeelingOption defines how peeled objects are handled during a
	// remote list.
	PeelingOption PeelingOption
	// ProxyOptions provides info required for connecting to a proxy.
	ProxyOptions transport.ProxyOptions
	// Timeout specifies the timeout in seconds for list operations
	Timeout int
}

// PeelingOption represents the different ways to handle peeled references.
//
// Peeled references represent the underlying object of an annotated
// (or signed) tag. Refer to upstream documentation for more info:
// https://github.com/git/git/blob/master/Documentation/technical/reftable.txt
type PeelingOption uint8

const (
	// IgnorePeeled ignores all peeled reference names. This is the default behavior.
	IgnorePeeled PeelingOption = 0
	// OnlyPeeled returns only peeled reference names.
	OnlyPeeled PeelingOption = 1
	// AppendPeeled appends peeled reference names to the reference list.
	AppendPeeled PeelingOption = 2
)

// CleanOptions describes how a clean should be performed.
type CleanOptions struct {
	Dir bool
}

// GrepOptions describes how a grep should be performed.
type GrepOptions struct {
	// Patterns are compiled Regexp objects to be matched.
	Patterns []*regexp.Regexp
	// InvertMatch selects non-matching lines.
	InvertMatch bool
	// CommitHash is the hash of the commit from which worktree should be derived.
	CommitHash plumbing.Hash
	// ReferenceName is the branch or tag name from which worktree should be derived.
	ReferenceName plumbing.ReferenceName
	// PathSpecs are compiled Regexp objects of pathspec to use in the matching.
	PathSpecs []*regexp.Regexp
}

var (
	ErrHashOrReference = errors.New("ambiguous options, only one of CommitHash or ReferenceName can be passed")
)

// Validate validates the fields and sets the default values.
//
// TODO: deprecate in favor of Validate(r *Repository) in v6.
func (o *GrepOptions) Validate(w *Worktree) error {
	return o.validate(w.r)
}

func (o *GrepOptions) validate(r *Repository) error {
	if !o.CommitHash.IsZero() && o.ReferenceName != "" {
		return ErrHashOrReference
	}

	// If none of CommitHash and ReferenceName are provided, set commit hash of
	// the repository's head.
	if o.CommitHash.IsZero() && o.ReferenceName == "" {
		ref, err := r.Head()
		if err != nil {
			return err
		}
		o.CommitHash = ref.Hash()
	}

	return nil
}

// PlainOpenOptions describes how opening a plain repository should be
// performed.
type PlainOpenOptions struct {
	// DetectDotGit defines whether parent directories should be
	// walked until a .git directory or file is found.
	DetectDotGit bool
	// Enable .git/commondir support (see https://git-scm.com/docs/gitrepository-layout#Documentation/gitrepository-layout.txt).
	// NOTE: This option will only work with the filesystem storage.
	EnableDotGitCommonDir bool
}

// Validate validates the fields and sets the default values.
func (o *PlainOpenOptions) Validate() error { return nil }

type PlainInitOptions struct {
	InitOptions
	// Determines if the repository will have a worktree (non-bare) or not (bare).
	Bare         bool
	ObjectFormat formatcfg.ObjectFormat
}

// Validate validates the fields and sets the default values.
func (o *PlainInitOptions) Validate() error { return nil }

var (
	ErrNoRestorePaths = errors.New("you must specify path(s) to restore")
)

// RestoreOptions describes how a restore should be performed.
type RestoreOptions struct {
	// Marks to restore the content in the index
	Staged bool
	// Marks to restore the content of the working tree
	Worktree bool
	// List of file paths that will be restored
	Files []string
}

// Validate validates the fields and sets the default values.
func (o *RestoreOptions) Validate() error {
	if len(o.Files) == 0 {
		return ErrNoRestorePaths
	}

	return nil
}
