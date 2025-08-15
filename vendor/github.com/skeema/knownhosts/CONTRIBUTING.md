# Contributing to skeema/knownhosts

Thank you for your interest in contributing! This document provides guidelines for submitting pull requests.

### Link to an issue

Before starting the pull request process, initial discussion should take place on a GitHub issue first. For bug reports, the issue should track the open bug and confirm it is reproducible. For feature requests, the issue should cover why the feature is necessary.

In the issue comments, discuss your suggested approach for a fix/implementation, and please wait to get feedback before opening a pull request.

### Test coverage

In general, please provide reasonably thorough test coverage. Whenever possible, your PR should aim to match or improve the overall test coverage percentage of the package. You can run tests and check coverage locally using `go test -cover`. We also have CI automation in GitHub Actions which will comment on each pull request with a coverage percentage.

That said, it is fine to submit an initial draft / work-in-progress PR without coverage, if you are waiting on implementation feedback before writing the tests.

We intentionally avoid hard-coding SSH keys or known_hosts files into the test logic. Instead, the tests generate new keys and then use them to generate a known_hosts file, which is then cached/reused for that overall test run, in order to keep performance reasonable.

### Documentation

Exported types require doc comments. The linter CI step will catch this if missing.

### Backwards compatibility

Because this package is imported by [nearly 7000 repos on GitHub](https://github.com/skeema/knownhosts/network/dependents), we must be very strict about backwards compatibility of exported symbols and function signatures.

Backwards compatibility can be very tricky in some situations. In this case, a maintainer may need to add additional commits to your branch to adjust the approach. Please do not take offense if this occurs; it is sometimes simply faster to implement a refactor on our end directly. When the PR/branch is merged, a merge commit will be used, to ensure your commits appear as-is in the repo history and are still properly credited to you.

### Avoid rewriting core x/crypto/ssh/knownhosts logic

skeema/knownhosts is intended to be a relatively thin *wrapper* around x/crypto/ssh/knownhosts, without duplicating or re-implementing the core known_hosts file parsing and host key handling logic. Importers of this package should be confident that it can be used as a nearly-drop-in replacement for x/crypto/ssh/knownhosts without introducing substantial risk, security flaws, parser differentials, or unexpected behavior changes.

To solve shortcomings in x/crypto/ssh/knownhosts, we try to come up with workarounds that still utilize x/crypto/ssh/knownhosts functionality whenever possible.

Some bugs in x/crypto/ssh/knownhosts do require re-reading the known_hosts file here to solve, but we make that *optional* by offering separate constructors/types with and without that behavior.

