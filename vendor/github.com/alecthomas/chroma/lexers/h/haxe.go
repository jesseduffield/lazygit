package h

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Haxe lexer.
var Haxe = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Haxe",
		Aliases:   []string{"hx", "haxe", "hxsl"},
		Filenames: []string{"*.hx", "*.hxsl"},
		MimeTypes: []string{"text/haxe", "text/x-haxe", "text/x-hx"},
		DotAll:    true,
	},
	haxeRules,
))

func haxeRules() Rules {
	return Rules{
		"root": {
			Include("spaces"),
			Include("meta"),
			{`(?:package)\b`, KeywordNamespace, Push("semicolon", "package")},
			{`(?:import)\b`, KeywordNamespace, Push("semicolon", "import")},
			{`(?:using)\b`, KeywordNamespace, Push("semicolon", "using")},
			{`(?:extern|private)\b`, KeywordDeclaration, nil},
			{`(?:abstract)\b`, KeywordDeclaration, Push("abstract")},
			{`(?:class|interface)\b`, KeywordDeclaration, Push("class")},
			{`(?:enum)\b`, KeywordDeclaration, Push("enum")},
			{`(?:typedef)\b`, KeywordDeclaration, Push("typedef")},
			{`(?=.)`, Text, Push("expr-statement")},
		},
		"spaces": {
			{`\s+`, Text, nil},
			{`//[^\n\r]*`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`(#)(if|elseif|else|end|error)\b`, CommentPreproc, MutatorFunc(haxePreProcMutator)},
		},
		"string-single-interpol": {
			{`\$\{`, LiteralStringInterpol, Push("string-interpol-close", "expr")},
			{`\$\$`, LiteralStringEscape, nil},
			{`\$(?=(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+))`, LiteralStringInterpol, Push("ident")},
			Include("string-single"),
		},
		"string-single": {
			{`'`, LiteralStringSingle, Pop(1)},
			{`\\.`, LiteralStringEscape, nil},
			{`.`, LiteralStringSingle, nil},
		},
		"string-double": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`\\.`, LiteralStringEscape, nil},
			{`.`, LiteralStringDouble, nil},
		},
		"string-interpol-close": {
			{`\$(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, LiteralStringInterpol, nil},
			{`\}`, LiteralStringInterpol, Pop(1)},
		},
		"package": {
			Include("spaces"),
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, NameNamespace, nil},
			{`\.`, Punctuation, Push("import-ident")},
			Default(Pop(1)),
		},
		"import": {
			Include("spaces"),
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, NameNamespace, nil},
			{`\*`, Keyword, nil},
			{`\.`, Punctuation, Push("import-ident")},
			{`in`, KeywordNamespace, Push("ident")},
			Default(Pop(1)),
		},
		"import-ident": {
			Include("spaces"),
			{`\*`, Keyword, Pop(1)},
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, NameNamespace, Pop(1)},
		},
		"using": {
			Include("spaces"),
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, NameNamespace, nil},
			{`\.`, Punctuation, Push("import-ident")},
			Default(Pop(1)),
		},
		"preproc-error": {
			{`\s+`, CommentPreproc, nil},
			{`'`, LiteralStringSingle, Push("#pop", "string-single")},
			{`"`, LiteralStringDouble, Push("#pop", "string-double")},
			Default(Pop(1)),
		},
		"preproc-expr": {
			{`\s+`, CommentPreproc, nil},
			{`\!`, CommentPreproc, nil},
			{`\(`, CommentPreproc, Push("#pop", "preproc-parenthesis")},
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, CommentPreproc, Pop(1)},
			{`\.[0-9]+`, LiteralNumberFloat, nil},
			{`[0-9]+[eE][+\-]?[0-9]+`, LiteralNumberFloat, nil},
			{`[0-9]+\.[0-9]*[eE][+\-]?[0-9]+`, LiteralNumberFloat, nil},
			{`[0-9]+\.[0-9]+`, LiteralNumberFloat, nil},
			{`[0-9]+\.(?!(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)|\.\.)`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`'`, LiteralStringSingle, Push("#pop", "string-single")},
			{`"`, LiteralStringDouble, Push("#pop", "string-double")},
		},
		"preproc-parenthesis": {
			{`\s+`, CommentPreproc, nil},
			{`\)`, CommentPreproc, Pop(1)},
			Default(Push("preproc-expr-in-parenthesis")),
		},
		"preproc-expr-chain": {
			{`\s+`, CommentPreproc, nil},
			{`(?:%=|&=|\|=|\^=|\+=|\-=|\*=|/=|<<=|>\s*>\s*=|>\s*>\s*>\s*=|==|!=|<=|>\s*=|&&|\|\||<<|>>>|>\s*>|\.\.\.|<|>|%|&|\||\^|\+|\*|/|\-|=>|=)`, CommentPreproc, Push("#pop", "preproc-expr-in-parenthesis")},
			Default(Pop(1)),
		},
		"preproc-expr-in-parenthesis": {
			{`\s+`, CommentPreproc, nil},
			{`\!`, CommentPreproc, nil},
			{`\(`, CommentPreproc, Push("#pop", "preproc-expr-chain", "preproc-parenthesis")},
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, CommentPreproc, Push("#pop", "preproc-expr-chain")},
			{`\.[0-9]+`, LiteralNumberFloat, Push("#pop", "preproc-expr-chain")},
			{`[0-9]+[eE][+\-]?[0-9]+`, LiteralNumberFloat, Push("#pop", "preproc-expr-chain")},
			{`[0-9]+\.[0-9]*[eE][+\-]?[0-9]+`, LiteralNumberFloat, Push("#pop", "preproc-expr-chain")},
			{`[0-9]+\.[0-9]+`, LiteralNumberFloat, Push("#pop", "preproc-expr-chain")},
			{`[0-9]+\.(?!(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)|\.\.)`, LiteralNumberFloat, Push("#pop", "preproc-expr-chain")},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, Push("#pop", "preproc-expr-chain")},
			{`[0-9]+`, LiteralNumberInteger, Push("#pop", "preproc-expr-chain")},
			{`'`, LiteralStringSingle, Push("#pop", "preproc-expr-chain", "string-single")},
			{`"`, LiteralStringDouble, Push("#pop", "preproc-expr-chain", "string-double")},
		},
		"abstract": {
			Include("spaces"),
			Default(Pop(1), Push("abstract-body"), Push("abstract-relation"), Push("abstract-opaque"), Push("type-param-constraint"), Push("type-name")),
		},
		"abstract-body": {
			Include("spaces"),
			{`\{`, Punctuation, Push("#pop", "class-body")},
		},
		"abstract-opaque": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "parenthesis-close", "type")},
			Default(Pop(1)),
		},
		"abstract-relation": {
			Include("spaces"),
			{`(?:to|from)`, KeywordDeclaration, Push("type")},
			{`,`, Punctuation, nil},
			Default(Pop(1)),
		},
		"meta": {
			Include("spaces"),
			{`@`, NameDecorator, Push("meta-body", "meta-ident", "meta-colon")},
		},
		"meta-colon": {
			Include("spaces"),
			{`:`, NameDecorator, Pop(1)},
			Default(Pop(1)),
		},
		"meta-ident": {
			Include("spaces"),
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, NameDecorator, Pop(1)},
		},
		"meta-body": {
			Include("spaces"),
			{`\(`, NameDecorator, Push("#pop", "meta-call")},
			Default(Pop(1)),
		},
		"meta-call": {
			Include("spaces"),
			{`\)`, NameDecorator, Pop(1)},
			Default(Pop(1), Push("meta-call-sep"), Push("expr")),
		},
		"meta-call-sep": {
			Include("spaces"),
			{`\)`, NameDecorator, Pop(1)},
			{`,`, Punctuation, Push("#pop", "meta-call")},
		},
		"typedef": {
			Include("spaces"),
			Default(Pop(1), Push("typedef-body"), Push("type-param-constraint"), Push("type-name")),
		},
		"typedef-body": {
			Include("spaces"),
			{`=`, Operator, Push("#pop", "optional-semicolon", "type")},
		},
		"enum": {
			Include("spaces"),
			Default(Pop(1), Push("enum-body"), Push("bracket-open"), Push("type-param-constraint"), Push("type-name")),
		},
		"enum-body": {
			Include("spaces"),
			Include("meta"),
			{`\}`, Punctuation, Pop(1)},
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Push("enum-member", "type-param-constraint")},
		},
		"enum-member": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "semicolon", "flag", "function-param")},
			Default(Pop(1), Push("semicolon"), Push("flag")),
		},
		"class": {
			Include("spaces"),
			Default(Pop(1), Push("class-body"), Push("bracket-open"), Push("extends"), Push("type-param-constraint"), Push("type-name")),
		},
		"extends": {
			Include("spaces"),
			{`(?:extends|implements)\b`, KeywordDeclaration, Push("type")},
			{`,`, Punctuation, nil},
			Default(Pop(1)),
		},
		"bracket-open": {
			Include("spaces"),
			{`\{`, Punctuation, Pop(1)},
		},
		"bracket-close": {
			Include("spaces"),
			{`\}`, Punctuation, Pop(1)},
		},
		"class-body": {
			Include("spaces"),
			Include("meta"),
			{`\}`, Punctuation, Pop(1)},
			{`(?:static|public|private|override|dynamic|inline|macro)\b`, KeywordDeclaration, nil},
			Default(Push("class-member")),
		},
		"class-member": {
			Include("spaces"),
			{`(var)\b`, KeywordDeclaration, Push("#pop", "optional-semicolon", "var")},
			{`(function)\b`, KeywordDeclaration, Push("#pop", "optional-semicolon", "class-method")},
		},
		"function-local": {
			Include("spaces"),
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, NameFunction, Push("#pop", "optional-expr", "flag", "function-param", "parenthesis-open", "type-param-constraint")},
			Default(Pop(1), Push("optional-expr"), Push("flag"), Push("function-param"), Push("parenthesis-open"), Push("type-param-constraint")),
		},
		"optional-expr": {
			Include("spaces"),
			Include("expr"),
			Default(Pop(1)),
		},
		"class-method": {
			Include("spaces"),
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, NameFunction, Push("#pop", "optional-expr", "flag", "function-param", "parenthesis-open", "type-param-constraint")},
		},
		"function-param": {
			Include("spaces"),
			{`\)`, Punctuation, Pop(1)},
			{`\?`, Punctuation, nil},
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Push("#pop", "function-param-sep", "assign", "flag")},
		},
		"function-param-sep": {
			Include("spaces"),
			{`\)`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "function-param")},
		},
		"prop-get-set": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "parenthesis-close", "prop-get-set-opt", "comma", "prop-get-set-opt")},
			Default(Pop(1)),
		},
		"prop-get-set-opt": {
			Include("spaces"),
			{`(?:default|null|never|dynamic|get|set)\b`, Keyword, Pop(1)},
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Text, Pop(1)},
		},
		"expr-statement": {
			Include("spaces"),
			Default(Pop(1), Push("optional-semicolon"), Push("expr")),
		},
		"expr": {
			Include("spaces"),
			{`@`, NameDecorator, Push("#pop", "optional-expr", "meta-body", "meta-ident", "meta-colon")},
			{`(?:\+\+|\-\-|~(?!/)|!|\-)`, Operator, nil},
			{`\(`, Punctuation, Push("#pop", "expr-chain", "parenthesis")},
			{`(?:static|public|private|override|dynamic|inline)\b`, KeywordDeclaration, nil},
			{`(?:function)\b`, KeywordDeclaration, Push("#pop", "expr-chain", "function-local")},
			{`\{`, Punctuation, Push("#pop", "expr-chain", "bracket")},
			{`(?:true|false|null)\b`, KeywordConstant, Push("#pop", "expr-chain")},
			{`(?:this)\b`, Keyword, Push("#pop", "expr-chain")},
			{`(?:cast)\b`, Keyword, Push("#pop", "expr-chain", "cast")},
			{`(?:try)\b`, Keyword, Push("#pop", "catch", "expr")},
			{`(?:var)\b`, KeywordDeclaration, Push("#pop", "var")},
			{`(?:new)\b`, Keyword, Push("#pop", "expr-chain", "new")},
			{`(?:switch)\b`, Keyword, Push("#pop", "switch")},
			{`(?:if)\b`, Keyword, Push("#pop", "if")},
			{`(?:do)\b`, Keyword, Push("#pop", "do")},
			{`(?:while)\b`, Keyword, Push("#pop", "while")},
			{`(?:for)\b`, Keyword, Push("#pop", "for")},
			{`(?:untyped|throw)\b`, Keyword, nil},
			{`(?:return)\b`, Keyword, Push("#pop", "optional-expr")},
			{`(?:macro)\b`, Keyword, Push("#pop", "macro")},
			{`(?:continue|break)\b`, Keyword, Pop(1)},
			{`(?:\$\s*[a-z]\b|\$(?!(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)))`, Name, Push("#pop", "dollar")},
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Push("#pop", "expr-chain")},
			{`\.[0-9]+`, LiteralNumberFloat, Push("#pop", "expr-chain")},
			{`[0-9]+[eE][+\-]?[0-9]+`, LiteralNumberFloat, Push("#pop", "expr-chain")},
			{`[0-9]+\.[0-9]*[eE][+\-]?[0-9]+`, LiteralNumberFloat, Push("#pop", "expr-chain")},
			{`[0-9]+\.[0-9]+`, LiteralNumberFloat, Push("#pop", "expr-chain")},
			{`[0-9]+\.(?!(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)|\.\.)`, LiteralNumberFloat, Push("#pop", "expr-chain")},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, Push("#pop", "expr-chain")},
			{`[0-9]+`, LiteralNumberInteger, Push("#pop", "expr-chain")},
			{`'`, LiteralStringSingle, Push("#pop", "expr-chain", "string-single-interpol")},
			{`"`, LiteralStringDouble, Push("#pop", "expr-chain", "string-double")},
			{`~/(\\\\|\\/|[^/\n])*/[gimsu]*`, LiteralStringRegex, Push("#pop", "expr-chain")},
			{`\[`, Punctuation, Push("#pop", "expr-chain", "array-decl")},
		},
		"expr-chain": {
			Include("spaces"),
			{`(?:\+\+|\-\-)`, Operator, nil},
			{`(?:%=|&=|\|=|\^=|\+=|\-=|\*=|/=|<<=|>\s*>\s*=|>\s*>\s*>\s*=|==|!=|<=|>\s*=|&&|\|\||<<|>>>|>\s*>|\.\.\.|<|>|%|&|\||\^|\+|\*|/|\-|=>|=)`, Operator, Push("#pop", "expr")},
			{`(?:in)\b`, Keyword, Push("#pop", "expr")},
			{`\?`, Operator, Push("#pop", "expr", "ternary", "expr")},
			{`(\.)((?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+))`, ByGroups(Punctuation, Name), nil},
			{`\[`, Punctuation, Push("array-access")},
			{`\(`, Punctuation, Push("call")},
			Default(Pop(1)),
		},
		"macro": {
			Include("spaces"),
			Include("meta"),
			{`:`, Punctuation, Push("#pop", "type")},
			{`(?:extern|private)\b`, KeywordDeclaration, nil},
			{`(?:abstract)\b`, KeywordDeclaration, Push("#pop", "optional-semicolon", "abstract")},
			{`(?:class|interface)\b`, KeywordDeclaration, Push("#pop", "optional-semicolon", "macro-class")},
			{`(?:enum)\b`, KeywordDeclaration, Push("#pop", "optional-semicolon", "enum")},
			{`(?:typedef)\b`, KeywordDeclaration, Push("#pop", "optional-semicolon", "typedef")},
			Default(Pop(1), Push("expr")),
		},
		"macro-class": {
			{`\{`, Punctuation, Push("#pop", "class-body")},
			Include("class"),
		},
		"cast": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "parenthesis-close", "cast-type", "expr")},
			Default(Pop(1), Push("expr")),
		},
		"cast-type": {
			Include("spaces"),
			{`,`, Punctuation, Push("#pop", "type")},
			Default(Pop(1)),
		},
		"catch": {
			Include("spaces"),
			{`(?:catch)\b`, Keyword, Push("expr", "function-param", "parenthesis-open")},
			Default(Pop(1)),
		},
		"do": {
			Include("spaces"),
			Default(Pop(1), Push("do-while"), Push("expr")),
		},
		"do-while": {
			Include("spaces"),
			{`(?:while)\b`, Keyword, Push("#pop", "parenthesis", "parenthesis-open")},
		},
		"while": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "expr", "parenthesis")},
		},
		"for": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "expr", "parenthesis")},
		},
		"if": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "else", "optional-semicolon", "expr", "parenthesis")},
		},
		"else": {
			Include("spaces"),
			{`(?:else)\b`, Keyword, Push("#pop", "expr")},
			Default(Pop(1)),
		},
		"switch": {
			Include("spaces"),
			Default(Pop(1), Push("switch-body"), Push("bracket-open"), Push("expr")),
		},
		"switch-body": {
			Include("spaces"),
			{`(?:case|default)\b`, Keyword, Push("case-block", "case")},
			{`\}`, Punctuation, Pop(1)},
		},
		"case": {
			Include("spaces"),
			{`:`, Punctuation, Pop(1)},
			Default(Pop(1), Push("case-sep"), Push("case-guard"), Push("expr")),
		},
		"case-sep": {
			Include("spaces"),
			{`:`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "case")},
		},
		"case-guard": {
			Include("spaces"),
			{`(?:if)\b`, Keyword, Push("#pop", "parenthesis", "parenthesis-open")},
			Default(Pop(1)),
		},
		"case-block": {
			Include("spaces"),
			{`(?!(?:case|default)\b|\})`, Keyword, Push("expr-statement")},
			Default(Pop(1)),
		},
		"new": {
			Include("spaces"),
			Default(Pop(1), Push("call"), Push("parenthesis-open"), Push("type")),
		},
		"array-decl": {
			Include("spaces"),
			{`\]`, Punctuation, Pop(1)},
			Default(Pop(1), Push("array-decl-sep"), Push("expr")),
		},
		"array-decl-sep": {
			Include("spaces"),
			{`\]`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "array-decl")},
		},
		"array-access": {
			Include("spaces"),
			Default(Pop(1), Push("array-access-close"), Push("expr")),
		},
		"array-access-close": {
			Include("spaces"),
			{`\]`, Punctuation, Pop(1)},
		},
		"comma": {
			Include("spaces"),
			{`,`, Punctuation, Pop(1)},
		},
		"colon": {
			Include("spaces"),
			{`:`, Punctuation, Pop(1)},
		},
		"semicolon": {
			Include("spaces"),
			{`;`, Punctuation, Pop(1)},
		},
		"optional-semicolon": {
			Include("spaces"),
			{`;`, Punctuation, Pop(1)},
			Default(Pop(1)),
		},
		"ident": {
			Include("spaces"),
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Pop(1)},
		},
		"dollar": {
			Include("spaces"),
			{`\{`, Punctuation, Push("#pop", "expr-chain", "bracket-close", "expr")},
			Default(Pop(1), Push("expr-chain")),
		},
		"type-name": {
			Include("spaces"),
			{`_*[A-Z]\w*`, Name, Pop(1)},
		},
		"type-full-name": {
			Include("spaces"),
			{`\.`, Punctuation, Push("ident")},
			Default(Pop(1)),
		},
		"type": {
			Include("spaces"),
			{`\?`, Punctuation, nil},
			{`(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Push("#pop", "type-check", "type-full-name")},
			{`\{`, Punctuation, Push("#pop", "type-check", "type-struct")},
			{`\(`, Punctuation, Push("#pop", "type-check", "type-parenthesis")},
		},
		"type-parenthesis": {
			Include("spaces"),
			Default(Pop(1), Push("parenthesis-close"), Push("type")),
		},
		"type-check": {
			Include("spaces"),
			{`->`, Punctuation, Push("#pop", "type")},
			{`<(?!=)`, Punctuation, Push("type-param")},
			Default(Pop(1)),
		},
		"type-struct": {
			Include("spaces"),
			{`\}`, Punctuation, Pop(1)},
			{`\?`, Punctuation, nil},
			{`>`, Punctuation, Push("comma", "type")},
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Push("#pop", "type-struct-sep", "type", "colon")},
			Include("class-body"),
		},
		"type-struct-sep": {
			Include("spaces"),
			{`\}`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "type-struct")},
		},
		"type-param-type": {
			{`\.[0-9]+`, LiteralNumberFloat, Pop(1)},
			{`[0-9]+[eE][+\-]?[0-9]+`, LiteralNumberFloat, Pop(1)},
			{`[0-9]+\.[0-9]*[eE][+\-]?[0-9]+`, LiteralNumberFloat, Pop(1)},
			{`[0-9]+\.[0-9]+`, LiteralNumberFloat, Pop(1)},
			{`[0-9]+\.(?!(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)|\.\.)`, LiteralNumberFloat, Pop(1)},
			{`0x[0-9a-fA-F]+`, LiteralNumberHex, Pop(1)},
			{`[0-9]+`, LiteralNumberInteger, Pop(1)},
			{`'`, LiteralStringSingle, Push("#pop", "string-single")},
			{`"`, LiteralStringDouble, Push("#pop", "string-double")},
			{`~/(\\\\|\\/|[^/\n])*/[gim]*`, LiteralStringRegex, Pop(1)},
			{`\[`, Operator, Push("#pop", "array-decl")},
			Include("type"),
		},
		"type-param": {
			Include("spaces"),
			Default(Pop(1), Push("type-param-sep"), Push("type-param-type")),
		},
		"type-param-sep": {
			Include("spaces"),
			{`>`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "type-param")},
		},
		"type-param-constraint": {
			Include("spaces"),
			{`<(?!=)`, Punctuation, Push("#pop", "type-param-constraint-sep", "type-param-constraint-flag", "type-name")},
			Default(Pop(1)),
		},
		"type-param-constraint-sep": {
			Include("spaces"),
			{`>`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "type-param-constraint-sep", "type-param-constraint-flag", "type-name")},
		},
		"type-param-constraint-flag": {
			Include("spaces"),
			{`:`, Punctuation, Push("#pop", "type-param-constraint-flag-type")},
			Default(Pop(1)),
		},
		"type-param-constraint-flag-type": {
			Include("spaces"),
			{`\(`, Punctuation, Push("#pop", "type-param-constraint-flag-type-sep", "type")},
			Default(Pop(1), Push("type")),
		},
		"type-param-constraint-flag-type-sep": {
			Include("spaces"),
			{`\)`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("type")},
		},
		"parenthesis": {
			Include("spaces"),
			Default(Pop(1), Push("parenthesis-close"), Push("flag"), Push("expr")),
		},
		"parenthesis-open": {
			Include("spaces"),
			{`\(`, Punctuation, Pop(1)},
		},
		"parenthesis-close": {
			Include("spaces"),
			{`\)`, Punctuation, Pop(1)},
		},
		"var": {
			Include("spaces"),
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Text, Push("#pop", "var-sep", "assign", "flag", "prop-get-set")},
		},
		"var-sep": {
			Include("spaces"),
			{`,`, Punctuation, Push("#pop", "var")},
			Default(Pop(1)),
		},
		"assign": {
			Include("spaces"),
			{`=`, Operator, Push("#pop", "expr")},
			Default(Pop(1)),
		},
		"flag": {
			Include("spaces"),
			{`:`, Punctuation, Push("#pop", "type")},
			Default(Pop(1)),
		},
		"ternary": {
			Include("spaces"),
			{`:`, Operator, Pop(1)},
		},
		"call": {
			Include("spaces"),
			{`\)`, Punctuation, Pop(1)},
			Default(Pop(1), Push("call-sep"), Push("expr")),
		},
		"call-sep": {
			Include("spaces"),
			{`\)`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "call")},
		},
		"bracket": {
			Include("spaces"),
			{`(?!(?:\$\s*[a-z]\b|\$(?!(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+))))(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Push("#pop", "bracket-check")},
			{`'`, LiteralStringSingle, Push("#pop", "bracket-check", "string-single")},
			{`"`, LiteralStringDouble, Push("#pop", "bracket-check", "string-double")},
			Default(Pop(1), Push("block")),
		},
		"bracket-check": {
			Include("spaces"),
			{`:`, Punctuation, Push("#pop", "object-sep", "expr")},
			Default(Pop(1), Push("block"), Push("optional-semicolon"), Push("expr-chain")),
		},
		"block": {
			Include("spaces"),
			{`\}`, Punctuation, Pop(1)},
			Default(Push("expr-statement")),
		},
		"object": {
			Include("spaces"),
			{`\}`, Punctuation, Pop(1)},
			Default(Pop(1), Push("object-sep"), Push("expr"), Push("colon"), Push("ident-or-string")),
		},
		"ident-or-string": {
			Include("spaces"),
			{`(?!(?:function|class|static|var|if|else|while|do|for|break|return|continue|extends|implements|import|switch|case|default|public|private|try|untyped|catch|new|this|throw|extern|enum|in|interface|cast|override|dynamic|typedef|package|inline|using|null|true|false|abstract)\b)(?:_*[a-z]\w*|_+[0-9]\w*|_*[A-Z]\w*|_+|\$\w+)`, Name, Pop(1)},
			{`'`, LiteralStringSingle, Push("#pop", "string-single")},
			{`"`, LiteralStringDouble, Push("#pop", "string-double")},
		},
		"object-sep": {
			Include("spaces"),
			{`\}`, Punctuation, Pop(1)},
			{`,`, Punctuation, Push("#pop", "object")},
		},
	}
}

func haxePreProcMutator(state *LexerState) error {
	stack, ok := state.Get("haxe-pre-proc").([][]string)
	if !ok {
		stack = [][]string{}
	}

	proc := state.Groups[2]
	switch proc {
	case "if":
		stack = append(stack, state.Stack)
	case "else", "elseif":
		if len(stack) > 0 {
			state.Stack = stack[len(stack)-1]
		}
	case "end":
		stack = stack[:len(stack)-1]
	}

	if proc == "if" || proc == "elseif" {
		state.Stack = append(state.Stack, "preproc-expr")
	}

	if proc == "error" {
		state.Stack = append(state.Stack, "preproc-error")
	}
	state.Set("haxe-pre-proc", stack)
	return nil
}
