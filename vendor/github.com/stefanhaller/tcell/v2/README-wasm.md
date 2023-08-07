# WASM for _Tcell_

You can build _Tcell_ project into a webpage by compiling it slightly differently. This will result in a _Tcell_ project you can embed into another html page, or use as a standalone page.

## Building your project

WASM needs special build flags in order to work. You can build it by executing
```sh
GOOS=js GOARCH=wasm go build -o yourfile.wasm
```

## Additional files

You also need 5 other files in the same directory as the wasm. Four (`tcell.html`, `tcell.js`, `termstyle.css`, and `beep.wav`) are provided in the `webfiles` directory. The last one, `wasm_exec.js`, can be copied from GOROOT into the current directory by executing
```sh
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./
```

In `tcell.js`, you also need to change the constant
```js
const wasmFilePath = "yourfile.wasm"
```
to the file you outputed to when building.

## Displaying your project

### Standalone

You can see the project (with an white background around the terminal) by serving the directory. You can do this using any framework, including another golang project:

```golang
// server.go

package main

import (
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080",
		http.FileServer(http.Dir("/path/to/dir/to/serve")),
	))
}

```

To see the webpage with this example, you can type in `localhost:8080/tcell.html` into your browser while `server.go` is running.

### Embedding
It is recomended to use an iframe if you want to embed the app into a webpage:
```html
<iframe src="tcell.html" title="Tcell app"></iframe>
```

## Other considerations

### Accessing files

`io.Open(filename)` and other related functions for reading file systems do not work; use `http.Get(filename)` instead.