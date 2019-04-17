package runewidth

import (
	"testing"
	"unicode/utf8"
)

var benchSink int

func benchTable(b *testing.B, tbl table) int {
	n := 0
	for i := 0; i < b.N; i++ {
		for r := rune(0); r <= utf8.MaxRune; r++ {
			if inTable(r, tbl) {
				n++
			}
		}
	}
	return n
}

func BenchmarkTablePrivate(b *testing.B) {
	benchSink = benchTable(b, private)
}
func BenchmarkTableNonprint(b *testing.B) {
	benchSink = benchTable(b, nonprint)
}
func BenchmarkTableCombining(b *testing.B) {
	benchSink = benchTable(b, combining)
}
func BenchmarkTableDoublewidth(b *testing.B) {
	benchSink = benchTable(b, doublewidth)
}
func BenchmarkTableAmbiguous(b *testing.B) {
	benchSink = benchTable(b, ambiguous)
}
func BenchmarkTableEmoji(b *testing.B) {
	benchSink = benchTable(b, emoji)
}
func BenchmarkTableNotassigned(b *testing.B) {
	benchSink = benchTable(b, notassigned)
}
func BenchmarkTableNeutral(b *testing.B) {
	benchSink = benchTable(b, neutral)
}
