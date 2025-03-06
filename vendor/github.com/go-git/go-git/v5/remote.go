package git

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/internal/url"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/go-git/go-git/v5/utils/ioutil"
)

var (
	NoErrAlreadyUpToDate     = errors.New("already up-to-date")
	ErrDeleteRefNotSupported = errors.New("server does not support delete-refs")
	ErrForceNeeded           = errors.New("some refs were not updated")
	ErrExactSHA1NotSupported = errors.New("server does not support exact SHA1 refspec")
	ErrEmptyUrls             = errors.New("URLs cannot be empty")
)

type NoMatchingRefSpecError struct {
	refSpec config.RefSpec
}

func (e NoMatchingRefSpecError) Error() string {
	return fmt.Sprintf("couldn't find remote ref %q", e.refSpec.Src())
}

func (e NoMatchingRefSpecError) Is(target error) bool {
	_, ok := target.(NoMatchingRefSpecError)
	return ok
}

const (
	// This describes the maximum number of commits to walk when
	// computing the haves to send to a server, for each ref in the
	// repo containing this remote, when not using the multi-ack
	// protocol.  Setting this to 0 means there is no limit.
	maxHavesToVisitPerRef = 100

	// peeledSuffix is the suffix used to build peeled reference names.
	peeledSuffix = "^{}"
)

// Remote represents a connection to a remote repository.
type Remote struct {
	c *config.RemoteConfig
	s storage.Storer
}

// NewRemote creates a new Remote.
// The intended purpose is to use the Remote for tasks such as listing remote references (like using git ls-remote).
// Otherwise Remotes should be created via the use of a Repository.
func NewRemote(s storage.Storer, c *config.RemoteConfig) *Remote {
	return &Remote{s: s, c: c}
}

// Config returns the RemoteConfig object used to instantiate this Remote.
func (r *Remote) Config() *config.RemoteConfig {
	return r.c
}

func (r *Remote) String() string {
	var fetch, push string
	if len(r.c.URLs) > 0 {
		fetch = r.c.URLs[0]
		push = r.c.URLs[len(r.c.URLs)-1]
	}

	return fmt.Sprintf("%s\t%s (fetch)\n%[1]s\t%[3]s (push)", r.c.Name, fetch, push)
}

// Push performs a push to the remote. Returns NoErrAlreadyUpToDate if the
// remote was already up-to-date.
func (r *Remote) Push(o *PushOptions) error {
	return r.PushContext(context.Background(), o)
}

// PushContext performs a push to the remote. Returns NoErrAlreadyUpToDate if
// the remote was already up-to-date.
//
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects the
// transport operations.
func (r *Remote) PushContext(ctx context.Context, o *PushOptions) (err error) {
	if err := o.Validate(); err != nil {
		return err
	}

	if o.RemoteName != r.c.Name {
		return fmt.Errorf("remote names don't match: %s != %s", o.RemoteName, r.c.Name)
	}

	if o.RemoteURL == "" && len(r.c.URLs) > 0 {
		o.RemoteURL = r.c.URLs[len(r.c.URLs)-1]
	}

	s, err := newSendPackSession(o.RemoteURL, o.Auth, o.InsecureSkipTLS, o.CABundle, o.ProxyOptions)
	if err != nil {
		return err
	}

	defer ioutil.CheckClose(s, &err)

	ar, err := s.AdvertisedReferencesContext(ctx)
	if err != nil {
		return err
	}

	remoteRefs, err := ar.AllReferences()
	if err != nil {
		return err
	}

	if err := r.checkRequireRemoteRefs(o.RequireRemoteRefs, remoteRefs); err != nil {
		return err
	}

	isDelete := false
	allDelete := true
	for _, rs := range o.RefSpecs {
		if rs.IsDelete() {
			isDelete = true
		} else {
			allDelete = false
		}
		if isDelete && !allDelete {
			break
		}
	}

	if isDelete && !ar.Capabilities.Supports(capability.DeleteRefs) {
		return ErrDeleteRefNotSupported
	}

	if o.Force {
		for i := 0; i < len(o.RefSpecs); i++ {
			rs := &o.RefSpecs[i]
			if !rs.IsForceUpdate() && !rs.IsDelete() {
				o.RefSpecs[i] = config.RefSpec("+" + rs.String())
			}
		}
	}

	localRefs, err := r.references()
	if err != nil {
		return err
	}

	req, err := r.newReferenceUpdateRequest(o, localRefs, remoteRefs, ar)
	if err != nil {
		return err
	}

	if len(req.Commands) == 0 {
		return NoErrAlreadyUpToDate
	}

	objects := objectsToPush(req.Commands)

	haves, err := referencesToHashes(remoteRefs)
	if err != nil {
		return err
	}

	stop, err := r.s.Shallow()
	if err != nil {
		return err
	}

	// if we have shallow we should include this as part of the objects that
	// we are aware.
	haves = append(haves, stop...)

	var hashesToPush []plumbing.Hash
	// Avoid the expensive revlist operation if we're only doing deletes.
	if !allDelete {
		if url.IsLocalEndpoint(o.RemoteURL) {
			// If we're are pushing to a local repo, it might be much
			// faster to use a local storage layer to get the commits
			// to ignore, when calculating the object revlist.
			localStorer := filesystem.NewStorage(
				osfs.New(o.RemoteURL), cache.NewObjectLRUDefault())
			hashesToPush, err = revlist.ObjectsWithStorageForIgnores(
				r.s, localStorer, objects, haves)
		} else {
			hashesToPush, err = revlist.Objects(r.s, objects, haves)
		}
		if err != nil {
			return err
		}
	}

	if len(hashesToPush) == 0 {
		allDelete = true
		for _, command := range req.Commands {
			if command.Action() != packp.Delete {
				allDelete = false
				break
			}
		}
	}

	rs, err := pushHashes(ctx, s, r.s, req, hashesToPush, r.useRefDeltas(ar), allDelete)
	if err != nil {
		return err
	}

	if rs != nil {
		if err = rs.Error(); err != nil {
			return err
		}
	}

	return r.updateRemoteReferenceStorage(req)
}

func (r *Remote) useRefDeltas(ar *packp.AdvRefs) bool {
	return !ar.Capabilities.Supports(capability.OFSDelta)
}

func (r *Remote) addReachableTags(localRefs []*plumbing.Reference, remoteRefs storer.ReferenceStorer, req *packp.ReferenceUpdateRequest) error {
	tags := make(map[plumbing.Reference]struct{})
	// get a list of all tags locally
	for _, ref := range localRefs {
		if strings.HasPrefix(string(ref.Name()), "refs/tags") {
			tags[*ref] = struct{}{}
		}
	}

	remoteRefIter, err := remoteRefs.IterReferences()
	if err != nil {
		return err
	}

	// remove any that are already on the remote
	if err := remoteRefIter.ForEach(func(reference *plumbing.Reference) error {
		delete(tags, *reference)
		return nil
	}); err != nil {
		return err
	}

	for tag := range tags {
		tagObject, err := object.GetObject(r.s, tag.Hash())
		var tagCommit *object.Commit
		if err != nil {
			return fmt.Errorf("get tag object: %w", err)
		}

		if tagObject.Type() != plumbing.TagObject {
			continue
		}

		annotatedTag, ok := tagObject.(*object.Tag)
		if !ok {
			return errors.New("could not get annotated tag object")
		}

		tagCommit, err = object.GetCommit(r.s, annotatedTag.Target)
		if err != nil {
			return fmt.Errorf("get annotated tag commit: %w", err)
		}

		// only include tags that are reachable from one of the refs
		// already being pushed
		for _, cmd := range req.Commands {
			if tag.Name() == cmd.Name {
				continue
			}

			if strings.HasPrefix(cmd.Name.String(), "refs/tags") {
				continue
			}

			c, err := object.GetCommit(r.s, cmd.New)
			if err != nil {
				return fmt.Errorf("get commit %v: %w", cmd.Name, err)
			}

			if isAncestor, err := tagCommit.IsAncestor(c); err == nil && isAncestor {
				req.Commands = append(req.Commands, &packp.Command{Name: tag.Name(), New: tag.Hash()})
			}
		}
	}

	return nil
}

func (r *Remote) newReferenceUpdateRequest(
	o *PushOptions,
	localRefs []*plumbing.Reference,
	remoteRefs storer.ReferenceStorer,
	ar *packp.AdvRefs,
) (*packp.ReferenceUpdateRequest, error) {
	req := packp.NewReferenceUpdateRequestFromCapabilities(ar.Capabilities)

	if o.Progress != nil {
		req.Progress = o.Progress
		if ar.Capabilities.Supports(capability.Sideband64k) {
			_ = req.Capabilities.Set(capability.Sideband64k)
		} else if ar.Capabilities.Supports(capability.Sideband) {
			_ = req.Capabilities.Set(capability.Sideband)
		}
	}

	if ar.Capabilities.Supports(capability.PushOptions) {
		_ = req.Capabilities.Set(capability.PushOptions)
		for k, v := range o.Options {
			req.Options = append(req.Options, &packp.Option{Key: k, Value: v})
		}
	}

	if o.Atomic && ar.Capabilities.Supports(capability.Atomic) {
		_ = req.Capabilities.Set(capability.Atomic)
	}

	if err := r.addReferencesToUpdate(o.RefSpecs, localRefs, remoteRefs, req, o.Prune, o.ForceWithLease); err != nil {

		return nil, err
	}

	if o.FollowTags {
		if err := r.addReachableTags(localRefs, remoteRefs, req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (r *Remote) updateRemoteReferenceStorage(
	req *packp.ReferenceUpdateRequest,
) error {

	for _, spec := range r.c.Fetch {
		for _, c := range req.Commands {
			if !spec.Match(c.Name) {
				continue
			}

			local := spec.Dst(c.Name)
			ref := plumbing.NewHashReference(local, c.New)
			switch c.Action() {
			case packp.Create, packp.Update:
				if err := r.s.SetReference(ref); err != nil {
					return err
				}
			case packp.Delete:
				if err := r.s.RemoveReference(local); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// FetchContext fetches references along with the objects necessary to complete
// their histories.
//
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
//
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects the
// transport operations.
func (r *Remote) FetchContext(ctx context.Context, o *FetchOptions) error {
	_, err := r.fetch(ctx, o)
	return err
}

// Fetch fetches references along with the objects necessary to complete their
// histories.
//
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
func (r *Remote) Fetch(o *FetchOptions) error {
	return r.FetchContext(context.Background(), o)
}

func (r *Remote) fetch(ctx context.Context, o *FetchOptions) (sto storer.ReferenceStorer, err error) {
	if o.RemoteName == "" {
		o.RemoteName = r.c.Name
	}

	if err = o.Validate(); err != nil {
		return nil, err
	}

	if len(o.RefSpecs) == 0 {
		o.RefSpecs = r.c.Fetch
	}

	if o.RemoteURL == "" {
		o.RemoteURL = r.c.URLs[0]
	}

	s, err := newUploadPackSession(o.RemoteURL, o.Auth, o.InsecureSkipTLS, o.CABundle, o.ProxyOptions)
	if err != nil {
		return nil, err
	}

	defer ioutil.CheckClose(s, &err)

	ar, err := s.AdvertisedReferencesContext(ctx)
	if err != nil {
		return nil, err
	}

	req, err := r.newUploadPackRequest(o, ar)
	if err != nil {
		return nil, err
	}

	if err := r.isSupportedRefSpec(o.RefSpecs, ar); err != nil {
		return nil, err
	}

	remoteRefs, err := ar.AllReferences()
	if err != nil {
		return nil, err
	}

	localRefs, err := r.references()
	if err != nil {
		return nil, err
	}

	refs, specToRefs, err := calculateRefs(o.RefSpecs, remoteRefs, o.Tags)
	if err != nil {
		return nil, err
	}

	if !req.Depth.IsZero() {
		req.Shallows, err = r.s.Shallow()
		if err != nil {
			return nil, fmt.Errorf("existing checkout is not shallow")
		}
	}

	req.Wants, err = getWants(r.s, refs, o.Depth)
	if len(req.Wants) > 0 {
		req.Haves, err = getHaves(localRefs, remoteRefs, r.s, o.Depth)
		if err != nil {
			return nil, err
		}

		if err = r.fetchPack(ctx, o, s, req); err != nil {
			return nil, err
		}
	}

	var updatedPrune bool
	if o.Prune {
		updatedPrune, err = r.pruneRemotes(o.RefSpecs, localRefs, remoteRefs)
		if err != nil {
			return nil, err
		}
	}

	updated, err := r.updateLocalReferenceStorage(o.RefSpecs, refs, remoteRefs, specToRefs, o.Tags, o.Force)
	if err != nil {
		return nil, err
	}

	if !updated {
		updated, err = depthChanged(req.Shallows, r.s)
		if err != nil {
			return nil, fmt.Errorf("error checking depth change: %v", err)
		}
	}

	if !updated && !updatedPrune {
		// No references updated, but may have fetched new objects, check if we now have any of our wants
		for _, hash := range req.Wants {
			exists, _ := objectExists(r.s, hash)
			if exists {
				updated = true
				break
			}
		}

		if !updated {
			return remoteRefs, NoErrAlreadyUpToDate
		}
	}

	return remoteRefs, nil
}

func depthChanged(before []plumbing.Hash, s storage.Storer) (bool, error) {
	after, err := s.Shallow()
	if err != nil {
		return false, err
	}

	if len(before) != len(after) {
		return true, nil
	}

	bm := make(map[plumbing.Hash]bool, len(before))
	for _, b := range before {
		bm[b] = true
	}
	for _, a := range after {
		if _, ok := bm[a]; !ok {
			return true, nil
		}
	}

	return false, nil
}

func newUploadPackSession(url string, auth transport.AuthMethod, insecure bool, cabundle []byte, proxyOpts transport.ProxyOptions) (transport.UploadPackSession, error) {
	c, ep, err := newClient(url, insecure, cabundle, proxyOpts)
	if err != nil {
		return nil, err
	}

	return c.NewUploadPackSession(ep, auth)
}

func newSendPackSession(url string, auth transport.AuthMethod, insecure bool, cabundle []byte, proxyOpts transport.ProxyOptions) (transport.ReceivePackSession, error) {
	c, ep, err := newClient(url, insecure, cabundle, proxyOpts)
	if err != nil {
		return nil, err
	}

	return c.NewReceivePackSession(ep, auth)
}

func newClient(url string, insecure bool, cabundle []byte, proxyOpts transport.ProxyOptions) (transport.Transport, *transport.Endpoint, error) {
	ep, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, nil, err
	}
	ep.InsecureSkipTLS = insecure
	ep.CaBundle = cabundle
	ep.Proxy = proxyOpts

	c, err := client.NewClient(ep)
	if err != nil {
		return nil, nil, err
	}

	return c, ep, err
}

func (r *Remote) fetchPack(ctx context.Context, o *FetchOptions, s transport.UploadPackSession,
	req *packp.UploadPackRequest) (err error) {

	reader, err := s.UploadPack(ctx, req)
	if err != nil {
		if errors.Is(err, transport.ErrEmptyUploadPackRequest) {
			// XXX: no packfile provided, everything is up-to-date.
			return nil
		}
		return err
	}

	defer ioutil.CheckClose(reader, &err)

	if err = r.updateShallow(o, reader); err != nil {
		return err
	}

	if err = packfile.UpdateObjectStorage(r.s,
		buildSidebandIfSupported(req.Capabilities, reader, o.Progress),
	); err != nil {
		return err
	}

	return err
}

func (r *Remote) pruneRemotes(specs []config.RefSpec, localRefs []*plumbing.Reference, remoteRefs memory.ReferenceStorage) (bool, error) {
	var updatedPrune bool
	for _, spec := range specs {
		rev := spec.Reverse()
		for _, ref := range localRefs {
			if !rev.Match(ref.Name()) {
				continue
			}
			_, err := remoteRefs.Reference(rev.Dst(ref.Name()))
			if errors.Is(err, plumbing.ErrReferenceNotFound) {
				updatedPrune = true
				err := r.s.RemoveReference(ref.Name())
				if err != nil {
					return false, err
				}
			}
		}
	}
	return updatedPrune, nil
}

func (r *Remote) addReferencesToUpdate(
	refspecs []config.RefSpec,
	localRefs []*plumbing.Reference,
	remoteRefs storer.ReferenceStorer,
	req *packp.ReferenceUpdateRequest,
	prune bool,
	forceWithLease *ForceWithLease,
) error {
	// This references dictionary will be used to search references by name.
	refsDict := make(map[string]*plumbing.Reference)
	for _, ref := range localRefs {
		refsDict[ref.Name().String()] = ref
	}

	for _, rs := range refspecs {
		if rs.IsDelete() {
			if err := r.deleteReferences(rs, remoteRefs, refsDict, req, false); err != nil {
				return err
			}
		} else {
			err := r.addOrUpdateReferences(rs, localRefs, refsDict, remoteRefs, req, forceWithLease)
			if err != nil {
				return err
			}

			if prune {
				if err := r.deleteReferences(rs, remoteRefs, refsDict, req, true); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *Remote) addOrUpdateReferences(
	rs config.RefSpec,
	localRefs []*plumbing.Reference,
	refsDict map[string]*plumbing.Reference,
	remoteRefs storer.ReferenceStorer,
	req *packp.ReferenceUpdateRequest,
	forceWithLease *ForceWithLease,
) error {
	// If it is not a wildcard refspec we can directly search for the reference
	// in the references dictionary.
	if !rs.IsWildcard() {
		ref, ok := refsDict[rs.Src()]
		if !ok {
			commit, err := object.GetCommit(r.s, plumbing.NewHash(rs.Src()))
			if err == nil {
				return r.addCommit(rs, remoteRefs, commit.Hash, req)
			}
			return nil
		}

		return r.addReferenceIfRefSpecMatches(rs, remoteRefs, ref, req, forceWithLease)
	}

	for _, ref := range localRefs {
		err := r.addReferenceIfRefSpecMatches(rs, remoteRefs, ref, req, forceWithLease)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Remote) deleteReferences(rs config.RefSpec,
	remoteRefs storer.ReferenceStorer,
	refsDict map[string]*plumbing.Reference,
	req *packp.ReferenceUpdateRequest,
	prune bool) error {
	iter, err := remoteRefs.IterReferences()
	if err != nil {
		return err
	}

	return iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}

		if prune {
			rs := rs.Reverse()
			if !rs.Match(ref.Name()) {
				return nil
			}

			if _, ok := refsDict[rs.Dst(ref.Name()).String()]; ok {
				return nil
			}
		} else if rs.Dst("") != ref.Name() {
			return nil
		}

		cmd := &packp.Command{
			Name: ref.Name(),
			Old:  ref.Hash(),
			New:  plumbing.ZeroHash,
		}
		req.Commands = append(req.Commands, cmd)
		return nil
	})
}

func (r *Remote) addCommit(rs config.RefSpec,
	remoteRefs storer.ReferenceStorer, localCommit plumbing.Hash,
	req *packp.ReferenceUpdateRequest) error {

	if rs.IsWildcard() {
		return errors.New("can't use wildcard together with hash refspecs")
	}

	cmd := &packp.Command{
		Name: rs.Dst(""),
		Old:  plumbing.ZeroHash,
		New:  localCommit,
	}
	remoteRef, err := remoteRefs.Reference(cmd.Name)
	if err == nil {
		if remoteRef.Type() != plumbing.HashReference {
			// TODO: check actual git behavior here
			return nil
		}

		cmd.Old = remoteRef.Hash()
	} else if err != plumbing.ErrReferenceNotFound {
		return err
	}
	if cmd.Old == cmd.New {
		return nil
	}
	if !rs.IsForceUpdate() {
		if err := checkFastForwardUpdate(r.s, remoteRefs, cmd); err != nil {
			return err
		}
	}

	req.Commands = append(req.Commands, cmd)
	return nil
}

func (r *Remote) addReferenceIfRefSpecMatches(rs config.RefSpec,
	remoteRefs storer.ReferenceStorer, localRef *plumbing.Reference,
	req *packp.ReferenceUpdateRequest, forceWithLease *ForceWithLease) error {

	if localRef.Type() != plumbing.HashReference {
		return nil
	}

	if !rs.Match(localRef.Name()) {
		return nil
	}

	cmd := &packp.Command{
		Name: rs.Dst(localRef.Name()),
		Old:  plumbing.ZeroHash,
		New:  localRef.Hash(),
	}

	remoteRef, err := remoteRefs.Reference(cmd.Name)
	if err == nil {
		if remoteRef.Type() != plumbing.HashReference {
			// TODO: check actual git behavior here
			return nil
		}

		cmd.Old = remoteRef.Hash()
	} else if err != plumbing.ErrReferenceNotFound {
		return err
	}

	if cmd.Old == cmd.New {
		return nil
	}

	if forceWithLease != nil {
		if err = r.checkForceWithLease(localRef, cmd, forceWithLease); err != nil {
			return err
		}
	} else if !rs.IsForceUpdate() {
		if err := checkFastForwardUpdate(r.s, remoteRefs, cmd); err != nil {
			return err
		}
	}

	req.Commands = append(req.Commands, cmd)
	return nil
}

func (r *Remote) checkForceWithLease(localRef *plumbing.Reference, cmd *packp.Command, forceWithLease *ForceWithLease) error {
	remotePrefix := fmt.Sprintf("refs/remotes/%s/", r.Config().Name)

	ref, err := storer.ResolveReference(
		r.s,
		plumbing.ReferenceName(remotePrefix+strings.Replace(localRef.Name().String(), "refs/heads/", "", -1)))
	if err != nil {
		return err
	}

	if forceWithLease.RefName.String() == "" || (forceWithLease.RefName == cmd.Name) {
		expectedOID := ref.Hash()

		if !forceWithLease.Hash.IsZero() {
			expectedOID = forceWithLease.Hash
		}

		if cmd.Old != expectedOID {
			return fmt.Errorf("non-fast-forward update: %s", cmd.Name.String())
		}
	}

	return nil
}

func (r *Remote) references() ([]*plumbing.Reference, error) {
	var localRefs []*plumbing.Reference

	iter, err := r.s.IterReferences()
	if err != nil {
		return nil, err
	}

	for {
		ref, err := iter.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		localRefs = append(localRefs, ref)
	}

	return localRefs, nil
}

func getRemoteRefsFromStorer(remoteRefStorer storer.ReferenceStorer) (
	map[plumbing.Hash]bool, error) {
	remoteRefs := map[plumbing.Hash]bool{}
	iter, err := remoteRefStorer.IterReferences()
	if err != nil {
		return nil, err
	}
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}
		remoteRefs[ref.Hash()] = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	return remoteRefs, nil
}

// getHavesFromRef populates the given `haves` map with the given
// reference, and up to `maxHavesToVisitPerRef` ancestor commits.
func getHavesFromRef(
	ref *plumbing.Reference,
	remoteRefs map[plumbing.Hash]bool,
	s storage.Storer,
	haves map[plumbing.Hash]bool,
	depth int,
) error {
	h := ref.Hash()
	if haves[h] {
		return nil
	}

	commit, err := object.GetCommit(s, h)
	if err != nil {
		if !errors.Is(err, plumbing.ErrObjectNotFound) {
			// Ignore the error if this isn't a commit.
			haves[ref.Hash()] = true
		}
		return nil
	}

	// Until go-git supports proper commit negotiation during an
	// upload pack request, include up to `maxHavesToVisitPerRef`
	// commits from the history of each ref.
	walker := object.NewCommitPreorderIter(commit, haves, nil)
	toVisit := maxHavesToVisitPerRef
	// But only need up to the requested depth
	if depth > 0 && depth < maxHavesToVisitPerRef {
		toVisit = depth
	}
	// It is safe to ignore any error here as we are just trying to find the references that we already have
	// An example of a legitimate failure is we have a shallow clone and don't have the previous commit(s)
	_ = walker.ForEach(func(c *object.Commit) error {
		haves[c.Hash] = true
		toVisit--
		// If toVisit starts out at 0 (indicating there is no
		// max), then it will be negative here and we won't stop
		// early.
		if toVisit == 0 || remoteRefs[c.Hash] {
			return storer.ErrStop
		}
		return nil
	})

	return nil
}

func getHaves(
	localRefs []*plumbing.Reference,
	remoteRefStorer storer.ReferenceStorer,
	s storage.Storer,
	depth int,
) ([]plumbing.Hash, error) {
	haves := map[plumbing.Hash]bool{}

	// Build a map of all the remote references, to avoid loading too
	// many parent commits for references we know don't need to be
	// transferred.
	remoteRefs, err := getRemoteRefsFromStorer(remoteRefStorer)
	if err != nil {
		return nil, err
	}

	for _, ref := range localRefs {
		if haves[ref.Hash()] {
			continue
		}

		if ref.Type() != plumbing.HashReference {
			continue
		}

		err = getHavesFromRef(ref, remoteRefs, s, haves, depth)
		if err != nil {
			return nil, err
		}
	}

	var result []plumbing.Hash
	for h := range haves {
		result = append(result, h)
	}

	return result, nil
}

const refspecAllTags = "+refs/tags/*:refs/tags/*"

func calculateRefs(
	spec []config.RefSpec,
	remoteRefs storer.ReferenceStorer,
	tagMode TagMode,
) (memory.ReferenceStorage, [][]*plumbing.Reference, error) {
	if tagMode == AllTags {
		spec = append(spec, refspecAllTags)
	}

	refs := make(memory.ReferenceStorage)
	// list of references matched for each spec
	specToRefs := make([][]*plumbing.Reference, len(spec))
	for i := range spec {
		var err error
		specToRefs[i], err = doCalculateRefs(spec[i], remoteRefs, refs)
		if err != nil {
			return nil, nil, err
		}
	}

	return refs, specToRefs, nil
}

func doCalculateRefs(
	s config.RefSpec,
	remoteRefs storer.ReferenceStorer,
	refs memory.ReferenceStorage,
) ([]*plumbing.Reference, error) {
	var refList []*plumbing.Reference

	if s.IsExactSHA1() {
		ref := plumbing.NewHashReference(s.Dst(""), plumbing.NewHash(s.Src()))

		refList = append(refList, ref)
		return refList, refs.SetReference(ref)
	}

	var matched bool
	onMatched := func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.SymbolicReference {
			target, err := storer.ResolveReference(remoteRefs, ref.Name())
			if err != nil {
				return err
			}

			ref = plumbing.NewHashReference(ref.Name(), target.Hash())
		}

		if ref.Type() != plumbing.HashReference {
			return nil
		}

		matched = true
		refList = append(refList, ref)
		return refs.SetReference(ref)
	}

	var ret error
	if s.IsWildcard() {
		iter, err := remoteRefs.IterReferences()
		if err != nil {
			return nil, err
		}
		ret = iter.ForEach(func(ref *plumbing.Reference) error {
			if !s.Match(ref.Name()) {
				return nil
			}

			return onMatched(ref)
		})
	} else {
		var resolvedRef *plumbing.Reference
		src := s.Src()
		resolvedRef, ret = expand_ref(remoteRefs, plumbing.ReferenceName(src))
		if ret == nil {
			ret = onMatched(resolvedRef)
		}
	}

	if !matched && !s.IsWildcard() {
		return nil, NoMatchingRefSpecError{refSpec: s}
	}

	return refList, ret
}

func getWants(localStorer storage.Storer, refs memory.ReferenceStorage, depth int) ([]plumbing.Hash, error) {
	// If depth is anything other than 1 and the repo has shallow commits then just because we have the commit
	// at the reference doesn't mean that we don't still need to fetch the parents
	shallow := false
	if depth != 1 {
		if s, _ := localStorer.Shallow(); len(s) > 0 {
			shallow = true
		}
	}

	wants := map[plumbing.Hash]bool{}
	for _, ref := range refs {
		hash := ref.Hash()
		exists, err := objectExists(localStorer, ref.Hash())
		if err != nil {
			return nil, err
		}

		if !exists || shallow {
			wants[hash] = true
		}
	}

	var result []plumbing.Hash
	for h := range wants {
		result = append(result, h)
	}

	return result, nil
}

func objectExists(s storer.EncodedObjectStorer, h plumbing.Hash) (bool, error) {
	_, err := s.EncodedObject(plumbing.AnyObject, h)
	if err == plumbing.ErrObjectNotFound {
		return false, nil
	}

	return true, err
}

func checkFastForwardUpdate(s storer.EncodedObjectStorer, remoteRefs storer.ReferenceStorer, cmd *packp.Command) error {
	if cmd.Old == plumbing.ZeroHash {
		_, err := remoteRefs.Reference(cmd.Name)
		if err == plumbing.ErrReferenceNotFound {
			return nil
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("non-fast-forward update: %s", cmd.Name.String())
	}

	ff, err := isFastForward(s, cmd.Old, cmd.New, nil)
	if err != nil {
		return err
	}

	if !ff {
		return fmt.Errorf("non-fast-forward update: %s", cmd.Name.String())
	}

	return nil
}

func isFastForward(s storer.EncodedObjectStorer, old, new plumbing.Hash, earliestShallow *plumbing.Hash) (bool, error) {
	c, err := object.GetCommit(s, new)
	if err != nil {
		return false, err
	}

	parentsToIgnore := []plumbing.Hash{}
	if earliestShallow != nil {
		earliestCommit, err := object.GetCommit(s, *earliestShallow)
		if err != nil {
			return false, err
		}

		parentsToIgnore = earliestCommit.ParentHashes
	}

	found := false
	// stop iterating at the earliest shallow commit, ignoring its parents
	// note: when pull depth is smaller than the number of new changes on the remote, this fails due to missing parents.
	//       as far as i can tell, without the commits in-between the shallow pull and the earliest shallow, there's no
	//       real way of telling whether it will be a fast-forward merge.
	iter := object.NewCommitPreorderIter(c, nil, parentsToIgnore)
	err = iter.ForEach(func(c *object.Commit) error {
		if c.Hash != old {
			return nil
		}

		found = true
		return storer.ErrStop
	})
	return found, err
}

func (r *Remote) newUploadPackRequest(o *FetchOptions,
	ar *packp.AdvRefs) (*packp.UploadPackRequest, error) {

	req := packp.NewUploadPackRequestFromCapabilities(ar.Capabilities)

	if o.Depth != 0 {
		req.Depth = packp.DepthCommits(o.Depth)
		if err := req.Capabilities.Set(capability.Shallow); err != nil {
			return nil, err
		}
	}

	if o.Progress == nil && ar.Capabilities.Supports(capability.NoProgress) {
		if err := req.Capabilities.Set(capability.NoProgress); err != nil {
			return nil, err
		}
	}

	isWildcard := true
	for _, s := range o.RefSpecs {
		if !s.IsWildcard() {
			isWildcard = false
			break
		}
	}

	if isWildcard && o.Tags == TagFollowing && ar.Capabilities.Supports(capability.IncludeTag) {
		if err := req.Capabilities.Set(capability.IncludeTag); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (r *Remote) isSupportedRefSpec(refs []config.RefSpec, ar *packp.AdvRefs) error {
	var containsIsExact bool
	for _, ref := range refs {
		if ref.IsExactSHA1() {
			containsIsExact = true
		}
	}

	if !containsIsExact {
		return nil
	}

	if ar.Capabilities.Supports(capability.AllowReachableSHA1InWant) ||
		ar.Capabilities.Supports(capability.AllowTipSHA1InWant) {
		return nil
	}

	return ErrExactSHA1NotSupported
}

func buildSidebandIfSupported(l *capability.List, reader io.Reader, p sideband.Progress) io.Reader {
	var t sideband.Type

	switch {
	case l.Supports(capability.Sideband):
		t = sideband.Sideband
	case l.Supports(capability.Sideband64k):
		t = sideband.Sideband64k
	default:
		return reader
	}

	d := sideband.NewDemuxer(t, reader)
	d.Progress = p

	return d
}

func (r *Remote) updateLocalReferenceStorage(
	specs []config.RefSpec,
	fetchedRefs, remoteRefs memory.ReferenceStorage,
	specToRefs [][]*plumbing.Reference,
	tagMode TagMode,
	force bool,
) (updated bool, err error) {
	isWildcard := true
	forceNeeded := false

	for i, spec := range specs {
		if !spec.IsWildcard() {
			isWildcard = false
		}

		for _, ref := range specToRefs[i] {
			if ref.Type() != plumbing.HashReference {
				continue
			}

			localName := spec.Dst(ref.Name())
			// If localName doesn't start with "refs/" then treat as a branch.
			if !strings.HasPrefix(localName.String(), "refs/") {
				localName = plumbing.NewBranchReferenceName(localName.String())
			}
			old, _ := storer.ResolveReference(r.s, localName)
			new := plumbing.NewHashReference(localName, ref.Hash())

			// If the ref exists locally as a non-tag and force is not
			// specified, only update if the new ref is an ancestor of the old
			if old != nil && !old.Name().IsTag() && !force && !spec.IsForceUpdate() {
				ff, err := isFastForward(r.s, old.Hash(), new.Hash(), nil)
				if err != nil {
					return updated, err
				}

				if !ff {
					forceNeeded = true
					continue
				}
			}

			refUpdated, err := checkAndUpdateReferenceStorerIfNeeded(r.s, new, old)
			if err != nil {
				return updated, err
			}

			if refUpdated {
				updated = true
			}
		}
	}

	if tagMode == NoTags {
		return updated, nil
	}

	tags := fetchedRefs
	if isWildcard {
		tags = remoteRefs
	}
	tagUpdated, err := r.buildFetchedTags(tags)
	if err != nil {
		return updated, err
	}

	if tagUpdated {
		updated = true
	}

	if forceNeeded {
		err = ErrForceNeeded
	}

	return
}

func (r *Remote) buildFetchedTags(refs memory.ReferenceStorage) (updated bool, err error) {
	for _, ref := range refs {
		if !ref.Name().IsTag() {
			continue
		}

		_, err := r.s.EncodedObject(plumbing.AnyObject, ref.Hash())
		if err == plumbing.ErrObjectNotFound {
			continue
		}

		if err != nil {
			return false, err
		}

		refUpdated, err := updateReferenceStorerIfNeeded(r.s, ref)
		if err != nil {
			return updated, err
		}

		if refUpdated {
			updated = true
		}
	}

	return
}

// List the references on the remote repository.
// The provided Context must be non-nil. If the context expires before the
// operation is complete, an error is returned. The context only affects to the
// transport operations.
func (r *Remote) ListContext(ctx context.Context, o *ListOptions) (rfs []*plumbing.Reference, err error) {
	return r.list(ctx, o)
}

func (r *Remote) List(o *ListOptions) (rfs []*plumbing.Reference, err error) {
	timeout := o.Timeout
	// Default to the old hardcoded 10s value if a timeout is not explicitly set.
	if timeout == 0 {
		timeout = 10
	}
	if timeout < 0 {
		return nil, fmt.Errorf("invalid timeout: %d", timeout)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	return r.ListContext(ctx, o)
}

func (r *Remote) list(ctx context.Context, o *ListOptions) (rfs []*plumbing.Reference, err error) {
	if r.c == nil || len(r.c.URLs) == 0 {
		return nil, ErrEmptyUrls
	}

	s, err := newUploadPackSession(r.c.URLs[0], o.Auth, o.InsecureSkipTLS, o.CABundle, o.ProxyOptions)
	if err != nil {
		return nil, err
	}

	defer ioutil.CheckClose(s, &err)

	ar, err := s.AdvertisedReferencesContext(ctx)
	if err != nil {
		return nil, err
	}

	allRefs, err := ar.AllReferences()
	if err != nil {
		return nil, err
	}

	refs, err := allRefs.IterReferences()
	if err != nil {
		return nil, err
	}

	var resultRefs []*plumbing.Reference
	if o.PeelingOption == AppendPeeled || o.PeelingOption == IgnorePeeled {
		err = refs.ForEach(func(ref *plumbing.Reference) error {
			resultRefs = append(resultRefs, ref)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if o.PeelingOption == AppendPeeled || o.PeelingOption == OnlyPeeled {
		for k, v := range ar.Peeled {
			resultRefs = append(resultRefs, plumbing.NewReferenceFromStrings(k+"^{}", v.String()))
		}
	}

	return resultRefs, nil
}

func objectsToPush(commands []*packp.Command) []plumbing.Hash {
	objects := make([]plumbing.Hash, 0, len(commands))
	for _, cmd := range commands {
		if cmd.New == plumbing.ZeroHash {
			continue
		}
		objects = append(objects, cmd.New)
	}
	return objects
}

func referencesToHashes(refs storer.ReferenceStorer) ([]plumbing.Hash, error) {
	iter, err := refs.IterReferences()
	if err != nil {
		return nil, err
	}

	var hs []plumbing.Hash
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}

		hs = append(hs, ref.Hash())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return hs, nil
}

func pushHashes(
	ctx context.Context,
	sess transport.ReceivePackSession,
	s storage.Storer,
	req *packp.ReferenceUpdateRequest,
	hs []plumbing.Hash,
	useRefDeltas bool,
	allDelete bool,
) (*packp.ReportStatus, error) {
	rd, wr := io.Pipe()

	config, err := s.Config()
	if err != nil {
		return nil, err
	}

	// Set buffer size to 1 so the error message can be written when
	// ReceivePack fails. Otherwise the goroutine will be blocked writing
	// to the channel.
	done := make(chan error, 1)

	if !allDelete {
		req.Packfile = rd
		go func() {
			e := packfile.NewEncoder(wr, s, useRefDeltas)
			if _, err := e.Encode(hs, config.Pack.Window); err != nil {
				done <- wr.CloseWithError(err)
				return
			}

			done <- wr.Close()
		}()
	} else {
		close(done)
	}

	rs, err := sess.ReceivePack(ctx, req)
	if err != nil {
		// close the pipe to unlock encode write
		_ = rd.Close()
		return nil, err
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return rs, nil
}

func (r *Remote) updateShallow(o *FetchOptions, resp *packp.UploadPackResponse) error {
	if o.Depth == 0 || len(resp.Shallows) == 0 {
		return nil
	}

	shallows, err := r.s.Shallow()
	if err != nil {
		return err
	}

outer:
	for _, s := range resp.Shallows {
		for _, oldS := range shallows {
			if s == oldS {
				continue outer
			}
		}
		shallows = append(shallows, s)
	}

	return r.s.SetShallow(shallows)
}

func (r *Remote) checkRequireRemoteRefs(requires []config.RefSpec, remoteRefs storer.ReferenceStorer) error {
	for _, require := range requires {
		if require.IsWildcard() {
			return fmt.Errorf("wildcards not supported in RequireRemoteRefs, got %s", require.String())
		}

		name := require.Dst("")
		remote, err := remoteRefs.Reference(name)
		if err != nil {
			return fmt.Errorf("remote ref %s required to be %s but is absent", name.String(), require.Src())
		}

		var requireHash string
		if require.IsExactSHA1() {
			requireHash = require.Src()
		} else {
			target, err := storer.ResolveReference(remoteRefs, plumbing.ReferenceName(require.Src()))
			if err != nil {
				return fmt.Errorf("could not resolve ref %s in RequireRemoteRefs", require.Src())
			}
			requireHash = target.Hash().String()
		}

		if remote.Hash().String() != requireHash {
			return fmt.Errorf("remote ref %s required to be %s but is %s", name.String(), requireHash, remote.Hash().String())
		}
	}
	return nil
}
