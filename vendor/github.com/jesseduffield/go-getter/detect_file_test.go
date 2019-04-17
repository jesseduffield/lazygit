package getter

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type fileTest struct {
	in, pwd, out string
	err          bool
}

var fileTests = []fileTest{
	{"./foo", "/pwd", "file:///pwd/foo", false},
	{"./foo?foo=bar", "/pwd", "file:///pwd/foo?foo=bar", false},
	{"foo", "/pwd", "file:///pwd/foo", false},
}

var unixFileTests = []fileTest{
	{"./foo", "test-fixtures/detect-file-symlink-pwd/syml/pwd",
		"test-fixtures/detect-file-symlink-pwd/real/foo", false},

	{"/foo", "/pwd", "file:///foo", false},
	{"/foo?bar=baz", "/pwd", "file:///foo?bar=baz", false},
}

var winFileTests = []fileTest{
	{"/foo", "/pwd", "file:///pwd/foo", false},
	{`C:\`, `/pwd`, `file://C:/`, false},
	{`C:\?bar=baz`, `/pwd`, `file://C:/?bar=baz`, false},
}

func TestFileDetector(t *testing.T) {
	if runtime.GOOS == "windows" {
		fileTests = append(fileTests, winFileTests...)
	} else {
		fileTests = append(fileTests, unixFileTests...)
	}

	// Get the pwd
	pwdRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	pwdRoot, err = filepath.Abs(pwdRoot)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	f := new(FileDetector)
	for i, tc := range fileTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			pwd := tc.pwd
			if !filepath.IsAbs(pwd) {
				pwd = filepath.Join(pwdRoot, pwd)
			}

			out, ok, err := f.Detect(tc.in, pwd)
			if err != nil {
				t.Fatalf("err: %s", err)
			}
			if !ok {
				t.Fatal("not ok")
			}

			expected := tc.out
			if !strings.HasPrefix(expected, "file://") {
				expected = "file://" + filepath.Join(pwdRoot, expected)
			}

			if out != expected {
				t.Fatalf("input: %q\npwd: %q\nexpected: %q\nbad output: %#v",
					tc.in, pwd, expected, out)
			}
		})
	}
}

var noPwdFileTests = []fileTest{
	{in: "./foo", pwd: "", out: "", err: true},
	{in: "foo", pwd: "", out: "", err: true},
}

var noPwdUnixFileTests = []fileTest{
	{in: "/foo", pwd: "", out: "file:///foo", err: false},
}

var noPwdWinFileTests = []fileTest{
	{in: "/foo", pwd: "", out: "", err: true},
	{in: `C:\`, pwd: ``, out: `file://C:/`, err: false},
}

func TestFileDetector_noPwd(t *testing.T) {
	if runtime.GOOS == "windows" {
		noPwdFileTests = append(noPwdFileTests, noPwdWinFileTests...)
	} else {
		noPwdFileTests = append(noPwdFileTests, noPwdUnixFileTests...)
	}

	f := new(FileDetector)
	for i, tc := range noPwdFileTests {
		out, ok, err := f.Detect(tc.in, tc.pwd)
		if err != nil != tc.err {
			t.Fatalf("%d: err: %s", i, err)
		}
		if !ok {
			t.Fatal("not ok")
		}

		if out != tc.out {
			t.Fatalf("%d: bad: %#v", i, out)
		}
	}
}
