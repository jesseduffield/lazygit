package object

import "bytes"

const (
	signatureTypeUnknown signatureType = iota
	signatureTypeOpenPGP
	signatureTypeX509
	signatureTypeSSH
)

var (
	// openPGPSignatureFormat is the format of an OpenPGP signature.
	openPGPSignatureFormat = signatureFormat{
		[]byte("-----BEGIN PGP SIGNATURE-----"),
		[]byte("-----BEGIN PGP MESSAGE-----"),
	}
	// x509SignatureFormat is the format of an X509 signature, which is
	// a PKCS#7 (S/MIME) signature.
	x509SignatureFormat = signatureFormat{
		[]byte("-----BEGIN CERTIFICATE-----"),
		[]byte("-----BEGIN SIGNED MESSAGE-----"),
	}

	// sshSignatureFormat is the format of an SSH signature.
	sshSignatureFormat = signatureFormat{
		[]byte("-----BEGIN SSH SIGNATURE-----"),
	}
)

var (
	// knownSignatureFormats is a map of known signature formats, indexed by
	// their signatureType.
	knownSignatureFormats = map[signatureType]signatureFormat{
		signatureTypeOpenPGP: openPGPSignatureFormat,
		signatureTypeX509:    x509SignatureFormat,
		signatureTypeSSH:     sshSignatureFormat,
	}
)

// signatureType represents the type of the signature.
type signatureType int8

// signatureFormat represents the beginning of a signature.
type signatureFormat [][]byte

// typeForSignature returns the type of the signature based on its format.
func typeForSignature(b []byte) signatureType {
	for t, i := range knownSignatureFormats {
		for _, begin := range i {
			if bytes.HasPrefix(b, begin) {
				return t
			}
		}
	}
	return signatureTypeUnknown
}

// parseSignedBytes returns the position of the last signature block found in
// the given bytes. If no signature block is found, it returns -1.
//
// When multiple signature blocks are found, the position of the last one is
// returned. Any tailing bytes after this signature block start should be
// considered part of the signature.
//
// Given this, it would be safe to use the returned position to split the bytes
// into two parts: the first part containing the message, the second part
// containing the signature.
//
// Example:
//
//	message := []byte(`Message with signature
//
//	-----BEGIN SSH SIGNATURE-----
//	...`)
//
//	var signature string
//	if pos, _ := parseSignedBytes(message); pos != -1 {
//		signature = string(message[pos:])
//		message = message[:pos]
//	}
//
// This logic is on par with git's gpg-interface.c:parse_signed_buffer().
// https://github.com/git/git/blob/7c2ef319c52c4997256f5807564523dfd4acdfc7/gpg-interface.c#L668
func parseSignedBytes(b []byte) (int, signatureType) {
	var n, match = 0, -1
	var t signatureType
	for n < len(b) {
		var i = b[n:]
		if st := typeForSignature(i); st != signatureTypeUnknown {
			match = n
			t = st
		}
		if eol := bytes.IndexByte(i, '\n'); eol >= 0 {
			n += eol + 1
			continue
		}
		// If we reach this point, we've reached the end.
		break
	}
	return match, t
}
