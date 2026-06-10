package helpers

import (
	"path/filepath"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/tasks"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

var lazygitEditURLRegexp = regexp.MustCompile(`^lazygit-edit://(.+?):(\d+)$`)

type StagingHelper struct {
	c            *HelperCommon
	windowHelper *WindowHelper
}

func NewStagingHelper(
	c *HelperCommon,
	windowHelper *WindowHelper,
) *StagingHelper {
	return &StagingHelper{
		c:            c,
		windowHelper: windowHelper,
	}
}

// NOTE: used from outside this file
func (self *StagingHelper) RefreshStagingPanel(focusOpts types.OnFocusOpts) {
	secondaryFocused := self.secondaryStagingFocused()
	mainFocused := self.mainStagingFocused()

	// this method could be called when the staging panel is not being used,
	// in which case we don't want to do anything.
	if !mainFocused && !secondaryFocused {
		return
	}

	mainSelectedLineIdx := -1
	mainSelectedRealLineIdx := -1
	secondarySelectedLineIdx := -1
	secondarySelectedRealLineIdx := -1
	if focusOpts.ClickedViewLineIdx > 0 {
		if secondaryFocused {
			secondarySelectedLineIdx = focusOpts.ClickedViewLineIdx
			secondarySelectedRealLineIdx = focusOpts.ClickedViewRealLineIdx
		} else {
			mainSelectedLineIdx = focusOpts.ClickedViewLineIdx
			mainSelectedRealLineIdx = focusOpts.ClickedViewRealLineIdx
		}
	}

	mainContext := self.c.Contexts().Staging
	secondaryContext := self.c.Contexts().StagingSecondary

	var file *models.File
	node := self.c.Contexts().Files.GetSelected()
	if node != nil {
		file = node.File
	}

	if file == nil || (!file.HasUnstagedChanges && !file.HasStagedChanges) {
		self.handleStagingEscape()
		return
	}

	mainDiff := self.c.Git().WorkingTree.WorktreeFileDiff(file, true, false)
	secondaryDiff := self.c.Git().WorkingTree.WorktreeFileDiff(file, true, true)

	// grabbing locks here and releasing before we finish the function
	// because pushing say the secondary context could mean entering this function
	// again, and we don't want to have a deadlock
	mainContext.GetMutex().Lock()
	secondaryContext.GetMutex().Lock()

	hunkMode := self.c.UserConfig().Gui.UseHunkModeInStagingView
	mainContext.SetState(
		patch_exploring.NewState(mainDiff, mainSelectedLineIdx, mainSelectedRealLineIdx, focusOpts.ClickedViewRealLineIsDeletion, mainContext.GetView(), mainContext.GetState(), hunkMode, focusOpts.SelectLineInDefaultMode),
	)

	secondaryContext.SetState(
		patch_exploring.NewState(secondaryDiff, secondarySelectedLineIdx, secondarySelectedRealLineIdx, focusOpts.ClickedViewRealLineIsDeletion, secondaryContext.GetView(), secondaryContext.GetState(), hunkMode, focusOpts.SelectLineInDefaultMode),
	)

	mainState := mainContext.GetState()
	secondaryState := secondaryContext.GetState()

	mainContent := mainContext.GetContentToRender()
	secondaryContent := secondaryContext.GetContentToRender()

	mainContext.GetMutex().Unlock()
	secondaryContext.GetMutex().Unlock()

	if mainState == nil && secondaryState == nil {
		self.handleStagingEscape()
		return
	}

	if mainState == nil && !secondaryFocused {
		self.c.Context().Push(secondaryContext, focusOpts)
		return
	}

	if secondaryState == nil && secondaryFocused {
		self.c.Context().Push(mainContext, focusOpts)
		return
	}

	if secondaryFocused {
		self.c.Contexts().StagingSecondary.FocusSelection()
	} else {
		self.c.Contexts().Staging.FocusSelection()
	}

	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Staging,
		Main: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(mainContent),
			Title: self.c.Tr.UnstagedChanges,
		},
		Secondary: &types.ViewUpdateOpts{
			Task:  types.NewRenderStringWithoutScrollTask(secondaryContent),
			Title: self.c.Tr.StagedChanges,
		},
	})
}

func (self *StagingHelper) handleStagingEscape() {
	self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
}

func (self *StagingHelper) secondaryStagingFocused() bool {
	return self.c.Context().CurrentStatic().GetKey() == self.c.Contexts().StagingSecondary.GetKey()
}

func (self *StagingHelper) mainStagingFocused() bool {
	return self.c.Context().CurrentStatic().GetKey() == self.c.Contexts().Staging.GetKey()
}

// GetDiffLineInfo recovers the patch-space identity — (file, type, new-line,
// old-line) — of a rendered diff row, given the window showing the diff and the
// (wrapped) view line index. It is the single seam the focused main view and
// patch explorer consumers go through to act on the line the user is pointing
// at, and the strategy behind it is swappable (see diff-line-metadata-notes.md).
//
// It tries three backends in order of fidelity. First, mechanism #2: per-line
// OSC metadata emitted by a patched pager (delta), which carries the side
// directly and so serves the renderings #1 can't parse — delta's default mode,
// --line-numbers, diff-so-fancy. Failing that, mechanism #1: parsing the
// decolorized view buffer, which serves the structure-preserving renderings (no
// pager, git diff --color, delta --color-only, diff-so-fancy --patch). Failing
// that, delta's lazygit-edit:// hyperlinks; the hyperlink can't convey the side,
// so its result is reported as a non-deletion content line.
func (self *StagingHelper) GetDiffLineInfo(windowName string, viewLineIdx int) (types.DiffLineInfo, bool) {
	v, _ := self.c.GocuiGui().View(self.windowHelper.GetViewNameForWindow(windowName))
	if v == nil {
		return types.DiffLineInfo{}, false
	}
	return self.GetDiffLineInfoForView(v, viewLineIdx)
}

// GetDiffLineInfoForView is GetDiffLineInfo against a specific view rather than
// one looked up by window. It is used to read the identity of the line the patch
// explorer currently has selected when escaping back to the focused main view,
// where we hold the explorer's view directly.
func (self *StagingHelper) GetDiffLineInfoForView(v *gocui.View, viewLineIdx int) (types.DiffLineInfo, bool) {
	// A click/cursor lands on a (wrapped) view line; resolve it to the unwrapped
	// buffer line all three backends key off, then read that buffer line's content.
	bufferLineIdx, ok := v.BufferLineForViewLine(viewLineIdx)
	if !ok {
		return types.DiffLineInfo{}, false
	}
	return self.diffLineInfoFromContents(v.DiffLineContents(), bufferLineIdx)
}

// FindDiffLine returns the index of the first buffer line, at or after from,
// whose patch identity matches target (see types.DiffLineInfo.SamePatchLine), or
// -1 if none does. It is the inverse direction of the diff-line primitive:
// instead of resolving the line under a cursor, it scans a rendered diff for a
// known identity. The escape restore uses it to locate, in the focused main view
// it is returning to, the line the patch explorer had selected — scanning the
// re-render's content as it loads off-screen. (from lets the caller resume the
// scan as more content arrives without re-checking lines already scanned; the
// backends still see all of contents, so a buffer-parse match can look back to
// its file/hunk headers.)
func (self *StagingHelper) FindDiffLine(contents []gocui.DiffLineContent, target types.DiffLineInfo, from int) int {
	for i := from; i < len(contents); i++ {
		if info, ok := self.diffLineInfoFromContents(contents, i); ok && info.SamePatchLine(target) {
			return i
		}
	}
	return -1
}

// RestoreFocusedMainViewOnEscape arranges, when escaping a patch explorer back to
// the focused main view it was entered from, for that view to re-render and then
// land on the line the explorer currently has selected. After staging or dropping
// hunks the explorer's selection auto-advances, so the line the user ended up on
// is more useful to return to than the one they entered on — and, since the diff
// has changed, more reliable than replaying a numeric scroll/index.
//
// It reads that line's patch identity from the explorer now (explorerView,
// explorerSelectedLineIdx), then installs a restore that, as the main view
// re-renders, scans the incoming content for the row matching that identity and —
// once it and a screenful below have loaded — swaps in, scrolls there and selects
// it. If the identity can't be read, or the line never turns up in the re-render
// (the content changed out from under it), no restore is applied and the view
// just re-renders normally.
func (self *StagingHelper) RestoreFocusedMainViewOnEscape(explorerView, mainView *gocui.View, explorerSelectedLineIdx int) {
	target, ok := self.GetDiffLineInfoForView(explorerView, explorerSelectedLineIdx)
	if !ok {
		return
	}

	manager := self.c.GetViewBufferManagerForView(mainView)
	if manager == nil {
		return
	}

	// targetBufferLine is found once: either incrementally as the content loads
	// (scanned tracks how far we've scanned, so each line is checked once), or — if
	// the incremental scan can't resolve it — once more on the complete content
	// when we swap in (see Apply).
	targetBufferLine := -1
	scanned := 0
	manager.SetRestoreForNextTask(&tasks.RenderRestore{
		FirstPaintReady: func() bool {
			// Scan the incoming content for the target as it loads, so we can paint
			// at the saved position as soon as it's reachable. This finds the line
			// for backends that resolve a row on its own — OSC metadata, lazygit-edit
			// hyperlinks. The buffer-parse backend can't resolve here: it parses whole
			// hunks against their @@ lengths, and the trailing hunk is incomplete
			// while loading, so the parse is rejected as not-well-formed until the
			// diff is fully read. That case is handled in Apply instead.
			if targetBufferLine == -1 {
				contents := mainView.OffscreenDiffLineContents()
				targetBufferLine = self.FindDiffLine(contents, target, scanned)
				scanned = len(contents)
				if targetBufferLine == -1 {
					return false
				}
			}
			return mainView.OffscreenLineCount() >= targetBufferLine+mainView.InnerHeight()
		},
		Apply: func() {
			// The off-screen render has just been swapped in, so the displayed buffer
			// now holds it. If the incremental scan never found the target — the
			// common case for buffer-parse, which only becomes well-formed once the
			// whole diff has loaded, at which point the swap happens at end of input —
			// scan the now-complete content once more.
			if targetBufferLine == -1 {
				targetBufferLine = self.FindDiffLine(mainView.DiffLineContents(), target, 0)
			}
			// If the target line still isn't there (the content changed and the line
			// is gone), leave the scroll and selection as they are rather than showing
			// a selection on a line that no longer means what it did.
			if targetBufferLine != -1 {
				if viewLine, ok := mainView.ViewLineForBufferLine(targetBufferLine); ok {
					// scrollIntoView centres the line if it's off-screen, and leaves the
					// scroll untouched if it's already visible — so for the common
					// unchanged-content escape (the placeholder is the same content at
					// the same scroll) nothing moves.
					mainView.FocusPoint(0, viewLine, true)
					mainView.Highlight = true
					mainView.HighlightInactive = false
				}
			}
			manager.ClearRestoreForNextTask()
		},
	})
}

// diffLineInfoFromContents recovers the patch-space identity of the buffer line
// at idx within a snapshot of a diff's per-line content (see gocui.DiffLineContent).
// It is the single resolver behind both directions of the diff-line primitive —
// the forward consumers (click/enter/edit/PR, via GetDiffLineInfo on the displayed
// view) and the inverse identity scan that finds a target line in a focused main
// view as it re-renders (escape restore, via the loading off-screen content). It
// tries three backends in order of fidelity:
//
//   - mechanism #2: per-line OSC metadata emitted by a patched pager (delta),
//     which carries the side directly and so serves the renderings #1 can't parse
//     (delta's default mode, --line-numbers, diff-so-fancy);
//   - mechanism #1: parsing the decolorized buffer, which serves structure-
//     preserving renderings (no pager, git diff --color, delta --color-only);
//   - delta's lazygit-edit:// hyperlinks; the hyperlink can't convey the side, so
//     its result is reported as a non-deletion content line.
func (self *StagingHelper) diffLineInfoFromContents(contents []gocui.DiffLineContent, idx int) (types.DiffLineInfo, bool) {
	if idx < 0 || idx >= len(contents) {
		return types.DiffLineInfo{}, false
	}

	if info, ok := self.diffLineInfoFromMetadata(contents[idx].Metadata); ok {
		return info, true
	}
	if info, ok := self.diffLineInfoFromBuffer(contents, idx); ok {
		return info, true
	}
	return self.diffLineInfoFromHyperlink(contents[idx].Hyperlink)
}

// diffLineInfoFromMetadata reads mechanism #2's per-line OSC metadata. The
// payload is positional and ';'-delimited — version;type;new-line;old-line;file
// — with the file last so it may itself contain ';', and old-line empty unless
// the line is a deletion. See diff-line-metadata-notes.md §9.2.
func (self *StagingHelper) diffLineInfoFromMetadata(payload string) (types.DiffLineInfo, bool) {
	if payload == "" {
		return types.DiffLineInfo{}, false
	}

	parsed, ok := parseDiffLineMetadata(payload)
	if !ok {
		return types.DiffLineInfo{}, false
	}

	// The pager may emit an absolute or a repo-relative path (whichever is
	// convenient for it); normalize to the absolute path the consumers expect.
	path := parsed.RelPath
	if !filepath.IsAbs(path) {
		path = filepath.Join(self.c.Git().RepoPaths.WorktreePath(), path)
	}

	return types.DiffLineInfo{
		Path:    path,
		Type:    parsed.Type,
		NewLine: parsed.NewLine,
		OldLine: parsed.OldLine,
	}, true
}

func (self *StagingHelper) diffLineInfoFromBuffer(contents []gocui.DiffLineContent, idx int) (types.DiffLineInfo, bool) {
	texts := make([]string, len(contents))
	for i, c := range contents {
		texts[i] = c.Text
	}

	parsed, ok := parseDiffLineFromBuffer(texts, idx)
	if !ok {
		return types.DiffLineInfo{}, false
	}

	return types.DiffLineInfo{
		Path:    filepath.Join(self.c.Git().RepoPaths.WorktreePath(), parsed.RelPath),
		Type:    parsed.Type,
		NewLine: parsed.NewLine,
		OldLine: parsed.OldLine,
	}, true
}

func (self *StagingHelper) diffLineInfoFromHyperlink(hyperlink string) (types.DiffLineInfo, bool) {
	matches := lazygitEditURLRegexp.FindStringSubmatch(hyperlink)
	if matches == nil {
		return types.DiffLineInfo{}, false
	}

	return types.DiffLineInfo{
		// delta emits an absolute path here, which is what the consumers want.
		Path: matches[1],
		// The hyperlink carries no side, so it can't distinguish a deletion from
		// an addition or context line; report it as a plain content line.
		Type:    types.DiffLineOther,
		NewLine: utils.MustConvertToInt(matches[2]),
	}, true
}
