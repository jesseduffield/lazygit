package main

import (
	"github.com/mgutz/goa"
	f "github.com/mgutz/goa/filter"
	"github.com/mgutz/str"

	do "gopkg.in/godo.v2"
	"gopkg.in/godo.v2/util"
)

// Project is local project.
func tasks(p *do.Project) {
	p.Task("default", do.S{"readme"}, nil)

	p.Task("install", nil, func(c *do.Context) {
		c.Run("go get github.com/robertkrimen/godocdown/godocdown")
	})

	p.Task("lint", nil, func(c *do.Context) {
		c.Run("golint .")
		c.Run("gofmt -w -s .")
		c.Run("go vet .")
		c.Run("go test")
	})

	p.Task("readme", nil, func(c *do.Context) {
		c.Run("godocdown -output README.md")

		packageName, _ := util.PackageName("doc.go")

		// add godoc link
		goa.Pipe(
			f.Load("./README.md"),
			f.Str(str.ReplaceF("--", "\n[godoc](https://godoc.org/"+packageName+")\n", 1)),
			f.Write(),
		)
	}).Src("**/*.go")

	p.Task("test", nil, func(c *do.Context) {
		c.Run("go test")
	})
}

func main() {
	do.Godo(tasks)
}
