package utils

import (
	"log"
	"os"
	"path/filepath"
)

// Min returns the minimum of two integers
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum of two integers
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Clamp returns a value x restricted between min and max
func Clamp(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

// GetLazyRootDirectory finds a lazy project root directory.
//
// It's used for cheatsheet scripts and integration tests. Not to be confused with finding the
// root directory of _any_ random repo.
func GetLazyRootDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		_, err := os.Stat(filepath.Join(path, ".git"))
		if err == nil {
			return path
		}

		if !os.IsNotExist(err) {
			panic(err)
		}

		path = filepath.Dir(path)

		if path == "/" {
			log.Fatal("must run in lazy project folder or child folder")
		}
	}
}
