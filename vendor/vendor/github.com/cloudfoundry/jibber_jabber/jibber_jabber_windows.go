// +build windows

package jibber_jabber

import (
	"errors"
	"syscall"
	"unsafe"
)

const LOCALE_NAME_MAX_LENGTH uint32 = 85

var SUPPORTED_LOCALES = map[uintptr]string{
	0x0407: "de-DE",
	0x0409: "en-US",
	0x0c0a: "es-ES", //or is it 0x040a
	0x040c: "fr-FR",
	0x0410: "it-IT",
	0x0411: "ja-JA",
	0x0412: "ko_KR",
	0x0416: "pt-BR",
	//0x0419: "ru_RU", - Will add support for Russian when nicksnyder/go-i18n supports Russian
	0x0804: "zh-CN",
	0x0c04: "zh-HK",
	0x0404: "zh-TW",
}

func getWindowsLocaleFrom(sysCall string) (locale string, err error) {
	buffer := make([]uint16, LOCALE_NAME_MAX_LENGTH)

	dll := syscall.MustLoadDLL("kernel32")
	proc := dll.MustFindProc(sysCall)
	r, _, dllError := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(LOCALE_NAME_MAX_LENGTH))
	if r == 0 {
		err = errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE + ":\n" + dllError.Error())
		return
	}

	locale = syscall.UTF16ToString(buffer)

	return
}

func getAllWindowsLocaleFrom(sysCall string) (string, error) {
	dll, err := syscall.LoadDLL("kernel32")
	if err != nil {
		return "", errors.New("Could not find kernel32 dll")
	}

	proc, err := dll.FindProc(sysCall)
	if err != nil {
		return "", err
	}

	locale, _, dllError := proc.Call()
	if locale == 0 {
		return "", errors.New(COULD_NOT_DETECT_PACKAGE_ERROR_MESSAGE + ":\n" + dllError.Error())
	}

	return SUPPORTED_LOCALES[locale], nil
}

func getWindowsLocale() (locale string, err error) {
	dll, err := syscall.LoadDLL("kernel32")
	if err != nil {
		return "", errors.New("Could not find kernel32 dll")
	}

	proc, err := dll.FindProc("GetVersion")
	if err != nil {
		return "", err
	}

	v, _, _ := proc.Call()
	windowsVersion := byte(v)
	isVistaOrGreater := (windowsVersion >= 6)

	if isVistaOrGreater {
		locale, err = getWindowsLocaleFrom("GetUserDefaultLocaleName")
		if err != nil {
			locale, err = getWindowsLocaleFrom("GetSystemDefaultLocaleName")
		}
	} else if !isVistaOrGreater {
		locale, err = getAllWindowsLocaleFrom("GetUserDefaultLCID")
		if err != nil {
			locale, err = getAllWindowsLocaleFrom("GetSystemDefaultLCID")
		}
	} else {
		panic(v)
	}
	return
}
func DetectIETF() (locale string, err error) {
	locale, err = getWindowsLocale()
	return
}

func DetectLanguage() (language string, err error) {
	windows_locale, err := getWindowsLocale()
	if err == nil {
		language, _ = splitLocale(windows_locale)
	}

	return
}

func DetectTerritory() (territory string, err error) {
	windows_locale, err := getWindowsLocale()
	if err == nil {
		_, territory = splitLocale(windows_locale)
	}

	return
}
