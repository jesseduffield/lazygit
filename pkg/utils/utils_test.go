package utils

import "testing"

var testCases = []struct {
	Input	 []byte
	Expected []byte
}{
	{
		// \r\n
		Input: []byte{97, 115, 100, 102, 13, 10},
		Expected: []byte{97, 115, 100, 102, 10},
	},
	{
		// \r
		Input: []byte{97, 115, 100, 102, 13},
		Expected: []byte{97, 115, 100, 102, 10},
	},
	{
		// \n
		Input: []byte{97, 115, 100, 102, 10},
		Expected: []byte{97, 115, 100, 102, 10},
	},

}

func TestNormalizeLinefeeds(t *testing.T) {

	for _, tc := range testCases {
		if NormalizeLinefeeds(string(tc.Input)) != string(tc.Expected) {
			t.Error("Error")
		}
	}
}
