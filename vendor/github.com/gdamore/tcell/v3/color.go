// Copyright 2025 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import (
	ic "image/color"

	"github.com/gdamore/tcell/v3/color"
)

// Note that the entire contents of this file should be considered deprecated
// in favor of the color subpackage.  (This also means others can use the color
// package without importing the entirety of tcell into their binaries.)

// Color represents a color.  The low numeric values are the same as used
// by ECMA-48, and beyond that XTerm.
//
// Note that on various terminals colors may be approximated however, or
// not supported at all.  If no suitable representation for a color is known,
// the library will simply not set any color, deferring to whatever default
// attributes the terminal uses.
type Color = color.Color

const (
	// ColorDefault is used to leave the Color unchanged from whatever
	// system or terminal default may exist.  It's also the zero value.
	ColorDefault = color.Default

	// ColorValid is used to indicate the color value is actually
	// valid (initialized).
	// Deprecated: Use color.IsValid instead.
	ColorValid = color.IsValid

	// ColorIsRGB is used to indicate that the numeric value is not
	// a known color constant, but rather an RGB value.
	// Deprecated: Use color.IsRGB instead.
	ColorIsRGB = color.IsRGB

	// ColorSpecial is a flag used to indicate that the values have
	// special meaning, and live outside of the color space(s).
	// Deprecated.
	ColorSpecial = color.IsSpecial
)

// Note that the order of these options is important -- it follows the
// definitions used by ECMA and XTerm.  Hence any further named colors
// must begin at a value not less than 256.
//
// Deprecated: Use color.XXX symbols instead.
const (
	ColorBlack                = color.XTerm0
	ColorMaroon               = color.XTerm1
	ColorGreen                = color.XTerm2
	ColorOlive                = color.XTerm3
	ColorNavy                 = color.XTerm4
	ColorPurple               = color.XTerm5
	ColorTeal                 = color.XTerm6
	ColorSilver               = color.XTerm7
	ColorGray                 = color.XTerm8
	ColorRed                  = color.XTerm9
	ColorLime                 = color.XTerm10
	ColorYellow               = color.XTerm11
	ColorBlue                 = color.XTerm12
	ColorFuchsia              = color.XTerm13
	ColorAqua                 = color.XTerm14
	ColorWhite                = color.XTerm15
	Color16                   = color.XTerm16
	Color17                   = color.XTerm17
	Color18                   = color.XTerm18
	Color19                   = color.XTerm19
	Color20                   = color.XTerm20
	Color21                   = color.XTerm21
	Color22                   = color.XTerm22
	Color23                   = color.XTerm23
	Color24                   = color.XTerm24
	Color25                   = color.XTerm25
	Color26                   = color.XTerm26
	Color27                   = color.XTerm27
	Color28                   = color.XTerm28
	Color29                   = color.XTerm29
	Color30                   = color.XTerm30
	Color31                   = color.XTerm31
	Color32                   = color.XTerm32
	Color33                   = color.XTerm33
	Color34                   = color.XTerm34
	Color35                   = color.XTerm35
	Color36                   = color.XTerm36
	Color37                   = color.XTerm37
	Color38                   = color.XTerm38
	Color39                   = color.XTerm39
	Color40                   = color.XTerm40
	Color41                   = color.XTerm41
	Color42                   = color.XTerm42
	Color43                   = color.XTerm43
	Color44                   = color.XTerm44
	Color45                   = color.XTerm45
	Color46                   = color.XTerm46
	Color47                   = color.XTerm47
	Color48                   = color.XTerm48
	Color49                   = color.XTerm49
	Color50                   = color.XTerm50
	Color51                   = color.XTerm51
	Color52                   = color.XTerm52
	Color53                   = color.XTerm53
	Color54                   = color.XTerm54
	Color55                   = color.XTerm55
	Color56                   = color.XTerm56
	Color57                   = color.XTerm57
	Color58                   = color.XTerm58
	Color59                   = color.XTerm59
	Color60                   = color.XTerm60
	Color61                   = color.XTerm61
	Color62                   = color.XTerm62
	Color63                   = color.XTerm63
	Color64                   = color.XTerm64
	Color65                   = color.XTerm65
	Color66                   = color.XTerm66
	Color67                   = color.XTerm67
	Color68                   = color.XTerm68
	Color69                   = color.XTerm69
	Color70                   = color.XTerm70
	Color71                   = color.XTerm71
	Color72                   = color.XTerm72
	Color73                   = color.XTerm73
	Color74                   = color.XTerm74
	Color75                   = color.XTerm75
	Color76                   = color.XTerm76
	Color77                   = color.XTerm77
	Color78                   = color.XTerm78
	Color79                   = color.XTerm79
	Color80                   = color.XTerm80
	Color81                   = color.XTerm81
	Color82                   = color.XTerm82
	Color83                   = color.XTerm83
	Color84                   = color.XTerm84
	Color85                   = color.XTerm85
	Color86                   = color.XTerm86
	Color87                   = color.XTerm87
	Color88                   = color.XTerm88
	Color89                   = color.XTerm89
	Color90                   = color.XTerm90
	Color91                   = color.XTerm91
	Color92                   = color.XTerm92
	Color93                   = color.XTerm93
	Color94                   = color.XTerm94
	Color95                   = color.XTerm95
	Color96                   = color.XTerm96
	Color97                   = color.XTerm97
	Color98                   = color.XTerm98
	Color99                   = color.XTerm99
	Color100                  = color.XTerm100
	Color101                  = color.XTerm101
	Color102                  = color.XTerm102
	Color103                  = color.XTerm103
	Color104                  = color.XTerm104
	Color105                  = color.XTerm105
	Color106                  = color.XTerm106
	Color107                  = color.XTerm107
	Color108                  = color.XTerm108
	Color109                  = color.XTerm109
	Color110                  = color.XTerm110
	Color111                  = color.XTerm111
	Color112                  = color.XTerm112
	Color113                  = color.XTerm113
	Color114                  = color.XTerm114
	Color115                  = color.XTerm115
	Color116                  = color.XTerm116
	Color117                  = color.XTerm117
	Color118                  = color.XTerm118
	Color119                  = color.XTerm119
	Color120                  = color.XTerm120
	Color121                  = color.XTerm121
	Color122                  = color.XTerm122
	Color123                  = color.XTerm123
	Color124                  = color.XTerm124
	Color125                  = color.XTerm125
	Color126                  = color.XTerm126
	Color127                  = color.XTerm127
	Color128                  = color.XTerm128
	Color129                  = color.XTerm129
	Color130                  = color.XTerm130
	Color131                  = color.XTerm131
	Color132                  = color.XTerm132
	Color133                  = color.XTerm133
	Color134                  = color.XTerm134
	Color135                  = color.XTerm135
	Color136                  = color.XTerm136
	Color137                  = color.XTerm137
	Color138                  = color.XTerm138
	Color139                  = color.XTerm139
	Color140                  = color.XTerm140
	Color141                  = color.XTerm141
	Color142                  = color.XTerm142
	Color143                  = color.XTerm143
	Color144                  = color.XTerm144
	Color145                  = color.XTerm145
	Color146                  = color.XTerm146
	Color147                  = color.XTerm147
	Color148                  = color.XTerm148
	Color149                  = color.XTerm149
	Color150                  = color.XTerm150
	Color151                  = color.XTerm151
	Color152                  = color.XTerm152
	Color153                  = color.XTerm153
	Color154                  = color.XTerm154
	Color155                  = color.XTerm155
	Color156                  = color.XTerm156
	Color157                  = color.XTerm157
	Color158                  = color.XTerm158
	Color159                  = color.XTerm159
	Color160                  = color.XTerm160
	Color161                  = color.XTerm161
	Color162                  = color.XTerm162
	Color163                  = color.XTerm163
	Color164                  = color.XTerm164
	Color165                  = color.XTerm165
	Color166                  = color.XTerm166
	Color167                  = color.XTerm167
	Color168                  = color.XTerm168
	Color169                  = color.XTerm169
	Color170                  = color.XTerm170
	Color171                  = color.XTerm171
	Color172                  = color.XTerm172
	Color173                  = color.XTerm173
	Color174                  = color.XTerm174
	Color175                  = color.XTerm175
	Color176                  = color.XTerm176
	Color177                  = color.XTerm177
	Color178                  = color.XTerm178
	Color179                  = color.XTerm179
	Color180                  = color.XTerm180
	Color181                  = color.XTerm181
	Color182                  = color.XTerm182
	Color183                  = color.XTerm183
	Color184                  = color.XTerm184
	Color185                  = color.XTerm185
	Color186                  = color.XTerm186
	Color187                  = color.XTerm187
	Color188                  = color.XTerm188
	Color189                  = color.XTerm189
	Color190                  = color.XTerm190
	Color191                  = color.XTerm191
	Color192                  = color.XTerm192
	Color193                  = color.XTerm193
	Color194                  = color.XTerm194
	Color195                  = color.XTerm195
	Color196                  = color.XTerm196
	Color197                  = color.XTerm197
	Color198                  = color.XTerm198
	Color199                  = color.XTerm199
	Color200                  = color.XTerm200
	Color201                  = color.XTerm201
	Color202                  = color.XTerm202
	Color203                  = color.XTerm203
	Color204                  = color.XTerm204
	Color205                  = color.XTerm205
	Color206                  = color.XTerm206
	Color207                  = color.XTerm207
	Color208                  = color.XTerm208
	Color209                  = color.XTerm209
	Color210                  = color.XTerm210
	Color211                  = color.XTerm211
	Color212                  = color.XTerm212
	Color213                  = color.XTerm213
	Color214                  = color.XTerm214
	Color215                  = color.XTerm215
	Color216                  = color.XTerm216
	Color217                  = color.XTerm217
	Color218                  = color.XTerm218
	Color219                  = color.XTerm219
	Color220                  = color.XTerm220
	Color221                  = color.XTerm221
	Color222                  = color.XTerm222
	Color223                  = color.XTerm223
	Color224                  = color.XTerm224
	Color225                  = color.XTerm225
	Color226                  = color.XTerm226
	Color227                  = color.XTerm227
	Color228                  = color.XTerm228
	Color229                  = color.XTerm229
	Color230                  = color.XTerm230
	Color231                  = color.XTerm231
	Color232                  = color.XTerm232
	Color233                  = color.XTerm233
	Color234                  = color.XTerm234
	Color235                  = color.XTerm235
	Color236                  = color.XTerm236
	Color237                  = color.XTerm237
	Color238                  = color.XTerm238
	Color239                  = color.XTerm239
	Color240                  = color.XTerm240
	Color241                  = color.XTerm241
	Color242                  = color.XTerm242
	Color243                  = color.XTerm243
	Color244                  = color.XTerm244
	Color245                  = color.XTerm245
	Color246                  = color.XTerm246
	Color247                  = color.XTerm247
	Color248                  = color.XTerm248
	Color249                  = color.XTerm249
	Color250                  = color.XTerm250
	Color251                  = color.XTerm251
	Color252                  = color.XTerm252
	Color253                  = color.XTerm253
	Color254                  = color.XTerm254
	Color255                  = color.XTerm255
	ColorAliceBlue            = color.AliceBlue
	ColorAntiqueWhite         = color.AntiqueWhite
	ColorAquaMarine           = color.AquaMarine
	ColorAzure                = color.Azure
	ColorBeige                = color.Beige
	ColorBisque               = color.Bisque
	ColorBlanchedAlmond       = color.BlanchedAlmond
	ColorBlueViolet           = color.BlueViolet
	ColorBrown                = color.Brown
	ColorBurlyWood            = color.BurlyWood
	ColorCadetBlue            = color.CadetBlue
	ColorChartreuse           = color.Chartreuse
	ColorChocolate            = color.Chocolate
	ColorCoral                = color.Coral
	ColorCornflowerBlue       = color.CornflowerBlue
	ColorCornsilk             = color.Cornsilk
	ColorCrimson              = color.Crimson
	ColorDarkBlue             = color.DarkBlue
	ColorDarkCyan             = color.DarkCyan
	ColorDarkGoldenrod        = color.DarkGoldenrod
	ColorDarkGray             = color.DarkGray
	ColorDarkGreen            = color.DarkGreen
	ColorDarkKhaki            = color.DarkKhaki
	ColorDarkMagenta          = color.DarkMagenta
	ColorDarkOliveGreen       = color.DarkOliveGreen
	ColorDarkOrange           = color.DarkOrange
	ColorDarkOrchid           = color.DarkOrchid
	ColorDarkRed              = color.DarkRed
	ColorDarkSalmon           = color.DarkSalmon
	ColorDarkSeaGreen         = color.DarkSeaGreen
	ColorDarkSlateBlue        = color.DarkSlateBlue
	ColorDarkSlateGray        = color.DarkSlateGray
	ColorDarkTurquoise        = color.DarkTurquoise
	ColorDarkViolet           = color.DarkViolet
	ColorDeepPink             = color.DeepPink
	ColorDeepSkyBlue          = color.DeepSkyBlue
	ColorDimGray              = color.DimGray
	ColorDodgerBlue           = color.DodgerBlue
	ColorFireBrick            = color.FireBrick
	ColorFloralWhite          = color.FloralWhite
	ColorForestGreen          = color.ForestGreen
	ColorGainsboro            = color.Gainsboro
	ColorGhostWhite           = color.GhostWhite
	ColorGold                 = color.Gold
	ColorGoldenrod            = color.Goldenrod
	ColorGreenYellow          = color.GreenYellow
	ColorHoneydew             = color.Honeydew
	ColorHotPink              = color.HotPink
	ColorIndianRed            = color.IndianRed
	ColorIndigo               = color.Indigo
	ColorIvory                = color.Ivory
	ColorKhaki                = color.Khaki
	ColorLavender             = color.Lavender
	ColorLavenderBlush        = color.LavenderBlush
	ColorLawnGreen            = color.LawnGreen
	ColorLemonChiffon         = color.LemonChiffon
	ColorLightBlue            = color.LightBlue
	ColorLightCoral           = color.LightCoral
	ColorLightCyan            = color.LightCyan
	ColorLightGoldenrodYellow = color.LightGoldenrodYellow
	ColorLightGray            = color.LightGray
	ColorLightGreen           = color.LightGreen
	ColorLightPink            = color.LightPink
	ColorLightSalmon          = color.LightSalmon
	ColorLightSeaGreen        = color.LightSeaGreen
	ColorLightSkyBlue         = color.LightSkyBlue
	ColorLightSlateGray       = color.LightSlateGray
	ColorLightSteelBlue       = color.LightSteelBlue
	ColorLightYellow          = color.LightYellow
	ColorLimeGreen            = color.LimeGreen
	ColorLinen                = color.Linen
	ColorMediumAquamarine     = color.MediumAquamarine
	ColorMediumBlue           = color.MediumBlue
	ColorMediumOrchid         = color.MediumOrchid
	ColorMediumPurple         = color.MediumPurple
	ColorMediumSeaGreen       = color.MediumSeaGreen
	ColorMediumSlateBlue      = color.MediumSlateBlue
	ColorMediumSpringGreen    = color.MediumSpringGreen
	ColorMediumTurquoise      = color.MediumTurquoise
	ColorMediumVioletRed      = color.MediumVioletRed
	ColorMidnightBlue         = color.MidnightBlue
	ColorMintCream            = color.MintCream
	ColorMistyRose            = color.MistyRose
	ColorMoccasin             = color.Moccasin
	ColorNavajoWhite          = color.NavajoWhite
	ColorOldLace              = color.OldLace
	ColorOliveDrab            = color.OliveDrab
	ColorOrange               = color.Orange
	ColorOrangeRed            = color.OrangeRed
	ColorOrchid               = color.Orchid
	ColorPaleGoldenrod        = color.PaleGoldenrod
	ColorPaleGreen            = color.PaleGreen
	ColorPaleTurquoise        = color.PaleTurquoise
	ColorPaleVioletRed        = color.PaleVioletRed
	ColorPapayaWhip           = color.PapayaWhip
	ColorPeachPuff            = color.PeachPuff
	ColorPeru                 = color.Peru
	ColorPink                 = color.Pink
	ColorPlum                 = color.Plum
	ColorPowderBlue           = color.PowderBlue
	ColorRebeccaPurple        = color.RebeccaPurple
	ColorRosyBrown            = color.RosyBrown
	ColorRoyalBlue            = color.RoyalBlue
	ColorSaddleBrown          = color.SaddleBrown
	ColorSalmon               = color.Salmon
	ColorSandyBrown           = color.SandyBrown
	ColorSeaGreen             = color.SeaGreen
	ColorSeashell             = color.Seashell
	ColorSienna               = color.Sienna
	ColorSkyblue              = color.Skyblue
	ColorSlateBlue            = color.SlateBlue
	ColorSlateGray            = color.SlateGray
	ColorSnow                 = color.Snow
	ColorSpringGreen          = color.SpringGreen
	ColorSteelBlue            = color.SteelBlue
	ColorTan                  = color.Tan
	ColorThistle              = color.Thistle
	ColorTomato               = color.Tomato
	ColorTurquoise            = color.Turquoise
	ColorViolet               = color.Violet
	ColorWheat                = color.Wheat
	ColorWhiteSmoke           = color.WhiteSmoke
	ColorYellowGreen          = color.YellowGreen
)

// These are aliases for the color gray, because some of us spell
// it as grey. Deprecated: Use color values.
const (
	ColorGrey           = color.Gray
	ColorDimGrey        = color.DimGray
	ColorDarkGrey       = color.DarkGray
	ColorDarkSlateGrey  = color.DarkSlateGray
	ColorLightGrey      = color.LightGray
	ColorLightSlateGrey = color.LightSlateGray
	ColorSlateGrey      = color.SlateGray
)

// ColorValues maps color constants to their RGB values.
var ColorValues = color.ColorValues

// Special colors.
const (
	// ColorReset is used to indicate that the color should use the
	// vanilla terminal colors.  (Basically go back to the defaults.)
	// Deprecated: Use color.Reset.
	ColorReset = color.Reset

	// ColorNone indicates that we should not change the color from
	// whatever is already displayed.  This can only be used in limited
	// circumstances.
	// Deprecated: Use color.None.
	ColorNone = color.None
)

// ColorNames holds the written names of colors. Useful to present a list of
// recognized named colors. Deprecated: Use color.Names.
var ColorNames = color.Names

// NewRGBColor returns a new color with the given red, green, and blue values.
// Each value must be represented in the range 0-255.
// Deprecated: Use color.NewRGBColor.
func NewRGBColor(r, g, b int32) Color {
	return color.NewRGBColor(r, g, b)
}

// NewHexColor returns a color using the given 24-bit RGB value.
// Deprecated: Use color.NewHexColor.
func NewHexColor(v int32) Color {
	return color.NewHexColor(v)
}

// GetColor creates a Color from a color name (W3C name). A hex value may
// be supplied as a string in the format "#ffffff".
// Deprecated: Use color.GetColor.
func GetColor(name string) Color {
	return color.GetColor(name)
}

// PaletteColor creates a color based on the palette index.
// Deprecated: Use color.PaletteColor.
func PaletteColor(index int) Color {
	return color.PaletteColor(index)
}

// FromImageColor converts an image/color.Color into Color.
// Deprecated: Use color.FromImageColor.
func FromImageColor(imageColor ic.Color) Color {
	return color.FromImageColor(imageColor)
}

// FindColor attempts to find a given color, or the best match possible for it,
// from the palette given.  This is an expensive operation, so results should
// be cached by the caller.
// Deprecated: Use color.Find.
func FindColor(c Color, palette []Color) Color {
	return color.Find(c, palette)
}
