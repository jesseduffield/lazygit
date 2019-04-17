package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/text/language"
)

type testCase struct {
	name           string
	inFiles        map[string][]byte
	sourceLanguage language.Tag
	outFiles       map[string][]byte
	deleteFiles    []string
}

func expectFile(s string) []byte {
	// Trimming leading newlines gives nicer formatting for file literals in test cases.
	return bytes.TrimLeft([]byte(s), "\n")
}

func TestMerge(t *testing.T) {
	testCases := []*testCase{
		{
			name:           "single identity",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"one.en-US.toml": []byte("1HelloMessage = \"Hello\"\n"),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": []byte("1HelloMessage = \"Hello\"\n"),
			},
		},
		{
			name:           "plural identity",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"active.en-US.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread email"
other = "{{.Count}} unread emails"
`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread email"
other = "{{.Count}} unread emails"
`),
			},
		},
		{
			name:           "migrate source lang from v1 format",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"one.en-US.json": []byte(`[
	{
		"id": "simple",
		"translation": "simple translation"
	},
	{
		"id": "everything",
		"translation": {
			"zero": "zero translation",
			"one": "one translation",
			"two": "two translation",
			"few": "few translation",
			"many": "many translation",
			"other": "other translation"
		}
	}
]`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
simple = "simple translation"

[everything]
few = "few translation"
many = "many translation"
one = "one translation"
other = "other translation"
two = "two translation"
zero = "zero translation"
`),
			},
		},
		{
			name:           "migrate source lang from v1 flat format",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"one.en-US.json": []byte(`{
	"simple": {
		"other": "simple translation"
	},
	"everything": {
		"zero": "zero translation",
		"one": "one translation",
		"two": "two translation",
		"few": "few translation",
		"many": "many translation",
		"other": "other translation"
	}
}`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
simple = "simple translation"

[everything]
few = "few translation"
many = "many translation"
one = "one translation"
other = "other translation"
two = "two translation"
zero = "zero translation"
`),
			},
		},
		{
			name:           "merge source files",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"one.en-US.toml": []byte("1HelloMessage = \"Hello\"\n"),
				"two.en-US.toml": []byte("2GoodbyeMessage = \"Goodbye\"\n"),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": []byte("1HelloMessage = \"Hello\"\n2GoodbyeMessage = \"Goodbye\"\n"),
			},
		},
		{
			name:           "missing hash",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
1HelloMessage = "Hello"
`),
				"es-ES.toml": []byte(`
[1HelloMessage]
other = "Hola"
`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
1HelloMessage = "Hello"
`),
				"active.es-ES.toml": expectFile(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"
`),
			},
		},
		{
			name:           "add single translation",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
1HelloMessage = "Hello"
2GoodbyeMessage = "Goodbye"
`),
				"es-ES.toml": []byte(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"
`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
1HelloMessage = "Hello"
2GoodbyeMessage = "Goodbye"
`),
				"active.es-ES.toml": expectFile(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"
`),
				"translate.es-ES.toml": expectFile(`
[2GoodbyeMessage]
hash = "sha1-b5b29c53e3c71cb9c6581ab053d7758fab8ca24d"
other = "Goodbye"
`),
			},
		},
		{
			name:           "remove single translation",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
1HelloMessage = "Hello"
`),
				"es-ES.toml": []byte(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"

[2GoodbyeMessage]
hash = "sha1-b5b29c53e3c71cb9c6581ab053d7758fab8ca24d"
other = "Goodbye"
`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
1HelloMessage = "Hello"
`),
				"active.es-ES.toml": expectFile(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"
`),
			},
		},
		{
			name:           "edit single translation",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
1HelloMessage = "Hi"
`),
				"es-ES.toml": []byte(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"
`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
1HelloMessage = "Hi"
`),
				"translate.es-ES.toml": expectFile(`
[1HelloMessage]
hash = "sha1-94dd9e08c129c785f7f256e82fbe0a30e6d1ae40"
other = "Hi"
`),
			},
		},
		{
			name:           "add plural translation",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread email"
other = "{{.Count}} unread emails"
`),
				"es-ES.toml": nil,
				"ar-AR.toml": nil,
				"zh-CN.toml": nil,
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread email"
other = "{{.Count}} unread emails"
`),
				"translate.es-ES.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
`),
				"translate.ar-AR.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
few = "{{.Count}} unread emails"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
many = "{{.Count}} unread emails"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
two = "{{.Count}} unread emails"
zero = "{{.Count}} unread emails"
`),
				"translate.zh-CN.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
other = "{{.Count}} unread emails"
`),
			},
		},
		{
			name:           "remove plural translation",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
1HelloMessage = "Hello"
`),
				"es-ES.toml": []byte(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"

[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
`),
				"ar-AR.toml": []byte(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hello"

[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
few = "{{.Count}} unread emails"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
many = "{{.Count}} unread emails"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
two = "{{.Count}} unread emails"
zero = "{{.Count}} unread emails"
`),
				"zh-CN.toml": []byte(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hello"

[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
other = "{{.Count}} unread emails"
`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
1HelloMessage = "Hello"
`),
				"active.es-ES.toml": expectFile(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hola"
`),
				"active.ar-AR.toml": expectFile(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hello"
`),
				"active.zh-CN.toml": expectFile(`
[1HelloMessage]
hash = "sha1-f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
other = "Hello"
`),
			},
		},
		{
			name:           "edit plural translation",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread emails!"
other = "{{.Count}} unread emails!"
`),
				"es-ES.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
`),
				"ar-AR.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
few = "{{.Count}} unread emails"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
many = "{{.Count}} unread emails"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
two = "{{.Count}} unread emails"
zero = "{{.Count}} unread emails"
`),
				"zh-CN.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
other = "{{.Count}} unread emails"
`),
			},
			deleteFiles: []string{
				"active.es-ES.toml",
				"active.ar-AR.toml",
				"active.zh-CN.toml",
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread emails!"
other = "{{.Count}} unread emails!"
`),
				"translate.es-ES.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-92a24983c5bbc0c42462cdc252dca68ebdb46501"
one = "{{.Count}} unread emails!"
other = "{{.Count}} unread emails!"
`),
				"translate.ar-AR.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
few = "{{.Count}} unread emails!"
hash = "sha1-92a24983c5bbc0c42462cdc252dca68ebdb46501"
many = "{{.Count}} unread emails!"
one = "{{.Count}} unread emails!"
other = "{{.Count}} unread emails!"
two = "{{.Count}} unread emails!"
zero = "{{.Count}} unread emails!"
`),
				"translate.zh-CN.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-92a24983c5bbc0c42462cdc252dca68ebdb46501"
other = "{{.Count}} unread emails!"
`),
			},
		},
		{
			name:           "merge plural translation",
			sourceLanguage: language.AmericanEnglish,
			inFiles: map[string][]byte{
				"en-US.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
`),
				"zero.ar-AR.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
zero = "{{.Count}} unread emails"
`),
				"one.ar-AR.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
one = "{{.Count}} unread emails"
`),
				"two.ar-AR.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
two = "{{.Count}} unread emails"
`),
				"few.ar-AR.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
few = "{{.Count}} unread emails"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
`),
				"many.ar-AR.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
many = "{{.Count}} unread emails"
`),
				"other.ar-AR.toml": []byte(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
other = "{{.Count}} unread emails"
`),
			},
			outFiles: map[string][]byte{
				"active.en-US.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
`),
				"active.ar-AR.toml": expectFile(`
[UnreadEmails]
description = "Message that tells the user how many unread emails they have"
few = "{{.Count}} unread emails"
hash = "sha1-5afbc91dfedb9755627655c365eb47a89e541099"
many = "{{.Count}} unread emails"
one = "{{.Count}} unread emails"
other = "{{.Count}} unread emails"
two = "{{.Count}} unread emails"
zero = "{{.Count}} unread emails"
`),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			indir := mustTempDir("TestMergeCommandIn")
			defer os.RemoveAll(indir)
			outdir := mustTempDir("TestMergeCommandOut")
			defer os.RemoveAll(outdir)

			infiles := make([]string, 0, len(testCase.inFiles))
			for name, content := range testCase.inFiles {
				path := filepath.Join(indir, name)
				infiles = append(infiles, path)
				if err := ioutil.WriteFile(path, content, 0666); err != nil {
					t.Fatal(err)
				}
			}

			for _, name := range testCase.deleteFiles {
				path := filepath.Join(outdir, name)
				if err := ioutil.WriteFile(path, []byte(`this file should get deleted`), 0666); err != nil {
					t.Fatal(err)
				}
			}

			args := append([]string{"merge", "-sourceLanguage", testCase.sourceLanguage.String(), "-outdir", outdir}, infiles...)
			if code := testableMain(args); code != 0 {
				t.Fatalf("expected exit code 0; got %d\n", code)
			}

			files, err := ioutil.ReadDir(outdir)
			if err != nil {
				t.Fatal(err)
			}

			// Verify that all actual files have expected contents.
			actualFiles := make(map[string]struct{}, len(files))
			for _, f := range files {
				actualFiles[f.Name()] = struct{}{}
				if f.IsDir() {
					t.Errorf("found unexpected dir %s", f.Name())
					continue
				}
				path := filepath.Join(outdir, f.Name())
				actual, err := ioutil.ReadFile(path)
				if err != nil {
					t.Error(err)
					continue
				}
				expected, ok := testCase.outFiles[f.Name()]
				if !ok {
					t.Errorf("found unexpected file %s with contents:\n%s\n", f.Name(), actual)
					continue
				}
				if !bytes.Equal(actual, expected) {
					t.Errorf("unexpected contents %s\ngot\n%s\nexpected\n%s", f.Name(), actual, expected)
					continue
				}
			}

			// Verify that all expected files are accounted for.
			for name := range testCase.outFiles {
				if _, ok := actualFiles[name]; !ok {
					t.Errorf("did not find expected file %s", name)
				}
			}
		})
	}
}
