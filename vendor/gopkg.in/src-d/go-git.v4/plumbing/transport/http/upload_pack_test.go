package http

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/protocol/packp"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/test"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type UploadPackSuite struct {
	test.UploadPackSuite
	BaseSuite
}

var _ = Suite(&UploadPackSuite{})

func (s *UploadPackSuite) SetUpSuite(c *C) {
	s.BaseSuite.SetUpTest(c)
	s.UploadPackSuite.Client = DefaultClient
	s.UploadPackSuite.Endpoint = s.prepareRepository(c, fixtures.Basic().One(), "basic.git")
	s.UploadPackSuite.EmptyEndpoint = s.prepareRepository(c, fixtures.ByTag("empty").One(), "empty.git")
	s.UploadPackSuite.NonExistentEndpoint = s.newEndpoint(c, "non-existent.git")
}

// Overwritten, different behaviour for HTTP.
func (s *UploadPackSuite) TestAdvertisedReferencesNotExists(c *C) {
	r, err := s.Client.NewUploadPackSession(s.NonExistentEndpoint, s.EmptyAuth)
	c.Assert(err, IsNil)
	info, err := r.AdvertisedReferences()
	c.Assert(err, Equals, transport.ErrRepositoryNotFound)
	c.Assert(info, IsNil)
}

func (s *UploadPackSuite) TestuploadPackRequestToReader(c *C) {
	r := packp.NewUploadPackRequest()
	r.Wants = append(r.Wants, plumbing.NewHash("d82f291cde9987322c8a0c81a325e1ba6159684c"))
	r.Wants = append(r.Wants, plumbing.NewHash("2b41ef280fdb67a9b250678686a0c3e03b0a9989"))
	r.Haves = append(r.Haves, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))

	sr, err := uploadPackRequestToReader(r)
	c.Assert(err, IsNil)
	b, _ := ioutil.ReadAll(sr)
	c.Assert(string(b), Equals,
		"0032want 2b41ef280fdb67a9b250678686a0c3e03b0a9989\n"+
			"0032want d82f291cde9987322c8a0c81a325e1ba6159684c\n0000"+
			"0032have 6ecf0ef2c2dffb796033e5a02219af86ec6584e5\n"+
			"0009done\n",
	)
}

func (s *UploadPackSuite) prepareRepository(c *C, f *fixtures.Fixture, name string) *transport.Endpoint {
	fs := f.DotGit()

	err := fixtures.EnsureIsBare(fs)
	c.Assert(err, IsNil)

	path := filepath.Join(s.base, name)
	err = os.Rename(fs.Root(), path)
	c.Assert(err, IsNil)

	return s.newEndpoint(c, name)
}

func (s *UploadPackSuite) newEndpoint(c *C, name string) *transport.Endpoint {
	ep, err := transport.NewEndpoint(fmt.Sprintf("http://localhost:%d/%s", s.port, name))
	c.Assert(err, IsNil)

	return ep
}

func (s *UploadPackSuite) TestAdvertisedReferencesRedirectPath(c *C) {
	endpoint, _ := transport.NewEndpoint("https://gitlab.com/gitlab-org/gitter/webapp")

	session, err := s.Client.NewUploadPackSession(endpoint, s.EmptyAuth)
	c.Assert(err, IsNil)

	info, err := session.AdvertisedReferences()
	c.Assert(err, IsNil)
	c.Assert(info, NotNil)

	url := session.(*upSession).endpoint.String()
	c.Assert(url, Equals, "https://gitlab.com/gitlab-org/gitter/webapp.git")
}

func (s *UploadPackSuite) TestAdvertisedReferencesRedirectSchema(c *C) {
	endpoint, _ := transport.NewEndpoint("http://github.com/git-fixtures/basic")

	session, err := s.Client.NewUploadPackSession(endpoint, s.EmptyAuth)
	c.Assert(err, IsNil)

	info, err := session.AdvertisedReferences()
	c.Assert(err, IsNil)
	c.Assert(info, NotNil)

	url := session.(*upSession).endpoint.String()
	c.Assert(url, Equals, "https://github.com/git-fixtures/basic")
}
