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

The tests live in pkg/integration/tests. Each test has two important steps: the setup step and the run step.

### Setup step

In the setup step, we prepare a repo with shell commands, for example, creating a merge conflict that will need to be resolved upon opening lazygit. This is all done via the `shell` argument.

### Run step

The run step has four arguments passed in:

1. `shell`
2. `input`
3. `assert`
4. `keys`

`shell` we've already seen in the setup step. The reason it's passed into the run step is that we may want to emulate background events. For example, the user modifying a file outside of lazygit.

`input` is for driving the gui by pressing certain keys, selecting list items, etc.

`assert` is for asserting on the state of the lazygit session. When you call a method on `assert`, the assert struct will wait for the assertion to hold true and then continue (failing the test after a timeout). For this reason, assertions have two purposes: one is to ensure the test fails as soon as something unexpected happens, but another is to allow lazygit to process a keypress before you follow up with more keypresses. If you input a bunch of keypresses too quickly lazygit might get confused.

### Tips

Try to do as much setup work as possible in your setup step. For example, if all you're testing is that the user is able to resolve merge conflicts, create the merge conflicts in the setup step. On the other hand, if you're testing to see that lazygit can warn the user about merge conflicts after an attempted merge, it's fine to wait until the run step to actually create the conflicts. If the run step is focused on the thing you're trying to test, the test will run faster and its intent will be clearer.

Use assertions to ensure that lazygit has processed all your keybindings so far. For example, if you press 'n' on a branch to create a new branch, assert that the confirmation view is now focused.

If you find yourself doing something frequently in a test, consider making it a method in one of the helper arguments. For example, instead of calling `input.PressKey(keys.Universal.Confirm)` in 100 places, it's better to have a method `input.Confirm()`. This is not to say that everything should be made into a method on the input struct: just things that are particularly common in tests.

## Running tests

There are three ways to invoke a test:

1. go run cmd/integration_test/main.go cli [--slow or --sandbox] [testname or testpath...]
2. go run cmd/integration_test/main.go tui
3. go test pkg/integration/clients/go_test.go

The first, the test runner, is for directly running a test from the command line. If you pass no arguments, it runs all tests.
The second, the TUI, is for running tests from a terminal UI where it's easier to find a test and run it without having to copy it's name and paste it into the terminal. This is the easiest approach by far.
The third, the go-test command, intended only for use in CI, to be run along with the other `go test` tests. This runs the tests in headless mode so there's no visual output.

The name of a test is based on its path, so the name of the test at `pkg/integration/tests/commit/new_branch.go` is commit/new_branch. So to run it with our test runner you would run `go run cmd/integration_test/main.go cli commit/new_branch`.

You can pass the KEY_PRESS_DELAY env var to the test runner in order to set a delay in milliseconds between keypresses, which helps for watching a test at a realistic speed to understand what it's doing. Or you can pass the '--slow' flag which sets a pre-set 'slow' key delay. In the tui you can press 't' to run the test in slow mode.

### Running tests in VSCode

If you've opened an integration test file in your editor you can run that file by bringing up the command panel with `cmd+shift+p` and typing 'run task', then selecting the test task you want to run 

![image](https://user-images.githubusercontent.com/8456633/201500427-b86e129f-5f35-4d55-b7bd-fff5d8e4a04e.png)
![image](https://user-images.githubusercontent.com/8456633/201500431-903deb8c-c210-4054-8514-ab7088c7a839.png)
The test will run in a VSCode terminal:
![image](https://user-images.githubusercontent.com/8456633/201500446-b87abf11-9653-438f-8a9a-e0bf8abdb7ee.png)

### Snapshots

At the moment (this is subject to change) each test has a snapshot repo created after running for the first time. These snapshots live in `test/integration_new`, in folders named 'expected' (alongside the 'actual' folders which contain the resulting repo from the last test run). Whenever you run a test, the resultant repo will be compared against the snapshot repo and if they're different, you'll be asked whether you want to update the snapshot. If you want to update a snapshot without being prompted you can pass MODE=update to the test runner.

### Sandbox mode

Say you want to do a manual test of how lazygit handles merge-conflicts, but you can't be bothered actually finding a way to create merge conflicts in a repo. To make your life easier, you can simply run a merge-conflicts test in sandbox mode, meaning the setup step is run for you, and then instead of the test driving the lazygit session, you're allowed to drive it yourself.

To run a test in sandbox mode you can press 's' on a test in the test TUI or in the test runner pass MODE=sandbox or the --sandbox argument.

## Migration process

At the time of writing, most tests are created under an old approach, where you would record yourself in a lazygit session and then the test would replay the keybindings with the same timestamps. This old approach is great for writing tests quickly, but is much harder to maintain. It has to rely entirely on snapshots to determining if a test passes or fails, and can't do assertions along the way. It's also harder to grok what's the intention behind certain actions that take place within the test (e.g. was the recorder intentionally switching to another panel or was that just a misclick?).

At the moment, all the deprecated test code lives in pkg/integration/deprecated. Hopefully in the very near future we migrate everything across so that we don't need to maintain two systems.

We should never write any new tests under the old method, and if a given test breaks because of new functionality, it's best to simply rewrite it under the new approach. If you want to run a test for the sake of watching what it does so that you can transcribe it into the new approach, you can run:

```
go run pkg/integration/deprecated/cmd/tui/main.go
```

The tests in the old format live in test/integration. In the old format, test definitions are co-located with the snapshots. The setup step is done in a `setup.sh` shell script and the `recording.json` file contains the recorded keypresses to be replayed during the test.

If you have rewritten an integration test under the new pattern, be sure to delete the old integration test directory.
