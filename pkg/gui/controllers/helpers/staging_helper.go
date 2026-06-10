package helpers

import (
	"path/filepath"
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/patch_exploring"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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

	// A click/cursor lands on a (wrapped) view line; resolve it to the unwrapped
	// buffer line all three backends key off, then read that buffer line's content.
	bufferLineIdx, ok := v.BufferLineForViewLine(viewLineIdx)
	if !ok {
		return types.DiffLineInfo{}, false
	}
	return self.diffLineInfoFromContents(v.DiffLineContents(), bufferLineIdx)
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
