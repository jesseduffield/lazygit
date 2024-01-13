package parser

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type linkReferenceParagraphTransformer struct {
}

// LinkReferenceParagraphTransformer is a ParagraphTransformer implementation
// that parses and extracts link reference from paragraphs.
var LinkReferenceParagraphTransformer = &linkReferenceParagraphTransformer{}

func (p *linkReferenceParagraphTransformer) Transform(node *ast.Paragraph, reader text.Reader, pc Context) {
	lines := node.Lines()
	block := text.NewBlockReader(reader.Source(), lines)
	removes := [][2]int{}
	for {
		start, end := parseLinkReferenceDefinition(block, pc)
		if start > -1 {
			if start == end {
				end++
			}
			removes = append(removes, [2]int{start, end})
			continue
		}
		break
	}

	offset := 0
	for _, remove := range removes {
		if lines.Len() == 0 {
			break
		}
		s := lines.Sliced(remove[1]-offset, lines.Len())
		lines.SetSliced(0, remove[0]-offset)
		lines.AppendAll(s)
		offset = remove[1]
	}

	if lines.Len() == 0 {
		t := ast.NewTextBlock()
		t.SetBlankPreviousLines(node.HasBlankPreviousLines())
		node.Parent().ReplaceChild(node.Parent(), node, t)
		return
	}

	node.SetLines(lines)
}

func parseLinkReferenceDefinition(block text.Reader, pc Context) (int, int) {
	block.SkipSpaces()
	line, _ := block.PeekLine()
	if line == nil {
		return -1, -1
	}
	startLine, _ := block.Position()
	width, pos := util.IndentWidth(line, 0)
	if width > 3 {
		return -1, -1
	}
	if width != 0 {
		pos++
	}
	if line[pos] != '[' {
		return -1, -1
	}
	block.Advance(pos + 1)
	segments, found := block.FindClosure('[', ']', linkFindClosureOptions)
	if !found {
		return -1, -1
	}
	var label []byte
	if segments.Len() == 1 {
		label = block.Value(segments.At(0))
	} else {
		for i := 0; i < segments.Len(); i++ {
			s := segments.At(i)
			label = append(label, block.Value(s)...)
		}
	}
	if util.IsBlank(label) {
		return -1, -1
	}
	if block.Peek() != ':' {
		return -1, -1
	}
	block.Advance(1)
	block.SkipSpaces()
	destination, ok := parseLinkDestination(block)
	if !ok {
		return -1, -1
	}
	line, _ = block.PeekLine()
	isNewLine := line == nil || util.IsBlank(line)

	endLine, _ := block.Position()
	_, spaces, _ := block.SkipSpaces()
	opener := block.Peek()
	if opener != '"' && opener != '\'' && opener != '(' {
		if !isNewLine {
			return -1, -1
		}
		ref := NewReference(label, destination, nil)
		pc.AddReference(ref)
		return startLine, endLine + 1
	}
	if spaces == 0 {
		return -1, -1
	}
	block.Advance(1)
	closer := opener
	if opener == '(' {
		closer = ')'
	}
	segments, found = block.FindClosure(opener, closer, linkFindClosureOptions)
	if !found {
		if !isNewLine {
			return -1, -1
		}
		ref := NewReference(label, destination, nil)
		pc.AddReference(ref)
		block.AdvanceLine()
		return startLine, endLine + 1
	}
	var title []byte
	if segments.Len() == 1 {
		title = block.Value(segments.At(0))
	} else {
		for i := 0; i < segments.Len(); i++ {
			s := segments.At(i)
			title = append(title, block.Value(s)...)
		}
	}

	line, _ = block.PeekLine()
	if line != nil && !util.IsBlank(line) {
		if !isNewLine {
			return -1, -1
		}
		ref := NewReference(label, destination, title)
		pc.AddReference(ref)
		return startLine, endLine
	}

	endLine, _ = block.Position()
	ref := NewReference(label, destination, title)
	pc.AddReference(ref)
	return startLine, endLine + 1
}
