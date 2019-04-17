package ssh

import (
	"testing"

	"github.com/kevinburke/ssh_config"

	"golang.org/x/crypto/ssh"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

func Test(t *testing.T) { TestingT(t) }

func (s *SuiteCommon) TestOverrideConfig(c *C) {
	config := &ssh.ClientConfig{
		User: "foo",
		Auth: []ssh.AuthMethod{
			ssh.Password("yourpassword"),
		},
		HostKeyCallback: ssh.FixedHostKey(nil),
	}

	target := &ssh.ClientConfig{}
	overrideConfig(config, target)

	c.Assert(target.User, Equals, "foo")
	c.Assert(target.Auth, HasLen, 1)
	c.Assert(target.HostKeyCallback, NotNil)
}

func (s *SuiteCommon) TestOverrideConfigKeep(c *C) {
	config := &ssh.ClientConfig{
		User: "foo",
	}

	target := &ssh.ClientConfig{
		User: "bar",
	}

	overrideConfig(config, target)
	c.Assert(target.User, Equals, "foo")
}

func (s *SuiteCommon) TestDefaultSSHConfig(c *C) {
	defer func() {
		DefaultSSHConfig = ssh_config.DefaultUserSettings
	}()

	DefaultSSHConfig = &mockSSHConfig{map[string]map[string]string{
		"github.com": {
			"Hostname": "foo.local",
			"Port":     "42",
		},
	}}

	ep, err := transport.NewEndpoint("git@github.com:foo/bar.git")
	c.Assert(err, IsNil)

	cmd := &command{endpoint: ep}
	c.Assert(cmd.getHostWithPort(), Equals, "foo.local:42")
}

func (s *SuiteCommon) TestDefaultSSHConfigNil(c *C) {
	defer func() {
		DefaultSSHConfig = ssh_config.DefaultUserSettings
	}()

	DefaultSSHConfig = nil

	ep, err := transport.NewEndpoint("git@github.com:foo/bar.git")
	c.Assert(err, IsNil)

	cmd := &command{endpoint: ep}
	c.Assert(cmd.getHostWithPort(), Equals, "github.com:22")
}

func (s *SuiteCommon) TestDefaultSSHConfigWildcard(c *C) {
	defer func() {
		DefaultSSHConfig = ssh_config.DefaultUserSettings
	}()

	DefaultSSHConfig = &mockSSHConfig{Values: map[string]map[string]string{
		"*": {
			"Port": "42",
		},
	}}

	ep, err := transport.NewEndpoint("git@github.com:foo/bar.git")
	c.Assert(err, IsNil)

	cmd := &command{endpoint: ep}
	c.Assert(cmd.getHostWithPort(), Equals, "github.com:22")
}

type mockSSHConfig struct {
	Values map[string]map[string]string
}

func (c *mockSSHConfig) Get(alias, key string) string {
	a, ok := c.Values[alias]
	if !ok {
		return c.Values["*"][key]
	}

	return a[key]
}
