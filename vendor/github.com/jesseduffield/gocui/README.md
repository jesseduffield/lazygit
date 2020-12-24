# GOCUI - Go Console User Interface
[![CircleCI](https://circleci.com/gh/awesome-gocui/gocui/tree/master.svg?style=svg)](https://circleci.com/gh/awesome-gocui/gocui/tree/master)
[![CodeCov](https://codecov.io/gh/awesome-gocui/gocui/branch/master/graph/badge.svg)](https://codecov.io/gh/awesome-gocui/gocui)
[![Go Report Card](https://goreportcard.com/badge/github.com/awesome-gocui/gocui)](https://goreportcard.com/report/github.com/awesome-gocui/gocui)
[![GolangCI](https://golangci.com/badges/github.com/awesome-gocui/gocui.svg)](https://golangci.com/badges/github.com/awesome-gocui/gocui.svg)
[![GoDoc](https://godoc.org/github.com/awesome-gocui/gocui?status.svg)](https://godoc.org/github.com/awesome-gocui/gocui)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/awesome-gocui/gocui.svg)

Minimalist Go package aimed at creating Console User Interfaces.
A community fork based on the amazing work of [jroimartin](https://github.com/jroimartin/gocui)

## Features

* Minimalist API.
* Views (the "windows" in the GUI) implement the interface io.ReadWriter.
* Support for overlapping views.
* The GUI can be modified at runtime (concurrent-safe).
* Global and view-level keybindings.
* Mouse support.
* Colored text.
* Customizable editing mode.
* Easy to build reusable widgets, complex layouts...

## About fork

This fork has many improvements over the original work from [jroimartin](https://github.com/jroimartin/gocui).

* Written ontop of TCell
* Better wide character support
* Support for 1 Line height views
* Better support for running in docker container
* Customize frame colors
* Improved code comments and quality
* Many small improvements
* Change Visibility of views

For information about this org see: [awesome-gocui/about](https://github.com/awesome-gocui/about).

## Installation

Execute:

```
$ go get github.com/awesome-gocui/gocui
```

## Documentation

Execute:

```
$ go doc github.com/awesome-gocui/gocui
```

Or visit [godoc.org](https://godoc.org/github.com/awesome-gocui/gocui) to read it
online.

## Example
See the [_example](./_example/) folder for more examples

```go
package main

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		if _, err := g.SetCurrentView("hello"); err != nil {
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
* [lazygit](https://github.com/jesseduffield/lazygit): simple terminal UI for git commands.
* [lazydocker](https://github.com/jesseduffield/lazydocker): The lazier way to manage everything docker.

Note: if your project is not listed here, let us know! :)
