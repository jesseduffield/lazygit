package gitconfig

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGlobal(t *testing.T) {
	RegisterTestingT(t)

	reset := withGlobalGitConfigFile(`
[user]
    name  = deeeet
    email = deeeet@example.com
[github]
    user = ghdeeeet
`)
	defer reset()

	var err error
	username, err := Global("user.name")
	Expect(err).NotTo(HaveOccurred())
	Expect(username).To(Equal("deeeet"))

	email, err := Global("user.email")
	Expect(err).NotTo(HaveOccurred())
	Expect(email).To(Equal("deeeet@example.com"))

	githubuser, err := Global("github.user")
	Expect(err).NotTo(HaveOccurred())
	Expect(githubuser).To(Equal("ghdeeeet"))

	nothing, err := Local("nothing.return")
	Expect(err).To(HaveOccurred())
	Expect(nothing).To(Equal(""))
}

func TestEntire(t *testing.T) {
	RegisterTestingT(t)

	includeFilePath := includeGitConfigFile(`
[user]
    name  = deeeet
    email = deeeet@example.com
	`)

	content := fmt.Sprintf(`
[include]
    path = %s
`, includeFilePath)

	reset := withGlobalGitConfigFile(content)
	defer reset()

	var err error
	username, err := Entire("user.name")
	Expect(err).NotTo(HaveOccurred())
	Expect(username).To(Equal("deeeet"))

	email, err := Entire("user.email")
	Expect(err).NotTo(HaveOccurred())
	Expect(email).To(Equal("deeeet@example.com"))

	nothing, err := Local("nothing.return")
	Expect(err).To(HaveOccurred())
	Expect(nothing).To(Equal(""))
}

func TestLocal(t *testing.T) {
	RegisterTestingT(t)

	reset := withLocalGitConfigFile("remote.origin.url", "git@github.com:tcnksm/go-test-gitconfig.git")
	defer reset()

	var err error
	url, err := Local("remote.origin.url")
	Expect(err).NotTo(HaveOccurred())
	Expect(url).To(Equal("git@github.com:tcnksm/go-test-gitconfig.git"))

	nothing, err := Local("nothing.return")
	Expect(err).To(HaveOccurred())
	Expect(err.Error()).To(Equal("the key `nothing.return` is not found"))
	Expect(nothing).To(Equal(""))
}

func TestUsername(t *testing.T) {
	RegisterTestingT(t)

	reset := withGlobalGitConfigFile(`
[user]
    name  = taichi
    email = taichi@example.com
`)
	defer reset()

	var err error
	username, err := Username()
	Expect(err).NotTo(HaveOccurred())
	Expect(username).To(Equal("taichi"))
}

func TestEmail(t *testing.T) {
	RegisterTestingT(t)

	reset := withGlobalGitConfigFile(`
[user]
    name  = taichi
    email = taichi@example.com
`)
	defer reset()

	var err error
	username, err := Email()
	Expect(err).NotTo(HaveOccurred())
	Expect(username).To(Equal("taichi@example.com"))
}

func TestGithubToken(t *testing.T) {
	RegisterTestingT(t)

	reset := withGlobalGitConfigFile(`
[github]
    token  = 16c999e8c71134401a78d4d46435517b2271d6ac
`)
	defer reset()

	var err error
	token, err := GithubToken()
	Expect(err).NotTo(HaveOccurred())
	Expect(token).To(Equal("16c999e8c71134401a78d4d46435517b2271d6ac"))
}

func TestOriginURL(t *testing.T) {
	RegisterTestingT(t)

	reset := withLocalGitConfigFile("remote.origin.url", "git@github.com:taichi/gitconfig.git")
	defer reset()

	var err error
	url, err := OriginURL()
	Expect(err).NotTo(HaveOccurred())
	Expect(url).To(Equal("git@github.com:taichi/gitconfig.git"))
}

func TestRepository(t *testing.T) {
	RegisterTestingT(t)

	reset := withLocalGitConfigFile("remote.origin.url", "git@github.com:taichi/gitconfig.git")
	defer reset()

	var err error
	repository, err := Repository()

	Expect(err).NotTo(HaveOccurred())
	Expect(repository).To(Equal("gitconfig"))
}

func TestRetrieveRepoName(t *testing.T) {
	RegisterTestingT(t)

	repo := retrieveRepoName("https://github.com/tcnksm/ghr.git")
	Expect(repo).To(Equal("ghr"))

	repo = retrieveRepoName("https://github.com/tcnksm/ghr")
	Expect(repo).To(Equal("ghr"))

	repo = retrieveRepoName("git@github.com:taichi/gitconfig.git")
	Expect(repo).To(Equal("gitconfig"))
}
