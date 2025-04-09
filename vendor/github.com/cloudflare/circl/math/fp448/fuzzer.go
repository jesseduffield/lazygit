//go:build gofuzz
// +build gofuzz

// How to run the fuzzer:
//
//	$ go get -u github.com/dvyukov/go-fuzz/go-fuzz
//	$ go get -u github.com/dvyukov/go-fuzz/go-fuzz-build
//	$ go-fuzz-build -libfuzzer -func FuzzReduction -o lib.a
//	$ clang -fsanitize=fuzzer lib.a -o fu.exe
//	$ ./fu.exe
package fp448

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/cloudflare/circl/internal/conv"
)

// FuzzReduction is a fuzzer target for red64 function, which reduces t
// (112 bits) to a number t' (56 bits) congruent modulo p448.
func FuzzReduction(data []byte) int {
	if len(data) != 2*Size {
		return -1
	}
	var got, want Elt
	var lo, hi [7]uint64
	a := data[:Size]
	b := data[Size:]
	lo[0] = binary.LittleEndian.Uint64(a[0*8 : 1*8])
	lo[1] = binary.LittleEndian.Uint64(a[1*8 : 2*8])
	lo[2] = binary.LittleEndian.Uint64(a[2*8 : 3*8])
	lo[3] = binary.LittleEndian.Uint64(a[3*8 : 4*8])
	lo[4] = binary.LittleEndian.Uint64(a[4*8 : 5*8])
	lo[5] = binary.LittleEndian.Uint64(a[5*8 : 6*8])
	lo[6] = binary.LittleEndian.Uint64(a[6*8 : 7*8])

	hi[0] = binary.LittleEndian.Uint64(b[0*8 : 1*8])
	hi[1] = binary.LittleEndian.Uint64(b[1*8 : 2*8])
	hi[2] = binary.LittleEndian.Uint64(b[2*8 : 3*8])
	hi[3] = binary.LittleEndian.Uint64(b[3*8 : 4*8])
	hi[4] = binary.LittleEndian.Uint64(b[4*8 : 5*8])
	hi[5] = binary.LittleEndian.Uint64(b[5*8 : 6*8])
	hi[6] = binary.LittleEndian.Uint64(b[6*8 : 7*8])

	red64(&got, &lo, &hi)

	t := conv.BytesLe2BigInt(data[:2*Size])

	two448 := big.NewInt(1)
	two448.Lsh(two448, 448) // 2^448
	mask448 := big.NewInt(1)
	mask448.Sub(two448, mask448) // 2^448-1
	two224plus1 := big.NewInt(1)
	two224plus1.Lsh(two224plus1, 224)
	two224plus1.Add(two224plus1, big.NewInt(1)) // 2^224+1

	var loBig, hiBig big.Int
	for t.Cmp(two448) >= 0 {
		loBig.And(t, mask448)
		hiBig.Rsh(t, 448)
		t.Mul(&hiBig, two224plus1)
		t.Add(t, &loBig)
	}
	conv.BigInt2BytesLe(want[:], t)

	if got != want {
		fmt.Printf("in:   %v\n", conv.BytesLe2BigInt(data[:2*Size]))
		fmt.Printf("got:  %v\n", got)
		fmt.Printf("want: %v\n", want)
		panic("error found")
	}
	return 1
}
