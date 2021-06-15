package command

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func Detect(cwd string, args []string, generateMode bool) ([]Invocation, error) {
	if generateMode {
		return generateModeInvocations(cwd)
	}

	file := os.Getenv("GOFILE")
	var lineno int
	if goline, err := strconv.Atoi(os.Getenv("GOLINE")); err == nil {
		lineno = goline
	}

	i, err := NewInvocation(file, lineno, args)
	if err != nil {
		return nil, err
	}
	return []Invocation{i}, nil
}

type Invocation struct {
	Args []string
	Line int
	File string
}

func NewInvocation(file string, line int, args []string) (Invocation, error) {
	if len(args) < 1 {
		return Invocation{}, fmt.Errorf("%s:%v an invocation of counterfeiter must have arguments", file, line)
	}
	i := Invocation{
		File: file,
		Line: line,
		Args: args,
	}
	return i, nil
}

func generateModeInvocations(cwd string) ([]Invocation, error) {
	var result []Invocation
	// Find all the go files
	pkg, err := build.ImportDir(cwd, build.IgnoreVendor)
	if err != nil {
		return nil, err
	}

	gofiles := make([]string, 0, len(pkg.GoFiles)+len(pkg.CgoFiles)+len(pkg.TestGoFiles)+len(pkg.XTestGoFiles))
	gofiles = append(gofiles, pkg.GoFiles...)
	gofiles = append(gofiles, pkg.CgoFiles...)
	gofiles = append(gofiles, pkg.TestGoFiles...)
	gofiles = append(gofiles, pkg.XTestGoFiles...)
	sort.Strings(gofiles)

	for _, file := range gofiles {
		invocations, err := invocationsInFile(cwd, file)
		if err != nil {
			return nil, err
		}
		result = append(result, invocations...)
	}

	return result, nil
}

func invocationsInFile(dir string, file string) ([]Invocation, error) {
	str, err := ioutil.ReadFile(filepath.Join(dir, file))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(str), "\n")

	var result []Invocation
	line := 0
	for i := range lines {
		line++
		args, ok := matchForString(lines[i])
		if !ok {
			continue
		}
		inv, err := NewInvocation(file, line, args)
		if err != nil {
			return nil, err
		}

		result = append(result, inv)
	}

	return result, nil
}

const generateDirectivePrefix = "//counterfeiter:generate "

func matchForString(s string) ([]string, bool) {
	if !strings.HasPrefix(s, generateDirectivePrefix) {
		return nil, false
	}
	return stringToArgs(s[len(generateDirectivePrefix):]), true
}

func stringToArgs(s string) []string {
	a := strings.Fields(s)
	result := []string{
		"counterfeiter",
	}
	result = append(result, a...)
	return result
}
