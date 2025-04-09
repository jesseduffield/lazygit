package goldilocks

import fp "github.com/cloudflare/circl/math/fp448"

func (Curve) pull(P *twistPoint) *Point      { return twistCurve{}.push(P) }
func (twistCurve) pull(P *Point) *twistPoint { return Curve{}.push(P) }

// push sends a point on the Goldilocks curve to a point on the twist curve.
func (Curve) push(P *Point) *twistPoint {
	Q := &twistPoint{}
	Px, Py, Pz := &P.x, &P.y, &P.z
	a, b, c, d, e, f, g, h := &Q.x, &Q.y, &Q.z, &fp.Elt{}, &Q.ta, &Q.x, &Q.y, &Q.tb
	fp.Add(e, Px, Py)  // x+y
	fp.Sqr(a, Px)      // A = x^2
	fp.Sqr(b, Py)      // B = y^2
	fp.Sqr(c, Pz)      // z^2
	fp.Add(c, c, c)    // C = 2*z^2
	*d = *a            // D = A
	fp.Sqr(e, e)       // (x+y)^2
	fp.Sub(e, e, a)    // (x+y)^2-A
	fp.Sub(e, e, b)    // E = (x+y)^2-A-B
	fp.Add(h, b, d)    // H = B+D
	fp.Sub(g, b, d)    // G = B-D
	fp.Sub(f, c, h)    // F = C-H
	fp.Mul(&Q.z, f, g) // Z = F * G
	fp.Mul(&Q.x, e, f) // X = E * F
	fp.Mul(&Q.y, g, h) // Y = G * H, // T = E * H
	return Q
}

// push sends a point on the twist curve to a point on the Goldilocks curve.
func (twistCurve) push(P *twistPoint) *Point {
	Q := &Point{}
	Px, Py, Pz := &P.x, &P.y, &P.z
	a, b, c, d, e, f, g, h := &Q.x, &Q.y, &Q.z, &fp.Elt{}, &Q.ta, &Q.x, &Q.y, &Q.tb
	fp.Add(e, Px, Py)  // x+y
	fp.Sqr(a, Px)      // A = x^2
	fp.Sqr(b, Py)      // B = y^2
	fp.Sqr(c, Pz)      // z^2
	fp.Add(c, c, c)    // C = 2*z^2
	fp.Neg(d, a)       // D = -A
	fp.Sqr(e, e)       // (x+y)^2
	fp.Sub(e, e, a)    // (x+y)^2-A
	fp.Sub(e, e, b)    // E = (x+y)^2-A-B
	fp.Add(h, b, d)    // H = B+D
	fp.Sub(g, b, d)    // G = B-D
	fp.Sub(f, c, h)    // F = C-H
	fp.Mul(&Q.z, f, g) // Z = F * G
	fp.Mul(&Q.x, e, f) // X = E * F
	fp.Mul(&Q.y, g, h) // Y = G * H, // T = E * H
	return Q
}
