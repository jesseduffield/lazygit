package commands

import (
	"sort"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type fileInfo struct {
	mode                int // one of WHOLE/PART
	includedLineIndices []int
	diff                string
}

type applyPatchFunc func(patch string, flags ...string) error

// PatchManager manages the building of a patch for a commit to be applied to another commit (or the working tree, or removed from the current commit)
type PatchManager struct {
	CommitSha   string
	fileInfoMap map[string]*fileInfo
	Log         *logrus.Entry
	ApplyPatch  applyPatchFunc
}

// NewPatchManager returns a new PatchModifier
func NewPatchManager(log *logrus.Entry, applyPatch applyPatchFunc) *PatchManager {
	return &PatchManager{
		Log:        log,
		ApplyPatch: applyPatch,
	}
}

// NewPatchManager returns a new PatchModifier
func (p *PatchManager) Start(commitSha string, diffMap map[string]string) {
	p.CommitSha = commitSha
	p.fileInfoMap = map[string]*fileInfo{}
	for filename, diff := range diffMap {
		p.fileInfoMap[filename] = &fileInfo{
			mode: UNSELECTED,
			diff: diff,
		}
	}
}

func (p *PatchManager) AddFile(filename string) {
	p.fileInfoMap[filename].mode = WHOLE
	p.fileInfoMap[filename].includedLineIndices = nil
}

func (p *PatchManager) RemoveFile(filename string) {
	p.fileInfoMap[filename].mode = UNSELECTED
	p.fileInfoMap[filename].includedLineIndices = nil
}

func (p *PatchManager) ToggleFileWhole(filename string) {
	info := p.fileInfoMap[filename]
	switch info.mode {
	case UNSELECTED:
		p.AddFile(filename)
	case WHOLE:
		p.RemoveFile(filename)
	case PART:
		p.AddFile(filename)
	}
}

func getIndicesForRange(first, last int) []int {
	indices := []int{}
	for i := first; i <= last; i++ {
		indices = append(indices, i)
	}
	return indices
}

func (p *PatchManager) AddFileLineRange(filename string, firstLineIdx, lastLineIdx int) {
	info := p.fileInfoMap[filename]
	info.mode = PART
	info.includedLineIndices = utils.UnionInt(info.includedLineIndices, getIndicesForRange(firstLineIdx, lastLineIdx))
}

func (p *PatchManager) RemoveFileLineRange(filename string, firstLineIdx, lastLineIdx int) {
	info := p.fileInfoMap[filename]
	info.mode = PART
	info.includedLineIndices = utils.DifferenceInt(info.includedLineIndices, getIndicesForRange(firstLineIdx, lastLineIdx))
	if len(info.includedLineIndices) == 0 {
		p.RemoveFile(filename)
	}
}

func (p *PatchManager) RenderPlainPatchForFile(filename string, reverse bool, keepOriginalHeader bool) string {
	info := p.fileInfoMap[filename]
	if info == nil {
		return ""
	}

	switch info.mode {
	case WHOLE:
		// use the whole diff
		// the reverse flag is only for part patches so we're ignoring it here
		return info.diff
	case PART:
		// generate a new diff with just the selected lines
		m := NewPatchModifier(p.Log, filename, info.diff)
		return m.ModifiedPatchForLines(info.includedLineIndices, reverse, keepOriginalHeader)
	default:
		return ""
	}
}

func (p *PatchManager) RenderPatchForFile(filename string, plain bool, reverse bool, keepOriginalHeader bool) string {
	patch := p.RenderPlainPatchForFile(filename, reverse, keepOriginalHeader)
	if plain {
		return patch
	}
	parser, err := NewPatchParser(p.Log, patch)
	if err != nil {
		// swallowing for now
		return ""
	}
	// not passing included lines because we don't want to see them in the secondary panel
	return parser.Render(-1, -1, nil)
}

func (p *PatchManager) RenderEachFilePatch(plain bool) []string {
	// sort files by name then iterate through and render each patch
	filenames := make([]string, len(p.fileInfoMap))
	index := 0
	for filename := range p.fileInfoMap {
		filenames[index] = filename
		index++
	}

	sort.Strings(filenames)
	output := []string{}
	for _, filename := range filenames {
		patch := p.RenderPatchForFile(filename, plain, false, true)
		if patch != "" {
			output = append(output, patch)
		}
	}

	return output
}

func (p *PatchManager) RenderAggregatedPatchColored(plain bool) string {
	result := ""
	for _, patch := range p.RenderEachFilePatch(plain) {
		if patch != "" {
			result += patch + "\n"
		}
	}
	return result
}

func (p *PatchManager) GetFileStatus(filename string) int {
	info := p.fileInfoMap[filename]
	if info == nil {
		return UNSELECTED
	}
	return info.mode
}

func (p *PatchManager) GetFileIncLineIndices(filename string) []int {
	info := p.fileInfoMap[filename]
	if info == nil {
		return []int{}
	}
	return info.includedLineIndices
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
			if err = p.ApplyPatch(patch, applyFlags...); err != nil {
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
	p.CommitSha = ""
	p.fileInfoMap = map[string]*fileInfo{}
}

func (p *PatchManager) CommitSelected() bool {
	return p.CommitSha != ""
}

func (p *PatchManager) IsEmpty() bool {
	for _, fileInfo := range p.fileInfoMap {
		if fileInfo.mode == WHOLE || (fileInfo.mode == PART && len(fileInfo.includedLineIndices) > 0) {
			return false
		}
	}

	return true
}
