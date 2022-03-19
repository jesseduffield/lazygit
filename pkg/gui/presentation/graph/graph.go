package graph

import (
	"runtime"
	"strings"
	"sync"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
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
	fromPos int
	toPos   int
	fromSha string
	toSha   string
	kind    PipeKind
	style   style.TextStyle
}

var highlightStyle = style.FgLightWhite.SetBold()

func ContainsCommitSha(pipes []*Pipe, sha string) bool {
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

func GetPipeSets(commits []*models.Commit, getStyle func(c *models.Commit) style.TextStyle) [][]*Pipe {
	if len(commits) == 0 {
		return nil
	}

	pipes := []*Pipe{{fromPos: 0, toPos: 0, fromSha: "START", toSha: commits[0].Sha, kind: STARTS, style: style.FgDefault}}

	return slices.Map(commits, func(commit *models.Commit) []*Pipe {
		pipes = getNextPipes(pipes, commit, getStyle)
		return pipes
	})
}

func RenderAux(pipeSets [][]*Pipe, commits []*models.Commit, selectedCommitSha string) []string {
	maxProcs := runtime.GOMAXPROCS(0)

	// splitting up the rendering of the graph into multiple goroutines allows us to render the graph in parallel
	chunks := make([][]string, maxProcs)
	perProc := len(pipeSets) / maxProcs

	wg := sync.WaitGroup{}
	wg.Add(maxProcs)

	for i := 0; i < maxProcs; i++ {
		i := i
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
				line := renderPipeSet(pipeSet, selectedCommitSha, prevCommit)
				innerLines = append(innerLines, line)
			}
			chunks[i] = innerLines
			wg.Done()
		}()
	}

	wg.Wait()

	return slices.Flatten(chunks)
}

func getNextPipes(prevPipes []*Pipe, commit *models.Commit, getStyle func(c *models.Commit) style.TextStyle) []*Pipe {
	maxPos := slices.MaxBy(prevPipes, func(pipe *Pipe) int {
		return pipe.toPos
	})

	// a pipe that terminated in the previous line has no bearing on the current line
	// so we'll filter those out
	currentPipes := slices.Filter(prevPipes, func(pipe *Pipe) bool {
		return pipe.kind != TERMINATES
	})

	newPipes := make([]*Pipe, 0, len(currentPipes)+len(commit.Parents))
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
	takenSpots := set.New[int]()
	// a traversed spot is one where a current pipe is starting on, ending on, or passing through
	traversedSpots := set.New[int]()

	if len(commit.Parents) > 0 {
		newPipes = append(newPipes, &Pipe{
			fromPos: pos,
			toPos:   pos,
			fromSha: commit.Sha,
			toSha:   commit.Parents[0],
			kind:    STARTS,
			style:   getStyle(commit),
		})
	}

	traversedSpotsForContinuingPipes := set.New[int]()
	for _, pipe := range currentPipes {
		if !equalHashes(pipe.toSha, commit.Sha) {
			traversedSpotsForContinuingPipes.Add(pipe.toPos)
		}
	}

	getNextAvailablePosForContinuingPipe := func() int {
		i := 0
		for {
			if !traversedSpots.Includes(i) {
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
			if !takenSpots.Includes(i) && !traversedSpotsForContinuingPipes.Includes(i) {
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
			traversedSpots.Add(i)
		}
		takenSpots.Add(to)
	}

	for _, pipe := range currentPipes {
		if equalHashes(pipe.toSha, commit.Sha) {
			// terminating here
			newPipes = append(newPipes, &Pipe{
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
			newPipes = append(newPipes, &Pipe{
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
			newPipes = append(newPipes, &Pipe{
				fromPos: pos,
				toPos:   availablePos,
				fromSha: commit.Sha,
				toSha:   parent,
				kind:    STARTS,
				style:   getStyle(commit),
			})

			takenSpots.Add(availablePos)
		}
	}

	for _, pipe := range currentPipes {
		if !equalHashes(pipe.toSha, commit.Sha) && pipe.toPos > pos {
			// continuing on, potentially moving left to fill in a blank spot
			last := pipe.toPos
			for i := pipe.toPos; i > pos; i-- {
				if takenSpots.Includes(i) || traversedSpots.Includes(i) {
					break
				} else {
					last = i
				}
			}
			newPipes = append(newPipes, &Pipe{
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
	slices.SortFunc(newPipes, func(a, b *Pipe) bool {
		if a.toPos == b.toPos {
			return a.kind < b.kind
		}
		return a.toPos < b.toPos
	})

	return newPipes
}

func renderPipeSet(
	pipes []*Pipe,
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

	cells := slices.Map(lo.Range(maxPos+1), func(i int) *Cell {
		return &Cell{cellType: CONNECTION, style: style.FgDefault}
	})

	renderPipe := func(pipe *Pipe, style style.TextStyle, overrideRightStyle bool) {
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
	selectedPipes, nonSelectedPipes := slices.Partition(pipes, func(pipe *Pipe) bool {
		return highlight && equalHashes(pipe.fromSha, selectedCommitSha)
	})

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

	// using a string builder here for the sake of performance
	writer := &strings.Builder{}
	writer.Grow(len(cells) * 2)
	for _, cell := range cells {
		cell.render(writer)
	}
	return writer.String()
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
