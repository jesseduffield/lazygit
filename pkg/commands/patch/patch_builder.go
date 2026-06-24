package patch

import (
	"os"
	"sort"
	"strings"

	"github.com/jesseduffield/generics/maps"
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
	loadFileDiffFunc func(from string, to string, reverse bool, filename string, previousPath string, plain bool) (string, error)
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

	// newTempDir creates a fresh temp dir for the current patch, into which the patch is
	// materialized as two file trees so it can be rendered through any pager (see
	// PatchCommands.WriteCustomPatchDiffTrees). The patch builder owns its lifetime: a dir
	// is created on Start and removed on Reset. Nil in tests that don't render the patch.
	newTempDir func() (string, error)
	tempDir    string

	// generation is bumped on every change to the patch's contents, so a consumer that
	// materializes the patch (the secondary pane's two file trees) can tell when it's stale
	// and rebuild only then — covering every path that mutates the patch (the focused main
	// view and the old explorer alike) without rebuilding on mere navigation.
	generation int
}

func NewPatchBuilder(log *logrus.Entry, loadFileDiff loadFileDiffFunc, newTempDir func() (string, error)) *PatchBuilder {
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
	p.removeTempDir()
	if p.newTempDir != nil {
		if dir, err := p.newTempDir(); err != nil {
			p.Log.Error(err)
		} else {
			p.tempDir = dir
		}
	}

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

// TempDir is the directory the current patch is materialized into for rendering, or "" if
// none was created. See PatchCommands.WriteCustomPatchDiffTrees.
func (p *PatchBuilder) TempDir() string {
	return p.tempDir
}

// Generation is bumped each time the patch's contents change; see the field comment.
func (p *PatchBuilder) Generation() int {
	return p.generation
}

func (p *PatchBuilder) removeTempDir() {
	if p.tempDir != "" {
		_ = os.RemoveAll(p.tempDir)
		p.tempDir = ""
	}
}

// ActiveFilenames returns the files currently part of the patch (mode != UNSELECTED), in
// sorted order — the files to materialize when rendering the patch.
func (p *PatchBuilder) ActiveFilenames() []string {
	filenames := make([]string, 0, len(p.fileInfoMap))
	for filename, info := range p.fileInfoMap {
		if info.mode != UNSELECTED {
			filenames = append(filenames, filename)
		}
	}
	sort.Strings(filenames)
	return filenames
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

	diff, err := p.loadFileDiff(from, to, reverse, filename, previousPath, true)
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

func (p *PatchBuilder) GetFileIncLineIndices(filename string, previousPath string) ([]int, error) {
	info, err := p.getFileInfo(filename, previousPath)
	if err != nil {
		return nil, err
	}
	return info.includedLineIndices, nil
}

// LineIdentity identifies a change line (an addition or deletion) by its file line
// number and whether it's a deletion, independently of the line's index in the parsed
// patch. It is the identity the diff-line metadata resolves a rendered row to (see
// types.DiffLineInfo.PatchSelectLine), and lets the focused main view toggle patch
// membership and drive the inclusion gutter without dealing in patch-line indices —
// which differ between the raw diff and however a pager renders it.
type LineIdentity struct {
	LineNumber int
	IsDeletion bool
}

// changeLineIndexByIdentity scans a parsed diff and returns, for each change line, its
// index in the patch keyed by the line's identity. An addition is keyed by its new-file
// line number, a deletion by its old-file line number, so each change line has a
// distinct identity (two consecutive deletions share a new-file number but differ in
// the old-file one).
func changeLineIndexByIdentity(parsed *Patch) map[LineIdentity]int {
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

// PatchLineIndicesForLines maps the given change-line identities to their indices in
// filename's parsed diff — the index form that AddFileLineRange / RemoveFileLineRange
// and GetFileIncLineIndices work in. Identities that don't correspond to a change line
// (e.g. a context line) are skipped. It is how the focused main view, which knows a
// selection only as metadata identities, drives patch building.
func (p *PatchBuilder) PatchLineIndicesForLines(filename string, lines []LineIdentity) ([]int, error) {
	info, err := p.getFileInfo(filename, "")
	if err != nil {
		return nil, err
	}
	byIdentity := changeLineIndexByIdentity(Parse(info.diff))
	indices := make([]int, 0, len(lines))
	for _, line := range lines {
		if idx, ok := byIdentity[line]; ok {
			indices = append(indices, idx)
		}
	}
	return indices, nil
}

// IncludedChangeLineIndices returns the patch-line indices of the change lines (additions
// and deletions) currently included in the patch for filename, in ascending order. These
// are exactly the change lines the aggregated patch renders for the file, in the same
// order, so the k-th change line shown in the custom-patch (secondary) view corresponds to
// the k-th index here. That correspondence lets the focused main view remove a selection
// from the patch by its ordinal among the shown change lines, sidestepping the line-number
// renumbering the aggregated patch applies (which makes matching by identity unreliable for
// additions). Empty when the file isn't part of the patch.
func (p *PatchBuilder) IncludedChangeLineIndices(filename string) []int {
	info, ok := p.fileInfoMap[filename]
	if !ok || info.mode == UNSELECTED {
		return nil
	}
	lines := Parse(info.diff).Lines()
	included := append([]int{}, info.includedLineIndices...)
	sort.Ints(included)
	result := make([]int, 0, len(included))
	for _, idx := range included {
		if idx >= 0 && idx < len(lines) && (lines[idx].IsAddition() || lines[idx].IsDeletion()) {
			result = append(result, idx)
		}
	}
	return result
}

// IncludedLineIdentities returns the identities of the change lines currently included
// in the patch for filename — the identity space the inclusion gutter matches rendered
// rows against. Empty when the file isn't part of the patch.
func (p *PatchBuilder) IncludedLineIdentities(filename string) []LineIdentity {
	info, ok := p.fileInfoMap[filename]
	if !ok || info.mode == UNSELECTED {
		return nil
	}
	includedIdx := make(map[int]bool, len(info.includedLineIndices))
	for _, idx := range info.includedLineIndices {
		includedIdx[idx] = true
	}
	var identities []LineIdentity
	for identity, idx := range changeLineIndexByIdentity(Parse(info.diff)) {
		if includedIdx[idx] {
			identities = append(identities, identity)
		}
	}
	return identities
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
