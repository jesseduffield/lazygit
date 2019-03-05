# lazygit [![CircleCI](https://circleci.com/gh/jesseduffield/lazygit.svg?style=svg)](https://circleci.com/gh/jesseduffield/lazygit) [![codecov](https://codecov.io/gh/jesseduffield/lazygit/branch/master/graph/badge.svg)](https://codecov.io/gh/jesseduffield/lazygit) [![Go Report Card](https://goreportcard.com/badge/github.com/jesseduffield/lazygit)](https://goreportcard.com/report/github.com/jesseduffield/lazygit) [![GolangCI](https://golangci.com/badges/github.com/jesseduffield/lazygit.svg)](https://golangci.com) [![GoDoc](https://godoc.org/github.com/jesseduffield/lazygit?status.svg)](http://godoc.org/github.com/jesseduffield/lazygit) [![GitHub tag](https://img.shields.io/github/tag/jesseduffield/lazygit.svg)]()

A simple terminal UI for git commands, written in Go with the [gocui](https://github.com/jroimartin/gocui 'gocui') library.

Are YOU tired of typing every git command directly into the terminal, but you're
too stubborn to use Sourcetree because you'll never forgive Atlassian for making
Jira? This is the app for you!

![Gif](/docs/resources/lazygit-example.gif)

- [Installation](https://github.com/jesseduffield/lazygit#installation)
- [Usage](https://github.com/jesseduffield/lazygit#usage),
  [Keybindings](https://github.com/jesseduffield/lazygit/blob/master/docs/Keybindings_en.md)
- [Cool Features](https://github.com/jesseduffield/lazygit#cool-features)
- [Contributing](https://github.com/jesseduffield/lazygit#contributing)
- [Video Tutorial](https://youtu.be/VDXvbHZYeKY)
- [Twitch Stream](https://www.twitch.tv/jesseduffield)

## Installation

### Homebrew

```sh
brew tap jesseduffield/lazygit
brew install lazygit
```

### Ubuntu

Packages for Ubuntu 16.04, 18.04 and 18.10 are available via [Launchpad PPA](https://launchpad.net/~lazygit-team).

**Release builds**

Built from git tags. Supposed to be more stable.

```sh
sudo add-apt-repository ppa:lazygit-team/release
sudo apt-get update
sudo apt-get install lazygit
```

**Daily builds**

Built from master branch once in 24 hours (or more sometimes).

```sh
sudo add-apt-repository ppa:lazygit-team/daily
sudo apt-get update
sudo apt-get install lazygit
```

### Void Linux

Packages for Void Linux are available in the distro repo

They follow upstream latest releases

```sh
sudo xbps-install -S lazygit
```

### Arch Linux

Packages for Arch Linux are available via AUR (Arch User Repository).

There are two packages. The stable one which is built with the latest release
and the git version which builds from the most recent commit.

- Stable: https://aur.archlinux.org/packages/lazygit/
- Development: https://aur.archlinux.org/packages/lazygit-git/

Instruction of how to install AUR content can be found here:
https://wiki.archlinux.org/index.php/Arch_User_Repository

### Conda

Released versions are available for different platforms, see https://anaconda.org/conda-forge/lazygit

```sh
conda install -c conda-forge lazygit
```

### Binary Release (Windows/Linux/OSX)

You can download a binary release [here](https://github.com/jesseduffield/lazygit/releases).

### Go

```sh
go get github.com/jesseduffield/lazygit
```

Please note:
If you get an error claiming that lazygit cannot be found or is not defined, you
may need to add `~/go/bin` to your \$PATH (MacOS/Linux), or `%HOME%\go\bin`
(Windows). Not to be mistaked for `C:\Go\bin` (which is for Go's own binaries,
not apps like Lazygit).

## Usage

Call `lazygit` in your terminal inside a git repository. If you want, you can
also add an alias for this with `echo "alias lg='lazygit'" >> ~/.zshrc` (or
whichever rc file you're using).

- Basic video tutorial [here](https://youtu.be/VDXvbHZYeKY).
- List of keybindings
  [here](/docs/Keybindings_en.md).

## Cool features

- Adding files easily
- Resolving merge conflicts
- Easily check out recent branches
- Scroll through logs/diffs of branches/commits/stash
- Quick pushing/pulling
- Squash down and rename commits

### Resolving merge conflicts

![Gif](/docs/resources/resolving-merge-conflicts.gif)

### Interactive Rebasing

![Interactive Rebasing](/docs/resources/interactive-rebase.png)

## Contributing

We love your input! Please check out the [contributing guide](CONTRIBUTING.md).
For contributor discussion about things not better discussed here in the repo, join the slack channel

[![Slack](/docs/resources/slack_rgb.png)](https://join.slack.com/t/lazygit/shared_invite/enQtNDE3MjIwNTYyMDA0LTM3Yjk3NzdiYzhhNTA1YjM4Y2M4MWNmNDBkOTI0YTE4YjQ1ZmI2YWRhZTgwNjg2YzhhYjg3NDBlMmQyMTI5N2M)

## Donate

If you would like to support the development of lazygit, please donate

[![Donate](https://d1iczxrky3cnb2.cloudfront.net/button-medium-blue.png)](https://donorbox.org/lazygit)

## Work in progress

This is still a work in progress so there's still bugs to iron out and as this
is my first project in Go the code could no doubt use an increase in quality,
but I'll be improving on it whenever I find the time. If you have any feedback
feel free to [raise an issue](https://github.com/jesseduffield/lazygit/issues)/[submit a PR](https://github.com/jesseduffield/lazygit/pulls).

## Social

If you want to see what I (Jesse) am up to in terms of development, follow me on
[twitter](https://twitter.com/DuffieldJesse) or watch me program on
[twitch](https://www.twitch.tv/jesseduffield).

## Alternatives

If you find that lazygit doesn't quite satisfy your requirements, these may be a better fit:

- [tig](https://github.com/jonas/tig)
