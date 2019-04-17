package main

import (
	"os"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// An example of how to create and remove branches or any other kind of reference.
func main() {
	CheckArgs("<url>", "<directory>")
	url, directory := os.Args[1], os.Args[2]

	// Clone the given repository to the given directory
	Info("git clone %s %s", url, directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: url,
	})
	CheckIfError(err)

	// Create a new branch to the current HEAD
	Info("git branch my-branch")

	headRef, err := r.Head()
	CheckIfError(err)

	// Create a new plumbing.HashReference object with the name of the branch
	// and the hash from the HEAD. The reference name should be a full reference
	// name and not an abbreviated one, as is used on the git cli.
	//
	// For tags we should use `refs/tags/%s` instead of `refs/heads/%s` used
	// for branches.
	ref := plumbing.NewHashReference("refs/heads/my-branch", headRef.Hash())

	// The created reference is saved in the storage.
	err = r.Storer.SetReference(ref)
	CheckIfError(err)

	// Or deleted from it.
	Info("git branch -D my-branch")
	err = r.Storer.RemoveReference(ref.Name())
	CheckIfError(err)
}
