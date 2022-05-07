package mergeconflicts

import (
	"sync"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

// State represents the selection state of the merge conflict context.
type State struct {
	sync.Mutex

	// path of the file with the conflicts
	path string

	// This is a stack of the file content. It is used to undo changes.
	// The last item is the current file content.
	contents []string

	conflicts []*mergeConflict
	// this is the index of the above `conflicts` field which is currently selected
	conflictIndex int

	// this is the index of the selected conflict's available selections slice e.g. [TOP, MIDDLE, BOTTOM]
	// We use this to know which hunk of the conflict is selected.
	selectionIndex int
}

func NewState() *State {
	return &State{
		Mutex:          sync.Mutex{},
		conflictIndex:  0,
		selectionIndex: 0,
		conflicts:      []*mergeConflict{},
		contents:       []string{},
	}
}

func (s *State) setConflictIndex(index int) {
	if len(s.conflicts) == 0 {
		s.conflictIndex = 0
	} else {
		s.conflictIndex = utils.Clamp(index, 0, len(s.conflicts)-1)
	}
	s.setSelectionIndex(s.selectionIndex)
}

func (s *State) setSelectionIndex(index int) {
	if selections := s.availableSelections(); len(selections) != 0 {
		s.selectionIndex = utils.Clamp(index, 0, len(selections)-1)
	}
}

func (s *State) SelectNextConflictHunk() {
	s.setSelectionIndex(s.selectionIndex + 1)
}

func (s *State) SelectPrevConflictHunk() {
	s.setSelectionIndex(s.selectionIndex - 1)
}

func (s *State) SelectNextConflict() {
	s.setConflictIndex(s.conflictIndex + 1)
}

func (s *State) SelectPrevConflict() {
	s.setConflictIndex(s.conflictIndex - 1)
}

func (s *State) currentConflict() *mergeConflict {
	if len(s.conflicts) == 0 {
		return nil
	}

	return s.conflicts[s.conflictIndex]
}

// this is for starting a new merge conflict session
func (s *State) SetContent(content string, path string) {
	if content == s.GetContent() && path == s.path {
		return
	}

	s.path = path
	s.contents = []string{}
	s.PushContent(content)
}

// this is for when you've resolved a conflict. This allows you to undo to a previous
// state
func (s *State) PushContent(content string) {
	s.contents = append(s.contents, content)
	s.setConflicts(findConflicts(content))
}

func (s *State) GetContent() string {
	if len(s.contents) == 0 {
		return ""
	}

	return s.contents[len(s.contents)-1]
}

func (s *State) GetPath() string {
	return s.path
}

func (s *State) Undo() bool {
	if len(s.contents) <= 1 {
		return false
	}

	s.contents = s.contents[:len(s.contents)-1]

	newContent := s.GetContent()
	// We could be storing the old conflicts and selected index on a stack too.
	s.setConflicts(findConflicts(newContent))

	return true
}

func (s *State) setConflicts(conflicts []*mergeConflict) {
	s.conflicts = conflicts
	s.setConflictIndex(s.conflictIndex)
}

func (s *State) NoConflicts() bool {
	return len(s.conflicts) == 0
}

func (s *State) Selection() Selection {
	if selections := s.availableSelections(); len(selections) > 0 {
		return selections[s.selectionIndex]
	}
	return TOP
}

func (s *State) availableSelections() []Selection {
	if conflict := s.currentConflict(); conflict != nil {
		return availableSelections(conflict)
	}
	return nil
}

func (s *State) AllConflictsResolved() bool {
	return len(s.conflicts) == 0
}

func (s *State) Reset() {
	s.contents = []string{}
	s.path = ""
}

func (s *State) Active() bool {
	return s.path != ""
}

func (s *State) GetConflictMiddle() int {
	currentConflict := s.currentConflict()

	if currentConflict == nil {
		return 0
	}

	return currentConflict.target
}

func (s *State) ContentAfterConflictResolve(selection Selection) (bool, string, error) {
	conflict := s.currentConflict()
	if conflict == nil {
		return false, "", nil
	}

	content := ""
	err := utils.ForEachLineInFile(s.path, func(line string, i int) {
		if selection.isIndexToKeep(conflict, i) {
			content += line
		}
	})
	if err != nil {
		return false, "", err
	}

	return true, content, nil
}

func (s *State) GetSelectedLine() int {
	conflict := s.currentConflict()
	if conflict == nil {
		return 1
	}
	selection := s.Selection()
	startIndex, _ := selection.bounds(conflict)
	return startIndex + 1
}
