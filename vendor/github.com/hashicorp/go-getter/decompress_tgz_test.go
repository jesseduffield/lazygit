package getter

import (
	"path/filepath"
	"testing"
)

func TestTarGzipDecompressor(t *testing.T) {

	multiplePaths := []string{"dir/", "dir/test2", "test1"}
	orderingPaths := []string{"workers/", "workers/mq/", "workers/mq/__init__.py"}

	cases := []TestDecompressCase{
		{
			"empty.tar.gz",
			false,
			true,
			nil,
			"",
			nil,
		},

		{
			"single.tar.gz",
			false,
			false,
			nil,
			"d3b07384d113edec49eaa6238ad5ff00",
			nil,
		},

		{
			"single.tar.gz",
			true,
			false,
			[]string{"file"},
			"",
			nil,
		},

		{
			"multiple.tar.gz",
			true,
			false,
			[]string{"file1", "file2"},
			"",
			nil,
		},

		{
			"multiple.tar.gz",
			false,
			true,
			nil,
			"",
			nil,
		},

		{
			"multiple_dir.tar.gz",
			true,
			false,
			multiplePaths,
			"",
			nil,
		},

		// Tests when the file is listed before the parent folder
		{
			"ordering.tar.gz",
			true,
			false,
			orderingPaths,
			"",
			nil,
		},

		// Tests that a tar.gz can't contain references with "..".
		// GNU `tar` also disallows this.
		{
			"outside_parent.tar.gz",
			true,
			true,
			nil,
			"",
			nil,
		},
	}

	for i, tc := range cases {
		cases[i].Input = filepath.Join("./test-fixtures", "decompress-tgz", tc.Input)
	}

	TestDecompressor(t, new(TarGzipDecompressor), cases)
}
