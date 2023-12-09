package types

import "github.com/jesseduffield/gocui"

type Views struct {
	Status         *gocui.View
	Submodules     *gocui.View
	Files          *gocui.View
	Branches       *gocui.View
	Remotes        *gocui.View
	Worktrees      *gocui.View
	Tags           *gocui.View
	RemoteBranches *gocui.View
	ReflogCommits  *gocui.View
	Commits        *gocui.View
	Stash          *gocui.View

	Main                   *gocui.View
	Secondary              *gocui.View
	Staging                *gocui.View
	StagingSecondary       *gocui.View
	PatchBuilding          *gocui.View
	PatchBuildingSecondary *gocui.View
	MergeConflicts         *gocui.View

	Options           *gocui.View
	Confirmation      *gocui.View
	Menu              *gocui.View
	CommitMessage     *gocui.View
	CommitDescription *gocui.View
	CommitFiles       *gocui.View
	SubCommits        *gocui.View
	Information       *gocui.View
	AppStatus         *gocui.View
	Search            *gocui.View
	SearchPrefix      *gocui.View
	StatusSpacer1     *gocui.View
	StatusSpacer2     *gocui.View
	Limit             *gocui.View
	Suggestions       *gocui.View
	Tooltip           *gocui.View
	Extras            *gocui.View

	// for playing the easter egg snake game
	Snake *gocui.View
}
