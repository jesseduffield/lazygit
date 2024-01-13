// Copyright (c) 2019, David Kitchen <david@buro9.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the organisation (Microcosm) nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package css

import (
	"regexp"
	"strings"
)

var (
	defaultStyleHandlers = map[string]func(string) bool{
		"align-content":              AlignContentHandler,
		"align-items":                AlignItemsHandler,
		"align-self":                 AlignSelfHandler,
		"all":                        AllHandler,
		"animation":                  AnimationHandler,
		"animation-delay":            AnimationDelayHandler,
		"animation-direction":        AnimationDirectionHandler,
		"animation-duration":         AnimationDurationHandler,
		"animation-fill-mode":        AnimationFillModeHandler,
		"animation-iteration-count":  AnimationIterationCountHandler,
		"animation-name":             AnimationNameHandler,
		"animation-play-state":       AnimationPlayStateHandler,
		"animation-timing-function":  TimingFunctionHandler,
		"backface-visibility":        BackfaceVisibilityHandler,
		"background":                 BackgroundHandler,
		"background-attachment":      BackgroundAttachmentHandler,
		"background-blend-mode":      BackgroundBlendModeHandler,
		"background-clip":            BackgroundClipHandler,
		"background-color":           ColorHandler,
		"background-image":           ImageHandler,
		"background-origin":          BackgroundOriginHandler,
		"background-position":        BackgroundPositionHandler,
		"background-repeat":          BackgroundRepeatHandler,
		"background-size":            BackgroundSizeHandler,
		"border":                     BorderHandler,
		"border-bottom":              BorderSideHandler,
		"border-bottom-color":        ColorHandler,
		"border-bottom-left-radius":  BorderSideRadiusHandler,
		"border-bottom-right-radius": BorderSideRadiusHandler,
		"border-bottom-style":        BorderSideStyleHandler,
		"border-bottom-width":        BorderSideWidthHandler,
		"border-collapse":            BorderCollapseHandler,
		"border-color":               ColorHandler,
		"border-image":               BorderImageHandler,
		"border-image-outset":        BorderImageOutsetHandler,
		"border-image-repeat":        BorderImageRepeatHandler,
		"border-image-slice":         BorderImageSliceHandler,
		"border-image-source":        ImageHandler,
		"border-image-width":         BorderImageWidthHandler,
		"border-left":                BorderSideHandler,
		"border-left-color":          ColorHandler,
		"border-left-style":          BorderSideStyleHandler,
		"border-left-width":          BorderSideWidthHandler,
		"border-radius":              BorderRadiusHandler,
		"border-right":               BorderSideHandler,
		"border-right-color":         ColorHandler,
		"border-right-style":         BorderSideStyleHandler,
		"border-right-width":         BorderSideWidthHandler,
		"border-spacing":             BorderSpacingHandler,
		"border-style":               BorderStyleHandler,
		"border-top":                 BorderSideHandler,
		"border-top-color":           ColorHandler,
		"border-top-left-radius":     BorderSideRadiusHandler,
		"border-top-right-radius":    BorderSideRadiusHandler,
		"border-top-style":           BorderSideStyleHandler,
		"border-top-width":           BorderSideWidthHandler,
		"border-width":               BorderWidthHandler,
		"bottom":                     SideHandler,
		"box-decoration-break":       BoxDecorationBreakHandler,
		"box-shadow":                 BoxShadowHandler,
		"box-sizing":                 BoxSizingHandler,
		"break-after":                BreakBeforeAfterHandler,
		"break-before":               BreakBeforeAfterHandler,
		"break-inside":               BreakInsideHandler,
		"caption-side":               CaptionSideHandler,
		"caret-color":                CaretColorHandler,
		"clear":                      ClearHandler,
		"clip":                       ClipHandler,
		"color":                      ColorHandler,
		"column-count":               ColumnCountHandler,
		"column-fill":                ColumnFillHandler,
		"column-gap":                 ColumnGapHandler,
		"column-rule":                ColumnRuleHandler,
		"column-rule-color":          ColorHandler,
		"column-rule-style":          BorderSideStyleHandler,
		"column-rule-width":          ColumnRuleWidthHandler,
		"column-span":                ColumnSpanHandler,
		"column-width":               ColumnWidthHandler,
		"columns":                    ColumnsHandler,
		"cursor":                     CursorHandler,
		"direction":                  DirectionHandler,
		"display":                    DisplayHandler,
		"empty-cells":                EmptyCellsHandler,
		"filter":                     FilterHandler,
		"flex":                       FlexHandler,
		"flex-basis":                 FlexBasisHandler,
		"flex-direction":             FlexDirectionHandler,
		"flex-flow":                  FlexFlowHandler,
		"flex-grow":                  FlexGrowHandler,
		"flex-shrink":                FlexGrowHandler,
		"flex-wrap":                  FlexWrapHandler,
		"float":                      FloatHandler,
		"font":                       FontHandler,
		"font-family":                FontFamilyHandler,
		"font-kerning":               FontKerningHandler,
		"font-language-override":     FontLanguageOverrideHandler,
		"font-size":                  FontSizeHandler,
		"font-size-adjust":           FontSizeAdjustHandler,
		"font-stretch":               FontStretchHandler,
		"font-style":                 FontStyleHandler,
		"font-synthesis":             FontSynthesisHandler,
		"font-variant":               FontVariantHandler,
		"font-variant-caps":          FontVariantCapsHandler,
		"font-variant-position":      FontVariantPositionHandler,
		"font-weight":                FontWeightHandler,
		"grid":                       GridHandler,
		"grid-area":                  GridAreaHandler,
		"grid-auto-columns":          GridAutoColumnsHandler,
		"grid-auto-flow":             GridAutoFlowHandler,
		"grid-auto-rows":             GridAutoColumnsHandler,
		"grid-column":                GridColumnHandler,
		"grid-column-end":            GridAxisStartEndHandler,
		"grid-column-gap":            LengthHandler,
		"grid-column-start":          GridAxisStartEndHandler,
		"grid-gap":                   GridGapHandler,
		"grid-row":                   GridRowHandler,
		"grid-row-end":               GridAxisStartEndHandler,
		"grid-row-gap":               LengthHandler,
		"grid-row-start":             GridAxisStartEndHandler,
		"grid-template":              GridTemplateHandler,
		"grid-template-areas":        GridTemplateAreasHandler,
		"grid-template-columns":      GridTemplateColumnsHandler,
		"grid-template-rows":         GridTemplateRowsHandler,
		"hanging-punctuation":        HangingPunctuationHandler,
		"height":                     HeightHandler,
		"hyphens":                    HyphensHandler,
		"image-rendering":            ImageRenderingHandler,
		"isolation":                  IsolationHandler,
		"justify-content":            JustifyContentHandler,
		"left":                       SideHandler,
		"letter-spacing":             LetterSpacingHandler,
		"line-break":                 LineBreakHandler,
		"line-height":                LineHeightHandler,
		"list-style":                 ListStyleHandler,
		"list-style-image":           ImageHandler,
		"list-style-position":        ListStylePositionHandler,
		"list-style-type":            ListStyleTypeHandler,
		"margin":                     MarginHandler,
		"margin-bottom":              MarginSideHandler,
		"margin-left":                MarginSideHandler,
		"margin-right":               MarginSideHandler,
		"margin-top":                 MarginSideHandler,
		"max-height":                 MaxHeightWidthHandler,
		"max-width":                  MaxHeightWidthHandler,
		"min-height":                 MinHeightWidthHandler,
		"min-width":                  MinHeightWidthHandler,
		"mix-blend-mode":             MixBlendModeHandler,
		"object-fit":                 ObjectFitHandler,
		"object-position":            ObjectPositionHandler,
		"opacity":                    OpacityHandler,
		"order":                      OrderHandler,
		"orphans":                    OrphansHandler,
		"outline":                    OutlineHandler,
		"outline-color":              ColorHandler,
		"outline-offset":             OutlineOffsetHandler,
		"outline-style":              OutlineStyleHandler,
		"outline-width":              OutlineWidthHandler,
		"overflow":                   OverflowHandler,
		"overflow-wrap":              OverflowWrapHandler,
		"overflow-x":                 OverflowXYHandler,
		"overflow-y":                 OverflowXYHandler,
		"padding":                    PaddingHandler,
		"padding-bottom":             PaddingSideHandler,
		"padding-left":               PaddingSideHandler,
		"padding-right":              PaddingSideHandler,
		"padding-top":                PaddingSideHandler,
		"page-break-after":           PageBreakBeforeAfterHandler,
		"page-break-before":          PageBreakBeforeAfterHandler,
		"page-break-inside":          PageBreakInsideHandler,
		"perspective":                PerspectiveHandler,
		"perspective-origin":         PerspectiveOriginHandler,
		"pointer-events":             PointerEventsHandler,
		"position":                   PositionHandler,
		"quotes":                     QuotesHandler,
		"resize":                     ResizeHandler,
		"right":                      SideHandler,
		"scroll-behavior":            ScrollBehaviorHandler,
		"tab-size":                   TabSizeHandler,
		"table-layout":               TableLayoutHandler,
		"text-align":                 TextAlignHandler,
		"text-align-last":            TextAlignLastHandler,
		"text-combine-upright":       TextCombineUprightHandler,
		"text-decoration":            TextDecorationHandler,
		"text-decoration-color":      ColorHandler,
		"text-decoration-line":       TextDecorationLineHandler,
		"text-decoration-style":      TextDecorationStyleHandler,
		"text-indent":                TextIndentHandler,
		"text-justify":               TextJustifyHandler,
		"text-orientation":           TextOrientationHandler,
		"text-overflow":              TextOverflowHandler,
		"text-shadow":                TextShadowHandler,
		"text-transform":             TextTransformHandler,
		"top":                        SideHandler,
		"transform":                  TransformHandler,
		"transform-origin":           TransformOriginHandler,
		"transform-style":            TransformStyleHandler,
		"transition":                 TransitionHandler,
		"transition-delay":           TransitionDelayHandler,
		"transition-duration":        TransitionDurationHandler,
		"transition-property":        TransitionPropertyHandler,
		"transition-timing-function": TimingFunctionHandler,
		"unicode-bidi":               UnicodeBidiHandler,
		"user-select":                UserSelectHandler,
		"vertical-align":             VerticalAlignHandler,
		"visibility":                 VisiblityHandler,
		"white-space":                WhiteSpaceHandler,
		"widows":                     OrphansHandler,
		"width":                      WidthHandler,
		"word-break":                 WordBreakHandler,
		"word-spacing":               WordSpacingHandler,
		"word-wrap":                  WordWrapHandler,
		"writing-mode":               WritingModeHandler,
		"z-index":                    ZIndexHandler,
	}
	colorValues = []string{"initial", "inherit", "aliceblue", "antiquewhite",
		"aqua", "aquamarine", "azure", "beige", "bisque", "black",
		"blanchedalmond", "blue", "blueviolet", "brown", "burlywood",
		"cadetblue", "chartreuse", "chocolate", "coral", "cornflowerblue",
		"cornsilk", "crimson", "cyan", "darkblue", "darkcyan", "darkgoldenrod",
		"darkgray", "darkgrey", "darkgreen", "darkkhaki", "darkmagenta",
		"darkolivegreen", "darkorange", "darkorchid", "darkred", "darksalmon",
		"darkseagreen", "darkslateblue", "darkslategrey", "darkslategray",
		"darkturquoise", "darkviolet", "deeppink", "deepskyblue", "dimgray",
		"dimgrey", "dodgerblue", "firebrick", "floralwhite", "forestgreen",
		"fuchsia", "gainsboro", "ghostwhite", "gold", "goldenrod", "gray",
		"grey", "green", "greenyellow", "honeydew", "hotpink", "indianred",
		"indigo", "ivory", "khaki", "lavender", "lavenderblush",
		"lemonchiffon", "lightblue", "lightcoral", "lightcyan",
		"lightgoldenrodyellow", "lightgray", "lightgrey", "lightgreen",
		"lightpink", "lightsalmon", "lightseagreen", "lightskyblue",
		"lightslategray", "lightslategrey", "lightsteeelblue", "lightyellow",
		"lime", "limegreen", "linen", "magenta", "maroon", "mediumaquamarine",
		"mediumblue", "mediumorchid", "mediumpurple", "mediumseagreen",
		"mediumslateblue", "mediumspringgreen", "mediumturquoise",
		"mediumvioletred", "midnightblue", "mintcream", "mistyrose",
		"moccasin", "navajowhite", "navy", "oldlace", "olive", "olivedrab",
		"orange", "orangered", "orchid", "palegoldenrod", "palegreen",
		"paleturquoise", "palevioletred", "papayawhip", "peachpuff", "peru",
		"pink", "plum", "powderblue", "purple", "rebeccapurple", "red",
		"rosybrown", "royalblue", "saddlebrown", "salmon", "sandybrown",
		"seagreen", "seashell", "sienna", "silver", "skyblue", "slateblue",
		"slategray", "slategrey", "snow", "springgreen", "steelblue", "tan",
		"teal", "thistle", "tomato", "turquoise", "violet", "wheat", "white",
		"whitesmoke", "yellow", "yellowgreen"}

	Alpha             = regexp.MustCompile(`^[a-z]+$`)
	Blur              = regexp.MustCompile(`^blur\([0-9]+px\)$`)
	BrightnessCont    = regexp.MustCompile(`^(brightness|contrast)\([0-9]+\%\)$`)
	Count             = regexp.MustCompile(`^[0-9]+[\.]?[0-9]*$`)
	CubicBezier       = regexp.MustCompile(`^cubic-bezier\(([ ]*(0(.[0-9]+)?|1(.0)?),){3}[ ]*(0(.[0-9]+)?|1)\)$`)
	Digits            = regexp.MustCompile(`^digits [2-4]$`)
	DropShadow        = regexp.MustCompile(`drop-shadow\(([-]?[0-9]+px) ([-]?[0-9]+px)( [-]?[0-9]+px)?( ([-]?[0-9]+px))?`)
	Font              = regexp.MustCompile(`^('[a-z \-]+'|[a-z \-]+)$`)
	Grayscale         = regexp.MustCompile(`^grayscale\(([0-9]{1,2}|100)%\)$`)
	GridTemplateAreas = regexp.MustCompile(`^['"]?[a-z ]+['"]?$`)
	HexRGB            = regexp.MustCompile(`^#([0-9a-f]{3}|[0-9a-f]{6}|[0-9a-f]{8})$`)
	HSL               = regexp.MustCompile(`^hsl\([ ]*([012]?[0-9]{1,2}|3[0-5][0-9]|360),[ ]*([0-9]{0,2}|100)\%,[ ]*([0-9]{0,2}|100)\%\)$`)
	HSLA              = regexp.MustCompile(`^hsla\(([ ]*[012]?[0-9]{1,2}|3[0-5][0-9]|360),[ ]*([0-9]{0,2}|100)\%,[ ]*([0-9]{0,2}|100)\%,[ ]*(1|1\.0|0|(0\.[0-9]+))\)$`)
	HueRotate         = regexp.MustCompile(`^hue-rotate\(([12]?[0-9]{1,2}|3[0-5][0-9]|360)?\)$`)
	Invert            = regexp.MustCompile(`^invert\(([0-9]{1,2}|100)%\)$`)
	Length            = regexp.MustCompile(`^[\-]?([0-9]+|[0-9]*[\.][0-9]+)(%|cm|mm|in|px|pt|pc|em|ex|ch|rem|vw|vh|vmin|vmax|deg|rad|turn)?$`)
	Matrix            = regexp.MustCompile(`^matrix\(([ ]*[0-9]+[\.]?[0-9]*,){5}([ ]*[0-9]+[\.]?[0-9]*)\)$`)
	Matrix3D          = regexp.MustCompile(`^matrix3d\(([ ]*[0-9]+[\.]?[0-9]*,){15}([ ]*[0-9]+[\.]?[0-9]*)\)$`)
	NegTime           = regexp.MustCompile(`^[\-]?[0-9]+[\.]?[0-9]*(s|ms)?$`)
	Numeric           = regexp.MustCompile(`^[0-9]+$`)
	NumericDecimal    = regexp.MustCompile(`^[0-9\.]+$`)
	Opactiy           = regexp.MustCompile(`^opacity\(([0-9]{1,2}|100)%\)$`)
	Perspective       = regexp.MustCompile(`perspective\(`)
	Position          = regexp.MustCompile(`^[\-]*[0-9]+[cm|mm|in|px|pt|pc\%]* [[\-]*[0-9]+[cm|mm|in|px|pt|pc\%]*]*$`)
	Opacity           = regexp.MustCompile(`^(0[.]?[0-9]*)|(1.0)$`)
	QuotedAlpha       = regexp.MustCompile(`^["'][a-z]+["']$`)
	Quotes            = regexp.MustCompile(`^([ ]*["'][\x{0022}\x{0027}\x{2039}\x{2039}\x{203A}\x{00AB}\x{00BB}\x{2018}\x{2019}\x{201C}-\x{201E}]["'] ["'][\x{0022}\x{0027}\x{2039}\x{2039}\x{203A}\x{00AB}\x{00BB}\x{2018}\x{2019}\x{201C}-\x{201E}]["'])+$`)
	Rect              = regexp.MustCompile(`^rect\([0-9]+px,[ ]*[0-9]+px,[ ]*[0-9]+px,[ ]*[0-9]+px\)$`)
	RGB               = regexp.MustCompile(`^rgb\(([ ]*((([0-9]{1,2}|100)\%)|(([01]?[0-9]{1,2})|(2[0-4][0-9])|(25[0-5]))),){2}([ ]*((([0-9]{1,2}|100)\%)|(([01]?[0-9]{1,2})|(2[0-4][0-9])|(25[0-5]))))\)$`)
	RGBA              = regexp.MustCompile(`^rgba\(([ ]*((([0-9]{1,2}|100)\%)|(([01]?[0-9]{1,2})|(2[0-4][0-9])|(25[0-5]))),){3}[ ]*(1(\.0)?|0|(0\.[0-9]+))\)$`)
	Rotate            = regexp.MustCompile(`^rotate(x|y|z)?\(([12]?|3[0-5][0-9]|360)\)$`)
	Rotate3D          = regexp.MustCompile(`^rotate3d\(([ ]?(1(\.0)?|0\.[0-9]+),){3}([12]?|3[0-5][0-9]|360)\)$`)
	Saturate          = regexp.MustCompile(`^saturate\([0-9]+%\)$`)
	Sepia             = regexp.MustCompile(`^sepia\(([0-9]{1,2}|100)%\)$`)
	Skew              = regexp.MustCompile(`skew(x|y)?\(`)
	Span              = regexp.MustCompile(`^span [0-9]+$`)
	Steps             = regexp.MustCompile(`^steps\([ ]*[0-9]+([ ]*,[ ]*(start|end)?)\)$`)
	Time              = regexp.MustCompile(`^[0-9]+[\.]?[0-9]*(s|ms)?$`)
	TransitionProp    = regexp.MustCompile(`^([a-zA-Z]+,[ ]?)*[a-zA-Z]+$`)
	TranslateScale    = regexp.MustCompile(`(translate|translate3d|translatex|translatey|translatez|scale|scale3d|scalex|scaley|scalez)\(`)
	URL               = regexp.MustCompile(`^url\([\"\']?((https|http)[a-z0-9\.\\/_:]+[\"\']?)\)$`)
	ZIndex            = regexp.MustCompile(`^[\-]?[0-9]+$`)
)

func multiSplit(value string, seps ...string) []string {
	curArray := []string{value}
	for _, i := range seps {
		newArray := []string{}
		for _, j := range curArray {
			newArray = append(newArray, strings.Split(j, i)...)
		}
		curArray = newArray
	}
	return curArray
}

func recursiveCheck(value []string, funcs []func(string) bool) bool {
	for i := 0; i < len(value); i++ {
		tempVal := strings.Join(value[:i+1], " ")
		for _, j := range funcs {
			if j(tempVal) && (len(value[i+1:]) == 0 || recursiveCheck(value[i+1:], funcs)) {
				return true
			}
		}
	}
	return false
}

func in(value []string, arr []string) bool {
	for _, i := range value {
		foundString := false
		for _, j := range arr {
			if j == i {
				foundString = true
			}
		}
		if !foundString {
			return false
		}
	}
	return true
}

func splitValues(value string) []string {
	values := strings.Split(value, ",")
	for _, strippedValue := range values {
		strippedValue = strings.ToLower(strings.TrimSpace(strippedValue))
	}
	return values
}

func GetDefaultHandler(attr string) func(string) bool {

	if defaultStyleHandlers[attr] != nil {
		return defaultStyleHandlers[attr]
	}
	return BaseHandler
}

func BaseHandler(value string) bool {
	return false
}

func AlignContentHandler(value string) bool {
	values := []string{"stretch", "center", "flex-start",
		"flex-end", "space-between", "space-around", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AlignItemsHandler(value string) bool {
	values := []string{"stretch", "center", "flex-start",
		"flex-end", "baseline", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AlignSelfHandler(value string) bool {
	values := []string{"auto", "stretch", "center", "flex-start",
		"flex-end", "baseline", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AllHandler(value string) bool {
	values := []string{"initial", "inherit", "unset"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AnimationHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		AnimationNameHandler,
		AnimationDurationHandler,
		TimingFunctionHandler,
		AnimationDelayHandler,
		AnimationIterationCountHandler,
		AnimationDirectionHandler,
		AnimationFillModeHandler,
		AnimationPlayStateHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func AnimationDelayHandler(value string) bool {
	if NegTime.MatchString(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AnimationDirectionHandler(value string) bool {
	values := []string{"normal", "reverse", "alternate", "alternate-reverse", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AnimationDurationHandler(value string) bool {
	if Time.MatchString(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AnimationFillModeHandler(value string) bool {
	values := []string{"none", "forwards", "backwards", "both", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AnimationIterationCountHandler(value string) bool {
	if Count.MatchString(value) {
		return true
	}
	values := []string{"infinite", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func AnimationNameHandler(value string) bool {
	return Alpha.MatchString(value)
}

func AnimationPlayStateHandler(value string) bool {
	values := []string{"paused", "running", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TimingFunctionHandler(value string) bool {
	values := []string{"linear", "ease", "ease-in", "ease-out", "ease-in-out", "step-start", "step-end", "initial", "inherit"}
	splitVals := splitValues(value)
	if in(splitVals, values) {
		return true
	}
	if CubicBezier.MatchString(value) {
		return true
	}
	return Steps.MatchString(value)
}

func BackfaceVisibilityHandler(value string) bool {
	values := []string{"visible", "hidden", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BackgroundHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	newSplitVals := []string{}
	for _, i := range splitVals {
		if len(strings.Split(i, "/")) == 2 {
			newSplitVals = append(newSplitVals, strings.Split(i, "/")...)
		} else {
			newSplitVals = append(newSplitVals, i)
		}
	}
	usedFunctions := []func(string) bool{
		ColorHandler,
		ImageHandler,
		BackgroundPositionHandler,
		BackgroundSizeHandler,
		BackgroundRepeatHandler,
		BackgroundOriginHandler,
		BackgroundClipHandler,
		BackgroundAttachmentHandler,
	}
	return recursiveCheck(newSplitVals, usedFunctions)
}

func BackgroundAttachmentHandler(value string) bool {
	values := []string{"scroll", "fixed", "local", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BackgroundClipHandler(value string) bool {
	values := []string{"border-box", "padding-box", "content-box", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BackgroundBlendModeHandler(value string) bool {
	values := []string{"normal", "multiply", "screen", "overlay", "darken",
		"lighten", "color-dodge", "saturation", "color", "luminosity"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ImageHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	splitVals := splitValues(value)
	if in(splitVals, values) {
		return true
	}
	return URL.MatchString(value)
}

func BackgroundOriginHandler(value string) bool {
	values := []string{"padding-box", "border-box", "content-box", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BackgroundPositionHandler(value string) bool {
	splitVals := strings.Split(value, ";")
	values := []string{"left", "left top", "left bottom", "right", "right top", "right bottom", "right center", "center top", "center center", "center bottom", "center", "top", "bottom", "initial", "inherit"}
	if in(splitVals, values) {
		return true
	}
	return Position.MatchString(value)
}

func BackgroundRepeatHandler(value string) bool {
	values := []string{"repeat", "repeat-x", "repeat-y", "no-repeat", "space", "round", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BackgroundSizeHandler(value string) bool {
	splitVals := strings.Split(value, " ")
	values := []string{"auto", "cover", "contain", "initial", "inherit"}
	if in(splitVals, values) {
		return true
	}
	if len(splitVals) > 0 && LengthHandler(splitVals[0]) {
		if len(splitVals) < 2 || (len(splitVals) == 2 && LengthHandler(splitVals[1])) {
			return true
		}
	}
	return false
}

func BorderHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := multiSplit(value, " ", "/")
	usedFunctions := []func(string) bool{
		BorderWidthHandler,
		BorderStyleHandler,
		ColorHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func BorderSideHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		BorderSideWidthHandler,
		BorderSideStyleHandler,
		ColorHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func BorderSideRadiusHandler(value string) bool {
	splitVals := strings.Split(value, " ")
	valid := true
	for _, i := range splitVals {
		if !LengthHandler(i) {
			valid = false
			break
		}
	}
	if valid {
		return true
	}
	splitVals = splitValues(value)
	values := []string{"initial", "inherit"}
	return in(splitVals, values)
}

func BorderSideStyleHandler(value string) bool {
	values := []string{"none", "hidden", "dotted", "dashed", "solid", "double", "groove", "ridge", "inset", "outset", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BorderSideWidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	splitVals := strings.Split(value, ";")
	values := []string{"medium", "thin", "thick", "initial", "inherit"}
	return in(splitVals, values)
}

func BorderCollapseHandler(value string) bool {
	values := []string{"separate", "collapse", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BorderImageHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := multiSplit(value, " ", " / ")
	usedFunctions := []func(string) bool{
		ImageHandler,
		BorderImageSliceHandler,
		BorderImageWidthHandler,
		BorderImageOutsetHandler,
		BorderImageRepeatHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func BorderImageOutsetHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BorderImageRepeatHandler(value string) bool {
	values := []string{"stretch", "repeat", "round", "space", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BorderImageSliceHandler(value string) bool {
	values := []string{"fill", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 4 {
		return false
	}
	usedFunctions := []func(string) bool{
		LengthHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func BorderImageWidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BorderRadiusHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 4 {
		return false
	}
	usedFunctions := []func(string) bool{
		LengthHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func BorderSpacingHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 2 {
		return false
	}
	usedFunctions := []func(string) bool{
		LengthHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func BorderStyleHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 4 {
		return false
	}
	usedFunctions := []func(string) bool{
		BorderSideStyleHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func BorderWidthHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 4 {
		return false
	}
	usedFunctions := []func(string) bool{
		BorderSideWidthHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func SideHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"auto", "inherit", "unset"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BoxDecorationBreakHandler(value string) bool {
	values := []string{"slice", "clone", "initial", "initial", "inherit", "unset"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BoxShadowHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	commaSplitVals := strings.Split(value, ",")
	for _, val := range commaSplitVals {
		splitVals := strings.Split(val, " ")
		if len(splitVals) > 6 || len(splitVals) < 2 {
			return false
		}
		if !LengthHandler(splitVals[0]) {
			return false
		}
		if !LengthHandler(splitVals[1]) {
			return false
		}
		usedFunctions := []func(string) bool{
			LengthHandler,
			ColorHandler,
		}
		if len(splitVals) > 2 && !recursiveCheck(splitVals[2:], usedFunctions) {
			return false
		}
	}
	return true
}

func BoxSizingHandler(value string) bool {
	values := []string{"slicontent-box", "border-box", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BreakBeforeAfterHandler(value string) bool {
	values := []string{"auto", "avoid", "always", "all", "avoid-page", "page", "left", "right", "recto", "verso", "avoid-column", "column", "avoid-region", "region"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func BreakInsideHandler(value string) bool {
	values := []string{"auto", "avoid", "avoid-page", "avoid-column", "avoid-region"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func CaptionSideHandler(value string) bool {
	values := []string{"top", "bottom", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func CaretColorHandler(value string) bool {
	splitVals := splitValues(value)
	if in(splitVals, colorValues) {
		return true
	}
	if HexRGB.MatchString(value) {
		return true
	}
	if RGB.MatchString(value) {
		return true
	}
	if RGBA.MatchString(value) {
		return true
	}
	if HSL.MatchString(value) {
		return true
	}
	return HSLA.MatchString(value)
}

func ClearHandler(value string) bool {
	values := []string{"none", "left", "right", "both", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ClipHandler(value string) bool {
	if Rect.MatchString(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ColorHandler(value string) bool {
	splitVals := splitValues(value)
	if in(splitVals, colorValues) {
		return true
	}
	if HexRGB.MatchString(value) {
		return true
	}
	if RGB.MatchString(value) {
		return true
	}
	if RGBA.MatchString(value) {
		return true
	}
	if HSL.MatchString(value) {
		return true
	}
	return HSLA.MatchString(value)
}

func ColumnCountHandler(value string) bool {
	if Numeric.MatchString(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ColumnFillHandler(value string) bool {
	values := []string{"balance", "auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ColumnGapHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"normal", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ColumnRuleHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		ColumnRuleWidthHandler,
		BorderSideStyleHandler,
		ColorHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func ColumnRuleWidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	splitVals := strings.Split(value, ";")
	values := []string{"medium", "thin", "thick", "initial", "inherit"}
	return in(splitVals, values)
}

func ColumnSpanHandler(value string) bool {
	values := []string{"none", "all", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ColumnWidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	splitVals := strings.Split(value, ";")
	values := []string{"auto", "initial", "inherit"}
	return in(splitVals, values)
}

func ColumnsHandler(value string) bool {
	values := []string{"auto", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		ColumnWidthHandler,
		ColumnCountHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func CursorHandler(value string) bool {
	values := []string{"alias", "all-scroll", "auto", "cell", "context-menu", "col-resize", "copy", "crosshair", "default", "e-resize", "ew-resize", "grab", "grabbing", "help", "move", "n-resize", "ne-resize", "nesw-resize", "ns-resize", "nw-resize", "nwse-resize", "no-drop", "none", "not-allowed", "pointer", "progress", "row-resize", "s-resize", "se-resize", "sw-resize", "text", "vertical-text", "w-resize", "wait", "zoom-in", "zoom-out", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func DirectionHandler(value string) bool {
	values := []string{"ltr", "rtl", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func DisplayHandler(value string) bool {
	values := []string{"inline", "block", "contents", "flex", "grid", "inline-block", "inline-flex", "inline-grid", "inline-table", "list-item", "run-in", "table", "table-caption", "table-column-group", "table-header-group", "table-footer-group", "table-row-group", "table-cell", "table-column", "table-row", "none", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func EmptyCellsHandler(value string) bool {
	values := []string{"show", "hide", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FilterHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	splitVals := splitValues(value)
	if in(splitVals, values) {
		return true
	}
	if Blur.MatchString(value) {
		return true
	}
	if BrightnessCont.MatchString(value) {
		return true
	}
	if DropShadow.MatchString(value) {
		return true
	}
	colorValue := strings.TrimSuffix(string(DropShadow.ReplaceAll([]byte(value), []byte{})), ")")
	if ColorHandler(colorValue) {
		return true
	}
	if Grayscale.MatchString(value) {
		return true
	}
	if HueRotate.MatchString(value) {
		return true
	}
	if Invert.MatchString(value) {
		return true
	}
	if Opacity.MatchString(value) {
		return true
	}
	if Saturate.MatchString(value) {
		return true
	}
	return Sepia.MatchString(value)
}

func FlexHandler(value string) bool {
	values := []string{"auto", "initial", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		FlexGrowHandler,
		FlexBasisHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func FlexBasisHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	splitVals := strings.Split(value, ";")
	values := []string{"auto", "initial", "inherit"}
	return in(splitVals, values)
}

func FlexDirectionHandler(value string) bool {
	values := []string{"row", "row-reverse", "column", "column-reverse", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FlexFlowHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		FlexDirectionHandler,
		FlexWrapHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func FlexGrowHandler(value string) bool {
	if NumericDecimal.MatchString(value) {
		return true
	}
	splitVals := strings.Split(value, ";")
	values := []string{"initial", "inherit"}
	return in(splitVals, values)
}

func FlexWrapHandler(value string) bool {
	values := []string{"nowrap", "wrap", "wrap-reverse", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FloatHandler(value string) bool {
	values := []string{"none", "left", "right", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontHandler(value string) bool {
	values := []string{"caption", "icon", "menu", "message-box", "small-caption", "status-bar", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	newSplitVals := []string{}
	for _, i := range splitVals {
		if len(strings.Split(i, "/")) == 2 {
			newSplitVals = append(newSplitVals, strings.Split(i, "/")...)
		} else {
			newSplitVals = append(newSplitVals, i)
		}
	}
	usedFunctions := []func(string) bool{
		FontStyleHandler,
		FontVariantHandler,
		FontWeightHandler,
		FontSizeHandler,
		FontFamilyHandler,
	}
	return recursiveCheck(newSplitVals, usedFunctions)
}

func FontFamilyHandler(value string) bool {
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	if in(splitVals, values) {
		return true
	}
	for _, i := range splitVals {
		i = strings.TrimSpace(i)
		if Font.FindString(i) != i {
			return false
		}
	}
	return true
}

func FontKerningHandler(value string) bool {
	values := []string{"auto", "normal", "none"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontLanguageOverrideHandler(value string) bool {
	return Alpha.MatchString(value)
}

func FontSizeHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"medium", "xx-small", "x-small", "small", "large", "x-large", "xx-large", "smaller", "larger", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontSizeAdjustHandler(value string) bool {
	if Count.MatchString(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontStretchHandler(value string) bool {
	values := []string{"ultra-condensed", "extra-condensed", "condensed", "semi-condensed", "normal", "semi-expanded", "expanded", "extra-expanded", "ultra-expanded", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontStyleHandler(value string) bool {
	values := []string{"normal", "italic", "oblique", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontSynthesisHandler(value string) bool {
	values := []string{"none", "style", "weight"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontVariantCapsHandler(value string) bool {
	values := []string{"normal", "small-caps", "all-small-caps", "petite-caps", "all-petite-caps", "unicase", "titling-caps"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontVariantHandler(value string) bool {
	values := []string{"normal", "small-caps", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontVariantPositionHandler(value string) bool {
	values := []string{"normal", "sub", "super"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func FontWeightHandler(value string) bool {
	values := []string{"normal", "bold", "bolder", "lighter", "100", "200", "300", "400", "500", "600", "700", "800", "900", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func GridHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	newSplitVals := []string{}
	for _, i := range splitVals {
		if i != "/" {
			newSplitVals = append(newSplitVals, i)
		}
	}
	usedFunctions := []func(string) bool{
		GridTemplateRowsHandler,
		GridTemplateColumnsHandler,
		GridTemplateAreasHandler,
		GridAutoColumnsHandler,
		GridAutoFlowHandler,
	}
	return recursiveCheck(newSplitVals, usedFunctions)
}

func GridAreaHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " / ")
	usedFunctions := []func(string) bool{
		GridAxisStartEndHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func GridAutoColumnsHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"auto", "max-content", "min-content", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func GridAutoFlowHandler(value string) bool {
	values := []string{"row", "column", "dense", "row dense", "column dense"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func GridColumnHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " / ")
	if len(splitVals) > 2 {
		return false
	}
	usedFunctions := []func(string) bool{
		GridAxisStartEndHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func GridColumnGapHandler(value string) bool {
	return LengthHandler(value)
}

func LengthHandler(value string) bool {
	return Length.MatchString(value)
}

func LineBreakHandler(value string) bool {
	values := []string{"auto", "loose", "normal", "strict"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func GridAxisStartEndHandler(value string) bool {
	if Numeric.MatchString(value) {
		return true
	}
	if Span.MatchString(value) {
		return true
	}
	values := []string{"auto"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func GridGapHandler(value string) bool {
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 2 {
		return false
	}
	usedFunctions := []func(string) bool{
		GridColumnGapHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func GridRowHandler(value string) bool {
	splitVals := strings.Split(value, " / ")
	if len(splitVals) > 2 {
		return false
	}
	usedFunctions := []func(string) bool{
		GridAxisStartEndHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func GridTemplateHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " / ")
	if len(splitVals) > 2 {
		return false
	}
	usedFunctions := []func(string) bool{
		GridTemplateColumnsHandler,
		GridTemplateRowsHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func GridTemplateAreasHandler(value string) bool {
	values := []string{"none"}
	if in([]string{value}, values) {
		return true
	}
	return GridTemplateAreas.MatchString(value)
}

func GridTemplateColumnsHandler(value string) bool {
	splitVals := strings.Split(value, " ")
	values := []string{"none", "auto", "max-content", "min-content", "initial", "inherit"}
	for _, val := range splitVals {
		if LengthHandler(val) {
			continue
		}
		valArr := []string{val}
		if !in(valArr, values) {
			return false
		}
	}
	return true
}

func GridTemplateRowsHandler(value string) bool {
	splitVals := strings.Split(value, " ")
	values := []string{"none", "auto", "max-content", "min-content"}
	for _, val := range splitVals {
		if LengthHandler(val) {
			continue
		}
		valArr := []string{val}
		if !in(valArr, values) {
			return false
		}
	}
	return true
}

func HangingPunctuationHandler(value string) bool {
	values := []string{"none", "first", "last", "allow-end", "force-end", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func HeightHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func HyphensHandler(value string) bool {
	values := []string{"none", "manual", "auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ImageRenderingHandler(value string) bool {
	values := []string{"auto", "smooth", "high-quality", "crisp-edges", "pixelated"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func IsolationHandler(value string) bool {
	values := []string{"auto", "isolate", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func JustifyContentHandler(value string) bool {
	values := []string{"flex-start", "flex-end", "center", "space-between", "space-around", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func LetterSpacingHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"normal", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func LineHeightHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"normal", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ListStyleHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		ListStyleTypeHandler,
		ListStylePositionHandler,
		ImageHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func ListStylePositionHandler(value string) bool {
	values := []string{"inside", "outside", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ListStyleTypeHandler(value string) bool {
	values := []string{"disc", "armenian", "circle", "cjk-ideographic", "decimal", "decimal-leading-zero", "georgian", "hebrew", "hiragana", "hiragana-iroha", "katakana", "katakana-iroha", "lower-alpha", "lower-greek", "lower-latin", "lower-roman", "none", "square", "upper-alpha", "upper-greek", "upper-latin", "upper-roman", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func MarginHandler(value string) bool {
	values := []string{"auto", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		MarginSideHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func MarginSideHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func MaxHeightWidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"none", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func MinHeightWidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func MixBlendModeHandler(value string) bool {
	values := []string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "difference", "exclusion", "hue", "saturation", "color", "luminosity"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ObjectFitHandler(value string) bool {
	values := []string{"fill", "contain", "cover", "none", "scale-down", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ObjectPositionHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 2 {
		return false
	}
	usedFunctions := []func(string) bool{
		LengthHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func OpacityHandler(value string) bool {
	if Opacity.MatchString(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OrderHandler(value string) bool {
	if Numeric.MatchString(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OutlineHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		ColorHandler,
		OutlineWidthHandler,
		OutlineStyleHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func OutlineOffsetHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OutlineStyleHandler(value string) bool {
	values := []string{"none", "hidden", "dotted", "dashed", "solid", "double", "groove", "ridge", "inset", "outset", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OutlineWidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"medium", "thin", "thick", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OverflowHandler(value string) bool {
	values := []string{"visible", "hidden", "scroll", "auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OverflowXYHandler(value string) bool {
	values := []string{"visible", "hidden", "scroll", "auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OverflowWrapHandler(value string) bool {
	values := []string{"normal", "break-word", "anywhere"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func OrphansHandler(value string) bool {
	return Numeric.MatchString(value)
}

func PaddingHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	if len(splitVals) > 4 {
		return false
	}
	usedFunctions := []func(string) bool{
		PaddingSideHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func PaddingSideHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func PageBreakBeforeAfterHandler(value string) bool {
	values := []string{"auto", "always", "avoid", "left", "right", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func PageBreakInsideHandler(value string) bool {
	values := []string{"auto", "avoid", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func PerspectiveHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"none", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func PerspectiveOriginHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	xValues := []string{"left", "center", "right"}
	yValues := []string{"top", "center", "bottom"}
	if len(splitVals) > 1 {
		if !in([]string{splitVals[0]}, xValues) && !LengthHandler(splitVals[0]) {
			return false
		}
		return in([]string{splitVals[1]}, yValues) || LengthHandler(splitVals[1])
	} else if len(splitVals) == 1 {
		return in(splitVals, xValues) || in(splitVals, yValues) || LengthHandler(splitVals[0])
	}
	return false
}

func PointerEventsHandler(value string) bool {
	values := []string{"auto", "none", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func PositionHandler(value string) bool {
	values := []string{"static", "absolute", "fixed", "relative", "sticky", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func QuotesHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	splitVals := splitValues(value)
	if in(splitVals, values) {
		return true
	}
	return Quotes.MatchString(value)
}

func ResizeHandler(value string) bool {
	values := []string{"none", "both", "horizontal", "vertical", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ScrollBehaviorHandler(value string) bool {
	values := []string{"auto", "smooth", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TabSizeHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TableLayoutHandler(value string) bool {
	values := []string{"auto", "fixed", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextAlignHandler(value string) bool {
	values := []string{"left", "right", "center", "justify", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextAlignLastHandler(value string) bool {
	values := []string{"auto", "left", "right", "center", "justify", "start", "end", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextCombineUprightHandler(value string) bool {
	values := []string{"none", "all"}
	splitVals := splitValues(value)
	if in(splitVals, values) {
		return true
	}
	return Digits.MatchString(value)
}

func TextDecorationHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		TextDecorationStyleHandler,
		ColorHandler,
		TextDecorationLineHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func TextDecorationLineHandler(value string) bool {
	values := []string{"none", "underline", "overline", "line-through", "initial", "inherit"}
	splitVals := strings.Split(value, " ")
	return in(splitVals, values)
}

func TextDecorationStyleHandler(value string) bool {
	values := []string{"solid", "double", "dotted", "dashed", "wavy", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextIndentHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextJustifyHandler(value string) bool {
	values := []string{"auto", "inter-word", "inter-character", "none", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextOverflowHandler(value string) bool {
	if QuotedAlpha.MatchString(value) {
		return true
	}
	values := []string{"clip", "ellipsis", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextOrientationHandler(value string) bool {
	values := []string{"mixed", "upright", "sideways", "sideways-right"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TextShadowHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	commaSplitVals := strings.Split(value, ",")
	for _, val := range commaSplitVals {
		splitVals := strings.Split(val, " ")
		if len(splitVals) > 6 || len(splitVals) < 2 {
			return false
		}
		if !LengthHandler(splitVals[0]) {
			return false
		}
		if !LengthHandler(splitVals[1]) {
			return false
		}
		usedFunctions := []func(string) bool{
			LengthHandler,
			ColorHandler,
		}
		if len(splitVals) > 2 && !recursiveCheck(splitVals[2:], usedFunctions) {
			return false
		}
	}
	return true
}

func TextTransformHandler(value string) bool {
	values := []string{"none", "capitalize", "uppercase", "lowercase", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TransformHandler(value string) bool {
	values := []string{"none", "initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	if Matrix.MatchString(value) {
		return true
	}
	if Matrix3D.MatchString(value) {
		return true
	}
	subValue := string(TranslateScale.ReplaceAll([]byte(value), []byte{}))
	trimValue := strings.Split(strings.TrimSuffix(subValue, ")"), ",")
	valid := true
	for _, i := range trimValue {
		if !LengthHandler(strings.TrimSpace(i)) {
			valid = false
			break
		}
	}
	if valid && trimValue != nil {
		return true
	}
	if Rotate.MatchString(value) {
		return true
	}
	if Rotate3D.MatchString(value) {
		return true
	}
	subValue = string(Skew.ReplaceAll([]byte(value), []byte{}))
	subValue = strings.TrimSuffix(subValue, ")")
	trimValue = strings.Split(subValue, ",")
	valid = true
	for _, i := range trimValue {
		if !LengthHandler(strings.TrimSpace(i)) {
			valid = false
			break
		}
	}
	if valid {
		return true
	}
	subValue = string(Perspective.ReplaceAll([]byte(value), []byte{}))
	subValue = strings.TrimSuffix(subValue, ")")
	return LengthHandler(subValue)
}

func TransformOriginHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	xValues := []string{"left", "center", "right"}
	yValues := []string{"top", "center", "bottom"}
	if len(splitVals) > 2 {
		if !in([]string{splitVals[0]}, xValues) && !LengthHandler(splitVals[0]) {
			return false
		}
		if !in([]string{splitVals[1]}, yValues) && !LengthHandler(splitVals[1]) {
			return false
		}
		return LengthHandler(splitVals[2])
	} else if len(splitVals) > 1 {
		if !in([]string{splitVals[0]}, xValues) && !LengthHandler(splitVals[0]) {
			return false
		}
		return in([]string{splitVals[1]}, yValues) || LengthHandler(splitVals[1])
	} else if len(splitVals) == 1 {
		return in(splitVals, xValues) || in(splitVals, yValues) || LengthHandler(splitVals[0])
	}
	return false
}

func TransformStyleHandler(value string) bool {
	values := []string{"flat", "preserve-3d", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TransitionHandler(value string) bool {
	values := []string{"initial", "inherit"}
	if in([]string{value}, values) {
		return true
	}
	splitVals := strings.Split(value, " ")
	usedFunctions := []func(string) bool{
		TransitionPropertyHandler,
		TransitionDurationHandler,
		TimingFunctionHandler,
		TransitionDelayHandler,
		ColorHandler,
	}
	return recursiveCheck(splitVals, usedFunctions)
}

func TransitionDelayHandler(value string) bool {
	if Time.MatchString(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TransitionDurationHandler(value string) bool {
	if Time.MatchString(value) {
		return true
	}
	values := []string{"initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func TransitionPropertyHandler(value string) bool {
	if TransitionProp.MatchString(value) {
		return true
	}
	values := []string{"none", "all", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func UnicodeBidiHandler(value string) bool {
	values := []string{"normal", "embed", "bidi-override", "isolate", "isolate-override", "plaintext", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func UserSelectHandler(value string) bool {
	values := []string{"auto", "none", "text", "all"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func VerticalAlignHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"baseline", "sub", "super", "top", "text-top", "middle", "bottom", "text-bottom", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func VisiblityHandler(value string) bool {
	values := []string{"visible", "hidden", "collapse", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func WhiteSpaceHandler(value string) bool {
	values := []string{"normal", "nowrap", "pre", "pre-line", "pre-wrap", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func WidthHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func WordSpacingHandler(value string) bool {
	if LengthHandler(value) {
		return true
	}
	values := []string{"normal", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func WordBreakHandler(value string) bool {
	values := []string{"normal", "break-all", "keep-all", "break-word", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func WordWrapHandler(value string) bool {
	values := []string{"normal", "break-word", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func WritingModeHandler(value string) bool {
	values := []string{"horizontal-tb", "vertical-rl", "vertical-lr"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}

func ZIndexHandler(value string) bool {
	if ZIndex.MatchString(value) {
		return true
	}
	values := []string{"auto", "initial", "inherit"}
	splitVals := splitValues(value)
	return in(splitVals, values)
}
