package http

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

func Test(t *testing.T) { TestingT(t) }

type ClientSuite struct {
	Endpoint  *transport.Endpoint
	EmptyAuth transport.AuthMethod
}

var _ = Suite(&ClientSuite{})

func (s *ClientSuite) SetUpSuite(c *C) {
	var err error
	s.Endpoint, err = transport.NewEndpoint(
		"https://github.com/git-fixtures/basic",
	)
	c.Assert(err, IsNil)
}

func (s *UploadPackSuite) TestNewClient(c *C) {
	roundTripper := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cl := &http.Client{Transport: roundTripper}
	r, ok := NewClient(cl).(*client)
	c.Assert(ok, Equals, true)
	c.Assert(r.c, Equals, cl)
}

func (s *ClientSuite) TestNewBasicAuth(c *C) {
	a := &BasicAuth{"foo", "qux"}

	c.Assert(a.Name(), Equals, "http-basic-auth")
	c.Assert(a.String(), Equals, "http-basic-auth - foo:*******")
}

func (s *ClientSuite) TestNewTokenAuth(c *C) {
	a := &TokenAuth{"OAUTH-TOKEN-TEXT"}

	c.Assert(a.Name(), Equals, "http-token-auth")
	c.Assert(a.String(), Equals, "http-token-auth - *******")

	// Check header is set correctly
	req, err := http.NewRequest("GET", "https://github.com/git-fixtures/basic", nil)
	c.Assert(err, Equals, nil)
	a.setAuth(req)
	c.Assert(req.Header.Get("Authorization"), Equals, "Bearer OAUTH-TOKEN-TEXT")
}

func (s *ClientSuite) TestNewErrOK(c *C) {
	res := &http.Response{StatusCode: http.StatusOK}
	err := NewErr(res)
	c.Assert(err, IsNil)
}

func (s *ClientSuite) TestNewErrUnauthorized(c *C) {
	s.testNewHTTPError(c, http.StatusUnauthorized, "authentication required")
}

func (s *ClientSuite) TestNewErrForbidden(c *C) {
	s.testNewHTTPError(c, http.StatusForbidden, "authorization failed")
}

func (s *ClientSuite) TestNewErrNotFound(c *C) {
	s.testNewHTTPError(c, http.StatusNotFound, "repository not found")
}

func (s *ClientSuite) TestNewHTTPError40x(c *C) {
	s.testNewHTTPError(c, http.StatusPaymentRequired,
		"unexpected client error.*")
}

func (s *ClientSuite) testNewHTTPError(c *C, code int, msg string) {
	req, _ := http.NewRequest("GET", "foo", nil)
	res := &http.Response{
		StatusCode: code,
		Request:    req,
	}

	err := NewErr(res)
	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, msg)
}

func (s *ClientSuite) TestSetAuth(c *C) {
	auth := &BasicAuth{}
	r, err := DefaultClient.NewUploadPackSession(s.Endpoint, auth)
	c.Assert(err, IsNil)
	c.Assert(auth, Equals, r.(*upSession).auth)
}

type mockAuth struct{}

func (*mockAuth) Name() string   { return "" }
func (*mockAuth) String() string { return "" }

func (s *ClientSuite) TestSetAuthWrongType(c *C) {
	_, err := DefaultClient.NewUploadPackSession(s.Endpoint, &mockAuth{})
	c.Assert(err, Equals, transport.ErrInvalidAuthMethod)
}

type BaseSuite struct {
	fixtures.Suite

	base string
	host string
	port int
}

func (s *BaseSuite) SetUpTest(c *C) {
	l, err := net.Listen("tcp", "localhost:0")
	c.Assert(err, IsNil)

	base, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("go-git-http-%d", s.port))
	c.Assert(err, IsNil)

	s.port = l.Addr().(*net.TCPAddr).Port
	s.base = filepath.Join(base, s.host)

	err = os.MkdirAll(s.base, 0755)
	c.Assert(err, IsNil)

	cmd := exec.Command("git", "--exec-path")
	out, err := cmd.CombinedOutput()
	c.Assert(err, IsNil)

	server := &http.Server{
		Handler: &cgi.Handler{
			Path: filepath.Join(strings.Trim(string(out), "\n"), "git-http-backend"),
			Env:  []string{"GIT_HTTP_EXPORT_ALL=true", fmt.Sprintf("GIT_PROJECT_ROOT=%s", s.base)},
		},
	}
	go func() {
		log.Fatal(server.Serve(l))
	}()
}

func (s *BaseSuite) prepareRepository(c *C, f *fixtures.Fixture, name string) *transport.Endpoint {
	fs := f.DotGit()

	err := fixtures.EnsureIsBare(fs)
	c.Assert(err, IsNil)

	path := filepath.Join(s.base, name)
	err = os.Rename(fs.Root(), path)
	c.Assert(err, IsNil)

	return s.newEndpoint(c, name)
}

func (s *BaseSuite) newEndpoint(c *C, name string) *transport.Endpoint {
	ep, err := transport.NewEndpoint(fmt.Sprintf("http://localhost:%d/%s", s.port, name))
	c.Assert(err, IsNil)

	return ep
}

func (s *BaseSuite) TearDownTest(c *C) {
	err := os.RemoveAll(s.base)
	c.Assert(err, IsNil)
}
