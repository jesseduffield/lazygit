# Emoji
Emoji is a simple golang package.

![master workflow](https://github.com/kyokomi/emoji/actions/workflows/go.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/kyokomi/emoji/badge.png?branch=master)](https://coveralls.io/r/kyokomi/emoji?branch=master)
[![GoDoc](https://pkg.go.dev/badge/github.com/kyokomi/emoji.svg)](https://pkg.go.dev/github.com/kyokomi/emoji/v2)

Get it:

```
go get github.com/kyokomi/emoji/v2
```

Import it:

```
import (
	"github.com/kyokomi/emoji/v2"
)
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/kyokomi/emoji/v2"
)

func main() {
	fmt.Println("Hello World Emoji!")

	emoji.Println(":beer: Beer!!!")

	pizzaMessage := emoji.Sprint("I like a :pizza: and :sushi:!!")
	fmt.Println(pizzaMessage)
}
```

## Demo

![demo](screen/image.png)

## Reference

- [unicode Emoji Charts](http://www.unicode.org/emoji/charts/emoji-list.html)

## License

[MIT](https://github.com/kyokomi/emoji/blob/master/LICENSE)
