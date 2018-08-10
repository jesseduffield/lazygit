package main

import (
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing"
)

// context:
// we want to only show 'safe' branches (ones that haven't e.g. been deleted)
// which `git branch -a` gives us, but we also want the recency data that
// git reflog gives us.
// So we get the HEAD, then append get the reflog branches that intersect with
// our safe branches, then add the remaining safe branches, ensuring uniqueness
// along the way

type branchListBuilder struct{}

func newBranchListBuilder() *branchListBuilder {
	return &branchListBuilder{}
}

func (b *branchListBuilder) obtainCurrentBranch() Branch {
	// Using git-go whenever possible
	head, err := r.Head()
	if err != nil {
		panic(err)
	}
	name := head.Name().Short()
	return Branch{Name: name, Recency: "  *"}
}

func (*branchListBuilder) obtainReflogBranches() []Branch {
	branches := make([]Branch, 0)
	rawString, err := runDirectCommand("git reflog -n100 --pretty='%cr|%gs' --grep-reflog='checkout: moving' HEAD")
	if err != nil {
		return branches
	}

	branchLines := splitLines(rawString)
	for _, line := range branchLines {
		timeNumber, timeUnit, branchName := branchInfoFromLine(line)
		timeUnit = abbreviatedTimeUnit(timeUnit)
		branch := Branch{Name: branchName, Recency: timeNumber + timeUnit}
		branches = append(branches, branch)
	}
	return branches
}

func (b *branchListBuilder) obtainSafeBranches() []Branch {
	branches := make([]Branch, 0)

	bIter, err := r.Branches()
	if err != nil {
		panic(err)
	}
	err = bIter.ForEach(func(b *plumbing.Reference) error {
		name := b.Name().Short()
		branches = append(branches, Branch{Name: name})
		return nil
	})

	return branches
}

func (b *branchListBuilder) appendNewBranches(finalBranches, newBranches, existingBranches []Branch, included bool) []Branch {
	for _, newBranch := range newBranches {
		if included == branchIncluded(newBranch.Name, existingBranches) {
			finalBranches = append(finalBranches, newBranch)
		}

	}
	return finalBranches
}

func (b *branchListBuilder) build() []Branch {
	branches := make([]Branch, 0)
	head := b.obtainCurrentBranch()
	validBranches := b.obtainSafeBranches()
	reflogBranches := b.obtainReflogBranches()
	reflogBranches = uniqueByName(append([]Branch{head}, reflogBranches...))

	branches = b.appendNewBranches(branches, reflogBranches, validBranches, true)
	branches = b.appendNewBranches(branches, validBranches, branches, false)

	return branches
}

func uniqueByName(branches []Branch) []Branch {
	finalBranches := make([]Branch, 0)
	for _, branch := range branches {
		if branchIncluded(branch.Name, finalBranches) {
			continue
		}
		finalBranches = append(finalBranches, branch)
	}
	return finalBranches
}

// A line will have the form '10 days ago master' so we need to strip out the
// useful information from that into timeNumber, timeUnit, and branchName
func branchInfoFromLine(line string) (string, string, string) {
	r := regexp.MustCompile("\\|.*\\s")
	line = r.ReplaceAllString(line, " ")
	words := strings.Split(line, " ")
	return words[0], words[1], words[3]
}

func abbreviatedTimeUnit(timeUnit string) string {
	r := regexp.MustCompile("s$")
	timeUnit = r.ReplaceAllString(timeUnit, "")
	timeUnitMap := map[string]string{
		"hour":   "h",
		"minute": "m",
		"second": "s",
		"week":   "w",
		"year":   "y",
		"day":    "d",
		"month":  "m",
	}
	return timeUnitMap[timeUnit]
}
