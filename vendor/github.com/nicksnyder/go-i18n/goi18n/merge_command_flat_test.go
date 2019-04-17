package main

import "testing"

func TestMergeExecuteFlat(t *testing.T) {
	files := []string{
		"testdata/input/flat/en-us.one.yaml",
		"testdata/input/flat/en-us.two.json",
		"testdata/input/flat/fr-fr.json",
		"testdata/input/flat/ar-ar.one.toml",
		"testdata/input/flat/ar-ar.two.json",
	}
	testFlatMergeExecute(t, files)
}

func testFlatMergeExecute(t *testing.T, files []string) {
	resetDir(t, "testdata/output/flat")

	mc := &mergeCommand{
		translationFiles: files,
		sourceLanguage:   "en-us",
		outdir:           "testdata/output/flat",
		format:           "json",
		flat:             true,
	}
	if err := mc.execute(); err != nil {
		t.Fatal(err)
	}

	expectEqualFiles(t, "testdata/output/flat/en-us.all.json", "testdata/expected/flat/en-us.all.json")
	expectEqualFiles(t, "testdata/output/flat/ar-ar.all.json", "testdata/expected/flat/ar-ar.all.json")
	expectEqualFiles(t, "testdata/output/flat/fr-fr.all.json", "testdata/expected/flat/fr-fr.all.json")
	expectEqualFiles(t, "testdata/output/flat/en-us.untranslated.json", "testdata/expected/flat/en-us.untranslated.json")
	expectEqualFiles(t, "testdata/output/flat/ar-ar.untranslated.json", "testdata/expected/flat/ar-ar.untranslated.json")
	expectEqualFiles(t, "testdata/output/flat/fr-fr.untranslated.json", "testdata/expected/flat/fr-fr.untranslated.json")
}
