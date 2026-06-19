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

// findResolvedDiffLine returns the index of the first line, at or after from, that
// match accepts as the target, or -1 if none does. It is the inverse direction of the
// diff-line primitive: instead of resolving the line under a cursor, it scans a
// rendered diff (resolved once via resolveDiffLines) for a known identity. The
// position restore uses it to locate, in the re-rendered view, the line it wants to
// land on; match is how that restore decides a row is the target (see diffLineMatch).
func findResolvedDiffLine(resolved []resolvedDiffLine, target types.DiffLineInfo, match diffLineMatch, from int) int {
	for i := from; i < len(resolved); i++ {
		if resolved[i].ok && match(target, resolved[i].info) {
			return i
		}
	}
	return -1
}

// diffLineMatch decides whether a re-rendered row is the target a restore is looking
// for. matchByPatchLine is currently the only one — the post-stage reveal, the other
// case where a line's number can change across the re-render, matches by ordinal
// instead (see RevealSelectionAfterStaging).
type diffLineMatch func(target, row types.DiffLineInfo) bool

// matchByPatchLine matches by source-line number (DiffLineInfo.SamePatchLine). It is
// for restores whose target keeps its source-line number across the re-render — the
// escape restore and the context-size preserve, which re-render the same
// staged/unstaged state, so only context lines (not change lines) can come or go.
func matchByPatchLine(target, row types.DiffLineInfo) bool {
	return row.SamePatchLine(target)
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

	self.restoreDiffLinePositionOnRerender(mainView, []diffLineAnchor{{identity: target}}, matchByPatchLine, func(_ diffLineAnchor, viewLine int) {
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
// "advance to the next change" the staging view does. Install it before triggering the
// re-render: it rides the next render of targetView, like RestoreFocusedMainViewOnEscape.
// place re-selects the landed line in the current select mode.
//
// It works by preserving the selection's ordinal among change lines, exactly as the
// staging view preserves its patch-line index: the acted-on line's ordinal is read
// from sourceView (the pane the user acted in) before the op, and after the re-render
// the change line at that ordinal in targetView is selected. Because the op removes the
// acted-on change line(s), that ordinal then holds the next surviving change — the next
// line of the same block when one line was staged, the next block when a whole block
// was. This sidesteps matching by line number, which can't identify the target across a
// stage: a deletion shares its new-file number with the rest of its block, and in the
// staged pane the new-file number is the index, which the op shifts.
//
// sourceView and targetView are usually the same pane, but staging or unstaging can move
// the acted-on side to the other pane (e.g. unstaging the first hunk of an only-staged
// file splits it, pushing the staged remainder to the secondary half); the ordinal is
// preserved across that move since both panes show the same side.
func (self *StagingHelper) RevealSelectionAfterStaging(sourceView *gocui.View, targetView *gocui.View, firstLine int, place func(viewLine int)) {
	ordinal, ok := self.changeLineOrdinal(sourceView, firstLine)
	if !ok {
		return
	}
	self.revealChangeLineAtOrdinal(targetView, ordinal, place)
}

// changeLineOrdinal returns how many change lines precede the (change) line at
// viewLine in view's displayed diff — i.e. that line's index among the diff's change
// lines. ok is false when viewLine maps to no buffer line.
func (self *StagingHelper) changeLineOrdinal(view *gocui.View, viewLine int) (int, bool) {
	bufferLine, ok := view.BufferLineForViewLine(viewLine)
	if !ok {
		return 0, false
	}
	resolved := self.resolveDiffLines(view.DiffLineContents())
	ordinal := 0
	for i := 0; i < bufferLine && i < len(resolved); i++ {
		if resolved[i].ok && resolved[i].info.IsChange() {
			ordinal++
		}
	}
	return ordinal, true
}

// diffLineAnchor is a candidate line for a restore to land on after a re-render: a
// patch identity to find the line by, plus the screen row (offset from the view's
// top) it was on when captured, for a restore that wants to put it back there.
type diffLineAnchor struct {
	identity types.DiffLineInfo
	row      int
}

// installDiffLineRestore is the shared core of the position restores. It sets a
// restore on view's render manager that, as view next re-renders, locates a target
// buffer line and places it: findEarly is given the rows that have loaded since the
// previous call (resolvable on their own via OSC metadata or a lazygit-edit hyperlink)
// so the target can be found and painted early when possible; findComplete is given
// the whole resolved diff at the EOF swap for the definitive target (and the
// buffer-parse fallback, which only resolves once whole hunks have loaded). Each
// returns the anchor to hand back to place — the matched candidate for an identity
// restore, a zero anchor for a positional one — and the target buffer line. place
// receives that anchor and the target's view line; if neither finder produces a
// target, place isn't called and the view re-renders normally.
//
// The find runs against the off-screen buffer before the swap, so the (possibly
// whole-diff) scan happens while the previous content is still displayed; otherwise
// the new content would be drawn at the old scroll for the scan's duration.
func (self *StagingHelper) installDiffLineRestore(
	view *gocui.View,
	findEarly func(rows []gocui.DiffLineContent, offset int) (diffLineAnchor, int, bool),
	findComplete func(resolved []resolvedDiffLine) (diffLineAnchor, int, bool),
	place func(anchor diffLineAnchor, viewLine int),
) {
	// Get-or-create: the target pane may not have rendered yet (the secondary half
	// when a stage/unstage first splits the diff), so the restore has to be set on a
	// manager that the upcoming render will then reuse.
	manager := self.c.GetOrCreateViewBufferManagerForView(view)
	if manager == nil {
		return
	}

	// scanned tracks how far findEarly has looked, so each line is checked once
	// (keeping the incremental scan O(n) over the whole load).
	primaryAnchor := diffLineAnchor{}
	primaryBufferLine := -1
	scanned := 0
	manager.SetRestoreForNextTask(&tasks.RenderRestore{
		FirstPaintReady: func() bool {
			if primaryBufferLine == -1 {
				newRows := view.OffscreenDiffLineContentsFrom(scanned)
				if anchor, line, ok := findEarly(newRows, scanned); ok {
					primaryAnchor, primaryBufferLine = anchor, line
				}
				scanned += len(newRows)
				if primaryBufferLine == -1 {
					return false
				}
			}
			return view.OffscreenLineCount() >= primaryBufferLine+view.InnerHeight()
		},
		Apply: func(swapIn func()) {
			anchor, bufferLine := primaryAnchor, primaryBufferLine
			if bufferLine == -1 {
				if a, line, ok := findComplete(self.resolveDiffLines(view.OffscreenDiffLineContents())); ok {
					anchor, bufferLine = a, line
				}
			}

			swapIn()

			if bufferLine != -1 {
				if viewLine, ok := view.ViewLineForBufferLine(bufferLine); ok {
					place(anchor, viewLine)
				}
			}
			manager.ClearRestoreForNextTask()
		},
	})
}

// restoreDiffLinePositionOnRerender installs a restore that, as view next re-renders,
// lands on a row that match accepts as one of the given candidate identities, and calls
// place with that candidate and its view line. candidates are in priority order
// (nearest first); the restore lands on the first the re-render still contains, found
// incrementally for the nearest and at the EOF swap for any farther fallback (candidates
// aren't in load order, so a nearer one can load after a farther one and we mustn't
// commit prematurely). If none turns up, place isn't called. match decides what "still
// contains" means, since the stable notion of identity differs by what changed in the
// re-render (see diffLineMatch).
//
// It is the identity-matched restore behind escaping back to the focused main view (one
// candidate — the line the patch explorer had selected) and preserving a diff view's
// position across a -U context-size change (several candidates around the anchor, placed
// back where they were; see PreserveDiffPositionOnRerender). The post-stage reveal uses
// the positional restore instead (see revealChangeLineAtOrdinal).
func (self *StagingHelper) restoreDiffLinePositionOnRerender(view *gocui.View, candidates []diffLineAnchor, match diffLineMatch, place func(anchor diffLineAnchor, viewLine int)) {
	if len(candidates) == 0 {
		return
	}
	self.installDiffLineRestore(view,
		func(rows []gocui.DiffLineContent, offset int) (diffLineAnchor, int, bool) {
			for j, content := range rows {
				if info, ok := self.diffLineInfoPerRow(content); ok && match(candidates[0].identity, info) {
					return candidates[0], offset + j, true
				}
			}
			return diffLineAnchor{}, 0, false
		},
		func(resolved []resolvedDiffLine) (diffLineAnchor, int, bool) {
			for _, candidate := range candidates {
				if line := findResolvedDiffLine(resolved, candidate.identity, match, 0); line != -1 {
					return candidate, line, true
				}
			}
			return diffLineAnchor{}, 0, false
		},
		place,
	)
}

// revealChangeLineAtOrdinal installs a restore that, as view next re-renders, selects
// the change line at the given ordinal among the diff's change lines (clamped to the
// last when the re-render has fewer). It is the positional restore behind the post-stage
// reveal: see RevealSelectionAfterStaging for why an ordinal, not an identity.
func (self *StagingHelper) revealChangeLineAtOrdinal(view *gocui.View, ordinal int, place func(viewLine int)) {
	seen := 0 // change lines counted so far by the incremental scan
	self.installDiffLineRestore(view,
		func(rows []gocui.DiffLineContent, offset int) (diffLineAnchor, int, bool) {
			for j, content := range rows {
				if info, ok := self.diffLineInfoPerRow(content); ok && info.IsChange() {
					if seen == ordinal {
						return diffLineAnchor{}, offset + j, true
					}
					seen++
				}
			}
			return diffLineAnchor{}, 0, false
		},
		func(resolved []resolvedDiffLine) (diffLineAnchor, int, bool) {
			last, count := -1, 0
			for i, r := range resolved {
				if r.ok && r.info.IsChange() {
					last = i
					if count == ordinal {
						return diffLineAnchor{}, i, true
					}
					count++
				}
			}
			// Fewer change lines than the ordinal (the acted-on line was at or near the
			// end): clamp to the last surviving change.
			return diffLineAnchor{}, last, last != -1
		},
		func(_ diffLineAnchor, viewLine int) { place(viewLine) },
	)
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

	self.restoreDiffLinePositionOnRerender(view, self.nearbyDiffLines(view, anchorViewLine), matchByPatchLine, func(anchor diffLineAnchor, viewLine int) {
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
