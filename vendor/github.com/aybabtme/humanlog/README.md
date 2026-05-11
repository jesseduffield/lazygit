# humanlog

Read logs from `stdin` and prints them back to `stdout`, but prettier.

# Using it

[Grab a release](https://github.com/aybabtme/humanlog/releases) or :

## With Go installed
```bash
$ go get -u github.com/aybabtme/humanlog/...
```

## On linux

```bash
wget -qO- https://github.com/aybabtme/humanlog/releases/download/0.4.0/humanlog_Linux_x86_64.tar.gz | tar xvz
```

## On OS X

```bash
brew tap aybabtme/homebrew-tap
brew install humanlog
```

# Example

If you emit logs in JSON or in [`logfmt`](https://brandur.org/logfmt), you will enjoy pretty logs when those
entries are encountered by `humanlog`. Unrecognized lines are left unchanged.

```
$ humanlog < /var/log/logfile.log
```

![2__fish___users_antoine_gocode_src_github_com_aybabtme_humanlog__fish_](https://cloud.githubusercontent.com/assets/1189716/4328545/f2330bb4-3f86-11e4-8242-4f49f6ae9efc.png)

# Contributing

How to help:

* __support more log formats__: by submitting `human.Handler` implementations.
* __live querying__: add support for filtering in log output in real time.
* __charting__: some key-values have semantics that could be charted in real time. For
instance, durations, frequency of numeric values, etc. See the [l2met][] project.

# Usage

```
NAME:
   humanlog - reads structured logs from stdin, makes them pretty on stdout!

USAGE:
   humanlog [global options] command [command options] [arguments...]

VERSION:
   0.4.0

AUTHOR:
  Antoine Grondin - <antoine@digitalocean.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --skip '--skip option --skip option'   keys to skip when parsing a log entry
   --keep '--keep option --keep option'   keys to keep when parsing a log entry
   --sort-longest       sort by longest key after having sorted lexicographically
   --skip-unchanged        skip keys that have the same value than the previous entry
   --truncate           truncates values that are longer than --truncate-length
   --truncate-length '15'     truncate values that are longer than this length
   --help, -h           show help
   --version, -v        print the version
```
[l2met]: https://github.com/ryandotsmith/l2met
