// Package goldilocks provides elliptic curve operations over the goldilocks curve.
package goldilocks

import fp "github.com/cloudflare/circl/math/fp448"

// Curve is the Goldilocks curve x^2+y^2=z^2-39081x^2y^2.
type Curve struct{}

// Identity returns the identity point.
func (Curve) Identity() *Point {
	return &Point{
		y: fp.One(),
		z: fp.One(),
	}
}

// IsOnCurve returns true if the point lies on the curve.
func (Curve) IsOnCurve(P *Point) bool {
	x2, y2, t, t2, z2 := &fp.Elt{}, &fp.Elt{}, &fp.Elt{}, &fp.Elt{}, &fp.Elt{}
	rhs, lhs := &fp.Elt{}, &fp.Elt{}
	// Check z != 0
	eq0 := !fp.IsZero(&P.z)

	fp.Mul(t, &P.ta, &P.tb)  // t = ta*tb
	fp.Sqr(x2, &P.x)         // x^2
	fp.Sqr(y2, &P.y)         // y^2
	fp.Sqr(z2, &P.z)         // z^2
	fp.Sqr(t2, t)            // t^2
	fp.Add(lhs, x2, y2)      // x^2 + y^2
	fp.Mul(rhs, t2, &paramD) // dt^2
	fp.Add(rhs, rhs, z2)     // z^2 + dt^2
	fp.Sub(lhs, lhs, rhs)    // x^2 + y^2 - (z^2 + dt^2)
	eq1 := fp.IsZero(lhs)

	fp.Mul(lhs, &P.x, &P.y) // xy
	fp.Mul(rhs, t, &P.z)    // tz
	fp.Sub(lhs, lhs, rhs)   // xy - tz
	eq2 := fp.IsZero(lhs)

	return eq0 && eq1 && eq2
}

// Generator returns the generator point.
func (Curve) Generator() *Point {
	return &Point{
		x:  genX,
		y:  genY,
		z:  fp.One(),
		ta: genX,
		tb: genY,
	}
}

// Order returns the number of points in the prime subgroup.
func (Curve) Order() Scalar { return order }

// Double returns 2P.
func (Curve) Double(P *Point) *Point { R := *P; R.Double(); return &R }

// Add returns P+Q.
func (Curve) Add(P, Q *Point) *Point { R := *P; R.Add(Q); return &R }

// ScalarMult returns kP. This function runs in constant time.
func (e Curve) ScalarMult(k *Scalar, P *Point) *Point {
	k4 := &Scalar{}
	k4.divBy4(k)
	return e.pull(twistCurve{}.ScalarMult(k4, e.push(P)))
}

// ScalarBaseMult returns kG where G is the generator point. This function runs in constant time.
func (e Curve) ScalarBaseMult(k *Scalar) *Point {
	k4 := &Scalar{}
	k4.divBy4(k)
	return e.pull(twistCurve{}.ScalarBaseMult(k4))
}

// CombinedMult returns mG+nP, where G is the generator point. This function is non-constant time.
func (e Curve) CombinedMult(m, n *Scalar, P *Point) *Point {
	m4 := &Scalar{}
	n4 := &Scalar{}
	m4.divBy4(m)
	n4.divBy4(n)
	return e.pull(twistCurve{}.CombinedMult(m4, n4, twistCurve{}.pull(P)))
}
