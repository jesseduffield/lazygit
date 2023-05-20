//go:build aix
// +build aix

package pty

const (
	TIOCGWINSZ = 0
	TIOCSWINSZ = 0
)

func ioctl(fd, cmd, ptr uintptr) error {
	return ErrUnsupported
}
