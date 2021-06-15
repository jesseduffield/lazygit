package generator

import (
	"strings"
)

// Returns is a slice of Return.
type Returns []Return

// Return is the result of a method's invocation.
type Return struct {
	Name string
	Type string
}

// HasLength is true if there are returns, else false.
func (r Returns) HasLength() bool {
	return len(r) > 0
}

// WithPrefix builds a string representing the parameters returned from a
// function, and adds a prefix to each.
func (r Returns) WithPrefix(p string) string {
	if len(r) == 0 {
		return ""
	}

	rets := []string{}
	for i := range r {
		if p == "" {
			rets = append(rets, unexport(r[i].Name))
		} else {
			rets = append(rets, p+unexport(r[i].Name))
		}
	}
	return strings.Join(rets, ", ")
}

// AsArgs builds a string representing the arguments passed to a function.
func (r Returns) AsArgs() string {
	if len(r) == 0 {
		return ""
	}

	rets := []string{}
	for i := range r {
		rets = append(rets, r[i].Type)
	}
	return strings.Join(rets, ", ")
}

// AsNamedArgsWithTypes builds a string representing a function's named
// arguments, with associated types.
func (r Returns) AsNamedArgsWithTypes() string {
	if len(r) == 0 {
		return ""
	}

	rets := []string{}
	for i := range r {
		rets = append(rets, unexport(r[i].Name)+" "+r[i].Type)
	}
	return strings.Join(rets, ", ")
}

// AsNamedArgs builds a string representing a function's named arguments.
func (r Returns) AsNamedArgs() string {
	if len(r) == 0 {
		return ""
	}

	rets := []string{}
	for i := range r {
		rets = append(rets, unexport(r[i].Name))
	}
	return strings.Join(rets, ", ")
}

// AsReturnSignature builds a string representing signature for the returns of
// a function.
func (r Returns) AsReturnSignature() string {
	if len(r) == 0 {
		return ""
	}
	if len(r) == 1 {
		return r[0].Type
	}
	result := "("
	for i := range r {
		result = result + r[i].Type
		if i < len(r) {
			result = result + ", "
		}
	}
	result = result + ")"
	return result
}
