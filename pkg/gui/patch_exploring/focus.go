package patch_exploring

func calculateOrigin(currentOrigin int, bufferHeight int, numLines int, firstLineIdx int, lastLineIdx int, selectedLineIdx int, mode selectMode) int {
	needToSeeIdx, wantToSeeIdx := getNeedAndWantLineIdx(firstLineIdx, lastLineIdx, selectedLineIdx, mode)

	return calculateNewOriginWithNeededAndWantedIdx(currentOrigin, bufferHeight, numLines, needToSeeIdx, wantToSeeIdx)
}

// we want to scroll our origin so that the index we need to see is in view
// and the other index we want to see (e.g. the other side of a line range)
// is as close to being in view as possible.
func calculateNewOriginWithNeededAndWantedIdx(currentOrigin int, bufferHeight int, numLines int, needToSeeIdx int, wantToSeeIdx int) int {
	origin := currentOrigin
	if needToSeeIdx < currentOrigin || needToSeeIdx > currentOrigin+bufferHeight {
		origin = max(min(needToSeeIdx-bufferHeight/2, numLines-bufferHeight-1), 0)
	}

	bottom := origin + bufferHeight

	if wantToSeeIdx < origin {
		requiredChange := origin - wantToSeeIdx
		allowedChange := bottom - needToSeeIdx
		return origin - min(requiredChange, allowedChange)
	} else if wantToSeeIdx > origin+bufferHeight {
		requiredChange := wantToSeeIdx - bottom
		allowedChange := needToSeeIdx - origin
		return origin + min(requiredChange, allowedChange)
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
