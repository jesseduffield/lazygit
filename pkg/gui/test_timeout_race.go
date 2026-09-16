//go:build race

package gui

// The race detector makes everything run several times slower, so the
// recording watchdog needs a correspondingly longer timeout; otherwise it
// fires on tests that are merely slow under -race rather than actually stuck.
// The `race` build tag is set automatically when the binary is built with
// -race, so this can't drift out of sync with the actual build.
const testTimeoutMultiplier = 4
