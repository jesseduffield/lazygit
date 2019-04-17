package getter

// GitDetector implements Detector to detect Git SSH URLs such as
// git@host.com:dir1/dir2 and converts them to proper URLs.
type GitDetector struct{}

func (d *GitDetector) Detect(src, _ string) (string, bool, error) {
	if len(src) == 0 {
		return "", false, nil
	}

	u, err := detectSSH(src)
	if err != nil {
		return "", true, err
	}
	if u == nil {
		return "", false, nil
	}

	// We require the username to be "git" to assume that this is a Git URL
	if u.User.Username() != "git" {
		return "", false, nil
	}

	return "git::" + u.String(), true, nil
}
