// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chacha20poly1305

import (
	"bytes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	mathrand "math/rand"
	"strconv"
	"testing"
)

func TestVectors(t *testing.T) {
	for i, test := range chacha20Poly1305Tests {
		key, _ := hex.DecodeString(test.key)
		nonce, _ := hex.DecodeString(test.nonce)
		ad, _ := hex.DecodeString(test.aad)
		plaintext, _ := hex.DecodeString(test.plaintext)

		var (
			aead cipher.AEAD
			err  error
		)
		switch len(nonce) {
		case NonceSize:
			aead, err = New(key)
		case NonceSizeX:
			aead, err = NewX(key)
		default:
			t.Fatalf("#%d: wrong nonce length: %d", i, len(nonce))
		}
		if err != nil {
			t.Fatal(err)
		}

		ct := aead.Seal(nil, nonce, plaintext, ad)
		if ctHex := hex.EncodeToString(ct); ctHex != test.out {
			t.Errorf("#%d: got %s, want %s", i, ctHex, test.out)
			continue
		}

		plaintext2, err := aead.Open(nil, nonce, ct, ad)
		if err != nil {
			t.Errorf("#%d: Open failed", i)
			continue
		}

		if !bytes.Equal(plaintext, plaintext2) {
			t.Errorf("#%d: plaintext's don't match: got %x vs %x", i, plaintext2, plaintext)
			continue
		}

		if len(ad) > 0 {
			alterAdIdx := mathrand.Intn(len(ad))
			ad[alterAdIdx] ^= 0x80
			if _, err := aead.Open(nil, nonce, ct, ad); err == nil {
				t.Errorf("#%d: Open was successful after altering additional data", i)
			}
			ad[alterAdIdx] ^= 0x80
		}

		alterNonceIdx := mathrand.Intn(aead.NonceSize())
		nonce[alterNonceIdx] ^= 0x80
		if _, err := aead.Open(nil, nonce, ct, ad); err == nil {
			t.Errorf("#%d: Open was successful after altering nonce", i)
		}
		nonce[alterNonceIdx] ^= 0x80

		alterCtIdx := mathrand.Intn(len(ct))
		ct[alterCtIdx] ^= 0x80
		if _, err := aead.Open(nil, nonce, ct, ad); err == nil {
			t.Errorf("#%d: Open was successful after altering ciphertext", i)
		}
		ct[alterCtIdx] ^= 0x80
	}
}

func TestRandom(t *testing.T) {
	// Some random tests to verify Open(Seal) == Plaintext
	f := func(t *testing.T, nonceSize int) {
		for i := 0; i < 256; i++ {
			var nonce = make([]byte, nonceSize)
			var key [32]byte

			al := mathrand.Intn(128)
			pl := mathrand.Intn(16384)
			ad := make([]byte, al)
			plaintext := make([]byte, pl)
			cryptorand.Read(key[:])
			cryptorand.Read(nonce[:])
			cryptorand.Read(ad)
			cryptorand.Read(plaintext)

			var (
				aead cipher.AEAD
				err  error
			)
			switch len(nonce) {
			case NonceSize:
				aead, err = New(key[:])
			case NonceSizeX:
				aead, err = NewX(key[:])
			default:
				t.Fatalf("#%d: wrong nonce length: %d", i, len(nonce))
			}
			if err != nil {
				t.Fatal(err)
			}

			ct := aead.Seal(nil, nonce[:], plaintext, ad)

			plaintext2, err := aead.Open(nil, nonce[:], ct, ad)
			if err != nil {
				t.Errorf("Random #%d: Open failed", i)
				continue
			}

			if !bytes.Equal(plaintext, plaintext2) {
				t.Errorf("Random #%d: plaintext's don't match: got %x vs %x", i, plaintext2, plaintext)
				continue
			}

			if len(ad) > 0 {
				alterAdIdx := mathrand.Intn(len(ad))
				ad[alterAdIdx] ^= 0x80
				if _, err := aead.Open(nil, nonce[:], ct, ad); err == nil {
					t.Errorf("Random #%d: Open was successful after altering additional data", i)
				}
				ad[alterAdIdx] ^= 0x80
			}

			alterNonceIdx := mathrand.Intn(aead.NonceSize())
			nonce[alterNonceIdx] ^= 0x80
			if _, err := aead.Open(nil, nonce[:], ct, ad); err == nil {
				t.Errorf("Random #%d: Open was successful after altering nonce", i)
			}
			nonce[alterNonceIdx] ^= 0x80

			alterCtIdx := mathrand.Intn(len(ct))
			ct[alterCtIdx] ^= 0x80
			if _, err := aead.Open(nil, nonce[:], ct, ad); err == nil {
				t.Errorf("Random #%d: Open was successful after altering ciphertext", i)
			}
			ct[alterCtIdx] ^= 0x80
		}
	}
	t.Run("Standard", func(t *testing.T) { f(t, NonceSize) })
	t.Run("X", func(t *testing.T) { f(t, NonceSizeX) })
}

func benchamarkChaCha20Poly1305Seal(b *testing.B, buf []byte, nonceSize int) {
	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	var key [32]byte
	var nonce = make([]byte, nonceSize)
	var ad [13]byte
	var out []byte

	var aead cipher.AEAD
	switch len(nonce) {
	case NonceSize:
		aead, _ = New(key[:])
	case NonceSizeX:
		aead, _ = NewX(key[:])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out = aead.Seal(out[:0], nonce[:], buf[:], ad[:])
	}
}

func benchamarkChaCha20Poly1305Open(b *testing.B, buf []byte, nonceSize int) {
	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	var key [32]byte
	var nonce = make([]byte, nonceSize)
	var ad [13]byte
	var ct []byte
	var out []byte

	var aead cipher.AEAD
	switch len(nonce) {
	case NonceSize:
		aead, _ = New(key[:])
	case NonceSizeX:
		aead, _ = NewX(key[:])
	}
	ct = aead.Seal(ct[:0], nonce[:], buf[:], ad[:])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, _ = aead.Open(out[:0], nonce[:], ct[:], ad[:])
	}
}

func BenchmarkChacha20Poly1305(b *testing.B) {
	for _, length := range []int{64, 1350, 8 * 1024} {
		b.Run("Open-"+strconv.Itoa(length), func(b *testing.B) {
			benchamarkChaCha20Poly1305Open(b, make([]byte, length), NonceSize)
		})
		b.Run("Seal-"+strconv.Itoa(length), func(b *testing.B) {
			benchamarkChaCha20Poly1305Seal(b, make([]byte, length), NonceSize)
		})

		b.Run("Open-"+strconv.Itoa(length)+"-X", func(b *testing.B) {
			benchamarkChaCha20Poly1305Open(b, make([]byte, length), NonceSizeX)
		})
		b.Run("Seal-"+strconv.Itoa(length)+"-X", func(b *testing.B) {
			benchamarkChaCha20Poly1305Seal(b, make([]byte, length), NonceSizeX)
		})
	}
}

var key = make([]byte, KeySize)

func ExampleNewX() {
	aead, err := NewX(key)
	if err != nil {
		log.Fatalln("Failed to instantiate XChaCha20-Poly1305:", err)
	}

	for _, msg := range []string{
		"Attack at dawn.",
		"The eagle has landed.",
		"Gophers, gophers, gophers everywhere!",
	} {
		// Encryption.
		nonce := make([]byte, NonceSizeX)
		if _, err := cryptorand.Read(nonce); err != nil {
			panic(err)
		}
		ciphertext := aead.Seal(nil, nonce, []byte(msg), nil)

		// Decryption.
		plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			log.Fatalln("Failed to decrypt or authenticate message:", err)
		}

		fmt.Printf("%s\n", plaintext)
	}

	// Output: Attack at dawn.
	// The eagle has landed.
	// Gophers, gophers, gophers everywhere!
}
