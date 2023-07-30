package oscommands

import (
	"os/exec"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
)

// A command object is a general way to represent a command to be run on the
// command line.
type ICmdObj interface {
	GetCmd() *exec.Cmd
	// outputs string representation of command. Note that if the command was built
	// using NewFromArgs, the output won't be quite the same as what you would type
	// into a terminal e.g. 'sh -c git commit' as opposed to 'sh -c "git commit"'
	ToString() string

	// outputs args vector e.g. ["git", "commit", "-m", "my message"]
	Args() []string

	AddEnvVars(...string) ICmdObj
	GetEnvVars() []string

	// sets the working directory
	SetWd(string) ICmdObj

	// runs the command and returns an error if any
	Run() error
	// runs the command and returns the output as a string, and an error if any
	RunWithOutput() (string, error)
	// runs the command and returns stdout and stderr as a string, and an error if any
	RunWithOutputs() (string, string, error)
	// runs the command and runs a callback function on each line of the output. If the callback returns true for the boolean value, we kill the process and return.
	RunAndProcessLines(onLine func(line string) (bool, error)) error

	// Be calling DontLog(), we're saying that once we call Run(), we don't want to
	// log the command in the UI (it'll still be logged in the log file). The general rule
	// is that if a command doesn't change the git state (e.g. read commands like `git diff`)
	// then we don't want to log it. If we are changing something (e.g. `git add .`) then
	// we do. The only exception is if we're running a command in the background periodically
	// like `git fetch`, which technically does mutate stuff but isn't something we need
	// to notify the user about.
	DontLog() ICmdObj

	// This returns false if DontLog() was called
	ShouldLog() bool

	// when you call this, then call Run(), we'll stream the output to the cmdWriter (i.e. the command log panel)
	StreamOutput() ICmdObj
	// returns true if StreamOutput() was called
	ShouldStreamOutput() bool

	// if you call this before ShouldStreamOutput we'll consider an error with no
	// stderr content as a non-error. Not yet supported for Run or RunWithOutput (
	// but adding support is trivial)
	IgnoreEmptyError() ICmdObj
	// returns true if IgnoreEmptyError() was called
	ShouldIgnoreEmptyError() bool

	PromptOnCredentialRequest(task gocui.Task) ICmdObj
	FailOnCredentialRequest() ICmdObj

	WithMutex(mutex *deadlock.Mutex) ICmdObj
	Mutex() *deadlock.Mutex

	GetCredentialStrategy() CredentialStrategy
	GetTask() gocui.Task

	Clone() ICmdObj
}

type CmdObj struct {
	cmd *exec.Cmd

	runner ICmdObjRunner

	// see DontLog()
	dontLog bool

	// see StreamOutput()
	streamOutput bool

	// see IgnoreEmptyError()
	ignoreEmptyError bool

	// if set to true, it means we might be asked to enter a username/password by this command.
	credentialStrategy CredentialStrategy
	task               gocui.Task

	// can be set so that we don't run certain commands simultaneously
	mutex *deadlock.Mutex
}

type CredentialStrategy int

const (
	// do not expect a credential request. If we end up getting one
	// we'll be in trouble because the command will hang indefinitely
	NONE CredentialStrategy = iota
	// expect a credential request and if we get one, prompt the user to enter their username/password
	PROMPT
	// in this case we will check for a credential request (i.e. the command pauses to ask for
	// username/password) and if we get one, we just submit a newline, forcing the
	// command to fail. We use this e.g. for a background `git fetch` to prevent it
	// from hanging indefinitely.
	FAIL
)

var _ ICmdObj = &CmdObj{}

func (self *CmdObj) GetCmd() *exec.Cmd {
	return self.cmd
}

func (self *CmdObj) ToString() string {
	// if a given arg contains a space, we need to wrap it in quotes
	quotedArgs := lo.Map(self.cmd.Args, func(arg string, _ int) string {
		if strings.Contains(arg, " ") {
			return `"` + arg + `"`
		}
		return arg
	})

	return strings.Join(quotedArgs, " ")
}

func (self *CmdObj) Args() []string {
	return self.cmd.Args
}

func (self *CmdObj) AddEnvVars(vars ...string) ICmdObj {
	self.cmd.Env = append(self.cmd.Env, vars...)

	return self
}

func (self *CmdObj) GetEnvVars() []string {
	return self.cmd.Env
}

func (self *CmdObj) SetWd(wd string) ICmdObj {
	self.cmd.Dir = wd

	return self
}

func (self *CmdObj) DontLog() ICmdObj {
	self.dontLog = true
	return self
}

func (self *CmdObj) ShouldLog() bool {
	return !self.dontLog
}

func (self *CmdObj) StreamOutput() ICmdObj {
	self.streamOutput = true

	return self
}

func (self *CmdObj) ShouldStreamOutput() bool {
	return self.streamOutput
}

func (self *CmdObj) IgnoreEmptyError() ICmdObj {
	self.ignoreEmptyError = true

	return self
}

func (self *CmdObj) Mutex() *deadlock.Mutex {
	return self.mutex
}

func (self *CmdObj) WithMutex(mutex *deadlock.Mutex) ICmdObj {
	self.mutex = mutex

	return self
}

func (self *CmdObj) ShouldIgnoreEmptyError() bool {
	return self.ignoreEmptyError
}

func (self *CmdObj) Run() error {
	return self.runner.Run(self)
}

func (self *CmdObj) RunWithOutput() (string, error) {
	return self.runner.RunWithOutput(self)
}

func (self *CmdObj) RunWithOutputs() (string, string, error) {
	return self.runner.RunWithOutputs(self)
}

func (self *CmdObj) RunAndProcessLines(onLine func(line string) (bool, error)) error {
	return self.runner.RunAndProcessLines(self, onLine)
}

func (self *CmdObj) PromptOnCredentialRequest(task gocui.Task) ICmdObj {
	self.credentialStrategy = PROMPT
	self.task = task

	return self
}

func (self *CmdObj) FailOnCredentialRequest() ICmdObj {
	self.credentialStrategy = FAIL

	return self
}

func (self *CmdObj) GetCredentialStrategy() CredentialStrategy {
	return self.credentialStrategy
}

func (self *CmdObj) GetTask() gocui.Task {
	return self.task
}

func (self *CmdObj) Clone() ICmdObj {
	clone := &CmdObj{}
	*clone = *self
	clone.cmd = cloneCmd(self.cmd)
	return clone
}

func cloneCmd(cmd *exec.Cmd) *exec.Cmd {
	clone := &exec.Cmd{}
	*clone = *cmd

	return clone
}
