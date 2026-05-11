// Copyright (c) 2014 The go-patricia AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package patricia

import (
	"io"
	"sort"
)

type childList interface {
	length() int
	head() *Trie
	add(child *Trie) childList
	remove(b byte)
	replace(b byte, child *Trie)
	next(b byte) *Trie
	combinedMask() uint64
	getChildren() []*Trie
	walk(prefix *Prefix, visitor VisitorFunc) error
	print(w io.Writer, indent int)
	clone() childList
	total() int
}

type tries []*Trie

func (t tries) Len() int {
	return len(t)
}

func (t tries) Less(i, j int) bool {
	strings := sort.StringSlice{string(t[i].prefix), string(t[j].prefix)}
	return strings.Less(0, 1)
}

func (t tries) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type childContainer struct {
	char byte
	node *Trie
}
type superDenseChildList struct {
	children []childContainer
}

func newSuperDenseChildList() childList {
	return &superDenseChildList{
		make([]childContainer, 0),
	}
}

func (list *superDenseChildList) length() int {
	return len(list.children)
}

func (list *superDenseChildList) head() *Trie {
	if len(list.children) > 0 {
		return list.children[0].node
	}
	return nil
}

func (list *superDenseChildList) add(child *Trie) childList {
	char := child.prefix[0]
	list.children = append(list.children, childContainer{
		char,
		child,
	})
	return list
}

func (list *superDenseChildList) remove(b byte) {
	children := list.children
	for i := 0; i < len(children); i++ {
		if children[i].char == b {
			// children[i] = children[len(children)-1]
			// children = children[:len(children)-1]
			// children = append(children[:i], children[i+1:]...)
			newChildren := make([]childContainer, len(children)-1)
			// copy the elements over to avoid "memory leaks"
			copy(newChildren, children[:i])
			copy(newChildren[i:], children[i+1:])
			list.children = newChildren

			// list.children = make([]childContainer, len(children))
			// copy(list.children, children)
			// children = nil

			return
		}
	}
}

func (list *superDenseChildList) replace(b byte, child *Trie) {
	children := list.children
	for i := 0; i < len(list.children); i++ {
		if children[i].char == b {
			children[i].node = child
			return
		}
	}
}

func (list *superDenseChildList) next(b byte) *Trie {
	children := list.children
	for i := 0; i < len(list.children); i++ {
		if children[i].char == b {
			return children[i].node
		}
	}
	return nil
}

func (list superDenseChildList) combinedMask() uint64 {
	var mask uint64
	for _, child := range list.children {
		// fmt.Printf("child = %+v\n", child)
		mask |= child.node.mask
	}
	return mask
}

func (list *superDenseChildList) getChildren() []*Trie {
	children := make([]*Trie, 0, len(list.children))
	for _, child := range list.children {
		children = append(children, child.node)
	}
	return children
}

func (list *superDenseChildList) walk(prefix *Prefix, visitor VisitorFunc) error {
	for _, child := range list.children {
		node := child.node
		*prefix = append(*prefix, node.prefix...)
		if node.item != nil {
			if err := visitor(*prefix, node.item); err != nil {
				if err == SkipSubtree {
					*prefix = (*prefix)[:len(*prefix)-len(node.prefix)]
					continue
				}
				*prefix = (*prefix)[:len(*prefix)-len(node.prefix)]
				return err
			}
		}

		err := node.children.walk(prefix, visitor)
		*prefix = (*prefix)[:len(*prefix)-len(node.prefix)]
		if err != nil {
			return err
		}
	}

	return nil
}

func (list *superDenseChildList) print(w io.Writer, indent int) {
	for _, child := range list.children {
		child.node.print(w, indent)
	}
}

func (list *superDenseChildList) clone() childList {
	clones := make([]childContainer, len(list.children))

	for i := 0; i < len(list.children); i++ {
		child := list.children[i]
		clones[i] = childContainer{child.char, child.node.Clone()}
	}

	return &superDenseChildList{
		clones,
	}
}

func (list *superDenseChildList) total() int {
	return len(list.children)
}
