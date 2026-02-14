package ai

import (
	"regexp"
	"strings"
)

const maxDiffSize = 50000

var lockFilePatterns = []*regexp.Regexp{
	regexp.MustCompile(`package-lock\.json$`),
	regexp.MustCompile(`yarn\.lock$`),
	regexp.MustCompile(`pnpm-lock\.yaml$`),
	regexp.MustCompile(`composer\.lock$`),
	regexp.MustCompile(`Gemfile\.lock$`),
	regexp.MustCompile(`Pipfile\.lock$`),
	regexp.MustCompile(`poetry\.lock$`),
	regexp.MustCompile(`(?i)cargo\.lock$`),
	regexp.MustCompile(`go\.sum$`),
	regexp.MustCompile(`mix\.lock$`),
}

type ContextBuilder struct{}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{}
}

func (cb *ContextBuilder) BuildContext(diff string) string {
	filtered := cb.filterLockFiles(diff)

	if len(filtered) > maxDiffSize {
		filtered = filtered[:maxDiffSize] + "\n... (diff truncated)"
	}

	return filtered
}

func (cb *ContextBuilder) filterLockFiles(diff string) string {
	var result strings.Builder
	lines := strings.Split(diff, "\n")

	inLockFile := false

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			inLockFile = cb.isLockFile(line)
			if !inLockFile {
				result.WriteString(line)
				result.WriteString("\n")
			}
			continue
		}

		if !inLockFile {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String()
}

func (cb *ContextBuilder) isLockFile(diffHeaderLine string) bool {
	for _, pattern := range lockFilePatterns {
		if pattern.MatchString(diffHeaderLine) {
			return true
		}
	}
	return false
}
