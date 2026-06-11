package helpers

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// diffFilePrefix marks the start of a file's section in a (possibly multi-file)
// unified diff.
const diffFilePrefix = "diff --git "

// parsedDiffLine is what a backend recovers about a rendered diff row. RelPath
// is the path as the source spells it — repo-relative from the diff header for
// the buffer parser (#1), or whatever the pager emitted (possibly absolute) for
// the OSC metadata (#2); the caller turns it into the absolute path of
// types.DiffLineInfo.
type parsedDiffLine struct {
	RelPath string
	Type    types.DiffLineType
	NewLine int
	OldLine int
}

// parseDiffLineFromBuffer recovers a rendered diff row's patch-space identity by
// parsing the decolorized diff buffer (mechanism #1 in diff-line-metadata-notes.md).
//
// bufferLines is the full unwrapped view buffer; targetIdx is the buffer line to
// resolve. A commit diff spans multiple files, so we first split on the
// "diff --git" boundaries to isolate the file section containing targetIdx, then
// reuse patch.Parse on that single-file section: its patch line indices line up
// 1:1 with the section's buffer lines, so the type and old/new line numbers fall
// straight out of the patch arithmetic.
//
// ok is false when the buffer isn't a parseable unified diff at targetIdx (e.g.
// a pager that restructures the diff, like delta's default mode), so the caller
// can fall back to another backend.
func parseDiffLineFromBuffer(bufferLines []string, targetIdx int) (parsedDiffLine, bool) {
	if targetIdx < 0 || targetIdx >= len(bufferLines) {
		return parsedDiffLine{}, false
	}

	// Find the file section containing the target: the nearest "diff --git" at or
	// above it, up to the next one (or the end of the buffer).
	fileStart := -1
	for i := targetIdx; i >= 0; i-- {
		if strings.HasPrefix(bufferLines[i], diffFilePrefix) {
			fileStart = i
			break
		}
	}
	if fileStart == -1 {
		return parsedDiffLine{}, false
	}
	fileEnd := len(bufferLines)
	for i := fileStart + 1; i < len(bufferLines); i++ {
		if strings.HasPrefix(bufferLines[i], diffFilePrefix) {
			fileEnd = i
			break
		}
	}

	fileLines := bufferLines[fileStart:fileEnd]
	relPath := pathFromDiffHeader(fileLines)
	if relPath == "" {
		return parsedDiffLine{}, false
	}

	p := patch.Parse(strings.Join(fileLines, "\n"))
	// Bail if the body doesn't match the hunk headers: the rendering restructured
	// the diff (e.g. delta's line-number gutters push the +/- marker off the
	// start of the line, so every body line reads as context), and trusting the
	// mis-parse would land us on the wrong line. Better to fall back.
	if !p.IsWellFormed() {
		return parsedDiffLine{}, false
	}
	patchLines := p.Lines()
	patchLineIdx := targetIdx - fileStart
	if patchLineIdx < 0 || patchLineIdx >= len(patchLines) {
		return parsedDiffLine{}, false
	}

	result := parsedDiffLine{
		RelPath: relPath,
		Type:    diffLineTypeForKind(patchLines[patchLineIdx].Kind),
		NewLine: p.LineNumberOfLine(patchLineIdx),
	}
	if result.Type == types.DiffLineDeleted {
		result.OldLine = p.OldLineNumberOfLine(patchLineIdx)
	}
	return result, true
}

func diffLineTypeForKind(kind patch.PatchLineKind) types.DiffLineType {
	switch kind {
	case patch.PATCH_HEADER:
		return types.DiffLineFileHeader
	case patch.HUNK_HEADER:
		return types.DiffLineHunkHeader
	case patch.ADDITION:
		return types.DiffLineAdded
	case patch.DELETION:
		return types.DiffLineDeleted
	case patch.CONTEXT:
		return types.DiffLineContext
	default:
		return types.DiffLineOther
	}
}

// pathFromDiffHeader extracts the new-file path of a single file's diff section.
// It prefers the "+++ b/<path>" line (falling back to "--- a/<path>" when the
// new path is /dev/null, i.e. a deleted file), and as a last resort the
// "diff --git" line. Paths with characters git C-quotes are not handled (this is
// a prototype); the common unquoted case is.
func pathFromDiffHeader(fileLines []string) string {
	var oldPath, newPath string
	for _, line := range fileLines {
		if strings.HasPrefix(line, "@@") {
			break // past the header
		}
		switch {
		case strings.HasPrefix(line, "+++ "):
			newPath = stripDiffPathPrefix(strings.TrimPrefix(line, "+++ "))
		case strings.HasPrefix(line, "--- "):
			oldPath = stripDiffPathPrefix(strings.TrimPrefix(line, "--- "))
		}
	}

	if newPath != "" && newPath != "/dev/null" {
		return newPath
	}
	if oldPath != "" && oldPath != "/dev/null" {
		return oldPath
	}
	return pathFromDiffGitLine(fileLines[0])
}

// stripDiffPathPrefix removes git's default a/ or b/ diff path prefix if present.
func stripDiffPathPrefix(path string) string {
	if strings.HasPrefix(path, "a/") || strings.HasPrefix(path, "b/") {
		return path[2:]
	}
	return path
}

// pathFromDiffGitLine extracts the new-file path from a "diff --git a/X b/X"
// line, used only when the +++/--- lines are absent.
func pathFromDiffGitLine(line string) string {
	rest := strings.TrimPrefix(line, diffFilePrefix)
	if idx := strings.LastIndex(rest, " b/"); idx != -1 {
		return rest[idx+len(" b/"):]
	}
	return ""
}

// parseDiffLineMetadata parses mechanism #2's OSC 1717 payload (v1):
// version;type;new-line;old-line;file — positional and ';'-delimited, with the
// file last (so it may itself contain ';') and old-line empty unless the line is
// a deletion. See diff-line-metadata-notes.md §9.2. ok is false for a payload of
// an unknown version or shape, so the caller can fall back to another backend.
func parseDiffLineMetadata(payload string) (parsedDiffLine, bool) {
	fields := strings.SplitN(payload, ";", 5)
	if len(fields) < 5 || fields[0] != "1" {
		return parsedDiffLine{}, false
	}

	lineType, ok := diffLineTypeFromMetadata(fields[1])
	if !ok {
		return parsedDiffLine{}, false
	}

	newLine, err := strconv.Atoi(fields[2])
	if err != nil {
		return parsedDiffLine{}, false
	}

	oldLine := 0
	if fields[3] != "" {
		if oldLine, err = strconv.Atoi(fields[3]); err != nil {
			return parsedDiffLine{}, false
		}
	}

	return parsedDiffLine{RelPath: fields[4], Type: lineType, NewLine: newLine, OldLine: oldLine}, true
}

func diffLineTypeFromMetadata(typeField string) (types.DiffLineType, bool) {
	switch typeField {
	case "c":
		return types.DiffLineContext, true
	case "a":
		return types.DiffLineAdded, true
	case "d":
		return types.DiffLineDeleted, true
	default:
		return types.DiffLineOther, false
	}
}
