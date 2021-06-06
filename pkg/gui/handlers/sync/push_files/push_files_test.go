package push_files_test

import (
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	. "github.com/jesseduffield/lazygit/pkg/gui/types"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/gui/handlers/sync/push_files"
	. "github.com/jesseduffield/lazygit/pkg/gui/handlers/sync/push_files/push_filesfakes"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var tr = i18n.EnglishTranslationSet()

var stubRegularPush = func(
	branchName string,
	force bool, upstream,
	args string,
	promptUserForCredential func(string) string,
) error {
	Expect(branchName).To(Equal("mybranch"))
	Expect(force).To(BeFalse())
	Expect(upstream).To(BeEmpty())
	Expect(args).To(BeEmpty())
	return nil
}

var stubForcePush = func(
	branchName string,
	force bool, upstream,
	args string,
	promptUserForCredential func(string) string,
) error {
	Expect(branchName).To(Equal("mybranch"))
	Expect(force).To(BeTrue())
	Expect(upstream).To(BeEmpty())
	Expect(args).To(BeEmpty())
	return nil
}

var _ = Describe("PushFiles", func() {
	var (
		gui                      *FakeGui
		handler                  *PushFilesHandler
		gitCommand               *FakeIGitCommand
		forcePushEnabledInConfig bool
	)

	BeforeEach(func() {
		gui = &FakeGui{}
		gitCommand = &FakeIGitCommand{}
	})

	JustBeforeEach(func() {
		userConfig := &config.UserConfig{}
		userConfig.Git.DisableForcePushing = !forcePushEnabledInConfig
		gui.GetUserConfigReturns(userConfig)
		gui.GetGitCommandReturns(gitCommand)
		gitCommand.WithSpanReturns(gitCommand)
		gui.GetTrStub = func() *i18n.TranslationSet { return &tr }
		gui.PopupPanelFocusedReturns(false)
		gui.WithPopupWaitingStatusStub = func(message string, f func() error) error {
			Expect(message).To(Equal("Pushing..."))
			return f()
		}

		handler = &PushFilesHandler{
			Gui: gui,
		}
	})

	Context("When branch is not tracking a remote", func() {
		JustBeforeEach(func() {
			gui.CurrentBranchReturns(
				&models.Branch{Pushables: "?", Pullables: "?", Name: "mybranch"},
			)
		})

		Context("When branch is tracking a remote in the git config", func() {
			JustBeforeEach(func() {
				gitCommand.FindRemoteForBranchInConfigStub = func(branchName string) (string, error) {
					Expect(branchName).To(Equal("mybranch"))
					return "remoteName", nil
				}
			})

			It("should invoke a push to the upstream", func() {
				gitCommand.PushStub = func(
					branchName string,
					force bool, upstream,
					args string,
					promptUserForCredential func(string) string,
				) error {
					Expect(branchName).To(Equal("mybranch"))
					Expect(force).To(BeFalse())
					Expect(upstream).To(BeEmpty())
					Expect(args).To(Equal("remoteName mybranch"))
					return nil
				}

				err := handler.Run()
				Expect(err).To(BeNil())

				Expect(gitCommand.PushCallCount()).To(Equal(1))
			})
		})

		Context("When branch is not tracking a remote in the git config", func() {
			var pushToCurrent bool

			JustBeforeEach(func() {
				gitCommand.FindRemoteForBranchInConfigStub = func(branchName string) (string, error) {
					Expect(branchName).To(Equal("mybranch"))
					return "", nil
				}

				gitCommand.GetPushToCurrentReturns(pushToCurrent)
			})

			Context("When push-to-current is configured", func() {
				BeforeEach(func() {
					pushToCurrent = true
				})

				It("should invoke a push to the upstream", func() {
					gitCommand.PushStub = func(
						branchName string,
						force bool, upstream,
						args string,
						promptUserForCredential func(string) string,
					) error {
						Expect(branchName).To(Equal("mybranch"))
						Expect(force).To(BeFalse())
						Expect(upstream).To(BeEmpty())
						Expect(args).To(Equal("--set-upstream"))
						return nil
					}

					err := handler.Run()
					Expect(err).To(BeNil())

					Expect(gitCommand.PushCallCount()).To(Equal(1))
				})
			})

			Context("When push-to-current is not configured", func() {
				BeforeEach(func() {
					pushToCurrent = false
				})

				It("should prompt to set upstream and then push", func() {
					gitCommand.PushStub = func(
						branchName string,
						force bool, upstream,
						args string,
						promptUserForCredential func(string) string,
					) error {
						Expect(branchName).To(Equal("mybranch"))
						Expect(force).To(BeFalse())
						Expect(upstream).To(Equal("origin mybranch"))
						// TODO: see if we should be passing --set-upstream here
						Expect(args).To(Equal(""))
						return nil
					}

					gui.PromptStub = func(opts PromptOpts) error {
						Expect(opts.Title).To(Equal("Enter upstream as '<remote> <branchname>'"))
						// pressing enter without modifying the content
						return opts.HandleConfirm(opts.InitialContent)
					}

					err := handler.Run()
					Expect(err).To(BeNil())

					Expect(gitCommand.PushCallCount()).To(Equal(1))
					Expect(gui.PromptCallCount()).To(Equal(1))
				})
			})
		})
	})

	Context("When branch is tracking a remote", func() {
		Context("When branch has no commits to pull", func() {
			JustBeforeEach(func() {
				gui.CurrentBranchReturns(
					&models.Branch{Pushables: "0", Pullables: "0", Name: "mybranch"},
				)
			})

			It("should invoke a regular push", func() {
				gitCommand.PushStub = stubRegularPush

				err := handler.Run()
				Expect(err).To(BeNil())

				Expect(gitCommand.PushCallCount()).To(Equal(1))
			})

			Context("When push fails and requires push", func() {
				Context("When force push is disabled in the user config", func() {

				})

				Context("When force push is enabled in the user config", func() {
					Context("When user confirms to force push", func() {

					})

					Context("When user does not confirm to force push", func() {

					})
				})
			})
		})

		Context("When branch has commits to pull", func() {
			JustBeforeEach(func() {
				gui.CurrentBranchReturns(
					&models.Branch{Pushables: "0", Pullables: "1", Name: "mybranch"},
				)
			})

			Context("When force pushing is disabled in the config", func() {
				BeforeEach(func() {
					forcePushEnabledInConfig = false
				})

				It("should display an error", func() {
					gui.CreateErrorPanelStub = func(message string) error {
						Expect(message).To(ContainSubstring("you've disabled force pushing"))
						return nil
					}

					err := handler.Run()
					Expect(err).To(BeNil())

					Expect(gitCommand.PushCallCount()).To(Equal(0))
				})
			})

			Context("When force pushing is enabled in the config", func() {
				BeforeEach(func() {
					forcePushEnabledInConfig = true
				})

				Context("When user does not confirm to push", func() {
					It("should not push at all", func() {
						gui.AskStub = func(opts AskOpts) error { return nil }

						err := handler.Run()
						Expect(err).To(BeNil())

						Expect(gitCommand.PushCallCount()).To(Equal(0))
					})
				})

				Context("When user does confirm to push", func() {
					It("should force push", func() {
						gitCommand.PushStub = stubForcePush

						gui.AskStub = func(opts AskOpts) error {
							return opts.HandleConfirm()
						}

						err := handler.Run()
						Expect(err).To(BeNil())

						Expect(gitCommand.PushCallCount()).To(Equal(1))
					})
				})
			})
		})
	})
})
