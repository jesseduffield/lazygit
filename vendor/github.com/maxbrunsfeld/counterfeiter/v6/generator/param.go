package generator

import "strings"

// Params is a slice of Param.
type Params []Param

// Param is an argument to a function.
type Param struct {
	Name       string
	Type       string
	IsVariadic bool
	IsSlice    bool
}

// Slices returns those params that are a slice.
func (p Params) Slices() Params {
	var result Params
	for i := range p {
		if p[i].IsSlice {
			result = append(result, p[i])
		}
	}
	return result
}

// HasLength returns true if there are params. It returns false if there are no
// params.
func (p Params) HasLength() bool {
	return len(p) > 0
}

// WithPrefix builds a string representing a functions parameters, and adds a
// prefix to each.
func (p Params) WithPrefix(prefix string) string {
	if len(p) == 0 {
		return ""
	}

	params := []string{}
	for i := range p {
		if prefix == "" {
			params = append(params, unexport(p[i].Name))
		} else {
			params = append(params, prefix+unexport(p[i].Name))
		}
	}
	return strings.Join(params, ", ")
}

// AsArgs builds a string that represents the parameters to a function as
// arguments to a function invocation.
func (p Params) AsArgs() string {
	if len(p) == 0 {
		return ""
	}

	params := []string{}
	for i := range p {
		params = append(params, p[i].Type)
	}
	return strings.Join(params, ", ")
}

// AsNamedArgsWithTypes builds a string that represents parameters as named
// arugments to a function, with associated types.
func (p Params) AsNamedArgsWithTypes() string {
	if len(p) == 0 {
		return ""
	}

	params := []string{}
	for i := range p {
		params = append(params, unexport(p[i].Name)+" "+p[i].Type)
	}
	return strings.Join(params, ", ")
}

// AsNamedArgs builds a string that represents parameters as named arguments.
func (p Params) AsNamedArgs() string {
	if len(p) == 0 {
		return ""
	}

	params := []string{}
	for i := range p {
		if p[i].IsSlice {
			params = append(params, unexport(p[i].Name)+"Copy")
		} else {
			params = append(params, unexport(p[i].Name))
		}
	}
	return strings.Join(params, ", ")
}

// AsNamedArgsForInvocation builds a string that represents a function's
// arguments as required for invocation of the function.
func (p Params) AsNamedArgsForInvocation() string {
	if len(p) == 0 {
		return ""
	}

	params := []string{}
	for i := range p {
		if p[i].IsVariadic {
			params = append(params, unexport(p[i].Name)+"...")
		} else {
			params = append(params, unexport(p[i].Name))
		}
	}
	return strings.Join(params, ", ")
}

// AsReturnSignature builds a string representing signature for the params of
// a function.
func (p Params) AsReturnSignature() string {
	if len(p) == 0 {
		return ""
	}
	if len(p) == 1 {
		if p[0].IsVariadic {
			return strings.Replace(p[0].Type, "...", "[]", -1)
		}
		return p[0].Type
	}
	result := "("
	for i := range p {
		t := p[i].Type
		if p[i].IsVariadic {
			t = strings.Replace(t, "...", "[]", -1)
		}
		result = result + t
		if i < len(p) {
			result = result + ", "
		}
	}
	result = result + ")"
	return result
}
