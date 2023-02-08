package git_commands

import "fmt"

type CloneCommands struct {
	*GitCommon
}

func NewCloneCommands(gitCommon *GitCommon) *CloneCommands {
	return &CloneCommands{
		GitCommon: gitCommon,
	}
}

func (self *CloneCommands) Clone(url string, destination string) error {
	return self.cmd.New(fmt.Sprintf("git clone %s %s", url, destination)).Run()
}
