package mergeconflicts

import (
	"sync"

	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type State struct {
	sync.Mutex

	conflicts []*mergeConflict
	// this is the index of the above `conflicts` field which is currently selected
	conflictIndex int

	// this is the index of the selected conflict's available selections slice e.g. [TOP, MIDDLE, BOTTOM]
	// We use this to know which hunk of the conflict is selected.
	selectionIndex int

	// this allows us to undo actions
	EditHistory *stack.Stack
}

func NewState() *State {
	return &State{
		Mutex:          sync.Mutex{},
		conflictIndex:  0,
		selectionIndex: 0,
		conflicts:      []*mergeConflict{},
		EditHistory:    stack.New(),
	}
}

func (s *State) setConflictIndex(index int) {
	if len(s.conflicts) == 0 {
		s.conflictIndex = 0
	} else {
		s.conflictIndex = clamp(index, 0, len(s.conflicts)-1)
	}
	s.setSelectionIndex(s.selectionIndex)
}

func (s *State) setSelectionIndex(index int) {
	if selections := s.availableSelections(); len(selections) != 0 {
		s.selectionIndex = clamp(index, 0, len(selections)-1)
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

func (s *State) PushFileSnapshot(content string) {
	s.EditHistory.Push(content)
}

func (s *State) PopFileSnapshot() (string, bool) {
	if s.EditHistory.Len() == 0 {
		return "", false
	}

	return s.EditHistory.Pop().(string), true
}

func (s *State) currentConflict() *mergeConflict {
	if len(s.conflicts) == 0 {
		return nil
	}

	return s.conflicts[s.conflictIndex]
}

func (s *State) SetConflictsFromCat(cat string) {
	s.setConflicts(findConflicts(cat))
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

func (s *State) IsFinalConflict() bool {
	return len(s.conflicts) == 1
}

func (s *State) Reset() {
	s.EditHistory = stack.New()
}

func (s *State) GetConflictMiddle() int {
	return s.currentConflict().target
}

func (s *State) ContentAfterConflictResolve(
	path string,
	selection Selection,
) (bool, string, error) {
	conflict := s.currentConflict()
	if conflict == nil {
		return false, "", nil
	}

	content := ""
	err := utils.ForEachLineInFile(path, func(line string, i int) {
		if selection.isIndexToKeep(conflict, i) {
			content += line
		}
	})

	if err != nil {
		return false, "", err
	}

	return true, content, nil
}

func clamp(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}
