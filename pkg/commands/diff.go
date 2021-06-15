package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

// WorktreeFileDiff returns the diff of a file
func (c *Git) WorktreeFileDiff(file *models.File, plain bool, cached bool) string {
	// for now we assume an error means the file was deleted
	s, _ := c.RunWithOutput(c.WorktreeFileDiffCmdObj(file, plain, cached))
	return s
}

func (c *Git) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool) ICmdObj {
	path := c.GetOS().Quote(node.GetPath())

	var colorArg string
	if plain {
		colorArg = "never"
	} else {
		colorArg = c.ColorArg()
	}

	trackedArg := "--"
	if !node.GetIsTracked() && !node.GetHasStagedChanges() && !cached {
		trackedArg = "--no-index -- /dev/null"
	}

	cachedArg := ""
	if cached {
		cachedArg = " --cached"
	}

	return BuildGitCmdObjFromStr(
		fmt.Sprintf(
			"diff --submodule --no-ext-diff%s --color=%s %s %s",
			cachedArg,
			colorArg,
			trackedArg,
			path,
		),
	)
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (c *Git) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	return c.RunWithOutput(c.ShowFileDiffCmdObj(from, to, reverse, fileName, plain, false))
}

// we may just want to always hide renames, or always show renames. I've combined two functions that were both identical except for that flag here, but I doubt that flag was particularly important.
func (c *Git) ShowFileDiffCmdObj(from string, to string, reverse bool, path string, plain bool, showRenames bool) ICmdObj {
	colorArg := c.ColorArg()
	if plain {
		colorArg = "never"
	}

	noRenamesArg := ""
	if !showRenames {
		noRenamesArg = " --no-renames"
	}

	return BuildGitCmdObjFromStr(
		fmt.Sprintf(
			"diff --submodule --no-ext-diff%s --color=%s %s",
			noRenamesArg,
			colorArg,
			c.DiffEndArgs(from, to, reverse, path),
		),
	)
}

// we've got this as a separate function because the GUI wants to display this part to the user
func (c *Git) DiffEndArgs(from string, to string, reverse bool, path string) string {
	output := from
	if to != "" {
		output += " " + to
	}

	if reverse {
		output += " -R"
	}

	if path != "" {
		output += " -- " + path
	}

	return output
}
