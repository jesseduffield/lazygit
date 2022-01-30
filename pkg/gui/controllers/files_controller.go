package controllers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type FilesController struct {
	// I've said publicly that I'm against single-letter variable names but in this
	// case I would actually prefer a _zero_ letter variable name in the form of
	// struct embedding, but Go does not allow hiding public fields in an embedded struct
	// to the client
	c          *types.ControllerCommon
	getContext func() *context.WorkingTreeContext
	getFiles   func() []*models.File
	git        *commands.GitCommand
	os         *oscommands.OSCommand

	getSelectedFileNode    func() *filetree.FileNode
	getContexts            func() context.ContextTree
	enterSubmodule         func(submodule *models.SubmoduleConfig) error
	getSubmodules          func() []*models.SubmoduleConfig
	setCommitMessage       func(message string)
	getCheckedOutBranch    func() *models.Branch
	withGpgHandling        func(cmdObj oscommands.ICmdObj, waitingStatus string, onSuccess func() error) error
	getFailedCommitMessage func() string
	getCommits             func() []*models.Commit
	getSelectedPath        func() string
	switchToMergeFn        func(path string) error
	suggestionsHelper      ISuggestionsHelper
	refsHelper             IRefsHelper
	filesHelper            IFileHelper
	workingTreeHelper      IWorkingTreeHelper
}

var _ types.IController = &FilesController{}

func NewFilesController(
	c *types.ControllerCommon,
	getContext func() *context.WorkingTreeContext,
	getFiles func() []*models.File,
	git *commands.GitCommand,
	os *oscommands.OSCommand,
	getSelectedFileNode func() *filetree.FileNode,
	allContexts func() context.ContextTree,
	enterSubmodule func(submodule *models.SubmoduleConfig) error,
	getSubmodules func() []*models.SubmoduleConfig,
	setCommitMessage func(message string),
	withGpgHandling func(cmdObj oscommands.ICmdObj, waitingStatus string, onSuccess func() error) error,
	getFailedCommitMessage func() string,
	getCommits func() []*models.Commit,
	getSelectedPath func() string,
	switchToMergeFn func(path string) error,
	suggestionsHelper ISuggestionsHelper,
	refsHelper IRefsHelper,
	filesHelper IFileHelper,
	workingTreeHelper IWorkingTreeHelper,
) *FilesController {
	return &FilesController{
		c:                      c,
		getContext:             getContext,
		getFiles:               getFiles,
		git:                    git,
		os:                     os,
		getSelectedFileNode:    getSelectedFileNode,
		getContexts:            allContexts,
		enterSubmodule:         enterSubmodule,
		getSubmodules:          getSubmodules,
		setCommitMessage:       setCommitMessage,
		withGpgHandling:        withGpgHandling,
		getFailedCommitMessage: getFailedCommitMessage,
		getCommits:             getCommits,
		getSelectedPath:        getSelectedPath,
		switchToMergeFn:        switchToMergeFn,
		suggestionsHelper:      suggestionsHelper,
		refsHelper:             refsHelper,
		filesHelper:            filesHelper,
		workingTreeHelper:      workingTreeHelper,
	}
}

func (self *FilesController) Keybindings(getKey func(key string) interface{}, config config.KeybindingConfig, guards types.KeybindingGuards) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         getKey(config.Universal.Select),
			Handler:     self.checkSelectedFileNode(self.press),
			Description: self.c.Tr.LcToggleStaged,
		},
		{
			Key:     gocui.MouseLeft,
			Handler: func() error { return self.getContext().HandleClick(self.checkSelectedFileNode(self.press)) },
		},
		{
			Key:         getKey("<c-b>"), // TODO: softcode
			Handler:     self.handleStatusFilterPressed,
			Description: self.c.Tr.LcFileFilter,
		},
		{
			Key:         getKey(config.Files.CommitChanges),
			Handler:     self.HandleCommitPress,
			Description: self.c.Tr.CommitChanges,
		},
		{
			Key:         getKey(config.Files.CommitChangesWithoutHook),
			Handler:     self.HandleWIPCommitPress,
			Description: self.c.Tr.LcCommitChangesWithoutHook,
		},
		{
			Key:         getKey(config.Files.AmendLastCommit),
			Handler:     self.handleAmendCommitPress,
			Description: self.c.Tr.AmendLastCommit,
		},
		{
			Key:         getKey(config.Files.CommitChangesWithEditor),
			Handler:     self.HandleCommitEditorPress,
			Description: self.c.Tr.CommitChangesWithEditor,
		},
		{
			Key:         getKey(config.Universal.Edit),
			Handler:     self.checkSelectedFileNode(self.edit),
			Description: self.c.Tr.LcEditFile,
		},
		{
			Key:         getKey(config.Universal.OpenFile),
			Handler:     self.Open,
			Description: self.c.Tr.LcOpenFile,
		},
		{
			Key:         getKey(config.Files.IgnoreFile),
			Handler:     self.checkSelectedFileNode(self.ignore),
			Description: self.c.Tr.LcIgnoreFile,
		},
		{
			Key:         getKey(config.Universal.Remove),
			Handler:     self.checkSelectedFileNode(self.remove),
			Description: self.c.Tr.LcViewDiscardOptions,
			OpensMenu:   true,
		},
		{
			Key:         getKey(config.Files.RefreshFiles),
			Handler:     self.refresh,
			Description: self.c.Tr.LcRefreshFiles,
		},
		{
			Key:         getKey(config.Files.StashAllChanges),
			Handler:     self.stash,
			Description: self.c.Tr.LcStashAllChanges,
		},
		{
			Key:         getKey(config.Files.ViewStashOptions),
			Handler:     self.createStashMenu,
			Description: self.c.Tr.LcViewStashOptions,
			OpensMenu:   true,
		},
		{
			Key:         getKey(config.Files.ToggleStagedAll),
			Handler:     self.stageAll,
			Description: self.c.Tr.LcToggleStagedAll,
		},
		{
			Key:         getKey(config.Universal.GoInto),
			Handler:     self.enter,
			Description: self.c.Tr.FileEnter,
		},
		{
			ViewName:    "",
			Key:         getKey(config.Universal.ExecuteCustomCommand),
			Handler:     self.handleCustomCommand,
			Description: self.c.Tr.LcExecuteCustomCommand,
		},
		{
			Key:         getKey(config.Commits.ViewResetOptions),
			Handler:     self.createResetMenu,
			Description: self.c.Tr.LcViewResetToUpstreamOptions,
			OpensMenu:   true,
		},
		// here
		{
			Key:         getKey(config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.LcToggleTreeView,
		},
		{
			Key:         getKey(config.Files.OpenMergeTool),
			Handler:     self.OpenMergeTool,
			Description: self.c.Tr.LcOpenMergeTool,
		},
	}

	return append(bindings, self.getContext().Keybindings(getKey, config, guards)...)
}

func (self *FilesController) press(node *filetree.FileNode) error {
	if node.IsLeaf() {
		file := node.File

		if file.HasInlineMergeConflicts {
			return self.c.PushContext(self.getContexts().Merging)
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

	return self.getContext().HandleFocus()
}

func (self *FilesController) checkSelectedFileNode(callback func(*filetree.FileNode) error) func() error {
	return func() error {
		node := self.getSelectedFileNode()
		if node == nil {
			return nil
		}

		return callback(node)
	}
}

func (self *FilesController) Context() types.Context {
	return self.getContext()
}

func (self *FilesController) getSelectedFile() *models.File {
	node := self.getSelectedFileNode()
	if node == nil {
		return nil
	}
	return node.File
}

func (self *FilesController) enter() error {
	return self.EnterFile(types.OnFocusOpts{ClickedViewName: "", ClickedViewLineIdx: -1})
}

func (self *FilesController) EnterFile(opts types.OnFocusOpts) error {
	node := self.getSelectedFileNode()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return self.handleToggleDirCollapsed()
	}

	file := node.File

	submoduleConfigs := self.getSubmodules()
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

	return self.c.PushContext(self.getContexts().Staging, opts)
}

func (self *FilesController) allFilesStaged() bool {
	for _, file := range self.getFiles() {
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

	return self.getContexts().Files.HandleFocus()
}

func (self *FilesController) ignore(node *filetree.FileNode) error {
	if node.GetPath() == ".gitignore" {
		return self.c.ErrorMsg("Cannot ignore .gitignore")
	}

	unstageFiles := func() error {
		return node.ForEachFile(func(file *models.File) error {
			if file.HasStagedChanges {
				if err := self.git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if node.GetIsTracked() {
		return self.c.Ask(types.AskOpts{
			Title:  self.c.Tr.IgnoreTracked,
			Prompt: self.c.Tr.IgnoreTrackedPrompt,
			HandleConfirm: func() error {
				self.c.LogAction(self.c.Tr.Actions.IgnoreFile)
				// not 100% sure if this is necessary but I'll assume it is
				if err := unstageFiles(); err != nil {
					return err
				}

				if err := self.git.WorkingTree.RemoveTrackedFiles(node.GetPath()); err != nil {
					return err
				}

				if err := self.git.WorkingTree.Ignore(node.GetPath()); err != nil {
					return err
				}
				return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
			},
		})
	}

	self.c.LogAction(self.c.Tr.Actions.IgnoreFile)

	if err := unstageFiles(); err != nil {
		return err
	}

	if err := self.git.WorkingTree.Ignore(node.GetPath()); err != nil {
		return self.c.Error(err)
	}

	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
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
	noStagedFiles := !self.workingTreeHelper.AnyStagedFiles()
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

	if len(self.getFiles()) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.workingTreeHelper.AnyStagedFiles() {
		return self.promptToStageAllAndRetry(self.HandleCommitPress)
	}

	failedCommitMessage := self.getFailedCommitMessage()
	if len(failedCommitMessage) > 0 {
		self.setCommitMessage(failedCommitMessage)
	} else {
		commitPrefixConfig := self.commitPrefixConfigForRepo()
		if commitPrefixConfig != nil {
			prefixPattern := commitPrefixConfig.Pattern
			prefixReplace := commitPrefixConfig.Replace
			rgx, err := regexp.Compile(prefixPattern)
			if err != nil {
				return self.c.ErrorMsg(fmt.Sprintf("%s: %s", self.c.Tr.LcCommitPrefixPatternError, err.Error()))
			}
			prefix := rgx.ReplaceAllString(self.getCheckedOutBranch().Name, prefixReplace)
			self.setCommitMessage(prefix)
		}
	}

	if err := self.c.PushContext(self.getContexts().CommitMessage); err != nil {
		return err
	}

	return nil
}

func (self *FilesController) promptToStageAllAndRetry(retry func() error) error {
	return self.c.Ask(types.AskOpts{
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
	if len(self.getFiles()) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.workingTreeHelper.AnyStagedFiles() {
		return self.promptToStageAllAndRetry(self.handleAmendCommitPress)
	}

	if len(self.getCommits()) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoCommitToAmend)
	}

	return self.c.Ask(types.AskOpts{
		Title:  strings.Title(self.c.Tr.AmendLastCommit),
		Prompt: self.c.Tr.SureToAmend,
		HandleConfirm: func() error {
			cmdObj := self.git.Commit.AmendHeadCmdObj()
			self.c.LogAction(self.c.Tr.Actions.AmendCommit)
			return self.withGpgHandling(cmdObj, self.c.Tr.AmendingStatus, nil)
		},
	})
}

// HandleCommitEditorPress - handle when the user wants to commit changes via
// their editor rather than via the popup panel
func (self *FilesController) HandleCommitEditorPress() error {
	if len(self.getFiles()) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.workingTreeHelper.AnyStagedFiles() {
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
				DisplayString: self.c.Tr.FilterStagedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayStaged)
				},
			},
			{
				DisplayString: self.c.Tr.FilterUnstagedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayUnstaged)
				},
			},
			{
				DisplayString: self.c.Tr.ResetCommitFilterState,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayAll)
				},
			},
		},
	})
}

func (self *FilesController) setStatusFiltering(filter filetree.FileTreeDisplayFilter) error {
	self.getContext().FileTreeViewModel.SetFilter(filter)
	return self.c.PostRefreshUpdate(self.getContext())
}

func (self *FilesController) edit(node *filetree.FileNode) error {
	if node.File == nil {
		return self.c.ErrorMsg(self.c.Tr.ErrCannotEditDirectory)
	}

	return self.filesHelper.EditFile(node.GetPath())
}

func (self *FilesController) Open() error {
	node := self.getSelectedFileNode()
	if node == nil {
		return nil
	}

	return self.filesHelper.OpenFile(node.GetPath())
}

func (self *FilesController) switchToMerge() error {
	file := self.getSelectedFile()
	if file == nil {
		return nil
	}

	return self.switchToMergeFn(file.Name)
}

func (self *FilesController) handleCustomCommand() error {
	return self.c.Prompt(types.PromptOpts{
		Title:               self.c.Tr.CustomCommand,
		FindSuggestionsFunc: self.suggestionsHelper.GetCustomCommandsHistorySuggestionsFunc(),
		HandleConfirm: func(command string) error {
			self.c.GetAppState().CustomCommandsHistory = utils.Limit(
				utils.Uniq(
					append(self.c.GetAppState().CustomCommandsHistory, command),
				),
				1000,
			)

			err := self.c.SaveAppState()
			if err != nil {
				self.c.Log.Error(err)
			}

			self.c.LogAction(self.c.Tr.Actions.CustomCommand)
			return self.c.RunSubprocessAndRefresh(
				self.os.Cmd.NewShell(command),
			)
		},
	})
}

func (self *FilesController) createStashMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.LcStashOptions,
		Items: []*types.MenuItem{
			{
				DisplayString: self.c.Tr.LcStashAllChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StashAllChanges)
					return self.handleStashSave(self.git.Stash.Save)
				},
			},
			{
				DisplayString: self.c.Tr.LcStashStagedChanges,
				OnPress: func() error {
					self.c.LogAction(self.c.Tr.Actions.StashStagedChanges)
					return self.handleStashSave(self.git.Stash.SaveStagedChanges)
				},
			},
		},
	})
}

func (self *FilesController) stash() error {
	return self.handleStashSave(self.git.Stash.Save)
}

func (self *FilesController) createResetMenu() error {
	return self.refsHelper.CreateGitResetMenu("@{upstream}")
}

func (self *FilesController) handleToggleDirCollapsed() error {
	node := self.getSelectedFileNode()
	if node == nil {
		return nil
	}

	self.getContext().FileTreeViewModel.ToggleCollapsed(node.GetPath())

	if err := self.c.PostRefreshUpdate(self.getContexts().Files); err != nil {
		self.c.Log.Error(err)
	}

	return nil
}

func (self *FilesController) toggleTreeView() error {
	// get path of currently selected file
	path := self.getSelectedPath()

	self.getContext().FileTreeViewModel.ToggleShowTree()

	// find that same node in the new format and move the cursor to it
	if path != "" {
		self.getContext().FileTreeViewModel.ExpandToPath(path)
		index, found := self.getContext().FileTreeViewModel.GetIndexForPath(path)
		if found {
			self.getContext().GetPanelState().SetSelectedLineIdx(index)
		}
	}

	return self.c.PostRefreshUpdate(self.getContext())
}

func (self *FilesController) OpenMergeTool() error {
	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.MergeToolTitle,
		Prompt: self.c.Tr.MergeToolPrompt,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.OpenMergeTool)
			return self.c.RunSubprocessAndRefresh(
				self.git.WorkingTree.OpenMergeToolCmdObj(),
			)
		},
	})
}

func (self *FilesController) ResetSubmodule(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.LcResettingSubmoduleStatus, func() error {
		self.c.LogAction(self.c.Tr.Actions.ResetSubmodule)

		file := self.workingTreeHelper.FileForSubmodule(submodule)
		if file != nil {
			if err := self.git.WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return self.c.Error(err)
			}
		}

		if err := self.git.Submodule.Stash(submodule); err != nil {
			return self.c.Error(err)
		}
		if err := self.git.Submodule.Reset(submodule); err != nil {
			return self.c.Error(err)
		}

		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.SUBMODULES}})
	})
}

func (self *FilesController) handleStashSave(stashFunc func(message string) error) error {
	if !self.workingTreeHelper.IsWorkingTreeDirty() {
		return self.c.ErrorMsg(self.c.Tr.NoTrackedStagedFilesStash)
	}

	return self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.StashChanges,
		HandleConfirm: func(stashComment string) error {
			if err := stashFunc(stashComment); err != nil {
				return self.c.Error(err)
			}
			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
		},
	})
}
