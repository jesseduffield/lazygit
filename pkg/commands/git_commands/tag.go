package git_commands

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type TagCommands struct {
	*GitCommon
}

func NewTagCommands(gitCommon *GitCommon) *TagCommands {
	return &TagCommands{
		GitCommon: gitCommon,
	}
}

func (self *TagCommands) CreateLightweightObj(tagName string, ref string, force bool) *oscommands.CmdObj {
	cmdArgs := NewGitCmd("tag").
		ArgIf(force, "--force").
		Arg("--", tagName).
		ArgIf(len(ref) > 0, ref).
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *TagCommands) CreateAnnotatedObj(tagName, ref, msg string, force bool) *oscommands.CmdObj {
	cmdArgs := NewGitCmd("tag").Arg(tagName).
		ArgIf(force, "--force").
		ArgIf(len(ref) > 0, ref).
		Arg("-m", msg).
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *TagCommands) HasTag(tagName string) bool {
	cmdArgs := NewGitCmd("show-ref").
		Arg("--tags", "--quiet", "--verify", "--").
		Arg("refs/tags/" + tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run() == nil
}

func (self *TagCommands) LocalDelete(tagName string) error {
	cmdArgs := NewGitCmd("tag").Arg("-d", tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *TagCommands) Push(task gocui.Task, remoteName string, tagName string) error {
	cmdArgs := NewGitCmd("push").Arg(remoteName, "tag", tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}

// Return info about an annotated tag in the format:
//
//	Tagger:     tagger name <tagger email>
//	TaggerDate: tagger date
//
//	Tag message
//
// Should only be called for annotated tags.
func (self *TagCommands) ShowAnnotationInfo(tagName string) (string, error) {
	cmdArgs := NewGitCmd("for-each-ref").
		Arg("--format=Tagger:     %(taggername) %(taggeremail)%0aTaggerDate: %(taggerdate)%0a%0a%(contents)").
		Arg("refs/tags/" + tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).RunWithOutput()
}

func (self *TagCommands) IsTagAnnotated(tagName string) (bool, error) {
	cmdArgs := NewGitCmd("cat-file").
		Arg("-t").
		Arg("refs/tags/" + tagName).
		ToArgv()

	output, err := self.cmd.New(cmdArgs).RunWithOutput()
	return strings.TrimSpace(output) == "tag", err
}
