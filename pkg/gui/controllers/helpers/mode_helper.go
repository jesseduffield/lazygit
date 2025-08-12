package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

type ModeHelper struct {
	c *HelperCommon

	diffHelper           *DiffHelper
	patchBuildingHelper  *PatchBuildingHelper
	cherryPickHelper     *CherryPickHelper
	mergeAndRebaseHelper *MergeAndRebaseHelper
	bisectHelper         *BisectHelper
	suppressRebasingMode bool
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
				return !self.suppressRebasingMode && self.c.Git().Status.WorkingTreeState().Any()
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

func (self *ModeHelper) ClearFiltering() error {
	selectedCommitHash := self.c.Contexts().LocalCommits.GetSelectedCommitHash()
	self.c.Modes().Filtering.Reset()
	if self.c.State().GetRepoState().GetScreenMode() == types.SCREEN_HALF {
		self.c.State().GetRepoState().SetScreenMode(types.SCREEN_NORMAL)
	}

	self.c.Refresh(types.RefreshOptions{
		Scope: ScopesToRefreshWhenFilteringModeChanges(),
		Then: func() {
			// Find the commit that was last selected in filtering mode, and select it again after refreshing
			if !self.c.Contexts().LocalCommits.SelectCommitByHash(selectedCommitHash) {
				// If we couldn't find it (either because no commit was selected
				// in filtering mode, or because the commit is outside the
				// initial 300 range), go back to the commit that was selected
				// before we entered filtering
				self.c.Contexts().LocalCommits.SelectCommitByHash(self.c.Modes().Filtering.GetSelectedCommitHash())
			}

			self.c.PostRefreshUpdate(self.c.Contexts().LocalCommits)
		},
	})
	return nil
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

func (self *ModeHelper) SetSuppressRebasingMode(value bool) {
	self.suppressRebasingMode = value
}
