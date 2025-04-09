// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto/dsa"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	_ "crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/ecdh"
	"github.com/ProtonMail/go-crypto/openpgp/ecdsa"
	"github.com/ProtonMail/go-crypto/openpgp/ed25519"
	"github.com/ProtonMail/go-crypto/openpgp/ed448"
	"github.com/ProtonMail/go-crypto/openpgp/eddsa"
	"github.com/ProtonMail/go-crypto/openpgp/elgamal"
	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/internal/algorithm"
	"github.com/ProtonMail/go-crypto/openpgp/internal/ecc"
	"github.com/ProtonMail/go-crypto/openpgp/internal/encoding"
	"github.com/ProtonMail/go-crypto/openpgp/x25519"
	"github.com/ProtonMail/go-crypto/openpgp/x448"
)

// PublicKey represents an OpenPGP public key. See RFC 4880, section 5.5.2.
type PublicKey struct {
	Version      int
	CreationTime time.Time
	PubKeyAlgo   PublicKeyAlgorithm
	PublicKey    interface{} // *rsa.PublicKey, *dsa.PublicKey, *ecdsa.PublicKey or *eddsa.PublicKey, *x25519.PublicKey, *x448.PublicKey, *ed25519.PublicKey, *ed448.PublicKey
	Fingerprint  []byte
	KeyId        uint64
	IsSubkey     bool

	// RFC 4880 fields
	n, e, p, q, g, y encoding.Field

	// RFC 6637 fields
	// oid contains the OID byte sequence identifying the elliptic curve used
	oid encoding.Field

	// kdf stores key derivation function parameters
	// used for ECDH encryption. See RFC 6637, Section 9.
	kdf encoding.Field
}

// UpgradeToV5 updates the version of the key to v5, and updates all necessary
// fields.
func (pk *PublicKey) UpgradeToV5() {
	pk.Version = 5
	pk.setFingerprintAndKeyId()
}

// UpgradeToV6 updates the version of the key to v6, and updates all necessary
// fields.
func (pk *PublicKey) UpgradeToV6() error {
	pk.Version = 6
	pk.setFingerprintAndKeyId()
	return pk.checkV6Compatibility()
}

// signingKey provides a convenient abstraction over signature verification
// for v3 and v4 public keys.
type signingKey interface {
	SerializeForHash(io.Writer) error
	SerializeSignaturePrefix(io.Writer) error
	serializeWithoutHeaders(io.Writer) error
}

// NewRSAPublicKey returns a PublicKey that wraps the given rsa.PublicKey.
func NewRSAPublicKey(creationTime time.Time, pub *rsa.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoRSA,
		PublicKey:    pub,
		n:            new(encoding.MPI).SetBig(pub.N),
		e:            new(encoding.MPI).SetBig(big.NewInt(int64(pub.E))),
	}

	pk.setFingerprintAndKeyId()
	return pk
}

// NewDSAPublicKey returns a PublicKey that wraps the given dsa.PublicKey.
func NewDSAPublicKey(creationTime time.Time, pub *dsa.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoDSA,
		PublicKey:    pub,
		p:            new(encoding.MPI).SetBig(pub.P),
		q:            new(encoding.MPI).SetBig(pub.Q),
		g:            new(encoding.MPI).SetBig(pub.G),
		y:            new(encoding.MPI).SetBig(pub.Y),
	}

	pk.setFingerprintAndKeyId()
	return pk
}

// NewElGamalPublicKey returns a PublicKey that wraps the given elgamal.PublicKey.
func NewElGamalPublicKey(creationTime time.Time, pub *elgamal.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoElGamal,
		PublicKey:    pub,
		p:            new(encoding.MPI).SetBig(pub.P),
		g:            new(encoding.MPI).SetBig(pub.G),
		y:            new(encoding.MPI).SetBig(pub.Y),
	}

	pk.setFingerprintAndKeyId()
	return pk
}

func NewECDSAPublicKey(creationTime time.Time, pub *ecdsa.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoECDSA,
		PublicKey:    pub,
		p:            encoding.NewMPI(pub.MarshalPoint()),
	}

	curveInfo := ecc.FindByCurve(pub.GetCurve())
	if curveInfo == nil {
		panic("unknown elliptic curve")
	}
	pk.oid = curveInfo.Oid
	pk.setFingerprintAndKeyId()
	return pk
}

func NewECDHPublicKey(creationTime time.Time, pub *ecdh.PublicKey) *PublicKey {
	var pk *PublicKey
	var kdf = encoding.NewOID([]byte{0x1, pub.Hash.Id(), pub.Cipher.Id()})
	pk = &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoECDH,
		PublicKey:    pub,
		p:            encoding.NewMPI(pub.MarshalPoint()),
		kdf:          kdf,
	}

	curveInfo := ecc.FindByCurve(pub.GetCurve())

	if curveInfo == nil {
		panic("unknown elliptic curve")
	}

	pk.oid = curveInfo.Oid
	pk.setFingerprintAndKeyId()
	return pk
}

func NewEdDSAPublicKey(creationTime time.Time, pub *eddsa.PublicKey) *PublicKey {
	curveInfo := ecc.FindByCurve(pub.GetCurve())
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoEdDSA,
		PublicKey:    pub,
		oid:          curveInfo.Oid,
		// Native point format, see draft-koch-eddsa-for-openpgp-04, Appendix B
		p: encoding.NewMPI(pub.MarshalPoint()),
	}

	pk.setFingerprintAndKeyId()
	return pk
}

func NewX25519PublicKey(creationTime time.Time, pub *x25519.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoX25519,
		PublicKey:    pub,
	}

	pk.setFingerprintAndKeyId()
	return pk
}

func NewX448PublicKey(creationTime time.Time, pub *x448.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoX448,
		PublicKey:    pub,
	}

	pk.setFingerprintAndKeyId()
	return pk
}

func NewEd25519PublicKey(creationTime time.Time, pub *ed25519.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoEd25519,
		PublicKey:    pub,
	}

	pk.setFingerprintAndKeyId()
	return pk
}

func NewEd448PublicKey(creationTime time.Time, pub *ed448.PublicKey) *PublicKey {
	pk := &PublicKey{
		Version:      4,
		CreationTime: creationTime,
		PubKeyAlgo:   PubKeyAlgoEd448,
		PublicKey:    pub,
	}

	pk.setFingerprintAndKeyId()
	return pk
}

func (pk *PublicKey) parse(r io.Reader) (err error) {
	// RFC 4880, section 5.5.2
	var buf [6]byte
	_, err = readFull(r, buf[:])
	if err != nil {
		return
	}

	pk.Version = int(buf[0])
	if pk.Version != 4 && pk.Version != 5 && pk.Version != 6 {
		return errors.UnsupportedError("public key version " + strconv.Itoa(int(buf[0])))
	}

	if V5Disabled && pk.Version == 5 {
		return errors.UnsupportedError("support for parsing v5 entities is disabled; build with `-tags v5` if needed")
	}

	if pk.Version >= 5 {
		// Read the four-octet scalar octet count
		// The count is not used in this implementation
		var n [4]byte
		_, err = readFull(r, n[:])
		if err != nil {
			return
		}
	}
	pk.CreationTime = time.Unix(int64(uint32(buf[1])<<24|uint32(buf[2])<<16|uint32(buf[3])<<8|uint32(buf[4])), 0)
	pk.PubKeyAlgo = PublicKeyAlgorithm(buf[5])
	// Ignore four-ocet length
	switch pk.PubKeyAlgo {
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		err = pk.parseRSA(r)
	case PubKeyAlgoDSA:
		err = pk.parseDSA(r)
	case PubKeyAlgoElGamal:
		err = pk.parseElGamal(r)
	case PubKeyAlgoECDSA:
		err = pk.parseECDSA(r)
	case PubKeyAlgoECDH:
		err = pk.parseECDH(r)
	case PubKeyAlgoEdDSA:
		err = pk.parseEdDSA(r)
	case PubKeyAlgoX25519:
		err = pk.parseX25519(r)
	case PubKeyAlgoX448:
		err = pk.parseX448(r)
	case PubKeyAlgoEd25519:
		err = pk.parseEd25519(r)
	case PubKeyAlgoEd448:
		err = pk.parseEd448(r)
	default:
		err = errors.UnsupportedError("public key type: " + strconv.Itoa(int(pk.PubKeyAlgo)))
	}
	if err != nil {
		return
	}

	pk.setFingerprintAndKeyId()
	return
}

func (pk *PublicKey) setFingerprintAndKeyId() {
	// RFC 4880, section 12.2
	if pk.Version >= 5 {
		fingerprint := sha256.New()
		if err := pk.SerializeForHash(fingerprint); err != nil {
			// Should not happen for a hash.
			panic(err)
		}
		pk.Fingerprint = make([]byte, 32)
		copy(pk.Fingerprint, fingerprint.Sum(nil))
		pk.KeyId = binary.BigEndian.Uint64(pk.Fingerprint[:8])
	} else {
		fingerprint := sha1.New()
		if err := pk.SerializeForHash(fingerprint); err != nil {
			// Should not happen for a hash.
			panic(err)
		}
		pk.Fingerprint = make([]byte, 20)
		copy(pk.Fingerprint, fingerprint.Sum(nil))
		pk.KeyId = binary.BigEndian.Uint64(pk.Fingerprint[12:20])
	}
}

func (pk *PublicKey) checkV6Compatibility() error {
	// Implementations MUST NOT accept or generate version 6 key material using the deprecated OIDs.
	switch pk.PubKeyAlgo {
	case PubKeyAlgoECDH:
		curveInfo := ecc.FindByOid(pk.oid)
		if curveInfo == nil {
			return errors.UnsupportedError(fmt.Sprintf("unknown oid: %x", pk.oid))
		}
		if curveInfo.GenName == ecc.Curve25519GenName {
			return errors.StructuralError("cannot generate v6 key with deprecated OID: Curve25519Legacy")
		}
	case PubKeyAlgoEdDSA:
		return errors.StructuralError("cannot generate v6 key with deprecated algorithm: EdDSALegacy")
	}
	return nil
}

// parseRSA parses RSA public key material from the given Reader. See RFC 4880,
// section 5.5.2.
func (pk *PublicKey) parseRSA(r io.Reader) (err error) {
	pk.n = new(encoding.MPI)
	if _, err = pk.n.ReadFrom(r); err != nil {
		return
	}
	pk.e = new(encoding.MPI)
	if _, err = pk.e.ReadFrom(r); err != nil {
		return
	}

	if len(pk.e.Bytes()) > 3 {
		err = errors.UnsupportedError("large public exponent")
		return
	}
	rsa := &rsa.PublicKey{
		N: new(big.Int).SetBytes(pk.n.Bytes()),
		E: 0,
	}
	for i := 0; i < len(pk.e.Bytes()); i++ {
		rsa.E <<= 8
		rsa.E |= int(pk.e.Bytes()[i])
	}
	pk.PublicKey = rsa
	return
}

// parseDSA parses DSA public key material from the given Reader. See RFC 4880,
// section 5.5.2.
func (pk *PublicKey) parseDSA(r io.Reader) (err error) {
	pk.p = new(encoding.MPI)
	if _, err = pk.p.ReadFrom(r); err != nil {
		return
	}
	pk.q = new(encoding.MPI)
	if _, err = pk.q.ReadFrom(r); err != nil {
		return
	}
	pk.g = new(encoding.MPI)
	if _, err = pk.g.ReadFrom(r); err != nil {
		return
	}
	pk.y = new(encoding.MPI)
	if _, err = pk.y.ReadFrom(r); err != nil {
		return
	}

	dsa := new(dsa.PublicKey)
	dsa.P = new(big.Int).SetBytes(pk.p.Bytes())
	dsa.Q = new(big.Int).SetBytes(pk.q.Bytes())
	dsa.G = new(big.Int).SetBytes(pk.g.Bytes())
	dsa.Y = new(big.Int).SetBytes(pk.y.Bytes())
	pk.PublicKey = dsa
	return
}

// parseElGamal parses ElGamal public key material from the given Reader. See
// RFC 4880, section 5.5.2.
func (pk *PublicKey) parseElGamal(r io.Reader) (err error) {
	pk.p = new(encoding.MPI)
	if _, err = pk.p.ReadFrom(r); err != nil {
		return
	}
	pk.g = new(encoding.MPI)
	if _, err = pk.g.ReadFrom(r); err != nil {
		return
	}
	pk.y = new(encoding.MPI)
	if _, err = pk.y.ReadFrom(r); err != nil {
		return
	}

	elgamal := new(elgamal.PublicKey)
	elgamal.P = new(big.Int).SetBytes(pk.p.Bytes())
	elgamal.G = new(big.Int).SetBytes(pk.g.Bytes())
	elgamal.Y = new(big.Int).SetBytes(pk.y.Bytes())
	pk.PublicKey = elgamal
	return
}

// parseECDSA parses ECDSA public key material from the given Reader. See
// RFC 6637, Section 9.
func (pk *PublicKey) parseECDSA(r io.Reader) (err error) {
	pk.oid = new(encoding.OID)
	if _, err = pk.oid.ReadFrom(r); err != nil {
		return
	}

	curveInfo := ecc.FindByOid(pk.oid)
	if curveInfo == nil {
		return errors.UnsupportedError(fmt.Sprintf("unknown oid: %x", pk.oid))
	}

	pk.p = new(encoding.MPI)
	if _, err = pk.p.ReadFrom(r); err != nil {
		return
	}

	c, ok := curveInfo.Curve.(ecc.ECDSACurve)
	if !ok {
		return errors.UnsupportedError(fmt.Sprintf("unsupported oid: %x", pk.oid))
	}

	ecdsaKey := ecdsa.NewPublicKey(c)
	err = ecdsaKey.UnmarshalPoint(pk.p.Bytes())
	pk.PublicKey = ecdsaKey

	return
}

// parseECDH parses ECDH public key material from the given Reader. See
// RFC 6637, Section 9.
func (pk *PublicKey) parseECDH(r io.Reader) (err error) {
	pk.oid = new(encoding.OID)
	if _, err = pk.oid.ReadFrom(r); err != nil {
		return
	}

	curveInfo := ecc.FindByOid(pk.oid)
	if curveInfo == nil {
		return errors.UnsupportedError(fmt.Sprintf("unknown oid: %x", pk.oid))
	}

	if pk.Version == 6 && curveInfo.GenName == ecc.Curve25519GenName {
		// Implementations MUST NOT accept or generate version 6 key material using the deprecated OIDs.
		return errors.StructuralError("cannot read v6 key with deprecated OID: Curve25519Legacy")
	}

	pk.p = new(encoding.MPI)
	if _, err = pk.p.ReadFrom(r); err != nil {
		return
	}
	pk.kdf = new(encoding.OID)
	if _, err = pk.kdf.ReadFrom(r); err != nil {
		return
	}

	c, ok := curveInfo.Curve.(ecc.ECDHCurve)
	if !ok {
		return errors.UnsupportedError(fmt.Sprintf("unsupported oid: %x", pk.oid))
	}

	if kdfLen := len(pk.kdf.Bytes()); kdfLen < 3 {
		return errors.UnsupportedError("unsupported ECDH KDF length: " + strconv.Itoa(kdfLen))
	}
	if reserved := pk.kdf.Bytes()[0]; reserved != 0x01 {
		return errors.UnsupportedError("unsupported KDF reserved field: " + strconv.Itoa(int(reserved)))
	}
	kdfHash, ok := algorithm.HashById[pk.kdf.Bytes()[1]]
	if !ok {
		return errors.UnsupportedError("unsupported ECDH KDF hash: " + strconv.Itoa(int(pk.kdf.Bytes()[1])))
	}
	kdfCipher, ok := algorithm.CipherById[pk.kdf.Bytes()[2]]
	if !ok {
		return errors.UnsupportedError("unsupported ECDH KDF cipher: " + strconv.Itoa(int(pk.kdf.Bytes()[2])))
	}

	ecdhKey := ecdh.NewPublicKey(c, kdfHash, kdfCipher)
	err = ecdhKey.UnmarshalPoint(pk.p.Bytes())
	pk.PublicKey = ecdhKey

	return
}

func (pk *PublicKey) parseEdDSA(r io.Reader) (err error) {
	if pk.Version == 6 {
		// Implementations MUST NOT accept or generate version 6 key material using the deprecated OIDs.
		return errors.StructuralError("cannot generate v6 key with deprecated algorithm: EdDSALegacy")
	}

	pk.oid = new(encoding.OID)
	if _, err = pk.oid.ReadFrom(r); err != nil {
		return
	}

	curveInfo := ecc.FindByOid(pk.oid)
	if curveInfo == nil {
		return errors.UnsupportedError(fmt.Sprintf("unknown oid: %x", pk.oid))
	}

	c, ok := curveInfo.Curve.(ecc.EdDSACurve)
	if !ok {
		return errors.UnsupportedError(fmt.Sprintf("unsupported oid: %x", pk.oid))
	}

	pk.p = new(encoding.MPI)
	if _, err = pk.p.ReadFrom(r); err != nil {
		return
	}

	if len(pk.p.Bytes()) == 0 {
		return errors.StructuralError("empty EdDSA public key")
	}

	pub := eddsa.NewPublicKey(c)

	switch flag := pk.p.Bytes()[0]; flag {
	case 0x04:
		// TODO: see _grcy_ecc_eddsa_ensure_compact in grcypt
		return errors.UnsupportedError("unsupported EdDSA compression: " + strconv.Itoa(int(flag)))
	case 0x40:
		err = pub.UnmarshalPoint(pk.p.Bytes())
	default:
		return errors.UnsupportedError("unsupported EdDSA compression: " + strconv.Itoa(int(flag)))
	}

	pk.PublicKey = pub
	return
}

func (pk *PublicKey) parseX25519(r io.Reader) (err error) {
	point := make([]byte, x25519.KeySize)
	_, err = io.ReadFull(r, point)
	if err != nil {
		return
	}
	pub := &x25519.PublicKey{
		Point: point,
	}
	pk.PublicKey = pub
	return
}

func (pk *PublicKey) parseX448(r io.Reader) (err error) {
	point := make([]byte, x448.KeySize)
	_, err = io.ReadFull(r, point)
	if err != nil {
		return
	}
	pub := &x448.PublicKey{
		Point: point,
	}
	pk.PublicKey = pub
	return
}

func (pk *PublicKey) parseEd25519(r io.Reader) (err error) {
	point := make([]byte, ed25519.PublicKeySize)
	_, err = io.ReadFull(r, point)
	if err != nil {
		return
	}
	pub := &ed25519.PublicKey{
		Point: point,
	}
	pk.PublicKey = pub
	return
}

func (pk *PublicKey) parseEd448(r io.Reader) (err error) {
	point := make([]byte, ed448.PublicKeySize)
	_, err = io.ReadFull(r, point)
	if err != nil {
		return
	}
	pub := &ed448.PublicKey{
		Point: point,
	}
	pk.PublicKey = pub
	return
}

// SerializeForHash serializes the PublicKey to w with the special packet
// header format needed for hashing.
func (pk *PublicKey) SerializeForHash(w io.Writer) error {
	if err := pk.SerializeSignaturePrefix(w); err != nil {
		return err
	}
	return pk.serializeWithoutHeaders(w)
}

// SerializeSignaturePrefix writes the prefix for this public key to the given Writer.
// The prefix is used when calculating a signature over this public key. See
// RFC 4880, section 5.2.4.
func (pk *PublicKey) SerializeSignaturePrefix(w io.Writer) error {
	var pLength = pk.algorithmSpecificByteCount()
	// version, timestamp, algorithm
	pLength += versionSize + timestampSize + algorithmSize
	if pk.Version >= 5 {
		// key octet count (4).
		pLength += 4
		_, err := w.Write([]byte{
			// When a v4 signature is made over a key, the hash data starts with the octet 0x99, followed by a two-octet length
			// of the key, and then the body of the key packet. When a v6 signature is made over a key, the hash data starts
			// with the salt, then octet 0x9B, followed by a four-octet length of the key, and then the body of the key packet.
			0x95 + byte(pk.Version),
			byte(pLength >> 24),
			byte(pLength >> 16),
			byte(pLength >> 8),
			byte(pLength),
		})
		return err
	}
	if _, err := w.Write([]byte{0x99, byte(pLength >> 8), byte(pLength)}); err != nil {
		return err
	}
	return nil
}

func (pk *PublicKey) Serialize(w io.Writer) (err error) {
	length := uint32(versionSize + timestampSize + algorithmSize) // 6 byte header
	length += pk.algorithmSpecificByteCount()
	if pk.Version >= 5 {
		length += 4 // octet key count
	}
	packetType := packetTypePublicKey
	if pk.IsSubkey {
		packetType = packetTypePublicSubkey
	}
	err = serializeHeader(w, packetType, int(length))
	if err != nil {
		return
	}
	return pk.serializeWithoutHeaders(w)
}

func (pk *PublicKey) algorithmSpecificByteCount() uint32 {
	length := uint32(0)
	switch pk.PubKeyAlgo {
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		length += uint32(pk.n.EncodedLength())
		length += uint32(pk.e.EncodedLength())
	case PubKeyAlgoDSA:
		length += uint32(pk.p.EncodedLength())
		length += uint32(pk.q.EncodedLength())
		length += uint32(pk.g.EncodedLength())
		length += uint32(pk.y.EncodedLength())
	case PubKeyAlgoElGamal:
		length += uint32(pk.p.EncodedLength())
		length += uint32(pk.g.EncodedLength())
		length += uint32(pk.y.EncodedLength())
	case PubKeyAlgoECDSA:
		length += uint32(pk.oid.EncodedLength())
		length += uint32(pk.p.EncodedLength())
	case PubKeyAlgoECDH:
		length += uint32(pk.oid.EncodedLength())
		length += uint32(pk.p.EncodedLength())
		length += uint32(pk.kdf.EncodedLength())
	case PubKeyAlgoEdDSA:
		length += uint32(pk.oid.EncodedLength())
		length += uint32(pk.p.EncodedLength())
	case PubKeyAlgoX25519:
		length += x25519.KeySize
	case PubKeyAlgoX448:
		length += x448.KeySize
	case PubKeyAlgoEd25519:
		length += ed25519.PublicKeySize
	case PubKeyAlgoEd448:
		length += ed448.PublicKeySize
	default:
		panic("unknown public key algorithm")
	}
	return length
}

// serializeWithoutHeaders marshals the PublicKey to w in the form of an
// OpenPGP public key packet, not including the packet header.
func (pk *PublicKey) serializeWithoutHeaders(w io.Writer) (err error) {
	t := uint32(pk.CreationTime.Unix())
	if _, err = w.Write([]byte{
		byte(pk.Version),
		byte(t >> 24), byte(t >> 16), byte(t >> 8), byte(t),
		byte(pk.PubKeyAlgo),
	}); err != nil {
		return
	}

	if pk.Version >= 5 {
		n := pk.algorithmSpecificByteCount()
		if _, err = w.Write([]byte{
			byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n),
		}); err != nil {
			return
		}
	}

	switch pk.PubKeyAlgo {
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		if _, err = w.Write(pk.n.EncodedBytes()); err != nil {
			return
		}
		_, err = w.Write(pk.e.EncodedBytes())
		return
	case PubKeyAlgoDSA:
		if _, err = w.Write(pk.p.EncodedBytes()); err != nil {
			return
		}
		if _, err = w.Write(pk.q.EncodedBytes()); err != nil {
			return
		}
		if _, err = w.Write(pk.g.EncodedBytes()); err != nil {
			return
		}
		_, err = w.Write(pk.y.EncodedBytes())
		return
	case PubKeyAlgoElGamal:
		if _, err = w.Write(pk.p.EncodedBytes()); err != nil {
			return
		}
		if _, err = w.Write(pk.g.EncodedBytes()); err != nil {
			return
		}
		_, err = w.Write(pk.y.EncodedBytes())
		return
	case PubKeyAlgoECDSA:
		if _, err = w.Write(pk.oid.EncodedBytes()); err != nil {
			return
		}
		_, err = w.Write(pk.p.EncodedBytes())
		return
	case PubKeyAlgoECDH:
		if _, err = w.Write(pk.oid.EncodedBytes()); err != nil {
			return
		}
		if _, err = w.Write(pk.p.EncodedBytes()); err != nil {
			return
		}
		_, err = w.Write(pk.kdf.EncodedBytes())
		return
	case PubKeyAlgoEdDSA:
		if _, err = w.Write(pk.oid.EncodedBytes()); err != nil {
			return
		}
		_, err = w.Write(pk.p.EncodedBytes())
		return
	case PubKeyAlgoX25519:
		publicKey := pk.PublicKey.(*x25519.PublicKey)
		_, err = w.Write(publicKey.Point)
		return
	case PubKeyAlgoX448:
		publicKey := pk.PublicKey.(*x448.PublicKey)
		_, err = w.Write(publicKey.Point)
		return
	case PubKeyAlgoEd25519:
		publicKey := pk.PublicKey.(*ed25519.PublicKey)
		_, err = w.Write(publicKey.Point)
		return
	case PubKeyAlgoEd448:
		publicKey := pk.PublicKey.(*ed448.PublicKey)
		_, err = w.Write(publicKey.Point)
		return
	}
	return errors.InvalidArgumentError("bad public-key algorithm")
}

// CanSign returns true iff this public key can generate signatures
func (pk *PublicKey) CanSign() bool {
	return pk.PubKeyAlgo != PubKeyAlgoRSAEncryptOnly && pk.PubKeyAlgo != PubKeyAlgoElGamal && pk.PubKeyAlgo != PubKeyAlgoECDH
}

// VerifyHashTag returns nil iff sig appears to be a plausible signature of the data
// hashed into signed, based solely on its HashTag. signed is mutated by this call.
func VerifyHashTag(signed hash.Hash, sig *Signature) (err error) {
	if sig.Version == 5 && (sig.SigType == 0x00 || sig.SigType == 0x01) {
		sig.AddMetadataToHashSuffix()
	}
	signed.Write(sig.HashSuffix)
	hashBytes := signed.Sum(nil)
	if hashBytes[0] != sig.HashTag[0] || hashBytes[1] != sig.HashTag[1] {
		return errors.SignatureError("hash tag doesn't match")
	}
	return nil
}

// VerifySignature returns nil iff sig is a valid signature, made by this
// public key, of the data hashed into signed. signed is mutated by this call.
func (pk *PublicKey) VerifySignature(signed hash.Hash, sig *Signature) (err error) {
	if !pk.CanSign() {
		return errors.InvalidArgumentError("public key cannot generate signatures")
	}
	if sig.Version == 5 && (sig.SigType == 0x00 || sig.SigType == 0x01) {
		sig.AddMetadataToHashSuffix()
	}
	signed.Write(sig.HashSuffix)
	hashBytes := signed.Sum(nil)
	// see discussion https://github.com/ProtonMail/go-crypto/issues/107
	if sig.Version >= 5 && (hashBytes[0] != sig.HashTag[0] || hashBytes[1] != sig.HashTag[1]) {
		return errors.SignatureError("hash tag doesn't match")
	}

	if pk.PubKeyAlgo != sig.PubKeyAlgo {
		return errors.InvalidArgumentError("public key and signature use different algorithms")
	}

	switch pk.PubKeyAlgo {
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly:
		rsaPublicKey, _ := pk.PublicKey.(*rsa.PublicKey)
		err = rsa.VerifyPKCS1v15(rsaPublicKey, sig.Hash, hashBytes, padToKeySize(rsaPublicKey, sig.RSASignature.Bytes()))
		if err != nil {
			return errors.SignatureError("RSA verification failure")
		}
		return nil
	case PubKeyAlgoDSA:
		dsaPublicKey, _ := pk.PublicKey.(*dsa.PublicKey)
		// Need to truncate hashBytes to match FIPS 186-3 section 4.6.
		subgroupSize := (dsaPublicKey.Q.BitLen() + 7) / 8
		if len(hashBytes) > subgroupSize {
			hashBytes = hashBytes[:subgroupSize]
		}
		if !dsa.Verify(dsaPublicKey, hashBytes, new(big.Int).SetBytes(sig.DSASigR.Bytes()), new(big.Int).SetBytes(sig.DSASigS.Bytes())) {
			return errors.SignatureError("DSA verification failure")
		}
		return nil
	case PubKeyAlgoECDSA:
		ecdsaPublicKey := pk.PublicKey.(*ecdsa.PublicKey)
		if !ecdsa.Verify(ecdsaPublicKey, hashBytes, new(big.Int).SetBytes(sig.ECDSASigR.Bytes()), new(big.Int).SetBytes(sig.ECDSASigS.Bytes())) {
			return errors.SignatureError("ECDSA verification failure")
		}
		return nil
	case PubKeyAlgoEdDSA:
		eddsaPublicKey := pk.PublicKey.(*eddsa.PublicKey)
		if !eddsa.Verify(eddsaPublicKey, hashBytes, sig.EdDSASigR.Bytes(), sig.EdDSASigS.Bytes()) {
			return errors.SignatureError("EdDSA verification failure")
		}
		return nil
	case PubKeyAlgoEd25519:
		ed25519PublicKey := pk.PublicKey.(*ed25519.PublicKey)
		if !ed25519.Verify(ed25519PublicKey, hashBytes, sig.EdSig) {
			return errors.SignatureError("Ed25519 verification failure")
		}
		return nil
	case PubKeyAlgoEd448:
		ed448PublicKey := pk.PublicKey.(*ed448.PublicKey)
		if !ed448.Verify(ed448PublicKey, hashBytes, sig.EdSig) {
			return errors.SignatureError("ed448 verification failure")
		}
		return nil
	default:
		return errors.SignatureError("Unsupported public key algorithm used in signature")
	}
}

// keySignatureHash returns a Hash of the message that needs to be signed for
// pk to assert a subkey relationship to signed.
func keySignatureHash(pk, signed signingKey, hashFunc hash.Hash) (h hash.Hash, err error) {
	h = hashFunc

	// RFC 4880, section 5.2.4
	err = pk.SerializeForHash(h)
	if err != nil {
		return nil, err
	}

	err = signed.SerializeForHash(h)
	return
}

// VerifyKeyHashTag returns nil iff sig appears to be a plausible signature over this
// primary key and subkey, based solely on its HashTag.
func (pk *PublicKey) VerifyKeyHashTag(signed *PublicKey, sig *Signature) error {
	preparedHash, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	h, err := keySignatureHash(pk, signed, preparedHash)
	if err != nil {
		return err
	}
	return VerifyHashTag(h, sig)
}

// VerifyKeySignature returns nil iff sig is a valid signature, made by this
// public key, of signed.
func (pk *PublicKey) VerifyKeySignature(signed *PublicKey, sig *Signature) error {
	preparedHash, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	h, err := keySignatureHash(pk, signed, preparedHash)
	if err != nil {
		return err
	}
	if err = pk.VerifySignature(h, sig); err != nil {
		return err
	}

	if sig.FlagSign {
		// Signing subkeys must be cross-signed. See
		// https://www.gnupg.org/faq/subkey-cross-certify.html.
		if sig.EmbeddedSignature == nil {
			return errors.StructuralError("signing subkey is missing cross-signature")
		}
		preparedHashEmbedded, err := sig.EmbeddedSignature.PrepareVerify()
		if err != nil {
			return err
		}
		// Verify the cross-signature. This is calculated over the same
		// data as the main signature, so we cannot just recursively
		// call signed.VerifyKeySignature(...)
		if h, err = keySignatureHash(pk, signed, preparedHashEmbedded); err != nil {
			return errors.StructuralError("error while hashing for cross-signature: " + err.Error())
		}
		if err := signed.VerifySignature(h, sig.EmbeddedSignature); err != nil {
			return errors.StructuralError("error while verifying cross-signature: " + err.Error())
		}
	}

	return nil
}

func keyRevocationHash(pk signingKey, hashFunc hash.Hash) (err error) {
	return pk.SerializeForHash(hashFunc)
}

// VerifyRevocationHashTag returns nil iff sig appears to be a plausible signature
// over this public key, based solely on its HashTag.
func (pk *PublicKey) VerifyRevocationHashTag(sig *Signature) (err error) {
	preparedHash, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	if err = keyRevocationHash(pk, preparedHash); err != nil {
		return err
	}
	return VerifyHashTag(preparedHash, sig)
}

// VerifyRevocationSignature returns nil iff sig is a valid signature, made by this
// public key.
func (pk *PublicKey) VerifyRevocationSignature(sig *Signature) (err error) {
	preparedHash, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	if err = keyRevocationHash(pk, preparedHash); err != nil {
		return err
	}
	return pk.VerifySignature(preparedHash, sig)
}

// VerifySubkeyRevocationSignature returns nil iff sig is a valid subkey revocation signature,
// made by this public key, of signed.
func (pk *PublicKey) VerifySubkeyRevocationSignature(sig *Signature, signed *PublicKey) (err error) {
	preparedHash, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	h, err := keySignatureHash(pk, signed, preparedHash)
	if err != nil {
		return err
	}
	return pk.VerifySignature(h, sig)
}

// userIdSignatureHash returns a Hash of the message that needs to be signed
// to assert that pk is a valid key for id.
func userIdSignatureHash(id string, pk *PublicKey, h hash.Hash) (err error) {

	// RFC 4880, section 5.2.4
	if err := pk.SerializeSignaturePrefix(h); err != nil {
		return err
	}
	if err := pk.serializeWithoutHeaders(h); err != nil {
		return err
	}

	var buf [5]byte
	buf[0] = 0xb4
	buf[1] = byte(len(id) >> 24)
	buf[2] = byte(len(id) >> 16)
	buf[3] = byte(len(id) >> 8)
	buf[4] = byte(len(id))
	h.Write(buf[:])
	h.Write([]byte(id))

	return nil
}

// directKeySignatureHash returns a Hash of the message that needs to be signed.
func directKeySignatureHash(pk *PublicKey, h hash.Hash) (err error) {
	return pk.SerializeForHash(h)
}

// VerifyUserIdHashTag returns nil iff sig appears to be a plausible signature over this
// public key and UserId, based solely on its HashTag
func (pk *PublicKey) VerifyUserIdHashTag(id string, sig *Signature) (err error) {
	preparedHash, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	err = userIdSignatureHash(id, pk, preparedHash)
	if err != nil {
		return err
	}
	return VerifyHashTag(preparedHash, sig)
}

// VerifyUserIdSignature returns nil iff sig is a valid signature, made by this
// public key, that id is the identity of pub.
func (pk *PublicKey) VerifyUserIdSignature(id string, pub *PublicKey, sig *Signature) (err error) {
	h, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	if err := userIdSignatureHash(id, pub, h); err != nil {
		return err
	}
	return pk.VerifySignature(h, sig)
}

// VerifyDirectKeySignature returns nil iff sig is a valid signature, made by this
// public key.
func (pk *PublicKey) VerifyDirectKeySignature(sig *Signature) (err error) {
	h, err := sig.PrepareVerify()
	if err != nil {
		return err
	}
	if err := directKeySignatureHash(pk, h); err != nil {
		return err
	}
	return pk.VerifySignature(h, sig)
}

// KeyIdString returns the public key's fingerprint in capital hex
// (e.g. "6C7EE1B8621CC013").
func (pk *PublicKey) KeyIdString() string {
	return fmt.Sprintf("%016X", pk.KeyId)
}

// KeyIdShortString returns the short form of public key's fingerprint
// in capital hex, as shown by gpg --list-keys (e.g. "621CC013").
// This function will return the full key id for v5 and v6 keys
// since the short key id is undefined for them.
func (pk *PublicKey) KeyIdShortString() string {
	if pk.Version >= 5 {
		return pk.KeyIdString()
	}
	return fmt.Sprintf("%X", pk.Fingerprint[16:20])
}

// BitLength returns the bit length for the given public key.
func (pk *PublicKey) BitLength() (bitLength uint16, err error) {
	switch pk.PubKeyAlgo {
	case PubKeyAlgoRSA, PubKeyAlgoRSAEncryptOnly, PubKeyAlgoRSASignOnly:
		bitLength = pk.n.BitLength()
	case PubKeyAlgoDSA:
		bitLength = pk.p.BitLength()
	case PubKeyAlgoElGamal:
		bitLength = pk.p.BitLength()
	case PubKeyAlgoECDSA:
		bitLength = pk.p.BitLength()
	case PubKeyAlgoECDH:
		bitLength = pk.p.BitLength()
	case PubKeyAlgoEdDSA:
		bitLength = pk.p.BitLength()
	case PubKeyAlgoX25519:
		bitLength = x25519.KeySize * 8
	case PubKeyAlgoX448:
		bitLength = x448.KeySize * 8
	case PubKeyAlgoEd25519:
		bitLength = ed25519.PublicKeySize * 8
	case PubKeyAlgoEd448:
		bitLength = ed448.PublicKeySize * 8
	default:
		err = errors.InvalidArgumentError("bad public-key algorithm")
	}
	return
}

// Curve returns the used elliptic curve of this public key.
// Returns an error if no elliptic curve is used.
func (pk *PublicKey) Curve() (curve Curve, err error) {
	switch pk.PubKeyAlgo {
	case PubKeyAlgoECDSA, PubKeyAlgoECDH, PubKeyAlgoEdDSA:
		curveInfo := ecc.FindByOid(pk.oid)
		if curveInfo == nil {
			return "", errors.UnsupportedError(fmt.Sprintf("unknown oid: %x", pk.oid))
		}
		curve = Curve(curveInfo.GenName)
	case PubKeyAlgoEd25519, PubKeyAlgoX25519:
		curve = Curve25519
	case PubKeyAlgoEd448, PubKeyAlgoX448:
		curve = Curve448
	default:
		err = errors.InvalidArgumentError("public key does not operate with an elliptic curve")
	}
	return
}

// KeyExpired returns whether sig is a self-signature of a key that has
// expired or is created in the future.
func (pk *PublicKey) KeyExpired(sig *Signature, currentTime time.Time) bool {
	if pk.CreationTime.Unix() > currentTime.Unix() {
		return true
	}
	if sig.KeyLifetimeSecs == nil || *sig.KeyLifetimeSecs == 0 {
		return false
	}
	expiry := pk.CreationTime.Add(time.Duration(*sig.KeyLifetimeSecs) * time.Second)
	return currentTime.Unix() > expiry.Unix()
}
