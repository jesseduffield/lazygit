// Package renderer renders the given AST to certain formats.
package renderer

import (
	"bufio"
	"io"
	"sync"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// A Config struct is a data structure that holds configuration of the Renderer.
type Config struct {
	Options       map[OptionName]interface{}
	NodeRenderers util.PrioritizedSlice
}

// NewConfig returns a new Config
func NewConfig() *Config {
	return &Config{
		Options:       map[OptionName]interface{}{},
		NodeRenderers: util.PrioritizedSlice{},
	}
}

// An OptionName is a name of the option.
type OptionName string

// An Option interface is a functional option type for the Renderer.
type Option interface {
	SetConfig(*Config)
}

type withNodeRenderers struct {
	value []util.PrioritizedValue
}

func (o *withNodeRenderers) SetConfig(c *Config) {
	c.NodeRenderers = append(c.NodeRenderers, o.value...)
}

// WithNodeRenderers is a functional option that allow you to add
// NodeRenderers to the renderer.
func WithNodeRenderers(ps ...util.PrioritizedValue) Option {
	return &withNodeRenderers{ps}
}

type withOption struct {
	name  OptionName
	value interface{}
}

func (o *withOption) SetConfig(c *Config) {
	c.Options[o.name] = o.value
}

// WithOption is a functional option that allow you to set
// an arbitrary option to the parser.
func WithOption(name OptionName, value interface{}) Option {
	return &withOption{name, value}
}

// A SetOptioner interface sets given option to the object.
type SetOptioner interface {
	// SetOption sets given option to the object.
	// Unacceptable options may be passed.
	// Thus implementations must ignore unacceptable options.
	SetOption(name OptionName, value interface{})
}

// NodeRendererFunc is a function that renders a given node.
type NodeRendererFunc func(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error)

// A NodeRenderer interface offers NodeRendererFuncs.
type NodeRenderer interface {
	// RendererFuncs registers NodeRendererFuncs to given NodeRendererFuncRegisterer.
	RegisterFuncs(NodeRendererFuncRegisterer)
}

// A NodeRendererFuncRegisterer registers
type NodeRendererFuncRegisterer interface {
	// Register registers given NodeRendererFunc to this object.
	Register(ast.NodeKind, NodeRendererFunc)
}

// A Renderer interface renders given AST node to given
// writer with given Renderer.
type Renderer interface {
	Render(w io.Writer, source []byte, n ast.Node) error

	// AddOptions adds given option to this renderer.
	AddOptions(...Option)
}

type renderer struct {
	config               *Config
	options              map[OptionName]interface{}
	nodeRendererFuncsTmp map[ast.NodeKind]NodeRendererFunc
	maxKind              int
	nodeRendererFuncs    []NodeRendererFunc
	initSync             sync.Once
}

// NewRenderer returns a new Renderer with given options.
func NewRenderer(options ...Option) Renderer {
	config := NewConfig()
	for _, opt := range options {
		opt.SetConfig(config)
	}

	r := &renderer{
		options:              map[OptionName]interface{}{},
		config:               config,
		nodeRendererFuncsTmp: map[ast.NodeKind]NodeRendererFunc{},
	}

	return r
}

func (r *renderer) AddOptions(opts ...Option) {
	for _, opt := range opts {
		opt.SetConfig(r.config)
	}
}

func (r *renderer) Register(kind ast.NodeKind, v NodeRendererFunc) {
	r.nodeRendererFuncsTmp[kind] = v
	if int(kind) > r.maxKind {
		r.maxKind = int(kind)
	}
}

// Render renders the given AST node to the given writer with the given Renderer.
func (r *renderer) Render(w io.Writer, source []byte, n ast.Node) error {
	r.initSync.Do(func() {
		r.options = r.config.Options
		r.config.NodeRenderers.Sort()
		l := len(r.config.NodeRenderers)
		for i := l - 1; i >= 0; i-- {
			v := r.config.NodeRenderers[i]
			nr, _ := v.Value.(NodeRenderer)
			if se, ok := v.Value.(SetOptioner); ok {
				for oname, ovalue := range r.options {
					se.SetOption(oname, ovalue)
				}
			}
			nr.RegisterFuncs(r)
		}
		r.nodeRendererFuncs = make([]NodeRendererFunc, r.maxKind+1)
		for kind, nr := range r.nodeRendererFuncsTmp {
			r.nodeRendererFuncs[kind] = nr
		}
		r.config = nil
		r.nodeRendererFuncsTmp = nil
	})
	writer, ok := w.(util.BufWriter)
	if !ok {
		writer = bufio.NewWriter(w)
	}
	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error
		f := r.nodeRendererFuncs[n.Kind()]
		if f != nil {
			s, err = f(writer, source, n, entering)
		}
		return s, err
	})
	if err != nil {
		return err
	}
	return writer.Flush()
}
