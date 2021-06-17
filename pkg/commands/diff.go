package commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

//counterfeiter:generate . IDiffMgr
type IDiffMgr interface {
	WorktreeFileDiff(file *models.File, plain bool, cached bool) string
	WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool) ICmdObj
	ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error)
	ShowFileDiffCmdObj(from string, to string, reverse bool, path string, plain bool, showRenames bool) ICmdObj
	DiffEndArgs(from string, to string, reverse bool, path string) string
	LoadDiffFiles(from string, to string, reverse bool) ([]*models.CommitFile, error)
}

type DiffMgr struct {
	*MgrCtx

	diffFilesLoader *DiffFilesLoader
}

func NewDiffMgr(
	mgrCtx *MgrCtx,
) *DiffMgr {
	return &DiffMgr{
		MgrCtx:          mgrCtx,
		diffFilesLoader: NewDiffFilesLoader(mgrCtx),
	}
}

func (c *DiffMgr) LoadDiffFiles(from string, to string, reverse bool) ([]*models.CommitFile, error) {
	return c.diffFilesLoader.Load(from, to, reverse)
}

// WorktreeFileDiff returns the diff of a file
func (c *DiffMgr) WorktreeFileDiff(file *models.File, plain bool, cached bool) string {
	// for now we assume an error means the file was deleted
	s, _ := c.RunWithOutput(c.WorktreeFileDiffCmdObj(file, plain, cached))
	return s
}

func (c *DiffMgr) WorktreeFileDiffCmdObj(node models.IFile, plain bool, cached bool) ICmdObj {
	path := c.Quote(node.GetPath())

	var colorArg string
	if plain {
		colorArg = "never"
	} else {
		colorArg = c.config.ColorArg()
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
func (c *DiffMgr) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	return c.RunWithOutput(c.ShowFileDiffCmdObj(from, to, reverse, fileName, plain, false))
}

// we may just want to always hide renames, or always show renames. I've combined two functions that were both identical except for that flag here, but I doubt that flag was particularly important.
func (c *DiffMgr) ShowFileDiffCmdObj(from string, to string, reverse bool, path string, plain bool, showRenames bool) ICmdObj {
	colorArg := c.config.ColorArg()
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
func (c *DiffMgr) DiffEndArgs(from string, to string, reverse bool, path string) string {
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
