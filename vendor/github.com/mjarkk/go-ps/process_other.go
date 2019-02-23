// +build netbsd,arm openbsd,amd64 openbsd,386 freebsd,386 netbsd,amd64 freebsd,arm netbsd,386

package ps

import "errors"

// processes is a fallback function the real functions are in the files spesific to an OS and platform
func processes() ([]Process, error) {
	return []Process{}, errors.New("OS or Platform not supported")
}

// findProcess is a fallback function the real functions are in the files spesific to an OS
func findProcess(pid int) (Process, error) {
	return nil, errors.New("OS or Platform not supported")
}

// supports returns false because of the fallback functions above
func supported() bool {
	return false
}
