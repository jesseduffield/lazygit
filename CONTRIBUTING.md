# Contributing

♥ We love pull requests from everyone !

When contributing to this repository, please first discuss the change you wish
to make via issue, email, or any other method with the owners of this repository
before making a change.

## PR walkthrough

[This video](https://www.youtube.com/watch?v=kNavnhzZHtk) walks through the process of adding a small feature to lazygit. If you have no idea where to start, watching that video is a good first step.

## All code changes happen through Pull Requests

Pull requests are the best way to propose changes to the codebase. We actively
welcome your pull requests:

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.
3. If you've added code that need documentation, update the documentation.
4. Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
5. Issue that pull request!

If you've never written Go in your life, then join the club! Lazygit was the maintainer's first Go program, and most contributors have never used Go before. Go is widely considered an easy-to-learn language, so if you're looking for an open source project to gain dev experience, you've come to the right place.

## Code of conduct

Please note by participating in this project, you agree to abide by the [code of conduct].

[code of conduct]: https://github.com/jesseduffield/lazygit/blob/master/CODE-OF-CONDUCT.md

## Any contributions you make will be under the MIT Software License

In short, when you submit code changes, your submissions are understood to be
under the same [MIT License](http://choosealicense.com/licenses/mit/) that
covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using Github's [issues](https://github.com/jesseduffield/lazygit/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new
issue](https://github.com/jesseduffield/lazygit/issues/new); it's that easy!

## Go

This project is written in Go. Go is an opinionated language with strict idioms, but some of those idioms are a little extreme. Some things we do differently:

1. There is no shame in using `self` as a receiver name in a struct method. In fact we encourage it
2. There is no shame in prefixing an interface with 'I' instead of suffixing with 'er' when there are several methods on the interface.
3. If a struct implements an interface, we make it explicit with something like:

```go
var _ MyInterface = &MyStruct{}
```

This makes the intent clearer and means that if we fail to satisfy the interface we'll get an error in the file that needs fixing.

### Code Formatting

To check code formatting [gofumpt](https://pkg.go.dev/mvdan.cc/gofumpt#section-readme) (which is a bit stricter than [gofmt](https://pkg.go.dev/cmd/gofmt)) is used.
VSCode will format the code correctly if you tell the Go extension to use `gofumpt` via your [`settings.json`](https://code.visualstudio.com/docs/getstarted/settings#_settingsjson)
by setting [`formatting.gofumpt`](https://github.com/golang/tools/blob/master/gopls/doc/settings.md#gofumpt-bool) to `true`:

```jsonc
// .vscode/settings.json
{
    "gopls": {
        "formatting.gofumpt": true
    }
}
```

## Internationalisation

Boy that's a hard word to spell. Anyway, lazygit is translated into several languages within the pkg/i18n package. If you need to render text to the user, you should add a new field to the TranslationSet struct in `pkg/i18n/english.go` and add the actual content within the `EnglishTranslationSet()` method in the same file. Although it is appreciated if you translate the text into other languages, it's not expected of you (google translate will likely do a bad job anyway!).

## Debugging

The easiest way to debug lazygit is to have two terminal tabs open at once: one for running lazygit (via `go run main.go -debug` in the project root) and one for viewing lazygit's logs (which can be done via `go run main.go --logs` or just `lazygit --logs`).

From most places in the codebase you have access to a logger e.g. `gui.Log.Warn("blah")`.

If you find that the existing logs are too noisy, you can set the log level with e.g. `LOG_LEVEL=warn go run main.go -debug` and then only use `Warn` logs yourself.

If you need to log from code in the vendor directory (e.g. the `gocui` package), you won't have access to the logger, but you can easily add logging support by adding the following:
```go
func newLogger() *logrus.Entry {
	// REPLACE THE BELOW PATH WITH YOUR ACTUAL LOG PATH (YOU'LL SEE THIS PRINTED WHEN YOU RUN `lazygit --logs`
	logPath := "/Users/jesseduffield/Library/Application Support/jesseduffield/lazygit/development.log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("unable to log to file")
	}
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	logger.SetOutput(file)
	return logger.WithFields(logrus.Fields{})
}

var Log = newLogger()
...
Log.Warn("blah")
```

If you keep having to do some setup steps to reproduce an issue, read the Testing section below to see how to create an integration test by recording a lazygit session. It's pretty easy!

### VSCode debugger

If you want to trigger a debug session from VSCode, you can use the following snippet. Note that the `console` key is, at the time of writing, still an experimental feature.

```jsonc
// .vscode/launch.json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "debug lazygit",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "main.go",
      "args": [
        "--debug"
      ],
      "console": "externalTerminal" // <-- you need this to actually see the lazygit UI in a window while debugging
    }
  ]
}
```

## Testing

Lazygit has two kinds of tests: unit tests and integration tests. Unit tests go in files that end in `_test.go`, and are written in Go. Lazygit has its own integration test system where you can build a sandbox repo with a shell script, record yourself doing something, and commit the resulting repo snapshot. It's pretty damn cool! To learn more see [here](https://github.com/jesseduffield/lazygit/blob/master/docs/Integration_Tests.md)

## Updating Gocui

Sometimes you will need to make a change in the gocui fork (https://github.com/jesseduffield/gocui). Gocui is the package responsible for rendering windows and handling user input. Here's the typical process to follow:

1. Make the changes in gocui inside the vendor directory so it's easy to test against lazygit
2. Copy the changes over to the actual gocui repo (clone it if you haven't already, and use the `awesome` branch, not `master`)
3. Raise a PR on the gocui repo with your changes
4. After that PR is merged, make a PR in lazygit bumping the gocui version. You can bump the version by running the following at the lazygit repo root:

```sh
./scripts/bump_gocui.sh
```

5. Raise a PR in lazygit with those changes

## Improvements

If you can think of any way to improve these docs let us know.
