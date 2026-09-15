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

package color

import (
	"fmt"
	ic "image/color"
	"strconv"
)

// Color represents a color.  The low numeric values are the same as used
// by ECMA-48, and beyond that XTerm.
//
// For Color names we use the W3C approved color names.
//
// Note that on various terminals colors may be approximated however, or
// not supported at all.  If no suitable representation for a color is known,
// the library will simply not set any color, deferring to whatever default
// attributes the terminal uses.
type Color uint32

const (
	// Default is used to leave the Color unchanged from whatever
	// system or terminal default may exist.  It's also the zero value.
	Default Color = 0

	// IsValid is used to indicate the color value is actually
	// valid (initialized).  This is useful to permit the zero value
	// to be treated as the default. This should not be used
	// directly by applications.  Use the IsValid method on Color instead.
	IsValid Color = 1 << 31

	// IsRGB is used to indicate that the numeric value is not
	// a known color constant, but rather an RGB value.  The lower
	// order 3 bytes are RGB.  This should not be used directly,
	// instead use the Color.IsRGB method.
	IsRGB Color = 1 << 30

	// IsSpecial is a flag used to indicate that the values have
	// special meaning, and live outside of the color space(s).
	// This should not be used directly by applications.
	IsSpecial Color = 1 << 29
)

// Note that the order of these options is important -- it follows the
// definitions used by ECMA and XTerm.  Hence any further named colors
// must begin at a value not less than 256.
const (
	Black   = XTerm0
	Maroon  = XTerm1
	Green   = XTerm2
	Olive   = XTerm3
	Navy    = XTerm4
	Purple  = XTerm5
	Teal    = XTerm6
	Silver  = XTerm7
	Gray    = XTerm8
	Red     = XTerm9
	Lime    = XTerm10
	Yellow  = XTerm11
	Blue    = XTerm12
	Fuchsia = XTerm13
	Aqua    = XTerm14
	White   = XTerm15
)

const (
	XTerm0 = IsValid + iota
	XTerm1
	XTerm2
	XTerm3
	XTerm4
	XTerm5
	XTerm6
	XTerm7
	XTerm8
	XTerm9
	XTerm10
	XTerm11
	XTerm12
	XTerm13
	XTerm14
	XTerm15
	XTerm16
	XTerm17
	XTerm18
	XTerm19
	XTerm20
	XTerm21
	XTerm22
	XTerm23
	XTerm24
	XTerm25
	XTerm26
	XTerm27
	XTerm28
	XTerm29
	XTerm30
	XTerm31
	XTerm32
	XTerm33
	XTerm34
	XTerm35
	XTerm36
	XTerm37
	XTerm38
	XTerm39
	XTerm40
	XTerm41
	XTerm42
	XTerm43
	XTerm44
	XTerm45
	XTerm46
	XTerm47
	XTerm48
	XTerm49
	XTerm50
	XTerm51
	XTerm52
	XTerm53
	XTerm54
	XTerm55
	XTerm56
	XTerm57
	XTerm58
	XTerm59
	XTerm60
	XTerm61
	XTerm62
	XTerm63
	XTerm64
	XTerm65
	XTerm66
	XTerm67
	XTerm68
	XTerm69
	XTerm70
	XTerm71
	XTerm72
	XTerm73
	XTerm74
	XTerm75
	XTerm76
	XTerm77
	XTerm78
	XTerm79
	XTerm80
	XTerm81
	XTerm82
	XTerm83
	XTerm84
	XTerm85
	XTerm86
	XTerm87
	XTerm88
	XTerm89
	XTerm90
	XTerm91
	XTerm92
	XTerm93
	XTerm94
	XTerm95
	XTerm96
	XTerm97
	XTerm98
	XTerm99
	XTerm100
	XTerm101
	XTerm102
	XTerm103
	XTerm104
	XTerm105
	XTerm106
	XTerm107
	XTerm108
	XTerm109
	XTerm110
	XTerm111
	XTerm112
	XTerm113
	XTerm114
	XTerm115
	XTerm116
	XTerm117
	XTerm118
	XTerm119
	XTerm120
	XTerm121
	XTerm122
	XTerm123
	XTerm124
	XTerm125
	XTerm126
	XTerm127
	XTerm128
	XTerm129
	XTerm130
	XTerm131
	XTerm132
	XTerm133
	XTerm134
	XTerm135
	XTerm136
	XTerm137
	XTerm138
	XTerm139
	XTerm140
	XTerm141
	XTerm142
	XTerm143
	XTerm144
	XTerm145
	XTerm146
	XTerm147
	XTerm148
	XTerm149
	XTerm150
	XTerm151
	XTerm152
	XTerm153
	XTerm154
	XTerm155
	XTerm156
	XTerm157
	XTerm158
	XTerm159
	XTerm160
	XTerm161
	XTerm162
	XTerm163
	XTerm164
	XTerm165
	XTerm166
	XTerm167
	XTerm168
	XTerm169
	XTerm170
	XTerm171
	XTerm172
	XTerm173
	XTerm174
	XTerm175
	XTerm176
	XTerm177
	XTerm178
	XTerm179
	XTerm180
	XTerm181
	XTerm182
	XTerm183
	XTerm184
	XTerm185
	XTerm186
	XTerm187
	XTerm188
	XTerm189
	XTerm190
	XTerm191
	XTerm192
	XTerm193
	XTerm194
	XTerm195
	XTerm196
	XTerm197
	XTerm198
	XTerm199
	XTerm200
	XTerm201
	XTerm202
	XTerm203
	XTerm204
	XTerm205
	XTerm206
	XTerm207
	XTerm208
	XTerm209
	XTerm210
	XTerm211
	XTerm212
	XTerm213
	XTerm214
	XTerm215
	XTerm216
	XTerm217
	XTerm218
	XTerm219
	XTerm220
	XTerm221
	XTerm222
	XTerm223
	XTerm224
	XTerm225
	XTerm226
	XTerm227
	XTerm228
	XTerm229
	XTerm230
	XTerm231
	XTerm232
	XTerm233
	XTerm234
	XTerm235
	XTerm236
	XTerm237
	XTerm238
	XTerm239
	XTerm240
	XTerm241
	XTerm242
	XTerm243
	XTerm244
	XTerm245
	XTerm246
	XTerm247
	XTerm248
	XTerm249
	XTerm250
	XTerm251
	XTerm252
	XTerm253
	XTerm254
	XTerm255
	AliceBlue            = IsRGB | IsValid | 0xF0F8FF
	AntiqueWhite         = IsRGB | IsValid | 0xFAEBD7
	AquaMarine           = IsRGB | IsValid | 0x7FFFD4
	Azure                = IsRGB | IsValid | 0xF0FFFF
	Beige                = IsRGB | IsValid | 0xF5F5DC
	Bisque               = IsRGB | IsValid | 0xFFE4C4
	BlanchedAlmond       = IsRGB | IsValid | 0xFFEBCD
	BlueViolet           = IsRGB | IsValid | 0x8A2BE2
	Brown                = IsRGB | IsValid | 0xA52A2A
	BurlyWood            = IsRGB | IsValid | 0xDEB887
	CadetBlue            = IsRGB | IsValid | 0x5F9EA0
	Chartreuse           = IsRGB | IsValid | 0x7FFF00
	Chocolate            = IsRGB | IsValid | 0xD2691E
	Coral                = IsRGB | IsValid | 0xFF7F50
	CornflowerBlue       = IsRGB | IsValid | 0x6495ED
	Cornsilk             = IsRGB | IsValid | 0xFFF8DC
	Crimson              = IsRGB | IsValid | 0xDC143C
	DarkBlue             = IsRGB | IsValid | 0x00008B
	DarkCyan             = IsRGB | IsValid | 0x008B8B
	DarkGoldenrod        = IsRGB | IsValid | 0xB8860B
	DarkGray             = IsRGB | IsValid | 0xA9A9A9
	DarkGreen            = IsRGB | IsValid | 0x006400
	DarkKhaki            = IsRGB | IsValid | 0xBDB76B
	DarkMagenta          = IsRGB | IsValid | 0x8B008B
	DarkOliveGreen       = IsRGB | IsValid | 0x556B2F
	DarkOrange           = IsRGB | IsValid | 0xFF8C00
	DarkOrchid           = IsRGB | IsValid | 0x9932CC
	DarkRed              = IsRGB | IsValid | 0x8B0000
	DarkSalmon           = IsRGB | IsValid | 0xE9967A
	DarkSeaGreen         = IsRGB | IsValid | 0x8FBC8F
	DarkSlateBlue        = IsRGB | IsValid | 0x483D8B
	DarkSlateGray        = IsRGB | IsValid | 0x2F4F4F
	DarkTurquoise        = IsRGB | IsValid | 0x00CED1
	DarkViolet           = IsRGB | IsValid | 0x9400D3
	DeepPink             = IsRGB | IsValid | 0xFF1493
	DeepSkyBlue          = IsRGB | IsValid | 0x00BFFF
	DimGray              = IsRGB | IsValid | 0x696969
	DodgerBlue           = IsRGB | IsValid | 0x1E90FF
	FireBrick            = IsRGB | IsValid | 0xB22222
	FloralWhite          = IsRGB | IsValid | 0xFFFAF0
	ForestGreen          = IsRGB | IsValid | 0x228B22
	Gainsboro            = IsRGB | IsValid | 0xDCDCDC
	GhostWhite           = IsRGB | IsValid | 0xF8F8FF
	Gold                 = IsRGB | IsValid | 0xFFD700
	Goldenrod            = IsRGB | IsValid | 0xDAA520
	GreenYellow          = IsRGB | IsValid | 0xADFF2F
	Honeydew             = IsRGB | IsValid | 0xF0FFF0
	HotPink              = IsRGB | IsValid | 0xFF69B4
	IndianRed            = IsRGB | IsValid | 0xCD5C5C
	Indigo               = IsRGB | IsValid | 0x4B0082
	Ivory                = IsRGB | IsValid | 0xFFFFF0
	Khaki                = IsRGB | IsValid | 0xF0E68C
	Lavender             = IsRGB | IsValid | 0xE6E6FA
	LavenderBlush        = IsRGB | IsValid | 0xFFF0F5
	LawnGreen            = IsRGB | IsValid | 0x7CFC00
	LemonChiffon         = IsRGB | IsValid | 0xFFFACD
	LightBlue            = IsRGB | IsValid | 0xADD8E6
	LightCoral           = IsRGB | IsValid | 0xF08080
	LightCyan            = IsRGB | IsValid | 0xE0FFFF
	LightGoldenrodYellow = IsRGB | IsValid | 0xFAFAD2
	LightGray            = IsRGB | IsValid | 0xD3D3D3
	LightGreen           = IsRGB | IsValid | 0x90EE90
	LightPink            = IsRGB | IsValid | 0xFFB6C1
	LightSalmon          = IsRGB | IsValid | 0xFFA07A
	LightSeaGreen        = IsRGB | IsValid | 0x20B2AA
	LightSkyBlue         = IsRGB | IsValid | 0x87CEFA
	LightSlateGray       = IsRGB | IsValid | 0x778899
	LightSteelBlue       = IsRGB | IsValid | 0xB0C4DE
	LightYellow          = IsRGB | IsValid | 0xFFFFE0
	LimeGreen            = IsRGB | IsValid | 0x32CD32
	Linen                = IsRGB | IsValid | 0xFAF0E6
	MediumAquamarine     = IsRGB | IsValid | 0x66CDAA
	MediumBlue           = IsRGB | IsValid | 0x0000CD
	MediumOrchid         = IsRGB | IsValid | 0xBA55D3
	MediumPurple         = IsRGB | IsValid | 0x9370DB
	MediumSeaGreen       = IsRGB | IsValid | 0x3CB371
	MediumSlateBlue      = IsRGB | IsValid | 0x7B68EE
	MediumSpringGreen    = IsRGB | IsValid | 0x00FA9A
	MediumTurquoise      = IsRGB | IsValid | 0x48D1CC
	MediumVioletRed      = IsRGB | IsValid | 0xC71585
	MidnightBlue         = IsRGB | IsValid | 0x191970
	MintCream            = IsRGB | IsValid | 0xF5FFFA
	MistyRose            = IsRGB | IsValid | 0xFFE4E1
	Moccasin             = IsRGB | IsValid | 0xFFE4B5
	NavajoWhite          = IsRGB | IsValid | 0xFFDEAD
	OldLace              = IsRGB | IsValid | 0xFDF5E6
	OliveDrab            = IsRGB | IsValid | 0x6B8E23
	Orange               = IsRGB | IsValid | 0xFFA500
	OrangeRed            = IsRGB | IsValid | 0xFF4500
	Orchid               = IsRGB | IsValid | 0xDA70D6
	PaleGoldenrod        = IsRGB | IsValid | 0xEEE8AA
	PaleGreen            = IsRGB | IsValid | 0x98FB98
	PaleTurquoise        = IsRGB | IsValid | 0xAFEEEE
	PaleVioletRed        = IsRGB | IsValid | 0xDB7093
	PapayaWhip           = IsRGB | IsValid | 0xFFEFD5
	PeachPuff            = IsRGB | IsValid | 0xFFDAB9
	Peru                 = IsRGB | IsValid | 0xCD853F
	Pink                 = IsRGB | IsValid | 0xFFC0CB
	Plum                 = IsRGB | IsValid | 0xDDA0DD
	PowderBlue           = IsRGB | IsValid | 0xB0E0E6
	RebeccaPurple        = IsRGB | IsValid | 0x663399
	RosyBrown            = IsRGB | IsValid | 0xBC8F8F
	RoyalBlue            = IsRGB | IsValid | 0x4169E1
	SaddleBrown          = IsRGB | IsValid | 0x8B4513
	Salmon               = IsRGB | IsValid | 0xFA8072
	SandyBrown           = IsRGB | IsValid | 0xF4A460
	SeaGreen             = IsRGB | IsValid | 0x2E8B57
	Seashell             = IsRGB | IsValid | 0xFFF5EE
	Sienna               = IsRGB | IsValid | 0xA0522D
	Skyblue              = IsRGB | IsValid | 0x87CEEB
	SlateBlue            = IsRGB | IsValid | 0x6A5ACD
	SlateGray            = IsRGB | IsValid | 0x708090
	Snow                 = IsRGB | IsValid | 0xFFFAFA
	SpringGreen          = IsRGB | IsValid | 0x00FF7F
	SteelBlue            = IsRGB | IsValid | 0x4682B4
	Tan                  = IsRGB | IsValid | 0xD2B48C
	Thistle              = IsRGB | IsValid | 0xD8BFD8
	Tomato               = IsRGB | IsValid | 0xFF6347
	Turquoise            = IsRGB | IsValid | 0x40E0D0
	Violet               = IsRGB | IsValid | 0xEE82EE
	Wheat                = IsRGB | IsValid | 0xF5DEB3
	WhiteSmoke           = IsRGB | IsValid | 0xF5F5F5
	YellowGreen          = IsRGB | IsValid | 0x9ACD32
)

// These are aliases for the color gray, because some of us spell
// it as grey.
const (
	Grey           = Gray
	DimGrey        = DimGray
	DarkGrey       = DarkGray
	DarkSlateGrey  = DarkSlateGray
	LightGrey      = LightGray
	LightSlateGrey = LightSlateGray
	SlateGrey      = SlateGray
)

// ColorValues maps color constants to their RGB values.
var ColorValues = map[Color]int32{
	Black:                0x000000,
	Maroon:               0x800000,
	Green:                0x008000,
	Olive:                0x808000,
	Navy:                 0x000080,
	Purple:               0x800080,
	Teal:                 0x008080,
	Silver:               0xC0C0C0,
	Gray:                 0x808080,
	Red:                  0xFF0000,
	Lime:                 0x00FF00,
	Yellow:               0xFFFF00,
	Blue:                 0x0000FF,
	Fuchsia:              0xFF00FF,
	Aqua:                 0x00FFFF,
	White:                0xFFFFFF,
	XTerm16:              0x000000, // black
	XTerm17:              0x00005F,
	XTerm18:              0x000087,
	XTerm19:              0x0000AF,
	XTerm20:              0x0000D7,
	XTerm21:              0x0000FF, // blue
	XTerm22:              0x005F00,
	XTerm23:              0x005F5F,
	XTerm24:              0x005F87,
	XTerm25:              0x005FAF,
	XTerm26:              0x005FD7,
	XTerm27:              0x005FFF,
	XTerm28:              0x008700,
	XTerm29:              0x00875F,
	XTerm30:              0x008787,
	XTerm31:              0x0087Af,
	XTerm32:              0x0087D7,
	XTerm33:              0x0087FF,
	XTerm34:              0x00AF00,
	XTerm35:              0x00AF5F,
	XTerm36:              0x00AF87,
	XTerm37:              0x00AFAF,
	XTerm38:              0x00AFD7,
	XTerm39:              0x00AFFF,
	XTerm40:              0x00D700,
	XTerm41:              0x00D75F,
	XTerm42:              0x00D787,
	XTerm43:              0x00D7AF,
	XTerm44:              0x00D7D7,
	XTerm45:              0x00D7FF,
	XTerm46:              0x00FF00, // lime
	XTerm47:              0x00FF5F,
	XTerm48:              0x00FF87,
	XTerm49:              0x00FFAF,
	XTerm50:              0x00FFd7,
	XTerm51:              0x00FFFF, // aqua
	XTerm52:              0x5F0000,
	XTerm53:              0x5F005F,
	XTerm54:              0x5F0087,
	XTerm55:              0x5F00AF,
	XTerm56:              0x5F00D7,
	XTerm57:              0x5F00FF,
	XTerm58:              0x5F5F00,
	XTerm59:              0x5F5F5F,
	XTerm60:              0x5F5F87,
	XTerm61:              0x5F5FAF,
	XTerm62:              0x5F5FD7,
	XTerm63:              0x5F5FFF,
	XTerm64:              0x5F8700,
	XTerm65:              0x5F875F,
	XTerm66:              0x5F8787,
	XTerm67:              0x5F87AF,
	XTerm68:              0x5F87D7,
	XTerm69:              0x5F87FF,
	XTerm70:              0x5FAF00,
	XTerm71:              0x5FAF5F,
	XTerm72:              0x5FAF87,
	XTerm73:              0x5FAFAF,
	XTerm74:              0x5FAFD7,
	XTerm75:              0x5FAFFF,
	XTerm76:              0x5FD700,
	XTerm77:              0x5FD75F,
	XTerm78:              0x5FD787,
	XTerm79:              0x5FD7AF,
	XTerm80:              0x5FD7D7,
	XTerm81:              0x5FD7FF,
	XTerm82:              0x5FFF00,
	XTerm83:              0x5FFF5F,
	XTerm84:              0x5FFF87,
	XTerm85:              0x5FFFAF,
	XTerm86:              0x5FFFD7,
	XTerm87:              0x5FFFFF,
	XTerm88:              0x870000,
	XTerm89:              0x87005F,
	XTerm90:              0x870087,
	XTerm91:              0x8700AF,
	XTerm92:              0x8700D7,
	XTerm93:              0x8700FF,
	XTerm94:              0x875F00,
	XTerm95:              0x875F5F,
	XTerm96:              0x875F87,
	XTerm97:              0x875FAF,
	XTerm98:              0x875FD7,
	XTerm99:              0x875FFF,
	XTerm100:             0x878700,
	XTerm101:             0x87875F,
	XTerm102:             0x878787,
	XTerm103:             0x8787AF,
	XTerm104:             0x8787D7,
	XTerm105:             0x8787FF,
	XTerm106:             0x87AF00,
	XTerm107:             0x87AF5F,
	XTerm108:             0x87AF87,
	XTerm109:             0x87AFAF,
	XTerm110:             0x87AFD7,
	XTerm111:             0x87AFFF,
	XTerm112:             0x87D700,
	XTerm113:             0x87D75F,
	XTerm114:             0x87D787,
	XTerm115:             0x87D7AF,
	XTerm116:             0x87D7D7,
	XTerm117:             0x87D7FF,
	XTerm118:             0x87FF00,
	XTerm119:             0x87FF5F,
	XTerm120:             0x87FF87,
	XTerm121:             0x87FFAF,
	XTerm122:             0x87FFD7,
	XTerm123:             0x87FFFF,
	XTerm124:             0xAF0000,
	XTerm125:             0xAF005F,
	XTerm126:             0xAF0087,
	XTerm127:             0xAF00AF,
	XTerm128:             0xAF00D7,
	XTerm129:             0xAF00FF,
	XTerm130:             0xAF5F00,
	XTerm131:             0xAF5F5F,
	XTerm132:             0xAF5F87,
	XTerm133:             0xAF5FAF,
	XTerm134:             0xAF5FD7,
	XTerm135:             0xAF5FFF,
	XTerm136:             0xAF8700,
	XTerm137:             0xAF875F,
	XTerm138:             0xAF8787,
	XTerm139:             0xAF87AF,
	XTerm140:             0xAF87D7,
	XTerm141:             0xAF87FF,
	XTerm142:             0xAFAF00,
	XTerm143:             0xAFAF5F,
	XTerm144:             0xAFAF87,
	XTerm145:             0xAFAFAF,
	XTerm146:             0xAFAFD7,
	XTerm147:             0xAFAFFF,
	XTerm148:             0xAFD700,
	XTerm149:             0xAFD75F,
	XTerm150:             0xAFD787,
	XTerm151:             0xAFD7AF,
	XTerm152:             0xAFD7D7,
	XTerm153:             0xAFD7FF,
	XTerm154:             0xAFFF00,
	XTerm155:             0xAFFF5F,
	XTerm156:             0xAFFF87,
	XTerm157:             0xAFFFAF,
	XTerm158:             0xAFFFD7,
	XTerm159:             0xAFFFFF,
	XTerm160:             0xD70000,
	XTerm161:             0xD7005F,
	XTerm162:             0xD70087,
	XTerm163:             0xD700AF,
	XTerm164:             0xD700D7,
	XTerm165:             0xD700FF,
	XTerm166:             0xD75F00,
	XTerm167:             0xD75F5F,
	XTerm168:             0xD75F87,
	XTerm169:             0xD75FAF,
	XTerm170:             0xD75FD7,
	XTerm171:             0xD75FFF,
	XTerm172:             0xD78700,
	XTerm173:             0xD7875F,
	XTerm174:             0xD78787,
	XTerm175:             0xD787AF,
	XTerm176:             0xD787D7,
	XTerm177:             0xD787FF,
	XTerm178:             0xD7AF00,
	XTerm179:             0xD7AF5F,
	XTerm180:             0xD7AF87,
	XTerm181:             0xD7AFAF,
	XTerm182:             0xD7AFD7,
	XTerm183:             0xD7AFFF,
	XTerm184:             0xD7D700,
	XTerm185:             0xD7D75F,
	XTerm186:             0xD7D787,
	XTerm187:             0xD7D7AF,
	XTerm188:             0xD7D7D7,
	XTerm189:             0xD7D7FF,
	XTerm190:             0xD7FF00,
	XTerm191:             0xD7FF5F,
	XTerm192:             0xD7FF87,
	XTerm193:             0xD7FFAF,
	XTerm194:             0xD7FFD7,
	XTerm195:             0xD7FFFF,
	XTerm196:             0xFF0000, // red
	XTerm197:             0xFF005F,
	XTerm198:             0xFF0087,
	XTerm199:             0xFF00AF,
	XTerm200:             0xFF00D7,
	XTerm201:             0xFF00FF, // fuchsia
	XTerm202:             0xFF5F00,
	XTerm203:             0xFF5F5F,
	XTerm204:             0xFF5F87,
	XTerm205:             0xFF5FAF,
	XTerm206:             0xFF5FD7,
	XTerm207:             0xFF5FFF,
	XTerm208:             0xFF8700,
	XTerm209:             0xFF875F,
	XTerm210:             0xFF8787,
	XTerm211:             0xFF87AF,
	XTerm212:             0xFF87D7,
	XTerm213:             0xFF87FF,
	XTerm214:             0xFFAF00,
	XTerm215:             0xFFAF5F,
	XTerm216:             0xFFAF87,
	XTerm217:             0xFFAFAF,
	XTerm218:             0xFFAFD7,
	XTerm219:             0xFFAFFF,
	XTerm220:             0xFFD700,
	XTerm221:             0xFFD75F,
	XTerm222:             0xFFD787,
	XTerm223:             0xFFD7AF,
	XTerm224:             0xFFD7D7,
	XTerm225:             0xFFD7FF,
	XTerm226:             0xFFFF00, // yellow
	XTerm227:             0xFFFF5F,
	XTerm228:             0xFFFF87,
	XTerm229:             0xFFFFAF,
	XTerm230:             0xFFFFD7,
	XTerm231:             0xFFFFFF, // white
	XTerm232:             0x080808,
	XTerm233:             0x121212,
	XTerm234:             0x1C1C1C,
	XTerm235:             0x262626,
	XTerm236:             0x303030,
	XTerm237:             0x3A3A3A,
	XTerm238:             0x444444,
	XTerm239:             0x4E4E4E,
	XTerm240:             0x585858,
	XTerm241:             0x626262,
	XTerm242:             0x6C6C6C,
	XTerm243:             0x767676,
	XTerm244:             0x808080, // grey
	XTerm245:             0x8A8A8A,
	XTerm246:             0x949494,
	XTerm247:             0x9E9E9E,
	XTerm248:             0xA8A8A8,
	XTerm249:             0xB2B2B2,
	XTerm250:             0xBCBCBC,
	XTerm251:             0xC6C6C6,
	XTerm252:             0xD0D0D0,
	XTerm253:             0xDADADA,
	XTerm254:             0xE4E4E4,
	XTerm255:             0xEEEEEE,
	AliceBlue:            0xF0F8FF,
	AntiqueWhite:         0xFAEBD7,
	AquaMarine:           0x7FFFD4,
	Azure:                0xF0FFFF,
	Beige:                0xF5F5DC,
	Bisque:               0xFFE4C4,
	BlanchedAlmond:       0xFFEBCD,
	BlueViolet:           0x8A2BE2,
	Brown:                0xA52A2A,
	BurlyWood:            0xDEB887,
	CadetBlue:            0x5F9EA0,
	Chartreuse:           0x7FFF00,
	Chocolate:            0xD2691E,
	Coral:                0xFF7F50,
	CornflowerBlue:       0x6495ED,
	Cornsilk:             0xFFF8DC,
	Crimson:              0xDC143C,
	DarkBlue:             0x00008B,
	DarkCyan:             0x008B8B,
	DarkGoldenrod:        0xB8860B,
	DarkGray:             0xA9A9A9,
	DarkGreen:            0x006400,
	DarkKhaki:            0xBDB76B,
	DarkMagenta:          0x8B008B,
	DarkOliveGreen:       0x556B2F,
	DarkOrange:           0xFF8C00,
	DarkOrchid:           0x9932CC,
	DarkRed:              0x8B0000,
	DarkSalmon:           0xE9967A,
	DarkSeaGreen:         0x8FBC8F,
	DarkSlateBlue:        0x483D8B,
	DarkSlateGray:        0x2F4F4F,
	DarkTurquoise:        0x00CED1,
	DarkViolet:           0x9400D3,
	DeepPink:             0xFF1493,
	DeepSkyBlue:          0x00BFFF,
	DimGray:              0x696969,
	DodgerBlue:           0x1E90FF,
	FireBrick:            0xB22222,
	FloralWhite:          0xFFFAF0,
	ForestGreen:          0x228B22,
	Gainsboro:            0xDCDCDC,
	GhostWhite:           0xF8F8FF,
	Gold:                 0xFFD700,
	Goldenrod:            0xDAA520,
	GreenYellow:          0xADFF2F,
	Honeydew:             0xF0FFF0,
	HotPink:              0xFF69B4,
	IndianRed:            0xCD5C5C,
	Indigo:               0x4B0082,
	Ivory:                0xFFFFF0,
	Khaki:                0xF0E68C,
	Lavender:             0xE6E6FA,
	LavenderBlush:        0xFFF0F5,
	LawnGreen:            0x7CFC00,
	LemonChiffon:         0xFFFACD,
	LightBlue:            0xADD8E6,
	LightCoral:           0xF08080,
	LightCyan:            0xE0FFFF,
	LightGoldenrodYellow: 0xFAFAD2,
	LightGray:            0xD3D3D3,
	LightGreen:           0x90EE90,
	LightPink:            0xFFB6C1,
	LightSalmon:          0xFFA07A,
	LightSeaGreen:        0x20B2AA,
	LightSkyBlue:         0x87CEFA,
	LightSlateGray:       0x778899,
	LightSteelBlue:       0xB0C4DE,
	LightYellow:          0xFFFFE0,
	LimeGreen:            0x32CD32,
	Linen:                0xFAF0E6,
	MediumAquamarine:     0x66CDAA,
	MediumBlue:           0x0000CD,
	MediumOrchid:         0xBA55D3,
	MediumPurple:         0x9370DB,
	MediumSeaGreen:       0x3CB371,
	MediumSlateBlue:      0x7B68EE,
	MediumSpringGreen:    0x00FA9A,
	MediumTurquoise:      0x48D1CC,
	MediumVioletRed:      0xC71585,
	MidnightBlue:         0x191970,
	MintCream:            0xF5FFFA,
	MistyRose:            0xFFE4E1,
	Moccasin:             0xFFE4B5,
	NavajoWhite:          0xFFDEAD,
	OldLace:              0xFDF5E6,
	OliveDrab:            0x6B8E23,
	Orange:               0xFFA500,
	OrangeRed:            0xFF4500,
	Orchid:               0xDA70D6,
	PaleGoldenrod:        0xEEE8AA,
	PaleGreen:            0x98FB98,
	PaleTurquoise:        0xAFEEEE,
	PaleVioletRed:        0xDB7093,
	PapayaWhip:           0xFFEFD5,
	PeachPuff:            0xFFDAB9,
	Peru:                 0xCD853F,
	Pink:                 0xFFC0CB,
	Plum:                 0xDDA0DD,
	PowderBlue:           0xB0E0E6,
	RebeccaPurple:        0x663399,
	RosyBrown:            0xBC8F8F,
	RoyalBlue:            0x4169E1,
	SaddleBrown:          0x8B4513,
	Salmon:               0xFA8072,
	SandyBrown:           0xF4A460,
	SeaGreen:             0x2E8B57,
	Seashell:             0xFFF5EE,
	Sienna:               0xA0522D,
	Skyblue:              0x87CEEB,
	SlateBlue:            0x6A5ACD,
	SlateGray:            0x708090,
	Snow:                 0xFFFAFA,
	SpringGreen:          0x00FF7F,
	SteelBlue:            0x4682B4,
	Tan:                  0xD2B48C,
	Thistle:              0xD8BFD8,
	Tomato:               0xFF6347,
	Turquoise:            0x40E0D0,
	Violet:               0xEE82EE,
	Wheat:                0xF5DEB3,
	WhiteSmoke:           0xF5F5F5,
	YellowGreen:          0x9ACD32,
}

// Special colors.
const (
	// Reset is used to indicate that the color should use the
	// vanilla terminal colors.  (Basically go back to the defaults.)
	Reset = IsSpecial | iota

	// None indicates that we should not change the color from
	// whatever is already displayed.  This can only be used in limited
	// circumstances.
	None
)

// Names holds the written names of colors. Useful to present a list of
// recognized named colors.
var Names = map[string]Color{
	"black":                Black,
	"maroon":               Maroon,
	"green":                Green,
	"olive":                Olive,
	"navy":                 Navy,
	"purple":               Purple,
	"teal":                 Teal,
	"silver":               Silver,
	"gray":                 Gray,
	"red":                  Red,
	"lime":                 Lime,
	"yellow":               Yellow,
	"blue":                 Blue,
	"fuchsia":              Fuchsia,
	"aqua":                 Aqua,
	"white":                White,
	"aliceblue":            AliceBlue,
	"antiquewhite":         AntiqueWhite,
	"aquamarine":           AquaMarine,
	"azure":                Azure,
	"beige":                Beige,
	"bisque":               Bisque,
	"blanchedalmond":       BlanchedAlmond,
	"blueviolet":           BlueViolet,
	"brown":                Brown,
	"burlywood":            BurlyWood,
	"cadetblue":            CadetBlue,
	"chartreuse":           Chartreuse,
	"chocolate":            Chocolate,
	"coral":                Coral,
	"cornflowerblue":       CornflowerBlue,
	"cornsilk":             Cornsilk,
	"crimson":              Crimson,
	"darkblue":             DarkBlue,
	"darkcyan":             DarkCyan,
	"darkgoldenrod":        DarkGoldenrod,
	"darkgray":             DarkGray,
	"darkgreen":            DarkGreen,
	"darkkhaki":            DarkKhaki,
	"darkmagenta":          DarkMagenta,
	"darkolivegreen":       DarkOliveGreen,
	"darkorange":           DarkOrange,
	"darkorchid":           DarkOrchid,
	"darkred":              DarkRed,
	"darksalmon":           DarkSalmon,
	"darkseagreen":         DarkSeaGreen,
	"darkslateblue":        DarkSlateBlue,
	"darkslategray":        DarkSlateGray,
	"darkturquoise":        DarkTurquoise,
	"darkviolet":           DarkViolet,
	"deeppink":             DeepPink,
	"deepskyblue":          DeepSkyBlue,
	"dimgray":              DimGray,
	"dodgerblue":           DodgerBlue,
	"firebrick":            FireBrick,
	"floralwhite":          FloralWhite,
	"forestgreen":          ForestGreen,
	"gainsboro":            Gainsboro,
	"ghostwhite":           GhostWhite,
	"gold":                 Gold,
	"goldenrod":            Goldenrod,
	"greenyellow":          GreenYellow,
	"honeydew":             Honeydew,
	"hotpink":              HotPink,
	"indianred":            IndianRed,
	"indigo":               Indigo,
	"ivory":                Ivory,
	"khaki":                Khaki,
	"lavender":             Lavender,
	"lavenderblush":        LavenderBlush,
	"lawngreen":            LawnGreen,
	"lemonchiffon":         LemonChiffon,
	"lightblue":            LightBlue,
	"lightcoral":           LightCoral,
	"lightcyan":            LightCyan,
	"lightgoldenrodyellow": LightGoldenrodYellow,
	"lightgray":            LightGray,
	"lightgreen":           LightGreen,
	"lightpink":            LightPink,
	"lightsalmon":          LightSalmon,
	"lightseagreen":        LightSeaGreen,
	"lightskyblue":         LightSkyBlue,
	"lightslategray":       LightSlateGray,
	"lightsteelblue":       LightSteelBlue,
	"lightyellow":          LightYellow,
	"limegreen":            LimeGreen,
	"linen":                Linen,
	"mediumaquamarine":     MediumAquamarine,
	"mediumblue":           MediumBlue,
	"mediumorchid":         MediumOrchid,
	"mediumpurple":         MediumPurple,
	"mediumseagreen":       MediumSeaGreen,
	"mediumslateblue":      MediumSlateBlue,
	"mediumspringgreen":    MediumSpringGreen,
	"mediumturquoise":      MediumTurquoise,
	"mediumvioletred":      MediumVioletRed,
	"midnightblue":         MidnightBlue,
	"mintcream":            MintCream,
	"mistyrose":            MistyRose,
	"moccasin":             Moccasin,
	"navajowhite":          NavajoWhite,
	"oldlace":              OldLace,
	"olivedrab":            OliveDrab,
	"orange":               Orange,
	"orangered":            OrangeRed,
	"orchid":               Orchid,
	"palegoldenrod":        PaleGoldenrod,
	"palegreen":            PaleGreen,
	"paleturquoise":        PaleTurquoise,
	"palevioletred":        PaleVioletRed,
	"papayawhip":           PapayaWhip,
	"peachpuff":            PeachPuff,
	"peru":                 Peru,
	"pink":                 Pink,
	"plum":                 Plum,
	"powderblue":           PowderBlue,
	"rebeccapurple":        RebeccaPurple,
	"rosybrown":            RosyBrown,
	"royalblue":            RoyalBlue,
	"saddlebrown":          SaddleBrown,
	"salmon":               Salmon,
	"sandybrown":           SandyBrown,
	"seagreen":             SeaGreen,
	"seashell":             Seashell,
	"sienna":               Sienna,
	"skyblue":              Skyblue,
	"slateblue":            SlateBlue,
	"slategray":            SlateGray,
	"snow":                 Snow,
	"springgreen":          SpringGreen,
	"steelblue":            SteelBlue,
	"tan":                  Tan,
	"thistle":              Thistle,
	"tomato":               Tomato,
	"turquoise":            Turquoise,
	"violet":               Violet,
	"wheat":                Wheat,
	"whitesmoke":           WhiteSmoke,
	"yellowgreen":          YellowGreen,
	"grey":                 Gray,
	"dimgrey":              DimGray,
	"darkgrey":             DarkGray,
	"darkslategrey":        DarkSlateGray,
	"lightgrey":            LightGray,
	"lightslategrey":       LightSlateGray,
	"slategrey":            SlateGray,
}

// Valid indicates the color is a valid value (has been set).
func (c Color) Valid() bool {
	return c&IsValid != 0
}

// IsRGB is true if the color is an RGB specific value.
func (c Color) IsRGB() bool {
	return c&(IsValid|IsRGB) == (IsValid | IsRGB)
}

// CSS returns the CSS hex string ( #ABCDEF ) if valid
// if not a valid color returns empty string
func (c Color) CSS() string {
	if !c.Valid() {
		return ""
	}
	return fmt.Sprintf("#%06X", c.Hex())
}

// String implements fmt.Stringer to return either the
// W3C name if it has one or the CSS hex string '#ABCDEF'
func (c Color) String() string {
	if !c.Valid() {
		switch c {
		case None:
			return "none"
		case Default:
			return "default"
		case Reset:
			return "reset"
		}
		return ""
	}
	return c.Name(true)
}

// Name returns W3C name or an empty string if no arguments
// if passed true as an argument it will falls back to
// the CSS hex string if no W3C name found '#ABCDEF'
func (c Color) Name(css ...bool) string {
	for name, hex := range Names {
		if c == hex {
			return name
		}
	}
	if len(css) > 0 && css[0] {
		return c.CSS()
	}
	return ""
}

// Hex returns the color's hexadecimal RGB 24-bit value with each component
// consisting of a single byte, R << 16 | G << 8 | B.  If the color
// is unknown or unset, -1 is returned.
func (c Color) Hex() int32 {
	if !c.Valid() {
		return -1
	}
	if c&IsRGB != 0 {
		return int32(c & 0xffffff)
	}
	if v, ok := ColorValues[c]; ok {
		return v
	}
	return -1
}

// RGB returns the red, green, and blue components of the color, with
// each component represented as a value 0-255.  In the event that the
// color cannot be broken up (not set usually), -1 is returned for each value.
func (c Color) RGB() (int32, int32, int32) {
	v := c.Hex()
	if v < 0 {
		return -1, -1, -1
	}
	return (v >> 16) & 0xff, (v >> 8) & 0xff, v & 0xff
}

// TrueColor returns the true color (RGB) version of the provided color.
// This is useful for ensuring color accuracy when using named colors.
// This will override terminal theme colors.
func (c Color) TrueColor() Color {
	if !c.Valid() {
		return Default
	}
	if c&IsRGB != 0 {
		return c | IsValid
	}
	if hex := c.Hex(); hex < 0 {
		return Default
	} else {
		return Color(hex) | IsRGB | IsValid
	}
}

// RGBA makes these colors directly usable as imageColor colors.
// The values are scaled only to 16 bits.  Invalid colors are returned
// with all values being zero (notably the alpha is zero, so fully transparent),
// otherwise the alpha channel is set to 0xffff (fully opaque).
func (c Color) RGBA() (r, g, b, a uint32) {
	if !c.Valid() {
		return 0, 0, 0, 0
	}
	r1, g1, b1 := c.RGB()
	r = uint32(r1)
	g = uint32(g1)
	b = uint32(b1)
	r = r | r<<8
	g = g | g<<8
	b = b | b<<8
	a = 0xffff
	return r, g, b, a
}

// NewRGBColor returns a new color with the given red, green, and blue values.
// Each value must be represented in the range 0-255.
func NewRGBColor(r, g, b int32) Color {
	return NewHexColor(((r & 0xff) << 16) | ((g & 0xff) << 8) | (b & 0xff))
}

// NewHexColor returns a color using the given 24-bit RGB value.
func NewHexColor(v int32) Color {
	return IsRGB | Color(v) | IsValid
}

// GetColor creates a Color from a color name (W3C name). A hex value may
// be supplied as a string in the format "#ffffff".
func GetColor(name string) Color {
	if c, ok := Names[name]; ok {
		return c
	}
	if len(name) == 7 && name[0] == '#' {
		if v, e := strconv.ParseInt(name[1:], 16, 32); e == nil {
			return NewHexColor(int32(v))
		}
	}
	return Default
}

// PaletteColor creates a color based on the palette index.
func PaletteColor(index int) Color {
	return Color(index) | IsValid
}

// FromImageColor converts an image/color.Color into Color.
// The alpha value is limited to just zero and non-zero, so it should
// be tracked separately if full detail is needed. (A zero alpha
// becomes the default color, which means no color change at all.)
func FromImageColor(imageColor ic.Color) Color {
	r, g, b, a := imageColor.RGBA()
	if a == 0 {
		return Default
	}
	// NOTE image/color.Color RGB values range is [0, 0xFFFF] as uint32
	return NewRGBColor(int32(r>>8), int32(g>>8), int32(b>>8))
}
