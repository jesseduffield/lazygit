// +build darwin freebsd linux netbsd openbsd

package jibber_jabber

import (
	"errors"
	"os"
	"strings"
)

func getLangFromEnv() (locale string) {
	locale = os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LANG")
	}
	return
}

func getUnixLocale() (unix_locale string, err error) {
	unix_locale = getLangFromEnv()
	if unix_locale == "" {
		err = errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE)
	}

	return
}

func DetectIETF() (locale string, err error) {
	unix_locale, err := getUnixLocale()
	if err == nil {
		language, territory := splitLocale(unix_locale)
		locale = language
		if territory != "" {
			locale = strings.Join([]string{language, territory}, "-")
		}
	}

	return
}

func DetectLanguage() (language string, err error) {
	unix_locale, err := getUnixLocale()
	if err == nil {
		language, _ = splitLocale(unix_locale)
	}

	return
}

func DetectTerritory() (territory string, err error) {
	unix_locale, err := getUnixLocale()
	if err == nil {
		_, territory = splitLocale(unix_locale)
	}

	return
}
