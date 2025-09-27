// Package knownhosts is a thin wrapper around golang.org/x/crypto/ssh/knownhosts,
// adding the ability to obtain the list of host key algorithms for a known host.
package knownhosts

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strings"

	"golang.org/x/crypto/ssh"
	xknownhosts "golang.org/x/crypto/ssh/knownhosts"
)

// HostKeyDB wraps logic in golang.org/x/crypto/ssh/knownhosts with additional
// behaviors, such as the ability to perform host key/algorithm lookups from
// known_hosts entries.
type HostKeyDB struct {
	callback   ssh.HostKeyCallback
	isCert     map[string]bool // keyed by "filename:line"
	isWildcard map[string]bool // keyed by "filename:line"
}

// NewDB creates a HostKeyDB from the given OpenSSH known_hosts file(s). It
// reads and parses the provided files one additional time (beyond logic in
// golang.org/x/crypto/ssh/knownhosts) in order to:
//
//   - Handle CA lines properly and return ssh.CertAlgo* values when calling the
//     HostKeyAlgorithms method, for use in ssh.ClientConfig.HostKeyAlgorithms
//   - Allow * wildcards in hostnames to match on non-standard ports, providing
//     a workaround for https://github.com/golang/go/issues/52056 in order to
//     align with OpenSSH's wildcard behavior
//
// When supplying multiple files, their order does not matter.
func NewDB(files ...string) (*HostKeyDB, error) {
	cb, err := xknownhosts.New(files...)
	if err != nil {
		return nil, err
	}
	hkdb := &HostKeyDB{
		callback:   cb,
		isCert:     make(map[string]bool),
		isWildcard: make(map[string]bool),
	}

	// Re-read each file a single time, looking for @cert-authority lines. The
	// logic for reading the file is designed to mimic hostKeyDB.Read from
	// golang.org/x/crypto/ssh/knownhosts
	for _, filename := range files {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Bytes()
			line = bytes.TrimSpace(line)
			// Does the line start with "@cert-authority" followed by whitespace?
			if len(line) > 15 && bytes.HasPrefix(line, []byte("@cert-authority")) && (line[15] == ' ' || line[15] == '\t') {
				mapKey := fmt.Sprintf("%s:%d", filename, lineNum)
				hkdb.isCert[mapKey] = true
				line = bytes.TrimSpace(line[16:])
			}
			// truncate line to just the host pattern field
			if i := bytes.IndexAny(line, "\t "); i >= 0 {
				line = line[:i]
			}
			// Does the host pattern contain a * wildcard and no specific port?
			if i := bytes.IndexRune(line, '*'); i >= 0 && !bytes.Contains(line[i:], []byte("]:")) {
				mapKey := fmt.Sprintf("%s:%d", filename, lineNum)
				hkdb.isWildcard[mapKey] = true
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("knownhosts: %s:%d: %w", filename, lineNum, err)
		}
	}
	return hkdb, nil
}

// HostKeyCallback returns an ssh.HostKeyCallback. This can be used directly in
// ssh.ClientConfig.HostKeyCallback, as shown in the example for NewDB.
// Alternatively, you can wrap it with an outer callback to potentially handle
// appending a new entry to the known_hosts file; see example in WriteKnownHost.
func (hkdb *HostKeyDB) HostKeyCallback() ssh.HostKeyCallback {
	// Either NewDB found no wildcard host patterns, or hkdb was created from
	// HostKeyCallback.ToDB in which case we didn't scan known_hosts for them:
	// return the callback (which came from x/crypto/ssh/knownhosts) as-is
	if len(hkdb.isWildcard) == 0 {
		return hkdb.callback
	}

	// If we scanned for wildcards and found at least one, return a wrapped
	// callback with extra behavior: if the host lookup found no matches, and the
	// host arg had a non-standard port, re-do the lookup on standard port 22. If
	// that second call returns a *xknownhosts.KeyError, filter down any resulting
	// Want keys to known wildcard entries.
	f := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		callbackErr := hkdb.callback(hostname, remote, key)
		if callbackErr == nil || IsHostKeyChanged(callbackErr) { // hostname has known_host entries as-is
			return callbackErr
		}
		justHost, port, splitErr := net.SplitHostPort(hostname)
		if splitErr != nil || port == "" || port == "22" { // hostname already using standard port
			return callbackErr
		}
		// If we reach here, the port was non-standard and no known_host entries
		// were found for the non-standard port. Try again with standard port.
		if tcpAddr, ok := remote.(*net.TCPAddr); ok && tcpAddr.Port != 22 {
			remote = &net.TCPAddr{
				IP:   tcpAddr.IP,
				Port: 22,
				Zone: tcpAddr.Zone,
			}
		}
		callbackErr = hkdb.callback(justHost+":22", remote, key)
		var keyErr *xknownhosts.KeyError
		if errors.As(callbackErr, &keyErr) && len(keyErr.Want) > 0 {
			wildcardKeys := make([]xknownhosts.KnownKey, 0, len(keyErr.Want))
			for _, wantKey := range keyErr.Want {
				if hkdb.isWildcard[fmt.Sprintf("%s:%d", wantKey.Filename, wantKey.Line)] {
					wildcardKeys = append(wildcardKeys, wantKey)
				}
			}
			callbackErr = &xknownhosts.KeyError{
				Want: wildcardKeys,
			}
		}
		return callbackErr
	}
	return ssh.HostKeyCallback(f)
}

// PublicKey wraps ssh.PublicKey with an additional field, to identify
// whether the key corresponds to a certificate authority.
type PublicKey struct {
	ssh.PublicKey
	Cert bool
}

// HostKeys returns a slice of known host public keys for the supplied host:port
// found in the known_hosts file(s), or an empty slice if the host is not
// already known. For hosts that have multiple known_hosts entries (for
// different key types), the result will be sorted by known_hosts filename and
// line number.
// If hkdb was originally created by calling NewDB, the Cert boolean field of
// each result entry reports whether the key corresponded to a @cert-authority
// line. If hkdb was NOT obtained from NewDB, then Cert will always be false.
func (hkdb *HostKeyDB) HostKeys(hostWithPort string) (keys []PublicKey) {
	var keyErr *xknownhosts.KeyError
	placeholderAddr := &net.TCPAddr{IP: []byte{0, 0, 0, 0}}
	placeholderPubKey := &fakePublicKey{}
	var kkeys []xknownhosts.KnownKey
	callback := hkdb.HostKeyCallback()
	if hkcbErr := callback(hostWithPort, placeholderAddr, placeholderPubKey); errors.As(hkcbErr, &keyErr) {
		kkeys = append(kkeys, keyErr.Want...)
		knownKeyLess := func(i, j int) bool {
			if kkeys[i].Filename < kkeys[j].Filename {
				return true
			}
			return (kkeys[i].Filename == kkeys[j].Filename && kkeys[i].Line < kkeys[j].Line)
		}
		sort.Slice(kkeys, knownKeyLess)
		keys = make([]PublicKey, len(kkeys))
		for n := range kkeys {
			keys[n] = PublicKey{
				PublicKey: kkeys[n].Key,
			}
			if len(hkdb.isCert) > 0 {
				keys[n].Cert = hkdb.isCert[fmt.Sprintf("%s:%d", kkeys[n].Filename, kkeys[n].Line)]
			}
		}
	}
	return keys
}

// HostKeyAlgorithms returns a slice of host key algorithms for the supplied
// host:port found in the known_hosts file(s), or an empty slice if the host
// is not already known. The result may be used in ssh.ClientConfig's
// HostKeyAlgorithms field, either as-is or after filtering (if you wish to
// ignore or prefer particular algorithms). For hosts that have multiple
// known_hosts entries (of different key types), the result will be sorted by
// known_hosts filename and line number.
// If hkdb was originally created by calling NewDB, any @cert-authority lines
// in the known_hosts file will properly be converted to the corresponding
// ssh.CertAlgo* values.
func (hkdb *HostKeyDB) HostKeyAlgorithms(hostWithPort string) (algos []string) {
	// We ensure that algos never contains duplicates. This is done for robustness
	// even though currently golang.org/x/crypto/ssh/knownhosts never exposes
	// multiple keys of the same type. This way our behavior here is unaffected
	// even if https://github.com/golang/go/issues/28870 is implemented, for
	// example by https://github.com/golang/crypto/pull/254.
	hostKeys := hkdb.HostKeys(hostWithPort)
	seen := make(map[string]struct{}, len(hostKeys))
	addAlgo := func(typ string, cert bool) {
		if cert {
			typ = keyTypeToCertAlgo(typ)
		}
		if _, already := seen[typ]; !already {
			algos = append(algos, typ)
			seen[typ] = struct{}{}
		}
	}
	for _, key := range hostKeys {
		typ := key.Type()
		if typ == ssh.KeyAlgoRSA {
			// KeyAlgoRSASHA256 and KeyAlgoRSASHA512 are only public key algorithms,
			// not public key formats, so they can't appear as a PublicKey.Type.
			// The corresponding PublicKey.Type is KeyAlgoRSA. See RFC 8332, Section 2.
			addAlgo(ssh.KeyAlgoRSASHA512, key.Cert)
			addAlgo(ssh.KeyAlgoRSASHA256, key.Cert)
		}
		addAlgo(typ, key.Cert)
	}
	return algos
}

func keyTypeToCertAlgo(keyType string) string {
	switch keyType {
	case ssh.KeyAlgoRSA:
		return ssh.CertAlgoRSAv01
	case ssh.KeyAlgoRSASHA256:
		return ssh.CertAlgoRSASHA256v01
	case ssh.KeyAlgoRSASHA512:
		return ssh.CertAlgoRSASHA512v01
	case ssh.KeyAlgoDSA:
		return ssh.CertAlgoDSAv01
	case ssh.KeyAlgoECDSA256:
		return ssh.CertAlgoECDSA256v01
	case ssh.KeyAlgoSKECDSA256:
		return ssh.CertAlgoSKECDSA256v01
	case ssh.KeyAlgoECDSA384:
		return ssh.CertAlgoECDSA384v01
	case ssh.KeyAlgoECDSA521:
		return ssh.CertAlgoECDSA521v01
	case ssh.KeyAlgoED25519:
		return ssh.CertAlgoED25519v01
	case ssh.KeyAlgoSKED25519:
		return ssh.CertAlgoSKED25519v01
	}
	return ""
}

// HostKeyCallback wraps ssh.HostKeyCallback with additional methods to
// perform host key and algorithm lookups from the known_hosts entries. It is
// otherwise identical to ssh.HostKeyCallback, and does not introduce any file-
// parsing behavior beyond what is in golang.org/x/crypto/ssh/knownhosts.
//
// In most situations, use HostKeyDB and its constructor NewDB instead of using
// the HostKeyCallback type. The HostKeyCallback type is only provided for
// backwards compatibility with older versions of this package, as well as for
// very strict situations where any extra known_hosts file-parsing is
// undesirable.
//
// Methods of HostKeyCallback do not provide any special treatment for
// @cert-authority lines, which will (incorrectly) look like normal non-CA host
// keys. Additionally, HostKeyCallback lacks the fix for applying * wildcard
// known_host entries to all ports, like OpenSSH's behavior.
type HostKeyCallback ssh.HostKeyCallback

// New creates a HostKeyCallback from the given OpenSSH known_hosts file(s). The
// returned value may be used in ssh.ClientConfig.HostKeyCallback by casting it
// to ssh.HostKeyCallback, or using its HostKeyCallback method. Otherwise, it
// operates the same as the New function in golang.org/x/crypto/ssh/knownhosts.
// When supplying multiple files, their order does not matter.
//
// In most situations, you should avoid this function, as the returned value
// lacks several enhanced behaviors. See doc comment for HostKeyCallback for
// more information. Instead, most callers should use NewDB to create a
// HostKeyDB, which includes these enhancements.
func New(files ...string) (HostKeyCallback, error) {
	cb, err := xknownhosts.New(files...)
	return HostKeyCallback(cb), err
}

// HostKeyCallback simply casts the receiver back to ssh.HostKeyCallback, for
// use in ssh.ClientConfig.HostKeyCallback.
func (hkcb HostKeyCallback) HostKeyCallback() ssh.HostKeyCallback {
	return ssh.HostKeyCallback(hkcb)
}

// ToDB converts the receiver into a HostKeyDB. However, the returned HostKeyDB
// lacks the enhanced behaviors described in the doc comment for NewDB: proper
// CA support, and wildcard matching on nonstandard ports.
//
// It is generally preferable to create a HostKeyDB by using NewDB. The ToDB
// method is only provided for situations in which the calling code needs to
// make the extra NewDB behaviors optional / user-configurable, perhaps for
// reasons of performance or code trust (since NewDB reads the known_host file
// an extra time, which may be undesirable in some strict situations). This way,
// callers can conditionally create a non-enhanced HostKeyDB by using New and
// ToDB. See code example.
func (hkcb HostKeyCallback) ToDB() *HostKeyDB {
	// This intentionally leaves the isCert and isWildcard map fields as nil, as
	// there is no way to retroactively populate them from just a HostKeyCallback.
	// Methods of HostKeyDB will skip any related enhanced behaviors accordingly.
	return &HostKeyDB{callback: ssh.HostKeyCallback(hkcb)}
}

// HostKeys returns a slice of known host public keys for the supplied host:port
// found in the known_hosts file(s), or an empty slice if the host is not
// already known. For hosts that have multiple known_hosts entries (for
// different key types), the result will be sorted by known_hosts filename and
// line number.
// In the returned values, there is no way to distinguish between CA keys
// (known_hosts lines beginning with @cert-authority) and regular keys. To do
// so, see NewDB and HostKeyDB.HostKeys instead.
func (hkcb HostKeyCallback) HostKeys(hostWithPort string) []ssh.PublicKey {
	annotatedKeys := hkcb.ToDB().HostKeys(hostWithPort)
	rawKeys := make([]ssh.PublicKey, len(annotatedKeys))
	for n, ak := range annotatedKeys {
		rawKeys[n] = ak.PublicKey
	}
	return rawKeys
}

// HostKeyAlgorithms returns a slice of host key algorithms for the supplied
// host:port found in the known_hosts file(s), or an empty slice if the host
// is not already known. The result may be used in ssh.ClientConfig's
// HostKeyAlgorithms field, either as-is or after filtering (if you wish to
// ignore or prefer particular algorithms). For hosts that have multiple
// known_hosts entries (for different key types), the result will be sorted by
// known_hosts filename and line number.
// The returned values will not include ssh.CertAlgo* values. If any
// known_hosts lines had @cert-authority prefixes, their original key algo will
// be returned instead. For proper CA support, see NewDB and
// HostKeyDB.HostKeyAlgorithms instead.
func (hkcb HostKeyCallback) HostKeyAlgorithms(hostWithPort string) (algos []string) {
	return hkcb.ToDB().HostKeyAlgorithms(hostWithPort)
}

// HostKeyAlgorithms is a convenience function for performing host key algorithm
// lookups on an ssh.HostKeyCallback directly. It is intended for use in code
// paths that stay with the New method of golang.org/x/crypto/ssh/knownhosts
// rather than this package's New or NewDB methods.
// The returned values will not include ssh.CertAlgo* values. If any
// known_hosts lines had @cert-authority prefixes, their original key algo will
// be returned instead. For proper CA support, see NewDB and
// HostKeyDB.HostKeyAlgorithms instead.
func HostKeyAlgorithms(cb ssh.HostKeyCallback, hostWithPort string) []string {
	return HostKeyCallback(cb).HostKeyAlgorithms(hostWithPort)
}

// IsHostKeyChanged returns a boolean indicating whether the error indicates
// the host key has changed. It is intended to be called on the error returned
// from invoking a host key callback, to check whether an SSH host is known.
func IsHostKeyChanged(err error) bool {
	var keyErr *xknownhosts.KeyError
	return errors.As(err, &keyErr) && len(keyErr.Want) > 0
}

// IsHostUnknown returns a boolean indicating whether the error represents an
// unknown host. It is intended to be called on the error returned from invoking
// a host key callback to check whether an SSH host is known.
func IsHostUnknown(err error) bool {
	var keyErr *xknownhosts.KeyError
	return errors.As(err, &keyErr) && len(keyErr.Want) == 0
}

// Normalize normalizes an address into the form used in known_hosts. This
// implementation includes a fix for https://github.com/golang/go/issues/53463
// and will omit brackets around ipv6 addresses on standard port 22.
func Normalize(address string) string {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		host = address
		port = "22"
	}
	entry := host
	if port != "22" {
		entry = "[" + entry + "]:" + port
	} else if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		entry = entry[1 : len(entry)-1]
	}
	return entry
}

// Line returns a line to append to the known_hosts files. This implementation
// uses the local patched implementation of Normalize in order to solve
// https://github.com/golang/go/issues/53463.
func Line(addresses []string, key ssh.PublicKey) string {
	var trimmed []string
	for _, a := range addresses {
		trimmed = append(trimmed, Normalize(a))
	}

	return strings.Join([]string{
		strings.Join(trimmed, ","),
		key.Type(),
		base64.StdEncoding.EncodeToString(key.Marshal()),
	}, " ")
}

// WriteKnownHost writes a known_hosts line to w for the supplied hostname,
// remote, and key. This is useful when writing a custom hostkey callback which
// wraps a callback obtained from this package to provide additional known_hosts
// management functionality. The hostname, remote, and key typically correspond
// to the callback's args. This function does not support writing
// @cert-authority lines.
func WriteKnownHost(w io.Writer, hostname string, remote net.Addr, key ssh.PublicKey) error {
	// Always include hostname; only also include remote if it isn't a zero value
	// and doesn't normalize to the same string as hostname.
	hostnameNormalized := Normalize(hostname)
	if strings.ContainsAny(hostnameNormalized, "\t ") {
		return fmt.Errorf("knownhosts: hostname '%s' contains spaces", hostnameNormalized)
	}
	addresses := []string{hostnameNormalized}
	remoteStrNormalized := Normalize(remote.String())
	if remoteStrNormalized != "[0.0.0.0]:0" && remoteStrNormalized != hostnameNormalized &&
		!strings.ContainsAny(remoteStrNormalized, "\t ") {
		addresses = append(addresses, remoteStrNormalized)
	}
	line := Line(addresses, key) + "\n"
	_, err := w.Write([]byte(line))
	return err
}

// WriteKnownHostCA writes a @cert-authority line to w for the supplied host
// name/pattern and key.
func WriteKnownHostCA(w io.Writer, hostPattern string, key ssh.PublicKey) error {
	encodedKey := base64.StdEncoding.EncodeToString(key.Marshal())
	_, err := fmt.Fprintf(w, "@cert-authority %s %s %s\n", hostPattern, key.Type(), encodedKey)
	return err
}

// fakePublicKey is used as part of the work-around for
// https://github.com/golang/go/issues/29286
type fakePublicKey struct{}

func (fakePublicKey) Type() string {
	return "fake-public-key"
}
func (fakePublicKey) Marshal() []byte {
	return []byte("fake public key")
}
func (fakePublicKey) Verify(_ []byte, _ *ssh.Signature) error {
	return errors.New("Verify called on placeholder key")
}
