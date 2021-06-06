package push_files_test

import (
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/gui/handlers/sync/push_files"
	. "github.com/jesseduffield/lazygit/pkg/gui/handlers/sync/push_files/push_filesfakes"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var tr = i18n.EnglishTranslationSet()

var _ = Describe("PushFiles", func() {
	var (
		gui     *FakeGui
		handler *PushFilesHandler
	)

	BeforeEach(func() {
		gui = &FakeGui{}
		handler = &PushFilesHandler{
			Gui: gui,
		}
	})

	Context("When able to push unforcefully", func() {
		It("should invoke a regular push", func() {
			gui.GetTrStub = func() *i18n.TranslationSet { return &tr }
			gui.PopupPanelFocusedReturns(false)
			gui.CurrentBranchReturns(&models.Branch{Pushables: "0", Pullables: "0", Name: "mybranch"})
			gui.WithPopupWaitingStatusStub = func(message string, f func() error) error {
				Expect(message).To(Equal("Pushing..."))
				return f()
			}

			testGitCommandObj := &FakeIGitCommand{}

			gui.GetGitCommandReturns(testGitCommandObj)

			testGitCommandObj.WithSpanReturns(testGitCommandObj)
			testGitCommandObj.PushStub = func(
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

			err := handler.Run()
			Expect(err).To(BeNil())
		})
	})
})
