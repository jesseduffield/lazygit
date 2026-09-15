# Mergo

[![GitHub release][5]][6]
[![GoCard][7]][8]
[![Test status][1]][2]
[![OpenSSF Scorecard][21]][22]
[![OpenSSF Best Practices][19]][20]
[![Coverage status][9]][10]
[![Sourcegraph][11]][12]
[![FOSSA status][13]][14]

[![GoDoc][3]][4]
[![Become my sponsor][15]][16]
[![Tidelift][17]][18]

[1]: https://github.com/imdario/mergo/workflows/tests/badge.svg?branch=master
[2]: https://github.com/imdario/mergo/actions/workflows/tests.yml
[3]: https://godoc.org/github.com/imdario/mergo?status.svg
[4]: https://godoc.org/github.com/imdario/mergo
[5]: https://img.shields.io/github/release/imdario/mergo.svg
[6]: https://github.com/imdario/mergo/releases
[7]: https://goreportcard.com/badge/imdario/mergo
[8]: https://goreportcard.com/report/github.com/imdario/mergo
[9]: https://coveralls.io/repos/github/imdario/mergo/badge.svg?branch=master
[10]: https://coveralls.io/github/imdario/mergo?branch=master
[11]: https://sourcegraph.com/github.com/imdario/mergo/-/badge.svg
[12]: https://sourcegraph.com/github.com/imdario/mergo?badge
[13]: https://app.fossa.io/api/projects/git%2Bgithub.com%2Fimdario%2Fmergo.svg?type=shield
[14]: https://app.fossa.io/projects/git%2Bgithub.com%2Fimdario%2Fmergo?ref=badge_shield
[15]: https://img.shields.io/github/sponsors/imdario
[16]: https://github.com/sponsors/imdario
[17]: https://tidelift.com/badges/package/go/github.com%2Fimdario%2Fmergo
[18]: https://tidelift.com/subscription/pkg/go-github.com-imdario-mergo
[19]: https://bestpractices.coreinfrastructure.org/projects/7177/badge
[20]: https://bestpractices.coreinfrastructure.org/projects/7177
[21]: https://api.securityscorecards.dev/projects/github.com/imdario/mergo/badge
[22]: https://api.securityscorecards.dev/projects/github.com/imdario/mergo

A helper to merge structs and maps in Golang. Useful for configuration default values, avoiding messy if-statements.

Mergo merges same-type structs and maps by setting default values in zero-value fields. Mergo won't merge unexported (private) fields. It will do recursively any exported one. It also won't merge structs inside maps (because they are not addressable using Go reflection).

Also a lovely [comune](http://en.wikipedia.org/wiki/Mergo) (municipality) in the Province of Ancona in the Italian region of Marche.

## Status

Mergo is stable and frozen, ready for production. Check a short list of the projects using at large scale it [here](https://github.com/imdario/mergo#mergo-in-the-wild).

No new features are accepted. They will be considered for a future v2 that improves the implementation and fixes bugs for corner cases.

### Important notes

#### 1.0.0

In [1.0.0](//github.com/imdario/mergo/releases/tag/1.0.0) Mergo moves to a vanity URL `dario.cat/mergo`. No more v1 versions will be released.

If the vanity URL is causing issues in your project due to a dependency pulling Mergo - it isn't a direct dependency in your project - it is recommended to use [replace](https://github.com/golang/go/wiki/Modules#when-should-i-use-the-replace-directive) to pin the version to the last one with the old import URL:

```
replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.16
```

#### 0.3.9

Please keep in mind that a problematic PR broke [0.3.9](//github.com/imdario/mergo/releases/tag/0.3.9). I reverted it in [0.3.10](//github.com/imdario/mergo/releases/tag/0.3.10), and I consider it stable but not bug-free. Also, this version adds support for go modules.

Keep in mind that in [0.3.2](//github.com/imdario/mergo/releases/tag/0.3.2), Mergo changed `Merge()`and `Map()` signatures to support [transformers](#transformers). I added an optional/variadic argument so that it won't break the existing code.

If you were using Mergo before April 6th, 2015, please check your project works as intended after updating your local copy with ```go get -u dario.cat/mergo```. I apologize for any issue caused by its previous behavior and any future bug that Mergo could cause in existing projects after the change (release 0.2.0).

### Donations

If Mergo is useful to you, consider buying me a coffee, a beer, or making a monthly donation to allow me to keep building great free software. :heart_eyes:

<a href="https://liberapay.com/dario/donate"><img alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg"></a>
<a href='https://github.com/sponsors/imdario' target='_blank'><img alt="Become my sponsor" src="https://img.shields.io/github/sponsors/imdario?style=for-the-badge" /></a>

### Mergo in the wild

Mergo is used by [thousands](https://deps.dev/go/dario.cat%2Fmergo/v1.0.0/dependents) [of](https://deps.dev/go/github.com%2Fimdario%2Fmergo/v0.3.16/dependents) [projects](https://deps.dev/go/github.com%2Fimdario%2Fmergo/v0.3.12), including:

* [containerd/containerd](https://github.com/containerd/containerd)
* [datadog/datadog-agent](https://github.com/datadog/datadog-agent)
* [docker/cli/](https://github.com/docker/cli/)
* [goreleaser/goreleaser](https://github.com/goreleaser/goreleaser)
* [go-micro/go-micro](https://github.com/go-micro/go-micro)
* [grafana/loki](https://github.com/grafana/loki)
* [masterminds/sprig](github.com/Masterminds/sprig)
* [moby/moby](https://github.com/moby/moby)
* [slackhq/nebula](https://github.com/slackhq/nebula)
* [volcano-sh/volcano](https://github.com/volcano-sh/volcano)

## Install

    go get dario.cat/mergo

    // use in your .go code
    import (
        "dario.cat/mergo"
    )

## Usage

You can only merge same-type structs with exported fields initialized as zero value of their type and same-types maps. Mergo won't merge unexported (private) fields but will do recursively any exported one. It won't merge empty structs value as [they are zero values](https://golang.org/ref/spec#The_zero_value) too. Also, maps will be merged recursively except for structs inside maps (because they are not addressable using Go reflection).

```go
if err := mergo.Merge(&dst, src); err != nil {
    // ...
}
```

Also, you can merge overwriting values using the transformer `WithOverride`.

```go
if err := mergo.Merge(&dst, src, mergo.WithOverride); err != nil {
    // ...
}
```

If you need to override pointers, so the source pointer's value is assigned to the destination's pointer, you must use `WithoutDereference`:

```go
package main

import (
	"fmt"

	"dario.cat/mergo"
)

type Foo struct {
	A *string
	B int64
}

func main() {
	first := "first"
	second := "second"
	src := Foo{
		A: &first,
		B: 2,
	}

	dest := Foo{
		A: &second,
		B: 1,
	}

	mergo.Merge(&dest, src, mergo.WithOverride, mergo.WithoutDereference)
}
```

Additionally, you can map a `map[string]interface{}` to a struct (and otherwise, from struct to map), following the same restrictions as in `Merge()`. Keys are capitalized to find each corresponding exported field.

```go
if err := mergo.Map(&dst, srcMap); err != nil {
    // ...
}
```

Warning: if you map a struct to map, it won't do it recursively. Don't expect Mergo to map struct members of your struct as `map[string]interface{}`. They will be just assigned as values.

Here is a nice example:

```go
package main

import (
	"fmt"
	"dario.cat/mergo"
)

type Foo struct {
	A string
	B int64
}

func main() {
	src := Foo{
		A: "one",
		B: 2,
	}
	dest := Foo{
		A: "two",
	}
	mergo.Merge(&dest, src)
	fmt.Println(dest)
	// Will print
	// {two 2}
}
```

### Transformers

Transformers allow to merge specific types differently than in the default behavior. In other words, now you can customize how some types are merged. For example, `time.Time` is a struct; it doesn't have zero value but IsZero can return true because it has fields with zero value. How can we merge a non-zero `time.Time`?

```go
package main

import (
	"fmt"
	"dario.cat/mergo"
    "reflect"
    "time"
)

type timeTransformer struct {
}

func (t timeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}

type Snapshot struct {
	Time time.Time
	// ...
}

func main() {
	src := Snapshot{time.Now()}
	dest := Snapshot{}
	mergo.Merge(&dest, src, mergo.WithTransformers(timeTransformer{}))
	fmt.Println(dest)
	// Will print
	// { 2018-01-12 01:15:00 +0000 UTC m=+0.000000001 }
}
```

## Contact me

If I can help you, you have an idea or you are using Mergo in your projects, don't hesitate to drop me a line (or a pull request): [@im_dario](https://twitter.com/im_dario)

## About

Written by [Dario Castañé](http://dario.im).

## License

[BSD 3-Clause](http://opensource.org/licenses/BSD-3-Clause) license, as [Go language](http://golang.org/LICENSE).

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fimdario%2Fmergo.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fimdario%2Fmergo?ref=badge_large)
