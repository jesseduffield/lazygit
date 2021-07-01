package commands_test

import (
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommitsLoader", func() {
	var (
		commander     *FakeICommander
		CommitsLoader *CommitsLoader
	)

	BeforeEach(func() {
		commander = NewFakeCommander()

		mgrCtx := NewFakeMgrCtx(commander, nil, nil)

		statusMgr := NewStatusMgr(mgrCtx)

		CommitsLoader = NewCommitsLoader(mgrCtx, statusMgr)
	})

	Describe("Load", func() {
		It("returns commits", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				{
					cmdStr:    "git merge-base HEAD HEAD@{u}",
					outputStr: "2d42f9649d4326c7adb78997b8e67b3435d60114",
					outputErr: nil,
				},
				{
					cmdStr:    "git symbolic-ref --short HEAD",
					outputStr: "current-branch",
					outputErr: nil,
				},
				{
					cmdStr:    "git merge-base HEAD master",
					outputStr: "9fdf92b226032d39503dbf40ef931d5d017b4235",
					outputErr: nil,
				},
			}, func() {
				commits, err := CommitsLoader.Load(LoadCommitsOptions{RefName: "HEAD"})
				Expect(err).To(BeNil())
				Expect(commits).To(Equal(
					[]*models.Commit{},
				))
			})
		})
	})
})
