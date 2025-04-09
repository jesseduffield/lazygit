package s2k

import "crypto"

// Config collects configuration parameters for s2k key-stretching
// transformations. A nil *Config is valid and results in all default
// values.
type Config struct {
	// S2K (String to Key) mode, used for key derivation in the context of secret key encryption
	// and passphrase-encrypted data. Either s2k.Argon2S2K or s2k.IteratedSaltedS2K may be used.
	// If the passphrase is a high-entropy key, indicated by setting PassphraseIsHighEntropy to true,
	// s2k.SaltedS2K can also be used.
	// Note: Argon2 is the strongest option but not all OpenPGP implementations are compatible with it
	//(pending standardisation).
	// 0 (simple), 1(salted), 3(iterated), 4(argon2)
	// 2(reserved) 100-110(private/experimental).
	S2KMode Mode
	// Only relevant if S2KMode is not set to s2k.Argon2S2K.
	// Hash is the default hash function to be used. If
	// nil, SHA256 is used.
	Hash crypto.Hash
	// Argon2 parameters for S2K (String to Key).
	// Only relevant if S2KMode is set to s2k.Argon2S2K.
	// If nil, default parameters are used.
	// For more details on the choice of parameters, see https://tools.ietf.org/html/rfc9106#section-4.
	Argon2Config *Argon2Config
	// Only relevant if S2KMode is set to s2k.IteratedSaltedS2K.
	// Iteration count for Iterated S2K (String to Key). It
	// determines the strength of the passphrase stretching when
	// the said passphrase is hashed to produce a key. S2KCount
	// should be between 65536 and 65011712, inclusive. If Config
	// is nil or S2KCount is 0, the value 16777216 used. Not all
	// values in the above range can be represented. S2KCount will
	// be rounded up to the next representable value if it cannot
	// be encoded exactly. When set, it is strongly encrouraged to
	// use a value that is at least 65536. See RFC 4880 Section
	// 3.7.1.3.
	S2KCount int
	// Indicates whether the passphrase passed by the application is a
	// high-entropy key (e.g. it's randomly generated or derived from
	// another passphrase using a strong key derivation function).
	// When true, allows the S2KMode to be s2k.SaltedS2K.
	// When the passphrase is not a high-entropy key, using SaltedS2K is
	// insecure, and not allowed by draft-ietf-openpgp-crypto-refresh-08.
	PassphraseIsHighEntropy bool
}

// Argon2Config stores the Argon2 parameters
// A nil *Argon2Config is valid and results in all default
type Argon2Config struct {
	NumberOfPasses      uint8
	DegreeOfParallelism uint8
	// Memory specifies the desired Argon2 memory usage in kibibytes.
	// For example memory=64*1024 sets the memory cost to ~64 MB.
	Memory uint32
}

func (c *Config) Mode() Mode {
	if c == nil {
		return IteratedSaltedS2K
	}
	return c.S2KMode
}

func (c *Config) hash() crypto.Hash {
	if c == nil || uint(c.Hash) == 0 {
		return crypto.SHA256
	}

	return c.Hash
}

func (c *Config) Argon2() *Argon2Config {
	if c == nil || c.Argon2Config == nil {
		return nil
	}
	return c.Argon2Config
}

// EncodedCount get encoded count
func (c *Config) EncodedCount() uint8 {
	if c == nil || c.S2KCount == 0 {
		return 224 // The common case. Corresponding to 16777216
	}

	i := c.S2KCount

	switch {
	case i < 65536:
		i = 65536
	case i > 65011712:
		i = 65011712
	}

	return encodeCount(i)
}

func (c *Argon2Config) Passes() uint8 {
	if c == nil || c.NumberOfPasses == 0 {
		return 3
	}
	return c.NumberOfPasses
}

func (c *Argon2Config) Parallelism() uint8 {
	if c == nil || c.DegreeOfParallelism == 0 {
		return 4
	}
	return c.DegreeOfParallelism
}

func (c *Argon2Config) EncodedMemory() uint8 {
	if c == nil || c.Memory == 0 {
		return 16 // 64 MiB of RAM
	}

	memory := c.Memory
	lowerBound := uint32(c.Parallelism()) * 8
	upperBound := uint32(2147483648)

	switch {
	case memory < lowerBound:
		memory = lowerBound
	case memory > upperBound:
		memory = upperBound
	}

	return encodeMemory(memory, c.Parallelism())
}
