package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type ModeHelper struct {
	c *HelperCommon

	diffHelper                   *DiffHelper
	patchBuildingHelper          *PatchBuildingHelper
	cherryPickHelper             *CherryPickHelper
	mergeAndRebaseHelper         *MergeAndRebaseHelper
	bisectHelper                 *BisectHelper
	suppressWorkingTreeStateMode bool
}

func NewModeHelper(
	c *HelperCommon,
	diffHelper *DiffHelper,
	patchBuildingHelper *PatchBuildingHelper,
	cherryPickHelper *CherryPickHelper,
	mergeAndRebaseHelper *MergeAndRebaseHelper,
	bisectHelper *BisectHelper,
) *ModeHelper {
	return &ModeHelper{
		c:                    c,
		diffHelper:           diffHelper,
		patchBuildingHelper:  patchBuildingHelper,
		cherryPickHelper:     cherryPickHelper,
		mergeAndRebaseHelper: mergeAndRebaseHelper,
		bisectHelper:         bisectHelper,
	}
}

type ModeStatus struct {
	IsActive    func() bool
	InfoLabel   func() string
	CancelLabel func() string
	Reset       func() error
}

func (self *ModeHelper) Statuses() []ModeStatus {
	return []ModeStatus{
		{
			IsActive: self.c.Modes().Diffing.Active,
			InfoLabel: func() string {
				return self.withResetButton(
					fmt.Sprintf(
						"%s %s",
						self.c.Tr.ShowingGitDiff,
						"git diff "+strings.Join(self.diffHelper.DiffArgs(), " "),
					),
					style.FgMagenta,
				)
			},
			CancelLabel: func() string {
				return self.c.Tr.CancelDiffingMode
			},
			Reset: self.diffHelper.ExitDiffMode,
		},
		{
			IsActive: self.c.Git().Patch.PatchBuilder.Active,
			InfoLabel: func() string {
				return self.withResetButton(self.c.Tr.BuildingPatch, style.FgYellow.SetBold())
			},
			CancelLabel: func() string {
				return self.c.Tr.ExitCustomPatchBuilder
			},
			Reset: self.patchBuildingHelper.Reset,
		},
		{
			IsActive: self.c.Modes().Filtering.Active,
			InfoLabel: func() string {
				filterContent := lo.Ternary(self.c.Modes().Filtering.GetPath() != "", self.c.Modes().Filtering.GetPath(), self.c.Modes().Filtering.GetAuthor())
				return self.withResetButton(
					fmt.Sprintf(
						"%s '%s'",
						self.c.Tr.FilteringBy,
						filterContent,
					),
					style.FgRed,
				)
			},
			CancelLabel: func() string {
				return self.c.Tr.ExitFilterMode
			},
			Reset: self.ExitFilterMode,
		},
		{
			IsActive: self.c.Modes().MarkedBaseCommit.Active,
			InfoLabel: func() string {
				return self.withResetButton(
					self.c.Tr.MarkedBaseCommitStatus,
					style.FgCyan,
				)
			},
			CancelLabel: func() string {
				return self.c.Tr.CancelMarkedBaseCommit
			},
			Reset: self.mergeAndRebaseHelper.ResetMarkedBaseCommit,
		},
		{
			IsActive: self.c.Modes().CherryPicking.Active,
			InfoLabel: func() string {
				copiedCount := len(self.c.Modes().CherryPicking.CherryPickedCommits)
				text := self.c.Tr.CommitsCopied
				if copiedCount == 1 {
					text = self.c.Tr.CommitCopied
				}

				return self.withResetButton(
					fmt.Sprintf(
						"%d %s",
						copiedCount,
						text,
					),
					style.FgCyan,
				)
			},
			CancelLabel: func() string {
				return self.c.Tr.ResetCherryPickShort
			},
			Reset: self.cherryPickHelper.Reset,
		},
		{
			IsActive: func() bool {
				return !self.suppressWorkingTreeStateMode && self.c.Git().Status.WorkingTreeState().Any()
			},
			InfoLabel: func() string {
				workingTreeState := self.c.Git().Status.WorkingTreeState()
				return self.withResetButton(
					workingTreeState.Title(self.c.Tr), style.FgYellow,
				)
			},
			CancelLabel: func() string {
				return fmt.Sprintf(self.c.Tr.AbortTitle, self.c.Git().Status.WorkingTreeState().CommandName())
			},
			Reset: self.mergeAndRebaseHelper.AbortMergeOrRebaseWithConfirm,
		},
		{
			IsActive: func() bool {
				return self.c.Model().BisectInfo.Started()
			},
			InfoLabel: func() string {
				return self.withResetButton(self.c.Tr.Bisect.Bisecting, style.FgGreen)
			},
			CancelLabel: func() string {
				return self.c.Tr.Actions.ResetBisect
			},
			Reset: self.bisectHelper.Reset,
		},
	}
}

func (self *ModeHelper) withResetButton(content string, textStyle style.TextStyle) string {
	return textStyle.Sprintf(
		"%s %s",
		content,
		style.AttrUnderline.Sprint(self.c.Tr.ResetInParentheses),
	)
}

func (self *ModeHelper) GetActiveMode() (ModeStatus, bool) {
	return lo.Find(self.Statuses(), func(mode ModeStatus) bool {
		return mode.IsActive()
	})
}

func (self *ModeHelper) IsAnyModeActive() bool {
	return lo.SomeBy(self.Statuses(), func(mode ModeStatus) bool {
		return mode.IsActive()
	})
}

func (self *ModeHelper) ExitFilterMode() error {
	return self.ClearFiltering()
}

func (self *ModeHelper) SetFilteringPath(path string) error {
	return self.setFiltering(func() {
		self.c.Modes().Filtering.SetPath(path)
	})
}

func (self *ModeHelper) SetFilteringAuthor(author string) error {
	return self.setFiltering(func() {
		self.c.Modes().Filtering.SetAuthor(author)
	})
}

func (self *ModeHelper) setFiltering(setFilter func()) error {
	return self.changeFiltering(
		func() {
			// Whatever we were filtering by before is replaced, not added to
			self.c.Modes().Filtering.Reset()
			setFilter()
			self.c.Modes().Filtering.SetSelectedCommitHash(
				self.c.Contexts().LocalCommits.GetSelectedCommitHash())
		},
		func() {
			self.c.Contexts().LocalCommits.SetSelection(0)
		},
	)
}

func (self *ModeHelper) ClearFiltering() error {
	selectedCommitHash := self.c.Contexts().LocalCommits.GetSelectedCommitHash()

	return self.changeFiltering(
		self.c.Modes().Filtering.Reset,
		func() {
			// Find the commit that was last selected in filtering mode, and select it again after refreshing
			if !self.c.Contexts().LocalCommits.SelectCommitByHash(selectedCommitHash) {
				// If we couldn't find it (either because no commit was selected
				// in filtering mode, or because the commit is outside the
				// initial 300 range), go back to the commit that was selected
				// before we entered filtering
				self.c.Contexts().LocalCommits.SelectCommitByHash(self.c.Modes().Filtering.GetSelectedCommitHash())
			}
		},
	)
}

// changeFiltering applies a change to the filtering mode: setFilter mutates the
// mode, then the views whose contents depend on the filter are reloaded, and
// selectCommit puts the selection where it belongs in the reloaded commit list.
//
// Reloading the commit list can take seconds in a big repo, so it happens on a
// worker with a waiting status. Everything the user can see of the change waits
// for it: the screen mode, the focused panel and the reloaded lists all land in
// the same frame, from the refresh's Then, rather than framing an unfiltered
// list as if it were the filtered one. Until then the pre-change state stays on
// screen, and it stays consistent, because the only thing that has changed
// behind it is the filter that the reload is in the middle of applying. The one
// thing that can't wait is the mode indicator in the information panel: the
// filter has to be set before the reload can use it, so the indicator leads the
// lists by however long the reload takes.
//
// Input is blocked for the duration: the keys the user presses arrive after the
// change, which is where they meant them to go, and it keeps a second filter
// change from racing this one — they would both refresh with whichever filter
// happened to be set when their git commands ran.
func (self *ModeHelper) changeFiltering(setFilter func(), selectCommit func()) error {
	setFilter()

	filtering := self.c.Modes().Filtering.Active()
	message := lo.Ternary(filtering, self.c.Tr.ApplyingFilterStatus, self.c.Tr.RemovingFilterStatus)

	return self.c.WithWaitingStatusBlockingInput(types.WaitingStatusOpts{Message: message}, func(gocui.Task) error {
		self.c.RefreshFromWorker(types.RefreshOptions{
			Scope:          ScopesToRefreshWhenFilteringModeChanges(),
			BatchUIUpdates: true,
			Then: func() error {
				repoState := self.c.State().GetRepoState()
				if filtering {
					if repoState.GetScreenMode() == types.SCREEN_NORMAL {
						repoState.SetScreenMode(types.SCREEN_HALF)
					}
					self.c.Context().Push(self.c.Contexts().LocalCommits, types.OnFocusOpts{})
				} else if repoState.GetScreenMode() == types.SCREEN_HALF {
					repoState.SetScreenMode(types.SCREEN_NORMAL)
				}

				selectCommit()
				self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
				return nil
			},
		})
		return nil
	})
}

// Stashes really only need to be refreshed when filtering by path, not by author, but it's too much
// work to distinguish this, and refreshing stashes is fast, so we don't bother
func ScopesToRefreshWhenFilteringModeChanges() []types.RefreshableView {
	return []types.RefreshableView{
		types.COMMITS,
		types.SUB_COMMITS,
		types.REFLOG,
		types.STASH,
	}
}

func (self *ModeHelper) SetSuppressWorkingTreeStateMode(value bool) {
	self.suppressWorkingTreeStateMode = value
}
