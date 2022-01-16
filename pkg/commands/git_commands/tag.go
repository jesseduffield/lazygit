package git_commands

import (
	"fmt"
)

type TagCommands struct {
	*GitCommon
}

func NewTagCommands(gitCommon *GitCommon) *TagCommands {
	return &TagCommands{
		GitCommon: gitCommon,
	}
}

func (self *TagCommands) CreateLightweight(tagName string, commitSha string) error {
	return self.cmd.New(fmt.Sprintf("git tag -- %s %s", self.cmd.Quote(tagName), commitSha)).Run()
}

func (self *TagCommands) CreateAnnotated(tagName, commitSha, msg string) error {
	return self.cmd.New(fmt.Sprintf("git tag %s %s -m %s", tagName, commitSha, self.cmd.Quote(msg))).Run()
}

func (self *TagCommands) Delete(tagName string) error {
	return self.cmd.New(fmt.Sprintf("git tag -d %s", self.cmd.Quote(tagName))).Run()
}

func (self *TagCommands) Push(remoteName string, tagName string) error {
	return self.cmd.New(fmt.Sprintf("git push %s %s", self.cmd.Quote(remoteName), self.cmd.Quote(tagName))).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}
