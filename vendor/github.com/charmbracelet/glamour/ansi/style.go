package ansi

// Chroma holds all the chroma settings.
type Chroma struct {
	Text                StylePrimitive `json:"text,omitempty"`
	Error               StylePrimitive `json:"error,omitempty"`
	Comment             StylePrimitive `json:"comment,omitempty"`
	CommentPreproc      StylePrimitive `json:"comment_preproc,omitempty"`
	Keyword             StylePrimitive `json:"keyword,omitempty"`
	KeywordReserved     StylePrimitive `json:"keyword_reserved,omitempty"`
	KeywordNamespace    StylePrimitive `json:"keyword_namespace,omitempty"`
	KeywordType         StylePrimitive `json:"keyword_type,omitempty"`
	Operator            StylePrimitive `json:"operator,omitempty"`
	Punctuation         StylePrimitive `json:"punctuation,omitempty"`
	Name                StylePrimitive `json:"name,omitempty"`
	NameBuiltin         StylePrimitive `json:"name_builtin,omitempty"`
	NameTag             StylePrimitive `json:"name_tag,omitempty"`
	NameAttribute       StylePrimitive `json:"name_attribute,omitempty"`
	NameClass           StylePrimitive `json:"name_class,omitempty"`
	NameConstant        StylePrimitive `json:"name_constant,omitempty"`
	NameDecorator       StylePrimitive `json:"name_decorator,omitempty"`
	NameException       StylePrimitive `json:"name_exception,omitempty"`
	NameFunction        StylePrimitive `json:"name_function,omitempty"`
	NameOther           StylePrimitive `json:"name_other,omitempty"`
	Literal             StylePrimitive `json:"literal,omitempty"`
	LiteralNumber       StylePrimitive `json:"literal_number,omitempty"`
	LiteralDate         StylePrimitive `json:"literal_date,omitempty"`
	LiteralString       StylePrimitive `json:"literal_string,omitempty"`
	LiteralStringEscape StylePrimitive `json:"literal_string_escape,omitempty"`
	GenericDeleted      StylePrimitive `json:"generic_deleted,omitempty"`
	GenericEmph         StylePrimitive `json:"generic_emph,omitempty"`
	GenericInserted     StylePrimitive `json:"generic_inserted,omitempty"`
	GenericStrong       StylePrimitive `json:"generic_strong,omitempty"`
	GenericSubheading   StylePrimitive `json:"generic_subheading,omitempty"`
	Background          StylePrimitive `json:"background,omitempty"`
}

// StylePrimitive holds all the basic style settings.
type StylePrimitive struct {
	BlockPrefix     string  `json:"block_prefix,omitempty"`
	BlockSuffix     string  `json:"block_suffix,omitempty"`
	Prefix          string  `json:"prefix,omitempty"`
	Suffix          string  `json:"suffix,omitempty"`
	Color           *string `json:"color,omitempty"`
	BackgroundColor *string `json:"background_color,omitempty"`
	Underline       *bool   `json:"underline,omitempty"`
	Bold            *bool   `json:"bold,omitempty"`
	Upper           *bool   `json:"upper,omitempty"`
	Lower           *bool   `json:"lower,omitempty"`
	Title           *bool   `json:"title,omitempty"`
	Italic          *bool   `json:"italic,omitempty"`
	CrossedOut      *bool   `json:"crossed_out,omitempty"`
	Faint           *bool   `json:"faint,omitempty"`
	Conceal         *bool   `json:"conceal,omitempty"`
	Overlined       *bool   `json:"overlined,omitempty"`
	Inverse         *bool   `json:"inverse,omitempty"`
	Blink           *bool   `json:"blink,omitempty"`
	Format          string  `json:"format,omitempty"`
}

// StyleTask holds the style settings for a task item.
type StyleTask struct {
	StylePrimitive
	Ticked   string `json:"ticked,omitempty"`
	Unticked string `json:"unticked,omitempty"`
}

// StyleBlock holds the basic style settings for block elements.
type StyleBlock struct {
	StylePrimitive
	Indent      *uint   `json:"indent,omitempty"`
	IndentToken *string `json:"indent_token,omitempty"`
	Margin      *uint   `json:"margin,omitempty"`
}

// StyleCodeBlock holds the style settings for a code block.
type StyleCodeBlock struct {
	StyleBlock
	Theme  string  `json:"theme,omitempty"`
	Chroma *Chroma `json:"chroma,omitempty"`
}

// StyleList holds the style settings for a list.
type StyleList struct {
	StyleBlock
	LevelIndent uint `json:"level_indent,omitempty"`
}

// StyleTable holds the style settings for a table.
type StyleTable struct {
	StyleBlock
	CenterSeparator *string `json:"center_separator,omitempty"`
	ColumnSeparator *string `json:"column_separator,omitempty"`
	RowSeparator    *string `json:"row_separator,omitempty"`
}

// StyleConfig is used to configure the styling behavior of an ANSIRenderer.
type StyleConfig struct {
	Document   StyleBlock `json:"document,omitempty"`
	BlockQuote StyleBlock `json:"block_quote,omitempty"`
	Paragraph  StyleBlock `json:"paragraph,omitempty"`
	List       StyleList  `json:"list,omitempty"`

	Heading StyleBlock `json:"heading,omitempty"`
	H1      StyleBlock `json:"h1,omitempty"`
	H2      StyleBlock `json:"h2,omitempty"`
	H3      StyleBlock `json:"h3,omitempty"`
	H4      StyleBlock `json:"h4,omitempty"`
	H5      StyleBlock `json:"h5,omitempty"`
	H6      StyleBlock `json:"h6,omitempty"`

	Text           StylePrimitive `json:"text,omitempty"`
	Strikethrough  StylePrimitive `json:"strikethrough,omitempty"`
	Emph           StylePrimitive `json:"emph,omitempty"`
	Strong         StylePrimitive `json:"strong,omitempty"`
	HorizontalRule StylePrimitive `json:"hr,omitempty"`

	Item        StylePrimitive `json:"item,omitempty"`
	Enumeration StylePrimitive `json:"enumeration,omitempty"`
	Task        StyleTask      `json:"task,omitempty"`

	Link     StylePrimitive `json:"link,omitempty"`
	LinkText StylePrimitive `json:"link_text,omitempty"`

	Image     StylePrimitive `json:"image,omitempty"`
	ImageText StylePrimitive `json:"image_text,omitempty"`

	Code      StyleBlock     `json:"code,omitempty"`
	CodeBlock StyleCodeBlock `json:"code_block,omitempty"`

	Table StyleTable `json:"table,omitempty"`

	DefinitionList        StyleBlock     `json:"definition_list,omitempty"`
	DefinitionTerm        StylePrimitive `json:"definition_term,omitempty"`
	DefinitionDescription StylePrimitive `json:"definition_description,omitempty"`

	HTMLBlock StyleBlock `json:"html_block,omitempty"`
	HTMLSpan  StyleBlock `json:"html_span,omitempty"`
}

func cascadeStyles(toBlock bool, s ...StyleBlock) StyleBlock {
	var r StyleBlock

	for _, v := range s {
		r = cascadeStyle(r, v, toBlock)
	}
	return r
}

func cascadeStyle(parent StyleBlock, child StyleBlock, toBlock bool) StyleBlock {
	s := child

	s.Color = parent.Color
	s.BackgroundColor = parent.BackgroundColor
	s.Underline = parent.Underline
	s.Bold = parent.Bold
	s.Upper = parent.Upper
	s.Title = parent.Title
	s.Lower = parent.Lower
	s.Italic = parent.Italic
	s.CrossedOut = parent.CrossedOut
	s.Faint = parent.Faint
	s.Conceal = parent.Conceal
	s.Overlined = parent.Overlined
	s.Inverse = parent.Inverse
	s.Blink = parent.Blink

	if toBlock {
		s.Indent = parent.Indent
		s.Margin = parent.Margin
		s.BlockPrefix = parent.BlockPrefix
		s.BlockSuffix = parent.BlockSuffix
		s.Prefix = parent.Prefix
		s.Suffix = parent.Suffix
	}

	if child.Color != nil {
		s.Color = child.Color
	}
	if child.BackgroundColor != nil {
		s.BackgroundColor = child.BackgroundColor
	}
	if child.Indent != nil {
		s.Indent = child.Indent
	}
	if child.Margin != nil {
		s.Margin = child.Margin
	}
	if child.Underline != nil {
		s.Underline = child.Underline
	}
	if child.Bold != nil {
		s.Bold = child.Bold
	}
	if child.Upper != nil {
		s.Upper = child.Upper
	}
	if child.Lower != nil {
		s.Lower = child.Lower
	}
	if child.Title != nil {
		s.Title = child.Title
	}
	if child.Italic != nil {
		s.Italic = child.Italic
	}
	if child.CrossedOut != nil {
		s.CrossedOut = child.CrossedOut
	}
	if child.Faint != nil {
		s.Faint = child.Faint
	}
	if child.Conceal != nil {
		s.Conceal = child.Conceal
	}
	if child.Overlined != nil {
		s.Overlined = child.Overlined
	}
	if child.Inverse != nil {
		s.Inverse = child.Inverse
	}
	if child.Blink != nil {
		s.Blink = child.Blink
	}
	if child.BlockPrefix != "" {
		s.BlockPrefix = child.BlockPrefix
	}
	if child.BlockSuffix != "" {
		s.BlockSuffix = child.BlockSuffix
	}
	if child.Prefix != "" {
		s.Prefix = child.Prefix
	}
	if child.Suffix != "" {
		s.Suffix = child.Suffix
	}
	if child.Format != "" {
		s.Format = child.Format
	}

	return s
}
