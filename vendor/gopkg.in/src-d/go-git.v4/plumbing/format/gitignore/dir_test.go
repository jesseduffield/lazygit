package gitignore

import (
	"os"
	"os/user"
	"strconv"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
)

type MatcherSuite struct {
	GFS  billy.Filesystem // git repository root
	RFS  billy.Filesystem // root that contains user home
	MCFS billy.Filesystem // root that contains user home, but missing ~/.gitconfig
	MEFS billy.Filesystem // root that contains user home, but missing excludesfile entry
	MIFS billy.Filesystem // root that contains user home, but missing .gitnignore

	SFS  billy.Filesystem // root that contains /etc/gitconfig
}

var _ = Suite(&MatcherSuite{})

func (s *MatcherSuite) SetUpTest(c *C) {
	// setup generic git repository root
	fs := memfs.New()
	f, err := fs.Create(".gitignore")
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("vendor/g*/\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	err = fs.MkdirAll("vendor", os.ModePerm)
	c.Assert(err, IsNil)
	f, err = fs.Create("vendor/.gitignore")
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("!github.com/\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	fs.MkdirAll("another", os.ModePerm)
	fs.MkdirAll("vendor/github.com", os.ModePerm)
	fs.MkdirAll("vendor/gopkg.in", os.ModePerm)

	s.GFS = fs

	// setup root that contains user home
	usr, err := user.Current()
	c.Assert(err, IsNil)

	fs = memfs.New()
	err = fs.MkdirAll(usr.HomeDir, os.ModePerm)
	c.Assert(err, IsNil)

	f, err = fs.Create(fs.Join(usr.HomeDir, gitconfigFile))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("[core]\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("	excludesfile = " + strconv.Quote(fs.Join(usr.HomeDir, ".gitignore_global")) + "\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	f, err = fs.Create(fs.Join(usr.HomeDir, ".gitignore_global"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("# IntelliJ\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte(".idea/\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("*.iml\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	s.RFS = fs

	// root that contains user home, but missing ~/.gitconfig
	fs = memfs.New()
	err = fs.MkdirAll(usr.HomeDir, os.ModePerm)
	c.Assert(err, IsNil)

	f, err = fs.Create(fs.Join(usr.HomeDir, ".gitignore_global"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("# IntelliJ\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte(".idea/\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("*.iml\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	s.MCFS = fs

	// setup root that contains user home, but missing excludesfile entry
	fs = memfs.New()
	err = fs.MkdirAll(usr.HomeDir, os.ModePerm)
	c.Assert(err, IsNil)

	f, err = fs.Create(fs.Join(usr.HomeDir, gitconfigFile))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("[core]\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	f, err = fs.Create(fs.Join(usr.HomeDir, ".gitignore_global"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("# IntelliJ\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte(".idea/\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("*.iml\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	s.MEFS = fs

	// setup root that contains user home, but missing .gitnignore
	fs = memfs.New()
	err = fs.MkdirAll(usr.HomeDir, os.ModePerm)
	c.Assert(err, IsNil)

	f, err = fs.Create(fs.Join(usr.HomeDir, gitconfigFile))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("[core]\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("	excludesfile = " + strconv.Quote(fs.Join(usr.HomeDir, ".gitignore_global")) + "\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	s.MIFS = fs

	// setup root that contains user home
	fs = memfs.New()
	err = fs.MkdirAll("etc", os.ModePerm)
	c.Assert(err, IsNil)

	f, err = fs.Create(systemFile)
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("[core]\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("	excludesfile = /etc/gitignore_global\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	f, err = fs.Create("/etc/gitignore_global")
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("# IntelliJ\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte(".idea/\n"))
	c.Assert(err, IsNil)
	_, err = f.Write([]byte("*.iml\n"))
	c.Assert(err, IsNil)
	err = f.Close()
	c.Assert(err, IsNil)

	s.SFS = fs
}

func (s *MatcherSuite) TestDir_ReadPatterns(c *C) {
	ps, err := ReadPatterns(s.GFS, nil)
	c.Assert(err, IsNil)
	c.Assert(ps, HasLen, 2)

	m := NewMatcher(ps)
	c.Assert(m.Match([]string{"vendor", "gopkg.in"}, true), Equals, true)
	c.Assert(m.Match([]string{"vendor", "github.com"}, true), Equals, false)
}

func (s *MatcherSuite) TestDir_LoadGlobalPatterns(c *C) {
	ps, err := LoadGlobalPatterns(s.RFS)
	c.Assert(err, IsNil)
	c.Assert(ps, HasLen, 2)

	m := NewMatcher(ps)
	c.Assert(m.Match([]string{"go-git.v4.iml"}, true), Equals, true)
	c.Assert(m.Match([]string{".idea"}, true), Equals, true)
}

func (s *MatcherSuite) TestDir_LoadGlobalPatternsMissingGitconfig(c *C) {
	ps, err := LoadGlobalPatterns(s.MCFS)
	c.Assert(err, IsNil)
	c.Assert(ps, HasLen, 0)
}

func (s *MatcherSuite) TestDir_LoadGlobalPatternsMissingExcludesfile(c *C) {
	ps, err := LoadGlobalPatterns(s.MEFS)
	c.Assert(err, IsNil)
	c.Assert(ps, HasLen, 0)
}

func (s *MatcherSuite) TestDir_LoadGlobalPatternsMissingGitignore(c *C) {
	ps, err := LoadGlobalPatterns(s.MIFS)
	c.Assert(err, IsNil)
	c.Assert(ps, HasLen, 0)
}

func (s *MatcherSuite) TestDir_LoadSystemPatterns(c *C) {
	ps, err := LoadSystemPatterns(s.SFS)
	c.Assert(err, IsNil)
	c.Assert(ps, HasLen, 2)

	m := NewMatcher(ps)
	c.Assert(m.Match([]string{"go-git.v4.iml"}, true), Equals, true)
	c.Assert(m.Match([]string{".idea"}, true), Equals, true)
}
