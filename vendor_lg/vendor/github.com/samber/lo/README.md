# lo

[![tag](https://img.shields.io/github/tag/samber/lo.svg)](https://github.com/samber/lo/releases)
[![GoDoc](https://godoc.org/github.com/samber/lo?status.svg)](https://pkg.go.dev/github.com/samber/lo)
![Build Status](https://github.com/samber/lo/actions/workflows/go.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/lo)](https://goreportcard.com/report/github.com/samber/lo)
[![codecov](https://codecov.io/gh/samber/lo/branch/master/graph/badge.svg)](https://codecov.io/gh/samber/lo)

✨ **`samber/lo` is a Lodash-style Go library based on Go 1.18+ Generics.**

This project started as an experiment with the new generics implementation. It may look like [Lodash](https://github.com/lodash/lodash) in some aspects. I used to code with the fantastic ["go-funk"](https://github.com/thoas/go-funk) package, but "go-funk" uses reflection and therefore is not typesafe.

As expected, benchmarks demonstrate that generics will be much faster than implementations based on the "reflect" package. Benchmarks also show similar performance gains compared to pure `for` loops. [See below](#-benchmark).

In the future, 5 to 10 helpers will overlap with those coming into the Go standard library (under package names `slices` and `maps`). I feel this library is legitimate and offers many more valuable abstractions.

**See also:**

- [samber/do](https://github.com/samber/do): A dependency injection toolkit based on Go 1.18+ Generics
- [samber/mo](https://github.com/samber/mo): Monads based on Go 1.18+ Generics (Option, Result, Either...)

**Why this name?**

I wanted a **short name**, similar to "Lodash" and no Go package currently uses this name.

## 🚀 Install

```sh
go get github.com/samber/lo@v1
```

This library is v1 and follows SemVer strictly.

No breaking changes will be made to exported APIs before v2.0.0.

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

- [Filter](#filter)
- [Map](#map)
- [FilterMap](#filtermap)
- [FlatMap](#flatmap)
- [Reduce](#reduce)
- [ReduceRight](#reduceright)
- [ForEach](#foreach)
- [Times](#times)
- [Uniq](#uniq)
- [UniqBy](#uniqby)
- [GroupBy](#groupby)
- [Chunk](#chunk)
- [PartitionBy](#partitionby)
- [Flatten](#flatten)
- [Shuffle](#shuffle)
- [Reverse](#reverse)
- [Fill](#fill)
- [Repeat](#repeat)
- [RepeatBy](#repeatby)
- [KeyBy](#keyby)
- [Associate / SliceToMap](#associate-alias-slicetomap)
- [Drop](#drop)
- [DropRight](#dropright)
- [DropWhile](#dropwhile)
- [DropRightWhile](#droprightwhile)
- [Reject](#reject)
- [Count](#count)
- [CountBy](#countby)
- [Subset](#subset)
- [Slice](#slice)
- [Replace](#replace)
- [ReplaceAll](#replaceall)
- [Compact](#compact)
- [IsSorted](#issorted)
- [IsSortedByKey](#issortedbykey)

Supported helpers for maps:

- [Keys](#keys)
- [Values](#values)
- [PickBy](#pickby)
- [PickByKeys](#pickbykeys)
- [PickByValues](#pickbyvalues)
- [OmitBy](#omitby)
- [OmitByKeys](#omitbykeys)
- [OmitByValues](#omitbyvalues)
- [Entries / ToPairs](#entries-alias-topairs)
- [FromEntries / FromPairs](#fromentries-alias-frompairs)
- [Invert](#invert)
- [Assign (merge of maps)](#assign)
- [MapKeys](#mapkeys)
- [MapValues](#mapvalues)
- [MapToSlice](#maptoslice)

Supported math helpers:

- [Range / RangeFrom / RangeWithSteps](#range--rangefrom--rangewithsteps)
- [Clamp](#clamp)
- [SumBy](#sumby)

Supported helpers for strings:

- [Substring](#substring)
- [ChunkString](#chunkstring)
- [RuneLength](#runelength)

Supported helpers for tuples:

- [T2 -> T9](#t2---t9)
- [Unpack2 -> Unpack9](#unpack2---unpack9)
- [Zip2 -> Zip9](#zip2---zip9)
- [Unzip2 -> Unzip9](#unzip2---unzip9)

Supported helpers for channels:

- [ChannelDispatcher](#channeldispatcher)
- [SliceToChannel](#slicetochannel)
- [Generator](#generator)
- [Batch](#batch)
- [BatchWithTimeout](#batchwithtimeout)

Supported intersection helpers:

- [Contains](#contains)
- [ContainsBy](#containsby)
- [Every](#every)
- [EveryBy](#everyby)
- [Some](#some)
- [SomeBy](#someby)
- [None](#none)
- [NoneBy](#noneby)
- [Intersect](#intersect)
- [Difference](#difference)
- [Union](#union)
- [Without](#without)
- [WithoutEmpty](#withoutempty)

Supported search helpers:

- [IndexOf](#indexof)
- [LastIndexOf](#lastindexof)
- [Find](#find)
- [FindIndexOf](#findindexof)
- [FindLastIndexOf](#findlastindexof)
- [FindKey](#findkey)
- [FindKeyBy](#findkeyby)
- [FindUniques](#finduniques)
- [FindUniquesBy](#finduniquesby)
- [FindDuplicates](#findduplicates)
- [FindDuplicatesBy](#findduplicatesby)
- [Min](#min)
- [MinBy](#minby)
- [Max](#max)
- [MaxBy](#maxby)
- [Last](#last)
- [Nth](#nth)
- [Sample](#sample)
- [Samples](#samples)

Conditional helpers:

- [Ternary](#ternary)
- [TernaryF](#ternaryf)
- [If / ElseIf / Else](#if--elseif--else)
- [Switch / Case / Default](#switch--case--default)

Type manipulation helpers:

- [ToPtr](#toptr)
- [FromPtr](#fromptr)
- [FromPtrOr](#fromptror)
- [ToSlicePtr](#tosliceptr)
- [ToAnySlice](#toanyslice)
- [FromAnySlice](#fromanyslice)
- [Empty](#empty)
- [IsEmpty](#isempty)
- [IsNotEmpty](#isnotempty)
- [Coalesce](#coalesce)

Function helpers:

- [Partial](#partial)

Concurrency helpers:

- [Attempt](#attempt)
- [AttemptWithDelay](#attemptwithdelay)
- [Debounce](#debounce)
- [Synchronize](#synchronize)
- [Async](#async)

Error handling:

- [Validate](#validate)
- [Must](#must)
- [Try](#try)
- [Try1 -> Try6](#try0-6)
- [TryOr](#tryor)
- [TryOr1 -> TryOr6](#tryor0-6)
- [TryCatch](#trycatch)
- [TryWithErrorValue](#trywitherrorvalue)
- [TryCatchWithErrorValue](#trycatchwitherrorvalue)
- [ErrorsAs](#errorsas)

Constraints:

- Clonable

### Filter

Iterates over a collection and returns an array of all the elements the predicate function returns `true` for.

```go
even := lo.Filter[int]([]int{1, 2, 3, 4}, func(x int, index int) bool {
    return x%2 == 0
})
// []int{2, 4}
```

[[play](https://go.dev/play/p/Apjg3WeSi7K)]

### Map

Manipulates a slice of one type and transforms it into a slice of another type:

```go
import "github.com/samber/lo"

lo.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, index int) string {
    return strconv.FormatInt(x, 10)
})
// []string{"1", "2", "3", "4"}
```

[[play](https://go.dev/play/p/OkPcYAhBo0D)]

Parallel processing: like `lo.Map()`, but the mapper function is called in a goroutine. Results are returned in the same order.

```go
import lop "github.com/samber/lo/parallel"

lop.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
    return strconv.FormatInt(x, 10)
})
// []string{"1", "2", "3", "4"}
```

### FilterMap

Returns a slice which obtained after both filtering and mapping using the given callback function.

The callback function should return two values: the result of the mapping operation and whether the result element should be included or not.

```go
matching := lo.FilterMap[string, string]([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (string, bool) {
    if strings.HasSuffix(x, "pu") {
        return "xpu", true
    }
    return "", false
})
// []string{"xpu", "xpu"}
```

[[play](https://go.dev/play/p/-AuYXfy7opz)]

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

[[play](https://go.dev/play/p/YSoYmQTA8-U)]

### Reduce

Reduces a collection to a single value. The value is calculated by accumulating the result of running each element in the collection through an accumulator function. Each successive invocation is supplied with the return value returned by the previous call.

```go
sum := lo.Reduce[int, int]([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
    return agg + item
}, 0)
// 10
```

[[play](https://go.dev/play/p/R4UHXZNaaUG)]

### ReduceRight

Like `lo.Reduce` except that it iterates over elements of collection from right to left.

```go
result := lo.ReduceRight[[]int, []int]([][]int{{0, 1}, {2, 3}, {4, 5}}, func(agg []int, item []int, _ int) []int {
	  return append(agg, item...)
}, []int{})
// []int{4, 5, 2, 3, 0, 1}
```

[[play](https://go.dev/play/p/Fq3W70l7wXF)]

### ForEach

Iterates over elements of a collection and invokes the function over each element.

```go
import "github.com/samber/lo"

lo.ForEach[string]([]string{"hello", "world"}, func(x string, _ int) {
    println(x)
})
// prints "hello\nworld\n"
```

[[play](https://go.dev/play/p/oofyiUPRf8t)]

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

[[play](https://go.dev/play/p/vgQj3Glr6lT)]

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

[[play](https://go.dev/play/p/DTzbeXZ6iEN)]

### UniqBy

Returns a duplicate-free version of an array, in which only the first occurrence of each element is kept. The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is invoked for each element in array to generate the criterion by which uniqueness is computed.

```go
uniqValues := lo.UniqBy[int, int]([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// []int{0, 1, 2}
```

[[play](https://go.dev/play/p/g42Z3QSb53u)]

### GroupBy

Returns an object composed of keys generated from the results of running each element of collection through iteratee.

```go
import lo "github.com/samber/lo"

groups := lo.GroupBy[int, int]([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// map[int][]int{0: []int{0, 3}, 1: []int{1, 4}, 2: []int{2, 5}}
```

[[play](https://go.dev/play/p/XnQBd_v6brd)]

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

[[play](https://go.dev/play/p/EeKl0AuTehH)]

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

[[play](https://go.dev/play/p/NfQ_nGjkgXW)]

Parallel processing: like `lo.PartitionBy()`, but callback is called in goroutine. Results are returned in the same order.

```go
import lop "github.com/samber/lo/parallel"

partitions := lop.PartitionBy[int, string]([]int{-2, -1, 0, 1, 2, 3, 4, 5}, func(x int) string {
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

[[play](https://go.dev/play/p/rbp9ORaMpjw)]

### Shuffle

Returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm.

```go
randomOrder := lo.Shuffle[int]([]int{0, 1, 2, 3, 4, 5})
// []int{1, 4, 0, 3, 5, 2}
```

[[play](https://go.dev/play/p/Qp73bnTDnc7)]

### Reverse

Reverses array so that the first element becomes the last, the second element becomes the second to last, and so on.

```go
reverseOrder := lo.Reverse[int]([]int{0, 1, 2, 3, 4, 5})
// []int{5, 4, 3, 2, 1, 0}
```

[[play](https://go.dev/play/p/fhUMLvZ7vS6)]

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

[[play](https://go.dev/play/p/VwR34GzqEub)]

### Repeat

Builds a slice with N copies of initial value.

```go
type foo struct {
	bar string
}

func (f foo) Clone() foo {
	return foo{f.bar}
}

slice := lo.Repeat[foo](2, foo{"a"})
// []foo{foo{"a"}, foo{"a"}}
```

[[play](https://go.dev/play/p/g3uHXbmc3b6)]

### RepeatBy

Builds a slice with values returned by N calls of callback.

```go
slice := lo.RepeatBy[string](0, func (i int) string {
    return strconv.FormatInt(int64(math.Pow(float64(i), 2)), 10)
})
// []int{}

slice := lo.RepeatBy[string](5, func(i int) string {
    return strconv.FormatInt(int64(math.Pow(float64(i), 2)), 10)
})
// []int{0, 1, 4, 9, 16}
```

[[play](https://go.dev/play/p/ozZLCtX_hNU)]

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
result := lo.KeyBy[string, Character](characters, func(char Character) string {
    return string(rune(char.code))
})
//map[a:{dir:left code:97} d:{dir:right code:100}]
```

[[play](https://go.dev/play/p/mdaClUAT-zZ)]

### Associate (alias: SliceToMap)

Returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
If any of two pairs would have the same key the last one gets added to the map.

The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.

```go
in := []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}}

aMap := lo.Associate[*foo, string, int](in, func (f *foo) (string, int) {
	return f.baz, f.bar
})
// map[string][int]{ "apple":1, "banana":2 }
```

[[play](https://go.dev/play/p/WHa2CfMO3Lr)]

### Drop

Drops n elements from the beginning of a slice or array.

```go
l := lo.Drop[int]([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{2, 3, 4, 5}
```

[[play](https://go.dev/play/p/JswS7vXRJP2)]

### DropRight

Drops n elements from the end of a slice or array.

```go
l := lo.DropRight[int]([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{0, 1, 2, 3}
```

[[play](https://go.dev/play/p/GG0nXkSJJa3)]

### DropWhile

Drop elements from the beginning of a slice or array while the predicate returns true.

```go
l := lo.DropWhile[string]([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
	return len(val) <= 2
})
// []string{"aaa", "aa", "aa"}
```

[[play](https://go.dev/play/p/7gBPYw2IK16)]

### DropRightWhile

Drop elements from the end of a slice or array while the predicate returns true.

```go
l := lo.DropRightWhile[string]([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
	return len(val) <= 2
})
// []string{"a", "aa", "aaa"}
```

[[play](https://go.dev/play/p/3-n71oEC0Hz)]

### Reject

The opposite of Filter, this method returns the elements of collection that predicate does not return truthy for.

```go
odd := lo.Reject[int]([]int{1, 2, 3, 4}, func(x int, _ int) bool {
    return x%2 == 0
})
// []int{1, 3}
```

[[play](https://go.dev/play/p/YkLMODy1WEL)]

### Count

Counts the number of elements in the collection that compare equal to value.

```go
count := lo.Count[int]([]int{1, 5, 1}, 1)
// 2
```

[[play](https://go.dev/play/p/Y3FlK54yveC)]

### CountBy

Counts the number of elements in the collection for which predicate is true.

```go
count := lo.CountBy[int]([]int{1, 5, 1}, func(i int) bool {
    return i < 4
})
// 2
```

[[play](https://go.dev/play/p/ByQbNYQQi4X)]

### Subset

Returns a copy of a slice from `offset` up to `length` elements. Like `slice[start:start+length]`, but does not panic on overflow.

```go
in := []int{0, 1, 2, 3, 4}

sub := lo.Subset(in, 2, 3)
// []int{2, 3, 4}

sub := lo.Subset(in, -4, 3)
// []int{1, 2, 3}

sub := lo.Subset(in, -2, math.MaxUint)
// []int{3, 4}
```

[[play](https://go.dev/play/p/tOQu1GhFcog)]

### Slice

Returns a copy of a slice from `start` up to, but not including `end`. Like `slice[start:end]`, but does not panic on overflow.

```go
in := []int{0, 1, 2, 3, 4}

slice := lo.Slice(in, 0, 5)
// []int{0, 1, 2, 3, 4}

slice := lo.Slice(in, 2, 3)
// []int{2}

slice := lo.Slice(in, 2, 6)
// []int{2, 3, 4}

slice := lo.Slice(in, 4, 3)
// []int{}
```

[[play](https://go.dev/play/p/8XWYhfMMA1h)]

### Replace

Returns a copy of the slice with the first n non-overlapping instances of old replaced by new.

```go
in := []int{0, 1, 0, 1, 2, 3, 0}

slice := lo.Replace(in, 0, 42, 1)
// []int{42, 1, 0, 1, 2, 3, 0}

slice := lo.Replace(in, -1, 42, 1)
// []int{0, 1, 0, 1, 2, 3, 0}

slice := lo.Replace(in, 0, 42, 2)
// []int{42, 1, 42, 1, 2, 3, 0}

slice := lo.Replace(in, 0, 42, -1)
// []int{42, 1, 42, 1, 2, 3, 42}
```

[[play](https://go.dev/play/p/XfPzmf9gql6)]

### ReplaceAll

Returns a copy of the slice with all non-overlapping instances of old replaced by new.

```go
in := []int{0, 1, 0, 1, 2, 3, 0}

slice := lo.ReplaceAll(in, 0, 42)
// []int{42, 1, 42, 1, 2, 3, 42}

slice := lo.ReplaceAll(in, -1, 42)
// []int{0, 1, 0, 1, 2, 3, 0}
```

[[play](https://go.dev/play/p/a9xZFUHfYcV)]

### Compact

Returns a slice of all non-zero elements.

```go
in := []string{"", "foo", "", "bar", ""}

slice := lo.Compact[string](in)
// []string{"foo", "bar"}
```

[[play](https://go.dev/play/p/tXiy-iK6PAc)]

### IsSorted

Checks if a slice is sorted.

```go
slice := lo.IsSorted([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
// true
```

[[play](https://go.dev/play/p/mc3qR-t4mcx)]

### IsSortedByKey

Checks if a slice is sorted by iteratee.

```go
slice := lo.IsSortedByKey([]string{"a", "bb", "ccc"}, func(s string) int {
    return len(s)
})
// true
```

[[play](https://go.dev/play/p/wiG6XyBBu49)]

### Keys

Creates an array of the map keys.

```go
keys := lo.Keys[string, int](map[string]int{"foo": 1, "bar": 2})
// []string{"foo", "bar"}
```

[[play](https://go.dev/play/p/Uu11fHASqrU)]

### Values

Creates an array of the map values.

```go
values := lo.Values[string, int](map[string]int{"foo": 1, "bar": 2})
// []int{1, 2}
```

[[play](https://go.dev/play/p/nnRTQkzQfF6)]

### PickBy

Returns same map type filtered by given predicate.

```go
m := lo.PickBy[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
    return value%2 == 1
})
// map[string]int{"foo": 1, "baz": 3}
```

[[play](https://go.dev/play/p/kdg8GR_QMmf)]

### PickByKeys

Returns same map type filtered by given keys.

```go
m := lo.PickByKeys[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})
// map[string]int{"foo": 1, "baz": 3}
```

[[play](https://go.dev/play/p/R1imbuci9qU)]

### PickByValues

Returns same map type filtered by given values.

```go
m := lo.PickByValues[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})
// map[string]int{"foo": 1, "baz": 3}
```

[[play](https://go.dev/play/p/1zdzSvbfsJc)]

### OmitBy

Returns same map type filtered by given predicate.

```go
m := lo.OmitBy[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
    return value%2 == 1
})
// map[string]int{"bar": 2}
```

[[play](https://go.dev/play/p/EtBsR43bdsd)]

### OmitByKeys

Returns same map type filtered by given keys.

```go
m := lo.OmitByKeys[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})
// map[string]int{"bar": 2}
```

[[play](https://go.dev/play/p/t1QjCrs-ysk)]

### OmitByValues

Returns same map type filtered by given values.

```go
m := lo.OmitByValues[string, int](map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})
// map[string]int{"bar": 2}
```

[[play](https://go.dev/play/p/9UYZi-hrs8j)]

### Entries (alias: ToPairs)

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

[[play](https://go.dev/play/p/3Dhgx46gawJ)]

### FromEntries (alias: FromPairs)

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

[[play](https://go.dev/play/p/oIr5KHFGCEN)]

### Invert

Creates a map composed of the inverted keys and values. If map contains duplicate values, subsequent values overwrite property assignments of previous values.

```go
m1 := lo.Invert[string, int](map[string]int{"a": 1, "b": 2})
// map[int]string{1: "a", 2: "b"}

m2 := lo.Invert[string, int](map[string]int{"a": 1, "b": 2, "c": 1})
// map[int]string{1: "c", 2: "b"}
```

[[play](https://go.dev/play/p/rFQ4rak6iA1)]

### Assign

Merges multiple maps from left to right.

```go
mergedMaps := lo.Assign[string, int](
    map[string]int{"a": 1, "b": 2},
    map[string]int{"b": 3, "c": 4},
)
// map[string]int{"a": 1, "b": 3, "c": 4}
```

[[play](https://go.dev/play/p/VhwfJOyxf5o)]

### MapKeys

Manipulates a map keys and transforms it to a map of another type.

```go
m2 := lo.MapKeys[int, int, string](map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, v int) string {
    return strconv.FormatInt(int64(v), 10)
})
// map[string]int{"1": 1, "2": 2, "3": 3, "4": 4}
```

[[play](https://go.dev/play/p/9_4WPIqOetJ)]

### MapValues

Manipulates a map values and transforms it to a map of another type.

```go
m1 := map[int]int64{1: 1, 2: 2, 3: 3}

m2 := lo.MapValues[int, int64, string](m1, func(x int64, _ int) string {
	return strconv.FormatInt(x, 10)
})
// map[int]string{1: "1", 2: "2", 3: "3"}
```

[[play](https://go.dev/play/p/T_8xAfvcf0W)]

### MapToSlice

Transforms a map into a slice based on specific iteratee.

```go
m := map[int]int64{1: 4, 2: 5, 3: 6}

s := lo.MapToSlice(m, func(k int, v int64) string {
    return fmt.Sprintf("%d_%d", k, v)
})
// []string{"1_4", "2_5", "3_6"}
```

[[play](https://go.dev/play/p/ZuiCZpDt6LD)]

### Range / RangeFrom / RangeWithSteps

Creates an array of numbers (positive and/or negative) progressing from start up to, but not including end.

```go
result := Range(4)
// [0, 1, 2, 3]

result := Range(-4)
// [0, -1, -2, -3]

result := RangeFrom(1, 5)
// [1, 2, 3, 4, 5]

result := RangeFrom[float64](1.0, 5)
// [1.0, 2.0, 3.0, 4.0, 5.0]

result := RangeWithSteps(0, 20, 5)
// [0, 5, 10, 15]

result := RangeWithSteps[float32](-1.0, -4.0, -1.0)
// [-1.0, -2.0, -3.0]

result := RangeWithSteps(1, 4, -1)
// []

result := Range(0)
// []
```

[[play](https://go.dev/play/p/0r6VimXAi9H)]

### Clamp

Clamps number within the inclusive lower and upper bounds.

```go
r1 := lo.Clamp(0, -10, 10)
// 0

r2 := lo.Clamp(-42, -10, 10)
// -10

r3 := lo.Clamp(42, -10, 10)
// 10
```

[[play](https://go.dev/play/p/RU4lJNC2hlI)]

### SumBy

Summarizes the values in a collection using the given return value from the iteration function.
If collection is empty 0 is returned.

```go
strings := []string{"foo", "bar"}
sum := lo.SumBy(strings, func(item string) int {
    return len(item)
})
// 6
```

[[play](https://go.dev/play/p/Dz_a_7jN_ca)]

### Substring

Return part of a string.

```go
sub := lo.Substring("hello", 2, 3)
// "llo"

sub := lo.Substring("hello", -4, 3)
// "ell"

sub := lo.Substring("hello", -2, math.MaxUint)
// "lo"
```

[[play](https://go.dev/play/p/TQlxQi82Lu1)]

### ChunkString

Returns an array of strings split into groups the length of size. If array can't be split evenly, the final chunk will be the remaining elements.

```go
lo.ChunkString("123456", 2)
// []string{"12", "34", "56"}

lo.ChunkString("1234567", 2)
// []string{"12", "34", "56", "7"}

lo.ChunkString("", 2)
// []string{""}

lo.ChunkString("1", 2)
// []string{"1"}
```

[[play](https://go.dev/play/p/__FLTuJVz54)]

### RuneLength

An alias to utf8.RuneCountInString which returns the number of runes in string.

```go
sub := lo.RuneLength("hellô")
// 5

sub := len("hellô")
// 6
```

[[play](https://go.dev/play/p/tuhgW_lWY8l)]

### T2 -> T9

Creates a tuple from a list of values.

```go
tuple1 := lo.T2("x", 1)
// Tuple2[string, int]{A: "x", B: 1}

func example() (string, int) { return "y", 2 }
tuple2 := lo.T2(example())
// Tuple2[string, int]{A: "y", B: 2}
```

[[play](https://go.dev/play/p/IllL3ZO4BQm)]

### Unpack2 -> Unpack9

Returns values contained in tuple.

```go
r1, r2 := lo.Unpack2[string, int](lo.Tuple2[string, int]{"a", 1})
// "a", 1
```

[[play](https://go.dev/play/p/xVP_k0kJ96W)]

### Zip2 -> Zip9

Zip creates a slice of grouped elements, the first of which contains the first elements of the given arrays, the second of which contains the second elements of the given arrays, and so on.

When collections have different size, the Tuple attributes are filled with zero value.

```go
tuples := lo.Zip2[string, int]([]string{"a", "b"}, []int{1, 2})
// []Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}}
```

[[play](https://go.dev/play/p/jujaA6GaJTp)]

### Unzip2 -> Unzip9

Unzip accepts an array of grouped elements and creates an array regrouping the elements to their pre-zip configuration.

```go
a, b := lo.Unzip2[string, int]([]Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}})
// []string{"a", "b"}
// []int{1, 2}
```

[[play](https://go.dev/play/p/ciHugugvaAW)]

### ChannelDispatcher

Distributes messages from input channels into N child channels. Close events are propagated to children.

Underlying channels can have a fixed buffer capacity or be unbuffered when cap is 0.

```go
ch := make(chan int, 42)
for i := 0; i <= 10; i++ {
    ch <- i
}

children := lo.ChannelDispatcher(ch, 5, 10, DispatchingStrategyRoundRobin[int])
// []<-chan int{...}

consumer := func(c <-chan int) {
    for {
        msg, ok := <-c
        if !ok {
            println("closed")
            break
        }

        println(msg)
    }
}

for i := range children {
    go consumer(children[i])
}
```

Many distributions strategies are available:

- [lo.DispatchingStrategyRoundRobin](./channel.go): Distributes messages in a rotating sequential manner.
- [lo.DispatchingStrategyRandom](./channel.go): Distributes messages in a random manner.
- [lo.DispatchingStrategyWeightedRandom](./channel.go): Distributes messages in a weighted manner.
- [lo.DispatchingStrategyFirst](./channel.go): Distributes messages in the first non-full channel.
- [lo.DispatchingStrategyLeast](./channel.go): Distributes messages in the emptiest channel.
- [lo.DispatchingStrategyMost](./channel.go): Distributes to the fulliest channel.

Some strategies bring fallback, in order to favor non-blocking behaviors. See implementations.

For custom strategies, just implement the `lo.DispatchingStrategy` prototype:

```go
type DispatchingStrategy[T any] func(message T, messageIndex uint64, channels []<-chan T) int
```

Eg:

```go
type Message struct {
    TenantID uuid.UUID
}

func hash(id uuid.UUID) int {
	h := fnv.New32a()
	h.Write([]byte(id.String()))
	return int(h.Sum32())
}

// Routes messages per TenantID.
customStrategy := func(message pubsub.AMQPSubMessage, messageIndex uint64, channels []<-chan pubsub.AMQPSubMessage) int {
    destination := hash(message.TenantID) % len(channels)

    // check if channel is full
    if len(channels[destination]) < cap(channels[destination]) {
        return destination
    }

    // fallback when child channel is full
    return utils.DispatchingStrategyRoundRobin(message, uint64(destination), channels)
}

children := lo.ChannelDispatcher(ch, 5, 10, customStrategy)
...
```

### SliceToChannel

Returns a read-only channels of collection elements. Channel is closed after last element. Channel capacity can be customized.

```go
list := []int{1, 2, 3, 4, 5}

for v := range lo.SliceToChannel(2, list) {
    println(v)		
}
// prints 1, then 2, then 3, then 4, then 5
```

### Generator

Implements the generator design pattern. Channel is closed after last element. Channel capacity can be customized.

```go
generator := func(yield func(int)) {
    yield(1)
    yield(2)
    yield(3)
}

for v := range lo.Generator(2, generator) {
    println(v)
}
// prints 1, then 2, then 3
```

### Batch

Creates a slice of n elements from a channel. Returns the slice, the slice length, the read time and the channel status (opened/closed).

```go
ch := lo.SliceToChannel(2, []int{1, 2, 3, 4, 5})

items1, length1, duration1, ok1 := lo.Batch(ch, 3)
// []int{1, 2, 3}, 3, 0s, true
items2, length2, duration2, ok2 := lo.Batch(ch, 3)
// []int{4, 5}, 2, 0s, false
```

Example: RabbitMQ consumer 👇

```go
ch := readFromQueue()

for {
    // read 1k items
    items, length, _, ok := lo.Batch(ch, 1000)

    // do batching stuff

    if !ok {
        break
    }
}
```

### BatchWithTimeout

Creates a slice of n elements from a channel, with timeout. Returns the slice, the slice length, the read time and the channel status (opened/closed).

```go
generator := func(yield func(int)) {
    for i := 0; i < 5; i++ {
        yield(i)
        time.Sleep(35*time.Millisecond)
    }
}

ch := lo.Generator(0, generator)

items1, length1, duration1, ok1 := lo.BatchWithTimeout(ch, 3, 100*time.Millisecond)
// []int{1, 2}, 2, 100ms, true
items2, length2, duration2, ok2 := lo.BatchWithTimeout(ch, 3, 100*time.Millisecond)
// []int{3, 4, 5}, 3, 75ms, true
items3, length3, duration2, ok3 := lo.BatchWithTimeout(ch, 3, 100*time.Millisecond)
// []int{}, 0, 10ms, false
```

Example: RabbitMQ consumer 👇

```go
ch := readFromQueue()

for {
    // read 1k items
    // wait up to 1 second
    items, length, _, ok := lo.BatchWithTimeout(ch, 1000, 1*time.Second)

    // do batching stuff

    if !ok {
        break
    }
}
```

Example: Multithreaded RabbitMQ consumer 👇

```go
ch := readFromQueue()

// 5 workers
// prefetch 1k messages per worker
children := lo.ChannelDispatcher(ch, 5, 1000, DispatchingStrategyFirst[int])

consumer := func(c <-chan int) {
    for {
        // read 1k items
        // wait up to 1 second
        items, length, _, ok := lo.BatchWithTimeout(ch, 1000, 1*time.Second)

        // do batching stuff

        if !ok {
            break
        }
    }
}

for i := range children {
    go consumer(children[i])
}
```

### Contains

Returns true if an element is present in a collection.

```go
present := lo.Contains[int]([]int{0, 1, 2, 3, 4, 5}, 5)
// true
```

### ContainsBy

Returns true if the predicate function returns `true`.

```go
present := lo.ContainsBy[int]([]int{0, 1, 2, 3, 4, 5}, func(x int) bool {
    return x == 3
})
// true
```

### Every

Returns true if all elements of a subset are contained into a collection or if the subset is empty.

```go
ok := lo.Every[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// true

ok := lo.Every[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// false
```

### EveryBy

Returns true if the predicate returns true for all of the elements in the collection or if the collection is empty.

```go
b := EveryBy[int]([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 5
})
// true
```

### Some

Returns true if at least 1 element of a subset is contained into a collection.
If the subset is empty Some returns false.

```go
ok := lo.Some[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// true

ok := lo.Some[int]([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// false
```

### SomeBy

Returns true if the predicate returns true for any of the elements in the collection.
If the collection is empty SomeBy returns false.

```go
b := SomeBy[int]([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 3
})
// true
```

### None

Returns true if no element of a subset are contained into a collection or if the subset is empty.

```go
b := None[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// false
b := None[int]([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// true
```

### NoneBy

Returns true if the predicate returns true for none of the elements in the collection or if the collection is empty.

```go
b := NoneBy[int]([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 0
})
// true
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

left, right := lo.Difference[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
// []int{}, []int{}
```

### Union

Returns all distinct elements from both collections. Result will not change the order of elements relatively.

```go
union := lo.Union[int]([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 10})
// []int{0, 1, 2, 3, 4, 5, 10}
```

### Without

Returns slice excluding all given values.

```go
subset := lo.Without[int]([]int{0, 2, 10}, 2)
// []int{0, 10}

subset := lo.Without[int]([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
// []int{10}
```

### WithoutEmpty

Returns slice excluding empty values.

```go
subset := lo.WithoutEmpty[int]([]int{0, 2, 10})
// []int{2, 10}
```

### IndexOf

Returns the index at which the first occurrence of a value is found in an array or return -1 if the value cannot be found.

```go
found := lo.IndexOf[int]([]int{0, 1, 2, 1, 2, 3}, 2)
// 2

notFound := lo.IndexOf[int]([]int{0, 1, 2, 1, 2, 3}, 6)
// -1
```

### LastIndexOf

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

### FindIndexOf

FindIndexOf searches an element in a slice based on a predicate and returns the index and true. It returns -1 and false if the element is not found.

```go
str, index, ok := lo.FindIndexOf[string]([]string{"a", "b", "a", "b"}, func(i string) bool {
    return i == "b"
})
// "b", 1, true

str, index, ok := lo.FindIndexOf[string]([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", -1, false
```

### FindLastIndexOf

FindLastIndexOf searches an element in a slice based on a predicate and returns the index and true. It returns -1 and false if the element is not found.

```go
str, index, ok := lo.FindLastIndexOf[string]([]string{"a", "b", "a", "b"}, func(i string) bool {
    return i == "b"
})
// "b", 4, true

str, index, ok := lo.FindLastIndexOf[string]([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", -1, false
```

### FindKey

Returns the key of the first value matching.

```go
result1, ok1 := lo.FindKey(map[string]int{"foo": 1, "bar": 2, "baz": 3}, 2)
// "bar", true

result2, ok2 := lo.FindKey(map[string]int{"foo": 1, "bar": 2, "baz": 3}, 42)
// "", false

type test struct {
    foobar string
}
result3, ok3 := lo.FindKey(map[string]test{"foo": test{"foo"}, "bar": test{"bar"}, "baz": test{"baz"}}, test{"foo"})
// "foo", true
```

### FindKeyBy

Returns the key of the first element predicate returns truthy for.

```go
result1, ok1 := lo.FindKeyBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(k string, v int) bool {
    return k == "foo"
})
// "foo", true

result2, ok2 := lo.FindKeyBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(k string, v int) bool {
    return false
})
// "", false
```

### FindUniques

Returns a slice with all the unique elements of the collection. The order of result values is determined by the order they occur in the array.

```go
uniqueValues := lo.FindUniques[int]([]int{1, 2, 2, 1, 2, 3})
// []int{3}
```

### FindUniquesBy

Returns a slice with all the unique elements of the collection. The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is invoked for each element in array to generate the criterion by which uniqueness is computed.

```go
uniqueValues := lo.FindUniquesBy[int, int]([]int{3, 4, 5, 6, 7}, func(i int) int {
    return i%3
})
// []int{5}
```

### FindDuplicates

Returns a slice with the first occurence of each duplicated elements of the collection. The order of result values is determined by the order they occur in the array.

```go
duplicatedValues := lo.FindDuplicates[int]([]int{1, 2, 2, 1, 2, 3})
// []int{1, 2}
```

### FindDuplicatesBy

Returns a slice with the first occurence of each duplicated elements of the collection. The order of result values is determined by the order they occur in the array. It accepts `iteratee` which is invoked for each element in array to generate the criterion by which uniqueness is computed.

```go
duplicatedValues := lo.FindDuplicatesBy[int, int]([]int{3, 4, 5, 6, 7}, func(i int) int {
    return i%3
})
// []int{3, 4}
```

### Min

Search the minimum value of a collection.

```go
min := lo.Min[int]([]int{1, 2, 3})
// 1

min := lo.Min[int]([]int{})
// 0
```

### MinBy

Search the minimum value of a collection using the given comparison function.
If several values of the collection are equal to the smallest value, returns the first such value.

```go
min := lo.MinBy[string]([]string{"s1", "string2", "s3"}, func(item string, min string) bool {
    return len(item) < len(min)
})
// "s1"

min := lo.MinBy[string]([]string{}, func(item string, min string) bool {
    return len(item) < len(min)
})
// ""
```

### Max

Search the maximum value of a collection.

```go
max := lo.Max[int]([]int{1, 2, 3})
// 3

max := lo.Max[int]([]int{})
// 0
```

### MaxBy

Search the maximum value of a collection using the given comparison function.
If several values of the collection are equal to the greatest value, returns the first such value.

```go
max := lo.MaxBy[string]([]string{"string1", "s2", "string3"}, func(item string, max string) bool {
    return len(item) > len(max)
})
// "string1"

max := lo.MaxBy[string]([]string{}, func(item string, max string) bool {
    return len(item) > len(max)
})
// ""
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

[[play](https://go.dev/play/p/t-D7WBL44h2)]

### TernaryF

A 1 line if/else statement whose options are functions.

```go
result := lo.TernaryF[string](true, func() string { return "a" }, func() string { return "b" })
// "a"

result := lo.TernaryF[string](false, func() string { return "a" }, func() string { return "b" })
// "b"
```

Useful to avoid nil-pointer dereferencing in intializations, or avoid running unnecessary code

```go
var s *string

someStr := TernaryF[string](s == nil, func() string { return uuid.New().String() }, func() string { return *s })
// ef782193-c30c-4e2e-a7ae-f8ab5e125e02
```

[[play](https://go.dev/play/p/AO4VW20JoqM)]

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

Using callbacks:

```go
result := lo.IfF[int](true, func () int {
        return 1
    }).
    ElseIfF(false, func () int {
        return 2
    }).
    ElseF(func () int {
        return 3
    })
// 1
```

Mixed:

```go
result := lo.IfF[int](true, func () int {
        return 1
    }).
    Else(42)
// 1
```

[[play](https://go.dev/play/p/WSw3ApMxhyW)]

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

Mixed:

```go
result := lo.Switch[int, string](1).
    CaseF(1, func() string {
        return "1"
    }).
    Default("42")
// "1"
```

[[play](https://go.dev/play/p/TGbKUMAeRUd)]

### ToPtr

Returns a pointer copy of value.

```go
ptr := lo.ToPtr[string]("hello world")
// *string{"hello world"}
```

### FromPtr

Returns the pointer value or empty.

```go
str := "hello world"
value := lo.FromPtr[string](&str)
// "hello world"

value := lo.FromPtr[string](nil)
// ""
```

### FromPtrOr

Returns the pointer value or the fallback value.

```go
str := "hello world"
value := lo.FromPtrOr[string](&str, "empty")
// "hello world"

value := lo.FromPtrOr[string](nil, "empty")
// "empty"
```

### ToSlicePtr

Returns a slice of pointer copy of value.

```go
ptr := lo.ToSlicePtr[string]([]string{"hello", "world"})
// []*string{"hello", "world"}
```

### ToAnySlice

Returns a slice with all elements mapped to `any` type.

```go
elements := lo.ToAnySlice[int]([]int{1, 5, 1})
// []any{1, 5, 1}
```

### FromAnySlice

Returns an `any` slice with all elements mapped to a type. Returns false in case of type conversion failure.

```go
elements, ok := lo.FromAnySlice[string]([]any{"foobar", 42})
// []string{}, false

elements, ok := lo.FromAnySlice[string]([]any{"foobar", "42"})
// []string{"foobar", "42"}, true
```

### Empty

Returns an empty value.

```go
lo.Empty[int]()
// 0
lo.Empty[string]()
// ""
lo.Empty[bool]()
// false
```

### IsEmpty

Returns true if argument is a zero value.

```go
lo.IsEmpty[int](0)
// true
lo.IsEmpty[int](42)
// false

lo.IsEmpty[string]("")
// true
lo.IsEmpty[bool]("foobar")
// false

type test struct {
    foobar string
}

lo.IsEmpty[test](test{foobar: ""})
// true
lo.IsEmpty[test](test{foobar: "foobar"})
// false
```

### IsNotEmpty

Returns true if argument is a zero value.

```go
lo.IsNotEmpty[int](0)
// false
lo.IsNotEmpty[int](42)
// true

lo.IsNotEmpty[string]("")
// false
lo.IsNotEmpty[bool]("foobar")
// true

type test struct {
    foobar string
}

lo.IsNotEmpty[test](test{foobar: ""})
// false
lo.IsNotEmpty[test](test{foobar: "foobar"})
// true
```

### Coalesce

Returns the first non-empty arguments. Arguments must be comparable.

```go
result, ok := lo.Coalesce(0, 1, 2, 3)
// 1 true

result, ok := lo.Coalesce("")
// "" false

var nilStr *string
str := "foobar"
result, ok := lo.Coalesce[*string](nil, nilStr, &str)
// &"foobar" true
```

### Partial

Returns new function that, when called, has its first argument set to the provided value.

```go
add := func(x, y int) int { return x + y }
f := lo.Partial(add, 5)

f(10)
// 15

f(42)
// 47
```

### Attempt

Invokes a function N times until it returns valid output. Returning either the caught error or nil. When first argument is less than `1`, the function runs until a successful response is returned.

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

For more advanced retry strategies (delay, exponential backoff...), please take a look on [cenkalti/backoff](https://github.com/cenkalti/backoff).

[[play](https://go.dev/play/p/3ggJZ2ZKcMj)]

### AttemptWithDelay

Invokes a function N times until it returns valid output, with a pause between each call. Returning either the caught error or nil.

When first argument is less than `1`, the function runs until a successful response is returned.

```go
iter, duration, err := lo.AttemptWithDelay(5, 2*time.Second, func(i int, duration time.Duration) error {
    if i == 2 {
        return nil
    }

    return fmt.Errorf("failed")
})
// 3
// ~ 4 seconds
// nil
```

For more advanced retry strategies (delay, exponential backoff...), please take a look on [cenkalti/backoff](https://github.com/cenkalti/backoff).

[[play](https://go.dev/play/p/tVs6CygC7m1)]

### Debounce

`NewDebounce` creates a debounced instance that delays invoking functions given until after wait milliseconds have elapsed, until `cancel` is called.

```go
f := func() {
    println("Called once after 100ms when debounce stopped invoking!")
}

debounce, cancel := lo.NewDebounce(100 * time.Millisecond, f)
for j := 0; j < 10; j++ {
    debounce()
}

time.Sleep(1 * time.Second)
cancel()
```

[[play](https://go.dev/play/p/mz32VMK2nqe)]

### Synchronize

Wraps the underlying callback in a mutex. It receives an optional mutex.

```go
s := lo.Synchronize()

for i := 0; i < 10; i++ {
    go s.Do(func () {
        println("will be called sequentially")
    })
}
```

It is equivalent to:

```go
mu := sync.Mutex{}

func foobar() {
    mu.Lock()
    defer mu.Unlock()

    // ...
}
```

### Async

Executes a function in a goroutine and returns the result in a channel.

```go
ch := lo.Async(func() error { time.Sleep(10 * time.Second); return nil })
// chan error (nil)
```

### Async{0->6}

Executes a function in a goroutine and returns the result in a channel.
For function with multiple return values, the results will be returned as a tuple inside the channel.
For function without return, struct{} will be returned in the channel.

```go
ch := lo.Async0(func() { time.Sleep(10 * time.Second) })
// chan struct{}

ch := lo.Async1(func() int {
  time.Sleep(10 * time.Second);
  return 42
})
// chan int (42)

ch := lo.Async2(func() (int, string) {
  time.Sleep(10 * time.Second);
  return 42, "Hello"
})
// chan lo.Tuple2[int, string] ({42, "Hello"})
```

### Validate

Helper function that creates an error when a condition is not met.

```go
slice := []string{"a"}
val := lo.Validate(len(slice) == 0, "Slice should be empty but contains %v", slice)
// error("Slice should be empty but contains [a]")

slice := []string{}
val := lo.Validate(len(slice) == 0, "Slice should be empty but contains %v", slice)
// nil
```

[[play](https://go.dev/play/p/vPyh51XpCBt)]

### Must

Wraps a function call to panics if second argument is `error` or `false`, returns the value otherwise.

```go
val := lo.Must(time.Parse("2006-01-02", "2022-01-15"))
// 2022-01-15

val := lo.Must(time.Parse("2006-01-02", "bad-value"))
// panics
```

[[play](https://go.dev/play/p/TMoWrRp3DyC)]

### Must{0->6}

Must\* has the same behavior than Must, but returns multiple values.

```go
func example0() (error)
func example1() (int, error)
func example2() (int, string, error)
func example3() (int, string, time.Date, error)
func example4() (int, string, time.Date, bool, error)
func example5() (int, string, time.Date, bool, float64, error)
func example6() (int, string, time.Date, bool, float64, byte, error)

lo.Must0(example0())
val1 := lo.Must1(example1())    // alias to Must
val1, val2 := lo.Must2(example2())
val1, val2, val3 := lo.Must3(example3())
val1, val2, val3, val4 := lo.Must4(example4())
val1, val2, val3, val4, val5 := lo.Must5(example5())
val1, val2, val3, val4, val5, val6 := lo.Must6(example6())
```

You can wrap functions like `func (...) (..., ok bool)`.

```go
// math.Signbit(float64) bool
lo.Must0(math.Signbit(v))

// bytes.Cut([]byte,[]byte) ([]byte, []byte, bool)
before, after := lo.Must2(bytes.Cut(s, sep))
```

You can give context to the panic message by adding some printf-like arguments.

```go
val, ok := lo.Find(myString, func(i string) bool {
    return i == requiredChar
})
lo.Must0(ok, "'%s' must always contain '%s'", myString, requiredChar)

list := []int{0, 1, 2}
item := 5
lo.Must0(lo.Contains[int](list, item), "'%s' must always contain '%s'", list, item)
...
```

[[play](https://go.dev/play/p/TMoWrRp3DyC)]

### Try

Calls the function and return false in case of error and on panic.

```go
ok := lo.Try(func() error {
    panic("error")
    return nil
})
// false

ok := lo.Try(func() error {
    return nil
})
// true

ok := lo.Try(func() error {
    return fmt.Errorf("error")
})
// false
```

[[play](https://go.dev/play/p/mTyyWUvn9u4)]

### Try{0->6}

The same behavior than `Try`, but callback returns 2 variables.

```go
ok := lo.Try2(func() (string, error) {
    panic("error")
    return "", nil
})
// false
```

[[play](https://go.dev/play/p/mTyyWUvn9u4)]

### TryOr

Calls the function and return a default value in case of error and on panic.

```go
str, ok := lo.TryOr(func() (string, error) {
    panic("error")
    return "hello", nil
}, "world")
// world
// false

ok := lo.TryOr(func() error {
    return "hello", nil
}, "world")
// hello
// true

ok := lo.TryOr(func() error {
    return "hello", fmt.Errorf("error")
}, "world")
// world
// false
```

[[play](https://go.dev/play/p/B4F7Wg2Zh9X)]

### TryOr{0->6}

The same behavior than `TryOr`, but callback returns 2 variables.

```go
str, nbr, ok := lo.TryOr2(func() (string, int, error) {
    panic("error")
    return "hello", 42, nil
}, "world", 21)
// world
// 21
// false
```

[[play](https://go.dev/play/p/B4F7Wg2Zh9X)]

### TryWithErrorValue

The same behavior than `Try`, but also returns value passed to panic.

```go
err, ok := lo.TryWithErrorValue(func() error {
    panic("error")
    return nil
})
// "error", false
```

[[play](https://go.dev/play/p/Kc7afQIT2Fs)]

### TryCatch

The same behavior than `Try`, but calls the catch function in case of error.

```go
caught := false

ok := lo.TryCatch(func() error {
    panic("error")
    return nil
}, func() {
    caught = true
})
// false
// caught == true
```

[[play](https://go.dev/play/p/PnOON-EqBiU)]

### TryCatchWithErrorValue

The same behavior than `TryWithErrorValue`, but calls the catch function in case of error.

```go
caught := false

ok := lo.TryCatchWithErrorValue(func() error {
    panic("error")
    return nil
}, func(val any) {
    caught = val == "error"
})
// false
// caught == true
```

[[play](https://go.dev/play/p/8Pc9gwX_GZO)]

### ErrorsAs

A shortcut for:

```go
err := doSomething()

var rateLimitErr *RateLimitError
if ok := errors.As(err, &rateLimitErr); ok {
    // retry later
}
```

1 line `lo` helper:

```go
err := doSomething()

if rateLimitErr, ok := lo.ErrorsAs[*RateLimitError](err); ok {
    // retry later
}
```

[[play](https://go.dev/play/p/8wk5rH8UfrE)]

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

- `lo.Map` is way faster (x7) than `go-funk`, a reflection-based Map implementation.
- `lo.Map` have the same allocation profile than `for`.
- `lo.Map` is 4% slower than `for`.
- `lop.Map` is slower than `lo.Map` because it implies more memory allocation and locks. `lop.Map` will be useful for long-running callbacks, such as i/o bound processing.
- `for` beats other implementations for memory and CPU.

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/lo)
- Fix [open issues](https://github.com/samber/lo/issues) or request new features

Don't hesitate ;)

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
