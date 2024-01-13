package ansi

import (
	"bytes"
)

// BlockStack is a stack of block elements, used to calculate the current
// indentation & margin level during the rendering process.
type BlockStack []BlockElement

// Len returns the length of the stack.
func (s *BlockStack) Len() int {
	return len(*s)
}

// Push appends an item to the stack.
func (s *BlockStack) Push(e BlockElement) {
	*s = append(*s, e)
}

// Pop removes the last item on the stack.
func (s *BlockStack) Pop() {
	stack := *s
	if len(stack) == 0 {
		return
	}

	stack = stack[0 : len(stack)-1]
	*s = stack
}

// Indent returns the current indentation level of all elements in the stack.
func (s BlockStack) Indent() uint {
	var i uint

	for _, v := range s {
		if v.Style.Indent == nil {
			continue
		}
		i += *v.Style.Indent
	}

	return i
}

// Margin returns the current margin level of all elements in the stack.
func (s BlockStack) Margin() uint {
	var i uint

	for _, v := range s {
		if v.Style.Margin == nil {
			continue
		}
		i += *v.Style.Margin
	}

	return i
}

// Width returns the available rendering width
func (s BlockStack) Width(ctx RenderContext) uint {
	if s.Indent()+s.Margin()*2 > uint(ctx.options.WordWrap) {
		return 0
	}
	return uint(ctx.options.WordWrap) - s.Indent() - s.Margin()*2
}

// Parent returns the current BlockElement's parent.
func (s BlockStack) Parent() BlockElement {
	if len(s) == 1 {
		return BlockElement{
			Block: &bytes.Buffer{},
		}
	}

	return s[len(s)-2]
}

// Current returns the current BlockElement.
func (s BlockStack) Current() BlockElement {
	if len(s) == 0 {
		return BlockElement{
			Block: &bytes.Buffer{},
		}
	}

	return s[len(s)-1]
}

// With returns a StylePrimitive that inherits the current BlockElement's style.
func (s BlockStack) With(child StylePrimitive) StylePrimitive {
	sb := StyleBlock{}
	sb.StylePrimitive = child
	return cascadeStyle(s.Current().Style, sb, false).StylePrimitive
}
