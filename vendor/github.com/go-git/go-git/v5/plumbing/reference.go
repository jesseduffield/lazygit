package plumbing

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	refPrefix       = "refs/"
	refHeadPrefix   = refPrefix + "heads/"
	refTagPrefix    = refPrefix + "tags/"
	refRemotePrefix = refPrefix + "remotes/"
	refNotePrefix   = refPrefix + "notes/"
	symrefPrefix    = "ref: "
)

// RefRevParseRules are a set of rules to parse references into short names, or expand into a full reference.
// These are the same rules as used by git in shorten_unambiguous_ref and expand_ref.
// See: https://github.com/git/git/blob/e0aaa1b6532cfce93d87af9bc813fb2e7a7ce9d7/refs.c#L417
var RefRevParseRules = []string{
	"%s",
	"refs/%s",
	"refs/tags/%s",
	"refs/heads/%s",
	"refs/remotes/%s",
	"refs/remotes/%s/HEAD",
}

var (
	ErrReferenceNotFound = errors.New("reference not found")

	// ErrInvalidReferenceName is returned when a reference name is invalid.
	ErrInvalidReferenceName = errors.New("invalid reference name")
)

// ReferenceType reference type's
type ReferenceType int8

const (
	InvalidReference  ReferenceType = 0
	HashReference     ReferenceType = 1
	SymbolicReference ReferenceType = 2
)

func (r ReferenceType) String() string {
	switch r {
	case InvalidReference:
		return "invalid-reference"
	case HashReference:
		return "hash-reference"
	case SymbolicReference:
		return "symbolic-reference"
	}

	return ""
}

// ReferenceName reference name's
type ReferenceName string

// NewBranchReferenceName returns a reference name describing a branch based on
// his short name.
func NewBranchReferenceName(name string) ReferenceName {
	return ReferenceName(refHeadPrefix + name)
}

// NewNoteReferenceName returns a reference name describing a note based on his
// short name.
func NewNoteReferenceName(name string) ReferenceName {
	return ReferenceName(refNotePrefix + name)
}

// NewRemoteReferenceName returns a reference name describing a remote branch
// based on his short name and the remote name.
func NewRemoteReferenceName(remote, name string) ReferenceName {
	return ReferenceName(refRemotePrefix + fmt.Sprintf("%s/%s", remote, name))
}

// NewRemoteHEADReferenceName returns a reference name describing a the HEAD
// branch of a remote.
func NewRemoteHEADReferenceName(remote string) ReferenceName {
	return ReferenceName(refRemotePrefix + fmt.Sprintf("%s/%s", remote, HEAD))
}

// NewTagReferenceName returns a reference name describing a tag based on short
// his name.
func NewTagReferenceName(name string) ReferenceName {
	return ReferenceName(refTagPrefix + name)
}

// IsBranch check if a reference is a branch
func (r ReferenceName) IsBranch() bool {
	return strings.HasPrefix(string(r), refHeadPrefix)
}

// IsNote check if a reference is a note
func (r ReferenceName) IsNote() bool {
	return strings.HasPrefix(string(r), refNotePrefix)
}

// IsRemote check if a reference is a remote
func (r ReferenceName) IsRemote() bool {
	return strings.HasPrefix(string(r), refRemotePrefix)
}

// IsTag check if a reference is a tag
func (r ReferenceName) IsTag() bool {
	return strings.HasPrefix(string(r), refTagPrefix)
}

func (r ReferenceName) String() string {
	return string(r)
}

// Short returns the short name of a ReferenceName
func (r ReferenceName) Short() string {
	s := string(r)
	res := s
	for _, format := range RefRevParseRules[1:] {
		_, err := fmt.Sscanf(s, format, &res)
		if err == nil {
			continue
		}
	}

	return res
}

var (
	ctrlSeqs = regexp.MustCompile(`[\000-\037\177]`)
)

// Validate validates a reference name.
// This follows the git-check-ref-format rules.
// See https://git-scm.com/docs/git-check-ref-format
//
// It is important to note that this function does not check if the reference
// exists in the repository.
// It only checks if the reference name is valid.
// This functions does not support the --refspec-pattern, --normalize, and
// --allow-onelevel options.
//
// Git imposes the following rules on how references are named:
//
//  1. They can include slash / for hierarchical (directory) grouping, but no
//     slash-separated component can begin with a dot . or end with the
//     sequence .lock.
//  2. They must contain at least one /. This enforces the presence of a
//     category like heads/, tags/ etc. but the actual names are not
//     restricted. If the --allow-onelevel option is used, this rule is
//     waived.
//  3. They cannot have two consecutive dots .. anywhere.
//  4. They cannot have ASCII control characters (i.e. bytes whose values are
//     lower than \040, or \177 DEL), space, tilde ~, caret ^, or colon :
//     anywhere.
//  5. They cannot have question-mark ?, asterisk *, or open bracket [
//     anywhere. See the --refspec-pattern option below for an exception to this
//     rule.
//  6. They cannot begin or end with a slash / or contain multiple consecutive
//     slashes (see the --normalize option below for an exception to this rule).
//  7. They cannot end with a dot ..
//  8. They cannot contain a sequence @{.
//  9. They cannot be the single character @.
//  10. They cannot contain a \.
func (r ReferenceName) Validate() error {
	s := string(r)
	if len(s) == 0 {
		return ErrInvalidReferenceName
	}

	// HEAD is a special case
	if r == HEAD {
		return nil
	}

	// rule 7
	if strings.HasSuffix(s, ".") {
		return ErrInvalidReferenceName
	}

	// rule 2
	parts := strings.Split(s, "/")
	if len(parts) < 2 {
		return ErrInvalidReferenceName
	}

	isBranch := r.IsBranch()
	isTag := r.IsTag()
	for i, part := range parts {
		// rule 6
		if len(part) == 0 {
			return ErrInvalidReferenceName
		}

		if strings.HasPrefix(part, ".") || // rule 1
			strings.Contains(part, "..") || // rule 3
			ctrlSeqs.MatchString(part) || // rule 4
			strings.ContainsAny(part, "~^:?*[ \t\n") || // rule 4 & 5
			strings.Contains(part, "@{") || // rule 8
			part == "@" || // rule 9
			strings.Contains(part, "\\") || // rule 10
			strings.HasSuffix(part, ".lock") { // rule 1
			return ErrInvalidReferenceName
		}

		if (isBranch || isTag) && strings.HasPrefix(part, "-") && (i == 2) { // branches & tags can't start with -
			return ErrInvalidReferenceName
		}
	}

	return nil
}

const (
	HEAD   ReferenceName = "HEAD"
	Master ReferenceName = "refs/heads/master"
	Main   ReferenceName = "refs/heads/main"
)

// Reference is a representation of git reference
type Reference struct {
	t      ReferenceType
	n      ReferenceName
	h      Hash
	target ReferenceName
}

// NewReferenceFromStrings creates a reference from name and target as string,
// the resulting reference can be a SymbolicReference or a HashReference base
// on the target provided
func NewReferenceFromStrings(name, target string) *Reference {
	n := ReferenceName(name)

	if strings.HasPrefix(target, symrefPrefix) {
		target := ReferenceName(target[len(symrefPrefix):])
		return NewSymbolicReference(n, target)
	}

	return NewHashReference(n, NewHash(target))
}

// NewSymbolicReference creates a new SymbolicReference reference
func NewSymbolicReference(n, target ReferenceName) *Reference {
	return &Reference{
		t:      SymbolicReference,
		n:      n,
		target: target,
	}
}

// NewHashReference creates a new HashReference reference
func NewHashReference(n ReferenceName, h Hash) *Reference {
	return &Reference{
		t: HashReference,
		n: n,
		h: h,
	}
}

// Type returns the type of a reference
func (r *Reference) Type() ReferenceType {
	return r.t
}

// Name returns the name of a reference
func (r *Reference) Name() ReferenceName {
	return r.n
}

// Hash returns the hash of a hash reference
func (r *Reference) Hash() Hash {
	return r.h
}

// Target returns the target of a symbolic reference
func (r *Reference) Target() ReferenceName {
	return r.target
}

// Strings dump a reference as a [2]string
func (r *Reference) Strings() [2]string {
	var o [2]string
	o[0] = r.Name().String()

	switch r.Type() {
	case HashReference:
		o[1] = r.Hash().String()
	case SymbolicReference:
		o[1] = symrefPrefix + r.Target().String()
	}

	return o
}

func (r *Reference) String() string {
	ref := ""
	switch r.Type() {
	case HashReference:
		ref = r.Hash().String()
	case SymbolicReference:
		ref = symrefPrefix + r.Target().String()
	default:
		return ""
	}

	name := r.Name().String()
	var v strings.Builder
	v.Grow(len(ref) + len(name) + 1)
	v.WriteString(ref)
	v.WriteString(" ")
	v.WriteString(name)
	return v.String()
}
