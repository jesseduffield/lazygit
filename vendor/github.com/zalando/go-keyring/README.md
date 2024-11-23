# Go Keyring library
[![Go Report Card](https://goreportcard.com/badge/github.com/zalando/go-keyring)](https://goreportcard.com/report/github.com/zalando/go-keyring)
[![GoDoc](https://godoc.org/github.com/zalando/go-keyring?status.svg)](https://godoc.org/github.com/zalando/go-keyring)

`go-keyring` is an OS-agnostic library for *setting*, *getting* and *deleting*
secrets from the system keyring. It supports **OS X**, **Linux/BSD (dbus)** and
**Windows**.

go-keyring was created after its authors searched for, but couldn't find, a better alternative. It aims to simplify
using statically linked binaries, which is cumbersome when relying on C bindings (as other keyring libraries do).

#### Potential Uses

If you're working with an application that needs to store user credentials
locally on the user's machine, go-keyring might come in handy. For instance, if you are writing a CLI for an API
that requires a username and password, you can store this information in the
keyring instead of having the user type it on every invocation.

## Dependencies

#### OS X

The OS X implementation depends on the `/usr/bin/security` binary for
interfacing with the OS X keychain. It should be available by default.

#### Linux and *BSD

The Linux and *BSD implementation depends on the [Secret Service][SecretService] dbus
interface, which is provided by [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring).

It's expected that the default collection `login` exists in the keyring, because
it's the default in most distros. If it doesn't exist, you can create it through the
keyring frontend program [Seahorse](https://wiki.gnome.org/Apps/Seahorse):

 * Open `seahorse`
 * Go to **File > New > Password Keyring**
 * Click **Continue**
 * When asked for a name, use: **login**

## Example Usage

How to *set* and *get* a secret from the keyring:

```go
package main

import (
    "log"

    "github.com/zalando/go-keyring"
)

func main() {
    service := "my-app"
    user := "anon"
    password := "secret"

    // set password
    err := keyring.Set(service, user, password)
    if err != nil {
        log.Fatal(err)
    }

    // get password
    secret, err := keyring.Get(service, user)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(secret)
}

```

## Tests
### Running tests
Running the tests is simple:

```
go test
```

Which OS you use *does* matter. If you're using **Linux** or **BSD**, it will
test the implementation in `keyring_unix.go`. If running the tests
on **OS X**, it will test the implementation in `keyring_darwin.go`.

### Mocking
If you need to mock the keyring behavior for testing on systems without a keyring implementation you can call `MockInit()` which will replace the OS defined provider with an in-memory one.

```go
package implementation

import (
    "testing"

    "github.com/zalando/go-keyring"
)

func TestMockedSetGet(t *testing.T) {
    keyring.MockInit()
    err := keyring.Set("service", "user", "password")
    if err != nil {
        t.Fatal(err)
    }

    p, err := keyring.Get("service", "user")
    if err != nil {
        t.Fatal(err)
    }

    if p != "password" {
        t.Error("password was not the expected string")
    }

}

```

## Contributing/TODO

We welcome contributions from the community; please use [CONTRIBUTING.md](CONTRIBUTING.md) as your guidelines for getting started. Here are some items that we'd love help with:

- The code base
- Better test coverage

Please use GitHub issues as the starting point for contributions, new ideas and/or bug reports.

## Contact

* E-Mail: team-teapot@zalando.de
* Security issues: Please send an email to the [maintainers](MAINTAINERS), and we'll try to get back to you within two workdays. If you don't hear back, send an email to team-teapot@zalando.de and someone will respond within five days max.

## Contributors

Thanks to:

- [your name here]

## License

See [LICENSE](LICENSE) file.


[SecretService]: https://specifications.freedesktop.org/secret-service/latest/
