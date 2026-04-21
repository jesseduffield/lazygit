package lo

// T2 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T2[A, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{A: a, B: b}
}

// T3 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T3[A, B, C any](a A, b B, c C) Tuple3[A, B, C] {
	return Tuple3[A, B, C]{A: a, B: b, C: c}
}

// T4 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T4[A, B, C, D any](a A, b B, c C, d D) Tuple4[A, B, C, D] {
	return Tuple4[A, B, C, D]{A: a, B: b, C: c, D: d}
}

// T5 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T5[A, B, C, D, E any](a A, b B, c C, d D, e E) Tuple5[A, B, C, D, E] {
	return Tuple5[A, B, C, D, E]{A: a, B: b, C: c, D: d, E: e}
}

// T6 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T6[A, B, C, D, E, F any](a A, b B, c C, d D, e E, f F) Tuple6[A, B, C, D, E, F] {
	return Tuple6[A, B, C, D, E, F]{A: a, B: b, C: c, D: d, E: e, F: f}
}

// T7 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T7[A, B, C, D, E, F, G any](a A, b B, c C, d D, e E, f F, g G) Tuple7[A, B, C, D, E, F, G] {
	return Tuple7[A, B, C, D, E, F, G]{A: a, B: b, C: c, D: d, E: e, F: f, G: g}
}

// T8 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T8[A, B, C, D, E, F, G, H any](a A, b B, c C, d D, e E, f F, g G, h H) Tuple8[A, B, C, D, E, F, G, H] {
	return Tuple8[A, B, C, D, E, F, G, H]{A: a, B: b, C: c, D: d, E: e, F: f, G: g, H: h}
}

// T9 creates a tuple from a list of values.
// Play: https://go.dev/play/p/IllL3ZO4BQm
func T9[A, B, C, D, E, F, G, H, I any](a A, b B, c C, d D, e E, f F, g G, h H, i I) Tuple9[A, B, C, D, E, F, G, H, I] {
	return Tuple9[A, B, C, D, E, F, G, H, I]{A: a, B: b, C: c, D: d, E: e, F: f, G: g, H: h, I: i}
}

// Unpack2 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack2[A, B any](tuple Tuple2[A, B]) (A, B) {
	return tuple.A, tuple.B
}

// Unpack3 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack3[A, B, C any](tuple Tuple3[A, B, C]) (A, B, C) {
	return tuple.A, tuple.B, tuple.C
}

// Unpack4 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack4[A, B, C, D any](tuple Tuple4[A, B, C, D]) (A, B, C, D) {
	return tuple.A, tuple.B, tuple.C, tuple.D
}

// Unpack5 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack5[A, B, C, D, E any](tuple Tuple5[A, B, C, D, E]) (A, B, C, D, E) {
	return tuple.A, tuple.B, tuple.C, tuple.D, tuple.E
}

// Unpack6 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack6[A, B, C, D, E, F any](tuple Tuple6[A, B, C, D, E, F]) (A, B, C, D, E, F) {
	return tuple.A, tuple.B, tuple.C, tuple.D, tuple.E, tuple.F
}

// Unpack7 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack7[A, B, C, D, E, F, G any](tuple Tuple7[A, B, C, D, E, F, G]) (A, B, C, D, E, F, G) {
	return tuple.A, tuple.B, tuple.C, tuple.D, tuple.E, tuple.F, tuple.G
}

// Unpack8 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack8[A, B, C, D, E, F, G, H any](tuple Tuple8[A, B, C, D, E, F, G, H]) (A, B, C, D, E, F, G, H) {
	return tuple.A, tuple.B, tuple.C, tuple.D, tuple.E, tuple.F, tuple.G, tuple.H
}

// Unpack9 returns values contained in a tuple.
// Play: https://go.dev/play/p/xVP_k0kJ96W
func Unpack9[A, B, C, D, E, F, G, H, I any](tuple Tuple9[A, B, C, D, E, F, G, H, I]) (A, B, C, D, E, F, G, H, I) {
	return tuple.A, tuple.B, tuple.C, tuple.D, tuple.E, tuple.F, tuple.G, tuple.H, tuple.I
}

// Zip2 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip2[A, B any](a []A, b []B) []Tuple2[A, B] {
	size := uint(Max([]int{len(a), len(b)}))

	result := make([]Tuple2[A, B], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
	}

	return result
}

// Zip3 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip3[A, B, C any](a []A, b []B, c []C) []Tuple3[A, B, C] {
	size := uint(Max([]int{len(a), len(b), len(c)}))

	result := make([]Tuple3[A, B, C], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
		result[index].C = NthOrEmpty(c, index)
	}

	return result
}

// Zip4 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip4[A, B, C, D any](a []A, b []B, c []C, d []D) []Tuple4[A, B, C, D] {
	size := uint(Max([]int{len(a), len(b), len(c), len(d)}))

	result := make([]Tuple4[A, B, C, D], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
		result[index].C = NthOrEmpty(c, index)
		result[index].D = NthOrEmpty(d, index)
	}

	return result
}

// Zip5 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip5[A, B, C, D, E any](a []A, b []B, c []C, d []D, e []E) []Tuple5[A, B, C, D, E] {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e)}))

	result := make([]Tuple5[A, B, C, D, E], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
		result[index].C = NthOrEmpty(c, index)
		result[index].D = NthOrEmpty(d, index)
		result[index].E = NthOrEmpty(e, index)
	}

	return result
}

// Zip6 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip6[A, B, C, D, E, F any](a []A, b []B, c []C, d []D, e []E, f []F) []Tuple6[A, B, C, D, E, F] {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f)}))

	result := make([]Tuple6[A, B, C, D, E, F], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
		result[index].C = NthOrEmpty(c, index)
		result[index].D = NthOrEmpty(d, index)
		result[index].E = NthOrEmpty(e, index)
		result[index].F = NthOrEmpty(f, index)
	}

	return result
}

// Zip7 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip7[A, B, C, D, E, F, G any](a []A, b []B, c []C, d []D, e []E, f []F, g []G) []Tuple7[A, B, C, D, E, F, G] {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g)}))

	result := make([]Tuple7[A, B, C, D, E, F, G], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
		result[index].C = NthOrEmpty(c, index)
		result[index].D = NthOrEmpty(d, index)
		result[index].E = NthOrEmpty(e, index)
		result[index].F = NthOrEmpty(f, index)
		result[index].G = NthOrEmpty(g, index)
	}

	return result
}

// Zip8 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip8[A, B, C, D, E, F, G, H any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H) []Tuple8[A, B, C, D, E, F, G, H] {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g), len(h)}))

	result := make([]Tuple8[A, B, C, D, E, F, G, H], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
		result[index].C = NthOrEmpty(c, index)
		result[index].D = NthOrEmpty(d, index)
		result[index].E = NthOrEmpty(e, index)
		result[index].F = NthOrEmpty(f, index)
		result[index].G = NthOrEmpty(g, index)
		result[index].H = NthOrEmpty(h, index)
	}

	return result
}

// Zip9 creates a slice of grouped elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/jujaA6GaJTp
func Zip9[A, B, C, D, E, F, G, H, I any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, i []I) []Tuple9[A, B, C, D, E, F, G, H, I] {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g), len(h), len(i)}))

	result := make([]Tuple9[A, B, C, D, E, F, G, H, I], size)

	for index := uint(0); index < size; index++ {
		result[index].A = NthOrEmpty(a, index)
		result[index].B = NthOrEmpty(b, index)
		result[index].C = NthOrEmpty(c, index)
		result[index].D = NthOrEmpty(d, index)
		result[index].E = NthOrEmpty(e, index)
		result[index].F = NthOrEmpty(f, index)
		result[index].G = NthOrEmpty(g, index)
		result[index].H = NthOrEmpty(h, index)
		result[index].I = NthOrEmpty(i, index)
	}

	return result
}

// ZipBy2 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/wlHur6yO8rR
func ZipBy2[A, B, Out any](a []A, b []B, iteratee func(a A, b B) Out) []Out {
	size := uint(Max([]int{len(a), len(b)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
		)
	}

	return result
}

// ZipBy3 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/j9maveOnSQX
func ZipBy3[A, B, C, Out any](a []A, b []B, c []C, iteratee func(a A, b B, c C) Out) []Out {
	size := uint(Max([]int{len(a), len(b), len(c)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
		)
	}

	return result
}

// ZipBy4 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/Y1eF2Ke0Ayz
func ZipBy4[A, B, C, D, Out any](a []A, b []B, c []C, d []D, iteratee func(a A, b B, c C, d D) Out) []Out {
	size := uint(Max([]int{len(a), len(b), len(c), len(d)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
		)
	}

	return result
}

// ZipBy5 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/SLynyalh5Oa
func ZipBy5[A, B, C, D, E, Out any](a []A, b []B, c []C, d []D, e []E, iteratee func(a A, b B, c C, d D, e E) Out) []Out {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
		)
	}

	return result
}

// ZipBy6 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/IK6KVgw9e-S
func ZipBy6[A, B, C, D, E, F, Out any](a []A, b []B, c []C, d []D, e []E, f []F, iteratee func(a A, b B, c C, d D, e E, f F) Out) []Out {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
		)
	}

	return result
}

// ZipBy7 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/4uW6a2vXh8w
func ZipBy7[A, B, C, D, E, F, G, Out any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, iteratee func(a A, b B, c C, d D, e E, f F, g G) Out) []Out {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
			NthOrEmpty(g, index),
		)
	}

	return result
}

// ZipBy8 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/tk8xW7XzY4v
func ZipBy8[A, B, C, D, E, F, G, H, Out any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, iteratee func(a A, b B, c C, d D, e E, f F, g G, h H) Out) []Out {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g), len(h)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
			NthOrEmpty(g, index),
			NthOrEmpty(h, index),
		)
	}

	return result
}

// ZipBy9 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// Play: https://go.dev/play/p/VGqjDmQ9YqX
func ZipBy9[A, B, C, D, E, F, G, H, I, Out any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, i []I, iteratee func(a A, b B, c C, d D, e E, f F, g G, h H, i I) Out) []Out {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g), len(h), len(i)}))

	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		result[index] = iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
			NthOrEmpty(g, index),
			NthOrEmpty(h, index),
			NthOrEmpty(i, index),
		)
	}

	return result
}

// ZipByErr2 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr2[A, B, Out any](a []A, b []B, iteratee func(a A, b B) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// ZipByErr3 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr3[A, B, C, Out any](a []A, b []B, c []C, iteratee func(a A, b B, c C) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b), len(c)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// ZipByErr4 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr4[A, B, C, D, Out any](a []A, b []B, c []C, d []D, iteratee func(a A, b B, c C, d D) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b), len(c), len(d)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// ZipByErr5 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr5[A, B, C, D, E, Out any](a []A, b []B, c []C, d []D, e []E, iteratee func(a A, b B, c C, d D, e E) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// ZipByErr6 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr6[A, B, C, D, E, F, Out any](a []A, b []B, c []C, d []D, e []E, f []F, iteratee func(a A, b B, c C, d D, e E, f F) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// ZipByErr7 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr7[A, B, C, D, E, F, G, Out any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, iteratee func(a A, b B, c C, d D, e E, f F, g G) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
			NthOrEmpty(g, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// ZipByErr8 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr8[A, B, C, D, E, F, G, H, Out any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, iteratee func(a A, b B, c C, d D, e E, f F, g G, h H) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g), len(h)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
			NthOrEmpty(g, index),
			NthOrEmpty(h, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// ZipByErr9 creates a slice of transformed elements, the first of which contains the first elements
// of the given slices, the second of which contains the second elements of the given slices, and so on.
// When collections are different sizes, the Tuple attributes are filled with zero value.
// It returns the first error returned by the iteratee.
func ZipByErr9[A, B, C, D, E, F, G, H, I, Out any](a []A, b []B, c []C, d []D, e []E, f []F, g []G, h []H, i []I, iteratee func(a A, b B, c C, d D, e E, f F, g G, h H, i I) (Out, error)) ([]Out, error) {
	size := uint(Max([]int{len(a), len(b), len(c), len(d), len(e), len(f), len(g), len(h), len(i)}))
	result := make([]Out, size)

	for index := uint(0); index < size; index++ {
		r, err := iteratee(
			NthOrEmpty(a, index),
			NthOrEmpty(b, index),
			NthOrEmpty(c, index),
			NthOrEmpty(d, index),
			NthOrEmpty(e, index),
			NthOrEmpty(f, index),
			NthOrEmpty(g, index),
			NthOrEmpty(h, index),
			NthOrEmpty(i, index),
		)
		if err != nil {
			return nil, err
		}
		result[index] = r
	}

	return result, nil
}

// Unzip2 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip2[A, B any](tuples []Tuple2[A, B]) ([]A, []B) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
	}

	return r1, r2
}

// Unzip3 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip3[A, B, C any](tuples []Tuple3[A, B, C]) ([]A, []B, []C) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
	}

	return r1, r2, r3
}

// Unzip4 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip4[A, B, C, D any](tuples []Tuple4[A, B, C, D]) ([]A, []B, []C, []D) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
		r4 = append(r4, tuples[i].D)
	}

	return r1, r2, r3, r4
}

// Unzip5 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip5[A, B, C, D, E any](tuples []Tuple5[A, B, C, D, E]) ([]A, []B, []C, []D, []E) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
		r4 = append(r4, tuples[i].D)
		r5 = append(r5, tuples[i].E)
	}

	return r1, r2, r3, r4, r5
}

// Unzip6 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip6[A, B, C, D, E, F any](tuples []Tuple6[A, B, C, D, E, F]) ([]A, []B, []C, []D, []E, []F) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
		r4 = append(r4, tuples[i].D)
		r5 = append(r5, tuples[i].E)
		r6 = append(r6, tuples[i].F)
	}

	return r1, r2, r3, r4, r5, r6
}

// Unzip7 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip7[A, B, C, D, E, F, G any](tuples []Tuple7[A, B, C, D, E, F, G]) ([]A, []B, []C, []D, []E, []F, []G) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
		r4 = append(r4, tuples[i].D)
		r5 = append(r5, tuples[i].E)
		r6 = append(r6, tuples[i].F)
		r7 = append(r7, tuples[i].G)
	}

	return r1, r2, r3, r4, r5, r6, r7
}

// Unzip8 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip8[A, B, C, D, E, F, G, H any](tuples []Tuple8[A, B, C, D, E, F, G, H]) ([]A, []B, []C, []D, []E, []F, []G, []H) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)
	r8 := make([]H, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
		r4 = append(r4, tuples[i].D)
		r5 = append(r5, tuples[i].E)
		r6 = append(r6, tuples[i].F)
		r7 = append(r7, tuples[i].G)
		r8 = append(r8, tuples[i].H)
	}

	return r1, r2, r3, r4, r5, r6, r7, r8
}

// Unzip9 accepts a slice of grouped elements and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/ciHugugvaAW
func Unzip9[A, B, C, D, E, F, G, H, I any](tuples []Tuple9[A, B, C, D, E, F, G, H, I]) ([]A, []B, []C, []D, []E, []F, []G, []H, []I) {
	size := len(tuples)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)
	r8 := make([]H, 0, size)
	r9 := make([]I, 0, size)

	for i := range tuples {
		r1 = append(r1, tuples[i].A)
		r2 = append(r2, tuples[i].B)
		r3 = append(r3, tuples[i].C)
		r4 = append(r4, tuples[i].D)
		r5 = append(r5, tuples[i].E)
		r6 = append(r6, tuples[i].F)
		r7 = append(r7, tuples[i].G)
		r8 = append(r8, tuples[i].H)
		r9 = append(r9, tuples[i].I)
	}

	return r1, r2, r3, r4, r5, r6, r7, r8, r9
}

// UnzipBy2 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/tN8yqaRZz0r
func UnzipBy2[In, A, B any](items []In, iteratee func(In) (a A, b B)) ([]A, []B) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)

	for i := range items {
		a, b := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
	}

	return r1, r2
}

// UnzipBy3 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/36ITO2DlQq1
func UnzipBy3[In, A, B, C any](items []In, iteratee func(In) (a A, b B, c C)) ([]A, []B, []C) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)

	for i := range items {
		a, b, c := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
	}

	return r1, r2, r3
}

// UnzipBy4 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/zJ6qY1dD1rL
func UnzipBy4[In, A, B, C, D any](items []In, iteratee func(In) (a A, b B, c C, d D)) ([]A, []B, []C, []D) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)

	for i := range items {
		a, b, c, d := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
	}

	return r1, r2, r3, r4
}

// UnzipBy5 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/3f7jKkV9xZt
func UnzipBy5[In, A, B, C, D, E any](items []In, iteratee func(In) (a A, b B, c C, d D, e E)) ([]A, []B, []C, []D, []E) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)

	for i := range items {
		a, b, c, d, e := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
	}

	return r1, r2, r3, r4, r5
}

// UnzipBy6 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/8Y1b7tKu2pL
func UnzipBy6[In, A, B, C, D, E, F any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F)) ([]A, []B, []C, []D, []E, []F) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)

	for i := range items {
		a, b, c, d, e, f := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
	}

	return r1, r2, r3, r4, r5, r6
}

// UnzipBy7 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/7j1kLmVn3pM
func UnzipBy7[In, A, B, C, D, E, F, G any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F, g G)) ([]A, []B, []C, []D, []E, []F, []G) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)

	for i := range items {
		a, b, c, d, e, f, g := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
		r7 = append(r7, g)
	}

	return r1, r2, r3, r4, r5, r6, r7
}

// UnzipBy8 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/1n2k3L4m5N6
func UnzipBy8[In, A, B, C, D, E, F, G, H any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F, g G, h H)) ([]A, []B, []C, []D, []E, []F, []G, []H) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)
	r8 := make([]H, 0, size)

	for i := range items {
		a, b, c, d, e, f, g, h := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
		r7 = append(r7, g)
		r8 = append(r8, h)
	}

	return r1, r2, r3, r4, r5, r6, r7, r8
}

// UnzipBy9 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// Play: https://go.dev/play/p/7o8p9q0r1s2
func UnzipBy9[In, A, B, C, D, E, F, G, H, I any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F, g G, h H, i I)) ([]A, []B, []C, []D, []E, []F, []G, []H, []I) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)
	r8 := make([]H, 0, size)
	r9 := make([]I, 0, size)

	for i := range items {
		a, b, c, d, e, f, g, h, i := iteratee(items[i])
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
		r7 = append(r7, g)
		r8 = append(r8, h)
		r9 = append(r9, i)
	}

	return r1, r2, r3, r4, r5, r6, r7, r8, r9
}

// UnzipByErr2 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr2[In, A, B any](items []In, iteratee func(In) (a A, b B, err error)) ([]A, []B, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)

	for i := range items {
		a, b, err := iteratee(items[i])
		if err != nil {
			return nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
	}

	return r1, r2, nil
}

// UnzipByErr3 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr3[In, A, B, C any](items []In, iteratee func(In) (a A, b B, c C, err error)) ([]A, []B, []C, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)

	for i := range items {
		a, b, c, err := iteratee(items[i])
		if err != nil {
			return nil, nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
	}

	return r1, r2, r3, nil
}

// UnzipByErr4 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr4[In, A, B, C, D any](items []In, iteratee func(In) (a A, b B, c C, d D, err error)) ([]A, []B, []C, []D, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)

	for i := range items {
		a, b, c, d, err := iteratee(items[i])
		if err != nil {
			return nil, nil, nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
	}

	return r1, r2, r3, r4, nil
}

// UnzipByErr5 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr5[In, A, B, C, D, E any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, err error)) ([]A, []B, []C, []D, []E, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)

	for i := range items {
		a, b, c, d, e, err := iteratee(items[i])
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
	}

	return r1, r2, r3, r4, r5, nil
}

// UnzipByErr6 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr6[In, A, B, C, D, E, F any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F, err error)) ([]A, []B, []C, []D, []E, []F, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)

	for i := range items {
		a, b, c, d, e, f, err := iteratee(items[i])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
	}

	return r1, r2, r3, r4, r5, r6, nil
}

// UnzipByErr7 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr7[In, A, B, C, D, E, F, G any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F, g G, err error)) ([]A, []B, []C, []D, []E, []F, []G, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)

	for i := range items {
		a, b, c, d, e, f, g, err := iteratee(items[i])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
		r7 = append(r7, g)
	}

	return r1, r2, r3, r4, r5, r6, r7, nil
}

// UnzipByErr8 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr8[In, A, B, C, D, E, F, G, H any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F, g G, h H, err error)) ([]A, []B, []C, []D, []E, []F, []G, []H, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)
	r8 := make([]H, 0, size)

	for i := range items {
		a, b, c, d, e, f, g, h, err := iteratee(items[i])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
		r7 = append(r7, g)
		r8 = append(r8, h)
	}

	return r1, r2, r3, r4, r5, r6, r7, r8, nil
}

// UnzipByErr9 iterates over a collection and creates a slice regrouping the elements
// to their pre-zip configuration.
// It returns the first error returned by the iteratee.
func UnzipByErr9[In, A, B, C, D, E, F, G, H, I any](items []In, iteratee func(In) (a A, b B, c C, d D, e E, f F, g G, h H, i I, err error)) ([]A, []B, []C, []D, []E, []F, []G, []H, []I, error) {
	size := len(items)
	r1 := make([]A, 0, size)
	r2 := make([]B, 0, size)
	r3 := make([]C, 0, size)
	r4 := make([]D, 0, size)
	r5 := make([]E, 0, size)
	r6 := make([]F, 0, size)
	r7 := make([]G, 0, size)
	r8 := make([]H, 0, size)
	r9 := make([]I, 0, size)

	for i := range items {
		a, b, c, d, e, f, g, h, i, err := iteratee(items[i])
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
		}
		r1 = append(r1, a)
		r2 = append(r2, b)
		r3 = append(r3, c)
		r4 = append(r4, d)
		r5 = append(r5, e)
		r6 = append(r6, f)
		r7 = append(r7, g)
		r8 = append(r8, h)
		r9 = append(r9, i)
	}

	return r1, r2, r3, r4, r5, r6, r7, r8, r9, nil
}

// CrossJoin2 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/3VFppyL9FDU
func CrossJoin2[A, B any](listA []A, listB []B) []Tuple2[A, B] {
	return CrossJoinBy2(listA, listB, T2[A, B])
}

// CrossJoin3 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/2WGeHyJj4fK
func CrossJoin3[A, B, C any](listA []A, listB []B, listC []C) []Tuple3[A, B, C] {
	return CrossJoinBy3(listA, listB, listC, T3[A, B, C])
}

// CrossJoin4 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/6XhKjLmMnNp
func CrossJoin4[A, B, C, D any](listA []A, listB []B, listC []C, listD []D) []Tuple4[A, B, C, D] {
	return CrossJoinBy4(listA, listB, listC, listD, T4[A, B, C, D])
}

// CrossJoin5 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/7oPqRsTuVwX
func CrossJoin5[A, B, C, D, E any](listA []A, listB []B, listC []C, listD []D, listE []E) []Tuple5[A, B, C, D, E] {
	return CrossJoinBy5(listA, listB, listC, listD, listE, T5[A, B, C, D, E])
}

// CrossJoin6 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/8yZ1aB2cD3e
func CrossJoin6[A, B, C, D, E, F any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F) []Tuple6[A, B, C, D, E, F] {
	return CrossJoinBy6(listA, listB, listC, listD, listE, listF, T6[A, B, C, D, E, F])
}

// CrossJoin7 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/9f4g5h6i7j8
func CrossJoin7[A, B, C, D, E, F, G any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G) []Tuple7[A, B, C, D, E, F, G] {
	return CrossJoinBy7(listA, listB, listC, listD, listE, listF, listG, T7[A, B, C, D, E, F, G])
}

// CrossJoin8 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/0k1l2m3n4o5
func CrossJoin8[A, B, C, D, E, F, G, H any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, listH []H) []Tuple8[A, B, C, D, E, F, G, H] {
	return CrossJoinBy8(listA, listB, listC, listD, listE, listF, listG, listH, T8[A, B, C, D, E, F, G, H])
}

// CrossJoin9 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/6p7q8r9s0t1
func CrossJoin9[A, B, C, D, E, F, G, H, I any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, listH []H, listI []I) []Tuple9[A, B, C, D, E, F, G, H, I] {
	return CrossJoinBy9(listA, listB, listC, listD, listE, listF, listG, listH, listI, T9[A, B, C, D, E, F, G, H, I])
}

// CrossJoinBy2 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/8Y7btpvuA-C
func CrossJoinBy2[A, B, Out any](listA []A, listB []B, transform func(a A, b B) Out) []Out {
	size := len(listA) * len(listB)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			result = append(result, transform(a, b))
		}
	}

	return result
}

// CrossJoinBy3 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/3z4y5x6w7v8
func CrossJoinBy3[A, B, C, Out any](listA []A, listB []B, listC []C, transform func(a A, b B, c C) Out) []Out {
	size := len(listA) * len(listB) * len(listC)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				result = append(result, transform(a, b, c))
			}
		}
	}

	return result
}

// CrossJoinBy4 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/8b9c0d1e2f3
func CrossJoinBy4[A, B, C, D, Out any](listA []A, listB []B, listC []C, listD []D, transform func(a A, b B, c C, d D) Out) []Out {
	size := len(listA) * len(listB) * len(listC) * len(listD)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					result = append(result, transform(a, b, c, d))
				}
			}
		}
	}

	return result
}

// CrossJoinBy5 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/4g5h6i7j8k9
func CrossJoinBy5[A, B, C, D, E, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, transform func(a A, b B, c C, d D, e E) Out) []Out {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						result = append(result, transform(a, b, c, d, e))
					}
				}
			}
		}
	}

	return result
}

// CrossJoinBy6 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/1l2m3n4o5p6
func CrossJoinBy6[A, B, C, D, E, F, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, transform func(a A, b B, c C, d D, e E, f F) Out) []Out {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							result = append(result, transform(a, b, c, d, e, f))
						}
					}
				}
			}
		}
	}

	return result
}

// CrossJoinBy7 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/7q8r9s0t1u2
func CrossJoinBy7[A, B, C, D, E, F, G, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, transform func(a A, b B, c C, d D, e E, f F, g G) Out) []Out {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF) * len(listG)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							for _, g := range listG {
								result = append(result, transform(a, b, c, d, e, f, g))
							}
						}
					}
				}
			}
		}
	}

	return result
}

// CrossJoinBy8 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/3v4w5x6y7z8
func CrossJoinBy8[A, B, C, D, E, F, G, H, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, listH []H, transform func(a A, b B, c C, d D, e E, f F, g G, h H) Out) []Out {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF) * len(listG) * len(listH)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							for _, g := range listG {
								for _, h := range listH {
									result = append(result, transform(a, b, c, d, e, f, g, h))
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

// CrossJoinBy9 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// Play: https://go.dev/play/p/9a0b1c2d3e4
func CrossJoinBy9[A, B, C, D, E, F, G, H, I, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, listH []H, listI []I, transform func(a A, b B, c C, d D, e E, f F, g G, h H, i I) Out) []Out {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF) * len(listG) * len(listH) * len(listI)
	if size == 0 {
		return []Out{}
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							for _, g := range listG {
								for _, h := range listH {
									for _, i := range listI {
										result = append(result, transform(a, b, c, d, e, f, g, h, i))
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

// CrossJoinByErr2 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr2[A, B, Out any](listA []A, listB []B, transform func(a A, b B) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			r, err := transform(a, b)
			if err != nil {
				return nil, err
			}
			result = append(result, r)
		}
	}

	return result, nil
}

// CrossJoinByErr3 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr3[A, B, C, Out any](listA []A, listB []B, listC []C, transform func(a A, b B, c C) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB) * len(listC)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				r, err := transform(a, b, c)
				if err != nil {
					return nil, err
				}
				result = append(result, r)
			}
		}
	}

	return result, nil
}

// CrossJoinByErr4 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr4[A, B, C, D, Out any](listA []A, listB []B, listC []C, listD []D, transform func(a A, b B, c C, d D) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB) * len(listC) * len(listD)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					r, err := transform(a, b, c, d)
					if err != nil {
						return nil, err
					}
					result = append(result, r)
				}
			}
		}
	}

	return result, nil
}

// CrossJoinByErr5 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr5[A, B, C, D, E, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, transform func(a A, b B, c C, d D, e E) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						r, err := transform(a, b, c, d, e)
						if err != nil {
							return nil, err
						}
						result = append(result, r)
					}
				}
			}
		}
	}

	return result, nil
}

// CrossJoinByErr6 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr6[A, B, C, D, E, F, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, transform func(a A, b B, c C, d D, e E, f F) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							r, err := transform(a, b, c, d, e, f)
							if err != nil {
								return nil, err
							}
							result = append(result, r)
						}
					}
				}
			}
		}
	}

	return result, nil
}

// CrossJoinByErr7 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr7[A, B, C, D, E, F, G, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, transform func(a A, b B, c C, d D, e E, f F, g G) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF) * len(listG)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							for _, g := range listG {
								r, err := transform(a, b, c, d, e, f, g)
								if err != nil {
									return nil, err
								}
								result = append(result, r)
							}
						}
					}
				}
			}
		}
	}

	return result, nil
}

// CrossJoinByErr8 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr8[A, B, C, D, E, F, G, H, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, listH []H, transform func(a A, b B, c C, d D, e E, f F, g G, h H) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF) * len(listG) * len(listH)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							for _, g := range listG {
								for _, h := range listH {
									r, err := transform(a, b, c, d, e, f, g, h)
									if err != nil {
										return nil, err
									}
									result = append(result, r)
								}
							}
						}
					}
				}
			}
		}
	}

	return result, nil
}

// CrossJoinByErr9 combines every item from one list with every item from others.
// It is the cartesian product of lists received as arguments. The transform function
// is used to create the output values.
// Returns an empty list if a list is empty.
// It returns the first error returned by the transform function.
func CrossJoinByErr9[A, B, C, D, E, F, G, H, I, Out any](listA []A, listB []B, listC []C, listD []D, listE []E, listF []F, listG []G, listH []H, listI []I, transform func(a A, b B, c C, d D, e E, f F, g G, h H, i I) (Out, error)) ([]Out, error) {
	size := len(listA) * len(listB) * len(listC) * len(listD) * len(listE) * len(listF) * len(listG) * len(listH) * len(listI)
	if size == 0 {
		return []Out{}, nil
	}

	result := make([]Out, 0, size)

	for _, a := range listA {
		for _, b := range listB {
			for _, c := range listC {
				for _, d := range listD {
					for _, e := range listE {
						for _, f := range listF {
							for _, g := range listG {
								for _, h := range listH {
									for _, i := range listI {
										r, err := transform(a, b, c, d, e, f, g, h, i)
										if err != nil {
											return nil, err
										}
										result = append(result, r)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return result, nil
}
