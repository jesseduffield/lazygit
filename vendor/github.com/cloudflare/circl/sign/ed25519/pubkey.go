//go:build go1.13
// +build go1.13

package ed25519

import cryptoEd25519 "crypto/ed25519"

// PublicKey is the type of Ed25519 public keys.
type PublicKey cryptoEd25519.PublicKey
