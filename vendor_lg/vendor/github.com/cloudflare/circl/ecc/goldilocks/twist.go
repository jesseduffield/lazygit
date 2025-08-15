package goldilocks

import (
	"crypto/subtle"
	"math/bits"

	"github.com/cloudflare/circl/internal/conv"
	"github.com/cloudflare/circl/math"
	fp "github.com/cloudflare/circl/math/fp448"
)

// twistCurve is -x^2+y^2=1-39082x^2y^2 and is 4-isogenous to Goldilocks.
type twistCurve struct{}

// Identity returns the identity point.
func (twistCurve) Identity() *twistPoint {
	return &twistPoint{
		y: fp.One(),
		z: fp.One(),
	}
}

// subYDiv16 update x = (x - y) / 16.
func subYDiv16(x *scalar64, y int64) {
	s := uint64(y >> 63)
	x0, b0 := bits.Sub64((*x)[0], uint64(y), 0)
	x1, b1 := bits.Sub64((*x)[1], s, b0)
	x2, b2 := bits.Sub64((*x)[2], s, b1)
	x3, b3 := bits.Sub64((*x)[3], s, b2)
	x4, b4 := bits.Sub64((*x)[4], s, b3)
	x5, b5 := bits.Sub64((*x)[5], s, b4)
	x6, _ := bits.Sub64((*x)[6], s, b5)
	x[0] = (x0 >> 4) | (x1 << 60)
	x[1] = (x1 >> 4) | (x2 << 60)
	x[2] = (x2 >> 4) | (x3 << 60)
	x[3] = (x3 >> 4) | (x4 << 60)
	x[4] = (x4 >> 4) | (x5 << 60)
	x[5] = (x5 >> 4) | (x6 << 60)
	x[6] = (x6 >> 4)
}

func recodeScalar(d *[113]int8, k *Scalar) {
	var k64 scalar64
	k64.fromScalar(k)
	for i := 0; i < 112; i++ {
		d[i] = int8((k64[0] & 0x1f) - 16)
		subYDiv16(&k64, int64(d[i]))
	}
	d[112] = int8(k64[0])
}

// ScalarMult returns kP.
func (e twistCurve) ScalarMult(k *Scalar, P *twistPoint) *twistPoint {
	var TabP [8]preTwistPointProy
	var S preTwistPointProy
	var d [113]int8

	var isZero int
	if k.IsZero() {
		isZero = 1
	}
	subtle.ConstantTimeCopy(isZero, k[:], order[:])

	minusK := *k
	isEven := 1 - int(k[0]&0x1)
	minusK.Neg()
	subtle.ConstantTimeCopy(isEven, k[:], minusK[:])
	recodeScalar(&d, k)

	P.oddMultiples(TabP[:])
	Q := e.Identity()
	for i := 112; i >= 0; i-- {
		Q.Double()
		Q.Double()
		Q.Double()
		Q.Double()
		mask := d[i] >> 7
		absDi := (d[i] + mask) ^ mask
		inx := int32((absDi - 1) >> 1)
		sig := int((d[i] >> 7) & 0x1)
		for j := range TabP {
			S.cmov(&TabP[j], uint(subtle.ConstantTimeEq(inx, int32(j))))
		}
		S.cneg(sig)
		Q.mixAdd(&S)
	}
	Q.cneg(uint(isEven))
	return Q
}

const (
	omegaFix = 7
	omegaVar = 5
)

// CombinedMult returns mG+nP.
func (e twistCurve) CombinedMult(m, n *Scalar, P *twistPoint) *twistPoint {
	nafFix := math.OmegaNAF(conv.BytesLe2BigInt(m[:]), omegaFix)
	nafVar := math.OmegaNAF(conv.BytesLe2BigInt(n[:]), omegaVar)

	if len(nafFix) > len(nafVar) {
		nafVar = append(nafVar, make([]int32, len(nafFix)-len(nafVar))...)
	} else if len(nafFix) < len(nafVar) {
		nafFix = append(nafFix, make([]int32, len(nafVar)-len(nafFix))...)
	}

	var TabQ [1 << (omegaVar - 2)]preTwistPointProy
	P.oddMultiples(TabQ[:])
	Q := e.Identity()
	for i := len(nafFix) - 1; i >= 0; i-- {
		Q.Double()
		// Generator point
		if nafFix[i] != 0 {
			idxM := absolute(nafFix[i]) >> 1
			R := tabVerif[idxM]
			if nafFix[i] < 0 {
				R.neg()
			}
			Q.mixAddZ1(&R)
		}
		// Variable input point
		if nafVar[i] != 0 {
			idxN := absolute(nafVar[i]) >> 1
			S := TabQ[idxN]
			if nafVar[i] < 0 {
				S.neg()
			}
			Q.mixAdd(&S)
		}
	}
	return Q
}

// absolute returns always a positive value.
func absolute(x int32) int32 {
	mask := x >> 31
	return (x + mask) ^ mask
}
