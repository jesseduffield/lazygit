// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	_ "crypto/sha512" // need for EC keys
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
)

// jwsEncodeJSON signs claimset using provided key and a nonce.
// The result is serialized in JSON format.
// See https://tools.ietf.org/html/rfc7515#section-7.
func jwsEncodeJSON(claimset interface{}, key crypto.Signer, nonce string) ([]byte, error) {
	jwk, err := jwkEncode(key.Public())
	if err != nil {
		return nil, err
	}
	alg, sha := jwsHasher(key.Public())
	if alg == "" || !sha.Available() {
		return nil, ErrUnsupportedKey
	}
	phead := fmt.Sprintf(`{"alg":%q,"jwk":%s,"nonce":%q}`, alg, jwk, nonce)
	phead = base64.RawURLEncoding.EncodeToString([]byte(phead))
	cs, err := json.Marshal(claimset)
	if err != nil {
		return nil, err
	}
	payload := base64.RawURLEncoding.EncodeToString(cs)
	hash := sha.New()
	hash.Write([]byte(phead + "." + payload))
	sig, err := jwsSign(key, sha, hash.Sum(nil))
	if err != nil {
		return nil, err
	}

	enc := struct {
		Protected string `json:"protected"`
		Payload   string `json:"payload"`
		Sig       string `json:"signature"`
	}{
		Protected: phead,
		Payload:   payload,
		Sig:       base64.RawURLEncoding.EncodeToString(sig),
	}
	return json.Marshal(&enc)
}

// jwkEncode encodes public part of an RSA or ECDSA key into a JWK.
// The result is also suitable for creating a JWK thumbprint.
// https://tools.ietf.org/html/rfc7517
func jwkEncode(pub crypto.PublicKey) (string, error) {
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		// https://tools.ietf.org/html/rfc7518#section-6.3.1
		n := pub.N
		e := big.NewInt(int64(pub.E))
		// Field order is important.
		// See https://tools.ietf.org/html/rfc7638#section-3.3 for details.
		return fmt.Sprintf(`{"e":"%s","kty":"RSA","n":"%s"}`,
			base64.RawURLEncoding.EncodeToString(e.Bytes()),
			base64.RawURLEncoding.EncodeToString(n.Bytes()),
		), nil
	case *ecdsa.PublicKey:
		// https://tools.ietf.org/html/rfc7518#section-6.2.1
		p := pub.Curve.Params()
		n := p.BitSize / 8
		if p.BitSize%8 != 0 {
			n++
		}
		x := pub.X.Bytes()
		if n > len(x) {
			x = append(make([]byte, n-len(x)), x...)
		}
		y := pub.Y.Bytes()
		if n > len(y) {
			y = append(make([]byte, n-len(y)), y...)
		}
		// Field order is important.
		// See https://tools.ietf.org/html/rfc7638#section-3.3 for details.
		return fmt.Sprintf(`{"crv":"%s","kty":"EC","x":"%s","y":"%s"}`,
			p.Name,
			base64.RawURLEncoding.EncodeToString(x),
			base64.RawURLEncoding.EncodeToString(y),
		), nil
	}
	return "", ErrUnsupportedKey
}

// jwsSign signs the digest using the given key.
// The hash is unused for ECDSA keys.
//
// Note: non-stdlib crypto.Signer implementations are expected to return
// the signature in the format as specified in RFC7518.
// See https://tools.ietf.org/html/rfc7518 for more details.
func jwsSign(key crypto.Signer, hash crypto.Hash, digest []byte) ([]byte, error) {
	if key, ok := key.(*ecdsa.PrivateKey); ok {
		// The key.Sign method of ecdsa returns ASN1-encoded signature.
		// So, we use the package Sign function instead
		// to get R and S values directly and format the result accordingly.
		r, s, err := ecdsa.Sign(rand.Reader, key, digest)
		if err != nil {
			return nil, err
		}
		rb, sb := r.Bytes(), s.Bytes()
		size := key.Params().BitSize / 8
		if size%8 > 0 {
			size++
		}
		sig := make([]byte, size*2)
		copy(sig[size-len(rb):], rb)
		copy(sig[size*2-len(sb):], sb)
		return sig, nil
	}
	return key.Sign(rand.Reader, digest, hash)
}

// jwsHasher indicates suitable JWS algorithm name and a hash function
// to use for signing a digest with the provided key.
// It returns ("", 0) if the key is not supported.
func jwsHasher(pub crypto.PublicKey) (string, crypto.Hash) {
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return "RS256", crypto.SHA256
	case *ecdsa.PublicKey:
		switch pub.Params().Name {
		case "P-256":
			return "ES256", crypto.SHA256
		case "P-384":
			return "ES384", crypto.SHA384
		case "P-521":
			return "ES512", crypto.SHA512
		}
	}
	return "", 0
}

// JWKThumbprint creates a JWK thumbprint out of pub
// as specified in https://tools.ietf.org/html/rfc7638.
func JWKThumbprint(pub crypto.PublicKey) (string, error) {
	jwk, err := jwkEncode(pub)
	if err != nil {
		return "", err
	}
	b := sha256.Sum256([]byte(jwk))
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}
