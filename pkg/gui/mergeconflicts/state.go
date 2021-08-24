package mergeconflicts

import (
	"sync"

	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func clamp(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

type Selection int

const (
	TOP Selection = iota
	MIDDLE
	BOTTOM
	ALL
)

// mergeConflict : A git conflict with a start, ancestor (if exists), target, and end corresponding to line
// numbers in the file where the conflict markers appear
type mergeConflict struct {
	start    int
	ancestor int
	target   int
	end      int
}

func (c *mergeConflict) hasAncestor() bool {
	return c.ancestor >= 0
}

func (c *mergeConflict) availableSelections() []Selection {
	if c.hasAncestor() {
		return []Selection{TOP, MIDDLE, BOTTOM}
	} else {
		return []Selection{TOP, BOTTOM}
	}
}

type State struct {
	sync.Mutex
	conflictIndex  int
	selectionIndex int
	conflicts      []*mergeConflict
	EditHistory    *stack.Stack
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
		s.conflictIndex = clamp(index, 0, len(s.conflicts)-1)
	} else {
		s.conflictIndex = 0
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
		return conflict.availableSelections()
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

func (s *State) ContentAfterConflictResolve(path string, selection Selection) (bool, string, error) {
	conflict := s.currentConflict()
	if conflict == nil {
		return false, "", nil
	}

	content := ""
	err := utils.ForEachLineInFile(path, func(line string, i int) {
		if !isIndexToDelete(i, conflict, selection) {
			content += line
		}
	})

	if err != nil {
		return false, "", err
	}

	return true, content, nil
}

func isIndexToDelete(i int, conflict *mergeConflict, selection Selection) bool {
	if i < conflict.start || conflict.end < i {
		return false
	}

	isMarkerLine :=
		i == conflict.start ||
			i == conflict.ancestor ||
			i == conflict.target ||
			i == conflict.end

	var isWantedContent bool
	switch selection {
	case TOP:
		if conflict.hasAncestor() {
			isWantedContent = conflict.start < i && i < conflict.ancestor
		} else {
			isWantedContent = conflict.start < i && i < conflict.target
		}
	case MIDDLE:
		isWantedContent = conflict.ancestor < i && i < conflict.target
	case BOTTOM:
		isWantedContent = conflict.target < i && i < conflict.end
	case ALL:
		isWantedContent = true
	}

	return isMarkerLine || !isWantedContent
}
