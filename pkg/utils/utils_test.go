package utils

import "testing"

var testCases = []struct {
	Input    []byte
	Expected []byte
}{
	{
		// \r\n
		Input:    []byte{97, 115, 100, 102, 13, 10},
		Expected: []byte{97, 115, 100, 102},
	},
	{
		// \r
		Input:    []byte{97, 115, 100, 102, 13},
		Expected: []byte{97, 115, 100, 102},
	},
	{
		// \n
		Input:    []byte{97, 115, 100, 102, 10},
		Expected: []byte{97, 115, 100, 102, 10},
	},
}

func TestNormalizeLinefeeds(t *testing.T) {
	for _, tc := range testCases {
		input := NormalizeLinefeeds(string(tc.Input))
		expected := string(tc.Expected)
		if input != expected {
			t.Error("Expected " + expected + ", got " + input)
		}
	}
}
