package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Powershell lexer.
var Powershell = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "PowerShell",
		Aliases:         []string{"powershell", "posh", "ps1", "psm1", "psd1"},
		Filenames:       []string{"*.ps1", "*.psm1", "*.psd1"},
		MimeTypes:       []string{"text/x-powershell"},
		DotAll:          true,
		CaseInsensitive: true,
	},
	powershellRules,
))

func powershellRules() Rules {
	return Rules{
		"root": {
			{`\(`, Punctuation, Push("child")},
			{`\s+`, Text, nil},
			{`^(\s*#[#\s]*)(\.(?:component|description|example|externalhelp|forwardhelpcategory|forwardhelptargetname|functionality|inputs|link|notes|outputs|parameter|remotehelprunspace|role|synopsis))([^\n]*$)`, ByGroups(Comment, LiteralStringDoc, Comment), nil},
			{`#[^\n]*?$`, Comment, nil},
			{`(&lt;|<)#`, CommentMultiline, Push("multline")},
			{`(?i)([A-Z]:)`, Name, nil},
			{`@"\n`, LiteralStringHeredoc, Push("heredoc-double")},
			{`@'\n.*?\n'@`, LiteralStringHeredoc, nil},
			{"`[\\'\"$@-]", Punctuation, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`'([^']|'')*'`, LiteralStringSingle, nil},
			{`(\$|@@|@)((global|script|private|env):)?\w+`, NameVariable, nil},
			{`[a-z]\w*-[a-z]\w*\b`, NameBuiltin, nil},
			{`(while|validateset|validaterange|validatepattern|validatelength|validatecount|until|trap|switch|return|ref|process|param|parameter|in|if|global:|function|foreach|for|finally|filter|end|elseif|else|dynamicparam|do|default|continue|cmdletbinding|break|begin|alias|\?|%|#script|#private|#local|#global|mandatory|parametersetname|position|valuefrompipeline|valuefrompipelinebypropertyname|valuefromremainingarguments|helpmessage|try|catch|throw)\b`, Keyword, nil},
			{`-(and|as|band|bnot|bor|bxor|casesensitive|ccontains|ceq|cge|cgt|cle|clike|clt|cmatch|cne|cnotcontains|cnotlike|cnotmatch|contains|creplace|eq|exact|f|file|ge|gt|icontains|ieq|ige|igt|ile|ilike|ilt|imatch|ine|inotcontains|inotlike|inotmatch|ireplace|is|isnot|le|like|lt|match|ne|not|notcontains|notlike|notmatch|or|regex|replace|wildcard)\b`, Operator, nil},
			{`(ac|asnp|cat|cd|cfs|chdir|clc|clear|clhy|cli|clp|cls|clv|cnsn|compare|copy|cp|cpi|cpp|curl|cvpa|dbp|del|diff|dir|dnsn|ebp|echo|epal|epcsv|epsn|erase|etsn|exsn|fc|fhx|fl|foreach|ft|fw|gal|gbp|gc|gci|gcm|gcs|gdr|ghy|gi|gjb|gl|gm|gmo|gp|gps|gpv|group|gsn|gsnp|gsv|gu|gv|gwmi|h|history|icm|iex|ihy|ii|ipal|ipcsv|ipmo|ipsn|irm|ise|iwmi|iwr|kill|lp|ls|man|md|measure|mi|mount|move|mp|mv|nal|ndr|ni|nmo|npssc|nsn|nv|ogv|oh|popd|ps|pushd|pwd|r|rbp|rcjb|rcsn|rd|rdr|ren|ri|rjb|rm|rmdir|rmo|rni|rnp|rp|rsn|rsnp|rujb|rv|rvpa|rwmi|sajb|sal|saps|sasv|sbp|sc|select|set|shcm|si|sl|sleep|sls|sort|sp|spjb|spps|spsv|start|sujb|sv|swmi|tee|trcm|type|wget|where|wjb|write)\s`, NameBuiltin, nil},
			{"\\[[a-z_\\[][\\w. `,\\[\\]]*\\]", NameConstant, nil},
			{`-[a-z_]\w*`, Name, nil},
			{`\w+`, Name, nil},
			{"[.,;@{}\\[\\]$()=+*/\\\\&%!~?^`|<>-]|::", Punctuation, nil},
		},
		"child": {
			{`\)`, Punctuation, Pop(1)},
			Include("root"),
		},
		"multline": {
			{`[^#&.]+`, CommentMultiline, nil},
			{`#(>|&gt;)`, CommentMultiline, Pop(1)},
			{`\.(component|description|example|externalhelp|forwardhelpcategory|forwardhelptargetname|functionality|inputs|link|notes|outputs|parameter|remotehelprunspace|role|synopsis)`, LiteralStringDoc, nil},
			{`[#&.]`, CommentMultiline, nil},
		},
		"string": {
			{"`[0abfnrtv'\\\"$`]", LiteralStringEscape, nil},
			{"[^$`\"]+", LiteralStringDouble, nil},
			{`\$\(`, Punctuation, Push("child")},
			{`""`, LiteralStringDouble, nil},
			{"[`$]", LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
		"heredoc-double": {
			{`\n"@`, LiteralStringHeredoc, Pop(1)},
			{`\$\(`, Punctuation, Push("child")},
			{`[^@\n]+"]`, LiteralStringHeredoc, nil},
			{`.`, LiteralStringHeredoc, nil},
		},
	}
}
