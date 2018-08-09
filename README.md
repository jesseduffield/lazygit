# lazygit [![Go Report Card](https://goreportcard.com/badge/github.com/jesseduffield/lazygit)](https://goreportcard.com/report/github.com/jesseduffield/lazygit)

A simple terminal UI for git commands, written in Go with the [gocui](https://github.com/jroimartin/gocui "gocui") library.

Are YOU tired of typing every git command directly into the terminal, but you're too stubborn to use Sourcetree because you'll never forgive Atlassian for making Jira? This is the app for you!

[Tutorial](https://www.youtube.com/watch?v=VDXvbHZYeKY)

![Gif](https://image.ibb.co/mmeXho/optimisedgif.gif)

## Installation

### Homebrew
```sh
brew tap jesseduffield/lazygit
brew install lazygit
```

### Ubuntu
Packages for Ubuntu 16.04, 18.04 and 18.10 are available via Launchpad PPA.

They are built daily, straight from master branch.

```sh
sudo add-apt-repository ppa:lazygit-team/daily
sudo apt-get update
sudo apt-get install lazygit
```

### Arch Linux
Packages for Arch Linux are available via AUR (Arch User Repository).

There are two packages. The stable which is built with the last releases and the git version which takes the last commit and build based on.

Stable:
https://aur.archlinux.org/packages/lazygit/

Development:
https://aur.archlinux.org/packages/lazygit-git/

Instruction of how to install AUR content can be found here:
https://wiki.archlinux.org/index.php/Arch_User_Repository

### Binary Release (Windows/Linux/OSX)
You can download a binary release [here](https://github.com/jesseduffield/lazygit/releases)

### Go
```sh
go get github.com/jesseduffield/lazygit
```

Please note:
If you get an error claiming that lazygit cannot be found or is not defined, you may need to add `~/go/bin` to your $PATH (MacOS/Linux), or `%HOME%\go\bin` (Windows). Not to be mistaked for `C:\Go\bin` (which is for Go's own binaries, not apps like Lazygit)


## Usage
Call `lazygit` in your terminal inside a git repository.
If you want, you can also add an alias for this with `echo "alias lg='lazygit'" >> ~/.zshrc` (or whichever rc file you're using).
Basic tutorial [Here](https://www.youtube.com/watch?v=VDXvbHZYeKY)

[Keybindings](https://github.com/jesseduffield/lazygit/blob/master/docs/Keybindings.md)

## Cool features
- Adding files easily
- Resolving merge conflicts
- Easily check out recent branches
- Scroll through logs/diffs of branches/commits/stash
- Quick pushing/pulling
- Squash down and rename commits

### Resolving merge conflicts
![Gif](https://image.ibb.co/iyxUTT/shortermerging.gif)

### Viewing commit diffs
![Viewing Commit Diffs](https://image.ibb.co/gPD02o/capture.png)

## Milestones
- [x] Easy Installation (homebrew, release binaries)
- [ ] Configurable Keybindings
- [ ] Configurable Color Themes
- [ ] Spawning Subprocesses (help needed - have a look at https://github.com/jesseduffield/lazygit/pull/18)
- [ ] Maintainability
- [ ] Performance
- [ ] i18n

## Contributing
We love your input! Please check out the [contributing guide](CONTRIBUTING.md).

## Work in progress
This is still a work in progress so there's still bugs to iron out and as this is my first project in Go the code could no doubt use an increase in quality, but I'll be improving on it whenever I find the time. If you have any feedback feel free to [raise an issue](https://github.com/jesseduffield/lazygit/issues)/[submit a PR](https://github.com/jesseduffield/lazygit/pulls).
