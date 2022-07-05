# How To Make And Run Integration Tests For lazygit

Integration tests are located in `test/integration`. Each test will run a bash script to prepare a test repo, then replay a recorded lazygit session from within that repo, and then the resultant repo will be compared to an expected repo that was created upon the initial recording. Each integration test lives in its own directory, and the name of the directory becomes the name of the test. Within the directory must be the following files:

### `test.json`

An example of a `test.json` is:

```
{ "description": "stage a file and commit the change", "speed": 20 }
```

The `speed` key refers to the playback speed as a multiple of the original recording speed. So 20 means the test will run 20 times faster than the original recording speed. If a test fails for a given speed, it will drop the speed and re-test, until finally attempting the test at the original speed. If you omit the speed, it will default to 10.

### `setup.sh`

This is a bash script containing the instructions for creating the test repo from scratch. For example:

```
#!/bin/sh

cd $1

git init

git config user.email "CI@example.com"
git config user.name "CI"

echo test1 > myfile1
git add .
git commit -am "myfile1"
```

Be sure to:

- ensure that by the end of the test you've got at least one commit in the repo, as we've had issues in the past when that wasn't the case.
- set the git user email and name as above so that your own user details aren't included in the snapshot.

## Running tests

### From a TUI

You can run/record/sandbox tests via a TUI with the following command:

```
go run test/lazyintegration/main.go
```

This TUI makes much of the following documentation redundant, but feel free to read through anyway!

### From command line

To run all tests - assuming you're at the project root:

```
go test ./pkg/gui/
```

To run them in parallel

```
PARALLEL=true go test ./pkg/gui
```

To run a single test

```
go test ./pkg/gui -run /<test name>
# For example, to run the `tags` test:
go test ./pkg/gui -run /tags
```

To run a test at a certain speed

```
SPEED=2 go test ./pkg/gui -run /<test name>
```

To update a snapshot

```
MODE=updateSnapshot go test ./pkg/gui -run /<test name>
```

## Creating a new test

To create a new test:

1. Copy and paste an existing test directory and rename the new directory to whatever you want the test name to be. Update the test.json file's description to describe your test.
2. Update the `setup.sh` any way you like
3. If you want to have a config folder for just that test, create a `config` directory to contain a `config.yml` and optionally a `state.yml` file. Otherwise, the `test/default_test_config` directory will be used.
4. From the lazygit root directory, run:

```
MODE=record go test ./pkg/gui -run /<test name>
```

5. Feel free to re-attempt recording as many times as you like. In the absence of a proper testing framework, the more deliberate your keypresses, the better!
6. Once satisfied with the recording, stage all the newly created files: `test.json`, `setup.sh`, `recording.json` and the `expected` directory that contains a copy of the repo you created.

The resulting directory will look like:

```
actual/ (the resulting repo(s) after running the test, ignored by git)
expected/ (the 'snapshot' repo(s))
config/ (need not be present)
test.json
setup.sh
recording.json
```

## Sandboxing

The integration tests serve a secondary purpose of providing a setup for easy sandboxing. If you want to run a test in sandbox mode (meaning the session won't be recorded and we won't create/update snapshots), go:

```
MODE=sandbox go test ./pkg/gui -run /<test name>
```

## Feedback

If you think this process can be improved, let me know! It shouldn't be too hard to change things.
