package commands_test

import (
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BranchesMgr", func() {
	var (
		commander   *FakeICommander
		gitconfig   *FakeIGitConfig
		branchesMgr *BranchesMgr
	)

	BeforeEach(func() {
		commander = NewFakeCommander()
		gitconfig = &FakeIGitConfig{}
		gitconfig.ColorArgCalls(func() string { return "always" })
		branchesMgr = NewBranchesMgr(commander, gitconfig)
	})

	Describe("NewBranch", func() {
		It("runs expected command", func() {
			commander.RunCalls(func(cmdObj ICmdObj) error {
				Expect(cmdObj.ToString()).To(Equal("git checkout -b newName master"))

				return nil
			})

			branchesMgr.NewBranch("newName", "master")
		})
	})
})
