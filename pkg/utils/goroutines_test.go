package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// If a goroutine panics without restoring the terminal, the panic is printed
// to a terminal that is still in raw mode and on the alternate screen, and the
// terminal needs a reset afterwards. Go has no process-wide hook for panics,
// so every goroutine has to restore the terminal itself. This test checks that
// they all do:
//
//   - a go statement runs utils.Safe (gocui can't import utils, so it uses its
//     own safe function)
//   - the function given to time.AfterFunc is wrapped with utils.SafeFunc
//   - the function given to a Go method, such as errgroup.Group.Go or
//     sync.WaitGroup.Go, is wrapped with utils.SafeFunc or
//     utils.SafeFuncWithError
func TestGoroutinesAreGuardedAgainstPanics(t *testing.T) {
	root := filepath.Join("..", "..")
	fset := token.NewFileSet()

	paths := []string{"main.go"}
	for _, dir := range []string{"pkg", "cmd"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				relPath, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				paths = append(paths, relPath)
			}
			return nil
		})
		assert.NoError(t, err)
	}

	for _, path := range paths {
		src, err := os.ReadFile(filepath.Join(root, path))
		if !assert.NoError(t, err) {
			continue
		}
		file, err := parser.ParseFile(fset, path, src, 0)
		if !assert.NoError(t, err) {
			continue
		}

		// Inside package utils, its functions are called without the package name
		qualify := func(name string) string {
			if file.Name.Name == "utils" {
				return name
			}
			return "utils." + name
		}
		safe := qualify("Safe")
		if file.Name.Name == "gocui" {
			safe = "safe"
		}
		safeFunc := qualify("SafeFunc")
		safeFuncWithError := qualify("SafeFuncWithError")

		ast.Inspect(file, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.GoStmt:
				if !isCallTo(node.Call, safe) {
					t.Errorf("%s: go statement doesn't run %s", fset.Position(node.Pos()), safe)
				}
			case *ast.CallExpr:
				if types.ExprString(node.Fun) == "time.AfterFunc" && len(node.Args) == 2 &&
					!isCallTo(node.Args[1], safeFunc) {
					t.Errorf("%s: function given to time.AfterFunc isn't wrapped with %s",
						fset.Position(node.Pos()), safeFunc)
				}
				if selector, ok := node.Fun.(*ast.SelectorExpr); ok && selector.Sel.Name == "Go" && len(node.Args) == 1 &&
					!isCallTo(node.Args[0], safeFunc, safeFuncWithError) {
					t.Errorf("%s: function given to %s isn't wrapped with %s or %s",
						fset.Position(node.Pos()), types.ExprString(node.Fun), safeFunc, safeFuncWithError)
				}
			}
			return true
		})
	}
}

// Whether expr is a call to one of the given functions, which are spelled
// like "utils.Safe".
func isCallTo(expr ast.Expr, functionNames ...string) bool {
	call, ok := expr.(*ast.CallExpr)
	return ok && lo.Contains(functionNames, types.ExprString(call.Fun))
}
