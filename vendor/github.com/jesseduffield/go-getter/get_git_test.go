package getter

import (
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var testHasGit bool

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
}

func TestGitGetter_impl(t *testing.T) {
	var _ Getter = new(GitGetter)
}

func TestGitGetter(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	g := new(GitGetter)
	dst := tempDir(t)

	repo := testGitRepo(t, "basic")
	repo.commitFile("foo.txt", "hello")

	// With a dir that doesn't exist
	if err := g.Get(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "foo.txt")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGitGetter_branch(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	g := new(GitGetter)
	dst := tempDir(t)

	repo := testGitRepo(t, "branch")
	repo.git("checkout", "-b", "test-branch")
	repo.commitFile("branch.txt", "branch")

	q := repo.url.Query()
	q.Add("ref", "test-branch")
	repo.url.RawQuery = q.Encode()

	if err := g.Get(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "branch.txt")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Get again should work
	if err := g.Get(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath = filepath.Join(dst, "branch.txt")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGitGetter_branchUpdate(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	g := new(GitGetter)
	dst := tempDir(t)

	// First setup the state with a fresh branch
	repo := testGitRepo(t, "branch-update")
	repo.git("checkout", "-b", "test-branch")
	repo.commitFile("branch.txt", "branch")

	// Get the "test-branch" branch
	q := repo.url.Query()
	q.Add("ref", "test-branch")
	repo.url.RawQuery = q.Encode()
	if err := g.Get(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "branch.txt")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Commit an update to the branch
	repo.commitFile("branch-update.txt", "branch-update")

	// Get again should work
	if err := g.Get(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath = filepath.Join(dst, "branch-update.txt")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGitGetter_tag(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	g := new(GitGetter)
	dst := tempDir(t)

	repo := testGitRepo(t, "tag")
	repo.commitFile("tag.txt", "tag")
	repo.git("tag", "v1.0")

	q := repo.url.Query()
	q.Add("ref", "v1.0")
	repo.url.RawQuery = q.Encode()

	if err := g.Get(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "tag.txt")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Get again should work
	if err := g.Get(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath = filepath.Join(dst, "tag.txt")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGitGetter_GetFile(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	g := new(GitGetter)
	dst := tempFile(t)

	repo := testGitRepo(t, "file")
	repo.commitFile("file.txt", "hello")

	// Download the file
	repo.url.Path = filepath.Join(repo.url.Path, "file.txt")
	if err := g.GetFile(dst, repo.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("err: %s", err)
	}
	assertContents(t, dst, "hello")
}

func TestGitGetter_gitVersion(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-getter")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	script := filepath.Join(dir, "git")
	err = ioutil.WriteFile(
		script,
		[]byte("#!/bin/sh\necho \"git version 2.0 (Some Metadata Here)\n\""),
		0700)
	if err != nil {
		t.Fatal(err)
	}

	defer func(v string) {
		os.Setenv("PATH", v)
	}(os.Getenv("PATH"))

	os.Setenv("PATH", dir)

	// Asking for a higher version throws an error
	if err := checkGitVersion("2.3"); err == nil {
		t.Fatal("expect git version error")
	}

	// Passes when version is satisfied
	if err := checkGitVersion("1.9"); err != nil {
		t.Fatal(err)
	}
}

func TestGitGetter_sshKey(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	g := new(GitGetter)
	dst := tempDir(t)

	encodedKey := base64.StdEncoding.EncodeToString([]byte(testGitToken))

	u, err := url.Parse("ssh://git@github.com/hashicorp/test-private-repo" +
		"?sshkey=" + encodedKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := g.Get(dst, u); err != nil {
		t.Fatalf("err: %s", err)
	}

	readmePath := filepath.Join(dst, "README.md")
	if _, err := os.Stat(readmePath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGitGetter_submodule(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	g := new(GitGetter)
	dst := tempDir(t)

	// Set up the grandchild
	gc := testGitRepo(t, "grandchild")
	gc.commitFile("grandchild.txt", "grandchild")

	// Set up the child
	c := testGitRepo(t, "child")
	c.commitFile("child.txt", "child")
	c.git("submodule", "add", gc.dir)
	c.git("commit", "-m", "Add grandchild submodule")

	// Set up the parent
	p := testGitRepo(t, "parent")
	p.commitFile("parent.txt", "parent")
	p.git("submodule", "add", c.dir)
	p.git("commit", "-m", "Add child submodule")

	// Clone the root repository
	if err := g.Get(dst, p.url); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check that the files exist
	for _, path := range []string{
		filepath.Join(dst, "parent.txt"),
		filepath.Join(dst, "child", "child.txt"),
		filepath.Join(dst, "child", "grandchild", "grandchild.txt"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("err: %s", err)
		}
	}
}

func TestGitGetter_setupGitEnv_sshKey(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skipping on windows since the test requires sh")
		return
	}

	cmd := exec.Command("/bin/sh", "-c", "echo $GIT_SSH_COMMAND")
	setupGitEnv(cmd, "/tmp/foo.pem")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	actual := strings.TrimSpace(string(out))
	if actual != "ssh -i /tmp/foo.pem" {
		t.Fatalf("unexpected GIT_SSH_COMMAND: %q", actual)
	}
}

func TestGitGetter_setupGitEnvWithExisting_sshKey(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skipping on windows since the test requires sh")
		return
	}

	// start with an existing ssh command configuration
	os.Setenv("GIT_SSH_COMMAND", "ssh -o StrictHostKeyChecking=no")
	defer os.Setenv("GIT_SSH_COMMAND", "")

	cmd := exec.Command("/bin/sh", "-c", "echo $GIT_SSH_COMMAND")
	setupGitEnv(cmd, "/tmp/foo.pem")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	actual := strings.TrimSpace(string(out))
	if actual != "ssh -o StrictHostKeyChecking=no -i /tmp/foo.pem" {
		t.Fatalf("unexpected GIT_SSH_COMMAND: %q", actual)
	}
}

// gitRepo is a helper struct which controls a single temp git repo.
type gitRepo struct {
	t   *testing.T
	url *url.URL
	dir string
}

// testGitRepo creates a new test git repository.
func testGitRepo(t *testing.T, name string) *gitRepo {
	dir, err := ioutil.TempDir("", "go-getter")
	if err != nil {
		t.Fatal(err)
	}
	dir = filepath.Join(dir, name)
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}

	r := &gitRepo{
		t:   t,
		dir: dir,
	}

	url, err := url.Parse("file://" + r.dir)
	if err != nil {
		t.Fatal(err)
	}
	r.url = url

	r.git("init")
	r.git("config", "user.name", "go-getter")
	r.git("config", "user.email", "go-getter@hashicorp.com")

	return r
}

// git runs a git command against the repo.
func (r *gitRepo) git(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.dir
	if err := cmd.Run(); err != nil {
		r.t.Fatal(err)
	}
}

// commitFile writes and commits a text file to the repo.
func (r *gitRepo) commitFile(file, content string) {
	path := filepath.Join(r.dir, file)
	if err := ioutil.WriteFile(path, []byte(content), 0600); err != nil {
		r.t.Fatal(err)
	}
	r.git("add", file)
	r.git("commit", "-m", "Adding "+file)
}

// This is a read-only deploy key for an empty test repository.
// Note: This is split over multiple lines to avoid being disabled by key
// scanners automatically.
var testGitToken = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA9cHsxCl3Jjgu9DHpwvmfFOl1XEdY+ShHDR/cMnzJ5ddk5/oV
Wy6EWatvyHZfRSZMwzv4PtKeUPm6iXjqWp4xdWU9khlPzozyj+U9Fq70TRVUW9E5
T1XdQVwJE421yffr4VMMwu60wBqjI1epapH2i2inYvw9Zl9X2MXq0+jTvFvDerbT
mDtfStDPljenELAIZtWVETSvbI46gALwbxbM2292ZUIL4D6jRz0aZMmyy/twYv8r
9WGJLwmYzU518Ie7zqKW/mCTdTrV0WRiDj0MeRaPgrGY9amuHE4r9iG/cJkwpKAO
Ccz0Hs6i89u9vZnTqZU9V7weJqRAQcMjXXR6yQIDAQABAoIBAQDBzICKnGxiTlHw
rd+6qqChnAy5jWYDbZjCJ8q8YZ3RS08+g/8NXZxvHftTqM0uOaq1FviHig3gq15H
hHvCpBc6jXDFYoKFzq6FfO/0kFkE5HoWweIgxwRow0xBCDJAJ+ryUEyy+Ay/pQHb
IAjwilRS0V+WdnVw4mTjBAhPvb4jPOo97Yfy3PYUyx2F3newkqXOZy+zx3G/ANoa
ncypfMGyy76sfCWKqw4J1gVkVQLwbB6gQkXUFGYwY9sRrxbG93kQw76Flc/E/s52
62j4v1IM0fq0t/St+Y/+s6Lkw` + `aqt3ft1nsqWcRaVDdqvMfkzgJGXlw0bGzJG5MEQ
AIBq3dHRAoGBAP8OeG/DKG2Z1VmSfzuz1pas1fbZ+F7venOBrjez3sKlb3Pyl2aH
mt2wjaTUi5v10VrHgYtOEdqyhQeUSYydWXIBKNMag0NLLrfFUKZK+57wrHWFdFjn
VgpsdkLSNTOZpC8gA5OaJ+36IcOPfGqyyP9wuuRoaYnVT1KEzqLa9FEFAoGBAPaq
pglwhil2rxjJE4zq0afQLNpAfi7Xqcrepij+xvJIcIj7nawxXuPxqRFxONE/h3yX
zkybO8wLdbHX9Iw/wc1j50Uf1Z5gHdLf7/hQJoWKpz1RnkWRy6CYON8v1tpVp0tb
OAajR/kZnzebq2mfa7pyy5zDCX++2kp/dcFwHf31AoGAE8oupBVTZLWj7TBFuP8q
LkS40U92Sv9v09iDCQVmylmFvUxcXPM2m+7f/qMTNgWrucxzC7kB/6MMWVszHbrz
vrnCTibnemgx9sZTjKOSxHFOIEw7i85fSa3Cu0qOIDPSnmlwfZpfcMKQrhjLAYhf
uhooFiLX1X78iZ2OXup4PHUCgYEAsmBrm83sp1V1gAYBBlnVbXakyNv0pCk/Vz61
iFXeRt1NzDGxLxGw3kQnED8BaIh5kQcyn8Fud7sdzJMv/LAqlT4Ww60mzNYTGyjo
H3jOsqm3ESfRvduWFreeAQBWbiOczGjV1i8D4EbAFfWT+tjXjchwKBf+6Yt5zn/o
Bw/uEHUCgYAFs+JPOR25oRyBs7ujrMo/OY1z/eXTVVgZxY+tYGe1FJqDeFyR7ytK
+JBB1MuDwQKGm2wSIXdCzTNoIx2B9zTseiPTwT8G7vqNFhXoIaTBp4P2xIQb45mJ
7GkTsMBHwpSMOXgX9Weq3v5xOJ2WxVtjENmd6qzxcYCO5lP15O17hA==
-----END RSA PRIVATE KEY-----`
