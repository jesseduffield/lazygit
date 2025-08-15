// Package math provides some utility functions for big integers.
package math

import "math/big"

// SignedDigit obtains the signed-digit recoding of n and returns a list L of
// digits such that n = sum( L[i]*2^(i*(w-1)) ), and each L[i] is an odd number
// in the set {±1, ±3, ..., ±2^(w-1)-1}. The third parameter ensures that the
// output has ceil(l/(w-1)) digits.
//
// Restrictions:
//   - n is odd and n > 0.
//   - 1 < w < 32.
//   - l >= bit length of n.
//
// References:
//   - Alg.6 in "Exponent Recoding and Regular Exponentiation Algorithms"
//     by Joye-Tunstall. http://doi.org/10.1007/978-3-642-02384-2_21
//   - Alg.6 in "Selecting Elliptic Curves for Cryptography: An Efficiency and
//     Security Analysis" by Bos et al. http://doi.org/10.1007/s13389-015-0097-y
func SignedDigit(n *big.Int, w, l uint) []int32 {
	if n.Sign() <= 0 || n.Bit(0) == 0 {
		panic("n must be non-zero, odd, and positive")
	}
	if w <= 1 || w >= 32 {
		panic("Verify that 1 < w < 32")
	}
	if uint(n.BitLen()) > l {
		panic("n is too big to fit in l digits")
	}
	lenN := (l + (w - 1) - 1) / (w - 1) // ceil(l/(w-1))
	L := make([]int32, lenN+1)
	var k, v big.Int
	k.Set(n)

	var i uint
	for i = 0; i < lenN; i++ {
		words := k.Bits()
		value := int32(words[0] & ((1 << w) - 1))
		value -= int32(1) << (w - 1)
		L[i] = value
		v.SetInt64(int64(value))
		k.Sub(&k, &v)
		k.Rsh(&k, w-1)
	}
	L[i] = int32(k.Int64())
	return L
}

// OmegaNAF obtains the window-w Non-Adjacent Form of a positive number n and
// 1 < w < 32. The returned slice L holds n = sum( L[i]*2^i ).
//
// Reference:
//   - Alg.9 "Efficient arithmetic on Koblitz curves" by Solinas.
//     http://doi.org/10.1023/A:1008306223194
func OmegaNAF(n *big.Int, w uint) (L []int32) {
	if n.Sign() < 0 {
		panic("n must be positive")
	}
	if w <= 1 || w >= 32 {
		panic("Verify that 1 < w < 32")
	}

	L = make([]int32, n.BitLen()+1)
	var k, v big.Int
	k.Set(n)

	i := 0
	for ; k.Sign() > 0; i++ {
		value := int32(0)
		if k.Bit(0) == 1 {
			words := k.Bits()
			value = int32(words[0] & ((1 << w) - 1))
			if value >= (int32(1) << (w - 1)) {
				value -= int32(1) << w
			}
			v.SetInt64(int64(value))
			k.Sub(&k, &v)
		}
		L[i] = value
		k.Rsh(&k, 1)
	}
	return L[:i]
}
