package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type TagCommands struct {
	*common.Common

	cmd oscommands.ICmdObjBuilder
}

func NewTagCommands(common *common.Common, cmd oscommands.ICmdObjBuilder) *TagCommands {
	return &TagCommands{
		Common: common,
		cmd:    cmd,
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
	return self.cmd.New(fmt.Sprintf("git push %s %s", self.cmd.Quote(remoteName), self.cmd.Quote(tagName))).PromptOnCredentialRequest().Run()
}
