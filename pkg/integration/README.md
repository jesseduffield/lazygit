# Integration Tests

The pkg/integration package is for integration testing: that is, actually running a real lazygit session and having a robot pretend to be a human user and then making assertions that everything works as expected.

TL;DR: integration tests live in pkg/integration/tests. Run integration tests with:

```sh
go run cmd/integration_test/main.go tui
```

or

```sh
go run cmd/integration_test/main.go cli [--slow or --sandbox] [testname or testpath...]
```

## Writing tests

The tests live in pkg/integration/tests. Each test is registered in `pkg/integration/tests/test_list.go` which is an auto-generated file. You can re-generate that file by running `go generate ./...` at the root of the Lazygit repo.

Each test has two important steps: the setup step and the run step.

### Setup step

In the setup step, we prepare a repo with shell commands, for example, creating a merge conflict that will need to be resolved upon opening lazygit. This is all done via the `shell` argument.

### Run step

The run step has two arguments passed in:

1. `t` (the test driver)
2. `keys`

`t` is for driving the gui by pressing certain keys, selecting list items, etc.
`keys` is for use when getting the test to press a particular key e.g. `t.Views().Commits().Focus().PressKey(keys.Universal.Confirm)`

## Running tests

There are three ways to invoke a test:

1. go run cmd/integration_test/main.go cli [--slow or --sandbox] [testname or testpath...]
2. go run cmd/integration_test/main.go tui
3. go test pkg/integration/clients/*.go

The first, the test runner, is for directly running a test from the command line. If you pass no arguments, it runs all tests.
The second, the TUI, is for running tests from a terminal UI where it's easier to find a test and run it without having to copy it's name and paste it into the terminal. This is the easiest approach by far.
The third, the go-test command, intended only for use in CI, to be run along with the other `go test` tests. This runs the tests in headless mode so there's no visual output.

The name of a test is based on its path, so the name of the test at `pkg/integration/tests/commit/new_branch.go` is commit/new_branch. So to run it with our test runner you would run `go run cmd/integration_test/main.go cli commit/new_branch`.

You can pass the KEY_PRESS_DELAY env var to the test runner in order to set a delay in milliseconds between keypresses, which helps for watching a test at a realistic speed to understand what it's doing. Or you can pass the '--slow' flag which sets a pre-set 'slow' key delay. In the tui you can press 't' to run the test in slow mode.

The resultant repo will be stored in `test/results`, so if you're not sure what went wrong you can go there and inspect the repo.

### Running tests in VSCode

If you've opened an integration test file in your editor you can run that file by bringing up the command panel with `cmd+shift+p` and typing 'run task', then selecting the test task you want to run

![image](https://user-images.githubusercontent.com/8456633/201500427-b86e129f-5f35-4d55-b7bd-fff5d8e4a04e.png)
![image](https://user-images.githubusercontent.com/8456633/201500431-903deb8c-c210-4054-8514-ab7088c7a839.png)
The test will run in a VSCode terminal:
![image](https://user-images.githubusercontent.com/8456633/201500446-b87abf11-9653-438f-8a9a-e0bf8abdb7ee.png)

### Sandbox mode

Say you want to do a manual test of how lazygit handles merge-conflicts, but you can't be bothered actually finding a way to create merge conflicts in a repo. To make your life easier, you can simply run a merge-conflicts test in sandbox mode, meaning the setup step is run for you, and then instead of the test driving the lazygit session, you're allowed to drive it yourself.

To run a test in sandbox mode you can press 's' on a test in the test TUI or in the test runner pass the --sandbox argument.

## Tips for writing tests

### Handle most setup in the `shell` part of the test

Try to do as much setup work as possible in your setup step. For example, if all you're testing is that the user is able to resolve merge conflicts, create the merge conflicts in the setup step. On the other hand, if you're testing to see that lazygit can warn the user about merge conflicts after an attempted merge, it's fine to wait until the run step to actually create the conflicts. If the run step is focused on the thing you're trying to test, the test will run faster and its intent will be clearer.

### Create helper functions for (very) frequently used test logic

If within a test directory you find several tests need to share some logic, you can create a file called `shared.go` in that directory to hold shared helper functions (see `pkg/integration/tests/filter_by_path/shared.go` for an example).

If you need to share test logic across test directories you can put helper functions in the `tests/shared` package. If you find yourself frequently doing the same thing from within a test across test directories, for example, responding a particular popup, consider adding a helper method to `pkg/integration/components/common.go`. If you look around the code in the `components` directory you may find another place that's sensible to put your helper function.

### Don't do too much in one test

If you're testing different pieces of functionality, it's better to test them in isolation using multiple short tests, compared to one larger longer test. Sometimes it's appropriate to have a longer test which tests how various different pieces interact, but err on the side of keeping things short.
