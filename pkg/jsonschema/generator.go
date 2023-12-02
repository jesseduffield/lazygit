//go:build ignore

package main

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/jsonschema"
)

func main() {
	fmt.Printf("Generating jsonschema in %s...\n", jsonschema.GetSchemaDir())
	jsonschema.GenerateSchema()
}
