package logrus

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	current := len(handlers)

	var results []string

	h1 := func() { results = append(results, "first") }
	h2 := func() { results = append(results, "second") }

	RegisterExitHandler(h1)
	RegisterExitHandler(h2)

	if len(handlers) != current+2 {
		t.Fatalf("expected %d handlers, got %d", current+2, len(handlers))
	}

	runHandlers()

	if len(results) != 2 {
		t.Fatalf("expected 2 handlers to be run, ran %d", len(results))
	}

	if results[0] != "first" {
		t.Fatal("expected handler h1 to be run first, but it wasn't")
	}

	if results[1] != "second" {
		t.Fatal("expected handler h2 to be run second, but it wasn't")
	}
}

func TestDefer(t *testing.T) {
	current := len(handlers)

	var results []string

	h1 := func() { results = append(results, "first") }
	h2 := func() { results = append(results, "second") }

	DeferExitHandler(h1)
	DeferExitHandler(h2)

	if len(handlers) != current+2 {
		t.Fatalf("expected %d handlers, got %d", current+2, len(handlers))
	}

	runHandlers()

	if len(results) != 2 {
		t.Fatalf("expected 2 handlers to be run, ran %d", len(results))
	}

	if results[0] != "second" {
		t.Fatal("expected handler h2 to be run first, but it wasn't")
	}

	if results[1] != "first" {
		t.Fatal("expected handler h1 to be run second, but it wasn't")
	}
}

func TestHandler(t *testing.T) {
	testprog := testprogleader
	testprog = append(testprog, getPackage()...)
	testprog = append(testprog, testprogtrailer...)
	tempDir, err := ioutil.TempDir("", "test_handler")
	if err != nil {
		log.Fatalf("can't create temp dir. %q", err)
	}
	defer os.RemoveAll(tempDir)

	gofile := filepath.Join(tempDir, "gofile.go")
	if err := ioutil.WriteFile(gofile, testprog, 0666); err != nil {
		t.Fatalf("can't create go file. %q", err)
	}

	outfile := filepath.Join(tempDir, "outfile.out")
	arg := time.Now().UTC().String()
	err = exec.Command("go", "run", gofile, outfile, arg).Run()
	if err == nil {
		t.Fatalf("completed normally, should have failed")
	}

	data, err := ioutil.ReadFile(outfile)
	if err != nil {
		t.Fatalf("can't read output file %s. %q", outfile, err)
	}

	if string(data) != arg {
		t.Fatalf("bad data. Expected %q, got %q", data, arg)
	}
}

// getPackage returns the name of the current package, which makes running this
// test in a fork simpler
func getPackage() []byte {
	pc, _, _, _ := runtime.Caller(0)
	fullFuncName := runtime.FuncForPC(pc).Name()
	idx := strings.LastIndex(fullFuncName, ".")
	return []byte(fullFuncName[:idx]) // trim off function details
}

var testprogleader = []byte(`
// Test program for atexit, gets output file and data as arguments and writes
// data to output file in atexit handler.
package main

import (
	"`)
var testprogtrailer = []byte(
	`"
	"flag"
	"fmt"
	"io/ioutil"
)

var outfile = ""
var data = ""

func handler() {
	ioutil.WriteFile(outfile, []byte(data), 0666)
}

func badHandler() {
	n := 0
	fmt.Println(1/n)
}

func main() {
	flag.Parse()
	outfile = flag.Arg(0)
	data = flag.Arg(1)

	logrus.RegisterExitHandler(handler)
	logrus.RegisterExitHandler(badHandler)
	logrus.Fatal("Bye bye")
}
`)
