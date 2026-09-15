//go:build go1.24

package deadlock

import (
	"unsafe"
	"weak"
)

type beforeAfter struct {
	before weak.Pointer[byte]
	after  weak.Pointer[byte]
}

// ptrFromInterface extracts the data pointer from an interface{} value.
// An interface (eface) is {type *_type, data unsafe.Pointer}; we grab the second word.
func ptrFromInterface(i interface{}) *byte {
	type eface struct {
		_    uintptr
		data unsafe.Pointer
	}
	return (*byte)((*eface)(unsafe.Pointer(&i)).data)
}

func newBeforeAfter(before, after interface{}) beforeAfter {
	return beforeAfter{
		before: weak.Make(ptrFromInterface(before)),
		after:  weak.Make(ptrFromInterface(after)),
	}
}
