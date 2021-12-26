package utils

import "regexp"

func FindNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	if len(match) == 0 {
		return nil
	}

	results := map[string]string{}
	for i, value := range match[1:] {
		results[regex.SubexpNames()[i+1]] = value
	}
	return results
}
