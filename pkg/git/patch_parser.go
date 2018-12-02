package git

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type PatchParser struct {
	Log *logrus.Entry
}

// NewPatchParser builds a new branch list builder
func NewPatchParser(log *logrus.Entry) (*PatchParser, error) {
	return &PatchParser{
		Log: log,
	}, nil
}

func (p *PatchParser) ParsePatch(patch string) ([]int, []int, error) {
	lines := strings.Split(patch, "\n")
	hunkStarts := []int{}
	stageableLines := []int{}
	headerLength := 4
	for offsetIndex, line := range lines[headerLength:] {
		index := offsetIndex + headerLength
		if strings.HasPrefix(line, "@@") {
			hunkStarts = append(hunkStarts, index)
		}
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "+") {
			stageableLines = append(stageableLines, index)
		}
	}
	return hunkStarts, stageableLines, nil
}
