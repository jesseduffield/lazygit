package integration_tests

import (
	"github.com/jesseduffield/lazygit/pkg/integration/integration_tests/branch"
	"github.com/jesseduffield/lazygit/pkg/integration/integration_tests/commit"
	"github.com/jesseduffield/lazygit/pkg/integration/integration_tests/interactive_rebase"

	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

// Here is where we lists the actual tests that will run. When you create a new test,
// be sure to add it to this list.

var Tests = []types.Test{
	commit.Commit,
	commit.NewBranch,
	branch.Suggestions,
	interactive_rebase.One,
}
