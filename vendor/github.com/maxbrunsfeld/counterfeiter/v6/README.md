# `counterfeiter` [![CircleCI](https://circleci.com/gh/maxbrunsfeld/counterfeiter.svg?style=svg)](https://circleci.com/gh/maxbrunsfeld/counterfeiter) [![Build status](https://ci.appveyor.com/api/projects/status/0j2v7pt06lp9yanm/branch/master?svg=true)](https://ci.appveyor.com/project/maxbrunsfeld/counterfeiter/branch/master)

When writing unit-tests for an object, it is often useful to have fake implementations
of the object's collaborators. In go, such fake implementations cannot be generated
automatically at runtime, and writing them by hand can be quite arduous.

`counterfeiter` allows you to simply generate test doubles for a given interface.

### Supported Versions Of `go`

`counterfeiter` follows the [support policy of `go` itself](https://golang.org/doc/devel/release.html#policy):

> Each major Go release is supported until there are two newer major releases. For example, Go 1.5 was supported until the Go 1.7 release, and Go 1.6 was supported until the Go 1.8 release. We fix critical problems, including [critical security problems](https://golang.org/security), in supported releases as needed by issuing minor revisions (for example, Go 1.6.1, Go 1.6.2, and so on).

If you are having problems with `counterfeiter` and are not using a supported version of go, please update to use a supported version of go before opening an issue.

### Using `counterfeiter`

⚠️ Please use [`go modules`](https://blog.golang.org/using-go-modules) when working with counterfeiter.

Typically, `counterfeiter` is used in `go generate` directives. It can be frustrating when you change your interface declaration and suddenly all of your generated code is suddenly out-of-date. The best practice here is to use the [`go generate` command](https://blog.golang.org/generate) to make it easier to keep your test doubles up to date.

#### Step 1 - Create `tools.go`

You can take a dependency on tools by creating a `tools.go` file, as described in [How can I track tool dependencies for a module?](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module). This ensures that everyone working with your module is using the same version of each tool you use.

```shell
$ cat tools/tools.go
```

```go
// +build tools

package tools

import (
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
)

// This file imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
```

#### Step 2a - Add `go:generate` Directives

You can add directives right next to your interface definitions (or not), in any `.go` file in your module.

```shell
$ cat myinterface.go
```

```go
package foo

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . MySpecialInterface

type MySpecialInterface interface {
	DoThings(string, uint64) (int, error)
}
```

```shell
$ go generate ./...
Writing `FakeMySpecialInterface` to `foofakes/fake_my_special_interface.go`... Done
```

#### Step 2b - Add `counterfeiter:generate` Directives

If you plan to have many directives in a single package, consider using this
option. You can add directives right next to your interface definitions
(or not), in any `.go` file in your module.

```shell
$ cat myinterface.go
```

```go
package foo

// You only need **one** of these per package!
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

// You will add lots of directives like these in the same package...
//counterfeiter:generate . MySpecialInterface
type MySpecialInterface interface {
	DoThings(string, uint64) (int, error)
}

// Like this...
//counterfeiter:generate . MyOtherInterface
type MyOtherInterface interface {
	DoOtherThings(string, uint64) (int, error)
}
```

```shell
$ go generate ./...
Writing `FakeMySpecialInterface` to `foofakes/fake_my_special_interface.go`... Done
Writing `FakeMyOtherInterface` to `foofakes/fake_my_other_interface.go`... Done
```

#### Step 3 - Run `go generate`

You can run `go generate` in the directory with your directive, or in the root of your module (to ensure you generate for all packages in your module):

```shell
$ go generate ./...
```

#### Invoking `counterfeiter` from the shell

You can use the following command to invoke `counterfeiter` from within a go module:

```shell
$ go run github.com/maxbrunsfeld/counterfeiter/v6

USAGE
	counterfeiter
		[-generate] [-o <output-path>] [-p] [--fake-name <fake-name>]
		[<source-path>] <interface> [-]
```

#### Installing `counterfeiter` to `$GOPATH/bin`

This is unnecessary if you're using the approach described above, but does allow you to invoke `counterfeiter` in your shell _outside_ of a module:

```shell
$ GO111MODULE=off go get -u github.com/maxbrunsfeld/counterfeiter
$ counterfeiter

USAGE
	counterfeiter
		[-generate] [-o <output-path>] [-p] [--fake-name <fake-name>]
		[<source-path>] <interface> [-]
```

### Generating Test Doubles

Given a path to a package and an interface name, you can generate a test double.

```shell
$ cat path/to/foo/file.go
```

```go
package foo

type MySpecialInterface interface {
		DoThings(string, uint64) (int, error)
}
```

```shell
$ go run github.com/maxbrunsfeld/counterfeiter/v6 path/to/foo MySpecialInterface
Wrote `FakeMySpecialInterface` to `path/to/foo/foofakes/fake_my_special_interface.go`
```

### Using Test Doubles In Your Tests

Instantiate fakes`:

```go
import "my-repo/path/to/foo/foofakes"

var fake = &foofakes.FakeMySpecialInterface{}
```

Fakes record the arguments they were called with:

```go
fake.DoThings("stuff", 5)

Expect(fake.DoThingsCallCount()).To(Equal(1))

str, num := fake.DoThingsArgsForCall(0)
Expect(str).To(Equal("stuff"))
Expect(num).To(Equal(uint64(5)))
```

You can stub their return values:

```go
fake.DoThingsReturns(3, errors.New("the-error"))

num, err := fake.DoThings("stuff", 5)
Expect(num).To(Equal(3))
Expect(err).To(Equal(errors.New("the-error")))
```

For more examples of using the `counterfeiter` API, look at [some of the provided examples](https://github.com/maxbrunsfeld/counterfeiter/blob/master/generated_fakes_test.go).

### Generating Test Doubles For Third Party Interfaces

For third party interfaces, you can specify the interface using the alternative syntax `<package>.<interface>`, for example:

```shell
$ go run github.com/maxbrunsfeld/counterfeiter/v6 github.com/go-redis/redis.Pipeliner
```

### Running The Tests For `counterfeiter`

If you want to run the tests for `counterfeiter` (perhaps, because you want to contribute a PR), all you have to do is run `scripts/ci.sh`.

### Contributions

So you want to contribute to `counterfeiter`! That's great, here's exactly what you should do:

- open a new github issue, describing your problem, or use case
- help us understand how you want to fix or extend `counterfeiter`
- write one or more unit tests for the behavior you want
- write the simplest code you can for the feature you're working on
- try to find any opportunities to refactor
- avoid writing code that isn't covered by unit tests

`counterfeiter` has a few high level goals for contributors to keep in mind

- keep unit-level test coverage as high as possible
- keep `main.go` as simple as possible
- avoid making the command line options any more complicated
- avoid making the internals of `counterfeiter` any more complicated

If you have any questions about how to contribute, rest assured that @tjarratt and other maintainers will work with you to ensure we make `counterfeiter` better, together. This project has largely been maintained by the community, and we greatly appreciate any PR (whether big or small).

### License

`counterfeiter` is MIT-licensed.
