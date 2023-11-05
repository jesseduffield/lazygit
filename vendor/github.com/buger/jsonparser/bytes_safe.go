// +build appengine appenginevm

package jsonparser

import (
	"strconv"
)

// See fastbytes_unsafe.go for explanation on why *[]byte is used (signatures must be consistent with those in that file)

func equalStr(b *[]byte, s string) bool {
	return string(*b) == s
}

func parseFloat(b *[]byte) (float64, error) {
	return strconv.ParseFloat(string(*b), 64)
}

func bytesToString(b *[]byte) string {
	return string(*b)
}

func StringToBytes(s string) []byte {
	return []byte(s)
}
