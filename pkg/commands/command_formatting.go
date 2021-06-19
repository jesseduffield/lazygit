package commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/utils"
)

// filters a string slice by which are present (non-empty strings), then joins with a space. Intended to be used for constructing CLI commands
func joinPresent(strs ...string) string {
	present := make([]string, 0, len(strs))
	for _, str := range strs {
		if str != "" {
			present = append(present, str)
		}
	}

	return strings.Join(present, " ")
}

func joinPresentKwargs(args map[string]bool) string {
	present := make([]string, 0, len(args))
	for value, include := range args {
		if include {
			present = append(present, value)
		}
	}

	utils.SortAlphabeticalInPlace(present)

	return strings.Join(present, " ")
}
