package graphemes

import (
	"bufio"

	"github.com/clipperhouse/uax29/v2/internal/iterators"
)

// is determines if lookup intersects propert(ies)
func (lookup property) is(properties property) bool {
	return (lookup & properties) != 0
}

const _Ignore = _Extend

// SplitFunc is a bufio.SplitFunc implementation of Unicode grapheme cluster segmentation, for use with bufio.Scanner.
//
// See https://unicode.org/reports/tr29/#Grapheme_Cluster_Boundaries.
var SplitFunc bufio.SplitFunc = splitFunc[[]byte]

func splitFunc[T iterators.Stringish](data T, atEOF bool) (advance int, token T, err error) {
	var empty T
	if len(data) == 0 {
		return 0, empty, nil
	}

	// These vars are stateful across loop iterations
	var pos int
	var lastExIgnore property = 0     // "last excluding ignored categories"
	var lastLastExIgnore property = 0 // "last one before that"
	var regionalIndicatorCount int

	// Rules are usually of the form Cat1 × Cat2; "current" refers to the first property
	// to the right of the ×, from which we look back or forward

	current, w := lookup(data[pos:])
	if w == 0 {
		if !atEOF {
			// Rune extends past current data, request more
			return 0, empty, nil
		}
		pos = len(data)
		return pos, data[:pos], nil
	}

	// https://unicode.org/reports/tr29/#GB1
	// Start of text always advances
	pos += w

	for {
		eot := pos == len(data) // "end of text"

		if eot {
			if !atEOF {
				// Token extends past current data, request more
				return 0, empty, nil
			}

			// https://unicode.org/reports/tr29/#GB2
			break
		}

		/*
			We've switched the evaluation order of GB1↓ and GB2↑. It's ok:
			because we've checked for len(data) at the top of this function,
			sot and eot are mutually exclusive, order doesn't matter.
		*/

		// Rules are usually of the form Cat1 × Cat2; "current" refers to the first property
		// to the right of the ×, from which we look back or forward

		// Remember previous properties to avoid lookups/lookbacks
		last := current
		if !last.is(_Ignore) {
			lastLastExIgnore = lastExIgnore
			lastExIgnore = last
		}

		current, w = lookup(data[pos:])
		if w == 0 {
			if atEOF {
				// Just return the bytes, we can't do anything with them
				pos = len(data)
				break
			}
			// Rune extends past current data, request more
			return 0, empty, nil
		}

		// Optimization: no rule can possibly apply
		if current|last == 0 { // i.e. both are zero
			break
		}

		// https://unicode.org/reports/tr29/#GB3
		if current.is(_LF) && last.is(_CR) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB4
		// https://unicode.org/reports/tr29/#GB5
		if (current | last).is(_Control | _CR | _LF) {
			break
		}

		// https://unicode.org/reports/tr29/#GB6
		if current.is(_L|_V|_LV|_LVT) && last.is(_L) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB7
		if current.is(_V|_T) && last.is(_LV|_V) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB8
		if current.is(_T) && last.is(_LVT|_T) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB9
		if current.is(_Extend | _ZWJ) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB9a
		if current.is(_SpacingMark) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB9b
		if last.is(_Prepend) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB9c
		// TODO(clipperhouse):
		// It appears to be added in Unicode 15.1.0:
		// https://unicode.org/versions/Unicode15.1.0/#Migration
		// This package currently supports Unicode 15.0.0, so
		// out of scope for now

		// https://unicode.org/reports/tr29/#GB11
		if current.is(_ExtendedPictographic) && last.is(_ZWJ) && lastLastExIgnore.is(_ExtendedPictographic) {
			pos += w
			continue
		}

		// https://unicode.org/reports/tr29/#GB12
		// https://unicode.org/reports/tr29/#GB13
		if (current & last).is(_RegionalIndicator) {
			regionalIndicatorCount++

			odd := regionalIndicatorCount%2 == 1
			if odd {
				pos += w
				continue
			}
		}

		// If we fall through all the above rules, it's a grapheme cluster break
		break
	}

	// Return token
	return pos, data[:pos], nil
}
