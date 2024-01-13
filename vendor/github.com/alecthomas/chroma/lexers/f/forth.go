package f

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Forth lexer.
var Forth = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Forth",
		Aliases:         []string{"forth"},
		Filenames:       []string{"*.frt", "*.fth", "*.fs"},
		MimeTypes:       []string{"application/x-forth"},
		CaseInsensitive: true,
	},
	forthRules,
))

func forthRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`\\.*?\n`, CommentSingle, nil},
			{`\([\s].*?\)`, CommentSingle, nil},
			{`(:|variable|constant|value|buffer:)(\s+)`, ByGroups(KeywordNamespace, Text), Push("worddef")},
			{`([.sc]")(\s+?)`, ByGroups(LiteralString, Text), Push("stringdef")},
			{`(blk|block|buffer|evaluate|flush|load|save-buffers|update|empty-buffers|list|refill|scr|thru|\#s|\*\/mod|\+loop|\/mod|0<|0=|1\+|1-|2!|2\*|2\/|2@|2drop|2dup|2over|2swap|>body|>in|>number|>r|\?dup|abort|abort\"|abs|accept|align|aligned|allot|and|base|begin|bl|c!|c,|c@|cell\+|cells|char|char\+|chars|constant|count|cr|create|decimal|depth|do|does>|drop|dup|else|emit|environment\?|evaluate|execute|exit|fill|find|fm\/mod|here|hold|i|if|immediate|invert|j|key|leave|literal|loop|lshift|m\*|max|min|mod|move|negate|or|over|postpone|quit|r>|r@|recurse|repeat|rot|rshift|s\"|s>d|sign|sm\/rem|source|space|spaces|state|swap|then|type|u\.|u\<|um\*|um\/mod|unloop|until|variable|while|word|xor|\[char\]|\[\'\]|@|!|\#|<\#|\#>|:|;|\+|-|\*|\/|,|<|>|\|1\+|1-|\.|\.r|0<>|0>|2>r|2r>|2r@|:noname|\?do|again|c\"|case|compile,|endcase|endof|erase|false|hex|marker|nip|of|pad|parse|pick|refill|restore-input|roll|save-input|source-id|to|true|tuck|u\.r|u>|unused|value|within|\[compile\]|\#tib|convert|expect|query|span|tib|2constant|2literal|2variable|d\+|d-|d\.|d\.r|d0<|d0=|d2\*|d2\/|d<|d=|d>s|dabs|dmax|dmin|dnegate|m\*\/|m\+|2rot|du<|catch|throw|abort|abort\"|at-xy|key\?|page|ekey|ekey>char|ekey\?|emit\?|ms|time&date|BIN|CLOSE-FILE|CREATE-FILE|DELETE-FILE|FILE-POSITION|FILE-SIZE|INCLUDE-FILE|INCLUDED|OPEN-FILE|R\/O|R\/W|READ-FILE|READ-LINE|REPOSITION-FILE|RESIZE-FILE|S\"|SOURCE-ID|W/O|WRITE-FILE|WRITE-LINE|FILE-STATUS|FLUSH-FILE|REFILL|RENAME-FILE|>float|d>f|f!|f\*|f\+|f-|f\/|f0<|f0=|f<|f>d|f@|falign|faligned|fconstant|fdepth|fdrop|fdup|fliteral|float\+|floats|floor|fmax|fmin|fnegate|fover|frot|fround|fswap|fvariable|represent|df!|df@|dfalign|dfaligned|dfloat\+|dfloats|f\*\*|f\.|fabs|facos|facosh|falog|fasin|fasinh|fatan|fatan2|fatanh|fcos|fcosh|fe\.|fexp|fexpm1|fln|flnp1|flog|fs\.|fsin|fsincos|fsinh|fsqrt|ftan|ftanh|f~|precision|set-precision|sf!|sf@|sfalign|sfaligned|sfloat\+|sfloats|\(local\)|to|locals\||allocate|free|resize|definitions|find|forth-wordlist|get-current|get-order|search-wordlist|set-current|set-order|wordlist|also|forth|only|order|previous|-trailing|\/string|blank|cmove|cmove>|compare|search|sliteral|.s|dump|see|words|;code|ahead|assembler|bye|code|cs-pick|cs-roll|editor|state|\[else\]|\[if\]|\[then\]|forget|defer|defer@|defer!|action-of|begin-structure|field:|buffer:|parse-name|buffer:|traverse-wordlist|n>r|nr>|2value|fvalue|name>interpret|name>compile|name>string|cfield:|end-structure)\s`, Keyword, nil},
			{`(\$[0-9A-F]+)`, LiteralNumberHex, nil},
			{`(\#|%|&|\-|\+)?[0-9]+`, LiteralNumberInteger, nil},
			{`(\#|%|&|\-|\+)?[0-9.]+`, KeywordType, nil},
			{`(@i|!i|@e|!e|pause|noop|turnkey|sleep|itype|icompare|sp@|sp!|rp@|rp!|up@|up!|>a|a>|a@|a!|a@+|a@-|>b|b>|b@|b!|b@+|b@-|find-name|1ms|sp0|rp0|\(evaluate\)|int-trap|int!)\s`, NameConstant, nil},
			{`(do-recognizer|r:fail|recognizer:|get-recognizers|set-recognizers|r:float|r>comp|r>int|r>post|r:name|r:word|r:dnum|r:num|recognizer|forth-recognizer|rec:num|rec:float|rec:word)\s`, NameDecorator, nil},
			{`(Evalue|Rvalue|Uvalue|Edefer|Rdefer|Udefer)(\s+)`, ByGroups(KeywordNamespace, Text), Push("worddef")},
			{`[^\s]+(?=[\s])`, NameFunction, nil},
		},
		"worddef": {
			{`\S+`, NameClass, Pop(1)},
		},
		"stringdef": {
			{`[^"]+`, LiteralString, Pop(1)},
		},
	}
}
