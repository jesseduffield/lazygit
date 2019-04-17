package getter

import (
	"fmt"
	"net/url"
	"strings"
)

// GitHubDetector implements Detector to detect GitHub URLs and turn
// them into URLs that the Git Getter can understand.
type GitHubDetector struct{}

func (d *GitHubDetector) Detect(src, _ string) (string, bool, error) {
	if len(src) == 0 {
		return "", false, nil
	}

	if strings.HasPrefix(src, "github.com/") {
		return d.detectHTTP(src)
	}

	return "", false, nil
}

func (d *GitHubDetector) detectHTTP(src string) (string, bool, error) {
	parts := strings.Split(src, "/")
	if len(parts) < 3 {
		return "", false, fmt.Errorf(
			"GitHub URLs should be github.com/username/repo")
	}

	urlStr := fmt.Sprintf("https://%s", strings.Join(parts[:3], "/"))
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", true, fmt.Errorf("error parsing GitHub URL: %s", err)
	}

	if !strings.HasSuffix(url.Path, ".git") {
		url.Path += ".git"
	}

	if len(parts) > 3 {
		url.Path += "//" + strings.Join(parts[3:], "/")
	}

	return "git::" + url.String(), true, nil
}
