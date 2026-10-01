package helpers

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// diffFilePrefix marks the start of a file's section in a (possibly multi-file)
// unified diff.
const diffFilePrefix = "diff --git "

// submodulePrefix opens the line a submodule's section of a diff starts with, which
// stands in for the "diff --git" header a file of the repo gets.
const submodulePrefix = "Submodule "

// submoduleSectionPattern matches that line and captures the submodule's path. git
// writes one of two kinds: the commit the submodule is checked out at has moved
// ("Submodule sub a32f27c..2d9f921:", with "..." in place of ".." where the move is
// no fast-forward, " (rewind)" where it goes backwards, and a message in brackets in
// place of the colon where the two commits can't both be read), or its working tree
// is dirty ("Submodule sub contains untracked content").
//
// The path is captured greedily: git writes it unquoted, so a path that itself ends
// in something reading like a range of commits is told apart by taking the last such
// range on the line.
var submoduleSectionPattern = regexp.MustCompile(
	`^Submodule (.+) (?:contains (?:untracked|modified) content|[0-9a-f]+\.{2,3}[0-9a-f]+(?: \(.*\))?:?)$`)

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
	r := parseFileSection(bufferLines[start:end], end == len(bufferLines))[targetIdx-start]
	return r.parsed, r.ok
}

// parseAllDiffLinesFromBuffer resolves every line of a (possibly multi-file)
// diff buffer in one pass, parsing each file section exactly once. It is the
// batch form of parseDiffLineFromBuffer, for callers that scan a whole buffer:
// resolving line by line would re-parse a section once per line of it — O(n²) on
// a large single-file diff — whereas this is O(n). The result is indexed 1:1
// with bufferLines; a line in an unparseable section, or above the first one, is
// left ok=false.
func parseAllDiffLinesFromBuffer(bufferLines []string) []bufferLineParse {
	result := make([]bufferLineParse, len(bufferLines))
	for i := 0; i < len(bufferLines); {
		if !startsFileSection(bufferLines[i]) {
			i++ // in no file section; leave it unresolved
			continue
		}
		end := fileSectionEnd(bufferLines, i)
		copy(result[i:end], parseFileSection(bufferLines[i:end], end == len(bufferLines)))
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

// renderingStatesDiffLines reports whether the renderer stated, for at least one row
// of the rendering, which diff line it shows. The version-only record a renderer
// announces the protocol with names no line, and doesn't count.
//
// The answer settles how the whole rendering is read. A renderer that states its
// lines lays the diff out as it likes, so its text is no unified diff and must not be
// parsed as one, not even for the rows it says nothing about. Such a row can read
// like a diff header when the file being diffed is itself a diff; the parser would
// take that for the start of a file section and place every untagged row below it
// in a file the diff doesn't have. So the rows a renderer leaves untagged (dividers,
// padding) have no identity, as the protocol has it. A rendering without any record
// is a diff that describes itself, and is parsed as one.
func renderingStatesDiffLines(contents []gocui.DiffLineContent) bool {
	return slices.ContainsFunc(contents, func(content gocui.DiffLineContent) bool {
		return slices.ContainsFunc(content.Metadata, func(record string) bool {
			_, ok := parseDiffLineMetadata(record)
			return ok
		})
	})
}

// parseDiffLineRecords parses the records a row carries, left to right, leaving out
// the ones we don't understand. A row carries more than one when the rendering puts
// two diff lines on it, as a side-by-side row does with a deletion and the addition
// replacing it.
func parseDiffLineRecords(metadata []string) []parsedDiffLine {
	parsed := make([]parsedDiffLine, 0, len(metadata))
	for _, record := range metadata {
		if line, ok := parseDiffLineMetadata(record); ok {
			parsed = append(parsed, line)
		}
	}
	return parsed
}

// parseDiffLineIdentities recovers, for every row of a rendering, the diff lines it
// shows, indexed 1:1 with contents; a row that shows none we can place gets an empty
// entry. The rendering is read the way renderingStatesDiffLines settles: by the
// renderer's records, every one a row carries, or else by parsing the rendering as a
// unified diff, where each row shows one line. Each file's section is parsed once;
// resolving row by row would re-run that parse once per row, O(n²) on a large
// single-file diff.
func parseDiffLineIdentities(contents []gocui.DiffLineContent) [][]parsedDiffLine {
	identities := make([][]parsedDiffLine, len(contents))
	if renderingStatesDiffLines(contents) {
		for i, content := range contents {
			if parsed := parseDiffLineRecords(content.Metadata); len(parsed) > 0 {
				identities[i] = parsed
			}
		}
		return identities
	}

	for i, parsed := range parseAllDiffLinesFromBuffer(diffLineTexts(contents)) {
		if parsed.ok {
			identities[i] = []parsedDiffLine{parsed.parsed}
		}
	}
	return identities
}

// fileSectionBounds returns the half-open range [start, end) of the file section
// containing targetIdx: the nearest line starting a section at or above it, up to
// where that section ends. start is -1 when targetIdx is above the first file
// section, or below the end of the last one that begins above it.
func fileSectionBounds(bufferLines []string, targetIdx int) (start, end int) {
	for start = targetIdx; start >= 0; start-- {
		if !startsFileSection(bufferLines[start]) {
			continue
		}
		if end = fileSectionEnd(bufferLines, start); targetIdx < end {
			return start, end
		}
		return -1, -1
	}
	return -1, -1
}

// fileSectionEnd returns the line the file section beginning at start ends before.
//
// A file's section runs to the next one, since every line between them is part of
// its diff. A submodule's runs only as far as what git writes for it — the line
// naming it, and the log of the commits it moved over — because the lines after
// that need not belong to any section at all. A diff renderer's output has no
// "diff --git" line to stop at, and the rows it puts between one file and the next
// belong to neither.
func fileSectionEnd(bufferLines []string, start int) int {
	if submodulePath(bufferLines[start]) != "" {
		end := start + 1
		for end < len(bufferLines) && isSubmoduleLogLine(bufferLines[end]) {
			end++
		}
		return end
	}

	for i := start + 1; i < len(bufferLines); i++ {
		if startsFileSection(bufferLines[i]) {
			return i
		}
	}
	return len(bufferLines)
}

// startsFileSection reports whether the line opens a section of a diff: git's header
// for a file of the repo, or the line a submodule's section begins with.
func startsFileSection(line string) bool {
	return strings.HasPrefix(line, diffFilePrefix) || submodulePath(line) != ""
}

// isSubmoduleLogLine reports whether the line is one of the commits git lists under
// a submodule's header, which it writes as two spaces, the direction the commit was
// moved in, and the commit's subject.
func isSubmoduleLogLine(line string) bool {
	return strings.HasPrefix(line, "  > ") || strings.HasPrefix(line, "  < ")
}

// submodulePath returns the submodule whose section the given line opens, and "" for
// every other line.
//
// A submodule gets no "diff --git" header and no hunks: git states which commits it
// moved between and lists them, so that one line is all there is to take the path
// from. The prefix is tested first so that the pattern is run over next to no lines
// of a diff.
func submodulePath(line string) string {
	if !strings.HasPrefix(line, submodulePrefix) {
		return ""
	}
	if match := submoduleSectionPattern.FindStringSubmatch(line); match != nil {
		return match[1]
	}
	return ""
}

// parseFileSection parses one file's diff section (fileLines, starting at the line
// that opens it) a single time and returns the identity of each of its
// lines, indexed 1:1 with fileLines. patch.Parse's line indices line up with the
// section's buffer lines, so the type and the old/new line numbers fall out of
// the patch arithmetic. A submodule's section has no hunks at all, so every row of
// it comes out as a header of the submodule, which is what they are: what git states
// there is which commits it moved between, not lines of a file.
//
// Every line is left ok=false when the section has no
// recoverable path or isn't a well-formed unified diff — the rendering
// restructured it, and acting on a mis-parse would land us on the wrong line, so
// the caller should fall back.
//
// endsTheBuffer says the section runs to the end of what we were given. That is
// where a diff we have only part of breaks off. A long one is read a screenful
// at a time and the rest as the user scrolls, so its last hunk holds fewer lines
// than its header declares until the reading is done. Insisting on the whole
// hunk there would leave every line of the file unresolved while the diff is the
// one on screen, so a section in that position is held to what has arrived.
func parseFileSection(fileLines []string, endsTheBuffer bool) []bufferLineParse {
	result := make([]bufferLineParse, len(fileLines))

	relPath := pathFromDiffHeader(fileLines)
	if relPath == "" {
		return result
	}
	p := patch.Parse(strings.Join(fileLines, "\n"))
	isWellFormed := p.IsWellFormed
	if endsTheBuffer {
		isWellFormed = p.IsWellFormedSoFar
	}
	if !isWellFormed() {
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

// pathFromDiffHeader extracts the new-file path of a single diff section. A
// submodule's section states its path in the line it opens with. For a file of the
// repo the path comes from the "+++ b/<path>" line, falling back to "--- a/<path>"
// when the new path is /dev/null (a deleted file), and to the "diff --git" line when
// there are no such lines at all (a pure rename, which has no hunks).
func pathFromDiffHeader(fileLines []string) string {
	if path := submodulePath(fileLines[0]); path != "" {
		return path
	}

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
