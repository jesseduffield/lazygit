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

func (self *TagCommands) CreateLightweight(tagName string, ref string) error {
	if len(ref) > 0 {
		return self.cmd.New(fmt.Sprintf("git tag -- %s %s", self.cmd.Quote(tagName), self.cmd.Quote(ref))).Run()
	} else {
		return self.cmd.New(fmt.Sprintf("git tag -- %s", self.cmd.Quote(tagName))).Run()
	}
}

func (self *TagCommands) CreateAnnotated(tagName, ref, msg string) error {
	if len(ref) > 0 {
		return self.cmd.New(fmt.Sprintf("git tag %s %s -m %s", self.cmd.Quote(tagName), self.cmd.Quote(ref), self.cmd.Quote(msg))).Run()
	} else {
		return self.cmd.New(fmt.Sprintf("git tag %s -m %s", self.cmd.Quote(tagName), self.cmd.Quote(msg))).Run()
	}
}

func (self *TagCommands) Delete(tagName string) error {
	return self.cmd.New(fmt.Sprintf("git tag -d %s", self.cmd.Quote(tagName))).Run()
}

func (self *TagCommands) Push(remoteName string, tagName string) error {
	return self.cmd.New(fmt.Sprintf("git push %s tag %s", self.cmd.Quote(remoteName), self.cmd.Quote(tagName))).PromptOnCredentialRequest().WithMutex(self.syncMutex).Run()
}
