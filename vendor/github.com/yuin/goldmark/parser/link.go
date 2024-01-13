package parser

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var linkLabelStateKey = NewContextKey()

type linkLabelState struct {
	ast.BaseInline

	Segment text.Segment

	IsImage bool

	Prev *linkLabelState

	Next *linkLabelState

	First *linkLabelState

	Last *linkLabelState
}

func newLinkLabelState(segment text.Segment, isImage bool) *linkLabelState {
	return &linkLabelState{
		Segment: segment,
		IsImage: isImage,
	}
}

func (s *linkLabelState) Text(source []byte) []byte {
	return s.Segment.Value(source)
}

func (s *linkLabelState) Dump(source []byte, level int) {
	fmt.Printf("%slinkLabelState: \"%s\"\n", strings.Repeat("    ", level), s.Text(source))
}

var kindLinkLabelState = ast.NewNodeKind("LinkLabelState")

func (s *linkLabelState) Kind() ast.NodeKind {
	return kindLinkLabelState
}

func linkLabelStateLength(v *linkLabelState) int {
	if v == nil || v.Last == nil || v.First == nil {
		return 0
	}
	return v.Last.Segment.Stop - v.First.Segment.Start
}

func pushLinkLabelState(pc Context, v *linkLabelState) {
	tlist := pc.Get(linkLabelStateKey)
	var list *linkLabelState
	if tlist == nil {
		list = v
		v.First = v
		v.Last = v
		pc.Set(linkLabelStateKey, list)
	} else {
		list = tlist.(*linkLabelState)
		l := list.Last
		list.Last = v
		l.Next = v
		v.Prev = l
	}
}

func removeLinkLabelState(pc Context, d *linkLabelState) {
	tlist := pc.Get(linkLabelStateKey)
	var list *linkLabelState
	if tlist == nil {
		return
	}
	list = tlist.(*linkLabelState)

	if d.Prev == nil {
		list = d.Next
		if list != nil {
			list.First = d
			list.Last = d.Last
			list.Prev = nil
			pc.Set(linkLabelStateKey, list)
		} else {
			pc.Set(linkLabelStateKey, nil)
		}
	} else {
		d.Prev.Next = d.Next
		if d.Next != nil {
			d.Next.Prev = d.Prev
		}
	}
	if list != nil && d.Next == nil {
		list.Last = d.Prev
	}
	d.Next = nil
	d.Prev = nil
	d.First = nil
	d.Last = nil
}

type linkParser struct {
}

var defaultLinkParser = &linkParser{}

// NewLinkParser return a new InlineParser that parses links.
func NewLinkParser() InlineParser {
	return defaultLinkParser
}

func (s *linkParser) Trigger() []byte {
	return []byte{'!', '[', ']'}
}

var linkBottom = NewContextKey()

func (s *linkParser) Parse(parent ast.Node, block text.Reader, pc Context) ast.Node {
	line, segment := block.PeekLine()
	if line[0] == '!' {
		if len(line) > 1 && line[1] == '[' {
			block.Advance(1)
			pc.Set(linkBottom, pc.LastDelimiter())
			return processLinkLabelOpen(block, segment.Start+1, true, pc)
		}
		return nil
	}
	if line[0] == '[' {
		pc.Set(linkBottom, pc.LastDelimiter())
		return processLinkLabelOpen(block, segment.Start, false, pc)
	}

	// line[0] == ']'
	tlist := pc.Get(linkLabelStateKey)
	if tlist == nil {
		return nil
	}
	last := tlist.(*linkLabelState).Last
	if last == nil {
		return nil
	}
	block.Advance(1)
	removeLinkLabelState(pc, last)
	// CommonMark spec says:
	//  > A link label can have at most 999 characters inside the square brackets.
	if linkLabelStateLength(tlist.(*linkLabelState)) > 998 {
		ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
		return nil
	}

	if !last.IsImage && s.containsLink(last) { // a link in a link text is not allowed
		ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
		return nil
	}

	c := block.Peek()
	l, pos := block.Position()
	var link *ast.Link
	var hasValue bool
	if c == '(' { // normal link
		link = s.parseLink(parent, last, block, pc)
	} else if c == '[' { // reference link
		link, hasValue = s.parseReferenceLink(parent, last, block, pc)
		if link == nil && hasValue {
			ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
			return nil
		}
	}

	if link == nil {
		// maybe shortcut reference link
		block.SetPosition(l, pos)
		ssegment := text.NewSegment(last.Segment.Stop, segment.Start)
		maybeReference := block.Value(ssegment)
		// CommonMark spec says:
		//  > A link label can have at most 999 characters inside the square brackets.
		if len(maybeReference) > 999 {
			ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
			return nil
		}

		ref, ok := pc.Reference(util.ToLinkReference(maybeReference))
		if !ok {
			ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
			return nil
		}
		link = ast.NewLink()
		s.processLinkLabel(parent, link, last, pc)
		link.Title = ref.Title()
		link.Destination = ref.Destination()
	}
	if last.IsImage {
		last.Parent().RemoveChild(last.Parent(), last)
		return ast.NewImage(link)
	}
	last.Parent().RemoveChild(last.Parent(), last)
	return link
}

func (s *linkParser) containsLink(n ast.Node) bool {
	if n == nil {
		return false
	}
	for c := n; c != nil; c = c.NextSibling() {
		if _, ok := c.(*ast.Link); ok {
			return true
		}
		if s.containsLink(c.FirstChild()) {
			return true
		}
	}
	return false
}

func processLinkLabelOpen(block text.Reader, pos int, isImage bool, pc Context) *linkLabelState {
	start := pos
	if isImage {
		start--
	}
	state := newLinkLabelState(text.NewSegment(start, pos+1), isImage)
	pushLinkLabelState(pc, state)
	block.Advance(1)
	return state
}

func (s *linkParser) processLinkLabel(parent ast.Node, link *ast.Link, last *linkLabelState, pc Context) {
	var bottom ast.Node
	if v := pc.Get(linkBottom); v != nil {
		bottom = v.(ast.Node)
	}
	pc.Set(linkBottom, nil)
	ProcessDelimiters(bottom, pc)
	for c := last.NextSibling(); c != nil; {
		next := c.NextSibling()
		parent.RemoveChild(parent, c)
		link.AppendChild(link, c)
		c = next
	}
}

var linkFindClosureOptions text.FindClosureOptions = text.FindClosureOptions{
	Nesting: false,
	Newline: true,
	Advance: true,
}

func (s *linkParser) parseReferenceLink(parent ast.Node, last *linkLabelState, block text.Reader, pc Context) (*ast.Link, bool) {
	_, orgpos := block.Position()
	block.Advance(1) // skip '['
	segments, found := block.FindClosure('[', ']', linkFindClosureOptions)
	if !found {
		return nil, false
	}

	var maybeReference []byte
	if segments.Len() == 1 { // avoid allocate a new byte slice
		maybeReference = block.Value(segments.At(0))
	} else {
		maybeReference = []byte{}
		for i := 0; i < segments.Len(); i++ {
			s := segments.At(i)
			maybeReference = append(maybeReference, block.Value(s)...)
		}
	}
	if util.IsBlank(maybeReference) { // collapsed reference link
		s := text.NewSegment(last.Segment.Stop, orgpos.Start-1)
		maybeReference = block.Value(s)
	}
	// CommonMark spec says:
	//  > A link label can have at most 999 characters inside the square brackets.
	if len(maybeReference) > 999 {
		return nil, true
	}

	ref, ok := pc.Reference(util.ToLinkReference(maybeReference))
	if !ok {
		return nil, true
	}

	link := ast.NewLink()
	s.processLinkLabel(parent, link, last, pc)
	link.Title = ref.Title()
	link.Destination = ref.Destination()
	return link, true
}

func (s *linkParser) parseLink(parent ast.Node, last *linkLabelState, block text.Reader, pc Context) *ast.Link {
	block.Advance(1) // skip '('
	block.SkipSpaces()
	var title []byte
	var destination []byte
	var ok bool
	if block.Peek() == ')' { // empty link like '[link]()'
		block.Advance(1)
	} else {
		destination, ok = parseLinkDestination(block)
		if !ok {
			return nil
		}
		block.SkipSpaces()
		if block.Peek() == ')' {
			block.Advance(1)
		} else {
			title, ok = parseLinkTitle(block)
			if !ok {
				return nil
			}
			block.SkipSpaces()
			if block.Peek() == ')' {
				block.Advance(1)
			} else {
				return nil
			}
		}
	}

	link := ast.NewLink()
	s.processLinkLabel(parent, link, last, pc)
	link.Destination = destination
	link.Title = title
	return link
}

func parseLinkDestination(block text.Reader) ([]byte, bool) {
	block.SkipSpaces()
	line, _ := block.PeekLine()
	if block.Peek() == '<' {
		i := 1
		for i < len(line) {
			c := line[i]
			if c == '\\' && i < len(line)-1 && util.IsPunct(line[i+1]) {
				i += 2
				continue
			} else if c == '>' {
				block.Advance(i + 1)
				return line[1:i], true
			}
			i++
		}
		return nil, false
	}
	opened := 0
	i := 0
	for i < len(line) {
		c := line[i]
		if c == '\\' && i < len(line)-1 && util.IsPunct(line[i+1]) {
			i += 2
			continue
		} else if c == '(' {
			opened++
		} else if c == ')' {
			opened--
			if opened < 0 {
				break
			}
		} else if util.IsSpace(c) {
			break
		}
		i++
	}
	block.Advance(i)
	return line[:i], len(line[:i]) != 0
}

func parseLinkTitle(block text.Reader) ([]byte, bool) {
	block.SkipSpaces()
	opener := block.Peek()
	if opener != '"' && opener != '\'' && opener != '(' {
		return nil, false
	}
	closer := opener
	if opener == '(' {
		closer = ')'
	}
	block.Advance(1)
	segments, found := block.FindClosure(opener, closer, linkFindClosureOptions)
	if found {
		if segments.Len() == 1 {
			return block.Value(segments.At(0)), true
		}
		var title []byte
		for i := 0; i < segments.Len(); i++ {
			s := segments.At(i)
			title = append(title, block.Value(s)...)
		}
		return title, true
	}
	return nil, false
}

func (s *linkParser) CloseBlock(parent ast.Node, block text.Reader, pc Context) {
	pc.Set(linkBottom, nil)
	tlist := pc.Get(linkLabelStateKey)
	if tlist == nil {
		return
	}
	for s := tlist.(*linkLabelState); s != nil; {
		next := s.Next
		removeLinkLabelState(pc, s)
		s.Parent().ReplaceChild(s.Parent(), s, ast.NewTextSegment(s.Segment))
		s = next
	}
}
