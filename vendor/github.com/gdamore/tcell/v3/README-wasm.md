# WASM for _Tcell_

You can build _Tcell_ project into a webpage by compiling it slightly differently. This will result in a _Tcell_ project you can embed into another html page, or use as a standalone page.

## Building your project

WASM needs special build flags in order to work. You can build it by executing
```sh
GOOS=js GOARCH=wasm go build -o yourfile.wasm
```

## Additional files

You also need the supporting web files in the same directory as the wasm. The files `tcell.html`, `tcell.js`, `termstyle.css`, and `beep.wav`, plus the `ghostty-web` directory, are provided in the `webfiles` directory. The last file, `wasm_exec.js`, can be copied from GOROOT into the current directory by executing
```sh
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ./
```

The web frontend uses `ghostty-web`. The required browser runtime files are vendored in `webfiles/ghostty-web` and must be copied alongside `tcell.js`; no npm, bundler, or external CDN is required. The vendored `ghostty-web` files are MIT licensed; see `webfiles/ghostty-web/LICENSE`.

```sh
cp -R webfiles/ghostty-web /path/to/dir/to/serve/
```

The vendored `ghostty-web.js` is intentionally browser-only. Its upstream Node `readFile` fallback import is removed so browser-oriented servers and bundlers such as Vite do not try to resolve a Node file-system shim; the bundled code loads `ghostty-vt.wasm` with `fetch`.

For example:

```sh
mkdir -p /tmp/tcell-wasm
cp webfiles/tcell.html webfiles/tcell.js webfiles/termstyle.css webfiles/beep.wav /tmp/tcell-wasm/
cp -R webfiles/ghostty-web /tmp/tcell-wasm/
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" /tmp/tcell-wasm/
GOOS=js GOARCH=wasm go build -o /tmp/tcell-wasm/main.wasm ./demos/unicode
python3 -m http.server -d /tmp/tcell-wasm 8080
```

In `tcell.js`, you also need to change the constant
```js
const wasmFilePath = "yourfile.wasm"
```
to the file you outputted to when building.

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
It is recommended to use an iframe if you want to embed the app into a webpage:
```html
<iframe src="tcell.html" title="Tcell app"></iframe>
```

### Sizing

By default the web terminal fits itself to the size of the `#terminal` element and reacts to container resizes. The bundled `termstyle.css` makes this full-page by default.

You can override the terminal cell dimensions explicitly in HTML:

```html
<pre id="terminal" data-cols="100" data-rows="30"></pre>
```

If only one of `data-cols` or `data-rows` is set, the other dimension remains reactive.

## Other considerations

### Accessing files

`io.Open(filename)` and other related functions for reading file systems do not work; use `http.Get(filename)` instead.

### Keyboard shortcuts

The browser may reserve some key combinations before JavaScript can see or cancel them. This is especially common for Meta/Command shortcuts on macOS, such as Command-L. Standalone Meta key events can be reported, but Meta-modified key combinations are browser-dependent and should not be relied upon in WASM web mode.
