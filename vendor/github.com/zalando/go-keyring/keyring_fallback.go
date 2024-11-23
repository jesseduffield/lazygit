package keyring

import (
	"errors"
	"runtime"
)

// All of the following methods error out on unsupported platforms
var ErrUnsupportedPlatform = errors.New("unsupported platform: " + runtime.GOOS)

type fallbackServiceProvider struct{}

func (fallbackServiceProvider) Set(service, user, pass string) error {
	return ErrUnsupportedPlatform
}

func (fallbackServiceProvider) Get(service, user string) (string, error) {
	return "", ErrUnsupportedPlatform
}

func (fallbackServiceProvider) Delete(service, user string) error {
	return ErrUnsupportedPlatform
}

func (fallbackServiceProvider) DeleteAll(service string) error {
	return ErrUnsupportedPlatform
}
