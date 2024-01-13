package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// CSS lexer.
var CSS = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "CSS",
		Aliases:   []string{"css"},
		Filenames: []string{"*.css"},
		MimeTypes: []string{"text/css"},
	},
	cssRules,
))

func cssRules() Rules {
	return Rules{
		"root": {
			Include("basics"),
		},
		"basics": {
			{`\s+`, Text, nil},
			{`/\*(?:.|\n)*?\*/`, Comment, nil},
			{`\{`, Punctuation, Push("content")},
			{`(\:{1,2})([\w-]+)`, ByGroups(Punctuation, NameDecorator), nil},
			{`(\.)([\w-]+)`, ByGroups(Punctuation, NameClass), nil},
			{`(\#)([\w-]+)`, ByGroups(Punctuation, NameNamespace), nil},
			{`(@)([\w-]+)`, ByGroups(Punctuation, Keyword), Push("atrule")},
			{`[\w-]+`, NameTag, nil},
			{`[~^*!%&$\[\]()<>|+=@:;,./?-]`, Operator, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
		},
		"atrule": {
			{`\{`, Punctuation, Push("atcontent")},
			{`;`, Punctuation, Pop(1)},
			Include("basics"),
		},
		"atcontent": {
			Include("basics"),
			{`\}`, Punctuation, Pop(2)},
		},
		"atparenthesis": {
			Include("common-values"),
			{`/\*(?:.|\n)*?\*/`, Comment, nil},
			Include("numeric-values"),
			{`[*+/-]`, Operator, nil},
			{`[,]`, Punctuation, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`[a-zA-Z_-]\w*`, Name, nil},
			{`\(`, Punctuation, Push("atparenthesis")},
			{`\)`, Punctuation, Pop(1)},
		},
		"content": {
			{`\s+`, Text, nil},
			{`\}`, Punctuation, Pop(1)},
			{`;`, Punctuation, nil},
			{`^@.*?$`, CommentPreproc, nil},
			{Words(``, ``, `-ms-`, `mso-`, `-moz-`, `-o-`, `-xv-`, `-atsc-`, `-wap-`, `-khtml-`, `-webkit-`, `prince-`, `-ah-`, `-hp-`, `-ro-`, `-rim-`, `-tc-`), KeywordPseudo, nil},
			{`(align-content|align-items|align-self|alignment-baseline|all|animation|animation-delay|animation-direction|animation-duration|animation-fill-mode|animation-iteration-count|animation-name|animation-play-state|animation-timing-function|appearance|azimuth|backface-visibility|background|background-attachment|background-blend-mode|background-clip|background-color|background-image|background-origin|background-position|background-repeat|background-size|baseline-shift|bookmark-label|bookmark-level|bookmark-state|border|border-bottom|border-bottom-color|border-bottom-left-radius|border-bottom-right-radius|border-bottom-style|border-bottom-width|border-boundary|border-collapse|border-color|border-image|border-image-outset|border-image-repeat|border-image-slice|border-image-source|border-image-width|border-left|border-left-color|border-left-style|border-left-width|border-radius|border-right|border-right-color|border-right-style|border-right-width|border-spacing|border-style|border-top|border-top-color|border-top-left-radius|border-top-right-radius|border-top-style|border-top-width|border-width|bottom|box-decoration-break|box-shadow|box-sizing|box-snap|box-suppress|break-after|break-before|break-inside|caption-side|caret|caret-animation|caret-color|caret-shape|chains|clear|clip|clip-path|clip-rule|color|color-interpolation-filters|column-count|column-fill|column-gap|column-rule|column-rule-color|column-rule-style|column-rule-width|column-span|column-width|columns|content|counter-increment|counter-reset|counter-set|crop|cue|cue-after|cue-before|cursor|direction|display|dominant-baseline|elevation|empty-cells|filter|flex|flex-basis|flex-direction|flex-flow|flex-grow|flex-shrink|flex-wrap|float|float-defer|float-offset|float-reference|flood-color|flood-opacity|flow|flow-from|flow-into|font|font-family|font-feature-settings|font-kerning|font-language-override|font-size|font-size-adjust|font-stretch|font-style|font-synthesis|font-variant|font-variant-alternates|font-variant-caps|font-variant-east-asian|font-variant-ligatures|font-variant-numeric|font-variant-position|font-weight|footnote-display|footnote-policy|glyph-orientation-vertical|grid|grid-area|grid-auto-columns|grid-auto-flow|grid-auto-rows|grid-column|grid-column-end|grid-column-gap|grid-column-start|grid-gap|grid-row|grid-row-end|grid-row-gap|grid-row-start|grid-template|grid-template-areas|grid-template-columns|grid-template-rows|hanging-punctuation|height|hyphenate-character|hyphenate-limit-chars|hyphenate-limit-last|hyphenate-limit-lines|hyphenate-limit-zone|hyphens|image-orientation|image-resolution|initial-letter|initial-letter-align|initial-letter-wrap|isolation|justify-content|justify-items|justify-self|left|letter-spacing|lighting-color|line-break|line-grid|line-height|line-snap|list-style|list-style-image|list-style-position|list-style-type|margin|margin-bottom|margin-left|margin-right|margin-top|marker-side|marquee-direction|marquee-loop|marquee-speed|marquee-style|mask|mask-border|mask-border-mode|mask-border-outset|mask-border-repeat|mask-border-slice|mask-border-source|mask-border-width|mask-clip|mask-composite|mask-image|mask-mode|mask-origin|mask-position|mask-repeat|mask-size|mask-type|max-height|max-lines|max-width|min-height|min-width|mix-blend-mode|motion|motion-offset|motion-path|motion-rotation|move-to|nav-down|nav-left|nav-right|nav-up|object-fit|object-position|offset-after|offset-before|offset-end|offset-start|opacity|order|orphans|outline|outline-color|outline-offset|outline-style|outline-width|overflow|overflow-style|overflow-wrap|overflow-x|overflow-y|padding|padding-bottom|padding-left|padding-right|padding-top|page|page-break-after|page-break-before|page-break-inside|page-policy|pause|pause-after|pause-before|perspective|perspective-origin|pitch|pitch-range|play-during|polar-angle|polar-distance|position|presentation-level|quotes|region-fragment|resize|rest|rest-after|rest-before|richness|right|rotation|rotation-point|ruby-align|ruby-merge|ruby-position|running|scroll-snap-coordinate|scroll-snap-destination|scroll-snap-points-x|scroll-snap-points-y|scroll-snap-type|shape-image-threshold|shape-inside|shape-margin|shape-outside|size|speak|speak-as|speak-header|speak-numeral|speak-punctuation|speech-rate|stress|string-set|tab-size|table-layout|text-align|text-align-last|text-combine-upright|text-decoration|text-decoration-color|text-decoration-line|text-decoration-skip|text-decoration-style|text-emphasis|text-emphasis-color|text-emphasis-position|text-emphasis-style|text-indent|text-justify|text-orientation|text-overflow|text-shadow|text-space-collapse|text-space-trim|text-spacing|text-transform|text-underline-position|text-wrap|top|transform|transform-origin|transform-style|transition|transition-delay|transition-duration|transition-property|transition-timing-function|unicode-bidi|user-select|vertical-align|visibility|voice-balance|voice-duration|voice-family|voice-pitch|voice-range|voice-rate|voice-stress|voice-volume|volume|white-space|widows|width|will-change|word-break|word-spacing|word-wrap|wrap-after|wrap-before|wrap-flow|wrap-inside|wrap-through|writing-mode|z-index)(\s*)(\:)`, ByGroups(Keyword, Text, Punctuation), Push("value-start")},
			{`(--[a-zA-Z_][\w-]*)(\s*)(\:)`, ByGroups(NameVariable, Text, Punctuation), Push("value-start")},
			{`([a-zA-Z_][\w-]*)(\s*)(\:)`, ByGroups(Name, Text, Punctuation), Push("value-start")},
			{`/\*(?:.|\n)*?\*/`, Comment, nil},
		},
		"value-start": {
			Include("common-values"),
			{Words(``, `\b`, `align-content`, `align-items`, `align-self`, `alignment-baseline`, `all`, `animation`, `animation-delay`, `animation-direction`, `animation-duration`, `animation-fill-mode`, `animation-iteration-count`, `animation-name`, `animation-play-state`, `animation-timing-function`, `appearance`, `azimuth`, `backface-visibility`, `background`, `background-attachment`, `background-blend-mode`, `background-clip`, `background-color`, `background-image`, `background-origin`, `background-position`, `background-repeat`, `background-size`, `baseline-shift`, `bookmark-label`, `bookmark-level`, `bookmark-state`, `border`, `border-bottom`, `border-bottom-color`, `border-bottom-left-radius`, `border-bottom-right-radius`, `border-bottom-style`, `border-bottom-width`, `border-boundary`, `border-collapse`, `border-color`, `border-image`, `border-image-outset`, `border-image-repeat`, `border-image-slice`, `border-image-source`, `border-image-width`, `border-left`, `border-left-color`, `border-left-style`, `border-left-width`, `border-radius`, `border-right`, `border-right-color`, `border-right-style`, `border-right-width`, `border-spacing`, `border-style`, `border-top`, `border-top-color`, `border-top-left-radius`, `border-top-right-radius`, `border-top-style`, `border-top-width`, `border-width`, `bottom`, `box-decoration-break`, `box-shadow`, `box-sizing`, `box-snap`, `box-suppress`, `break-after`, `break-before`, `break-inside`, `caption-side`, `caret`, `caret-animation`, `caret-color`, `caret-shape`, `chains`, `clear`, `clip`, `clip-path`, `clip-rule`, `color`, `color-interpolation-filters`, `column-count`, `column-fill`, `column-gap`, `column-rule`, `column-rule-color`, `column-rule-style`, `column-rule-width`, `column-span`, `column-width`, `columns`, `content`, `counter-increment`, `counter-reset`, `counter-set`, `crop`, `cue`, `cue-after`, `cue-before`, `cursor`, `direction`, `display`, `dominant-baseline`, `elevation`, `empty-cells`, `filter`, `flex`, `flex-basis`, `flex-direction`, `flex-flow`, `flex-grow`, `flex-shrink`, `flex-wrap`, `float`, `float-defer`, `float-offset`, `float-reference`, `flood-color`, `flood-opacity`, `flow`, `flow-from`, `flow-into`, `font`, `font-family`, `font-feature-settings`, `font-kerning`, `font-language-override`, `font-size`, `font-size-adjust`, `font-stretch`, `font-style`, `font-synthesis`, `font-variant`, `font-variant-alternates`, `font-variant-caps`, `font-variant-east-asian`, `font-variant-ligatures`, `font-variant-numeric`, `font-variant-position`, `font-weight`, `footnote-display`, `footnote-policy`, `glyph-orientation-vertical`, `grid`, `grid-area`, `grid-auto-columns`, `grid-auto-flow`, `grid-auto-rows`, `grid-column`, `grid-column-end`, `grid-column-gap`, `grid-column-start`, `grid-gap`, `grid-row`, `grid-row-end`, `grid-row-gap`, `grid-row-start`, `grid-template`, `grid-template-areas`, `grid-template-columns`, `grid-template-rows`, `hanging-punctuation`, `height`, `hyphenate-character`, `hyphenate-limit-chars`, `hyphenate-limit-last`, `hyphenate-limit-lines`, `hyphenate-limit-zone`, `hyphens`, `image-orientation`, `image-resolution`, `initial-letter`, `initial-letter-align`, `initial-letter-wrap`, `isolation`, `justify-content`, `justify-items`, `justify-self`, `left`, `letter-spacing`, `lighting-color`, `line-break`, `line-grid`, `line-height`, `line-snap`, `list-style`, `list-style-image`, `list-style-position`, `list-style-type`, `margin`, `margin-bottom`, `margin-left`, `margin-right`, `margin-top`, `marker-side`, `marquee-direction`, `marquee-loop`, `marquee-speed`, `marquee-style`, `mask`, `mask-border`, `mask-border-mode`, `mask-border-outset`, `mask-border-repeat`, `mask-border-slice`, `mask-border-source`, `mask-border-width`, `mask-clip`, `mask-composite`, `mask-image`, `mask-mode`, `mask-origin`, `mask-position`, `mask-repeat`, `mask-size`, `mask-type`, `max-height`, `max-lines`, `max-width`, `min-height`, `min-width`, `mix-blend-mode`, `motion`, `motion-offset`, `motion-path`, `motion-rotation`, `move-to`, `nav-down`, `nav-left`, `nav-right`, `nav-up`, `object-fit`, `object-position`, `offset-after`, `offset-before`, `offset-end`, `offset-start`, `opacity`, `order`, `orphans`, `outline`, `outline-color`, `outline-offset`, `outline-style`, `outline-width`, `overflow`, `overflow-style`, `overflow-wrap`, `overflow-x`, `overflow-y`, `padding`, `padding-bottom`, `padding-left`, `padding-right`, `padding-top`, `page`, `page-break-after`, `page-break-before`, `page-break-inside`, `page-policy`, `pause`, `pause-after`, `pause-before`, `perspective`, `perspective-origin`, `pitch`, `pitch-range`, `play-during`, `polar-angle`, `polar-distance`, `position`, `presentation-level`, `quotes`, `region-fragment`, `resize`, `rest`, `rest-after`, `rest-before`, `richness`, `right`, `rotation`, `rotation-point`, `ruby-align`, `ruby-merge`, `ruby-position`, `running`, `scroll-snap-coordinate`, `scroll-snap-destination`, `scroll-snap-points-x`, `scroll-snap-points-y`, `scroll-snap-type`, `shape-image-threshold`, `shape-inside`, `shape-margin`, `shape-outside`, `size`, `speak`, `speak-as`, `speak-header`, `speak-numeral`, `speak-punctuation`, `speech-rate`, `stress`, `string-set`, `tab-size`, `table-layout`, `text-align`, `text-align-last`, `text-combine-upright`, `text-decoration`, `text-decoration-color`, `text-decoration-line`, `text-decoration-skip`, `text-decoration-style`, `text-emphasis`, `text-emphasis-color`, `text-emphasis-position`, `text-emphasis-style`, `text-indent`, `text-justify`, `text-orientation`, `text-overflow`, `text-shadow`, `text-space-collapse`, `text-space-trim`, `text-spacing`, `text-transform`, `text-underline-position`, `text-wrap`, `top`, `transform`, `transform-origin`, `transform-style`, `transition`, `transition-delay`, `transition-duration`, `transition-property`, `transition-timing-function`, `unicode-bidi`, `user-select`, `vertical-align`, `visibility`, `voice-balance`, `voice-duration`, `voice-family`, `voice-pitch`, `voice-range`, `voice-rate`, `voice-stress`, `voice-volume`, `volume`, `white-space`, `widows`, `width`, `will-change`, `word-break`, `word-spacing`, `word-wrap`, `wrap-after`, `wrap-before`, `wrap-flow`, `wrap-inside`, `wrap-through`, `writing-mode`, `z-index`), Keyword, nil},
			{`\!important`, CommentPreproc, nil},
			{`/\*(?:.|\n)*?\*/`, Comment, nil},
			Include("numeric-values"),
			{`[~^*!%&<>|+=@:./?-]+`, Operator, nil},
			{`[\[\](),]+`, Punctuation, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`[a-zA-Z_][\w-]*`, Name, nil},
			{`;`, Punctuation, Pop(1)},
			{`\}`, Punctuation, Pop(2)},
		},
		"function-start": {
			Include("common-values"),
			{`/\*(?:.|\n)*?\*/`, Comment, nil},
			Include("numeric-values"),
			{`[*+/-]`, Operator, nil},
			{`[,]`, Punctuation, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
			{`[a-zA-Z_-]\w*`, Name, nil},
			{`\(`, Punctuation, Push("atparenthesis")},
			{`\)`, Punctuation, Pop(1)},
		},
		"common-values": {
			{`\s+`, Text, nil},
			{Words(``, ``, `-ms-`, `mso-`, `-moz-`, `-o-`, `-xv-`, `-atsc-`, `-wap-`, `-khtml-`, `-webkit-`, `prince-`, `-ah-`, `-hp-`, `-ro-`, `-rim-`, `-tc-`), KeywordPseudo, nil},
			Include("urls"),
			{`(attr|blackness|blend|blenda|blur|brightness|calc|circle|color-mod|contrast|counter|cubic-bezier|device-cmyk|drop-shadow|ellipse|gray|grayscale|hsl|hsla|hue|hue-rotate|hwb|image|inset|invert|lightness|linear-gradient|matrix|matrix3d|opacity|perspective|polygon|radial-gradient|rect|repeating-linear-gradient|repeating-radial-gradient|rgb|rgba|rotate|rotate3d|rotateX|rotateY|rotateZ|saturate|saturation|scale|scale3d|scaleX|scaleY|scaleZ|sepia|shade|skewX|skewY|steps|tint|toggle|translate|translate3d|translateX|translateY|translateZ|whiteness)(\()`, ByGroups(NameBuiltin, Punctuation), Push("function-start")},
			{`([a-zA-Z_][\w-]+)(\()`, ByGroups(NameFunction, Punctuation), Push("function-start")},
			{Words(``, `\b`, `absolute`, `alias`, `all`, `all-petite-caps`, `all-scroll`, `all-small-caps`, `allow-end`, `alpha`, `alternate`, `alternate-reverse`, `always`, `armenian`, `auto`, `avoid`, `avoid-column`, `avoid-page`, `backwards`, `balance`, `baseline`, `below`, `blink`, `block`, `bold`, `bolder`, `border-box`, `both`, `bottom`, `box-decoration`, `break-word`, `capitalize`, `cell`, `center`, `circle`, `clip`, `clone`, `close-quote`, `col-resize`, `collapse`, `color`, `color-burn`, `color-dodge`, `column`, `column-reverse`, `compact`, `condensed`, `contain`, `container`, `content-box`, `context-menu`, `copy`, `cover`, `crisp-edges`, `crosshair`, `currentColor`, `cursive`, `darken`, `dashed`, `decimal`, `decimal-leading-zero`, `default`, `descendants`, `difference`, `digits`, `disc`, `distribute`, `dot`, `dotted`, `double`, `double-circle`, `e-resize`, `each-line`, `ease`, `ease-in`, `ease-in-out`, `ease-out`, `edges`, `ellipsis`, `end`, `ew-resize`, `exclusion`, `expanded`, `extra-condensed`, `extra-expanded`, `fantasy`, `fill`, `fill-box`, `filled`, `first`, `fixed`, `flat`, `flex`, `flex-end`, `flex-start`, `flip`, `force-end`, `forwards`, `from-image`, `full-width`, `geometricPrecision`, `georgian`, `groove`, `hanging`, `hard-light`, `help`, `hidden`, `hide`, `horizontal`, `hue`, `icon`, `infinite`, `inherit`, `initial`, `ink`, `inline`, `inline-block`, `inline-flex`, `inline-table`, `inset`, `inside`, `inter-word`, `invert`, `isolate`, `italic`, `justify`, `large`, `larger`, `last`, `left`, `lighten`, `lighter`, `line-through`, `linear`, `list-item`, `local`, `loose`, `lower-alpha`, `lower-greek`, `lower-latin`, `lower-roman`, `lowercase`, `ltr`, `luminance`, `luminosity`, `mandatory`, `manipulation`, `manual`, `margin-box`, `match-parent`, `medium`, `mixed`, `monospace`, `move`, `multiply`, `n-resize`, `ne-resize`, `nesw-resize`, `no-close-quote`, `no-drop`, `no-open-quote`, `no-repeat`, `none`, `normal`, `not-allowed`, `nowrap`, `ns-resize`, `nw-resize`, `nwse-resize`, `objects`, `oblique`, `off`, `on`, `open`, `open-quote`, `optimizeLegibility`, `optimizeSpeed`, `outset`, `outside`, `over`, `overlay`, `overline`, `padding-box`, `page`, `pan-down`, `pan-left`, `pan-right`, `pan-up`, `pan-x`, `pan-y`, `paused`, `petite-caps`, `pixelated`, `pointer`, `preserve-3d`, `progress`, `proximity`, `relative`, `repeat`, `repeat no-repeat`, `repeat-x`, `repeat-y`, `reverse`, `ridge`, `right`, `round`, `row`, `row-resize`, `row-reverse`, `rtl`, `ruby`, `ruby-base`, `ruby-base-container`, `ruby-text`, `ruby-text-container`, `run-in`, `running`, `s-resize`, `sans-serif`, `saturation`, `scale-down`, `screen`, `scroll`, `se-resize`, `semi-condensed`, `semi-expanded`, `separate`, `serif`, `sesame`, `show`, `sideways`, `sideways-left`, `sideways-right`, `slice`, `small`, `small-caps`, `smaller`, `smooth`, `snap`, `soft-light`, `solid`, `space`, `space-around`, `space-between`, `spaces`, `square`, `start`, `static`, `step-end`, `step-start`, `sticky`, `stretch`, `strict`, `stroke-box`, `style`, `sw-resize`, `table`, `table-caption`, `table-cell`, `table-column`, `table-column-group`, `table-footer-group`, `table-header-group`, `table-row`, `table-row-group`, `text`, `thick`, `thin`, `titling-caps`, `to`, `top`, `triangle`, `ultra-condensed`, `ultra-expanded`, `under`, `underline`, `unicase`, `unset`, `upper-alpha`, `upper-latin`, `upper-roman`, `uppercase`, `upright`, `use-glyph-orientation`, `vertical`, `vertical-text`, `view-box`, `visible`, `w-resize`, `wait`, `wavy`, `weight`, `weight style`, `wrap`, `wrap-reverse`, `x-large`, `x-small`, `xx-large`, `xx-small`, `zoom-in`, `zoom-out`), KeywordConstant, nil},
			{Words(``, `\b`, `above`, `aural`, `behind`, `bidi-override`, `center-left`, `center-right`, `cjk-ideographic`, `continuous`, `crop`, `cross`, `embed`, `far-left`, `far-right`, `fast`, `faster`, `hebrew`, `high`, `higher`, `hiragana`, `hiragana-iroha`, `katakana`, `katakana-iroha`, `landscape`, `left-side`, `leftwards`, `level`, `loud`, `low`, `lower`, `message-box`, `middle`, `mix`, `narrower`, `once`, `portrait`, `right-side`, `rightwards`, `silent`, `slow`, `slower`, `small-caption`, `soft`, `spell-out`, `status-bar`, `super`, `text-bottom`, `text-top`, `wider`, `x-fast`, `x-high`, `x-loud`, `x-low`, `x-soft`, `yes`, `pre`, `pre-wrap`, `pre-line`), KeywordConstant, nil},
			{Words(``, `\b`, `aliceblue`, `antiquewhite`, `aqua`, `aquamarine`, `azure`, `beige`, `bisque`, `black`, `blanchedalmond`, `blue`, `blueviolet`, `brown`, `burlywood`, `cadetblue`, `chartreuse`, `chocolate`, `coral`, `cornflowerblue`, `cornsilk`, `crimson`, `cyan`, `darkblue`, `darkcyan`, `darkgoldenrod`, `darkgray`, `darkgreen`, `darkgrey`, `darkkhaki`, `darkmagenta`, `darkolivegreen`, `darkorange`, `darkorchid`, `darkred`, `darksalmon`, `darkseagreen`, `darkslateblue`, `darkslategray`, `darkslategrey`, `darkturquoise`, `darkviolet`, `deeppink`, `deepskyblue`, `dimgray`, `dimgrey`, `dodgerblue`, `firebrick`, `floralwhite`, `forestgreen`, `fuchsia`, `gainsboro`, `ghostwhite`, `gold`, `goldenrod`, `gray`, `green`, `greenyellow`, `grey`, `honeydew`, `hotpink`, `indianred`, `indigo`, `ivory`, `khaki`, `lavender`, `lavenderblush`, `lawngreen`, `lemonchiffon`, `lightblue`, `lightcoral`, `lightcyan`, `lightgoldenrodyellow`, `lightgray`, `lightgreen`, `lightgrey`, `lightpink`, `lightsalmon`, `lightseagreen`, `lightskyblue`, `lightslategray`, `lightslategrey`, `lightsteelblue`, `lightyellow`, `lime`, `limegreen`, `linen`, `magenta`, `maroon`, `mediumaquamarine`, `mediumblue`, `mediumorchid`, `mediumpurple`, `mediumseagreen`, `mediumslateblue`, `mediumspringgreen`, `mediumturquoise`, `mediumvioletred`, `midnightblue`, `mintcream`, `mistyrose`, `moccasin`, `navajowhite`, `navy`, `oldlace`, `olive`, `olivedrab`, `orange`, `orangered`, `orchid`, `palegoldenrod`, `palegreen`, `paleturquoise`, `palevioletred`, `papayawhip`, `peachpuff`, `peru`, `pink`, `plum`, `powderblue`, `purple`, `rebeccapurple`, `red`, `rosybrown`, `royalblue`, `saddlebrown`, `salmon`, `sandybrown`, `seagreen`, `seashell`, `sienna`, `silver`, `skyblue`, `slateblue`, `slategray`, `slategrey`, `snow`, `springgreen`, `steelblue`, `tan`, `teal`, `thistle`, `tomato`, `turquoise`, `violet`, `wheat`, `white`, `whitesmoke`, `yellow`, `yellowgreen`, `transparent`), KeywordConstant, nil},
		},
		"urls": {
			{`(url)(\()(".*?")(\))`, ByGroups(NameBuiltin, Punctuation, LiteralStringDouble, Punctuation), nil},
			{`(url)(\()('.*?')(\))`, ByGroups(NameBuiltin, Punctuation, LiteralStringSingle, Punctuation), nil},
			{`(url)(\()(.*?)(\))`, ByGroups(NameBuiltin, Punctuation, LiteralStringOther, Punctuation), nil},
		},
		"numeric-values": {
			{`\#[a-zA-Z0-9]{1,6}`, LiteralNumberHex, nil},
			{`[+\-]?[0-9]*[.][0-9]+`, LiteralNumberFloat, Push("numeric-end")},
			{`[+\-]?[0-9]+`, LiteralNumberInteger, Push("numeric-end")},
		},
		"numeric-end": {
			{Words(``, `\b`, `deg`, `grad`, `rad`, `turn`, `Hz`, `kHz`, `em`, `ex`, `ch`, `rem`, `vh`, `vw`, `vmin`, `vmax`, `px`, `mm`, `cm`, `in`, `pt`, `pc`, `q`, `dpi`, `dpcm`, `dppx`, `s`, `ms`), KeywordType, nil},
			{`%`, KeywordType, nil},
			Default(Pop(1)),
		},
	}
}
