package presentation

func GetSuggestionListDisplayStrings(suggestions []string) [][]string {
	lines := make([][]string, len(suggestions))

	for i := range suggestions {
		lines[i] = getSuggestionDisplayStrings(suggestions[i])
	}

	return lines
}

func getSuggestionDisplayStrings(suggestion string) []string {
	return []string{suggestion}
}
