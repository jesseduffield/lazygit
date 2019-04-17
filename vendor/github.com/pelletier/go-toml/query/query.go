package query

import (
	"time"

	"github.com/pelletier/go-toml"
)

// NodeFilterFn represents a user-defined filter function, for use with
// Query.SetFilter().
//
// The return value of the function must indicate if 'node' is to be included
// at this stage of the TOML path.  Returning true will include the node, and
// returning false will exclude it.
//
// NOTE: Care should be taken to write script callbacks such that they are safe
// to use from multiple goroutines.
type NodeFilterFn func(node interface{}) bool

// Result is the result of Executing a Query.
type Result struct {
	items     []interface{}
	positions []toml.Position
}

// appends a value/position pair to the result set.
func (r *Result) appendResult(node interface{}, pos toml.Position) {
	r.items = append(r.items, node)
	r.positions = append(r.positions, pos)
}

// Values is a set of values within a Result.  The order of values is not
// guaranteed to be in document order, and may be different each time a query is
// executed.
func (r Result) Values() []interface{} {
	return r.items
}

// Positions is a set of positions for values within a Result.  Each index
// in Positions() corresponds to the entry in Value() of the same index.
func (r Result) Positions() []toml.Position {
	return r.positions
}

// runtime context for executing query paths
type queryContext struct {
	result       *Result
	filters      *map[string]NodeFilterFn
	lastPosition toml.Position
}

// generic path functor interface
type pathFn interface {
	setNext(next pathFn)
	// it is the caller's responsibility to set the ctx.lastPosition before invoking call()
	// node can be one of: *toml.Tree, []*toml.Tree, or a scalar
	call(node interface{}, ctx *queryContext)
}

// A Query is the representation of a compiled TOML path.  A Query is safe
// for concurrent use by multiple goroutines.
type Query struct {
	root    pathFn
	tail    pathFn
	filters *map[string]NodeFilterFn
}

func newQuery() *Query {
	return &Query{
		root:    nil,
		tail:    nil,
		filters: &defaultFilterFunctions,
	}
}

func (q *Query) appendPath(next pathFn) {
	if q.root == nil {
		q.root = next
	} else {
		q.tail.setNext(next)
	}
	q.tail = next
	next.setNext(newTerminatingFn()) // init the next functor
}

// Compile compiles a TOML path expression. The returned Query can be used
// to match elements within a Tree and its descendants. See Execute.
func Compile(path string) (*Query, error) {
	return parseQuery(lexQuery(path))
}

// Execute executes a query against a Tree, and returns the result of the query.
func (q *Query) Execute(tree *toml.Tree) *Result {
	result := &Result{
		items:     []interface{}{},
		positions: []toml.Position{},
	}
	if q.root == nil {
		result.appendResult(tree, tree.GetPosition(""))
	} else {
		ctx := &queryContext{
			result:  result,
			filters: q.filters,
		}
		ctx.lastPosition = tree.Position()
		q.root.call(tree, ctx)
	}
	return result
}

// CompileAndExecute is a shorthand for Compile(path) followed by Execute(tree).
func CompileAndExecute(path string, tree *toml.Tree) (*Result, error) {
	query, err := Compile(path)
	if err != nil {
		return nil, err
	}
	return query.Execute(tree), nil
}

// SetFilter sets a user-defined filter function.  These may be used inside
// "?(..)" query expressions to filter TOML document elements within a query.
func (q *Query) SetFilter(name string, fn NodeFilterFn) {
	if q.filters == &defaultFilterFunctions {
		// clone the static table
		q.filters = &map[string]NodeFilterFn{}
		for k, v := range defaultFilterFunctions {
			(*q.filters)[k] = v
		}
	}
	(*q.filters)[name] = fn
}

var defaultFilterFunctions = map[string]NodeFilterFn{
	"tree": func(node interface{}) bool {
		_, ok := node.(*toml.Tree)
		return ok
	},
	"int": func(node interface{}) bool {
		_, ok := node.(int64)
		return ok
	},
	"float": func(node interface{}) bool {
		_, ok := node.(float64)
		return ok
	},
	"string": func(node interface{}) bool {
		_, ok := node.(string)
		return ok
	},
	"time": func(node interface{}) bool {
		_, ok := node.(time.Time)
		return ok
	},
	"bool": func(node interface{}) bool {
		_, ok := node.(bool)
		return ok
	},
}
