package helpers

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// A diff opens with a diffstat naming every file in it, above the diff of each of
// them. Here each of those names is made a link to where that file's diff begins, so
// that a file of a long diff can be gone to by clicking the line that names it.
//
// The names are recognized in the output as it is written to the pane, where they cost
// next to nothing to find. The diffstat is git's own text whichever renderer the diff
// goes through — delta, diff-so-fancy and difftastic all pass it on untouched — and it
// comes first, so the scan for it ends with it.

// DiffStatLinkScheme names a link to a file of the diff the pane is showing, as
// lazygit-edit names one that opens a file in the editor. The link is never handed to
// the terminal — gocui takes the escape sequence out of the content and gives the URL
// back when the cell it covers is clicked — so the path in it needs no escaping.
const DiffStatLinkScheme = "lazygit-diff-file://"

// diffStatEntryPattern matches a line of a diffstat and captures the path it states.
// Such a line holds the name of the file, padded out to the width of the longest, then
// the number of lines it changes (or "Bin" for a binary file) and the graph of them.
//
// The name is captured greedily, so that the separator found is the last one on the
// line rather than one in a file name that contains " | " itself.
var diffStatEntryPattern = regexp.MustCompile(`^ (.*[^ ]) +\| +(?:Bin|\d+)`)

// DiffStatLinkWriter hands a pane's content on to it, turning the file names in the
// diffstat the content opens with into links (see DiffStatLinkScheme).
type DiffStatLinkWriter struct {
	writer io.Writer

	// Whether the diffstat is still to come, is being written now, or is behind us. It
	// is behind us once a line comes that is no entry of it, or that the diff proper
	// begins with, and nothing past that is looked at. A name down there names the file
	// the reader is already in.
	//
	// It is atomic because the render is begun on the UI thread while the content of
	// it arrives on the goroutine reading the command's output.
	state atomic.Int32
}

type diffStatState int32

const (
	diffStatToCome diffStatState = iota
	inDiffStat
	diffStatDone
)

func NewDiffStatLinkWriter(writer io.Writer) *DiffStatLinkWriter {
	return &DiffStatLinkWriter{writer: writer}
}

// BeginRender starts a fresh render, whose own diffstat is the one to look for. It is
// called as the render is asked for, before any of it is written.
//
// linkFiles says whether this render is one whose file names lead anywhere: the pane's
// own diff, whose rows can be placed in the files they show. Either is passed through
// untouched without it. Content that is no diff of the panel's — a commit log, a
// message — has no diffstat in it, and a line of one that happens to read like an entry
// of a diffstat names no file to go to. A rendering whose rows nothing can place does
// have the files in it, but nothing to find the one a name stands for with.
func (self *DiffStatLinkWriter) BeginRender(linkFiles bool) {
	self.setState(lo.Ternary(linkFiles, diffStatToCome, diffStatDone))
}

func (self *DiffStatLinkWriter) getState() diffStatState {
	return diffStatState(self.state.Load())
}

func (self *DiffStatLinkWriter) setState(state diffStatState) {
	self.state.Store(int32(state))
}

func (self *DiffStatLinkWriter) Write(p []byte) (int, error) {
	linked := self.withFileNameLinked(p)

	written, err := self.writer.Write(linked)
	if err != nil {
		return 0, err
	}
	if written < len(linked) {
		return 0, io.ErrShortWrite
	}
	// The caller is owed an answer about what it gave us, not about what we passed on.
	return len(p), nil
}

// withFileNameLinked returns the given line of the render with the name in it linked,
// where the line is an entry of the diffstat.
func (self *DiffStatLinkWriter) withFileNameLinked(line []byte) []byte {
	state := self.getState()
	if state == diffStatDone {
		return line
	}

	if beginsTheDiffItself(line) {
		// A diffstat that hasn't come by now isn't coming: the pane is showing a diff
		// that was asked for without one.
		self.setState(diffStatDone)
		return line
	}

	match := diffStatEntry(line)
	if match == nil {
		if state == inDiffStat {
			self.setState(diffStatDone)
		}
		return line
	}
	self.setState(inDiffStat)

	start, end := match[2], match[3]
	// The link states the name as it reads on screen, so that a renderer that colors
	// the diffstat doesn't put escape sequences into it.
	name := utils.Decolorise(string(line[start:end]))
	linked := make([]byte, 0, len(line)+len(name)+32)
	linked = append(linked, line[:start]...)
	linked = append(linked, style.PrintHyperlink(string(line[start:end]), DiffStatLinkScheme+name)...)
	return append(linked, line[end:]...)
}

// diffLineRecordOpener opens an OSC 1717 record, ahead of the version whose fields the
// record states (see parseDiffLineMetadata).
const diffLineRecordOpener = "\x1b]1717;"

// beginsTheDiffItself reports whether the line is one of the diff proper rather than
// one of the diffstat above it: git's own header for a file, or a line a renderer
// states a record for.
func beginsTheDiffItself(line []byte) bool {
	return bytes.HasPrefix(line, []byte("diff --")) || statesADiffLine(line)
}

// statesADiffLine reports whether the line carries a record in which a diff renderer
// states which line of which file it is rendering. Every such record is about a line of
// the diff, and the diffstat comes before all of them.
//
// The version the record opens with is passed over rather than read. This asks which
// lines a renderer states records for, and the answer holds whichever version of the
// protocol it speaks. A record with nothing after the version is the handshake a
// renderer announces itself with, which is about no line, so the search goes on past it.
func statesADiffLine(line []byte) bool {
	for rest := line; ; {
		at := bytes.Index(rest, []byte(diffLineRecordOpener))
		if at == -1 {
			return false
		}
		rest = rest[at+len(diffLineRecordOpener):]

		digits := 0
		for digits < len(rest) && rest[digits] >= '0' && rest[digits] <= '9' {
			digits++
		}
		if digits > 0 && digits < len(rest) && rest[digits] == ';' {
			return true
		}
	}
}

// diffStatEntry matches line against diffStatEntryPattern, behind the two checks that
// answer for nearly every line of a diff without the pattern being run at all: an entry
// of a diffstat is indented by a space, and holds the separator. The indices it returns
// are into the whole line, whatever the match was made past.
func diffStatEntry(line []byte) []int {
	start := handshakeEnd(line)
	rest := line[start:]
	if len(rest) == 0 || rest[0] != ' ' || bytes.IndexByte(rest, '|') == -1 {
		return nil
	}

	match := diffStatEntryPattern.FindSubmatchIndex(rest)
	for i := range match {
		if match[i] >= 0 {
			match[i] += start
		}
	}
	return match
}

// handshakeEnd returns where the record a renderer announces itself with ends, for a
// line that opens with one, and 0 for every other line.
//
// A renderer writes the handshake before anything else and with no newline after it,
// so it lands at the start of the first line of its output. For a diff with nothing
// above its diffstat — the diff of a range of commits, or of a stash — that is the line
// naming the first file in it, and the entry begins after the record rather than at the
// start of the line.
func handshakeEnd(line []byte) int {
	if !bytes.HasPrefix(line, []byte(diffLineRecordOpener)) {
		return 0
	}

	after := len(diffLineRecordOpener)
	for after < len(line) && line[after] >= '0' && line[after] <= '9' {
		after++
	}

	// Either terminator ends a record. One that goes on into a field instead is about
	// a line of the diff, which is below the whole diffstat and no entry of it.
	switch {
	case after < len(line) && line[after] == '\x07':
		return after + 1
	case after+1 < len(line) && line[after] == '\x1b' && line[after+1] == '\\':
		return after + 2
	}
	return 0
}

// JumpToFileNamedInDiffStat goes to the file of the pane's diff that the given diffstat
// entry names, for a click on the link made for that entry. It lands the way picking
// the file from the menu of the diff's files does.
//
// The diff is read to the end first, as it is for that menu. The diffstat is on screen
// only while the view is at the top of the diff, so the file clicked is nearly always
// below the part of it that has been read.
func (self *DiffLineHelper) JumpToFileNamedInDiffStat(pane types.DiffPaneContext, entry string) {
	manager := self.c.GetViewBufferManagerForView(pane.GetView())
	if manager == nil {
		return
	}
	manager.ReadToEnd(func() {
		self.c.OnUIThread(func() error {
			self.jumpToFileNamedInDiffStat(pane, entry)
			return nil
		})
	})
}

func (self *DiffLineHelper) jumpToFileNamedInDiffStat(pane types.DiffPaneContext, entry string) {
	view := pane.GetView()
	worktreePath := self.c.Git().RepoPaths.WorktreePath()
	files := self.FilesInDiff(view)
	names := lo.Map(files, func(file string, _ int) string {
		return repoRelativePath(worktreePath, file)
	})

	index, ok := fileNamedByDiffStatEntry(entry, names)
	if !ok {
		self.c.ErrorToast(utils.ResolvePlaceholderString(
			self.c.Tr.NoFileInDiffNamed, map[string]string{"path": entry}))
		return
	}

	if target, ok := self.StartOfFileInDiff(view, files[index]); ok {
		self.PlaceNavigationTarget(pane, target, true)
	}
}

// fileNamedByDiffStatEntry returns which of the diff's files a diffstat entry names.
//
// An entry states the path as the diffstat has room for it rather than as git names
// the file. A path too long for the column is cut off on the left behind "...", and a
// rename is compacted to the "{old => new}" form. So the name is looked for among the
// files the diff turned out to hold, whole and then as the end of one, and is taken
// only where it names a single file.
func fileNamedByDiffStatEntry(entry string, paths []string) (int, bool) {
	name := renamedTo(strings.TrimSpace(entry))

	if index, ok := theOneMatching(paths, func(p string) bool { return p == name }); ok {
		return index, true
	}

	// Where the diffstat cut the path off, what is left is the end of it. The cut is at
	// a directory boundary where there is room for one, and inside the file name where
	// there isn't.
	tail := strings.TrimPrefix(name, "...")
	return theOneMatching(paths, func(p string) bool { return strings.HasSuffix(p, tail) })
}

// renamedTo returns the path a diffstat entry for a rename leaves the file at, and the
// entry itself for any other one. A rename states both paths, with whatever they have
// in common written once: "dir/{old => new}/file", or "old => new" where they share
// nothing. The part shared with the old path is gone along with the "{" when the entry
// is cut off on the left, which leaves a path to match the end of.
func renamedTo(entry string) string {
	const arrow = " => "
	at := strings.Index(entry, arrow)
	if at == -1 {
		return entry
	}

	shared := ""
	if brace := strings.Index(entry[:at], "{"); brace != -1 {
		shared = entry[:brace]
	}
	renamed := entry[at+len(arrow):]
	if closing := strings.Index(renamed, "}"); closing != -1 {
		return shared + renamed[:closing] + renamed[closing+1:]
	}
	return shared + renamed
}

// theOneMatching returns the index of the one element the predicate holds for, and
// false where it holds for none of them or for several.
func theOneMatching(paths []string, matches func(string) bool) (int, bool) {
	found := -1
	for i, candidate := range paths {
		if !matches(candidate) {
			continue
		}
		if found != -1 {
			return 0, false
		}
		found = i
	}
	return found, found != -1
}
