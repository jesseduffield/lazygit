package x448

import (
	"encoding/binary"
	"math/bits"

	"github.com/cloudflare/circl/math/fp448"
)

func doubleGeneric(x, z *fp448.Elt) {
	t0, t1 := &fp448.Elt{}, &fp448.Elt{}
	fp448.AddSub(x, z)
	fp448.Sqr(x, x)
	fp448.Sqr(z, z)
	fp448.Sub(t0, x, z)
	mulA24Generic(t1, t0)
	fp448.Add(t1, t1, z)
	fp448.Mul(x, x, z)
	fp448.Mul(z, t0, t1)
}

func diffAddGeneric(w *[5]fp448.Elt, b uint) {
	mu, x1, z1, x2, z2 := &w[0], &w[1], &w[2], &w[3], &w[4]
	fp448.Cswap(x1, x2, b)
	fp448.Cswap(z1, z2, b)
	fp448.AddSub(x1, z1)
	fp448.Mul(z1, z1, mu)
	fp448.AddSub(x1, z1)
	fp448.Sqr(x1, x1)
	fp448.Sqr(z1, z1)
	fp448.Mul(x1, x1, z2)
	fp448.Mul(z1, z1, x2)
}

func ladderStepGeneric(w *[5]fp448.Elt, b uint) {
	x1, x2, z2, x3, z3 := &w[0], &w[1], &w[2], &w[3], &w[4]
	t0 := &fp448.Elt{}
	t1 := &fp448.Elt{}
	fp448.AddSub(x2, z2)
	fp448.AddSub(x3, z3)
	fp448.Mul(t0, x2, z3)
	fp448.Mul(t1, x3, z2)
	fp448.AddSub(t0, t1)
	fp448.Cmov(x2, x3, b)
	fp448.Cmov(z2, z3, b)
	fp448.Sqr(x3, t0)
	fp448.Sqr(z3, t1)
	fp448.Mul(z3, x1, z3)
	fp448.Sqr(x2, x2)
	fp448.Sqr(z2, z2)
	fp448.Sub(t0, x2, z2)
	mulA24Generic(t1, t0)
	fp448.Add(t1, t1, z2)
	fp448.Mul(x2, x2, z2)
	fp448.Mul(z2, t0, t1)
}

func mulA24Generic(z, x *fp448.Elt) {
	const A24 = 39082
	const n = 8
	var xx [7]uint64
	for i := range xx {
		xx[i] = binary.LittleEndian.Uint64(x[i*n : (i+1)*n])
	}
	h0, l0 := bits.Mul64(xx[0], A24)
	h1, l1 := bits.Mul64(xx[1], A24)
	h2, l2 := bits.Mul64(xx[2], A24)
	h3, l3 := bits.Mul64(xx[3], A24)
	h4, l4 := bits.Mul64(xx[4], A24)
	h5, l5 := bits.Mul64(xx[5], A24)
	h6, l6 := bits.Mul64(xx[6], A24)

	l1, c0 := bits.Add64(h0, l1, 0)
	l2, c1 := bits.Add64(h1, l2, c0)
	l3, c2 := bits.Add64(h2, l3, c1)
	l4, c3 := bits.Add64(h3, l4, c2)
	l5, c4 := bits.Add64(h4, l5, c3)
	l6, c5 := bits.Add64(h5, l6, c4)
	l7, _ := bits.Add64(h6, 0, c5)

	l0, c0 = bits.Add64(l0, l7, 0)
	l1, c1 = bits.Add64(l1, 0, c0)
	l2, c2 = bits.Add64(l2, 0, c1)
	l3, c3 = bits.Add64(l3, l7<<32, c2)
	l4, c4 = bits.Add64(l4, 0, c3)
	l5, c5 = bits.Add64(l5, 0, c4)
	l6, l7 = bits.Add64(l6, 0, c5)

	xx[0], c0 = bits.Add64(l0, l7, 0)
	xx[1], c1 = bits.Add64(l1, 0, c0)
	xx[2], c2 = bits.Add64(l2, 0, c1)
	xx[3], c3 = bits.Add64(l3, l7<<32, c2)
	xx[4], c4 = bits.Add64(l4, 0, c3)
	xx[5], c5 = bits.Add64(l5, 0, c4)
	xx[6], _ = bits.Add64(l6, 0, c5)

	for i := range xx {
		binary.LittleEndian.PutUint64(z[i*n:(i+1)*n], xx[i])
	}
}
