package litter

import (
	"reflect"
)

// deInterface returns values inside of non-nil interfaces when possible.
// This is useful for data types like structs, arrays, slices, and maps which
// can contain varying types packed inside an interface.
func deInterface(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func isPointerValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	}
	return false
}

func isZeroValue(v reflect.Value) bool {
	return (isPointerValue(v) && v.IsNil()) ||
		(v.IsValid() && v.CanInterface() && reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()))
}
