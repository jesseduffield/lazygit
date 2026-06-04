package lo

// Entry defines a key/value pairs.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Tuple2 is a group of 2 elements (pair).
type Tuple2[A, B any] struct {
	A A
	B B
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/yrtn7QJTmL_E
func (t Tuple2[A, B]) Unpack() (A, B) {
	return t.A, t.B
}

// Tuple3 is a group of 3 elements.
type Tuple3[A, B, C any] struct {
	A A
	B B
	C C
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/yrtn7QJTmL_E
func (t Tuple3[A, B, C]) Unpack() (A, B, C) {
	return t.A, t.B, t.C
}

// Tuple4 is a group of 4 elements.
type Tuple4[A, B, C, D any] struct {
	A A
	B B
	C C
	D D
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/yrtn7QJTmL_E
func (t Tuple4[A, B, C, D]) Unpack() (A, B, C, D) {
	return t.A, t.B, t.C, t.D
}

// Tuple5 is a group of 5 elements.
type Tuple5[A, B, C, D, E any] struct {
	A A
	B B
	C C
	D D
	E E
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/7J4KrtgtK3M
func (t Tuple5[A, B, C, D, E]) Unpack() (A, B, C, D, E) {
	return t.A, t.B, t.C, t.D, t.E
}

// Tuple6 is a group of 6 elements.
type Tuple6[A, B, C, D, E, F any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/7J4KrtgtK3M
func (t Tuple6[A, B, C, D, E, F]) Unpack() (A, B, C, D, E, F) {
	return t.A, t.B, t.C, t.D, t.E, t.F
}

// Tuple7 is a group of 7 elements.
type Tuple7[A, B, C, D, E, F, G any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/Ow9Zgf_zeiA
func (t Tuple7[A, B, C, D, E, F, G]) Unpack() (A, B, C, D, E, F, G) {
	return t.A, t.B, t.C, t.D, t.E, t.F, t.G
}

// Tuple8 is a group of 8 elements.
type Tuple8[A, B, C, D, E, F, G, H any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
	H H
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/Ow9Zgf_zeiA
func (t Tuple8[A, B, C, D, E, F, G, H]) Unpack() (A, B, C, D, E, F, G, H) {
	return t.A, t.B, t.C, t.D, t.E, t.F, t.G, t.H
}

// Tuple9 is a group of 9 elements.
type Tuple9[A, B, C, D, E, F, G, H, I any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
	H H
	I I
}

// Unpack returns values contained in a tuple.
// Play: https://go.dev/play/p/Ow9Zgf_zeiA
func (t Tuple9[A, B, C, D, E, F, G, H, I]) Unpack() (A, B, C, D, E, F, G, H, I) {
	return t.A, t.B, t.C, t.D, t.E, t.F, t.G, t.H, t.I
}
