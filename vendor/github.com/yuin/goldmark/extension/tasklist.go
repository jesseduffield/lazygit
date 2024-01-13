package extension

import (
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"regexp"
)

var taskListRegexp = regexp.MustCompile(`^\[([\sxX])\]\s*`)

type taskCheckBoxParser struct {
}

var defaultTaskCheckBoxParser = &taskCheckBoxParser{}

// NewTaskCheckBoxParser returns a new  InlineParser that can parse
// checkboxes in list items.
// This parser must take precedence over the parser.LinkParser.
func NewTaskCheckBoxParser() parser.InlineParser {
	return defaultTaskCheckBoxParser
}

func (s *taskCheckBoxParser) Trigger() []byte {
	return []byte{'['}
}

func (s *taskCheckBoxParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	// Given AST structure must be like
	// - List
	//   - ListItem         : parent.Parent
	//     - TextBlock      : parent
	//       (current line)
	if parent.Parent() == nil || parent.Parent().FirstChild() != parent {
		return nil
	}

	if _, ok := parent.Parent().(*gast.ListItem); !ok {
		return nil
	}
	line, _ := block.PeekLine()
	m := taskListRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	value := line[m[2]:m[3]][0]
	block.Advance(m[1])
	checked := value == 'x' || value == 'X'
	return ast.NewTaskCheckBox(checked)
}

func (s *taskCheckBoxParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// TaskCheckBoxHTMLRenderer is a renderer.NodeRenderer implementation that
// renders checkboxes in list items.
type TaskCheckBoxHTMLRenderer struct {
	html.Config
}

// NewTaskCheckBoxHTMLRenderer returns a new TaskCheckBoxHTMLRenderer.
func NewTaskCheckBoxHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &TaskCheckBoxHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *TaskCheckBoxHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindTaskCheckBox, r.renderTaskCheckBox)
}

func (r *TaskCheckBoxHTMLRenderer) renderTaskCheckBox(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		return gast.WalkContinue, nil
	}
	n := node.(*ast.TaskCheckBox)

	if n.IsChecked {
		w.WriteString(`<input checked="" disabled="" type="checkbox"`)
	} else {
		w.WriteString(`<input disabled="" type="checkbox"`)
	}
	if r.XHTML {
		w.WriteString(" /> ")
	} else {
		w.WriteString("> ")
	}
	return gast.WalkContinue, nil
}

type taskList struct {
}

// TaskList is an extension that allow you to use GFM task lists.
var TaskList = &taskList{}

func (e *taskList) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTaskCheckBoxParser(), 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTaskCheckBoxHTMLRenderer(), 500),
	))
}
