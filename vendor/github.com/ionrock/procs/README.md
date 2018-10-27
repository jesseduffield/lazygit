# Procs

[![](https://travis-ci.org/ionrock/procs.svg?branch=master)](https://travis-ci.org/ionrock/procs)
[![Go Report Card](https://goreportcard.com/badge/github.com/ionrock/procs)](https://goreportcard.com/report/github.com/ionrock/procs)
[![GoDoc](https://godoc.org/github.com/ionrock/procs?status.svg)](https://godoc.org/github.com/ionrock/procs)

Procs is a library to make working with command line applications a
little nicer.

The primary use case is when you have to use a command line client in
place of an API. Often times you want to do things like output stdout
within your own logs or ensure that every time the command is called,
there are a standard set of flags that are used.

## Basic Usage

The majority of this functionality is intended to be included the
procs.Process.

### Defining a Command

A command can be defined by a string rather than a []string. Normally,
this also implies that the library will run the command in a shell,
exposing a potential man in the middle attack. Rather than using a
shell, procs [lexically
parses](https://github.com/flynn-archive/go-shlex) the command for the
different arguments. It also allows for pipes in order to string
commands together.

```go
p := procs.NewProcess("kubectl get events | grep dev")
```

You can also define a new `Process` by passing in predefined commands.

```go
cmds := []*exec.Cmd{
	exec.Command("kubectl", "get", "events"),
	exec.Command("grep", "dev"),
}

p := procs.Process{Cmds: cmds}
```

### Output Handling

One use case that is cumbersome is using the piped output from a
command. For example, lets say we wanted to start a couple commands
and have each command have its own prefix in stdout, while still
capturing the output of the command as-is.

```go
p := procs.NewProcess("cmd1")
p.OutputHandler = func(line string) string {
	fmt.Printf("cmd1 | %s\n")
	return line
}
out, _ := p.Run()
fmt.Println(out)
```

Whatever is returned from the `OutputHandler` will be in the buffered
output. In this way you can choose to filter or skip output buffering
completely.

You can also define a `ErrHandler` using the same signature to get the
same filtering for stderr.

### Environment Variables

Rather than use the `exec.Cmd` `[]string` environment variables, a
`procs.Process` uses a `map[string]string` for environment variables.

```go
p := procs.NewProcess("echo $FOO")
p.Env = map[string]string{"FOO": "foo"}
```

Also, environment variables defined by the `Process.Env` can be
expanded automatically using the `os.Expand` semantics and the
provided environment.

There is a `ParseEnv` function that can help to merge the parent
processes' environment with any new values.

```go
env := ParseEnv(os.Environ())
env["USER"] = "foo"
```

Finally, if you are building commands manually, the `Env` function can
take a `map[string]string` and convert it to a `[]string` for use with
an `exec.Cmd`. The `Env` function also accepts a `useEnv` bool to help
include the parent process environment.

```go
cmd := exec.Command("knife", "cookbook", "show", cb)
cmd.Env = Env(map[string]string{"USER": "knife-user"}, true)
```

## Example Applications

Take a look in the [`cmd`](./cmd/) dir for some simple applications
that use the library. You can also `make all` to build them. The
examples below assume you've built them locally.

### Prelog

The `prelog` command allows running a command and prefixing the output
with a value.

```bash
$ ./prelog -prefix foo -- echo 'hello world!'
Running the command
foo | hello world!
Accessing the output without a prefix.
hello world!
Running the command with Start / Wait
foo | hello world!
```

### Cmdtmpl

The `cmdtmpl` command uses the `procs.Builder` to create a command
based on some paramters. It will take a `data.yml` file and
`template.yml` file to create a command.

```bash
$ cat example/data.json
{
  "source": "https://my.example.org",
  "user": "foo",
  "model": "widget",
  "action": "create",
  "args": "-f new -i improved"
}
$ cat example/template.json
[
  "mysvc ${model} ${action} ${args}",
  "--endpoint ${source}",
  "--username ${user}"
]
$ ./cmdtmpl -data example/data.json -template example/template.json
Command: mysvc foo widget create -f new -i imporoved --endpoint https://my.example.org --username foo
$ ./cmdtmpl -data example/data.json -template example/template.json -field user=bar
Command: mysvc foo widget create -f new -i imporoved --endpoint https://my.example.org --username bar
```

### Procmon

The `procmon` command acts like
[foreman](https://github.com/ddollar/foreman) with the difference
being it uses a JSON file with key value pairs instead of a
Procfile. This example uses the `procs.Manager` to manage a set of
`procs.Processes`.

```bash
$ cat example/procfile.json
{
  "web": "python -m SimpleHTTPServer"
}
$ ./procmon -procfile example/procfile.json
web | Starting web with python -m SimpleHTTPServer
```

You can then access http://localhost:8000 to see the logs. You can
also kill the child process and see `procmon` recognizing it has
exited and exit itself.
