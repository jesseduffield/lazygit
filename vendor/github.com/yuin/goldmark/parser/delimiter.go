package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// A DelimiterProcessor interface provides a set of functions about
// Delimiter nodes.
type DelimiterProcessor interface {
	// IsDelimiter returns true if given character is a delimiter, otherwise false.
	IsDelimiter(byte) bool

	// CanOpenCloser returns true if given opener can close given closer, otherwise false.
	CanOpenCloser(opener, closer *Delimiter) bool

	// OnMatch will be called when new matched delimiter found.
	// OnMatch should return a new Node correspond to the matched delimiter.
	OnMatch(consumes int) ast.Node
}

// A Delimiter struct represents a delimiter like '*' of the Markdown text.
type Delimiter struct {
	ast.BaseInline

	Segment text.Segment

	// CanOpen is set true if this delimiter can open a span for a new node.
	// See https://spec.commonmark.org/0.30/#can-open-emphasis for details.
	CanOpen bool

	// CanClose is set true if this delimiter can close a span for a new node.
	// See https://spec.commonmark.org/0.30/#can-open-emphasis for details.
	CanClose bool

	// Length is a remaining length of this delimiter.
	Length int

	// OriginalLength is a original length of this delimiter.
	OriginalLength int

	// Char is a character of this delimiter.
	Char byte

	// PreviousDelimiter is a previous sibling delimiter node of this delimiter.
	PreviousDelimiter *Delimiter

	// NextDelimiter is a next sibling delimiter node of this delimiter.
	NextDelimiter *Delimiter

	// Processor is a DelimiterProcessor associated with this delimiter.
	Processor DelimiterProcessor
}

// Inline implements Inline.Inline.
func (d *Delimiter) Inline() {}

// Dump implements Node.Dump.
func (d *Delimiter) Dump(source []byte, level int) {
	fmt.Printf("%sDelimiter: \"%s\"\n", strings.Repeat("    ", level), string(d.Text(source)))
}

var kindDelimiter = ast.NewNodeKind("Delimiter")

// Kind implements Node.Kind
func (d *Delimiter) Kind() ast.NodeKind {
	return kindDelimiter
}

// Text implements Node.Text
func (d *Delimiter) Text(source []byte) []byte {
	return d.Segment.Value(source)
}

// ConsumeCharacters consumes delimiters.
func (d *Delimiter) ConsumeCharacters(n int) {
	d.Length -= n
	d.Segment = d.Segment.WithStop(d.Segment.Start + d.Length)
}

// CalcComsumption calculates how many characters should be used for opening
// a new span correspond to given closer.
func (d *Delimiter) CalcComsumption(closer *Delimiter) int {
	if (d.CanClose || closer.CanOpen) && (d.OriginalLength+closer.OriginalLength)%3 == 0 && closer.OriginalLength%3 != 0 {
		return 0
	}
	if d.Length >= 2 && closer.Length >= 2 {
		return 2
	}
	return 1
}

// NewDelimiter returns a new Delimiter node.
func NewDelimiter(canOpen, canClose bool, length int, char byte, processor DelimiterProcessor) *Delimiter {
	c := &Delimiter{
		BaseInline:        ast.BaseInline{},
		CanOpen:           canOpen,
		CanClose:          canClose,
		Length:            length,
		OriginalLength:    length,
		Char:              char,
		PreviousDelimiter: nil,
		NextDelimiter:     nil,
		Processor:         processor,
	}
	return c
}

// ScanDelimiter scans a delimiter by given DelimiterProcessor.
func ScanDelimiter(line []byte, before rune, min int, processor DelimiterProcessor) *Delimiter {
	i := 0
	c := line[i]
	j := i
	if !processor.IsDelimiter(c) {
		return nil
	}
	for ; j < len(line) && c == line[j]; j++ {
	}
	if (j - i) >= min {
		after := rune(' ')
		if j != len(line) {
			after = util.ToRune(line, j)
		}

		canOpen, canClose := false, false
		beforeIsPunctuation := util.IsPunctRune(before)
		beforeIsWhitespace := util.IsSpaceRune(before)
		afterIsPunctuation := util.IsPunctRune(after)
		afterIsWhitespace := util.IsSpaceRune(after)

		isLeft := !afterIsWhitespace &&
			(!afterIsPunctuation || beforeIsWhitespace || beforeIsPunctuation)
		isRight := !beforeIsWhitespace &&
			(!beforeIsPunctuation || afterIsWhitespace || afterIsPunctuation)

		if line[i] == '_' {
			canOpen = isLeft && (!isRight || beforeIsPunctuation)
			canClose = isRight && (!isLeft || afterIsPunctuation)
		} else {
			canOpen = isLeft
			canClose = isRight
		}
		return NewDelimiter(canOpen, canClose, j-i, c, processor)
	}
	return nil
}

// ProcessDelimiters processes the delimiter list in the context.
// Processing will be stop when reaching the bottom.
//
// If you implement an inline parser that can have other inline nodes as
// children, you should call this function when nesting span has closed.
func ProcessDelimiters(bottom ast.Node, pc Context) {
	lastDelimiter := pc.LastDelimiter()
	if lastDelimiter == nil {
		return
	}
	var closer *Delimiter
	if bottom != nil {
		if bottom != lastDelimiter {
			for c := lastDelimiter.PreviousSibling(); c != nil && c != bottom; {
				if d, ok := c.(*Delimiter); ok {
					closer = d
				}
				c = c.PreviousSibling()
			}
		}
	} else {
		closer = pc.FirstDelimiter()
	}
	if closer == nil {
		pc.ClearDelimiters(bottom)
		return
	}
	for closer != nil {
		if !closer.CanClose {
			closer = closer.NextDelimiter
			continue
		}
		consume := 0
		found := false
		maybeOpener := false
		var opener *Delimiter
		for opener = closer.PreviousDelimiter; opener != nil && opener != bottom; opener = opener.PreviousDelimiter {
			if opener.CanOpen && opener.Processor.CanOpenCloser(opener, closer) {
				maybeOpener = true
				consume = opener.CalcComsumption(closer)
				if consume > 0 {
					found = true
					break
				}
			}
		}
		if !found {
			next := closer.NextDelimiter
			if !maybeOpener && !closer.CanOpen {
				pc.RemoveDelimiter(closer)
			}
			closer = next
			continue
		}
		opener.ConsumeCharacters(consume)
		closer.ConsumeCharacters(consume)

		node := opener.Processor.OnMatch(consume)

		parent := opener.Parent()
		child := opener.NextSibling()

		for child != nil && child != closer {
			next := child.NextSibling()
			node.AppendChild(node, child)
			child = next
		}
		parent.InsertAfter(parent, opener, node)

		for c := opener.NextDelimiter; c != nil && c != closer; {
			next := c.NextDelimiter
			pc.RemoveDelimiter(c)
			c = next
		}

		if opener.Length == 0 {
			pc.RemoveDelimiter(opener)
		}

		if closer.Length == 0 {
			next := closer.NextDelimiter
			pc.RemoveDelimiter(closer)
			closer = next
		}
	}
	pc.ClearDelimiters(bottom)
}
