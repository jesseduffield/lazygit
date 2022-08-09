package tests

import (
	"github.com/jesseduffield/lazygit/pkg/integration/helpers"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/branch"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/commit"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/interactive_rebase"
)

// Here is where we lists the actual tests that will run. When you create a new test,
// be sure to add it to this list.

var Tests = []*helpers.Test{
	commit.Commit,
	commit.NewBranch,
	branch.Suggestions,
	interactive_rebase.One,
}
