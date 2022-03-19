# lo

![Build Status](https://github.com/samber/lo/actions/workflows/go.yml/badge.svg)
[![GoDoc](https://godoc.org/github.com/samber/lo?status.svg)](https://pkg.go.dev/github.com/samber/lo)
[![Go report](https://goreportcard.com/badge/github.com/samber/lo)](https://goreportcard.com/report/github.com/samber/lo)

✨ **`lo` is a Lodash-style Go library based on Go 1.18+ Generics.**

This project started as an experiment with the new generics implementation. It may look like [Lodash](https://github.com/lodash/lodash) in some aspects. I used to code with the fantastic ["go-funk"](https://github.com/thoas/go-funk) package, but "go-funk" uses reflection and therefore is not typesafe.

As expected, benchmarks demonstrate that generics will be much faster than implementations based on the "reflect" package. Benchmarks also show similar performance gains compared to pure `for` loops. [See below](#-benchmark).

In the future, 5 to 10 helpers will overlap with those coming into the Go standard library (under package names `slices` and `maps`). I feel this library is legitimate and offers many more valuable abstractions.

### Why this name?

I wanted a **short name**, similar to "Lodash" and no Go package currently uses this name.

## 🚀 Install

```sh
go get github.com/samber/lo
```

## 💡 Usage

You can import `lo` using:

```go
import (
    "github.com/samber/lo"
    lop "github.com/samber/lo/parallel"
)
```

Then use one of the helpers below:

```go
names := lo.Uniq[string]([]string{"Samuel", "Marc", "Samuel"})
// []string{"Samuel", "Marc"}
```

Most of the time, the compiler will be able to infer the type so that you can call: `lo.Uniq([]string{...})`.

## 🤠 Spec

GoDoc: [https://godoc.org/github.com/samber/lo](https://godoc.org/github.com/samber/lo)

Supported helpers for slices:

- Filter
- Map
- FlatMap
- Reduce
- ForEach
- Times
- Uniq
- UniqBy
- GroupBy
- Chunk
- PartitionBy
- Flatten
- Shuffle
- Reverse
- Fill
- Repeat
- KeyBy
- Drop
- DropRight
- DropWhile
- DropRightWhile

Supported helpers for maps:

- Keys
- Values
- Entries
- FromEntries
- Assign (merge of maps)
- MapValues

Supported helpers for tuples:

- Zip2 -> Zip9
- Unzip2 -> Unzip9

Supported intersection helpers:

- Contains
- ContainsBy
- Every
- Some
- Intersect
- Difference
- Union

Supported search helpers:

- IndexOf
- LastIndexOf
- Find
- Min
- Max
- Last
- Nth
- Sample
- Samples

Other functional programming helpers:

- Ternary (1 line if/else statement)
- If / ElseIf / Else
- Switch / Case / Default
- ToPtr
- ToSlicePtr
- Attempt
- Range / RangeFrom / RangeWithSteps

Constraints:

- Clonable

### Map

Manipulates a slice of one type and transforms it into a slice of another type:

```go
import "github.com/samber/lo"

lo.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
    return strconv.FormatInt(x, 10)
})
// []string{"1", "2", "3", "4"}
```

Parallel processing: like `lo.Map()`, but the mapper function is called in a goroutine. Results are returned in the same order.

```go
import lop "github.com/samber/lo/parallel"

lop.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
    return strconv.FormatInt(x, 10)
})
// []string{"1", "2", "3", "4"}
```

### FlatMap

Manipulates a slice and transforms and flattens it to a slice of another type.

```go
lo.FlatMap[int, string]([]int{0, 1, 2}, func(x int, _ int) []string {
	return []string{
		strconv.FormatInt(x, 10),
		strconv.FormatInt(x, 10),
	}
})
// []string{"0", "0", "1", "1", "2", "2"}
```

### Filter

Iterates over a collection and returns an array of all the elements the predicate function returns `true` for.

```go
even := lo.Filter[int]([]int{1, 2, 3, 4}, func(x int, _ int) bool {
    return x%2 == 0
})
// []int{2, 4}
```

### Contains

Returns true if an element is present in a collection.

```go
present := lo.Contains[int]([]int{0, 1, 2, 3, 4, 5}, 5)
// true
```

### Contains

Returns true if the predicate function returns `true`.

```go
present := lo.ContainsBy[int]([]int{0, 1, 2, 3, 4, 5}, func(x int) bool {
    return x == 3
})
// true
```

### Reduce

Reduces a collection to a single value. The value is calculated by accumulating the result of running each element in the collection through an accumulator function. Each successive invocation is supplied with the return value returned by the previous call.

```go
sum := lo.Reduce[int, int]([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
    return agg + item
}, 0)
// 10
```

### ForEach

Iterates over elements of a collection and invokes the function over each element.

```go
import "github.com/samber/lo"

lo.ForEach[string]([]string{"hello", "world"}, func(x string, _ int) {
    println(x)
})
// prints "hello\nworld\n"
```

Parallel processing: like `lo.ForEach()`, but the callback is called as a goroutine.

```go
import lop "github.com/samber/lo/parallel"

lop.ForEach[string]([]string{"hello", "world"}, func(x string, _ int) {
    println(x)
})
// prints "hello\nworld\n" or "world\nhello\n"
```

### Times

Times invokes the iteratee n times, returning an array of the results of each invocation. The iteratee is invoked with index as argument.

```go
import "github.com/samber/lo"

lo.Times[string](3, func(i int) string {
    return strconv.FormatInt(int64(i), 10)
})
// []string{"0", "1", "2"}
```

Parallel processing: like `lo.Times()`, but callback is called in goroutine.

```go
import lop "github.com/samber/lo/parallel"

lop.Times[string](3, func(i int) string {
    return strconv.FormatInt(int64(i), 10)
})
// []string{"0", "1", "2"}
```

### Uniq

Returns a duplicate-free version of an array, in which only the first occurrence of each element is kept. The order of result values is determined by the order they occur in the array.

```go
uniqValues := lo.Uniq[int]([]int{1, 2, 2, 1})
// []int{1, 2}
```

### UniqBy

Returns a duplicate-free version of an array, in which only the first occurrence of each element is kept. The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is invoked for each element in array to generate the criterion by which uniqueness is computed.

```go
uniqValues := lo.UniqBy[int, int]([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// []int{0, 1, 2}
```

### GroupBy

Returns an object composed of keys generated from the results of running each element of collection through iteratee.

```go
import lo "github.com/samber/lo"

groups := lo.GroupBy[int, int]([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// map[int][]int{0: []int{0, 3}, 1: []int{1, 4}, 2: []int{2, 5}}
```

Parallel processing: like `lo.GroupBy()`, but callback is called in goroutine.

```go
import lop "github.com/samber/lo/parallel"

lop.GroupBy[int, int]([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// map[int][]int{0: []int{0, 3}, 1: []int{1, 4}, 2: []int{2, 5}}
```

### Chunk

Returns an array of elements split into groups the length of size. If array can't be split evenly, the final chunk will be the remaining elements.

```go
lo.Chunk[int]([]int{0, 1, 2, 3, 4, 5}, 2)
// [][]int{{0, 1}, {2, 3}, {4, 5}}

lo.Chunk[int]([]int{0, 1, 2, 3, 4, 5, 6}, 2)
// [][]int{{0, 1}, {2, 3}, {4, 5}, {6}}

lo.Chunk[int]([]int{}, 2)
// [][]int{}

lo.Chunk[int]([]int{0}, 2)
// [][]int{{0}}
```

### PartitionBy

Returns an array of elements split into groups. The order of grouped values is determined by the order they occur in collection. The grouping is generated from the results of running each element of collection through iteratee.

```go
import lo "github.com/samber/lo"

partitions := lo.PartitionBy[int, string]([]int{-2, -1, 0, 1, 2, 3, 4, 5}, func(x int) string {
    if x < 0 {
        return "negative"
    } else if x%2 == 0 {
        return "even"
    }
    return "odd"
})
// [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}}
```

Parallel processing: like `lo.PartitionBy()`, but callback is called in goroutine. Results are returned in the same order.

```go
import lop "github.com/samber/lo/parallel"

partitions := lo.PartitionBy[int, string]([]int{-2, -1, 0, 1, 2, 3, 4, 5}, func(x int) string {
    if x < 0 {
        return "negative"
    } else if x%2 == 0 {
        return "even"
    }
    return "odd"
})
// [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}}
```

### Flatten

Returns an array a single level deep.

```go
flat := lo.Flatten[int]([][]int{{0, 1}, {2, 3, 4, 5}})
// []int{0, 1, 2, 3, 4, 5}
```

### Shuffle

Returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm.

```go
randomOrder := lo.Shuffle[int]([]int{0, 1, 2, 3, 4, 5})
// []int{0, 1, 2, 3, 4, 5}
```

### Reverse

Reverses array so that the first element becomes the last, the second element becomes the second to last, and so on.

```go
reverseOder := lo.Reverse[int]([]int{0, 1, 2, 3, 4, 5})
// []int{5, 4, 3, 2, 1, 0}
```

### Fill

Fills elements of array with `initial` value.

```go
type foo struct {
	bar string
}

func (f foo) Clone() foo {
	return foo{f.bar}
}

initializedSlice := lo.Fill[foo]([]foo{foo{"a"}, foo{"a"}}, foo{"b"})
// []foo{foo{"b"}, foo{"b"}}
```

### Repeat

Builds a slice with N copies of initial value.

```go
type foo struct {
	bar string
}

func (f foo) Clone() foo {
	return foo{f.bar}
}

initializedSlice := lo.Repeat[foo](2, foo{"a"})
// []foo{foo{"a"}, foo{"a"}}
```

### KeyBy

Transforms a slice or an array of structs to a map based on a pivot callback.

```go
m := lo.KeyBy[int, string]([]string{"a", "aa", "aaa"}, func(str string) int {
    return len(str)
})
// map[int]string{1: "a", 2: "aa", 3: "aaa"}

type Character struct {
	dir  string
	code int
}
characters := []Character{
    {dir: "left", code: 97},
    {dir: "right", code: 100},
}
result := KeyBy[Character, string](characters, func(char Character) string {
    return string(rune(char.code))
})
//map[a:{dir:left code:97} d:{dir:right code:100}]
```

### Drop

Drops n elements from the beginning of a slice or array.

```go
l := lo.Drop[int]([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{2, 3, 4, 5}
```

### DropRight

Drops n elements from the end of a slice or array.

```go
l := lo.DropRight[int]([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{0, 1, 2, 3}
```

### DropWhile

Drop elements from the beginning of a slice or array while the predicate returns true.

```go
l := lo.DropWhile[string]([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
	return len(val) <= 2
})
// []string{"aaa", "aa", "a"}
```

### DropRightWhile

Drop elements from the end of a slice or array while the predicate returns true.

```go
l := lo.DropRightWhile[string]([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
	return len(val) <= 2
})
// []string{"a", "aa", "aaa"}
```

### Keys

Creates an array of the map keys.

```go
keys := lo.Keys[string, int](map[string]int{"foo": 1, "bar": 2})
// []string{"bar", "foo"}
```

### Values

Creates an array of the map values.

```go
values := lo.Values[string, int](map[string]int{"foo": 1, "bar": 2})
// []int{1, 2}
```

### Entries

Transforms a map into array of key/value pairs.

```go
entries := lo.Entries[string, int](map[string]int{"foo": 1, "bar": 2})
// []lo.Entry[string, int]{
//     {
//         Key: "foo",
//         Value: 1,
//     },
//     {
//         Key: "bar",
//         Value: 2,
//     },
// }
```

### FromEntries

Transforms an array of key/value pairs into a map.

```go
m := lo.FromEntries[string, int]([]lo.Entry[string, int]{
    {
        Key: "foo",
        Value: 1,
    },
    {
        Key: "bar",
        Value: 2,
    },
})
// map[string]int{"foo": 1, "bar": 2}
```

### Assign

Merges multiple maps from left to right.

```go
mergedMaps := lo.Assign[string, int](
    map[string]int{"a": 1, "b": 2},
    map[string]int{"b": 3, "c": 4},
)
// map[string]int{"a": 1, "b": 3, "c": 4}
```

### MapValues

Manipulates a map values and transforms it to a map of another type.

```go
m1 := map[int]int64{1: 1, 2: 2, 3: 3}

m2 := lo.MapValues[int, int64, string](m, func(x int64, _ int) string {
	return strconv.FormatInt(x, 10)
})
// map[int]string{1: "1", 2: "2", 3: "3"}
```

### Zip2 -> Zip9

Zip creates a slice of grouped elements, the first of which contains the first elements of the given arrays, the second of which contains the second elements of the given arrays, and so on.

When collections have different size, the Tuple attributes are filled with zero value.

```go
tuples := lo.Zip2[string, int]([]string{"a", "b"}, []int{1, 2})
// []Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}}
```

### Unzip2 -> Unzip9

Unzip accepts an array of grouped elements and creates an array regrouping the elements to their pre-zip configuration.

```go
a, b := lo.Unzip2[string, int]([]Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}})
// []string{"a", "b"}
// []int{1, 2}
```

### Every

Returns true if all elements of a subset are contained into a collection.

```go
ok := lo.Every[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// true

ok := lo.Every[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// false
```

### Some

Returns true if at least 1 element of a subset is contained into a collection.

```go
ok := lo.Some[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// true

ok := lo.Some[int]([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// false
```

### Intersect

Returns the intersection between two collections.

```go
result1 := lo.Intersect[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// []int{0, 2}

result2 := lo.Intersect[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 6}
// []int{0}

result3 := lo.Intersect[int]([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// []int{}
```

### Difference

Returns the difference between two collections.

- The first value is the collection of element absent of list2.
- The second value is the collection of element absent of list1.

```go
left, right := lo.Difference[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
// []int{1, 3, 4, 5}, []int{6}

left, right := Difference[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
// []int{}, []int{}
```

### Union

Returns all distinct elements from both collections. Result will not change the order of elements relatively.

```go
union := lo.Union[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10})
// []int{0, 1, 2, 3, 4, 5, 10}
```

### IndexOf

Returns the index at which the first occurrence of a value is found in an array or return -1 if the value cannot be found.

```go
found := lo.IndexOf[int]([]int{0, 1, 2, 1, 2, 3}, 2)
// 2

notFound := lo.IndexOf[int]([]int{0, 1, 2, 1, 2, 3}, 6)
// -1
```

### LastIndex

Returns the index at which the last occurrence of a value is found in an array or return -1 if the value cannot be found.

```go
found := lo.LastIndexOf[int]([]int{0, 1, 2, 1, 2, 3}, 2)
// 4

notFound := lo.LastIndexOf[int]([]int{0, 1, 2, 1, 2, 3}, 6)
// -1
```

### Find

Search an element in a slice based on a predicate. It returns element and true if element was found.

```go
str, ok := lo.Find[string]([]string{"a", "b", "c", "d"}, func(i string) bool {
    return i == "b"
})
// "b", true

str, ok := lo.Find[string]([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", false
```

### Min

Search the minimum value of a collection.

```go
min := lo.Min[int]([]int{1, 2, 3})
// 1

min := lo.Min[int]([]int{})
// 0
```

### Max

Search the maximum value of a collection.

```go
max := lo.Max[int]([]int{1, 2, 3})
// 3

max := lo.Max[int]([]int{})
// 0
```

### Last

Returns the last element of a collection or error if empty.

```go
last, err := lo.Last[int]([]int{1, 2, 3})
// 3
```

### Nth

Returns the element at index `nth` of collection. If `nth` is negative, the nth element from the end is returned. An error is returned when nth is out of slice bounds.

```go
nth, err := lo.Nth[int]([]int{0, 1, 2, 3}, 2)
// 2

nth, err := lo.Nth[int]([]int{0, 1, 2, 3}, -2)
// 2
```

### Sample

Returns a random item from collection.

```go
lo.Sample[string]([]string{"a", "b", "c"})
// a random string from []string{"a", "b", "c"}

lo.Sample[string]([]string{})
// ""
```

### Samples

Returns N random unique items from collection.

```go
lo.Samples[string]([]string{"a", "b", "c"}, 3)
// []string{"a", "b", "c"} in random order
```

### Ternary

A 1 line if/else statement.

```go
result := lo.Ternary[string](true, "a", "b")
// "a"

result := lo.Ternary[string](false, "a", "b")
// "b"
```

### If / ElseIf / Else

```go
result := lo.If[int](true, 1).
    ElseIf(false, 2).
    Else(3)
// 1

result := lo.If[int](false, 1).
    ElseIf(true, 2).
    Else(3)
// 2

result := lo.If[int](false, 1).
    ElseIf(false, 2).
    Else(3)
// 3
```

### Switch / Case / Default

```go
result := lo.Switch[int, string](1).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "1"

result := lo.Switch[int, string](2).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "2"

result := lo.Switch[int, string](42).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "3"
```

Using callbacks:

```go
result := lo.Switch[int, string](1).
    CaseF(1, func() string {
        return "1"
    }).
    CaseF(2, func() string {
        return "2"
    }).
    DefaultF(func() string {
        return "3"
    })
// "1"
```

### ToPtr

Returns a pointer copy of value.

```go
ptr := lo.ToPtr[string]("hello world")
// *string{"hello world"}
```

### ToSlicePtr

Returns a slice of pointer copy of value.

```go
ptr := lo.ToSlicePtr[string]([]string{"hello", "world"})
// []*string{"hello", "world"}
```

### Attempt

Invokes a function N times until it returns valid output. Returning either the caught error or nil. When first argument is less than `1`, the function runs until a sucessfull response is returned.

```go
iter, err := lo.Attempt(42, func(i int) error {
    if i == 5 {
        return nil
    }

    return fmt.Errorf("failed")
})
// 6
// nil

iter, err := lo.Attempt(2, func(i int) error {
    if i == 5 {
        return nil
    }

    return fmt.Errorf("failed")
})
// 2
// error "failed"

iter, err := lo.Attempt(0, func(i int) error {
    if i < 42 {
        return fmt.Errorf("failed")
    }

    return nil
})
// 43
// nil
```

### Range / RangeFrom / RangeWithSteps

Creates an array of numbers (positive and/or negative) progressing from start up to, but not including end.

```go
result := Range(4)
// [0, 1, 2, 3]

result := Range(-4);
// [0, -1, -2, -3]

result := RangeFrom(1, 5);
// [1, 2, 3, 4]

result := RangeFrom[float64](1.0, 5);
// [1.0, 2.0, 3.0, 4.0]

result := RangeWithSteps(0, 20, 5);
// [0, 5, 10, 15]

result := RangeWithSteps[float32](-1.0, -4.0, -1.0);
// [-1.0, -2.0, -3.0]

result := RangeWithSteps(1, 4, -1);
// []

result := Range(0);
// []
```

For more advanced retry strategies (delay, exponential backoff...), please take a look on [cenkalti/backoff](https://github.com/cenkalti/backoff).

## 🛩 Benchmark

We executed a simple benchmark with the a dead-simple `lo.Map` loop:

See the full implementation [here](./benchmark_test.go).

```go
_ = lo.Map[int64](arr, func(x int64, i int) string {
    return strconv.FormatInt(x, 10)
})
```

**Result:**

Here is a comparison between `lo.Map`, `lop.Map`, `go-funk` library and a simple Go `for` loop.

```
$ go test -benchmem -bench ./...
goos: linux
goarch: amd64
pkg: github.com/samber/lo
cpu: Intel(R) Core(TM) i5-7267U CPU @ 3.10GHz
cpu: Intel(R) Core(TM) i7 CPU         920  @ 2.67GHz
BenchmarkMap/lo.Map-8         	       8	 132728237 ns/op	39998945 B/op	 1000002 allocs/op
BenchmarkMap/lop.Map-8        	       2	 503947830 ns/op	119999956 B/op	 3000007 allocs/op
BenchmarkMap/reflect-8        	       2	 826400560 ns/op	170326512 B/op	 4000042 allocs/op
BenchmarkMap/for-8            	       9	 126252954 ns/op	39998674 B/op	 1000001 allocs/op
PASS
ok  	github.com/samber/lo	6.657s
```

- `lo.Map` is way faster (x7) than `go-funk`, a relection-based Map implementation.
- `lo.Map` have the same allocation profile than `for`.
- `lo.Map` is 4% slower than `for`.
- `lop.Map` is slower than `lo.Map` because it implies more memory allocation and locks. `lop.Map` will be usefull for long-running callbacks, such as i/o bound processing.
- `for` beats other implementations for memory and CPU.

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/lo)
- Fix [open issues](https://github.com/samber/lo/issues) or request new features

Don't hesitate ;)

### Install go 1.18

```bash
make go1.18beta1
```

If your OS currently not default to Go 1.18, replace `BIN=go` by `BIN=go1.18beta1` in the Makefile.

### With Docker

```bash
docker-compose run --rm dev
```

### Without Docker

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Authors

- Samuel Berthe

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![support us](https://c5.patreon.com/external/logo/become_a_patron_button.png)](https://www.patreon.com/samber)

## 📝 License

Copyright © 2022 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
