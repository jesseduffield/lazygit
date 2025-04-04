//go:build !amd64 || purego
// +build !amd64 purego

package fp25519

func cmov(x, y *Elt, n uint)  { cmovGeneric(x, y, n) }
func cswap(x, y *Elt, n uint) { cswapGeneric(x, y, n) }
func add(z, x, y *Elt)        { addGeneric(z, x, y) }
func sub(z, x, y *Elt)        { subGeneric(z, x, y) }
func addsub(x, y *Elt)        { addsubGeneric(x, y) }
func mul(z, x, y *Elt)        { mulGeneric(z, x, y) }
func sqr(z, x *Elt)           { sqrGeneric(z, x) }
func modp(z *Elt)             { modpGeneric(z) }
