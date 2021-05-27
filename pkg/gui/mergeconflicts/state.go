package mergeconflicts

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/golang-collections/collections/stack"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type Selection int

const (
	TOP Selection = iota
	BOTTOM
	BOTH
)

// mergeConflict : A git conflict with a start middle and end corresponding to line
// numbers in the file where the conflict bars appear
type mergeConflict struct {
	start  int
	middle int
	end    int
}

type State struct {
	sync.Mutex
	conflictIndex int
	conflictTop   bool
	conflicts     []*mergeConflict
	EditHistory   *stack.Stack
}

func NewState() *State {
	return &State{
		Mutex:         sync.Mutex{},
		conflictIndex: 0,
		conflictTop:   true,
		conflicts:     []*mergeConflict{},
		EditHistory:   stack.New(),
	}
}

func (s *State) SelectTopOption() {
	s.conflictTop = true
}

func (s *State) SelectBottomOption() {
	s.conflictTop = false
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

func findConflicts(content string) []*mergeConflict {
	conflicts := make([]*mergeConflict, 0)

	if content == "" {
		return conflicts
	}

	var newConflict *mergeConflict
	for i, line := range utils.SplitLines(content) {
		trimmedLine := strings.TrimPrefix(line, "++")
		switch trimmedLine {
		case "<<<<<<< HEAD", "<<<<<<< MERGE_HEAD", "<<<<<<< Updated upstream", "<<<<<<< ours":
			newConflict = &mergeConflict{start: i}
		case "=======":
			newConflict.middle = i
		default:
			if strings.HasPrefix(trimmedLine, ">>>>>>> ") {
				newConflict.end = i
				conflicts = append(conflicts, newConflict)
			}
		}

	}
	return conflicts
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
	if s.conflictTop {
		return TOP
	} else {
		return BOTTOM
	}
}

func (s *State) IsFinalConflict() bool {
	return len(s.conflicts) == 1
}

func (s *State) Reset() {
	s.EditHistory = stack.New()
}

func (s *State) GetConflictMiddle() int {
	return s.currentConflict().middle
}

func (s *State) ContentAfterConflictResolve(path string, selection Selection) (bool, string, error) {
	conflict := s.currentConflict()
	if conflict == nil {
		return false, "", nil
	}

	file, err := os.Open(path)
	if err != nil {
		return false, "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	content := ""
	for i := 0; true; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if !isIndexToDelete(i, conflict, selection) {
			content += line
		}
	}

	return true, content, nil
}

func isIndexToDelete(i int, conflict *mergeConflict, selection Selection) bool {
	return i == conflict.middle ||
		i == conflict.start ||
		i == conflict.end ||
		selection != BOTH &&
			(selection == BOTTOM && i > conflict.start && i < conflict.middle) ||
		(selection == TOP && i > conflict.middle && i < conflict.end)
}
