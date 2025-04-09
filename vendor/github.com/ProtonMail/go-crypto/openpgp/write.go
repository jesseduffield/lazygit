// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package openpgp

import (
	"crypto"
	"hash"
	"io"
	"strconv"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/internal/algorithm"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

// DetachSign signs message with the private key from signer (which must
// already have been decrypted) and writes the signature to w.
// If config is nil, sensible defaults will be used.
func DetachSign(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) error {
	return detachSign(w, signer, message, packet.SigTypeBinary, config)
}

// ArmoredDetachSign signs message with the private key from signer (which
// must already have been decrypted) and writes an armored signature to w.
// If config is nil, sensible defaults will be used.
func ArmoredDetachSign(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) (err error) {
	return armoredDetachSign(w, signer, message, packet.SigTypeBinary, config)
}

// DetachSignText signs message (after canonicalising the line endings) with
// the private key from signer (which must already have been decrypted) and
// writes the signature to w.
// If config is nil, sensible defaults will be used.
func DetachSignText(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) error {
	return detachSign(w, signer, message, packet.SigTypeText, config)
}

// ArmoredDetachSignText signs message (after canonicalising the line endings)
// with the private key from signer (which must already have been decrypted)
// and writes an armored signature to w.
// If config is nil, sensible defaults will be used.
func ArmoredDetachSignText(w io.Writer, signer *Entity, message io.Reader, config *packet.Config) error {
	return armoredDetachSign(w, signer, message, packet.SigTypeText, config)
}

func armoredDetachSign(w io.Writer, signer *Entity, message io.Reader, sigType packet.SignatureType, config *packet.Config) (err error) {
	out, err := armor.Encode(w, SignatureType, nil)
	if err != nil {
		return
	}
	err = detachSign(out, signer, message, sigType, config)
	if err != nil {
		return
	}
	return out.Close()
}

func detachSign(w io.Writer, signer *Entity, message io.Reader, sigType packet.SignatureType, config *packet.Config) (err error) {
	signingKey, ok := signer.SigningKeyById(config.Now(), config.SigningKey())
	if !ok {
		return errors.InvalidArgumentError("no valid signing keys")
	}
	if signingKey.PrivateKey == nil {
		return errors.InvalidArgumentError("signing key doesn't have a private key")
	}
	if signingKey.PrivateKey.Encrypted {
		return errors.InvalidArgumentError("signing key is encrypted")
	}
	if _, ok := algorithm.HashToHashId(config.Hash()); !ok {
		return errors.InvalidArgumentError("invalid hash function")
	}

	sig := createSignaturePacket(signingKey.PublicKey, sigType, config)

	h, err := sig.PrepareSign(config)
	if err != nil {
		return
	}
	wrappedHash, err := wrapHashForSignature(h, sig.SigType)
	if err != nil {
		return
	}
	if _, err = io.Copy(wrappedHash, message); err != nil {
		return err
	}

	err = sig.Sign(h, signingKey.PrivateKey, config)
	if err != nil {
		return
	}

	return sig.Serialize(w)
}

// FileHints contains metadata about encrypted files. This metadata is, itself,
// encrypted.
type FileHints struct {
	// IsBinary can be set to hint that the contents are binary data.
	IsBinary bool
	// FileName hints at the name of the file that should be written. It's
	// truncated to 255 bytes if longer. It may be empty to suggest that the
	// file should not be written to disk. It may be equal to "_CONSOLE" to
	// suggest the data should not be written to disk.
	FileName string
	// ModTime contains the modification time of the file, or the zero time if not applicable.
	ModTime time.Time
}

// SymmetricallyEncrypt acts like gpg -c: it encrypts a file with a passphrase.
// The resulting WriteCloser must be closed after the contents of the file have
// been written.
// If config is nil, sensible defaults will be used.
func SymmetricallyEncrypt(ciphertext io.Writer, passphrase []byte, hints *FileHints, config *packet.Config) (plaintext io.WriteCloser, err error) {
	if hints == nil {
		hints = &FileHints{}
	}

	key, err := packet.SerializeSymmetricKeyEncrypted(ciphertext, passphrase, config)
	if err != nil {
		return
	}

	var w io.WriteCloser
	cipherSuite := packet.CipherSuite{
		Cipher: config.Cipher(),
		Mode:   config.AEAD().Mode(),
	}
	w, err = packet.SerializeSymmetricallyEncrypted(ciphertext, config.Cipher(), config.AEAD() != nil, cipherSuite, key, config)
	if err != nil {
		return
	}

	literalData := w
	if algo := config.Compression(); algo != packet.CompressionNone {
		var compConfig *packet.CompressionConfig
		if config != nil {
			compConfig = config.CompressionConfig
		}
		literalData, err = packet.SerializeCompressed(w, algo, compConfig)
		if err != nil {
			return
		}
	}

	var epochSeconds uint32
	if !hints.ModTime.IsZero() {
		epochSeconds = uint32(hints.ModTime.Unix())
	}
	return packet.SerializeLiteral(literalData, hints.IsBinary, hints.FileName, epochSeconds)
}

// intersectPreferences mutates and returns a prefix of a that contains only
// the values in the intersection of a and b. The order of a is preserved.
func intersectPreferences(a []uint8, b []uint8) (intersection []uint8) {
	var j int
	for _, v := range a {
		for _, v2 := range b {
			if v == v2 {
				a[j] = v
				j++
				break
			}
		}
	}

	return a[:j]
}

// intersectPreferences mutates and returns a prefix of a that contains only
// the values in the intersection of a and b. The order of a is preserved.
func intersectCipherSuites(a [][2]uint8, b [][2]uint8) (intersection [][2]uint8) {
	var j int
	for _, v := range a {
		for _, v2 := range b {
			if v[0] == v2[0] && v[1] == v2[1] {
				a[j] = v
				j++
				break
			}
		}
	}

	return a[:j]
}

func hashToHashId(h crypto.Hash) uint8 {
	v, ok := algorithm.HashToHashId(h)
	if !ok {
		panic("tried to convert unknown hash")
	}
	return v
}

// EncryptText encrypts a message to a number of recipients and, optionally,
// signs it. Optional information is contained in 'hints', also encrypted, that
// aids the recipients in processing the message. The resulting WriteCloser
// must be closed after the contents of the file have been written. If config
// is nil, sensible defaults will be used. The signing is done in text mode.
func EncryptText(ciphertext io.Writer, to []*Entity, signed *Entity, hints *FileHints, config *packet.Config) (plaintext io.WriteCloser, err error) {
	return encrypt(ciphertext, ciphertext, to, signed, hints, packet.SigTypeText, config)
}

// Encrypt encrypts a message to a number of recipients and, optionally, signs
// it. hints contains optional information, that is also encrypted, that aids
// the recipients in processing the message. The resulting WriteCloser must
// be closed after the contents of the file have been written.
// If config is nil, sensible defaults will be used.
func Encrypt(ciphertext io.Writer, to []*Entity, signed *Entity, hints *FileHints, config *packet.Config) (plaintext io.WriteCloser, err error) {
	return encrypt(ciphertext, ciphertext, to, signed, hints, packet.SigTypeBinary, config)
}

// EncryptSplit encrypts a message to a number of recipients and, optionally, signs
// it. hints contains optional information, that is also encrypted, that aids
// the recipients in processing the message. The resulting WriteCloser must
// be closed after the contents of the file have been written.
// If config is nil, sensible defaults will be used.
func EncryptSplit(keyWriter io.Writer, dataWriter io.Writer, to []*Entity, signed *Entity, hints *FileHints, config *packet.Config) (plaintext io.WriteCloser, err error) {
	return encrypt(keyWriter, dataWriter, to, signed, hints, packet.SigTypeBinary, config)
}

// EncryptTextSplit encrypts a message to a number of recipients and, optionally, signs
// it. hints contains optional information, that is also encrypted, that aids
// the recipients in processing the message. The resulting WriteCloser must
// be closed after the contents of the file have been written.
// If config is nil, sensible defaults will be used.
func EncryptTextSplit(keyWriter io.Writer, dataWriter io.Writer, to []*Entity, signed *Entity, hints *FileHints, config *packet.Config) (plaintext io.WriteCloser, err error) {
	return encrypt(keyWriter, dataWriter, to, signed, hints, packet.SigTypeText, config)
}

// writeAndSign writes the data as a payload package and, optionally, signs
// it. hints contains optional information, that is also encrypted,
// that aids the recipients in processing the message. The resulting
// WriteCloser must be closed after the contents of the file have been
// written. If config is nil, sensible defaults will be used.
func writeAndSign(payload io.WriteCloser, candidateHashes []uint8, signed *Entity, hints *FileHints, sigType packet.SignatureType, config *packet.Config) (plaintext io.WriteCloser, err error) {
	var signer *packet.PrivateKey
	if signed != nil {
		signKey, ok := signed.SigningKeyById(config.Now(), config.SigningKey())
		if !ok {
			return nil, errors.InvalidArgumentError("no valid signing keys")
		}
		signer = signKey.PrivateKey
		if signer == nil {
			return nil, errors.InvalidArgumentError("no private key in signing key")
		}
		if signer.Encrypted {
			return nil, errors.InvalidArgumentError("signing key must be decrypted")
		}
	}

	var hash crypto.Hash
	for _, hashId := range candidateHashes {
		if h, ok := algorithm.HashIdToHash(hashId); ok && h.Available() {
			hash = h
			break
		}
	}

	// If the hash specified by config is a candidate, we'll use that.
	if configuredHash := config.Hash(); configuredHash.Available() {
		for _, hashId := range candidateHashes {
			if h, ok := algorithm.HashIdToHash(hashId); ok && h == configuredHash {
				hash = h
				break
			}
		}
	}

	if hash == 0 {
		hashId := candidateHashes[0]
		name, ok := algorithm.HashIdToString(hashId)
		if !ok {
			name = "#" + strconv.Itoa(int(hashId))
		}
		return nil, errors.InvalidArgumentError("cannot encrypt because no candidate hash functions are compiled in. (Wanted " + name + " in this case.)")
	}

	var salt []byte
	if signer != nil {
		var opsVersion = 3
		if signer.Version == 6 {
			opsVersion = signer.Version
		}
		ops := &packet.OnePassSignature{
			Version:    opsVersion,
			SigType:    sigType,
			Hash:       hash,
			PubKeyAlgo: signer.PubKeyAlgo,
			KeyId:      signer.KeyId,
			IsLast:     true,
		}
		if opsVersion == 6 {
			ops.KeyFingerprint = signer.Fingerprint
			salt, err = packet.SignatureSaltForHash(hash, config.Random())
			if err != nil {
				return nil, err
			}
			ops.Salt = salt
		}
		if err := ops.Serialize(payload); err != nil {
			return nil, err
		}
	}

	if hints == nil {
		hints = &FileHints{}
	}

	w := payload
	if signer != nil {
		// If we need to write a signature packet after the literal
		// data then we need to stop literalData from closing
		// encryptedData.
		w = noOpCloser{w}

	}
	var epochSeconds uint32
	if !hints.ModTime.IsZero() {
		epochSeconds = uint32(hints.ModTime.Unix())
	}
	literalData, err := packet.SerializeLiteral(w, hints.IsBinary, hints.FileName, epochSeconds)
	if err != nil {
		return nil, err
	}

	if signer != nil {
		h, wrappedHash, err := hashForSignature(hash, sigType, salt)
		if err != nil {
			return nil, err
		}
		metadata := &packet.LiteralData{
			Format:   'u',
			FileName: hints.FileName,
			Time:     epochSeconds,
		}
		if hints.IsBinary {
			metadata.Format = 'b'
		}
		return signatureWriter{payload, literalData, hash, wrappedHash, h, salt, signer, sigType, config, metadata}, nil
	}
	return literalData, nil
}

// encrypt encrypts a message to a number of recipients and, optionally, signs
// it. hints contains optional information, that is also encrypted, that aids
// the recipients in processing the message. The resulting WriteCloser must
// be closed after the contents of the file have been written.
// If config is nil, sensible defaults will be used.
func encrypt(keyWriter io.Writer, dataWriter io.Writer, to []*Entity, signed *Entity, hints *FileHints, sigType packet.SignatureType, config *packet.Config) (plaintext io.WriteCloser, err error) {
	if len(to) == 0 {
		return nil, errors.InvalidArgumentError("no encryption recipient provided")
	}

	// These are the possible ciphers that we'll use for the message.
	candidateCiphers := []uint8{
		uint8(packet.CipherAES256),
		uint8(packet.CipherAES128),
	}

	// These are the possible hash functions that we'll use for the signature.
	candidateHashes := []uint8{
		hashToHashId(crypto.SHA256),
		hashToHashId(crypto.SHA384),
		hashToHashId(crypto.SHA512),
		hashToHashId(crypto.SHA3_256),
		hashToHashId(crypto.SHA3_512),
	}

	// Prefer GCM if everyone supports it
	candidateCipherSuites := [][2]uint8{
		{uint8(packet.CipherAES256), uint8(packet.AEADModeGCM)},
		{uint8(packet.CipherAES256), uint8(packet.AEADModeEAX)},
		{uint8(packet.CipherAES256), uint8(packet.AEADModeOCB)},
		{uint8(packet.CipherAES128), uint8(packet.AEADModeGCM)},
		{uint8(packet.CipherAES128), uint8(packet.AEADModeEAX)},
		{uint8(packet.CipherAES128), uint8(packet.AEADModeOCB)},
	}

	candidateCompression := []uint8{
		uint8(packet.CompressionNone),
		uint8(packet.CompressionZIP),
		uint8(packet.CompressionZLIB),
	}

	encryptKeys := make([]Key, len(to))

	// AEAD is used only if config enables it and every key supports it
	aeadSupported := config.AEAD() != nil

	for i := range to {
		var ok bool
		encryptKeys[i], ok = to[i].EncryptionKey(config.Now())
		if !ok {
			return nil, errors.InvalidArgumentError("cannot encrypt a message to key id " + strconv.FormatUint(to[i].PrimaryKey.KeyId, 16) + " because it has no valid encryption keys")
		}

		primarySelfSignature, _ := to[i].PrimarySelfSignature()
		if primarySelfSignature == nil {
			return nil, errors.InvalidArgumentError("entity without a self-signature")
		}

		if !primarySelfSignature.SEIPDv2 {
			aeadSupported = false
		}

		candidateCiphers = intersectPreferences(candidateCiphers, primarySelfSignature.PreferredSymmetric)
		candidateHashes = intersectPreferences(candidateHashes, primarySelfSignature.PreferredHash)
		candidateCipherSuites = intersectCipherSuites(candidateCipherSuites, primarySelfSignature.PreferredCipherSuites)
		candidateCompression = intersectPreferences(candidateCompression, primarySelfSignature.PreferredCompression)
	}

	// In the event that the intersection of supported algorithms is empty we use the ones
	// labelled as MUST that every implementation supports.
	if len(candidateCiphers) == 0 {
		// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#section-9.3
		candidateCiphers = []uint8{uint8(packet.CipherAES128)}
	}
	if len(candidateHashes) == 0 {
		// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#hash-algos
		candidateHashes = []uint8{hashToHashId(crypto.SHA256)}
	}
	if len(candidateCipherSuites) == 0 {
		// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#section-9.6
		candidateCipherSuites = [][2]uint8{{uint8(packet.CipherAES128), uint8(packet.AEADModeOCB)}}
	}

	cipher := packet.CipherFunction(candidateCiphers[0])
	aeadCipherSuite := packet.CipherSuite{
		Cipher: packet.CipherFunction(candidateCipherSuites[0][0]),
		Mode:   packet.AEADMode(candidateCipherSuites[0][1]),
	}

	// If the cipher specified by config is a candidate, we'll use that.
	configuredCipher := config.Cipher()
	for _, c := range candidateCiphers {
		cipherFunc := packet.CipherFunction(c)
		if cipherFunc == configuredCipher {
			cipher = cipherFunc
			break
		}
	}

	var symKey []byte
	if aeadSupported {
		symKey = make([]byte, aeadCipherSuite.Cipher.KeySize())
	} else {
		symKey = make([]byte, cipher.KeySize())
	}

	if _, err := io.ReadFull(config.Random(), symKey); err != nil {
		return nil, err
	}

	for _, key := range encryptKeys {
		if err := packet.SerializeEncryptedKeyAEAD(keyWriter, key.PublicKey, cipher, aeadSupported, symKey, config); err != nil {
			return nil, err
		}
	}

	var payload io.WriteCloser
	payload, err = packet.SerializeSymmetricallyEncrypted(dataWriter, cipher, aeadSupported, aeadCipherSuite, symKey, config)
	if err != nil {
		return
	}

	payload, err = handleCompression(payload, candidateCompression, config)
	if err != nil {
		return nil, err
	}

	return writeAndSign(payload, candidateHashes, signed, hints, sigType, config)
}

// Sign signs a message. The resulting WriteCloser must be closed after the
// contents of the file have been written.  hints contains optional information
// that aids the recipients in processing the message.
// If config is nil, sensible defaults will be used.
func Sign(output io.Writer, signed *Entity, hints *FileHints, config *packet.Config) (input io.WriteCloser, err error) {
	if signed == nil {
		return nil, errors.InvalidArgumentError("no signer provided")
	}

	// These are the possible hash functions that we'll use for the signature.
	candidateHashes := []uint8{
		hashToHashId(crypto.SHA256),
		hashToHashId(crypto.SHA384),
		hashToHashId(crypto.SHA512),
		hashToHashId(crypto.SHA3_256),
		hashToHashId(crypto.SHA3_512),
	}
	defaultHashes := candidateHashes[0:1]
	primarySelfSignature, _ := signed.PrimarySelfSignature()
	if primarySelfSignature == nil {
		return nil, errors.StructuralError("signed entity has no self-signature")
	}
	preferredHashes := primarySelfSignature.PreferredHash
	if len(preferredHashes) == 0 {
		preferredHashes = defaultHashes
	}
	candidateHashes = intersectPreferences(candidateHashes, preferredHashes)
	if len(candidateHashes) == 0 {
		return nil, errors.StructuralError("cannot sign because signing key shares no common algorithms with candidate hashes")
	}

	return writeAndSign(noOpCloser{output}, candidateHashes, signed, hints, packet.SigTypeBinary, config)
}

// signatureWriter hashes the contents of a message while passing it along to
// literalData. When closed, it closes literalData, writes a signature packet
// to encryptedData and then also closes encryptedData.
type signatureWriter struct {
	encryptedData io.WriteCloser
	literalData   io.WriteCloser
	hashType      crypto.Hash
	wrappedHash   hash.Hash
	h             hash.Hash
	salt          []byte // v6 only
	signer        *packet.PrivateKey
	sigType       packet.SignatureType
	config        *packet.Config
	metadata      *packet.LiteralData // V5 signatures protect document metadata
}

func (s signatureWriter) Write(data []byte) (int, error) {
	s.wrappedHash.Write(data)
	switch s.sigType {
	case packet.SigTypeBinary:
		return s.literalData.Write(data)
	case packet.SigTypeText:
		flag := 0
		return writeCanonical(s.literalData, data, &flag)
	}
	return 0, errors.UnsupportedError("unsupported signature type: " + strconv.Itoa(int(s.sigType)))
}

func (s signatureWriter) Close() error {
	sig := createSignaturePacket(&s.signer.PublicKey, s.sigType, s.config)
	sig.Hash = s.hashType
	sig.Metadata = s.metadata

	if err := sig.SetSalt(s.salt); err != nil {
		return err
	}

	if err := sig.Sign(s.h, s.signer, s.config); err != nil {
		return err
	}
	if err := s.literalData.Close(); err != nil {
		return err
	}
	if err := sig.Serialize(s.encryptedData); err != nil {
		return err
	}
	return s.encryptedData.Close()
}

func createSignaturePacket(signer *packet.PublicKey, sigType packet.SignatureType, config *packet.Config) *packet.Signature {
	sigLifetimeSecs := config.SigLifetime()
	return &packet.Signature{
		Version:           signer.Version,
		SigType:           sigType,
		PubKeyAlgo:        signer.PubKeyAlgo,
		Hash:              config.Hash(),
		CreationTime:      config.Now(),
		IssuerKeyId:       &signer.KeyId,
		IssuerFingerprint: signer.Fingerprint,
		Notations:         config.Notations(),
		SigLifetimeSecs:   &sigLifetimeSecs,
	}
}

// noOpCloser is like an ioutil.NopCloser, but for an io.Writer.
// TODO: we have two of these in OpenPGP packages alone. This probably needs
// to be promoted somewhere more common.
type noOpCloser struct {
	w io.Writer
}

func (c noOpCloser) Write(data []byte) (n int, err error) {
	return c.w.Write(data)
}

func (c noOpCloser) Close() error {
	return nil
}

func handleCompression(compressed io.WriteCloser, candidateCompression []uint8, config *packet.Config) (data io.WriteCloser, err error) {
	data = compressed
	confAlgo := config.Compression()
	if confAlgo == packet.CompressionNone {
		return
	}

	// Set algorithm labelled as MUST as fallback
	// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#section-9.4
	finalAlgo := packet.CompressionNone
	// if compression specified by config available we will use it
	for _, c := range candidateCompression {
		if uint8(confAlgo) == c {
			finalAlgo = confAlgo
			break
		}
	}

	if finalAlgo != packet.CompressionNone {
		var compConfig *packet.CompressionConfig
		if config != nil {
			compConfig = config.CompressionConfig
		}
		data, err = packet.SerializeCompressed(compressed, finalAlgo, compConfig)
		if err != nil {
			return
		}
	}
	return data, nil
}
