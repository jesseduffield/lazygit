package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type HelperCommon struct {
	*common.Common
	types.IGuiCommon
	IGetContexts
}

type IGetContexts interface {
	Contexts() *context.ContextTree
}

type Helpers struct {
	Refs           *RefsHelper
	Bisect         *BisectHelper
	Suggestions    *SuggestionsHelper
	Files          *FilesHelper
	WorkingTree    *WorkingTreeHelper
	Tags           *TagsHelper
	MergeAndRebase *MergeAndRebaseHelper
	MergeConflicts *MergeConflictsHelper
	CherryPick     *CherryPickHelper
	Host           *HostHelper
	PatchBuilding  *PatchBuildingHelper
	Staging        *StagingHelper
	GPG            *GpgHelper
	Upstream       *UpstreamHelper
	AmendHelper    *AmendHelper
	Commits        *CommitsHelper
	Snake          *SnakeHelper
	// lives in context package because our contexts need it to render to main
	Diff              *DiffHelper
	Repos             *ReposHelper
	RecordDirectory   *RecordDirectoryHelper
	Update            *UpdateHelper
	Window            *WindowHelper
	View              *ViewHelper
	Refresh           *RefreshHelper
	Confirmation      *ConfirmationHelper
	Mode              *ModeHelper
	AppStatus         *AppStatusHelper
	WindowArrangement *WindowArrangementHelper
	Search            *SearchHelper
	Worktree          *WorktreeHelper
}

func NewStubHelpers() *Helpers {
	return &Helpers{
		Refs:              &RefsHelper{},
		Bisect:            &BisectHelper{},
		Suggestions:       &SuggestionsHelper{},
		Files:             &FilesHelper{},
		WorkingTree:       &WorkingTreeHelper{},
		Tags:              &TagsHelper{},
		MergeAndRebase:    &MergeAndRebaseHelper{},
		MergeConflicts:    &MergeConflictsHelper{},
		CherryPick:        &CherryPickHelper{},
		Host:              &HostHelper{},
		PatchBuilding:     &PatchBuildingHelper{},
		Staging:           &StagingHelper{},
		GPG:               &GpgHelper{},
		Upstream:          &UpstreamHelper{},
		AmendHelper:       &AmendHelper{},
		Commits:           &CommitsHelper{},
		Snake:             &SnakeHelper{},
		Diff:              &DiffHelper{},
		Repos:             &ReposHelper{},
		RecordDirectory:   &RecordDirectoryHelper{},
		Update:            &UpdateHelper{},
		Window:            &WindowHelper{},
		View:              &ViewHelper{},
		Refresh:           &RefreshHelper{},
		Confirmation:      &ConfirmationHelper{},
		Mode:              &ModeHelper{},
		AppStatus:         &AppStatusHelper{},
		WindowArrangement: &WindowArrangementHelper{},
		Search:            &SearchHelper{},
		Worktree:          &WorktreeHelper{},
	}
}
