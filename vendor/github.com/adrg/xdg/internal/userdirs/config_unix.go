//go:build aix || dragonfly || freebsd || (js && wasm) || nacl || linux || netbsd || openbsd || solaris

package userdirs

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/adrg/xdg/internal/pathutil"
)

// ParseConfigFile parses the user directories config file at the
// specified location.
func ParseConfigFile(name string) (*Directories, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseConfig(f)
}

// ParseConfig parses the user directories config file contained in
// the provided reader.
func ParseConfig(r io.Reader) (*Directories, error) {
	dirs := &Directories{}
	fieldsMap := map[string]*string{
		EnvDesktopDir:     &dirs.Desktop,
		EnvDownloadDir:    &dirs.Download,
		EnvDocumentsDir:   &dirs.Documents,
		EnvMusicDir:       &dirs.Music,
		EnvPicturesDir:    &dirs.Pictures,
		EnvVideosDir:      &dirs.Videos,
		EnvTemplatesDir:   &dirs.Templates,
		EnvPublicShareDir: &dirs.PublicShare,
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		if !strings.HasPrefix(line, "XDG_") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			continue
		}

		// Parse key.
		field, ok := fieldsMap[strings.TrimSpace(parts[0])]
		if !ok {
			continue
		}

		// Parse value.
		runes := []rune(strings.TrimSpace(parts[1]))

		lenRunes := len(runes)
		if lenRunes <= 2 || runes[0] != '"' {
			continue
		}

		for i := 1; i < lenRunes; i++ {
			if runes[i] == '"' {
				*field = pathutil.ExpandHome(string(runes[1:i]))
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dirs, nil
}
