## encoding


[![Linux](https://img.shields.io/github/actions/workflow/status/gdamore/encoding/linux.yml?branch=main&logoColor=grey&logo=linux&label=)](https://github.com/gdamore/encoding/actions/workflows/linux.yml)
[![Windows](https://img.shields.io/github/actions/workflow/status/gdamore/encoding/windows.yml?branch=main&logoColor=grey&logo=windows&label=)](https://github.com/gdamore/encoding/actions/workflows/windows.yml)
[![Apache License](https://img.shields.io/github/license/gdamore/encoding.svg?logoColor=silver&logo=opensourceinitiative&color=blue&label=)](https://github.com/gdamore/encoding/blob/master/LICENSE)
[![Coverage](https://img.shields.io/codecov/c/github/gdamore/encoding?logoColor=grey&logo=codecov&label=)](https://codecov.io/gh/gdamore/encoding)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/gdamore/encoding)

Package encoding provides a number of encodings that are missing from the
standard Go [encoding]("https://godoc.org/golang.org/x/text/encoding") package.

We hope that we can contribute these to the standard Go library someday.  It
turns out that some of these are useful for dealing with I/O streams coming
from non-UTF friendly sources.

The UTF8 Encoder is also useful for situations where valid UTF-8 might be
carried in streams that contain non-valid UTF; in particular I use it for
helping me cope with terminals that embed escape sequences in otherwise
valid UTF-8.
