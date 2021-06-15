package generator

import (
	"fmt"
	"go/types"
	"strings"

	"golang.org/x/tools/go/types/typeutil"
)

func (f *Fake) addTypesForMethod(sig *types.Signature) {
	for i := 0; i < sig.Results().Len(); i++ {
		ret := sig.Results().At(i)
		f.addImportsFor(ret.Type())
	}
	for i := 0; i < sig.Params().Len(); i++ {
		param := sig.Params().At(i)
		f.addImportsFor(param.Type())
	}
}

func methodForSignature(sig *types.Signature, methodName string, imports Imports) Method {
	params := []Param{}
	for i := 0; i < sig.Params().Len(); i++ {
		param := sig.Params().At(i)
		isVariadic := i == sig.Params().Len()-1 && sig.Variadic()
		typ := types.TypeString(param.Type(), imports.AliasForPackage)
		if isVariadic {
			typ = "..." + typ[2:] // Change []string to ...string
		}
		p := Param{
			Name:       fmt.Sprintf("arg%v", i+1),
			Type:       typ,
			IsVariadic: isVariadic,
			IsSlice:    strings.HasPrefix(typ, "[]"),
		}
		params = append(params, p)
	}
	returns := []Return{}
	for i := 0; i < sig.Results().Len(); i++ {
		ret := sig.Results().At(i)
		r := Return{
			Name: fmt.Sprintf("result%v", i+1),
			Type: types.TypeString(ret.Type(), imports.AliasForPackage),
		}
		returns = append(returns, r)
	}
	return Method{
		Name:    methodName,
		Returns: returns,
		Params:  params,
	}
}

// interfaceMethodSet identifies the methods that are exported for a given
// interface.
func interfaceMethodSet(t types.Type) []*rawMethod {
	if t == nil {
		return nil
	}
	var result []*rawMethod
	methods := typeutil.IntuitiveMethodSet(t, nil)
	for i := range methods {
		if methods[i].Obj() == nil || methods[i].Type() == nil {
			continue
		}
		fun, ok := methods[i].Obj().(*types.Func)
		if !ok {
			continue
		}
		sig, ok := methods[i].Type().(*types.Signature)
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

func (f *Fake) loadMethods() {
	var methods []*rawMethod
	if f.Mode == Package {
		methods = packageMethodSet(f.Package)
	} else {
		if !f.IsInterface() || f.Target == nil || f.Target.Type() == nil {
			return
		}
		methods = interfaceMethodSet(f.Target.Type())
	}

	for i := range methods {
		f.addTypesForMethod(methods[i].Signature)
	}

	for i := range methods {
		method := methodForSignature(methods[i].Signature, methods[i].Func.Name(), f.Imports)
		f.Methods = append(f.Methods, method)
	}
}
