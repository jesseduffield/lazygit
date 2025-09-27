package patch

import (
	"sort"
	"strings"

	"github.com/jesseduffield/generics/maps"
	"github.com/samber/lo"
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
}

type (
	loadFileDiffFunc func(from string, to string, reverse bool, filename string, plain bool) (string, error)
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

	// loadFileDiff loads the diff of a file, for a given to (typically a commit hash)
	loadFileDiff loadFileDiffFunc
}

func NewPatchBuilder(log *logrus.Entry, loadFileDiff loadFileDiffFunc) *PatchBuilder {
	return &PatchBuilder{
		Log:          log,
		loadFileDiff: loadFileDiff,
	}
}

func (p *PatchBuilder) Start(from, to string, reverse bool, canRebase bool) {
	p.To = to
	p.From = from
	p.reverse = reverse
	p.CanRebase = canRebase
	p.fileInfoMap = map[string]*fileInfo{}
}

func (p *PatchBuilder) PatchToApply(reverse bool, turnAddedFilesIntoDiffAgainstEmptyFile bool) string {
	patch := ""

	for filename, info := range p.fileInfoMap {
		if info.mode == UNSELECTED {
			continue
		}

		patch += p.RenderPatchForFile(RenderPatchForFileOpts{
			Filename:                               filename,
			Plain:                                  true,
			Reverse:                                reverse,
			TurnAddedFilesIntoDiffAgainstEmptyFile: turnAddedFilesIntoDiffAgainstEmptyFile,
		})
	}

	return patch
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

func (p *PatchBuilder) AddFileWhole(filename string) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}

	p.addFileWhole(info)

	return nil
}

func (p *PatchBuilder) RemoveFile(filename string) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}

	p.removeFile(info)

	return nil
}

func (p *PatchBuilder) getFileInfo(filename string) (*fileInfo, error) {
	info, ok := p.fileInfoMap[filename]
	if ok {
		return info, nil
	}

	diff, err := p.loadFileDiff(p.From, p.To, p.reverse, filename, true)
	if err != nil {
		return nil, err
	}
	info = &fileInfo{
		mode: UNSELECTED,
		diff: diff,
	}

	p.fileInfoMap[filename] = info

	return info, nil
}

func (p *PatchBuilder) AddFileLineRange(filename string, lineIndices []int) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}
	info.mode = PART
	info.includedLineIndices = lo.Union(info.includedLineIndices, lineIndices)

	return nil
}

func (p *PatchBuilder) RemoveFileLineRange(filename string, lineIndices []int) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}
	info.mode = PART
	info.includedLineIndices, _ = lo.Difference(info.includedLineIndices, lineIndices)
	if len(info.includedLineIndices) == 0 {
		p.removeFile(info)
	}

	return nil
}

type RenderPatchForFileOpts struct {
	Filename                               string
	Plain                                  bool
	Reverse                                bool
	TurnAddedFilesIntoDiffAgainstEmptyFile bool
}

func (p *PatchBuilder) RenderPatchForFile(opts RenderPatchForFileOpts) string {
	info, err := p.getFileInfo(opts.Filename)
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
			IncludedLineIndices:                    info.includedLineIndices,
		})

	if opts.Plain {
		return patch.FormatPlain()
	}
	return patch.FormatView(FormatViewOpts{})
}

func (p *PatchBuilder) renderEachFilePatch(plain bool) []string {
	// sort files by name then iterate through and render each patch
	filenames := maps.Keys(p.fileInfoMap)

	sort.Strings(filenames)
	patches := lo.Map(filenames, func(filename string, _ int) string {
		return p.RenderPatchForFile(RenderPatchForFileOpts{
			Filename:                               filename,
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
	if parent != p.To {
		return UNSELECTED
	}

	info, ok := p.fileInfoMap[filename]
	if !ok {
		return UNSELECTED
	}

	return info.mode
}

func (p *PatchBuilder) GetFileIncLineIndices(filename string) ([]int, error) {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return nil, err
	}
	return info.includedLineIndices, nil
}

// clears the patch
func (p *PatchBuilder) Reset() {
	p.To = ""
	p.fileInfoMap = map[string]*fileInfo{}
}

func (p *PatchBuilder) Active() bool {
	return p.To != ""
}

func (p *PatchBuilder) IsEmpty() bool {
	for _, fileInfo := range p.fileInfoMap {
		if fileInfo.mode == WHOLE || (fileInfo.mode == PART && len(fileInfo.includedLineIndices) > 0) {
			return false
		}
	}

	return true
}

// if any of these things change we'll need to reset and start a new patch
func (p *PatchBuilder) NewPatchRequired(from string, to string, reverse bool) bool {
	return from != p.From || to != p.To || reverse != p.reverse
}

func (p *PatchBuilder) AllFilesInPatch() []string {
	return lo.Keys(p.fileInfoMap)
}
