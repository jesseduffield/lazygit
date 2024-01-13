package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Scss lexer.
var Scss = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "SCSS",
		Aliases:         []string{"scss"},
		Filenames:       []string{"*.scss"},
		MimeTypes:       []string{"text/x-scss"},
		NotMultiline:    true,
		DotAll:          true,
		CaseInsensitive: true,
	},
	scssRules,
))

func scssRules() Rules {
	cssProperties := []string{`src`, `align-content`, `align-items`, `align-self`, `all`, `animation`, `animation-delay`, `animation-direction`, `animation-duration`, `animation-fill-mode`, `animation-iteration-count`, `animation-name`, `animation-play-state`, `animation-timing-function`, `appearance`, `aspect-ratio`, `backface-visibility`, `background`, `background-attachment`, `background-blend-mode`, `background-clip`, `background-color`, `background-image`, `background-origin`, `background-position`, `background-repeat`, `background-size`, `block-size`, `border`, `border-block`, `border-block-color`, `border-block-end`, `border-block-end-color`, `border-block-end-style`, `border-block-end-width`, `border-block-start`, `border-block-start-color`, `border-block-start-style`, `border-block-start-width`, `border-block-style`, `border-block-width`, `border-bottom`, `border-bottom-color`, `border-bottom-left-radius`, `border-bottom-right-radius`, `border-bottom-style`, `border-bottom-width`, `border-collapse`, `border-color`, `border-end-end-radius`, `border-end-start-radius`, `border-image`, `border-image-outset`, `border-image-repeat`, `border-image-slice`, `border-image-source`, `border-image-width`, `border-inline`, `border-inline-color`, `border-inline-end`, `border-inline-end-color`, `border-inline-end-style`, `border-inline-end-width`, `border-inline-start`, `border-inline-start-color`, `border-inline-start-style`, `border-inline-start-width`, `border-inline-style`, `border-inline-width`, `border-left`, `border-left-color`, `border-left-style`, `border-left-width`, `border-radius`, `border-right`, `border-right-color`, `border-right-style`, `border-right-width`, `border-spacing`, `border-start-end-radius`, `border-start-start-radius`, `border-style`, `border-top`, `border-top-color`, `border-top-left-radius`, `border-top-right-radius`, `border-top-style`, `border-top-width`, `border-width`, `bottom`, `box-decoration-break`, `box-shadow`, `box-sizing`, `break-after`, `break-before`, `break-inside`, `caption-side`, `caret-color`, `clear`, `clip`, `clip-path`, `color`, `color-adjust`, `column-count`, `column-fill`, `column-gap`, `column-rule`, `column-rule-color`, `column-rule-style`, `column-rule-width`, `column-span`, `column-width`, `columns`, `content`, `content-visibility`, `counter-increment`, `counter-reset`, `counter-set`, `cursor`, `direction`, `display`, `empty-cells`, `filter`, `flex`, `flex-basis`, `flex-direction`, `flex-flow`, `flex-grow`, `flex-shrink`, `flex-wrap`, `float`, `font`, `font-family`, `font-feature-settings`, `font-kerning`, `font-language-override`, `font-optical-sizing`, `font-size`, `font-size-adjust`, `font-stretch`, `font-style`, `font-synthesis`, `font-variant`, `font-variant-caps`, `font-variant-east-asian`, `font-variant-emoji`, `font-variant-ligatures`, `font-variant-numeric`, `font-variant-position`, `font-weight`, `footnote-display`, `footnote-policy`, `forced-color-adjust`, `gap`, `grid`, `grid-area`, `grid-auto-columns`, `grid-auto-flow`, `grid-auto-rows`, `grid-column`, `grid-column-end`, `grid-column-start`, `grid-row`, `grid-row-end`, `grid-row-start`, `grid-template`, `grid-template-areas`, `grid-template-columns`, `grid-template-rows`, `hanging-punctuation`, `height`, `hyphens`, `image-orientation`, `image-rendering`, `image-resolution`, `initial-letter`, `initial-letter-align`, `initial-letter-wrap`, `inline-size`, `inline-sizing`, `inset`, `inset-block`, `inset-block-end`, `inset-block-start`, `inset-inline`, `inset-inline-end`, `inset-inline-start`, `isolation`, `justify-content`, `justify-items`, `justify-self`, `left`, `letter-spacing`, `lighting-color`, `line-break`, `line-clamp`, `line-grid`, `line-height`, `line-padding`, `line-snap`, `list-style`, `list-style-image`, `list-style-position`, `list-style-type`, `margin`, `margin-block`, `margin-block-end`, `margin-block-start`, `margin-bottom`, `margin-break`, `margin-inline`, `margin-inline-end`, `margin-inline-start`, `margin-left`, `margin-right`, `margin-top`, `margin-trim`, `mask`, `mask-border`, `mask-border-mode`, `mask-border-repeat`, `mask-border-slice`, `mask-border-source`, `mask-border-width`, `mask-clip`, `mask-composite`, `mask-image`, `mask-mode`, `mask-origin`, `mask-position`, `mask-repeat`, `mask-size`, `mask-type`, `max-block-size`, `max-height`, `max-lines`, `max-width`, `min-height`, `min-width`, `mix-blend-mode`, `nav-down`, `nav-left`, `nav-right`, `nav-up`, `object-fit`, `object-position`, `offset`, `offset-anchor`, `offset-distance`, `offset-path`, `offset-position`, `offset-rotate`, `opacity`, `order`, `orphans`, `outline`, `outline-color`, `outline-offset`, `outline-style`, `outline-width`, `overflow`, `overflow-anchor`, `overflow-block`, `overflow-clip-margin`, `overflow-inline`, `overflow-wrap`, `overflow-x`, `overflow-y`, `overscroll-behavior`, `overscroll-behavior-block`, `overscroll-behavior-inline`, `overscroll-behavior-x`, `overscroll-behavior-y`, `padding`, `padding-block`, `padding-block-end`, `padding-block-start`, `padding-bottom`, `padding-inline`, `padding-inline-end`, `padding-inline-start`, `padding-left`, `padding-right`, `padding-top`, `page`, `page-break-after`, `page-break-before`, `page-break-inside`, `page-orientation`, `perspective`, `perspective-origin`, `place-content`, `place-items`, `place-self`, `pointer-events`, `position`, `quotes`, `resize`, `right`, `rotate`, `row-gap`, `scale`, `scroll-behavior`, `scroll-margin`, `scroll-margin-block`, `scroll-margin-block-end`, `scroll-margin-block-start`, `scroll-margin-bottom`, `scroll-margin-inline`, `scroll-margin-inline-end`, `scroll-margin-inline-start`, `scroll-margin-left`, `scroll-margin-right`, `scroll-margin-top`, `scroll-padding`, `scroll-padding-block`, `scroll-padding-block-end`, `scroll-padding-block-start`, `scroll-padding-bottom`, `scroll-padding-inline`, `scroll-padding-inline-end`, `scroll-padding-inline-start`, `scroll-padding-left`, `scroll-padding-right`, `scroll-padding-top`, `scroll-snap-align`, `scroll-snap-stop`, `scroll-snap-type`, `scrollbar-color`, `scrollbar-gutter`, `scrollbar-width`, `shape-image-threshold`, `shape-inside`, `shape-margin`, `shape-outside`, `shape-padding`, `spatial-navigation-action`, `spatial-navigation-contain`, `spatial-navigation-function`, `string-set`, `tab-size`, `table-layout`, `text-align`, `text-align-all`, `text-align-last`, `text-combine-upright`, `text-decoration`, `text-decoration-color`, `text-decoration-line`, `text-decoration-skip`, `text-decoration-style`, `text-decoration-thickness`, `text-emphasis`, `text-emphasis-position`, `text-emphasis-style`, `text-group-align`, `text-indent`, `text-justify`, `text-orientation`, `text-overflow`, `text-rendering`, `text-shadow`, `text-size-adjust`, `text-space-trim`, `text-spacing`, `text-transform`, `text-underline-offset`, `text-underline-position`, `text-wrap`, `top`, `touch-action`, `transform`, `transform-box`, `transform-origin`, `transform-style`, `transition`, `transition-delay`, `transition-duration`, `transition-property`, `transition-timing-function`, `translate`, `unicode-bidi`, `user-select`, `vertical-align`, `visibility`, `white-space`, `widows`, `width`, `will-change`, `word-break`, `word-spacing`, `word-wrap`, `wrap-after`, `wrap-before`, `wrap-flow`, `wrap-inside`, `wrap-through`, `writing-mode`, `z-index`}

	cssPropertyValues := []string{
		`stretch`, `flex-start`, `flex-end`, `center`, `space-between`, `space-around`, `normal`, `baseline`, `first`, `last`, `space-evenly`, `start`, `end`, `safe`, `unsafe`, `self-start`, `self-end`, `auto`, `initial`, `inherit`, `unset`, `revert`, `reverse`, `alternate-reverse`, `alternate`, `none`, `forwards`, `both`, `backwards`, `infinite`, `running`, `paused`, `ease`, `linear`, `step-start`, `step-end`, `ease-in`, `ease-out`, `ease-in-out`, `cubic-bezier`, `textfield`, `menulist-button`, `searchfield`, `textarea`, `push-button`, `slider-horizontal`, `checkbox`, `radio`, `square-button`, `menulist`, `listbox`, `meter`, `progress-bar`, `button`, `visible`, `hidden`, `scroll`, `fixed`, `local`, `multiply`, `screen`, `overlay`, `darken`, `lighten`, `color-dodge`, `color-burn`, `hard-light`, `soft-light`, `difference`, `exclusion`, `hue`, `saturation`, `color`, `luminosity`, `border-box`, `padding-box`, `content-box`, `text`, `transparent`, `currentColor`, `url`, `element`, `image`, `image-set`, `cross-fade`, `top`, `bottom`, `left`, `right`, `repeat`, `no-repeat`, `space`, `round`, `repeat-y`, `repeat-x`, `cover`, `contain`, `max-content`, `min-content`, `available`, `fit-content`, `dotted`, `dashed`, `solid`, `double`, `groove`, `ridge`, `inset`, `outset`, `medium`, `thin`, `thick`, `separate`, `collapse`, `slice`, `clone`, `avoid`, `always`, `all`, `avoid-page`, `page`, `recto`, `verso`, `avoid-column`, `column`, `avoid-region`, `region`, `inline-start`, `inline-end`, `rect`, `circle`, `ellipse`, `polygon`, `path`, `margin-box`, `fill-box`, `stroke-box`, `view-box`, `economy`, `exact`, `balance`, `balance-all`, `open-quote`, `close-quote`, `no-open-quote`, `no-close-quote`, `contents`, `attr`, `target-counter`, `target-text`, `leader`, `default`, `context-menu`, `help`, `pointer`, `progress`, `wait`, `cell`, `crosshair`, `vertical-text`, `alias`, `copy`, `move`, `no-drop`, `not-allowed`, `grab`, `grabbing`, `all-scroll`, `col-resize`, `row-resize`, `n-resize`, `s-resize`, `w-resize`, `ne-resize`, `nw-resize`, `se-resize`, `sw-resize`, `ew-resize`, `ns-resize`, `nesw-resize`, `nwse-resize`, `zoom-in`, `zoom-out`, `ltr`, `rtl`, `inline`, `block`, `inline-block`, `inline-table`, `run-in`, `flow`, `flow-root`, `table`, `flex`, `grid`, `ruby`, `list-item`, `table-row-group`, `table-header-group`, `table-footer-group`, `table-row`, `table-cell`, `table-column-group`, `table-column`, `table-caption`, `ruby-base`, `ruby-text`, `ruby-base-container`, `ruby-text-container`, `inline-flex`, `inline-grid`, `show`, `hide`, `url;`, `blur`, `brightness`, `contrast`, `drop-shadow`, `grayscale`, `hue-rotate`, `invert`, `opacity`, `saturate`, `sepia`, `content`, `row`, `row-reverse`, `column-reverse`, `nowrap`, `wrap`, `wrap-reverse`, `block-start`, `block-end`, `caption`, `icon`, `menu`, `message-box`, `small-caption`, `status-bar`, `serif`, `sans-serif`, `cursive`, `fantasy`, `monospace`, `system-ui`, `emoji`, `math`, `fangsong`, `ui-serif`, `ui-sans-serif`, `ui-monospace`, `ui-rounded`, `xx-small`, `x-small`, `small`, `large`, `x-large`, `xx-large`, `xxx-large`, `smaller`, `larger`, `semi-condensed`, `condensed`, `extra-condensed`, `ultra-condensed`, `semi-expanded`, `expanded`, `extra-expanded`, `ultra-expanded`, `italic`, `weight`, `style`, `small-caps`, `all-small-caps`, `petite-caps`, `all-petite-caps`, `unicase`, `titling-caps`, `jis78`, `jis83`, `jis90`, `jis04`, `simplified`, `traditional`, `full-width`, `proportional-width`, `unicode`, `common-ligatures`, `no-common-ligatures`, `discretionary-ligatures`, `no-discretionary-ligatures`, `historical-ligatures`, `no-historical-ligatures`, `contextual`, `no-contextual`, `ordinal`, `slashed-zero`, `lining-nums`, `oldstyle-nums`, `proportional-nums`, `tabular-nums`, `diagonal-fractions`, `stacked-fractions`, `sub`, `super`, `bold`, `lighter`, `bolder`, `line`, `auto‑flow`, `minmax`, `auto;`, `dense`, `subgrid`, `masonry`, `force-end`, `allow-end`, `manual`, `from-image`, `smooth`, `high-quality`, `crisp-edges`, `pixelated`, `snap`, `alphabetic`, `ideographic`, `hebrew`, `hanging`, `isolate`, `legacy`, `white`, `strict`, `loose`, `anywhere`, `match-parent`, `create`, `outside`, `inside`, `disc`, `square`, `decimal`, `symbols`, `cjk-decimal`, `decimal-leading-zero`, `lower-roman`, `upper-roman`, `lower-greek`, `lower-alpha`, `lower-latin`, `upper-alpha`, `upper-latin`, `arabic-indic`, `armenian`, `bengali`, `cambodian`, `cjk-earthly-branch`, `cjk-heavenly-stem`, `cjk-ideographic`, `devanagari`, `ethiopic-numeric`, `georgian`, `gujarati`, `gurmukhi`, `hiragana`, `hiragana-iroha`, `japanese-formal`, `japanese-informal`, `kannada`, `katakana`, `katakana-iroha`, `khmer`, `korean-hangul-formal`, `korean-hanja-formal`, `korean-hanja-informal`, `lao`, `lower-armenian`, `malayalam`, `mongolian`, `myanmar`, `oriya`, `persian`, `simp-chinese-formal`, `simp-chinese-informal`, `tamil`, `telugu`, `thai`, `tibetan`, `trad-chinese-formal`, `trad-chinese-informal`, `upper-armenian`, `disclosure-open`, `disclosure-closed`, `keep`, `discard`, `in-flow`, `alpha`, `luminance`, `no-clip`, `add`, `subtract`, `intersect`, `exclude`, `match-source`, `current`, `fill`, `scale-down`, `ray`, `clip`, `break-word`, `upright`, `rotate-left`, `rotate-right`, `static`, `relative`, `absolute`, `sticky`, `horizontal`, `vertical`, `x`, `y`, `z`, `mandatory`, `proximity`, `dark`, `light`, `stable`, `force`, `outside-shape`, `display`, `focus`, `justify`, `justify-all`, `underline`, `overline`, `line-through`, `blink`, `spelling-error`, `grammar-error`, `wavy`, `from-font`, `over`, `under`, `filled`, `open`, `dot`, `double-circle`, `triangle`, `sesame`, `inter-word`, `inter-character`, `mixed`, `sideways`, `ellipsis`, `optimizeSpeed`, `optimizeLegibility`, `geometricPrecision`, `trim-inner`, `discard-before`, `discard-after`, `trim-start`, `space-start`, `space-first`, `trim-end`, `space-end`, `trim-adjacent`, `space-adjacent`, `no-compress`, `ideograph-alpha`, `ideograph-numeric`, `punctuation`, `capitalize`, `uppercase`, `lowercase`, `full-size-kana`, `pretty`, `pan-x`, `pan-left`, `pan-right`, `pan-y`, `pan-up`, `pan-down`, `manipulation`, `matrix`, `matrix3d`, `translate`, `translateX`, `translateY`, `translate3d`, `translateZ`, `scale`, `scaleX`, `scaleY`, `scaleZ`, `rotate`, `rotate3d`, `rotateX`, `rotateY`, `rotateZ`, `perspective`, `skew`, `skewX`, `skewY`, `flat`, `preserve-3d`, `border-bottom`, `embed`, `bidi-override`, `isolate-override`, `plaintext`, `middle`, `text-top`, `text-bottom`, `pre`, `pre-wrap`, `pre-line`, `scroll-position`, `transform`, `width`, `keep-all`, `break-all`, `avoid-line`, `avoid-flex`, `minimum`, `maximum`, `clear`, `horizontal-tb`, `vertical-rl`, `vertical-lr`, `sideways-rl`, `sideways-lr`,
	}

	cssColorNames := []string{`AliceBlue`, `AntiqueWhite`, `Aqua`, `Aquamarine`, `Azure`, `Beige`, `Bisque`, `Black`, `BlanchedAlmond`, `Blue`, `BlueViolet`, `Brown`, `BurlyWood`, `CadetBlue`, `Chartreuse`, `Chocolate`, `Coral`, `CornflowerBlue`, `Cornsilk`, `Crimson`, `Cyan`, `DarkBlue`, `DarkCyan`, `DarkGoldenRod`, `DarkGray`, `DarkGrey`, `DarkGreen`, `DarkKhaki`, `DarkMagenta`, `DarkOliveGreen`, `DarkOrange`, `DarkOrchid`, `DarkRed`, `DarkSalmon`, `DarkSeaGreen`, `DarkSlateBlue`, `DarkSlateGray`, `DarkSlateGrey`, `DarkTurquoise`, `DarkViolet`, `DeepPink`, `DeepSkyBlue`, `DimGray`, `DimGrey`, `DodgerBlue`, `FireBrick`, `FloralWhite`, `ForestGreen`, `Fuchsia`, `Gainsboro`, `GhostWhite`, `Gold`, `GoldenRod`, `Gray`, `Grey`, `Green`, `GreenYellow`, `HoneyDew`, `HotPink`, `IndianRed`, `Indigo`, `Ivory`, `Khaki`, `Lavender`, `LavenderBlush`, `LawnGreen`, `LemonChiffon`, `LightBlue`, `LightCoral`, `LightCyan`, `LightGoldenRodYellow`, `LightGray`, `LightGrey`, `LightGreen`, `LightPink`, `LightSalmon`, `LightSeaGreen`, `LightSkyBlue`, `LightSlateGray`, `LightSlateGrey`, `LightSteelBlue`, `LightYellow`, `Lime`, `LimeGreen`, `Linen`, `Magenta`, `Maroon`, `MediumAquaMarine`, `MediumBlue`, `MediumOrchid`, `MediumPurple`, `MediumSeaGreen`, `MediumSlateBlue`, `MediumSpringGreen`, `MediumTurquoise`, `MediumVioletRed`, `MidnightBlue`, `MintCream`, `MistyRose`, `Moccasin`, `NavajoWhite`, `Navy`, `OldLace`, `Olive`, `OliveDrab`, `Orange`, `OrangeRed`, `Orchid`, `PaleGoldenRod`, `PaleGreen`, `PaleTurquoise`, `PaleVioletRed`, `PapayaWhip`, `PeachPuff`, `Peru`, `Pink`, `Plum`, `PowderBlue`, `Purple`, `RebeccaPurple`, `Red`, `RosyBrown`, `RoyalBlue`, `SaddleBrown`, `Salmon`, `SandyBrown`, `SeaGreen`, `SeaShell`, `Sienna`, `Silver`, `SkyBlue`, `SlateBlue`, `SlateGray`, `SlateGrey`, `Snow`, `SpringGreen`, `SteelBlue`, `Tan`, `Teal`, `Thistle`, `Tomato`, `Turquoise`, `Violet`, `Wheat`, `White`, `WhiteSmoke`, `Yellow`, `YellowGreen`}

	const cssPropertyPattern = `(?<!:\s*)[^ \t:="\[;(),.]+(?=[ \t]*:(?:#{|[^{])+?[;,])`

	return Rules{
		"root": {
			{`[{}()]`, Punctuation, nil},
			{`\s+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`@import`, Keyword, Push("value")},
			{
				`(@(?:use|forward))(\s+)([^\s;]+)(?:(\s+)(as)(\s+)(\w+|\*))?(?:(\s+)(with))?`,
				ByGroups(Keyword, Text, UsingSelf("selector"), Text, Keyword, Text, NameNamespace, Text, Keyword),
				Push("value"),
			},
			{`@for`, Keyword, Push("for")},
			{`@each`, Keyword, Push("each")},
			{`@(debug|warn|else if|if|while|return)`, Keyword, Push("value")},
			{`(@(?:mixin|function))( [\w-]+)`, ByGroups(Keyword, NameFunction), Push("value")},
			{`(@include)( [\w-]+)`, ByGroups(Keyword, NameDecorator), Push("value")},
			{`@extend`, Keyword, Push("selector")},
			{`(@media)(\s+)`, ByGroups(Keyword, Text), Push("value")},
			{`@[\w-]+`, Keyword, Push("selector")},
			{`(\$[\w-]*\w)([ \t]*:)`, ByGroups(NameVariable, Operator), Push("value")},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{Words(`\b`, `(?=\s*:)`, cssProperties...), NameAttribute, Push("attr")},
			{cssPropertyPattern, NameAttribute, Push("attr")},
			Default(Push("selector")),
		},
		"attr": {
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`[ \t]*:`, Operator, Push("value")},
			Default(Pop(1)),
		},
		"inline-comment": {
			{`(\\#|#(?=[^{])|\*(?=[^/])|[^#*])+`, CommentMultiline, nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`\*/`, Comment, Pop(1)},
		},
		"value": {
			{`[ \t]+`, Text, nil},
			{`!(important|default|global)`, Keyword, nil},
			{`[!$][\w-]+`, NameVariable, nil},
			{`url\(`, LiteralStringOther, Push("string-url")},
			{`[a-z_-][\w-]*(?=\()`, NameFunction, nil},
			{Words(`\b`, `(?=\s*:)`, cssProperties...), NameAttribute, nil},
			{cssPropertyPattern, NameAttribute, nil},
			{Words(`\b`, `\b`, cssPropertyValues...), NameEntity, nil},
			{Words(`(?i)\b`, `\b`, cssColorNames...), NameConstant, nil},
			{`(true|false)`, NamePseudo, nil},
			{`(and|or|not)`, OperatorWord, nil},
			{`/\*`, CommentMultiline, Push("inline-comment")},
			{`//[^\n]*`, CommentSingle, nil},
			{`\#[a-z0-9]{1,6}`, LiteralNumberHex, nil},
			{`(-)?(\d+)(\%|[a-z]+)?`, ByGroups(Operator, LiteralNumberInteger, KeywordType), nil},
			{`(-?)(\d*\.\d+)(\%|[a-z]+)?`, ByGroups(Operator, LiteralNumberFloat, KeywordType), nil},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`[~^*!&%<>|+=@:,./?-]+`, Operator, nil},
			{`[\[\]()]+`, Punctuation, nil},
			{`"`, LiteralStringDouble, Push("string-double")},
			{`'`, LiteralStringSingle, Push("string-single")},
			{`[a-z_-][\w-]*`, Name, nil},
			{`\n`, Text, nil},
			{`[;{}]`, Punctuation, Pop(1)},
		},
		"interpolation": {
			{`\}`, LiteralStringInterpol, Pop(1)},
			Include("value"),
		},
		"selector": {
			{`[ \t]+`, Text, nil},
			{`\:`, NameDecorator, Push("pseudo-class")},
			{`\.`, NameClass, Push("class")},
			{`#\{`, LiteralStringInterpol, Push("interpolation")},
			{`\#`, NameNamespace, Push("id")},
			{`&`, Keyword, nil},
			{`[~^*!&\[\]()<>|+=@:,./?-]`, Operator, nil},
			{`(%)([\w-]+)`, ByGroups(Operator, NameClass), nil},
			{`"`, LiteralStringDouble, Push("string-double")},
			{`'`, LiteralStringSingle, Push("string-single")},
			{`\n`, Text, nil},
			{`[;{}]`, Punctuation, Pop(1)},
			{`[\w-]+`, NameTag, nil},
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
		"each": {
			{`in`, OperatorWord, nil},
			Include("value"),
		},
	}
}
