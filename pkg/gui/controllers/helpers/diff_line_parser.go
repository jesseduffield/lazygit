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

// bufferLineParse is the buffer-parse backend's result for one buffer line: the
// recovered identity and whether the line could be resolved (false for a line in
// an unparseable or restructured section, or outside any file section).
type bufferLineParse struct {
	parsed parsedDiffLine
	ok     bool
}

// parseDiffLineFromBuffer recovers a rendered diff row's patch-space identity by
// parsing the decolorized diff buffer (mechanism #1 in diff-line-metadata-notes.md).
//
// bufferLines is the full unwrapped view buffer; targetIdx is the buffer line to
// resolve. A commit diff spans multiple files, so we isolate the file section
// containing targetIdx and parse just that one (see parseFileSection). Use this for
// a single line (e.g. a click); to resolve every line of a buffer, prefer
// parseAllDiffLinesFromBuffer, which parses each section only once.
//
// ok is false when the buffer isn't a parseable unified diff at targetIdx (e.g.
// a pager that restructures the diff, like delta's default mode), so the caller
// can fall back to another backend.
func parseDiffLineFromBuffer(bufferLines []string, targetIdx int) (parsedDiffLine, bool) {
	if targetIdx < 0 || targetIdx >= len(bufferLines) {
		return parsedDiffLine{}, false
	}
	start, end := fileSectionBounds(bufferLines, targetIdx)
	if start == -1 {
		return parsedDiffLine{}, false
	}
	r := parseFileSection(bufferLines[start:end])[targetIdx-start]
	return r.parsed, r.ok
}

// parseAllDiffLinesFromBuffer resolves every line of a (possibly multi-file) diff
// buffer in one pass, parsing each file section exactly once. It is the batch form
// of parseDiffLineFromBuffer for callers that scan the whole buffer (the position-
// restore and navigation scans): resolving line-by-line would re-parse each section
// once per line — O(n²) on a large single-file diff — whereas this is O(n). The
// returned slice is indexed 1:1 with bufferLines; a line in an unparseable section,
// or before the first "diff --git", is left ok=false.
func parseAllDiffLinesFromBuffer(bufferLines []string) []bufferLineParse {
	result := make([]bufferLineParse, len(bufferLines))
	for i := 0; i < len(bufferLines); {
		if !strings.HasPrefix(bufferLines[i], diffFilePrefix) {
			i++ // not in a file section yet; leave it unresolved
			continue
		}
		_, end := fileSectionBounds(bufferLines, i)
		copy(result[i:end], parseFileSection(bufferLines[i:end]))
		i = end
	}
	return result
}

// fileSectionBounds returns the half-open range [start, end) of the file section
// containing targetIdx: the nearest "diff --git" at or above it, up to the next one
// (or the end of the buffer). start is -1 when targetIdx is before the first file
// section.
func fileSectionBounds(bufferLines []string, targetIdx int) (start, end int) {
	start = -1
	for i := targetIdx; i >= 0; i-- {
		if strings.HasPrefix(bufferLines[i], diffFilePrefix) {
			start = i
			break
		}
	}
	if start == -1 {
		return -1, -1
	}
	end = len(bufferLines)
	for i := start + 1; i < len(bufferLines); i++ {
		if strings.HasPrefix(bufferLines[i], diffFilePrefix) {
			end = i
			break
		}
	}
	return start, end
}

// parseFileSection parses one file's diff section (fileLines, starting at its
// "diff --git" line) a single time and returns the resolved identity for each of
// its lines, indexed 1:1 with fileLines. patch.Parse's line indices line up with
// the section's buffer lines, so the type and old/new line numbers fall straight
// out of the patch arithmetic. Every line is left ok=false when the section has no
// recoverable path or isn't a well-formed unified diff — the rendering restructured
// it (e.g. delta's line-number gutters push the +/- marker off the start of the
// line, so every body line reads as context), and trusting the mis-parse would land
// us on the wrong line, so the caller should fall back to another backend.
func parseFileSection(fileLines []string) []bufferLineParse {
	result := make([]bufferLineParse, len(fileLines))

	relPath := pathFromDiffHeader(fileLines)
	if relPath == "" {
		return result
	}
	p := patch.Parse(strings.Join(fileLines, "\n"))
	if !p.IsWellFormed() {
		return result
	}
	patchLines := p.Lines()
	for i := range fileLines {
		if i >= len(patchLines) {
			break
		}
		parsed := parsedDiffLine{
			RelPath: relPath,
			Type:    diffLineTypeForKind(patchLines[i].Kind),
			NewLine: p.LineNumberOfLine(i),
		}
		if parsed.Type == types.DiffLineDeleted {
			parsed.OldLine = p.OldLineNumberOfLine(i)
		}
		result[i] = bufferLineParse{parsed, true}
	}
	return result
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
// file last (so it may itself contain ';'), old-line empty unless the line is a
// deletion, and new-line empty on a file header (the one type that carries no
// line number). See diff-line-metadata-notes.md §9.2. ok is false for a payload
// of an unknown version or shape, so the caller can fall back to another backend.
func parseDiffLineMetadata(payload string) (parsedDiffLine, bool) {
	fields := strings.SplitN(payload, ";", 5)
	if len(fields) < 5 || fields[0] != "1" {
		return parsedDiffLine{}, false
	}

	lineType, ok := diffLineTypeFromMetadata(fields[1])
	if !ok {
		return parsedDiffLine{}, false
	}

	newLine := 0
	if fields[2] != "" {
		var err error
		if newLine, err = strconv.Atoi(fields[2]); err != nil {
			return parsedDiffLine{}, false
		}
	} else if lineType != types.DiffLineFileHeader {
		return parsedDiffLine{}, false
	}

	oldLine := 0
	if fields[3] != "" {
		var err error
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
	case "f":
		return types.DiffLineFileHeader, true
	case "h":
		return types.DiffLineHunkHeader, true
	default:
		return types.DiffLineOther, false
	}
}
