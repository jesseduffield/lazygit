// ubc package provides ways for SHA1 blocks to be checked for
// Unavoidable Bit Conditions that arise from crypto analysis attacks.
package ubc

//go:generate go run -C asm . -out ../ubc_amd64.s -pkg $GOPACKAGE
