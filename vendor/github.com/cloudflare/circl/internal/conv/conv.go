package conv

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/cryptobyte"
)

// BytesLe2Hex returns an hexadecimal string of a number stored in a
// little-endian order slice x.
func BytesLe2Hex(x []byte) string {
	b := &strings.Builder{}
	b.Grow(2*len(x) + 2)
	fmt.Fprint(b, "0x")
	if len(x) == 0 {
		fmt.Fprint(b, "00")
	}
	for i := len(x) - 1; i >= 0; i-- {
		fmt.Fprintf(b, "%02x", x[i])
	}
	return b.String()
}

// BytesLe2BigInt converts a little-endian slice x into a big-endian
// math/big.Int.
func BytesLe2BigInt(x []byte) *big.Int {
	n := len(x)
	b := new(big.Int)
	if len(x) > 0 {
		y := make([]byte, n)
		for i := 0; i < n; i++ {
			y[n-1-i] = x[i]
		}
		b.SetBytes(y)
	}
	return b
}

// BytesBe2Uint64Le converts a big-endian slice x to a little-endian slice of uint64.
func BytesBe2Uint64Le(x []byte) []uint64 {
	l := len(x)
	z := make([]uint64, (l+7)/8)
	blocks := l / 8
	for i := 0; i < blocks; i++ {
		z[i] = binary.BigEndian.Uint64(x[l-8*(i+1):])
	}
	remBytes := l % 8
	for i := 0; i < remBytes; i++ {
		z[blocks] |= uint64(x[l-1-8*blocks-i]) << uint(8*i)
	}
	return z
}

// BigInt2BytesLe stores a positive big.Int number x into a little-endian slice z.
// The slice is modified if the bitlength of x <= 8*len(z) (padding with zeros).
// If x does not fit in the slice or is negative, z is not modified.
func BigInt2BytesLe(z []byte, x *big.Int) {
	xLen := (x.BitLen() + 7) >> 3
	zLen := len(z)
	if zLen >= xLen && x.Sign() >= 0 {
		y := x.Bytes()
		for i := 0; i < xLen; i++ {
			z[i] = y[xLen-1-i]
		}
		for i := xLen; i < zLen; i++ {
			z[i] = 0
		}
	}
}

// Uint64Le2BigInt converts a little-endian slice x into a big number.
func Uint64Le2BigInt(x []uint64) *big.Int {
	n := len(x)
	b := new(big.Int)
	var bi big.Int
	for i := n - 1; i >= 0; i-- {
		bi.SetUint64(x[i])
		b.Lsh(b, 64)
		b.Add(b, &bi)
	}
	return b
}

// Uint64Le2BytesLe converts a little-endian slice x to a little-endian slice of bytes.
func Uint64Le2BytesLe(x []uint64) []byte {
	b := make([]byte, 8*len(x))
	n := len(x)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint64(b[i*8:], x[i])
	}
	return b
}

// Uint64Le2BytesBe converts a little-endian slice x to a big-endian slice of bytes.
func Uint64Le2BytesBe(x []uint64) []byte {
	b := make([]byte, 8*len(x))
	n := len(x)
	for i := 0; i < n; i++ {
		binary.BigEndian.PutUint64(b[i*8:], x[n-1-i])
	}
	return b
}

// Uint64Le2Hex returns an hexadecimal string of a number stored in a
// little-endian order slice x.
func Uint64Le2Hex(x []uint64) string {
	b := new(strings.Builder)
	b.Grow(16*len(x) + 2)
	fmt.Fprint(b, "0x")
	if len(x) == 0 {
		fmt.Fprint(b, "00")
	}
	for i := len(x) - 1; i >= 0; i-- {
		fmt.Fprintf(b, "%016x", x[i])
	}
	return b.String()
}

// BigInt2Uint64Le stores a positive big.Int number x into a little-endian slice z.
// The slice is modified if the bitlength of x <= 8*len(z) (padding with zeros).
// If x does not fit in the slice or is negative, z is not modified.
func BigInt2Uint64Le(z []uint64, x *big.Int) {
	xLen := (x.BitLen() + 63) >> 6 // number of 64-bit words
	zLen := len(z)
	if zLen >= xLen && x.Sign() > 0 {
		var y, yi big.Int
		y.Set(x)
		two64 := big.NewInt(1)
		two64.Lsh(two64, 64).Sub(two64, big.NewInt(1))
		for i := 0; i < xLen; i++ {
			yi.And(&y, two64)
			z[i] = yi.Uint64()
			y.Rsh(&y, 64)
		}
	}
	for i := xLen; i < zLen; i++ {
		z[i] = 0
	}
}

// MarshalBinary encodes a value into a byte array in a format readable by UnmarshalBinary.
func MarshalBinary(v cryptobyte.MarshalingValue) ([]byte, error) {
	const DefaultSize = 32
	b := cryptobyte.NewBuilder(make([]byte, 0, DefaultSize))
	b.AddValue(v)
	return b.Bytes()
}

// MarshalBinaryLen encodes a value into an array of n bytes in a format readable by UnmarshalBinary.
func MarshalBinaryLen(v cryptobyte.MarshalingValue, length uint) ([]byte, error) {
	b := cryptobyte.NewFixedBuilder(make([]byte, 0, length))
	b.AddValue(v)
	return b.Bytes()
}

// A UnmarshalingValue decodes itself from a cryptobyte.String and advances the pointer.
// It reports whether the read was successful.
type UnmarshalingValue interface {
	Unmarshal(*cryptobyte.String) bool
}

// UnmarshalBinary recovers a value from a byte array.
// It returns an error if the read was unsuccessful.
func UnmarshalBinary(v UnmarshalingValue, data []byte) (err error) {
	s := cryptobyte.String(data)
	if data == nil || !v.Unmarshal(&s) || !s.Empty() {
		err = fmt.Errorf("cannot read %T from input string", v)
	}
	return
}
