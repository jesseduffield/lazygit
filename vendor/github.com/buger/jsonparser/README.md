[![Go Report Card](https://goreportcard.com/badge/github.com/buger/jsonparser)](https://goreportcard.com/report/github.com/buger/jsonparser) ![License](https://img.shields.io/dub/l/vibe-d.svg)
# Alternative JSON parser for Go (10x times faster standard library)

It does not require you to know the structure of the payload (eg. create structs), and allows accessing fields by providing the path to them. It is up to **10 times faster** than standard `encoding/json` package (depending on payload size and usage), **allocates no memory**. See benchmarks below.

## Rationale
Originally I made this for a project that relies on a lot of 3rd party APIs that can be unpredictable and complex.
I love simplicity and prefer to avoid external dependecies. `encoding/json` requires you to know exactly your data structures, or if you prefer to use `map[string]interface{}` instead, it will be very slow and hard to manage.
I investigated what's on the market and found that most libraries are just wrappers around `encoding/json`, there is few options with own parsers (`ffjson`, `easyjson`), but they still requires you to create data structures.


Goal of this project is to push JSON parser to the performance limits and not sacrifice with compliance and developer user experience.

## Example
For the given JSON our goal is to extract the user's full name, number of github followers and avatar.

```go
import "github.com/buger/jsonparser"

...

data := []byte(`{
  "person": {
    "name": {
      "first": "Leonid",
      "last": "Bugaev",
      "fullName": "Leonid Bugaev"
    },
    "github": {
      "handle": "buger",
      "followers": 109
    },
    "avatars": [
      { "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
    ]
  },
  "company": {
    "name": "Acme"
  }
}`)

// You can specify key path by providing arguments to Get function
jsonparser.Get(data, "person", "name", "fullName")

// There is `GetInt` and `GetBoolean` helpers if you exactly know key data type
jsonparser.GetInt(data, "person", "github", "followers")

// When you try to get object, it will return you []byte slice pointer to data containing it
// In `company` it will be `{"name": "Acme"}`
jsonparser.Get(data, "company")

// If the key doesn't exist it will throw an error
var size int64
if value, err := jsonparser.GetInt(data, "company", "size"); err == nil {
  size = value
}

// You can use `ArrayEach` helper to iterate items [item1, item2 .... itemN]
jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	fmt.Println(jsonparser.Get(value, "url"))
}, "person", "avatars")

// Or use can access fields by index!
jsonparser.GetString(data, "person", "avatars", "[0]", "url")

// You can use `ObjectEach` helper to iterate objects { "key1":object1, "key2":object2, .... "keyN":objectN }
jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
        fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
	return nil
}, "person", "name")

// The most efficient way to extract multiple keys is `EachKey`

paths := [][]string{
  []string{"person", "name", "fullName"},
  []string{"person", "avatars", "[0]", "url"},
  []string{"company", "url"},
}
jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error){
  switch idx {
  case 0: // []string{"person", "name", "fullName"}
    ...
  case 1: // []string{"person", "avatars", "[0]", "url"}
    ...
  case 2: // []string{"company", "url"},
    ...
  }
}, paths...)

// For more information see docs below
```

## Reference

Library API is really simple. You just need the `Get` method to perform any operation. The rest is just helpers around it.

You also can view API at [godoc.org](https://godoc.org/github.com/buger/jsonparser)


### **`Get`**
```go
func Get(data []byte, keys ...string) (value []byte, dataType jsonparser.ValueType, offset int, err error)
```
Receives data structure, and key path to extract value from.

Returns:
* `value` - Pointer to original data structure containing key value, or just empty slice if nothing found or error
* `dataType` - 	Can be: `NotExist`, `String`, `Number`, `Object`, `Array`, `Boolean` or `Null`
* `offset` - Offset from provided data structure where key value ends. Used mostly internally, for example for `ArrayEach` helper.
* `err` - If the key is not found or any other parsing issue, it should return error. If key not found it also sets `dataType` to `NotExist`

Accepts multiple keys to specify path to JSON value (in case of quering nested structures).
If no keys are provided it will try to extract the closest JSON value (simple ones or object/array), useful for reading streams or arrays, see `ArrayEach` implementation.

Note that keys can be an array indexes: `jsonparser.GetInt("person", "avatars", "[0]", "url")`, pretty cool, yeah?

### **`GetString`**
```go
func GetString(data []byte, keys ...string) (val string, err error)
```
Returns strings properly handing escaped and unicode characters. Note that this will cause additional memory allocations.

### **`GetUnsafeString`**
If you need string in your app, and ready to sacrifice with support of escaped symbols in favor of speed. It returns string mapped to existing byte slice memory, without any allocations:
```go
s, _, := jsonparser.GetUnsafeString(data, "person", "name", "title")
switch s {
  case 'CEO':
    ...
  case 'Engineer'
    ...
  ...
}
```
Note that `unsafe` here means that your string will exist until GC will free underlying byte slice, for most of cases it means that you can use this string only in current context, and should not pass it anywhere externally: through channels or any other way.


### **`GetBoolean`**, **`GetInt`** and **`GetFloat`**
```go
func GetBoolean(data []byte, keys ...string) (val bool, err error)

func GetFloat(data []byte, keys ...string) (val float64, err error)

func GetInt(data []byte, keys ...string) (val int64, err error)
```
If you know the key type, you can use the helpers above.
If key data type do not match, it will return error.

### **`ArrayEach`**
```go
func ArrayEach(data []byte, cb func(value []byte, dataType jsonparser.ValueType, offset int, err error), keys ...string)
```
Needed for iterating arrays, accepts a callback function with the same return arguments as `Get`.

### **`ObjectEach`**
```go
func ObjectEach(data []byte, callback func(key []byte, value []byte, dataType ValueType, offset int) error, keys ...string) (err error)
```
Needed for iterating object, accepts a callback function. Example:
```go
var handler func([]byte, []byte, jsonparser.ValueType, int) error
handler = func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	//do stuff here
}
jsonparser.ObjectEach(myJson, handler)
```


### **`EachKey`**
```go
func EachKey(data []byte, cb func(idx int, value []byte, dataType jsonparser.ValueType, err error), paths ...[]string)
```
When you need to read multiple keys, and you do not afraid of low-level API `EachKey` is your friend. It read payload only single time, and calls callback function once path is found. For example when you call multiple times `Get`, it has to process payload multiple times, each time you call it. Depending on payload `EachKey` can be multiple times faster than `Get`. Path can use nested keys as well!

```go
paths := [][]string{
	[]string{"uuid"},
	[]string{"tz"},
	[]string{"ua"},
	[]string{"st"},
}
var data SmallPayload

jsonparser.EachKey(smallFixture, func(idx int, value []byte, vt jsonparser.ValueType, err error){
	switch idx {
	case 0:
		data.Uuid, _ = value
	case 1:
		v, _ := jsonparser.ParseInt(value)
		data.Tz = int(v)
	case 2:
		data.Ua, _ = value
	case 3:
		v, _ := jsonparser.ParseInt(value)
		data.St = int(v)
	}
}, paths...)
```

### **`Set`**
```go
func Set(data []byte, setValue []byte, keys ...string) (value []byte, err error)
```
Receives existing data structure, key path to set, and value to set at that key. *This functionality is experimental.*

Returns:
* `value` - Pointer to original data structure with updated or added key value.
* `err` - If any parsing issue, it should return error.

Accepts multiple keys to specify path to JSON value (in case of updating or creating  nested structures).

Note that keys can be an array indexes: `jsonparser.Set(data, []byte("http://github.com"), "person", "avatars", "[0]", "url")`

### **`Delete`**
```go
func Delete(data []byte, keys ...string) value []byte
```
Receives existing data structure, and key path to delete. *This functionality is experimental.*

Returns:
* `value` - Pointer to original data structure with key path deleted if it can be found. If there is no key path, then the whole data structure is deleted.

Accepts multiple keys to specify path to JSON value (in case of updating or creating  nested structures).

Note that keys can be an array indexes: `jsonparser.Delete(data, "person", "avatars", "[0]", "url")`


## What makes it so fast?
* It does not rely on `encoding/json`, `reflection` or `interface{}`, the only real package dependency is `bytes`.
* Operates with JSON payload on byte level, providing you pointers to the original data structure: no memory allocation.
* No automatic type conversions, by default everything is a []byte, but it provides you value type, so you can convert by yourself (there is few helpers included).
* Does not parse full record, only keys you specified


## Benchmarks

There are 3 benchmark types, trying to simulate real-life usage for small, medium and large JSON payloads.
For each metric, the lower value is better. Time/op is in nanoseconds. Values better than standard encoding/json marked as bold text.
Benchmarks run on standard Linode 1024 box.

Compared libraries:
* https://golang.org/pkg/encoding/json
* https://github.com/Jeffail/gabs
* https://github.com/a8m/djson
* https://github.com/bitly/go-simplejson
* https://github.com/antonholmquist/jason
* https://github.com/mreiferson/go-ujson
* https://github.com/ugorji/go/codec
* https://github.com/pquerna/ffjson
* https://github.com/mailru/easyjson
* https://github.com/buger/jsonparser

#### TLDR
If you want to skip next sections we have 2 winner: `jsonparser` and `easyjson`.
`jsonparser` is up to 10 times faster than standard `encoding/json` package (depending on payload size and usage), and almost infinitely (literally) better in memory consumption because it operates with data on byte level, and provide direct slice pointers.
`easyjson` wins in CPU in medium tests and frankly i'm impressed with this package: it is remarkable results considering that it is almost drop-in replacement for `encoding/json` (require some code generation).

It's hard to fully compare `jsonparser` and `easyjson` (or `ffson`), they a true parsers and fully process record, unlike `jsonparser` which parse only keys you specified.

If you searching for replacement of `encoding/json` while keeping structs, `easyjson` is an amazing choice. If you want to process dynamic JSON, have memory constrains, or more control over your data you should try `jsonparser`.

`jsonparser` performance heavily depends on usage, and it works best when you do not need to process full record, only some keys. The more calls you need to make, the slower it will be, in contrast `easyjson` (or `ffjson`, `encoding/json`) parser record only 1 time, and then you can make as many calls as you want.

With great power comes great responsibility! :)


#### Small payload

Each test processes 190 bytes of http log as a JSON record.
It should read multiple fields.
https://github.com/buger/jsonparser/blob/master/benchmark/benchmark_small_payload_test.go

Library | time/op | bytes/op | allocs/op 
 ------ | ------- | -------- | -------
encoding/json struct | 7879 | 880 | 18 
encoding/json interface{} | 8946 | 1521 | 38
Jeffail/gabs | 10053 | 1649 | 46
bitly/go-simplejson | 10128 | 2241 | 36 
antonholmquist/jason | 27152 | 7237 | 101 
github.com/ugorji/go/codec | 8806 | 2176 | 31 
mreiferson/go-ujson | **7008** | **1409** | 37 
a8m/djson | 3862 | 1249 | 30 
pquerna/ffjson | **3769** | **624** | **15** 
mailru/easyjson | **2002** | **192** | **9** 
buger/jsonparser | **1367** | **0** | **0** 
buger/jsonparser (EachKey API) | **809** | **0** | **0** 

Winners are ffjson, easyjson and jsonparser, where jsonparser is up to 9.8x faster than encoding/json and 4.6x faster than ffjson, and slightly faster than easyjson.
If you look at memory allocation, jsonparser has no rivals, as it makes no data copy and operates with raw []byte structures and pointers to it.

#### Medium payload

Each test processes a 2.4kb JSON record (based on Clearbit API).
It should read multiple nested fields and 1 array.

https://github.com/buger/jsonparser/blob/master/benchmark/benchmark_medium_payload_test.go

| Library | time/op | bytes/op | allocs/op |
| ------- | ------- | -------- | --------- |
| encoding/json struct | 57749 | 1336 | 29 |
| encoding/json interface{} | 79297 | 10627 | 215 |
| Jeffail/gabs | 83807 | 11202 | 235 |
| bitly/go-simplejson | 88187 | 17187 | 220 |
| antonholmquist/jason | 94099 | 19013 | 247 |
| github.com/ugorji/go/codec | 114719 | 6712 | 152 |
| mreiferson/go-ujson | **56972** | 11547 | 270 |
| a8m/djson | 28525 | 10196 | 198 | 
| pquerna/ffjson | **20298** | **856** | **20** |
| mailru/easyjson | **10512** | **336** | **12** |
| buger/jsonparser | **15955** | **0** | **0** |
| buger/jsonparser (EachKey API) | **8916** | **0** | **0** |

The difference between ffjson and jsonparser in CPU usage is smaller, while the memory consumption difference is growing. On the other hand `easyjson` shows remarkable performance for medium payload.

`gabs`, `go-simplejson` and `jason` are based on encoding/json and map[string]interface{} and actually only helpers for unstructured JSON, their performance correlate with `encoding/json interface{}`, and they will skip next round.
`go-ujson` while have its own parser, shows same performance as `encoding/json`, also skips next round. Same situation with `ugorji/go/codec`, but it showed unexpectedly bad performance for complex payloads.


#### Large payload

Each test processes a 24kb JSON record (based on Discourse API)
It should read 2 arrays, and for each item in array get a few fields.
Basically it means processing a full JSON file.

https://github.com/buger/jsonparser/blob/master/benchmark/benchmark_large_payload_test.go

| Library | time/op | bytes/op | allocs/op |
| --- | --- | --- | --- |
| encoding/json struct | 748336 | 8272 | 307 |
| encoding/json interface{} | 1224271 | 215425 | 3395 |
| a8m/djson | 510082 | 213682 | 2845 |
| pquerna/ffjson | **312271** | **7792** | **298** |
| mailru/easyjson | **154186** | **6992** | **288** |
| buger/jsonparser | **85308** | **0** | **0** |

`jsonparser` now is a winner, but do not forget that it is way more lightweight parser than `ffson` or `easyjson`, and they have to parser all the data, while `jsonparser` parse only what you need. All `ffjson`, `easysjon` and `jsonparser` have their own parsing code, and does not depend on `encoding/json` or `interface{}`, thats one of the reasons why they are so fast. `easyjson` also use a bit of `unsafe` package to reduce memory consuption (in theory it can lead to some unexpected GC issue, but i did not tested enough)

Also last benchmark did not included `EachKey` test, because in this particular case we need to read lot of Array values, and using `ArrayEach` is more efficient. 

## Questions and support

All bug-reports and suggestions should go though Github Issues.

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Added some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Development

All my development happens using Docker, and repo include some Make tasks to simplify development.

* `make build` - builds docker image, usually can be called only once
* `make test` - run tests
* `make fmt` - run go fmt
* `make bench` - run benchmarks (if you need to run only single benchmark modify `BENCHMARK` variable in make file)
* `make profile` - runs benchmark and generate 3 files-  `cpu.out`, `mem.mprof` and `benchmark.test` binary, which can be used for `go tool pprof`
* `make bash` - enter container (i use it for running `go tool pprof` above)
