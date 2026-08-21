package helpers

import (
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// diffFilePrefix marks the start of a file's section in a (possibly multi-file)
// unified diff.
const diffFilePrefix = "diff --git "

// parsedDiffLine is what the parser recovers about a row of a rendered diff.
// Path is the path as the diff header spells it, i.e. relative to the repo root;
// the caller turns it into the absolute path of types.DiffLineInfo.
type parsedDiffLine struct {
	Path    string
	Type    types.DiffLineType
	NewLine int
	OldLine int
}

// bufferLineParse is the parser's result for one buffer line: the recovered
// identity, and whether the line could be resolved at all (false for a line in
// an unparseable section, or outside any file section).
type bufferLineParse struct {
	parsed parsedDiffLine
	ok     bool
}

// parseDiffLineFromBuffer recovers the identity of a row of a rendered diff by
// parsing the view's decolorized contents.
//
// bufferLines is the full unwrapped view buffer; targetIdx is the buffer line to
// resolve. A commit's diff spans several files, so we isolate the file section
// containing targetIdx and parse just that one (see parseFileSection). Use this
// for a single line, e.g. the one under the cursor; to resolve every line of a
// buffer, use parseAllDiffLinesFromBuffer, which parses each section only once.
//
// ok is false when the buffer isn't a parseable unified diff at targetIdx,
// because the diff renderer restructured it, so that the caller can fall back.
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

// parseAllDiffLinesFromBuffer resolves every line of a (possibly multi-file)
// diff buffer in one pass, parsing each file section exactly once. It is the
// batch form of parseDiffLineFromBuffer, for callers that scan a whole buffer:
// resolving line by line would re-parse a section once per line of it — O(n²) on
// a large single-file diff — whereas this is O(n). The result is indexed 1:1
// with bufferLines; a line in an unparseable section, or before the first
// "diff --git", is left ok=false.
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

// diffLineTexts extracts the text of each rendered row — the material the buffer
// parser works on.
func diffLineTexts(contents []gocui.DiffLineContent) []string {
	texts := make([]string, len(contents))
	for i, content := range contents {
		texts[i] = content.Text
	}
	return texts
}

// fileSectionBounds returns the half-open range [start, end) of the file section
// containing targetIdx: the nearest "diff --git" at or above it, up to the next
// one (or the end of the buffer). start is -1 when targetIdx is before the first
// file section.
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
// "diff --git" line) a single time and returns the identity of each of its
// lines, indexed 1:1 with fileLines. patch.Parse's line indices line up with the
// section's buffer lines, so the type and the old/new line numbers fall out of
// the patch arithmetic. Every line is left ok=false when the section has no
// recoverable path or isn't a well-formed unified diff — the rendering
// restructured it, and acting on a mis-parse would land us on the wrong line, so
// the caller should fall back.
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
			Path:    relPath,
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
// It prefers the "+++ b/<path>" line, falling back to "--- a/<path>" when the
// new path is /dev/null (a deleted file), and to the "diff --git" line when
// there are no such lines at all (a pure rename, which has no hunks).
func pathFromDiffHeader(fileLines []string) string {
	var oldPath, newPath string
	for _, line := range fileLines {
		if strings.HasPrefix(line, "@@") {
			break // past the header
		}
		switch {
		case strings.HasPrefix(line, "+++ "):
			newPath = pathFromDiffHeaderField(strings.TrimPrefix(line, "+++ "))
		case strings.HasPrefix(line, "--- "):
			oldPath = pathFromDiffHeaderField(strings.TrimPrefix(line, "--- "))
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

// pathFromDiffHeaderField decodes one path field of a diff header — the part
// after "--- " or "+++ ", or one of the two paths on the "diff --git" line —
// into the repo-relative path it names.
//
// git spells such a field in three ways: plain; terminated by a tab, when the
// path contains a space; or C-quoted as a whole, when the path contains
// characters git won't print raw — which, with core.quotePath enabled (the
// default), includes every non-ASCII byte, so `café` arrives as
// `"b/caf\303\251"`. The quoting is Go's string syntax, octal escapes included,
// so strconv decodes it for us.
//
// Returns "" for a quoted field we can't decode: better to resolve nothing than
// to point a consumer at a path that doesn't exist.
func pathFromDiffHeaderField(field string) string {
	field = strings.TrimSuffix(field, "\t")

	if strings.HasPrefix(field, `"`) {
		unquoted, err := strconv.Unquote(field)
		if err != nil {
			return ""
		}
		field = unquoted
	}

	return stripDiffPathPrefix(field)
}

// stripDiffPathPrefix removes the a/ or b/ prefix git puts on the paths in a
// diff header. We ask git for these prefixes explicitly (diff.noprefix=false),
// so they are always there.
func stripDiffPathPrefix(path string) string {
	if strings.HasPrefix(path, "a/") || strings.HasPrefix(path, "b/") {
		return path[2:]
	}
	return path
}

// parseDiffLineMetadata parses the payload of an OSC 1717 record, in which a
// diff renderer states which line of which file it is rendering. The v1 payload
// is positional and ';'-delimited:
//
//	version;type;new-line;old-line;file
//
// The file comes last so that it may itself contain a ';'. The old-file line is
// empty unless the line is a deletion, the only kind that needs it, and the
// new-file line is empty on a file header, the one kind that has no line.
//
// ok is false for a payload of an unknown version or shape, so that the caller
// can fall back to reading the rendered text.
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

	return parsedDiffLine{Path: fields[4], Type: lineType, NewLine: newLine, OldLine: oldLine}, true
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

// pathFromDiffGitLine extracts the new-file path from a "diff --git a/X b/X"
// line, where the two paths are separated by a space and either may be quoted.
// A path containing " b/" (or ` "b/`) would defeat this, but the +++/--- lines
// are unambiguous and we only get here when they are absent.
func pathFromDiffGitLine(line string) string {
	rest := strings.TrimPrefix(line, diffFilePrefix)
	if idx := strings.LastIndex(rest, ` "b/`); idx != -1 {
		return pathFromDiffHeaderField(rest[idx+1:])
	}
	if idx := strings.LastIndex(rest, " b/"); idx != -1 {
		return pathFromDiffHeaderField(rest[idx+1:])
	}
	return ""
}
