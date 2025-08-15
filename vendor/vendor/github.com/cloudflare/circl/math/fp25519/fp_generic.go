package fp25519

import (
	"encoding/binary"
	"math/bits"
)

func cmovGeneric(x, y *Elt, n uint) {
	m := -uint64(n & 0x1)
	x0 := binary.LittleEndian.Uint64(x[0*8 : 1*8])
	x1 := binary.LittleEndian.Uint64(x[1*8 : 2*8])
	x2 := binary.LittleEndian.Uint64(x[2*8 : 3*8])
	x3 := binary.LittleEndian.Uint64(x[3*8 : 4*8])

	y0 := binary.LittleEndian.Uint64(y[0*8 : 1*8])
	y1 := binary.LittleEndian.Uint64(y[1*8 : 2*8])
	y2 := binary.LittleEndian.Uint64(y[2*8 : 3*8])
	y3 := binary.LittleEndian.Uint64(y[3*8 : 4*8])

	x0 = (x0 &^ m) | (y0 & m)
	x1 = (x1 &^ m) | (y1 & m)
	x2 = (x2 &^ m) | (y2 & m)
	x3 = (x3 &^ m) | (y3 & m)

	binary.LittleEndian.PutUint64(x[0*8:1*8], x0)
	binary.LittleEndian.PutUint64(x[1*8:2*8], x1)
	binary.LittleEndian.PutUint64(x[2*8:3*8], x2)
	binary.LittleEndian.PutUint64(x[3*8:4*8], x3)
}

func cswapGeneric(x, y *Elt, n uint) {
	m := -uint64(n & 0x1)
	x0 := binary.LittleEndian.Uint64(x[0*8 : 1*8])
	x1 := binary.LittleEndian.Uint64(x[1*8 : 2*8])
	x2 := binary.LittleEndian.Uint64(x[2*8 : 3*8])
	x3 := binary.LittleEndian.Uint64(x[3*8 : 4*8])

	y0 := binary.LittleEndian.Uint64(y[0*8 : 1*8])
	y1 := binary.LittleEndian.Uint64(y[1*8 : 2*8])
	y2 := binary.LittleEndian.Uint64(y[2*8 : 3*8])
	y3 := binary.LittleEndian.Uint64(y[3*8 : 4*8])

	t0 := m & (x0 ^ y0)
	t1 := m & (x1 ^ y1)
	t2 := m & (x2 ^ y2)
	t3 := m & (x3 ^ y3)
	x0 ^= t0
	x1 ^= t1
	x2 ^= t2
	x3 ^= t3
	y0 ^= t0
	y1 ^= t1
	y2 ^= t2
	y3 ^= t3

	binary.LittleEndian.PutUint64(x[0*8:1*8], x0)
	binary.LittleEndian.PutUint64(x[1*8:2*8], x1)
	binary.LittleEndian.PutUint64(x[2*8:3*8], x2)
	binary.LittleEndian.PutUint64(x[3*8:4*8], x3)

	binary.LittleEndian.PutUint64(y[0*8:1*8], y0)
	binary.LittleEndian.PutUint64(y[1*8:2*8], y1)
	binary.LittleEndian.PutUint64(y[2*8:3*8], y2)
	binary.LittleEndian.PutUint64(y[3*8:4*8], y3)
}

func addGeneric(z, x, y *Elt) {
	x0 := binary.LittleEndian.Uint64(x[0*8 : 1*8])
	x1 := binary.LittleEndian.Uint64(x[1*8 : 2*8])
	x2 := binary.LittleEndian.Uint64(x[2*8 : 3*8])
	x3 := binary.LittleEndian.Uint64(x[3*8 : 4*8])

	y0 := binary.LittleEndian.Uint64(y[0*8 : 1*8])
	y1 := binary.LittleEndian.Uint64(y[1*8 : 2*8])
	y2 := binary.LittleEndian.Uint64(y[2*8 : 3*8])
	y3 := binary.LittleEndian.Uint64(y[3*8 : 4*8])

	z0, c0 := bits.Add64(x0, y0, 0)
	z1, c1 := bits.Add64(x1, y1, c0)
	z2, c2 := bits.Add64(x2, y2, c1)
	z3, c3 := bits.Add64(x3, y3, c2)

	z0, c0 = bits.Add64(z0, (-c3)&38, 0)
	z1, c1 = bits.Add64(z1, 0, c0)
	z2, c2 = bits.Add64(z2, 0, c1)
	z3, c3 = bits.Add64(z3, 0, c2)
	z0, _ = bits.Add64(z0, (-c3)&38, 0)

	binary.LittleEndian.PutUint64(z[0*8:1*8], z0)
	binary.LittleEndian.PutUint64(z[1*8:2*8], z1)
	binary.LittleEndian.PutUint64(z[2*8:3*8], z2)
	binary.LittleEndian.PutUint64(z[3*8:4*8], z3)
}

func subGeneric(z, x, y *Elt) {
	x0 := binary.LittleEndian.Uint64(x[0*8 : 1*8])
	x1 := binary.LittleEndian.Uint64(x[1*8 : 2*8])
	x2 := binary.LittleEndian.Uint64(x[2*8 : 3*8])
	x3 := binary.LittleEndian.Uint64(x[3*8 : 4*8])

	y0 := binary.LittleEndian.Uint64(y[0*8 : 1*8])
	y1 := binary.LittleEndian.Uint64(y[1*8 : 2*8])
	y2 := binary.LittleEndian.Uint64(y[2*8 : 3*8])
	y3 := binary.LittleEndian.Uint64(y[3*8 : 4*8])

	z0, c0 := bits.Sub64(x0, y0, 0)
	z1, c1 := bits.Sub64(x1, y1, c0)
	z2, c2 := bits.Sub64(x2, y2, c1)
	z3, c3 := bits.Sub64(x3, y3, c2)

	z0, c0 = bits.Sub64(z0, (-c3)&38, 0)
	z1, c1 = bits.Sub64(z1, 0, c0)
	z2, c2 = bits.Sub64(z2, 0, c1)
	z3, c3 = bits.Sub64(z3, 0, c2)
	z0, _ = bits.Sub64(z0, (-c3)&38, 0)

	binary.LittleEndian.PutUint64(z[0*8:1*8], z0)
	binary.LittleEndian.PutUint64(z[1*8:2*8], z1)
	binary.LittleEndian.PutUint64(z[2*8:3*8], z2)
	binary.LittleEndian.PutUint64(z[3*8:4*8], z3)
}

func addsubGeneric(x, y *Elt) {
	z := &Elt{}
	addGeneric(z, x, y)
	subGeneric(y, x, y)
	*x = *z
}

func mulGeneric(z, x, y *Elt) {
	x0 := binary.LittleEndian.Uint64(x[0*8 : 1*8])
	x1 := binary.LittleEndian.Uint64(x[1*8 : 2*8])
	x2 := binary.LittleEndian.Uint64(x[2*8 : 3*8])
	x3 := binary.LittleEndian.Uint64(x[3*8 : 4*8])

	y0 := binary.LittleEndian.Uint64(y[0*8 : 1*8])
	y1 := binary.LittleEndian.Uint64(y[1*8 : 2*8])
	y2 := binary.LittleEndian.Uint64(y[2*8 : 3*8])
	y3 := binary.LittleEndian.Uint64(y[3*8 : 4*8])

	yi := y0
	h0, l0 := bits.Mul64(x0, yi)
	h1, l1 := bits.Mul64(x1, yi)
	h2, l2 := bits.Mul64(x2, yi)
	h3, l3 := bits.Mul64(x3, yi)

	z0 := l0
	a0, c0 := bits.Add64(h0, l1, 0)
	a1, c1 := bits.Add64(h1, l2, c0)
	a2, c2 := bits.Add64(h2, l3, c1)
	a3, _ := bits.Add64(h3, 0, c2)

	yi = y1
	h0, l0 = bits.Mul64(x0, yi)
	h1, l1 = bits.Mul64(x1, yi)
	h2, l2 = bits.Mul64(x2, yi)
	h3, l3 = bits.Mul64(x3, yi)

	z1, c0 := bits.Add64(a0, l0, 0)
	h0, c1 = bits.Add64(h0, l1, c0)
	h1, c2 = bits.Add64(h1, l2, c1)
	h2, c3 := bits.Add64(h2, l3, c2)
	h3, _ = bits.Add64(h3, 0, c3)

	a0, c0 = bits.Add64(a1, h0, 0)
	a1, c1 = bits.Add64(a2, h1, c0)
	a2, c2 = bits.Add64(a3, h2, c1)
	a3, _ = bits.Add64(0, h3, c2)

	yi = y2
	h0, l0 = bits.Mul64(x0, yi)
	h1, l1 = bits.Mul64(x1, yi)
	h2, l2 = bits.Mul64(x2, yi)
	h3, l3 = bits.Mul64(x3, yi)

	z2, c0 := bits.Add64(a0, l0, 0)
	h0, c1 = bits.Add64(h0, l1, c0)
	h1, c2 = bits.Add64(h1, l2, c1)
	h2, c3 = bits.Add64(h2, l3, c2)
	h3, _ = bits.Add64(h3, 0, c3)

	a0, c0 = bits.Add64(a1, h0, 0)
	a1, c1 = bits.Add64(a2, h1, c0)
	a2, c2 = bits.Add64(a3, h2, c1)
	a3, _ = bits.Add64(0, h3, c2)

	yi = y3
	h0, l0 = bits.Mul64(x0, yi)
	h1, l1 = bits.Mul64(x1, yi)
	h2, l2 = bits.Mul64(x2, yi)
	h3, l3 = bits.Mul64(x3, yi)

	z3, c0 := bits.Add64(a0, l0, 0)
	h0, c1 = bits.Add64(h0, l1, c0)
	h1, c2 = bits.Add64(h1, l2, c1)
	h2, c3 = bits.Add64(h2, l3, c2)
	h3, _ = bits.Add64(h3, 0, c3)

	z4, c0 := bits.Add64(a1, h0, 0)
	z5, c1 := bits.Add64(a2, h1, c0)
	z6, c2 := bits.Add64(a3, h2, c1)
	z7, _ := bits.Add64(0, h3, c2)

	red64(z, z0, z1, z2, z3, z4, z5, z6, z7)
}

func sqrGeneric(z, x *Elt) {
	x0 := binary.LittleEndian.Uint64(x[0*8 : 1*8])
	x1 := binary.LittleEndian.Uint64(x[1*8 : 2*8])
	x2 := binary.LittleEndian.Uint64(x[2*8 : 3*8])
	x3 := binary.LittleEndian.Uint64(x[3*8 : 4*8])

	h0, a0 := bits.Mul64(x0, x1)
	h1, l1 := bits.Mul64(x0, x2)
	h2, l2 := bits.Mul64(x0, x3)
	h3, l3 := bits.Mul64(x3, x1)
	h4, l4 := bits.Mul64(x3, x2)
	h, l := bits.Mul64(x1, x2)

	a1, c0 := bits.Add64(l1, h0, 0)
	a2, c1 := bits.Add64(l2, h1, c0)
	a3, c2 := bits.Add64(l3, h2, c1)
	a4, c3 := bits.Add64(l4, h3, c2)
	a5, _ := bits.Add64(h4, 0, c3)

	a2, c0 = bits.Add64(a2, l, 0)
	a3, c1 = bits.Add64(a3, h, c0)
	a4, c2 = bits.Add64(a4, 0, c1)
	a5, c3 = bits.Add64(a5, 0, c2)
	a6, _ := bits.Add64(0, 0, c3)

	a0, c0 = bits.Add64(a0, a0, 0)
	a1, c1 = bits.Add64(a1, a1, c0)
	a2, c2 = bits.Add64(a2, a2, c1)
	a3, c3 = bits.Add64(a3, a3, c2)
	a4, c4 := bits.Add64(a4, a4, c3)
	a5, c5 := bits.Add64(a5, a5, c4)
	a6, _ = bits.Add64(a6, a6, c5)

	b1, b0 := bits.Mul64(x0, x0)
	b3, b2 := bits.Mul64(x1, x1)
	b5, b4 := bits.Mul64(x2, x2)
	b7, b6 := bits.Mul64(x3, x3)

	b1, c0 = bits.Add64(b1, a0, 0)
	b2, c1 = bits.Add64(b2, a1, c0)
	b3, c2 = bits.Add64(b3, a2, c1)
	b4, c3 = bits.Add64(b4, a3, c2)
	b5, c4 = bits.Add64(b5, a4, c3)
	b6, c5 = bits.Add64(b6, a5, c4)
	b7, _ = bits.Add64(b7, a6, c5)

	red64(z, b0, b1, b2, b3, b4, b5, b6, b7)
}

func modpGeneric(x *Elt) {
	x0 := binary.LittleEndian.Uint64(x[0*8 : 1*8])
	x1 := binary.LittleEndian.Uint64(x[1*8 : 2*8])
	x2 := binary.LittleEndian.Uint64(x[2*8 : 3*8])
	x3 := binary.LittleEndian.Uint64(x[3*8 : 4*8])

	// CX = C[255] ? 38 : 19
	cx := uint64(19) << (x3 >> 63)
	// PUT BIT 255 IN CARRY FLAG AND CLEAR
	x3 &^= 1 << 63

	x0, c0 := bits.Add64(x0, cx, 0)
	x1, c1 := bits.Add64(x1, 0, c0)
	x2, c2 := bits.Add64(x2, 0, c1)
	x3, _ = bits.Add64(x3, 0, c2)

	// TEST FOR BIT 255 AGAIN; ONLY TRIGGERED ON OVERFLOW MODULO 2^255-19
	// cx = C[255] ? 0 : 19
	cx = uint64(19) &^ (-(x3 >> 63))
	// CLEAR BIT 255
	x3 &^= 1 << 63

	x0, c0 = bits.Sub64(x0, cx, 0)
	x1, c1 = bits.Sub64(x1, 0, c0)
	x2, c2 = bits.Sub64(x2, 0, c1)
	x3, _ = bits.Sub64(x3, 0, c2)

	binary.LittleEndian.PutUint64(x[0*8:1*8], x0)
	binary.LittleEndian.PutUint64(x[1*8:2*8], x1)
	binary.LittleEndian.PutUint64(x[2*8:3*8], x2)
	binary.LittleEndian.PutUint64(x[3*8:4*8], x3)
}

func red64(z *Elt, x0, x1, x2, x3, x4, x5, x6, x7 uint64) {
	h0, l0 := bits.Mul64(x4, 38)
	h1, l1 := bits.Mul64(x5, 38)
	h2, l2 := bits.Mul64(x6, 38)
	h3, l3 := bits.Mul64(x7, 38)

	l1, c0 := bits.Add64(h0, l1, 0)
	l2, c1 := bits.Add64(h1, l2, c0)
	l3, c2 := bits.Add64(h2, l3, c1)
	l4, _ := bits.Add64(h3, 0, c2)

	l0, c0 = bits.Add64(l0, x0, 0)
	l1, c1 = bits.Add64(l1, x1, c0)
	l2, c2 = bits.Add64(l2, x2, c1)
	l3, c3 := bits.Add64(l3, x3, c2)
	l4, _ = bits.Add64(l4, 0, c3)

	_, l4 = bits.Mul64(l4, 38)
	l0, c0 = bits.Add64(l0, l4, 0)
	z1, c1 := bits.Add64(l1, 0, c0)
	z2, c2 := bits.Add64(l2, 0, c1)
	z3, c3 := bits.Add64(l3, 0, c2)
	z0, _ := bits.Add64(l0, (-c3)&38, 0)

	binary.LittleEndian.PutUint64(z[0*8:1*8], z0)
	binary.LittleEndian.PutUint64(z[1*8:2*8], z1)
	binary.LittleEndian.PutUint64(z[2*8:3*8], z2)
	binary.LittleEndian.PutUint64(z[3*8:4*8], z3)
}
