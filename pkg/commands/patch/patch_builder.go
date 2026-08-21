package patch

import (
	"os"
	"sort"
	"strings"

	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/generics/set"
	"github.com/samber/lo"
	"github.com/sasha-s/go-deadlock"
	"github.com/sirupsen/logrus"
)

type PatchStatus int

const (
	// UNSELECTED is for when the commit file has not been added to the patch in any way
	UNSELECTED PatchStatus = iota
	// WHOLE is for when you want to add the whole diff of a file to the patch,
	// including e.g. if it was deleted
	WHOLE
	// PART is for when you're only talking about specific lines that have been modified
	PART
)

type fileInfo struct {
	mode                PatchStatus
	includedLineIndices []int
	diff                string
	// For a renamed file, the path it was renamed from; empty otherwise. We
	// need to keep hold of it so we can re-render the file's patch (which is
	// keyed by the new path) without the caller having to supply it again.
	previousPath string
}

type (
	loadFileDiffFunc func(from string, to string, reverse bool, filename string, previousPath string) (string, error)
)

// PatchBuilder manages the building of a patch for a commit to be applied to another commit (or the working tree, or removed from the current commit). We also support building patches from things like stashes, for which there is less flexibility
type PatchBuilder struct {
	// To is the commit hash if we're dealing with files of a commit, or a stash ref for a stash
	To      string
	From    string
	reverse bool

	// CanRebase tells us whether we're allowed to modify our commits. CanRebase should be true for commits of the currently checked out branch and false for everything else
	// TODO: move this out into a proper mode struct in the gui package: it doesn't really belong here
	CanRebase bool

	// fileInfoMap starts empty but you add files to it as you go along
	fileInfoMap map[string]*fileInfo
	Log         *logrus.Entry

	// mutex guards the fields that a git worker can mutate (via Reset, at the
	// end of a patch-consuming operation) while the UI thread reads them to
	// render — chiefly To and the fileInfoMap pointer. The map's *entries* are
	// only ever touched on the UI thread, so we only hold the lock long enough
	// to read or swap the fields, never across the git I/O in getFileInfo.
	mutex deadlock.Mutex

	// loadFileDiff loads the diff of a file, for a given to (typically a commit hash)
	loadFileDiff loadFileDiffFunc

	// newTempDir makes a directory for the current patch to be materialized into, as
	// two file trees that can be diffed against each other and so rendered like any
	// other diff (see PatchCommands.WriteCustomPatchDiffTrees). Its lifetime is the
	// patch's: made when one is started, removed when it is given up.
	newTempDir func() (string, error)
	tempDir    string

	// generation counts the changes made to the patch, so that whoever materializes it
	// can tell whether what they last built still describes it — and rebuild only then,
	// rather than on every render of it.
	generation int
}

func NewPatchBuilder(
	log *logrus.Entry, loadFileDiff loadFileDiffFunc, newTempDir func() (string, error),
) *PatchBuilder {
	return &PatchBuilder{
		Log:          log,
		loadFileDiff: loadFileDiff,
		newTempDir:   newTempDir,
	}
}

func (p *PatchBuilder) Start(from, to string, reverse bool, canRebase bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.generation++
	p.makeTempDir()

	p.To = to
	p.From = from
	p.reverse = reverse
	p.CanRebase = canRebase
	p.fileInfoMap = map[string]*fileInfo{}
}

// snapshotFileInfoMap returns the current fileInfoMap under the lock. The map's
// entries are only mutated on the UI thread, so callers can read the returned
// map without holding the lock; the lock only serializes the pointer swap that
// Reset/Start do (potentially from a git worker) against these reads.
func (p *PatchBuilder) snapshotFileInfoMap() map[string]*fileInfo {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.fileInfoMap
}

// TempDir is the directory the patch is materialized into for rendering, and "" when
// there is none — no patch, or a directory we failed to make.
func (p *PatchBuilder) TempDir() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.tempDir
}

// Generation says which version of the patch this is; see the field.
func (p *PatchBuilder) Generation() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.generation
}

// makeTempDir replaces the directory the patch is materialized into with a fresh one.
// Only call this with the lock held.
func (p *PatchBuilder) makeTempDir() {
	p.removeTempDir()
	if p.newTempDir == nil {
		return
	}
	dir, err := p.newTempDir()
	if err != nil {
		p.Log.Error(err)
		return
	}
	p.tempDir = dir
}

// removeTempDir takes the patch's materialized form away with the patch. Only call this
// with the lock held.
func (p *PatchBuilder) removeTempDir() {
	if p.tempDir == "" {
		return
	}
	if err := os.RemoveAll(p.tempDir); err != nil {
		p.Log.Error(err)
	}
	p.tempDir = ""
}

// PatchFile is what materializing the patch needs to know about one of its files: where
// the patch expects to find it, and where its content before the patch comes from.
type PatchFile struct {
	// Path is the name the patch knows the file by: for a renamed file, the name it had
	// before where the patch carries the rename, and the name it was renamed to where the
	// patch keeps only a content change and leaves the rename behind.
	Path string
	// ContentPath is where the file's content before the patch is to be found in the
	// commit the patch is built from — for a renamed file always the name it had there,
	// whatever the patch calls it.
	ContentPath string
}

// FilesInPatch says which files the patch touches, in a stable order, and where each of
// them comes from.
func (p *PatchBuilder) FilesInPatch() []PatchFile {
	fileInfoMap := p.snapshotFileInfoMap()

	filenames := maps.Keys(fileInfoMap)
	sort.Strings(filenames)

	files := make([]PatchFile, 0, len(filenames))
	for _, filename := range filenames {
		info := fileInfoMap[filename]
		if info.mode == UNSELECTED {
			continue
		}
		file := PatchFile{Path: filename, ContentPath: filename}
		if info.previousPath != "" {
			file.ContentPath = info.previousPath
			if info.mode == WHOLE {
				file.Path = info.previousPath
			}
		}
		files = append(files, file)
	}
	return files
}

func (p *PatchBuilder) PatchToApply(reverse bool, turnAddedFilesIntoDiffAgainstEmptyFile bool) string {
	var patch strings.Builder

	for filename, info := range p.snapshotFileInfoMap() {
		if info.mode == UNSELECTED {
			continue
		}

		patch.WriteString(p.RenderPatchForFile(RenderPatchForFileOpts{
			Filename:                               filename,
			PreviousPath:                           info.previousPath,
			Plain:                                  true,
			Reverse:                                reverse,
			TurnAddedFilesIntoDiffAgainstEmptyFile: turnAddedFilesIntoDiffAgainstEmptyFile,
		}))
	}

	return patch.String()
}

func (p *PatchBuilder) addFileWhole(info *fileInfo) {
	if info.mode != WHOLE {
		info.mode = WHOLE
		lineCount := len(strings.Split(info.diff, "\n"))
		// add every line index
		// TODO: add tests and then use lo.Range to simplify
		info.includedLineIndices = make([]int, lineCount)
		for i := range lineCount {
			info.includedLineIndices[i] = i
		}
	}
}

func (p *PatchBuilder) removeFile(info *fileInfo) {
	info.mode = UNSELECTED
	info.includedLineIndices = nil
}

func (p *PatchBuilder) AddFileWhole(filename string, previousPath string) error {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return err
	}

	p.generation++
	p.addFileWhole(info)

	return nil
}

func (p *PatchBuilder) RemoveFile(filename string, previousPath string) error {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return err
	}

	p.generation++
	p.removeFile(info)

	return nil
}

func (p *PatchBuilder) getFileInfo(filename string, previousPath string) (*fileInfo, error) {
	p.mutex.Lock()
	fileInfoMap := p.fileInfoMap
	from, to, reverse := p.From, p.To, p.reverse
	p.mutex.Unlock()

	info, ok := fileInfoMap[filename]
	if ok {
		return info, nil
	}

	diff, err := p.loadFileDiff(from, to, reverse, filename, previousPath)
	if err != nil {
		return nil, err
	}
	info = &fileInfo{
		mode:         UNSELECTED,
		diff:         diff,
		previousPath: previousPath,
	}

	fileInfoMap[filename] = info

	return info, nil
}

func (p *PatchBuilder) AddFileLineRange(filename string, previousPath string, lineIndices []int) error {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return err
	}
	p.generation++
	info.mode = PART
	info.includedLineIndices = lo.Union(info.includedLineIndices, lineIndices)

	return nil
}

func (p *PatchBuilder) RemoveFileLineRange(filename string, previousPath string, lineIndices []int) error {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return err
	}
	p.generation++
	info.mode = PART
	info.includedLineIndices, _ = lo.Difference(info.includedLineIndices, lineIndices)
	if len(info.includedLineIndices) == 0 {
		p.removeFile(info)
	}

	return nil
}

type RenderPatchForFileOpts struct {
	Filename                               string
	PreviousPath                           string
	Plain                                  bool
	Reverse                                bool
	TurnAddedFilesIntoDiffAgainstEmptyFile bool
}

func (p *PatchBuilder) RenderPatchForFile(opts RenderPatchForFileOpts) string {
	info, err := p.getFileInfo(opts.Filename, opts.PreviousPath)
	if err != nil {
		p.Log.Error(err)
		return ""
	}

	if info.mode == UNSELECTED {
		return ""
	}

	if info.mode == WHOLE && opts.Plain {
		// Use the whole diff (spares us parsing it and then formatting it).
		// TODO: see if this is actually noticeably faster.
		// The reverse flag is only for part patches so we're ignoring it here.
		return info.diff
	}

	patch := Parse(info.diff).
		Transform(TransformOpts{
			Reverse:                                opts.Reverse,
			TurnAddedFilesIntoDiffAgainstEmptyFile: opts.TurnAddedFilesIntoDiffAgainstEmptyFile,
			// For a partial selection of a renamed file we keep only the
			// content change and drop the rename, so that the rename stays in
			// the commit. A whole-file selection keeps the rename (and short-
			// circuits before this for plain output).
			StripRename:         info.mode == PART && info.previousPath != "",
			IncludedLineIndices: info.includedLineIndices,
		})

	if opts.Plain {
		return patch.FormatPlain()
	}
	return patch.FormatView(FormatViewOpts{})
}

func (p *PatchBuilder) renderEachFilePatch(plain bool) []string {
	fileInfoMap := p.snapshotFileInfoMap()

	// sort files by name then iterate through and render each patch
	filenames := maps.Keys(fileInfoMap)

	sort.Strings(filenames)
	patches := lo.Map(filenames, func(filename string, _ int) string {
		return p.RenderPatchForFile(RenderPatchForFileOpts{
			Filename:                               filename,
			PreviousPath:                           fileInfoMap[filename].previousPath,
			Plain:                                  plain,
			Reverse:                                false,
			TurnAddedFilesIntoDiffAgainstEmptyFile: true,
		})
	})
	output := lo.Filter(patches, func(patch string, _ int) bool {
		return patch != ""
	})

	return output
}

func (p *PatchBuilder) RenderAggregatedPatch(plain bool) string {
	return strings.Join(p.renderEachFilePatch(plain), "")
}

func (p *PatchBuilder) GetFileStatus(filename string, parent string) PatchStatus {
	p.mutex.Lock()
	to := p.To
	fileInfoMap := p.fileInfoMap
	p.mutex.Unlock()

	if parent != to {
		return UNSELECTED
	}

	info, ok := fileInfoMap[filename]
	if !ok {
		return UNSELECTED
	}

	return info.mode
}

// LineIdentity says which change line of a file is meant — the line number it has on
// the side it belongs to, and whether it is a deletion — without reference to where
// that line sits in the file's parsed diff.
//
// It is how a diff shown in the main view speaks about its lines: what a rendered row
// resolves to is a line of a file, while the index of that line in the diff depends on
// how much of the diff is being shown and in what order a renderer laid it out.
type LineIdentity struct {
	LineNumber int
	IsDeletion bool
}

// ChangeLineIndexByIdentity indexes a parsed diff's change lines by their identity. An
// addition is numbered in the new file and a deletion in the old one, which is what
// keeps two consecutive deletions — one new-file position between them — apart.
func ChangeLineIndexByIdentity(parsed *Patch) map[LineIdentity]int {
	byIdentity := map[LineIdentity]int{}
	for idx, line := range parsed.Lines() {
		switch {
		case line.IsAddition():
			byIdentity[LineIdentity{parsed.LineNumberOfLine(idx), false}] = idx
		case line.IsDeletion():
			byIdentity[LineIdentity{parsed.OldLineNumberOfLine(idx), true}] = idx
		}
	}
	return byIdentity
}

// ChangeLineIndicesForLines maps the given change lines of a parsed diff to their
// indices in it. A line that names no change line of the diff — a context line, or a
// line that isn't in the diff at all — contributes nothing.
func ChangeLineIndicesForLines(parsed *Patch, lines []LineIdentity) []int {
	byIdentity := ChangeLineIndexByIdentity(parsed)
	indices := make([]int, 0, len(lines))
	for _, line := range lines {
		if idx, ok := byIdentity[line]; ok {
			indices = append(indices, idx)
		}
	}
	return indices
}

// PatchLineIndicesForLines maps change lines of filename to their indices in that
// file's diff, which is what the patch is built in terms of.
func (p *PatchBuilder) PatchLineIndicesForLines(
	filename string, previousPath string, lines []LineIdentity,
) ([]int, error) {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return nil, err
	}

	return ChangeLineIndicesForLines(Parse(info.diff), lines), nil
}

// SelectionRepresentsWholeFile says whether lines select the entirety of a diff that
// consists of one solid block of additions or deletions, with no context. These are
// the added and deleted files whose file operation must travel with their contents.
func (p *PatchBuilder) SelectionRepresentsWholeFile(
	filename string, previousPath string, lines []LineIdentity,
) (bool, error) {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return false, err
	}

	parsed := Parse(info.diff)
	if !parsed.IsSingleHunkForWholeFile() {
		return false, nil
	}
	all := ChangeLineIndexByIdentity(parsed)
	selected := set.NewFromSlice(lines)
	return lo.EveryBy(maps.Keys(all), func(identity LineIdentity) bool {
		return selected.Includes(identity)
	}), nil
}

// IncludedLineIdentities says which change lines of filename are in the patch, as the
// identities a diff of that file shown anywhere can be compared against. Empty for a
// file that is no part of the patch.
func (p *PatchBuilder) IncludedLineIdentities(filename string) []LineIdentity {
	info, ok := p.snapshotFileInfoMap()[filename]
	if !ok || info.mode == UNSELECTED {
		return nil
	}

	included := set.NewFromSlice(info.includedLineIndices)
	identities := []LineIdentity{}
	for identity, idx := range ChangeLineIndexByIdentity(Parse(info.diff)) {
		if included.Includes(idx) {
			identities = append(identities, identity)
		}
	}
	return identities
}

// IncludedChangeLineIndices says which of filename's change lines are in the patch, as
// their indices in the file's diff and in the order the file has them.
//
// It is how a line of the patch as it is shown names the line of the diff it came from:
// all that can be said about a line of the patch is which of the file's changes it is,
// its line numbers being the patch's own — a patch that leaves an earlier addition out
// numbers everything after it differently from the diff it was built from.
func (p *PatchBuilder) IncludedChangeLineIndices(filename string) []int {
	info, ok := p.snapshotFileInfoMap()[filename]
	if !ok || info.mode == UNSELECTED {
		return nil
	}

	included := set.NewFromSlice(info.includedLineIndices)
	indices := []int{}
	for idx, line := range Parse(info.diff).Lines() {
		if (line.IsAddition() || line.IsDeletion()) && included.Includes(idx) {
			indices = append(indices, idx)
		}
	}
	return indices
}

func (p *PatchBuilder) GetFileIncLineIndices(filename string, previousPath string) ([]int, error) {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return nil, err
	}
	return info.includedLineIndices, nil
}

// clears the patch
func (p *PatchBuilder) Reset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.generation++
	p.removeTempDir()

	p.To = ""
	p.fileInfoMap = map[string]*fileInfo{}
}

func (p *PatchBuilder) Active() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.To != ""
}

func (p *PatchBuilder) IsEmpty() bool {
	for _, fileInfo := range p.snapshotFileInfoMap() {
		if fileInfo.mode == WHOLE || (fileInfo.mode == PART && len(fileInfo.includedLineIndices) > 0) {
			return false
		}
	}

	return true
}

// if any of these things change we'll need to reset and start a new patch
func (p *PatchBuilder) NewPatchRequired(from string, to string, reverse bool) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return from != p.From || to != p.To || reverse != p.reverse
}

func (p *PatchBuilder) AllFilesInPatch() []string {
	return lo.Keys(p.snapshotFileInfoMap())
}
