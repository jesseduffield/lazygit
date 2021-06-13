package oscommands

import (
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

type CmdLogEntry struct {
	// e.g. 'git commit -m "haha"'
	cmdStr string
	// Span is something like 'Staging File'. Multiple commands can be grouped under the same
	// span
	span string

	// sometimes our command is direct like 'git commit', and sometimes it's a
	// command to remove a file but through Go's standard library rather than the
	// command line
	commandLine bool
}

func (e CmdLogEntry) GetCmdStr() string {
	return e.cmdStr
}

func (e CmdLogEntry) GetSpan() string {
	return e.span
}

func (e CmdLogEntry) GetCommandLine() bool {
	return e.commandLine
}

func NewCmdLogEntry(cmdStr string, span string, commandLine bool) CmdLogEntry {
	return CmdLogEntry{cmdStr: cmdStr, span: span, commandLine: commandLine}
}

func NewCmdLogEntryFromCmdObj(cmdObj ICmdObj, span string) CmdLogEntry {
	return CmdLogEntry{cmdStr: cmdObj.ToString(), span: span, commandLine: true}
}
