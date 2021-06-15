package generator

import (
	"go/types"

	"golang.org/x/tools/go/packages"
)

type rawMethod struct {
	Func      *types.Func
	Signature *types.Signature
}

// packageMethodSet identifies the functions that are exported from a given
// package.
func packageMethodSet(p *packages.Package) []*rawMethod {
	if p == nil || p.Types == nil || p.Types.Scope() == nil {
		return nil
	}
	var result []*rawMethod
	scope := p.Types.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if !obj.Exported() {
			continue // skip unexported names
		}
		fun, ok := obj.(*types.Func)
		if !ok {
			continue
		}
		sig, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}
		result = append(result, &rawMethod{
			Func:      fun,
			Signature: sig,
		})
	}

	return result
}
