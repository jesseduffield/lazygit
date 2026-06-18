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

// ChangeLinesInViewRange resolves every change line (addition/deletion) in the
// inclusive view-line range [first, last] of the diff shown in the given window, as
// the patch-space identities to stage. It is the range form of GetDiffLineInfo,
// behind staging a selection from the focused main view. Context and header lines
// are skipped (Transform emits context regardless of the included set, so only
// change lines need collecting — see §21.3), and view lines that wrap to the same
// buffer line are de-duplicated.
//
// A side-by-side rendering carries more than one record per row — a deletion on the
// left, the addition replacing it on the right — and staging includes both (you
// can't stage one side of a side-by-side row; accepted restriction). So each row's
// metadata payloads are all resolved; rows without metadata (no pager, or the
// buffer-parse / hyperlink backends) fall back to their single resolved record.
func (self *StagingHelper) ChangeLinesInViewRange(windowName string, first int, last int) []types.DiffLineInfo {
	v, _ := self.c.GocuiGui().View(self.windowHelper.GetViewNameForWindow(windowName))
	if v == nil {
		return nil
	}

	resolved := self.resolveDiffLines(v.DiffLineContents())
	payloadsByLine := v.DiffLineMetadataPayloads()
	var infos []types.DiffLineInfo
	lastBufferLine := -1
	for viewLine := first; viewLine <= last; viewLine++ {
		bufferLine, ok := v.BufferLineForViewLine(viewLine)
		if !ok || bufferLine == lastBufferLine {
			continue
		}
		lastBufferLine = bufferLine

		if bufferLine < len(payloadsByLine) && len(payloadsByLine[bufferLine]) > 0 {
			for _, payload := range payloadsByLine[bufferLine] {
				if info, ok := self.diffLineInfoFromMetadata(payload); ok && info.IsChange() {
					infos = append(infos, info)
				}
			}
		} else if bufferLine < len(resolved) && resolved[bufferLine].ok && resolved[bufferLine].info.IsChange() {
			infos = append(infos, resolved[bufferLine].info)
		}
	}
	return infos
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

// findResolvedDiffLine returns the index of the first line, at or after from, whose
// patch identity matches target (see types.DiffLineInfo.SamePatchLine), or -1 if
// none does. It is the inverse direction of the diff-line primitive: instead of
// resolving the line under a cursor, it scans a rendered diff (resolved once via
// resolveDiffLines) for a known identity. The position restore uses it to locate,
// in the re-rendered view, the line it wants to land on.
func findResolvedDiffLine(resolved []resolvedDiffLine, target types.DiffLineInfo, from int) int {
	for i := from; i < len(resolved); i++ {
		if resolved[i].ok && resolved[i].info.SamePatchLine(target) {
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

	self.restoreDiffLinePositionOnRerender(mainView, []diffLineAnchor{{identity: target}}, func(_ diffLineAnchor, viewLine int) {
		// scrollIntoView centres the line if it's off-screen, and leaves the scroll
		// untouched if it's already visible — so for the common unchanged-content
		// escape (the placeholder is the same content at the same scroll) nothing
		// moves.
		mainView.FocusPoint(0, viewLine, true)
		mainView.Highlight = true
		mainView.HighlightInactive = false
	})
}

// RevealSelectionAfterStaging arranges, after staging or unstaging re-renders the
// focused main view's diff, for the selection to land on the change nearest the one
// just acted on rather than at a stale (and possibly off-content) position — the same
// "advance to the next change" the staging view does. firstLine and lastLine bound
// the just-staged selection in the current (pre-staging) diff; place positions and
// re-selects the landed line in the current select mode. Install it before triggering
// the re-render: it rides the next render of view, like RestoreFocusedMainViewOnEscape.
//
// The candidates, in priority order, are the change block after the selection, then
// the one before it (both survive staging — only the selection itself was staged),
// then the selection's own first line, which only turns up when the whole side was
// staged and the view flips to show the other (staged) side. If none survives, place
// isn't called and the view re-renders without moving the selection.
func (self *StagingHelper) RevealSelectionAfterStaging(view *gocui.View, firstLine int, lastLine int, place func(viewLine int)) {
	var candidates []diffLineAnchor
	addByViewLine := func(viewLine int) {
		if info, ok := self.GetDiffLineInfoForView(view, viewLine); ok {
			candidates = append(candidates, diffLineAnchor{identity: info})
		}
	}

	if next, ok := self.AdjacentChangeBlock(view, lastLine, true); ok {
		addByViewLine(next)
	}
	if prev, ok := self.AdjacentChangeBlock(view, firstLine, false); ok {
		addByViewLine(prev)
	}
	addByViewLine(firstLine)

	self.restoreDiffLinePositionOnRerender(view, candidates, func(_ diffLineAnchor, viewLine int) {
		place(viewLine)
	})
}

// diffLineAnchor is a candidate line for a restore to land on after a re-render: a
// patch identity to find the line by, plus the screen row (offset from the view's
// top) it was on when captured, for a restore that wants to put it back there.
type diffLineAnchor struct {
	identity types.DiffLineInfo
	row      int
}

// restoreDiffLinePositionOnRerender installs a restore on view's render manager so
// that, as view next re-renders, it lands on a row matching one of the given
// candidate identities, and calls place with that candidate and its view line to
// position and/or select it. candidates are in priority order (nearest first); the
// restore lands on the first one that the re-render still contains. If none turns up
// (the content changed out from under all of them) place is not called and the view
// just re-renders normally.
//
// It is the context-neutral core behind both restoring the focused main view on
// escape (one candidate — the line the patch explorer had selected — placed by
// scrolling to and selecting it) and preserving a diff view's position when its -U
// context size changes (several candidates around the anchor, placed back where they
// were; see PreserveDiffPositionOnRerender).
func (self *StagingHelper) restoreDiffLinePositionOnRerender(view *gocui.View, candidates []diffLineAnchor, place func(anchor diffLineAnchor, viewLine int)) {
	if len(candidates) == 0 {
		return
	}
	manager := self.c.GetViewBufferManagerForView(view)
	if manager == nil {
		return
	}

	// The primary (nearest) candidate is found incrementally as the content loads
	// (scanned tracks how far we've scanned, so each line is checked once), letting
	// us paint early when it survives — the common case. A fallback to a farther
	// candidate is resolved only on the complete content at the EOF swap (see Apply),
	// because candidates aren't in load order, so a nearer one can load after a
	// farther one and we mustn't commit to the farther one prematurely.
	primaryBufferLine := -1
	scanned := 0
	manager.SetRestoreForNextTask(&tasks.RenderRestore{
		FirstPaintReady: func() bool {
			// Scan the rows that have loaded since we last looked for the primary
			// candidate, resolving each on its own (diffLineInfoPerRow: OSC metadata
			// or a lazygit-edit hyperlink). The buffer-parse backend can't resolve
			// here — it parses whole hunks against their @@ lengths, and the trailing
			// hunk is incomplete while loading, so the parse is rejected until the diff
			// is fully read; that (and any fallback) is handled in Apply. Scanning only
			// the new rows each time keeps this O(n) over the whole load.
			if primaryBufferLine == -1 {
				newRows := view.OffscreenDiffLineContentsFrom(scanned)
				for j, content := range newRows {
					if info, ok := self.diffLineInfoPerRow(content); ok && info.SamePatchLine(candidates[0].identity) {
						primaryBufferLine = scanned + j
						break
					}
				}
				scanned += len(newRows)
				if primaryBufferLine == -1 {
					return false
				}
			}
			return view.OffscreenLineCount() >= primaryBufferLine+view.InnerHeight()
		},
		Apply: func(swapIn func()) {
			// Find the candidate to land on, then swap the off-screen content in and
			// place it. The find runs against the *off-screen* buffer, before the swap,
			// so the (possibly whole-diff) scan happens while the previous content is
			// still displayed — otherwise the new content would be drawn at the old
			// scroll for the duration of the scan. Land on the nearest candidate the
			// content still contains: the primary one if the incremental scan found it,
			// otherwise scan the now-complete content in priority order (the common path
			// for buffer-parse, which only becomes well-formed once the whole diff has
			// loaded, so its swap is at end of input and the off-screen buffer is whole).
			matched, bufferLine := -1, primaryBufferLine
			if bufferLine != -1 {
				matched = 0
			} else {
				resolved := self.resolveDiffLines(view.OffscreenDiffLineContents())
				for i, candidate := range candidates {
					if line := findResolvedDiffLine(resolved, candidate.identity, 0); line != -1 {
						matched, bufferLine = i, line
						break
					}
				}
			}

			swapIn()

			// If no candidate is there (the content changed and they're all gone),
			// leave the scroll and selection as they are rather than acting on a line
			// that no longer means what it did.
			if matched != -1 {
				if viewLine, ok := view.ViewLineForBufferLine(bufferLine); ok {
					place(candidates[matched], viewLine)
				}
			}
			manager.ClearRestoreForNextTask()
		},
	})
}

// PreserveDiffPositionOnRerender remembers where a diff view is anchored and, when
// it next re-renders, restores that position instead of letting it jump to the top.
// It is the diff-line primitive's context-change consumer: the increase/decrease
// context-size keybindings re-render the diff with a different git command, which
// otherwise resets the scroll. Call it on the view about to be re-rendered, right
// before triggering the re-render.
//
// The anchor is the selected line if a selection is showing, otherwise the line in
// the middle of the visible content (so what stays put is the line you're most likely
// looking at, not the top edge). We'd like to keep the anchor line itself, but it may
// not survive the re-render: a context line vanishes when the context size shrinks. So
// we collect the
// lines around the anchor as fallbacks (nearest first; see nearbyDiffLines) and land
// on the nearest one that survives — the anchor itself when it does, otherwise the
// closest surviving line, which minimises scrolling. Each is put back at the same
// screen row it was on. A showing selection is re-established on the landed line;
// otherwise the view stays in scroll mode.
func (self *StagingHelper) PreserveDiffPositionOnRerender(view *gocui.View) {
	// If the view isn't the one currently shown in its window it won't be the one
	// re-rendered (e.g. the merge-conflicts view occupies the main window), so a
	// restore set on it would linger and wrongly suppress a later render's scroll
	// reset. There's nothing to preserve then.
	if !view.Visible {
		return
	}

	showSelection := view.Highlight
	anchorViewLine := view.MiddleVisibleLineIdx()
	if showSelection {
		anchorViewLine = view.SelectedLineIdx()
	}

	self.restoreDiffLinePositionOnRerender(view, self.nearbyDiffLines(view, anchorViewLine), func(anchor diffLineAnchor, viewLine int) {
		// Put the landed line back on the screen row it was captured on, clamped into
		// the view in case it was off-screen (a fallback line can be), so the restore
		// always lands somewhere visible.
		row := max(0, min(anchor.row, view.InnerHeight()-1))
		view.SetOrigin(0, viewLine-row)
		if showSelection {
			// SetOrigin already placed the row, so don't scroll again; just move the
			// cursor onto it and turn the selection back on.
			view.FocusPoint(0, viewLine, false)
			view.Highlight = true
			view.HighlightInactive = false
		}
	})
}

// nearbyDiffLines collects the resolvable diff lines around the given view line in
// view's displayed diff, as restore candidates ordered by proximity (the anchor line
// itself first, then outward, preferring at-or-below on ties), each tagged with the
// screen row it was on. Expansion in each direction stops once it reaches a change
// line (addition/deletion), which always survives a context-size change — so the list
// always ends in a guaranteed survivor, while still offering the nearer context lines
// as preferred candidates. The restore (see restoreDiffLinePositionOnRerender) lands
// on the first of these the re-render still contains.
func (self *StagingHelper) nearbyDiffLines(view *gocui.View, anchorViewLine int) []diffLineAnchor {
	contents := view.DiffLineContents()
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return nil
	}
	resolved := self.resolveDiffLines(contents)
	originY := view.OriginY()

	var anchors []diffLineAnchor
	// collect appends the line at bufferLine if it resolves, and reports whether it's
	// a change line (a guaranteed survivor, so expansion past it isn't needed).
	collect := func(bufferLine int) (atChangeLine bool) {
		r := resolved[bufferLine]
		if !r.ok {
			return false
		}
		if viewLine, ok := view.ViewLineForBufferLine(bufferLine); ok {
			anchors = append(anchors, diffLineAnchor{identity: r.info, row: viewLine - originY})
		}
		return r.info.IsChange()
	}

	if collect(anchor) {
		return anchors
	}
	below, above := anchor+1, anchor-1
	for below < len(contents) || above >= 0 {
		if below < len(contents) {
			if collect(below) {
				below = len(contents) // stop expanding past the first change line below
			} else {
				below++
			}
		}
		if above >= 0 {
			if collect(above) {
				above = -1 // stop expanding past the first change line above
			} else {
				above--
			}
		}
	}
	return anchors
}

// resolvedDiffLine is one line's resolved patch identity plus whether it could be
// resolved, the element of the table resolveDiffLines produces.
type resolvedDiffLine struct {
	info types.DiffLineInfo
	ok   bool
}

// resolveDiffLines recovers the patch-space identity of every line of a diff's
// per-line content (see gocui.DiffLineContent) in one pass, indexed 1:1 with
// contents. It is the batch form of the inverse diff-line primitive, used by the
// whole-buffer scans — the position restore (nearbyDiffLines / Apply) and the
// focused main view's navigation. Resolving line-by-line would re-run the buffer-
// parse backend's whole-section parse once per line (O(n²) on a large single-file
// diff); this parses the buffer-parse backend once for the whole buffer
// (parseAllDiffLinesFromBuffer) and applies the per-line backend precedence
// (resolveDiffLine) on top.
func (self *StagingHelper) resolveDiffLines(contents []gocui.DiffLineContent) []resolvedDiffLine {
	bufferParsed := parseAllDiffLinesFromBuffer(diffLineTexts(contents))
	resolved := make([]resolvedDiffLine, len(contents))
	for i, content := range contents {
		info, ok := self.resolveDiffLine(content, bufferParsed[i])
		resolved[i] = resolvedDiffLine{info, ok}
	}
	return resolved
}

// diffLineInfoFromContents recovers the patch-space identity of a single buffer
// line at idx (see resolveDiffLine for the backend precedence). It is the forward
// consumers' resolver (click/enter/edit/PR, via GetDiffLineInfo on the displayed
// view); to resolve every line of a buffer, use resolveDiffLines, which parses the
// buffer-parse backend once rather than once per line.
func (self *StagingHelper) diffLineInfoFromContents(contents []gocui.DiffLineContent, idx int) (types.DiffLineInfo, bool) {
	if idx < 0 || idx >= len(contents) {
		return types.DiffLineInfo{}, false
	}
	parsed, ok := parseDiffLineFromBuffer(diffLineTexts(contents), idx)
	return self.resolveDiffLine(contents[idx], bufferLineParse{parsed, ok})
}

// resolveDiffLine applies the diff-line backends to one rendered row in order of
// fidelity, given the row's content and the buffer-parse backend's already-computed
// result for it:
//
//   - mechanism #2: per-line OSC metadata emitted by a patched pager (delta),
//     which carries the side directly and so serves the renderings #1 can't parse
//     (delta's default mode, --line-numbers, diff-so-fancy);
//   - mechanism #1: parsing the decolorized buffer (bufferParsed), which serves
//     structure-preserving renderings (no pager, git diff --color, delta --color-only);
//   - delta's lazygit-edit:// hyperlinks; the hyperlink can't convey the side, so
//     its result is reported as a non-deletion content line.
func (self *StagingHelper) resolveDiffLine(content gocui.DiffLineContent, bufferParsed bufferLineParse) (types.DiffLineInfo, bool) {
	if info, ok := self.diffLineInfoFromMetadata(content.Metadata); ok {
		return info, true
	}
	if bufferParsed.ok {
		return self.diffLineInfoFromParsed(bufferParsed.parsed), true
	}
	return self.diffLineInfoFromHyperlink(content.Hyperlink)
}

// diffLineInfoPerRow resolves a single rendered row using only the backends that can
// place it from material carried on the row itself — OSC metadata or a lazygit-edit
// hyperlink. It omits the buffer-parse backend, which needs the surrounding
// file/hunk lines and so can't resolve a row while the diff is still streaming in;
// the incremental restore scan uses this, and resolves the buffer-parse case once
// the diff is complete (see restoreDiffLinePositionOnRerender's Apply).
func (self *StagingHelper) diffLineInfoPerRow(content gocui.DiffLineContent) (types.DiffLineInfo, bool) {
	if info, ok := self.diffLineInfoFromMetadata(content.Metadata); ok {
		return info, true
	}
	return self.diffLineInfoFromHyperlink(content.Hyperlink)
}

// diffLineTexts extracts the decolorized text of each row — the input the
// buffer-parse backend (mechanism #1) works on.
func diffLineTexts(contents []gocui.DiffLineContent) []string {
	texts := make([]string, len(contents))
	for i, content := range contents {
		texts[i] = content.Text
	}
	return texts
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

// diffLineInfoFromParsed turns the buffer parser's repo-relative result into the
// absolute-path identity the consumers expect (mechanism #1).
func (self *StagingHelper) diffLineInfoFromParsed(parsed parsedDiffLine) types.DiffLineInfo {
	return types.DiffLineInfo{
		Path:    filepath.Join(self.c.Git().RepoPaths.WorktreePath(), parsed.RelPath),
		Type:    parsed.Type,
		NewLine: parsed.NewLine,
		OldLine: parsed.OldLine,
	}
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
