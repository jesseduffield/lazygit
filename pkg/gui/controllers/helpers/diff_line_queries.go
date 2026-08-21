package helpers

import (
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// The questions a diff view can be asked about what it is showing — where the change
// lines are, which block or file a row belongs to — answered in the view-line terms a
// cursor and a click speak. They are all built on the identities recovered in
// diff_line_helper.go, which is where the answering stops and the recovering starts.

// DiffLinesInViewRange returns the identity of every diff line shown by the rows in
// the inclusive view-line range [first, last] of view's rendered diff, in display
// order. Rows whose identity can't be recovered are left out, as are the wrapped
// segments of a row already counted.
//
// A row can show more than one diff line — a side-by-side rendering puts a deletion
// beside the addition replacing it — and all of them are reported: what the user
// pointed at is the row, so everything on it is selected.
func (self *DiffLineHelper) DiffLinesInViewRange(view *gocui.View, first int, last int) []types.DiffLineInfo {
	identities := self.resolveDiffLineIdentities(view.DiffLineContents())

	infos := []types.DiffLineInfo{}
	previousBufferLine := -1
	for viewLine := first; viewLine <= last; viewLine++ {
		bufferLine, ok := view.BufferLineForViewLine(viewLine)
		if !ok || bufferLine == previousBufferLine || bufferLine >= len(identities) {
			continue
		}
		previousBufferLine = bufferLine
		infos = append(infos, identities[bufferLine]...)
	}
	return self.inRepoTerms(view, infos)
}

// ChangeLineOrdinals says, for each of the given change lines, which of its file's
// changes it is in the given diff — its place among them, counted from the top of the
// file — keyed by file. Lines the diff doesn't have are left out.
//
// It is how a line is named in something built out of a diff rather than being that diff:
// the custom patch holds the lines it was given in the order the file has them, so a
// place among a file's changes is a line of the patch.
func (self *DiffLineHelper) ChangeLineOrdinals(
	diff string, infos []types.DiffLineInfo,
) map[string][]int {
	ordinals := map[patchLine]int{}
	counts := map[string]int{}
	for _, parsed := range parseAllDiffLinesFromBuffer(strings.Split(diff, "\n")) {
		if !parsed.ok {
			continue
		}
		info := self.diffLineInfo(parsed.parsed)
		if !info.IsChange() {
			continue
		}
		ordinals[patchLineOf(info)] = counts[info.Path]
		counts[info.Path]++
	}

	ordinalsByPath := map[string][]int{}
	for _, info := range infos {
		if ordinal, ok := ordinals[patchLineOf(info)]; ok {
			ordinalsByPath[info.Path] = append(ordinalsByPath[info.Path], ordinal)
		}
	}
	return ordinalsByPath
}

// inRepoTerms brings the paths of lines recovered from a view into the repo's terms.
//
// They are in them already for a diff of the repo's own files. The pane previewing the
// custom patch, though, shows a diff of the two trees the patch was materialized into: a
// diff renderer states the path it was handed there, which is under the tree's own name,
// while the diff's text names the trees where an ordinary diff has git's a/ and b/
// prefixes and so needs nothing.
func (self *DiffLineHelper) inRepoTerms(view *gocui.View, infos []types.DiffLineInfo) []types.DiffLineInfo {
	if !self.ShowsCustomPatch(view) {
		return infos
	}

	worktreePath := self.c.Git().RepoPaths.WorktreePath()
	treesDir := self.c.Git().Patch.PatchBuilder.TempDir()
	return lo.Map(infos, func(info types.DiffLineInfo, _ int) types.DiffLineInfo {
		info.Path = repoPathOfTreePath(info.Path, treesDir, worktreePath)
		return info
	})
}

// repoPathOfTreePath maps a path under one of the trees the custom patch was materialized
// into to the file of the repo it stands for: the path is the tree's name followed by the
// file's own, stated either against the directory holding the trees or against the repo,
// depending on how the renderer that stated it was given it.
func repoPathOfTreePath(path string, treesDir string, worktreePath string) string {
	root := worktreePath
	if treesDir != "" && strings.HasPrefix(path, treesDir+string(filepath.Separator)) {
		root = treesDir
	}
	relativePath, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}

	segments := strings.Split(filepath.ToSlash(relativePath), "/")
	if len(segments) > 1 && (segments[0] == "a" || segments[0] == "b") {
		relativePath = filepath.Join(segments[1:]...)
	}
	return filepath.Join(worktreePath, relativePath)
}

// ChangeLinesInViewRange returns the change lines — the additions and deletions —
// among the diff lines shown by the rows in the inclusive view-line range. Those are
// the lines a patch is built from: a patch carries whatever context it needs around
// them by itself, so a selection contributes only its changes.
func (self *DiffLineHelper) ChangeLinesInViewRange(view *gocui.View, first int, last int) []types.DiffLineInfo {
	return lo.Filter(self.DiffLinesInViewRange(view, first, last),
		func(info types.DiffLineInfo, _ int) bool { return info.IsChange() })
}

// changeLines resolves view's rendered diff to one flag per buffer line: whether
// that row is a change line (an addition or a deletion), as opposed to context, a
// header, or a row whose identity couldn't be recovered. Those are the rows a
// selection is anchored on and navigation moves between.
func (self *DiffLineHelper) changeLines(view *gocui.View) []bool {
	resolved := self.resolveDiffLines(view.DiffLineContents())
	isChange := make([]bool, len(resolved))
	for i, r := range resolved {
		isChange[i] = r.ok && r.info.IsChange()
	}
	return isChange
}

// FirstChangeLineInView returns the view line of the first change line on screen. It
// is where the selection goes when the main view is focused by keyboard: focusing a
// diff you are reading points at something in it without moving it, so the search
// stops at the bottom of the viewport rather than going after a change further down.
// ok is false when the viewport holds no change line — scrolled into a long stretch
// of context, or past the last change.
func (self *DiffLineHelper) FirstChangeLineInView(view *gocui.View) (int, bool) {
	top, bottom, ok := visibleBufferLines(view)
	if !ok {
		return 0, false
	}

	isChange := self.changeLines(view)
	for i := top; i <= min(bottom, len(isChange)-1); i++ {
		if isChange[i] {
			return view.ViewLineForBufferLine(i)
		}
	}
	return 0, false
}

// FirstChangeBlockInView returns the view line of the first change block on screen:
// the first one that *begins* in the viewport, and failing that the one that reaches
// into the viewport from above, whose start is off screen. That order is what hunk
// mode wants of the block it offers up on focus — a block whose beginning the user can
// see, rather than the tail of one they have scrolled past the start of, with the one
// bleeding in from above kept as the answer for a change too long to fit on screen,
// where there is no other. ok is false when the viewport shows no change line.
func (self *DiffLineHelper) FirstChangeBlockInView(view *gocui.View) (int, bool) {
	top, bottom, ok := visibleBufferLines(view)
	if !ok {
		return 0, false
	}

	isChange := self.changeLines(view)
	for i := top; i <= min(bottom, len(isChange)-1); i++ {
		if isChange[i] && (i == 0 || !isChange[i-1]) {
			return view.ViewLineForBufferLine(i)
		}
	}
	// A block covering the top line is one that began above it: nothing else can put a
	// change there once no block starts on screen.
	if top < len(isChange) && isChange[top] {
		return view.ViewLineForBufferLine(top)
	}
	return 0, false
}

// visibleBufferLines returns the first and last line of view's content that the
// viewport shows any part of, for the queries that only care about what the user can
// see. The last line is the one at the bottom edge, or the content's last when the
// content ends above it. ok is false for a view showing no content at all.
func visibleBufferLines(view *gocui.View) (int, int, bool) {
	top, ok := view.BufferLineForViewLine(view.OriginY())
	if !ok {
		return 0, 0, false
	}

	lastVisible := min(view.OriginY()+view.InnerHeight(), view.ViewLinesHeight()) - 1
	bottom, ok := view.BufferLineForViewLine(lastVisible)
	if !ok {
		return top, top, true
	}
	return top, bottom, true
}

// ViewHasChangeLines reports whether view's rendered diff holds any change line at
// all, i.e. whether there is anything to select. It is false over a non-diff
// placeholder, and over a diff with nothing in it — an empty commit, a binary file —
// which are the cases where the focused main view shows no selection.
func (self *DiffLineHelper) ViewHasChangeLines(view *gocui.View) bool {
	return lo.Contains(self.changeLines(view), true)
}

// IsChangeLine reports whether the given view line of view's rendered diff is a
// change line rather than context, a header, or an unresolvable row — i.e. whether
// pointing at it points at something a patch could be built from.
func (self *DiffLineHelper) IsChangeLine(view *gocui.View, viewLineIdx int) bool {
	info, ok := self.GetDiffLineInfo(view, viewLineIdx)
	return ok && info.IsChange()
}

// IsSingleHunkForWholeFile reports whether the file the given change line belongs to
// is shown as one solid block of changes — every row of its diff a change of the same
// kind, no context — which is what a newly added or deleted file looks like. That is
// the case where widening the selection to the change block would select the file
// entire, so hunk mode drops to a single line there instead. It asks of a rendered
// diff the question patch.Patch.IsSingleHunkForWholeFile asks of a patch.
//
// It says false while the diff is still being read in, since the rows that would
// answer otherwise — a context line, a change of the other kind — may not have
// arrived yet. That errs towards hunk mode, which is what the user asked for.
func (self *DiffLineHelper) IsSingleHunkForWholeFile(view *gocui.View, changeViewLine int) bool {
	if manager := self.c.GetViewBufferManagerForView(view); manager != nil && manager.IsLoading() {
		return false
	}

	anchor, ok := view.BufferLineForViewLine(changeViewLine)
	if !ok {
		return false
	}
	resolved := self.resolveDiffLines(view.DiffLineContents())
	if anchor >= len(resolved) || !resolved[anchor].ok {
		return false
	}

	// The question is per file: a commit's diff may hold a newly added file next to an
	// edited one.
	path := resolved[anchor].info.Path
	kind := resolved[anchor].info.Type
	for _, row := range resolved {
		if !row.ok || row.info.Path != path {
			continue
		}
		if row.info.Type == types.DiffLineContext {
			return false
		}
		if row.info.IsChange() && row.info.Type != kind {
			return false
		}
	}
	return true
}

// ChangeBlockBounds returns the inclusive view-line range of the change block to
// select in hunk mode around anchorViewLine. A change block is lazygit's notion of a
// hunk — a run of consecutive added or deleted lines bounded by context, of which a
// single git @@ hunk may hold several. When the anchor is context, the block used is
// the first at or below it, or — with nothing below, the cursor sitting past the last
// change — the nearest above, so that hunk mode always has a block to select. ok is
// false only when the diff holds no change line at all.
func (self *DiffLineHelper) ChangeBlockBounds(view *gocui.View, anchorViewLine int) (int, int, bool) {
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, 0, false
	}

	isChange := self.changeLines(view)
	start := anchor
	for start < len(isChange) && !isChange[start] {
		start++
	}
	if start >= len(isChange) {
		for start = min(anchor, len(isChange)-1); start >= 0 && !isChange[start]; start-- {
		}
		if start < 0 {
			return 0, 0, false
		}
	}
	end := start
	for start > 0 && isChange[start-1] {
		start--
	}
	for end < len(isChange)-1 && isChange[end+1] {
		end++
	}

	startView, startOk := view.ViewLineForBufferLine(start)
	// The block's last line goes to its last view line, so that a line the view
	// wrapped is highlighted to its end rather than only where it begins.
	endView, endOk := view.LastViewLineForBufferLine(end)
	if !startOk || !endOk {
		return 0, 0, false
	}
	return startView, endView, true
}

// AdjacentChangeBlock returns the view line to move to for next/previous change-block
// navigation in view's rendered diff, starting from anchorViewLine. A change block is
// lazygit's notion of a hunk (see ChangeBlockBounds). forward=true targets the start
// of the next block, forward=false the start of the previous one — from mid-block that
// means the previous block, rather than the one we are in. ok is false when there's no
// further block, so the caller leaves the view where it is.
func (self *DiffLineHelper) AdjacentChangeBlock(view *gocui.View, anchorViewLine int, forward bool) (int, bool) {
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, false
	}

	target, ok := changeBlockStart(self.changeLines(view), anchor, forward)
	if !ok {
		return 0, false
	}
	return view.ViewLineForBufferLine(target)
}

// AdjacentFile returns the view line to move to for next/previous file navigation in
// view's (possibly multi-file) rendered diff, starting from anchorViewLine: the first
// located row of the neighbouring file, found where the rows' file changes. ok is
// false at the first or last file.
func (self *DiffLineHelper) AdjacentFile(view *gocui.View, anchorViewLine int, forward bool) (int, bool) {
	anchor, ok := view.BufferLineForViewLine(anchorViewLine)
	if !ok {
		return 0, false
	}

	target, ok := fileStart(self.filePaths(view), anchor, forward)
	if !ok {
		return 0, false
	}
	return view.ViewLineForBufferLine(target)
}

// filePaths resolves view's rendered diff to the path each buffer line belongs to,
// empty for a row whose identity couldn't be recovered.
func (self *DiffLineHelper) filePaths(view *gocui.View) []string {
	resolved := self.resolveDiffLines(view.DiffLineContents())
	paths := make([]string, len(resolved))
	for i, row := range resolved {
		if row.ok {
			paths[i] = row.info.Path
		}
	}
	return paths
}

// fileStart finds, in a diff whose lines carry the file path they belong to (empty for
// a row no backend could place), the first located row of the file adjacent to `from`
// in the given direction — the row file navigation lands on. It is the pure index
// arithmetic behind AdjacentFile.
//
// A file is identified by its path, so we look for where the path changes, skipping
// rows that carry none: those are the blank separator rows between files, or the
// header rows of a diff renderer that doesn't state which file its headers belong to.
// So the landing row is the file's header wherever the source says so — a parseable
// buffer, or a renderer that tags its headers — and the file's first content line
// otherwise, which is an accepted degradation.
func fileStart(paths []string, from int, forward bool) (int, bool) {
	anchorPath, ok := anchorFilePath(paths, from)
	if !ok {
		return 0, false
	}

	if forward {
		for i := from; i < len(paths); i++ {
			if paths[i] != "" && paths[i] != anchorPath {
				return i, true
			}
		}
		return 0, false
	}

	// Walk back past the current file (its rows and any unlocated ones) to the previous
	// file's last located row, then back over that whole file, landing on its first.
	i := from
	for i >= 0 && (paths[i] == "" || paths[i] == anchorPath) {
		i--
	}
	if i < 0 {
		return 0, false
	}
	prevPath := paths[i]
	for i > 0 && (paths[i-1] == "" || paths[i-1] == prevPath) {
		i--
	}
	for paths[i] != prevPath {
		i++
	}
	return i, true
}

// anchorFilePath returns the path of the file the anchor sits in: the first row at or
// below it that carries a path — the file whose content is at or below the top of the
// view — falling back to the nearest above when there is nothing below. Scanning down
// first matters because the anchor is often a file-header row that carries no path of
// its own, whose nearest tagged row above is the *previous* file's content; taking
// that would make next-file navigation jump back into the file just left, so a second
// press wouldn't advance. ok is false when no row carries a path.
func anchorFilePath(paths []string, from int) (string, bool) {
	if from < 0 {
		return "", false
	}
	for i := from; i < len(paths); i++ {
		if paths[i] != "" {
			return paths[i], true
		}
	}
	for i := min(from, len(paths)) - 1; i >= 0; i-- {
		if paths[i] != "" {
			return paths[i], true
		}
	}
	return "", false
}

// changeBlockStart finds, in a diff whose lines are flagged by isChange, the first
// line of the change block adjacent to `from` in the given direction. It is the pure
// index arithmetic behind AdjacentChangeBlock.
func changeBlockStart(isChange []bool, from int, forward bool) (int, bool) {
	if from < 0 || from >= len(isChange) {
		return 0, false
	}

	if forward {
		i := from
		for i < len(isChange) && isChange[i] { // leave the current block
			i++
		}
		for i < len(isChange) && !isChange[i] { // skip the separating context
			i++
		}
		if i == len(isChange) {
			return 0, false
		}
		return i, true
	}

	i := from
	for i >= 0 && isChange[i] { // leave the current block
		i--
	}
	for i >= 0 && !isChange[i] { // skip context, landing on the previous block's last line
		i--
	}
	if i < 0 {
		return 0, false
	}
	for i > 0 && isChange[i-1] { // walk back to that block's first line
		i--
	}
	return i, true
}
