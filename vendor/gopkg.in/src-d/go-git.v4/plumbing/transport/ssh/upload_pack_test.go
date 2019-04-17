package ssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/test"

	"github.com/gliderlabs/ssh"
	stdssh "golang.org/x/crypto/ssh"
	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
)

type UploadPackSuite struct {
	test.UploadPackSuite
	fixtures.Suite

	port int
	base string
}

var _ = Suite(&UploadPackSuite{})

func (s *UploadPackSuite) SetUpSuite(c *C) {
	s.Suite.SetUpSuite(c)

	l, err := net.Listen("tcp", "localhost:0")
	c.Assert(err, IsNil)

	s.port = l.Addr().(*net.TCPAddr).Port
	s.base, err = ioutil.TempDir(os.TempDir(), fmt.Sprintf("go-git-ssh-%d", s.port))
	c.Assert(err, IsNil)

	DefaultAuthBuilder = func(user string) (AuthMethod, error) {
		return &Password{User: user}, nil
	}

	s.UploadPackSuite.Client = NewClient(&stdssh.ClientConfig{
		HostKeyCallback: stdssh.InsecureIgnoreHostKey(),
	})

	s.UploadPackSuite.Endpoint = s.prepareRepository(c, fixtures.Basic().One(), "basic.git")
	s.UploadPackSuite.EmptyEndpoint = s.prepareRepository(c, fixtures.ByTag("empty").One(), "empty.git")
	s.UploadPackSuite.NonExistentEndpoint = s.newEndpoint(c, "non-existent.git")

	server := &ssh.Server{Handler: handlerSSH}
	go func() {
		log.Fatal(server.Serve(l))
	}()
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
	ep, err := transport.NewEndpoint(fmt.Sprintf(
		"ssh://git@localhost:%d/%s/%s", s.port, filepath.ToSlash(s.base), name,
	))

	c.Assert(err, IsNil)
	return ep
}

func handlerSSH(s ssh.Session) {
	cmd, stdin, stderr, stdout, err := buildCommand(s.Command())
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, s)
	}()

	go func() {
		defer stderr.Close()
		io.Copy(s.Stderr(), stderr)
	}()

	defer stdout.Close()
	io.Copy(s, stdout)

	if err := cmd.Wait(); err != nil {
		return
	}
}

func buildCommand(c []string) (cmd *exec.Cmd, stdin io.WriteCloser, stderr, stdout io.ReadCloser, err error) {
	if len(c) != 2 {
		err = fmt.Errorf("invalid command")
		return
	}

	// fix for Windows environments
	path := strings.Replace(c[1], "/C:/", "C:/", 1)

	cmd = exec.Command(c[0], path)
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	stdin, err = cmd.StdinPipe()
	if err != nil {
		return
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		return
	}

	return
}
