/*
Package fuzzy provides fuzzy string matching optimized
for filenames and code symbols in the style of Sublime Text,
VSCode, IntelliJ IDEA et al.
*/
package fuzzy

import (
	"sort"
	"unicode"
	"unicode/utf8"
)

// Match represents a matched string.
type Match struct {
	// The matched string.
	Str string
	// The index of the matched string in the supplied slice.
	Index int
	// The indexes of matched characters. Useful for highlighting matches.
	MatchedIndexes []int
	// Score used to rank matches
	Score int
}

const (
	firstCharMatchBonus            = 10
	matchFollowingSeparatorBonus   = 20
	camelCaseMatchBonus            = 20
	adjacentMatchBonus             = 5
	unmatchedLeadingCharPenalty    = -5
	maxUnmatchedLeadingCharPenalty = -15
)

var separators = []rune("/-_ .\\")

// Matches is a slice of Match structs
type Matches []Match

func (a Matches) Len() int           { return len(a) }
func (a Matches) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Matches) Less(i, j int) bool { return a[i].Score >= a[j].Score }

// Source represents an abstract source of a list of strings. Source must be iterable type such as a slice.
// The source will be iterated over till Len() with String(i) being called for each element where i is the
// index of the element. You can find a working example in the README.
type Source interface {
	// The string to be matched at position i.
	String(i int) string
	// The length of the source. Typically is the length of the slice of things that you want to match.
	Len() int
}

type stringSource []string

func (ss stringSource) String(i int) string {
	return ss[i]
}

func (ss stringSource) Len() int { return len(ss) }

/*
Find looks up pattern in data and returns matches
in descending order of match quality. Match quality
is determined by a set of bonus and penalty rules.

The following types of matches apply a bonus:

* The first character in the pattern matches the first character in the match string.

* The matched character is camel cased.

* The matched character follows a separator such as an underscore character.

* The matched character is adjacent to a previous match.

Penalties are applied for every character in the search string that wasn't matched and all leading
characters upto the first match.
*/
func Find(pattern string, data []string) Matches {
	return FindFrom(pattern, stringSource(data))
}

/*
FindFrom is an alternative implementation of Find using a Source
instead of a list of strings.
*/
func FindFrom(pattern string, data Source) Matches {
	if len(pattern) == 0 {
		return nil
	}
	runes := []rune(pattern)
	var matches Matches
	var matchedIndexes []int
	for i := 0; i < data.Len(); i++ {
		var match Match
		match.Str = data.String(i)
		match.Index = i
		if matchedIndexes != nil {
			match.MatchedIndexes = matchedIndexes
		} else {
			match.MatchedIndexes = make([]int, 0, len(runes))
		}
		var score int
		patternIndex := 0
		bestScore := -1
		matchedIndex := -1
		currAdjacentMatchBonus := 0
		var last rune
		var lastIndex int
		nextc, nextSize := utf8.DecodeRuneInString(data.String(i))
		var candidate rune
		var candidateSize int
		for j := 0; j < len(data.String(i)); j += candidateSize {
			candidate, candidateSize = nextc, nextSize
			if equalFold(candidate, runes[patternIndex]) {
				score = 0
				if j == 0 {
					score += firstCharMatchBonus
				}
				if unicode.IsLower(last) && unicode.IsUpper(candidate) {
					score += camelCaseMatchBonus
				}
				if j != 0 && isSeparator(last) {
					score += matchFollowingSeparatorBonus
				}
				if len(match.MatchedIndexes) > 0 {
					lastMatch := match.MatchedIndexes[len(match.MatchedIndexes)-1]
					bonus := adjacentCharBonus(lastIndex, lastMatch, currAdjacentMatchBonus)
					score += bonus
					// adjacent matches are incremental and keep increasing based on previous adjacent matches
					// thus we need to maintain the current match bonus
					currAdjacentMatchBonus += bonus
				}
				if score > bestScore {
					bestScore = score
					matchedIndex = j
				}
			}
			var nextp rune
			if patternIndex < len(runes)-1 {
				nextp = runes[patternIndex+1]
			}
			if j+candidateSize < len(data.String(i)) {
				if data.String(i)[j+candidateSize] < utf8.RuneSelf { // Fast path for ASCII
					nextc, nextSize = rune(data.String(i)[j+candidateSize]), 1
				} else {
					nextc, nextSize = utf8.DecodeRuneInString(data.String(i)[j+candidateSize:])
				}
			} else {
				nextc, nextSize = 0, 0
			}
			// We apply the best score when we have the next match coming up or when the search string has ended.
			// Tracking when the next match is coming up allows us to exhaustively find the best match and not necessarily
			// the first match.
			// For example given the pattern "tk" and search string "The Black Knight", exhaustively matching allows us
			// to match the second k thus giving this string a higher score.
			if equalFold(nextp, nextc) || nextc == 0 {
				if matchedIndex > -1 {
					if len(match.MatchedIndexes) == 0 {
						penalty := matchedIndex * unmatchedLeadingCharPenalty
						bestScore += max(penalty, maxUnmatchedLeadingCharPenalty)
					}
					match.Score += bestScore
					match.MatchedIndexes = append(match.MatchedIndexes, matchedIndex)
					score = 0
					bestScore = -1
					patternIndex++
				}
			}
			lastIndex = j
			last = candidate
		}
		// apply penalty for each unmatched character
		penalty := len(match.MatchedIndexes) - len(data.String(i))
		match.Score += penalty
		if len(match.MatchedIndexes) == len(runes) {
			matches = append(matches, match)
			matchedIndexes = nil
		} else {
			matchedIndexes = match.MatchedIndexes[:0] // Recycle match index slice
		}
	}
	sort.Stable(matches)
	return matches
}

// Taken from strings.EqualFold
func equalFold(tr, sr rune) bool {
	if tr == sr {
		return true
	}
	if tr < sr {
		tr, sr = sr, tr
	}
	// Fast check for ASCII.
	if tr < utf8.RuneSelf {
		// ASCII, and sr is upper case.  tr must be lower case.
		if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
			return true
		}
		return false
	}

	// General case. SimpleFold(x) returns the next equivalent rune > x
	// or wraps around to smaller values.
	r := unicode.SimpleFold(sr)
	for r != sr && r < tr {
		r = unicode.SimpleFold(r)
	}
	return r == tr
}

func adjacentCharBonus(i int, lastMatch int, currentBonus int) int {
	if lastMatch == i {
		return currentBonus*2 + adjacentMatchBonus
	}
	return 0
}

func isSeparator(s rune) bool {
	for _, sep := range separators {
		if s == sep {
			return true
		}
	}
	return false
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
