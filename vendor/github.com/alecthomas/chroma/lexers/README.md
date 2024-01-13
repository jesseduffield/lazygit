# Lexer tests

The tests in this directory feed a known input `testdata/<name>.actual` into the parser for `<name>` and check
that its output matches `<name>.exported`.

It is also possible to perform several tests on a same parser `<name>`, by placing know inputs `*.actual` into a
directory `testdata/<name>/`.

## Running the tests

Run the tests as normal:
```go
go test ./lexers
```

## Update existing tests
When you add a new test data file (`*.actual`), you need to regenerate all tests. That's how Chroma creates the `*.expected` test file based on the corresponding lexer.

To regenerate all tests, type in your terminal:

```go
RECORD=true go test ./lexers
```

This first sets the `RECORD` environment variable to `true`. Then it runs `go test` on the `./lexers` directory of the Chroma project.

(That environment variable tells Chroma it needs to output test data. After running `go test ./lexers` you can remove or reset that variable.)

### Windows users
Windows users will find that the `RECORD=true go test ./lexers` command fails in both the standard command prompt terminal and in PowerShell.

Instead we have to perform both steps separately:

- Set the `RECORD` environment variable to `true`.
	+ In the regular command prompt window, the `set` command sets an environment variable for the current session: `set RECORD=true`. See [this page](https://superuser.com/questions/212150/how-to-set-env-variable-in-windows-cmd-line) for more.
	+ In PowerShell, you can use the `$env:RECORD = 'true'` command for that. See [this article](https://mcpmag.com/articles/2019/03/28/environment-variables-in-powershell.aspx) for more.
	+ You can also make a persistent environment variable by hand in the Windows computer settings. See [this article](https://www.computerhope.com/issues/ch000549.htm) for how.
- When the environment variable is set, run `go tests ./lexers`.

Chroma will now regenerate the test files and print its results to the console window.
