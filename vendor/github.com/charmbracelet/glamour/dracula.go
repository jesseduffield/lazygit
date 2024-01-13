package glamour

import "github.com/charmbracelet/glamour/ansi"

var DraculaStyleConfig = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
			Color:       stringPtr("#f8f8f2"),
		},
		Margin: uintPtr(2),
	},
	BlockQuote: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color:  stringPtr("#f1fa8c"),
			Italic: boolPtr(true),
		},
		Indent: uintPtr(2),
	},
	List: ansi.StyleList{
		LevelIndent: 2,
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: stringPtr("#f8f8f2"),
			},
		},
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       stringPtr("#bd93f9"),
			Bold:        boolPtr(true),
		},
	},
	H1: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "# ",
		},
	},
	H2: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "## ",
		},
	},
	H3: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "### ",
		},
	},
	H4: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "#### ",
		},
	},
	H5: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "##### ",
		},
	},
	H6: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Prefix: "###### ",
		},
	},
	Strikethrough: ansi.StylePrimitive{
		CrossedOut: boolPtr(true),
	},
	Emph: ansi.StylePrimitive{
		Color:  stringPtr("#f1fa8c"),
		Italic: boolPtr(true),
	},
	Strong: ansi.StylePrimitive{
		Bold:  boolPtr(true),
		Color: stringPtr("#ffb86c"),
	},
	HorizontalRule: ansi.StylePrimitive{
		Color:  stringPtr("#6272A4"),
		Format: "\n--------\n",
	},
	Item: ansi.StylePrimitive{
		BlockPrefix: "• ",
	},
	Enumeration: ansi.StylePrimitive{
		BlockPrefix: ". ",
		Color:       stringPtr("#8be9fd"),
	},
	Task: ansi.StyleTask{
		StylePrimitive: ansi.StylePrimitive{},
		Ticked:         "[✓] ",
		Unticked:       "[ ] ",
	},
	Link: ansi.StylePrimitive{
		Color:     stringPtr("#8be9fd"),
		Underline: boolPtr(true),
	},
	LinkText: ansi.StylePrimitive{
		Color: stringPtr("#ff79c6"),
	},
	Image: ansi.StylePrimitive{
		Color:     stringPtr("#8be9fd"),
		Underline: boolPtr(true),
	},
	ImageText: ansi.StylePrimitive{
		Color:  stringPtr("#ff79c6"),
		Format: "Image: {{.text}} →",
	},
	Code: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			Color: stringPtr("#50fa7b"),
		},
	},
	CodeBlock: ansi.StyleCodeBlock{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: stringPtr("#ffb86c"),
			},
			Margin: uintPtr(2),
		},
		Chroma: &ansi.Chroma{
			Text: ansi.StylePrimitive{
				Color: stringPtr("#f8f8f2"),
			},
			Error: ansi.StylePrimitive{
				Color:           stringPtr("#f8f8f2"),
				BackgroundColor: stringPtr("#ff5555"),
			},
			Comment: ansi.StylePrimitive{
				Color: stringPtr("#6272A4"),
			},
			CommentPreproc: ansi.StylePrimitive{
				Color: stringPtr("#ff79c6"),
			},
			Keyword: ansi.StylePrimitive{
				Color: stringPtr("#ff79c6"),
			},
			KeywordReserved: ansi.StylePrimitive{
				Color: stringPtr("#ff79c6"),
			},
			KeywordNamespace: ansi.StylePrimitive{
				Color: stringPtr("#ff79c6"),
			},
			KeywordType: ansi.StylePrimitive{
				Color: stringPtr("#8be9fd"),
			},
			Operator: ansi.StylePrimitive{
				Color: stringPtr("#ff79c6"),
			},
			Punctuation: ansi.StylePrimitive{
				Color: stringPtr("#f8f8f2"),
			},
			Name: ansi.StylePrimitive{
				Color: stringPtr("#8be9fd"),
			},
			NameBuiltin: ansi.StylePrimitive{
				Color: stringPtr("#8be9fd"),
			},
			NameTag: ansi.StylePrimitive{
				Color: stringPtr("#ff79c6"),
			},
			NameAttribute: ansi.StylePrimitive{
				Color: stringPtr("#50fa7b"),
			},
			NameClass: ansi.StylePrimitive{
				Color: stringPtr("#8be9fd"),
			},
			NameConstant: ansi.StylePrimitive{
				Color: stringPtr("#bd93f9"),
			},
			NameDecorator: ansi.StylePrimitive{
				Color: stringPtr("#50fa7b"),
			},
			NameFunction: ansi.StylePrimitive{
				Color: stringPtr("#50fa7b"),
			},
			LiteralNumber: ansi.StylePrimitive{
				Color: stringPtr("#6EEFC0"),
			},
			LiteralString: ansi.StylePrimitive{
				Color: stringPtr("#f1fa8c"),
			},
			LiteralStringEscape: ansi.StylePrimitive{
				Color: stringPtr("#ff79c6"),
			},
			GenericDeleted: ansi.StylePrimitive{
				Color: stringPtr("#ff5555"),
			},
			GenericEmph: ansi.StylePrimitive{
				Color:  stringPtr("#f1fa8c"),
				Italic: boolPtr(true),
			},
			GenericInserted: ansi.StylePrimitive{
				Color: stringPtr("#50fa7b"),
			},
			GenericStrong: ansi.StylePrimitive{
				Color: stringPtr("#ffb86c"),
				Bold:  boolPtr(true),
			},
			GenericSubheading: ansi.StylePrimitive{
				Color: stringPtr("#bd93f9"),
			},
			Background: ansi.StylePrimitive{
				BackgroundColor: stringPtr("#282a36"),
			},
		},
	},
	Table: ansi.StyleTable{
		StyleBlock: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{},
		},
		CenterSeparator: stringPtr("┼"),
		ColumnSeparator: stringPtr("│"),
		RowSeparator:    stringPtr("─"),
	},
	DefinitionDescription: ansi.StylePrimitive{
		BlockPrefix: "\n🠶 ",
	},
}
