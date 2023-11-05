package jsonparser

import (
	bio "bytes"
)

// minInt64 '-9223372036854775808' is the smallest representable number in int64
const minInt64 = `9223372036854775808`

// About 2x faster then strconv.ParseInt because it only supports base 10, which is enough for JSON
func parseInt(bytes []byte) (v int64, ok bool, overflow bool) {
	if len(bytes) == 0 {
		return 0, false, false
	}

	var neg bool = false
	if bytes[0] == '-' {
		neg = true
		bytes = bytes[1:]
	}

	var b int64 = 0
	for _, c := range bytes {
		if c >= '0' && c <= '9' {
			b = (10 * v) + int64(c-'0')
		} else {
			return 0, false, false
		}
		if overflow = (b < v); overflow {
			break
		}
		v = b
	}

	if overflow {
		if neg && bio.Equal(bytes, []byte(minInt64)) {
			return b, true, false
		}
		return 0, false, true
	}

	if neg {
		return -v, true, false
	} else {
		return v, true, false
	}
}
