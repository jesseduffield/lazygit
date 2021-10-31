package graph

import (
	"sort"
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type PipeKind uint8

const (
	TERMINATES PipeKind = iota
	STARTS
	CONTINUES
)

type Pipe struct {
	fromPos int
	toPos   int
	fromSha string
	toSha   string
	kind    PipeKind
	style   style.TextStyle
}

var highlightStyle = style.FgLightWhite.SetBold()

func ContainsCommitSha(pipes []Pipe, sha string) bool {
	for _, pipe := range pipes {
		if equalHashes(pipe.fromSha, sha) {
			return true
		}
	}
	return false
}

func (self Pipe) left() int {
	return utils.Min(self.fromPos, self.toPos)
}

func (self Pipe) right() int {
	return utils.Max(self.fromPos, self.toPos)
}

func RenderCommitGraph(commits []*models.Commit, selectedCommitSha string, getStyle func(c *models.Commit) style.TextStyle) []string {
	pipeSets := GetPipeSets(commits, getStyle)
	if len(pipeSets) == 0 {
		return nil
	}

	lines := RenderAux(pipeSets, commits, selectedCommitSha)

	return lines
}

func GetPipeSets(commits []*models.Commit, getStyle func(c *models.Commit) style.TextStyle) [][]Pipe {
	if len(commits) == 0 {
		return nil
	}

	pipes := []Pipe{{fromPos: 0, toPos: 0, fromSha: "START", toSha: commits[0].Sha, kind: STARTS, style: style.FgDefault}}

	pipeSets := [][]Pipe{}
	for _, commit := range commits {
		pipes = getNextPipes(pipes, commit, getStyle)
		pipeSets = append(pipeSets, pipes)
	}

	return pipeSets
}

func RenderAux(pipeSets [][]Pipe, commits []*models.Commit, selectedCommitSha string) []string {
	lines := make([]string, len(pipeSets))
	wg := sync.WaitGroup{}
	wg.Add(len(pipeSets))
	for i, pipeSet := range pipeSets {
		i := i
		pipeSet := pipeSet
		go func() {
			defer wg.Done()
			var prevCommit *models.Commit
			if i > 0 {
				prevCommit = commits[i-1]
			}
			line := renderPipeSet(pipeSet, selectedCommitSha, prevCommit)
			lines[i] = line
		}()
	}
	wg.Wait()
	return lines
}

func getNextPipes(prevPipes []Pipe, commit *models.Commit, getStyle func(c *models.Commit) style.TextStyle) []Pipe {
	currentPipes := make([]Pipe, 0, len(prevPipes))
	maxPos := 0
	for _, pipe := range prevPipes {
		// a pipe that terminated in the previous line has no bearing on the current line
		// so we'll filter those out
		if pipe.kind != TERMINATES {
			currentPipes = append(currentPipes, pipe)
		}
		maxPos = utils.Max(maxPos, pipe.toPos)
	}

	newPipes := make([]Pipe, 0, len(currentPipes)+len(commit.Parents))
	// start by assuming that we've got a brand new commit not related to any preceding commit.
	// (this only happens when we're doing `git log --all`). These will be tacked onto the far end.
	pos := maxPos + 1
	for _, pipe := range currentPipes {
		if equalHashes(pipe.toSha, commit.Sha) {
			// turns out this commit does have a descendant so we'll place it right under the first instance
			pos = pipe.toPos
			break
		}
	}

	// a taken spot is one where a current pipe is ending on
	takenSpots := make(map[int]bool)
	// a traversed spot is one where a current pipe is starting on, ending on, or passing through
	traversedSpots := make(map[int]bool)

	if len(commit.Parents) > 0 {
		newPipes = append(newPipes, Pipe{
			fromPos: pos,
			toPos:   pos,
			fromSha: commit.Sha,
			toSha:   commit.Parents[0],
			kind:    STARTS,
			style:   getStyle(commit),
		})
	}

	traversedSpotsForContinuingPipes := make(map[int]bool)
	for _, pipe := range currentPipes {
		if !equalHashes(pipe.toSha, commit.Sha) {
			traversedSpotsForContinuingPipes[pipe.toPos] = true
		}
	}

	getNextAvailablePosForContinuingPipe := func() int {
		i := 0
		for {
			if !traversedSpots[i] {
				return i
			}
			i++
		}
	}

	getNextAvailablePosForNewPipe := func() int {
		i := 0
		for {
			// a newly created pipe is not allowed to end on a spot that's already taken,
			// nor on a spot that's been traversed by a continuing pipe.
			if !takenSpots[i] && !traversedSpotsForContinuingPipes[i] {
				return i
			}
			i++
		}
	}

	traverse := func(from, to int) {
		left, right := from, to
		if left > right {
			left, right = right, left
		}
		for i := left; i <= right; i++ {
			traversedSpots[i] = true
		}
		takenSpots[to] = true
	}

	for _, pipe := range currentPipes {
		if equalHashes(pipe.toSha, commit.Sha) {
			// terminating here
			newPipes = append(newPipes, Pipe{
				fromPos: pipe.toPos,
				toPos:   pos,
				fromSha: pipe.fromSha,
				toSha:   pipe.toSha,
				kind:    TERMINATES,
				style:   pipe.style,
			})
			traverse(pipe.toPos, pos)
		} else if pipe.toPos < pos {
			// continuing here
			availablePos := getNextAvailablePosForContinuingPipe()
			newPipes = append(newPipes, Pipe{
				fromPos: pipe.toPos,
				toPos:   availablePos,
				fromSha: pipe.fromSha,
				toSha:   pipe.toSha,
				kind:    CONTINUES,
				style:   pipe.style,
			})
			traverse(pipe.toPos, availablePos)
		}
	}

	if commit.IsMerge() {
		for _, parent := range commit.Parents[1:] {
			availablePos := getNextAvailablePosForNewPipe()
			// need to act as if continuing pipes are going to continue on the same line.
			newPipes = append(newPipes, Pipe{
				fromPos: pos,
				toPos:   availablePos,
				fromSha: commit.Sha,
				toSha:   parent,
				kind:    STARTS,
				style:   getStyle(commit),
			})

			takenSpots[availablePos] = true
		}
	}

	for _, pipe := range currentPipes {
		if !equalHashes(pipe.toSha, commit.Sha) && pipe.toPos > pos {
			// continuing on, potentially moving left to fill in a blank spot
			last := pipe.toPos
			for i := pipe.toPos; i > pos; i-- {
				if takenSpots[i] || traversedSpots[i] {
					break
				} else {
					last = i
				}
			}
			newPipes = append(newPipes, Pipe{
				fromPos: pipe.toPos,
				toPos:   last,
				fromSha: pipe.fromSha,
				toSha:   pipe.toSha,
				kind:    CONTINUES,
				style:   pipe.style,
			})
			traverse(pipe.toPos, last)
		}
	}

	// not efficient but doing it for now: sorting my pipes by toPos, then by kind
	sort.Slice(newPipes, func(i, j int) bool {
		if newPipes[i].toPos == newPipes[j].toPos {
			return newPipes[i].kind < newPipes[j].kind
		}
		return newPipes[i].toPos < newPipes[j].toPos
	})

	return newPipes
}

func renderPipeSet(
	pipes []Pipe,
	selectedCommitSha string,
	prevCommit *models.Commit,
) string {
	maxPos := 0
	commitPos := 0
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

	cells := make([]*Cell, maxPos+1)
	for i := range cells {
		cells[i] = &Cell{cellType: CONNECTION, style: style.FgDefault}
	}

	renderPipe := func(pipe Pipe, style style.TextStyle, overrideRightStyle bool) {
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
	if prevCommit != nil && equalHashes(prevCommit.Sha, selectedCommitSha) {
		highlight = false
		for _, pipe := range pipes {
			if equalHashes(pipe.fromSha, selectedCommitSha) && (pipe.kind != TERMINATES || pipe.fromPos != pipe.toPos) {
				highlight = true
			}
		}
	}

	// so we have our commit pos again, now it's time to build the cells.
	// we'll handle the one that's sourced from our selected commit last so that it can override the other cells.
	selectedPipes := []Pipe{}
	// pre-allocating this one because most of the time we'll only have non-selected pipes
	nonSelectedPipes := make([]Pipe, 0, len(pipes))

	for _, pipe := range pipes {
		if highlight && equalHashes(pipe.fromSha, selectedCommitSha) {
			selectedPipes = append(selectedPipes, pipe)
		} else {
			nonSelectedPipes = append(nonSelectedPipes, pipe)
		}
	}

	for _, pipe := range nonSelectedPipes {
		if pipe.kind == STARTS {
			renderPipe(pipe, pipe.style, true)
		}
	}

	for _, pipe := range nonSelectedPipes {
		if pipe.kind != STARTS && !(pipe.kind == TERMINATES && pipe.fromPos == commitPos && pipe.toPos == commitPos) {
			renderPipe(pipe, pipe.style, false)
		}
	}

	for _, pipe := range selectedPipes {
		for i := pipe.left(); i <= pipe.right(); i++ {
			cells[i].reset()
		}
	}
	for _, pipe := range selectedPipes {
		renderPipe(pipe, highlightStyle, true)
		if pipe.toPos == commitPos {
			cells[pipe.toPos].setStyle(highlightStyle)
		}
	}

	cType := COMMIT
	if isMerge {
		cType = MERGE
	}

	cells[commitPos].setType(cType)

	renderedCells := make([]string, len(cells))
	for i, cell := range cells {
		renderedCells[i] = cell.render()
	}
	return strings.Join(renderedCells, "")
}

func equalHashes(a, b string) bool {
	// if our selectedCommitSha is an empty string we treat that as meaning there is no selected commit sha
	if a == "" || b == "" {
		return false
	}

	length := utils.Min(len(a), len(b))
	// parent hashes are only stored up to 20 characters for some reason so we'll truncate to that for comparison
	return a[:length] == b[:length]
}
