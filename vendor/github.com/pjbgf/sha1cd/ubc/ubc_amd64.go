//go:build !noasm && gc && amd64
// +build !noasm,gc,amd64

package ubc

func CalculateDvMaskAMD64(W [80]uint32) uint32

// Check takes as input an expanded message block and verifies the unavoidable bitconditions
// for all listed DVs. It returns a dvmask where each bit belonging to a DV is set if all
// unavoidable bitconditions for that DV have been met.
// Thus, one needs to do the recompression check for each DV that has its bit set.
func CalculateDvMask(W [80]uint32) uint32 {
	return CalculateDvMaskAMD64(W)
}
