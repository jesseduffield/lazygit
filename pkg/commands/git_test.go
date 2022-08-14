package commands

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-errors/errors"
	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sasha-s/go-deadlock"
	"github.com/stretchr/testify/assert"
)

type fileInfoMock struct {
	name        string
	size        int64
	fileMode    os.FileMode
	fileModTime time.Time
	isDir       bool
	sys         interface{}
}

// Name is a function.
func (f fileInfoMock) Name() string {
	return f.name
}

// Size is a function.
func (f fileInfoMock) Size() int64 {
	return f.size
}

// Mode is a function.
func (f fileInfoMock) Mode() os.FileMode {
	return f.fileMode
}

// ModTime is a function.
func (f fileInfoMock) ModTime() time.Time {
	return f.fileModTime
}

// IsDir is a function.
func (f fileInfoMock) IsDir() bool {
	return f.isDir
}

// Sys is a function.
func (f fileInfoMock) Sys() interface{} {
	return f.sys
}

// TestNavigateToRepoRootDirectory is a function.
func TestNavigateToRepoRootDirectory(t *testing.T) {
	type scenario struct {
		testName string
		stat     func(string) (os.FileInfo, error)
		chdir    func(string) error
		test     func(error)
	}

	scenarios := []scenario{
		{
			"Navigate to git repository",
			func(string) (os.FileInfo, error) {
				return fileInfoMock{isDir: true}, nil
			},
			func(string) error {
				return nil
			},
			func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			"An error occurred when getting path information",
			func(string) (os.FileInfo, error) {
				return nil, fmt.Errorf("An error occurred")
			},
			func(string) error {
				return nil
			},
			func(err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "An error occurred")
			},
		},
		{
			"An error occurred when trying to move one path backward",
			func(string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			func(string) error {
				return fmt.Errorf("An error occurred")
			},
			func(err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "An error occurred")
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			s.test(navigateToRepoRootDirectory(s.stat, s.chdir))
		})
	}
}

// TestSetupRepository is a function.
func TestSetupRepository(t *testing.T) {
	type scenario struct {
		testName          string
		openGitRepository func(string, *gogit.PlainOpenOptions) (*gogit.Repository, error)
		errorStr          string
		options           gogit.PlainOpenOptions
		test              func(*gogit.Repository, error)
	}

	scenarios := []scenario{
		{
			"A gitconfig parsing error occurred",
			func(string, *gogit.PlainOpenOptions) (*gogit.Repository, error) {
				return nil, fmt.Errorf(`unquoted '\' must be followed by new line`)
			},
			"error translated",
			gogit.PlainOpenOptions{},
			func(r *gogit.Repository, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "error translated")
			},
		},
		{
			"A gogit error occurred",
			func(string, *gogit.PlainOpenOptions) (*gogit.Repository, error) {
				return nil, fmt.Errorf("Error from inside gogit")
			},
			"",
			gogit.PlainOpenOptions{},
			func(r *gogit.Repository, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "Error from inside gogit")
			},
		},
		{
			"Setup done properly",
			func(string, *gogit.PlainOpenOptions) (*gogit.Repository, error) {
				assert.NoError(t, os.RemoveAll("/tmp/lazygit-test"))
				r, err := gogit.PlainInit("/tmp/lazygit-test", false)
				assert.NoError(t, err)
				return r, nil
			},
			"",
			gogit.PlainOpenOptions{},
			func(r *gogit.Repository, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, r)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			s.test(setupRepository(s.openGitRepository, s.options, s.errorStr))
		})
	}
}

// TestNewGitCommand is a function.
func TestNewGitCommand(t *testing.T) {
	actual, err := os.Getwd()
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, os.Chdir(actual))
	}()

	type scenario struct {
		testName string
		setup    func()
		test     func(*GitCommand, error)
	}

	scenarios := []scenario{
		{
			"An error occurred, folder doesn't contains a git repository",
			func() {
				assert.NoError(t, os.Chdir("/tmp"))
			},
			func(gitCmd *GitCommand, err error) {
				assert.Error(t, err)
				assert.Regexp(t, `Must open lazygit in a git repository`, err.Error())
			},
		},
		{
			"New GitCommand object created",
			func() {
				assert.NoError(t, os.RemoveAll("/tmp/lazygit-test"))
				_, err := gogit.PlainInit("/tmp/lazygit-test", false)
				assert.NoError(t, err)
				assert.NoError(t, os.Chdir("/tmp/lazygit-test"))
			},
			func(gitCmd *GitCommand, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			s.setup()
			s.test(
				NewGitCommand(utils.NewDummyCommon(),
					oscommands.NewDummyOSCommand(),
					git_config.NewFakeGitConfig(nil),
					&deadlock.Mutex{},
				))
		})
	}
}

func TestFindDotGitDir(t *testing.T) {
	type scenario struct {
		testName string
		stat     func(string) (os.FileInfo, error)
		readFile func(filename string) ([]byte, error)
		test     func(string, error)
	}

	scenarios := []scenario{
		{
			".git is a directory",
			func(dotGit string) (os.FileInfo, error) {
				assert.Equal(t, ".git", dotGit)
				return os.Stat("testdata/a_dir")
			},
			func(dotGit string) ([]byte, error) {
				assert.Fail(t, "readFile should not be called if .git is a directory")
				return nil, nil
			},
			func(gitDir string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, ".git", gitDir)
			},
		},
		{
			".git is a file",
			func(dotGit string) (os.FileInfo, error) {
				assert.Equal(t, ".git", dotGit)
				return os.Stat("testdata/a_file")
			},
			func(dotGit string) ([]byte, error) {
				assert.Equal(t, ".git", dotGit)
				return []byte("gitdir: blah\n"), nil
			},
			func(gitDir string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "blah", gitDir)
			},
		},
		{
			"os.Stat returns an error",
			func(dotGit string) (os.FileInfo, error) {
				assert.Equal(t, ".git", dotGit)
				return nil, errors.New("error")
			},
			func(dotGit string) ([]byte, error) {
				assert.Fail(t, "readFile should not be called os.Stat returns an error")
				return nil, nil
			},
			func(gitDir string, err error) {
				assert.Error(t, err)
			},
		},
		{
			"readFile returns an error",
			func(dotGit string) (os.FileInfo, error) {
				assert.Equal(t, ".git", dotGit)
				return os.Stat("testdata/a_file")
			},
			func(dotGit string) ([]byte, error) {
				return nil, errors.New("error")
			},
			func(gitDir string, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			s.test(findDotGitDir(s.stat, s.readFile))
		})
	}
}
