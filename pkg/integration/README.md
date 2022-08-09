# Integration Tests

There's a lot happening in this package so it's worth a proper explanation.

This package is for integration testing: that is, actually running a real lazygit session and having a robot pretend to be a human user and then making assertions that everything works as expected.

There are three ways to invoke a test:

1. go run pkg/integration/runner/main.go commit/new_branch
2. go test pkg/integration/integration_test.go
3.
