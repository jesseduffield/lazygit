# lazygit [![CircleCI](https://circleci.com/gh/jesseduffield/lazygit.svg?style=svg)](https://circleci.com/gh/jesseduffield/lazygit) [![codecov](https://codecov.io/gh/jesseduffield/lazygit/branch/master/graph/badge.svg)](https://codecov.io/gh/jesseduffield/lazygit) [![Go Report Card](https://goreportcard.com/badge/github.com/jesseduffield/lazygit)](https://goreportcard.com/report/github.com/jesseduffield/lazygit) [![GolangCI](https://golangci.com/badges/github.com/jesseduffield/lazygit.svg)](https://golangci.com) [![GoDoc](https://godoc.org/github.com/jesseduffield/lazygit?status.svg)](http://godoc.org/github.com/jesseduffield/lazygit) [![GitHub tag](https://img.shields.io/github/tag/jesseduffield/lazygit.svg)]()

A simple terminal UI for git commands, written in Go with the [gocui](https://github.com/jroimartin/gocui 'gocui') library.

Rant time: You've heard it before, git is _powerful_, but what good is that power when everything is so damn hard to do? Interactive rebasing requires you to edit a goddamn TODO file in your editor? *Are you kidding me?* To stage part of a file you need to use a command line program stepping through each hunk and if a hunk can't be split down any further but contains code you don't want to stage, bad luck? *Are you KIDDING me?!* Sometimes you get asked to stash your changes when switching branches only to realise that after you switch and unstash that there weren't even any conflicts and it would have been fine to just checkout the branch directly? *YOU HAVE GOT TO BE KIDDING ME!* 

If you're a mere mortal like me and you're tired of hearing how powerful git is when in your daily life it's a powerful pain in your ass, lazygit might be for you.

![Gif](/docs/resources/lazygit-example.gif)

- [Installation](https://github.com/jesseduffield/lazygit#installation)
- [Usage](https://github.com/jesseduffield/lazygit#usage),
  [Keybindings](/docs/keybindings)
- [Cool Features](https://github.com/jesseduffield/lazygit#cool-features)
- [Contributing](https://github.com/jesseduffield/lazygit#contributing)
- [Video Tutorial](https://youtu.be/VDXvbHZYeKY)
- [Rebase Magic Video Tutorial](https://youtu.be/4XaToVut_hs)
- [Twitch Stream](https://www.twitch.tv/jesseduffield)

[<img src="https://i.imgur.com/sVEktDn.png">](https://youtu.be/CPLdltN7wgE)

Github Sponsors is matching all donations dollar-for-dollar for 12 months so if you're feeling generous consider [sponsoring me](https://github.com/sponsors/jesseduffield)

## Installation

### Homebrew
Normally the lazygit formula can be found in the Homebrew core but we suggest you tap our formula to get the frequently updated one. It works with Linux, too.

Tap:
```
brew install jesseduffield/lazygit/lazygit
```

Core:

```
brew install lazygit
```

### MacPorts
Latest version built from github releases.
Tap:
```
sudo port install lazygit
```

### Ubuntu

Packages for Ubuntu are available via [Launchpad PPA](https://launchpad.net/~lazygit-team).

```sh
sudo add-apt-repository ppa:lazygit-team/release
sudo apt-get update
sudo apt-get install lazygit
```

### Void Linux

Packages for Void Linux are available in the distro repo

They follow upstream latest releases

```sh
sudo xbps-install -S lazygit
```
### Scoop (Windows)

You can install `lazygit` using [scoop](https://scoop.sh/). It's in the `extras` bucket:

```sh
# Add the extras bucket
scoop bucket add extras

# Install lazygit
scoop install lazygit
```

### Arch Linux

Packages for Arch Linux are available via AUR (Arch User Repository).

There are two packages. The stable one which is built with the latest release
and the git version which builds from the most recent commit.

- Stable: https://aur.archlinux.org/packages/lazygit/
- Development: https://aur.archlinux.org/packages/lazygit-git/

Instruction of how to install AUR content can be found here:
https://wiki.archlinux.org/index.php/Arch_User_Repository

### Fedora and CentOS 7

Packages for Fedora and CentOS 7 are available via [Copr](https://copr.fedorainfracloud.org/coprs/atim/lazygit/) (Cool Other Package Repo).

```sh
sudo dnf copr enable atim/lazygit -y
sudo dnf install lazygit
```

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
- Rebase Magic tutorial [here](https://youtu.be/4XaToVut_hs)
- List of keybindings
  [here](/docs/keybindings).
  
## Changing Directory On Exit

If you change repos in lazygit and want your shell to change directory into that repo on exiting lazygit, add this to your `~/.zshrc` (or other rc file):
```
lg()
{
    export LAZYGIT_NEW_DIR_FILE=~/.lazygit/newdir

    lazygit "$@"

    if [ -f $LAZYGIT_NEW_DIR_FILE ]; then
            cd "$(cat $LAZYGIT_NEW_DIR_FILE)"
            rm -f $LAZYGIT_NEW_DIR_FILE > /dev/null
    fi
}
```
Then `source ~/.zshrc` and from now on when you call `lg` and exit you'll switch directories to whatever you were in inside lazyigt. To override this behaviour you can exit using `shift+Q` rather than just `q`.

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

If you would like to support the development of lazygit, consider [sponsoring me](https://github.com/sponsors/jesseduffield) (github is matching all donations dollar-for-dollar for 12 months)

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
