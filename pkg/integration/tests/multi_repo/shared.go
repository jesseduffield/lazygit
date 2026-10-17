package multi_repo

import (
	"github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Git finds the lazygit source repo above the test directory unless the ceiling
// stops it, and then lazygit would never see a non-repo start directory.
var ceilingEnv = map[string]string{"GIT_CEILING_DIRECTORIES": "{{actualPath}}"}

// setupWorkspace creates actual/workspace with three repos and leaves the shell
// in it, so that lazygit starts there. alpha is on branch feature, beta has an
// untracked file, and org/gamma is one level deeper.
func setupWorkspace(shell *components.Shell) {
	shell.Chdir("..")
	shell.CreateDir("workspace")
	shell.Chdir("workspace")

	for _, name := range []string{"alpha", "beta", "org/gamma"} {
		shell.CreateDir(name)
		shell.Chdir(name)
		shell.Init()
		shell.EmptyCommit("initial")
		if name == "alpha" {
			shell.NewBranch("feature")
		}
		if name == "beta" {
			shell.CreateFile("untracked.txt", "content")
		}
		shell.Chdir(relativeBack(name))
	}
}

func relativeBack(name string) string {
	if name == "org/gamma" {
		return "../.."
	}
	return ".."
}
