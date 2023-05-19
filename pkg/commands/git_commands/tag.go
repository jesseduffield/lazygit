package git_commands

type TagCommands struct {
	*GitCommon
}

func NewTagCommands(gitCommon *GitCommon) *TagCommands {
	return &TagCommands{
		GitCommon: gitCommon,
	}
}

func (self *TagCommands) CreateLightweight(tagName string, ref string) error {
	cmdStr := NewGitCmd("tag").Arg("--", self.cmd.Quote(tagName)).
		ArgIf(len(ref) > 0, self.cmd.Quote(ref)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *TagCommands) CreateAnnotated(tagName, ref, msg string) error {
	cmdStr := NewGitCmd("tag").Arg(self.cmd.Quote(tagName)).
		ArgIf(len(ref) > 0, self.cmd.Quote(ref)).
		Arg("-m", self.cmd.Quote(msg)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *TagCommands) Delete(tagName string) error {
	cmdStr := NewGitCmd("tag").Arg("-d", self.cmd.Quote(tagName)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *TagCommands) Push(remoteName string, tagName string) error {
	cmdStr := NewGitCmd("push").Arg(self.cmd.Quote(remoteName), "tag", self.cmd.Quote(tagName)).
		ToString()

	return self.cmd.New(cmdStr).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}
