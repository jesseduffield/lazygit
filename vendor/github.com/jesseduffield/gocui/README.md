# GOCUI - Go Console User Interface

[![GoDoc](https://godoc.org/github.com/jroimartin/gocui?status.svg)](https://godoc.org/github.com/jroimartin/gocui)

Minimalist Go package aimed at creating Console User Interfaces.

## Features

* Minimalist API.
* Views (the "windows" in the GUI) implement the interface io.ReadWriter.
* Support for overlapping views.
* The GUI can be modified at runtime (concurrent-safe).
* Global and view-level keybindings.
* Mouse support.
* Colored text.
* Customizable edition mode.
* Easy to build reusable widgets, complex layouts...

## Installation

Execute:

```
$ go get github.com/jroimartin/gocui
```

## Documentation

Execute:

```
$ go doc github.com/jroimartin/gocui
```

Or visit [godoc.org](https://godoc.org/github.com/jroimartin/gocui) to read it
online.

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
```

## Screenshots

![r2cui](https://cloud.githubusercontent.com/assets/1223476/19418932/63645052-93ce-11e6-867c-da5e97e37237.png)

![_examples/demo.go](https://cloud.githubusercontent.com/assets/1223476/5992750/720b84f0-aa36-11e4-88ec-296fa3247b52.png)

![_examples/dynamic.go](https://cloud.githubusercontent.com/assets/1223476/5992751/76ad5cc2-aa36-11e4-8204-6a90269db827.png)

## Projects using gocui

* [komanda-cli](https://github.com/mephux/komanda-cli): IRC Client For Developers.
* [vuls](https://github.com/future-architect/vuls): Agentless vulnerability scanner for Linux/FreeBSD.
* [wuzz](https://github.com/asciimoo/wuzz): Interactive cli tool for HTTP inspection.
* [httplab](https://github.com/gchaincl/httplab): Interactive web server.
* [domainr](https://github.com/MichaelThessel/domainr): Tool that checks the availability of domains based on keywords.
* [gotime](https://github.com/nanohard/gotime): Time tracker for projects and tasks.
* [claws](https://github.com/thehowl/claws): Interactive command line client for testing websockets.
* [terminews](http://github.com/antavelos/terminews): Terminal based RSS reader.
* [diagram](https://github.com/esimov/diagram): Tool to convert ascii arts into hand drawn diagrams.
* [pody](https://github.com/JulienBreux/pody): CLI app to manage Pods in a Kubernetes cluster.
* [kubexp](https://github.com/alitari/kubexp): Kubernetes client.
* [kcli](https://github.com/cswank/kcli): Tool for inspecting kafka topics/partitions/messages.
* [fac](https://github.com/mkchoi212/fac): git merge conflict resolver
* [jsonui](https://github.com/gulyasm/jsonui): Interactive JSON explorer for your terminal.
* [cointop](https://github.com/miguelmota/cointop): Interactive terminal based UI application for tracking cryptocurrencies.

Note: if your project is not listed here, let us know! :)
