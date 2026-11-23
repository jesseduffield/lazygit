package terminfo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// termCache is the terminfo cache.
var termCache = struct {
	sync.RWMutex

	db map[string]*Terminfo
}{
	db: make(map[string]*Terminfo),
}

// Load follows the behavior described in terminfo(5) to find correct the
// terminfo file using the name, reads the file and then returns a Terminfo
// struct that describes the file. It returns ErrFileNotFound when no
// directory holds the file.
//
// Load caches each entry that it reads, and returns the same *Terminfo to
// every caller. A caller must not change the returned value or its maps,
// because the change is then seen by every other caller.
func Load(name string) (*Terminfo, error) {
	if name == "" {
		return nil, ErrEmptyTermName
	}
	termCache.RLock()
	ti, ok := termCache.db[name]
	termCache.RUnlock()
	if ok {
		return ti, nil
	}
	var checkDirs []string
	// check $TERMINFO
	if dir := os.Getenv("TERMINFO"); dir != "" {
		checkDirs = append(checkDirs, dir)
	}
	// check $HOME/.terminfo, as ncurses does, and skip it when $HOME is not set
	if home := os.Getenv("HOME"); home != "" {
		checkDirs = append(checkDirs, filepath.Join(home, ".terminfo"))
	}
	// check $TERMINFO_DIRS
	if dirs := os.Getenv("TERMINFO_DIRS"); dirs != "" {
		checkDirs = append(checkDirs, strings.Split(dirs, ":")...)
	}
	// check fallback directories
	checkDirs = append(checkDirs, "/etc/terminfo", "/lib/terminfo", "/usr/share/terminfo")
	for _, dir := range checkDirs {
		ti, err := Open(dir, name)
		if err != nil && !errors.Is(err, ErrFileNotFound) && !os.IsNotExist(err) {
			return nil, err
		} else if err == nil {
			return ti, nil
		}
	}
	return nil, ErrFileNotFound
}

// LoadFromEnv loads the terminal info based on the name contained in
// environment variable TERM.
func LoadFromEnv() (*Terminfo, error) {
	return Load(os.Getenv("TERM"))
}
