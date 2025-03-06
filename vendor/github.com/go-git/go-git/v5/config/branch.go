package config

import (
	"errors"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	format "github.com/go-git/go-git/v5/plumbing/format/config"
)

var (
	errBranchEmptyName     = errors.New("branch config: empty name")
	errBranchInvalidMerge  = errors.New("branch config: invalid merge")
	errBranchInvalidRebase = errors.New("branch config: rebase must be one of 'true' or 'interactive'")
)

// Branch contains information on the
// local branches and which remote to track
type Branch struct {
	// Name of branch
	Name string
	// Remote name of remote to track
	Remote string
	// Merge is the local refspec for the branch
	Merge plumbing.ReferenceName
	// Rebase instead of merge when pulling. Valid values are
	// "true" and "interactive".  "false" is undocumented and
	// typically represented by the non-existence of this field
	Rebase string
	// Description explains what the branch is for.
	// Multi-line explanations may be used.
	//
	// Original git command to edit:
	//	git branch --edit-description
	Description string

	raw *format.Subsection
}

// Validate validates fields of branch
func (b *Branch) Validate() error {
	if b.Name == "" {
		return errBranchEmptyName
	}

	if b.Merge != "" && !b.Merge.IsBranch() {
		return errBranchInvalidMerge
	}

	if b.Rebase != "" &&
		b.Rebase != "true" &&
		b.Rebase != "interactive" &&
		b.Rebase != "false" {
		return errBranchInvalidRebase
	}

	return plumbing.NewBranchReferenceName(b.Name).Validate()
}

func (b *Branch) marshal() *format.Subsection {
	if b.raw == nil {
		b.raw = &format.Subsection{}
	}

	b.raw.Name = b.Name

	if b.Remote == "" {
		b.raw.RemoveOption(remoteSection)
	} else {
		b.raw.SetOption(remoteSection, b.Remote)
	}

	if b.Merge == "" {
		b.raw.RemoveOption(mergeKey)
	} else {
		b.raw.SetOption(mergeKey, string(b.Merge))
	}

	if b.Rebase == "" {
		b.raw.RemoveOption(rebaseKey)
	} else {
		b.raw.SetOption(rebaseKey, b.Rebase)
	}

	if b.Description == "" {
		b.raw.RemoveOption(descriptionKey)
	} else {
		desc := quoteDescription(b.Description)
		b.raw.SetOption(descriptionKey, desc)
	}

	return b.raw
}

// hack to trigger conditional quoting in the
// plumbing/format/config/Encoder.encodeOptions
//
// Current Encoder implementation uses Go %q format if value contains a backslash character,
// which is not consistent with reference git implementation.
// git just replaces newline characters with \n, while Encoder prints them directly.
// Until value quoting fix, we should escape description value by replacing newline characters with \n.
func quoteDescription(desc string) string {
	return strings.ReplaceAll(desc, "\n", `\n`)
}

func (b *Branch) unmarshal(s *format.Subsection) error {
	b.raw = s

	b.Name = b.raw.Name
	b.Remote = b.raw.Options.Get(remoteSection)
	b.Merge = plumbing.ReferenceName(b.raw.Options.Get(mergeKey))
	b.Rebase = b.raw.Options.Get(rebaseKey)
	b.Description = unquoteDescription(b.raw.Options.Get(descriptionKey))

	return b.Validate()
}

// hack to enable conditional quoting in the
// plumbing/format/config/Encoder.encodeOptions
// goto quoteDescription for details.
func unquoteDescription(desc string) string {
	return strings.ReplaceAll(desc, `\n`, "\n")
}
