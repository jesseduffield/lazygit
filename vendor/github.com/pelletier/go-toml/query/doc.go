// Package query performs JSONPath-like queries on a TOML document.
//
// The query path implementation is based loosely on the JSONPath specification:
// http://goessner.net/articles/JsonPath/.
//
// The idea behind a query path is to allow quick access to any element, or set
// of elements within TOML document, with a single expression.
//
//   result, err := query.CompileAndExecute("$.foo.bar.baz", tree)
//
// This is roughly equivalent to:
//
//   next := tree.Get("foo")
//   if next != nil {
//     next = next.Get("bar")
//     if next != nil {
//       next = next.Get("baz")
//     }
//   }
//   result := next
//
// err is nil if any parsing exception occurs.
//
// If no node in the tree matches the query, result will simply contain an empty list of
// items.
//
// As illustrated above, the query path is much more efficient, especially since
// the structure of the TOML file can vary.  Rather than making assumptions about
// a document's structure, a query allows the programmer to make structured
// requests into the document, and get zero or more values as a result.
//
// Query syntax
//
// The syntax of a query begins with a root token, followed by any number
// sub-expressions:
//
//   $
//                    Root of the TOML tree.  This must always come first.
//   .name
//                    Selects child of this node, where 'name' is a TOML key
//                    name.
//   ['name']
//                    Selects child of this node, where 'name' is a string
//                    containing a TOML key name.
//   [index]
//                    Selcts child array element at 'index'.
//   ..expr
//                    Recursively selects all children, filtered by an a union,
//                    index, or slice expression.
//   ..*
//                    Recursive selection of all nodes at this point in the
//                    tree.
//   .*
//                    Selects all children of the current node.
//   [expr,expr]
//                    Union operator - a logical 'or' grouping of two or more
//                    sub-expressions: index, key name, or filter.
//   [start:end:step]
//                    Slice operator - selects array elements from start to
//                    end-1, at the given step.  All three arguments are
//                    optional.
//   [?(filter)]
//                    Named filter expression - the function 'filter' is
//                    used to filter children at this node.
//
// Query Indexes And Slices
//
// Index expressions perform no bounds checking, and will contribute no
// values to the result set if the provided index or index range is invalid.
// Negative indexes represent values from the end of the array, counting backwards.
//
//   // select the last index of the array named 'foo'
//   query.CompileAndExecute("$.foo[-1]", tree)
//
// Slice expressions are supported, by using ':' to separate a start/end index pair.
//
//   // select up to the first five elements in the array
//   query.CompileAndExecute("$.foo[0:5]", tree)
//
// Slice expressions also allow negative indexes for the start and stop
// arguments.
//
//   // select all array elements.
//   query.CompileAndExecute("$.foo[0:-1]", tree)
//
// Slice expressions may have an optional stride/step parameter:
//
//   // select every other element
//   query.CompileAndExecute("$.foo[0:-1:2]", tree)
//
// Slice start and end parameters are also optional:
//
//   // these are all equivalent and select all the values in the array
//   query.CompileAndExecute("$.foo[:]", tree)
//   query.CompileAndExecute("$.foo[0:]", tree)
//   query.CompileAndExecute("$.foo[:-1]", tree)
//   query.CompileAndExecute("$.foo[0:-1:]", tree)
//   query.CompileAndExecute("$.foo[::1]", tree)
//   query.CompileAndExecute("$.foo[0::1]", tree)
//   query.CompileAndExecute("$.foo[:-1:1]", tree)
//   query.CompileAndExecute("$.foo[0:-1:1]", tree)
//
// Query Filters
//
// Query filters are used within a Union [,] or single Filter [] expression.
// A filter only allows nodes that qualify through to the next expression,
// and/or into the result set.
//
//   // returns children of foo that are permitted by the 'bar' filter.
//   query.CompileAndExecute("$.foo[?(bar)]", tree)
//
// There are several filters provided with the library:
//
//   tree
//          Allows nodes of type Tree.
//   int
//          Allows nodes of type int64.
//   float
//          Allows nodes of type float64.
//   string
//          Allows nodes of type string.
//   time
//          Allows nodes of type time.Time.
//   bool
//          Allows nodes of type bool.
//
// Query Results
//
// An executed query returns a Result object.  This contains the nodes
// in the TOML tree that qualify the query expression.  Position information
// is also available for each value in the set.
//
//   // display the results of a query
//   results := query.CompileAndExecute("$.foo.bar.baz", tree)
//   for idx, value := results.Values() {
//       fmt.Println("%v: %v", results.Positions()[idx], value)
//   }
//
// Compiled Queries
//
// Queries may be executed directly on a Tree object, or compiled ahead
// of time and executed discretely.  The former is more convenient, but has the
// penalty of having to recompile the query expression each time.
//
//   // basic query
//   results := query.CompileAndExecute("$.foo.bar.baz", tree)
//
//   // compiled query
//   query, err := toml.Compile("$.foo.bar.baz")
//   results := query.Execute(tree)
//
//   // run the compiled query again on a different tree
//   moreResults := query.Execute(anotherTree)
//
// User Defined Query Filters
//
// Filter expressions may also be user defined by using the SetFilter()
// function on the Query object.  The function must return true/false, which
// signifies if the passed node is kept or discarded, respectively.
//
//   // create a query that references a user-defined filter
//   query, _ := query.Compile("$[?(bazOnly)]")
//
//   // define the filter, and assign it to the query
//   query.SetFilter("bazOnly", func(node interface{}) bool{
//       if tree, ok := node.(*Tree); ok {
//           return tree.Has("baz")
//       }
//       return false  // reject all other node types
//   })
//
//   // run the query
//   query.Execute(tree)
//
package query
