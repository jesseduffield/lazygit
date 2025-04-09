package ed448

import (
	"crypto/rand"
	"encoding/asn1"

	"github.com/cloudflare/circl/sign"
)

var sch sign.Scheme = &scheme{}

// Scheme returns a signature interface.
func Scheme() sign.Scheme { return sch }

type scheme struct{}

func (*scheme) Name() string          { return "Ed448" }
func (*scheme) PublicKeySize() int    { return PublicKeySize }
func (*scheme) PrivateKeySize() int   { return PrivateKeySize }
func (*scheme) SignatureSize() int    { return SignatureSize }
func (*scheme) SeedSize() int         { return SeedSize }
func (*scheme) TLSIdentifier() uint   { return 0x0808 }
func (*scheme) SupportsContext() bool { return true }
func (*scheme) Oid() asn1.ObjectIdentifier {
	return asn1.ObjectIdentifier{1, 3, 101, 113}
}

func (*scheme) GenerateKey() (sign.PublicKey, sign.PrivateKey, error) {
	return GenerateKey(rand.Reader)
}

func (*scheme) Sign(
	sk sign.PrivateKey,
	message []byte,
	opts *sign.SignatureOpts,
) []byte {
	priv, ok := sk.(PrivateKey)
	if !ok {
		panic(sign.ErrTypeMismatch)
	}
	ctx := ""
	if opts != nil {
		ctx = opts.Context
	}
	return Sign(priv, message, ctx)
}

func (*scheme) Verify(
	pk sign.PublicKey,
	message, signature []byte,
	opts *sign.SignatureOpts,
) bool {
	pub, ok := pk.(PublicKey)
	if !ok {
		panic(sign.ErrTypeMismatch)
	}
	ctx := ""
	if opts != nil {
		ctx = opts.Context
	}
	return Verify(pub, message, signature, ctx)
}

func (*scheme) DeriveKey(seed []byte) (sign.PublicKey, sign.PrivateKey) {
	privateKey := NewKeyFromSeed(seed)
	publicKey := make(PublicKey, PublicKeySize)
	copy(publicKey, privateKey[SeedSize:])
	return publicKey, privateKey
}

func (*scheme) UnmarshalBinaryPublicKey(buf []byte) (sign.PublicKey, error) {
	if len(buf) < PublicKeySize {
		return nil, sign.ErrPubKeySize
	}
	pub := make(PublicKey, PublicKeySize)
	copy(pub, buf[:PublicKeySize])
	return pub, nil
}

func (*scheme) UnmarshalBinaryPrivateKey(buf []byte) (sign.PrivateKey, error) {
	if len(buf) < PrivateKeySize {
		return nil, sign.ErrPrivKeySize
	}
	priv := make(PrivateKey, PrivateKeySize)
	copy(priv, buf[:PrivateKeySize])
	return priv, nil
}
