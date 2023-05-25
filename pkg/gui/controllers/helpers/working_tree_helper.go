package helpers

import (
	"fmt"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	if err := self.prepareFilesForCommit(); err != nil {
		return self.c.Error(err)
	}

	if len(self.c.Model().Files) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.AnyStagedFiles() {
		return self.PromptToStageAllAndRetry(self.HandleCommitPress)
	}

	return self.commitsHelper.OpenCommitMessagePanel(
		&OpenCommitMessagePanelOpts{
			CommitIndex:     context.NoCommitIndex,
			InitialMessage:  initialMessage,
			Title:           self.c.Tr.CommitSummary,
			PreserveMessage: true,
			OnConfirm:       self.handleCommit,
		},
	)
}

func (self *WorkingTreeHelper) handleCommit(message string) error {
	cmdObj := self.c.Git().Commit.CommitCmdObj(message)
	self.c.LogAction(self.c.Tr.Actions.Commit)
	_ = self.commitsHelper.PopCommitMessageContexts()
	return self.gpgHelper.WithGpgHandling(cmdObj, self.c.Tr.CommittingStatus, func() error {
		self.commitsHelper.OnCommitSuccess()
		return nil
	})
}

// HandleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (self *WorkingTreeHelper) HandleCommitEditorPress() error {
	if len(self.c.Model().Files) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.AnyStagedFiles() {
		return self.PromptToStageAllAndRetry(self.HandleCommitEditorPress)
	}

	self.c.LogAction(self.c.Tr.Actions.Commit)
	return self.c.RunSubprocessAndRefresh(
		self.c.Git().Commit.CommitEditorCmdObj(),
	)
}

func (self *WorkingTreeHelper) HandleWIPCommitPress() error {
	skipHookPrefix := self.c.UserConfig.Git.SkipHookPrefix
	if skipHookPrefix == "" {
		return self.c.ErrorMsg(self.c.Tr.SkipHookPrefixNotConfigured)
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
			rgx, err := regexp.Compile(prefixPattern)
			if err != nil {
				return self.c.ErrorMsg(fmt.Sprintf("%s: %s", self.c.Tr.CommitPrefixPatternError, err.Error()))
			}
			prefix := rgx.ReplaceAllString(self.refHelper.GetCheckedOutRef().Name, prefixReplace)
			message = prefix
		}
	}

	return self.HandleCommitPressWithMessage(message)
}

func (self *WorkingTreeHelper) PromptToStageAllAndRetry(retry func() error) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.NoFilesStagedTitle,
		Prompt: self.c.Tr.NoFilesStagedPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
			if err := self.c.Git().WorkingTree.StageAll(); err != nil {
				return self.c.Error(err)
			}
			if err := self.syncRefresh(); err != nil {
				return self.c.Error(err)
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
	cfg, ok := self.c.UserConfig.Git.CommitPrefixes[utils.GetCurrentRepoName()]
	if !ok {
		return nil
	}

	return &cfg
}
