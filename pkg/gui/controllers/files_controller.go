package controllers

import (
	"errors"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type FilesController struct {
	baseController // nolint: unused
	*ListControllerTrait[*filetree.FileNode]
	c *ControllerCommon
}

var _ types.IController = &FilesController{}

func NewFilesController(
	c *ControllerCommon,
) *FilesController {
	return &FilesController{
		c: c,
		ListControllerTrait: NewListControllerTrait[*filetree.FileNode](
			c,
			c.Contexts().Files,
			c.Contexts().Files.GetSelected,
			c.Contexts().Files.GetSelectedItems,
		),
	}
}

func (self *FilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItems(self.press),
			GetDisabledReason: self.require(self.itemsSelected()),
			Description:       self.c.Tr.Stage,
			Tooltip:           self.c.Tr.StageTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.OpenStatusFilter),
			Handler:     self.handleStatusFilterPressed,
			Description: self.c.Tr.FileFilter,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.CopyFileInfoToClipboard),
			Handler:     self.openCopyMenu,
			Description: self.c.Tr.CopyToClipboardMenu,
			OpensMenu:   true,
		},
		{
			Key:             opts.GetKey(opts.Config.Files.CommitChanges),
			Handler:         self.c.Helpers().WorkingTree.HandleCommitPress,
			Description:     self.c.Tr.Commit,
			Tooltip:         self.c.Tr.CommitTooltip,
			DisplayOnScreen: true,
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
			Key:         opts.GetKey(opts.Config.Files.FindBaseCommitForFixup),
			Handler:     self.c.Helpers().FixupHelper.HandleFindBaseCommitForFixupPress,
			Description: self.c.Tr.FindBaseCommitForFixup,
			Tooltip:     self.c.Tr.FindBaseCommitForFixupTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Edit),
			Handler:           self.withItems(self.edit),
			GetDisabledReason: self.require(self.itemsSelected(self.canEditFiles)),
			Description:       self.c.Tr.Edit,
			Tooltip:           self.c.Tr.EditFileTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:           self.Open,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenFile,
			Tooltip:           self.c.Tr.OpenFileTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Files.IgnoreFile),
			Handler:           self.withItem(self.ignoreOrExcludeMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Actions.IgnoreExcludeFile,
			OpensMenu:         true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.RefreshFiles),
			Handler:     self.refresh,
			Description: self.c.Tr.RefreshFiles,
		},
		{
			Key:             opts.GetKey(opts.Config.Files.StashAllChanges),
			Handler:         self.stash,
			Description:     self.c.Tr.Stash,
			Tooltip:         self.c.Tr.StashTooltip,
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ViewStashOptions),
			Handler:     self.createStashMenu,
			Description: self.c.Tr.ViewStashOptions,
			Tooltip:     self.c.Tr.ViewStashOptionsTooltip,
			OpensMenu:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleStagedAll),
			Handler:     self.toggleStagedAll,
			Description: self.c.Tr.ToggleStagedAll,
			Tooltip:     self.c.Tr.ToggleStagedAllTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Handler:           self.enter,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.FileEnter,
			Tooltip:           self.c.Tr.FileEnterTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.remove),
			GetDisabledReason: self.require(self.itemsSelected(self.canRemove)),
			Description:       self.c.Tr.Discard,
			Tooltip:           self.c.Tr.DiscardFileChangesTooltip,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Key:         opts.GetKey(opts.Config.Commits.ViewResetOptions),
			Handler:     self.createResetToUpstreamMenu,
			Description: self.c.Tr.ViewResetToUpstreamOptions,
			OpensMenu:   true,
		},
		{
			Key:             opts.GetKey(opts.Config.Files.ViewResetOptions),
			Handler:         self.createResetMenu,
			Description:     self.c.Tr.Reset,
			Tooltip:         self.c.Tr.FileResetOptionsTooltip,
			OpensMenu:       true,
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.ToggleTreeView,
			Tooltip:     self.c.Tr.ToggleTreeViewTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler:           self.withItem(self.openDiffTool),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.OpenMergeTool),
			Handler:     self.c.Helpers().WorkingTree.OpenMergeTool,
			Description: self.c.Tr.OpenMergeTool,
			Tooltip:     self.c.Tr.OpenMergeToolTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.Fetch),
			Handler:     self.fetch,
			Description: self.c.Tr.Fetch,
			Tooltip:     self.c.Tr.FetchTooltip,
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

func (self *FilesController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			node := self.context().GetSelected()

			if node == nil {
				self.c.RenderToMainViews(types.RefreshMainOpts{
					Pair: self.c.MainViewPairs().Normal,
					Main: &types.ViewUpdateOpts{
						Title:    self.c.Tr.DiffTitle,
						SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
						Task:     types.NewRenderStringTask(self.c.Tr.NoChangedFiles),
					},
				})
				return
			}

			if node.File != nil && node.File.HasInlineMergeConflicts {
				hasConflicts, err := self.c.Helpers().MergeConflicts.SetMergeState(node.GetPath())
				if err != nil {
					return
				}

				if hasConflicts {
					self.c.Helpers().MergeConflicts.Render()
					return
				}
			}

			self.c.Helpers().MergeConflicts.ResetMergeState()

			pair := self.c.MainViewPairs().Normal
			if node.File != nil {
				pair = self.c.MainViewPairs().Staging
			}

			split := self.c.UserConfig().Gui.SplitDiff == "always" || (node.GetHasUnstagedChanges() && node.GetHasStagedChanges())
			mainShowsStaged := !split && node.GetHasStagedChanges()

			cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(node, false, mainShowsStaged)
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
				cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(node, false, true)

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

			self.c.RenderToMainViews(refreshOpts)
		})
	}
}

func (self *FilesController) GetOnClick() func() error {
	return self.withItemGraceful(func(node *filetree.FileNode) error {
		return self.press([]*filetree.FileNode{node})
	})
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
func (self *FilesController) optimisticChange(nodes []*filetree.FileNode, optimisticChangeFn func(*models.File) bool) error {
	rerender := false

	for _, node := range nodes {
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
	}

	if rerender {
		if err := self.c.PostRefreshUpdate(self.c.Contexts().Files); err != nil {
			return err
		}
	}

	return nil
}

func (self *FilesController) pressWithLock(selectedNodes []*filetree.FileNode) error {
	// Obtaining this lock because optimistic rendering requires us to mutate
	// the files in our model.
	self.c.Mutexes().RefreshingFilesMutex.Lock()
	defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

	for _, node := range selectedNodes {
		// if any files within have inline merge conflicts we can't stage or unstage,
		// or it'll end up with those >>>>>> lines actually staged
		if node.GetHasInlineMergeConflicts() {
			return errors.New(self.c.Tr.ErrStageDirWithInlineMergeConflicts)
		}
	}

	toPaths := func(nodes []*filetree.FileNode) []string {
		return lo.Map(nodes, func(node *filetree.FileNode, _ int) string {
			return node.Path
		})
	}

	selectedNodes = normalisedSelectedNodes(selectedNodes)

	// If any node has unstaged changes, we'll stage all the selected unstaged nodes (staging already staged deleted files/folders would fail).
	// Otherwise, we unstage all the selected nodes.
	unstagedSelectedNodes := filterNodesHaveUnstagedChanges(selectedNodes)

	if len(unstagedSelectedNodes) > 0 {
		self.c.LogAction(self.c.Tr.Actions.StageFile)

		if err := self.optimisticChange(unstagedSelectedNodes, self.optimisticStage); err != nil {
			return err
		}

		if err := self.c.Git().WorkingTree.StageFiles(toPaths(unstagedSelectedNodes)); err != nil {
			return err
		}
	} else {
		self.c.LogAction(self.c.Tr.Actions.UnstageFile)

		if err := self.optimisticChange(selectedNodes, self.optimisticUnstage); err != nil {
			return err
		}

		// need to partition the paths into tracked and untracked (where we assume directories are tracked). Then we'll run the commands separately.
		trackedNodes, untrackedNodes := utils.Partition(selectedNodes, func(node *filetree.FileNode) bool {
			// We treat all directories as tracked. I'm not actually sure why we do this but
			// it's been the existing behaviour for a while and nobody has complained
			return !node.IsFile() || node.GetIsTracked()
		})

		if len(untrackedNodes) > 0 {
			if err := self.c.Git().WorkingTree.UnstageUntrackedFiles(toPaths(untrackedNodes)); err != nil {
				return err
			}
		}

		if len(trackedNodes) > 0 {
			if err := self.c.Git().WorkingTree.UnstageTrackedFiles(toPaths(trackedNodes)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *FilesController) press(nodes []*filetree.FileNode) error {
	if err := self.pressWithLock(nodes); err != nil {
		return err
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Mode: types.ASYNC}); err != nil {
		return err
	}

	self.context().HandleFocus(types.OnFocusOpts{})
	return nil
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
		return errors.New(self.c.Tr.FileStagingRequirements)
	}

	self.c.Context().Push(self.c.Contexts().Staging, opts)
	return nil
}

func (self *FilesController) toggleStagedAll() error {
	if err := self.toggleStagedAllWithLock(); err != nil {
		return err
	}

	if err := self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Mode: types.ASYNC}); err != nil {
		return err
	}

	self.context().HandleFocus(types.OnFocusOpts{})
	return nil
}

func (self *FilesController) toggleStagedAllWithLock() error {
	self.c.Mutexes().RefreshingFilesMutex.Lock()
	defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

	root := self.context().FileTreeViewModel.GetRoot()

	// if any files within have inline merge conflicts we can't stage or unstage,
	// or it'll end up with those >>>>>> lines actually staged
	if root.GetHasInlineMergeConflicts() {
		return errors.New(self.c.Tr.ErrStageDirWithInlineMergeConflicts)
	}

	if root.GetHasUnstagedChanges() {
		self.c.LogAction(self.c.Tr.Actions.StageAllFiles)

		if err := self.optimisticChange([]*filetree.FileNode{root}, self.optimisticStage); err != nil {
			return err
		}

		if err := self.c.Git().WorkingTree.StageAll(); err != nil {
			return err
		}
	} else {
		self.c.LogAction(self.c.Tr.Actions.UnstageAllFiles)

		if err := self.optimisticChange([]*filetree.FileNode{root}, self.optimisticUnstage); err != nil {
			return err
		}

		if err := self.c.Git().WorkingTree.UnstageAll(); err != nil {
			return err
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
		self.c.Confirm(types.ConfirmOpts{
			Title:  trText,
			Prompt: trPrompt,
			HandleConfirm: func() error {
				return self.ignoreOrExcludeTracked(node, trAction, f)
			},
		})

		return nil
	}
	return self.ignoreOrExcludeUntracked(node, trAction, f)
}

func (self *FilesController) ignore(node *filetree.FileNode) error {
	if node.GetPath() == ".gitignore" {
		return errors.New(self.c.Tr.Actions.IgnoreFileErr)
	}
	return self.ignoreOrExcludeFile(node, self.c.Tr.IgnoreTracked, self.c.Tr.IgnoreTrackedPrompt, self.c.Tr.Actions.IgnoreExcludeFile, self.c.Git().WorkingTree.Ignore)
}

func (self *FilesController) exclude(node *filetree.FileNode) error {
	if node.GetPath() == ".gitignore" {
		return errors.New(self.c.Tr.Actions.ExcludeGitIgnoreErr)
	}

	return self.ignoreOrExcludeFile(node, self.c.Tr.ExcludeTracked, self.c.Tr.ExcludeTrackedPrompt, self.c.Tr.Actions.ExcludeFile, self.c.Git().WorkingTree.Exclude)
}

func (self *FilesController) ignoreOrExcludeMenu(node *filetree.FileNode) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.IgnoreExcludeFile,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.IgnoreFile},
				OnPress: func() error {
					if err := self.ignore(node); err != nil {
						return err
					}
					return nil
				},
				Key: 'i',
			},
			{
				LabelColumns: []string{self.c.Tr.ExcludeFile},
				OnPress: func() error {
					if err := self.exclude(node); err != nil {
						return err
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
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.AmendLastCommitTitle,
		Prompt: self.c.Tr.SureToAmend,
		HandleConfirm: func() error {
			return self.c.Helpers().WorkingTree.WithEnsureCommitableFiles(func() error {
				if len(self.c.Model().Commits) == 0 {
					return errors.New(self.c.Tr.NoCommitToAmend)
				}

				return self.c.Helpers().AmendHelper.AmendHead()
			})
		},
	})

	return nil
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
				Key: 's',
			},
			{
				Label: self.c.Tr.FilterUnstagedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayUnstaged)
				},
				Key: 'u',
			},
			{
				Label: self.c.Tr.ResetFilter,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayAll)
				},
				Key: 'r',
			},
		},
	})
}

func (self *FilesController) setStatusFiltering(filter filetree.FileTreeDisplayFilter) error {
	self.context().FileTreeViewModel.SetStatusFilter(filter)
	return self.c.PostRefreshUpdate(self.context())
}

func (self *FilesController) edit(nodes []*filetree.FileNode) error {
	return self.c.Helpers().Files.EditFiles(lo.FilterMap(nodes,
		func(node *filetree.FileNode, _ int) (string, bool) {
			return node.GetPath(), node.IsFile()
		}))
}

func (self *FilesController) canEditFiles(nodes []*filetree.FileNode) *types.DisabledReason {
	if lo.NoneBy(nodes, func(node *filetree.FileNode) bool { return node.IsFile() }) {
		return &types.DisabledReason{
			Text:             self.c.Tr.ErrCannotEditDirectory,
			ShowErrorInPanel: true,
		}
	}

	return nil
}

func (self *FilesController) Open() error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	return self.c.Helpers().Files.OpenFile(node.GetPath())
}

func (self *FilesController) openDiffTool(node *filetree.FileNode) error {
	fromCommit := ""
	reverse := false
	if self.c.Modes().Diffing.Active() {
		fromCommit = self.c.Modes().Diffing.Ref
		reverse = self.c.Modes().Diffing.Reverse
	}
	return self.c.RunSubprocessAndRefresh(
		self.c.Git().Diff.OpenDiffToolCmdObj(
			git_commands.DiffToolCmdOptions{
				Filepath:    node.Path,
				FromCommit:  fromCommit,
				ToCommit:    "",
				Reverse:     reverse,
				IsDirectory: !node.IsFile(),
				Staged:      !node.GetHasUnstagedChanges(),
			}),
	)
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
						return errors.New(self.c.Tr.NoFilesToStash)
					}
					return self.handleStashSave(self.c.Git().Stash.Push, self.c.Tr.Actions.StashAllChanges)
				},
				Key: 'a',
			},
			{
				Label: self.c.Tr.StashAllChangesKeepIndex,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
						return errors.New(self.c.Tr.NoFilesToStash)
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
						return errors.New(self.c.Tr.NoTrackedStagedFilesStash)
					}
					return self.handleStashSave(self.c.Git().Stash.SaveStagedChanges, self.c.Tr.Actions.StashStagedChanges)
				},
				Key: 's',
			},
			{
				Label: self.c.Tr.StashUnstagedChanges,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirty() {
						return errors.New(self.c.Tr.NoFilesToStash)
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

func (self *FilesController) openCopyMenu() error {
	node := self.context().GetSelected()

	copyNameItem := &types.MenuItem{
		Label: self.c.Tr.CopyFileName,
		OnPress: func() error {
			if err := self.c.OS().CopyToClipboard(node.Name()); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FileNameCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Key:            'n',
	}
	copyPathItem := &types.MenuItem{
		Label: self.c.Tr.CopyFilePath,
		OnPress: func() error {
			if err := self.c.OS().CopyToClipboard(node.Path); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FilePathCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Key:            'p',
	}
	copyFileDiffItem := &types.MenuItem{
		Label:   self.c.Tr.CopySelectedDiff,
		Tooltip: self.c.Tr.CopyFileDiffTooltip,
		OnPress: func() error {
			path := self.context().GetSelectedPath()
			hasStaged := self.hasPathStagedChanges(node)
			diff, err := self.c.Git().Diff.GetPathDiff(path, hasStaged)
			if err != nil {
				return err
			}
			if err := self.c.OS().CopyToClipboard(diff); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FileDiffCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected(
			func(file *filetree.FileNode) *types.DisabledReason {
				if !node.GetHasStagedOrTrackedChanges() {
					return &types.DisabledReason{Text: self.c.Tr.NoContentToCopyError}
				}
				return nil
			},
		))(),
		Key: 's',
	}
	copyAllDiff := &types.MenuItem{
		Label:   self.c.Tr.CopyAllFilesDiff,
		Tooltip: self.c.Tr.CopyFileDiffTooltip,
		OnPress: func() error {
			hasStaged := self.c.Helpers().WorkingTree.AnyStagedFiles()
			diff, err := self.c.Git().Diff.GetAllDiff(hasStaged)
			if err != nil {
				return err
			}
			if err := self.c.OS().CopyToClipboard(diff); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.AllFilesDiffCopiedToast)
			return nil
		},
		DisabledReason: self.require(
			func() *types.DisabledReason {
				if !self.anyStagedOrTrackedFile() {
					return &types.DisabledReason{Text: self.c.Tr.NoContentToCopyError}
				}
				return nil
			},
		)(),
		Key: 'a',
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CopyToClipboardMenu,
		Items: []*types.MenuItem{
			copyNameItem,
			copyPathItem,
			copyFileDiffItem,
			copyAllDiff,
		},
	})
}

func (self *FilesController) anyStagedOrTrackedFile() bool {
	if !self.c.Helpers().WorkingTree.AnyStagedFiles() {
		return self.c.Helpers().WorkingTree.AnyTrackedFiles()
	}
	return true
}

func (self *FilesController) hasPathStagedChanges(node *filetree.FileNode) bool {
	return node.SomeFile(func(t *models.File) bool {
		return t.HasStagedChanges
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
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.StashChanges,
		HandleConfirm: func(stashComment string) error {
			self.c.LogAction(action)

			if err := stashFunc(stashComment); err != nil {
				return err
			}
			return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
		},
	})

	return nil
}

func (self *FilesController) onClickMain(opts gocui.ViewMouseBindingOpts) error {
	return self.EnterFile(types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: opts.Y})
}

func (self *FilesController) onClickSecondary(opts gocui.ViewMouseBindingOpts) error {
	return self.EnterFile(types.OnFocusOpts{ClickedWindowName: "secondary", ClickedViewLineIdx: opts.Y})
}

func (self *FilesController) fetch() error {
	return self.c.WithWaitingStatus(self.c.Tr.FetchingStatus, func(task gocui.Task) error {
		if err := self.fetchAux(task); err != nil {
			return err
		}
		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	})
}

func (self *FilesController) fetchAux(task gocui.Task) (err error) {
	self.c.LogAction("Fetch")
	err = self.c.Git().Sync.Fetch(task)

	if err != nil && strings.Contains(err.Error(), "exit status 128") {
		return errors.New(self.c.Tr.PassUnameWrong)
	}

	_ = self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.BRANCHES, types.COMMITS, types.REMOTES, types.TAGS}, Mode: types.ASYNC})

	return err
}

// Couldn't think of a better term than 'normalised'. Alas.
// The idea is that when you select a range of nodes, you will often have both
// a node and its parent node selected. If we are trying to discard changes to the
// selected nodes, we'll get an error if we try to discard the child after the parent.
// So we just need to filter out any nodes from the selection that are descendants
// of other nodes
func normalisedSelectedNodes(selectedNodes []*filetree.FileNode) []*filetree.FileNode {
	return lo.Filter(selectedNodes, func(node *filetree.FileNode, _ int) bool {
		return !isDescendentOfSelectedNodes(node, selectedNodes)
	})
}

func isDescendentOfSelectedNodes(node *filetree.FileNode, selectedNodes []*filetree.FileNode) bool {
	for _, selectedNode := range selectedNodes {
		if selectedNode.IsFile() {
			continue
		}

		selectedNodePath := selectedNode.GetPath()
		nodePath := node.GetPath()

		if strings.HasPrefix(nodePath, selectedNodePath+"/") {
			return true
		}
	}
	return false
}

func someNodesHaveUnstagedChanges(nodes []*filetree.FileNode) bool {
	return lo.SomeBy(nodes, (*filetree.FileNode).GetHasUnstagedChanges)
}

func someNodesHaveStagedChanges(nodes []*filetree.FileNode) bool {
	return lo.SomeBy(nodes, (*filetree.FileNode).GetHasStagedChanges)
}

func filterNodesHaveUnstagedChanges(nodes []*filetree.FileNode) []*filetree.FileNode {
	return lo.Filter(nodes, func(node *filetree.FileNode, _ int) bool {
		return node.GetHasUnstagedChanges()
	})
}

func (self *FilesController) canRemove(selectedNodes []*filetree.FileNode) *types.DisabledReason {
	submodules := self.c.Model().Submodules
	submoduleCount := lo.CountBy(selectedNodes, func(node *filetree.FileNode) bool {
		return node.File != nil && node.File.IsSubmodule(submodules)
	})
	if submoduleCount > 0 && len(selectedNodes) > 1 {
		return &types.DisabledReason{Text: self.c.Tr.RangeSelectNotSupportedForSubmodules}
	}

	return nil
}

func (self *FilesController) remove(selectedNodes []*filetree.FileNode) error {
	submodules := self.c.Model().Submodules

	// If we have one submodule then we must only have one submodule or `canRemove` would have
	// returned an error
	firstNode := selectedNodes[0]
	if firstNode.File != nil && firstNode.File.IsSubmodule(submodules) {
		submodule := firstNode.File.SubmoduleConfig(submodules)

		menuItems := []*types.MenuItem{
			{
				Label: self.c.Tr.SubmoduleStashAndReset,
				OnPress: func() error {
					return self.ResetSubmodule(submodule)
				},
			},
		}

		return self.c.Menu(types.CreateMenuOptions{Title: firstNode.GetPath(), Items: menuItems})
	}

	selectedNodes = normalisedSelectedNodes(selectedNodes)

	discardAllChangesItem := types.MenuItem{
		Label: self.c.Tr.DiscardAllChanges,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInFile)

			if self.context().IsSelectingRange() {
				defer self.context().CancelRangeSelect()
			}

			for _, node := range selectedNodes {
				if err := self.c.Git().WorkingTree.DiscardAllDirChanges(node); err != nil {
					return err
				}
			}

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.WORKTREES}})
		},
		Key: self.c.KeybindingsOpts().GetKey(self.c.UserConfig().Keybinding.Files.ConfirmDiscard),
		Tooltip: utils.ResolvePlaceholderString(
			self.c.Tr.DiscardAllTooltip,
			map[string]string{
				"path": self.formattedPaths(selectedNodes),
			},
		),
	}

	discardUnstagedChangesItem := types.MenuItem{
		Label: self.c.Tr.DiscardUnstagedChanges,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.DiscardAllUnstagedChangesInFile)

			if self.context().IsSelectingRange() {
				defer self.context().CancelRangeSelect()
			}

			for _, node := range selectedNodes {
				if err := self.c.Git().WorkingTree.DiscardUnstagedDirChanges(node); err != nil {
					return err
				}
			}

			return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.WORKTREES}})
		},
		Key: 'u',
		Tooltip: utils.ResolvePlaceholderString(
			self.c.Tr.DiscardUnstagedTooltip,
			map[string]string{
				"path": self.formattedPaths(selectedNodes),
			},
		),
	}

	if !someNodesHaveStagedChanges(selectedNodes) || !someNodesHaveUnstagedChanges(selectedNodes) {
		discardUnstagedChangesItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.DiscardUnstagedDisabled}
	}

	menuItems := []*types.MenuItem{
		&discardAllChangesItem,
		&discardUnstagedChangesItem,
	}

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.DiscardChangesTitle, Items: menuItems})
}

func (self *FilesController) ResetSubmodule(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.ResettingSubmoduleStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.ResetSubmodule)

		file := self.c.Helpers().WorkingTree.FileForSubmodule(submodule)
		if file != nil {
			if err := self.c.Git().WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return err
			}
		}

		if err := self.c.Git().Submodule.Stash(submodule); err != nil {
			return err
		}
		if err := self.c.Git().Submodule.Reset(submodule); err != nil {
			return err
		}

		return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.SUBMODULES}})
	})
}

func (self *FilesController) formattedPaths(nodes []*filetree.FileNode) string {
	return utils.FormatPaths(lo.Map(nodes, func(node *filetree.FileNode, _ int) string {
		return node.GetPath()
	}))
}
