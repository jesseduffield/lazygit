package controllers

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// togglePatchFromFocusedMainView toggles the selected diff line(s) — a single line, a
// range, or a hunk — into or out of the custom patch being built for (from, to,
// reverse). It is the patch-building counterpart of the files panel's staging handler,
// shared by the panels that build a patch from a focused main view: the commit files
// panel (the per-file diff) and the commits / sub-commits / stash panels (the
// whole-commit diff). The selection is resolved to change-line identities (shared with
// staging) and mapped to the patch builder's per-file line indices.
//
// Building a patch needs the patch builder started for this target; if a patch for a
// different target is active, we confirm before discarding it (as entering the patch
// builder does). The diff itself is unchanged by a toggle, so refresh — supplied by the
// caller, since what needs re-rendering differs per panel — re-renders the same diff
// command (keeping its scroll) and the secondary patch view, and the inclusion gutter
// rides that re-render. That re-render (and its layout re-wrap when the secondary view
// first appears) is in view-line space, which moves the selection, so we re-establish it
// afterwards by its width-independent change-line ordinal.
func togglePatchFromFocusedMainView(
	c *ControllerCommon,
	mainViewName string,
	firstLineIdx int,
	lastLineIdx int,
	from string,
	to string,
	reverse bool,
	canRebase bool,
	refresh func(),
) error {
	infos := c.Helpers().Staging.ChangeLinesInViewRange(mainViewName, firstLineIdx, lastLineIdx)
	if len(infos) == 0 {
		return nil
	}

	patchBuilder := c.Git().Patch.PatchBuilder
	mustDiscardPatch := patchBuilder.Active() && patchBuilder.NewPatchRequired(from, to, reverse)
	return c.ConfirmIf(mustDiscardPatch, types.ConfirmOpts{
		Title:  c.Tr.DiscardPatch,
		Prompt: c.Tr.DiscardPatchConfirm,
		HandleConfirm: func() error {
			if mustDiscardPatch {
				patchBuilder.Reset()
			}
			if !patchBuilder.Active() {
				patchBuilder.Start(from, to, reverse, canRebase)
			}

			if err := togglePatchLines(c, infos); err != nil {
				return err
			}

			refresh()
			revealSelectionAfterPatchToggle(c, mainViewName, firstLineIdx)
			return nil
		},
	})
}

// togglePatchLines toggles the change lines identified by infos into or out of the
// custom patch. The direction is decided once, from the first selected change line (in
// view order) — already in the patch means remove the whole selection, otherwise add it
// — and applied to every file the selection spans, mirroring how the patch explorer
// toggles a range.
func togglePatchLines(c *ControllerCommon, infos []types.DiffLineInfo) error {
	patchBuilder := c.Git().Patch.PatchBuilder

	// Group the selected change lines by file (in view order), then resolve each file's
	// identities to patch-line indices in one pass.
	var order []string
	identitiesByFile := map[string][]patch.LineIdentity{}
	for _, info := range infos {
		filename := patchFilename(c, info.Path)
		if filename == "" {
			continue
		}
		if _, seen := identitiesByFile[filename]; !seen {
			order = append(order, filename)
		}
		lineNumber, isDeletion := info.PatchSelectLine()
		identitiesByFile[filename] = append(identitiesByFile[filename],
			patch.LineIdentity{LineNumber: lineNumber, IsDeletion: isDeletion})
	}
	if len(order) == 0 {
		return nil
	}

	indicesByFile := map[string][]int{}
	for filename, identities := range identitiesByFile {
		indices, err := patchBuilder.PatchLineIndicesForLines(filename, identities)
		if err != nil {
			return err
		}
		indicesByFile[filename] = indices
	}

	// Decide the direction from the first selected change line.
	firstFile := order[0]
	included, err := patchBuilder.GetFileIncLineIndices(firstFile, "")
	if err != nil {
		return err
	}
	remove := len(indicesByFile[firstFile]) > 0 && lo.Contains(included, indicesByFile[firstFile][0])
	toggle := patchBuilder.AddFileLineRange
	if remove {
		toggle = patchBuilder.RemoveFileLineRange
	}

	for _, filename := range order {
		indices := indicesByFile[filename]
		if len(indices) == 0 {
			continue
		}
		if err := toggle(filename, "", indices); err != nil {
			return err
		}
	}
	return nil
}

// revealSelectionAfterPatchToggle re-establishes the selection after a toggle's
// re-render. Toggling doesn't change the diff, so the same content is still there — but
// the re-render (and the layout re-wrap when the secondary view first appears) is in
// view-line space, which moves the selection. We preserve the selection's change-line
// ordinal, which is unchanged (the line isn't consumed, unlike staging) and
// width-independent, re-expanding the hunk in hunk mode. A range collapses to a line,
// as a staged range does.
func revealSelectionAfterPatchToggle(c *ControllerCommon, mainViewName string, firstLineIdx int) {
	mainContext := c.Contexts().Normal
	if mainViewName == c.Contexts().NormalSecondary.GetViewName() {
		mainContext = c.Contexts().NormalSecondary
	}
	view := mainContext.GetView()

	sel := mainContext.DiffSelectState()
	if sel.Mode == context.DiffSelectModeRange {
		sel.Mode = context.DiffSelectModeLine
		sel.RangeIsSticky = false
	}
	mode := sel.Mode

	c.Helpers().Staging.RevealSelectionAfterStaging(view, view, firstLineIdx, func(viewLine int) {
		if mode == context.DiffSelectModeHunk {
			selectDiffHunk(c, mainContext, viewLine)
		} else {
			view.CancelRangeSelect()
			showSelectionAtLine(view, viewLine, true)
		}
	})
}

// patchFilename maps a diff line's absolute path to the key the patch builder stores
// the file under — its repo-relative, slash-separated path.
func patchFilename(c *ControllerCommon, absPath string) string {
	relPath, err := filepath.Rel(c.Git().RepoPaths.WorktreePath(), absPath)
	if err != nil {
		return ""
	}
	return filepath.ToSlash(relPath)
}
