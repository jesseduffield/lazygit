//go:build windows
// +build windows

package pty

import (
	"context"
	"io"
	"os"
	"syscall"

	_ "unsafe" // for go:linkname
)

// copied from os/exec.Cmd for platform compatibility
// we need to use startupInfoEx for pty support, but os/exec.Cmd only have
// support for startupInfo on windows, so we have to rewrite some internal
// logic for windows while keep its behavior compatible with other platforms.

// cmd represents an external command being prepared or run.
//
// A cmd cannot be reused after calling its Run, Output or CombinedOutput
// methods.
//go:linkname cmd os/exec.Cmd
type cmd struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Env specifies the environment of the process.
	// Each entry is of the form "key=value".
	// If Env is nil, the new process uses the current process's
	// environment.
	// If Env contains duplicate environment keys, only the last
	// value in the slice for each duplicate key is used.
	// As a special case on Windows, SYSTEMROOT is always added if
	// missing and not explicitly set to the empty string.
	Env []string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, Run runs the command in the
	// calling process's current directory.
	Dir string

	// Stdin specifies the process's standard input.
	//
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// If Stdin is an *os.File, the process's standard input is connected
	// directly to that file.
	//
	// Otherwise, during the execution of the command a separate
	// goroutine reads from Stdin and delivers that data to the command
	// over a pipe. In this case, Wait does not complete until the goroutine
	// stops copying, either because it has reached the end of Stdin
	// (EOF or a read error) or because writing to the pipe returned an error.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If either is an *os.File, the corresponding output from the process
	// is connected directly to that file.
	//
	// Otherwise, during the execution of the command a separate goroutine
	// reads from the process over a pipe and delivers that data to the
	// corresponding Writer. In this case, Wait does not complete until the
	// goroutine reaches EOF or encounters an error.
	//
	// If Stdout and Stderr are the same writer, and have a type that can
	// be compared with ==, at most one goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer

	// ExtraFiles specifies additional open files to be inherited by the
	// new process. It does not include standard input, standard output, or
	// standard error. If non-nil, entry i becomes file descriptor 3+i.
	//
	// ExtraFiles is not supported on Windows.
	ExtraFiles []*os.File

	// SysProcAttr holds optional, operating system-specific attributes.
	// Run passes it to os.StartProcess as the os.ProcAttr's Sys field.
	SysProcAttr *syscall.SysProcAttr

	// Process is the underlying process, once started.
	Process *os.Process

	// ProcessState contains information about an exited process,
	// available after a call to Wait or Run.
	ProcessState *os.ProcessState

	ctx             context.Context // nil means none
	lookPathErr     error           // LookPath error, if any.
	finished        bool            // when Wait was called
	childFiles      []*os.File
	closeAfterStart []io.Closer
	closeAfterWait  []io.Closer
	goroutine       []func() error
	errch           chan error // one send per goroutine
	waitDone        chan struct{}
}

//go:linkname _cmd_closeDescriptors os/exec.(*Cmd).closeDescriptors
func _cmd_closeDescriptors(c *cmd, closers []io.Closer)

//go:linkname _cmd_envv os/exec.(*Cmd).envv
func _cmd_envv(c *cmd) ([]string, error)

//go:linkname _cmd_argv os/exec.(*Cmd).argv
func _cmd_argv(c *cmd) []string

//go:linkname lookExtensions os/exec.lookExtensions
func lookExtensions(path, dir string) (string, error)

//go:linkname dedupEnv os/exec.dedupEnv
func dedupEnv(env []string) []string

//go:linkname addCriticalEnv os/exec.addCriticalEnv
func addCriticalEnv(env []string) []string

//go:linkname newProcess os.newProcess
func newProcess(pid int, handle uintptr) *os.Process

//go:linkname execEnvDefault internal/syscall/execenv.Default
func execEnvDefault(sys *syscall.SysProcAttr) (env []string, err error)

//go:linkname createEnvBlock syscall.createEnvBlock
func createEnvBlock(envv []string) *uint16

//go:linkname makeCmdLine syscall.makeCmdLine
func makeCmdLine(args []string) string

//go:linkname joinExeDirAndFName syscall.joinExeDirAndFName
func joinExeDirAndFName(dir, p string) (name string, err error)
