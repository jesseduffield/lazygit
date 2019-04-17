package getter

import (
	"path/filepath"
	"testing"
)

func TestBzip2Decompressor(t *testing.T) {
	cases := []TestDecompressCase{
		{
			"single.bz2",
			false,
			false,
			nil,
			"d3b07384d113edec49eaa6238ad5ff00",
			nil,
		},

		{
			"single.bz2",
			true,
			true,
			nil,
			"",
			nil,
		},
	}

	for i, tc := range cases {
		cases[i].Input = filepath.Join("./test-fixtures", "decompress-bz2", tc.Input)
	}

	TestDecompressor(t, new(Bzip2Decompressor), cases)
}
