# Extending go-git

`go-git` was built in a highly extensible manner, which enables some of its functionalities to be changed or extended without the need of changing its codebase. Here are the key extensibility features:

## Dot Git Storers

Dot git storers are the components responsible for storing the Git internal files, including objects and references.

The built-in storer implementations include [memory](storage/memory) and [filesystem](storage/filesystem). The `memory` storer stores all the data in memory, and its use look like this:

```go
	r, err := git.Init(memory.NewStorage(), nil)
```

The `filesystem` storer stores the data in the OS filesystem, and can be used as follows:

```go
    r, err := git.Init(filesystem.NewStorage(osfs.New("/tmp/foo")), nil)
```

New implementations can be created by implementing the [storage.Storer interface](storage/storer.go#L16).

## Filesystem

Git repository worktrees are managed using a filesystem abstraction based on [go-billy](https://github.com/go-git/go-billy). The Git operations will take place against the specific filesystem implementation. Initialising a repository in Memory can be done as follows:

```go
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
```

The same operation can be done against the OS filesystem:

```go
    fs := osfs.New("/tmp/foo")
    r, err := git.Init(memory.NewStorage(), fs)
```

New filesystems (e.g. cloud based storage) could be created by implementing `go-billy`'s [Filesystem interface](https://github.com/go-git/go-billy/blob/326c59f064021b821a55371d57794fbfb86d4cb3/fs.go#L52).

## Transport Schemes

Git supports various transport schemes, including `http`, `https`, `ssh`, `git`, `file`. `go-git` defines the [transport.Transport interface](plumbing/transport/common.go#L48) to represent them.

The built-in implementations can be replaced by calling `client.InstallProtocol`.

An example of changing the built-in `https` implementation to skip TLS could look like this:

```go
	customClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client.InstallProtocol("https", githttp.NewClient(customClient))
```

Some internal implementations enables code reuse amongst the different transport implementations. Some of these may be made public in the future (e.g. `plumbing/transport/internal/common`).

## Cache

Several different operations across `go-git` lean on caching of objects in order to achieve optimal performance. The caching functionality is defined by the [cache.Object interface](plumbing/cache/common.go#L17).

Two built-in implementations are `cache.ObjectLRU` and `cache.BufferLRU`. However, the caching functionality can be customized by implementing the interface `cache.Object` interface.

## Hash

`go-git` uses the `crypto.Hash` interface to represent hash functions. The built-in implementations are `github.com/pjbgf/sha1cd` for SHA1 and Go's `crypto/SHA256`.

The default hash functions can be changed by calling `hash.RegisterHash`.
```go
    func init() {
        hash.RegisterHash(crypto.SHA1, sha1.New)
    }
```

New `SHA1` or `SHA256` hash functions that implement the `hash.RegisterHash` interface can be registered by calling `RegisterHash`.
