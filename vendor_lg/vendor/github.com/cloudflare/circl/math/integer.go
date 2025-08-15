package math

import "math/bits"

// NextPow2 finds the next power of two (N=2^k, k>=0) greater than n.
// If n is already a power of two, then this function returns n, and log2(n).
func NextPow2(n uint) (N uint, k uint) {
	if bits.OnesCount(n) == 1 {
		k = uint(bits.TrailingZeros(n))
		N = n
	} else {
		k = uint(bits.Len(n))
		N = uint(1) << k
	}
	return
}
