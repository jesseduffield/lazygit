package daemon

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/env"
)

// Sometimes lazygit will be invoked in daemon mode from a parent lazygit process.
// We do this when git lets us supply a program to run within a git command.
// For example, if we want to ensure that a git command doesn't hang due to
// waiting for an editor to save a commit message, we can tell git to invoke lazygit
// as the editor via 'GIT_EDITOR=lazygit', and use the env var
// 'LAZYGIT_DAEMON_KIND=EXIT_IMMEDIATELY' to specify that we want to run lazygit
// as a daemon which simply exits immediately. Any additional arguments we want
// to pass to a daemon can be done via other env vars.

type DaemonKind string

const (
	InteractiveRebase DaemonKind = "INTERACTIVE_REBASE"
	ExitImmediately   DaemonKind = "EXIT_IMMEDIATELY"
)

const (
	DaemonKindEnvKey string = "LAZYGIT_DAEMON_KIND"
	RebaseTODOEnvKey string = "LAZYGIT_REBASE_TODO"
)

type Daemon interface {
	Run() error
}

func Handle(common *common.Common) {
	d := getDaemon(common)
	if d == nil {
		return
	}

	if err := d.Run(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func InDaemonMode() bool {
	return getDaemonKind() != ""
}

func getDaemon(common *common.Common) Daemon {
	switch getDaemonKind() {
	case InteractiveRebase:
		return &rebaseDaemon{c: common}
	case ExitImmediately:
		return &exitImmediatelyDaemon{c: common}
	}

	return nil
}

func getDaemonKind() DaemonKind {
	return DaemonKind(os.Getenv(DaemonKindEnvKey))
}

type rebaseDaemon struct {
	c *common.Common
}

func (self *rebaseDaemon) Run() error {
	self.c.Log.Info("Lazygit invoked as interactive rebase demon")
	self.c.Log.Info("args: ", os.Args)

	if strings.HasSuffix(os.Args[1], "git-rebase-todo") {
		if err := ioutil.WriteFile(os.Args[1], []byte(os.Getenv(RebaseTODOEnvKey)), 0o644); err != nil {
			return err
		}
	} else if strings.HasSuffix(os.Args[1], filepath.Join(gitDir(), "COMMIT_EDITMSG")) { // TODO: test
		// if we are rebasing and squashing, we'll see a COMMIT_EDITMSG
		// but in this case we don't need to edit it, so we'll just return
	} else {
		self.c.Log.Info("Lazygit demon did not match on any use cases")
	}

	return nil
}

func gitDir() string {
	dir := env.GetGitDirEnv()
	if dir == "" {
		return ".git"
	}
	return dir
}

type exitImmediatelyDaemon struct {
	c *common.Common
}

func (self *exitImmediatelyDaemon) Run() error {
	return nil
}
