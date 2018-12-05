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
	pastHeader := false
	for index, line := range lines {
		if strings.HasPrefix(line, "@@") {
			pastHeader = true
			hunkStarts = append(hunkStarts, index)
		}
		if pastHeader && (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "+")) {
			stageableLines = append(stageableLines, index)
		}
	}
	p.Log.WithField("staging", "staging").Info(stageableLines)
	return hunkStarts, stageableLines, nil
}
