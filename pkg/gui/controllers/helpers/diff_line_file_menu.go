package helpers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// OpenJumpToFileMenu offers the files of the diff the given pane is showing in a menu,
// so that one of them can be gone to directly rather than by stepping through the diff
// a file at a time. Picking one goes to where that file's diff begins, the way stepping
// to it with next-file does: the selection moves there in a focused pane, and an
// unfocused one scrolls the file to the top.
//
// The diff is read to the end before the menu is built: a file below the part of it that
// has been read so far is in neither the list nor the view, and reaching the far end of
// a long diff is what the menu is for.
func (self *DiffLineHelper) OpenJumpToFileMenu(pane types.DiffPaneContext, title string) error {
	manager := self.c.GetViewBufferManagerForView(pane.GetView())
	if manager == nil {
		return nil
	}
	manager.ReadToEnd(func() {
		self.c.OnUIThread(func() error { return self.showJumpToFileMenu(pane, title) })
	})
	return nil
}

// showJumpToFileMenu offers the diff's files by the paths git names them by. Each item
// names its file rather than the row that file begins at, so that a diff re-rendered
// while the menu is up is jumped into at the row the file begins at now.
//
// A menu offering the one file of a single-file diff would be a menu with nothing to
// choose, so it says what it found instead. It says it here rather than as the key's
// disabled reason because how many files there are is only known once the diff has been
// read to the end, which is too much to do for every keypress that asks whether a key
// applies.
func (self *DiffLineHelper) showJumpToFileMenu(pane types.DiffPaneContext, title string) error {
	view := pane.GetView()
	files := self.FilesInDiff(view)
	if len(files) == 0 {
		return nil
	}
	if len(files) == 1 {
		self.c.ErrorToast(self.c.Tr.DisabledMenuItemPrefix + self.c.Tr.OnlyOneFileInDiff)
		return nil
	}

	worktreePath := self.c.Git().RepoPaths.WorktreePath()
	menuItems := lo.Map(files, func(path string, _ int) *types.MenuItem {
		label := repoRelativePath(worktreePath, path)
		if label == "" {
			label = path
		}
		return &types.MenuItem{
			Label: label,
			OnPress: func() error {
				if target, ok := self.StartOfFileInDiff(view, path); ok {
					self.PlaceNavigationTarget(pane, target, true)
				}
				return nil
			},
		}
	})

	return self.c.Menu(types.CreateMenuOptions{
		Title:           title,
		Items:           menuItems,
		FilterAsYouType: true,
	})
}
