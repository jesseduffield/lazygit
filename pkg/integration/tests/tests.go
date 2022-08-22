package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/bisect"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/branch"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/commit"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/custom_commands"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/interactive_rebase"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

// Here is where we lists the actual tests that will run. When you create a new test,
// be sure to add it to this list.

var tests = []*components.IntegrationTest{
	commit.Commit,
	commit.NewBranch,
	branch.Suggestions,
	branch.Delete,
	interactive_rebase.One,
	custom_commands.Basic,
	custom_commands.MultiplePrompts,
	custom_commands.MenuFromCommand,
	bisect.Basic,
	bisect.FromOtherBranch,
}

func GetTests() []*components.IntegrationTest {
	// first we ensure that each test in this directory has actually been added to the above list.
	testCount := 0

	testNamesSet := set.NewFromSlice(slices.Map(
		tests,
		func(test *components.IntegrationTest) string {
			return test.Name()
		},
	))

	missingTestNames := []string{}

	if err := filepath.Walk(filepath.Join(utils.GetLazygitRootDirectory(), "pkg/integration/tests"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// ignoring this current file
			if filepath.Base(path) == "tests.go" {
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
		panic(fmt.Sprintf("The following tests are missing from the list of tests: %s. You need to add them to `pkg/integration/tests/tests.go`.", strings.Join(missingTestNames, ", ")))
	}

	if testCount > len(tests) {
		panic("you have not added all of the tests to the tests list in `pkg/integration/tests/tests.go`")
	} else if testCount < len(tests) {
		panic("There are more tests in `pkg/integration/tests/tests.go` than there are test files in the tests directory. Ensure that you only have one test per file and you haven't included the same test twice in the tests list.")
	}

	return tests
}
