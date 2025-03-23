package jsonschema

import (
	"fmt"
	"io/fs"
	gopath "path"
	"path/filepath"
	"reflect"
	"strings"

	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
)

type commentOptions struct {
	fullObjectText bool // use the first sentence only?
}

// CommentOption allows for special configuration options when preparing Go
// source files for comment extraction.
type CommentOption func(*commentOptions)

// WithFullComment will configure the comment extraction to process to use an
// object type's full comment text instead of just the synopsis.
func WithFullComment() CommentOption {
	return func(o *commentOptions) {
		o.fullObjectText = true
	}
}

// AddGoComments will update the reflectors comment map with all the comments
// found in the provided source directories including sub-directories, in order to
// generate a dictionary of comments associated with Types and Fields. The results
// will be added to the `Reflect.CommentMap` ready to use with Schema "description"
// fields.
//
// The `go/parser` library is used to extract all the comments and unfortunately doesn't
// have a built-in way to determine the fully qualified name of a package. The `base`
// parameter, the URL used to import that package, is thus required to be able to match
// reflected types.
//
// When parsing type comments, by default we use the `go/doc`'s Synopsis method to extract
// the first phrase only. Field comments, which tend to be much shorter, will include everything.
// This behavior can be changed by using the `WithFullComment` option.
func (r *Reflector) AddGoComments(base, path string, opts ...CommentOption) error {
	if r.CommentMap == nil {
		r.CommentMap = make(map[string]string)
	}
	co := new(commentOptions)
	for _, opt := range opts {
		opt(co)
	}

	return r.extractGoComments(base, path, r.CommentMap, co)
}

func (r *Reflector) extractGoComments(base, path string, commentMap map[string]string, opts *commentOptions) error {
	fset := token.NewFileSet()
	dict := make(map[string][]*ast.Package)
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			for _, v := range d {
				// paths may have multiple packages, like for tests
				k := gopath.Join(base, path)
				dict[k] = append(dict[k], v)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for pkg, p := range dict {
		for _, f := range p {
			gtxt := ""
			typ := ""
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					typ = x.Name.String()
					if !ast.IsExported(typ) {
						typ = ""
					} else {
						txt := x.Doc.Text()
						if txt == "" && gtxt != "" {
							txt = gtxt
							gtxt = ""
						}
						if !opts.fullObjectText {
							txt = doc.Synopsis(txt)
						}
						commentMap[fmt.Sprintf("%s.%s", pkg, typ)] = strings.TrimSpace(txt)
					}
				case *ast.Field:
					txt := x.Doc.Text()
					if txt == "" {
						txt = x.Comment.Text()
					}
					if typ != "" && txt != "" {
						for _, n := range x.Names {
							if ast.IsExported(n.String()) {
								k := fmt.Sprintf("%s.%s.%s", pkg, typ, n)
								commentMap[k] = strings.TrimSpace(txt)
							}
						}
					}
				case *ast.GenDecl:
					// remember for the next type
					gtxt = x.Doc.Text()
				}
				return true
			})
		}
	}

	return nil
}

func (r *Reflector) lookupComment(t reflect.Type, name string) string {
	if r.LookupComment != nil {
		if comment := r.LookupComment(t, name); comment != "" {
			return comment
		}
	}

	if r.CommentMap == nil {
		return ""
	}

	n := fullyQualifiedTypeName(t)
	if name != "" {
		n = n + "." + name
	}

	return r.CommentMap[n]
}
