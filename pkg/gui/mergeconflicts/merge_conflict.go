package mergeconflicts

// mergeConflict : A git conflict with a start, ancestor (if exists), target, and end corresponding to line
// numbers in the file where the conflict markers appear.
// If no ancestor is present (i.e. we're not using the diff3 algorithm), then
// the `ancestor` field's value will be -1
type mergeConflict struct {
	start    int
	ancestor int
	target   int
	end      int
}

func (c *mergeConflict) hasAncestor() bool {
	return c.ancestor >= 0
}

func (c *mergeConflict) isMarkerLine(i int) bool {
	return i == c.start ||
		i == c.ancestor ||
		i == c.target ||
		i == c.end
}

type Selection int

const (
	TOP Selection = iota
	MIDDLE
	BOTTOM
	ALL
)

func (s Selection) isIndexToKeep(conflict *mergeConflict, i int) bool {
	// we're only handling one conflict at a time so any lines outside this
	// conflict we'll keep
	if i < conflict.start || conflict.end < i {
		return true
	}

	if conflict.isMarkerLine(i) {
		return false
	}

	return s.selected(conflict, i)
}

func (s Selection) bounds(c *mergeConflict) (int, int) {
	switch s {
	case TOP:
		if c.hasAncestor() {
			return c.start, c.ancestor
		} else {
			return c.start, c.target
		}
	case MIDDLE:
		return c.ancestor, c.target
	case BOTTOM:
		return c.target, c.end
	case ALL:
		return c.start, c.end
	}

	panic("unexpected selection for merge conflict")
}

func (s Selection) selected(c *mergeConflict, idx int) bool {
	start, end := s.bounds(c)
	return start < idx && idx < end
}

func availableSelections(c *mergeConflict) []Selection {
	if c.hasAncestor() {
		return []Selection{TOP, MIDDLE, BOTTOM}
	} else {
		return []Selection{TOP, BOTTOM}
	}
}
