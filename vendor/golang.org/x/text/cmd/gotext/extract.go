// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/text/internal"
	"golang.org/x/text/language"
	"golang.org/x/text/message/pipeline"
)

// TODO:
// - merge information into existing files
// - handle different file formats (PO, XLIFF)
// - handle features (gender, plural)
// - message rewriting

var (
	srcLang *string
	lang    *string
)

func init() {
	srcLang = cmdExtract.Flag.String("srclang", "en-US", "the source-code language")
	lang = cmdExtract.Flag.String("lang", "en-US", "comma-separated list of languages to process")
}

var cmdExtract = &Command{
	Run:       runExtract,
	UsageLine: "extract <package>*",
	Short:     "extracts strings to be translated from code",
}

func runExtract(cmd *Command, args []string) error {
	tag, err := language.Parse(*srcLang)
	if err != nil {
		return wrap(err, "")
	}
	config := &pipeline.Config{
		SourceLanguage: tag,
		Packages:       args,
	}
	out, err := pipeline.Extract(config)

	data, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		return wrap(err, "")
	}
	os.MkdirAll(*dir, 0755)
	// TODO: this file can probably go if we replace the extract + generate
	// cycle with a init once and update cycle.
	file := filepath.Join(*dir, extractFile)
	if err := ioutil.WriteFile(file, data, 0644); err != nil {
		return wrap(err, "could not create file")
	}

	langs := append(getLangs(), tag)
	langs = internal.UniqueTags(langs)
	for _, tag := range langs {
		// TODO: inject translations from existing files to avoid retranslation.
		out.Language = tag
		data, err := json.MarshalIndent(out, "", "    ")
		if err != nil {
			return wrap(err, "JSON marshal failed")
		}
		file := filepath.Join(*dir, tag.String(), outFile)
		if err := os.MkdirAll(filepath.Dir(file), 0750); err != nil {
			return wrap(err, "dir create failed")
		}
		if err := ioutil.WriteFile(file, data, 0740); err != nil {
			return wrap(err, "write failed")
		}
	}
	return nil
}
