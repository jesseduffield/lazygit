package ed25519

import (
	"crypto/subtle"
	"encoding/binary"
	"math/bits"

	"github.com/cloudflare/circl/internal/conv"
	"github.com/cloudflare/circl/math"
	fp "github.com/cloudflare/circl/math/fp25519"
)

var paramD = fp.Elt{
	0xa3, 0x78, 0x59, 0x13, 0xca, 0x4d, 0xeb, 0x75,
	0xab, 0xd8, 0x41, 0x41, 0x4d, 0x0a, 0x70, 0x00,
	0x98, 0xe8, 0x79, 0x77, 0x79, 0x40, 0xc7, 0x8c,
	0x73, 0xfe, 0x6f, 0x2b, 0xee, 0x6c, 0x03, 0x52,
}

// mLSBRecoding parameters.
const (
	fxT        = 257
	fxV        = 2
	fxW        = 3
	fx2w1      = 1 << (uint(fxW) - 1)
	numWords64 = (paramB * 8 / 64)
)

// mLSBRecoding is the odd-only modified LSB-set.
//
// Reference:
//
//	"Efficient and secure algorithms for GLV-based scalar multiplication and
//	 their implementation on GLV–GLS curves" by (Faz-Hernandez et al.)
//	 http://doi.org/10.1007/s13389-014-0085-7.
func mLSBRecoding(L []int8, k []byte) {
	const ee = (fxT + fxW*fxV - 1) / (fxW * fxV)
	const dd = ee * fxV
	const ll = dd * fxW
	if len(L) == (ll + 1) {
		var m [numWords64 + 1]uint64
		for i := 0; i < numWords64; i++ {
			m[i] = binary.LittleEndian.Uint64(k[8*i : 8*i+8])
		}
		condAddOrderN(&m)
		L[dd-1] = 1
		for i := 0; i < dd-1; i++ {
			kip1 := (m[(i+1)/64] >> (uint(i+1) % 64)) & 0x1
			L[i] = int8(kip1<<1) - 1
		}
		{ // right-shift by d
			right := uint(dd % 64)
			left := uint(64) - right
			lim := ((numWords64+1)*64 - dd) / 64
			j := dd / 64
			for i := 0; i < lim; i++ {
				m[i] = (m[i+j] >> right) | (m[i+j+1] << left)
			}
			m[lim] = m[lim+j] >> right
		}
		for i := dd; i < ll; i++ {
			L[i] = L[i%dd] * int8(m[0]&0x1)
			div2subY(m[:], int64(L[i]>>1), numWords64)
		}
		L[ll] = int8(m[0])
	}
}

// absolute returns always a positive value.
func absolute(x int32) int32 {
	mask := x >> 31
	return (x + mask) ^ mask
}

// condAddOrderN updates x = x+order if x is even, otherwise x remains unchanged.
func condAddOrderN(x *[numWords64 + 1]uint64) {
	isOdd := (x[0] & 0x1) - 1
	c := uint64(0)
	for i := 0; i < numWords64; i++ {
		orderWord := binary.LittleEndian.Uint64(order[8*i : 8*i+8])
		o := isOdd & orderWord
		x0, c0 := bits.Add64(x[i], o, c)
		x[i] = x0
		c = c0
	}
	x[numWords64], _ = bits.Add64(x[numWords64], 0, c)
}

// div2subY update x = (x/2) - y.
func div2subY(x []uint64, y int64, l int) {
	s := uint64(y >> 63)
	for i := 0; i < l-1; i++ {
		x[i] = (x[i] >> 1) | (x[i+1] << 63)
	}
	x[l-1] = (x[l-1] >> 1)

	b := uint64(0)
	x0, b0 := bits.Sub64(x[0], uint64(y), b)
	x[0] = x0
	b = b0
	for i := 1; i < l-1; i++ {
		x0, b0 := bits.Sub64(x[i], s, b)
		x[i] = x0
		b = b0
	}
	x[l-1], _ = bits.Sub64(x[l-1], s, b)
}

func (P *pointR1) fixedMult(scalar []byte) {
	if len(scalar) != paramB {
		panic("wrong scalar size")
	}
	const ee = (fxT + fxW*fxV - 1) / (fxW * fxV)
	const dd = ee * fxV
	const ll = dd * fxW

	L := make([]int8, ll+1)
	mLSBRecoding(L[:], scalar)
	S := &pointR3{}
	P.SetIdentity()
	for ii := ee - 1; ii >= 0; ii-- {
		P.double()
		for j := 0; j < fxV; j++ {
			dig := L[fxW*dd-j*ee+ii-ee]
			for i := (fxW-1)*dd - j*ee + ii - ee; i >= (2*dd - j*ee + ii - ee); i = i - dd {
				dig = 2*dig + L[i]
			}
			idx := absolute(int32(dig))
			sig := L[dd-j*ee+ii-ee]
			Tabj := &tabSign[fxV-j-1]
			for k := 0; k < fx2w1; k++ {
				S.cmov(&Tabj[k], subtle.ConstantTimeEq(int32(k), idx))
			}
			S.cneg(subtle.ConstantTimeEq(int32(sig), -1))
			P.mixAdd(S)
		}
	}
}

const (
	omegaFix = 7
	omegaVar = 5
)

// doubleMult returns P=mG+nQ.
func (P *pointR1) doubleMult(Q *pointR1, m, n []byte) {
	nafFix := math.OmegaNAF(conv.BytesLe2BigInt(m), omegaFix)
	nafVar := math.OmegaNAF(conv.BytesLe2BigInt(n), omegaVar)

	if len(nafFix) > len(nafVar) {
		nafVar = append(nafVar, make([]int32, len(nafFix)-len(nafVar))...)
	} else if len(nafFix) < len(nafVar) {
		nafFix = append(nafFix, make([]int32, len(nafVar)-len(nafFix))...)
	}

	var TabQ [1 << (omegaVar - 2)]pointR2
	Q.oddMultiples(TabQ[:])
	P.SetIdentity()
	for i := len(nafFix) - 1; i >= 0; i-- {
		P.double()
		// Generator point
		if nafFix[i] != 0 {
			idxM := absolute(nafFix[i]) >> 1
			R := tabVerif[idxM]
			if nafFix[i] < 0 {
				R.neg()
			}
			P.mixAdd(&R)
		}
		// Variable input point
		if nafVar[i] != 0 {
			idxN := absolute(nafVar[i]) >> 1
			S := TabQ[idxN]
			if nafVar[i] < 0 {
				S.neg()
			}
			P.add(&S)
		}
	}
}
