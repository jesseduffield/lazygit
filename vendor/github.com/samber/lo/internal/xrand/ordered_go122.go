//go:build go1.22

package xrand

import "math/rand/v2"

// Shuffle returns a slice of shuffled values. Uses the Fisher-Yates shuffle algorithm.
func Shuffle(n int, swap func(i, j int)) {
	rand.Shuffle(n, swap)
}

// IntN returns, as an int, a pseudo-random number in the half-open interval [0,n)
// from the default Source.
// It panics if n <= 0.
func IntN(n int) int {
	return rand.IntN(n)
}

// Int64 returns a non-negative pseudo-random 63-bit integer as an int64
// from the default Source.
func Int64() int64 {
	return rand.Int64()
}
