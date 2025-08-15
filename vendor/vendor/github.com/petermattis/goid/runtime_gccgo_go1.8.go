//go:build gccgo && go1.8
// +build gccgo,go1.8

package goid

// https://github.com/gcc-mirror/gcc/blob/releases/gcc-7/libgo/go/runtime/runtime2.go#L329-L354

type g struct {
	_panic       uintptr
	_defer       uintptr
	m            uintptr
	syscallsp    uintptr
	syscallpc    uintptr
	param        uintptr
	atomicstatus uint32
	goid         int64 // Here it is!
}
