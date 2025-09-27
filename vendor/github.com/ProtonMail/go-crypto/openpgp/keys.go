// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openpgp

import (
	goerrors "errors"
	"fmt"
	"io"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// PublicKeyType is the armor type for a PGP public key.
var PublicKeyType = "PGP PUBLIC KEY BLOCK"

// PrivateKeyType is the armor type for a PGP private key.
var PrivateKeyType = "PGP PRIVATE KEY BLOCK"

// An Entity represents the components of an OpenPGP key: a primary public key
// (which must be a signing key), one or more identities claimed by that key,
// and zero or more subkeys, which may be encryption keys.
type Entity struct {
	PrimaryKey    *packet.PublicKey
	PrivateKey    *packet.PrivateKey
	Identities    map[string]*Identity // indexed by Identity.Name
	Revocations   []*packet.Signature
	Subkeys       []Subkey
	SelfSignature *packet.Signature   // Direct-key self signature of the PrimaryKey (contains primary key properties in v6)
	Signatures    []*packet.Signature // all (potentially unverified) self-signatures, revocations, and third-party signatures
}

// An Identity represents an identity claimed by an Entity and zero or more
// assertions by other entities about that claim.
type Identity struct {
	Name          string // by convention, has the form "Full Name (comment) <email@example.com>"
	UserId        *packet.UserId
	SelfSignature *packet.Signature
	Revocations   []*packet.Signature
	Signatures    []*packet.Signature // all (potentially unverified) self-signatures, revocations, and third-party signatures
}

// A Subkey is an additional public key in an Entity. Subkeys can be used for
// encryption.
type Subkey struct {
	PublicKey   *packet.PublicKey
	PrivateKey  *packet.PrivateKey
	Sig         *packet.Signature
	Revocations []*packet.Signature
}

// A Key identifies a specific public key in an Entity. This is either the
// Entity's primary key or a subkey.
type Key struct {
	Entity        *Entity
	PublicKey     *packet.PublicKey
	PrivateKey    *packet.PrivateKey
	SelfSignature *packet.Signature
	Revocations   []*packet.Signature
}

// A KeyRing provides access to public and private keys.
type KeyRing interface {
	// KeysById returns the set of keys that have the given key id.
	KeysById(id uint64) []Key
	// KeysByIdAndUsage returns the set of keys with the given id
	// that also meet the key usage given by requiredUsage.
	// The requiredUsage is expressed as the bitwise-OR of
	// packet.KeyFlag* values.
	KeysByIdUsage(id uint64, requiredUsage byte) []Key
	// DecryptionKeys returns all private keys that are valid for
	// decryption.
	DecryptionKeys() []Key
}

// PrimaryIdentity returns an Identity, preferring non-revoked identities,
// identities marked as primary, or the latest-created identity, in that order.
func (e *Entity) PrimaryIdentity() *Identity {
	var primaryIdentity *Identity
	for _, ident := range e.Identities {
		if shouldPreferIdentity(primaryIdentity, ident) {
			primaryIdentity = ident
		}
	}
	return primaryIdentity
}

func shouldPreferIdentity(existingId, potentialNewId *Identity) bool {
	if existingId == nil {
		return true
	}

	if len(existingId.Revocations) > len(potentialNewId.Revocations) {
		return true
	}

	if len(existingId.Revocations) < len(potentialNewId.Revocations) {
		return false
	}

	if existingId.SelfSignature == nil {
		return true
	}

	if existingId.SelfSignature.IsPrimaryId != nil && *existingId.SelfSignature.IsPrimaryId &&
		!(potentialNewId.SelfSignature.IsPrimaryId != nil && *potentialNewId.SelfSignature.IsPrimaryId) {
		return false
	}

	if !(existingId.SelfSignature.IsPrimaryId != nil && *existingId.SelfSignature.IsPrimaryId) &&
		potentialNewId.SelfSignature.IsPrimaryId != nil && *potentialNewId.SelfSignature.IsPrimaryId {
		return true
	}

	return potentialNewId.SelfSignature.CreationTime.After(existingId.SelfSignature.CreationTime)
}

// EncryptionKey returns the best candidate Key for encrypting a message to the
// given Entity.
func (e *Entity) EncryptionKey(now time.Time) (Key, bool) {
	// Fail to find any encryption key if the...
	primarySelfSignature, primaryIdentity := e.PrimarySelfSignature()
	if primarySelfSignature == nil || // no self-signature found
		e.PrimaryKey.KeyExpired(primarySelfSignature, now) || // primary key has expired
		e.Revoked(now) || // primary key has been revoked
		primarySelfSignature.SigExpired(now) || // user ID or or direct self-signature has expired
		(primaryIdentity != nil && primaryIdentity.Revoked(now)) { // user ID has been revoked (for v4 keys)
		return Key{}, false
	}

	// Iterate the keys to find the newest, unexpired one
	candidateSubkey := -1
	var maxTime time.Time
	for i, subkey := range e.Subkeys {
		if subkey.Sig.FlagsValid &&
			subkey.Sig.FlagEncryptCommunications &&
			subkey.PublicKey.PubKeyAlgo.CanEncrypt() &&
			!subkey.PublicKey.KeyExpired(subkey.Sig, now) &&
			!subkey.Sig.SigExpired(now) &&
			!subkey.Revoked(now) &&
			(maxTime.IsZero() || subkey.Sig.CreationTime.After(maxTime)) {
			candidateSubkey = i
			maxTime = subkey.Sig.CreationTime
		}
	}

	if candidateSubkey != -1 {
		subkey := e.Subkeys[candidateSubkey]
		return Key{e, subkey.PublicKey, subkey.PrivateKey, subkey.Sig, subkey.Revocations}, true
	}

	// If we don't have any subkeys for encryption and the primary key
	// is marked as OK to encrypt with, then we can use it.
	if primarySelfSignature.FlagsValid && primarySelfSignature.FlagEncryptCommunications &&
		e.PrimaryKey.PubKeyAlgo.CanEncrypt() {
		return Key{e, e.PrimaryKey, e.PrivateKey, primarySelfSignature, e.Revocations}, true
	}

	return Key{}, false
}

// CertificationKey return the best candidate Key for certifying a key with this
// Entity.
func (e *Entity) CertificationKey(now time.Time) (Key, bool) {
	return e.CertificationKeyById(now, 0)
}

// CertificationKeyById return the Key for key certification with this
// Entity and keyID.
func (e *Entity) CertificationKeyById(now time.Time, id uint64) (Key, bool) {
	return e.signingKeyByIdUsage(now, id, packet.KeyFlagCertify)
}

// SigningKey return the best candidate Key for signing a message with this
// Entity.
func (e *Entity) SigningKey(now time.Time) (Key, bool) {
	return e.SigningKeyById(now, 0)
}

// SigningKeyById return the Key for signing a message with this
// Entity and keyID.
func (e *Entity) SigningKeyById(now time.Time, id uint64) (Key, bool) {
	return e.signingKeyByIdUsage(now, id, packet.KeyFlagSign)
}

func (e *Entity) signingKeyByIdUsage(now time.Time, id uint64, flags int) (Key, bool) {
	// Fail to find any signing key if the...
	primarySelfSignature, primaryIdentity := e.PrimarySelfSignature()
	if primarySelfSignature == nil || // no self-signature found
		e.PrimaryKey.KeyExpired(primarySelfSignature, now) || // primary key has expired
		e.Revoked(now) || // primary key has been revoked
		primarySelfSignature.SigExpired(now) || // user ID or direct self-signature has expired
		(primaryIdentity != nil && primaryIdentity.Revoked(now)) { // user ID has been revoked (for v4 keys)
		return Key{}, false
	}

	// Iterate the keys to find the newest, unexpired one
	candidateSubkey := -1
	var maxTime time.Time
	for idx, subkey := range e.Subkeys {
		if subkey.Sig.FlagsValid &&
			(flags&packet.KeyFlagCertify == 0 || subkey.Sig.FlagCertify) &&
			(flags&packet.KeyFlagSign == 0 || subkey.Sig.FlagSign) &&
			subkey.PublicKey.PubKeyAlgo.CanSign() &&
			!subkey.PublicKey.KeyExpired(subkey.Sig, now) &&
			!subkey.Sig.SigExpired(now) &&
			!subkey.Revoked(now) &&
			(maxTime.IsZero() || subkey.Sig.CreationTime.After(maxTime)) &&
			(id == 0 || subkey.PublicKey.KeyId == id) {
			candidateSubkey = idx
			maxTime = subkey.Sig.CreationTime
		}
	}

	if candidateSubkey != -1 {
		subkey := e.Subkeys[candidateSubkey]
		return Key{e, subkey.PublicKey, subkey.PrivateKey, subkey.Sig, subkey.Revocations}, true
	}

	// If we don't have any subkeys for signing and the primary key
	// is marked as OK to sign with, then we can use it.
	if primarySelfSignature.FlagsValid &&
		(flags&packet.KeyFlagCertify == 0 || primarySelfSignature.FlagCertify) &&
		(flags&packet.KeyFlagSign == 0 || primarySelfSignature.FlagSign) &&
		e.PrimaryKey.PubKeyAlgo.CanSign() &&
		(id == 0 || e.PrimaryKey.KeyId == id) {
		return Key{e, e.PrimaryKey, e.PrivateKey, primarySelfSignature, e.Revocations}, true
	}

	// No keys with a valid Signing Flag or no keys matched the id passed in
	return Key{}, false
}

func revoked(revocations []*packet.Signature, now time.Time) bool {
	for _, revocation := range revocations {
		if revocation.RevocationReason != nil && *revocation.RevocationReason == packet.KeyCompromised {
			// If the key is compromised, the key is considered revoked even before the revocation date.
			return true
		}
		if !revocation.SigExpired(now) {
			return true
		}
	}
	return false
}

// Revoked returns whether the entity has any direct key revocation signatures.
// Note that third-party revocation signatures are not supported.
// Note also that Identity and Subkey revocation should be checked separately.
func (e *Entity) Revoked(now time.Time) bool {
	return revoked(e.Revocations, now)
}

// EncryptPrivateKeys encrypts all non-encrypted keys in the entity with the same key
// derived from the provided passphrase. Public keys and dummy keys are ignored,
// and don't cause an error to be returned.
func (e *Entity) EncryptPrivateKeys(passphrase []byte, config *packet.Config) error {
	var keysToEncrypt []*packet.PrivateKey
	// Add entity private key to encrypt.
	if e.PrivateKey != nil && !e.PrivateKey.Dummy() && !e.PrivateKey.Encrypted {
		keysToEncrypt = append(keysToEncrypt, e.PrivateKey)
	}

	// Add subkeys to encrypt.
	for _, sub := range e.Subkeys {
		if sub.PrivateKey != nil && !sub.PrivateKey.Dummy() && !sub.PrivateKey.Encrypted {
			keysToEncrypt = append(keysToEncrypt, sub.PrivateKey)
		}
	}
	return packet.EncryptPrivateKeys(keysToEncrypt, passphrase, config)
}

// DecryptPrivateKeys decrypts all encrypted keys in the entity with the given passphrase.
// Avoids recomputation of similar s2k key derivations. Public keys and dummy keys are ignored,
// and don't cause an error to be returned.
func (e *Entity) DecryptPrivateKeys(passphrase []byte) error {
	var keysToDecrypt []*packet.PrivateKey
	// Add entity private key to decrypt.
	if e.PrivateKey != nil && !e.PrivateKey.Dummy() && e.PrivateKey.Encrypted {
		keysToDecrypt = append(keysToDecrypt, e.PrivateKey)
	}

	// Add subkeys to decrypt.
	for _, sub := range e.Subkeys {
		if sub.PrivateKey != nil && !sub.PrivateKey.Dummy() && sub.PrivateKey.Encrypted {
			keysToDecrypt = append(keysToDecrypt, sub.PrivateKey)
		}
	}
	return packet.DecryptPrivateKeys(keysToDecrypt, passphrase)
}

// Revoked returns whether the identity has been revoked by a self-signature.
// Note that third-party revocation signatures are not supported.
func (i *Identity) Revoked(now time.Time) bool {
	return revoked(i.Revocations, now)
}

// Revoked returns whether the subkey has been revoked by a self-signature.
// Note that third-party revocation signatures are not supported.
func (s *Subkey) Revoked(now time.Time) bool {
	return revoked(s.Revocations, now)
}

// Revoked returns whether the key or subkey has been revoked by a self-signature.
// Note that third-party revocation signatures are not supported.
// Note also that Identity revocation should be checked separately.
// Normally, it's not necessary to call this function, except on keys returned by
// KeysById or KeysByIdUsage.
func (key *Key) Revoked(now time.Time) bool {
	return revoked(key.Revocations, now)
}

// An EntityList contains one or more Entities.
type EntityList []*Entity

// KeysById returns the set of keys that have the given key id.
func (el EntityList) KeysById(id uint64) (keys []Key) {
	for _, e := range el {
		if e.PrimaryKey.KeyId == id {
			selfSig, _ := e.PrimarySelfSignature()
			keys = append(keys, Key{e, e.PrimaryKey, e.PrivateKey, selfSig, e.Revocations})
		}

		for _, subKey := range e.Subkeys {
			if subKey.PublicKey.KeyId == id {
				keys = append(keys, Key{e, subKey.PublicKey, subKey.PrivateKey, subKey.Sig, subKey.Revocations})
			}
		}
	}
	return
}

// KeysByIdAndUsage returns the set of keys with the given id that also meet
// the key usage given by requiredUsage.  The requiredUsage is expressed as
// the bitwise-OR of packet.KeyFlag* values.
func (el EntityList) KeysByIdUsage(id uint64, requiredUsage byte) (keys []Key) {
	for _, key := range el.KeysById(id) {
		if requiredUsage != 0 {
			if key.SelfSignature == nil || !key.SelfSignature.FlagsValid {
				continue
			}

			var usage byte
			if key.SelfSignature.FlagCertify {
				usage |= packet.KeyFlagCertify
			}
			if key.SelfSignature.FlagSign {
				usage |= packet.KeyFlagSign
			}
			if key.SelfSignature.FlagEncryptCommunications {
				usage |= packet.KeyFlagEncryptCommunications
			}
			if key.SelfSignature.FlagEncryptStorage {
				usage |= packet.KeyFlagEncryptStorage
			}
			if usage&requiredUsage != requiredUsage {
				continue
			}
		}

		keys = append(keys, key)
	}
	return
}

// DecryptionKeys returns all private keys that are valid for decryption.
func (el EntityList) DecryptionKeys() (keys []Key) {
	for _, e := range el {
		for _, subKey := range e.Subkeys {
			if subKey.PrivateKey != nil && subKey.Sig.FlagsValid && (subKey.Sig.FlagEncryptStorage || subKey.Sig.FlagEncryptCommunications) {
				keys = append(keys, Key{e, subKey.PublicKey, subKey.PrivateKey, subKey.Sig, subKey.Revocations})
			}
		}
	}
	return
}

// ReadArmoredKeyRing reads one or more public/private keys from an armor keyring file.
func ReadArmoredKeyRing(r io.Reader) (EntityList, error) {
	block, err := armor.Decode(r)
	if err == io.EOF {
		return nil, errors.InvalidArgumentError("no armored data found")
	}
	if err != nil {
		return nil, err
	}
	if block.Type != PublicKeyType && block.Type != PrivateKeyType {
		return nil, errors.InvalidArgumentError("expected public or private key block, got: " + block.Type)
	}

	return ReadKeyRing(block.Body)
}

// ReadKeyRing reads one or more public/private keys. Unsupported keys are
// ignored as long as at least a single valid key is found.
func ReadKeyRing(r io.Reader) (el EntityList, err error) {
	packets := packet.NewReader(r)
	var lastUnsupportedError error

	for {
		var e *Entity
		e, err = ReadEntity(packets)
		if err != nil {
			// TODO: warn about skipped unsupported/unreadable keys
			if _, ok := err.(errors.UnsupportedError); ok {
				lastUnsupportedError = err
				err = readToNextPublicKey(packets)
			} else if _, ok := err.(errors.StructuralError); ok {
				// Skip unreadable, badly-formatted keys
				lastUnsupportedError = err
				err = readToNextPublicKey(packets)
			}
			if err == io.EOF {
				err = nil
				break
			}
			if err != nil {
				el = nil
				break
			}
		} else {
			el = append(el, e)
		}
	}

	if len(el) == 0 && err == nil {
		err = lastUnsupportedError
	}
	return
}

// readToNextPublicKey reads packets until the start of the entity and leaves
// the first packet of the new entity in the Reader.
func readToNextPublicKey(packets *packet.Reader) (err error) {
	var p packet.Packet
	for {
		p, err = packets.Next()
		if err == io.EOF {
			return
		} else if err != nil {
			if _, ok := err.(errors.UnsupportedError); ok {
				continue
			}
			return
		}

		if pk, ok := p.(*packet.PublicKey); ok && !pk.IsSubkey {
			packets.Unread(p)
			return
		}
	}
}

// ReadEntity reads an entity (public key, identities, subkeys etc) from the
// given Reader.
func ReadEntity(packets *packet.Reader) (*Entity, error) {
	e := new(Entity)
	e.Identities = make(map[string]*Identity)

	p, err := packets.Next()
	if err != nil {
		return nil, err
	}

	var ok bool
	if e.PrimaryKey, ok = p.(*packet.PublicKey); !ok {
		if e.PrivateKey, ok = p.(*packet.PrivateKey); !ok {
			packets.Unread(p)
			return nil, errors.StructuralError("first packet was not a public/private key")
		}
		e.PrimaryKey = &e.PrivateKey.PublicKey
	}

	if !e.PrimaryKey.PubKeyAlgo.CanSign() {
		return nil, errors.StructuralError("primary key cannot be used for signatures")
	}

	var revocations []*packet.Signature
	var directSignatures []*packet.Signature
EachPacket:
	for {
		p, err := packets.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch pkt := p.(type) {
		case *packet.UserId:
			if err := addUserID(e, packets, pkt); err != nil {
				return nil, err
			}
		case *packet.Signature:
			if pkt.SigType == packet.SigTypeKeyRevocation {
				revocations = append(revocations, pkt)
			} else if pkt.SigType == packet.SigTypeDirectSignature {
				directSignatures = append(directSignatures, pkt)
			}
			// Else, ignoring the signature as it does not follow anything
			// we would know to attach it to.
		case *packet.PrivateKey:
			if !pkt.IsSubkey {
				packets.Unread(p)
				break EachPacket
			}
			err = addSubkey(e, packets, &pkt.PublicKey, pkt)
			if err != nil {
				return nil, err
			}
		case *packet.PublicKey:
			if !pkt.IsSubkey {
				packets.Unread(p)
				break EachPacket
			}
			err = addSubkey(e, packets, pkt, nil)
			if err != nil {
				return nil, err
			}
		default:
			// we ignore unknown packets.
		}
	}

	if len(e.Identities) == 0 && e.PrimaryKey.Version < 6 {
		return nil, errors.StructuralError(fmt.Sprintf("v%d entity without any identities", e.PrimaryKey.Version))
	}

	// An implementation MUST ensure that a valid direct-key signature is present before using a v6 key.
	if e.PrimaryKey.Version == 6 {
		if len(directSignatures) == 0 {
			return nil, errors.StructuralError("v6 entity without a valid direct-key signature")
		}
		// Select main direct key signature.
		var mainDirectKeySelfSignature *packet.Signature
		for _, directSignature := range directSignatures {
			if directSignature.SigType == packet.SigTypeDirectSignature &&
				directSignature.CheckKeyIdOrFingerprint(e.PrimaryKey) &&
				(mainDirectKeySelfSignature == nil ||
					directSignature.CreationTime.After(mainDirectKeySelfSignature.CreationTime)) {
				mainDirectKeySelfSignature = directSignature
			}
		}
		if mainDirectKeySelfSignature == nil {
			return nil, errors.StructuralError("no valid direct-key self-signature for v6 primary key found")
		}
		// Check that the main self-signature is valid.
		err = e.PrimaryKey.VerifyDirectKeySignature(mainDirectKeySelfSignature)
		if err != nil {
			return nil, errors.StructuralError("invalid direct-key self-signature for v6 primary key")
		}
		e.SelfSignature = mainDirectKeySelfSignature
		e.Signatures = directSignatures
	}

	for _, revocation := range revocations {
		err = e.PrimaryKey.VerifyRevocationSignature(revocation)
		if err == nil {
			e.Revocations = append(e.Revocations, revocation)
		} else {
			// TODO: RFC 4880 5.2.3.15 defines revocation keys.
			return nil, errors.StructuralError("revocation signature signed by alternate key")
		}
	}

	return e, nil
}

func addUserID(e *Entity, packets *packet.Reader, pkt *packet.UserId) error {
	// Make a new Identity object, that we might wind up throwing away.
	// We'll only add it if we get a valid self-signature over this
	// userID.
	identity := new(Identity)
	identity.Name = pkt.Id
	identity.UserId = pkt

	for {
		p, err := packets.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		sig, ok := p.(*packet.Signature)
		if !ok {
			packets.Unread(p)
			break
		}

		if sig.SigType != packet.SigTypeGenericCert &&
			sig.SigType != packet.SigTypePersonaCert &&
			sig.SigType != packet.SigTypeCasualCert &&
			sig.SigType != packet.SigTypePositiveCert &&
			sig.SigType != packet.SigTypeCertificationRevocation {
			return errors.StructuralError("user ID signature with wrong type")
		}

		if sig.CheckKeyIdOrFingerprint(e.PrimaryKey) {
			if err = e.PrimaryKey.VerifyUserIdSignature(pkt.Id, e.PrimaryKey, sig); err != nil {
				return errors.StructuralError("user ID self-signature invalid: " + err.Error())
			}
			if sig.SigType == packet.SigTypeCertificationRevocation {
				identity.Revocations = append(identity.Revocations, sig)
			} else if identity.SelfSignature == nil || sig.CreationTime.After(identity.SelfSignature.CreationTime) {
				identity.SelfSignature = sig
			}
			identity.Signatures = append(identity.Signatures, sig)
			e.Identities[pkt.Id] = identity
		} else {
			identity.Signatures = append(identity.Signatures, sig)
		}
	}

	return nil
}

func addSubkey(e *Entity, packets *packet.Reader, pub *packet.PublicKey, priv *packet.PrivateKey) error {
	var subKey Subkey
	subKey.PublicKey = pub
	subKey.PrivateKey = priv

	for {
		p, err := packets.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.StructuralError("subkey signature invalid: " + err.Error())
		}

		sig, ok := p.(*packet.Signature)
		if !ok {
			packets.Unread(p)
			break
		}

		if sig.SigType != packet.SigTypeSubkeyBinding && sig.SigType != packet.SigTypeSubkeyRevocation {
			return errors.StructuralError("subkey signature with wrong type")
		}

		if err := e.PrimaryKey.VerifyKeySignature(subKey.PublicKey, sig); err != nil {
			return errors.StructuralError("subkey signature invalid: " + err.Error())
		}

		switch sig.SigType {
		case packet.SigTypeSubkeyRevocation:
			subKey.Revocations = append(subKey.Revocations, sig)
		case packet.SigTypeSubkeyBinding:
			if subKey.Sig == nil || sig.CreationTime.After(subKey.Sig.CreationTime) {
				subKey.Sig = sig
			}
		}
	}

	if subKey.Sig == nil {
		return errors.StructuralError("subkey packet not followed by signature")
	}

	e.Subkeys = append(e.Subkeys, subKey)

	return nil
}

// SerializePrivate serializes an Entity, including private key material, but
// excluding signatures from other entities, to the given Writer.
// Identities and subkeys are re-signed in case they changed since NewEntry.
// If config is nil, sensible defaults will be used.
func (e *Entity) SerializePrivate(w io.Writer, config *packet.Config) (err error) {
	if e.PrivateKey.Dummy() {
		return errors.ErrDummyPrivateKey("dummy private key cannot re-sign identities")
	}
	return e.serializePrivate(w, config, true)
}

// SerializePrivateWithoutSigning serializes an Entity, including private key
// material, but excluding signatures from other entities, to the given Writer.
// Self-signatures of identities and subkeys are not re-signed. This is useful
// when serializing GNU dummy keys, among other things.
// If config is nil, sensible defaults will be used.
func (e *Entity) SerializePrivateWithoutSigning(w io.Writer, config *packet.Config) (err error) {
	return e.serializePrivate(w, config, false)
}

func (e *Entity) serializePrivate(w io.Writer, config *packet.Config, reSign bool) (err error) {
	if e.PrivateKey == nil {
		return goerrors.New("openpgp: private key is missing")
	}
	err = e.PrivateKey.Serialize(w)
	if err != nil {
		return
	}
	for _, revocation := range e.Revocations {
		err := revocation.Serialize(w)
		if err != nil {
			return err
		}
	}
	for _, directSignature := range e.Signatures {
		err := directSignature.Serialize(w)
		if err != nil {
			return err
		}
	}
	for _, ident := range e.Identities {
		err = ident.UserId.Serialize(w)
		if err != nil {
			return
		}
		if reSign {
			if ident.SelfSignature == nil {
				return goerrors.New("openpgp: can't re-sign identity without valid self-signature")
			}
			err = ident.SelfSignature.SignUserId(ident.UserId.Id, e.PrimaryKey, e.PrivateKey, config)
			if err != nil {
				return
			}
		}
		for _, sig := range ident.Signatures {
			err = sig.Serialize(w)
			if err != nil {
				return err
			}
		}
	}
	for _, subkey := range e.Subkeys {
		err = subkey.PrivateKey.Serialize(w)
		if err != nil {
			return
		}
		if reSign {
			err = subkey.Sig.SignKey(subkey.PublicKey, e.PrivateKey, config)
			if err != nil {
				return
			}
			if subkey.Sig.EmbeddedSignature != nil {
				err = subkey.Sig.EmbeddedSignature.CrossSignKey(subkey.PublicKey, e.PrimaryKey,
					subkey.PrivateKey, config)
				if err != nil {
					return
				}
			}
		}
		for _, revocation := range subkey.Revocations {
			err := revocation.Serialize(w)
			if err != nil {
				return err
			}
		}
		err = subkey.Sig.Serialize(w)
		if err != nil {
			return
		}
	}
	return nil
}

// Serialize writes the public part of the given Entity to w, including
// signatures from other entities. No private key material will be output.
func (e *Entity) Serialize(w io.Writer) error {
	err := e.PrimaryKey.Serialize(w)
	if err != nil {
		return err
	}
	for _, revocation := range e.Revocations {
		err := revocation.Serialize(w)
		if err != nil {
			return err
		}
	}
	for _, directSignature := range e.Signatures {
		err := directSignature.Serialize(w)
		if err != nil {
			return err
		}
	}
	for _, ident := range e.Identities {
		err = ident.UserId.Serialize(w)
		if err != nil {
			return err
		}
		for _, sig := range ident.Signatures {
			err = sig.Serialize(w)
			if err != nil {
				return err
			}
		}
	}
	for _, subkey := range e.Subkeys {
		err = subkey.PublicKey.Serialize(w)
		if err != nil {
			return err
		}
		for _, revocation := range subkey.Revocations {
			err := revocation.Serialize(w)
			if err != nil {
				return err
			}
		}
		err = subkey.Sig.Serialize(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// SignIdentity adds a signature to e, from signer, attesting that identity is
// associated with e. The provided identity must already be an element of
// e.Identities and the private key of signer must have been decrypted if
// necessary.
// If config is nil, sensible defaults will be used.
func (e *Entity) SignIdentity(identity string, signer *Entity, config *packet.Config) error {
	certificationKey, ok := signer.CertificationKey(config.Now())
	if !ok {
		return errors.InvalidArgumentError("no valid certification key found")
	}

	if certificationKey.PrivateKey.Encrypted {
		return errors.InvalidArgumentError("signing Entity's private key must be decrypted")
	}

	ident, ok := e.Identities[identity]
	if !ok {
		return errors.InvalidArgumentError("given identity string not found in Entity")
	}

	sig := createSignaturePacket(certificationKey.PublicKey, packet.SigTypeGenericCert, config)

	signingUserID := config.SigningUserId()
	if signingUserID != "" {
		if _, ok := signer.Identities[signingUserID]; !ok {
			return errors.InvalidArgumentError("signer identity string not found in signer Entity")
		}
		sig.SignerUserId = &signingUserID
	}

	if err := sig.SignUserId(identity, e.PrimaryKey, certificationKey.PrivateKey, config); err != nil {
		return err
	}
	ident.Signatures = append(ident.Signatures, sig)
	return nil
}

// RevokeKey generates a key revocation signature (packet.SigTypeKeyRevocation) with the
// specified reason code and text (RFC4880 section-5.2.3.23).
// If config is nil, sensible defaults will be used.
func (e *Entity) RevokeKey(reason packet.ReasonForRevocation, reasonText string, config *packet.Config) error {
	revSig := createSignaturePacket(e.PrimaryKey, packet.SigTypeKeyRevocation, config)
	revSig.RevocationReason = &reason
	revSig.RevocationReasonText = reasonText

	if err := revSig.RevokeKey(e.PrimaryKey, e.PrivateKey, config); err != nil {
		return err
	}
	e.Revocations = append(e.Revocations, revSig)
	return nil
}

// RevokeSubkey generates a subkey revocation signature (packet.SigTypeSubkeyRevocation) for
// a subkey with the specified reason code and text (RFC4880 section-5.2.3.23).
// If config is nil, sensible defaults will be used.
func (e *Entity) RevokeSubkey(sk *Subkey, reason packet.ReasonForRevocation, reasonText string, config *packet.Config) error {
	if err := e.PrimaryKey.VerifyKeySignature(sk.PublicKey, sk.Sig); err != nil {
		return errors.InvalidArgumentError("given subkey is not associated with this key")
	}

	revSig := createSignaturePacket(e.PrimaryKey, packet.SigTypeSubkeyRevocation, config)
	revSig.RevocationReason = &reason
	revSig.RevocationReasonText = reasonText

	if err := revSig.RevokeSubkey(sk.PublicKey, e.PrivateKey, config); err != nil {
		return err
	}

	sk.Revocations = append(sk.Revocations, revSig)
	return nil
}

func (e *Entity) primaryDirectSignature() *packet.Signature {
	return e.SelfSignature
}

// PrimarySelfSignature searches the entity for the self-signature that stores key preferences.
// For V4 keys, returns the self-signature of the primary identity, and the identity.
// For V6 keys, returns the latest valid direct-key self-signature, and no identity (nil).
// This self-signature is to be used to check the key expiration,
// algorithm preferences, and so on.
func (e *Entity) PrimarySelfSignature() (*packet.Signature, *Identity) {
	if e.PrimaryKey.Version == 6 {
		return e.primaryDirectSignature(), nil
	}
	primaryIdentity := e.PrimaryIdentity()
	if primaryIdentity == nil {
		return nil, nil
	}
	return primaryIdentity.SelfSignature, primaryIdentity
}
