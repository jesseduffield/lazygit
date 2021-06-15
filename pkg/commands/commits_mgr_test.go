package commands_test

import (
	"github.com/go-errors/errors"
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommitsMgr", func() {
	var (
		commander  *FakeICommander
		gitconfig  *FakeIGitConfig
		commitsMgr *CommitsMgr
	)

	BeforeEach(func() {
		commander = NewFakeCommander()
		gitconfig = &FakeIGitConfig{}
		gitconfig.ColorArgCalls(func() string { return "always" })
		commitsMgr = NewCommitsMgr(commander, gitconfig)
	})

	Describe("RewordHead", func() {
		It("runs expected command", func() {
			commander.RunCalls(func(cmdObj ICmdObj) error {
				Expect(cmdObj.ToString()).To(Equal("git commit --allow-empty --amend --only -m \"newName\""))

				return nil
			})

			commitsMgr.RewordHead("newName")
		})
	})

	DescribeTable("CommitCmdObj",
		func(message string, flags string, expected string) {
			commander.BuildGitCmdObjFromStrCalls(func(command string) ICmdObj {
				Expect(command).To(Equal(expected))

				return nil
			})
		},
		Entry(
			"with message",
			"my message",
			"",
			"commit -m \"my message\"",
		),
		Entry(
			"with additional flags",
			"my message",
			"--flag",
			"commit --flag -m \"my message\"",
		),
		Entry(
			"with multiline message",
			"line one\nline two",
			"--flag",
			"commit --flag -m \"line one\" -m \"line two\"",
		),
	)

	Describe("GetHeadMessage", func() {
		It("runs expected command and trims output", func() {
			commander.RunWithOutputCalls(func(obj ICmdObj) (string, error) {
				Expect(obj.ToString()).To(Equal("git log -1 --pretty=%s"))

				return "blah blah\n", nil
			})

			message, err := commitsMgr.GetHeadMessage()

			Expect(message).To(Equal("blah blah"))
			Expect(err).To(BeNil())
		})

		It("returns error if one occurs", func() {
			commander.RunWithOutputReturns("", errors.New("my error"))

			message, err := commitsMgr.GetHeadMessage()

			Expect(message).To(Equal(""))
			Expect(err).To(MatchError("my error"))
		})
	})

	Describe("GetMessageFirstLine", func() {
		It("returns first line", func() {
			commander.RunWithOutputCalls(func(obj ICmdObj) (string, error) {
				Expect(obj.ToString()).To(Equal("git show --no-patch --pretty=format:%s abc123"))

				return "firstline", nil
			})

			message, err := commitsMgr.GetMessageFirstLine("abc123")

			Expect(message).To(Equal("firstline"))
			Expect(err).To(BeNil())
		})

		It("bubbles up error", func() {
			commander.RunWithOutputReturns("", errors.New("my error"))

			message, err := commitsMgr.GetMessageFirstLine("abc123")

			Expect(message).To(Equal(""))
			Expect(err).To(MatchError("my error"))
		})
	})

	Describe("AmendHead", func() {
		It("runs command", func() {
			commander.RunCalls(func(obj ICmdObj) error {
				Expect(obj.ToString()).To(Equal("git commit --amend --no-edit --allow-empty"))

				return nil
			})

			err := commitsMgr.AmendHead()

			Expect(err).To(BeNil())
		})
	})

	Describe("AmendHeadCmdObj", func() {
		It("returns command object", func() {
			obj := commitsMgr.AmendHeadCmdObj()
			Expect(obj.ToString()).To(Equal("git commit --amend --no-edit --allow-empty"))
		})
	})

	Describe("ShowCmdObj", func() {
		It("returns command object", func() {
			obj := commitsMgr.ShowCmdObj("abc123", "path")
			Expect(obj.ToString()).To(Equal("git show --submodule --color=always --no-renames --stat -p abc123 -- \"path\""))
		})

		It("handles lack of a path", func() {
			obj := commitsMgr.ShowCmdObj("abc123", "")
			Expect(obj.ToString()).To(Equal("git show --submodule --color=always --no-renames --stat -p abc123"))
		})
	})

	Describe("Revert", func() {
		It("runs command", func() {
			commander.RunCalls(func(cmdObj ICmdObj) error {
				Expect(cmdObj.ToString()).To(Equal("git revert abc123"))

				return nil
			})

			err := commitsMgr.Revert("abc123")
			Expect(err).To(BeNil())
		})
	})

	Describe("RevertMerge", func() {
		It("runs command", func() {
			commander.RunCalls(func(cmdObj ICmdObj) error {
				Expect(cmdObj.ToString()).To(Equal("git revert abc123 -m 1"))

				return nil
			})

			err := commitsMgr.RevertMerge("abc123", 1)
			Expect(err).To(BeNil())
		})
	})

	Describe("CreateFixupCommit", func() {
		It("runs command", func() {
			commander.RunCalls(func(cmdObj ICmdObj) error {
				Expect(cmdObj.ToString()).To(Equal("git commit --fixup=abc123"))

				return nil
			})

			err := commitsMgr.CreateFixupCommit("abc123")
			Expect(err).To(BeNil())
		})
	})
})
