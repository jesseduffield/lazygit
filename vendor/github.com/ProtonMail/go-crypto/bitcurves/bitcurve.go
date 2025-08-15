package bitcurves

// Copyright 2010 The Go Authors. All rights reserved.
// Copyright 2011 ThePiachu. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bitelliptic implements several Koblitz elliptic curves over prime
// fields.

// This package operates, internally, on Jacobian coordinates. For a given
// (x, y) position on the curve, the Jacobian coordinates are (x1, y1, z1)
// where x = x1/z1² and y = y1/z1³. The greatest speedups come when the whole
// calculation can be performed within the transform (as in ScalarMult and
// ScalarBaseMult). But even for Add and Double, it's faster to apply and
// reverse the transform than to operate in affine coordinates.

import (
	"crypto/elliptic"
	"io"
	"math/big"
	"sync"
)

// A BitCurve represents a Koblitz Curve with a=0.
// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw.html
type BitCurve struct {
	Name    string
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the BitCurve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
}

// Params returns the parameters of the given BitCurve (see BitCurve struct)
func (bitCurve *BitCurve) Params() (cp *elliptic.CurveParams) {
	cp = new(elliptic.CurveParams)
	cp.Name = bitCurve.Name
	cp.P = bitCurve.P
	cp.N = bitCurve.N
	cp.Gx = bitCurve.Gx
	cp.Gy = bitCurve.Gy
	cp.BitSize = bitCurve.BitSize
	return cp
}

// IsOnCurve returns true if the given (x,y) lies on the BitCurve.
func (bitCurve *BitCurve) IsOnCurve(x, y *big.Int) bool {
	// y² = x³ + b
	y2 := new(big.Int).Mul(y, y) //y²
	y2.Mod(y2, bitCurve.P)       //y²%P

	x3 := new(big.Int).Mul(x, x) //x²
	x3.Mul(x3, x)                //x³

	x3.Add(x3, bitCurve.B) //x³+B
	x3.Mod(x3, bitCurve.P) //(x³+B)%P

	return x3.Cmp(y2) == 0
}

// affineFromJacobian reverses the Jacobian transform. See the comment at the
// top of the file.
func (bitCurve *BitCurve) affineFromJacobian(x, y, z *big.Int) (xOut, yOut *big.Int) {
	if z.Cmp(big.NewInt(0)) == 0 {
		panic("bitcurve: Can't convert to affine with Jacobian Z = 0")
	}
	// x = YZ^2 mod P
	zinv := new(big.Int).ModInverse(z, bitCurve.P)
	zinvsq := new(big.Int).Mul(zinv, zinv)

	xOut = new(big.Int).Mul(x, zinvsq)
	xOut.Mod(xOut, bitCurve.P)
	// y = YZ^3 mod P
	zinvsq.Mul(zinvsq, zinv)
	yOut = new(big.Int).Mul(y, zinvsq)
	yOut.Mod(yOut, bitCurve.P)
	return xOut, yOut
}

// Add returns the sum of (x1,y1) and (x2,y2)
func (bitCurve *BitCurve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	z := new(big.Int).SetInt64(1)
	x, y, z := bitCurve.addJacobian(x1, y1, z, x2, y2, z)
	return bitCurve.affineFromJacobian(x, y, z)
}

// addJacobian takes two points in Jacobian coordinates, (x1, y1, z1) and
// (x2, y2, z2) and returns their sum, also in Jacobian form.
func (bitCurve *BitCurve) addJacobian(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1 := new(big.Int).Mul(z1, z1)
	z1z1.Mod(z1z1, bitCurve.P)
	z2z2 := new(big.Int).Mul(z2, z2)
	z2z2.Mod(z2z2, bitCurve.P)

	u1 := new(big.Int).Mul(x1, z2z2)
	u1.Mod(u1, bitCurve.P)
	u2 := new(big.Int).Mul(x2, z1z1)
	u2.Mod(u2, bitCurve.P)
	h := new(big.Int).Sub(u2, u1)
	if h.Sign() == -1 {
		h.Add(h, bitCurve.P)
	}
	i := new(big.Int).Lsh(h, 1)
	i.Mul(i, i)
	j := new(big.Int).Mul(h, i)

	s1 := new(big.Int).Mul(y1, z2)
	s1.Mul(s1, z2z2)
	s1.Mod(s1, bitCurve.P)
	s2 := new(big.Int).Mul(y2, z1)
	s2.Mul(s2, z1z1)
	s2.Mod(s2, bitCurve.P)
	r := new(big.Int).Sub(s2, s1)
	if r.Sign() == -1 {
		r.Add(r, bitCurve.P)
	}
	r.Lsh(r, 1)
	v := new(big.Int).Mul(u1, i)

	x3 := new(big.Int).Set(r)
	x3.Mul(x3, x3)
	x3.Sub(x3, j)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	x3.Mod(x3, bitCurve.P)

	y3 := new(big.Int).Set(r)
	v.Sub(v, x3)
	y3.Mul(y3, v)
	s1.Mul(s1, j)
	s1.Lsh(s1, 1)
	y3.Sub(y3, s1)
	y3.Mod(y3, bitCurve.P)

	z3 := new(big.Int).Add(z1, z2)
	z3.Mul(z3, z3)
	z3.Sub(z3, z1z1)
	if z3.Sign() == -1 {
		z3.Add(z3, bitCurve.P)
	}
	z3.Sub(z3, z2z2)
	if z3.Sign() == -1 {
		z3.Add(z3, bitCurve.P)
	}
	z3.Mul(z3, h)
	z3.Mod(z3, bitCurve.P)

	return x3, y3, z3
}

// Double returns 2*(x,y)
func (bitCurve *BitCurve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	z1 := new(big.Int).SetInt64(1)
	return bitCurve.affineFromJacobian(bitCurve.doubleJacobian(x1, y1, z1))
}

// doubleJacobian takes a point in Jacobian coordinates, (x, y, z), and
// returns its double, also in Jacobian form.
func (bitCurve *BitCurve) doubleJacobian(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l

	a := new(big.Int).Mul(x, x) //X1²
	b := new(big.Int).Mul(y, y) //Y1²
	c := new(big.Int).Mul(b, b) //B²

	d := new(big.Int).Add(x, b) //X1+B
	d.Mul(d, d)                 //(X1+B)²
	d.Sub(d, a)                 //(X1+B)²-A
	d.Sub(d, c)                 //(X1+B)²-A-C
	d.Mul(d, big.NewInt(2))     //2*((X1+B)²-A-C)

	e := new(big.Int).Mul(big.NewInt(3), a) //3*A
	f := new(big.Int).Mul(e, e)             //E²

	x3 := new(big.Int).Mul(big.NewInt(2), d) //2*D
	x3.Sub(f, x3)                            //F-2*D
	x3.Mod(x3, bitCurve.P)

	y3 := new(big.Int).Sub(d, x3)                  //D-X3
	y3.Mul(e, y3)                                  //E*(D-X3)
	y3.Sub(y3, new(big.Int).Mul(big.NewInt(8), c)) //E*(D-X3)-8*C
	y3.Mod(y3, bitCurve.P)

	z3 := new(big.Int).Mul(y, z) //Y1*Z1
	z3.Mul(big.NewInt(2), z3)    //3*Y1*Z1
	z3.Mod(z3, bitCurve.P)

	return x3, y3, z3
}

// TODO: double check if it is okay
// ScalarMult returns k*(Bx,By) where k is a number in big-endian form.
func (bitCurve *BitCurve) ScalarMult(Bx, By *big.Int, k []byte) (*big.Int, *big.Int) {
	// We have a slight problem in that the identity of the group (the
	// point at infinity) cannot be represented in (x, y) form on a finite
	// machine. Thus the standard add/double algorithm has to be tweaked
	// slightly: our initial state is not the identity, but x, and we
	// ignore the first true bit in |k|.  If we don't find any true bits in
	// |k|, then we return nil, nil, because we cannot return the identity
	// element.

	Bz := new(big.Int).SetInt64(1)
	x := Bx
	y := By
	z := Bz

	seenFirstTrue := false
	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			if seenFirstTrue {
				x, y, z = bitCurve.doubleJacobian(x, y, z)
			}
			if byte&0x80 == 0x80 {
				if !seenFirstTrue {
					seenFirstTrue = true
				} else {
					x, y, z = bitCurve.addJacobian(Bx, By, Bz, x, y, z)
				}
			}
			byte <<= 1
		}
	}

	if !seenFirstTrue {
		return nil, nil
	}

	return bitCurve.affineFromJacobian(x, y, z)
}

// ScalarBaseMult returns k*G, where G is the base point of the group and k is
// an integer in big-endian form.
func (bitCurve *BitCurve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	return bitCurve.ScalarMult(bitCurve.Gx, bitCurve.Gy, k)
}

var mask = []byte{0xff, 0x1, 0x3, 0x7, 0xf, 0x1f, 0x3f, 0x7f}

// TODO: double check if it is okay
// GenerateKey returns a public/private key pair. The private key is generated
// using the given reader, which must return random data.
func (bitCurve *BitCurve) GenerateKey(rand io.Reader) (priv []byte, x, y *big.Int, err error) {
	byteLen := (bitCurve.BitSize + 7) >> 3
	priv = make([]byte, byteLen)

	for x == nil {
		_, err = io.ReadFull(rand, priv)
		if err != nil {
			return
		}
		// We have to mask off any excess bits in the case that the size of the
		// underlying field is not a whole number of bytes.
		priv[0] &= mask[bitCurve.BitSize%8]
		// This is because, in tests, rand will return all zeros and we don't
		// want to get the point at infinity and loop forever.
		priv[1] ^= 0x42
		x, y = bitCurve.ScalarBaseMult(priv)
	}
	return
}

// Marshal converts a point into the form specified in section 4.3.6 of ANSI
// X9.62.
func (bitCurve *BitCurve) Marshal(x, y *big.Int) []byte {
	byteLen := (bitCurve.BitSize + 7) >> 3

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point

	xBytes := x.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	yBytes := y.Bytes()
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)
	return ret
}

// Unmarshal converts a point, serialised by Marshal, into an x, y pair. On
// error, x = nil.
func (bitCurve *BitCurve) Unmarshal(data []byte) (x, y *big.Int) {
	byteLen := (bitCurve.BitSize + 7) >> 3
	if len(data) != 1+2*byteLen {
		return
	}
	if data[0] != 4 { // uncompressed form
		return
	}
	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
	return
}

//curve parameters taken from:
//http://www.secg.org/collateral/sec2_final.pdf

var initonce sync.Once
var secp160k1 *BitCurve
var secp192k1 *BitCurve
var secp224k1 *BitCurve
var secp256k1 *BitCurve

func initAll() {
	initS160()
	initS192()
	initS224()
	initS256()
}

func initS160() {
	// See SEC 2 section 2.4.1
	secp160k1 = new(BitCurve)
	secp160k1.Name = "secp160k1"
	secp160k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFAC73", 16)
	secp160k1.N, _ = new(big.Int).SetString("0100000000000000000001B8FA16DFAB9ACA16B6B3", 16)
	secp160k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000007", 16)
	secp160k1.Gx, _ = new(big.Int).SetString("3B4C382CE37AA192A4019E763036F4F5DD4D7EBB", 16)
	secp160k1.Gy, _ = new(big.Int).SetString("938CF935318FDCED6BC28286531733C3F03C4FEE", 16)
	secp160k1.BitSize = 160
}

func initS192() {
	// See SEC 2 section 2.5.1
	secp192k1 = new(BitCurve)
	secp192k1.Name = "secp192k1"
	secp192k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFEE37", 16)
	secp192k1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFE26F2FC170F69466A74DEFD8D", 16)
	secp192k1.B, _ = new(big.Int).SetString("000000000000000000000000000000000000000000000003", 16)
	secp192k1.Gx, _ = new(big.Int).SetString("DB4FF10EC057E9AE26B07D0280B7F4341DA5D1B1EAE06C7D", 16)
	secp192k1.Gy, _ = new(big.Int).SetString("9B2F2F6D9C5628A7844163D015BE86344082AA88D95E2F9D", 16)
	secp192k1.BitSize = 192
}

func initS224() {
	// See SEC 2 section 2.6.1
	secp224k1 = new(BitCurve)
	secp224k1.Name = "secp224k1"
	secp224k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFE56D", 16)
	secp224k1.N, _ = new(big.Int).SetString("010000000000000000000000000001DCE8D2EC6184CAF0A971769FB1F7", 16)
	secp224k1.B, _ = new(big.Int).SetString("00000000000000000000000000000000000000000000000000000005", 16)
	secp224k1.Gx, _ = new(big.Int).SetString("A1455B334DF099DF30FC28A169A467E9E47075A90F7E650EB6B7A45C", 16)
	secp224k1.Gy, _ = new(big.Int).SetString("7E089FED7FBA344282CAFBD6F7E319F7C0B0BD59E2CA4BDB556D61A5", 16)
	secp224k1.BitSize = 224
}

func initS256() {
	// See SEC 2 section 2.7.1
	secp256k1 = new(BitCurve)
	secp256k1.Name = "secp256k1"
	secp256k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	secp256k1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	secp256k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	secp256k1.Gx, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	secp256k1.Gy, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	secp256k1.BitSize = 256
}

// S160 returns a BitCurve which implements secp160k1 (see SEC 2 section 2.4.1)
func S160() *BitCurve {
	initonce.Do(initAll)
	return secp160k1
}

// S192 returns a BitCurve which implements secp192k1 (see SEC 2 section 2.5.1)
func S192() *BitCurve {
	initonce.Do(initAll)
	return secp192k1
}

// S224 returns a BitCurve which implements secp224k1 (see SEC 2 section 2.6.1)
func S224() *BitCurve {
	initonce.Do(initAll)
	return secp224k1
}

// S256 returns a BitCurve which implements bitcurves (see SEC 2 section 2.7.1)
func S256() *BitCurve {
	initonce.Do(initAll)
	return secp256k1
}
