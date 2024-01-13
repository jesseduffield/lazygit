// +build !appengine,!js

package util

import (
	"reflect"
	"unsafe"
)

// BytesToReadOnlyString returns a string converted from given bytes.
func BytesToReadOnlyString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToReadOnlyBytes returns bytes converted from given string.
func StringToReadOnlyBytes(s string) (bs []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}
