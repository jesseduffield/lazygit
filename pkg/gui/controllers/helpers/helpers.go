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
	GPG            *GpgHelper
	Upstream       *UpstreamHelper
	AmendHelper    *AmendHelper
	Commits        *CommitsHelper
}

func NewStubHelpers() *Helpers {
	return &Helpers{
		Refs:           &RefsHelper{},
		Bisect:         &BisectHelper{},
		Suggestions:    &SuggestionsHelper{},
		Files:          &FilesHelper{},
		WorkingTree:    &WorkingTreeHelper{},
		Tags:           &TagsHelper{},
		MergeAndRebase: &MergeAndRebaseHelper{},
		MergeConflicts: &MergeConflictsHelper{},
		CherryPick:     &CherryPickHelper{},
		Host:           &HostHelper{},
		PatchBuilding:  &PatchBuildingHelper{},
		GPG:            &GpgHelper{},
		Upstream:       &UpstreamHelper{},
		AmendHelper:    &AmendHelper{},
		Commits:        &CommitsHelper{},
	}
}
