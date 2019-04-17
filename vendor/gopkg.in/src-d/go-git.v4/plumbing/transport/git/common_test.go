package git

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	fixtures.Suite

	base   string
	port   int
	daemon *exec.Cmd
}

func (s *BaseSuite) SetUpTest(c *C) {
	if runtime.GOOS == "windows" {
		c.Skip(`git for windows has issues with write operations through git:// protocol.
		See https://github.com/git-for-windows/git/issues/907`)
	}

	var err error
	s.port, err = freePort()
	c.Assert(err, IsNil)

	s.base, err = ioutil.TempDir(os.TempDir(), fmt.Sprintf("go-git-protocol-%d", s.port))
	c.Assert(err, IsNil)
}

func (s *BaseSuite) StartDaemon(c *C) {
	s.daemon = exec.Command(
		"git",
		"daemon",
		fmt.Sprintf("--base-path=%s", s.base),
		"--export-all",
		"--enable=receive-pack",
		"--reuseaddr",
		fmt.Sprintf("--port=%d", s.port),
		// Unless max-connections is limited to 1, a git-receive-pack
		// might not be seen by a subsequent operation.
		"--max-connections=1",
	)

	// Environment must be inherited in order to acknowledge GIT_EXEC_PATH if set.
	s.daemon.Env = os.Environ()

	err := s.daemon.Start()
	c.Assert(err, IsNil)

	// Connections might be refused if we start sending request too early.
	time.Sleep(time.Millisecond * 500)
}

func (s *BaseSuite) newEndpoint(c *C, name string) *transport.Endpoint {
	ep, err := transport.NewEndpoint(fmt.Sprintf("git://localhost:%d/%s", s.port, name))
	c.Assert(err, IsNil)

	return ep
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

func (s *BaseSuite) TearDownTest(c *C) {
	_ = s.daemon.Process.Signal(os.Kill)
	_ = s.daemon.Wait()

	err := os.RemoveAll(s.base)
	c.Assert(err, IsNil)
}

func freePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	return l.Addr().(*net.TCPAddr).Port, l.Close()
}
