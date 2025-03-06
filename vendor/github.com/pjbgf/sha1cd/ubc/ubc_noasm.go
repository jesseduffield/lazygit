//go:build !amd64 || noasm || !gc
// +build !amd64 noasm !gc

package ubc

// Check takes as input an expanded message block and verifies the unavoidable bitconditions
// for all listed DVs. It returns a dvmask where each bit belonging to a DV is set if all
// unavoidable bitconditions for that DV have been met.
// Thus, one needs to do the recompression check for each DV that has its bit set.
func CalculateDvMask(W [80]uint32) uint32 {
	return CalculateDvMaskGeneric(W)
}
