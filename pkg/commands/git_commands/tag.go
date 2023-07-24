package git_commands

import "github.com/jesseduffield/gocui"

type TagCommands struct {
	*GitCommon
}

func NewTagCommands(gitCommon *GitCommon) *TagCommands {
	return &TagCommands{
		GitCommon: gitCommon,
	}
}

func (self *TagCommands) CreateLightweight(tagName string, ref string, force bool) error {
	cmdArgs := NewGitCmd("tag").
		ArgIf(force, "--force").
		Arg("--", tagName).
		ArgIf(len(ref) > 0, ref).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *TagCommands) CreateAnnotated(tagName, ref, msg string, force bool) error {
	cmdArgs := NewGitCmd("tag").Arg(tagName).
		ArgIf(force, "--force").
		ArgIf(len(ref) > 0, ref).
		Arg("-m", msg).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *TagCommands) HasTag(tagName string) bool {
	cmdArgs := NewGitCmd("show-ref").
		Arg("--tags", "--quiet", "--verify", "--").
		Arg("refs/tags/" + tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run() == nil
}

func (self *TagCommands) Delete(tagName string) error {
	cmdArgs := NewGitCmd("tag").Arg("-d", tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *TagCommands) Push(task gocui.Task, remoteName string, tagName string) error {
	cmdArgs := NewGitCmd("push").Arg(remoteName, "tag", tagName).
		ToArgv()

	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).WithMutex(self.syncMutex).Run()
}
