package patch_exploring

import "github.com/jesseduffield/lazygit/pkg/utils"

func calculateOrigin(currentOrigin int, bufferHeight int, firstLineIdx int, lastLineIdx int, selectedLineIdx int, mode selectMode) int {
	needToSeeIdx, wantToSeeIdx := getNeedAndWantLineIdx(firstLineIdx, lastLineIdx, selectedLineIdx, mode)

	return calculateNewOriginWithNeededAndWantedIdx(currentOrigin, bufferHeight, needToSeeIdx, wantToSeeIdx)
}

// we want to scroll our origin so that the index we need to see is in view
// and the other index we want to see (e.g. the other side of a line range)
// is in as close to being in view as possible.
func calculateNewOriginWithNeededAndWantedIdx(currentOrigin int, bufferHeight int, needToSeeIdx int, wantToSeeIdx int) int {
	origin := currentOrigin
	if needToSeeIdx < currentOrigin {
		origin = needToSeeIdx
	} else if needToSeeIdx > currentOrigin+bufferHeight {
		origin = needToSeeIdx - bufferHeight
	}

	bottom := origin + bufferHeight

	if wantToSeeIdx < origin {
		requiredChange := origin - wantToSeeIdx
		allowedChange := bottom - needToSeeIdx
		return origin - utils.Min(requiredChange, allowedChange)
	} else if wantToSeeIdx > origin+bufferHeight {
		requiredChange := wantToSeeIdx - bottom
		allowedChange := needToSeeIdx - origin
		return origin + utils.Min(requiredChange, allowedChange)
	} else {
		return origin
	}
}

func getNeedAndWantLineIdx(firstLineIdx int, lastLineIdx int, selectedLineIdx int, mode selectMode) (int, int) {
	switch mode {
	case LINE:
		return selectedLineIdx, selectedLineIdx
	case RANGE:
		if selectedLineIdx == firstLineIdx {
			return firstLineIdx, lastLineIdx
		} else {
			return lastLineIdx, firstLineIdx
		}
	case HUNK:
		return firstLineIdx, lastLineIdx
	default:
		// we should never land here
		panic("unknown mode")
	}
}
