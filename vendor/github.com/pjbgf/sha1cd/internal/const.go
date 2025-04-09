package shared

const (
	// Constants for the SHA-1 hash function.
	K0 = 0x5A827999
	K1 = 0x6ED9EBA1
	K2 = 0x8F1BBCDC
	K3 = 0xCA62C1D6

	// Initial values for the buffer variables: h0, h1, h2, h3, h4.
	Init0 = 0x67452301
	Init1 = 0xEFCDAB89
	Init2 = 0x98BADCFE
	Init3 = 0x10325476
	Init4 = 0xC3D2E1F0

	// Initial values for the temporary variables (ihvtmp0, ihvtmp1, ihvtmp2, ihvtmp3, ihvtmp4) during the SHA recompression step.
	InitTmp0 = 0xD5
	InitTmp1 = 0x394
	InitTmp2 = 0x8152A8
	InitTmp3 = 0x0
	InitTmp4 = 0xA7ECE0

	// SHA1 contains 2 buffers, each based off 5 32-bit words.
	WordBuffers = 5

	// The output of SHA1 is 20 bytes (160 bits).
	Size = 20

	// Rounds represents the number of steps required to process each chunk.
	Rounds = 80

	// SHA1 processes the input data in chunks. Each chunk contains 64 bytes.
	Chunk = 64

	// The number of pre-step compression state to store.
	// Currently there are 3 pre-step compression states required: 0, 58, 65.
	PreStepState = 3

	Magic         = "shacd\x01"
	MarshaledSize = len(Magic) + 5*4 + Chunk + 8
)
