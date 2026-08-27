// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

// Package format exposes gofumpt's formatting in an API similar to go/format.
// In general, the APIs are only guaranteed to work well when the input source
// is in canonical gofmt format.
package format

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	goversion "go/version"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/ast/astutil"

	"mvdan.cc/gofumpt/internal/govendor/go/format"
	"mvdan.cc/gofumpt/internal/version"
)

// Options is the set of formatting options which affect gofumpt.
type Options struct {
	// LangVersion is the Go version a piece of code is written in.
	// The version is used to decide whether to apply formatting
	// rules which require new language features.
	// When empty, a default of go1 is assumed.
	// Otherwise, the version must satisfy [go/version.IsValid].
	//
	// When formatting a Go module, LangVersion should typically be
	//
	//     go list -m -f {{.GoVersion}}
	//
	// with a "go" prefix, or the equivalent from `go mod edit -json`.
	LangVersion string

	// ModulePath corresponds to the Go module path which contains the source
	// code being formatted. When formatting a Go module, ModulePath should be
	//
	//     go list -m -f {{.Path}}
	//
	// or the equivalent from `go mod edit -json`.
	//
	// ModulePath is used for formatting decisions like what import paths are
	// considered to be not part of the standard library. When empty, the source
	// is formatted as if it weren't inside a module.
	ModulePath string

	// ExtraRules enables all extra formatting rules, such as grouping function
	// parameters with repeated types together.
	//
	// Deprecated: use [Options.Extra] instead.
	ExtraRules bool

	// Extra allows enabling extra formatting rules which are disabled by default.
	Extra Extra
}

// Extra is the set of extra formatting rules which are available.
//
// As the formatter evolves, we might add or remove boolean fields here.
// Go API users who wish to avoid build errors in such cases
// can use the string API in [Extra.Set].
type Extra struct {
	// TODO: should we have "All" to turn them all on,
	// akin to how the CLI has -extra=true for historical reasons?
	// I lean against it, as it should be a conscious choice to turn on
	// each of these extra rules, and we should be able to add more rules
	// without fear of causing unexpected changes for users.

	// GroupParams groups function parameters with repeated types.
	GroupParams bool

	// ClotheReturns clothes naked returns in functions with named results.
	ClotheReturns bool

	// BalanceCalls places a multi-line call's closing parenthesis on its
	// own line when the opening parenthesis ends a line.
	BalanceCalls bool
}

func (e *Extra) String() string {
	var active []string
	if e.GroupParams {
		active = append(active, "group_params")
	}
	if e.ClotheReturns {
		active = append(active, "clothe_returns")
	}
	if e.BalanceCalls {
		active = append(active, "balance_calls")
	}
	return strings.Join(active, ",")
}

func (e *Extra) Set(v string) error {
	if v == "true" {
		e.GroupParams = true
		e.ClotheReturns = true
		e.BalanceCalls = true
		return nil
	}
	*e = Extra{}
	if v == "false" {
		return nil
	}
	for s := range strings.SplitSeq(v, ",") {
		switch s {
		case "group_params":
			e.GroupParams = true
		case "clothe_returns":
			e.ClotheReturns = true
		case "balance_calls":
			e.BalanceCalls = true
		default:
			return fmt.Errorf("unknown rule: %q", s)
		}
	}
	return nil
}

func (e *Extra) IsBoolFlag() bool { return true }

// Source formats src in gofumpt's format, assuming that src holds a valid Go
// source file.
func Source(src []byte, opts Options) ([]byte, error) {
	fset := token.NewFileSet()

	// Ensure our parsed files never start with base 1,
	// to ensure that using token.NoPos+1 will panic.
	fset.AddFile("gofumpt_base.go", 1, 10)

	file, err := parser.ParseFile(fset, "", src, parser.SkipObjectResolution|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	File(fset, file, opts)

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// File modifies a file and fset in place to follow gofumpt's format. The
// changes might include manipulating adding or removing newlines in fset,
// modifying the position of nodes, or modifying literal values.
func File(fset *token.FileSet, file *ast.File, opts Options) {
	simplify(file)

	if opts.ExtraRules {
		opts.Extra.Set("true") // enable all the extra rules
	}

	if opts.LangVersion == "" {
		opts.LangVersion = "go1"
	} else {
		lang := goversion.Lang(opts.LangVersion)
		if lang == "" {
			panic(fmt.Sprintf("invalid Go version: %q", opts.LangVersion))
		}
		opts.LangVersion = lang
	}
	f := &fumpter{
		file:    fset.File(file.Pos()),
		fset:    fset,
		astFile: file,
		Options: opts,

		minSplitFactor: 0.4,
	}
	var topFuncType *ast.FuncType
	pre := func(c *astutil.Cursor) bool {
		f.applyPre(c)
		switch node := c.Node().(type) {
		case *ast.FuncDecl:
			topFuncType = node.Type
			f.parentFuncTypes = append(f.parentFuncTypes, node.Type)
		case *ast.FuncLit:
			f.parentFuncTypes = append(f.parentFuncTypes, node.Type)
		case *ast.FieldList:
			ft, _ := c.Parent().(*ast.FuncType)
			if ft == nil || ft != topFuncType {
				break
			}

			// For top-level function declaration parameters,
			// require the line split to be longer.
			// This avoids func lines which are a bit too short,
			// and allows func lines which are a bit longer.
			//
			// We don't just increase longLineLimit,
			// as we still want splits at around the same place.
			if ft.Params == node {
				f.minSplitFactor = 0.6
			}

			// Don't split result parameters into multiple lines,
			// as that can be easily confused for input parameters.
			// TODO: consider the same for single-line func calls in
			// if statements.
			// TODO: perhaps just use a higher factor, like 0.8.
			if ft.Results == node {
				f.minSplitFactor = 1000
			}
		case *ast.BlockStmt:
			f.blockLevel++
		}
		return true
	}
	post := func(c *astutil.Cursor) bool {
		f.applyPost(c)

		// Reset minSplitFactor and blockLevel.
		switch node := c.Node().(type) {
		case *ast.FuncDecl, *ast.FuncLit:
			f.parentFuncTypes = f.parentFuncTypes[:len(f.parentFuncTypes)-1]
		case *ast.FuncType:
			if node == topFuncType {
				f.minSplitFactor = 0.4
			}
		case *ast.BlockStmt:
			f.blockLevel--
		}
		return true
	}
	astutil.Apply(file, pre, post)
}

// Multiline nodes which could easily fit on a single line under this many bytes
// may be collapsed onto a single line.
const shortLineLimit = 60

// Single-line nodes which take over this many bytes, and could easily be split
// into two lines of at least its minSplitFactor factor, may be split.
const longLineLimit = 100

var rxOctalInteger = regexp.MustCompile(`\A0[0-7_]+\z`)

type fumpter struct {
	Options

	file *token.File
	fset *token.FileSet

	astFile *ast.File

	// blockLevel is the number of indentation blocks we're currently under.
	// It is used to approximate the levels of indentation a line will end
	// up with.
	blockLevel int

	minSplitFactor float64

	// parentFuncTypes is a stack of parent function types,
	// used to determine return type information when clothing naked returns.
	parentFuncTypes []*ast.FuncType
}

func (f *fumpter) commentsBetween(p1, p2 token.Pos) []*ast.CommentGroup {
	comments := f.astFile.Comments
	i1 := sort.Search(len(comments), func(i int) bool {
		return comments[i].Pos() >= p1
	})
	comments = comments[i1:]
	i2 := sort.Search(len(comments), func(i int) bool {
		return comments[i].Pos() >= p2
	})
	comments = comments[:i2]
	return comments
}

func (f *fumpter) inlineComment(pos token.Pos) *ast.Comment {
	comments := f.astFile.Comments
	i := sort.Search(len(comments), func(i int) bool {
		return comments[i].Pos() >= pos
	})
	if i >= len(comments) {
		return nil
	}
	line := f.Line(pos)
	for _, comment := range comments[i].List {
		if f.Line(comment.Pos()) == line {
			return comment
		}
	}
	return nil
}

// addNewline is a hack to let us force a newline at a certain position.
func (f *fumpter) addNewline(at token.Pos) {
	offset := f.Offset(at)

	lines := f.file.Lines()
	i, exists := slices.BinarySearch(lines, offset)
	if exists {
		// This newline already exists; do nothing. Duplicate
		// newlines can't exist.
		return
	}
	lines = slices.Insert(lines, i, offset)
	if !f.file.SetLines(lines) {
		panic(fmt.Sprintf("could not set lines to %v", lines))
	}
}

// removeLines removes all newlines between two positions, so that they end
// up on the same line.
func (f *fumpter) removeLines(fromLine, toLine int) {
	for fromLine < toLine {
		f.file.MergeLine(fromLine)
		toLine--
	}
}

// removeLinesBetween is like removeLines, but it leaves one newline between the
// two positions.
func (f *fumpter) removeLinesBetween(from, to token.Pos) {
	f.removeLines(f.Line(from)+1, f.Line(to))
}

// removeParens unwraps a single-spec var group like "var (\n\tx = 1\n)" into a
// lone "var x = 1". It only acts on such groups without a doc comment.
func (f *fumpter) removeParens(node *ast.GenDecl) {
	if node.Tok != token.VAR || len(node.Specs) != 1 ||
		!node.Lparen.IsValid() || node.Doc != nil {
		return
	}
	specPos := node.Specs[0].Pos()
	specEnd := node.Specs[0].End()

	if len(f.commentsBetween(node.TokPos, specPos)) > 0 {
		// If the single spec has a comment on the line above,
		// the comment must go before the entire declaration now.
		node.TokPos = specPos
	} else {
		f.removeLines(f.Line(node.TokPos), f.Line(specPos))
	}
	if len(f.commentsBetween(specEnd, node.Rparen)) > 0 {
		// Leave one newline to not force a comment on the next line to
		// become an inline comment.
		f.removeLines(f.Line(specEnd)+1, f.Line(node.Rparen))
	} else {
		f.removeLines(f.Line(specEnd), f.Line(node.Rparen))
	}

	// Remove the parentheses. go/printer will automatically
	// get rid of the newlines.
	node.Lparen = token.NoPos
	node.Rparen = token.NoPos
}

func (f *fumpter) Position(p token.Pos) token.Position {
	return f.file.PositionFor(p, false)
}

func (f *fumpter) Line(p token.Pos) int {
	return f.Position(p).Line
}

func (f *fumpter) Offset(p token.Pos) int {
	return f.file.Offset(p)
}

type byteCounter int

func (b *byteCounter) Write(p []byte) (n int, err error) {
	*b += byteCounter(len(p))
	return len(p), nil
}

func (f *fumpter) printLength(node ast.Node) int {
	var count byteCounter
	if err := format.Node(&count, f.fset, node); err != nil {
		panic(fmt.Sprintf("unexpected print error: %v", err))
	}

	// Add the space taken by an inline comment.
	if c := f.inlineComment(node.End()); c != nil {
		fmt.Fprintf(&count, " %s", c.Text)
	}

	// Add an approximation of the indentation level. We can't know the
	// number of tabs go/printer will add ahead of time. Trying to print the
	// entire top-level declaration would tell us that, but then it's near
	// impossible to reliably find our node again.
	return int(count) + (f.blockLevel * 8)
}

func (f *fumpter) lineEnd(line int) token.Pos {
	if line < 1 {
		panic("illegal line number")
	}
	total := f.file.LineCount()
	if line > total {
		panic("illegal line number")
	}
	if line == total {
		return f.astFile.End()
	}
	return f.file.LineStart(line+1) - 1
}

// rxCommentDirective covers all common Go comment directives, such as:
//
//	//go:          | standard Go directives, like go:noinline
//	//some-words:  | similar to the syntax above, like lint:ignore or go-sumtype:decl
//	//export       | to mark cgo funcs for exporting
//	//extern       | C function declarations for gccgo
//	//line         | inserted line information for cmd/compile
//	//noinspection | noinspection directive for GoLand and friends
//	//nolint       | nolint directive for golangci
//	//#nosec       | #nosec directive for gosec
//	//NOSONAR      | NOSONAR directive for SonarQube
//	//sys(nb)?     | syscall function wrapper prototypes
var rxCommentDirective = regexp.MustCompile(
	`^(?:` +
		// Patterns directly from https://go.dev/doc/comment#syntax.
		// Note that we adjust the first pattern to allow for //go-sumtype:decl,
		// which is a tool that existed before the Go convention was documented.
		`[a-z0-9-]+:[a-z0-9]` +
		`|export ` +
		`|extern ` +
		`|line ` +
		// Third-party patterns; we generally assume they end with a word boundary.
		`|no(?:inspection|lint)\b` +
		`|#nosec\b` +
		`|NOSONAR\b` +
		`|sys(?:nb)?\b` +
		`)`)

// rxShebangComment matches a shebang like `//usr/bin/env go run`.
var rxShebangComment = regexp.MustCompile(`^//[^ /].*\bbin/`)

// commentGroupLooksLikeCode reports whether the lines of a //-style comment
// group parse as Go statements with at least one non-trivial statement.
// A bare identifier path or label is treated as trivial, since prose like
// "// foo" or "// TODO: bar" parses but is not commented-out code.
func commentGroupLooksLikeCode(group *ast.CommentGroup) bool {
	src := "package p\nfunc _() {\n" + group.Text() + "}\n"
	// AllErrors avoids the parser's panic/recover bailout on too many errors,
	// which crashes under tinygo's Wasm target as it lacks recover support.
	file, err := parser.ParseFile(token.NewFileSet(), "", src, parser.SkipObjectResolution|parser.AllErrors)
	if err != nil {
		return false
	}
	fn, _ := file.Decls[0].(*ast.FuncDecl)
	if fn == nil || fn.Body == nil {
		return false
	}
	for _, stmt := range fn.Body.List {
		if !isTrivialStmt(stmt) {
			return true
		}
	}
	return false
}

func isTrivialStmt(stmt ast.Stmt) bool {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		return isIdentPath(s.X)
	case *ast.LabeledStmt:
		return isTrivialStmt(s.Stmt)
	case *ast.EmptyStmt:
		return true
	}
	return false
}

func isIdentPath(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return isIdentPath(e.X)
	}
	return false
}

func (f *fumpter) applyPre(c *astutil.Cursor) {
	f.splitLongLine(c)

	switch node := c.Node().(type) {
	case *ast.File:
		// Unwrap single-spec var groups before the joining below,
		// so an adjacent var line and var group merge in one pass.
		for _, decl := range node.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok {
				f.removeParens(decl)
			}
		}

		// Join contiguous lone var/const/import lines.
		// Abort if there are empty lines in between,
		// including a leading comment if it's a directive.
		newDecls := make([]ast.Decl, 0, len(node.Decls))
		for i := 0; i < len(node.Decls); {
			newDecls = append(newDecls, node.Decls[i])
			start, ok := node.Decls[i].(*ast.GenDecl)
			if !ok || isCgoImport(start) || containsAnyDirective(start.Doc) {
				i++
				continue
			}
			lastPos := start.Pos()
			merged := false
		contLoop:
			for i++; i < len(node.Decls); {
				cont, ok := node.Decls[i].(*ast.GenDecl)
				if !ok || cont.Tok != start.Tok || cont.Lparen != token.NoPos || isCgoImport(cont) {
					break
				}
				// Are there things between these two declarations? e.g. empty lines, comments, directives
				// If so, break the chain on empty lines and directives, continue below for comments.
				if f.Line(lastPos) < f.Line(cont.Pos())-1 {
					// break on empty line
					if cont.Doc == nil {
						break
					}
					// break on directive
					for i, comment := range cont.Doc.List {
						if f.Line(comment.Slash) != f.Line(lastPos)+1+i || rxCommentDirective.MatchString(strings.TrimPrefix(comment.Text, "//")) {
							break contLoop
						}
					}
					// continue below for comments
				}

				start.Specs = append(start.Specs, cont.Specs...)
				merged = true
				end := cont.End()
				if c := f.inlineComment(cont.End()); c != nil {
					// don't move an inline comment outside
					end = c.End()
				}
				// Point Rparen at the last content character, like a real
				// ')', so start.End() stays on the content's final line and
				// the empty-line separator below is idempotent in one pass.
				start.Rparen = end - 1
				lastPos = cont.Pos()
				i++
			}
			// Re-sort imports in the new group so the output is idempotent.
			// Set Lparen so ast.SortImports doesn't skip the merged decl.
			if merged && start.Tok == token.IMPORT {
				start.Lparen = start.TokPos + token.Pos(len("import"))
				ast.SortImports(f.fset, f.astFile)
			}
		}
		node.Decls = newDecls

		// Multiline top-level declarations should be separated by an
		// empty line.
		// Do this after the joining of lone declarations above,
		// as joining single-line declarations makes then multi-line.
		var lastMulti bool
		var lastEnd token.Pos
		for _, decl := range node.Decls {
			pos := decl.Pos()
			// Trailing inline comments on lastEnd's line belong to the
			// previous decl and extend its effective end.
			effectiveEnd := lastEnd
			lastEndLine := f.Line(lastEnd)
			for _, cg := range f.commentsBetween(lastEnd, pos) {
				if f.Line(cg.Pos()) != lastEndLine {
					pos = cg.Pos()
					break
				}
				effectiveEnd = cg.End()
			}

			// Note that we want End-1, as End is the character after the node.
			multi := f.Line(pos) < f.Line(decl.End()-1)
			// A func declaration which fits on a single source line may
			// still be printed across multiple lines: go/printer's funcBody
			// breaks the body onto its own lines once header+body exceeds
			// 100 bytes. Approximate that with the source byte length.
			if fn, _ := decl.(*ast.FuncDecl); fn != nil && !multi && fn.Body != nil &&
				f.Offset(fn.End())-f.Offset(fn.Pos()) > 100 {
				multi = true
			}
			if multi && lastMulti && f.Line(effectiveEnd)+1 == f.Line(pos) {
				f.addNewline(effectiveEnd)
			}

			lastMulti = multi
			lastEnd = decl.End()
		}

		// Comments aren't nodes, so they're not walked by default.
	groupLoop:
		for _, group := range node.Comments {
			for _, comment := range group.List {
				// Leave shebang lines like `//usr/bin/env go run` alone.
				if f.Line(comment.Slash) == 1 && rxShebangComment.MatchString(comment.Text) {
					continue groupLoop
				}
				if comment.Text == "//gofumpt:diagnose" || strings.HasPrefix(comment.Text, "//gofumpt:diagnose ") {
					slc := []string{
						"//gofumpt:diagnose",
						"version:",
						version.String(""),
						"flags:",
						"-lang=" + f.LangVersion,
						"-modpath=" + f.ModulePath,
					}
					if s := f.Extra.String(); s != "" {
						slc = append(slc, "-extra="+s)
					}
					comment.Text = strings.Join(slc, " ")
				}
				body := strings.TrimPrefix(comment.Text, "//")
				if body == comment.Text {
					// /*-style comment
					continue groupLoop
				}
				if rxCommentDirective.MatchString(body) {
					// this line is a directive
					continue groupLoop
				}
				r, _ := utf8.DecodeRuneInString(body)
				if !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsSpace(r) {
					// this line could be code like "//{"
					continue groupLoop
				}
			}
			if commentGroupLooksLikeCode(group) {
				continue groupLoop
			}
			// If none of the comment group's lines look like a
			// directive or code, add spaces, if needed.
			for _, comment := range group.List {
				body := strings.TrimPrefix(comment.Text, "//")
				r, _ := utf8.DecodeRuneInString(body)
				if !unicode.IsSpace(r) {
					comment.Text = "// " + body
				}
			}
		}

	case *ast.DeclStmt:
		decl, ok := node.Decl.(*ast.GenDecl)
		if !ok || decl.Tok != token.VAR || len(decl.Specs) != 1 {
			break // e.g. const name = "value"
		}
		spec := decl.Specs[0].(*ast.ValueSpec)
		if spec.Type != nil {
			break // e.g. var name Type
		}
		tok := token.ASSIGN
		names := make([]ast.Expr, len(spec.Names))
		for i, name := range spec.Names {
			names[i] = name
			if name.Name != "_" {
				tok = token.DEFINE
			}
		}
		c.Replace(&ast.AssignStmt{
			Lhs: names,
			Tok: tok,
			Rhs: spec.Values,
		})

	case *ast.GenDecl:
		if node.Tok == token.IMPORT && node.Lparen.IsValid() {
			f.joinStdImports(node)
		}

		// Single var declarations shouldn't use parentheses, unless
		// there's a comment on the grouped declaration.
		f.removeParens(node)

	case *ast.InterfaceType:
		if len(node.Methods.List) > 0 {
			method := node.Methods.List[0]
			removeToPos := method.Pos()
			if comments := f.commentsBetween(node.Interface, method.Pos()); len(comments) > 0 {
				// only remove leading line upto the first comment
				removeToPos = comments[0].Pos()
			}
			// remove leading lines if they exist
			f.removeLines(f.Line(node.Interface)+1, f.Line(removeToPos))
		}

	case *ast.BlockStmt:
		f.stmts(node.List)
		comments := f.commentsBetween(node.Lbrace, node.Rbrace)
		if len(node.List) == 0 && len(comments) == 0 {
			f.removeLinesBetween(node.Lbrace, node.Rbrace)
			break
		}

		var sign *ast.FuncType
		var cond ast.Expr
		switch parent := c.Parent().(type) {
		case *ast.FuncDecl:
			sign = parent.Type
		case *ast.FuncLit:
			sign = parent.Type
		case *ast.IfStmt:
			cond = parent.Cond
		case *ast.ForStmt:
			cond = parent.Cond
		}

		if len(node.List) > 1 && sign == nil {
			// only if we have a single statement, or if
			// it's a func body.
			break
		}
		var bodyPos, bodyEnd token.Pos

		if len(node.List) > 0 {
			bodyPos = node.List[0].Pos()
			bodyEnd = node.List[len(node.List)-1].End()
		}
		if len(comments) > 0 {
			if pos := comments[0].Pos(); !bodyPos.IsValid() || pos < bodyPos {
				bodyPos = pos
			}
			if pos := comments[len(comments)-1].End(); !bodyPos.IsValid() || pos > bodyEnd {
				bodyEnd = pos
			}
		}

		f.removeLinesBetween(bodyEnd, node.Rbrace)

		if cond != nil && f.Line(cond.Pos()) != f.Line(cond.End()) {
			// The body is preceded by a multi-line condition, so an
			// empty line can help readability.
			return
		}
		if sign != nil {
			endLine := f.Line(sign.End())

			if f.Line(sign.Pos()) != endLine {
				handleMultiLine := func(fl *ast.FieldList) {
					// Refuse to insert a newline before the closing token
					// if the list is empty or all in one line.
					if fl == nil || len(fl.List) == 0 {
						return
					}
					fieldOpeningLine := f.Line(fl.Opening)
					fieldClosingLine := f.Line(fl.Closing)
					if fieldOpeningLine == fieldClosingLine {
						return
					}

					lastFieldEnd := fl.List[len(fl.List)-1].End()
					lastFieldLine := f.Line(lastFieldEnd)
					isLastFieldOnFieldClosingLine := lastFieldLine == fieldClosingLine
					isLastFieldOnSigClosingLine := lastFieldLine == endLine

					var isLastCommentGrpOnFieldClosingLine, isLastCommentGrpOnSigClosingLine bool
					if comments := f.commentsBetween(lastFieldEnd, fl.Closing); len(comments) > 0 {
						lastCommentGrp := comments[len(comments)-1]
						lastCommentGrpLine := f.Line(lastCommentGrp.End())

						isLastCommentGrpOnFieldClosingLine = lastCommentGrpLine == fieldClosingLine
						isLastCommentGrpOnSigClosingLine = lastCommentGrpLine == endLine
					}

					// is there a comment grp/last field, field closing and sig closing on the same line?
					if (isLastFieldOnFieldClosingLine && isLastFieldOnSigClosingLine) ||
						(isLastCommentGrpOnFieldClosingLine && isLastCommentGrpOnSigClosingLine) {
						fl.Closing += 1
						f.addNewline(fl.Closing)
					}
				}
				handleMultiLine(sign.Params)
				if sign.Results != nil && len(sign.Results.List) > 0 {
					lastResultLine := f.Line(sign.Results.List[len(sign.Results.List)-1].End())
					isLastResultOnParamClosingLine := sign.Params != nil && lastResultLine == f.Line(sign.Params.Closing)
					if !isLastResultOnParamClosingLine {
						handleMultiLine(sign.Results)
					}
				}
			}
		}

		f.removeLinesBetween(node.Lbrace, bodyPos)

	case *ast.CaseClause:
		f.stmts(node.Body)
		openLine := f.Line(node.Case)
		closeLine := f.Line(node.Colon)
		if openLine == closeLine {
			// nothing to do
			break
		}
		if len(f.commentsBetween(node.Case, node.Colon)) > 0 {
			// don't move comments
			break
		}
		// check the length excluding the body
		nodeWithoutBody := &ast.CaseClause{
			Case:  node.Case,
			List:  node.List,
			Colon: node.Colon,
		}
		if f.printLength(nodeWithoutBody) > shortLineLimit {
			// too long to collapse
			break
		}
		f.removeLines(openLine, closeLine)

	case *ast.CommClause:
		f.stmts(node.Body)

	case *ast.FieldList:
		numFields := node.NumFields()
		comments := f.commentsBetween(node.Pos(), node.End())

		if numFields == 0 && len(comments) == 0 {
			// Empty field lists should not contain a newline.
			// Do not join the two lines if the first has an inline
			// comment, as that can result in broken formatting.
			openLine := f.Line(node.Pos())
			closeLine := f.Line(node.End())
			f.removeLines(openLine, closeLine)
		} else {
			// Remove lines before first comment/field and lines after last
			// comment/field
			var bodyPos, bodyEnd token.Pos
			if numFields > 0 {
				bodyPos = node.List[0].Pos()
				bodyEnd = node.List[len(node.List)-1].End()
			}
			if len(comments) > 0 {
				if pos := comments[0].Pos(); !bodyPos.IsValid() || pos < bodyPos {
					bodyPos = pos
				}
				if pos := comments[len(comments)-1].End(); !bodyPos.IsValid() || pos > bodyEnd {
					bodyEnd = pos
				}
			}
			f.removeLinesBetween(node.Pos(), bodyPos)
			f.removeLinesBetween(bodyEnd, node.End())
		}

		if !f.Extra.GroupParams {
			break
		}
		switch c.Parent().(type) {
		case *ast.FuncDecl, *ast.FuncType, *ast.InterfaceType:
			node.List = f.mergeAdjacentFields(node.List)
			c.Replace(node)
		case *ast.StructType:
			// Do not merge adjacent fields in structs.
		}

	case *ast.ParenExpr:
		// Unwrap any chain of redundant inner parens first,
		// since astutil.Apply does not walk replacement nodes.
		node.X = ast.Unparen(node.X)
		if f.canRemoveParens(node) {
			c.Replace(node.X)
		}

	case *ast.BasicLit:
		// Octal number literals were introduced in Go 1.13.
		if goversion.Compare(f.LangVersion, "go1.13") >= 0 {
			if node.Kind == token.INT && rxOctalInteger.MatchString(node.Value) {
				node.Value = "0o" + node.Value[1:]
				c.Replace(node)
			}
		}

	case *ast.AssignStmt:
		// Only remove lines between the assignment token and the right-hand side
		// for simple single-value assignments. Skip multi-value assignments and
		// binary expressions like long string concatenations, where a line break
		// after the assignment token can improve readability.
		if len(node.Rhs) == 1 {
			if _, ok := node.Rhs[0].(*ast.BinaryExpr); !ok {
				f.removeLines(f.Line(node.TokPos), f.Line(node.Rhs[0].Pos()))
			}
		}

	case *ast.ReturnStmt:
		if len(node.Results) > 0 {
			break
		}
		if !f.Extra.ClotheReturns {
			break
		}
		results := f.parentFuncTypes[len(f.parentFuncTypes)-1].Results
		if results.NumFields() == 0 {
			break
		}

		// The function has return values; let's clothe the return.
		node.Results = make([]ast.Expr, 0, results.NumFields())
	nameLoop:
		for _, result := range results.List {
			for _, ident := range result.Names {
				name := ident.Name
				if name == "_" { // we can't handle blank names just yet
					node.Results = nil
					break nameLoop
				}
				node.Results = append(node.Results, &ast.Ident{
					// Use the Pos of the return statement, to not interfere with comment placement.
					NamePos: node.Pos(),
					Name:    name,
				})
			}
		}
		if len(node.Results) > 0 {
			c.Replace(node)
		}
	}
}

func (f *fumpter) applyPost(c *astutil.Cursor) {
	switch node := c.Node().(type) {
	// Adding newlines to composite literals happens as a "post" step, so
	// that we can take into account whether "pre" steps added any newlines
	// that would affect us here.
	case *ast.CompositeLit:
		if len(node.Elts) == 0 {
			// doesn't have elements
			break
		}
		openLine := f.Line(node.Lbrace)
		closeLine := f.Line(node.Rbrace)
		if openLine == closeLine {
			// all in a single line
			break
		}

		newlineAroundElems := false
		newlineBetweenElems := false
		lastEnd := node.Lbrace
		lastLine := openLine
		for i, elem := range node.Elts {
			pos := elem.Pos()
			comments := f.commentsBetween(lastEnd, pos)
			if len(comments) > 0 {
				pos = comments[0].Pos()
			}
			if curLine := f.Line(pos); curLine > lastLine {
				if i == 0 {
					newlineAroundElems = true

					// remove leading lines if they exist
					f.removeLines(openLine+1, curLine)
				} else {
					newlineBetweenElems = true
				}
			}
			lastEnd = elem.End()
			lastLine = f.Line(lastEnd)
		}
		if closeLine > lastLine {
			newlineAroundElems = true
		}

		if newlineBetweenElems || newlineAroundElems {
			first := node.Elts[0]
			if openLine == f.Line(first.Pos()) {
				// We want the newline right after the brace.
				f.addNewline(node.Lbrace + 1)
				closeLine = f.Line(node.Rbrace)
			}
			last := node.Elts[len(node.Elts)-1]
			if closeLine == f.Line(last.End()) {
				// We want the newline right before the brace.
				f.addNewline(node.Rbrace)
			}
		}

		// If there's a newline between any consecutive elements, there
		// must be a newline between all composite literal elements.
		if !newlineBetweenElems {
			break
		}
		for i1, elem1 := range node.Elts {
			i2 := i1 + 1
			if i2 >= len(node.Elts) {
				break
			}
			elem2 := node.Elts[i2]
			// TODO: do we care about &{}?
			_, ok1 := elem1.(*ast.CompositeLit)
			_, ok2 := elem2.(*ast.CompositeLit)
			if !ok1 && !ok2 {
				continue
			}
			if f.Line(elem1.End()) == f.Line(elem2.Pos()) {
				f.addNewline(elem1.End())
			}
		}

	// In a multi-line call, if the opening parenthesis is at the end of a
	// line, the closing parenthesis should be at the start of a line.
	// See https://github.com/mvdan/gofumpt/issues/74.
	case *ast.CallExpr:
		if !f.Extra.BalanceCalls {
			break
		}
		if len(node.Args) == 0 {
			break
		}
		openLine := f.Line(node.Lparen)
		closeLine := f.Line(node.Rparen)
		if openLine == closeLine {
			break
		}
		firstLine := f.Line(node.Args[0].Pos())
		lastEnd := node.Args[len(node.Args)-1].End()
		if comment := f.inlineComment(lastEnd); comment != nil {
			lastEnd = comment.End()
		}
		lastLine := f.Line(lastEnd)
		openAtEOL := openLine != firstLine
		closeAtBOL := closeLine != lastLine
		if openAtEOL && !closeAtBOL {
			f.addNewline(node.Rparen)
		}
	}
}

func (f *fumpter) splitLongLine(c *astutil.Cursor) {
	if os.Getenv("GOFUMPT_SPLIT_LONG_LINES") != "on" {
		// By default, this feature is turned off.
		// Turn it on by setting GOFUMPT_SPLIT_LONG_LINES=on.
		return
	}
	node := c.Node()
	if node == nil {
		return
	}

	newlinePos := node.Pos()
	start := f.Position(node.Pos())
	end := f.Position(node.End())

	// If the node is already split in multiple lines, there's nothing to do.
	if start.Line != end.Line {
		return
	}

	// Only split at the start of the current node if it's part of a list.
	if _, ok := c.Parent().(*ast.BinaryExpr); ok {
		// Chains of binary expressions are considered lists, too.
	} else if c.Index() >= 0 {
		// For the rest of the nodes, we're in a list if c.Index() >= 0.
	} else {
		return
	}

	// Like in printLength, add an approximation of the indentation level.
	// Since any existing tabs were already counted as one column, multiply
	// the level by 7.
	startCol := start.Column + f.blockLevel*7
	endCol := end.Column + f.blockLevel*7

	// If this is a composite literal,
	// and we were going to insert a newline before the entire literal,
	// insert the newline before the first element instead.
	// Since we'll add a newline after the last element too,
	// this format is generally going to be nicer.
	if comp := isComposite(node); comp != nil && len(comp.Elts) > 0 {
		newlinePos = comp.Elts[0].Pos()
	}

	// If this is a function call,
	// and we were to add a newline before the first argument,
	// prefer adding the newline before the entire call.
	// End-of-line parentheses aren't very nice, as we don't put their
	// counterparts at the start of a line too.
	// We do this by using the average of the two starting positions.
	if call, _ := node.(*ast.CallExpr); call != nil && len(call.Args) > 0 {
		first := f.Position(call.Args[0].Pos())
		startCol += (first.Column - start.Column) / 2
	}

	// If the start position is too short, we definitely won't split the line.
	if startCol <= shortLineLimit {
		return
	}

	lineEnd := f.Position(f.lineEnd(start.Line))

	// firstLength and secondLength are the split line lengths, excluding
	// indentation.
	firstLength := start.Column - f.blockLevel
	if firstLength < 0 {
		panic("negative length")
	}
	secondLength := lineEnd.Column - start.Column
	if secondLength < 0 {
		panic("negative length")
	}

	// If the line ends past the long line limit,
	// and both splits are estimated to take at least minSplitFactor of the limit,
	// then split the line.
	minSplitLength := int(f.minSplitFactor * longLineLimit)
	if endCol > longLineLimit &&
		firstLength >= minSplitLength && secondLength >= minSplitLength {
		f.addNewline(newlinePos)
	}
}

// canRemoveParens reports whether the parentheses around node are definitely
// useless and can be safely removed without changing intent.
func (f *fumpter) canRemoveParens(node *ast.ParenExpr) bool {
	// Don't drop parens which contain comments,
	// as the printer may not place them well without the parens.
	if len(f.commentsBetween(node.Lparen, node.Rparen)) > 0 {
		return false
	}
	return !keepParens(node.X, true)
}

// keepParens reports whether the parentheses directly around expr should be
// kept: around binary, unary, and type expressions for readability and for
// conversions like `(<-chan T)(v)`, but only when outermost; and around an
// expression whose leftmost operand is a composite literal, whose brace would
// otherwise open an if, for, or switch body.
func keepParens(expr ast.Expr, outermost bool) bool {
	switch expr := expr.(type) {
	case *ast.CompositeLit:
		return true
	case *ast.CallExpr:
		return keepParens(expr.Fun, false)
	case *ast.SelectorExpr:
		return keepParens(expr.X, false)
	case *ast.IndexExpr:
		return keepParens(expr.X, false)
	case *ast.IndexListExpr:
		return keepParens(expr.X, false)
	case *ast.SliceExpr:
		return keepParens(expr.X, false)
	case *ast.TypeAssertExpr:
		return keepParens(expr.X, false)
	case *ast.BinaryExpr, *ast.UnaryExpr, *ast.StarExpr,
		*ast.ChanType, *ast.ArrayType, *ast.MapType,
		*ast.FuncType, *ast.InterfaceType, *ast.StructType:
		return outermost
	}
	return false
}

func isComposite(node ast.Node) *ast.CompositeLit {
	switch node := node.(type) {
	case *ast.CompositeLit:
		return node
	case *ast.UnaryExpr:
		return isComposite(node.X) // e.g. &T{}
	default:
		return nil
	}
}

func (f *fumpter) stmts(list []ast.Stmt) {
	for i, stmt := range list {
		ifs, ok := stmt.(*ast.IfStmt)
		if !ok || i < 1 {
			continue // not an if following another statement
		}
		as, ok := list[i-1].(*ast.AssignStmt)
		if !ok || (as.Tok != token.DEFINE && as.Tok != token.ASSIGN) ||
			!identEqual(as.Lhs[len(as.Lhs)-1], "err") {
			continue // not ", err :=" nor ", err ="
		}
		be, ok := ifs.Cond.(*ast.BinaryExpr)
		if !ok || ifs.Init != nil || ifs.Else != nil {
			continue // complex if
		}
		if be.Op != token.NEQ || !identEqual(be.X, "err") ||
			!identEqual(be.Y, "nil") {
			continue // not "err != nil"
		}
		f.removeLinesBetween(as.End(), ifs.Pos())
	}
}

func identEqual(expr ast.Expr, name string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == name
}

// isCgoImport returns true if the declaration is simply:
//
//	import "C"
//
// or the equivalent:
//
//	import `C`
//
// Note that parentheses do not affect the result.
func isCgoImport(decl *ast.GenDecl) bool {
	if decl.Tok != token.IMPORT || len(decl.Specs) != 1 {
		return false
	}
	spec := decl.Specs[0].(*ast.ImportSpec)
	v, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		panic(err) // should never error
	}
	return v == "C"
}

// joinStdImports ensures that all standard library imports are together and at
// the top of the imports list.
func (f *fumpter) joinStdImports(d *ast.GenDecl) {
	var std, other []ast.Spec
	firstGroup := true
	lastEnd := d.Pos()
	needsSort := false

	// If ModulePath is "foo/bar", we assume "foo/..." is not part of std.
	// Users shouldn't declare modules that may collide with std this way,
	// but historically some private codebases have done so.
	// This is a relatively harmless way to make gofumpt compatible with them,
	// as it changes nothing for the common external module paths.
	var modulePrefix string
	if f.ModulePath == "" {
		// Nothing to do.
	} else if i := strings.IndexByte(f.ModulePath, '/'); i != -1 {
		// ModulePath is "foo/bar", so we use "foo" as the prefix.
		modulePrefix = f.ModulePath[:i]
	} else {
		// ModulePath is "foo", so we use "foo" as the prefix.
		modulePrefix = f.ModulePath
	}

	for i, spec := range d.Specs {
		spec := spec.(*ast.ImportSpec)
		if coms := f.commentsBetween(lastEnd, spec.Pos()); len(coms) > 0 {
			lastEnd = coms[len(coms)-1].End()
		}
		if i > 0 && firstGroup && f.Line(spec.Pos()) > f.Line(lastEnd)+1 {
			firstGroup = false
		} else {
			// We're still in the first group, update lastEnd.
			lastEnd = spec.End()
		}

		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			panic(err) // should never error
		}
		periodIndex := strings.IndexByte(path, '.')
		slashIndex := strings.IndexByte(path, '/')
		switch {
		// Imports with a period in the first path element are third party.
		// Note that this includes "foo.com" and excludes "foo/bar.com/baz".
		case periodIndex > 0 && (slashIndex == -1 || periodIndex < slashIndex),

			// "test" and "example" are reserved as per golang.org/issue/37641.
			strings.HasPrefix(path, "test/"),
			strings.HasPrefix(path, "example/"),

			// See if we match modulePrefix; see its documentation above.
			// We match either exactly or with a slash suffix,
			// so that the prefix "foo" for "foo/..." does not match "foobar".
			path == modulePrefix || strings.HasPrefix(path, modulePrefix+"/"),

			// To be conservative, if an import has a name or an inline
			// comment, and isn't part of the top group, treat it as non-std.
			!firstGroup && (spec.Name != nil || spec.Comment != nil):
			other = append(other, spec)
			continue
		}

		// If we're moving this std import further up, reset its
		// position, to avoid breaking comments.
		if !firstGroup || len(other) > 0 {
			setPos(reflect.ValueOf(spec), d.Pos())
			needsSort = true
		}
		std = append(std, spec)
	}
	// Ensure there is an empty line between std imports and other imports.
	if len(std) > 0 && len(other) > 0 && f.Line(std[len(std)-1].End())+1 >= f.Line(other[0].Pos()) {
		// We add two newlines, as that's necessary in some edge cases.
		// For example, if the std and non-std imports were together and
		// without indentation, adding one newline isn't enough. Two
		// empty lines will be printed as one by go/printer, anyway.
		f.addNewline(other[0].Pos() - 1)
		f.addNewline(other[0].Pos())
	}
	// Finally, join the imports, keeping std at the top.
	d.Specs = append(std, other...)

	// If we moved any std imports to the first group, we need to sort them
	// again.
	if needsSort {
		ast.SortImports(f.fset, f.astFile)
	}
}

// mergeAdjacentFields returns fields with adjacent fields merged if possible.
func (f *fumpter) mergeAdjacentFields(fields []*ast.Field) []*ast.Field {
	// If there are less than two fields then there is nothing to merge.
	if len(fields) < 2 {
		return fields
	}

	// Otherwise, iterate over adjacent pairs of fields, merging if possible,
	// and mutating fields. Elements of fields may be mutated (if merged with
	// following fields), discarded (if merged with a preceding field), or left
	// unchanged.
	i := 0
	for j := 1; j < len(fields); j++ {
		if f.shouldMergeAdjacentFields(fields[i], fields[j]) {
			fields[i].Names = append(fields[i].Names, fields[j].Names...)
		} else {
			i++
			fields[i] = fields[j]
		}
	}
	return fields[:i+1]
}

func (f *fumpter) shouldMergeAdjacentFields(f1, f2 *ast.Field) bool {
	if len(f1.Names) == 0 || len(f2.Names) == 0 {
		// Both must have names for the merge to work.
		return false
	}
	if f.Line(f1.Pos()) != f.Line(f2.Pos()) {
		// Trust the user if they used separate lines.
		return false
	}

	// Only merge if the types that the syntax nodes represent are equal,
	// e.g. two *ast.Ident nodes "int" are equal, but the two *ast.Ident nodes
	// "string" and "bool" are not. We use reflection to quickly discard most cases.
	//
	// We use an empty [token.FileSet] so that positions are ignored when printing,
	// and two syntax nodes with different uses of newlines end up the same.
	//
	// Note that we could in theory use go/types here, but in practice gofumpt
	// needs to be fast, hence it shouldn't rely on expensive typechecking.
	if reflect.TypeOf(f1.Type) != reflect.TypeOf(f2.Type) {
		return false
	}
	emptyFset := token.NewFileSet()
	var b1, b2 bytes.Buffer
	if err := format.Node(&b1, emptyFset, f1.Type); err != nil {
		return false
	}
	if err := format.Node(&b2, emptyFset, f2.Type); err != nil {
		return false
	}
	return bytes.Equal(b1.Bytes(), b2.Bytes())
}

var posType = reflect.TypeFor[token.Pos]()

// setPos recursively sets all position fields in the node v to pos.
func setPos(v reflect.Value, pos token.Pos) {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if !v.IsValid() {
		return
	}
	if v.Type() == posType {
		v.Set(reflect.ValueOf(pos))
	}
	if v.Kind() == reflect.Struct {
		for i := range v.NumField() {
			setPos(v.Field(i), pos)
		}
	}
}

func containsAnyDirective(group *ast.CommentGroup) bool {
	if group == nil {
		return false
	}
	for _, comment := range group.List {
		body := strings.TrimPrefix(comment.Text, "//")
		if rxCommentDirective.MatchString(body) {
			return true
		}
	}
	return false
}
