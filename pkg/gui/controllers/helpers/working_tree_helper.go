package helpers

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type WorkingTreeHelper struct {
	c                    *HelperCommon
	refHelper            *RefsHelper
	commitsHelper        *CommitsHelper
	gpgHelper            *GpgHelper
	mergeAndRebaseHelper *MergeAndRebaseHelper
}

func NewWorkingTreeHelper(
	c *HelperCommon,
	refHelper *RefsHelper,
	commitsHelper *CommitsHelper,
	gpgHelper *GpgHelper,
	mergeAndRebaseHelper *MergeAndRebaseHelper,
) *WorkingTreeHelper {
	return &WorkingTreeHelper{
		c:                    c,
		refHelper:            refHelper,
		commitsHelper:        commitsHelper,
		gpgHelper:            gpgHelper,
		mergeAndRebaseHelper: mergeAndRebaseHelper,
	}
}

func (self *WorkingTreeHelper) AnyStagedFiles() bool {
	return AnyStagedFiles(self.c.Model().Files)
}

func AnyStagedFiles(files []*models.File) bool {
	return lo.SomeBy(files, func(f *models.File) bool { return f.HasStagedChanges })
}

func (self *WorkingTreeHelper) AnyStagedFilesExceptSubmodules() bool {
	return AnyStagedFilesExceptSubmodules(self.c.Model().Files, self.c.Model().Submodules)
}

func AnyStagedFilesExceptSubmodules(files []*models.File, submoduleConfigs []*models.SubmoduleConfig) bool {
	return lo.SomeBy(files, func(f *models.File) bool { return f.HasStagedChanges && !f.IsSubmodule(submoduleConfigs) })
}

func (self *WorkingTreeHelper) AnyTrackedFiles() bool {
	return AnyTrackedFiles(self.c.Model().Files)
}

func AnyTrackedFiles(files []*models.File) bool {
	return lo.SomeBy(files, func(f *models.File) bool { return f.Tracked })
}

func (self *WorkingTreeHelper) AnyTrackedFilesExceptSubmodules() bool {
	return AnyTrackedFilesExceptSubmodules(self.c.Model().Files, self.c.Model().Submodules)
}

func AnyTrackedFilesExceptSubmodules(files []*models.File, submoduleConfigs []*models.SubmoduleConfig) bool {
	return lo.SomeBy(files, func(f *models.File) bool { return f.Tracked && !f.IsSubmodule(submoduleConfigs) })
}

func (self *WorkingTreeHelper) IsWorkingTreeDirtyExceptSubmodules() bool {
	return IsWorkingTreeDirtyExceptSubmodules(self.c.Model().Files, self.c.Model().Submodules)
}

func IsWorkingTreeDirtyExceptSubmodules(files []*models.File, submoduleConfigs []*models.SubmoduleConfig) bool {
	return AnyStagedFilesExceptSubmodules(files, submoduleConfigs) || AnyTrackedFilesExceptSubmodules(files, submoduleConfigs)
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
	self.c.LogAction(self.c.Tr.Actions.OpenMergeTool)
	return self.c.RunSubprocessAndRefresh(
		self.c.Git().WorkingTree.OpenMergeToolCmdObj(),
	)
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
			self.commitsHelper.ClearPreservedCommitMessage()
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
	self.commitsHelper.ClearPreservedCommitMessage()

	self.c.LogAction(self.c.Tr.Actions.Commit)
	return self.c.RunSubprocessAndRefresh(
		self.c.Git().Commit.CommitInEditorWithMessageFileCmdObj(filepath, forceSkipHooks),
	)
}

// HandleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (self *WorkingTreeHelper) HandleCommitEditorPress() error {
	return self.WithEnsureCommittableFiles(func() error {
		// See reasoning in switchFromCommitMessagePanelToEditor for why it makes sense
		// to clear this message before calling into the editor
		self.commitsHelper.ClearPreservedCommitMessage()

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
			if err := self.c.Git().WorkingTree.StageAll(false); err != nil {
				return err
			}
			self.syncRefresh()

			return retry()
		},
	})

	return nil
}

// for when you need to refetch files before continuing an action. Runs synchronously.
func (self *WorkingTreeHelper) syncRefresh() {
	self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES}})
}

func (self *WorkingTreeHelper) prepareFilesForCommit() error {
	noStagedFiles := !self.AnyStagedFiles()
	if noStagedFiles && self.c.UserConfig().Gui.SkipNoStagedFilesWarning {
		self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
		err := self.c.Git().WorkingTree.StageAll(false)
		if err != nil {
			return err
		}

		self.syncRefresh()
	}

	return nil
}

func (self *WorkingTreeHelper) commitPrefixConfigsForRepo() []config.CommitPrefixConfig {
	cfg, ok := self.c.UserConfig().Git.CommitPrefixes[self.c.Git().RepoPaths.RepoName()]
	if ok {
		return append(cfg, self.c.UserConfig().Git.CommitPrefix...)
	}

	return self.c.UserConfig().Git.CommitPrefix
}

func (self *WorkingTreeHelper) mergeFile(filepath string, strategy string) (string, error) {
	if self.c.Git().Version.IsOlderThan(2, 43, 0) {
		return self.mergeFileWithTempFiles(filepath, strategy)
	}

	return self.mergeFileWithObjectIDs(filepath, strategy)
}

func (self *WorkingTreeHelper) mergeFileWithTempFiles(filepath string, strategy string) (string, error) {
	showToTempFile := func(stage int, label string) (string, error) {
		output, err := self.c.Git().WorkingTree.ShowFileAtStage(filepath, stage)
		if err != nil {
			return "", err
		}

		f, err := os.CreateTemp(self.c.GetConfig().GetTempDir(), "mergefile-"+label+"-*")
		if err != nil {
			return "", err
		}
		defer f.Close()

		if _, err := f.Write([]byte(output)); err != nil {
			return "", err
		}

		return f.Name(), nil
	}

	baseFilepath, err := showToTempFile(1, "base")
	if err != nil {
		return "", err
	}
	defer os.Remove(baseFilepath)

	oursFilepath, err := showToTempFile(2, "ours")
	if err != nil {
		return "", err
	}
	defer os.Remove(oursFilepath)

	theirsFilepath, err := showToTempFile(3, "theirs")
	if err != nil {
		return "", err
	}
	defer os.Remove(theirsFilepath)

	return self.c.Git().WorkingTree.MergeFileForFiles(strategy, oursFilepath, baseFilepath, theirsFilepath)
}

func (self *WorkingTreeHelper) mergeFileWithObjectIDs(filepath, strategy string) (string, error) {
	baseID, err := self.c.Git().WorkingTree.ObjectIDAtStage(filepath, 1)
	if err != nil {
		return "", err
	}

	oursID, err := self.c.Git().WorkingTree.ObjectIDAtStage(filepath, 2)
	if err != nil {
		return "", err
	}

	theirsID, err := self.c.Git().WorkingTree.ObjectIDAtStage(filepath, 3)
	if err != nil {
		return "", err
	}

	return self.c.Git().WorkingTree.MergeFileForObjectIDs(strategy, oursID, baseID, theirsID)
}

func (self *WorkingTreeHelper) CreateMergeConflictMenu(selectedFilepaths []string) error {
	onMergeStrategySelected := func(strategy string) error {
		for _, filepath := range selectedFilepaths {
			output, err := self.mergeFile(filepath, strategy)
			if err != nil {
				return err
			}

			if err = os.WriteFile(filepath, []byte(output), 0o644); err != nil {
				return err
			}
		}

		err := self.c.Git().WorkingTree.StageFiles(selectedFilepaths, nil)
		self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES}})
		return err
	}

	cmdColor := style.FgBlue
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.MergeConflictOptionsTitle,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{
					self.c.Tr.UseCurrentChanges,
					cmdColor.Sprint("git merge-file --ours"),
				},
				OnPress: func() error {
					return onMergeStrategySelected("--ours")
				},
				Key: 'c',
			},
			{
				LabelColumns: []string{
					self.c.Tr.UseIncomingChanges,
					cmdColor.Sprint("git merge-file --theirs"),
				},
				OnPress: func() error {
					return onMergeStrategySelected("--theirs")
				},
				Key: 'i',
			},
			{
				LabelColumns: []string{
					self.c.Tr.UseBothChanges,
					cmdColor.Sprint("git merge-file --union"),
				},
				OnPress: func() error {
					return onMergeStrategySelected("--union")
				},
				Key: 'b',
			},
			{
				LabelColumns: []string{
					self.c.Tr.OpenMergeTool,
					cmdColor.Sprint("git mergetool"),
				},
				OnPress: self.OpenMergeTool,
				Key:     'm',
			},
		},
	})
}
