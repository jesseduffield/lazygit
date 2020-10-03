# XDG [![Build status](https://ci.appveyor.com/api/projects/status/9eoupq9jgsu2p0jv?svg=true)](https://ci.appveyor.com/project/dixonwille/xdg) [![Build Status](https://travis-ci.org/OpenPeeDeeP/xdg.svg?branch=master)](https://travis-ci.org/OpenPeeDeeP/xdg) [![Go Report Card](https://goreportcard.com/badge/github.com/OpenPeeDeeP/xdg)](https://goreportcard.com/report/github.com/OpenPeeDeeP/xdg) [![GoDoc](https://godoc.org/github.com/OpenPeeDeeP/xdg?status.svg)](https://godoc.org/github.com/OpenPeeDeeP/xdg) [![codecov](https://codecov.io/gh/OpenPeeDeeP/xdg/branch/master/graph/badge.svg)](https://codecov.io/gh/OpenPeeDeeP/xdg)

A cross platform package that tries to follow [XDG Standard](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) when possible. Since XDG is linux specific, I am only able to follow standards to the T on linux. But for the other operating systems I am finding similar best practice locations for the files.

## Locations Per OS

The following table shows what is used if the envrionment variable is not set. If the variable is set then this package uses that. Linux follows the default standards. Mac does when it comes to the home directory but for system wide it uses the standard `/Library/Application Support`. As for Windows, the variable defaults are just other environment variables set up by the operation system.

> When creating `XDG` application the `Vendor` and `Application` names are appeneded to the end of the path to keep projects unique.

|  | Linux(and BSD) | Mac | Windows |
| ---: | :---: | :---: | :---: |
| `XDG_DATA_DIRS` | [`/usr/local/share`, `/usr/share`] | [`/Library/Application Support`] | `%PROGRAMDATA%` |
| `XDG_DATA_HOME` | `~/.local/share` | `~/Library/Application Support` | `%APPDATA%` |
| `XDG_CONFIG_DIRS` | [`/etc/xdg`] | [`/Library/Application Support`] | `%PROGRAMDATA%` |
| `XDG_CONFIG_HOME` | `~/.config` | `~/Library/Application Support` | `%APPDATA%` |
| `XDG_CACHE_HOME` | `~/.cache` | `~/Library/Caches` | `%LOCALAPPDATA%` |

## Notes

- This package does not merge files if they exist across different directories.
- The `Query` methods search through the system variables, `DIRS`, first (when using environment variables first in the variable has presidence). It then checks home variables, `HOME`.
- This package will not create any directories for you. In the standard, it states the following:

> If, when attempting to write a file, the destination directory is non-existant an attempt should be made to create it with permission `0700`. If the destination directory exists already the permissions should not be changed. The application should be prepared to handle the case where the file could not be written, either because the directory was non-existant and could not be created, or for any other reason. In such case it may chose to present an error message to the user.
