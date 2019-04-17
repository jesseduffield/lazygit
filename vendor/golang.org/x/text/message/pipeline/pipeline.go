// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pipeline provides tools for creating translation pipelines.
//
// NOTE: UNDER DEVELOPMENT. API MAY CHANGE.
package pipeline

import (
	"fmt"
	"go/build"
	"go/parser"
	"log"

	"golang.org/x/tools/go/loader"
)

const (
	extractFile  = "extracted.gotext.json"
	outFile      = "out.gotext.json"
	gotextSuffix = ".gotext.json"
)

// NOTE: The command line tool already prefixes with "gotext:".
var (
	wrap = func(err error, msg string) error {
		return fmt.Errorf("%s: %v", msg, err)
	}
	wrapf = func(err error, msg string, args ...interface{}) error {
		return wrap(err, fmt.Sprintf(msg, args...))
	}
	errorf = fmt.Errorf
)

// TODO: don't log.
func logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func loadPackages(conf *loader.Config, args []string) (*loader.Program, error) {
	if len(args) == 0 {
		args = []string{"."}
	}

	conf.Build = &build.Default
	conf.ParserMode = parser.ParseComments

	// Use the initial packages from the command line.
	args, err := conf.FromArgs(args, false)
	if err != nil {
		return nil, wrap(err, "loading packages failed")
	}

	// Load, parse and type-check the whole program.
	return conf.Load()
}
