package lo

import (
	"time"
)

// Duration returns the time taken to execute a function.
// Play: https://go.dev/play/p/HQfbBbAXaFP
func Duration(callback func()) time.Duration {
	return Duration0(callback)
}

// Duration0 returns the time taken to execute a function.
// Play: https://go.dev/play/p/HQfbBbAXaFP
func Duration0(callback func()) time.Duration {
	start := time.Now()
	callback()
	return time.Since(start)
}

// Duration1 returns the time taken to execute a function.
// Play: https://go.dev/play/p/HQfbBbAXaFP
func Duration1[A any](callback func() A) (A, time.Duration) {
	start := time.Now()
	a := callback()
	return a, time.Since(start)
}

// Duration2 returns the time taken to execute a function.
// Play: https://go.dev/play/p/HQfbBbAXaFP
func Duration2[A, B any](callback func() (A, B)) (A, B, time.Duration) {
	start := time.Now()
	a, b := callback()
	return a, b, time.Since(start)
}

// Duration3 returns the time taken to execute a function.
// Play: https://go.dev/play/p/xr863iwkAxQ
func Duration3[A, B, C any](callback func() (A, B, C)) (A, B, C, time.Duration) {
	start := time.Now()
	a, b, c := callback()
	return a, b, c, time.Since(start)
}

// Duration4 returns the time taken to execute a function.
// Play: https://go.dev/play/p/xr863iwkAxQ
func Duration4[A, B, C, D any](callback func() (A, B, C, D)) (A, B, C, D, time.Duration) {
	start := time.Now()
	a, b, c, d := callback()
	return a, b, c, d, time.Since(start)
}

// Duration5 returns the time taken to execute a function.
// Play: https://go.dev/play/p/xr863iwkAxQ
func Duration5[A, B, C, D, E any](callback func() (A, B, C, D, E)) (A, B, C, D, E, time.Duration) {
	start := time.Now()
	a, b, c, d, e := callback()
	return a, b, c, d, e, time.Since(start)
}

// Duration6 returns the time taken to execute a function.
// Play: https://go.dev/play/p/mR4bTQKO-Tf
func Duration6[A, B, C, D, E, F any](callback func() (A, B, C, D, E, F)) (A, B, C, D, E, F, time.Duration) {
	start := time.Now()
	a, b, c, d, e, f := callback()
	return a, b, c, d, e, f, time.Since(start)
}

// Duration7 returns the time taken to execute a function.
// Play: https://go.dev/play/p/jgIAcBWWInS
func Duration7[A, B, C, D, E, F, G any](callback func() (A, B, C, D, E, F, G)) (A, B, C, D, E, F, G, time.Duration) {
	start := time.Now()
	a, b, c, d, e, f, g := callback()
	return a, b, c, d, e, f, g, time.Since(start)
}

// Duration8 returns the time taken to execute a function.
// Play: https://go.dev/play/p/T8kxpG1c5Na
func Duration8[A, B, C, D, E, F, G, H any](callback func() (A, B, C, D, E, F, G, H)) (A, B, C, D, E, F, G, H, time.Duration) {
	start := time.Now()
	a, b, c, d, e, f, g, h := callback()
	return a, b, c, d, e, f, g, h, time.Since(start)
}

// Duration9 returns the time taken to execute a function.
// Play: https://go.dev/play/p/bg9ix2VrZ0j
func Duration9[A, B, C, D, E, F, G, H, I any](callback func() (A, B, C, D, E, F, G, H, I)) (A, B, C, D, E, F, G, H, I, time.Duration) {
	start := time.Now()
	a, b, c, d, e, f, g, h, i := callback()
	return a, b, c, d, e, f, g, h, i, time.Since(start)
}

// Duration10 returns the time taken to execute a function.
// Play: https://go.dev/play/p/Y3n7oJXqJbk
func Duration10[A, B, C, D, E, F, G, H, I, J any](callback func() (A, B, C, D, E, F, G, H, I, J)) (A, B, C, D, E, F, G, H, I, J, time.Duration) {
	start := time.Now()
	a, b, c, d, e, f, g, h, i, j := callback()
	return a, b, c, d, e, f, g, h, i, j, time.Since(start)
}
