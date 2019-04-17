package getter

import (
	"path/filepath"
	"testing"
)

func TestXzDecompressor(t *testing.T) {
	cases := []TestDecompressCase{
		{
			"single.xz",
			false,
			false,
			nil,
			"d3b07384d113edec49eaa6238ad5ff00",
			nil,
		},

		{
			"single.xz",
			true,
			true,
			nil,
			"",
			nil,
		},
	}

	for i, tc := range cases {
		cases[i].Input = filepath.Join("./test-fixtures", "decompress-xz", tc.Input)
	}

	TestDecompressor(t, new(XzDecompressor), cases)
}
