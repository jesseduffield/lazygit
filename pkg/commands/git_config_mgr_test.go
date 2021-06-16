package commands_test

import (
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("GitConfigMgr", func() {
	var (
		commander         *FakeICommander
		userConfig        *config.UserConfig
		gitconfig         *GitConfigMgr
		getGitConfigValue func(string) (string, error)
	)

	BeforeEach(func() {
		commander = NewFakeCommander()
		userConfig = &config.UserConfig{}
		getGitConfigValue = func(string) (string, error) { return "", nil }
		gitconfig = NewGitConfigMgr(commander, userConfig, ".", getGitConfigValue, nil)
	})

	Describe("UsingGpg", func() {
		Context("user has overridden GPG in their config", func() {
			BeforeEach(func() {
				userConfig.Git.OverrideGpg = true
			})

			It("returns false", func() {
				Expect(gitconfig.UsingGpg()).To(BeFalse())
			})
		})

		Context("user has not overridden GPG in their config", func() {
			BeforeEach(func() {
				userConfig.Git.OverrideGpg = false
			})

			DescribeTable("CommitCmdObj",
				func(result string, expected bool) {
					getGitConfigValue = func(input string) (string, error) {
						Expect(input).To(Equal("commit.gpgsign"))
						return result, nil
					}

					gitconfig = NewGitConfigMgr(commander, userConfig, ".", getGitConfigValue, nil)
					Expect(gitconfig.UsingGpg()).To(Equal(expected))
				},
				Entry("when returning true", "true", true),
				Entry("when returning 1", "1", true),
				Entry("when returning yes", "yes", true),
				Entry("when returning True (capitalised)", "True", true),
				Entry("when returning false", "false", false),
			)
		})
	})
})
