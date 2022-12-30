package helpers

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
	Snake          *SnakeHelper
	// lives in context package because our contexts need it to render to main
	Diff            *DiffHelper
	Repos           *ReposHelper
	RecordDirectory *RecordDirectoryHelper
	Update          *UpdateHelper
	Window          *WindowHelper
	View            *ViewHelper
	Refresh         *RefreshHelper
}

func NewStubHelpers() *Helpers {
	return &Helpers{
		Refs:            &RefsHelper{},
		Bisect:          &BisectHelper{},
		Suggestions:     &SuggestionsHelper{},
		Files:           &FilesHelper{},
		WorkingTree:     &WorkingTreeHelper{},
		Tags:            &TagsHelper{},
		MergeAndRebase:  &MergeAndRebaseHelper{},
		MergeConflicts:  &MergeConflictsHelper{},
		CherryPick:      &CherryPickHelper{},
		Host:            &HostHelper{},
		PatchBuilding:   &PatchBuildingHelper{},
		Staging:         &StagingHelper{},
		GPG:             &GpgHelper{},
		Upstream:        &UpstreamHelper{},
		AmendHelper:     &AmendHelper{},
		Snake:           &SnakeHelper{},
		Diff:            &DiffHelper{},
		Repos:           &ReposHelper{},
		RecordDirectory: &RecordDirectoryHelper{},
		Update:          &UpdateHelper{},
		Window:          &WindowHelper{},
		View:            &ViewHelper{},
		Refresh:         &RefreshHelper{},
	}
}
