package mergeconflicts

import (
	"sync"

	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type Selection int

const (
	TOP Selection = iota
	MIDDLE
	BOTTOM
	BOTH
)

// mergeConflict : A git conflict with a start, ancestor (if exists), target, and end corresponding to line
// numbers in the file where the conflict markers appear
type mergeConflict struct {
	start    int
	ancestor int
	target   int
	end      int
}

type State struct {
	sync.Mutex
	conflictIndex     int
	conflictSelection Selection
	conflicts         []*mergeConflict
	EditHistory       *stack.Stack
}

func NewState() *State {
	return &State{
		Mutex:             sync.Mutex{},
		conflictIndex:     0,
		conflictSelection: TOP,
		conflicts:         []*mergeConflict{},
		EditHistory:       stack.New(),
	}
}

func (s *State) SelectPrevConflictHunk() {
	switch s.conflictSelection {
	case MIDDLE:
		s.conflictSelection = TOP
	case BOTTOM:
		if s.currentConflict().ancestor >= 0 {
			s.conflictSelection = MIDDLE
		} else {
			s.conflictSelection = TOP
		}
	}
}

func (s *State) SelectNextConflictHunk() {
	switch s.conflictSelection {
	case TOP:
		if s.currentConflict().ancestor >= 0 {
			s.conflictSelection = MIDDLE
		} else {
			s.conflictSelection = BOTTOM
		}
	case MIDDLE:
		s.conflictSelection = BOTTOM
	}
}

func (s *State) SelectNextConflict() {
	if s.conflictIndex < len(s.conflicts)-1 {
		s.conflictIndex++
	}
}

func (s *State) SelectPrevConflict() {
	if s.conflictIndex > 0 {
		s.conflictIndex--
	}
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

	if s.conflictIndex > len(s.conflicts)-1 {
		s.conflictIndex = len(s.conflicts) - 1
	} else if s.conflictIndex < 0 {
		s.conflictIndex = 0
	}
}

func (s *State) NoConflicts() bool {
	return len(s.conflicts) == 0
}

func (s *State) Selection() Selection {
	return s.conflictSelection
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
	isMarkerLine :=
		i == conflict.start ||
			i == conflict.ancestor ||
			i == conflict.target ||
			i == conflict.end

	var isWantedContent bool
	switch selection {
	case TOP:
		if conflict.ancestor >= 0 {
			isWantedContent = conflict.start < i && i < conflict.ancestor
		} else {
			isWantedContent = conflict.start < i && i < conflict.target
		}
	case MIDDLE:
		isWantedContent = conflict.ancestor < i && i < conflict.target
	case BOTTOM:
		isWantedContent = conflict.target < i && i < conflict.end
	}
	return isMarkerLine || !isWantedContent
}
