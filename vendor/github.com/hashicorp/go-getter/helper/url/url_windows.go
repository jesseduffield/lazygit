package url

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

func parse(rawURL string) (*url.URL, error) {
	// Make sure we're using "/" since URLs are "/"-based.
	rawURL = filepath.ToSlash(rawURL)

	if len(rawURL) > 1 && rawURL[1] == ':' {
		// Assume we're dealing with a drive letter. In which case we
		// force the 'file' scheme to avoid "net/url" URL.String() prepending
		// our url with "./".
		rawURL = "file://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if len(u.Host) > 1 && u.Host[1] == ':' && strings.HasPrefix(rawURL, "file://") {
		// Assume we're dealing with a drive letter file path where the drive
		// letter has been parsed into the URL Host.
		u.Path = fmt.Sprintf("%s%s", u.Host, u.Path)
		u.Host = ""
	}

	// Remove leading slash for absolute file paths.
	if len(u.Path) > 2 && u.Path[0] == '/' && u.Path[2] == ':' {
		u.Path = u.Path[1:]
	}

	return u, err
}
