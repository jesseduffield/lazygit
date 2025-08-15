package x25519

import (
	"encoding/binary"
	"math/bits"

	fp "github.com/cloudflare/circl/math/fp25519"
)

func doubleGeneric(x, z *fp.Elt) {
	t0, t1 := &fp.Elt{}, &fp.Elt{}
	fp.AddSub(x, z)
	fp.Sqr(x, x)
	fp.Sqr(z, z)
	fp.Sub(t0, x, z)
	mulA24Generic(t1, t0)
	fp.Add(t1, t1, z)
	fp.Mul(x, x, z)
	fp.Mul(z, t0, t1)
}

func diffAddGeneric(w *[5]fp.Elt, b uint) {
	mu, x1, z1, x2, z2 := &w[0], &w[1], &w[2], &w[3], &w[4]
	fp.Cswap(x1, x2, b)
	fp.Cswap(z1, z2, b)
	fp.AddSub(x1, z1)
	fp.Mul(z1, z1, mu)
	fp.AddSub(x1, z1)
	fp.Sqr(x1, x1)
	fp.Sqr(z1, z1)
	fp.Mul(x1, x1, z2)
	fp.Mul(z1, z1, x2)
}

func ladderStepGeneric(w *[5]fp.Elt, b uint) {
	x1, x2, z2, x3, z3 := &w[0], &w[1], &w[2], &w[3], &w[4]
	t0 := &fp.Elt{}
	t1 := &fp.Elt{}
	fp.AddSub(x2, z2)
	fp.AddSub(x3, z3)
	fp.Mul(t0, x2, z3)
	fp.Mul(t1, x3, z2)
	fp.AddSub(t0, t1)
	fp.Cmov(x2, x3, b)
	fp.Cmov(z2, z3, b)
	fp.Sqr(x3, t0)
	fp.Sqr(z3, t1)
	fp.Mul(z3, x1, z3)
	fp.Sqr(x2, x2)
	fp.Sqr(z2, z2)
	fp.Sub(t0, x2, z2)
	mulA24Generic(t1, t0)
	fp.Add(t1, t1, z2)
	fp.Mul(x2, x2, z2)
	fp.Mul(z2, t0, t1)
}

func mulA24Generic(z, x *fp.Elt) {
	const A24 = 121666
	const n = 8
	var xx [4]uint64
	for i := range xx {
		xx[i] = binary.LittleEndian.Uint64(x[i*n : (i+1)*n])
	}

	h0, l0 := bits.Mul64(xx[0], A24)
	h1, l1 := bits.Mul64(xx[1], A24)
	h2, l2 := bits.Mul64(xx[2], A24)
	h3, l3 := bits.Mul64(xx[3], A24)

	var c3 uint64
	l1, c0 := bits.Add64(h0, l1, 0)
	l2, c1 := bits.Add64(h1, l2, c0)
	l3, c2 := bits.Add64(h2, l3, c1)
	l4, _ := bits.Add64(h3, 0, c2)
	_, l4 = bits.Mul64(l4, 38)
	l0, c0 = bits.Add64(l0, l4, 0)
	xx[1], c1 = bits.Add64(l1, 0, c0)
	xx[2], c2 = bits.Add64(l2, 0, c1)
	xx[3], c3 = bits.Add64(l3, 0, c2)
	xx[0], _ = bits.Add64(l0, (-c3)&38, 0)
	for i := range xx {
		binary.LittleEndian.PutUint64(z[i*n:(i+1)*n], xx[i])
	}
}
