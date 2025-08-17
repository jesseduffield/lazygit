// Package ecc implements a generic interface for ECDH, ECDSA, and EdDSA.
package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"io"
	"math/big"
)

type genericCurve struct {
	Curve elliptic.Curve
}

func NewGenericCurve(c elliptic.Curve) *genericCurve {
	return &genericCurve{
		Curve: c,
	}
}

func (c *genericCurve) GetCurveName() string {
	return c.Curve.Params().Name
}

func (c *genericCurve) MarshalBytePoint(point []byte) []byte {
	return point
}

func (c *genericCurve) UnmarshalBytePoint(point []byte) []byte {
	return point
}

func (c *genericCurve) MarshalIntegerPoint(x, y *big.Int) []byte {
	return elliptic.Marshal(c.Curve, x, y)
}

func (c *genericCurve) UnmarshalIntegerPoint(point []byte) (x, y *big.Int) {
	return elliptic.Unmarshal(c.Curve, point)
}

func (c *genericCurve) MarshalByteSecret(d []byte) []byte {
	return d
}

func (c *genericCurve) UnmarshalByteSecret(d []byte) []byte {
	return d
}

func (c *genericCurve) MarshalIntegerSecret(d *big.Int) []byte {
	return d.Bytes()
}

func (c *genericCurve) UnmarshalIntegerSecret(d []byte) *big.Int {
	return new(big.Int).SetBytes(d)
}

func (c *genericCurve) GenerateECDH(rand io.Reader) (point, secret []byte, err error) {
	secret, x, y, err := elliptic.GenerateKey(c.Curve, rand)
	if err != nil {
		return nil, nil, err
	}

	point = elliptic.Marshal(c.Curve, x, y)
	return point, secret, nil
}

func (c *genericCurve) GenerateECDSA(rand io.Reader) (x, y, secret *big.Int, err error) {
	priv, err := ecdsa.GenerateKey(c.Curve, rand)
	if err != nil {
		return
	}

	return priv.X, priv.Y, priv.D, nil
}

func (c *genericCurve) Encaps(rand io.Reader, point []byte) (ephemeral, sharedSecret []byte, err error) {
	xP, yP := elliptic.Unmarshal(c.Curve, point)
	if xP == nil {
		panic("invalid point")
	}

	d, x, y, err := elliptic.GenerateKey(c.Curve, rand)
	if err != nil {
		return nil, nil, err
	}

	vsG := elliptic.Marshal(c.Curve, x, y)
	zbBig, _ := c.Curve.ScalarMult(xP, yP, d)

	byteLen := (c.Curve.Params().BitSize + 7) >> 3
	zb := make([]byte, byteLen)
	zbBytes := zbBig.Bytes()
	copy(zb[byteLen-len(zbBytes):], zbBytes)

	return vsG, zb, nil
}

func (c *genericCurve) Decaps(ephemeral, secret []byte) (sharedSecret []byte, err error) {
	x, y := elliptic.Unmarshal(c.Curve, ephemeral)
	zbBig, _ := c.Curve.ScalarMult(x, y, secret)
	byteLen := (c.Curve.Params().BitSize + 7) >> 3
	zb := make([]byte, byteLen)
	zbBytes := zbBig.Bytes()
	copy(zb[byteLen-len(zbBytes):], zbBytes)

	return zb, nil
}

func (c *genericCurve) Sign(rand io.Reader, x, y, d *big.Int, hash []byte) (r, s *big.Int, err error) {
	priv := &ecdsa.PrivateKey{D: d, PublicKey: ecdsa.PublicKey{X: x, Y: y, Curve: c.Curve}}
	return ecdsa.Sign(rand, priv, hash)
}

func (c *genericCurve) Verify(x, y *big.Int, hash []byte, r, s *big.Int) bool {
	pub := &ecdsa.PublicKey{X: x, Y: y, Curve: c.Curve}
	return ecdsa.Verify(pub, hash, r, s)
}

func (c *genericCurve) validate(xP, yP *big.Int, secret []byte) error {
	// the public point should not be at infinity (0,0)
	zero := new(big.Int)
	if xP.Cmp(zero) == 0 && yP.Cmp(zero) == 0 {
		return errors.KeyInvalidError(fmt.Sprintf("ecc (%s): infinity point", c.Curve.Params().Name))
	}

	// re-derive the public point Q' = (X,Y) = dG
	// to compare to declared Q in public key
	expectedX, expectedY := c.Curve.ScalarBaseMult(secret)
	if xP.Cmp(expectedX) != 0 || yP.Cmp(expectedY) != 0 {
		return errors.KeyInvalidError(fmt.Sprintf("ecc (%s): invalid point", c.Curve.Params().Name))
	}

	return nil
}

func (c *genericCurve) ValidateECDSA(xP, yP *big.Int, secret []byte) error {
	return c.validate(xP, yP, secret)
}

func (c *genericCurve) ValidateECDH(point []byte, secret []byte) error {
	xP, yP := elliptic.Unmarshal(c.Curve, point)
	if xP == nil {
		return errors.KeyInvalidError(fmt.Sprintf("ecc (%s): invalid point", c.Curve.Params().Name))
	}

	return c.validate(xP, yP, secret)
}
