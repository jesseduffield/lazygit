//go:build !amd64 || noasm || !gc
// +build !amd64 noasm !gc

package sha1cd

func block(dig *digest, p []byte) {
	blockGeneric(dig, p)
}
