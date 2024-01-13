package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Sass lexer.
var Sass = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Sass",
		Aliases:         []string{"sass"},
		Filenames:       []string{"*.sass"},
		MimeTypes:       []string{"text/x-sass"},
		CaseInsensitive: true,
	},
	sassRules,
))

func sassRules() Rules {
	return Rules{
		// "root": {
		// },
		"root": {
			{`[ \t]*\n`, Text, nil},
			// { `[ \t]*`, ?? <function _indentation at 0x106932e18> ??, nil },
			// { `//[^\n]*`, ?? <function _starts_block.<locals>.callback at 0x106936048> ??, Push("root") },
			// { `/\*[^\n]*`, ?? <function _starts_block.<locals>.callback at 0x1069360d0> ??, Push("root") },
			{`@import`, Keyword, Push("import")},
			{`@for`, Keyword, Push("for")},
			{`@(debug|warn|if|while)`, Keyword, Push("value")},
			{`(@mixin)( [\w-]+)`, ByGroups(Keyword, NameFunction), Push("value")},
			{`(@include)( [\w-]+)`, ByGroups(Keyword, NameDecorator), Push("value")},
			{`@extend`, Keyword, Push("selector")},
			{`@[\w-]+`, Keyword, Push("selector")},
			{`=[\w-]+`, NameFunction, Push("value")},
			{`\+[\w-]+`, NameDecorator, Push("value")},
			{`([!$][\w-]\w*)([ \t]*(?:(?:\|\|)?=|:))`, ByGroups(NameVariable, Operator), Push("value")},
			{`:`, NameAttribute, Push("old-style-attr")},
			{`(?=.+?[=:]([^a-z]|$))`, NameAttribute, Push("new-style-attr")},
			Default(Push("selector")),
		},
		"single-comment": {
			{`.+`, CommentSingle, nil},
			{`\n`, Text, Push("root")},
		},
		"multi-comment": {
			{`.+`, CommentMultiline, nil},
			{`\n`, Text, Push("root")},
		},
		"import": {
			{`[ \t]+`, Text, nil},
			{`\S+`, LiteralString, nil},
			{`\n`, Text, Push("root")},
		},
		"old-style-attr": {
			{`[^\s:="\[]+`, NameAttribute, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`[ \t]*=`, Operator, Push("value")},
			Default(Push("value")),
		},
		"new-style-attr": {
			{`[^\s:="\[]+`, NameAttribute, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`[ \t]*[=:]`, Operator, Push("value")},
		},
		"inline-comment": {
			{`(\\#|#(?=[^\n{])|\*(?=[^\n/])|[^\n#*])+`, CommentMultiline, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`\*/`, Comment, Pop(1)},
		},
		"value": {
			{`[ \t]+`, Text, nil},
			{`[!$][\w-]+`, NameVariable, nil},
			{`url\(`, LiteralStringOther, Push("string-url")},
			{`[a-z_-][\w-]*(?=\()`, NameFunction, nil},
			{Words(``, `\b`, `align-content`, `align-items`, `align-self`, `alignment-baseline`, `all`, `animation`, `animation-delay`, `animation-direction`, `animation-duration`, `animation-fill-mode`, `animation-iteration-count`, `animation-name`, `animation-play-state`, `animation-timing-function`, `appearance`, `azimuth`, `backface-visibility`, `background`, `background-attachment`, `background-blend-mode`, `background-clip`, `background-color`, `background-image`, `background-origin`, `background-position`, `background-repeat`, `background-size`, `baseline-shift`, `bookmark-label`, `bookmark-level`, `bookmark-state`, `border`, `border-bottom`, `border-bottom-color`, `border-bottom-left-radius`, `border-bottom-right-radius`, `border-bottom-style`, `border-bottom-width`, `border-boundary`, `border-collapse`, `border-color`, `border-image`, `border-image-outset`, `border-image-repeat`, `border-image-slice`, `border-image-source`, `border-image-width`, `border-left`, `border-left-color`, `border-left-style`, `border-left-width`, `border-radius`, `border-right`, `border-right-color`, `border-right-style`, `border-right-width`, `border-spacing`, `border-style`, `border-top`, `border-top-color`, `border-top-left-radius`, `border-top-right-radius`, `border-top-style`, `border-top-width`, `border-width`, `bottom`, `box-decoration-break`, `box-shadow`, `box-sizing`, `box-snap`, `box-suppress`, `break-after`, `break-before`, `break-inside`, `caption-side`, `caret`, `caret-animation`, `caret-color`, `caret-shape`, `chains`, `clear`, `clip`, `clip-path`, `clip-rule`, `color`, `color-interpolation-filters`, `column-count`, `column-fill`, `column-gap`, `column-rule`, `column-rule-color`, `column-rule-style`, `column-rule-width`, `column-span`, `column-width`, `columns`, `content`, `counter-increment`, `counter-reset`, `counter-set`, `crop`, `cue`, `cue-after`, `cue-before`, `cursor`, `direction`, `display`, `dominant-baseline`, `elevation`, `empty-cells`, `filter`, `flex`, `flex-basis`, `flex-direction`, `flex-flow`, `flex-grow`, `flex-shrink`, `flex-wrap`, `float`, `float-defer`, `float-offset`, `float-reference`, `flood-color`, `flood-opacity`, `flow`, `flow-from`, `flow-into`, `font`, `font-family`, `font-feature-settings`, `font-kerning`, `font-language-override`, `font-size`, `font-size-adjust`, `font-stretch`, `font-style`, `font-synthesis`, `font-variant`, `font-variant-alternates`, `font-variant-caps`, `font-variant-east-asian`, `font-variant-ligatures`, `font-variant-numeric`, `font-variant-position`, `font-weight`, `footnote-display`, `footnote-policy`, `glyph-orientation-vertical`, `grid`, `grid-area`, `grid-auto-columns`, `grid-auto-flow`, `grid-auto-rows`, `grid-column`, `grid-column-end`, `grid-column-gap`, `grid-column-start`, `grid-gap`, `grid-row`, `grid-row-end`, `grid-row-gap`, `grid-row-start`, `grid-template`, `grid-template-areas`, `grid-template-columns`, `grid-template-rows`, `hanging-punctuation`, `height`, `hyphenate-character`, `hyphenate-limit-chars`, `hyphenate-limit-last`, `hyphenate-limit-lines`, `hyphenate-limit-zone`, `hyphens`, `image-orientation`, `image-resolution`, `initial-letter`, `initial-letter-align`, `initial-letter-wrap`, `isolation`, `justify-content`, `justify-items`, `justify-self`, `left`, `letter-spacing`, `lighting-color`, `line-break`, `line-grid`, `line-height`, `line-snap`, `list-style`, `list-style-image`, `list-style-position`, `list-style-type`, `margin`, `margin-bottom`, `margin-left`, `margin-right`, `margin-top`, `marker-side`, `marquee-direction`, `marquee-loop`, `marquee-speed`, `marquee-style`, `mask`, `mask-border`, `mask-border-mode`, `mask-border-outset`, `mask-border-repeat`, `mask-border-slice`, `mask-border-source`, `mask-border-width`, `mask-clip`, `mask-composite`, `mask-image`, `mask-mode`, `mask-origin`, `mask-position`, `mask-repeat`, `mask-size`, `mask-type`, `max-height`, `max-lines`, `max-width`, `min-height`, `min-width`, `mix-blend-mode`, `motion`, `motion-offset`, `motion-path`, `motion-rotation`, `move-to`, `nav-down`, `nav-left`, `nav-right`, `nav-up`, `object-fit`, `object-position`, `offset-after`, `offset-before`, `offset-end`, `offset-start`, `opacity`, `order`, `orphans`, `outline`, `outline-color`, `outline-offset`, `outline-style`, `outline-width`, `overflow`, `overflow-style`, `overflow-wrap`, `overflow-x`, `overflow-y`, `padding`, `padding-bottom`, `padding-left`, `padding-right`, `padding-top`, `page`, `page-break-after`, `page-break-before`, `page-break-inside`, `page-policy`, `pause`, `pause-after`, `pause-before`, `perspective`, `perspective-origin`, `pitch`, `pitch-range`, `play-during`, `polar-angle`, `polar-distance`, `position`, `presentation-level`, `quotes`, `region-fragment`, `resize`, `rest`, `rest-after`, `rest-before`, `richness`, `right`, `rotation`, `rotation-point`, `ruby-align`, `ruby-merge`, `ruby-position`, `running`, `scroll-snap-coordinate`, `scroll-snap-destination`, `scroll-snap-points-x`, `scroll-snap-points-y`, `scroll-snap-type`, `shape-image-threshold`, `shape-inside`, `shape-margin`, `shape-outside`, `size`, `speak`, `speak-as`, `speak-header`, `speak-numeral`, `speak-punctuation`, `speech-rate`, `stress`, `string-set`, `tab-size`, `table-layout`, `text-align`, `text-align-last`, `text-combine-upright`, `text-decoration`, `text-decoration-color`, `text-decoration-line`, `text-decoration-skip`, `text-decoration-style`, `text-emphasis`, `text-emphasis-color`, `text-emphasis-position`, `text-emphasis-style`, `text-indent`, `text-justify`, `text-orientation`, `text-overflow`, `text-shadow`, `text-space-collapse`, `text-space-trim`, `text-spacing`, `text-transform`, `text-underline-position`, `text-wrap`, `top`, `transform`, `transform-origin`, `transform-style`, `transition`, `transition-delay`, `transition-duration`, `transition-property`, `transition-timing-function`, `unicode-bidi`, `user-select`, `vertical-align`, `visibility`, `voice-balance`, `voice-duration`, `voice-family`, `voice-pitch`, `voice-range`, `voice-rate`, `voice-stress`, `voice-volume`, `volume`, `white-space`, `widows`, `width`, `will-change`, `word-break`, `word-spacing`, `word-wrap`, `wrap-after`, `wrap-before`, `wrap-flow`, `wrap-inside`, `wrap-through`, `writing-mode`, `z-index`, `above`, `absolute`, `always`, `armenian`, `aural`, `auto`, `avoid`, `baseline`, `behind`, `below`, `bidi-override`, `blink`, `block`, `bold`, `bolder`, `both`, `capitalize`, `center-left`, `center-right`, `center`, `circle`, `cjk-ideographic`, `close-quote`, `collapse`, `condensed`, `continuous`, `crop`, `crosshair`, `cross`, `cursive`, `dashed`, `decimal-leading-zero`, `decimal`, `default`, `digits`, `disc`, `dotted`, `double`, `e-resize`, `embed`, `extra-condensed`, `extra-expanded`, `expanded`, `fantasy`, `far-left`, `far-right`, `faster`, `fast`, `fixed`, `georgian`, `groove`, `hebrew`, `help`, `hidden`, `hide`, `higher`, `high`, `hiragana-iroha`, `hiragana`, `icon`, `inherit`, `inline-table`, `inline`, `inset`, `inside`, `invert`, `italic`, `justify`, `katakana-iroha`, `katakana`, `landscape`, `larger`, `large`, `left-side`, `leftwards`, `level`, `lighter`, `line-through`, `list-item`, `loud`, `lower-alpha`, `lower-greek`, `lower-roman`, `lowercase`, `ltr`, `lower`, `low`, `medium`, `message-box`, `middle`, `mix`, `monospace`, `n-resize`, `narrower`, `ne-resize`, `no-close-quote`, `no-open-quote`, `no-repeat`, `none`, `normal`, `nowrap`, `nw-resize`, `oblique`, `once`, `open-quote`, `outset`, `outside`, `overline`, `pointer`, `portrait`, `px`, `relative`, `repeat-x`, `repeat-y`, `repeat`, `rgb`, `ridge`, `right-side`, `rightwards`, `s-resize`, `sans-serif`, `scroll`, `se-resize`, `semi-condensed`, `semi-expanded`, `separate`, `serif`, `show`, `silent`, `slow`, `slower`, `small-caps`, `small-caption`, `smaller`, `soft`, `solid`, `spell-out`, `square`, `static`, `status-bar`, `super`, `sw-resize`, `table-caption`, `table-cell`, `table-column`, `table-column-group`, `table-footer-group`, `table-header-group`, `table-row`, `table-row-group`, `text`, `text-bottom`, `text-top`, `thick`, `thin`, `transparent`, `ultra-condensed`, `ultra-expanded`, `underline`, `upper-alpha`, `upper-latin`, `upper-roman`, `uppercase`, `url`, `visible`, `w-resize`, `wait`, `wider`, `x-fast`, `x-high`, `x-large`, `x-loud`, `x-low`, `x-small`, `x-soft`, `xx-large`, `xx-small`, `yes`), NameConstant, nil},
			{Words(``, `\b`, `aliceblue`, `antiquewhite`, `aqua`, `aquamarine`, `azure`, `beige`, `bisque`, `black`, `blanchedalmond`, `blue`, `blueviolet`, `brown`, `burlywood`, `cadetblue`, `chartreuse`, `chocolate`, `coral`, `cornflowerblue`, `cornsilk`, `crimson`, `cyan`, `darkblue`, `darkcyan`, `darkgoldenrod`, `darkgray`, `darkgreen`, `darkgrey`, `darkkhaki`, `darkmagenta`, `darkolivegreen`, `darkorange`, `darkorchid`, `darkred`, `darksalmon`, `darkseagreen`, `darkslateblue`, `darkslategray`, `darkslategrey`, `darkturquoise`, `darkviolet`, `deeppink`, `deepskyblue`, `dimgray`, `dimgrey`, `dodgerblue`, `firebrick`, `floralwhite`, `forestgreen`, `fuchsia`, `gainsboro`, `ghostwhite`, `gold`, `goldenrod`, `gray`, `green`, `greenyellow`, `grey`, `honeydew`, `hotpink`, `indianred`, `indigo`, `ivory`, `khaki`, `lavender`, `lavenderblush`, `lawngreen`, `lemonchiffon`, `lightblue`, `lightcoral`, `lightcyan`, `lightgoldenrodyellow`, `lightgray`, `lightgreen`, `lightgrey`, `lightpink`, `lightsalmon`, `lightseagreen`, `lightskyblue`, `lightslategray`, `lightslategrey`, `lightsteelblue`, `lightyellow`, `lime`, `limegreen`, `linen`, `magenta`, `maroon`, `mediumaquamarine`, `mediumblue`, `mediumorchid`, `mediumpurple`, `mediumseagreen`, `mediumslateblue`, `mediumspringgreen`, `mediumturquoise`, `mediumvioletred`, `midnightblue`, `mintcream`, `mistyrose`, `moccasin`, `navajowhite`, `navy`, `oldlace`, `olive`, `olivedrab`, `orange`, `orangered`, `orchid`, `palegoldenrod`, `palegreen`, `paleturquoise`, `palevioletred`, `papayawhip`, `peachpuff`, `peru`, `pink`, `plum`, `powderblue`, `purple`, `rebeccapurple`, `red`, `rosybrown`, `royalblue`, `saddlebrown`, `salmon`, `sandybrown`, `seagreen`, `seashell`, `sienna`, `silver`, `skyblue`, `slateblue`, `slategray`, `slategrey`, `snow`, `springgreen`, `steelblue`, `tan`, `teal`, `thistle`, `tomato`, `turquoise`, `violet`, `wheat`, `white`, `whitesmoke`, `yellow`, `yellowgreen`, `transparent`), NameEntity, nil},
			{Words(``, `\b`, `black`, `silver`, `gray`, `white`, `maroon`, `red`, `purple`, `fuchsia`, `green`, `lime`, `olive`, `yellow`, `navy`, `blue`, `teal`, `aqua`), NameBuiltin, nil},
			{`\!(important|default)`, NameException, nil},
			{`(true|false)`, NamePseudo, nil},
			{`(and|or|not)`, OperatorWord, nil},
			{`/\*`, CommentMultiline, Push("inline-comment")},
			{`//[^\n]*`, CommentSingle, nil},
			{`\#[a-z0-9]{1,6}`, LiteralNumberHex, nil},
			{`(-?\d+)(\%|[a-z]+)?`, ByGroups(LiteralNumberInteger, KeywordType), nil},
			{`(-?\d*\.\d+)(\%|[a-z]+)?`, ByGroups(LiteralNumberFloat, KeywordType), nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`[~^*!&%<>|+=@:,./?-]+`, Operator, nil},
			{`[\[\]()]+`, Punctuation, nil},
			{`"`, LiteralStringDouble, Push("string-double")},
			{`'`, LiteralStringSingle, Push("string-single")},
			{`[a-z_-][\w-]*`, Name, nil},
			{`\n`, Text, Push("root")},
		},
		"interpolation": {
			{`\}`, LiteralStringInterpol, Pop(1)},
			Include("value"),
		},
		"selector": {
			{`[ \t]+`, Text, nil},
			{`\:`, NameDecorator, Push("pseudo-class")},
			{`\.`, NameClass, Push("class")},
			{`\#`, NameNamespace, Push("id")},
			{`[\w-]+`, NameTag, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`&`, Keyword, nil},
			{`[~^*!&\[\]()<>|+=@:;,./?-]`, Operator, nil},
			{`"`, LiteralStringDouble, Push("string-double")},
			{`'`, LiteralStringSingle, Push("string-single")},
			{`\n`, Text, Push("root")},
		},
		"string-double": {
			{`(\\.|#(?=[^\n{])|[^\n"#])+`, LiteralStringDouble, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`"`, LiteralStringDouble, Pop(1)},
		},
		"string-single": {
			{`(\\.|#(?=[^\n{])|[^\n'#])+`, LiteralStringSingle, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`'`, LiteralStringSingle, Pop(1)},
		},
		"string-url": {
			{`(\\#|#(?=[^\n{])|[^\n#)])+`, LiteralStringOther, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`\)`, LiteralStringOther, Pop(1)},
		},
		"pseudo-class": {
			{`[\w-]+`, NameDecorator, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			Default(Pop(1)),
		},
		"class": {
			{`[\w-]+`, NameClass, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			Default(Pop(1)),
		},
		"id": {
			{`[\w-]+`, NameNamespace, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			Default(Pop(1)),
		},
		"for": {
			{`(from|to|through)`, OperatorWord, nil},
			Include("value"),
		},
	}
}
