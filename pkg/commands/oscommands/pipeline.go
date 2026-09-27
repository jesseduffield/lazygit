package oscommands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/samber/lo"
)

// Pipeline is a chain of running commands, each one's output feeding the next
// one's input, with the last one's output going somewhere the caller reads. Its
// method set is the one a render task expects of a command (see tasks.Cmd), so
// a pipeline can render a view just as a single command can.
type Pipeline struct {
	cmds   []*exec.Cmd
	cmdStr string
}

// StartPipeline starts the given commands wired A | B | C and returns the
// pipeline together with the reader for its output.
//
// Every command's stderr goes to that same output, so whatever a command
// complains about is part of what the caller reads. A diff renderer's error
// message belongs on screen with the diff it failed to render.
//
// Closing the reader is how a pipeline is brought down. The last command's next
// write fails, so it exits, and the failure travels back up the chain as each
// command in turn writes into a pipe whose reader is gone.
func (c *OSCommand) StartPipeline(cmdObjs ...*CmdObj) (*Pipeline, io.ReadCloser, error) {
	c.logPipeline(cmdObjs)

	cmds, parentEnds, err := wirePipeline(cmdObjs)
	if err != nil {
		return nil, nil, err
	}

	reader, writer, err := os.Pipe()
	if err != nil {
		closeAll(parentEnds)
		return nil, nil, err
	}
	for _, cmd := range cmds {
		cmd.Stderr = writer
	}
	cmds[len(cmds)-1].Stdout = writer
	parentEnds = append(parentEnds, writer)

	started, err := startPipeline(cmds, parentEnds)
	if err != nil {
		for _, cmd := range cmds[:started] {
			_ = cmd.Wait()
		}
		_ = reader.Close()
		return nil, nil, err
	}

	return &Pipeline{cmds: cmds, cmdStr: pipelineString(cmdObjs)}, reader, nil
}

func (self *Pipeline) String() string {
	return self.cmdStr
}

// Wait waits for every command to exit and reports the failure nearest the end
// of the pipeline. A command that fails leaves the ones before it writing into
// a pipe nobody reads, so their own broken-pipe failures are consequences of it
// rather than the cause worth reporting, while a command that fails early
// leaves the ones after it with nothing to read and no reason to fail at all.
func (self *Pipeline) Wait() error {
	var lastErr error
	for _, cmd := range self.cmds {
		if err := cmd.Wait(); err != nil {
			lastErr = fmt.Errorf("%s: %w", cmd.String(), err)
		}
	}

	return lastErr
}

// Terminate asks every command to stop, without waiting for any of them. On
// platforms where that does nothing, the pipeline comes down when its output
// reader is closed; see StartPipeline.
func (self *Pipeline) Terminate() error {
	var firstErr error
	for _, cmd := range self.cmds {
		if err := TerminateProcessGracefully(cmd.Process); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

// logPipeline enters a chain of commands into the command log, unless the first
// command was marked not to be logged; it speaks for the pipeline. A render runs
// its pipeline again on every selection change, so a caller has to be able to
// keep it out of the log.
func (c *OSCommand) logPipeline(cmdObjs []*CmdObj) {
	if cmdObjs[0].ShouldLog() {
		c.LogCommand(pipelineString(cmdObjs), true)
	}
}

// pipelineString names a chain of commands the way a shell would write it.
func pipelineString(cmdObjs []*CmdObj) string {
	return strings.Join(
		lo.Map(cmdObjs, func(cmdObj *CmdObj, _ int) string {
			return cmdObj.ToString()
		}),
		" | ",
	)
}

// wirePipeline connects each command's output to the next one's input, like
// A | B | C, and returns the commands along with the parent's ends of those
// pipes. The last command's output is left for the caller to direct.
//
// The parent's ends have to be closed once the commands are running.
// startPipeline does that; see there for why it matters.
func wirePipeline(cmdObjs []*CmdObj) ([]*exec.Cmd, []io.Closer, error) {
	cmds := lo.Map(cmdObjs, func(cmdObj *CmdObj, _ int) *exec.Cmd {
		return cmdObj.GetCmd()
	})

	parentEnds := []io.Closer{}
	for i := range len(cmds) - 1 {
		reader, writer, err := os.Pipe()
		if err != nil {
			closeAll(parentEnds)
			return nil, nil, err
		}

		cmds[i].Stdout = writer
		cmds[i+1].Stdin = reader
		parentEnds = append(parentEnds, reader, writer)
	}

	return cmds, parentEnds, nil
}

// startPipeline starts every command and reports how many it got going. Every
// one is started before any of them is waited for: waiting closes our end of
// the pipe that feeds the next command, and one that hasn't been started by
// then would inherit a closed stdin.
//
// Once they are all running, each of them holds its own ends of the pipes it
// reads and writes, and the parent lets go of its copies. Both directions
// matter. While the parent holds the read end of a link, a command writing
// into it never learns that the command meant to read it is gone, and keeps
// running after the pipeline has been brought down. While the parent holds the
// write end, the command reading it never reaches the end of its input.
//
// When a command fails to start, the ones already running are killed, since
// without the rest of the pipeline to drain them they could block forever
// writing to a full pipe. They still have to be reaped, so the count covers
// them too.
func startPipeline(cmds []*exec.Cmd, parentEnds []io.Closer) (int, error) {
	defer closeAll(parentEnds)

	for i, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			for _, started := range cmds[:i] {
				_ = started.Process.Kill()
			}

			return i, err
		}
	}

	return len(cmds), nil
}

func closeAll(closers []io.Closer) {
	for _, closer := range closers {
		_ = closer.Close()
	}
}
