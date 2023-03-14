//go:build ignore

// This file is invoked with `go generate ./...` and it generates the test_list.go file
// The test_list.go file is a list of all the integration tests.
// It's annoying to have to manually add an entry in that file for each test you
// create, so this generator is here to make the process easier.

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/samber/lo"
)

func main() {
	code := generateCode()

	formattedCode, err := format.Source(code)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("test_list.go", formattedCode, 0o644); err != nil {
		panic(err)
	}
}

func generateCode() []byte {
	// traverse parent directory to get all subling directories
	directories, err := ioutil.ReadDir("../tests")
	if err != nil {
		panic(err)
	}

	directories = lo.Filter(directories, func(file os.FileInfo, _ int) bool {
		// 'shared' is a special folder containing shared test code so we
		// ignore it here
		return file.IsDir() && file.Name() != "shared"
	})

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// THIS FILE IS AUTO-GENERATED. You can regenerate it by running `go generate ./...` at the root of the lazygit repo.\n\n")
	fmt.Fprintf(&buf, "package tests\n\n")
	fmt.Fprintf(&buf, "import (\n")
	fmt.Fprintf(&buf, "\t\"github.com/jesseduffield/lazygit/pkg/integration/components\"\n")
	for _, dir := range directories {
		fmt.Fprintf(&buf, "\t\"github.com/jesseduffield/lazygit/pkg/integration/tests/%s\"\n", dir.Name())
	}
	fmt.Fprintf(&buf, ")\n\n")
	fmt.Fprintf(&buf, "var tests = []*components.IntegrationTest{\n")
	for _, dir := range directories {
		appendDirTests(dir, &buf)
	}
	fmt.Fprintf(&buf, "}\n")

	return buf.Bytes()
}

func appendDirTests(dir fs.FileInfo, buf *bytes.Buffer) {
	files, err := ioutil.ReadDir(fmt.Sprintf("../tests/%s", dir.Name()))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		testName := snakeToPascal(
			strings.TrimSuffix(file.Name(), ".go"),
		)

		fileContents, err := ioutil.ReadFile(fmt.Sprintf("../tests/%s/%s", dir.Name(), file.Name()))
		if err != nil {
			panic(err)
		}

		fileContentsStr := string(fileContents)

		if !strings.Contains(fileContentsStr, "NewIntegrationTest(") {
			// the file does not define a test so it probably just contains shared test code
			continue
		}

		if !strings.Contains(fileContentsStr, fmt.Sprintf("var %s = NewIntegrationTest(NewIntegrationTestArgs{", testName)) {
			panic(fmt.Sprintf("expected test %s to be defined in file %s. Perhaps you misspelt it? The name of the test should be the name of the file but converted from snake_case to PascalCase", testName, file.Name()))
		}

		fmt.Fprintf(buf, "\t%s.%s,\n", dir.Name(), testName)
	}
}

// thanks ChatGPT
func snakeToPascal(s string) string {
	// Split the input string into words.
	words := strings.Split(s, "_")

	// Convert the first letter of each word to uppercase and concatenate them.
	var builder strings.Builder
	for _, w := range words {
		if len(w) > 0 {
			builder.WriteString(strings.ToUpper(w[:1]))
			builder.WriteString(w[1:])
		}
	}

	return builder.String()
}
