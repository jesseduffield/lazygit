package goldilocks

import (
	"fmt"

	fp "github.com/cloudflare/circl/math/fp448"
)

type twistPoint struct{ x, y, z, ta, tb fp.Elt }

type preTwistPointAffine struct{ addYX, subYX, dt2 fp.Elt }

type preTwistPointProy struct {
	preTwistPointAffine
	z2 fp.Elt
}

func (P *twistPoint) String() string {
	return fmt.Sprintf("x: %v\ny: %v\nz: %v\nta: %v\ntb: %v", P.x, P.y, P.z, P.ta, P.tb)
}

// cneg conditionally negates the point if b=1.
func (P *twistPoint) cneg(b uint) {
	t := &fp.Elt{}
	fp.Neg(t, &P.x)
	fp.Cmov(&P.x, t, b)
	fp.Neg(t, &P.ta)
	fp.Cmov(&P.ta, t, b)
}

// Double updates P with 2P.
func (P *twistPoint) Double() {
	// This is formula (7) from "Twisted Edwards Curves Revisited" by
	// Hisil H., Wong K.KH., Carter G., Dawson E. (2008)
	// https://doi.org/10.1007/978-3-540-89255-7_20
	Px, Py, Pz, Pta, Ptb := &P.x, &P.y, &P.z, &P.ta, &P.tb
	a, b, c, e, f, g, h := Px, Py, Pz, Pta, Px, Py, Ptb
	fp.Add(e, Px, Py) // x+y
	fp.Sqr(a, Px)     // A = x^2
	fp.Sqr(b, Py)     // B = y^2
	fp.Sqr(c, Pz)     // z^2
	fp.Add(c, c, c)   // C = 2*z^2
	fp.Add(h, a, b)   // H = A+B
	fp.Sqr(e, e)      // (x+y)^2
	fp.Sub(e, e, h)   // E = (x+y)^2-A-B
	fp.Sub(g, b, a)   // G = B-A
	fp.Sub(f, c, g)   // F = C-G
	fp.Mul(Pz, f, g)  // Z = F * G
	fp.Mul(Px, e, f)  // X = E * F
	fp.Mul(Py, g, h)  // Y = G * H, T = E * H
}

// mixAdd calculates P= P+Q, where Q is a precomputed point with Z_Q = 1.
func (P *twistPoint) mixAddZ1(Q *preTwistPointAffine) {
	fp.Add(&P.z, &P.z, &P.z) // D = 2*z1 (z2=1)
	P.coreAddition(Q)
}

// coreAddition calculates P=P+Q for curves with A=-1.
func (P *twistPoint) coreAddition(Q *preTwistPointAffine) {
	// This is the formula following (5) from "Twisted Edwards Curves Revisited" by
	// Hisil H., Wong K.KH., Carter G., Dawson E. (2008)
	// https://doi.org/10.1007/978-3-540-89255-7_20
	Px, Py, Pz, Pta, Ptb := &P.x, &P.y, &P.z, &P.ta, &P.tb
	addYX2, subYX2, dt2 := &Q.addYX, &Q.subYX, &Q.dt2
	a, b, c, d, e, f, g, h := Px, Py, &fp.Elt{}, Pz, Pta, Px, Py, Ptb
	fp.Mul(c, Pta, Ptb)  // t1 = ta*tb
	fp.Sub(h, Py, Px)    // y1-x1
	fp.Add(b, Py, Px)    // y1+x1
	fp.Mul(a, h, subYX2) // A = (y1-x1)*(y2-x2)
	fp.Mul(b, b, addYX2) // B = (y1+x1)*(y2+x2)
	fp.Mul(c, c, dt2)    // C = 2*D*t1*t2
	fp.Sub(e, b, a)      // E = B-A
	fp.Add(h, b, a)      // H = B+A
	fp.Sub(f, d, c)      // F = D-C
	fp.Add(g, d, c)      // G = D+C
	fp.Mul(Pz, f, g)     // Z = F * G
	fp.Mul(Px, e, f)     // X = E * F
	fp.Mul(Py, g, h)     // Y = G * H, T = E * H
}

func (P *preTwistPointAffine) neg() {
	P.addYX, P.subYX = P.subYX, P.addYX
	fp.Neg(&P.dt2, &P.dt2)
}

func (P *preTwistPointAffine) cneg(b int) {
	t := &fp.Elt{}
	fp.Cswap(&P.addYX, &P.subYX, uint(b))
	fp.Neg(t, &P.dt2)
	fp.Cmov(&P.dt2, t, uint(b))
}

func (P *preTwistPointAffine) cmov(Q *preTwistPointAffine, b uint) {
	fp.Cmov(&P.addYX, &Q.addYX, b)
	fp.Cmov(&P.subYX, &Q.subYX, b)
	fp.Cmov(&P.dt2, &Q.dt2, b)
}

// mixAdd calculates P= P+Q, where Q is a precomputed point with Z_Q != 1.
func (P *twistPoint) mixAdd(Q *preTwistPointProy) {
	fp.Mul(&P.z, &P.z, &Q.z2) // D = 2*z1*z2
	P.coreAddition(&Q.preTwistPointAffine)
}

// oddMultiples calculates T[i] = (2*i-1)P for 0 < i < len(T).
func (P *twistPoint) oddMultiples(T []preTwistPointProy) {
	if n := len(T); n > 0 {
		T[0].FromTwistPoint(P)
		_2P := *P
		_2P.Double()
		R := &preTwistPointProy{}
		R.FromTwistPoint(&_2P)
		for i := 1; i < n; i++ {
			P.mixAdd(R)
			T[i].FromTwistPoint(P)
		}
	}
}

// cmov conditionally moves Q into P if b=1.
func (P *preTwistPointProy) cmov(Q *preTwistPointProy, b uint) {
	P.preTwistPointAffine.cmov(&Q.preTwistPointAffine, b)
	fp.Cmov(&P.z2, &Q.z2, b)
}

// FromTwistPoint precomputes some coordinates of Q for missed addition.
func (P *preTwistPointProy) FromTwistPoint(Q *twistPoint) {
	fp.Add(&P.addYX, &Q.y, &Q.x)         // addYX = X + Y
	fp.Sub(&P.subYX, &Q.y, &Q.x)         // subYX = Y - X
	fp.Mul(&P.dt2, &Q.ta, &Q.tb)         // T = ta*tb
	fp.Mul(&P.dt2, &P.dt2, &paramDTwist) // D*T
	fp.Add(&P.dt2, &P.dt2, &P.dt2)       // dt2 = 2*D*T
	fp.Add(&P.z2, &Q.z, &Q.z)            // z2 = 2*Z
}
