package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestMergeExecuteJSON(t *testing.T) {
	files := []string{
		"testdata/input/en-us.one.json",
		"testdata/input/en-us.two.json",
		"testdata/input/fr-fr.json",
		"testdata/input/ar-ar.one.json",
		"testdata/input/ar-ar.two.json",
	}
	testMergeExecute(t, files)
}

func TestMergeExecuteYAML(t *testing.T) {
	files := []string{
		"testdata/input/yaml/en-us.one.yaml",
		"testdata/input/yaml/en-us.two.json",
		"testdata/input/yaml/fr-fr.json",
		"testdata/input/yaml/ar-ar.one.json",
		"testdata/input/yaml/ar-ar.two.json",
	}
	testMergeExecute(t, files)
}

func testMergeExecute(t *testing.T, files []string) {
	resetDir(t, "testdata/output")

	mc := &mergeCommand{
		translationFiles: files,
		sourceLanguage:   "en-us",
		outdir:           "testdata/output",
		format:           "json",
		flat:             false,
	}
	if err := mc.execute(); err != nil {
		t.Fatal(err)
	}

	expectEqualFiles(t, "testdata/output/en-us.all.json", "testdata/expected/en-us.all.json")
	expectEqualFiles(t, "testdata/output/ar-ar.all.json", "testdata/expected/ar-ar.all.json")
	expectEqualFiles(t, "testdata/output/fr-fr.all.json", "testdata/expected/fr-fr.all.json")
	expectEqualFiles(t, "testdata/output/en-us.untranslated.json", "testdata/expected/en-us.untranslated.json")
	expectEqualFiles(t, "testdata/output/ar-ar.untranslated.json", "testdata/expected/ar-ar.untranslated.json")
	expectEqualFiles(t, "testdata/output/fr-fr.untranslated.json", "testdata/expected/fr-fr.untranslated.json")
}

func resetDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatal(err)
	}
}

func expectEqualFiles(t *testing.T, expectedName, actualName string) {
	actual, err := ioutil.ReadFile(actualName)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(expectedName)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Errorf("contents of files did not match: %s, %s", expectedName, actualName)
	}
}
