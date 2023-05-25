package helpers

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/types/enums"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ModeHelper struct {
	c *HelperCommon

	diffHelper           *DiffHelper
	patchBuildingHelper  *PatchBuildingHelper
	cherryPickHelper     *CherryPickHelper
	mergeAndRebaseHelper *MergeAndRebaseHelper
	bisectHelper         *BisectHelper
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
	Description func() string
	Reset       func() error
}

func (self *ModeHelper) Statuses() []ModeStatus {
	return []ModeStatus{
		{
			IsActive: self.c.Modes().Diffing.Active,
			Description: func() string {
				return self.withResetButton(
					fmt.Sprintf(
						"%s %s",
						self.c.Tr.ShowingGitDiff,
						"git diff "+strings.Join(self.diffHelper.DiffArgs(), " "),
					),
					style.FgMagenta,
				)
			},
			Reset: self.diffHelper.ExitDiffMode,
		},
		{
			IsActive: self.c.Git().Patch.PatchBuilder.Active,
			Description: func() string {
				return self.withResetButton(self.c.Tr.BuildingPatch, style.FgYellow.SetBold())
			},
			Reset: self.patchBuildingHelper.Reset,
		},
		{
			IsActive: self.c.Modes().Filtering.Active,
			Description: func() string {
				return self.withResetButton(
					fmt.Sprintf(
						"%s '%s'",
						self.c.Tr.FilteringBy,
						self.c.Modes().Filtering.GetPath(),
					),
					style.FgRed,
				)
			},
			Reset: self.ExitFilterMode,
		},
		{
			IsActive: self.c.Modes().CherryPicking.Active,
			Description: func() string {
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
			Reset: self.cherryPickHelper.Reset,
		},
		{
			IsActive: func() bool {
				return self.c.Git().Status.WorkingTreeState() != enums.REBASE_MODE_NONE
			},
			Description: func() string {
				workingTreeState := self.c.Git().Status.WorkingTreeState()
				return self.withResetButton(
					presentation.FormatWorkingTreeStateTitle(self.c.Tr, workingTreeState), style.FgYellow,
				)
			},
			Reset: self.mergeAndRebaseHelper.AbortMergeOrRebaseWithConfirm,
		},
		{
			IsActive: func() bool {
				return self.c.Model().BisectInfo.Started()
			},
			Description: func() string {
				return self.withResetButton(self.c.Tr.Bisect.Bisecting, style.FgGreen)
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
	return slices.Find(self.Statuses(), func(mode ModeStatus) bool {
		return mode.IsActive()
	})
}

func (self *ModeHelper) IsAnyModeActive() bool {
	return slices.Some(self.Statuses(), func(mode ModeStatus) bool {
		return mode.IsActive()
	})
}

func (self *ModeHelper) ExitFilterMode() error {
	return self.ClearFiltering()
}

func (self *ModeHelper) ClearFiltering() error {
	self.c.Modes().Filtering.Reset()
	if self.c.State().GetRepoState().GetScreenMode() == types.SCREEN_HALF {
		self.c.State().GetRepoState().SetScreenMode(types.SCREEN_NORMAL)
	}

	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMITS}})
}
