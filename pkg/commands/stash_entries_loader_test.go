package commands_test

import (
	"github.com/go-errors/errors"
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StashEntriesLoader", func() {
	var (
		commander          *FakeICommander
		StashEntriesLoader *StashEntriesLoader
	)

	BeforeEach(func() {
		commander = NewFakeCommander()

		mgrCtx := NewFakeMgrCtx(commander, nil, nil)

		StashEntriesLoader = NewStashEntriesLoader(mgrCtx)
	})

	Describe("Load", func() {
		Context("not filtering by path", func() {
			Context("stash entries exist", func() {
				It("returns stash entries", func() {
					WithRunCalls(commander, []ExpectedRunCall{
						{
							cmdStr:    "git stash list --pretty='%gs'",
							outputStr: "WIP on mybranch: 55c6af2 foo\nWIP on master: bb86a3f bar",
							outputErr: nil,
						},
					}, func() {
						stashEntries := StashEntriesLoader.Load("")
						Expect(stashEntries).To(Equal(
							[]*models.StashEntry{
								{
									Index: 0,
									Name:  "WIP on mybranch: 55c6af2 foo",
								},
								{
									Index: 1,
									Name:  "WIP on master: bb86a3f bar",
								},
							},
						))
					})
				})
			})

			Context("no stash entries exist", func() {
				It("returns empty array", func() {
					WithRunCalls(commander, []ExpectedRunCall{
						{
							cmdStr:    "git stash list --pretty='%gs'",
							outputStr: "\n",
							outputErr: nil,
						},
					}, func() {
						stashEntries := StashEntriesLoader.Load("")
						Expect(stashEntries).To(Equal([]*models.StashEntry{}))
					})
				})
			})

			Context("error is raised by command", func() {
				It("returns empty array", func() {
					WithRunCalls(commander, []ExpectedRunCall{
						{
							cmdStr:    "git stash list --pretty='%gs'",
							outputStr: "",
							outputErr: errors.New("my error"),
						},
					}, func() {
						stashEntries := StashEntriesLoader.Load("")
						Expect(stashEntries).To(Equal([]*models.StashEntry{}))
					})
				})
			})
		})

		Context("filtering by path", func() {
			var output = `stash@{0}: On mybranch: foo

pkg/commands/loaders/files.go
stash@{1}: On patch-1: extras title

pkg/i18n/english.go
scripts/generate_cheatsheet.go
stash@{2}: On otherbranch: mocking

pkg/commands/interface.go
pkg/commands/loaders/files.go
pkg/gui/handlers/sync/push_files/mocks/Gui.go`

			Context("stash entries exist", func() {
				It("returns stash entries", func() {
					WithRunCalls(commander, []ExpectedRunCall{
						{
							cmdStr:    "git stash list --name-only",
							outputStr: output,
							outputErr: nil,
						},
					}, func() {
						stashEntries := StashEntriesLoader.Load("pkg/commands/loaders/files.go")
						Expect(stashEntries).To(Equal(
							[]*models.StashEntry{
								{
									Index: 0,
									Name:  "On mybranch: foo",
								},
								{
									Index: 2,
									Name:  "On otherbranch: mocking",
								},
							},
						))
					})
				})
			})

			Context("no stash entries exist", func() {
				It("returns empty array", func() {
					WithRunCalls(commander, []ExpectedRunCall{
						{
							cmdStr:    "git stash list --name-only",
							outputStr: "\n",
							outputErr: nil,
						},
					}, func() {
						stashEntries := StashEntriesLoader.Load("pkg/commands/loaders/files.go")
						Expect(stashEntries).To(Equal([]*models.StashEntry{}))
					})
				})
			})

			Context("error is raised by command", func() {
				// not sure if we should actually do this
				It("falls back to unfiltered search", func() {
					WithRunCalls(commander, []ExpectedRunCall{
						{
							cmdStr:    "git stash list --name-only",
							outputStr: "",
							outputErr: errors.New("my error"),
						},
						{
							cmdStr:    "git stash list --pretty='%gs'",
							outputStr: "WIP on mybranch: 55c6af2 foo\nWIP on master: bb86a3f bar",
							outputErr: nil,
						},
					}, func() {
						stashEntries := StashEntriesLoader.Load("pkg/commands/loaders/files.go")
						Expect(stashEntries).To(Equal([]*models.StashEntry{{
							Index: 0,
							Name:  "WIP on mybranch: 55c6af2 foo",
						},
							{
								Index: 1,
								Name:  "WIP on master: bb86a3f bar",
							},
						}))
					})
				})
			})
		})
	})
})
