package regexp2

import (
	"bytes"
	"errors"

	"github.com/dlclark/regexp2/syntax"
)

const (
	replaceSpecials     = 4
	replaceLeftPortion  = -1
	replaceRightPortion = -2
	replaceLastGroup    = -3
	replaceWholeString  = -4
)

// MatchEvaluator is a function that takes a match and returns a replacement string to be used
type MatchEvaluator func(Match) string

// Three very similar algorithms appear below: replace (pattern),
// replace (evaluator), and split.

// Replace Replaces all occurrences of the regex in the string with the
// replacement pattern.
//
// Note that the special case of no matches is handled on its own:
// with no matches, the input string is returned unchanged.
// The right-to-left case is split out because StringBuilder
// doesn't handle right-to-left string building directly very well.
func replace(regex *Regexp, data *syntax.ReplacerData, evaluator MatchEvaluator, input string, startAt, count int) (string, error) {
	if count < -1 {
		return "", errors.New("Count too small")
	}
	if count == 0 {
		return "", nil
	}

	m, err := regex.FindStringMatchStartingAt(input, startAt)

	if err != nil {
		return "", err
	}
	if m == nil {
		return input, nil
	}

	buf := &bytes.Buffer{}
	text := m.text

	if !regex.RightToLeft() {
		prevat := 0
		for m != nil {
			if m.Index != prevat {
				buf.WriteString(string(text[prevat:m.Index]))
			}
			prevat = m.Index + m.Length
			if evaluator == nil {
				replacementImpl(data, buf, m)
			} else {
				buf.WriteString(evaluator(*m))
			}

			count--
			if count == 0 {
				break
			}
			m, err = regex.FindNextMatch(m)
			if err != nil {
				return "", nil
			}
		}

		if prevat < len(text) {
			buf.WriteString(string(text[prevat:]))
		}
	} else {
		prevat := len(text)
		var al []string

		for m != nil {
			if m.Index+m.Length != prevat {
				al = append(al, string(text[m.Index+m.Length:prevat]))
			}
			prevat = m.Index
			if evaluator == nil {
				replacementImplRTL(data, &al, m)
			} else {
				al = append(al, evaluator(*m))
			}

			count--
			if count == 0 {
				break
			}
			m, err = regex.FindNextMatch(m)
			if err != nil {
				return "", nil
			}
		}

		if prevat > 0 {
			buf.WriteString(string(text[:prevat]))
		}

		for i := len(al) - 1; i >= 0; i-- {
			buf.WriteString(al[i])
		}
	}

	return buf.String(), nil
}

// Given a Match, emits into the StringBuilder the evaluated
// substitution pattern.
func replacementImpl(data *syntax.ReplacerData, buf *bytes.Buffer, m *Match) {
	for _, r := range data.Rules {

		if r >= 0 { // string lookup
			buf.WriteString(data.Strings[r])
		} else if r < -replaceSpecials { // group lookup
			m.groupValueAppendToBuf(-replaceSpecials-1-r, buf)
		} else {
			switch -replaceSpecials - 1 - r { // special insertion patterns
			case replaceLeftPortion:
				for i := 0; i < m.Index; i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceRightPortion:
				for i := m.Index + m.Length; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceLastGroup:
				m.groupValueAppendToBuf(m.GroupCount()-1, buf)
			case replaceWholeString:
				for i := 0; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			}
		}
	}
}

func replacementImplRTL(data *syntax.ReplacerData, al *[]string, m *Match) {
	l := *al
	buf := &bytes.Buffer{}

	for _, r := range data.Rules {
		buf.Reset()
		if r >= 0 { // string lookup
			l = append(l, data.Strings[r])
		} else if r < -replaceSpecials { // group lookup
			m.groupValueAppendToBuf(-replaceSpecials-1-r, buf)
			l = append(l, buf.String())
		} else {
			switch -replaceSpecials - 1 - r { // special insertion patterns
			case replaceLeftPortion:
				for i := 0; i < m.Index; i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceRightPortion:
				for i := m.Index + m.Length; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			case replaceLastGroup:
				m.groupValueAppendToBuf(m.GroupCount()-1, buf)
			case replaceWholeString:
				for i := 0; i < len(m.text); i++ {
					buf.WriteRune(m.text[i])
				}
			}
			l = append(l, buf.String())
		}
	}

	*al = l
}
