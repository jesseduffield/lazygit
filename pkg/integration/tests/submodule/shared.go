package submodule

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

func setupNestedSubmodules(shell *Shell) {
	// we're going to have a directory structure like this:
	// project
	//  - repo/modules/outerSubName/modules/innerSubName/
	//
	shell.CreateFileAndAdd("rootFile", "rootStuff")
	shell.Commit("initial repo commit")

	shell.Chdir("..")
	shell.CreateDir("innerSubmodule")
	shell.Chdir("innerSubmodule")
	shell.Init()
	shell.CreateFileAndAdd("inner", "inner")
	shell.Commit("initial inner commit")

	shell.Chdir("..")
	shell.CreateDir("outerSubmodule")
	shell.Chdir("outerSubmodule")
	shell.Init()
	shell.CreateFileAndAdd("outer", "outer")
	shell.Commit("initial outer commit")
	shell.CreateDir("modules")
	// the git config (-c) parameter below is required
	// to let git create a file-protocol/path submodule
	shell.RunCommand([]string{"git", "-c", "protocol.file.allow=always", "submodule", "add", "--name", "innerSubName", "../innerSubmodule", "modules/innerSubPath"})
	shell.Commit("add dependency as innerSubmodule")

	shell.Chdir("../repo")
	shell.CreateDir("modules")
	shell.RunCommand([]string{"git", "-c", "protocol.file.allow=always", "submodule", "add", "--name", "outerSubName", "../outerSubmodule", "modules/outerSubPath"})
	shell.Commit("add dependency as outerSubmodule")
	shell.RunCommand([]string{"git", "-c", "protocol.file.allow=always", "submodule", "update", "--init", "--recursive"})
}
