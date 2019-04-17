package git

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type RemoteSuite struct {
	BaseSuite
}

var _ = Suite(&RemoteSuite{})

func (s *RemoteSuite) TestFetchInvalidEndpoint(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "foo", URLs: []string{"http://\\"}})
	err := r.Fetch(&FetchOptions{RemoteName: "foo"})
	c.Assert(err, ErrorMatches, ".*invalid character.*")
}

func (s *RemoteSuite) TestFetchNonExistentEndpoint(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "foo", URLs: []string{"ssh://non-existent/foo.git"}})
	err := r.Fetch(&FetchOptions{})
	c.Assert(err, NotNil)
}

func (s *RemoteSuite) TestFetchInvalidSchemaEndpoint(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "foo", URLs: []string{"qux://foo"}})
	err := r.Fetch(&FetchOptions{})
	c.Assert(err, ErrorMatches, ".*unsupported scheme.*")
}

func (s *RemoteSuite) TestFetchInvalidFetchOptions(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "foo", URLs: []string{"qux://foo"}})
	invalid := config.RefSpec("^*$ñ")
	err := r.Fetch(&FetchOptions{RefSpecs: []config.RefSpec{invalid}})
	c.Assert(err, Equals, config.ErrRefSpecMalformedSeparator)
}

func (s *RemoteSuite) TestFetchWildcard(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetBasicLocalRepositoryURL()},
	})

	s.testFetch(c, r, &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/remotes/origin/master", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
		plumbing.NewReferenceFromStrings("refs/remotes/origin/branch", "e8d3ffab552895c19b9fcf7aa264d277cde33881"),
		plumbing.NewReferenceFromStrings("refs/tags/v1.0.0", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
	})
}

func (s *RemoteSuite) TestFetchWildcardTags(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetLocalRepositoryURL(fixtures.ByTag("tags").One())},
	})

	s.testFetch(c, r, &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/remotes/origin/master", "f7b877701fbf855b44c0a9e86f3fdce2c298b07f"),
		plumbing.NewReferenceFromStrings("refs/tags/annotated-tag", "b742a2a9fa0afcfa9a6fad080980fbc26b007c69"),
		plumbing.NewReferenceFromStrings("refs/tags/tree-tag", "152175bf7e5580299fa1f0ba41ef6474cc043b70"),
		plumbing.NewReferenceFromStrings("refs/tags/commit-tag", "ad7897c0fb8e7d9a9ba41fa66072cf06095a6cfc"),
		plumbing.NewReferenceFromStrings("refs/tags/blob-tag", "fe6cb94756faa81e5ed9240f9191b833db5f40ae"),
		plumbing.NewReferenceFromStrings("refs/tags/lightweight-tag", "f7b877701fbf855b44c0a9e86f3fdce2c298b07f"),
	})
}

func (s *RemoteSuite) TestFetch(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetLocalRepositoryURL(fixtures.ByTag("tags").One())},
	})

	s.testFetch(c, r, &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/remotes/origin/master"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/remotes/origin/master", "f7b877701fbf855b44c0a9e86f3fdce2c298b07f"),
	})
}

func (s *RemoteSuite) TestFetchNonExistantReference(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetLocalRepositoryURL(fixtures.ByTag("tags").One())},
	})

	err := r.Fetch(&FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/foo:refs/remotes/origin/foo"),
		},
	})

	c.Assert(err, ErrorMatches, "couldn't find remote ref.*")
}

func (s *RemoteSuite) TestFetchContext(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetLocalRepositoryURL(fixtures.ByTag("tags").One())},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := r.FetchContext(ctx, &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/remotes/origin/master"),
		},
	})
	c.Assert(err, NotNil)
}

func (s *RemoteSuite) TestFetchWithAllTags(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetLocalRepositoryURL(fixtures.ByTag("tags").One())},
	})

	s.testFetch(c, r, &FetchOptions{
		Tags: AllTags,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/remotes/origin/master"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/remotes/origin/master", "f7b877701fbf855b44c0a9e86f3fdce2c298b07f"),
		plumbing.NewReferenceFromStrings("refs/tags/annotated-tag", "b742a2a9fa0afcfa9a6fad080980fbc26b007c69"),
		plumbing.NewReferenceFromStrings("refs/tags/tree-tag", "152175bf7e5580299fa1f0ba41ef6474cc043b70"),
		plumbing.NewReferenceFromStrings("refs/tags/commit-tag", "ad7897c0fb8e7d9a9ba41fa66072cf06095a6cfc"),
		plumbing.NewReferenceFromStrings("refs/tags/blob-tag", "fe6cb94756faa81e5ed9240f9191b833db5f40ae"),
		plumbing.NewReferenceFromStrings("refs/tags/lightweight-tag", "f7b877701fbf855b44c0a9e86f3fdce2c298b07f"),
	})
}

func (s *RemoteSuite) TestFetchWithNoTags(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetLocalRepositoryURL(fixtures.ByTag("tags").One())},
	})

	s.testFetch(c, r, &FetchOptions{
		Tags: NoTags,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/remotes/origin/master", "f7b877701fbf855b44c0a9e86f3fdce2c298b07f"),
	})

}

func (s *RemoteSuite) TestFetchWithDepth(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetBasicLocalRepositoryURL()},
	})

	s.testFetch(c, r, &FetchOptions{
		Depth: 1,
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/remotes/origin/master", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
		plumbing.NewReferenceFromStrings("refs/remotes/origin/branch", "e8d3ffab552895c19b9fcf7aa264d277cde33881"),
		plumbing.NewReferenceFromStrings("refs/tags/v1.0.0", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
	})

	c.Assert(r.s.(*memory.Storage).Objects, HasLen, 18)
}

func (s *RemoteSuite) testFetch(c *C, r *Remote, o *FetchOptions, expected []*plumbing.Reference) {
	err := r.Fetch(o)
	c.Assert(err, IsNil)

	var refs int
	l, err := r.s.IterReferences()
	c.Assert(err, IsNil)
	l.ForEach(func(r *plumbing.Reference) error { refs++; return nil })

	c.Assert(refs, Equals, len(expected))

	for _, exp := range expected {
		r, err := r.s.Reference(exp.Name())
		c.Assert(err, IsNil)
		c.Assert(exp.String(), Equals, r.String())
	}
}

func (s *RemoteSuite) TestFetchWithProgress(c *C) {
	url := s.GetBasicLocalRepositoryURL()
	sto := memory.NewStorage()
	buf := bytes.NewBuffer(nil)

	r := newRemote(sto, &config.RemoteConfig{Name: "foo", URLs: []string{url}})

	refspec := config.RefSpec("+refs/heads/*:refs/remotes/origin/*")
	err := r.Fetch(&FetchOptions{
		RefSpecs: []config.RefSpec{refspec},
		Progress: buf,
	})

	c.Assert(err, IsNil)
	c.Assert(sto.Objects, HasLen, 31)

	c.Assert(buf.Len(), Not(Equals), 0)
}

type mockPackfileWriter struct {
	storage.Storer
	PackfileWriterCalled bool
}

func (m *mockPackfileWriter) PackfileWriter() (io.WriteCloser, error) {
	m.PackfileWriterCalled = true
	return m.Storer.(storer.PackfileWriter).PackfileWriter()
}

func (s *RemoteSuite) TestFetchWithPackfileWriter(c *C) {
	dir, err := ioutil.TempDir("", "fetch")
	c.Assert(err, IsNil)

	defer os.RemoveAll(dir) // clean up

	fss, err := filesystem.NewStorage(osfs.New(dir))
	c.Assert(err, IsNil)

	mock := &mockPackfileWriter{Storer: fss}

	url := s.GetBasicLocalRepositoryURL()
	r := newRemote(mock, &config.RemoteConfig{Name: "foo", URLs: []string{url}})

	refspec := config.RefSpec("+refs/heads/*:refs/remotes/origin/*")
	err = r.Fetch(&FetchOptions{
		RefSpecs: []config.RefSpec{refspec},
	})

	c.Assert(err, IsNil)

	var count int
	iter, err := mock.IterEncodedObjects(plumbing.AnyObject)
	c.Assert(err, IsNil)

	iter.ForEach(func(plumbing.EncodedObject) error {
		count++
		return nil
	})

	c.Assert(count, Equals, 31)
	c.Assert(mock.PackfileWriterCalled, Equals, true)
}

func (s *RemoteSuite) TestFetchNoErrAlreadyUpToDate(c *C) {
	url := s.GetBasicLocalRepositoryURL()
	s.doTestFetchNoErrAlreadyUpToDate(c, url)
}

func (s *RemoteSuite) TestFetchNoErrAlreadyUpToDateButStillUpdateLocalRemoteRefs(c *C) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{
		URLs: []string{s.GetBasicLocalRepositoryURL()},
	})

	o := &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
	}

	err := r.Fetch(o)
	c.Assert(err, IsNil)

	// Simulate an out of date remote ref even though we have the new commit locally
	r.s.SetReference(plumbing.NewReferenceFromStrings(
		"refs/remotes/origin/master", "918c48b83bd081e863dbe1b80f8998f058cd8294",
	))

	err = r.Fetch(o)
	c.Assert(err, IsNil)

	exp := plumbing.NewReferenceFromStrings(
		"refs/remotes/origin/master", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
	)

	ref, err := r.s.Reference("refs/remotes/origin/master")
	c.Assert(err, IsNil)
	c.Assert(exp.String(), Equals, ref.String())
}

func (s *RemoteSuite) TestFetchNoErrAlreadyUpToDateWithNonCommitObjects(c *C) {
	fixture := fixtures.ByTag("tags").One()
	url := s.GetLocalRepositoryURL(fixture)
	s.doTestFetchNoErrAlreadyUpToDate(c, url)
}

func (s *RemoteSuite) doTestFetchNoErrAlreadyUpToDate(c *C, url string) {
	r := newRemote(memory.NewStorage(), &config.RemoteConfig{URLs: []string{url}})

	o := &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/*:refs/remotes/origin/*"),
		},
	}

	err := r.Fetch(o)
	c.Assert(err, IsNil)
	err = r.Fetch(o)
	c.Assert(err, Equals, NoErrAlreadyUpToDate)
}

func (s *RemoteSuite) testFetchFastForward(c *C, sto storage.Storer) {
	r := newRemote(sto, &config.RemoteConfig{
		URLs: []string{s.GetBasicLocalRepositoryURL()},
	})

	s.testFetch(c, r, &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/heads/master"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/heads/master", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
	})

	// First make sure that we error correctly when a force is required.
	err := r.Fetch(&FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/heads/branch:refs/heads/master"),
		},
	})
	c.Assert(err, Equals, ErrForceNeeded)

	// And that forcing it fixes the problem.
	err = r.Fetch(&FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/branch:refs/heads/master"),
		},
	})
	c.Assert(err, IsNil)

	// Now test that a fast-forward, non-force fetch works.
	r.s.SetReference(plumbing.NewReferenceFromStrings(
		"refs/heads/master", "918c48b83bd081e863dbe1b80f8998f058cd8294",
	))
	s.testFetch(c, r, &FetchOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/heads/master:refs/heads/master"),
		},
	}, []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/heads/master", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
	})
}

func (s *RemoteSuite) TestFetchFastForwardMem(c *C) {
	s.testFetchFastForward(c, memory.NewStorage())
}

func (s *RemoteSuite) TestFetchFastForwardFS(c *C) {
	dir, err := ioutil.TempDir("", "fetch")
	c.Assert(err, IsNil)

	defer os.RemoveAll(dir) // clean up

	fss, err := filesystem.NewStorage(osfs.New(dir))
	c.Assert(err, IsNil)

	// This exercises `storage.filesystem.Storage.CheckAndSetReference()`.
	s.testFetchFastForward(c, fss)
}

func (s *RemoteSuite) TestString(c *C) {
	r := newRemote(nil, &config.RemoteConfig{
		Name: "foo",
		URLs: []string{"https://github.com/git-fixtures/basic.git"},
	})

	c.Assert(r.String(), Equals, ""+
		"foo\thttps://github.com/git-fixtures/basic.git (fetch)\n"+
		"foo\thttps://github.com/git-fixtures/basic.git (push)",
	)
}

func (s *RemoteSuite) TestPushToEmptyRepository(c *C) {
	url := c.MkDir()
	server, err := PlainInit(url, true)
	c.Assert(err, IsNil)

	srcFs := fixtures.Basic().One().DotGit()
	sto, err := filesystem.NewStorage(srcFs)
	c.Assert(err, IsNil)

	r := newRemote(sto, &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{url},
	})

	rs := config.RefSpec("refs/heads/*:refs/heads/*")
	err = r.Push(&PushOptions{
		RefSpecs: []config.RefSpec{rs},
	})
	c.Assert(err, IsNil)

	iter, err := r.s.IterReferences()
	c.Assert(err, IsNil)

	expected := make(map[string]string)
	iter.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsBranch() {
			return nil
		}

		expected[ref.Name().String()] = ref.Hash().String()
		return nil
	})
	c.Assert(err, IsNil)

	AssertReferences(c, server, expected)

}

func (s *RemoteSuite) TestPushContext(c *C) {
	url := c.MkDir()
	_, err := PlainInit(url, true)
	c.Assert(err, IsNil)

	fs := fixtures.ByURL("https://github.com/git-fixtures/tags.git").One().DotGit()
	sto, err := filesystem.NewStorage(fs)
	c.Assert(err, IsNil)

	r := newRemote(sto, &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{url},
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = r.PushContext(ctx, &PushOptions{
		RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
	})
	c.Assert(err, NotNil)
}

func (s *RemoteSuite) TestPushTags(c *C) {
	url := c.MkDir()
	server, err := PlainInit(url, true)
	c.Assert(err, IsNil)

	fs := fixtures.ByURL("https://github.com/git-fixtures/tags.git").One().DotGit()
	sto, err := filesystem.NewStorage(fs)
	c.Assert(err, IsNil)

	r := newRemote(sto, &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{url},
	})

	err = r.Push(&PushOptions{
		RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
	})
	c.Assert(err, IsNil)

	AssertReferences(c, server, map[string]string{
		"refs/tags/lightweight-tag": "f7b877701fbf855b44c0a9e86f3fdce2c298b07f",
		"refs/tags/annotated-tag":   "b742a2a9fa0afcfa9a6fad080980fbc26b007c69",
		"refs/tags/commit-tag":      "ad7897c0fb8e7d9a9ba41fa66072cf06095a6cfc",
		"refs/tags/blob-tag":        "fe6cb94756faa81e5ed9240f9191b833db5f40ae",
		"refs/tags/tree-tag":        "152175bf7e5580299fa1f0ba41ef6474cc043b70",
	})
}

func (s *RemoteSuite) TestPushNoErrAlreadyUpToDate(c *C) {
	fs := fixtures.Basic().One().DotGit()
	sto, err := filesystem.NewStorage(fs)
	c.Assert(err, IsNil)

	r := newRemote(sto, &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{fs.Root()},
	})

	err = r.Push(&PushOptions{
		RefSpecs: []config.RefSpec{"refs/heads/*:refs/heads/*"},
	})
	c.Assert(err, Equals, NoErrAlreadyUpToDate)
}

func (s *RemoteSuite) TestPushDeleteReference(c *C) {
	fs := fixtures.Basic().One().DotGit()
	sto, err := filesystem.NewStorage(fs)
	c.Assert(err, IsNil)

	r, err := PlainClone(c.MkDir(), true, &CloneOptions{
		URL: fs.Root(),
	})
	c.Assert(err, IsNil)

	remote, err := r.Remote(DefaultRemoteName)
	c.Assert(err, IsNil)

	err = remote.Push(&PushOptions{
		RefSpecs: []config.RefSpec{":refs/heads/branch"},
	})
	c.Assert(err, IsNil)

	_, err = sto.Reference(plumbing.ReferenceName("refs/heads/branch"))
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)

	_, err = r.Storer.Reference(plumbing.ReferenceName("refs/heads/branch"))
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)
}

func (s *RemoteSuite) TestPushRejectNonFastForward(c *C) {
	fs := fixtures.Basic().One().DotGit()
	server, err := filesystem.NewStorage(fs)
	c.Assert(err, IsNil)

	r, err := PlainClone(c.MkDir(), true, &CloneOptions{
		URL: fs.Root(),
	})
	c.Assert(err, IsNil)

	remote, err := r.Remote(DefaultRemoteName)
	c.Assert(err, IsNil)

	branch := plumbing.ReferenceName("refs/heads/branch")
	oldRef, err := server.Reference(branch)
	c.Assert(err, IsNil)
	c.Assert(oldRef, NotNil)

	err = remote.Push(&PushOptions{RefSpecs: []config.RefSpec{
		"refs/heads/master:refs/heads/branch",
	}})
	c.Assert(err, ErrorMatches, "non-fast-forward update: refs/heads/branch")

	newRef, err := server.Reference(branch)
	c.Assert(err, IsNil)
	c.Assert(newRef, DeepEquals, oldRef)
}

func (s *RemoteSuite) TestPushForce(c *C) {
	f := fixtures.Basic().One()
	sto, err := filesystem.NewStorage(f.DotGit())
	c.Assert(err, IsNil)

	dstFs := f.DotGit()
	dstSto, err := filesystem.NewStorage(dstFs)
	c.Assert(err, IsNil)

	url := dstFs.Root()
	r := newRemote(sto, &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{url},
	})

	oldRef, err := dstSto.Reference(plumbing.ReferenceName("refs/heads/branch"))
	c.Assert(err, IsNil)
	c.Assert(oldRef, NotNil)

	err = r.Push(&PushOptions{RefSpecs: []config.RefSpec{
		config.RefSpec("+refs/heads/master:refs/heads/branch"),
	}})
	c.Assert(err, IsNil)

	newRef, err := dstSto.Reference(plumbing.ReferenceName("refs/heads/branch"))
	c.Assert(err, IsNil)
	c.Assert(newRef, Not(DeepEquals), oldRef)
}

func (s *RemoteSuite) TestPushNewReference(c *C) {
	fs := fixtures.Basic().One().DotGit()
	url := c.MkDir()
	server, err := PlainClone(url, true, &CloneOptions{
		URL: fs.Root(),
	})
	c.Assert(err, IsNil)

	r, err := PlainClone(c.MkDir(), true, &CloneOptions{
		URL: url,
	})
	c.Assert(err, IsNil)

	remote, err := r.Remote(DefaultRemoteName)
	c.Assert(err, IsNil)

	ref, err := r.Reference(plumbing.ReferenceName("refs/heads/master"), true)
	c.Assert(err, IsNil)

	err = remote.Push(&PushOptions{RefSpecs: []config.RefSpec{
		"refs/heads/master:refs/heads/branch2",
	}})
	c.Assert(err, IsNil)

	AssertReferences(c, server, map[string]string{
		"refs/heads/branch2": ref.Hash().String(),
	})

	AssertReferences(c, r, map[string]string{
		"refs/remotes/origin/branch2": ref.Hash().String(),
	})
}

func (s *RemoteSuite) TestPushNewReferenceAndDeleteInBatch(c *C) {
	fs := fixtures.Basic().One().DotGit()
	url := c.MkDir()
	server, err := PlainClone(url, true, &CloneOptions{
		URL: fs.Root(),
	})
	c.Assert(err, IsNil)

	r, err := PlainClone(c.MkDir(), true, &CloneOptions{
		URL: url,
	})
	c.Assert(err, IsNil)

	remote, err := r.Remote(DefaultRemoteName)
	c.Assert(err, IsNil)

	ref, err := r.Reference(plumbing.ReferenceName("refs/heads/master"), true)
	c.Assert(err, IsNil)

	err = remote.Push(&PushOptions{RefSpecs: []config.RefSpec{
		"refs/heads/master:refs/heads/branch2",
		":refs/heads/branch",
	}})
	c.Assert(err, IsNil)

	AssertReferences(c, server, map[string]string{
		"refs/heads/branch2": ref.Hash().String(),
	})

	AssertReferences(c, r, map[string]string{
		"refs/remotes/origin/branch2": ref.Hash().String(),
	})

	_, err = server.Storer.Reference(plumbing.ReferenceName("refs/heads/branch"))
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)
}

func (s *RemoteSuite) TestPushInvalidEndpoint(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "foo", URLs: []string{"http://\\"}})
	err := r.Push(&PushOptions{RemoteName: "foo"})
	c.Assert(err, ErrorMatches, ".*invalid character.*")
}

func (s *RemoteSuite) TestPushNonExistentEndpoint(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "foo", URLs: []string{"ssh://non-existent/foo.git"}})
	err := r.Push(&PushOptions{})
	c.Assert(err, NotNil)
}

func (s *RemoteSuite) TestPushInvalidSchemaEndpoint(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "origin", URLs: []string{"qux://foo"}})
	err := r.Push(&PushOptions{})
	c.Assert(err, ErrorMatches, ".*unsupported scheme.*")
}

func (s *RemoteSuite) TestPushInvalidFetchOptions(c *C) {
	r := newRemote(nil, &config.RemoteConfig{Name: "foo", URLs: []string{"qux://foo"}})
	invalid := config.RefSpec("^*$ñ")
	err := r.Push(&PushOptions{RefSpecs: []config.RefSpec{invalid}})
	c.Assert(err, Equals, config.ErrRefSpecMalformedSeparator)
}

func (s *RemoteSuite) TestPushInvalidRefSpec(c *C) {
	r := newRemote(nil, &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{"some-url"},
	})

	rs := config.RefSpec("^*$**")
	err := r.Push(&PushOptions{
		RefSpecs: []config.RefSpec{rs},
	})
	c.Assert(err, Equals, config.ErrRefSpecMalformedSeparator)
}

func (s *RemoteSuite) TestPushWrongRemoteName(c *C) {
	r := newRemote(nil, &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{"some-url"},
	})

	err := r.Push(&PushOptions{
		RemoteName: "other-remote",
	})
	c.Assert(err, ErrorMatches, ".*remote names don't match.*")
}

func (s *RemoteSuite) TestGetHaves(c *C) {
	f := fixtures.Basic().One()
	sto, err := filesystem.NewStorage(f.DotGit())
	c.Assert(err, IsNil)

	var localRefs = []*plumbing.Reference{
		plumbing.NewReferenceFromStrings(
			"foo",
			"f7b877701fbf855b44c0a9e86f3fdce2c298b07f",
		),
		plumbing.NewReferenceFromStrings(
			"bar",
			"fe6cb94756faa81e5ed9240f9191b833db5f40ae",
		),
		plumbing.NewReferenceFromStrings(
			"qux",
			"f7b877701fbf855b44c0a9e86f3fdce2c298b07f",
		),
	}

	l, err := getHaves(localRefs, memory.NewStorage(), sto)
	c.Assert(err, IsNil)
	c.Assert(l, HasLen, 2)
}

func (s *RemoteSuite) TestList(c *C) {
	repo := fixtures.Basic().One()
	remote := newRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{repo.URL},
	})

	refs, err := remote.List(&ListOptions{})
	c.Assert(err, IsNil)

	expected := []*plumbing.Reference{
		plumbing.NewSymbolicReference("HEAD", "refs/heads/master"),
		plumbing.NewReferenceFromStrings("refs/heads/master", "6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
		plumbing.NewReferenceFromStrings("refs/heads/branch", "e8d3ffab552895c19b9fcf7aa264d277cde33881"),
		plumbing.NewReferenceFromStrings("refs/pull/1/head", "b8e471f58bcbca63b07bda20e428190409c2db47"),
		plumbing.NewReferenceFromStrings("refs/pull/2/head", "9632f02833b2f9613afb5e75682132b0b22e4a31"),
		plumbing.NewReferenceFromStrings("refs/pull/2/merge", "c37f58a130ca555e42ff96a071cb9ccb3f437504"),
	}
	c.Assert(len(refs), Equals, len(expected))
	for _, e := range expected {
		found := false
		for _, r := range refs {
			if r.Name() == e.Name() {
				found = true
				c.Assert(r, DeepEquals, e)
			}
		}
		c.Assert(found, Equals, true)
	}
}

func (s *RemoteSuite) TestUpdateShallows(c *C) {
	hashes := []plumbing.Hash{
		plumbing.NewHash("0000000000000000000000000000000000000001"),
		plumbing.NewHash("0000000000000000000000000000000000000002"),
		plumbing.NewHash("0000000000000000000000000000000000000003"),
		plumbing.NewHash("0000000000000000000000000000000000000004"),
		plumbing.NewHash("0000000000000000000000000000000000000005"),
		plumbing.NewHash("0000000000000000000000000000000000000006"),
	}

	tests := []struct {
		hashes []plumbing.Hash
		result []plumbing.Hash
	}{
		// add to empty shallows
		{hashes[0:2], hashes[0:2]},
		// add new hashes
		{hashes[2:4], hashes[0:4]},
		// add some hashes already in shallow list
		{hashes[2:6], hashes[0:6]},
		// add all hashes
		{hashes[0:6], hashes[0:6]},
		// add empty list
		{nil, hashes[0:6]},
	}

	remote := newRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: DefaultRemoteName,
	})

	shallows, err := remote.s.Shallow()
	c.Assert(err, IsNil)
	c.Assert(len(shallows), Equals, 0)

	resp := new(packp.UploadPackResponse)
	o := &FetchOptions{
		Depth: 1,
	}

	for _, t := range tests {
		resp.Shallows = t.hashes
		err = remote.updateShallow(o, resp)
		c.Assert(err, IsNil)

		shallow, err := remote.s.Shallow()
		c.Assert(err, IsNil)
		c.Assert(len(shallow), Equals, len(t.result))
		c.Assert(shallow, DeepEquals, t.result)
	}
}
