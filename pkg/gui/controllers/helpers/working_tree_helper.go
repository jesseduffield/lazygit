package helpers

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type WorkingTreeHelper struct {
	c             *HelperCommon
	refHelper     *RefsHelper
	commitsHelper *CommitsHelper
	gpgHelper     *GpgHelper
}

func NewWorkingTreeHelper(
	c *HelperCommon,
	refHelper *RefsHelper,
	commitsHelper *CommitsHelper,
	gpgHelper *GpgHelper,
) *WorkingTreeHelper {
	return &WorkingTreeHelper{
		c:             c,
		refHelper:     refHelper,
		commitsHelper: commitsHelper,
		gpgHelper:     gpgHelper,
	}
}

func (self *WorkingTreeHelper) AnyStagedFiles() bool {
	return AnyStagedFiles(self.c.Model().Files)
}

func AnyStagedFiles(files []*models.File) bool {
	return lo.SomeBy(files, func(f *models.File) bool { return f.HasStagedChanges })
}

func (self *WorkingTreeHelper) AnyTrackedFiles() bool {
	return AnyTrackedFiles(self.c.Model().Files)
}

func AnyTrackedFiles(files []*models.File) bool {
	return lo.SomeBy(files, func(f *models.File) bool { return f.Tracked })
}

func (self *WorkingTreeHelper) IsWorkingTreeDirty() bool {
	return IsWorkingTreeDirty(self.c.Model().Files)
}

func IsWorkingTreeDirty(files []*models.File) bool {
	return AnyStagedFiles(files) || AnyTrackedFiles(files)
}

func (self *WorkingTreeHelper) FileForSubmodule(submodule *models.SubmoduleConfig) *models.File {
	for _, file := range self.c.Model().Files {
		if file.IsSubmodule([]*models.SubmoduleConfig{submodule}) {
			return file
		}
	}

	return nil
}

func (self *WorkingTreeHelper) OpenMergeTool() error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.MergeToolTitle,
		Prompt: self.c.Tr.MergeToolPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.OpenMergeTool)
			return self.c.RunSubprocessAndRefresh(
				self.c.Git().WorkingTree.OpenMergeToolCmdObj(),
			)
		},
	})

	return nil
}

func (self *WorkingTreeHelper) HandleCommitPressWithMessage(initialMessage string, forceSkipHooks bool) error {
	return self.WithEnsureCommittableFiles(func() error {
		self.commitsHelper.OpenCommitMessagePanel(
			&OpenCommitMessagePanelOpts{
				CommitIndex:      context.NoCommitIndex,
				InitialMessage:   initialMessage,
				SummaryTitle:     self.c.Tr.CommitSummaryTitle,
				DescriptionTitle: self.c.Tr.CommitDescriptionTitle,
				PreserveMessage:  true,
				OnConfirm: func(summary string, description string) error {
					return self.handleCommit(summary, description, forceSkipHooks)
				},
				OnSwitchToEditor: func(filepath string) error {
					return self.switchFromCommitMessagePanelToEditor(filepath, forceSkipHooks)
				},
				ForceSkipHooks:  forceSkipHooks,
				SkipHooksPrefix: self.c.UserConfig().Git.SkipHookPrefix,
			},
		)

		return nil
	})
}

func (self *WorkingTreeHelper) handleCommit(summary string, description string, forceSkipHooks bool) error {
	cmdObj := self.c.Git().Commit.CommitCmdObj(summary, description, forceSkipHooks)
	self.c.LogAction(self.c.Tr.Actions.Commit)
	return self.gpgHelper.WithGpgHandling(cmdObj, git_commands.CommitGpgSign, self.c.Tr.CommittingStatus,
		func() error {
			self.commitsHelper.OnCommitSuccess()
			return nil
		}, nil)
}

func (self *WorkingTreeHelper) switchFromCommitMessagePanelToEditor(filepath string, forceSkipHooks bool) error {
	// We won't be able to tell whether the commit was successful, because
	// RunSubprocessAndRefresh doesn't return the error (it opens an error alert
	// itself and returns nil on error). But even if we could, we wouldn't have
	// access to the last message that the user typed, and it might be very
	// different from what was last in the commit panel. So the best we can do
	// here is to always clear the remembered commit message.
	self.commitsHelper.OnCommitSuccess()

	self.c.LogAction(self.c.Tr.Actions.Commit)
	return self.c.RunSubprocessAndRefresh(
		self.c.Git().Commit.CommitInEditorWithMessageFileCmdObj(filepath, forceSkipHooks),
	)
}

// HandleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (self *WorkingTreeHelper) HandleCommitEditorPress() error {
	return self.WithEnsureCommittableFiles(func() error {
		self.c.LogAction(self.c.Tr.Actions.Commit)
		return self.c.RunSubprocessAndRefresh(
			self.c.Git().Commit.CommitEditorCmdObj(),
		)
	})
}

func (self *WorkingTreeHelper) HandleWIPCommitPress() error {
	var initialMessage string
	preservedMessage := self.c.Contexts().CommitMessage.GetPreservedMessageAndLogError()
	if preservedMessage == "" {
		// Use the skipHook prefix only if we don't have a preserved message
		initialMessage = self.c.UserConfig().Git.SkipHookPrefix
	}
	return self.HandleCommitPressWithMessage(initialMessage, true)
}

func (self *WorkingTreeHelper) HandleCommitPress() error {
	message := self.c.Contexts().CommitMessage.GetPreservedMessageAndLogError()

	if message == "" {
		commitPrefixConfigs := self.commitPrefixConfigsForRepo()
		for _, commitPrefixConfig := range commitPrefixConfigs {
			prefixPattern := commitPrefixConfig.Pattern
			if prefixPattern == "" {
				continue
			}
			prefixReplace := commitPrefixConfig.Replace
			branchName := self.refHelper.GetCheckedOutRef().Name
			rgx, err := regexp.Compile(prefixPattern)
			if err != nil {
				return fmt.Errorf("%s: %s", self.c.Tr.CommitPrefixPatternError, err.Error())
			}

			if rgx.MatchString(branchName) {
				prefix := rgx.ReplaceAllString(branchName, prefixReplace)
				message = prefix
				break
			}
		}
	}

	return self.HandleCommitPressWithMessage(message, false)
}

func (self *WorkingTreeHelper) WithEnsureCommittableFiles(handler func() error) error {
	if err := self.prepareFilesForCommit(); err != nil {
		return err
	}

	if len(self.c.Model().Files) == 0 {
		return errors.New(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.AnyStagedFiles() {
		return self.promptToStageAllAndRetry(handler)
	}

	return handler()
}

func (self *WorkingTreeHelper) promptToStageAllAndRetry(retry func() error) error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.NoFilesStagedTitle,
		Prompt: self.c.Tr.NoFilesStagedPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
			if err := self.c.Git().WorkingTree.StageAll(); err != nil {
				return err
			}
			if err := self.syncRefresh(); err != nil {
				return err
			}

			return retry()
		},
	})

	return nil
}

// for when you need to refetch files before continuing an action. Runs synchronously.
func (self *WorkingTreeHelper) syncRefresh() error {
	return self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES}})
}

func (self *WorkingTreeHelper) prepareFilesForCommit() error {
	noStagedFiles := !self.AnyStagedFiles()
	if noStagedFiles && self.c.UserConfig().Gui.SkipNoStagedFilesWarning {
		self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
		err := self.c.Git().WorkingTree.StageAll()
		if err != nil {
			return err
		}

		return self.syncRefresh()
	}

	return nil
}

func (self *WorkingTreeHelper) commitPrefixConfigsForRepo() []config.CommitPrefixConfig {
	cfg, ok := self.c.UserConfig().Git.CommitPrefixes[self.c.Git().RepoPaths.RepoName()]
	if ok {
		return append(cfg, self.c.UserConfig().Git.CommitPrefix...)
	} else {
		return self.c.UserConfig().Git.CommitPrefix
	}
}
