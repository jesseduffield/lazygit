package git_commands

import (
	"fmt"
)

type BlameCommands struct {
	*GitCommon
}

func NewBlameCommands(gitCommon *GitCommon) *BlameCommands {
	return &BlameCommands{
		GitCommon: gitCommon,
	}
}

// Blame a range of lines. For each line, output the hash of the commit where
// the line last changed, then a space, then a description of the commit (author
// and date), another space, and then the line. For example:
//
//	ac90ebac688fe8bc2ffd922157a9d2c54681d2aa (Stefan Haller 2023-08-01 14:54:56 +0200 11) func NewBlameCommands(gitCommon *GitCommon) *BlameCommands {
//	ac90ebac688fe8bc2ffd922157a9d2c54681d2aa (Stefan Haller 2023-08-01 14:54:56 +0200 12) 	return &BlameCommands{
//	ac90ebac688fe8bc2ffd922157a9d2c54681d2aa (Stefan Haller 2023-08-01 14:54:56 +0200 13) 		GitCommon: gitCommon,
func (self *BlameCommands) BlameLineRange(filename string, commit string, firstLine int, numLines int) (string, error) {
	cmdArgs := NewGitCmd("blame").
		Arg("-l").
		Arg(fmt.Sprintf("-L%d,+%d", firstLine, numLines)).
		Arg(commit).
		Arg("--").
		Arg(filename)

	return self.cmd.New(cmdArgs.ToArgv()).RunWithOutput()
}
