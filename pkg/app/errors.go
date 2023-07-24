package app

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/samber/lo"
)

type errorMapping struct {
	originalError string
	newError      string
}

// knownError takes an error and tells us whether it's an error that we know about where we can print a nicely formatted version of it rather than panicking with a stack trace
func knownError(tr *i18n.TranslationSet, err error) (string, bool) {
	errorMessage := err.Error()

	knownErrorMessages := []string{tr.MinGitVersionError}

	if lo.Contains(knownErrorMessages, errorMessage) {
		return errorMessage, true
	}

	mappings := []errorMapping{
		{
			originalError: "fatal: not a git repository",
			newError:      tr.NotARepository,
		},
	}

	if mapping, ok := lo.Find(mappings, func(mapping errorMapping) bool {
		return strings.Contains(errorMessage, mapping.originalError)
	}); ok {
		return mapping.newError, true
	}

	return "", false
}
