package humanlog

import (
	"regexp"
)

// dcLogsPrefixRe parses out a prefix like 'web_1 | ' from docker-compose
// The regex exists of five parts:
// 1. An optional color terminal escape sequence
// 2. The name of the service
// 3. Any number of spaces, and a pipe symbol
// 4. An optional color reset escape sequence
// 5. The rest of the line
var dcLogsPrefixRe = regexp.MustCompile("^(?:\x1b\\[\\d+m)?(?P<service_name>[a-zA-Z0-9._-]+)\\s+\\|(?:\x1b\\[0m)? (?P<rest_of_line>.*)$")

type handler interface {
	TryHandle([]byte) bool
	setField(key, val []byte)
}

func tryDockerComposePrefix(d []byte, nextHandler handler) bool {
	if matches := dcLogsPrefixRe.FindSubmatch(d); matches != nil {
		if nextHandler.TryHandle(matches[2]) {
			nextHandler.setField([]byte(`service`), matches[1])
			return true
		}
	}
	return false
}
