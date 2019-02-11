package git

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-errors/errors"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
)

type PatchModifier struct {
	Log *logrus.Entry
	Tr  *i18n.Localizer
}

// NewPatchModifier builds a new branch list builder
func NewPatchModifier(log *logrus.Entry) (*PatchModifier, error) {
	return &PatchModifier{
		Log: log,
	}, nil
}

// ModifyPatchForHunk takes the original patch, which may contain several hunks,
// and removes any hunks that aren't the selected hunk
func (p *PatchModifier) ModifyPatchForHunk(patch string, hunkStarts []int, currentLine int) (string, error) {
	// get hunk start and end
	lines := strings.Split(patch, "\n")
	hunkStartIndex := utils.PrevIndex(hunkStarts, currentLine)
	hunkStart := hunkStarts[hunkStartIndex]
	nextHunkStartIndex := utils.NextIndex(hunkStarts, currentLine)
	var hunkEnd int
	if nextHunkStartIndex == 0 {
		hunkEnd = len(lines) - 1
	} else {
		hunkEnd = hunkStarts[nextHunkStartIndex]
	}

	headerLength, err := p.getHeaderLength(lines)
	if err != nil {
		return "", err
	}

	output := strings.Join(lines[0:headerLength], "\n") + "\n"
	output += strings.Join(lines[hunkStart:hunkEnd], "\n") + "\n"

	return output, nil
}

func (p *PatchModifier) getHeaderLength(patchLines []string) (int, error) {
	for index, line := range patchLines {
		if strings.HasPrefix(line, "@@") {
			return index, nil
		}
	}
	return 0, errors.New(p.Tr.SLocalize("CantFindHunks"))
}

// ModifyPatchForLine takes the original patch, which may contain several hunks,
// and the line number of the line we want to stage
func (p *PatchModifier) ModifyPatchForLine(patch string, lineNumber int) (string, error) {
	lines := strings.Split(patch, "\n")
	headerLength, err := p.getHeaderLength(lines)
	if err != nil {
		return "", err
	}
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

	return 0, errors.New(p.Tr.SLocalize("CantFindHunk"))
}

func (p *PatchModifier) getModifiedHunk(patchLines []string, hunkStart int, lineNumber int) ([]string, error) {
	lineChanges := 0
	// strip the hunk down to just the line we want to stage
	newHunk := []string{patchLines[hunkStart]}
	for offsetIndex, line := range patchLines[hunkStart+1:] {
		index := offsetIndex + hunkStart + 1
		if strings.HasPrefix(line, "@@") {
			newHunk = append(newHunk, "\n")
			break
		}
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
	re := regexp.MustCompile(`(\d+) @@`)
	prevLengthString := re.FindStringSubmatch(currentHeader)[1]

	prevLength, err := strconv.Atoi(prevLengthString)
	if err != nil {
		return "", err
	}
	re = regexp.MustCompile(`\d+ @@`)
	newLength := strconv.Itoa(prevLength + lineChanges)
	return re.ReplaceAllString(currentHeader, newLength+" @@"), nil
}
