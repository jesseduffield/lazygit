// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command gxz supports the compression and decompression of LZMA files.
//
// Use gxz -h to get information about supported flags.
package main

//go:generate xb cat -o licenses.go xzLicense:github.com/ulikunitz/xz/LICENSE goLicense:~/go/LICENSE
//go:generate xb version-file -o version.go

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"text/template"

	"github.com/ulikunitz/xz/internal/gflag"
	"github.com/ulikunitz/xz/internal/term"
	"github.com/ulikunitz/xz/internal/xlog"
)

const (
	usageStr = `Usage: gxz [OPTION]... [FILE]...
Compress or uncompress FILEs in the .lzma format (by default, compress FILES
in place).

  -c, --stdout      write to standard output and don't delete input files
  -d, --decompress  force decompression
  -f, --force       force overwrite of output file and compress links
  -F, --format <format>
                    Specify the file format to compress or decompress.
    auto            Default format for compression is xz. For decompression
	            the file content is used to identify the format.
    xz              The xz file format.
    lzma, alone     Compress to the .lzma file format.
  -h, --help        give this help
  -k, --keep        keep (don't delete) input files
  -L, --license     display software license
  -q, --quiet       suppress all warnings
  -v, --verbose     verbose mode
  -V, --version     display version string
  -z, --compress    force compression
  -0 ... -9         compression preset; default is 6
  --cpuprofile <file>
                    create a cpuprofile that can be used with go tool pprof

With no file, or when FILE is -, read standard input.

Report bugs using <https://github.com/ulikunitz/xz/issues>.
`
)

func usage(w io.Writer) {
	fmt.Fprint(w, usageStr)
}

func licenses(w io.Writer) {
	out := `
github.com/ulikunitz/xz -- xz for Go
====================================

{{.xz}}

Go Programming Language
=======================

The gxz program contains the packages gflag and xlog that are
extensions of packages from the Go standard library. The packages may
contain code from those packages.

{{.go}}
`
	out = strings.TrimLeft(out, " \n")
	tmpl, err := template.New("licenses").Parse(out)
	if err != nil {
		xlog.Panicf("error %s parsing licenses template", err)
	}
	lmap := map[string]string{
		"xz": strings.TrimSpace(xzLicense),
		"go": strings.TrimSpace(goLicense),
	}
	if err = tmpl.Execute(w, lmap); err != nil {
		xlog.Fatalf("error %s writing licenses template", err)
	}
}

type options struct {
	help       bool
	stdout     bool
	decompress bool
	force      bool
	format     string
	keep       bool
	license    bool
	version    bool
	quiet      int
	verbose    int
	preset     int
	cpuprofile string
}

func (o *options) Init() {
	if o.preset != 0 {
		xlog.Panicf("options are already initialized")
	}
	gflag.BoolVarP(&o.help, "help", "h", false, "")
	gflag.BoolVarP(&o.stdout, "stdout", "c", false, "")
	gflag.BoolVarP(&o.decompress, "decompress", "d", false, "")
	gflag.BoolVarP(&o.force, "force", "f", false, "")
	gflag.StringVarP(&o.format, "format", "F", "auto", "")
	gflag.BoolVarP(&o.keep, "keep", "k", false, "")
	gflag.BoolVarP(&o.license, "license", "L", false, "")
	gflag.BoolVarP(&o.version, "version", "V", false, "")
	gflag.CounterVarP(&o.quiet, "quiet", "q", 0, "")
	gflag.CounterVarP(&o.verbose, "verbose", "v", 0, "")
	gflag.PresetVar(&o.preset, 0, 9, 6, "")
	gflag.StringVarP(&o.cpuprofile, "cpuprofile", "", "", "")
}

// normalizeFormat normalizes the format field of options. If the
// function completes without error the format field will be "xz",
// "lzma" or "auto". The latter only if the option decompress is true.
func normalizeFormat(o *options) error {
	switch o.format {
	case "xz", "lzma":
	case "auto":
		if !o.decompress {
			o.format = "xz"
		}
	case "alone":
		o.format = "lzma"
	default:
		return fmt.Errorf("format %q unsupported", o.format)
	}
	return nil
}

func main() {
	// setup logger
	cmdName := filepath.Base(os.Args[0])
	xlog.SetPrefix(fmt.Sprintf("%s: ", cmdName))
	xlog.SetFlags(0)

	// initialize flags
	gflag.CommandLine = gflag.NewFlagSet(cmdName, gflag.ExitOnError)
	gflag.Usage = func() { usage(os.Stderr); os.Exit(1) }
	opts := options{}
	opts.Init()

	switch cmdName {
	case "lzma", "glzma":
		opts.format = "lzma"
	case "lzcat", "glzcat":
		opts.format = "lzma"
		fallthrough
	case "xzcat", "gxzcat":
		opts.stdout = true
		opts.decompress = true
	case "unlzma", "unglzma":
		opts.format = "lzma"
		fallthrough
	case "unxz", "ungxz":
		opts.decompress = true
	}
	gflag.Parse()

	if opts.help {
		usage(os.Stdout)
		os.Exit(0)
	}
	if opts.license {
		licenses(os.Stdout)
		os.Exit(0)
	}
	if opts.version {
		xlog.Printf("version %s\n", version)
		os.Exit(0)
	}

	flags := xlog.Flags()
	switch {
	case opts.verbose <= 0:
		flags |= xlog.Lnoprint | xlog.Lnodebug
	case opts.verbose == 1:
		flags |= xlog.Lnodebug
	}
	switch {
	case opts.quiet >= 2:
		flags |= xlog.Lnoprint | xlog.Lnowarn | xlog.Lnodebug
		flags |= xlog.Lnopanic | xlog.Lnofatal
	case opts.quiet == 1:
		flags |= xlog.Lnoprint | xlog.Lnowarn | xlog.Lnodebug
	}
	xlog.SetFlags(flags)

	if opts.cpuprofile != "" {
		f, err := os.Create(opts.cpuprofile)
		if err != nil {
			xlog.Fatal(err)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			xlog.Fatal(err)
		}
	}

	if err := normalizeFormat(&opts); err != nil {
		pprof.StopCPUProfile()
		xlog.Fatal(err)
	}

	var args []string
	if gflag.NArg() == 0 {
		opts.stdout = true
		args = []string{"-"}
	} else {
		args = gflag.Args()
	}

	if opts.stdout && !opts.decompress && !opts.force &&
		term.IsTerminal(os.Stdout.Fd()) {
		pprof.StopCPUProfile()
		xlog.Fatal(`Compressed data will not be written to a terminal
Use -f to force compression. For help type gxz -h.`)
	}

	exit := 0
	for _, arg := range args {
		if err := processFile(arg, &opts); err != nil {
			exit = 1
		}
	}

	pprof.StopCPUProfile()
	os.Exit(exit)
}
