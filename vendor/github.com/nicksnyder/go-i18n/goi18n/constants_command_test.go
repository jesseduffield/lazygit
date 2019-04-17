package main

import "testing"

func TestConstantsExecute(t *testing.T) {
	resetDir(t, "testdata/output")

	cc := &constantsCommand{
		translationFiles: []string{"testdata/input/en-us.constants.json"},
		packageName:      "R",
		outdir:           "testdata/output",
	}

	if err := cc.execute(); err != nil {
		t.Fatal(err)
	}

	expectEqualFiles(t, "testdata/output/R.go", "testdata/expected/R.go")
}

func TestToCamelCase(t *testing.T) {
	expectEqual := func(test, expected string) {
		result := toCamelCase(test)
		if result != expected {
			t.Fatalf("failed toCamelCase the test %s was expected %s but the result was %s", test, expected, result)
		}
	}

	expectEqual("", "")
	expectEqual("a", "A")
	expectEqual("_", "")
	expectEqual("__code__", "Code")
	expectEqual("test", "Test")
	expectEqual("test_one", "TestOne")
	expectEqual("test.two", "TestTwo")
	expectEqual("test_alpha_beta", "TestAlphaBeta")
	expectEqual("word  word", "WordWord")
	expectEqual("test_id", "TestID")
	expectEqual("tcp_name", "TCPName")
	expectEqual("こんにちは", "こんにちは")
	expectEqual("test_a", "TestA")
}
