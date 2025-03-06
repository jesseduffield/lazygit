package config

import (
	"errors"
	"strings"

	format "github.com/go-git/go-git/v5/plumbing/format/config"
)

var (
	errURLEmptyInsteadOf = errors.New("url config: empty insteadOf")
)

// Url defines Url rewrite rules
type URL struct {
	// Name new base url
	Name string
	// Any URL that starts with this value will be rewritten to start, instead, with <base>.
	// When more than one insteadOf strings match a given URL, the longest match is used.
	InsteadOf string

	// raw representation of the subsection, filled by marshal or unmarshal are
	// called.
	raw *format.Subsection
}

// Validate validates fields of branch
func (b *URL) Validate() error {
	if b.InsteadOf == "" {
		return errURLEmptyInsteadOf
	}

	return nil
}

const (
	insteadOfKey = "insteadOf"
)

func (u *URL) unmarshal(s *format.Subsection) error {
	u.raw = s

	u.Name = s.Name
	u.InsteadOf = u.raw.Option(insteadOfKey)
	return nil
}

func (u *URL) marshal() *format.Subsection {
	if u.raw == nil {
		u.raw = &format.Subsection{}
	}

	u.raw.Name = u.Name
	u.raw.SetOption(insteadOfKey, u.InsteadOf)

	return u.raw
}

func findLongestInsteadOfMatch(remoteURL string, urls map[string]*URL) *URL {
	var longestMatch *URL
	for _, u := range urls {
		if !strings.HasPrefix(remoteURL, u.InsteadOf) {
			continue
		}

		// according to spec if there is more than one match, take the logest
		if longestMatch == nil || len(longestMatch.InsteadOf) < len(u.InsteadOf) {
			longestMatch = u
		}
	}

	return longestMatch
}

func (u *URL) ApplyInsteadOf(url string) string {
	if !strings.HasPrefix(url, u.InsteadOf) {
		return url
	}

	return u.Name + url[len(u.InsteadOf):]
}
