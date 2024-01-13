// +build appengine js

package util

// BytesToReadOnlyString returns a string converted from given bytes.
func BytesToReadOnlyString(b []byte) string {
	return string(b)
}

// StringToReadOnlyBytes returns bytes converted from given string.
func StringToReadOnlyBytes(s string) []byte {
	return []byte(s)
}
