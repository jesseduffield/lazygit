package getter

import (
	"net/url"
	"os"
)

// FileGetter is a Getter implementation that will download a module from
// a file scheme.
type FileGetter struct {
	getter

	// Copy, if set to true, will copy data instead of using a symlink. If
	// false, attempts to symlink to speed up the operation and to lower the
	// disk space usage. If the symlink fails, may attempt to copy on windows.
	Copy bool
}

func (g *FileGetter) ClientMode(u *url.URL) (ClientMode, error) {
	path := u.Path
	if u.RawPath != "" {
		path = u.RawPath
	}

	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	// Check if the source is a directory.
	if fi.IsDir() {
		return ClientModeDir, nil
	}

	return ClientModeFile, nil
}
