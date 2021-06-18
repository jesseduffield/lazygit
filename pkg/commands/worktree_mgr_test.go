package commands_test

import (
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands/oscommandsfakes"
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorktreeMgr", func() {
	var (
		describedStruct *WorktreeMgr

		commander     *FakeICommander
		gitconfig     *FakeIGitConfigMgr
		userConfig    *config.UserConfig
		branchesMgr   *BranchesMgr
		submodulesMgr *SubmodulesMgr
		mgrCtx        *MgrCtx
		os            *oscommandsfakes.FakeIOS
	)

	BeforeEach(func() {
		commander = NewFakeCommander()
		gitconfig = &FakeIGitConfigMgr{}
		userConfig = &config.UserConfig{}
		gitconfig.GetUserConfigCalls(func() *config.UserConfig { return userConfig })
		gitconfig.ColorArgCalls(func() string { return "always" })

		os = &oscommandsfakes.FakeIOS{}
		mgrCtx = NewFakeMgrCtx(commander, gitconfig, os)

		statusMgr := NewStatusMgr(mgrCtx)
		branchesMgr = NewBranchesMgr(mgrCtx, statusMgr)
		submodulesMgr = NewSubmodulesMgr(mgrCtx)

		describedStruct = NewWorktreeMgr(mgrCtx, branchesMgr, submodulesMgr)
	})

	Describe("StageAll", func() {
		It("runs expected command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git add -A"),
			}, func() {
				err := describedStruct.StageAll()
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("UnstageAll", func() {
		It("runs expected command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git reset"),
			}, func() {
				err := describedStruct.UnstageAll()
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("StageFile", func() {
		It("runs expected command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git add -- \"myfile.go\""),
			}, func() {
				err := describedStruct.StageFile("myfile.go")
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("UnstageFile", func() {
		Context("with reset: false", func() {
			It("runs expected command", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall("git rm --cached --force -- \"myfile.go\""),
				}, func() {
					err := describedStruct.UnStageFile([]string{"myfile.go"}, false)
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with reset: true", func() {
			It("runs expected command", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall("git reset HEAD -- \"myfile.go\""),
				}, func() {
					err := describedStruct.UnStageFile([]string{"myfile.go"}, true)
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with multiple file names passed", func() {
			It("runs expected commands", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall("git reset HEAD -- \"myfile.go\""),
					SuccessCall("git reset HEAD -- \"myfile2.go\""),
				}, func() {
					err := describedStruct.UnStageFile([]string{"myfile.go", "myfile2.go"}, true)
					Expect(err).To(BeNil())
				})
			})
		})
	})

	Describe("DiscardAllFileChanges", func() {
		Context("with file with staged changes", func() {
			It("runs expected commands", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall(`git reset -- "myfile.go"`),
					SuccessCall(`git checkout -- "myfile.go"`),
				}, func() {
					err := describedStruct.DiscardAllFileChanges(&models.File{
						Name: "myfile.go", HasStagedChanges: true,
					})
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with file of 'AA' status", func() {
			It("runs expected commands", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall(`git checkout --ours --  "myfile.go"`),
					SuccessCall(`git add -- "myfile.go"`),
				}, func() {

					err := describedStruct.DiscardAllFileChanges(&models.File{
						Name: "myfile.go", HasStagedChanges: true, ShortStatus: "AA",
						HasMergeConflicts: true,
					})
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with file of 'DU' status", func() {
			It("runs expected commands", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall(`git rm "myfile.go"`),
				}, func() {
					err := describedStruct.DiscardAllFileChanges(&models.File{
						Name: "myfile.go", ShortStatus: "DU", HasMergeConflicts: true,
					})
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with file with merge conflicts", func() {
			It("runs expected commands", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall(`git reset -- "myfile.go"`),
					SuccessCall(`git checkout -- "myfile.go"`),
				}, func() {

					err := describedStruct.DiscardAllFileChanges(&models.File{
						Name: "myfile.go", HasMergeConflicts: true,
					})
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with file with 'DD' status", func() {
			It("runs expected commands", func() {
				WithRunCalls(commander, []ExpectedRunCall{
					SuccessCall(`git reset -- "myfile.go"`),
				}, func() {
					err := describedStruct.DiscardAllFileChanges(&models.File{
						Name: "myfile.go", ShortStatus: "DD", HasMergeConflicts: true,
					})
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with added file", func() {
			It("runs expected commands", func() {
				WithRunCalls(commander, []ExpectedRunCall{}, func() {
					os.RemoveFileCalls(func(s string) error {
						Expect(s).To(Equal("myfile.go"))
						return nil
					})

					err := describedStruct.DiscardAllFileChanges(&models.File{
						Name: "myfile.go", Added: true,
					})

					Expect(err).To(BeNil())
					Expect(os.RemoveFileCallCount()).To(Equal(1))
				})
			})
		})
	})
})
