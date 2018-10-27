package procs

import (
	"os"
	"strings"
)

// Builder helps construct commands using templates.
type Builder struct {
	Context   map[string]string
	Templates []string
}

func (b *Builder) getConfig(ctx map[string]string) func(string) string {
	return func(key string) string {
		if v, ok := ctx[key]; ok {
			return v
		}
		return ""
	}
}

func (b *Builder) expand(v string, ctx map[string]string) string {
	return os.Expand(v, b.getConfig(ctx))
}

// Command returns the result of the templates as a single string.
func (b *Builder) Command() string {
	parts := []string{}
	for _, t := range b.Templates {
		parts = append(parts, b.expand(t, b.Context))
	}

	return strings.Join(parts, " ")
}

// CommandContext returns the result of the templates as a single
// string, but allows providing an environment context as a
// map[string]string for expansions.
func (b *Builder) CommandContext(ctx map[string]string) string {
	// Build our environment context by starting with our Builder
	// context and overlay the passed in context map.
	env := make(map[string]string)
	for k, v := range b.Context {
		env[k] = b.expand(v, b.Context)
	}

	for k, v := range ctx {
		env[k] = b.expand(v, env)
	}

	parts := []string{}
	for _, t := range b.Templates {
		parts = append(parts, b.expand(t, env))
	}

	return strings.Join(parts, " ")
}
