//go:generate go run test_list_generator.go

package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/samber/lo"
)

func GetTests(lazygitRootDir string) []*components.IntegrationTest {
	// first we ensure that each test in this directory has actually been added to the above list.
	testCount := 0

	testNamesSet := set.NewFromSlice(lo.Map(
		tests,
		func(test *components.IntegrationTest, _ int) string {
			return test.Name()
		},
	))

	missingTestNames := []string{}

	if err := filepath.Walk(filepath.Join(lazygitRootDir, "pkg/integration/tests"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// ignoring non-test files
			if filepath.Base(path) == "tests.go" || filepath.Base(path) == "test_list.go" || filepath.Base(path) == "test_list_generator.go" {
				return nil
			}

			// the shared directory won't itself contain tests: only shared helper functions
			if filepath.Base(filepath.Dir(path)) == "shared" {
				return nil
			}

			// any file named shared.go will also be ignored, because those files are only used for shared helper functions
			if filepath.Base(path) == "shared.go" {
				return nil
			}

			nameFromPath := components.TestNameFromFilePath(path)
			if !testNamesSet.Includes(nameFromPath) {
				missingTestNames = append(missingTestNames, nameFromPath)
			}
			testCount++
		}
		return nil
	}); err != nil {
		panic(fmt.Sprintf("failed to walk tests: %v", err))
	}

	if len(missingTestNames) > 0 {
		panic(fmt.Sprintf("The following tests are missing from the list of tests: %s. You need to add them to `pkg/integration/tests/test_list.go`. Use `go generate ./...` to regenerate the tests list.", strings.Join(missingTestNames, ", ")))
	}

	if testCount > len(tests) {
		panic("you have not added all of the tests to the tests list in `pkg/integration/tests/test_list.go`. Use `go generate ./...` to regenerate the tests list.")
	} else if testCount < len(tests) {
		panic("There are more tests in `pkg/integration/tests/test_list.go` than there are test files in the tests directory. Ensure that you only have one test per file and you haven't included the same test twice in the tests list. Use `go generate ./...` to regenerate the tests list.")
	}

	return tests
}
