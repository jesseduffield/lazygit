package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// Unique eliminates the duplicate paths from the provided slice and returns
// the result. The paths are expanded using the `ExpandHome` function and only
// absolute paths are kept. The items in the output slice are in the order in
// which they occur in the input slice.
func Unique(paths []string) []string {
	var (
		uniq     []string
		registry = map[string]struct{}{}
	)

	for _, p := range paths {
		if p = ExpandHome(p); p != "" && filepath.IsAbs(p) {
			if _, ok := registry[p]; ok {
				continue
			}

			registry[p] = struct{}{}
			uniq = append(uniq, p)
		}
	}

	return uniq
}

// First returns the first absolute path from the provided slice.
// The paths in the input slice are expanded using the `ExpandHome` function.
func First(paths []string) string {
	for _, p := range paths {
		if p = ExpandHome(p); p != "" && filepath.IsAbs(p) {
			return p
		}
	}

	return ""
}

// Create returns a suitable location relative to which the file with the
// specified `name` can be written. The first path from the provided `paths`
// slice which is successfully created (or already exists) is used as a base
// path for the file. The `name` parameter should contain the name of the file
// which is going to be written in the location returned by this function, but
// it can also contain a set of parent directories, which will be created
// relative to the selected parent path.
func Create(name string, paths []string) (string, error) {
	searchedPaths := make([]string, 0, len(paths))
	for _, p := range paths {
		p = filepath.Join(p, name)

		dir := filepath.Dir(p)
		if Exists(dir) {
			return p, nil
		}
		if err := os.MkdirAll(dir, os.ModeDir|0o700); err == nil {
			return p, nil
		}

		searchedPaths = append(searchedPaths, dir)
	}

	return "", fmt.Errorf("could not create any of the following paths: %v",
		searchedPaths)
}

// Search searches for the file with the specified `name` in the provided
// slice of `paths`. The `name` parameter must contain the name of the file,
// but it can also contain a set of parent directories.
func Search(name string, paths []string) (string, error) {
	searchedPaths := make([]string, 0, len(paths))
	for _, p := range paths {
		p = filepath.Join(p, name)
		if Exists(p) {
			return p, nil
		}

		searchedPaths = append(searchedPaths, filepath.Dir(p))
	}

	return "", fmt.Errorf("could not locate `%s` in any of the following paths: %v",
		filepath.Base(name), searchedPaths)
}

// EnvPath returns the value of the environment variable with the specified
// `name` if it is an absolute path, or the first absolute fallback path.
// All paths are expanded using the `ExpandHome` function.
func EnvPath(name string, fallbackPaths ...string) string {
	dir := ExpandHome(os.Getenv(name))
	if dir != "" && filepath.IsAbs(dir) {
		return dir
	}

	return First(fallbackPaths)
}

// EnvPathList reads the value of the environment variable with the specified
// `name` and attempts to extract a list of absolute paths from it. If there
// are none, a list of absolute fallback paths is returned instead. Duplicate
// paths are removed from the returned slice. All paths are expanded using the
// `ExpandHome` function.
func EnvPathList(name string, fallbackPaths ...string) []string {
	dirs := Unique(filepath.SplitList(os.Getenv(name)))
	if len(dirs) != 0 {
		return dirs
	}

	return Unique(fallbackPaths)
}
