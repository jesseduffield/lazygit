package helpers

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type IWorkingTreeHelper interface {
	AnyStagedFiles() bool
	AnyTrackedFiles() bool
	IsWorkingTreeDirty() bool
	FileForSubmodule(submodule *models.SubmoduleConfig) *models.File
}

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
	for _, file := range self.c.Model().Files {
		if file.HasStagedChanges {
			return true
		}
	}
	return false
}

func (self *WorkingTreeHelper) AnyTrackedFiles() bool {
	for _, file := range self.c.Model().Files {
		if file.Tracked {
			return true
		}
	}
	return false
}

func (self *WorkingTreeHelper) IsWorkingTreeDirty() bool {
	return self.AnyStagedFiles() || self.AnyTrackedFiles()
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
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.MergeToolTitle,
		Prompt: self.c.Tr.MergeToolPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.OpenMergeTool)
			return self.c.RunSubprocessAndRefresh(
				self.c.Git().WorkingTree.OpenMergeToolCmdObj(),
			)
		},
	})
}

func (self *WorkingTreeHelper) HandleCommitPressWithMessage(initialMessage string) error {
	return self.WithEnsureCommitableFiles(func() error {
		return self.commitsHelper.OpenCommitMessagePanel(
			&OpenCommitMessagePanelOpts{
				CommitIndex:      context.NoCommitIndex,
				InitialMessage:   initialMessage,
				SummaryTitle:     self.c.Tr.CommitSummaryTitle,
				DescriptionTitle: self.c.Tr.CommitDescriptionTitle,
				PreserveMessage:  true,
				OnConfirm:        self.handleCommit,
				OnSwitchToEditor: self.switchFromCommitMessagePanelToEditor,
			},
		)
	})
}

func (self *WorkingTreeHelper) handleCommit(summary string, description string) error {
	cmdObj := self.c.Git().Commit.CommitCmdObj(summary, description)
	self.c.LogAction(self.c.Tr.Actions.Commit)
	return self.gpgHelper.WithGpgHandling(cmdObj, self.c.Tr.CommittingStatus, func() error {
		self.commitsHelper.OnCommitSuccess()
		return nil
	})
}

func (self *WorkingTreeHelper) switchFromCommitMessagePanelToEditor(filepath string) error {
	// We won't be able to tell whether the commit was successful, because
	// RunSubprocessAndRefresh doesn't return the error (it opens an error alert
	// itself and returns nil on error). But even if we could, we wouldn't have
	// access to the last message that the user typed, and it might be very
	// different from what was last in the commit panel. So the best we can do
	// here is to always clear the remembered commit message.
	self.commitsHelper.OnCommitSuccess()

	self.c.LogAction(self.c.Tr.Actions.Commit)
	return self.c.RunSubprocessAndRefresh(
		self.c.Git().Commit.CommitInEditorWithMessageFileCmdObj(filepath),
	)
}

// HandleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (self *WorkingTreeHelper) HandleCommitEditorPress() error {
	return self.WithEnsureCommitableFiles(func() error {
		self.c.LogAction(self.c.Tr.Actions.Commit)
		return self.c.RunSubprocessAndRefresh(
			self.c.Git().Commit.CommitEditorCmdObj(),
		)
	})
}

func (self *WorkingTreeHelper) HandleWIPCommitPress() error {
	skipHookPrefix := self.c.UserConfig.Git.SkipHookPrefix
	if skipHookPrefix == "" {
		return errors.New(self.c.Tr.SkipHookPrefixNotConfigured)
	}

	return self.HandleCommitPressWithMessage(skipHookPrefix)
}

func (self *WorkingTreeHelper) HandleCommitPress() error {
	message := self.c.Contexts().CommitMessage.GetPreservedMessage()

	if message == "" {
		commitPrefixConfig := self.commitPrefixConfigForRepo()
		if commitPrefixConfig != nil {
			prefixPattern := commitPrefixConfig.Pattern
			prefixReplace := commitPrefixConfig.Replace
			branchName := self.refHelper.GetCheckedOutRef().Name
			rgx, err := regexp.Compile(prefixPattern)
			if err != nil {
				return fmt.Errorf("%s: %s", self.c.Tr.CommitPrefixPatternError, err.Error())
			}

			if rgx.MatchString(branchName) {
				prefix := rgx.ReplaceAllString(branchName, prefixReplace)
				message = prefix
			}
		}
	}

	return self.HandleCommitPressWithMessage(message)
}

func (self *WorkingTreeHelper) WithEnsureCommitableFiles(handler func() error) error {
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
	return self.c.Confirm(types.ConfirmOpts{
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
}

// for when you need to refetch files before continuing an action. Runs synchronously.
func (self *WorkingTreeHelper) syncRefresh() error {
	return self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES}})
}

func (self *WorkingTreeHelper) prepareFilesForCommit() error {
	noStagedFiles := !self.AnyStagedFiles()
	if noStagedFiles && self.c.UserConfig.Gui.SkipNoStagedFilesWarning {
		self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
		err := self.c.Git().WorkingTree.StageAll()
		if err != nil {
			return err
		}

		return self.syncRefresh()
	}

	return nil
}

func (self *WorkingTreeHelper) commitPrefixConfigForRepo() *config.CommitPrefixConfig {
	cfg, ok := self.c.UserConfig.Git.CommitPrefixes[self.c.Git().RepoPaths.RepoName()]
	if ok {
		return &cfg
	}

	return self.c.UserConfig.Git.CommitPrefix
}
