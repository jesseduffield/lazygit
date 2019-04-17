// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const vfUsage = `xb version-file [options] <id>:<path>...

The command creates go file with a version constant. The version string
contains the contents of the VERSION environment variable or the output
of git describe.

   -h  prints this message and exits
   -p  package name (default main)
   -o  file name of output

`

func versionFileUsage(w io.Writer) {
	fmt.Fprint(w, vfUsage)
}

func versionFile() {
	cmdName := filepath.Base(os.Args[0])
	log.SetPrefix(fmt.Sprintf("%s: ", cmdName))
	log.SetFlags(0)

	flag.CommandLine = flag.NewFlagSet(cmdName, flag.ExitOnError)
	flag.Usage = func() { versionFileUsage(os.Stderr); os.Exit(1) }

	help := flag.Bool("h", false, "")
	pkg := flag.String("p", "main", "")
	out := flag.String("o", "", "")

	flag.Parse()

	if *help {
		versionFileUsage(os.Stdout)
		os.Exit(0)
	}

	if *pkg == "" {
		log.Fatal("option -p must not be empty")
	}

	var err error
	w := os.Stdout
	if *out != "" {
		if w, err = os.Create(*out); err != nil {
			log.Fatal(err)
		}
	}

	// get the version string
	version := os.Getenv("VERSION")
	if version == "" {
		b, err := exec.Command("git", "describe").Output()
		if err != nil {
			log.Fatalf("error %s while executing git describe", err)
		}
		version = string(b)
	}
	version = strings.TrimSpace(version)

	versionTmpl := `package main

const version = "{{.}}"
`
	tmpl := template.Must(template.New("version").Parse(versionTmpl))
	if err = tmpl.Execute(w, version); err != nil {
		log.Fatal(err)
	}
}
