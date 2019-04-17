package query

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

// base match
type matchBase struct {
	next pathFn
}

func (f *matchBase) setNext(next pathFn) {
	f.next = next
}

// terminating functor - gathers results
type terminatingFn struct {
	// empty
}

func newTerminatingFn() *terminatingFn {
	return &terminatingFn{}
}

func (f *terminatingFn) setNext(next pathFn) {
	// do nothing
}

func (f *terminatingFn) call(node interface{}, ctx *queryContext) {
	ctx.result.appendResult(node, ctx.lastPosition)
}

// match single key
type matchKeyFn struct {
	matchBase
	Name string
}

func newMatchKeyFn(name string) *matchKeyFn {
	return &matchKeyFn{Name: name}
}

func (f *matchKeyFn) call(node interface{}, ctx *queryContext) {
	if array, ok := node.([]*toml.Tree); ok {
		for _, tree := range array {
			item := tree.Get(f.Name)
			if item != nil {
				ctx.lastPosition = tree.GetPosition(f.Name)
				f.next.call(item, ctx)
			}
		}
	} else if tree, ok := node.(*toml.Tree); ok {
		item := tree.Get(f.Name)
		if item != nil {
			ctx.lastPosition = tree.GetPosition(f.Name)
			f.next.call(item, ctx)
		}
	}
}

// match single index
type matchIndexFn struct {
	matchBase
	Idx int
}

func newMatchIndexFn(idx int) *matchIndexFn {
	return &matchIndexFn{Idx: idx}
}

func (f *matchIndexFn) call(node interface{}, ctx *queryContext) {
	if arr, ok := node.([]interface{}); ok {
		if f.Idx < len(arr) && f.Idx >= 0 {
			if treesArray, ok := node.([]*toml.Tree); ok {
				if len(treesArray) > 0 {
					ctx.lastPosition = treesArray[0].Position()
				}
			}
			f.next.call(arr[f.Idx], ctx)
		}
	}
}

// filter by slicing
type matchSliceFn struct {
	matchBase
	Start, End, Step int
}

func newMatchSliceFn(start, end, step int) *matchSliceFn {
	return &matchSliceFn{Start: start, End: end, Step: step}
}

func (f *matchSliceFn) call(node interface{}, ctx *queryContext) {
	if arr, ok := node.([]interface{}); ok {
		// adjust indexes for negative values, reverse ordering
		realStart, realEnd := f.Start, f.End
		if realStart < 0 {
			realStart = len(arr) + realStart
		}
		if realEnd < 0 {
			realEnd = len(arr) + realEnd
		}
		if realEnd < realStart {
			realEnd, realStart = realStart, realEnd // swap
		}
		// loop and gather
		for idx := realStart; idx < realEnd; idx += f.Step {
			if treesArray, ok := node.([]*toml.Tree); ok {
				if len(treesArray) > 0 {
					ctx.lastPosition = treesArray[0].Position()
				}
			}
			f.next.call(arr[idx], ctx)
		}
	}
}

// match anything
type matchAnyFn struct {
	matchBase
}

func newMatchAnyFn() *matchAnyFn {
	return &matchAnyFn{}
}

func (f *matchAnyFn) call(node interface{}, ctx *queryContext) {
	if tree, ok := node.(*toml.Tree); ok {
		for _, k := range tree.Keys() {
			v := tree.Get(k)
			ctx.lastPosition = tree.GetPosition(k)
			f.next.call(v, ctx)
		}
	}
}

// filter through union
type matchUnionFn struct {
	Union []pathFn
}

func (f *matchUnionFn) setNext(next pathFn) {
	for _, fn := range f.Union {
		fn.setNext(next)
	}
}

func (f *matchUnionFn) call(node interface{}, ctx *queryContext) {
	for _, fn := range f.Union {
		fn.call(node, ctx)
	}
}

// match every single last node in the tree
type matchRecursiveFn struct {
	matchBase
}

func newMatchRecursiveFn() *matchRecursiveFn {
	return &matchRecursiveFn{}
}

func (f *matchRecursiveFn) call(node interface{}, ctx *queryContext) {
	originalPosition := ctx.lastPosition
	if tree, ok := node.(*toml.Tree); ok {
		var visit func(tree *toml.Tree)
		visit = func(tree *toml.Tree) {
			for _, k := range tree.Keys() {
				v := tree.Get(k)
				ctx.lastPosition = tree.GetPosition(k)
				f.next.call(v, ctx)
				switch node := v.(type) {
				case *toml.Tree:
					visit(node)
				case []*toml.Tree:
					for _, subtree := range node {
						visit(subtree)
					}
				}
			}
		}
		ctx.lastPosition = originalPosition
		f.next.call(tree, ctx)
		visit(tree)
	}
}

// match based on an externally provided functional filter
type matchFilterFn struct {
	matchBase
	Pos  toml.Position
	Name string
}

func newMatchFilterFn(name string, pos toml.Position) *matchFilterFn {
	return &matchFilterFn{Name: name, Pos: pos}
}

func (f *matchFilterFn) call(node interface{}, ctx *queryContext) {
	fn, ok := (*ctx.filters)[f.Name]
	if !ok {
		panic(fmt.Sprintf("%s: query context does not have filter '%s'",
			f.Pos.String(), f.Name))
	}
	switch castNode := node.(type) {
	case *toml.Tree:
		for _, k := range castNode.Keys() {
			v := castNode.Get(k)
			if fn(v) {
				ctx.lastPosition = castNode.GetPosition(k)
				f.next.call(v, ctx)
			}
		}
	case []*toml.Tree:
		for _, v := range castNode {
			if fn(v) {
				if len(castNode) > 0 {
					ctx.lastPosition = castNode[0].Position()
				}
				f.next.call(v, ctx)
			}
		}
	case []interface{}:
		for _, v := range castNode {
			if fn(v) {
				f.next.call(v, ctx)
			}
		}
	}
}
