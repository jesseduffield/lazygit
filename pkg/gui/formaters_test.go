package gui

import "testing"

func generateWord(amRunes int) []rune {
	charSlice := "abcåäö🐣⛽🌳"
	runeSlice := make([]rune, 0)
	for _, c := range charSlice {
		runeSlice = append(runeSlice, c)
	}
	ret := make([]rune, amRunes)
	for i := 0; i < amRunes; i++ {
		runeSliceLen := len(runeSlice)
		runeIdx := i % runeSliceLen
		ret[i] = []rune(runeSlice)[runeIdx]
	}
	return ret
}

func Test_formatGit72(t *testing.T) {
	t.Run("it should insert a newline if word is detected after col 72", func(t *testing.T) {
		w0 := generateWord(72)
		w1 := generateWord(1)
		given := append(w0, ' ')
		given = append(given, w1...)
		givenLen := len(given)
		//         w0  ' '  w1
		wantLen := 72 + 1 + 1
		if givenLen != wantLen {
			t.Fatalf("test setup failed, want: %v, got: %v", wantLen, givenLen)
		}

		gotRunes, err := formatGit72(given)
		if err != nil {
			t.Fatalf("failed to format: %v", err)
		}
		gotLen := len(gotRunes)
		// expect space rune to have been changed to a newline, keeping total length the same
		if wantLen != gotLen {
			t.Fatalf("expected: %v, got: %v", wantLen, gotLen)
		}

		// 73'd rune should be newline, as 72 runes are allowed per line
		want := '\n'
		got := gotRunes[72]
		if got != want {
			t.Fatalf("expected: %v, got: %v, given: %v, gotRunes: %v", want, got, given[68:wantLen], gotRunes[68:wantLen])
		}
	})

	// The vim git-plugin doesn't break mid-word, but rather breaks on the previous space
	// to the word which is overflowing. This generates a much more aesteticly pleasing format.
	t.Run("it should break on space behind overflowing word", func(t *testing.T) {
		w0 := generateWord(65)
		w1 := generateWord(10)
		given := append(w0, ' ')
		given = append(given, w1...)
		gotRunes, err := formatGit72(given)
		if err != nil {
			t.Fatalf("failed to format: %v", err)
		}
		if gotRunes[72] == '\n' {
			t.Fatal("expected break to not occur mid rune")
		}

		if gotRunes[65] != '\n' {
			t.Fatal("expected newline after previous space")
		}
	})
}
