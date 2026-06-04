
# lo - Iterate over slices, maps, channels...

[![tag](https://img.shields.io/github/tag/samber/lo.svg)](https://github.com/samber/lo/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/lo?status.svg)](https://pkg.go.dev/github.com/samber/lo)
![Build Status](https://github.com/samber/lo/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/lo)](https://goreportcard.com/report/github.com/samber/lo)
[![Coverage](https://img.shields.io/codecov/c/github/samber/lo)](https://codecov.io/gh/samber/lo)
[![Contributors](https://img.shields.io/github/contributors/samber/lo)](https://github.com/samber/lo/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/lo)](./LICENSE)

✨ **`samber/lo` is a Lodash-style Go library based on Go 1.18+ Generics.**

A utility library based on Go 1.18+ generics that makes it easier to work with slices, maps, strings, channels, and functions. It provides dozens of handy methods to simplify common coding tasks and improve code readability. It may look like [Lodash](https://github.com/lodash/lodash) in some aspects.

5 to 10 helpers may overlap with those from the Go standard library, in packages `slices` and `maps`. I feel this library is legitimate and offers many more valuable abstractions.

**See also:**

- [samber/ro](https://github.com/samber/ro): Reactive Programming for Go: declarative and composable API for event-driven applications
- [samber/do](https://github.com/samber/do): A dependency injection toolkit based on Go 1.18+ Generics
- [samber/mo](https://github.com/samber/mo): Monads based on Go 1.18+ Generics (Option, Result, Either...)

What makes it different from **samber/ro**?
- lo: synchronous helpers across finite sequences (maps, slices...)
- ro: processing of infinite data streams for event-driven scenarios

----

<div align="center">
  <sup><b>💖 Sponsored by:</b></sup>
  <br>
  <a href="https://www.dbos.dev/?utm_campaign=gh-smbr">
    <div>
	  <img width="200" alt="dbos" src="https://github.com/user-attachments/assets/d583cb62-7735-4d3c-beb7-e6ef1a5faf49" />
    </div>
    <div>
      DBOS - Durable workflow orchestration library for Go
    </div>
  </a>
</div>

----

**Why this name?**

I wanted a **short name**, similar to "Lodash", and no Go package uses this name.

![lo](docs/static/img/logo-full.png)

## 🚀 Install

```sh
go get github.com/samber/lo@v1
```

This library is v1 and follows SemVer strictly.

No breaking changes will be made to exported APIs before v2.0.0, except for experimental packages under `exp/`.

This library has no dependencies outside the Go standard library.

## 💡 Usage

You can import `lo` using:

```go
import (
    "github.com/samber/lo"
    lop "github.com/samber/lo/parallel"
    lom "github.com/samber/lo/mutable"
    loi "github.com/samber/lo/it"
)
```

Then use one of the helpers below:

```go
names := lo.Uniq([]string{"Samuel", "John", "Samuel"})
// []string{"Samuel", "John"}
```

### Tips for lazy developers

I cannot recommend it, but in case you are too lazy for repeating `lo.` everywhere, you can import the entire library into the namespace.

```go
import (
    . "github.com/samber/lo"
)
```

I take no responsibility for this junk. 😁 💩

## 🤠 Spec

GoDoc: [godoc.org/github.com/samber/lo](https://godoc.org/github.com/samber/lo)

Documentation: [lo.samber.dev](https://lo.samber.dev/docs/about)

Supported helpers for slices:

- [Filter](#filter)
- [Map](#map)
- [UniqMap](#uniqmap)
- [FilterMap](#filtermap)
- [FlatMap](#flatmap)
- [Reduce](#reduce)
- [ReduceRight](#reduceright)
- [ForEach](#foreach)
- [ForEachWhile](#foreachwhile)
- [Times](#times)
- [Uniq](#uniq)
- [UniqBy](#uniqby)
- [GroupBy](#groupby)
- [GroupByMap](#groupbymap)
- [Chunk](#chunk)
- [Window](#window)
- [Sliding](#sliding)
- [PartitionBy](#partitionby)
- [Flatten](#flatten)
- [Concat](#concat)
- [Interleave](#interleave)
- [Shuffle](#shuffle)
- [Reverse](#reverse)
- [Fill](#fill)
- [Repeat](#repeat)
- [RepeatBy](#repeatby)
- [KeyBy](#keyby)
- [SliceToMap / Associate](#slicetomap-alias-associate)
- [FilterSliceToMap](#filterslicetomap)
- [Keyify](#keyify)
- [Take](#take)
- [TakeWhile](#takewhile)
- [TakeFilter](#takefilter)
- [Drop](#drop)
- [DropRight](#dropright)
- [DropWhile](#dropwhile)
- [DropRightWhile](#droprightwhile)
- [DropByIndex](#DropByIndex)
- [Reject](#reject)
- [RejectMap](#rejectmap)
- [FilterReject](#filterreject)
- [Count](#count)
- [CountBy](#countby)
- [CountValues](#countvalues)
- [CountValuesBy](#countvaluesby)
- [Subset](#subset)
- [Slice](#slice)
- [Replace](#replace)
- [ReplaceAll](#replaceall)
- [Clone](#clone)
- [Compact](#compact)
- [IsSorted](#issorted)
- [IsSortedBy](#issortedby)
- [Splice](#Splice)
- [Cut](#Cut)
- [CutPrefix](#CutPrefix)
- [CutSuffix](#CutSuffix)
- [Trim](#Trim)
- [TrimLeft](#TrimLeft)
- [TrimPrefix](#TrimPrefix)
- [TrimRight](#TrimRight)
- [TrimSuffix](#TrimSuffix)

Supported helpers for maps:

- [Keys](#keys)
- [UniqKeys](#uniqkeys)
- [HasKey](#haskey)
- [ValueOr](#valueor)
- [Values](#values)
- [UniqValues](#uniqvalues)
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
- [ChunkEntries](#chunkentries)
- [MapKeys](#mapkeys)
- [MapValues](#mapvalues)
- [MapEntries](#mapentries)
- [MapToSlice](#maptoslice)
- [FilterMapToSlice](#FilterMapToSlice)
- [FilterKeys](#FilterKeys)
- [FilterValues](#FilterValues)

Supported math helpers:

- [Range / RangeFrom / RangeWithSteps](#range--rangefrom--rangewithsteps)
- [Clamp](#clamp)
- [Sum](#sum)
- [SumBy](#sumby)
- [Product](#product)
- [ProductBy](#productby)
- [Mean](#mean)
- [MeanBy](#meanby)
- [Mode](#mode)

Supported helpers for strings:

- [RandomString](#randomstring)
- [Substring](#substring)
- [ChunkString](#chunkstring)
- [RuneLength](#runelength)
- [PascalCase](#pascalcase)
- [CamelCase](#camelcase)
- [KebabCase](#kebabcase)
- [SnakeCase](#snakecase)
- [Words](#words)
- [Capitalize](#capitalize)
- [Ellipsis](#ellipsis)

Supported helpers for tuples:

- [T2 -> T9](#t2---t9)
- [Unpack2 -> Unpack9](#unpack2---unpack9)
- [Zip2 -> Zip9](#zip2---zip9)
- [ZipBy2 -> ZipBy9](#zipby2---zipby9)
- [Unzip2 -> Unzip9](#unzip2---unzip9)
- [UnzipBy2 -> UnzipBy9](#unzipby2---unzipby9)
- [CrossJoin2 -> CrossJoin2](#crossjoin2---crossjoin9)
- [CrossJoinBy2 -> CrossJoinBy2](#crossjoinby2---crossjoinby9)

Supported helpers for time and duration:

- [Duration](#duration)
- [Duration0 -> Duration10](#duration0---duration10)

Supported helpers for channels:

- [ChannelDispatcher](#channeldispatcher)
- [SliceToChannel](#slicetochannel)
- [ChannelToSlice](#channeltoslice)
- [Generator](#generator)
- [Buffer](#buffer)
- [BufferWithContext](#bufferwithcontext)
- [BufferWithTimeout](#bufferwithtimeout)
- [FanIn](#fanin)
- [FanOut](#fanout)

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
- [IntersectBy](#intersectby)
- [Difference](#difference)
- [Union](#union)
- [Without](#without)
- [WithoutBy](#withoutby)
- [WithoutEmpty](#withoutempty)
- [WithoutNth](#withoutnth)
- [ElementsMatch](#ElementsMatch)
- [ElementsMatchBy](#ElementsMatchBy)

Supported search helpers:

- [IndexOf](#indexof)
- [LastIndexOf](#lastindexof)
- [HasPrefix](#hasprefix)
- [HasSuffix](#hassuffix)
- [Find](#find)
- [FindIndexOf](#findindexof)
- [FindLastIndexOf](#findlastindexof)
- [FindOrElse](#findorelse)
- [FindKey](#findkey)
- [FindKeyBy](#findkeyby)
- [FindUniques](#finduniques)
- [FindUniquesBy](#finduniquesby)
- [FindDuplicates](#findduplicates)
- [FindDuplicatesBy](#findduplicatesby)
- [Min](#min)
- [MinIndex](#minindex)
- [MinBy](#minby)
- [MinIndexBy](#minindexby)
- [Earliest](#earliest)
- [EarliestBy](#earliestby)
- [Max](#max)
- [MaxIndex](#maxindex)
- [MaxBy](#maxby)
- [MaxIndexBy](#maxindexby)
- [Latest](#latest)
- [LatestBy](#latestby)
- [First](#first)
- [FirstOrEmpty](#FirstOrEmpty)
- [FirstOr](#FirstOr)
- [Last](#last)
- [LastOrEmpty](#LastOrEmpty)
- [LastOr](#LastOr)
- [Nth](#nth)
- [NthOr](#nthor)
- [NthOrEmpty](#nthorempty)
- [Sample](#sample)
- [SampleBy](#sampleby)
- [Samples](#samples)
- [SamplesBy](#samplesby)

Conditional helpers:

- [Ternary](#ternary)
- [TernaryF](#ternaryf)
- [If / ElseIf / Else](#if--elseif--else)
- [Switch / Case / Default](#switch--case--default)

Type manipulation helpers:

- [IsNil](#isnil)
- [IsNotNil](#isnotnil)
- [ToPtr](#toptr)
- [Nil](#nil)
- [EmptyableToPtr](#emptyabletoptr)
- [FromPtr](#fromptr)
- [FromPtrOr](#fromptror)
- [ToSlicePtr](#tosliceptr)
- [FromSlicePtr](#fromsliceptr)
- [FromSlicePtrOr](#fromsliceptror)
- [ToAnySlice](#toanyslice)
- [FromAnySlice](#fromanyslice)
- [Empty](#empty)
- [IsEmpty](#isempty)
- [IsNotEmpty](#isnotempty)
- [Coalesce](#coalesce)
- [CoalesceOrEmpty](#coalesceorempty)
- [CoalesceSlice](#coalesceslice)
- [CoalesceSliceOrEmpty](#coalescesliceorempty)
- [CoalesceMap](#coalescemap)
- [CoalesceMapOrEmpty](#coalescemaporempty)

Function helpers:

- [Partial](#partial)
- [Partial2 -> Partial5](#partial2---partial5)

Concurrency helpers:

- [Attempt](#attempt)
- [AttemptWhile](#attemptwhile)
- [AttemptWithDelay](#attemptwithdelay)
- [AttemptWhileWithDelay](#attemptwhilewithdelay)
- [Debounce](#debounce)
- [DebounceBy](#debounceby)
- [Throttle](#throttle)
- [ThrottleWithCount](#throttle)
- [ThrottleBy](#throttle)
- [ThrottleByWithCount](#throttle)
- [Synchronize](#synchronize)
- [Async](#async)
- [Async{0->6}](#async0-6)
- [Transaction](#transaction)
- [WaitFor](#waitfor)
- [WaitForWithContext](#waitforwithcontext)

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
- [Assert](#assert)
- [Assertf](#assertf)

Constraints:

- Clonable

### Filter

Iterates over a collection and returns a slice of all the elements the predicate function returns `true` for.

```go
even := lo.Filter([]int{1, 2, 3, 4}, func(x int, index int) bool {
    return x%2 == 0
})
// []int{2, 4}
```

```go
// Use FilterErr when the predicate can return an error
even, err := lo.FilterErr([]int{1, 2, 3, 4}, func(x int, _ int) (bool, error) {
    if x == 3 {
        return false, fmt.Errorf("number 3 is not allowed")
    }
    return x%2 == 0, nil
})
// []int(nil), error("number 3 is not allowed")
```

[[play](https://go.dev/play/p/Apjg3WeSi7K)]

Mutable: like `lo.Filter()`, but the slice is updated in place.

```go
import lom "github.com/samber/lo/mutable"

list := []int{1, 2, 3, 4}
newList := lom.Filter(list, func(x int) bool {
    return x%2 == 0
})

list
// []int{2, 4, 3, 4}

newList
// []int{2, 4}
```

### Map

Manipulates a slice of one type and transforms it into a slice of another type:

```go
import "github.com/samber/lo"

lo.Map([]int64{1, 2, 3, 4}, func(x int64, index int) string {
    return strconv.FormatInt(x, 10)
})
// []string{"1", "2", "3", "4"}
```

```go
// Use MapErr when the transform function can return an error
result, err := lo.MapErr([]int{1, 2, 3, 4}, func(x int, _ int) (string, error) {
    if x == 3 {
        return "", fmt.Errorf("number 3 is not allowed")
    }
    return strconv.Itoa(x), nil
})
// []string(nil), error("number 3 is not allowed")
```

[[play](https://go.dev/play/p/OkPcYAhBo0D)]

Parallel processing: like `lo.Map()`, but the transform function is called in a goroutine. Results are returned in the same order.

```go
import lop "github.com/samber/lo/parallel"

lop.Map([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
    return strconv.FormatInt(x, 10)
})
// []string{"1", "2", "3", "4"}
```

[[play](https://go.dev/play/p/sCJaB3quRMC)]

Mutable: like `lo.Map()`, but the slice is updated in place.

```go
import lom "github.com/samber/lo/mutable"

list := []int{1, 2, 3, 4}
lom.Map(list, func(x int) int {
    return x*2
})
// []int{2, 4, 6, 8}
```

[[play](https://go.dev/play/p/0jY3Z0B7O_5)]

### UniqMap

Manipulates a slice and transforms it to a slice of another type with unique values.

```go
type User struct {
    Name string
    Age  int
}
users := []User{{Name: "Alex", Age: 10}, {Name: "Alex", Age: 12}, {Name: "Bob", Age: 11}, {Name: "Alice", Age: 20}}

names := lo.UniqMap(users, func(u User, index int) string {
    return u.Name
})
// []string{"Alex", "Bob", "Alice"}
```

[[play](https://go.dev/play/p/fygzLBhvUdB)]

### FilterMap

Returns a slice obtained after both filtering and mapping using the given callback function.

The callback function should return two values: the result of the mapping operation and whether the result element should be included or not.

```go
matching := lo.FilterMap([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (string, bool) {
    if strings.HasSuffix(x, "pu") {
        return "xpu", true
    }
    return "", false
})
// []string{"xpu", "xpu"}
```

[[play](https://go.dev/play/p/-AuYXfy7opz)]

### FlatMap

Manipulates a slice and transforms and flattens it to a slice of another type. The transform function can either return a slice or a `nil`, and in the `nil` case no value is added to the final slice.

```go
lo.FlatMap([]int64{0, 1, 2}, func(x int64, _ int) []string {
    return []string{
        strconv.FormatInt(x, 10),
        strconv.FormatInt(x, 10),
    }
})
// []string{"0", "0", "1", "1", "2", "2"}
```

```go
// Use FlatMapErr when the transform function can return an error
result, err := lo.FlatMapErr([]int64{0, 1, 2, 3}, func(x int64, _ int) ([]string, error) {
    if x == 2 {
        return nil, fmt.Errorf("number 2 is not allowed")
    }
    return []string{strconv.FormatInt(x, 10), strconv.FormatInt(x, 10)}, nil
})
// []string(nil), error("number 2 is not allowed")
```

[[play](https://go.dev/play/p/YSoYmQTA8-U)]

### Reduce

Reduces a collection to a single value. The value is calculated by accumulating the result of running each element in the collection through an accumulator function. Each successive invocation is supplied with the return value returned by the previous call.

```go
sum := lo.Reduce([]int{1, 2, 3, 4}, func(agg int, item int, _ int) int {
    return agg + item
}, 0)
// 10
```

```go
// Use ReduceErr when the accumulator function can return an error
result, err := lo.ReduceErr([]int{1, 2, 3, 4}, func(agg int, item int, _ int) (int, error) {
    if item == 3 {
        return 0, fmt.Errorf("number 3 is not allowed")
    }
    return agg + item, nil
}, 0)
// 0, error("number 3 is not allowed")
```

[[play](https://go.dev/play/p/R4UHXZNaaUG)]

### ReduceRight

Like `lo.Reduce` except that it iterates over elements of collection from right to left.

```go
result := lo.ReduceRight([][]int{{0, 1}, {2, 3}, {4, 5}}, func(agg []int, item []int, _ int) []int {
    return append(agg, item...)
}, []int{})
// []int{4, 5, 2, 3, 0, 1}
```

```go
// Use ReduceRightErr when the accumulator function can return an error
result, err := lo.ReduceRightErr([]int{1, 2, 3, 4}, func(agg int, item int, _ int) (int, error) {
    if item == 2 {
        return 0, fmt.Errorf("number 2 is not allowed")
    }
    return agg + item, nil
}, 0)
// 0, error("number 2 is not allowed")
```

[[play](https://go.dev/play/p/Fq3W70l7wXF)]

### ForEach

Iterates over elements of a collection and invokes the function over each element.

```go
import "github.com/samber/lo"

lo.ForEach([]string{"hello", "world"}, func(x string, _ int) {
    println(x)
})
// prints "hello\nworld\n"
```

[[play](https://go.dev/play/p/oofyiUPRf8t)]

Parallel processing: like `lo.ForEach()`, but the callback is called as a goroutine.

```go
import lop "github.com/samber/lo/parallel"

lop.ForEach([]string{"hello", "world"}, func(x string, _ int) {
    println(x)
})
// prints "hello\nworld\n" or "world\nhello\n"
```

### ForEachWhile

Iterates over collection elements and invokes iteratee for each element collection return value decide to continue or break, like do while().

```go
list := []int64{1, 2, -42, 4}

lo.ForEachWhile(list, func(x int64, _ int) bool {
	if x < 0 {
		return false
	}
	fmt.Println(x)
	return true
})
// 1
// 2
```

[[play](https://go.dev/play/p/QnLGt35tnow)]

### Times

Times invokes the iteratee n times, returning a slice of the results of each invocation. The iteratee is invoked with index as argument.

```go
import "github.com/samber/lo"

lo.Times(3, func(i int) string {
    return strconv.FormatInt(int64(i), 10)
})
// []string{"0", "1", "2"}
```

[[play](https://go.dev/play/p/vgQj3Glr6lT)]

Parallel processing: like `lo.Times()`, but callback is called in goroutine.

```go
import lop "github.com/samber/lo/parallel"

lop.Times(3, func(i int) string {
    return strconv.FormatInt(int64(i), 10)
})
// []string{"0", "1", "2"}
```

### Uniq

Returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept. The order of result values is determined by the order they occur in the slice.

```go
uniqValues := lo.Uniq([]int{1, 2, 2, 1})
// []int{1, 2}
```

[[play](https://go.dev/play/p/DTzbeXZ6iEN)]

### UniqBy

Returns a duplicate-free version of a slice, in which only the first occurrence of each element is kept. The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is invoked for each element in the slice to generate the criterion by which uniqueness is computed.

```go
uniqValues := lo.UniqBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// []int{0, 1, 2}
```

```go
// Use UniqByErr when the iteratee function can return an error
result, err := lo.UniqByErr([]int{0, 1, 2, 3, 4, 5}, func(i int) (int, error) {
    if i == 3 {
        return 0, fmt.Errorf("number 3 is not allowed")
    }
    return i % 3, nil
})
// []int(nil), error("number 3 is not allowed")
```

[[play](https://go.dev/play/p/g42Z3QSb53u)]

### GroupBy

Returns an object composed of keys generated from the results of running each element of collection through iteratee.

```go
import lo "github.com/samber/lo"

groups := lo.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// map[int][]int{0: []int{0, 3}, 1: []int{1, 4}, 2: []int{2, 5}}
```

```go
// Use GroupByErr when the iteratee function can return an error
result, err := lo.GroupByErr([]int{0, 1, 2, 3, 4, 5}, func(i int) (int, error) {
    if i == 3 {
        return 0, fmt.Errorf("number 3 is not allowed")
    }
    return i % 3, nil
})
// map[int][]int(nil), error("number 3 is not allowed")
```

[[play](https://go.dev/play/p/XnQBd_v6brd)]

Parallel processing: like `lo.GroupBy()`, but callback is called in goroutine.

```go
import lop "github.com/samber/lo/parallel"

lop.GroupBy([]int{0, 1, 2, 3, 4, 5}, func(i int) int {
    return i%3
})
// map[int][]int{0: []int{0, 3}, 1: []int{1, 4}, 2: []int{2, 5}}
```

### GroupByMap

Returns an object composed of keys generated from the results of running each element of collection through iteratee.

```go
import lo "github.com/samber/lo"

groups := lo.GroupByMap([]int{0, 1, 2, 3, 4, 5}, func(i int) (int, int) {
    return i%3, i*2
})
// map[int][]int{0: []int{0, 6}, 1: []int{2, 8}, 2: []int{4, 10}}
```

```go
// Use GroupByMapErr when the transform function can return an error
result, err := lo.GroupByMapErr([]int{0, 1, 2, 3, 4, 5}, func(i int) (int, int, error) {
    if i == 3 {
        return 0, 0, fmt.Errorf("number 3 is not allowed")
    }
    return i % 3, i * 2, nil
})
// map[int][]int(nil), error("number 3 is not allowed")
```

[[play](https://go.dev/play/p/iMeruQ3_W80)]

### Chunk

Returns a slice of elements split into groups of length size. If the slice can't be split evenly, the final chunk will be the remaining elements.

```go
lo.Chunk([]int{0, 1, 2, 3, 4, 5}, 2)
// [][]int{{0, 1}, {2, 3}, {4, 5}}

lo.Chunk([]int{0, 1, 2, 3, 4, 5, 6}, 2)
// [][]int{{0, 1}, {2, 3}, {4, 5}, {6}}

lo.Chunk([]int{}, 2)
// [][]int{}

lo.Chunk([]int{0}, 2)
// [][]int{{0}}
```

[[play](https://go.dev/play/p/kEMkFbdu85g)]

### Window

Creates a slice of sliding windows of a given size. Each window shares size-1 elements with the previous one. This is equivalent to `Sliding(collection, size, 1)`.

```go
lo.Window([]int{1, 2, 3, 4, 5}, 3)
// [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}}

lo.Window([]float64{20, 22, 21, 23, 24}, 3)
// [][]float64{{20, 22, 21}, {22, 21, 23}, {21, 23, 24}}
```

### Sliding

Creates a slice of sliding windows of a given size with a given step. If step is equal to size, windows have no common elements (similar to Chunk). If step is less than size, windows share common elements.

```go
// Windows with shared elements (step < size)
lo.Sliding([]int{1, 2, 3, 4, 5, 6}, 3, 1)
// [][]int{{1, 2, 3}, {2, 3, 4}, {3, 4, 5}, {4, 5, 6}}

// Windows with no shared elements (step == size, like Chunk)
lo.Sliding([]int{1, 2, 3, 4, 5, 6}, 3, 3)
// [][]int{{1, 2, 3}, {4, 5, 6}}

// Step > size (skipping elements)
lo.Sliding([]int{1, 2, 3, 4, 5, 6, 7, 8}, 2, 3)
// [][]int{{1, 2}, {4, 5}, {7, 8}}
```

### PartitionBy

Returns a slice of elements split into groups. The order of grouped values is determined by the order they occur in collection. The grouping is generated from the results of running each element of collection through iteratee.

```go
import lo "github.com/samber/lo"

partitions := lo.PartitionBy([]int{-2, -1, 0, 1, 2, 3, 4, 5}, func(x int) string {
    if x < 0 {
        return "negative"
    } else if x%2 == 0 {
        return "even"
    }
    return "odd"
})
// [][]int{{-2, -1}, {0, 2, 4}, {1, 3, 5}}
```

```go
// Use PartitionByErr when the iteratee function can return an error
result, err := lo.PartitionByErr([]int{-2, -1, 0, 1, 2}, func(x int) (string, error) {
    if x == 0 {
        return "", fmt.Errorf("zero is not allowed")
    }
    if x < 0 {
        return "negative", nil
    } else if x%2 == 0 {
        return "even", nil
    }
    return "odd", nil
})
// [][]int(nil), error("zero is not allowed")
```

[[play](https://go.dev/play/p/NfQ_nGjkgXW)]

Parallel processing: like `lo.PartitionBy()`, but callback is called in goroutine. Results are returned in the same order.

```go
import lop "github.com/samber/lo/parallel"

partitions := lop.PartitionBy([]int{-2, -1, 0, 1, 2, 3, 4, 5}, func(x int) string {
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

Returns a slice a single level deep.

```go
flat := lo.Flatten([][]int{{0, 1}, {2, 3, 4, 5}})
// []int{0, 1, 2, 3, 4, 5}
```

[[play](https://go.dev/play/p/rbp9ORaMpjw)]

### Concat

Returns a new slice containing all the elements in collections. Concat conserves the order of the elements.

```go
slice := lo.Concat([]int{1, 2}, []int{3, 4})
// []int{1, 2, 3, 4}

slice := lo.Concat(nil, []int{1, 2}, nil, []int{3, 4}, nil)
// []int{1, 2, 3, 4}

slice := lo.Concat[int]()
// []int{}
```
### Interleave

Round-robin alternating input slices and sequentially appending value at index into result.

```go
interleaved := lo.Interleave([]int{1, 4, 7}, []int{2, 5, 8}, []int{3, 6, 9})
// []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

interleaved := lo.Interleave([]int{1}, []int{2, 5, 8}, []int{3, 6}, []int{4, 7, 9, 10})
// []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
```

[[play](https://go.dev/play/p/-RJkTLQEDVt)]

### Shuffle

Returns a slice of shuffled values. Uses the Fisher-Yates shuffle algorithm.

⚠️ This helper is **mutable**.

```go
import lom "github.com/samber/lo/mutable"

list := []int{0, 1, 2, 3, 4, 5}
lom.Shuffle(list)

list
// []int{1, 4, 0, 3, 5, 2}
```

[[play](https://go.dev/play/p/2xb3WdLjeSJ)]

### Reverse

Reverses a slice so that the first element becomes the last, the second element becomes the second to last, and so on.

⚠️ This helper is **mutable**.

```go
import lom "github.com/samber/lo/mutable"

list := []int{0, 1, 2, 3, 4, 5}
lom.Reverse(list)

list
// []int{5, 4, 3, 2, 1, 0}
```

[[play](https://go.dev/play/p/O-M5pmCRgzV)]

### Fill

Fills elements of a slice with `initial` value.

```go
type foo struct {
  bar string
}

func (f foo) Clone() foo {
  return foo{f.bar}
}

initializedSlice := lo.Fill([]foo{foo{"a"}, foo{"a"}}, foo{"b"})
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

slice := lo.Repeat(2, foo{"a"})
// []foo{foo{"a"}, foo{"a"}}
```

[[play](https://go.dev/play/p/g3uHXbmc3b6)]

### RepeatBy

Builds a slice with values returned by N calls of callback.

```go
slice := lo.RepeatBy(0, func (i int) string {
    return strconv.FormatInt(int64(math.Pow(float64(i), 2)), 10)
})
// []string{}

slice := lo.RepeatBy(5, func(i int) string {
    return strconv.FormatInt(int64(math.Pow(float64(i), 2)), 10)
})
// []string{"0", "1", "4", "9", "16"}
```

[[play](https://go.dev/play/p/ozZLCtX_hNU)]

With error handling:

```go
slice, err := lo.RepeatByErr(5, func(i int) (string, error) {
    if i == 3 {
        return "", fmt.Errorf("index 3 is not allowed")
    }
    return fmt.Sprintf("item-%d", i), nil
})
// []string(nil), error("index 3 is not allowed")
```

### KeyBy

Transforms a slice or a slice of structs to a map based on a pivot callback.

```go
m := lo.KeyBy([]string{"a", "aa", "aaa"}, func(str string) int {
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
result := lo.KeyBy(characters, func(char Character) string {
    return string(rune(char.code))
})
//map[a:{dir:left code:97} d:{dir:right code:100}]
```

```go
result, err := lo.KeyByErr([]string{"a", "aa", "aaa", ""}, func(str string) (int, error) {
    if str == "" {
        return 0, fmt.Errorf("empty string not allowed")
    }
    return len(str), nil
})
// map[int]string(nil), error("empty string not allowed")
```

[[play](https://go.dev/play/p/mdaClUAT-zZ)]

### SliceToMap (alias: Associate)

Returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
If any of two pairs have the same key the last one gets added to the map.

The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.

```go
in := []*foo{{baz: "apple", bar: 1}, {baz: "banana", bar: 2}}

aMap := lo.SliceToMap(in, func (f *foo) (string, int) {
    return f.baz, f.bar
})
// map[string][int]{ "apple":1, "banana":2 }
```

[[play](https://go.dev/play/p/WHa2CfMO3Lr)]

### FilterSliceToMap

Returns a map containing key-value pairs provided by transform function applied to elements of the given slice.

If any of two pairs have the same key the last one gets added to the map.

The order of keys in returned map is not specified and is not guaranteed to be the same from the original slice.

The third return value of the transform function is a boolean that indicates whether the key-value pair should be included in the map.

```go
list := []string{"a", "aa", "aaa"}

result := lo.FilterSliceToMap(list, func(str string) (string, int, bool) {
    return str, len(str), len(str) > 1
})
// map[string][int]{"aa":2 "aaa":3}
```

[[play](https://go.dev/play/p/2z0rDz2ZSGU)]

### Keyify

Returns a map with each unique element of the slice as a key.

```go
set := lo.Keyify([]int{1, 1, 2, 3, 4})
// map[int]struct{}{1:{}, 2:{}, 3:{}, 4:{}}
```

[[play](https://go.dev/play/p/RYhhM_csqIG)]

### Take

Takes the first n elements from a slice.

```go
l := lo.Take([]int{0, 1, 2, 3, 4, 5}, 3)
// []int{0, 1, 2}

l := lo.Take([]int{0, 1, 2}, 5)
// []int{0, 1, 2}
```

### TakeWhile

Takes elements from the beginning while the predicate returns true.

```go
l := lo.TakeWhile([]int{0, 1, 2, 3, 4, 5}, func(val int) bool {
    return val < 3
})
// []int{0, 1, 2}

l := lo.TakeWhile([]string{"a", "aa", "aaa", "aa"}, func(val string) bool {
    return len(val) <= 2
})
// []string{"a", "aa"}
```

### TakeFilter

Filters elements and takes the first n elements that match the predicate. Equivalent to calling Take(Filter(...)), but more efficient as it stops after finding n matches.

```go
l := lo.TakeFilter([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3, func(val int, index int) bool {
    return val%2 == 0
})
// []int{2, 4, 6}

l := lo.TakeFilter([]string{"a", "aa", "aaa", "aaaa"}, 2, func(val string, index int) bool {
    return len(val) > 1
})
// []string{"aa", "aaa"}
```

### Drop

Drops n elements from the beginning of a slice.

```go
l := lo.Drop([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{2, 3, 4, 5}
```

[[play](https://go.dev/play/p/JswS7vXRJP2)]

### DropRight

Drops n elements from the end of a slice.

```go
l := lo.DropRight([]int{0, 1, 2, 3, 4, 5}, 2)
// []int{0, 1, 2, 3}
```

[[play](https://go.dev/play/p/GG0nXkSJJa3)]

### DropWhile

Drop elements from the beginning of a slice while the predicate returns true.

```go
l := lo.DropWhile([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
    return len(val) <= 2
})
// []string{"aaa", "aa", "aa"}
```

[[play](https://go.dev/play/p/7gBPYw2IK16)]

### DropRightWhile

Drop elements from the end of a slice while the predicate returns true.

```go
l := lo.DropRightWhile([]string{"a", "aa", "aaa", "aa", "aa"}, func(val string) bool {
    return len(val) <= 2
})
// []string{"a", "aa", "aaa"}
```

[[play](https://go.dev/play/p/3-n71oEC0Hz)]

### DropByIndex

Drops elements from a slice by the index. A negative index will drop elements from the end of the slice.

```go
l := lo.DropByIndex([]int{0, 1, 2, 3, 4, 5}, 2, 4, -1)
// []int{0, 1, 3}
```

[[play](https://go.dev/play/p/JswS7vXRJP2)]

### Reject

The opposite of Filter, this method returns the elements of collection that predicate does not return true for.

```go
odd := lo.Reject([]int{1, 2, 3, 4}, func(x int, _ int) bool {
    return x%2 == 0
})
// []int{1, 3}
```

```go
// Use RejectErr when the predicate can return an error
odd, err := lo.RejectErr([]int{1, 2, 3, 4}, func(x int, _ int) (bool, error) {
    if x == 3 {
        return false, fmt.Errorf("number 3 is not allowed")
    }
    return x%2 == 0, nil
})
// []int(nil), error("number 3 is not allowed")
```

[[play](https://go.dev/play/p/YkLMODy1WEL)]

### RejectMap

The opposite of FilterMap, this method returns a slice obtained after both filtering and mapping using the given callback function.

The callback function should return two values:

- the result of the mapping operation and
- whether the result element should be included or not.

```go
items := lo.RejectMap([]int{1, 2, 3, 4}, func(x int, _ int) (int, bool) {
    return x*10, x%2 == 0
})
// []int{10, 30}
```

### FilterReject

Mixes Filter and Reject, this method returns two slices, one for the elements of collection that predicate returns true for and one for the elements that predicate does not return true for.

```go
kept, rejected := lo.FilterReject([]int{1, 2, 3, 4}, func(x int, _ int) bool {
    return x%2 == 0
})
// []int{2, 4}
// []int{1, 3}
```

### Count

Counts the number of elements in the collection that equal value.

```go
count := lo.Count([]int{1, 5, 1}, 1)
// 2
```

[[play](https://go.dev/play/p/Y3FlK54yveC)]

### CountBy

Counts the number of elements in the collection for which predicate is true.

```go
count := lo.CountBy([]int{1, 5, 1}, func(i int) bool {
    return i < 4
})
// 2
```

```go
// Use CountByErr when the predicate can return an error
count, err := lo.CountByErr([]int{1, 5, 1}, func(i int) (bool, error) {
    if i == 5 {
        return false, fmt.Errorf("5 not allowed")
    }
    return i < 4, nil
})
// 0, error("5 not allowed")
```

[[play](https://go.dev/play/p/ByQbNYQQi4X)]

### CountValues

Counts the number of each element in the collection.

```go
lo.CountValues([]int{})
// map[int]int{}

lo.CountValues([]int{1, 2})
// map[int]int{1: 1, 2: 1}

lo.CountValues([]int{1, 2, 2})
// map[int]int{1: 1, 2: 2}

lo.CountValues([]string{"foo", "bar", ""})
// map[string]int{"": 1, "foo": 1, "bar": 1}

lo.CountValues([]string{"foo", "bar", "bar"})
// map[string]int{"foo": 1, "bar": 2}
```

[[play](https://go.dev/play/p/-p-PyLT4dfy)]

### CountValuesBy

Counts the number of each element in the collection. It is equivalent to chaining lo.Map and lo.CountValues.

```go
isEven := func(v int) bool {
    return v%2==0
}

lo.CountValuesBy([]int{}, isEven)
// map[bool]int{}

lo.CountValuesBy([]int{1, 2}, isEven)
// map[bool]int{false: 1, true: 1}

lo.CountValuesBy([]int{1, 2, 2}, isEven)
// map[bool]int{false: 1, true: 2}

length := func(v string) int {
    return len(v)
}

lo.CountValuesBy([]string{"foo", "bar", ""}, length)
// map[int]int{0: 1, 3: 2}

lo.CountValuesBy([]string{"foo", "bar", "bar"}, length)
// map[int]int{3: 3}
```

[[play](https://go.dev/play/p/2U0dG1SnOmS)]

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

### Clone

Returns a shallow copy of the collection.

```go
in := []int{1, 2, 3, 4, 5}
cloned := lo.Clone(in)
// Verify it's a different slice by checking that modifying one doesn't affect the other
in[0] = 99
// cloned is []int{1, 2, 3, 4, 5}
```

[[play](https://go.dev/play/p/hgHmoOIxmuH)]

### Compact

Returns a slice of all non-zero elements.

```go
in := []string{"", "foo", "", "bar", ""}

slice := lo.Compact(in)
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

### IsSortedBy

Checks if a slice is sorted by iteratee.

```go
slice := lo.IsSortedBy([]string{"a", "bb", "ccc"}, func(s string) int {
    return len(s)
})
// true
```

[[play](https://go.dev/play/p/wiG6XyBBu49)]

### Splice

Splice inserts multiple elements at index i. A negative index counts back from the end of the slice. The helper is protected against overflow errors.

```go
result := lo.Splice([]string{"a", "b"}, 1, "1", "2")
// []string{"a", "1", "2", "b"}

// negative
result = lo.Splice([]string{"a", "b"}, -1, "1", "2")
// []string{"a", "1", "2", "b"}

// overflow
result = lo.Splice([]string{"a", "b"}, 42, "1", "2")
// []string{"a", "b", "1", "2"}
```

[[play](https://go.dev/play/p/G5_GhkeSUBA)]

### Cut

Slices collection around the first instance of separator, returning the part of collection before and after separator. The found result reports whether separator appears in collection. If separator does not appear in s, cut returns collection, empty slice of []T, false.

```go
actualLeft, actualRight, result = lo.Cut([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"b", "c", "d"})
// actualLeft: []string{"a"}
// actualRight: []string{"e", "f", "g"}
// result: true

result = lo.Cut([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"z"})
// actualLeft: []string{"a", "b", "c", "d", "e", "f", "g"}
// actualRight: []string{}
// result: false

result = lo.Cut([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"a", "b"})
// actualLeft: []string{}
// actualRight: []string{"c", "d", "e", "f", "g"}
// result: true
```

[[play](https://go.dev/play/p/GiL3qhpIP3f)]

### CutPrefix

Returns collection without the provided leading prefix []T and reports whether it found the prefix. If s doesn't start with prefix, CutPrefix returns collection, false. If prefix is the empty []T, CutPrefix returns collection, true.

```go
actualRight, result = lo.CutPrefix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"a", "b", "c"})
// actualRight: []string{"d", "e", "f", "g"}
// result: true

result = lo.CutPrefix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"b"})
// actualRight: []string{"a", "b", "c", "d", "e", "f", "g"}
// result: false

result = lo.CutPrefix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{})
// actualRight: []string{"a", "b", "c", "d", "e", "f", "g"}
// result: true
```

[[play](https://go.dev/play/p/7Plak4a1ICl)]

### CutSuffix

Returns collection without the provided ending suffix []T and reports whether it found the suffix. If it doesn't end with suffix, CutSuffix returns collection, false. If suffix is the empty []T, CutSuffix returns collection, true.

```go
actualLeft, result = lo.CutSuffix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"f", "g"})
// actualLeft: []string{"a", "b", "c", "d", "e"}
// result: true

actualLeft, result = lo.CutSuffix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{"b"})
// actualLeft: []string{"a", "b", "c", "d", "e", "f", "g"}
// result: false

actualLeft, result = lo.CutSuffix([]string{"a", "b", "c", "d", "e", "f", "g"}, []string{})
// actualLeft: []string{"a", "b", "c", "d", "e", "f", "g"}
// result: true
```

[[play](https://go.dev/play/p/7FKfBFvPTaT)]

### Trim

Removes all the leading and trailing cutset from the collection.

```go
result := lo.Trim([]int{0, 1, 2, 0, 3, 0}, []int{1, 0})
// []int{2, 0, 3}

result := lo.Trim([]string{"hello", "world", " "}, []string{" ", ""})
// []string{"hello", "world"}
```

[[play](https://go.dev/play/p/1an9mxLdRG5)]

### TrimLeft

Removes all the leading cutset from the collection.

```go
result := lo.TrimLeft([]int{0, 1, 2, 0, 3, 0}, []int{1, 0})
// []int{2, 0, 3, 0}

result := lo.TrimLeft([]string{"hello", "world", " "}, []string{" ", ""})
// []string{"hello", "world", " "}
```

[[play](https://go.dev/play/p/74aqfAYLmyi)]

### TrimPrefix

Removes all the leading prefix from the collection.

```go
result := lo.TrimPrefix([]int{1, 2, 1, 2, 3, 1, 2, 4}, []int{1, 2})
// []int{3, 1, 2, 4}

result := lo.TrimPrefix([]string{"hello", "world", "hello", "test"}, []string{"hello"})
// []string{"world", "hello", "test"}
```

[[play](https://go.dev/play/p/SHO6X-YegPg)]

### TrimRight

Removes all the trailing cutset from the collection.

```go
result := lo.TrimRight([]int{0, 1, 2, 0, 3, 0}, []int{0, 3})
// []int{0, 1, 2}

result := lo.TrimRight([]string{"hello", "world", "  "}, []string{" ", ""})
// []string{"hello", "world", ""}
```

[[play](https://go.dev/play/p/MRpAfR6sf0g)]

### TrimSuffix

Removes all the trailing suffix from the collection.

```go
result := lo.TrimSuffix([]int{1, 2, 3, 1, 2, 4, 2, 4, 2, 4}, []int{2, 4})
// []int{1, 2, 3, 1}

result := lo.TrimSuffix([]string{"hello", "world", "hello", "test"}, []string{"test"})
// []string{"hello", "world", "hello"}
```

[[play](https://go.dev/play/p/IjEUrV0iofq)]

### Keys

Creates a slice of the map keys.

Use the UniqKeys variant to deduplicate common keys.

```go
keys := lo.Keys(map[string]int{"foo": 1, "bar": 2})
// []string{"foo", "bar"}

keys := lo.Keys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []string{"foo", "bar", "baz"}

keys := lo.Keys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"bar": 3})
// []string{"foo", "bar", "bar"}
```

[[play](https://go.dev/play/p/Uu11fHASqrU)]

### UniqKeys

Creates a slice of unique map keys.

```go
keys := lo.UniqKeys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []string{"foo", "bar", "baz"}

keys := lo.UniqKeys(map[string]int{"foo": 1, "bar": 2}, map[string]int{"bar": 3})
// []string{"foo", "bar"}
```

[[play](https://go.dev/play/p/TPKAb6ILdHk)]

### HasKey

Returns whether the given key exists.

```go
exists := lo.HasKey(map[string]int{"foo": 1, "bar": 2}, "foo")
// true

exists := lo.HasKey(map[string]int{"foo": 1, "bar": 2}, "baz")
// false
```

[[play](https://go.dev/play/p/aVwubIvECqS)]

### Values

Creates a slice of the map values.

Use the UniqValues variant to deduplicate common values.

```go
values := lo.Values(map[string]int{"foo": 1, "bar": 2})
// []int{1, 2}

values := lo.Values(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []int{1, 2, 3}

values := lo.Values(map[string]int{"foo": 1, "bar": 2}, map[string]int{"bar": 2})
// []int{1, 2, 2}
```

[[play](https://go.dev/play/p/nnRTQkzQfF6)]

### UniqValues

Creates a slice of unique map values.

```go
values := lo.UniqValues(map[string]int{"foo": 1, "bar": 2})
// []int{1, 2}

values := lo.UniqValues(map[string]int{"foo": 1, "bar": 2}, map[string]int{"baz": 3})
// []int{1, 2, 3}

values := lo.UniqValues(map[string]int{"foo": 1, "bar": 2}, map[string]int{"bar": 2})
// []int{1, 2}
```

[[play](https://go.dev/play/p/nf6bXMh7rM3)]

### ValueOr

Returns the value of the given key or the fallback value if the key is not present.

```go
value := lo.ValueOr(map[string]int{"foo": 1, "bar": 2}, "foo", 42)
// 1

value := lo.ValueOr(map[string]int{"foo": 1, "bar": 2}, "baz", 42)
// 42
```

[[play](https://go.dev/play/p/bAq9mHErB4V)]

### PickBy

Returns same map type filtered by given predicate.

```go
m := lo.PickBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
    return value%2 == 1
})
// map[string]int{"foo": 1, "baz": 3}
```

```go
// Use PickByErr when the predicate can return an error
m, err := lo.PickByErr(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) (bool, error) {
    if key == "bar" {
        return false, fmt.Errorf("bar not allowed")
    }
    return value%2 == 1, nil
})
// map[string]int(nil), error("bar not allowed")
```

[[play](https://go.dev/play/p/kdg8GR_QMmf)]

### PickByKeys

Returns same map type filtered by given keys.

```go
m := lo.PickByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})
// map[string]int{"foo": 1, "baz": 3}
```

[[play](https://go.dev/play/p/R1imbuci9qU)]

### PickByValues

Returns same map type filtered by given values.

```go
m := lo.PickByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})
// map[string]int{"foo": 1, "baz": 3}
```

[[play](https://go.dev/play/p/1zdzSvbfsJc)]

### OmitBy

Returns same map type filtered by given predicate.

```go
m := lo.OmitBy(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) bool {
    return value%2 == 1
})
// map[string]int{"bar": 2}
```

```go
// Use OmitByErr when the predicate can return an error
m, err := lo.OmitByErr(map[string]int{"foo": 1, "bar": 2, "baz": 3}, func(key string, value int) (bool, error) {
    if key == "bar" {
        return false, fmt.Errorf("bar not allowed")
    }
    return value%2 == 1, nil
})
// map[string]int(nil), error("bar not allowed")
```

[[play](https://go.dev/play/p/EtBsR43bdsd)]

### OmitByKeys

Returns same map type filtered by given keys.

```go
m := lo.OmitByKeys(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []string{"foo", "baz"})
// map[string]int{"bar": 2}
```

[[play](https://go.dev/play/p/t1QjCrs-ysk)]

### OmitByValues

Returns same map type filtered by given values.

```go
m := lo.OmitByValues(map[string]int{"foo": 1, "bar": 2, "baz": 3}, []int{1, 3})
// map[string]int{"bar": 2}
```

[[play](https://go.dev/play/p/9UYZi-hrs8j)]

### Entries (alias: ToPairs)

Transforms a map into a slice of key/value pairs.

```go
entries := lo.Entries(map[string]int{"foo": 1, "bar": 2})
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

[[play](https://go.dev/play/p/_t4Xe34-Nl5)]

### FromEntries (alias: FromPairs)

Transforms a slice of key/value pairs into a map.

```go
m := lo.FromEntries([]lo.Entry[string, int]{
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
m1 := lo.Invert(map[string]int{"a": 1, "b": 2})
// map[int]string{1: "a", 2: "b"}

m2 := lo.Invert(map[string]int{"a": 1, "b": 2, "c": 1})
// map[int]string{1: "c", 2: "b"}
```

[[play](https://go.dev/play/p/rFQ4rak6iA1)]

### Assign

Merges multiple maps from left to right.

```go
mergedMaps := lo.Assign(
    map[string]int{"a": 1, "b": 2},
    map[string]int{"b": 3, "c": 4},
)
// map[string]int{"a": 1, "b": 3, "c": 4}
```

[[play](https://go.dev/play/p/VhwfJOyxf5o)]

### ChunkEntries

Splits a map into a slice of elements in groups of length equal to its size. If the map cannot be split evenly, the final chunk will contain the remaining elements.

```go
maps := lo.ChunkEntries(
    map[string]int{
        "a": 1,
        "b": 2,
        "c": 3,
        "d": 4,
        "e": 5,
    },
    3,
)
// []map[string]int{
//    {"a": 1, "b": 2, "c": 3},
//    {"d": 4, "e": 5},
// }
```
[[play](https://go.dev/play/p/X_YQL6mmoD-)]

### MapKeys

Manipulates map keys and transforms it to a map of another type.

```go
m2 := lo.MapKeys(map[int]int{1: 1, 2: 2, 3: 3, 4: 4}, func(_ int, v int) string {
    return strconv.FormatInt(int64(v), 10)
})
// map[string]int{"1": 1, "2": 2, "3": 3, "4": 4}
```

```go
// Use MapKeysErr when the iteratee can return an error
m2, err := lo.MapKeysErr(map[int]int{1: 1, 2: 2, 3: 3}, func(_ int, v int) (string, error) {
    if v == 2 {
        return "", fmt.Errorf("even number not allowed")
    }
    return strconv.FormatInt(int64(v), 10), nil
})
// map[string]int(nil), error("even number not allowed")
```

[[play](https://go.dev/play/p/9_4WPIqOetJ)]

### MapValues

Manipulates map values and transforms it to a map of another type.

```go
m1 := map[int]int64{1: 1, 2: 2, 3: 3}

m2 := lo.MapValues(m1, func(x int64, _ int) string {
    return strconv.FormatInt(x, 10)
})
// map[int]string{1: "1", 2: "2", 3: "3"}
```

```go
// Use MapValuesErr when the iteratee can return an error
m1 := map[int]int64{1: 1, 2: 2, 3: 3}
m2, err := lo.MapValuesErr(m1, func(x int64, _ int) (string, error) {
    if x == 2 {
        return "", fmt.Errorf("even number not allowed")
    }
    return strconv.FormatInt(x, 10), nil
})
// map[int]string(nil), error("even number not allowed")
```

[[play](https://go.dev/play/p/T_8xAfvcf0W)]

### MapEntries

Manipulates map entries and transforms it to a map of another type.

```go
in := map[string]int{"foo": 1, "bar": 2}

out := lo.MapEntries(in, func(k string, v int) (int, string) {
    return v,k
})
// map[int]string{1: "foo", 2: "bar"}
```

```go
// Use MapEntriesErr when the iteratee can return an error
in := map[string]int{"foo": 1, "bar": 2, "baz": 3}
out, err := lo.MapEntriesErr(in, func(k string, v int) (int, string, error) {
    if k == "bar" {
        return 0, "", fmt.Errorf("bar not allowed")
    }
    return v, k, nil
})
// map[int]string(nil), error("bar not allowed")
```

[[play](https://go.dev/play/p/VuvNQzxKimT)]

### MapToSlice

Transforms a map into a slice based on specified iteratee.

```go
m := map[int]int64{1: 4, 2: 5, 3: 6}

s := lo.MapToSlice(m, func(k int, v int64) string {
    return fmt.Sprintf("%d_%d", k, v)
})
// []string{"1_4", "2_5", "3_6"}
```

```go
// Use MapToSliceErr when the iteratee can return an error
m := map[int]int64{1: 4, 2: 5, 3: 6}
s, err := lo.MapToSliceErr(m, func(k int, v int64) (string, error) {
    if k == 2 {
        return "", fmt.Errorf("key 2 not allowed")
    }
    return fmt.Sprintf("%d_%d", k, v), nil
})
// []string(nil), error("key 2 not allowed")
```

[[play](https://go.dev/play/p/ZuiCZpDt6LD)]

### FilterMapToSlice

Transforms a map into a slice based on specified iteratee. The iteratee returns a value and a boolean. If the boolean is true, the value is added to the result slice.

If the boolean is false, the value is not added to the result slice. The order of the keys in the input map is not specified and the order of the keys in the output slice is not guaranteed.

```go
kv := map[int]int64{1: 1, 2: 2, 3: 3, 4: 4}

result := lo.FilterMapToSlice(kv, func(k int, v int64) (string, bool) {
    return fmt.Sprintf("%d_%d", k, v), k%2 == 0
})
// []{"2_2", "4_4"}
```

```go
kv := map[int]int64{1: 1, 2: 2, 3: 3, 4: 4}

result, err := lo.FilterMapToSliceErr(kv, func(k int, v int64) (string, bool, error) {
    if k == 3 {
        return "", false, fmt.Errorf("key 3 not allowed")
    }
    return fmt.Sprintf("%d_%d", k, v), k%2 == 0, nil
})
// []string(nil), error("key 3 not allowed")
```

### FilterKeys

Transforms a map into a slice based on predicate returns true for specific elements. It is a mix of `lo.Filter()` and `lo.Keys()`.

```go
kv := map[int]string{1: "foo", 2: "bar", 3: "baz"}

result := FilterKeys(kv, func(k int, v string) bool {
    return v == "foo"
})
// [1]
```

```go
// Use FilterKeysErr when the predicate can return an error
result, err := lo.FilterKeysErr(map[int]string{1: "foo", 2: "bar", 3: "baz"}, func(k int, v string) (bool, error) {
    if k == 3 {
        return false, fmt.Errorf("key 3 not allowed")
    }
    return v == "foo", nil
})
// []int(nil), error("key 3 not allowed")
```

[[play](https://go.dev/play/p/OFlKXlPrBAe)]

### FilterValues

Transforms a map into a slice based on predicate returns true for specific elements. It is a mix of `lo.Filter()` and `lo.Values()`.

```go
kv := map[int]string{1: "foo", 2: "bar", 3: "baz"}

result := FilterValues(kv, func(k int, v string) bool {
    return v == "foo"
})
// ["foo"]
```

```go
// Use FilterValuesErr when the predicate can return an error
result, err := lo.FilterValuesErr(map[int]string{1: "foo", 2: "bar", 3: "baz"}, func(k int, v string) (bool, error) {
    if k == 3 {
        return false, fmt.Errorf("key 3 not allowed")
    }
    return v == "foo", nil
})
// []string(nil), error("key 3 not allowed")
```

[[play](https://go.dev/play/p/YVD5r_h-LX-)]

### Range / RangeFrom / RangeWithSteps

Creates a slice of numbers (positive and/or negative) progressing from start up to, but not including end.

```go
result := lo.Range(4)
// [0, 1, 2, 3]

result := lo.Range(-4)
// [0, -1, -2, -3]

result := lo.RangeFrom(1, 5)
// [1, 2, 3, 4, 5]

result := lo.RangeFrom[float64](1.0, 5)
// [1.0, 2.0, 3.0, 4.0, 5.0]

result := lo.RangeWithSteps(0, 20, 5)
// [0, 5, 10, 15]

result := lo.RangeWithSteps[float32](-1.0, -4.0, -1.0)
// [-1.0, -2.0, -3.0]

result := lo.RangeWithSteps(1, 4, -1)
// []

result := lo.Range(0)
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

### Sum

Sums the values in a collection.

If collection is empty 0 is returned.

```go
list := []int{1, 2, 3, 4, 5}
sum := lo.Sum(list)
// 15
```

[[play](https://go.dev/play/p/upfeJVqs4Bt)]

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

With error handling:

```go
strings := []string{"foo", "bar", "baz"}
sum, err := lo.SumByErr(strings, func(item string) (int, error) {
    if item == "bar" {
        return 0, fmt.Errorf("invalid item: %s", item)
    }
    return len(item), nil
})
// sum: 3, err: invalid item: bar
```

### Product

Calculates the product of the values in a collection.

If collection is empty 0 is returned.

```go
list := []int{1, 2, 3, 4, 5}
product := lo.Product(list)
// 120
```

[[play](https://go.dev/play/p/2_kjM_smtAH)]

### ProductBy

Calculates the product of the values in a collection using the given return value from the iteration function.

If collection is empty 0 is returned.

```go
strings := []string{"foo", "bar"}
product := lo.ProductBy(strings, func(item string) int {
    return len(item)
})
// 9
```

```go
// Use ProductByErr when the transform function can return an error
strings := []string{"foo", "bar", "baz"}
product, err := lo.ProductByErr(strings, func(item string) (int, error) {
    if item == "bar" {
        return 0, fmt.Errorf("bar is not allowed")
    }
    return len(item), nil
})
// 3, error("bar is not allowed")
```

[[play](https://go.dev/play/p/wadzrWr9Aer)]

### Mean

Calculates the mean of a collection of numbers.

If collection is empty 0 is returned.

```go
mean := lo.Mean([]int{2, 3, 4, 5})
// 3

mean := lo.Mean([]float64{2, 3, 4, 5})
// 3.5

mean := lo.Mean([]float64{})
// 0
```

### MeanBy

Calculates the mean of a collection of numbers using the given return value from the iteration function.

If collection is empty 0 is returned.

```go
list := []string{"aa", "bbb", "cccc", "ddddd"}
mapper := func(item string) float64 {
    return float64(len(item))
}

mean := lo.MeanBy(list, mapper)
// 3.5

mean := lo.MeanBy([]float64{}, mapper)
// 0
```

```go
// Use MeanByErr when the transform function can return an error
list := []string{"aa", "bbb", "cccc", "ddddd"}
mean, err := lo.MeanByErr(list, func(item string) (float64, error) {
    if item == "cccc" {
        return 0, fmt.Errorf("cccc is not allowed")
    }
    return float64(len(item)), nil
})
// 0, error("cccc is not allowed")
```

[[play](https://go.dev/play/p/j7TsVwBOZ7P)]

### Mode

Calculates the mode (most frequent value) of a collection of numbers.

If multiple values have the same highest frequency, then multiple values are returned.

If the collection is empty, the zero value of `T[]` is returned.

```go
mode := lo.Mode([]int{2, 2, 3, 4})
// [2]

mode := lo.Mode([]float64{2, 2, 3, 3})
// [2, 3]

mode := lo.Mode([]float64{})
// []

mode := lo.Mode([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
// [1, 2, 3, 4, 5, 6, 7, 8, 9]
```

### RandomString

Returns a random string of the specified length and made of the specified charset.

```go
str := lo.RandomString(5, lo.LettersCharset)
// example: "eIGbt"
```

[[play](https://go.dev/play/p/rRseOQVVum4)]

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

Returns a slice of strings split into groups of length size. If the string can't be split evenly, the final chunk will be the remaining characters.

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

Note: `lo.ChunkString` and `lo.Chunk` functions behave inconsistently for empty input: `lo.ChunkString("", n)` returns `[""]` instead of `[]`. See [#788](https://github.com/samber/lo/issues/788).

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

### PascalCase

Converts string to pascal case.

```go
str := lo.PascalCase("hello_world")
// HelloWorld
```

[[play](https://go.dev/play/p/Dy_V_6DUYhe)]

### CamelCase

Converts string to camel case.

```go
str := lo.CamelCase("hello_world")
// helloWorld
```

[[play](https://go.dev/play/p/Go6aKwUiq59)]

### KebabCase

Converts string to kebab case.

```go
str := lo.KebabCase("helloWorld")
// hello-world
```

[[play](https://go.dev/play/p/96gT_WZnTVP)]

### SnakeCase

Converts string to snake case.

```go
str := lo.SnakeCase("HelloWorld")
// hello_world
```

[[play](https://go.dev/play/p/ziB0V89IeVH)]

### Words

Splits string into a slice of its words.

```go
str := lo.Words("helloWorld")
// []string{"hello", "world"}
```

[[play](https://go.dev/play/p/-f3VIQqiaVw)]

### Capitalize

Converts the first character of string to upper case and the remaining to lower case.

```go
str := lo.Capitalize("heLLO")
// Hello
```

[[play](https://go.dev/play/p/uLTZZQXqnsa)]

### Ellipsis

Trims and truncates a string to a specified length in runes (Unicode code points) and appends an ellipsis if truncated. Multi-byte characters such as emoji or CJK ideographs are never split in the middle.

```go
str := lo.Ellipsis("  Lorem Ipsum  ", 5)
// Lo...

str := lo.Ellipsis("Lorem Ipsum", 100)
// Lorem Ipsum

str := lo.Ellipsis("Lorem Ipsum", 3)
// ...

str := lo.Ellipsis("hello 世界! 你好", 8)
// hello...

str := lo.Ellipsis("🏠🐶🐱🌟", 4)
// 🏠🐶🐱🌟
```

[[play](https://go.dev/play/p/qE93rgqe1TW)]

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

Returns values contained in a tuple.

```go
r1, r2 := lo.Unpack2(lo.Tuple2[string, int]{"a", 1})
// "a", 1
```

Unpack is also available as a method of TupleX.

```go
tuple2 := lo.T2("a", 1)
a, b := tuple2.Unpack()
// "a", 1
```

[[play](https://go.dev/play/p/xVP_k0kJ96W)]

### Zip2 -> Zip9

Zip creates a slice of grouped elements, the first of which contains the first elements of the given slices, the second of which contains the second elements of the given slices, and so on.

When collections are different sizes, the Tuple attributes are filled with zero value.

```go
tuples := lo.Zip2([]string{"a", "b"}, []int{1, 2})
// []Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}}
```

[[play](https://go.dev/play/p/jujaA6GaJTp)]

### ZipBy2 -> ZipBy9

ZipBy creates a slice of transformed elements, the first of which contains the first elements of the given slices, the second of which contains the second elements of the given slices, and so on.

When collections are different sizes, the Tuple attributes are filled with zero value.

```go
items := lo.ZipBy2([]string{"a", "b"}, []int{1, 2}, func(a string, b int) string {
    return fmt.Sprintf("%s-%d", a, b)
})
// []string{"a-1", "b-2"}
```

With error handling:

```go
items, err := lo.ZipByErr2([]string{"a", "b"}, []int{1, 2}, func(a string, b int) (string, error) {
    if b == 2 {
        return "", fmt.Errorf("number 2 is not allowed")
    }
    return fmt.Sprintf("%s-%d", a, b), nil
})
// []string(nil), error("number 2 is not allowed")
```

### Unzip2 -> Unzip9

Unzip accepts a slice of grouped elements and creates a slice regrouping the elements to their pre-zip configuration.

```go
a, b := lo.Unzip2([]Tuple2[string, int]{{A: "a", B: 1}, {A: "b", B: 2}})
// []string{"a", "b"}
// []int{1, 2}
```

[[play](https://go.dev/play/p/ciHugugvaAW)]

### UnzipBy2 -> UnzipBy9

UnzipBy2 iterates over a collection and creates a slice regrouping the elements to their pre-zip configuration.

```go
a, b := lo.UnzipBy2([]string{"hello", "john", "doe"}, func(str string) (string, int) {
    return str, len(str)
})
// []string{"hello", "john", "doe"}
// []int{5, 4, 3}
```

```go
a, b, err := lo.UnzipByErr2([]string{"hello", "error", "world"}, func(str string) (string, int, error) {
    if str == "error" {
        return "", 0, fmt.Errorf("error string not allowed")
    }
    return str, len(str), nil
})
// []string{}
// []int{}
// error string not allowed
```

### CrossJoin2 -> CrossJoin9

Combines every item from one list with every item from others. It is the cartesian product of lists received as arguments. Returns an empty list if a list is empty.

```go
result := lo.CrossJoin2([]string{"hello", "john", "doe"}, []int{1, 2})
// lo.Tuple2{"hello", 1}
// lo.Tuple2{"hello", 2}
// lo.Tuple2{"john", 1}
// lo.Tuple2{"john", 2}
// lo.Tuple2{"doe", 1}
// lo.Tuple2{"doe", 2}
```

### CrossJoinBy2 -> CrossJoinBy9

Combines every item from one list with every item from others. It is the cartesian product of lists received as arguments. The transform function is used to create the output values. Returns an empty list if a list is empty.

```go
result := lo.CrossJoinBy2([]string{"hello", "john", "doe"}, []int{1, 2}, func(a A, b B) string {
    return fmt.Sprintf("%s - %d", a, b)
})
// "hello - 1"
// "hello - 2"
// "john - 1"
// "john - 2"
// "doe - 1"
// "doe - 2"
```

With error handling:

```go
result, err := lo.CrossJoinByErr2([]string{"hello", "john"}, []int{1, 2}, func(a string, b int) (string, error) {
    if a == "john" {
        return "", fmt.Errorf("john not allowed")
    }
    return fmt.Sprintf("%s - %d", a, b), nil
})
// []string(nil), error("john not allowed")
```

### Duration

Returns the time taken to execute a function.

```go
duration := lo.Duration(func() {
    // very long job
})
// 3s
```

[[play](https://go.dev/play/p/HQfbBbAXaFP)]

### Duration0 -> Duration10

Returns the time taken to execute a function.

```go
duration := lo.Duration0(func() {
    // very long job
})
// 3s

err, duration := lo.Duration1(func() error {
    // very long job
    return errors.New("an error")
})
// an error
// 3s

str, nbr, err, duration := lo.Duration3(func() (string, int, error) {
    // very long job
    return "hello", 42, nil
})
// hello
// 42
// nil
// 3s
```

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

[[play](https://go.dev/play/p/UZGu2wVg3J2)]

Many distributions strategies are available:

- [lo.DispatchingStrategyRoundRobin](./channel.go): Distributes messages in a rotating sequential manner.
- [lo.DispatchingStrategyRandom](./channel.go): Distributes messages in a random manner.
- [lo.DispatchingStrategyWeightedRandom](./channel.go): Distributes messages in a weighted manner.
- [lo.DispatchingStrategyFirst](./channel.go): Distributes messages in the first non-full channel.
- [lo.DispatchingStrategyLeast](./channel.go): Distributes messages in the emptiest channel.
- [lo.DispatchingStrategyMost](./channel.go): Distributes to the fullest channel.

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
customStrategy := func(message string, messageIndex uint64, channels []<-chan string) int {
    destination := hash(message) % len(channels)

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

Returns a read-only channel of collection elements. Channel is closed after last element. Channel capacity can be customized.

```go
list := []int{1, 2, 3, 4, 5}

for v := range lo.SliceToChannel(2, list) {
    println(v)
}
// prints 1, then 2, then 3, then 4, then 5
```

[[play](https://go.dev/play/p/lIbSY3QmiEg)]

### ChannelToSlice

Returns a slice built from channel items. Blocks until channel closes.

```go
list := []int{1, 2, 3, 4, 5}
ch := lo.SliceToChannel(2, list)

items := ChannelToSlice(ch)
// []int{1, 2, 3, 4, 5}
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

### Buffer

Creates a slice of n elements from a channel. Returns the slice, the slice length, the read time and the channel status (opened/closed).

```go
ch := lo.SliceToChannel(2, []int{1, 2, 3, 4, 5})

items1, length1, duration1, ok1 := lo.Buffer(ch, 3)
// []int{1, 2, 3}, 3, 0s, true
items2, length2, duration2, ok2 := lo.Buffer(ch, 3)
// []int{4, 5}, 2, 0s, false
```

Example: RabbitMQ consumer 👇

```go
ch := readFromQueue()

for {
    // read 1k items
    items, length, _, ok := lo.Buffer(ch, 1000)

    // do batching stuff

    if !ok {
        break
    }
}
```

### BufferWithContext

Creates a slice of n elements from a channel, with timeout. Returns the slice, the slice length, the read time and the channel status (opened/closed).

```go
ctx, cancel := context.WithCancel(context.TODO())
go func() {
    ch <- 0
    time.Sleep(10*time.Millisecond)
    ch <- 1
    time.Sleep(10*time.Millisecond)
    ch <- 2
    time.Sleep(10*time.Millisecond)
    ch <- 3
    time.Sleep(10*time.Millisecond)
    ch <- 4
    time.Sleep(10*time.Millisecond)
    cancel()
}()

items1, length1, duration1, ok1 := lo.BufferWithContext(ctx, ch, 3)
// []int{0, 1, 2}, 3, 20ms, true
items2, length2, duration2, ok2 := lo.BufferWithContext(ctx, ch, 3)
// []int{3, 4}, 2, 30ms, false
```

### BufferWithTimeout

Creates a slice of n elements from a channel, with timeout. Returns the slice, the slice length, the read time and the channel status (opened/closed).

```go
generator := func(yield func(int)) {
    for i := 0; i < 5; i++ {
        yield(i)
        time.Sleep(35*time.Millisecond)
    }
}

ch := lo.Generator(0, generator)

items1, length1, duration1, ok1 := lo.BufferWithTimeout(ch, 3, 100*time.Millisecond)
// []int{1, 2}, 2, 100ms, true
items2, length2, duration2, ok2 := lo.BufferWithTimeout(ch, 3, 100*time.Millisecond)
// []int{3, 4, 5}, 3, 75ms, true
items3, length3, duration2, ok3 := lo.BufferWithTimeout(ch, 3, 100*time.Millisecond)
// []int{}, 0, 10ms, false
```

Example: RabbitMQ consumer 👇

```go
ch := readFromQueue()

for {
    // read 1k items
    // wait up to 1 second
    items, length, _, ok := lo.BufferWithTimeout(ch, 1000, 1*time.Second)

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
children := lo.ChannelDispatcher(ch, 5, 1000, lo.DispatchingStrategyFirst[int])

consumer := func(c <-chan int) {
    for {
        // read 1k items
        // wait up to 1 second
        items, length, _, ok := lo.BufferWithTimeout(ch, 1000, 1*time.Second)

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

### FanIn

Merge messages from multiple input channels into a single buffered channel. Output messages have no priority. When all upstream channels reach EOF, downstream channel closes.

```go
stream1 := make(chan int, 42)
stream2 := make(chan int, 42)
stream3 := make(chan int, 42)

all := lo.FanIn(100, stream1, stream2, stream3)
// <-chan int
```

### FanOut

Broadcasts all the upstream messages to multiple downstream channels. When upstream channel reaches EOF, downstream channels close. If any downstream channels is full, broadcasting is paused.

```go
stream := make(chan int, 42)

all := lo.FanOut(5, 100, stream)
// [5]<-chan int
```

### Contains

Returns true if an element is present in a collection.

```go
present := lo.Contains([]int{0, 1, 2, 3, 4, 5}, 5)
// true
```

[[play](https://go.dev/play/p/W1EvyqY6t9j)]

### ContainsBy

Returns true if the predicate function returns `true`.

```go
present := lo.ContainsBy([]int{0, 1, 2, 3, 4, 5}, func(x int) bool {
    return x == 3
})
// true
```

### Every

Returns true if all elements of a subset are contained in a collection or if the subset is empty.

```go
ok := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// true

ok := lo.Every([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// false
```

### EveryBy

Returns true if the predicate returns true for all elements in the collection or if the collection is empty.

```go
b := EveryBy([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 5
})
// true
```

[[play](https://go.dev/play/p/dn1-vhHsq9x)]

### Some

Returns true if at least 1 element of a subset is contained in a collection.
If the subset is empty Some returns false.

```go
ok := lo.Some([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// true
```

[[play](https://go.dev/play/p/Lj4ceFkeT9V)]

ok := lo.Some([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// false
```

### SomeBy

Returns true if the predicate returns true for any of the elements in the collection.
If the collection is empty SomeBy returns false.

```go
b := SomeBy([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 3
})
// true
```

### None

Returns true if no element of a subset is contained in a collection or if the subset is empty.

```go
b := None([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// false
b := None([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// true
```

[[play](https://go.dev/play/p/fye7JsmxzPV)]

### NoneBy

Returns true if the predicate returns true for none of the elements in the collection or if the collection is empty.

```go
b := NoneBy([]int{1, 2, 3, 4}, func(x int) bool {
    return x < 0
})
// true
```

[[play](https://go.dev/play/p/O64WZ32H58S)]

### Intersect

Returns the intersection between collections.

```go
result1 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// []int{0, 2}

result2 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// []int{0}

result3 := lo.Intersect([]int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// []int{}

result4 := lo.Intersect([]int{0, 3, 5, 7}, []int{3, 5}, []int{0, 1, 2, 0, 3, 0})
// []int{3}
```

### IntersectBy

Returns the intersection between two collections using a custom key selector function.

```go
transform := func(v int) string {
    return strconv.Itoa(v)
}

result1 := lo.IntersectBy(transform, []int{0, 1, 2, 3, 4, 5}, []int{0, 2})
// []int{0, 2}

result2 := lo.IntersectBy(transform, []int{0, 1, 2, 3, 4, 5}, []int{0, 6})
// []int{0}

result3 := lo.IntersectBy(transform, []int{0, 1, 2, 3, 4, 5}, []int{-1, 6})
// []int{}

result4 := lo.IntersectBy(transform, []int{0, 3, 5, 7}, []int{3, 5}, []int{0, 1, 2, 0, 3, 0})
// []int{3}
```

### Difference

Returns the difference between two collections.

- The first value is the collection of elements absent from list2.
- The second value is the collection of elements absent from list1.

```go
left, right := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 2, 6})
// []int{1, 3, 4, 5}, []int{6}

left, right := lo.Difference([]int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5})
// []int{}, []int{}
```

[[play](https://go.dev/play/p/pKE-JgzqRpz)]

### Union

Returns all distinct elements from given collections. Result will not change the order of elements relatively.

```go
union := lo.Union([]int{0, 1, 2, 3, 4, 5}, []int{0, 2}, []int{0, 10})
// []int{0, 1, 2, 3, 4, 5, 10}
```

### Without

Returns a slice excluding all given values.

```go
subset := lo.Without([]int{0, 2, 10}, 2)
// []int{0, 10}

subset := lo.Without([]int{0, 2, 10}, 0, 1, 2, 3, 4, 5)
// []int{10}
```

### WithoutBy

Filters a slice by excluding elements whose extracted keys match any in the exclude list.

Returns a new slice containing only the elements whose keys are not in the exclude list.

```go
type User struct {
    ID int
    Name string
}

// original users
users := []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
    {ID: 3, Name: "Charlie"},
}

// extract function to get the user ID
getID := func(user User) int {
    return user.ID
}

// exclude users with IDs 2 and 3
excludedIDs := []int{2, 3}

// filtering users
filteredUsers := lo.WithoutBy(users, getID, excludedIDs...)
// []User[{ID: 1, Name: "Alice"}]
```

```go
// Use WithoutByErr when the iteratee can return an error
type struct User {
    ID int
    Name string
}

users := []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
    {ID: 3, Name: "Charlie"},
}

getID := func(user User) (int, error) {
    if user.ID == 2 {
        return 0, fmt.Errorf("Bob not allowed")
    }
    return user.ID, nil
}

filteredUsers, err := lo.WithoutByErr(users, getID, 2, 3)
// []User(nil), error("Bob not allowed")
```

### WithoutEmpty

Returns a slice excluding zero values.

```go
subset := lo.WithoutEmpty([]int{0, 2, 10})
// []int{2, 10}
```

### WithoutNth

Returns a slice excluding the nth value.

```go
subset := lo.WithoutNth([]int{-2, -1, 0, 1, 2}, 3, -42, 1)
// []int{-2, 0, 2}
```

### ElementsMatch

Returns true if lists contain the same set of elements (including empty set).

If there are duplicate elements, the number of occurrences in each list should match.

The order of elements is not checked.

```go
b := lo.ElementsMatch([]int{1, 1, 2}, []int{2, 1, 1})
// true
```

### ElementsMatchBy

Returns true if lists contain the same set of elements' keys (including empty set).

If there are duplicate keys, the number of occurrences in each list should match.

The order of elements is not checked.

```go
b := lo.ElementsMatchBy(
    []someType{a, b},
    []someType{b, a},
    func(item someType) string { return item.ID() },
)
// true
```

### IndexOf

Returns the index at which the first occurrence of a value is found in a slice or -1 if the value cannot be found.

```go
found := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
// 2

notFound := lo.IndexOf([]int{0, 1, 2, 1, 2, 3}, 6)
// -1
```

[[play](https://go.dev/play/p/Eo7W0lvKTky)]

### LastIndexOf

Returns the index at which the last occurrence of a value is found in a slice or -1 if the value cannot be found.

```go
found := lo.LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 2)
// 4

notFound := lo.LastIndexOf([]int{0, 1, 2, 1, 2, 3}, 6)
// -1
```

### HasPrefix

Returns true if the collection has the prefix.

```go
ok := lo.HasPrefix([]int{1, 2, 3, 4}, []int{42})
// false

ok := lo.HasPrefix([]int{1, 2, 3, 4}, []int{1, 2})
// true
```

[[play](https://go.dev/play/p/SrljzVDpMQM)]

### HasSuffix

Returns true if the collection has the suffix.

```go
ok := lo.HasSuffix([]int{1, 2, 3, 4}, []int{42})
// false

ok := lo.HasSuffix([]int{1, 2, 3, 4}, []int{3, 4})
// true
```

[[play](https://go.dev/play/p/bJeLetQNAON)]

### Find

Searches for an element in a slice based on a predicate. Returns element and true if element was found.

```go
str, ok := lo.Find([]string{"a", "b", "c", "d"}, func(i string) bool {
    return i == "b"
})
// "b", true

str, ok := lo.Find([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", false
```

```go
// Use FindErr when the predicate can return an error
str, err := lo.FindErr([]string{"a", "b", "c", "d"}, func(i string) (bool, error) {
    if i == "c" {
        return false, fmt.Errorf("c is not allowed")
    }
    return i == "b", nil
})
// "b", nil

str, err = lo.FindErr([]string{"a", "b", "c"}, func(i string) (bool, error) {
    if i == "b" {
        return false, fmt.Errorf("b is not allowed")
    }
    return i == "b", nil
})
// "", error("b is not allowed")
```

[[play](https://go.dev/play/p/Eo7W0lvKTky)]

### FindIndexOf

FindIndexOf searches for an element in a slice based on a predicate and returns the index and true. Returns -1 and false if the element is not found.

```go
str, index, ok := lo.FindIndexOf([]string{"a", "b", "a", "b"}, func(i string) bool {
    return i == "b"
})
// "b", 1, true

str, index, ok := lo.FindIndexOf([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", -1, false
```

[[play](https://go.dev/play/p/XWSEM4Ic_t0)]

### FindLastIndexOf

FindLastIndexOf searches for the last element in a slice based on a predicate and returns the index and true. Returns -1 and false if the element is not found.

```go
str, index, ok := lo.FindLastIndexOf([]string{"a", "b", "a", "b"}, func(i string) bool {
    return i == "b"
})
// "b", 4, true

str, index, ok := lo.FindLastIndexOf([]string{"foobar"}, func(i string) bool {
    return i == "b"
})
// "", -1, false
```

[[play](https://go.dev/play/p/dPiMRtJ6cUx)]

### FindOrElse

Searches for an element in a slice based on a predicate. Returns the element if found or a given fallback value otherwise.

```go
str := lo.FindOrElse([]string{"a", "b", "c", "d"}, "x", func(i string) bool {
    return i == "b"
})
// "b"

str := lo.FindOrElse([]string{"foobar"}, "x", func(i string) bool {
    return i == "b"
})
// "x"
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

Returns the key of the first element predicate returns true for.

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

Returns a slice with all the elements that appear in the collection only once. The order of result values is determined by the order they occur in the slice.

```go
uniqueValues := lo.FindUniques([]int{1, 2, 2, 1, 2, 3})
// []int{3}
```

### FindUniquesBy

Returns a slice with all the elements that appear in the collection only once. The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is invoked for each element in the slice to generate the criterion by which uniqueness is computed.

```go
uniqueValues := lo.FindUniquesBy([]int{3, 4, 5, 6, 7}, func(i int) int {
    return i%3
})
// []int{5}
```

### FindDuplicates

Returns a slice with the first occurrence of each duplicated element in the collection. The order of result values is determined by the order they occur in the slice.

```go
duplicatedValues := lo.FindDuplicates([]int{1, 2, 2, 1, 2, 3})
// []int{1, 2}
```

### FindDuplicatesBy

Returns a slice with the first occurrence of each duplicated element in the collection. The order of result values is determined by the order they occur in the slice. It accepts `iteratee` which is invoked for each element in the slice to generate the criterion by which uniqueness is computed.

```go
duplicatedValues := lo.FindDuplicatesBy([]int{3, 4, 5, 6, 7}, func(i int) int {
    return i%3
})
// []int{3, 4}
```

With error handling:

```go
duplicatedValues, err := lo.FindDuplicatesByErr([]int{3, 4, 5, 6, 7}, func(i int) (int, error) {
    if i == 5 {
        return 0, fmt.Errorf("number 5 is not allowed")
    }
    return i % 3, nil
})
// []int(nil), error("number 5 is not allowed")
```

### Min

Search the minimum value of a collection.

Returns zero value when the collection is empty.

```go
min := lo.Min([]int{1, 2, 3})
// 1

min := lo.Min([]int{})
// 0

min := lo.Min([]time.Duration{time.Second, time.Hour})
// 1s
```

[[play](https://go.dev/play/p/r6e-Z8JozS8)]

### MinIndex

Search the minimum value of a collection and the index of the minimum value.

Returns (zero value, -1) when the collection is empty.

```go
min, index := lo.MinIndex([]int{1, 2, 3})
// 1, 0

min, index := lo.MinIndex([]int{})
// 0, -1

min, index := lo.MinIndex([]time.Duration{time.Second, time.Hour})
// 1s, 0
```

### MinBy

Search the minimum value of a collection using the given comparison function.

If several values of the collection are equal to the smallest value, returns the first such value.

Returns zero value when the collection is empty.

```go
min := lo.MinBy([]string{"s1", "string2", "s3"}, func(item string, min string) bool {
    return len(item) < len(min)
})
// "s1"

min := lo.MinBy([]string{}, func(item string, min string) bool {
    return len(item) < len(min)
})
// ""
```

```go
// Use MinByErr when the comparison function can return an error
min, err := lo.MinByErr([]string{"s1", "string2", "s3"}, func(item string, min string) (bool, error) {
    if item == "string2" {
        return false, fmt.Errorf("string2 is not allowed")
    }
    return len(item) < len(min), nil
})
// "s1", error("string2 is not allowed")
```

### MinIndexBy

Search the minimum value of a collection using the given comparison function and the index of the minimum value.

If several values of the collection are equal to the smallest value, returns the first such value.

Returns (zero value, -1) when the collection is empty.

```go
min, index := lo.MinIndexBy([]string{"s1", "string2", "s3"}, func(item string, min string) bool {
    return len(item) < len(min)
})
// "s1", 0

min, index := lo.MinIndexBy([]string{}, func(item string, min string) bool {
    return len(item) < len(min)
})
// "", -1
```

```go
min, index, err := lo.MinIndexByErr([]string{"s1", "string2", "s3"}, func(item string, min string) (bool, error) {
    if item == "s2" {
        return false, fmt.Errorf("s2 is not allowed")
    }
    return len(item) < len(min), nil
})
// "s1", 0, error("s2 is not allowed")
```

### Earliest

Search the minimum time.Time of a collection.

Returns zero value when the collection is empty.

```go
earliest := lo.Earliest(time.Now(), time.Time{})
// 0001-01-01 00:00:00 +0000 UTC
```

### EarliestBy

Search the minimum time.Time of a collection using the given iteratee function.

Returns zero value when the collection is empty.

```go
type foo struct {
    bar time.Time
}

earliest := lo.EarliestBy([]foo{{time.Now()}, {}}, func(i foo) time.Time {
    return i.bar
})
// {bar:{2023-04-01 01:02:03 +0000 UTC}}
```

```go
// Use EarliestByErr when the iteratee function can return an error
earliest, err := lo.EarliestByErr([]foo{{time.Now()}, {}}, func(i foo) (time.Time, error) {
    if i.bar.IsZero() {
        return time.Time{}, fmt.Errorf("zero time not allowed")
    }
    return i.bar, nil
})
// {bar:{...}}, error("zero time not allowed")
```

### Max

Search the maximum value of a collection.

Returns zero value when the collection is empty.

```go
max := lo.Max([]int{1, 2, 3})
// 3

max := lo.Max([]int{})
// 0

max := lo.Max([]time.Duration{time.Second, time.Hour})
// 1h
```

### MaxIndex

Search the maximum value of a collection and the index of the maximum value.

Returns (zero value, -1) when the collection is empty.

```go
max, index := lo.MaxIndex([]int{1, 2, 3})
// 3, 2

max, index := lo.MaxIndex([]int{})
// 0, -1

max, index := lo.MaxIndex([]time.Duration{time.Second, time.Hour})
// 1h, 1
```

### MaxBy

Search the maximum value of a collection using the given comparison function.

If several values of the collection are equal to the greatest value, returns the first such value.

Returns zero value when the collection is empty.

```go
max := lo.MaxBy([]string{"string1", "s2", "string3"}, func(item string, max string) bool {
    return len(item) > len(max)
})
// "string1"

max := lo.MaxBy([]string{}, func(item string, max string) bool {
    return len(item) > len(max)
})
// ""
```

```go
// Use MaxByErr when the comparison function can return an error
max, err := lo.MaxByErr([]string{"string1", "s2", "string3"}, func(item string, max string) (bool, error) {
    if item == "s2" {
        return false, fmt.Errorf("s2 is not allowed")
    }
    return len(item) > len(max), nil
})
// "string1", error("s2 is not allowed")
```

[[play](https://go.dev/play/p/JW1qu-ECwF7)]

### MaxIndexBy

Search the maximum value of a collection using the given comparison function and the index of the maximum value.

If several values of the collection are equal to the greatest value, returns the first such value.

Returns (zero value, -1) when the collection is empty.

```go
max, index := lo.MaxIndexBy([]string{"string1", "s2", "string3"}, func(item string, max string) bool {
    return len(item) > len(max)
})
// "string1", 0

max, index := lo.MaxIndexBy([]string{}, func(item string, max string) bool {
    return len(item) > len(max)
})
// "", -1
```

```go
// Use MaxIndexByErr when the comparison function can return an error
max, index, err := lo.MaxIndexByErr([]string{"string1", "s2", "string3"}, func(item string, max string) (bool, error) {
    if item == "s2" {
        return false, fmt.Errorf("s2 is not allowed")
    }
    return len(item) > len(max), nil
})
// "string1", 0, error("s2 is not allowed")
```

[[play](https://go.dev/play/p/uaUszc-c9QK)]

### Latest

Search the maximum time.Time of a collection.

Returns zero value when the collection is empty.

```go
latest := lo.Latest(time.Now(), time.Time{})
// 2023-04-01 01:02:03 +0000 UTC
```

### LatestBy

Search the maximum time.Time of a collection using the given iteratee function.

Returns zero value when the collection is empty.

```go
type foo struct {
    bar time.Time
}

latest := lo.LatestBy([]foo{{time.Now()}, {}}, func(i foo) time.Time {
    return i.bar
})
// {bar:{2023-04-01 01:02:03 +0000 UTC}}
```

```go
// Use LatestByErr when the iteratee function can return an error
result, err := lo.LatestByErr([]foo{{time.Now()}, {}}, func(i foo) (time.Time, error) {
    if i.bar.IsZero() {
        return time.Time{}, fmt.Errorf("zero time not allowed")
    }
    return i.bar, nil
})
// foo{}, error("zero time not allowed")
```

### First

Returns the first element of a collection and check for availability of the first element.

```go
first, ok := lo.First([]int{1, 2, 3})
// 1, true

first, ok := lo.First([]int{})
// 0, false
```

### FirstOrEmpty

Returns the first element of a collection or zero value if empty.

```go
first := lo.FirstOrEmpty([]int{1, 2, 3})
// 1

first := lo.FirstOrEmpty([]int{})
// 0
```

### FirstOr

Returns the first element of a collection or the fallback value if empty.

```go
first := lo.FirstOr([]int{1, 2, 3}, 245)
// 1

first := lo.FirstOr([]int{}, 31)
// 31
```

### Last

Returns the last element of a collection or error if empty.

```go
last, ok := lo.Last([]int{1, 2, 3})
// 3
// true

last, ok := lo.Last([]int{})
// 0
// false
```

### LastOrEmpty

Returns the last element of a collection or zero value if empty.

```go
last := lo.LastOrEmpty([]int{1, 2, 3})
// 3

last := lo.LastOrEmpty([]int{})
// 0
```

### LastOr

Returns the last element of a collection or the fallback value if empty.

```go
last := lo.LastOr([]int{1, 2, 3}, 245)
// 3

last := lo.LastOr([]int{}, 31)
// 31
```

### Nth

Returns the element at index `nth` of collection. If `nth` is negative, the nth element from the end is returned. An error is returned when nth is out of slice bounds.

```go
nth, err := lo.Nth([]int{0, 1, 2, 3}, 2)
// 2

nth, err := lo.Nth([]int{0, 1, 2, 3}, -2)
// 2
```

### NthOr

Returns the element at index `nth` of the collection. If `nth` is negative, it returns the `nth` element from the end. If `nth` is out of slice bounds, it returns the provided fallback value
```go
nth := lo.NthOr([]int{10, 20, 30, 40, 50}, 2, -1)
// 30

nth := lo.NthOr([]int{10, 20, 30, 40, 50}, -1, -1)
// 50

nth := lo.NthOr([]int{10, 20, 30, 40, 50}, 5, -1)
// -1 (fallback value)
```

### NthOrEmpty

Returns the element at index `nth` of the collection. If `nth` is negative, it returns the `nth` element from the end. If `nth` is out of slice bounds, it returns the zero value for the element type (e.g., 0 for integers, "" for strings, etc).
``` go
nth := lo.NthOrEmpty([]int{10, 20, 30, 40, 50}, 2)
// 30

nth := lo.NthOrEmpty([]int{10, 20, 30, 40, 50}, -1)
// 50

nth := lo.NthOrEmpty([]int{10, 20, 30, 40, 50}, 5)
// 0 (zero value for int)

nth := lo.NthOrEmpty([]string{"apple", "banana", "cherry"}, 2)
// "cherry"

nth := lo.NthOrEmpty([]string{"apple", "banana", "cherry"}, 5)
// "" (zero value for string)
```

### Sample

Returns a random item from collection.

```go
lo.Sample([]string{"a", "b", "c"})
// a random string from []string{"a", "b", "c"}

lo.Sample([]string{})
// ""
```

[[play](https://go.dev/play/p/FYA45LcpfM2)]

### SampleBy

Returns a random item from collection, using a given random integer generator.

```go
import "math/rand"

r := rand.New(rand.NewSource(42))
lo.SampleBy([]string{"a", "b", "c"}, r.Intn)
// a random string from []string{"a", "b", "c"}, using a seeded random generator

lo.SampleBy([]string{}, r.Intn)
// ""
```

### Samples

Returns N random unique items from collection.

```go
lo.Samples([]string{"a", "b", "c"}, 3)
// []string{"a", "b", "c"} in random order
```

### SamplesBy

Returns N random unique items from collection, using a given random integer generator.

```go
r := rand.New(rand.NewSource(42))
lo.SamplesBy([]string{"a", "b", "c"}, 3, r.Intn)
// []string{"a", "b", "c"} in random order, using a seeded random generator
```

### Ternary

A single line if/else statement.

```go
result := lo.Ternary(true, "a", "b")
// "a"

result := lo.Ternary(false, "a", "b")
// "b"
```

Take care to avoid dereferencing potentially nil pointers in your A/B expressions, because they are both evaluated. See TernaryF to avoid this problem.

[[play](https://go.dev/play/p/t-D7WBL44h2)]

### TernaryF

A single line if/else statement whose options are functions.

```go
result := lo.TernaryF(true, func() string { return "a" }, func() string { return "b" })
// "a"

result := lo.TernaryF(false, func() string { return "a" }, func() string { return "b" })
// "b"
```

Useful to avoid nil-pointer dereferencing in initializations, or avoid running unnecessary code

```go
var s *string

someStr := TernaryF(s == nil, func() string { return uuid.New().String() }, func() string { return *s })
// ef782193-c30c-4e2e-a7ae-f8ab5e125e02
```

[[play](https://go.dev/play/p/AO4VW20JoqM)]

### If / ElseIf / Else

```go
result := lo.If(true, 1).
    ElseIf(false, 2).
    Else(3)
// 1

result := lo.If(false, 1).
    ElseIf(true, 2).
    Else(3)
// 2

result := lo.If(false, 1).
    ElseIf(false, 2).
    Else(3)
// 3
```

Using callbacks:

```go
result := lo.IfF(true, func () int {
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
result := lo.IfF(true, func () int {
        return 1
    }).
    Else(42)
// 1
```

[[play](https://go.dev/play/p/WSw3ApMxhyW)]

### Switch / Case / Default

```go
result := lo.Switch(1).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "1"

result := lo.Switch(2).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "2"

result := lo.Switch(42).
    Case(1, "1").
    Case(2, "2").
    Default("3")
// "3"
```

Using callbacks:

```go
result := lo.Switch(1).
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
result := lo.Switch(1).
    CaseF(1, func() string {
        return "1"
    }).
    Default("42")
// "1"
```

[[play](https://go.dev/play/p/TGbKUMAeRUd)]

### IsNil

Checks if a value is nil or if it's a reference type with a nil underlying value.

```go
var x int
lo.IsNil(x)
// false

var k struct{}
lo.IsNil(k)
// false

var i *int
lo.IsNil(i)
// true

var ifaceWithNilValue any = (*string)(nil)
lo.IsNil(ifaceWithNilValue)
// true
ifaceWithNilValue == nil
// false
```

### IsNotNil

Checks if a value is not nil or if it's not a reference type with a nil underlying value.

```go
var x int
lo.IsNotNil(x)
// true

var k struct{}
lo.IsNotNil(k)
// true

var i *int
lo.IsNotNil(i)
// false

var ifaceWithNilValue any = (*string)(nil)
lo.IsNotNil(ifaceWithNilValue)
// false
ifaceWithNilValue == nil
// true
```

### ToPtr

Returns a pointer copy of the value.

```go
ptr := lo.ToPtr("hello world")
// *string{"hello world"}
```

[[play](https://go.dev/play/p/P2sD0PMXw4F)]

### Nil

Returns a nil pointer of type.

```go
ptr := lo.Nil[float64]()
// nil
```

### EmptyableToPtr

Returns a pointer copy of value if it's nonzero.
Otherwise, returns nil pointer.

```go
ptr := lo.EmptyableToPtr(nil)
// nil

ptr := lo.EmptyableToPtr("")
// nil

ptr := lo.EmptyableToPtr([]int{})
// *[]int{}

ptr := lo.EmptyableToPtr("hello world")
// *string{"hello world"}
```

### FromPtr

Returns the pointer value or empty.

```go
str := "hello world"
value := lo.FromPtr(&str)
// "hello world"

value := lo.FromPtr(nil)
// ""
```

### FromPtrOr

Returns the pointer value or the fallback value.

```go
str := "hello world"
value := lo.FromPtrOr(&str, "empty")
// "hello world"

value := lo.FromPtrOr(nil, "empty")
// "empty"
```

### ToSlicePtr

Returns a slice of pointers to each value.

```go
ptr := lo.ToSlicePtr([]string{"hello", "world"})
// []*string{"hello", "world"}
```

### FromSlicePtr

Returns a slice with the pointer values.
Returns a zero value in case of a nil pointer element.

```go
str1 := "hello"
str2 := "world"

ptr := lo.FromSlicePtr[string]([]*string{&str1, &str2, nil})
// []string{"hello", "world", ""}

ptr := lo.Compact(
    lo.FromSlicePtr[string]([]*string{&str1, &str2, nil}),
)
// []string{"hello", "world"}
```

### FromSlicePtrOr

Returns a slice with the pointer values or the fallback value.

```go
str1 := "hello"
str2 := "world"

ptr := lo.FromSlicePtrOr([]*string{&str1, nil, &str2}, "fallback value")
// []string{"hello", "fallback value", "world"}
```

[[play](https://go.dev/play/p/CuXGVzo9G65)]

### ToAnySlice

Returns a slice with all elements mapped to `any` type.

```go
elements := lo.ToAnySlice([]int{1, 5, 1})
// []any{1, 5, 1}
```

### FromAnySlice

Returns a slice with all elements mapped to a type. Returns false in case of type conversion failure.

```go
elements, ok := lo.FromAnySlice([]any{"foobar", 42})
// []string{}, false

elements, ok := lo.FromAnySlice([]any{"foobar", "42"})
// []string{"foobar", "42"}, true
```

### Empty

Returns the [zero value](https://go.dev/ref/spec#The_zero_value).

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
lo.IsEmpty(0)
// true
lo.IsEmpty(42)
// false

lo.IsEmpty("")
// true
lo.IsEmpty("foobar")
// false

type test struct {
    foobar string
}

lo.IsEmpty(test{foobar: ""})
// true
lo.IsEmpty(test{foobar: "foobar"})
// false
```

### IsNotEmpty

Returns true if argument is a zero value.

```go
lo.IsNotEmpty(0)
// false
lo.IsNotEmpty(42)
// true

lo.IsNotEmpty("")
// false
lo.IsNotEmpty("foobar")
// true

type test struct {
    foobar string
}

lo.IsNotEmpty(test{foobar: ""})
// false
lo.IsNotEmpty(test{foobar: "foobar"})
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
result, ok := lo.Coalesce(nil, nilStr, &str)
// &"foobar" true
```

### CoalesceOrEmpty

Returns the first non-empty arguments. Arguments must be comparable.

```go
result := lo.CoalesceOrEmpty(0, 1, 2, 3)
// 1

result := lo.CoalesceOrEmpty("")
// ""

var nilStr *string
str := "foobar"
result := lo.CoalesceOrEmpty(nil, nilStr, &str)
// &"foobar"
```

### CoalesceSlice

Returns the first non-zero slice.

```go
result, ok := lo.CoalesceSlice([]int{1, 2, 3}, []int{4, 5, 6})
// [1, 2, 3]
// true

result, ok := lo.CoalesceSlice(nil, []int{})
// []
// true

result, ok := lo.CoalesceSlice([]int(nil))
// []
// false
```

### CoalesceSliceOrEmpty

Returns the first non-zero slice.

```go
result := lo.CoalesceSliceOrEmpty([]int{1, 2, 3}, []int{4, 5, 6})
// [1, 2, 3]

result := lo.CoalesceSliceOrEmpty(nil, []int{})
// []
```

### CoalesceMap

Returns the first non-zero map.

```go
result, ok := lo.CoalesceMap(map[string]int{"1": 1, "2": 2, "3": 3}, map[string]int{"4": 4, "5": 5, "6": 6})
// {"1": 1, "2": 2, "3": 3}
// true

result, ok := lo.CoalesceMap(nil, map[string]int{})
// {}
// true

result, ok := lo.CoalesceMap(map[string]int(nil))
// {}
// false
```

### CoalesceMapOrEmpty

Returns the first non-zero map.

```go
result := lo.CoalesceMapOrEmpty(map[string]int{"1": 1, "2": 2, "3": 3}, map[string]int{"4": 4, "5": 5, "6": 6})
// {"1": 1, "2": 2, "3": 3}

result := lo.CoalesceMapOrEmpty(nil, map[string]int{})
// {}
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

[[play](https://go.dev/play/p/Sy1gAQiQZ3v)]

### Partial2 -> Partial5

Returns new function that, when called, has its first argument set to the provided value.

```go
add := func(x, y, z int) int { return x + y + z }
f := lo.Partial2(add, 42)

f(10, 5)
// 57

f(42, -4)
// 80
```

[[play](https://go.dev/play/p/-xiPjy4JChJ)]

### Attempt

Invokes a function N times until it returns valid output. Returns either the caught error or nil.

When the first argument is less than `1`, the function runs until a successful response is returned.

```go
iter, err := lo.Attempt(42, func(i int) error {
    if i == 5 {
        return nil
    }

    return errors.New("failed")
})
// 6
// nil

iter, err := lo.Attempt(2, func(i int) error {
    if i == 5 {
        return nil
    }

    return errors.New("failed")
})
// 2
// error "failed"

iter, err := lo.Attempt(0, func(i int) error {
    if i < 42 {
        return errors.New("failed")
    }

    return nil
})
// 43
// nil
```

For more advanced retry strategies (delay, exponential backoff...), please take a look at [cenkalti/backoff](https://github.com/cenkalti/backoff).

[[play](https://go.dev/play/p/3ggJZ2ZKcMj)]

### AttemptWithDelay

Invokes a function N times until it returns valid output, with a pause between each call. Returns either the caught error or nil.

When the first argument is less than `1`, the function runs until a successful response is returned.

```go
iter, duration, err := lo.AttemptWithDelay(5, 2*time.Second, func(i int, duration time.Duration) error {
    if i == 2 {
        return nil
    }

    return errors.New("failed")
})
// 3
// ~ 4 seconds
// nil
```

For more advanced retry strategies (delay, exponential backoff...), please take a look at [cenkalti/backoff](https://github.com/cenkalti/backoff).

[[play](https://go.dev/play/p/tVs6CygC7m1)]

### AttemptWhile

Invokes a function N times until it returns valid output. Returns either the caught error or nil, along with a bool value to determine whether the function should be invoked again. It will terminate the invoke immediately if the second return value is false.

When the first argument is less than `1`, the function runs until a successful response is returned.

```go
count1, err1 := lo.AttemptWhile(5, func(i int) (error, bool) {
    err := doMockedHTTPRequest(i)
    if err != nil {
        if errors.Is(err, ErrBadRequest) { // let's assume ErrBadRequest is a critical error that needs to terminate the invoke
            return err, false // flag the second return value as false to terminate the invoke
        }

        return err, true
    }

    return nil, false
})
```

For more advanced retry strategies (delay, exponential backoff...), please take a look at [cenkalti/backoff](https://github.com/cenkalti/backoff).

[[play](https://go.dev/play/p/1VS7HxlYMOG)]

### AttemptWhileWithDelay

Invokes a function N times until it returns valid output, with a pause between each call. Returns either the caught error or nil, along with a bool value to determine whether the function should be invoked again. It will terminate the invoke immediately if the second return value is false.

When the first argument is less than `1`, the function runs until a successful response is returned.

```go
count1, time1, err1 := lo.AttemptWhileWithDelay(5, time.Millisecond, func(i int, d time.Duration) (error, bool) {
    err := doMockedHTTPRequest(i)
    if err != nil {
        if errors.Is(err, ErrBadRequest) { // let's assume ErrBadRequest is a critical error that needs to terminate the invoke
            return err, false // flag the second return value as false to terminate the invoke
        }

        return err, true
    }

    return nil, false
})
```

For more advanced retry strategies (delay, exponential backoff...), please take a look at [cenkalti/backoff](https://github.com/cenkalti/backoff).

[[play](https://go.dev/play/p/mhufUjJfLEF)]

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

### DebounceBy

`NewDebounceBy` creates a debounced instance for each distinct key, that delays invoking functions given until after wait milliseconds have elapsed, until `cancel` is called.

```go
f := func(key string, count int) {
    println(key + ": Called once after 100ms when debounce stopped invoking!")
}

debounce, cancel := lo.NewDebounceBy(100 * time.Millisecond, f)
for j := 0; j < 10; j++ {
    debounce("first key")
    debounce("second key")
}

time.Sleep(1 * time.Second)
cancel("first key")
cancel("second key")
```

[[play](https://go.dev/play/p/d3Vpt6pxhY8)]

### Throttle

Creates a throttled instance that invokes given functions only once in every interval.

This returns 2 functions, First one is throttled function and Second one is a function to reset interval.

```go
f := func() {
	println("Called once in every 100ms")
}

throttle, reset := lo.NewThrottle(100 * time.Millisecond, f)

for j := 0; j < 10; j++ {
	throttle()
	time.Sleep(30 * time.Millisecond)
}

reset()
throttle()
```

`NewThrottleWithCount` is NewThrottle with count limit, throttled function will be invoked count times in every interval.

```go
f := func() {
	println("Called three times in every 100ms")
}

throttle, reset := lo.NewThrottleWithCount(100 * time.Millisecond, f)

for j := 0; j < 10; j++ {
	throttle()
	time.Sleep(30 * time.Millisecond)
}

reset()
throttle()
```

`NewThrottleBy` and `NewThrottleByWithCount` are NewThrottle with sharding key, throttled function will be invoked count times in every interval.

```go
f := func(key string) {
	println(key, "Called three times in every 100ms")
}

throttle, reset := lo.NewThrottleByWithCount(100 * time.Millisecond, f)

for j := 0; j < 10; j++ {
	throttle("foo")
	time.Sleep(30 * time.Millisecond)
}

reset()
throttle()
```

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
For functions with multiple return values, the results will be returned as a tuple inside the channel.
For functions without return, struct{} will be returned in the channel.

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

### Transaction

Implements a Saga pattern.

```go
transaction := NewTransaction().
    Then(
        func(state int) (int, error) {
            fmt.Println("step 1")
            return state + 10, nil
        },
        func(state int) int {
            fmt.Println("rollback 1")
            return state - 10
        },
    ).
    Then(
        func(state int) (int, error) {
            fmt.Println("step 2")
            return state + 15, nil
        },
        func(state int) int {
            fmt.Println("rollback 2")
            return state - 15
        },
    ).
    Then(
        func(state int) (int, error) {
            fmt.Println("step 3")

            if true {
                return state, errors.New("error")
            }

            return state + 42, nil
        },
        func(state int) int {
            fmt.Println("rollback 3")
            return state - 42
        },
    )

_, _ = transaction.Process(-5)

// Output:
// step 1
// step 2
// step 3
// rollback 2
// rollback 1
```

### WaitFor

Runs periodically until a condition is validated.

```go
alwaysTrue := func(i int) bool { return true }
alwaysFalse := func(i int) bool { return false }
laterTrue := func(i int) bool {
    return i > 5
}

iterations, duration, ok := lo.WaitFor(alwaysTrue, 10*time.Millisecond, 2 * time.Millisecond)
// 1
// 1ms
// true

iterations, duration, ok := lo.WaitFor(alwaysFalse, 10*time.Millisecond, time.Millisecond)
// 10
// 10ms
// false

iterations, duration, ok := lo.WaitFor(laterTrue, 10*time.Millisecond, time.Millisecond)
// 7
// 7ms
// true

iterations, duration, ok := lo.WaitFor(laterTrue, 10*time.Millisecond, 5*time.Millisecond)
// 2
// 10ms
// false
```

[[play](https://go.dev/play/p/t_wTDmubbK3)]

### WaitForWithContext

Runs periodically until a condition is validated or context is invalid.

The condition receives also the context, so it can invalidate the process in the condition checker

```go
ctx := context.Background()

alwaysTrue := func(_ context.Context, i int) bool { return true }
alwaysFalse := func(_ context.Context, i int) bool { return false }
laterTrue := func(_ context.Context, i int) bool {
    return i >= 5
}

iterations, duration, ok := lo.WaitForWithContext(ctx, alwaysTrue, 10*time.Millisecond, 2 * time.Millisecond)
// 1
// 1ms
// true

iterations, duration, ok := lo.WaitForWithContext(ctx, alwaysFalse, 10*time.Millisecond, time.Millisecond)
// 10
// 10ms
// false

iterations, duration, ok := lo.WaitForWithContext(ctx, laterTrue, 10*time.Millisecond, time.Millisecond)
// 5
// 5ms
// true

iterations, duration, ok := lo.WaitForWithContext(ctx, laterTrue, 10*time.Millisecond, 5*time.Millisecond)
// 2
// 10ms
// false

expiringCtx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
iterations, duration, ok := lo.WaitForWithContext(expiringCtx, alwaysFalse, 100*time.Millisecond, time.Millisecond)
// 5
// 5.1ms
// false
```

[[play](https://go.dev/play/p/t_wTDmubbK3)]

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

Wraps a function call and panics if second argument is `error` or `false`, returns the value otherwise.

```go
val := lo.Must(time.Parse("2006-01-02", "2022-01-15"))
// 2022-01-15

val := lo.Must(time.Parse("2006-01-02", "bad-value"))
// panics
```

[[play](https://go.dev/play/p/TMoWrRp3DyC)]

### Must{0->6}

Must\* has the same behavior as Must but returns multiple values.

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
lo.Must0(lo.Contains(list, item), "'%s' must always contain '%s'", list, item)
...
```

[[play](https://go.dev/play/p/TMoWrRp3DyC)]

### Try

Calls the function and returns false in case of error and panic.

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
    return errors.New("error")
})
// false
```

[[play](https://go.dev/play/p/mTyyWUvn9u4)]

### Try{0->6}

The same behavior as `Try`, but the callback returns 2 variables.

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

str, ok := lo.TryOr(func() error {
    return "hello", nil
}, "world")
// hello
// true

str, ok := lo.TryOr(func() error {
    return "hello", errors.New("error")
}, "world")
// world
// false
```

[[play](https://go.dev/play/p/B4F7Wg2Zh9X)]

### TryOr{0->6}

The same behavior as `TryOr`, but the callback returns `X` variables.

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

The same behavior as `Try`, but also returns the value passed to panic.

```go
err, ok := lo.TryWithErrorValue(func() error {
    panic("error")
    return nil
})
// "error", false
```

[[play](https://go.dev/play/p/Kc7afQIT2Fs)]

### TryCatch

The same behavior as `Try`, but calls the catch function in case of error.

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

The same behavior as `TryWithErrorValue`, but calls the catch function in case of error.

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

single line `lo` helper:

```go
err := doSomething()

if rateLimitErr, ok := lo.ErrorsAs[*RateLimitError](err); ok {
    // retry later
}
```

[[play](https://go.dev/play/p/8wk5rH8UfrE)]

### Assert

Does nothing when the condition is `true`, otherwise it panics with an optional message.

Think twice before using it, given that [Go intentionally omits assertions from its standard library](https://go.dev/doc/faq#assertions).

```go
age := getUserAge()

lo.Assert(age >= 15)
```

```go
age := getUserAge()

lo.Assert(age >= 15, "user age must be >= 15")
```

[[play](https://go.dev/play/p/Xv8LLKBMNwI)]

### Assertf

Like `Assert`, but with `fmt.Printf`-like formatting.

Think twice before using it, given that [Go intentionally omits assertions from its standard library](https://go.dev/doc/faq#assertions).

```go
age := getUserAge()

lo.Assertf(age >= 15, "user age must be >= 15, got %d", age)
```

[[play](https://go.dev/play/p/TVPEmVcyrdY)]

## 🛩 Benchmark

We executed a simple benchmark with a dead-simple `lo.Map` loop:

See the full implementation [here](./map_benchmark_test.go).

```go
_ = lo.Map[int64](arr, func(x int64, i int) string {
    return strconv.FormatInt(x, 10)
})
```

**Result:**

Here is a comparison between `lo.Map`, `lop.Map`, `go-funk` library and a simple Go `for` loop.

```shell
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
- `lo.Map` has the same allocation profile as `for`.
- `lo.Map` is 4% slower than `for`.
- `lop.Map` is slower than `lo.Map` because it implies more memory allocation and locks. `lop.Map` is useful for long-running callbacks, such as i/o bound processing.
- `for` beats other implementations for memory and CPU.

## 🤝 Contributing

- Ping me on Twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/lo)
- Fix [open issues](https://github.com/samber/lo/issues) or request new features

Don't hesitate ;)

Helper naming: helpers must be self-explanatory and respect standards (other languages, libraries...). Feel free to suggest many names in your contributions.

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/lo)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2022 [Samuel Berthe](https://github.com/samber).

This project is under [MIT](./LICENSE) license.
