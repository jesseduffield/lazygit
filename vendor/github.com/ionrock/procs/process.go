// Procs is a library to make working with command line applications a
// little nicer.
//
// The goal is to expand on the os/exec package by providing some
// features usually accomplished in a shell, without having to resort to
// a shell. Procs also tries to make working with output simpler by
// providing a simple line handler API over working with io pipes.
//
// Finally, while the hope is that procs provides some convenience, it
// is also a goal to help make it easier to write more secure
// code. For example, avoiding a shell and the ability to manage the
// environment as a map[string]string are both measures that intend to
// make it easier to accomplish things like avoiding outputting
// secrets and opening the door for MITM attacks. With that said, it is
// always important to consider the security implications, especially
// when you are working with untrusted input or sensitive data.
package procs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

// OutHandler defines the interface for writing output handlers for
// Process objects.
type OutHandler func(string) string

// Process is intended to be used like exec.Cmd where possible.
type Process struct {
	// CmdString takes a string and parses it into the relevant cmds
	CmdString string

	// Cmds is the list of command delmited by pipes.
	Cmds []*exec.Cmd

	// Env provides a map[string]string that can mutated before
	// running a command.
	Env map[string]string

	// Dir defines the directory the command should run in. The
	// Default is the current dir.
	Dir string

	// OutputHandler can be defined to perform any sort of processing
	// on the output. The simple interface is to accept a string (a
	// line of output) and return a string that will be included in the
	// buffered output and/or output written to stdout.'
	//
	// For example defining the Process as:
	//
	//     prefix := "myapp"
	//     p := &procs.Process{
	//         OutputHandler: func(line string) string {
	//             return fmt.Sprintf("%s | %s", prefix, line)
	//         },
	//     }
	//
	// This would prefix the stdout lines with a "myapp | ".
	//
	// By the default, this function is nil and will be skipped, with
	// the unchanged line being added to the respective output buffer.
	OutputHandler OutHandler

	// ErrHandler is a OutputHandler for stderr.
	ErrHandler OutHandler

	// When no output is given, we'll buffer output in these vars.
	errBuffer bytes.Buffer
	outBuffer bytes.Buffer

	// When a output handler is provided, we ensure we're handling a
	// single line at at time.
	outputWait *sync.WaitGroup
}

// NewProcess creates a new *Process from a command string.
//
// It is assumed that the user will mutate the resulting *Process by
// setting the necessary attributes.
func NewProcess(command string) *Process {
	return &Process{CmdString: command}
}

// internal expand method to use the proc env.
func (p *Process) expand(s string) string {
	return os.Expand(s, func(key string) string {
		v, _ := p.Env[key]
		return v
	})
}

// addCmd adds a new command to the list of commands, ensuring the Dir
// and Env have been added to the underlying *exec.Cmd instances.
func (p *Process) addCmd(cmdparts []string) {
	var cmd *exec.Cmd
	if len(cmdparts) == 1 {
		cmd = exec.Command(cmdparts[0])
	} else {
		cmd = exec.Command(cmdparts[0], cmdparts[1:]...)
	}

	if p.Dir != "" {
		cmd.Dir = p.Dir
	}

	if p.Env != nil {
		env := []string{}
		for k, v := range p.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, p.expand(v)))
		}

		cmd.Env = env
	}

	p.Cmds = append(p.Cmds, cmd)
}

// findCmds parses the CmdString to find the commands that should be
// run by spliting the lexically parsed command by pipes ("|").
func (p *Process) findCmds() {
	// Skip if the cmd set is already set. This allows manual creation
	// of piped commands.
	if len(p.Cmds) > 0 {
		return
	}

	if p.CmdString == "" {
		return
	}

	parts := SplitCommand(p.CmdString)
	for i := range parts {
		parts[i] = p.expand(parts[i])
	}

	cmd := []string{}
	for _, part := range parts {
		if part == "|" {
			p.addCmd(cmd)
			cmd = []string{}
		} else {
			cmd = append(cmd, part)
		}
	}

	p.addCmd(cmd)
}

// lineReader takes will read a line in the io.Reader and write to the
// Process output buffer and use any OutputHandler that exists.
func (p *Process) lineReader(wg *sync.WaitGroup, r io.Reader, w *bytes.Buffer, handler OutHandler) {
	defer wg.Done()

	reader := bufio.NewReader(r)
	var buffer bytes.Buffer

	for {
		buf := make([]byte, 1024)

		n, err := reader.Read(buf)
		if err != nil {
			return
		}

		buf = buf[:n]

		for {
			i := bytes.IndexByte(buf, '\n')
			if i < 0 {
				break
			}

			buffer.Write(buf[0:i])
			outLine := buffer.String()
			if handler != nil {
				outLine = handler(outLine)
			}
			w.WriteString(outLine)
			buffer.Reset()
			buf = buf[i+1:]
		}
		buffer.Write(buf)
	}
}

// checkErr shortens the creation of the pipes by bailing out with a
// log.Fatal.
func checkErr(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (p *Process) setupPipes() error {
	last := len(p.Cmds) - 1

	if last != 0 {
		for i, cmd := range p.Cmds[:last] {
			var err error

			p.Cmds[i+1].Stdin, err = cmd.StdoutPipe()
			if err != nil {
				fmt.Printf("error creating stdout pipe: %s\n", err)
				return err
			}

			cmd.Stderr = &p.errBuffer
		}
	}

	cmd := p.Cmds[last]
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("error creating stdout pipe: %s\n", err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("error creating stderr pipe: %s\n", err)
		return err
	}

	p.outputWait = new(sync.WaitGroup)
	p.outputWait.Add(2)

	// These close the stdout/err channels
	go p.lineReader(p.outputWait, stdout, &p.outBuffer, p.OutputHandler)
	go p.lineReader(p.outputWait, stderr, &p.errBuffer, p.ErrHandler)

	return nil
}

// Run executes the cmds and returns the output as a string and any error.
func (p *Process) Run() error {
	if err := p.Start(); err != nil {
		return err
	}

	return p.Wait()
}

// Start will start the list of cmds.
func (p *Process) Start() error {
	p.findCmds()
	p.setupPipes()

	for i, cmd := range p.Cmds {
		err := cmd.Start()
		if err != nil {
			defer func() {
				for _, precmd := range p.Cmds[0:i] {
					precmd.Wait()
				}
			}()
			return err
		}
	}

	return nil
}

// Wait will block, waiting for the commands to finish.
func (p *Process) Wait() error {
	if p.outputWait != nil {
		p.outputWait.Wait()
	}

	var err error
	for _, cmd := range p.Cmds {
		err = cmd.Wait()
	}
	return err
}

// Stop tries to stop the process.
func (p *Process) Stop() error {
	for _, cmd := range p.Cmds {
		// ProcessState means it is already exited.
		if cmd.ProcessState != nil {
			continue
		}

		err := cmd.Process.Kill()
		if err != nil {
			return err
		}
	}

	return nil
}

// Output returns the buffered output as []byte.
func (p *Process) Output() ([]byte, error) {
	return p.outBuffer.Bytes(), nil
}

// ErrOutput returns the buffered stderr as []byte
func (p *Process) ErrOutput() ([]byte, error) {
	return p.errBuffer.Bytes(), nil
}
