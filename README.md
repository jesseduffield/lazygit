# lazygit
A simple terminal UI for git commands, written in Go with the [gocui](https://www.google.com "gocui") library.

are YOU tired of typing every git command directly into the terminal, but you're too stubborn to use Sourcetree because you'll never forgive Atlassian for making Jira? This is the app for you!

![Gif](https://image.ibb.co/mmeXho/optimisedgif.gif)

## Installation
`go get github.com/jesseduffield/lazygit`

then just call `lazygit` in your terminal inside a git repository

You can add an alias for this with `echo "alias lg='lazygit' > ~/.zshrc` (or whichever rc file you're using)

## Usage
All commands are listed at the bottom of the application

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

## Work in progress
This is still a work in progress so there's still bugs to iron out and as this is my first project in Go the code could no doubt use an increase in quality, but I'll be improving on it whenever I find the time. If you have any feedback feel free to raise an issue/submit a PR.
