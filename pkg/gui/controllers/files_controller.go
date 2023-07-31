package controllers

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type FilesController struct {
	baseController // nolint: unused
	c              *ControllerCommon
}

var _ types.IController = &FilesController{}

func NewFilesController(
	common *ControllerCommon,
) *FilesController {
	return &FilesController{
		c: common,
	}
}

func (self *FilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelectedFileNode(self.press),
			Description: self.c.Tr.ToggleStaged,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.OpenStatusFilter),
			Handler:     self.handleStatusFilterPressed,
			Description: self.c.Tr.FileFilter,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.CommitChanges),
			Handler:     self.c.Helpers().WorkingTree.HandleCommitPress,
			Description: self.c.Tr.CommitChanges,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.CommitChangesWithoutHook),
			Handler:     self.c.Helpers().WorkingTree.HandleWIPCommitPress,
			Description: self.c.Tr.CommitChangesWithoutHook,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.AmendLastCommit),
			Handler:     self.handleAmendCommitPress,
			Description: self.c.Tr.AmendLastCommit,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.CommitChangesWithEditor),
			Handler:     self.c.Helpers().WorkingTree.HandleCommitEditorPress,
			Description: self.c.Tr.CommitChangesWithEditor,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelectedFileNode(self.edit),
			Description: self.c.Tr.EditFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.Open,
			Description: self.c.Tr.OpenFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.IgnoreFile),
			Handler:     self.checkSelectedFileNode(self.ignoreOrExcludeMenu),
			Description: self.c.Tr.Actions.IgnoreExcludeFile,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.RefreshFiles),
			Handler:     self.refresh,
			Description: self.c.Tr.RefreshFiles,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.StashAllChanges),
			Handler:     self.stash,
			Description: self.c.Tr.StashAllChanges,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ViewStashOptions),
			Handler:     self.createStashMenu,
			Description: self.c.Tr.ViewStashOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleStagedAll),
			Handler:     self.toggleStagedAll,
			Description: self.c.Tr.ToggleStagedAll,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.enter,
			Description: self.c.Tr.FileEnter,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.createResetToUpstreamMenu,
			Description: self.c.Tr.ViewResetToUpstreamOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ViewResetOptions),
			Handler:     self.createResetMenu,
			Description: self.c.Tr.ViewResetOptions,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.ToggleTreeView,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.OpenMergeTool),
			Handler:     self.c.Helpers().WorkingTree.OpenMergeTool,
			Description: self.c.Tr.OpenMergeTool,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.Fetch),
			Handler:     self.fetch,
			Description: self.c.Tr.Fetch,
		},
	}
}

func (self *FilesController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    "main",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FocusedView: self.context().GetViewName(),
		},
		{
			ViewName:    "patchBuilding",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FocusedView: self.context().GetViewName(),
		},
		{
			ViewName:    "mergeConflicts",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FocusedView: self.context().GetViewName(),
		},
		{
			ViewName:    "secondary",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickSecondary,
			FocusedView: self.context().GetViewName(),
		},
		{
			ViewName:    "patchBuildingSecondary",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickSecondary,
			FocusedView: self.context().GetViewName(),
		},
	}
}

func (self *FilesController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			node := self.context().GetSelected()

			if node == nil {
				return self.c.RenderToMainViews(types.RefreshMainOpts{
					Pair: self.c.MainViewPairs().Normal,
					Main: &types.ViewUpdateOpts{
						Title:    self.c.Tr.DiffTitle,
						SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
						Task:     types.NewRenderStringTask(self.c.Tr.NoChangedFiles),
					},
				})
			}

			if node.File != nil && node.File.HasInlineMergeConflicts {
				hasConflicts, err := self.c.Helpers().MergeConflicts.SetMergeState(node.GetPath())
				if err != nil {
					return err
				}

				if hasConflicts {
					return self.c.Helpers().MergeConflicts.Render(false)
				}
			}

			self.c.Helpers().MergeConflicts.ResetMergeState()

			pair := self.c.MainViewPairs().Normal
			if node.File != nil {
				pair = self.c.MainViewPairs().Staging
			}

			split := self.c.UserConfig.Gui.SplitDiff == "always" || (node.GetHasUnstagedChanges() && node.GetHasStagedChanges())
			mainShowsStaged := !split && node.GetHasStagedChanges()

			cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(node, false, mainShowsStaged, self.c.GetAppState().IgnoreWhitespaceInDiffView)
			title := self.c.Tr.UnstagedChanges
			if mainShowsStaged {
				title = self.c.Tr.StagedChanges
			}
			refreshOpts := types.RefreshMainOpts{
				Pair: pair,
				Main: &types.ViewUpdateOpts{
					Task:     types.NewRunPtyTask(cmdObj.GetCmd()),
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Title:    title,
				},
			}

			if split {
				cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(node, false, true, self.c.GetAppState().IgnoreWhitespaceInDiffView)

				title := self.c.Tr.StagedChanges
				if mainShowsStaged {
					title = self.c.Tr.UnstagedChanges
				}

				refreshOpts.Secondary = &types.ViewUpdateOpts{
					Title:    title,
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Task:     types.NewRunPtyTask(cmdObj.GetCmd()),
				}
			}

			return self.c.RenderToMainViews(refreshOpts)
		})
	}
}

func (self *FilesController) GetOnClick() func() error {
	return self.checkSelectedFileNode(self.press)
}

// if we are dealing with a status for which there is no key in this map,
// then we won't optimistically render: we'll just let `git status` tell
// us what the new status is.
// There are no doubt more entries that could be added to these two maps.
var stageStatusMap = map[string]string{
	"??": "A ",
	" M": "M ",
	"MM": "M ",
	" D": "D ",
	" A": "A ",
	"AM": "A ",
	"MD": "D ",
}

var unstageStatusMap = map[string]string{
	"A ": "??",
	"M ": " M",
	"D ": " D",
}

func (self *FilesController) optimisticStage(file *models.File) bool {
	newShortStatus, ok := stageStatusMap[file.ShortStatus]
	if !ok {
		return false
	}

	models.SetStatusFields(file, newShortStatus)
	return true
}

func (self *FilesController) optimisticUnstage(file *models.File) bool {
	newShortStatus, ok := unstageStatusMap[file.ShortStatus]
	if !ok {
		return false
	}

	models.SetStatusFields(file, newShortStatus)
	return true
}

// Running a git add command followed by a git status command can take some time (e.g. 200ms).
// Given how often users stage/unstage files in Lazygit, we're adding some
// optimistic rendering to make things feel faster. When we go to stage
// a file, we'll first update that file's status in-memory, then re-render
// the files panel. Then we'll immediately do a proper git status call
// so that if the optimistic rendering got something wrong, it's quickly
// corrected.
func (self *FilesController) optimisticChange(node *filetree.FileNode, optimisticChangeFn func(*models.File) bool) error {
	rerender := false
	err := node.ForEachFile(func(f *models.File) error {
		// can't act on the file itself: we need to update the original model file
		for _, modelFile := range self.c.Model().Files {
			if modelFile.Name == f.Name {
				if optimisticChangeFn(modelFile) {
					rerender = true
				}
				break
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	if rerender {
		if err := self.c.PostRefreshUpdate(self.c.Contexts().Files); err != nil {
			return err
		}
	}

	return nil
}

func (self *FilesController) pressWithLock(node *filetree.FileNode) error {
	// Obtaining this lock because optimistic rendering requires us to mutate
	// the files in our model.
	self.c.Mutexes().RefreshingFilesMutex.Lock()
	defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

	if node.IsFile() {
		file := node.File

		if file.HasUnstagedChanges {
			self.c.LogAction(self.c.Tr.Actions.StageFile)

			if err := self.optimisticChange(node, self.optimisticStage); err != nil {
				return err
			}

			if err := self.c.Git().WorkingTree.StageFile(file.Name); err != nil {
				return self.c.Error(err)
			}
		} else {
			self.c.LogAction(self.c.Tr.Actions.UnstageFile)

			if err := self.optimisticChange(node, self.optimisticUnstage); err != nil {
				return err
			}

			if err := self.c.Git().WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
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

			if err := self.optimisticChange(node, self.optimisticStage); err != nil {
				return err
			}

			if err := self.c.Git().WorkingTree.StageFile(node.Path); err != nil {
				return self.c.Error(err)
			}
		} else {
			self.c.LogAction(self.c.Tr.Actions.UnstageFile)

			if err := self.optimisticChange(node, self.optimisticUnstage); err != nil {
				return err
			}

			// pretty sure it doesn't matter that we're always passing true here
			if err := self.c.Git().WorkingTree.UnStageFile([]string{node.Path}, true); err != nil {
				return self.c.Error(err)
			}
		}
	}

	return nil
}

func (self *FilesController) press(node *filetree.FileNode) error {
	if node.IsFile() && node.File.HasInlineMergeConflicts {
		return self.switchToMerge()
	}

	if err := self.pressWithLock(node); err != nil {
		return err
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Mode: types.ASYNC}); err != nil {
		return err
	}

	return self.context().HandleFocus(types.OnFocusOpts{})
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
	return self.c.Contexts().Files
}

func (self *FilesController) getSelectedFile() *models.File {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}
	return node.File
}

func (self *FilesController) enter() error {
	return self.EnterFile(types.OnFocusOpts{ClickedWindowName: "", ClickedViewLineIdx: -1})
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

	submoduleConfigs := self.c.Model().Submodules
	if file.IsSubmodule(submoduleConfigs) {
		submoduleConfig := file.SubmoduleConfig(submoduleConfigs)
		return self.c.Helpers().Repos.EnterSubmodule(submoduleConfig)
	}

	if file.HasInlineMergeConflicts {
		return self.switchToMerge()
	}
	if file.HasMergeConflicts {
		return self.c.ErrorMsg(self.c.Tr.FileStagingRequirements)
	}

	return self.c.PushContext(self.c.Contexts().Staging, opts)
}

func (self *FilesController) toggleStagedAll() error {
	if err := self.toggleStagedAllWithLock(); err != nil {
		return err
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Mode: types.ASYNC}); err != nil {
		return err
	}

	return self.context().HandleFocus(types.OnFocusOpts{})
}

func (self *FilesController) toggleStagedAllWithLock() error {
	self.c.Mutexes().RefreshingFilesMutex.Lock()
	defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

	root := self.context().FileTreeViewModel.GetRoot()

	// if any files within have inline merge conflicts we can't stage or unstage,
	// or it'll end up with those >>>>>> lines actually staged
	if root.GetHasInlineMergeConflicts() {
		return self.c.ErrorMsg(self.c.Tr.ErrStageDirWithInlineMergeConflicts)
	}

	if root.GetHasUnstagedChanges() {
		self.c.LogAction(self.c.Tr.Actions.StageAllFiles)

		if err := self.optimisticChange(root, self.optimisticStage); err != nil {
			return err
		}

		if err := self.c.Git().WorkingTree.StageAll(); err != nil {
			return self.c.Error(err)
		}
	} else {
		self.c.LogAction(self.c.Tr.Actions.UnstageAllFiles)

		if err := self.optimisticChange(root, self.optimisticUnstage); err != nil {
			return err
		}

		if err := self.c.Git().WorkingTree.UnstageAll(); err != nil {
			return self.c.Error(err)
		}
	}

	return nil
}

func (self *FilesController) unstageFiles(node *filetree.FileNode) error {
	return node.ForEachFile(func(file *models.File) error {
		if file.HasStagedChanges {
			if err := self.c.Git().WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
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

	if err := self.c.Git().WorkingTree.RemoveTrackedFiles(node.GetPath()); err != nil {
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
	err := self.ignoreOrExcludeFile(node, self.c.Tr.IgnoreTracked, self.c.Tr.IgnoreTrackedPrompt, self.c.Tr.Actions.IgnoreExcludeFile, self.c.Git().WorkingTree.Ignore)
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

	err := self.ignoreOrExcludeFile(node, self.c.Tr.ExcludeTracked, self.c.Tr.ExcludeTrackedPrompt, self.c.Tr.Actions.ExcludeFile, self.c.Git().WorkingTree.Exclude)
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
				LabelColumns: []string{self.c.Tr.IgnoreFile},
				OnPress: func() error {
					if err := self.ignore(node); err != nil {
						return self.c.Error(err)
					}
					return nil
				},
				Key: 'i',
			},
			{
				LabelColumns: []string{self.c.Tr.ExcludeFile},
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

func (self *FilesController) refresh() error {
	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
}

func (self *FilesController) handleAmendCommitPress() error {
	if len(self.c.Model().Files) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoFilesStagedTitle)
	}

	if !self.c.Helpers().WorkingTree.AnyStagedFiles() {
		return self.c.Helpers().WorkingTree.PromptToStageAllAndRetry(self.handleAmendCommitPress)
	}

	if len(self.c.Model().Commits) == 0 {
		return self.c.ErrorMsg(self.c.Tr.NoCommitToAmend)
	}

	return self.c.Helpers().AmendHelper.AmendHead()
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
				Label: self.c.Tr.ResetFilter,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayAll)
				},
			},
		},
	})
}

func (self *FilesController) setStatusFiltering(filter filetree.FileTreeDisplayFilter) error {
	self.context().FileTreeViewModel.SetStatusFilter(filter)
	return self.c.PostRefreshUpdate(self.context())
}

func (self *FilesController) edit(node *filetree.FileNode) error {
	if node.File == nil {
		return self.c.ErrorMsg(self.c.Tr.ErrCannotEditDirectory)
	}

	return self.c.Helpers().Files.EditFile(node.GetPath())
}

func (self *FilesController) Open() error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	return self.c.Helpers().Files.OpenFile(node.GetPath())
}

func (self *FilesController) switchToMerge() error {
	file := self.getSelectedFile()
	if file == nil {
		return nil
	}

	return self.c.Helpers().MergeConflicts.SwitchToMerge(file.Name)
}

func (self *FilesController) createStashMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.StashOptions,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.StashAllChanges,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
						return self.c.ErrorMsg(self.c.Tr.NoFilesToStash)
					}
					return self.handleStashSave(self.c.Git().Stash.Push, self.c.Tr.Actions.StashAllChanges)
				},
				Key: 'a',
			},
			{
				Label: self.c.Tr.StashAllChangesKeepIndex,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
						return self.c.ErrorMsg(self.c.Tr.NoFilesToStash)
					}
					// if there are no staged files it behaves the same as Stash.Save
					return self.handleStashSave(self.c.Git().Stash.StashAndKeepIndex, self.c.Tr.Actions.StashAllChangesKeepIndex)
				},
				Key: 'i',
			},
			{
				Label: self.c.Tr.StashIncludeUntrackedChanges,
				OnPress: func() error {
					return self.handleStashSave(self.c.Git().Stash.StashIncludeUntrackedChanges, self.c.Tr.Actions.StashIncludeUntrackedChanges)
				},
				Key: 'U',
			},
			{
				Label: self.c.Tr.StashStagedChanges,
				OnPress: func() error {
					// there must be something in staging otherwise the current implementation mucks the stash up
					if !self.c.Helpers().WorkingTree.AnyStagedFiles() {
						return self.c.ErrorMsg(self.c.Tr.NoTrackedStagedFilesStash)
					}
					return self.handleStashSave(self.c.Git().Stash.SaveStagedChanges, self.c.Tr.Actions.StashStagedChanges)
				},
				Key: 's',
			},
			{
				Label: self.c.Tr.StashUnstagedChanges,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
						return self.c.ErrorMsg(self.c.Tr.NoFilesToStash)
					}
					if self.c.Helpers().WorkingTree.AnyStagedFiles() {
						return self.handleStashSave(self.c.Git().Stash.StashUnstagedChanges, self.c.Tr.Actions.StashUnstagedChanges)
					}
					// ordinary stash
					return self.handleStashSave(self.c.Git().Stash.Push, self.c.Tr.Actions.StashUnstagedChanges)
				},
				Key: 'u',
			},
		},
	})
}

func (self *FilesController) stash() error {
	return self.handleStashSave(self.c.Git().Stash.Push, self.c.Tr.Actions.StashAllChanges)
}

func (self *FilesController) createResetToUpstreamMenu() error {
	return self.c.Helpers().Refs.CreateGitResetMenu("@{upstream}")
}

func (self *FilesController) handleToggleDirCollapsed() error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	self.context().FileTreeViewModel.ToggleCollapsed(node.GetPath())

	if err := self.c.PostRefreshUpdate(self.c.Contexts().Files); err != nil {
		self.c.Log.Error(err)
	}

	return nil
}

func (self *FilesController) toggleTreeView() error {
	self.context().FileTreeViewModel.ToggleShowTree()

	return self.c.PostRefreshUpdate(self.context())
}

func (self *FilesController) handleStashSave(stashFunc func(message string) error, action string) error {
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
	return self.EnterFile(types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: opts.Y})
}

func (self *FilesController) onClickSecondary(opts gocui.ViewMouseBindingOpts) error {
	return self.EnterFile(types.OnFocusOpts{ClickedWindowName: "secondary", ClickedViewLineIdx: opts.Y})
}

func (self *FilesController) fetch() error {
	return self.c.WithLoaderPanel(self.c.Tr.FetchWait, func(task gocui.Task) error {
		if err := self.fetchAux(task); err != nil {
			_ = self.c.Error(err)
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *FilesController) fetchAux(task gocui.Task) (err error) {
	self.c.LogAction("Fetch")
	err = self.c.Git().Sync.Fetch(task)

	if err != nil && strings.Contains(err.Error(), "exit status 128") {
		_ = self.c.ErrorMsg(self.c.Tr.PassUnameWrong)
	}

	_ = self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS}, Mode: types.ASYNC})

	return err
}
