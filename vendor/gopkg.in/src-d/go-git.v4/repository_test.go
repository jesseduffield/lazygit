package git

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-billy.v4/util"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type RepositorySuite struct {
	BaseSuite
}

var _ = Suite(&RepositorySuite{})

func (s *RepositorySuite) TestInit(c *C) {
	r, err := Init(memory.NewStorage(), memfs.New())
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Core.IsBare, Equals, false)
}

func (s *RepositorySuite) TestInitNonStandardDotGit(c *C) {
	dir, err := ioutil.TempDir("", "init-non-standard")
	c.Assert(err, IsNil)
	c.Assert(os.RemoveAll(dir), IsNil)

	fs := osfs.New(dir)
	dot, _ := fs.Chroot("storage")
	storage, err := filesystem.NewStorage(dot)
	c.Assert(err, IsNil)

	wt, _ := fs.Chroot("worktree")
	r, err := Init(storage, wt)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	f, err := fs.Open(fs.Join("worktree", ".git"))
	c.Assert(err, IsNil)

	all, err := ioutil.ReadAll(f)
	c.Assert(err, IsNil)
	c.Assert(string(all), Equals, fmt.Sprintf("gitdir: %s\n", filepath.Join("..", "storage")))

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Core.Worktree, Equals, filepath.Join("..", "worktree"))
}

func (s *RepositorySuite) TestInitStandardDotGit(c *C) {
	dir, err := ioutil.TempDir("", "init-standard")
	c.Assert(err, IsNil)
	c.Assert(os.RemoveAll(dir), IsNil)

	fs := osfs.New(dir)
	dot, _ := fs.Chroot(".git")
	storage, err := filesystem.NewStorage(dot)
	c.Assert(err, IsNil)

	r, err := Init(storage, fs)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	l, err := fs.ReadDir(".git")
	c.Assert(err, IsNil)
	c.Assert(len(l) > 0, Equals, true)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Core.Worktree, Equals, "")
}

func (s *RepositorySuite) TestInitBare(c *C) {
	r, err := Init(memory.NewStorage(), nil)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Core.IsBare, Equals, true)

}

func (s *RepositorySuite) TestInitAlreadyExists(c *C) {
	st := memory.NewStorage()

	r, err := Init(st, nil)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = Init(st, nil)
	c.Assert(err, Equals, ErrRepositoryAlreadyExists)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestOpen(c *C) {
	st := memory.NewStorage()

	r, err := Init(st, memfs.New())
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = Open(st, memfs.New())
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)
}

func (s *RepositorySuite) TestOpenBare(c *C) {
	st := memory.NewStorage()

	r, err := Init(st, nil)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = Open(st, nil)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)
}

func (s *RepositorySuite) TestOpenMissingWorktree(c *C) {
	st := memory.NewStorage()

	r, err := Init(st, memfs.New())
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = Open(st, nil)
	c.Assert(err, Equals, ErrWorktreeNotProvided)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestOpenNotExists(c *C) {
	r, err := Open(memory.NewStorage(), nil)
	c.Assert(err, Equals, ErrRepositoryNotExists)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestClone(c *C) {
	r, err := Clone(memory.NewStorage(), nil, &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	remotes, err := r.Remotes()
	c.Assert(err, IsNil)
	c.Assert(remotes, HasLen, 1)
}

func (s *RepositorySuite) TestCloneContext(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := CloneContext(ctx, memory.NewStorage(), nil, &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, NotNil)
}

func (s *RepositorySuite) TestCloneWithTags(c *C) {
	url := s.GetLocalRepositoryURL(
		fixtures.ByURL("https://github.com/git-fixtures/tags.git").One(),
	)

	r, err := Clone(memory.NewStorage(), nil, &CloneOptions{URL: url, Tags: NoTags})
	c.Assert(err, IsNil)

	remotes, err := r.Remotes()
	c.Assert(err, IsNil)
	c.Assert(remotes, HasLen, 1)

	i, err := r.References()
	c.Assert(err, IsNil)

	var count int
	i.ForEach(func(r *plumbing.Reference) error { count++; return nil })

	c.Assert(count, Equals, 3)
}

func (s *RepositorySuite) TestCreateRemoteAndRemote(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: "foo",
		URLs: []string{"http://foo/foo.git"},
	})

	c.Assert(err, IsNil)
	c.Assert(remote.Config().Name, Equals, "foo")

	alt, err := r.Remote("foo")
	c.Assert(err, IsNil)
	c.Assert(alt, Not(Equals), remote)
	c.Assert(alt.Config().Name, Equals, "foo")
}

func (s *RepositorySuite) TestCreateRemoteInvalid(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	remote, err := r.CreateRemote(&config.RemoteConfig{})

	c.Assert(err, Equals, config.ErrRemoteConfigEmptyName)
	c.Assert(remote, IsNil)
}

func (s *RepositorySuite) TestDeleteRemote(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	_, err := r.CreateRemote(&config.RemoteConfig{
		Name: "foo",
		URLs: []string{"http://foo/foo.git"},
	})

	c.Assert(err, IsNil)

	err = r.DeleteRemote("foo")
	c.Assert(err, IsNil)

	alt, err := r.Remote("foo")
	c.Assert(err, Equals, ErrRemoteNotFound)
	c.Assert(alt, IsNil)
}

func (s *RepositorySuite) TestCreateBranchAndBranch(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	testBranch := &config.Branch{
		Name:   "foo",
		Remote: "origin",
		Merge:  "refs/heads/foo",
	}
	err := r.CreateBranch(testBranch)

	c.Assert(err, IsNil)
	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(len(cfg.Branches), Equals, 1)
	branch := cfg.Branches["foo"]
	c.Assert(branch.Name, Equals, testBranch.Name)
	c.Assert(branch.Remote, Equals, testBranch.Remote)
	c.Assert(branch.Merge, Equals, testBranch.Merge)

	branch, err = r.Branch("foo")
	c.Assert(err, IsNil)
	c.Assert(branch.Name, Equals, testBranch.Name)
	c.Assert(branch.Remote, Equals, testBranch.Remote)
	c.Assert(branch.Merge, Equals, testBranch.Merge)
}

func (s *RepositorySuite) TestCreateBranchUnmarshal(c *C) {
	r, _ := Init(memory.NewStorage(), nil)

	expected := []byte(`[core]
	bare = true
[remote "foo"]
	url = http://foo/foo.git
	fetch = +refs/heads/*:refs/remotes/foo/*
[branch "foo"]
	remote = origin
	merge = refs/heads/foo
[branch "master"]
	remote = origin
	merge = refs/heads/master
`)

	_, err := r.CreateRemote(&config.RemoteConfig{
		Name: "foo",
		URLs: []string{"http://foo/foo.git"},
	})
	c.Assert(err, IsNil)
	testBranch1 := &config.Branch{
		Name:   "master",
		Remote: "origin",
		Merge:  "refs/heads/master",
	}
	testBranch2 := &config.Branch{
		Name:   "foo",
		Remote: "origin",
		Merge:  "refs/heads/foo",
	}
	err = r.CreateBranch(testBranch1)
	err = r.CreateBranch(testBranch2)

	c.Assert(err, IsNil)
	cfg, err := r.Config()
	c.Assert(err, IsNil)
	marshaled, err := cfg.Marshal()
	c.Assert(string(expected), Equals, string(marshaled))
}

func (s *RepositorySuite) TestBranchInvalid(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	branch, err := r.Branch("foo")

	c.Assert(err, NotNil)
	c.Assert(branch, IsNil)
}

func (s *RepositorySuite) TestCreateBranchInvalid(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.CreateBranch(&config.Branch{})

	c.Assert(err, NotNil)

	testBranch := &config.Branch{
		Name:   "foo",
		Remote: "origin",
		Merge:  "refs/heads/foo",
	}
	err = r.CreateBranch(testBranch)
	c.Assert(err, IsNil)
	err = r.CreateBranch(testBranch)
	c.Assert(err, NotNil)
}

func (s *RepositorySuite) TestDeleteBranch(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	testBranch := &config.Branch{
		Name:   "foo",
		Remote: "origin",
		Merge:  "refs/heads/foo",
	}
	err := r.CreateBranch(testBranch)

	c.Assert(err, IsNil)

	err = r.DeleteBranch("foo")
	c.Assert(err, IsNil)

	b, err := r.Branch("foo")
	c.Assert(err, Equals, ErrBranchNotFound)
	c.Assert(b, IsNil)

	err = r.DeleteBranch("foo")
	c.Assert(err, Equals, ErrBranchNotFound)
}

func (s *RepositorySuite) TestPlainInit(c *C) {
	dir, err := ioutil.TempDir("", "plain-init")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, true)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Core.IsBare, Equals, true)
}

func (s *RepositorySuite) TestPlainInitAlreadyExists(c *C) {
	dir, err := ioutil.TempDir("", "plain-init")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, true)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = PlainInit(dir, true)
	c.Assert(err, Equals, ErrRepositoryAlreadyExists)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestPlainOpen(c *C) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, false)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = PlainOpen(dir)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)
}

func (s *RepositorySuite) TestPlainOpenBare(c *C) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, true)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = PlainOpen(dir)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)
}

func (s *RepositorySuite) TestPlainOpenNotBare(c *C) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, false)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	r, err = PlainOpen(filepath.Join(dir, ".git"))
	c.Assert(err, Equals, ErrWorktreeNotProvided)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) testPlainOpenGitFile(c *C, f func(string, string) string) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, true)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	altDir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(altDir)

	err = ioutil.WriteFile(filepath.Join(altDir, ".git"), []byte(f(dir, altDir)), 0644)
	c.Assert(err, IsNil)

	r, err = PlainOpen(altDir)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)
}

func (s *RepositorySuite) TestPlainOpenBareAbsoluteGitDirFile(c *C) {
	s.testPlainOpenGitFile(c, func(dir, altDir string) string {
		return fmt.Sprintf("gitdir: %s\n", dir)
	})
}

func (s *RepositorySuite) TestPlainOpenBareAbsoluteGitDirFileNoEOL(c *C) {
	s.testPlainOpenGitFile(c, func(dir, altDir string) string {
		return fmt.Sprintf("gitdir: %s", dir)
	})
}

func (s *RepositorySuite) TestPlainOpenBareRelativeGitDirFile(c *C) {
	s.testPlainOpenGitFile(c, func(dir, altDir string) string {
		dir, err := filepath.Rel(altDir, dir)
		c.Assert(err, IsNil)
		return fmt.Sprintf("gitdir: %s\n", dir)
	})
}

func (s *RepositorySuite) TestPlainOpenBareRelativeGitDirFileNoEOL(c *C) {
	s.testPlainOpenGitFile(c, func(dir, altDir string) string {
		dir, err := filepath.Rel(altDir, dir)
		c.Assert(err, IsNil)
		return fmt.Sprintf("gitdir: %s\n", dir)
	})
}

func (s *RepositorySuite) TestPlainOpenBareRelativeGitDirFileTrailingGarbage(c *C) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, true)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	altDir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	err = ioutil.WriteFile(filepath.Join(altDir, ".git"), []byte(fmt.Sprintf("gitdir: %s\nTRAILING", altDir)), 0644)
	c.Assert(err, IsNil)

	r, err = PlainOpen(altDir)
	c.Assert(err, Equals, ErrRepositoryNotExists)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestPlainOpenBareRelativeGitDirFileBadPrefix(c *C) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	r, err := PlainInit(dir, true)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	altDir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	err = ioutil.WriteFile(filepath.Join(altDir, ".git"), []byte(fmt.Sprintf("xgitdir: %s\n", dir)), 0644)
	c.Assert(err, IsNil)

	r, err = PlainOpen(altDir)
	c.Assert(err, ErrorMatches, ".*gitdir.*")
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestPlainOpenNotExists(c *C) {
	r, err := PlainOpen("/not-exists/")
	c.Assert(err, Equals, ErrRepositoryNotExists)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestPlainOpenDetectDotGit(c *C) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	subdir := filepath.Join(dir, "a", "b")
	err = os.MkdirAll(subdir, 0755)
	c.Assert(err, IsNil)

	r, err := PlainInit(dir, false)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	opt := &PlainOpenOptions{DetectDotGit: true}
	r, err = PlainOpenWithOptions(subdir, opt)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)
}

func (s *RepositorySuite) TestPlainOpenNotExistsDetectDotGit(c *C) {
	dir, err := ioutil.TempDir("", "plain-open")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	opt := &PlainOpenOptions{DetectDotGit: true}
	r, err := PlainOpenWithOptions(dir, opt)
	c.Assert(err, Equals, ErrRepositoryNotExists)
	c.Assert(r, IsNil)
}

func (s *RepositorySuite) TestPlainClone(c *C) {
	r, err := PlainClone(c.MkDir(), false, &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	remotes, err := r.Remotes()
	c.Assert(err, IsNil)
	c.Assert(remotes, HasLen, 1)
	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Branches, HasLen, 1)
	c.Assert(cfg.Branches["master"].Name, Equals, "master")
}

func (s *RepositorySuite) TestPlainCloneContext(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := PlainCloneContext(ctx, c.MkDir(), false, &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, NotNil)
}

func (s *RepositorySuite) TestPlainCloneWithRecurseSubmodules(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	dir, err := ioutil.TempDir("", "plain-clone-submodule")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	path := fixtures.ByTag("submodule").One().Worktree().Root()
	r, err := PlainClone(dir, false, &CloneOptions{
		URL:               path,
		RecurseSubmodules: DefaultSubmoduleRecursionDepth,
	})

	c.Assert(err, IsNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Remotes, HasLen, 1)
	c.Assert(cfg.Branches, HasLen, 1)
	c.Assert(cfg.Submodules, HasLen, 2)
}

func (s *RepositorySuite) TestPlainCloneNoCheckout(c *C) {
	dir, err := ioutil.TempDir("", "plain-clone-no-checkout")
	c.Assert(err, IsNil)
	defer os.RemoveAll(dir)

	path := fixtures.ByTag("submodule").One().Worktree().Root()
	r, err := PlainClone(dir, false, &CloneOptions{
		URL:               path,
		NoCheckout:        true,
		RecurseSubmodules: DefaultSubmoduleRecursionDepth,
	})
	c.Assert(err, IsNil)

	h, err := r.Head()
	c.Assert(err, IsNil)
	c.Assert(h.Hash().String(), Equals, "b685400c1f9316f350965a5993d350bc746b0bf4")

	fi, err := osfs.New(dir).ReadDir("")
	c.Assert(err, IsNil)
	c.Assert(fi, HasLen, 1) // .git
}

func (s *RepositorySuite) TestFetch(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	_, err := r.CreateRemote(&config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{s.GetBasicLocalRepositoryURL()},
	})
	c.Assert(err, IsNil)
	c.Assert(r.Fetch(&FetchOptions{}), IsNil)

	remotes, err := r.Remotes()
	c.Assert(err, IsNil)
	c.Assert(remotes, HasLen, 1)

	_, err = r.Head()
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)

	branch, err := r.Reference("refs/remotes/origin/master", false)
	c.Assert(err, IsNil)
	c.Assert(branch, NotNil)
	c.Assert(branch.Type(), Equals, plumbing.HashReference)
	c.Assert(branch.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")
}

func (s *RepositorySuite) TestFetchContext(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	_, err := r.CreateRemote(&config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{s.GetBasicLocalRepositoryURL()},
	})
	c.Assert(err, IsNil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c.Assert(r.FetchContext(ctx, &FetchOptions{}), NotNil)
}

func (s *RepositorySuite) TestCloneWithProgress(c *C) {
	fs := memfs.New()

	buf := bytes.NewBuffer(nil)
	_, err := Clone(memory.NewStorage(), fs, &CloneOptions{
		URL:      s.GetBasicLocalRepositoryURL(),
		Progress: buf,
	})

	c.Assert(err, IsNil)
	c.Assert(buf.Len(), Not(Equals), 0)
}

func (s *RepositorySuite) TestCloneDeep(c *C) {
	fs := memfs.New()
	r, _ := Init(memory.NewStorage(), fs)

	head, err := r.Head()
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)
	c.Assert(head, IsNil)

	err = r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	remotes, err := r.Remotes()
	c.Assert(err, IsNil)
	c.Assert(remotes, HasLen, 1)

	head, err = r.Reference(plumbing.HEAD, false)
	c.Assert(err, IsNil)
	c.Assert(head, NotNil)
	c.Assert(head.Type(), Equals, plumbing.SymbolicReference)
	c.Assert(head.Target().String(), Equals, "refs/heads/master")

	branch, err := r.Reference(head.Target(), false)
	c.Assert(err, IsNil)
	c.Assert(branch, NotNil)
	c.Assert(branch.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")

	branch, err = r.Reference("refs/remotes/origin/master", false)
	c.Assert(err, IsNil)
	c.Assert(branch, NotNil)
	c.Assert(branch.Type(), Equals, plumbing.HashReference)
	c.Assert(branch.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")

	fi, err := fs.ReadDir("")
	c.Assert(err, IsNil)
	c.Assert(fi, HasLen, 8)
}

func (s *RepositorySuite) TestCloneConfig(c *C) {
	r, _ := Init(memory.NewStorage(), nil)

	head, err := r.Head()
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)
	c.Assert(head, IsNil)

	err = r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)

	c.Assert(cfg.Core.IsBare, Equals, true)
	c.Assert(cfg.Remotes, HasLen, 1)
	c.Assert(cfg.Remotes["origin"].Name, Equals, "origin")
	c.Assert(cfg.Remotes["origin"].URLs, HasLen, 1)
	c.Assert(cfg.Branches, HasLen, 1)
	c.Assert(cfg.Branches["master"].Name, Equals, "master")
}

func (s *RepositorySuite) TestCloneSingleBranchAndNonHEAD(c *C) {
	r, _ := Init(memory.NewStorage(), nil)

	head, err := r.Head()
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)
	c.Assert(head, IsNil)

	err = r.clone(context.Background(), &CloneOptions{
		URL:           s.GetBasicLocalRepositoryURL(),
		ReferenceName: plumbing.ReferenceName("refs/heads/branch"),
		SingleBranch:  true,
	})

	c.Assert(err, IsNil)

	remotes, err := r.Remotes()
	c.Assert(err, IsNil)
	c.Assert(remotes, HasLen, 1)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Branches, HasLen, 1)
	c.Assert(cfg.Branches["branch"].Name, Equals, "branch")
	c.Assert(cfg.Branches["branch"].Remote, Equals, "origin")
	c.Assert(cfg.Branches["branch"].Merge, Equals, plumbing.ReferenceName("refs/heads/branch"))

	head, err = r.Reference(plumbing.HEAD, false)
	c.Assert(err, IsNil)
	c.Assert(head, NotNil)
	c.Assert(head.Type(), Equals, plumbing.SymbolicReference)
	c.Assert(head.Target().String(), Equals, "refs/heads/branch")

	branch, err := r.Reference(head.Target(), false)
	c.Assert(err, IsNil)
	c.Assert(branch, NotNil)
	c.Assert(branch.Hash().String(), Equals, "e8d3ffab552895c19b9fcf7aa264d277cde33881")

	branch, err = r.Reference("refs/remotes/origin/branch", false)
	c.Assert(err, IsNil)
	c.Assert(branch, NotNil)
	c.Assert(branch.Type(), Equals, plumbing.HashReference)
	c.Assert(branch.Hash().String(), Equals, "e8d3ffab552895c19b9fcf7aa264d277cde33881")
}

func (s *RepositorySuite) TestCloneSingleBranch(c *C) {
	r, _ := Init(memory.NewStorage(), nil)

	head, err := r.Head()
	c.Assert(err, Equals, plumbing.ErrReferenceNotFound)
	c.Assert(head, IsNil)

	err = r.clone(context.Background(), &CloneOptions{
		URL:          s.GetBasicLocalRepositoryURL(),
		SingleBranch: true,
	})

	c.Assert(err, IsNil)

	remotes, err := r.Remotes()
	c.Assert(err, IsNil)
	c.Assert(remotes, HasLen, 1)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Branches, HasLen, 1)
	c.Assert(cfg.Branches["master"].Name, Equals, "master")
	c.Assert(cfg.Branches["master"].Remote, Equals, "origin")
	c.Assert(cfg.Branches["master"].Merge, Equals, plumbing.ReferenceName("refs/heads/master"))

	head, err = r.Reference(plumbing.HEAD, false)
	c.Assert(err, IsNil)
	c.Assert(head, NotNil)
	c.Assert(head.Type(), Equals, plumbing.SymbolicReference)
	c.Assert(head.Target().String(), Equals, "refs/heads/master")

	branch, err := r.Reference(head.Target(), false)
	c.Assert(err, IsNil)
	c.Assert(branch, NotNil)
	c.Assert(branch.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")

	branch, err = r.Reference("refs/remotes/origin/master", false)
	c.Assert(err, IsNil)
	c.Assert(branch, NotNil)
	c.Assert(branch.Type(), Equals, plumbing.HashReference)
	c.Assert(branch.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")
}

func (s *RepositorySuite) TestCloneDetachedHEAD(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL:           s.GetBasicLocalRepositoryURL(),
		ReferenceName: plumbing.ReferenceName("refs/tags/v1.0.0"),
	})
	c.Assert(err, IsNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Branches, HasLen, 0)

	head, err := r.Reference(plumbing.HEAD, false)
	c.Assert(err, IsNil)
	c.Assert(head, NotNil)
	c.Assert(head.Type(), Equals, plumbing.HashReference)
	c.Assert(head.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")

	count := 0
	objects, err := r.Objects()
	c.Assert(err, IsNil)
	objects.ForEach(func(object.Object) error { count++; return nil })
	c.Assert(count, Equals, 31)
}

func (s *RepositorySuite) TestCloneDetachedHEADAndShallow(c *C) {
	r, _ := Init(memory.NewStorage(), memfs.New())
	err := r.clone(context.Background(), &CloneOptions{
		URL:           s.GetBasicLocalRepositoryURL(),
		ReferenceName: plumbing.ReferenceName("refs/tags/v1.0.0"),
		Depth:         1,
	})

	c.Assert(err, IsNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Branches, HasLen, 0)

	head, err := r.Reference(plumbing.HEAD, false)
	c.Assert(err, IsNil)
	c.Assert(head, NotNil)
	c.Assert(head.Type(), Equals, plumbing.HashReference)
	c.Assert(head.Hash().String(), Equals, "6ecf0ef2c2dffb796033e5a02219af86ec6584e5")

	count := 0
	objects, err := r.Objects()
	c.Assert(err, IsNil)
	objects.ForEach(func(object.Object) error { count++; return nil })
	c.Assert(count, Equals, 15)
}

func (s *RepositorySuite) TestCloneDetachedHEADAnnotatedTag(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL:           s.GetLocalRepositoryURL(fixtures.ByTag("tags").One()),
		ReferenceName: plumbing.ReferenceName("refs/tags/annotated-tag"),
	})
	c.Assert(err, IsNil)

	cfg, err := r.Config()
	c.Assert(err, IsNil)
	c.Assert(cfg.Branches, HasLen, 0)

	head, err := r.Reference(plumbing.HEAD, false)
	c.Assert(err, IsNil)
	c.Assert(head, NotNil)
	c.Assert(head.Type(), Equals, plumbing.HashReference)
	c.Assert(head.Hash().String(), Equals, "f7b877701fbf855b44c0a9e86f3fdce2c298b07f")

	count := 0
	objects, err := r.Objects()
	c.Assert(err, IsNil)
	objects.ForEach(func(object.Object) error { count++; return nil })
	c.Assert(count, Equals, 7)
}

func (s *RepositorySuite) TestPush(c *C) {
	url := c.MkDir()
	server, err := PlainInit(url, true)
	c.Assert(err, IsNil)

	_, err = s.Repository.CreateRemote(&config.RemoteConfig{
		Name: "test",
		URLs: []string{url},
	})
	c.Assert(err, IsNil)

	err = s.Repository.Push(&PushOptions{
		RemoteName: "test",
	})
	c.Assert(err, IsNil)

	AssertReferences(c, server, map[string]string{
		"refs/heads/master": "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"refs/heads/branch": "e8d3ffab552895c19b9fcf7aa264d277cde33881",
	})

	AssertReferences(c, s.Repository, map[string]string{
		"refs/remotes/test/master": "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"refs/remotes/test/branch": "e8d3ffab552895c19b9fcf7aa264d277cde33881",
	})
}

func (s *RepositorySuite) TestPushContext(c *C) {
	url := c.MkDir()
	_, err := PlainInit(url, true)
	c.Assert(err, IsNil)

	_, err = s.Repository.CreateRemote(&config.RemoteConfig{
		Name: "foo",
		URLs: []string{url},
	})
	c.Assert(err, IsNil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = s.Repository.PushContext(ctx, &PushOptions{
		RemoteName: "foo",
	})
	c.Assert(err, NotNil)
}

// installPreReceiveHook installs a pre-receive hook in the .git
// directory at path which prints message m before exiting
// successfully.
func installPreReceiveHook(c *C, path, m string) {
	hooks := filepath.Join(path, "hooks")
	err := os.MkdirAll(hooks, 0777)
	c.Assert(err, IsNil)

	err = ioutil.WriteFile(filepath.Join(hooks, "pre-receive"), preReceiveHook(m), 0777)
	c.Assert(err, IsNil)
}

func (s *RepositorySuite) TestPushWithProgress(c *C) {
	url := c.MkDir()
	server, err := PlainInit(url, true)
	c.Assert(err, IsNil)

	m := "Receiving..."
	installPreReceiveHook(c, url, m)

	_, err = s.Repository.CreateRemote(&config.RemoteConfig{
		Name: "bar",
		URLs: []string{url},
	})
	c.Assert(err, IsNil)

	var p bytes.Buffer
	err = s.Repository.Push(&PushOptions{
		RemoteName: "bar",
		Progress:   &p,
	})
	c.Assert(err, IsNil)

	AssertReferences(c, server, map[string]string{
		"refs/heads/master": "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"refs/heads/branch": "e8d3ffab552895c19b9fcf7aa264d277cde33881",
	})

	c.Assert((&p).Bytes(), DeepEquals, []byte(m))
}

func (s *RepositorySuite) TestPushDepth(c *C) {
	url := c.MkDir()
	server, err := PlainClone(url, true, &CloneOptions{
		URL: fixtures.Basic().One().DotGit().Root(),
	})

	c.Assert(err, IsNil)

	r, err := Clone(memory.NewStorage(), memfs.New(), &CloneOptions{
		URL:   url,
		Depth: 1,
	})
	c.Assert(err, IsNil)

	err = util.WriteFile(r.wt, "foo", nil, 0755)
	c.Assert(err, IsNil)

	w, err := r.Worktree()
	c.Assert(err, IsNil)

	_, err = w.Add("foo")
	c.Assert(err, IsNil)

	hash, err := w.Commit("foo", &CommitOptions{
		Author:    defaultSignature(),
		Committer: defaultSignature(),
	})
	c.Assert(err, IsNil)

	err = r.Push(&PushOptions{})
	c.Assert(err, IsNil)

	AssertReferences(c, server, map[string]string{
		"refs/heads/master": hash.String(),
	})

	AssertReferences(c, r, map[string]string{
		"refs/remotes/origin/master": hash.String(),
	})
}

func (s *RepositorySuite) TestPushNonExistentRemote(c *C) {
	srcFs := fixtures.Basic().One().DotGit()
	sto, err := filesystem.NewStorage(srcFs)
	c.Assert(err, IsNil)

	r, err := Open(sto, srcFs)
	c.Assert(err, IsNil)

	err = r.Push(&PushOptions{RemoteName: "myremote"})
	c.Assert(err, ErrorMatches, ".*remote not found.*")
}

func (s *RepositorySuite) TestLog(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	cIter, err := r.Log(&LogOptions{
		From: plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47"),
	})

	c.Assert(err, IsNil)

	commitOrder := []plumbing.Hash{
		plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47"),
		plumbing.NewHash("b029517f6300c2da0f4b651b8642506cd6aaf45d"),
	}

	for _, o := range commitOrder {
		commit, err := cIter.Next()
		c.Assert(err, IsNil)
		c.Assert(commit.Hash, Equals, o)
	}
	_, err = cIter.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *RepositorySuite) TestLogHead(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	cIter, err := r.Log(&LogOptions{})

	c.Assert(err, IsNil)

	commitOrder := []plumbing.Hash{
		plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"),
		plumbing.NewHash("918c48b83bd081e863dbe1b80f8998f058cd8294"),
		plumbing.NewHash("af2d6a6954d532f8ffb47615169c8fdf9d383a1a"),
		plumbing.NewHash("1669dce138d9b841a518c64b10914d88f5e488ea"),
		plumbing.NewHash("35e85108805c84807bc66a02d91535e1e24b38b9"),
		plumbing.NewHash("b029517f6300c2da0f4b651b8642506cd6aaf45d"),
		plumbing.NewHash("a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69"),
		plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47"),
	}

	for _, o := range commitOrder {
		commit, err := cIter.Next()
		c.Assert(err, IsNil)
		c.Assert(commit.Hash, Equals, o)
	}
	_, err = cIter.Next()
	c.Assert(err, Equals, io.EOF)
}

func (s *RepositorySuite) TestLogError(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	_, err = r.Log(&LogOptions{
		From: plumbing.NewHash("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	})
	c.Assert(err, NotNil)
}

func (s *RepositorySuite) TestCommit(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	hash := plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47")
	commit, err := r.CommitObject(hash)
	c.Assert(err, IsNil)

	c.Assert(commit.Hash.IsZero(), Equals, false)
	c.Assert(commit.Hash, Equals, commit.ID())
	c.Assert(commit.Hash, Equals, hash)
	c.Assert(commit.Type(), Equals, plumbing.CommitObject)

	tree, err := commit.Tree()
	c.Assert(err, IsNil)
	c.Assert(tree.Hash.IsZero(), Equals, false)

	c.Assert(commit.Author.Email, Equals, "daniel@lordran.local")
}

func (s *RepositorySuite) TestCommits(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	count := 0
	commits, err := r.CommitObjects()
	c.Assert(err, IsNil)
	for {
		commit, err := commits.Next()
		if err != nil {
			break
		}

		count++
		c.Assert(commit.Hash.IsZero(), Equals, false)
		c.Assert(commit.Hash, Equals, commit.ID())
		c.Assert(commit.Type(), Equals, plumbing.CommitObject)
	}

	c.Assert(count, Equals, 9)
}

func (s *RepositorySuite) TestBlob(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})

	c.Assert(err, IsNil)

	blob, err := r.BlobObject(plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47"))
	c.Assert(err, NotNil)
	c.Assert(blob, IsNil)

	blobHash := plumbing.NewHash("9a48f23120e880dfbe41f7c9b7b708e9ee62a492")
	blob, err = r.BlobObject(blobHash)
	c.Assert(err, IsNil)

	c.Assert(blob.Hash.IsZero(), Equals, false)
	c.Assert(blob.Hash, Equals, blob.ID())
	c.Assert(blob.Hash, Equals, blobHash)
	c.Assert(blob.Type(), Equals, plumbing.BlobObject)
}

func (s *RepositorySuite) TestBlobs(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	count := 0
	blobs, err := r.BlobObjects()
	c.Assert(err, IsNil)
	for {
		blob, err := blobs.Next()
		if err != nil {
			break
		}

		count++
		c.Assert(blob.Hash.IsZero(), Equals, false)
		c.Assert(blob.Hash, Equals, blob.ID())
		c.Assert(blob.Type(), Equals, plumbing.BlobObject)
	}

	c.Assert(count, Equals, 10)
}

func (s *RepositorySuite) TestTagObject(c *C) {
	url := s.GetLocalRepositoryURL(
		fixtures.ByURL("https://github.com/git-fixtures/tags.git").One(),
	)

	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: url})
	c.Assert(err, IsNil)

	hash := plumbing.NewHash("ad7897c0fb8e7d9a9ba41fa66072cf06095a6cfc")
	tag, err := r.TagObject(hash)
	c.Assert(err, IsNil)

	c.Assert(tag.Hash.IsZero(), Equals, false)
	c.Assert(tag.Hash, Equals, hash)
	c.Assert(tag.Type(), Equals, plumbing.TagObject)
}

func (s *RepositorySuite) TestTags(c *C) {
	url := s.GetLocalRepositoryURL(
		fixtures.ByURL("https://github.com/git-fixtures/tags.git").One(),
	)

	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: url})
	c.Assert(err, IsNil)

	count := 0
	tags, err := r.Tags()
	c.Assert(err, IsNil)

	tags.ForEach(func(tag *plumbing.Reference) error {
		count++
		c.Assert(tag.Hash().IsZero(), Equals, false)
		c.Assert(tag.Name().IsTag(), Equals, true)
		return nil
	})

	c.Assert(count, Equals, 5)
}

func (s *RepositorySuite) TestBranches(c *C) {
	f := fixtures.ByURL("https://github.com/git-fixtures/root-references.git").One()
	sto, err := filesystem.NewStorage(f.DotGit())
	c.Assert(err, IsNil)
	r, err := Open(sto, f.DotGit())
	c.Assert(err, IsNil)

	count := 0
	branches, err := r.Branches()
	c.Assert(err, IsNil)

	branches.ForEach(func(branch *plumbing.Reference) error {
		count++
		c.Assert(branch.Hash().IsZero(), Equals, false)
		c.Assert(branch.Name().IsBranch(), Equals, true)
		return nil
	})

	c.Assert(count, Equals, 8)
}

func (s *RepositorySuite) TestNotes(c *C) {
	// TODO add fixture with Notes
	url := s.GetLocalRepositoryURL(
		fixtures.ByURL("https://github.com/git-fixtures/tags.git").One(),
	)

	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: url})
	c.Assert(err, IsNil)

	count := 0
	notes, err := r.Notes()
	c.Assert(err, IsNil)

	notes.ForEach(func(note *plumbing.Reference) error {
		count++
		c.Assert(note.Hash().IsZero(), Equals, false)
		c.Assert(note.Name().IsNote(), Equals, true)
		return nil
	})

	c.Assert(count, Equals, 0)
}

func (s *RepositorySuite) TestTree(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{
		URL: s.GetBasicLocalRepositoryURL(),
	})
	c.Assert(err, IsNil)

	invalidHash := plumbing.NewHash("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	tree, err := r.TreeObject(invalidHash)
	c.Assert(tree, IsNil)
	c.Assert(err, NotNil)

	hash := plumbing.NewHash("dbd3641b371024f44d0e469a9c8f5457b0660de1")
	tree, err = r.TreeObject(hash)
	c.Assert(err, IsNil)

	c.Assert(tree.Hash.IsZero(), Equals, false)
	c.Assert(tree.Hash, Equals, tree.ID())
	c.Assert(tree.Hash, Equals, hash)
	c.Assert(tree.Type(), Equals, plumbing.TreeObject)
	c.Assert(len(tree.Entries), Not(Equals), 0)
}

func (s *RepositorySuite) TestTrees(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	count := 0
	trees, err := r.TreeObjects()
	c.Assert(err, IsNil)
	for {
		tree, err := trees.Next()
		if err != nil {
			break
		}

		count++
		c.Assert(tree.Hash.IsZero(), Equals, false)
		c.Assert(tree.Hash, Equals, tree.ID())
		c.Assert(tree.Type(), Equals, plumbing.TreeObject)
		c.Assert(len(tree.Entries), Not(Equals), 0)
	}

	c.Assert(count, Equals, 12)
}

func (s *RepositorySuite) TestTagObjects(c *C) {
	url := s.GetLocalRepositoryURL(
		fixtures.ByURL("https://github.com/git-fixtures/tags.git").One(),
	)

	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: url})
	c.Assert(err, IsNil)

	count := 0
	tags, err := r.TagObjects()
	c.Assert(err, IsNil)

	tags.ForEach(func(tag *object.Tag) error {
		count++

		c.Assert(tag.Hash.IsZero(), Equals, false)
		c.Assert(tag.Type(), Equals, plumbing.TagObject)
		return nil
	})

	refs, _ := r.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		return nil
	})

	c.Assert(count, Equals, 4)
}

func (s *RepositorySuite) TestCommitIterClosePanic(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	commits, err := r.CommitObjects()
	c.Assert(err, IsNil)
	commits.Close()
}

func (s *RepositorySuite) TestRef(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	ref, err := r.Reference(plumbing.HEAD, false)
	c.Assert(err, IsNil)
	c.Assert(ref.Name(), Equals, plumbing.HEAD)

	ref, err = r.Reference(plumbing.HEAD, true)
	c.Assert(err, IsNil)
	c.Assert(ref.Name(), Equals, plumbing.ReferenceName("refs/heads/master"))
}

func (s *RepositorySuite) TestRefs(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)

	iter, err := r.References()
	c.Assert(err, IsNil)
	c.Assert(iter, NotNil)
}

func (s *RepositorySuite) TestObject(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	hash := plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5")
	o, err := r.Object(plumbing.CommitObject, hash)
	c.Assert(err, IsNil)

	c.Assert(o.ID().IsZero(), Equals, false)
	c.Assert(o.Type(), Equals, plumbing.CommitObject)
}

func (s *RepositorySuite) TestObjects(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	count := 0
	objects, err := r.Objects()
	c.Assert(err, IsNil)
	for {
		o, err := objects.Next()
		if err != nil {
			break
		}

		count++
		c.Assert(o.ID().IsZero(), Equals, false)
		c.Assert(o.Type(), Not(Equals), plumbing.AnyObject)
	}

	c.Assert(count, Equals, 31)
}

func (s *RepositorySuite) TestObjectNotFound(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: s.GetBasicLocalRepositoryURL()})
	c.Assert(err, IsNil)

	hash := plumbing.NewHash("0a3fb06ff80156fb153bcdcc58b5e16c2d27625c")
	tag, err := r.Object(plumbing.TagObject, hash)
	c.Assert(err, DeepEquals, plumbing.ErrObjectNotFound)
	c.Assert(tag, IsNil)
}

func (s *RepositorySuite) TestWorktree(c *C) {
	def := memfs.New()
	r, _ := Init(memory.NewStorage(), def)
	w, err := r.Worktree()
	c.Assert(err, IsNil)
	c.Assert(w.Filesystem, Equals, def)
}

func (s *RepositorySuite) TestWorktreeBare(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	w, err := r.Worktree()
	c.Assert(err, Equals, ErrIsBareRepository)
	c.Assert(w, IsNil)
}

func (s *RepositorySuite) TestResolveRevision(c *C) {
	f := fixtures.ByURL("https://github.com/git-fixtures/basic.git").One()
	sto, err := filesystem.NewStorage(f.DotGit())
	c.Assert(err, IsNil)
	r, err := Open(sto, f.DotGit())
	c.Assert(err, IsNil)

	datas := map[string]string{
		"HEAD":                       "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"heads/master":               "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"heads/master~1":             "918c48b83bd081e863dbe1b80f8998f058cd8294",
		"refs/heads/master":          "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"refs/heads/master~2^^~":     "b029517f6300c2da0f4b651b8642506cd6aaf45d",
		"refs/tags/v1.0.0":           "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"refs/remotes/origin/master": "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"refs/remotes/origin/HEAD":   "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"HEAD~2^^~":                  "b029517f6300c2da0f4b651b8642506cd6aaf45d",
		"HEAD~3^2":                   "a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69",
		"HEAD~3^2^0":                 "a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69",
		"HEAD~2^{/binary file}":      "35e85108805c84807bc66a02d91535e1e24b38b9",
		"HEAD~^{/!-some}":            "1669dce138d9b841a518c64b10914d88f5e488ea",
		"master":                     "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"branch":                     "e8d3ffab552895c19b9fcf7aa264d277cde33881",
		"v1.0.0":                     "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
		"branch~1":                   "918c48b83bd081e863dbe1b80f8998f058cd8294",
		"v1.0.0~1":                   "918c48b83bd081e863dbe1b80f8998f058cd8294",
		"master~1":                   "918c48b83bd081e863dbe1b80f8998f058cd8294",
		"918c48b83bd081e863dbe1b80f8998f058cd8294": "918c48b83bd081e863dbe1b80f8998f058cd8294",
	}

	for rev, hash := range datas {
		h, err := r.ResolveRevision(plumbing.Revision(rev))

		c.Assert(err, IsNil)
		c.Assert(h.String(), Equals, hash)
	}
}

func (s *RepositorySuite) TestResolveRevisionWithErrors(c *C) {
	url := s.GetLocalRepositoryURL(
		fixtures.ByURL("https://github.com/git-fixtures/basic.git").One(),
	)

	r, _ := Init(memory.NewStorage(), nil)
	err := r.clone(context.Background(), &CloneOptions{URL: url})
	c.Assert(err, IsNil)

	headRef, err := r.Head()
	c.Assert(err, IsNil)

	ref := plumbing.NewHashReference("refs/heads/918c48b83bd081e863dbe1b80f8998f058cd8294", headRef.Hash())
	err = r.Storer.SetReference(ref)
	c.Assert(err, IsNil)

	datas := map[string]string{
		"efs/heads/master~":                        "reference not found",
		"HEAD^3":                                   `Revision invalid : "3" found must be 0, 1 or 2 after "^"`,
		"HEAD^{/whatever}":                         `No commit message match regexp : "whatever"`,
		"4e1243bd22c66e76c2ba9eddc1f91394e57f9f83": "reference not found",
		"918c48b83bd081e863dbe1b80f8998f058cd8294": `refname "918c48b83bd081e863dbe1b80f8998f058cd8294" is ambiguous`,
	}

	for rev, rerr := range datas {
		_, err := r.ResolveRevision(plumbing.Revision(rev))

		c.Assert(err.Error(), Equals, rerr)
	}
}

func (s *RepositorySuite) testRepackObjects(
	c *C, deleteTime time.Time, expectedPacks int) {
	srcFs := fixtures.ByTag("unpacked").One().DotGit()
	var sto storage.Storer
	var err error
	sto, err = filesystem.NewStorage(srcFs)
	c.Assert(err, IsNil)

	los := sto.(storer.LooseObjectStorer)
	c.Assert(los, NotNil)

	numLooseStart := 0
	err = los.ForEachObjectHash(func(_ plumbing.Hash) error {
		numLooseStart++
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(numLooseStart > 0, Equals, true)

	pos := sto.(storer.PackedObjectStorer)
	c.Assert(los, NotNil)

	packs, err := pos.ObjectPacks()
	c.Assert(err, IsNil)
	numPacksStart := len(packs)
	c.Assert(numPacksStart > 1, Equals, true)

	r, err := Open(sto, srcFs)
	c.Assert(err, IsNil)
	c.Assert(r, NotNil)

	err = r.RepackObjects(&RepackConfig{
		OnlyDeletePacksOlderThan: deleteTime,
	})
	c.Assert(err, IsNil)

	numLooseEnd := 0
	err = los.ForEachObjectHash(func(_ plumbing.Hash) error {
		numLooseEnd++
		return nil
	})
	c.Assert(err, IsNil)
	c.Assert(numLooseEnd, Equals, 0)

	packs, err = pos.ObjectPacks()
	c.Assert(err, IsNil)
	numPacksEnd := len(packs)
	c.Assert(numPacksEnd, Equals, expectedPacks)
}

func (s *RepositorySuite) TestRepackObjects(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	s.testRepackObjects(c, time.Time{}, 1)
}

func (s *RepositorySuite) TestRepackObjectsWithNoDelete(c *C) {
	if testing.Short() {
		c.Skip("skipping test in short mode.")
	}

	s.testRepackObjects(c, time.Unix(0, 1), 3)
}

func ExecuteOnPath(c *C, path string, cmds ...string) error {
	for _, cmd := range cmds {
		err := executeOnPath(path, cmd)
		c.Assert(err, IsNil)
	}

	return nil
}

func executeOnPath(path, cmd string) error {
	args := strings.Split(cmd, " ")
	c := exec.Command(args[0], args[1:]...)
	c.Dir = path
	c.Env = os.Environ()

	buf := bytes.NewBuffer(nil)
	c.Stderr = buf
	c.Stdout = buf

	//defer func() { fmt.Println(buf.String()) }()

	return c.Run()
}

func (s *RepositorySuite) TestBrokenMultipleShallowFetch(c *C) {
	r, _ := Init(memory.NewStorage(), nil)
	_, err := r.CreateRemote(&config.RemoteConfig{
		Name: DefaultRemoteName,
		URLs: []string{s.GetBasicLocalRepositoryURL()},
	})
	c.Assert(err, IsNil)

	c.Assert(r.Fetch(&FetchOptions{
		Depth:    2,
		RefSpecs: []config.RefSpec{config.RefSpec("refs/heads/master:refs/heads/master")},
	}), IsNil)

	shallows, err := r.Storer.Shallow()
	c.Assert(err, IsNil)
	c.Assert(len(shallows), Equals, 1)

	ref, err := r.Reference("refs/heads/master", true)
	c.Assert(err, IsNil)
	cobj, err := r.CommitObject(ref.Hash())
	c.Assert(err, IsNil)
	c.Assert(cobj, NotNil)
	err = object.NewCommitPreorderIter(cobj, nil, nil).ForEach(func(c *object.Commit) error {
		for _, ph := range c.ParentHashes {
			for _, h := range shallows {
				if ph == h {
					return storer.ErrStop
				}
			}
		}

		return nil
	})
	c.Assert(err, IsNil)

	c.Assert(r.Fetch(&FetchOptions{
		Depth:    5,
		RefSpecs: []config.RefSpec{config.RefSpec("refs/heads/*:refs/heads/*")},
	}), IsNil)

	shallows, err = r.Storer.Shallow()
	c.Assert(err, IsNil)
	c.Assert(len(shallows), Equals, 3)

	ref, err = r.Reference("refs/heads/master", true)
	c.Assert(err, IsNil)
	cobj, err = r.CommitObject(ref.Hash())
	c.Assert(err, IsNil)
	c.Assert(cobj, NotNil)
	err = object.NewCommitPreorderIter(cobj, nil, nil).ForEach(func(c *object.Commit) error {
		for _, ph := range c.ParentHashes {
			for _, h := range shallows {
				if ph == h {
					return storer.ErrStop
				}
			}
		}

		return nil
	})
	c.Assert(err, IsNil)
}
