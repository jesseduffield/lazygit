# Contributing to bluemonday

Third-party patches are essential for keeping bluemonday secure and offering the features developers want. However there are a few guidelines that we need contributors to follow so that we can maintain the quality of work that developers who use bluemonday expect.

## Getting Started

* Make sure you have a [Github account](https://github.com/signup/free)

## Guidelines

1. Do not vendor dependencies. As a security package, were we to vendor dependencies the projects that then vendor bluemonday may not receive the latest security updates to the dependencies. By not vendoring dependencies the project that implements bluemonday will vendor the latest version of any dependent packages. Vendoring is a project problem, not a package problem. bluemonday will be tested against the latest version of dependencies periodically and during any PR/merge.
2. I do not care about spelling mistakes or whitespace and I do not believe that you should either. PRs therefore must be functional in their nature or be substantial and impactful if documentation or examples.

## Submitting an Issue

* Submit a ticket for your issue, assuming one does not already exist
* Clearly describe the issue including the steps to reproduce (with sample input and output) if it is a bug

If you are reporting a security flaw, you may expect that we will provide the code to fix it for you. Otherwise you may want to submit a pull request to ensure the resolution is applied sooner rather than later:

* Fork the repository on Github
* Issue a pull request containing code to resolve the issue

## Submitting a Pull Request

* Submit a ticket for your issue, assuming one does not already exist
* Describe the reason for the pull request and if applicable show some example inputs and outputs to demonstrate what the patch does
* Fork the repository on Github
* Before submitting the pull request you should
  1. Include tests for your patch, 1 test should encapsulate the entire patch and should refer to the Github issue
  1. If you have added new exposed/public functionality, you should ensure it is documented appropriately
  1. If you have added new exposed/public functionality, you should consider demonstrating how to use it within one of the helpers or shipped policies if appropriate or within a test if modifying a helper or policy is not appropriate
  1. Run all of the tests `go test -v ./...` or `make test` and ensure all tests pass
  1. Run gofmt `gofmt -w ./$*` or `make fmt`
  1. Run vet `go tool vet *.go` or `make vet` and resolve any issues
  1. Install golint using `go get -u github.com/golang/lint/golint` and run vet `golint *.go` or `make lint` and resolve every warning
* When submitting the pull request you should
  1. Note the issue(s) it resolves, i.e. `Closes #6` in the pull request comment to close issue #6 when the pull request is accepted

Once you have submitted a pull request, we *may* merge it without changes. If we have any comments or feedback, or need you to make changes to your pull request we will update the Github pull request or the associated issue. We expect responses from you within two weeks, and we may close the pull request is there is no activity.

### Contributor Licence Agreement

We haven't gone for the formal "Sign a Contributor Licence Agreement" thing that projects like [puppet](https://cla.puppetlabs.com/), [Mojito](https://developer.yahoo.com/cocktails/mojito/cla/) and companies like [Google](http://code.google.com/legal/individual-cla-v1.0.html) are using.

But we do need to know that we can accept and merge your contributions, so for now the act of contributing a pull request should be considered equivalent to agreeing to a contributor licence agreement, specifically:

You accept that the act of submitting code to the bluemonday project is to grant a copyright licence to the project that is perpetual, worldwide, non-exclusive, no-charge, royalty free and irrevocable.

You accept that all who comply with the licence of the project (BSD 3-clause) are permitted to use your contributions to the project.

You accept, and by submitting code do declare, that you have the legal right to grant such a licence to the project and that each of the contributions is your own original creation.
