// Copyright (C) 2019 ProtonTech AG

package packet

import "math/bits"

// CipherSuite contains a combination of Cipher and Mode
type CipherSuite struct {
	// The cipher function
	Cipher CipherFunction
	// The AEAD mode of operation.
	Mode AEADMode
}

// AEADConfig collects a number of AEAD parameters along with sensible defaults.
// A nil AEADConfig is valid and results in all default values.
type AEADConfig struct {
	// The AEAD mode of operation.
	DefaultMode AEADMode
	// Amount of octets in each chunk of data
	ChunkSize uint64
}

// Mode returns the AEAD mode of operation.
func (conf *AEADConfig) Mode() AEADMode {
	// If no preference is specified, OCB is used (which is mandatory to implement).
	if conf == nil || conf.DefaultMode == 0 {
		return AEADModeOCB
	}

	mode := conf.DefaultMode
	if mode != AEADModeEAX && mode != AEADModeOCB && mode != AEADModeGCM {
		panic("AEAD mode unsupported")
	}
	return mode
}

// ChunkSizeByte returns the byte indicating the chunk size. The effective
// chunk size is computed with the formula uint64(1) << (chunkSizeByte + 6)
// limit to 16 = 4 MiB
// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#section-5.13.2
func (conf *AEADConfig) ChunkSizeByte() byte {
	if conf == nil || conf.ChunkSize == 0 {
		return 12 // 1 << (12 + 6) == 262144 bytes
	}

	chunkSize := conf.ChunkSize
	exponent := bits.Len64(chunkSize) - 1
	switch {
	case exponent < 6:
		exponent = 6
	case exponent > 16:
		exponent = 16
	}

	return byte(exponent - 6)
}

// decodeAEADChunkSize returns the effective chunk size. In 32-bit systems, the
// maximum returned value is 1 << 30.
func decodeAEADChunkSize(c byte) int {
	size := uint64(1 << (c + 6))
	if size != uint64(int(size)) {
		return 1 << 30
	}
	return int(size)
}
