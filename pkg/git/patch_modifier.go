package git

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type PatchModifier struct {
	Log *logrus.Entry
}

// NewPatchModifier builds a new branch list builder
func NewPatchModifier(log *logrus.Entry) (*PatchModifier, error) {
	return &PatchModifier{
		Log: log,
	}, nil
}

// ModifyPatch takes the original patch, which may contain several hunks,
// and the line number of the line we want to stage
func (p *PatchModifier) ModifyPatch(patch string, lineNumber int) (string, error) {
	lines := strings.Split(patch, "\n")
	headerLength := 4
	output := strings.Join(lines[0:headerLength], "\n") + "\n"

	hunkStart, err := p.getHunkStart(lines, lineNumber)
	if err != nil {
		return "", err
	}

	hunk, err := p.getModifiedHunk(lines, hunkStart, lineNumber)
	if err != nil {
		return "", err
	}

	output += strings.Join(hunk, "\n")

	return output, nil
}

// getHunkStart returns the line number of the hunk we're going to be modifying
// in order to stage our line
func (p *PatchModifier) getHunkStart(patchLines []string, lineNumber int) (int, error) {
	// find the hunk that we're modifying
	hunkStart := 0
	for index, line := range patchLines {
		if strings.HasPrefix(line, "@@") {
			hunkStart = index
		}
		if index == lineNumber {
			return hunkStart, nil
		}
	}
	return 0, errors.New("Could not find hunk")
}

func (p *PatchModifier) getModifiedHunk(patchLines []string, hunkStart int, lineNumber int) ([]string, error) {
	lineChanges := 0
	// strip the hunk down to just the line we want to stage
	newHunk := []string{}
	for offsetIndex, line := range patchLines[hunkStart:] {
		index := offsetIndex + hunkStart
		if index != lineNumber {
			// we include other removals but treat them like context
			if strings.HasPrefix(line, "-") {
				newHunk = append(newHunk, " "+line[1:])
				lineChanges += 1
				continue
			}
			// we don't include other additions
			if strings.HasPrefix(line, "+") {
				lineChanges -= 1
				continue
			}
		}
		newHunk = append(newHunk, line)
	}

	var err error
	newHunk[0], err = p.updatedHeader(newHunk[0], lineChanges)
	if err != nil {
		return nil, err
	}

	return newHunk, nil
}

// updatedHeader returns the hunk header with the updated line range
// we need to update the hunk length to reflect the changes we made
// if the hunk has three additions but we're only staging one, then
// @@ -14,8 +14,11 @@ import (
// becomes
// @@ -14,8 +14,9 @@ import (
func (p *PatchModifier) updatedHeader(currentHeader string, lineChanges int) (string, error) {
	// current counter is the number after the second comma
	re := regexp.MustCompile(`^[^,]+,[^,]+,(\d+)`)
	prevLengthString := re.FindStringSubmatch(currentHeader)[1]

	prevLength, err := strconv.Atoi(prevLengthString)
	if err != nil {
		return "", err
	}
	re = regexp.MustCompile(`\d+ @@`)
	newLength := strconv.Itoa(prevLength + lineChanges)
	return re.ReplaceAllString(currentHeader, newLength+" @@"), nil
}
