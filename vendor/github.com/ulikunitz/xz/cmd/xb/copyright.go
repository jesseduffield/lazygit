// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const crUsageString = `xb copyright [options] <path>....

The xb copyright command adds a copyright remark to all go files below path.

  -h  prints this message and exits
`

func crUsage(w io.Writer) {
	fmt.Fprint(w, crUsageString)
}

const copyrightText = `
Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
`

func goComment(text string) string {
	buf := new(bytes.Buffer)
	scanner := bufio.NewScanner(strings.NewReader(text))
	var err error
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) == 0 {
			continue
		}
		if _, err = fmt.Fprintln(buf, "//", s); err != nil {
			panic(err)
		}
	}
	if err = scanner.Err(); err != nil {
		panic(err)
	}
	if _, err = fmt.Fprintln(buf); err != nil {
		panic(err)
	}
	return buf.String()
}

var goCopyright = goComment(copyrightText)

func addCopyright(path string) (err error) {
	log.Printf("adding copyright to %s", path)
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr := src.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()
	newPath := path + ".new"
	dst, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer func() {
		cerr := dst.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()
	out := bufio.NewWriter(dst)
	fmt.Fprint(out, goCopyright)
	scanner := bufio.NewScanner(src)
	line := 0
	del := false
	for scanner.Scan() {
		line++
		txt := scanner.Text()
		if line == 1 && strings.Contains(txt, "Copyright") {
			del = true
			continue
		}
		if del {
			s := strings.TrimSpace(txt)
			if len(s) == 0 {
				del = false
			}
			continue
		}
		fmt.Fprintln(out, txt)
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	if err = out.Flush(); err != nil {
		return
	}
	err = os.Rename(newPath, path)
	return
}

func walkCopyrights(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if !strings.HasSuffix(info.Name(), ".go") {
		return nil
	}
	return addCopyright(path)
}

func copyright() {
	cmdName := os.Args[0]
	log.SetPrefix(fmt.Sprintf("%s: ", cmdName))
	log.SetFlags(0)

	flag.CommandLine = flag.NewFlagSet(cmdName, flag.ExitOnError)
	flag.Usage = func() { crUsage(os.Stderr); os.Exit(1) }

	help := flag.Bool("h", false, "")

	flag.Parse()

	if *help {
		crUsage(os.Stdout)
		os.Exit(0)
	}

	for _, path := range flag.Args() {
		fi, err := os.Stat(path)
		if err != nil {
			log.Print(err)
			continue
		}
		if !fi.IsDir() {
			log.Printf("%s is not a directory", path)
			continue
		}
		if err = filepath.Walk(path, walkCopyrights); err != nil {
			log.Fatalf("%s error %s", path, err)
		}
	}
}
