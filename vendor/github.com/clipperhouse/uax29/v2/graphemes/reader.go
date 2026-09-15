// Package graphemes implements Unicode grapheme cluster boundaries: https://unicode.org/reports/tr29/#Grapheme_Cluster_Boundaries
package graphemes

import (
	"bufio"
	"io"
)

type Scanner struct {
	*bufio.Scanner
}

// FromReader returns a Scanner, to split graphemes per
// https://unicode.org/reports/tr29/#Grapheme_Cluster_Boundaries.
//
// It embeds a [bufio.Scanner], so you can use its methods.
//
// Iterate through graphemes by calling Scan() until false, then check Err().
func FromReader(r io.Reader) *Scanner {
	sc := bufio.NewScanner(r)
	sc.Split(SplitFunc)
	return &Scanner{
		Scanner: sc,
	}
}
