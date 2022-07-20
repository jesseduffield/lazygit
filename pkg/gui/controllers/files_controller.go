package controllers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type FilesController struct {
	baseController
	*controllerCommon

	enterSubmodule        func(submodule *models.SubmoduleConfig) error
	setCommitMessage      func(message string)
	getSavedCommitMessage func() string
	switchToMergeFn       func(path string) error
}

var _ types.IController = &FilesController{}

func NewFilesController(
	common *controllerCommon,
	enterSubmodule func(submodule *models.SubmoduleConfig) error,
	setCommitMessage func(message string),
	getSavedCommitMessage func() string,
	switchToMergeFn func(path string) error,
) *FilesController {
	return &FilesController{
		controllerCommon:      common,
		enterSubmodule:        enterSubmodule,
		setCommitMessage:      setCommitMessage,
		getSavedCommitMessage: getSavedCommitMessage,
		switchToMergeFn:       switchToMergeFn,
	}
}

func (self *FilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelectedFileNode(self.press),
			Description: self.c.Tr.LcToggleStaged,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.OpenStatusFilter),
			Handler:     self.handleStatusFilterPressed,
			Description: self.c.Tr.LcFileFilter,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.CommitChanges),
			Handler:     self.HandleCommitPress,
			Description: self.c.Tr.CommitChanges,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.CommitChangesWithoutHook),
			Handler:     self.HandleWIPCommitPress,
			Description: self.c.Tr.LcCommitChangesWithoutHook,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.AmendLastCommit),
			Handler:     self.handleAmendCommitPress,
			Description: self.c.Tr.AmendLastCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.CommitChangesWithEditor),
			Handler:     self.HandleCommitEditorPress,
			Description: self.c.Tr.CommitChangesWithEditor,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelectedFileNode(self.edit),
			Description: self.c.Tr.LcEditFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.Open,
			Description: self.c.Tr.LcOpenFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.IgnoreOrExcludeFile),
			Handler:     self.checkSelectedFileNode(self.ignoreOrExcludeMenu),
			Description: self.c.Tr.Actions.IgnoreExcludeFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.RefreshFiles),
			Handler:     self.refresh,
			Description: self.c.Tr.LcRefreshFiles,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.StashAllChanges),
			Handler:     self.stash,
			Description: self.c.Tr.LcStashAllChanges,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ViewStashOptions),
			Handler:     self.createStashMenu,
			Description: self.c.Tr.LcViewStashOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleStagedAll),
			Handler:     self.stageAll,
			Description: self.c.Tr.LcToggleStagedAll,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.enter,
			Description: self.c.Tr.FileEnter,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.createResetToUpstreamMenu,
			Description: self.c.Tr.LcViewResetToUpstreamOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ViewResetOptions),
			Handler:     self.createResetMenu,
			Description: self.c.Tr.LcViewResetOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.LcToggleTreeView,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.OpenMergeTool),
			Handler:     self.helpers.WorkingTree.OpenMergeTool,
			Description: self.c.Tr.LcOpenMergeTool,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.Fetch),
			Handler:     self.fetch,
			Description: self.c.Tr.LcFetch,
		},
	}
}

func (self *FilesController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    "main",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FromContext: string(self.context().GetKey()),
		},
		{
			ViewName:    "secondary",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickSecondary,
			FromContext: string(self.context().GetKey()),
		},
	}
}

func (self *FilesController) GetOnClick() func() error {
	return self.checkSelectedFileNode(self.press)
}

func (self *FilesController) press(node *filetree.FileNode) error {
	if node.IsLeaf() {
		file := node.File

		if file.HasInlineMergeConflicts {
			return self.c.PushContext(self.contexts.Merging)
		}

		if file.HasUnstagedChanges {
			self.c.LogAction(self.c.Tr.Actions.StageFile)
			if err := self.git.WorkingTree.StageFile(file.Name); err != nil {
				return self.c.Error(err)
			}
		} else {
			self.c.LogAction(self.c.Tr.Actions.UnstageFile)
			if err := self.git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return self.c.Error(err)
			}
		}
	} else {
		// if any files within have inline merge conflicts we can't stage or unstage,
		// or it'll end up with those >>>>>> lines actually staged
		if node.GetHasInlineMergeConflicts() {
			return self.c.ErrorMsg(self.c.Tr.ErrStageDirWithInlineMergeConflicts)
		}

		if node.GetHasUnstagedChanges() {
			self.c.LogAction(self.c.Tr.Actions.StageFile)
			if err := self.git.WorkingTree.StageFile(node.Path); err != nil {
				return self.c.Error(err)
			}
		} else {
			// pretty sure it doesn't matter that we're always passing true here
			self.c.LogAction(self.c.Tr.Actions.UnstageFile)
			if err := self.git.WorkingTree.UnStageFile([]string{node.Path}, true); err != nil {
				return self.c.Error(err)
			}
		}
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}}); err != nil {
		return err
	}

	return self.context().HandleFocus()
}

func (self *FilesController) checkSelectedFileNode(callback func(*filetree.FileNode) error) func() error {
	return func() error {
		node := self.context().GetSelected()
		if node == nil {
			return nil
		}

		return callback(node)
	}
}

func (self *FilesController) Context() types.Context {
	return self.context()
}

func (self *FilesController) context() *context.WorkingTreeContext {
	return self.contexts.Files
}

func (self *FilesController) getSelectedFile() *models.File {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}
	return node.File
}

func (self *FilesController) enter() error {
	return self.EnterFile(types.OnFocusOpts{ClickedViewName: "", ClickedViewLineIdx: -1})
}

func (self *FilesController) EnterFile(opts types.OnFocusOpts) error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return self.handleToggleDirCollapsed()
	}

	file := node.File

	submoduleConfigs := self.model.Submodules
	if file.IsSubmodule(submoduleConfigs) {
		submoduleConfig := file.SubmoduleConfig(submoduleConfigs)
		return self.enterSubmodule(submoduleConfig)
	}

	if file.HasInlineMergeConflicts {
		return self.switchToMerge()
	}
	if file.HasMergeConflicts {
		return self.c.ErrorMsg(self.c.Tr.FileStagingRequirements)
	}

	return self.c.PushContext(self.contexts.Staging, opts)
}

func (self *FilesController) allFilesStaged() bool {
	for _, file := range self.model.Files {
		if file.HasUnstagedChanges {
			return false
		}
	}
	return true
}

func (self *FilesController) stageAll() error {
	var err error
	if self.allFilesStaged() {
		self.c.LogAction(self.c.Tr.Actions.UnstageAllFiles)
		err = self.git.WorkingTree.UnstageAll()
	} else {
		self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
		err = self.git.WorkingTree.StageAll()
	}
	if err != nil {
		_ = self.c.Error(err)
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}}); err != nil {
		return err
	}

	return self.contexts.Files.HandleFocus()
}

func (self *FilesController) unstageFiles(node *filetree.FileNode) error {
	return node.ForEachFile(func(file *models.File) error {
		if file.HasStagedChanges {
			if err := self.git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return err
			}
		}

		return nil
	})
}

func (self *FilesController) ignoreOrExcludeTracked(node *filetree.FileNode, trAction string, f func(string) error) error {
	self.c.LogAction(trAction)
	// not 100% sure if this is necessary but I'll assume it is
	if err := self.unstageFiles(node); err != nil {
		return err
	}

	if err := self.git.WorkingTree.RemoveTrackedFiles(node.GetPath()); err != nil {
		return err
	}

	if err := f(node.GetPath()); err != nil {
		return err
	}

	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
}

func (self *FilesController) ignoreOrExcludeUntracked(node *filetree.FileNode, trAction string, f func(string) error) error {
	self.c.LogAction(trAction)

	if err := f(node.GetPath()); err != nil {
		return err
	}

	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
}

func (self *FilesController) ignoreOrExcludeFile(node *filetree.FileNode, trText string, trPrompt string, trAction string, f func(string) error) error {
	if node.GetIsTracked() {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  trText,
			Prompt: trPrompt,
			HandleConfirm: func() error {
				return self.ignoreOrExcludeTracked(node, trAction, f)
			},
		})
	}
	return self.ignoreOrExcludeUntracked(node, trAction, f)
}

func (self *FilesController) ignore(node *filetree.FileNode) error {
	if node.GetPath() == ".gitignore" {
		return self.c.ErrorMsg(self.c.Tr.Actions.IgnoreFileErr)
	}
	err := self.ignoreOrExcludeFile(node, self.c.Tr.IgnoreTracked, self.c.Tr.IgnoreTrackedPrompt, self.c.Tr.Actions.IgnoreExcludeFile, self.git.WorkingTree.Ignore)
	if err != nil {
		return err
	}

	return nil
}

func (self *FilesController) exclude(node *filetree.FileNode) error {
	if node.GetPath() == ".git/info/exclude" {
		return self.c.ErrorMsg(self.c.Tr.Actions.ExcludeFileErr)
	}

	if node.GetPath() == ".gitignore" {
		return self.c.ErrorMsg(self.c.Tr.Actions.ExcludeGitIgnoreErr)
	}

	err := self.ignoreOrExcludeFile(node, self.c.Tr.ExcludeTracked, self.c.Tr.ExcludeTrackedPrompt, self.c.Tr.Actions.ExcludeFile, self.git.WorkingTree.Exclude)
	if err != nil {
		return err
	}
	return nil
}

func (self *FilesController) ignoreOrExcludeMenu(node *filetree.FileNode) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.IgnoreExcludeFile,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.LcIgnoreFile},
				OnPress: func() error {
					if err := self.ignore(node); err != nil {
						return self.c.Error(err)
					}
					return nil
				},
				Key: 'i',
			},
			{
				LabelColumns: []string{self.c.Tr.LcExcludeFile},
				OnPress: func() error {
					if err := self.exclude(node); err != nil {
						return self.c.Error(err)
					}
					return nil
				},
				Key: 'e',
			},
		},
	})
}

func (self *FilesController) HandleWIPCommitPress() error {
	skipHookPrefix := self.c.UserConfig.Git.SkipHookPrefix
	if skipHookPrefix == "" {
		return self.c.ErrorMsg(self.c.Tr.SkipHookPrefixNotConfigured)
	}

	self.setCommitMessage(skipHookPrefix)

	return self.HandleCommitPress()
}

func (self *FilesController) commitPrefixConfigForRepo() *config.CommitPrefixConfig {
	cfg, ok := self.c.UserConfig.Git.CommitPrefixes[utils.GetCurrentRepoName()]
	if !ok {
		return nil
	}

	return &cfg
}

func (self *FilesController) prepareFilesForCommit() error {
	noStagedFiles := !self.helpers.WorkingTree.AnyStagedFiles()
	if noStagedFiles && self.c.UserConfig.Gui.SkipNoStagedFilesWarning {
		self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
		err := self.git.WorkingTree.StageAll()
		if err != nil {
			return err
		}

		return self.syncRefresh()
	}

	return nil
}

// for when you need to refetch files before continuing an action. Runs synchronously.
func (self *FilesController) syncRefresh() error {
	return self.c.Refresh(types.RefreshOptions{Mode: types.SYNC, Scope: []types.RefreshableView{types.FILES}})
}

func (self *FilesController) refresh() error {
	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
}

func (self *FilesController) HandleCommitPress() error {
	if err := self.prepareFilesForCommit(); err != nil {
		return self.c.Error(err)
	}

	if len(self.model.Files) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.helpers.WorkingTree.AnyStagedFiles() {
		return self.promptToStageAllAndRetry(self.HandleCommitPress)
	}

	savedCommitMessage := self.getSavedCommitMessage()
	if len(savedCommitMessage) > 0 {
		self.setCommitMessage(savedCommitMessage)
	} else {
		commitPrefixConfig := self.commitPrefixConfigForRepo()
		if commitPrefixConfig != nil {
			prefixPattern := commitPrefixConfig.Pattern
			prefixReplace := commitPrefixConfig.Replace
			rgx, err := regexp.Compile(prefixPattern)
			if err != nil {
				return self.c.ErrorMsg(fmt.Sprintf("%s: %s", self.c.Tr.LcCommitPrefixPatternError, err.Error()))
			}
			prefix := rgx.ReplaceAllString(self.helpers.Refs.GetCheckedOutRef().Name, prefixReplace)
			self.setCommitMessage(prefix)
		}
	}

	if err := self.c.PushContext(self.contexts.CommitMessage); err != nil {
		return err
	}

	return nil
}

func (self *FilesController) promptToStageAllAndRetry(retry func() error) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.NoFilesStagedTitle,
		Prompt: self.c.Tr.NoFilesStagedPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.StageAllFiles)
			if err := self.git.WorkingTree.StageAll(); err != nil {
				return self.c.Error(err)
			}
			if err := self.syncRefresh(); err != nil {
				return self.c.Error(err)
			}

			return retry()
		},
	})
}

func (self *FilesController) handleAmendCommitPress() error {
	if len(self.model.Files) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.helpers.WorkingTree.AnyStagedFiles() {
		return self.promptToStageAllAndRetry(self.handleAmendCommitPress)
	}

	if len(self.model.Commits) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoCommitToAmend)
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AmendLastCommitTitle,
		Prompt: self.c.Tr.SureToAmend,
		HandleConfirm: func() error {
			cmdObj := self.git.Commit.AmendHeadCmdObj()
			self.c.LogAction(self.c.Tr.Actions.AmendCommit)
			return self.helpers.GPG.WithGpgHandling(cmdObj, self.c.Tr.AmendingStatus, nil)
		},
	})
}

// HandleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (self *FilesController) HandleCommitEditorPress() error {
	if len(self.model.Files) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.helpers.WorkingTree.AnyStagedFiles() {
		return self.promptToStageAllAndRetry(self.HandleCommitEditorPress)
	}

	self.c.LogAction(self.c.Tr.Actions.Commit)
	return self.c.RunSubprocessAndRefresh(
		self.git.Commit.CommitEditorCmdObj(),
	)
}

func (self *FilesController) handleStatusFilterPressed() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.FilteringMenuTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.FilterStagedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayStaged)
				},
			},
			{
				Label: self.c.Tr.FilterUnstagedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayUnstaged)
				},
			},
			{
				Label: self.c.Tr.ResetCommitFilterState,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayAll)
				},
			},
		},
	})
}

func (self *FilesController) setStatusFiltering(filter filetree.FileTreeDisplayFilter) error {
	self.context().FileTreeViewModel.SetFilter(filter)
	return self.c.PostRefreshUpdate(self.context())
}

func (self *FilesController) edit(node *filetree.FileNode) error {
	if node.File == nil {
		return self.c.ErrorMsg(self.c.Tr.ErrCannotEditDirectory)
	}

	return self.helpers.Files.EditFile(node.GetPath())
}

func (self *FilesController) Open() error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	return self.helpers.Files.OpenFile(node.GetPath())
}

func (self *FilesController) switchToMerge() error {
	file := self.getSelectedFile()
	if file == nil {
		return nil
	}

	return self.switchToMergeFn(file.Name)
}

func (self *FilesController) createStashMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.LcStashOptions,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.LcStashAllChanges,
				OnPress: func() error {
					return self.handleStashSave(self.git.Stash.Save, self.c.Tr.Actions.StashAllChanges, self.c.Tr.NoFilesToStash)
				},
				Key: 'a',
			},
			{
				Label: self.c.Tr.LcStashAllChangesKeepIndex,
				OnPress: func() error {
					// if there are no staged files it behaves the same as Stash.Save
					return self.handleStashSave(self.git.Stash.StashAndKeepIndex, self.c.Tr.Actions.StashAllChangesKeepIndex, self.c.Tr.NoFilesToStash)
				},
				Key: 'i',
			},
			{
				Label: self.c.Tr.LcStashStagedChanges,
				OnPress: func() error {
					// there must be something in staging otherwise the current implementation mucks the stash up
					if !self.helpers.WorkingTree.AnyStagedFiles() {
						return self.c.ErrorMsg(self.c.Tr.NoTrackedStagedFilesStash)
					}
					return self.handleStashSave(self.git.Stash.SaveStagedChanges, self.c.Tr.Actions.StashStagedChanges, self.c.Tr.NoTrackedStagedFilesStash)
				},
				Key: 's',
			},
			{
				Label: self.c.Tr.LcStashUnstagedChanges,
				OnPress: func() error {
					if self.helpers.WorkingTree.AnyStagedFiles() {
						return self.handleStashSave(self.git.Stash.StashUnstagedChanges, self.c.Tr.Actions.StashUnstagedChanges, self.c.Tr.NoFilesToStash)
					}
					// ordinary stash
					return self.handleStashSave(self.git.Stash.Save, self.c.Tr.Actions.StashUnstagedChanges, self.c.Tr.NoFilesToStash)
				},
				Key: 'u',
			},
		},
	})
}

func (self *FilesController) stash() error {
	return self.handleStashSave(self.git.Stash.Save, self.c.Tr.Actions.StashAllChanges, self.c.Tr.NoTrackedStagedFilesStash)
}

func (self *FilesController) createResetToUpstreamMenu() error {
	return self.helpers.Refs.CreateGitResetMenu("@{upstream}")
}

func (self *FilesController) handleToggleDirCollapsed() error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	self.context().FileTreeViewModel.ToggleCollapsed(node.GetPath())

	if err := self.c.PostRefreshUpdate(self.contexts.Files); err != nil {
		self.c.Log.Error(err)
	}

	return nil
}

func (self *FilesController) toggleTreeView() error {
	self.context().FileTreeViewModel.ToggleShowTree()

	return self.c.PostRefreshUpdate(self.context())
}

func (self *FilesController) handleStashSave(stashFunc func(message string) error, action string, errorMsg string) error {
	if !self.helpers.WorkingTree.IsWorkingTreeDirty() {
		return self.c.ErrorMsg(errorMsg)
	}

	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.StashChanges,
		HandleConfirm: func(stashComment string) error {
			self.c.LogAction(action)

			if err := stashFunc(stashComment); err != nil {
				return self.c.Error(err)
			}
			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
		},
	})
}

func (self *FilesController) onClickMain(opts gocui.ViewMouseBindingOpts) error {
	return self.EnterFile(types.OnFocusOpts{ClickedViewName: "main", ClickedViewLineIdx: opts.Y})
}

func (self *FilesController) onClickSecondary(opts gocui.ViewMouseBindingOpts) error {
	return self.EnterFile(types.OnFocusOpts{ClickedViewName: "secondary", ClickedViewLineIdx: opts.Y})
}

func (self *FilesController) fetch() error {
	return self.c.WithLoaderPanel(self.c.Tr.FetchWait, func() error {
		if err := self.fetchAux(); err != nil {
			_ = self.c.Error(err)
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *FilesController) fetchAux() (err error) {
	self.c.LogAction("Fetch")
	err = self.git.Sync.Fetch(git_commands.FetchOptions{})

	if err != nil && strings.Contains(err.Error(), "exit status 128") {
		_ = self.c.ErrorMsg(self.c.Tr.PassUnameWrong)
	}

	_ = self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS}, Mode: types.ASYNC})

	return err
}
