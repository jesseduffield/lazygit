// Package str is a comprehensive set of string functions to build more
// Go awesomeness. Str complements Go's standard packages and does not duplicate
// functionality found in `strings` or  `strconv`.
//
// Str is based on plain functions instead of object-based methods,
// consistent with Go standard string packages.
//
//      str.Between("<a>foo</a>", "<a>", "</a>") == "foo"
//
// Str supports pipelining instead of chaining
//
//      s := str.Pipe("\nabcdef\n", Clean, BetweenF("a", "f"), ChompLeftF("bc"))
//
// User-defined filters can be added to the pipeline by inserting a function
// or closure that returns a function with this signature
//
//      func(string) string
//
package str
