package patch

import (
	"sort"
	"strings"

	"github.com/jesseduffield/generics/maps"
	"github.com/jesseduffield/generics/slices"
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
	applyPatchFunc   func(patch string, flags ...string) error
	loadFileDiffFunc func(from string, to string, reverse bool, filename string, plain bool) (string, error)
)

// PatchManager manages the building of a patch for a commit to be applied to another commit (or the working tree, or removed from the current commit). We also support building patches from things like stashes, for which there is less flexibility
type PatchManager struct {
	// To is the commit sha if we're dealing with files of a commit, or a stash ref for a stash
	To      string
	From    string
	reverse bool

	// CanRebase tells us whether we're allowed to modify our commits. CanRebase should be true for commits of the currently checked out branch and false for everything else
	// TODO: move this out into a proper mode struct in the gui package: it doesn't really belong here
	CanRebase bool

	// fileInfoMap starts empty but you add files to it as you go along
	fileInfoMap map[string]*fileInfo
	Log         *logrus.Entry
	applyPatch  applyPatchFunc

	// loadFileDiff loads the diff of a file, for a given to (typically a commit SHA)
	loadFileDiff loadFileDiffFunc
}

// NewPatchManager returns a new PatchManager
func NewPatchManager(log *logrus.Entry, applyPatch applyPatchFunc, loadFileDiff loadFileDiffFunc) *PatchManager {
	return &PatchManager{
		Log:          log,
		applyPatch:   applyPatch,
		loadFileDiff: loadFileDiff,
	}
}

// NewPatchManager returns a new PatchManager
func (p *PatchManager) Start(from, to string, reverse bool, canRebase bool) {
	p.To = to
	p.From = from
	p.reverse = reverse
	p.CanRebase = canRebase
	p.fileInfoMap = map[string]*fileInfo{}
}

func (p *PatchManager) addFileWhole(info *fileInfo) {
	info.mode = WHOLE
	lineCount := len(strings.Split(info.diff, "\n"))
	// add every line index
	// TODO: add tests and then use lo.Range to simplify
	info.includedLineIndices = make([]int, lineCount)
	for i := 0; i < lineCount; i++ {
		info.includedLineIndices[i] = i
	}
}

func (p *PatchManager) removeFile(info *fileInfo) {
	info.mode = UNSELECTED
	info.includedLineIndices = nil
}

func (p *PatchManager) AddFileWhole(filename string) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}

	p.addFileWhole(info)

	return nil
}

func (p *PatchManager) RemoveFile(filename string) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}

	p.removeFile(info)

	return nil
}

func getIndicesForRange(first, last int) []int {
	indices := []int{}
	for i := first; i <= last; i++ {
		indices = append(indices, i)
	}
	return indices
}

func (p *PatchManager) getFileInfo(filename string) (*fileInfo, error) {
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

func (p *PatchManager) AddFileLineRange(filename string, firstLineIdx, lastLineIdx int) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}
	info.mode = PART
	info.includedLineIndices = lo.Union(info.includedLineIndices, getIndicesForRange(firstLineIdx, lastLineIdx))

	return nil
}

func (p *PatchManager) RemoveFileLineRange(filename string, firstLineIdx, lastLineIdx int) error {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return err
	}
	info.mode = PART
	info.includedLineIndices, _ = lo.Difference(info.includedLineIndices, getIndicesForRange(firstLineIdx, lastLineIdx))
	if len(info.includedLineIndices) == 0 {
		p.removeFile(info)
	}

	return nil
}

func (p *PatchManager) renderPlainPatchForFile(filename string, reverse bool, keepOriginalHeader bool) string {
	info, err := p.getFileInfo(filename)
	if err != nil {
		p.Log.Error(err)
		return ""
	}

	switch info.mode {
	case WHOLE:
		// use the whole diff
		// the reverse flag is only for part patches so we're ignoring it here
		return info.diff
	case PART:
		// generate a new diff with just the selected lines
		return ModifiedPatchForLines(p.Log, filename, info.diff, info.includedLineIndices, reverse, keepOriginalHeader)
	default:
		return ""
	}
}

func (p *PatchManager) RenderPatchForFile(filename string, plain bool, reverse bool, keepOriginalHeader bool) string {
	patch := p.renderPlainPatchForFile(filename, reverse, keepOriginalHeader)
	if plain {
		return patch
	}
	parser := NewPatchParser(p.Log, patch)

	// not passing included lines because we don't want to see them in the secondary panel
	return parser.Render(false, -1, -1, nil)
}

func (p *PatchManager) renderEachFilePatch(plain bool) []string {
	// sort files by name then iterate through and render each patch
	filenames := maps.Keys(p.fileInfoMap)

	sort.Strings(filenames)
	patches := slices.Map(filenames, func(filename string) string {
		return p.RenderPatchForFile(filename, plain, false, true)
	})
	output := slices.Filter(patches, func(patch string) bool {
		return patch != ""
	})

	return output
}

func (p *PatchManager) RenderAggregatedPatchColored(plain bool) string {
	result := ""
	for _, patch := range p.renderEachFilePatch(plain) {
		if patch != "" {
			result += patch + "\n"
		}
	}
	return result
}

func (p *PatchManager) GetFileStatus(filename string, parent string) PatchStatus {
	if parent != p.To {
		return UNSELECTED
	}

	info, ok := p.fileInfoMap[filename]
	if !ok {
		return UNSELECTED
	}

	return info.mode
}

func (p *PatchManager) GetFileIncLineIndices(filename string) ([]int, error) {
	info, err := p.getFileInfo(filename)
	if err != nil {
		return nil, err
	}
	return info.includedLineIndices, nil
}

func (p *PatchManager) ApplyPatches(reverse bool) error {
	// for whole patches we'll apply the patch in reverse
	// but for part patches we'll apply a reverse patch forwards
	for filename, info := range p.fileInfoMap {
		if info.mode == UNSELECTED {
			continue
		}

		applyFlags := []string{"index", "3way"}
		reverseOnGenerate := false
		if reverse {
			if info.mode == WHOLE {
				applyFlags = append(applyFlags, "reverse")
			} else {
				reverseOnGenerate = true
			}
		}

		var err error
		// first run we try with the original header, then without
		for _, keepOriginalHeader := range []bool{true, false} {
			patch := p.RenderPatchForFile(filename, true, reverseOnGenerate, keepOriginalHeader)
			if patch == "" {
				continue
			}
			if err = p.applyPatch(patch, applyFlags...); err != nil {
				continue
			}
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// clears the patch
func (p *PatchManager) Reset() {
	p.To = ""
	p.fileInfoMap = map[string]*fileInfo{}
}

func (p *PatchManager) Active() bool {
	return p.To != ""
}

func (p *PatchManager) IsEmpty() bool {
	for _, fileInfo := range p.fileInfoMap {
		if fileInfo.mode == WHOLE || (fileInfo.mode == PART && len(fileInfo.includedLineIndices) > 0) {
			return false
		}
	}

	return true
}

// if any of these things change we'll need to reset and start a new patch
func (p *PatchManager) NewPatchRequired(from string, to string, reverse bool) bool {
	return from != p.From || to != p.To || reverse != p.reverse
}
