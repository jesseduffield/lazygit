package goldilocks

import (
	"crypto/subtle"

	mlsb "github.com/cloudflare/circl/math/mlsbset"
)

const (
	// MLSBRecoding parameters
	fxT   = 448
	fxV   = 2
	fxW   = 3
	fx2w1 = 1 << (uint(fxW) - 1)
)

// ScalarBaseMult returns kG where G is the generator point.
func (e twistCurve) ScalarBaseMult(k *Scalar) *twistPoint {
	m, err := mlsb.New(fxT, fxV, fxW)
	if err != nil {
		panic(err)
	}
	if m.IsExtended() {
		panic("not extended")
	}

	var isZero int
	if k.IsZero() {
		isZero = 1
	}
	subtle.ConstantTimeCopy(isZero, k[:], order[:])

	minusK := *k
	isEven := 1 - int(k[0]&0x1)
	minusK.Neg()
	subtle.ConstantTimeCopy(isEven, k[:], minusK[:])
	c, err := m.Encode(k[:])
	if err != nil {
		panic(err)
	}

	gP := c.Exp(groupMLSB{})
	P := gP.(*twistPoint)
	P.cneg(uint(isEven))
	return P
}

type groupMLSB struct{}

func (e groupMLSB) ExtendedEltP() mlsb.EltP      { return nil }
func (e groupMLSB) Sqr(x mlsb.EltG)              { x.(*twistPoint).Double() }
func (e groupMLSB) Mul(x mlsb.EltG, y mlsb.EltP) { x.(*twistPoint).mixAddZ1(y.(*preTwistPointAffine)) }
func (e groupMLSB) Identity() mlsb.EltG          { return twistCurve{}.Identity() }
func (e groupMLSB) NewEltP() mlsb.EltP           { return &preTwistPointAffine{} }
func (e groupMLSB) Lookup(a mlsb.EltP, v uint, s, u int32) {
	Tabj := &tabFixMult[v]
	P := a.(*preTwistPointAffine)
	for k := range Tabj {
		P.cmov(&Tabj[k], uint(subtle.ConstantTimeEq(int32(k), u)))
	}
	P.cneg(int(s >> 31))
}
