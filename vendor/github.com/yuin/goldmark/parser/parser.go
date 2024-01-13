// Package parser contains stuff that are related to parsing a Markdown text.
package parser

import (
	"fmt"
	"strings"
	"sync"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// A Reference interface represents a link reference in Markdown text.
type Reference interface {
	// String implements Stringer.
	String() string

	// Label returns a label of the reference.
	Label() []byte

	// Destination returns a destination(URL) of the reference.
	Destination() []byte

	// Title returns a title of the reference.
	Title() []byte
}

type reference struct {
	label       []byte
	destination []byte
	title       []byte
}

// NewReference returns a new Reference.
func NewReference(label, destination, title []byte) Reference {
	return &reference{label, destination, title}
}

func (r *reference) Label() []byte {
	return r.label
}

func (r *reference) Destination() []byte {
	return r.destination
}

func (r *reference) Title() []byte {
	return r.title
}

func (r *reference) String() string {
	return fmt.Sprintf("Reference{Label:%s, Destination:%s, Title:%s}", r.label, r.destination, r.title)
}

// An IDs interface is a collection of the element ids.
type IDs interface {
	// Generate generates a new element id.
	Generate(value []byte, kind ast.NodeKind) []byte

	// Put puts a given element id to the used ids table.
	Put(value []byte)
}

type ids struct {
	values map[string]bool
}

func newIDs() IDs {
	return &ids{
		values: map[string]bool{},
	}
}

func (s *ids) Generate(value []byte, kind ast.NodeKind) []byte {
	value = util.TrimLeftSpace(value)
	value = util.TrimRightSpace(value)
	result := []byte{}
	for i := 0; i < len(value); {
		v := value[i]
		l := util.UTF8Len(v)
		i += int(l)
		if l != 1 {
			continue
		}
		if util.IsAlphaNumeric(v) {
			if 'A' <= v && v <= 'Z' {
				v += 'a' - 'A'
			}
			result = append(result, v)
		} else if util.IsSpace(v) || v == '-' || v == '_' {
			result = append(result, '-')
		}
	}
	if len(result) == 0 {
		if kind == ast.KindHeading {
			result = []byte("heading")
		} else {
			result = []byte("id")
		}
	}
	if _, ok := s.values[util.BytesToReadOnlyString(result)]; !ok {
		s.values[util.BytesToReadOnlyString(result)] = true
		return result
	}
	for i := 1; ; i++ {
		newResult := fmt.Sprintf("%s-%d", result, i)
		if _, ok := s.values[newResult]; !ok {
			s.values[newResult] = true
			return []byte(newResult)
		}

	}
}

func (s *ids) Put(value []byte) {
	s.values[util.BytesToReadOnlyString(value)] = true
}

// ContextKey is a key that is used to set arbitrary values to the context.
type ContextKey int

// ContextKeyMax is a maximum value of the ContextKey.
var ContextKeyMax ContextKey

// NewContextKey return a new ContextKey value.
func NewContextKey() ContextKey {
	ContextKeyMax++
	return ContextKeyMax
}

// A Context interface holds a information that are necessary to parse
// Markdown text.
type Context interface {
	// String implements Stringer.
	String() string

	// Get returns a value associated with the given key.
	Get(ContextKey) interface{}

	// ComputeIfAbsent computes a value if a value associated with the given key is absent and returns the value.
	ComputeIfAbsent(ContextKey, func() interface{}) interface{}

	// Set sets the given value to the context.
	Set(ContextKey, interface{})

	// AddReference adds the given reference to this context.
	AddReference(Reference)

	// Reference returns (a reference, true) if a reference associated with
	// the given label exists, otherwise (nil, false).
	Reference(label string) (Reference, bool)

	// References returns a list of references.
	References() []Reference

	// IDs returns a collection of the element ids.
	IDs() IDs

	// BlockOffset returns a first non-space character position on current line.
	// This value is valid only for BlockParser.Open.
	// BlockOffset returns -1 if current line is blank.
	BlockOffset() int

	// BlockOffset sets a first non-space character position on current line.
	// This value is valid only for BlockParser.Open.
	SetBlockOffset(int)

	// BlockIndent returns an indent width on current line.
	// This value is valid only for BlockParser.Open.
	// BlockIndent returns -1 if current line is blank.
	BlockIndent() int

	// BlockIndent sets an indent width on current line.
	// This value is valid only for BlockParser.Open.
	SetBlockIndent(int)

	// FirstDelimiter returns a first delimiter of the current delimiter list.
	FirstDelimiter() *Delimiter

	// LastDelimiter returns a last delimiter of the current delimiter list.
	LastDelimiter() *Delimiter

	// PushDelimiter appends the given delimiter to the tail of the current
	// delimiter list.
	PushDelimiter(delimiter *Delimiter)

	// RemoveDelimiter removes the given delimiter from the current delimiter list.
	RemoveDelimiter(d *Delimiter)

	// ClearDelimiters clears the current delimiter list.
	ClearDelimiters(bottom ast.Node)

	// OpenedBlocks returns a list of nodes that are currently in parsing.
	OpenedBlocks() []Block

	// SetOpenedBlocks sets a list of nodes that are currently in parsing.
	SetOpenedBlocks([]Block)

	// LastOpenedBlock returns a last node that is currently in parsing.
	LastOpenedBlock() Block

	// IsInLinkLabel returns true if current position seems to be in link label.
	IsInLinkLabel() bool
}

// A ContextConfig struct is a data structure that holds configuration of the Context.
type ContextConfig struct {
	IDs IDs
}

// An ContextOption is a functional option type for the Context.
type ContextOption func(*ContextConfig)

// WithIDs is a functional option for the Context.
func WithIDs(ids IDs) ContextOption {
	return func(c *ContextConfig) {
		c.IDs = ids
	}
}

type parseContext struct {
	store         []interface{}
	ids           IDs
	refs          map[string]Reference
	blockOffset   int
	blockIndent   int
	delimiters    *Delimiter
	lastDelimiter *Delimiter
	openedBlocks  []Block
}

// NewContext returns a new Context.
func NewContext(options ...ContextOption) Context {
	cfg := &ContextConfig{
		IDs: newIDs(),
	}
	for _, option := range options {
		option(cfg)
	}

	return &parseContext{
		store:         make([]interface{}, ContextKeyMax+1),
		refs:          map[string]Reference{},
		ids:           cfg.IDs,
		blockOffset:   -1,
		blockIndent:   -1,
		delimiters:    nil,
		lastDelimiter: nil,
		openedBlocks:  []Block{},
	}
}

func (p *parseContext) Get(key ContextKey) interface{} {
	return p.store[key]
}

func (p *parseContext) ComputeIfAbsent(key ContextKey, f func() interface{}) interface{} {
	v := p.store[key]
	if v == nil {
		v = f()
		p.store[key] = v
	}
	return v
}

func (p *parseContext) Set(key ContextKey, value interface{}) {
	p.store[key] = value
}

func (p *parseContext) IDs() IDs {
	return p.ids
}

func (p *parseContext) BlockOffset() int {
	return p.blockOffset
}

func (p *parseContext) SetBlockOffset(v int) {
	p.blockOffset = v
}

func (p *parseContext) BlockIndent() int {
	return p.blockIndent
}

func (p *parseContext) SetBlockIndent(v int) {
	p.blockIndent = v
}

func (p *parseContext) LastDelimiter() *Delimiter {
	return p.lastDelimiter
}

func (p *parseContext) FirstDelimiter() *Delimiter {
	return p.delimiters
}

func (p *parseContext) PushDelimiter(d *Delimiter) {
	if p.delimiters == nil {
		p.delimiters = d
		p.lastDelimiter = d
	} else {
		l := p.lastDelimiter
		p.lastDelimiter = d
		l.NextDelimiter = d
		d.PreviousDelimiter = l
	}
}

func (p *parseContext) RemoveDelimiter(d *Delimiter) {
	if d.PreviousDelimiter == nil {
		p.delimiters = d.NextDelimiter
	} else {
		d.PreviousDelimiter.NextDelimiter = d.NextDelimiter
		if d.NextDelimiter != nil {
			d.NextDelimiter.PreviousDelimiter = d.PreviousDelimiter
		}
	}
	if d.NextDelimiter == nil {
		p.lastDelimiter = d.PreviousDelimiter
	}
	if p.delimiters != nil {
		p.delimiters.PreviousDelimiter = nil
	}
	if p.lastDelimiter != nil {
		p.lastDelimiter.NextDelimiter = nil
	}
	d.NextDelimiter = nil
	d.PreviousDelimiter = nil
	if d.Length != 0 {
		ast.MergeOrReplaceTextSegment(d.Parent(), d, d.Segment)
	} else {
		d.Parent().RemoveChild(d.Parent(), d)
	}
}

func (p *parseContext) ClearDelimiters(bottom ast.Node) {
	if p.lastDelimiter == nil {
		return
	}
	var c ast.Node
	for c = p.lastDelimiter; c != nil && c != bottom; {
		prev := c.PreviousSibling()
		if d, ok := c.(*Delimiter); ok {
			p.RemoveDelimiter(d)
		}
		c = prev
	}
}

func (p *parseContext) AddReference(ref Reference) {
	key := util.ToLinkReference(ref.Label())
	if _, ok := p.refs[key]; !ok {
		p.refs[key] = ref
	}
}

func (p *parseContext) Reference(label string) (Reference, bool) {
	v, ok := p.refs[label]
	return v, ok
}

func (p *parseContext) References() []Reference {
	ret := make([]Reference, 0, len(p.refs))
	for _, v := range p.refs {
		ret = append(ret, v)
	}
	return ret
}

func (p *parseContext) String() string {
	refs := []string{}
	for _, r := range p.refs {
		refs = append(refs, r.String())
	}

	return fmt.Sprintf("Context{Store:%#v, Refs:%s}", p.store, strings.Join(refs, ","))
}

func (p *parseContext) OpenedBlocks() []Block {
	return p.openedBlocks
}

func (p *parseContext) SetOpenedBlocks(v []Block) {
	p.openedBlocks = v
}

func (p *parseContext) LastOpenedBlock() Block {
	if l := len(p.openedBlocks); l != 0 {
		return p.openedBlocks[l-1]
	}
	return Block{}
}

func (p *parseContext) IsInLinkLabel() bool {
	tlist := p.Get(linkLabelStateKey)
	return tlist != nil
}

// State represents parser's state.
// State is designed to use as a bit flag.
type State int

const (
	none State = 1 << iota

	// Continue indicates parser can continue parsing.
	Continue

	// Close indicates parser cannot parse anymore.
	Close

	// HasChildren indicates parser may have child blocks.
	HasChildren

	// NoChildren indicates parser does not have child blocks.
	NoChildren

	// RequireParagraph indicates parser requires that the last node
	// must be a paragraph and is not converted to other nodes by
	// ParagraphTransformers.
	RequireParagraph
)

// A Config struct is a data structure that holds configuration of the Parser.
type Config struct {
	Options               map[OptionName]interface{}
	BlockParsers          util.PrioritizedSlice /*<BlockParser>*/
	InlineParsers         util.PrioritizedSlice /*<InlineParser>*/
	ParagraphTransformers util.PrioritizedSlice /*<ParagraphTransformer>*/
	ASTTransformers       util.PrioritizedSlice /*<ASTTransformer>*/
	EscapedSpace          bool
}

// NewConfig returns a new Config.
func NewConfig() *Config {
	return &Config{
		Options:               map[OptionName]interface{}{},
		BlockParsers:          util.PrioritizedSlice{},
		InlineParsers:         util.PrioritizedSlice{},
		ParagraphTransformers: util.PrioritizedSlice{},
		ASTTransformers:       util.PrioritizedSlice{},
	}
}

// An Option interface is a functional option type for the Parser.
type Option interface {
	SetParserOption(*Config)
}

// OptionName is a name of parser options.
type OptionName string

// Attribute is an option name that spacify attributes of elements.
const optAttribute OptionName = "Attribute"

type withAttribute struct {
}

func (o *withAttribute) SetParserOption(c *Config) {
	c.Options[optAttribute] = true
}

// WithAttribute is a functional option that enables custom attributes.
func WithAttribute() Option {
	return &withAttribute{}
}

// A Parser interface parses Markdown text into AST nodes.
type Parser interface {
	// Parse parses the given Markdown text into AST nodes.
	Parse(reader text.Reader, opts ...ParseOption) ast.Node

	// AddOption adds the given option to this parser.
	AddOptions(...Option)
}

// A SetOptioner interface sets the given option to the object.
type SetOptioner interface {
	// SetOption sets the given option to the object.
	// Unacceptable options may be passed.
	// Thus implementations must ignore unacceptable options.
	SetOption(name OptionName, value interface{})
}

// A BlockParser interface parses a block level element like Paragraph, List,
// Blockquote etc.
type BlockParser interface {
	// Trigger returns a list of characters that triggers Parse method of
	// this parser.
	// If Trigger returns a nil, Open will be called with any lines.
	Trigger() []byte

	// Open parses the current line and returns a result of parsing.
	//
	// Open must not parse beyond the current line.
	// If Open has been able to parse the current line, Open must advance a reader
	// position by consumed byte length.
	//
	// If Open has not been able to parse the current line, Open should returns
	// (nil, NoChildren). If Open has been able to parse the current line, Open
	// should returns a new Block node and returns HasChildren or NoChildren.
	Open(parent ast.Node, reader text.Reader, pc Context) (ast.Node, State)

	// Continue parses the current line and returns a result of parsing.
	//
	// Continue must not parse beyond the current line.
	// If Continue has been able to parse the current line, Continue must advance
	// a reader position by consumed byte length.
	//
	// If Continue has not been able to parse the current line, Continue should
	// returns Close. If Continue has been able to parse the current line,
	// Continue should returns (Continue | NoChildren) or
	// (Continue | HasChildren)
	Continue(node ast.Node, reader text.Reader, pc Context) State

	// Close will be called when the parser returns Close.
	Close(node ast.Node, reader text.Reader, pc Context)

	// CanInterruptParagraph returns true if the parser can interrupt paragraphs,
	// otherwise false.
	CanInterruptParagraph() bool

	// CanAcceptIndentedLine returns true if the parser can open new node when
	// the given line is being indented more than 3 spaces.
	CanAcceptIndentedLine() bool
}

// An InlineParser interface parses an inline level element like CodeSpan, Link etc.
type InlineParser interface {
	// Trigger returns a list of characters that triggers Parse method of
	// this parser.
	// Trigger characters must be a punctuation or a halfspace.
	// Halfspaces triggers this parser when character is any spaces characters or
	// a head of line
	Trigger() []byte

	// Parse parse the given block into an inline node.
	//
	// Parse can parse beyond the current line.
	// If Parse has been able to parse the current line, it must advance a reader
	// position by consumed byte length.
	Parse(parent ast.Node, block text.Reader, pc Context) ast.Node
}

// A CloseBlocker interface is a callback function that will be
// called when block is closed in the inline parsing.
type CloseBlocker interface {
	// CloseBlock will be called when a block is closed.
	CloseBlock(parent ast.Node, block text.Reader, pc Context)
}

// A ParagraphTransformer transforms parsed Paragraph nodes.
// For example, link references are searched in parsed Paragraphs.
type ParagraphTransformer interface {
	// Transform transforms the given paragraph.
	Transform(node *ast.Paragraph, reader text.Reader, pc Context)
}

// ASTTransformer transforms entire Markdown document AST tree.
type ASTTransformer interface {
	// Transform transforms the given AST tree.
	Transform(node *ast.Document, reader text.Reader, pc Context)
}

// DefaultBlockParsers returns a new list of default BlockParsers.
// Priorities of default BlockParsers are:
//
//	SetextHeadingParser, 100
//	ThematicBreakParser, 200
//	ListParser, 300
//	ListItemParser, 400
//	CodeBlockParser, 500
//	ATXHeadingParser, 600
//	FencedCodeBlockParser, 700
//	BlockquoteParser, 800
//	HTMLBlockParser, 900
//	ParagraphParser, 1000
func DefaultBlockParsers() []util.PrioritizedValue {
	return []util.PrioritizedValue{
		util.Prioritized(NewSetextHeadingParser(), 100),
		util.Prioritized(NewThematicBreakParser(), 200),
		util.Prioritized(NewListParser(), 300),
		util.Prioritized(NewListItemParser(), 400),
		util.Prioritized(NewCodeBlockParser(), 500),
		util.Prioritized(NewATXHeadingParser(), 600),
		util.Prioritized(NewFencedCodeBlockParser(), 700),
		util.Prioritized(NewBlockquoteParser(), 800),
		util.Prioritized(NewHTMLBlockParser(), 900),
		util.Prioritized(NewParagraphParser(), 1000),
	}
}

// DefaultInlineParsers returns a new list of default InlineParsers.
// Priorities of default InlineParsers are:
//
//	CodeSpanParser, 100
//	LinkParser, 200
//	AutoLinkParser, 300
//	RawHTMLParser, 400
//	EmphasisParser, 500
func DefaultInlineParsers() []util.PrioritizedValue {
	return []util.PrioritizedValue{
		util.Prioritized(NewCodeSpanParser(), 100),
		util.Prioritized(NewLinkParser(), 200),
		util.Prioritized(NewAutoLinkParser(), 300),
		util.Prioritized(NewRawHTMLParser(), 400),
		util.Prioritized(NewEmphasisParser(), 500),
	}
}

// DefaultParagraphTransformers returns a new list of default ParagraphTransformers.
// Priorities of default ParagraphTransformers are:
//
//	LinkReferenceParagraphTransformer, 100
func DefaultParagraphTransformers() []util.PrioritizedValue {
	return []util.PrioritizedValue{
		util.Prioritized(LinkReferenceParagraphTransformer, 100),
	}
}

// A Block struct holds a node and correspond parser pair.
type Block struct {
	// Node is a BlockNode.
	Node ast.Node
	// Parser is a BlockParser.
	Parser BlockParser
}

type parser struct {
	options               map[OptionName]interface{}
	blockParsers          [256][]BlockParser
	freeBlockParsers      []BlockParser
	inlineParsers         [256][]InlineParser
	closeBlockers         []CloseBlocker
	paragraphTransformers []ParagraphTransformer
	astTransformers       []ASTTransformer
	escapedSpace          bool
	config                *Config
	initSync              sync.Once
}

type withBlockParsers struct {
	value []util.PrioritizedValue
}

func (o *withBlockParsers) SetParserOption(c *Config) {
	c.BlockParsers = append(c.BlockParsers, o.value...)
}

// WithBlockParsers is a functional option that allow you to add
// BlockParsers to the parser.
func WithBlockParsers(bs ...util.PrioritizedValue) Option {
	return &withBlockParsers{bs}
}

type withInlineParsers struct {
	value []util.PrioritizedValue
}

func (o *withInlineParsers) SetParserOption(c *Config) {
	c.InlineParsers = append(c.InlineParsers, o.value...)
}

// WithInlineParsers is a functional option that allow you to add
// InlineParsers to the parser.
func WithInlineParsers(bs ...util.PrioritizedValue) Option {
	return &withInlineParsers{bs}
}

type withParagraphTransformers struct {
	value []util.PrioritizedValue
}

func (o *withParagraphTransformers) SetParserOption(c *Config) {
	c.ParagraphTransformers = append(c.ParagraphTransformers, o.value...)
}

// WithParagraphTransformers is a functional option that allow you to add
// ParagraphTransformers to the parser.
func WithParagraphTransformers(ps ...util.PrioritizedValue) Option {
	return &withParagraphTransformers{ps}
}

type withASTTransformers struct {
	value []util.PrioritizedValue
}

func (o *withASTTransformers) SetParserOption(c *Config) {
	c.ASTTransformers = append(c.ASTTransformers, o.value...)
}

// WithASTTransformers is a functional option that allow you to add
// ASTTransformers to the parser.
func WithASTTransformers(ps ...util.PrioritizedValue) Option {
	return &withASTTransformers{ps}
}

type withEscapedSpace struct {
}

func (o *withEscapedSpace) SetParserOption(c *Config) {
	c.EscapedSpace = true
}

// WithEscapedSpace is a functional option indicates that a '\' escaped half-space(0x20) should not trigger parsers.
func WithEscapedSpace() Option {
	return &withEscapedSpace{}
}

type withOption struct {
	name  OptionName
	value interface{}
}

func (o *withOption) SetParserOption(c *Config) {
	c.Options[o.name] = o.value
}

// WithOption is a functional option that allow you to set
// an arbitrary option to the parser.
func WithOption(name OptionName, value interface{}) Option {
	return &withOption{name, value}
}

// NewParser returns a new Parser with given options.
func NewParser(options ...Option) Parser {
	config := NewConfig()
	for _, opt := range options {
		opt.SetParserOption(config)
	}

	p := &parser{
		options: map[OptionName]interface{}{},
		config:  config,
	}

	return p
}

func (p *parser) AddOptions(opts ...Option) {
	for _, opt := range opts {
		opt.SetParserOption(p.config)
	}
}

func (p *parser) addBlockParser(v util.PrioritizedValue, options map[OptionName]interface{}) {
	bp, ok := v.Value.(BlockParser)
	if !ok {
		panic(fmt.Sprintf("%v is not a BlockParser", v.Value))
	}
	tcs := bp.Trigger()
	so, ok := v.Value.(SetOptioner)
	if ok {
		for oname, ovalue := range options {
			so.SetOption(oname, ovalue)
		}
	}
	if tcs == nil {
		p.freeBlockParsers = append(p.freeBlockParsers, bp)
	} else {
		for _, tc := range tcs {
			if p.blockParsers[tc] == nil {
				p.blockParsers[tc] = []BlockParser{}
			}
			p.blockParsers[tc] = append(p.blockParsers[tc], bp)
		}
	}
}

func (p *parser) addInlineParser(v util.PrioritizedValue, options map[OptionName]interface{}) {
	ip, ok := v.Value.(InlineParser)
	if !ok {
		panic(fmt.Sprintf("%v is not a InlineParser", v.Value))
	}
	tcs := ip.Trigger()
	so, ok := v.Value.(SetOptioner)
	if ok {
		for oname, ovalue := range options {
			so.SetOption(oname, ovalue)
		}
	}
	if cb, ok := ip.(CloseBlocker); ok {
		p.closeBlockers = append(p.closeBlockers, cb)
	}
	for _, tc := range tcs {
		if p.inlineParsers[tc] == nil {
			p.inlineParsers[tc] = []InlineParser{}
		}
		p.inlineParsers[tc] = append(p.inlineParsers[tc], ip)
	}
}

func (p *parser) addParagraphTransformer(v util.PrioritizedValue, options map[OptionName]interface{}) {
	pt, ok := v.Value.(ParagraphTransformer)
	if !ok {
		panic(fmt.Sprintf("%v is not a ParagraphTransformer", v.Value))
	}
	so, ok := v.Value.(SetOptioner)
	if ok {
		for oname, ovalue := range options {
			so.SetOption(oname, ovalue)
		}
	}
	p.paragraphTransformers = append(p.paragraphTransformers, pt)
}

func (p *parser) addASTTransformer(v util.PrioritizedValue, options map[OptionName]interface{}) {
	at, ok := v.Value.(ASTTransformer)
	if !ok {
		panic(fmt.Sprintf("%v is not a ASTTransformer", v.Value))
	}
	so, ok := v.Value.(SetOptioner)
	if ok {
		for oname, ovalue := range options {
			so.SetOption(oname, ovalue)
		}
	}
	p.astTransformers = append(p.astTransformers, at)
}

// A ParseConfig struct is a data structure that holds configuration of the Parser.Parse.
type ParseConfig struct {
	Context Context
}

// A ParseOption is a functional option type for the Parser.Parse.
type ParseOption func(c *ParseConfig)

// WithContext is a functional option that allow you to override
// a default context.
func WithContext(context Context) ParseOption {
	return func(c *ParseConfig) {
		c.Context = context
	}
}

func (p *parser) Parse(reader text.Reader, opts ...ParseOption) ast.Node {
	p.initSync.Do(func() {
		p.config.BlockParsers.Sort()
		for _, v := range p.config.BlockParsers {
			p.addBlockParser(v, p.config.Options)
		}
		for i := range p.blockParsers {
			if p.blockParsers[i] != nil {
				p.blockParsers[i] = append(p.blockParsers[i], p.freeBlockParsers...)
			}
		}

		p.config.InlineParsers.Sort()
		for _, v := range p.config.InlineParsers {
			p.addInlineParser(v, p.config.Options)
		}
		p.config.ParagraphTransformers.Sort()
		for _, v := range p.config.ParagraphTransformers {
			p.addParagraphTransformer(v, p.config.Options)
		}
		p.config.ASTTransformers.Sort()
		for _, v := range p.config.ASTTransformers {
			p.addASTTransformer(v, p.config.Options)
		}
		p.escapedSpace = p.config.EscapedSpace
		p.config = nil
	})
	c := &ParseConfig{}
	for _, opt := range opts {
		opt(c)
	}
	if c.Context == nil {
		c.Context = NewContext()
	}
	pc := c.Context
	root := ast.NewDocument()
	p.parseBlocks(root, reader, pc)

	blockReader := text.NewBlockReader(reader.Source(), nil)
	p.walkBlock(root, func(node ast.Node) {
		p.parseBlock(blockReader, node, pc)
	})
	for _, at := range p.astTransformers {
		at.Transform(root, reader, pc)
	}
	// root.Dump(reader.Source(), 0)
	return root
}

func (p *parser) transformParagraph(node *ast.Paragraph, reader text.Reader, pc Context) bool {
	for _, pt := range p.paragraphTransformers {
		pt.Transform(node, reader, pc)
		if node.Parent() == nil {
			return true
		}
	}
	return false
}

func (p *parser) closeBlocks(from, to int, reader text.Reader, pc Context) {
	blocks := pc.OpenedBlocks()
	for i := from; i >= to; i-- {
		node := blocks[i].Node
		blocks[i].Parser.Close(blocks[i].Node, reader, pc)
		paragraph, ok := node.(*ast.Paragraph)
		if ok && node.Parent() != nil {
			p.transformParagraph(paragraph, reader, pc)
		}
	}
	if from == len(blocks)-1 {
		blocks = blocks[0:to]
	} else {
		blocks = append(blocks[0:to], blocks[from+1:]...)
	}
	pc.SetOpenedBlocks(blocks)
}

type blockOpenResult int

const (
	paragraphContinuation blockOpenResult = iota + 1
	newBlocksOpened
	noBlocksOpened
)

func (p *parser) openBlocks(parent ast.Node, blankLine bool, reader text.Reader, pc Context) blockOpenResult {
	result := blockOpenResult(noBlocksOpened)
	continuable := false
	lastBlock := pc.LastOpenedBlock()
	if lastBlock.Node != nil {
		continuable = ast.IsParagraph(lastBlock.Node)
	}
retry:
	var bps []BlockParser
	line, _ := reader.PeekLine()
	w, pos := util.IndentWidth(line, reader.LineOffset())
	if w >= len(line) {
		pc.SetBlockOffset(-1)
		pc.SetBlockIndent(-1)
	} else {
		pc.SetBlockOffset(pos)
		pc.SetBlockIndent(w)
	}
	if line == nil || line[0] == '\n' {
		goto continuable
	}
	bps = p.freeBlockParsers
	if pos < len(line) {
		bps = p.blockParsers[line[pos]]
		if bps == nil {
			bps = p.freeBlockParsers
		}
	}
	if bps == nil {
		goto continuable
	}

	for _, bp := range bps {
		if continuable && result == noBlocksOpened && !bp.CanInterruptParagraph() {
			continue
		}
		if w > 3 && !bp.CanAcceptIndentedLine() {
			continue
		}
		lastBlock = pc.LastOpenedBlock()
		last := lastBlock.Node
		node, state := bp.Open(parent, reader, pc)
		if node != nil {
			// Parser requires last node to be a paragraph.
			// With table extension:
			//
			//     0
			//     -:
			//     -
			//
			// '-' on 3rd line seems a Setext heading because 1st and 2nd lines
			// are being paragraph when the Settext heading parser tries to parse the 3rd
			// line.
			// But 1st line and 2nd line are a table. Thus this paragraph will be transformed
			// by a paragraph transformer. So this text should be converted to a table and
			// an empty list.
			if state&RequireParagraph != 0 {
				if last == parent.LastChild() {
					// Opened paragraph may be transformed by ParagraphTransformers in
					// closeBlocks().
					lastBlock.Parser.Close(last, reader, pc)
					blocks := pc.OpenedBlocks()
					pc.SetOpenedBlocks(blocks[0 : len(blocks)-1])
					if p.transformParagraph(last.(*ast.Paragraph), reader, pc) {
						// Paragraph has been transformed.
						// So this parser is considered as failing.
						continuable = false
						goto retry
					}
				}
			}
			node.SetBlankPreviousLines(blankLine)
			if last != nil && last.Parent() == nil {
				lastPos := len(pc.OpenedBlocks()) - 1
				p.closeBlocks(lastPos, lastPos, reader, pc)
			}
			parent.AppendChild(parent, node)
			result = newBlocksOpened
			be := Block{node, bp}
			pc.SetOpenedBlocks(append(pc.OpenedBlocks(), be))
			if state&HasChildren != 0 {
				parent = node
				goto retry // try child block
			}
			break // no children, can not open more blocks on this line
		}
	}

continuable:
	if result == noBlocksOpened && continuable {
		state := lastBlock.Parser.Continue(lastBlock.Node, reader, pc)
		if state&Continue != 0 {
			result = paragraphContinuation
		}
	}
	return result
}

type lineStat struct {
	lineNum int
	level   int
	isBlank bool
}

func isBlankLine(lineNum, level int, stats []lineStat) bool {
	ret := true
	for i := len(stats) - 1 - level; i >= 0; i-- {
		ret = false
		s := stats[i]
		if s.lineNum == lineNum {
			if s.level < level && s.isBlank {
				return true
			} else if s.level == level {
				return s.isBlank
			}
		}
		if s.lineNum < lineNum {
			return ret
		}
	}
	return ret
}

func (p *parser) parseBlocks(parent ast.Node, reader text.Reader, pc Context) {
	pc.SetOpenedBlocks([]Block{})
	blankLines := make([]lineStat, 0, 128)
	isBlank := false
	for { // process blocks separated by blank lines
		_, lines, ok := reader.SkipBlankLines()
		if !ok {
			return
		}
		lineNum, _ := reader.Position()
		if lines != 0 {
			blankLines = blankLines[0:0]
			l := len(pc.OpenedBlocks())
			for i := 0; i < l; i++ {
				blankLines = append(blankLines, lineStat{lineNum - 1, i, lines != 0})
			}
		}
		isBlank = isBlankLine(lineNum-1, 0, blankLines)
		// first, we try to open blocks
		if p.openBlocks(parent, isBlank, reader, pc) != newBlocksOpened {
			return
		}
		reader.AdvanceLine()
		for { // process opened blocks line by line
			openedBlocks := pc.OpenedBlocks()
			l := len(openedBlocks)
			if l == 0 {
				break
			}
			lastIndex := l - 1
			for i := 0; i < l; i++ {
				be := openedBlocks[i]
				line, _ := reader.PeekLine()
				if line == nil {
					p.closeBlocks(lastIndex, 0, reader, pc)
					reader.AdvanceLine()
					return
				}
				lineNum, _ := reader.Position()
				blankLines = append(blankLines, lineStat{lineNum, i, util.IsBlank(line)})
				// If node is a paragraph, p.openBlocks determines whether it is continuable.
				// So we do not process paragraphs here.
				if !ast.IsParagraph(be.Node) {
					state := be.Parser.Continue(be.Node, reader, pc)
					if state&Continue != 0 {
						// When current node is a container block and has no children,
						// we try to open new child nodes
						if state&HasChildren != 0 && i == lastIndex {
							isBlank = isBlankLine(lineNum-1, i, blankLines)
							p.openBlocks(be.Node, isBlank, reader, pc)
							break
						}
						continue
					}
				}
				// current node may be closed or lazy continuation
				isBlank = isBlankLine(lineNum-1, i, blankLines)
				thisParent := parent
				if i != 0 {
					thisParent = openedBlocks[i-1].Node
				}
				lastNode := openedBlocks[lastIndex].Node
				result := p.openBlocks(thisParent, isBlank, reader, pc)
				if result != paragraphContinuation {
					// lastNode is a paragraph and was transformed by the paragraph
					// transformers.
					if openedBlocks[lastIndex].Node != lastNode {
						lastIndex--
					}
					p.closeBlocks(lastIndex, i, reader, pc)
				}
				break
			}

			reader.AdvanceLine()
		}
	}
}

func (p *parser) walkBlock(block ast.Node, cb func(node ast.Node)) {
	for c := block.FirstChild(); c != nil; c = c.NextSibling() {
		p.walkBlock(c, cb)
	}
	cb(block)
}

const (
	lineBreakHard uint8 = 1 << iota
	lineBreakSoft
	lineBreakVisible
)

func (p *parser) parseBlock(block text.BlockReader, parent ast.Node, pc Context) {
	if parent.IsRaw() {
		return
	}
	escaped := false
	source := block.Source()
	block.Reset(parent.Lines())
	for {
	retry:
		line, _ := block.PeekLine()
		if line == nil {
			break
		}
		lineLength := len(line)
		var lineBreakFlags uint8 = 0
		hasNewLine := line[lineLength-1] == '\n'
		if ((lineLength >= 3 && line[lineLength-2] == '\\' && line[lineLength-3] != '\\') || (lineLength == 2 && line[lineLength-2] == '\\')) && hasNewLine { // ends with \\n
			lineLength -= 2
			lineBreakFlags |= lineBreakHard | lineBreakVisible
		} else if ((lineLength >= 4 && line[lineLength-3] == '\\' && line[lineLength-2] == '\r' && line[lineLength-4] != '\\') || (lineLength == 3 && line[lineLength-3] == '\\' && line[lineLength-2] == '\r')) && hasNewLine { // ends with \\r\n
			lineLength -= 3
			lineBreakFlags |= lineBreakHard | lineBreakVisible
		} else if lineLength >= 3 && line[lineLength-3] == ' ' && line[lineLength-2] == ' ' && hasNewLine { // ends with [space][space]\n
			lineLength -= 3
			lineBreakFlags |= lineBreakHard
		} else if lineLength >= 4 && line[lineLength-4] == ' ' && line[lineLength-3] == ' ' && line[lineLength-2] == '\r' && hasNewLine { // ends with [space][space]\r\n
			lineLength -= 4
			lineBreakFlags |= lineBreakHard
		} else if hasNewLine {
			// If the line ends with a newline character, but it is not a hardlineBreak, then it is a softLinebreak
			// If the line ends with a hardlineBreak, then it cannot end with a softLinebreak
			// See https://spec.commonmark.org/0.30/#soft-line-breaks
			lineBreakFlags |= lineBreakSoft
		}

		l, startPosition := block.Position()
		n := 0
		for i := 0; i < lineLength; i++ {
			c := line[i]
			if c == '\n' {
				break
			}
			isSpace := util.IsSpace(c) && c != '\r' && c != '\n'
			isPunct := util.IsPunct(c)
			if (isPunct && !escaped) || isSpace && !(escaped && p.escapedSpace) || i == 0 {
				parserChar := c
				if isSpace || (i == 0 && !isPunct) {
					parserChar = ' '
				}
				ips := p.inlineParsers[parserChar]
				if ips != nil {
					block.Advance(n)
					n = 0
					savedLine, savedPosition := block.Position()
					if i != 0 {
						_, currentPosition := block.Position()
						ast.MergeOrAppendTextSegment(parent, startPosition.Between(currentPosition))
						_, startPosition = block.Position()
					}
					var inlineNode ast.Node
					for _, ip := range ips {
						inlineNode = ip.Parse(parent, block, pc)
						if inlineNode != nil {
							break
						}
						block.SetPosition(savedLine, savedPosition)
					}
					if inlineNode != nil {
						parent.AppendChild(parent, inlineNode)
						goto retry
					}
				}
			}
			if escaped {
				escaped = false
				n++
				continue
			}

			if c == '\\' {
				escaped = true
				n++
				continue
			}

			escaped = false
			n++
		}
		if n != 0 {
			block.Advance(n)
		}
		currentL, currentPosition := block.Position()
		if l != currentL {
			continue
		}
		diff := startPosition.Between(currentPosition)
		var text *ast.Text
		if lineBreakFlags&(lineBreakHard|lineBreakVisible) == lineBreakHard|lineBreakVisible {
			text = ast.NewTextSegment(diff)
		} else {
			text = ast.NewTextSegment(diff.TrimRightSpace(source))
		}
		text.SetSoftLineBreak(lineBreakFlags&lineBreakSoft != 0)
		text.SetHardLineBreak(lineBreakFlags&lineBreakHard != 0)
		parent.AppendChild(parent, text)
		block.AdvanceLine()
	}

	ProcessDelimiters(nil, pc)
	for _, ip := range p.closeBlockers {
		ip.CloseBlock(parent, block, pc)
	}
}
