package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// caddyfileCommon are the rules common to both of the lexer variants
func caddyfileCommonRules() Rules {
	return Rules{
		"site_block_common": {
			// Import keyword
			{`(import)(\s+)([^\s]+)`, ByGroups(Keyword, Text, NameVariableMagic), nil},
			// Matcher definition
			{`@[^\s]+(?=\s)`, NameDecorator, Push("matcher")},
			// Matcher token stub for docs
			{`\[\<matcher\>\]`, NameDecorator, Push("matcher")},
			// These cannot have matchers but may have things that look like
			// matchers in their arguments, so we just parse as a subdirective.
			{`try_files`, Keyword, Push("subdirective")},
			// These are special, they can nest more directives
			{`handle_errors|handle|route|handle_path|not`, Keyword, Push("nested_directive")},
			// Any other directive
			{`[^\s#]+`, Keyword, Push("directive")},
			Include("base"),
		},
		"matcher": {
			{`\{`, Punctuation, Push("block")},
			// Not can be one-liner
			{`not`, Keyword, Push("deep_not_matcher")},
			// Any other same-line matcher
			{`[^\s#]+`, Keyword, Push("arguments")},
			// Terminators
			{`\n`, Text, Pop(1)},
			{`\}`, Punctuation, Pop(1)},
			Include("base"),
		},
		"block": {
			{`\}`, Punctuation, Pop(2)},
			// Not can be one-liner
			{`not`, Keyword, Push("not_matcher")},
			// Any other subdirective
			{`[^\s#]+`, Keyword, Push("subdirective")},
			Include("base"),
		},
		"nested_block": {
			{`\}`, Punctuation, Pop(2)},
			// Matcher definition
			{`@[^\s]+(?=\s)`, NameDecorator, Push("matcher")},
			// Something that starts with literally < is probably a docs stub
			{`\<[^#]+\>`, Keyword, Push("nested_directive")},
			// Any other directive
			{`[^\s#]+`, Keyword, Push("nested_directive")},
			Include("base"),
		},
		"not_matcher": {
			{`\}`, Punctuation, Pop(2)},
			{`\{(?=\s)`, Punctuation, Push("block")},
			{`[^\s#]+`, Keyword, Push("arguments")},
			{`\s+`, Text, nil},
		},
		"deep_not_matcher": {
			{`\}`, Punctuation, Pop(2)},
			{`\{(?=\s)`, Punctuation, Push("block")},
			{`[^\s#]+`, Keyword, Push("deep_subdirective")},
			{`\s+`, Text, nil},
		},
		"directive": {
			{`\{(?=\s)`, Punctuation, Push("block")},
			Include("matcher_token"),
			Include("comments_pop_1"),
			{`\n`, Text, Pop(1)},
			Include("base"),
		},
		"nested_directive": {
			{`\{(?=\s)`, Punctuation, Push("nested_block")},
			Include("matcher_token"),
			Include("comments_pop_1"),
			{`\n`, Text, Pop(1)},
			Include("base"),
		},
		"subdirective": {
			{`\{(?=\s)`, Punctuation, Push("block")},
			Include("comments_pop_1"),
			{`\n`, Text, Pop(1)},
			Include("base"),
		},
		"arguments": {
			{`\{(?=\s)`, Punctuation, Push("block")},
			Include("comments_pop_2"),
			{`\\\n`, Text, nil}, // Skip escaped newlines
			{`\n`, Text, Pop(2)},
			Include("base"),
		},
		"deep_subdirective": {
			{`\{(?=\s)`, Punctuation, Push("block")},
			Include("comments_pop_3"),
			{`\n`, Text, Pop(3)},
			Include("base"),
		},
		"matcher_token": {
			{`@[^\s]+`, NameDecorator, Push("arguments")},         // Named matcher
			{`/[^\s]+`, NameDecorator, Push("arguments")},         // Path matcher
			{`\*`, NameDecorator, Push("arguments")},              // Wildcard path matcher
			{`\[\<matcher\>\]`, NameDecorator, Push("arguments")}, // Matcher token stub for docs
		},
		"comments": {
			{`^#.*\n`, CommentSingle, nil},   // Comment at start of line
			{`\s+#.*\n`, CommentSingle, nil}, // Comment preceded by whitespace
		},
		"comments_pop_1": {
			{`^#.*\n`, CommentSingle, Pop(1)},   // Comment at start of line
			{`\s+#.*\n`, CommentSingle, Pop(1)}, // Comment preceded by whitespace
		},
		"comments_pop_2": {
			{`^#.*\n`, CommentSingle, Pop(2)},   // Comment at start of line
			{`\s+#.*\n`, CommentSingle, Pop(2)}, // Comment preceded by whitespace
		},
		"comments_pop_3": {
			{`^#.*\n`, CommentSingle, Pop(3)},   // Comment at start of line
			{`\s+#.*\n`, CommentSingle, Pop(3)}, // Comment preceded by whitespace
		},
		"base": {
			Include("comments"),
			{`(on|off|first|last|before|after|internal|strip_prefix|strip_suffix|replace)\b`, NameConstant, nil},
			{`(https?://)?([a-z0-9.-]+)(:)([0-9]+)`, ByGroups(Name, Name, Punctuation, LiteralNumberInteger), nil},
			{`[a-z-]+/[a-z-+]+`, LiteralString, nil},
			{`[0-9]+[km]?\b`, LiteralNumberInteger, nil},
			{`\{[\w+.\$-]+\}`, LiteralStringEscape, nil}, // Placeholder
			{`\[(?=[^#{}$]+\])`, Punctuation, nil},
			{`\]|\|`, Punctuation, nil},
			{`[^\s#{}$\]]+`, LiteralString, nil},
			{`/[^\s#]*`, Name, nil},
			{`\s+`, Text, nil},
		},
	}
}

// Caddyfile lexer.
var Caddyfile = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Caddyfile",
		Aliases:   []string{"caddyfile", "caddy"},
		Filenames: []string{"Caddyfile*"},
		MimeTypes: []string{},
	},
	caddyfileRules,
))

func caddyfileRules() Rules {
	return Rules{
		"root": {
			Include("comments"),
			// Global options block
			{`^\s*(\{)\s*$`, ByGroups(Punctuation), Push("globals")},
			// Snippets
			{`(\([^\s#]+\))(\s*)(\{)`, ByGroups(NameVariableAnonymous, Text, Punctuation), Push("snippet")},
			// Site label
			{`[^#{(\s,]+`, GenericHeading, Push("label")},
			// Site label with placeholder
			{`\{[\w+.\$-]+\}`, LiteralStringEscape, Push("label")},
			{`\s+`, Text, nil},
		},
		"globals": {
			{`\}`, Punctuation, Pop(1)},
			{`[^\s#]+`, Keyword, Push("directive")},
			Include("base"),
		},
		"snippet": {
			{`\}`, Punctuation, Pop(1)},
			// Matcher definition
			{`@[^\s]+(?=\s)`, NameDecorator, Push("matcher")},
			// Any directive
			{`[^\s#]+`, Keyword, Push("directive")},
			Include("base"),
		},
		"label": {
			// Allow multiple labels, comma separated, newlines after
			// a comma means another label is coming
			{`,\s*\n?`, Text, nil},
			{` `, Text, nil},
			// Site label with placeholder
			{`\{[\w+.\$-]+\}`, LiteralStringEscape, nil},
			// Site label
			{`[^#{(\s,]+`, GenericHeading, nil},
			// Comment after non-block label (hack because comments end in \n)
			{`#.*\n`, CommentSingle, Push("site_block")},
			// Note: if \n, we'll never pop out of the site_block, it's valid
			{`\{(?=\s)|\n`, Punctuation, Push("site_block")},
		},
		"site_block": {
			{`\}`, Punctuation, Pop(2)},
			Include("site_block_common"),
		},
	}.Merge(caddyfileCommonRules())
}

// Caddyfile directive-only lexer.
var CaddyfileDirectives = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Caddyfile Directives",
		Aliases:   []string{"caddyfile-directives", "caddyfile-d", "caddy-d"},
		Filenames: []string{},
		MimeTypes: []string{},
	},
	caddyfileDirectivesRules,
))

func caddyfileDirectivesRules() Rules {
	return Rules{
		// Same as "site_block" in Caddyfile
		"root": {
			Include("site_block_common"),
		},
	}.Merge(caddyfileCommonRules())
}
