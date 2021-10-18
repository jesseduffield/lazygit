# Contributing

♥ We love pull requests from everyone !

When contributing to this repository, please first discuss the change you wish
to make via issue, email, or any other method with the owners of this repository
before making a change.

## So all code changes happen through Pull Requests

Pull requests are the best way to propose changes to the codebase. We actively
welcome your pull requests:

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.
3. If you've added code that need documentation, update the documentation.
4. Make sure your code follows the [effective go](https://golang.org/doc/effective_go.html) guidelines as much as possible.
5. Be sure to test your modifications.
6. Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
7. Issue that pull request!

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

## Updating Gocui

Sometimes you will need to make a change in the gocui fork (https://github.com/jesseduffield/gocui). Gocui is the package responsible for rending windows and handling user input. Here's the typical process to follow:

1. Make the changes in gocui inside the vendor directory so it's easy to test against lazygit
2. Copy the changes over to the actual gocui repo (clone it if you haven't already, and use the `awesome` branch, not `master`)
3. Raise a PR on the gocui repo with your changes
4. After that PR is merged, make a PR in lazygit bumping the gocui version. You can bump the version by running the following at the lazygit repo root:

```sh
./bump_gocui.sh
```

5. Raise a PR in lazygit with those changes
