// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const catUsageString = `xb cat [options] <id>:<path>...

This xb command puts the contents of the files given as relative paths to
the GOPATH variable as string constants into a go file. 

   -h  prints this message and exits
   -p  package name (default main)
   -o  file name of output

`

func catUsage(w io.Writer) {
	fmt.Fprint(w, catUsageString)
}

type gopath struct {
	p []string
	i int
}

func newGopath() *gopath {
	p := strings.Split(os.Getenv("GOPATH"), ":")
	return &gopath{p: p}
}

type cpair struct {
	id   string
	path string
}

func (p cpair) Read() (s string, err error) {
	var r io.ReadCloser
	if p.path == "-" {
		r = os.Stdin
	} else {
		if r, err = os.Open(p.path); err != nil {
			return
		}
	}
	defer func() {
		err = r.Close()
	}()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	s = string(b)
	return
}

func verifyPath(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", path)
	}
	return nil
}

func (gp *gopath) find(arg string) (p cpair, err error) {
	t := strings.SplitN(arg, ":", 2)
	switch len(t) {
	case 0:
		err = fmt.Errorf("empty argument not supported")
		return
	case 1:
		gp.i++
		p = cpair{fmt.Sprintf("gocat%d", gp.i), t[0]}
	case 2:
		p = cpair{t[0], t[1]}
	}
	if p.path == "-" {
		return
	}
	// substitute first ~ by $HOME
	p.path = strings.Replace(p.path, "~", os.Getenv("HOME"), 1)
	paths := make([]string, 0, len(gp.p)+1)
	if filepath.IsAbs(p.path) {
		paths = append(paths, filepath.Clean(p.path))
	} else {
		for _, q := range gp.p {
			u := filepath.Join(q, "src", p.path)
			paths = append(paths, filepath.Clean(u))
		}
		u := filepath.Join(".", p.path)
		paths = append(paths, filepath.Clean(u))
	}
	for _, u := range paths {
		if err = verifyPath(u); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return
		}
		p.path = u
		return
	}
	err = fmt.Errorf("file %s not found", p.path)
	return
}

// Gofile is used with the template gofileTmpl.
type Gofile struct {
	Pkg  string
	Cmap map[string]string
}

var gofileTmpl = `package {{.Pkg}}

{{range $k, $v := .Cmap}}const {{$k}} = ` + "`{{$v}}`\n{{end}}"

func cat() {
	var err error
	cmdName := filepath.Base(os.Args[0])
	log.SetPrefix(fmt.Sprintf("%s: ", cmdName))
	log.SetFlags(0)

	flag.CommandLine = flag.NewFlagSet(cmdName, flag.ExitOnError)
	flag.Usage = func() { catUsage(os.Stderr); os.Exit(1) }

	help := flag.Bool("h", false, "")
	pkg := flag.String("p", "main", "")
	out := flag.String("o", "", "")

	flag.Parse()

	if *help {
		catUsage(os.Stdout)
		os.Exit(0)
	}

	if *pkg == "" {
		log.Fatal("option -p must not be empty")
	}

	w := os.Stdout
	if *out != "" {
		if w, err = os.Create(*out); err != nil {
			log.Fatal(err)
		}
	}

	gp := newGopath()

	gofile := Gofile{
		Pkg:  *pkg,
		Cmap: make(map[string]string, len(flag.Args())),
	}
	for _, arg := range flag.Args() {
		p, err := gp.find(arg)
		if err != nil {
			log.Print(err)
			continue
		}
		s, err := p.Read()
		if err != nil {
			log.Print(err)
			continue
		}
		gofile.Cmap[p.id] = s
	}

	tmpl, err := template.New("gofile").Parse(gofileTmpl)
	if err != nil {
		log.Panicf("goFileTmpl error %s", err)
	}
	if err = tmpl.Execute(w, gofile); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
