package app

import (
	"fmt"
	"os"
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

	knownErrorMessages := []string{minGitVersionErrorMessage(tr)}

	if lo.Contains(knownErrorMessages, errorMessage) {
		return errorMessage, true
	}

	mappings := []errorMapping{
		{
			originalError: "fatal: not a git repository",
			newError:      tr.NotARepository,
		},
		{
			originalError: "getwd: no such file or directory",
			newError:      tr.WorkingDirectoryDoesNotExist,
		},
		{
			originalError: "terminal entry not found: term not set",
			newError:      tr.TermNotSet,
		},
		{
			originalError: "$TERM Not Found",
			newError:      fmt.Sprintf(tr.TermNotFound, os.Getenv("TERM")),
		},
	}

	if mapping, ok := lo.Find(mappings, func(mapping errorMapping) bool {
		return strings.Contains(errorMessage, mapping.originalError)
	}); ok {
		return mapping.newError, true
	}

	return "", false
}
