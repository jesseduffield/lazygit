package graph

import (
	"cmp"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type PipeKind uint8

const (
	TERMINATES PipeKind = iota
	STARTS
	CONTINUES
)

type Pipe struct {
	fromHash *string
	toHash   *string
	style    *style.TextStyle
	fromPos  int16
	toPos    int16
	kind     PipeKind
}

var (
	highlightStyle      = style.FgLightWhite.SetBold()
	EmptyTreeCommitHash = models.EmptyTreeCommitHash
	StartCommitHash     = "START"
)

func (self Pipe) left() int16 {
	return min(self.fromPos, self.toPos)
}

func (self Pipe) right() int16 {
	return max(self.fromPos, self.toPos)
}

func RenderCommitGraph(commits []*models.Commit, selectedCommitHashPtr *string, getStyle func(c *models.Commit) *style.TextStyle) []string {
	pipeSets := GetPipeSets(commits, getStyle)
	if len(pipeSets) == 0 {
		return nil
	}

	lines := RenderAux(pipeSets, commits, selectedCommitHashPtr)

	return lines
}

func GetPipeSets(commits []*models.Commit, getStyle func(c *models.Commit) *style.TextStyle) [][]Pipe {
	if len(commits) == 0 {
		return nil
	}

	pipes := []Pipe{{fromPos: 0, toPos: 0, fromHash: &StartCommitHash, toHash: commits[0].HashPtr(), kind: STARTS, style: &style.FgDefault}}

	return lo.Map(commits, func(commit *models.Commit, _ int) []Pipe {
		pipes = getNextPipes(pipes, commit, getStyle)
		return pipes
	})
}

func RenderAux(pipeSets [][]Pipe, commits []*models.Commit, selectedCommitHashPtr *string) []string {
	maxProcs := runtime.GOMAXPROCS(0)

	// splitting up the rendering of the graph into multiple goroutines allows us to render the graph in parallel
	chunks := make([][]string, maxProcs)
	perProc := len(pipeSets) / maxProcs

	wg := sync.WaitGroup{}
	wg.Add(maxProcs)

	for i := range maxProcs {
		go func() {
			from := i * perProc
			to := (i + 1) * perProc
			if i == maxProcs-1 {
				to = len(pipeSets)
			}
			innerLines := make([]string, 0, to-from)
			for j, pipeSet := range pipeSets[from:to] {
				k := from + j
				var prevCommit *models.Commit
				if k > 0 {
					prevCommit = commits[k-1]
				}
				line := renderPipeSet(pipeSet, selectedCommitHashPtr, prevCommit)
				innerLines = append(innerLines, line)
			}
			chunks[i] = innerLines
			wg.Done()
		}()
	}

	wg.Wait()

	return lo.Flatten(chunks)
}

func getNextPipes(prevPipes []Pipe, commit *models.Commit, getStyle func(c *models.Commit) *style.TextStyle) []Pipe {
	maxPos := int16(0)
	for _, pipe := range prevPipes {
		if pipe.toPos > maxPos {
			maxPos = pipe.toPos
		}
	}

	// a pipe that terminated in the previous line has no bearing on the current line
	// so we'll filter those out
	currentPipes := lo.Filter(prevPipes, func(pipe Pipe, _ int) bool {
		return pipe.kind != TERMINATES
	})

	newPipes := make([]Pipe, 0, len(currentPipes)+len(commit.ParentPtrs()))
	// start by assuming that we've got a brand new commit not related to any preceding commit.
	// (this only happens when we're doing `git log --all`). These will be tacked onto the far end.
	pos := maxPos + 1
	for _, pipe := range currentPipes {
		if equalHashes(pipe.toHash, commit.HashPtr()) {
			// turns out this commit does have a descendant so we'll place it right under the first instance
			pos = pipe.toPos
			break
		}
	}

	// a taken spot is one where a current pipe is ending on
	// Note: this set and similar ones below use int instead of int16 because
	// that's much more efficient. We cast the int16 values we store in these
	// sets to int on every access.
	takenSpots := set.New[int]()
	// a traversed spot is one where a current pipe is starting on, ending on, or passing through
	traversedSpots := set.New[int]()

	var toHash *string
	if commit.IsFirstCommit() {
		toHash = &EmptyTreeCommitHash
	} else {
		toHash = commit.ParentPtrs()[0]
	}
	newPipes = append(newPipes, Pipe{
		fromPos:  pos,
		toPos:    pos,
		fromHash: commit.HashPtr(),
		toHash:   toHash,
		kind:     STARTS,
		style:    getStyle(commit),
	})

	traversedSpotsForContinuingPipes := set.New[int]()
	for _, pipe := range currentPipes {
		if !equalHashes(pipe.toHash, commit.HashPtr()) {
			traversedSpotsForContinuingPipes.Add(int(pipe.toPos))
		}
	}

	getNextAvailablePosForContinuingPipe := func() int16 {
		i := int16(0)
		for {
			if !traversedSpots.Includes(int(i)) {
				return i
			}
			i++
		}
	}

	getNextAvailablePosForNewPipe := func() int16 {
		i := int16(0)
		for {
			// a newly created pipe is not allowed to end on a spot that's already taken,
			// nor on a spot that's been traversed by a continuing pipe.
			if !takenSpots.Includes(int(i)) && !traversedSpotsForContinuingPipes.Includes(int(i)) {
				return i
			}
			i++
		}
	}

	traverse := func(from, to int16) {
		left, right := from, to
		if left > right {
			left, right = right, left
		}
		for i := left; i <= right; i++ {
			traversedSpots.Add(int(i))
		}
		takenSpots.Add(int(to))
	}

	for _, pipe := range currentPipes {
		if equalHashes(pipe.toHash, commit.HashPtr()) {
			// terminating here
			newPipes = append(newPipes, Pipe{
				fromPos:  pipe.toPos,
				toPos:    pos,
				fromHash: pipe.fromHash,
				toHash:   pipe.toHash,
				kind:     TERMINATES,
				style:    pipe.style,
			})
			traverse(pipe.toPos, pos)
		} else if pipe.toPos < pos {
			// continuing here
			availablePos := getNextAvailablePosForContinuingPipe()
			newPipes = append(newPipes, Pipe{
				fromPos:  pipe.toPos,
				toPos:    availablePos,
				fromHash: pipe.fromHash,
				toHash:   pipe.toHash,
				kind:     CONTINUES,
				style:    pipe.style,
			})
			traverse(pipe.toPos, availablePos)
		}
	}

	if commit.IsMerge() {
		for _, parent := range commit.ParentPtrs()[1:] {
			availablePos := getNextAvailablePosForNewPipe()
			// need to act as if continuing pipes are going to continue on the same line.
			newPipes = append(newPipes, Pipe{
				fromPos:  pos,
				toPos:    availablePos,
				fromHash: commit.HashPtr(),
				toHash:   parent,
				kind:     STARTS,
				style:    getStyle(commit),
			})

			takenSpots.Add(int(availablePos))
		}
	}

	for _, pipe := range currentPipes {
		if !equalHashes(pipe.toHash, commit.HashPtr()) && pipe.toPos > pos {
			// continuing on, potentially moving left to fill in a blank spot
			last := pipe.toPos
			for i := pipe.toPos; i > pos; i-- {
				if takenSpots.Includes(int(i)) || traversedSpots.Includes(int(i)) {
					break
				}
				last = i
			}
			newPipes = append(newPipes, Pipe{
				fromPos:  pipe.toPos,
				toPos:    last,
				fromHash: pipe.fromHash,
				toHash:   pipe.toHash,
				kind:     CONTINUES,
				style:    pipe.style,
			})
			traverse(pipe.toPos, last)
		}
	}

	// not efficient but doing it for now: sorting my pipes by toPos, then by kind
	slices.SortFunc(newPipes, func(a, b Pipe) int {
		if a.toPos == b.toPos {
			return cmp.Compare(a.kind, b.kind)
		}
		return cmp.Compare(a.toPos, b.toPos)
	})

	return newPipes
}

func renderPipeSet(
	pipes []Pipe,
	selectedCommitHashPtr *string,
	prevCommit *models.Commit,
) string {
	maxPos := int16(0)
	commitPos := int16(0)
	startCount := 0
	for _, pipe := range pipes {
		if pipe.kind == STARTS {
			startCount++
			commitPos = pipe.fromPos
		} else if pipe.kind == TERMINATES {
			commitPos = pipe.toPos
		}

		if pipe.right() > maxPos {
			maxPos = pipe.right()
		}
	}
	isMerge := startCount > 1

	cells := lo.Map(lo.Range(int(maxPos)+1), func(i int, _ int) *Cell {
		return &Cell{cellType: CONNECTION, style: &style.FgDefault}
	})

	renderPipe := func(pipe *Pipe, style *style.TextStyle, overrideRightStyle bool) {
		left := pipe.left()
		right := pipe.right()

		if left != right {
			for i := left + 1; i < right; i++ {
				cells[i].setLeft(style).setRight(style, overrideRightStyle)
			}
			cells[left].setRight(style, overrideRightStyle)
			cells[right].setLeft(style)
		}

		if pipe.kind == STARTS || pipe.kind == CONTINUES {
			cells[pipe.toPos].setDown(style)
		}
		if pipe.kind == TERMINATES || pipe.kind == CONTINUES {
			cells[pipe.fromPos].setUp(style)
		}
	}

	// we don't want to highlight two commits if they're contiguous. We only want
	// to highlight multiple things if there's an actual visible pipe involved.
	highlight := true
	if prevCommit != nil && equalHashes(prevCommit.HashPtr(), selectedCommitHashPtr) {
		highlight = false
		for _, pipe := range pipes {
			if equalHashes(pipe.fromHash, selectedCommitHashPtr) && (pipe.kind != TERMINATES || pipe.fromPos != pipe.toPos) {
				highlight = true
			}
		}
	}

	// so we have our commit pos again, now it's time to build the cells.
	// we'll handle the one that's sourced from our selected commit last so that it can override the other cells.
	selectedPipes, nonSelectedPipes := utils.Partition(pipes, func(pipe Pipe) bool {
		return highlight && equalHashes(pipe.fromHash, selectedCommitHashPtr)
	})

	for _, pipe := range nonSelectedPipes {
		if pipe.kind == STARTS {
			renderPipe(&pipe, pipe.style, true)
		}
	}

	for _, pipe := range nonSelectedPipes {
		if pipe.kind != STARTS && !(pipe.kind == TERMINATES && pipe.fromPos == commitPos && pipe.toPos == commitPos) {
			renderPipe(&pipe, pipe.style, false)
		}
	}

	for _, pipe := range selectedPipes {
		for i := pipe.left(); i <= pipe.right(); i++ {
			cells[i].reset()
		}
	}
	for _, pipe := range selectedPipes {
		renderPipe(&pipe, &highlightStyle, true)
		if pipe.toPos == commitPos {
			cells[pipe.toPos].setStyle(&highlightStyle)
		}
	}

	cType := COMMIT
	if isMerge {
		cType = MERGE
	}

	cells[commitPos].setType(cType)

	// using a string builder here for the sake of performance
	writer := &strings.Builder{}
	writer.Grow(len(cells) * 2)
	for _, cell := range cells {
		cell.render(writer)
	}
	return writer.String()
}

func equalHashes(a, b *string) bool {
	// if our selectedCommitHashPtr is nil, there is no selected commit
	if a == nil || b == nil {
		return false
	}

	// We know that all hashes are stored in the pool, so we can compare their addresses
	return a == b
}
