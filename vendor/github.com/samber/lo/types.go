package lo

// Entry defines a key/value pairs.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Tuple2 is a group of 2 elements (pair).
type Tuple2[A any, B any] struct {
	A A
	B B
}

// Tuple3 is a group of 3 elements.
type Tuple3[A any, B any, C any] struct {
	A A
	B B
	C C
}

// Tuple4 is a group of 4 elements.
type Tuple4[A any, B any, C any, D any] struct {
	A A
	B B
	C C
	D D
}

// Tuple5 is a group of 5 elements.
type Tuple5[A any, B any, C any, D any, E any] struct {
	A A
	B B
	C C
	D D
	E E
}

// Tuple6 is a group of 6 elements.
type Tuple6[A any, B any, C any, D any, E any, F any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
}

// Tuple7 is a group of 7 elements.
type Tuple7[A any, B any, C any, D any, E any, F any, G any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
}

// Tuple8 is a group of 8 elements.
type Tuple8[A any, B any, C any, D any, E any, F any, G any, H any] struct {
	A A
	B B
	C C
	D D
	E E
	F F
	G G
	H H
}

// Tuple9 is a group of 9 elements.
type Tuple9[A any, B any, C any, D any, E any, F any, G any, H any, I any] struct {
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
