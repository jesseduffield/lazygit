package ed25519

import fp "github.com/cloudflare/circl/math/fp25519"

type (
	pointR1 struct{ x, y, z, ta, tb fp.Elt }
	pointR2 struct {
		pointR3
		z2 fp.Elt
	}
)
type pointR3 struct{ addYX, subYX, dt2 fp.Elt }

func (P *pointR1) neg() {
	fp.Neg(&P.x, &P.x)
	fp.Neg(&P.ta, &P.ta)
}

func (P *pointR1) SetIdentity() {
	P.x = fp.Elt{}
	fp.SetOne(&P.y)
	fp.SetOne(&P.z)
	P.ta = fp.Elt{}
	P.tb = fp.Elt{}
}

func (P *pointR1) toAffine() {
	fp.Inv(&P.z, &P.z)
	fp.Mul(&P.x, &P.x, &P.z)
	fp.Mul(&P.y, &P.y, &P.z)
	fp.Modp(&P.x)
	fp.Modp(&P.y)
	fp.SetOne(&P.z)
	P.ta = P.x
	P.tb = P.y
}

func (P *pointR1) ToBytes(k []byte) error {
	P.toAffine()
	var x [fp.Size]byte
	err := fp.ToBytes(k[:fp.Size], &P.y)
	if err != nil {
		return err
	}
	err = fp.ToBytes(x[:], &P.x)
	if err != nil {
		return err
	}
	b := x[0] & 1
	k[paramB-1] = k[paramB-1] | (b << 7)
	return nil
}

func (P *pointR1) FromBytes(k []byte) bool {
	if len(k) != paramB {
		panic("wrong size")
	}
	signX := k[paramB-1] >> 7
	copy(P.y[:], k[:fp.Size])
	P.y[fp.Size-1] &= 0x7F
	p := fp.P()
	if !isLessThan(P.y[:], p[:]) {
		return false
	}

	one, u, v := &fp.Elt{}, &fp.Elt{}, &fp.Elt{}
	fp.SetOne(one)
	fp.Sqr(u, &P.y)                // u = y^2
	fp.Mul(v, u, &paramD)          // v = dy^2
	fp.Sub(u, u, one)              // u = y^2-1
	fp.Add(v, v, one)              // v = dy^2+1
	isQR := fp.InvSqrt(&P.x, u, v) // x = sqrt(u/v)
	if !isQR {
		return false
	}
	fp.Modp(&P.x) // x = x mod p
	if fp.IsZero(&P.x) && signX == 1 {
		return false
	}
	if signX != (P.x[0] & 1) {
		fp.Neg(&P.x, &P.x)
	}
	P.ta = P.x
	P.tb = P.y
	fp.SetOne(&P.z)
	return true
}

// double calculates 2P for curves with A=-1.
func (P *pointR1) double() {
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

func (P *pointR1) mixAdd(Q *pointR3) {
	fp.Add(&P.z, &P.z, &P.z) // D = 2*z1
	P.coreAddition(Q)
}

func (P *pointR1) add(Q *pointR2) {
	fp.Mul(&P.z, &P.z, &Q.z2) // D = 2*z1*z2
	P.coreAddition(&Q.pointR3)
}

// coreAddition calculates P=P+Q for curves with A=-1.
func (P *pointR1) coreAddition(Q *pointR3) {
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

func (P *pointR1) oddMultiples(T []pointR2) {
	var R pointR2
	n := len(T)
	T[0].fromR1(P)
	_2P := *P
	_2P.double()
	R.fromR1(&_2P)
	for i := 1; i < n; i++ {
		P.add(&R)
		T[i].fromR1(P)
	}
}

func (P *pointR1) isEqual(Q *pointR1) bool {
	l, r := &fp.Elt{}, &fp.Elt{}
	fp.Mul(l, &P.x, &Q.z)
	fp.Mul(r, &Q.x, &P.z)
	fp.Sub(l, l, r)
	b := fp.IsZero(l)
	fp.Mul(l, &P.y, &Q.z)
	fp.Mul(r, &Q.y, &P.z)
	fp.Sub(l, l, r)
	b = b && fp.IsZero(l)
	fp.Mul(l, &P.ta, &P.tb)
	fp.Mul(l, l, &Q.z)
	fp.Mul(r, &Q.ta, &Q.tb)
	fp.Mul(r, r, &P.z)
	fp.Sub(l, l, r)
	b = b && fp.IsZero(l)
	return b
}

func (P *pointR3) neg() {
	P.addYX, P.subYX = P.subYX, P.addYX
	fp.Neg(&P.dt2, &P.dt2)
}

func (P *pointR2) fromR1(Q *pointR1) {
	fp.Add(&P.addYX, &Q.y, &Q.x)
	fp.Sub(&P.subYX, &Q.y, &Q.x)
	fp.Mul(&P.dt2, &Q.ta, &Q.tb)
	fp.Mul(&P.dt2, &P.dt2, &paramD)
	fp.Add(&P.dt2, &P.dt2, &P.dt2)
	fp.Add(&P.z2, &Q.z, &Q.z)
}

func (P *pointR3) cneg(b int) {
	t := &fp.Elt{}
	fp.Cswap(&P.addYX, &P.subYX, uint(b))
	fp.Neg(t, &P.dt2)
	fp.Cmov(&P.dt2, t, uint(b))
}

func (P *pointR3) cmov(Q *pointR3, b int) {
	fp.Cmov(&P.addYX, &Q.addYX, uint(b))
	fp.Cmov(&P.subYX, &Q.subYX, uint(b))
	fp.Cmov(&P.dt2, &Q.dt2, uint(b))
}
