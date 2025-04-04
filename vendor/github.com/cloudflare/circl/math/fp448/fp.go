// Package fp448 provides prime field arithmetic over GF(2^448-2^224-1).
package fp448

import (
	"errors"

	"github.com/cloudflare/circl/internal/conv"
)

// Size in bytes of an element.
const Size = 56

// Elt is a prime field element.
type Elt [Size]byte

func (e Elt) String() string { return conv.BytesLe2Hex(e[:]) }

// p is the prime modulus 2^448-2^224-1.
var p = Elt{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

// P returns the prime modulus 2^448-2^224-1.
func P() Elt { return p }

// ToBytes stores in b the little-endian byte representation of x.
func ToBytes(b []byte, x *Elt) error {
	if len(b) != Size {
		return errors.New("wrong size")
	}
	Modp(x)
	copy(b, x[:])
	return nil
}

// IsZero returns true if x is equal to 0.
func IsZero(x *Elt) bool { Modp(x); return *x == Elt{} }

// IsOne returns true if x is equal to 1.
func IsOne(x *Elt) bool { Modp(x); return *x == Elt{1} }

// SetOne assigns x=1.
func SetOne(x *Elt) { *x = Elt{1} }

// One returns the 1 element.
func One() (x Elt) { x = Elt{1}; return }

// Neg calculates z = -x.
func Neg(z, x *Elt) { Sub(z, &p, x) }

// Modp ensures that z is between [0,p-1].
func Modp(z *Elt) { Sub(z, z, &p) }

// InvSqrt calculates z = sqrt(x/y) iff x/y is a quadratic-residue. If so,
// isQR = true; otherwise, isQR = false, since x/y is a quadratic non-residue,
// and z = sqrt(-x/y).
func InvSqrt(z, x, y *Elt) (isQR bool) {
	// First note that x^(2(k+1)) = x^(p-1)/2 * x = legendre(x) * x
	// so that's x if x is a quadratic residue and -x otherwise.
	// Next, y^(6k+3) = y^(4k+2) * y^(2k+1) = y^(p-1) * y^((p-1)/2) = legendre(y).
	// So the z we compute satisfies z^2 y = x^(2(k+1)) y^(6k+3) = legendre(x)*legendre(y).
	// Thus if x and y are quadratic residues, then z is indeed sqrt(x/y).
	t0, t1 := &Elt{}, &Elt{}
	Mul(t0, x, y)         // x*y
	Sqr(t1, y)            // y^2
	Mul(t1, t0, t1)       // x*y^3
	powPminus3div4(z, t1) // (x*y^3)^k
	Mul(z, z, t0)         // z = x*y*(x*y^3)^k = x^(k+1) * y^(3k+1)

	// Check if x/y is a quadratic residue
	Sqr(t0, z)     // z^2
	Mul(t0, t0, y) // y*z^2
	Sub(t0, t0, x) // y*z^2-x
	return IsZero(t0)
}

// Inv calculates z = 1/x mod p.
func Inv(z, x *Elt) {
	// Calculates z = x^(4k+1) = x^(p-3+1) = x^(p-2) = x^-1, where k = (p-3)/4.
	t := &Elt{}
	powPminus3div4(t, x) // t = x^k
	Sqr(t, t)            // t = x^2k
	Sqr(t, t)            // t = x^4k
	Mul(z, t, x)         // z = x^(4k+1)
}

// powPminus3div4 calculates z = x^k mod p, where k = (p-3)/4.
func powPminus3div4(z, x *Elt) {
	x0, x1 := &Elt{}, &Elt{}
	Sqr(z, x)
	Mul(z, z, x)
	Sqr(x0, z)
	Mul(x0, x0, x)
	Sqr(z, x0)
	Sqr(z, z)
	Sqr(z, z)
	Mul(z, z, x0)
	Sqr(x1, z)
	for i := 0; i < 5; i++ {
		Sqr(x1, x1)
	}
	Mul(x1, x1, z)
	Sqr(z, x1)
	for i := 0; i < 11; i++ {
		Sqr(z, z)
	}
	Mul(z, z, x1)
	Sqr(z, z)
	Sqr(z, z)
	Sqr(z, z)
	Mul(z, z, x0)
	Sqr(x1, z)
	for i := 0; i < 26; i++ {
		Sqr(x1, x1)
	}
	Mul(x1, x1, z)
	Sqr(z, x1)
	for i := 0; i < 53; i++ {
		Sqr(z, z)
	}
	Mul(z, z, x1)
	Sqr(z, z)
	Sqr(z, z)
	Sqr(z, z)
	Mul(z, z, x0)
	Sqr(x1, z)
	for i := 0; i < 110; i++ {
		Sqr(x1, x1)
	}
	Mul(x1, x1, z)
	Sqr(z, x1)
	Mul(z, z, x)
	for i := 0; i < 223; i++ {
		Sqr(z, z)
	}
	Mul(z, z, x1)
}

// Cmov assigns y to x if n is 1.
func Cmov(x, y *Elt, n uint) { cmov(x, y, n) }

// Cswap interchanges x and y if n is 1.
func Cswap(x, y *Elt, n uint) { cswap(x, y, n) }

// Add calculates z = x+y mod p.
func Add(z, x, y *Elt) { add(z, x, y) }

// Sub calculates z = x-y mod p.
func Sub(z, x, y *Elt) { sub(z, x, y) }

// AddSub calculates (x,y) = (x+y mod p, x-y mod p).
func AddSub(x, y *Elt) { addsub(x, y) }

// Mul calculates z = x*y mod p.
func Mul(z, x, y *Elt) { mul(z, x, y) }

// Sqr calculates z = x^2 mod p.
func Sqr(z, x *Elt) { sqr(z, x) }
