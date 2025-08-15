# Jibber Jabber [![Build Status](https://travis-ci.org/cloudfoundry/jibber_jabber.svg?branch=master)](https://travis-ci.org/cloudfoundry/jibber_jabber)
Jibber Jabber is a GoLang Library that can be used to detect an operating system's current language.

### OS Support

OSX and Linux via the `LC_ALL` and `LANG` environment variables. These are standard variables that are used in ALL versions of UNIX for language detection.

Windows via [GetUserDefaultLocaleName](http://msdn.microsoft.com/en-us/library/windows/desktop/dd318136.aspx) and [GetSystemDefaultLocaleName](http://msdn.microsoft.com/en-us/library/windows/desktop/dd318122.aspx) system calls. These calls are supported in Windows Vista and up.

# Usage
Add the following line to your go `import`:

```
	"github.com/cloudfoundry/jibber_jabber"
```

### DetectIETF
`DetectIETF` will return the current locale as a string. The format of the locale will be the [ISO 639](http://en.wikipedia.org/wiki/ISO_639) two-letter language code, a DASH, then an [ISO 3166](http://en.wikipedia.org/wiki/ISO_3166-1) two-letter country code.

```
	userLocale, err := jibber_jabber.DetectIETF()
	println("Locale:", userLocale)
```

### DetectLanguage
`DetectLanguage` will return the current languge as a string. The format will be the [ISO 639](http://en.wikipedia.org/wiki/ISO_639) two-letter language code.

```
	userLanguage, err := jibber_jabber.DetectLanguage()
	println("Language:", userLanguage)
```

### DetectTerritory
`DetectTerritory` will return the current locale territory as a string. The format will be the [ISO 3166](http://en.wikipedia.org/wiki/ISO_3166-1) two-letter country code.

```
	localeTerritory, err := jibber_jabber.DetectTerritory()
	println("Territory:", localeTerritory)
```

### Errors
All the Detect commands will return an error if they are unable to read the Locale from the system.

For Windows, additional error information is provided due to the nature of the system call being used.
