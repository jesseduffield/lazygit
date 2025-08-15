package mlsbset

import "fmt"

// Power is a valid exponent produced by the MLSBSet encoding algorithm.
type Power struct {
	set Encoder // parameters of code.
	s   []int32 // set of signs.
	b   []int32 // set of digits.
	c   int     // carry is {0,1}.
}

// Exp is calculates x^k, where x is a predetermined element of a group G.
func (p *Power) Exp(G Group) EltG {
	a, b := G.Identity(), G.NewEltP()
	for e := int(p.set.p.E - 1); e >= 0; e-- {
		G.Sqr(a)
		for v := uint(0); v < p.set.p.V; v++ {
			sgnElt, idElt := p.Digit(v, uint(e))
			G.Lookup(b, v, sgnElt, idElt)
			G.Mul(a, b)
		}
	}
	if p.set.IsExtended() && p.c == 1 {
		G.Mul(a, G.ExtendedEltP())
	}
	return a
}

// Digit returns the (v,e)-th digit and its sign.
func (p *Power) Digit(v, e uint) (sgn, dig int32) {
	sgn = p.bit(0, v, e)
	dig = 0
	for i := p.set.p.W - 1; i > 0; i-- {
		dig = 2*dig + p.bit(i, v, e)
	}
	mask := dig >> 31
	dig = (dig + mask) ^ mask
	return sgn, dig
}

// bit returns the (w,v,e)-th bit of the code.
func (p *Power) bit(w, v, e uint) int32 {
	if !(w < p.set.p.W &&
		v < p.set.p.V &&
		e < p.set.p.E) {
		panic(fmt.Errorf("indexes outside (%v,%v,%v)", w, v, e))
	}
	if w == 0 {
		return p.s[p.set.p.E*v+e]
	}
	return p.b[p.set.p.D*(w-1)+p.set.p.E*v+e]
}

func (p *Power) String() string {
	dig := ""
	for j := uint(0); j < p.set.p.V; j++ {
		for i := uint(0); i < p.set.p.E; i++ {
			s, d := p.Digit(j, i)
			dig += fmt.Sprintf("(%2v,%2v) = %+2v %+2v\n", j, i, s, d)
		}
	}
	return fmt.Sprintf("len: %v\ncarry: %v\ndigits:\n%v", len(p.b)+len(p.s), p.c, dig)
}
