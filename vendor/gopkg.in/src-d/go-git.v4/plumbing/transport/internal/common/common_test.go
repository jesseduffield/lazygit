package common

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s *CommonSuite) TestIsRepoNotFoundErrorForUnknowSource(c *C) {
	msg := "unknown system is complaining of something very sad :("

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, false)
}

func (s *CommonSuite) TestIsRepoNotFoundErrorForGithub(c *C) {
	msg := fmt.Sprintf("%s : some error stuf", githubRepoNotFoundErr)

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, true)
}

func (s *CommonSuite) TestIsRepoNotFoundErrorForBitBucket(c *C) {
	msg := fmt.Sprintf("%s : some error stuf", bitbucketRepoNotFoundErr)

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, true)
}

func (s *CommonSuite) TestIsRepoNotFoundErrorForLocal(c *C) {
	msg := fmt.Sprintf("some error stuf : %s", localRepoNotFoundErr)

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, true)
}

func (s *CommonSuite) TestIsRepoNotFoundErrorForGitProtocolNotFound(c *C) {
	msg := fmt.Sprintf("%s : some error stuf", gitProtocolNotFoundErr)

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, true)
}

func (s *CommonSuite) TestIsRepoNotFoundErrorForGitProtocolNoSuch(c *C) {
	msg := fmt.Sprintf("%s : some error stuf", gitProtocolNoSuchErr)

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, true)
}

func (s *CommonSuite) TestIsRepoNotFoundErrorForGitProtocolAccessDenied(c *C) {
	msg := fmt.Sprintf("%s : some error stuf", gitProtocolAccessDeniedErr)

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, true)
}

func (s *CommonSuite) TestIsRepoNotFoundErrorForGogsAccessDenied(c *C) {
	msg := fmt.Sprintf("%s : some error stuf", gogsAccessDeniedErr)

	isRepoNotFound := isRepoNotFoundError(msg)

	c.Assert(isRepoNotFound, Equals, true)
}
