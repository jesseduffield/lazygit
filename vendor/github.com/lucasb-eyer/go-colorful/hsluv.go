package colorful

import "math"

// Source: https://github.com/hsluv/hsluv-go
// Under MIT License
// Modified so that Saturation and Luminance are in [0..1] instead of [0..100].

// HSLuv uses a rounded version of the D65. This has no impact on the final RGB
// values, but to keep high levels of accuracy for internal operations and when
// comparing to the test values, this modified white reference is used internally.
//
// See this GitHub thread for details on these values:
//
//	https://github.com/hsluv/hsluv/issues/79
var hSLuvD65 = [3]float64{0.95045592705167, 1.0, 1.089057750759878}

func LuvLChToHSLuv(l, c, h float64) (float64, float64, float64) {
	// [-1..1] but the code expects it to be [-100..100]
	c *= 100.0
	l *= 100.0

	var s, max float64
	if l > 99.9999999 || l < 0.00000001 {
		s = 0.0
	} else {
		max = maxChromaForLH(l, h)
		s = c / max * 100.0
	}
	return h, clamp01(s / 100.0), clamp01(l / 100.0)
}

func HSLuvToLuvLCh(h, s, l float64) (float64, float64, float64) {
	l *= 100.0
	s *= 100.0

	var c, max float64
	if l > 99.9999999 || l < 0.00000001 {
		c = 0.0
	} else {
		max = maxChromaForLH(l, h)
		c = max / 100.0 * s
	}

	// c is [-100..100], but for LCh it's supposed to be almost [-1..1]
	return clamp01(l / 100.0), c / 100.0, h
}

func LuvLChToHPLuv(l, c, h float64) (float64, float64, float64) {
	// [-1..1] but the code expects it to be [-100..100]
	c *= 100.0
	l *= 100.0

	var s, max float64
	if l > 99.9999999 || l < 0.00000001 {
		s = 0.0
	} else {
		max = maxSafeChromaForL(l)
		s = c / max * 100.0
	}
	return h, s / 100.0, l / 100.0
}

func HPLuvToLuvLCh(h, s, l float64) (float64, float64, float64) {
	// [-1..1] but the code expects it to be [-100..100]
	l *= 100.0
	s *= 100.0

	var c, max float64
	if l > 99.9999999 || l < 0.00000001 {
		c = 0.0
	} else {
		max = maxSafeChromaForL(l)
		c = max / 100.0 * s
	}
	return l / 100.0, c / 100.0, h
}

// HSLuv creates a new Color from values in the HSLuv color space.
// Hue in [0..360], a Saturation [0..1], and a Luminance (lightness) in [0..1].
//
// The returned color values are clamped (using .Clamped), so this will never output
// an invalid color.
func HSLuv(h, s, l float64) Color {
	// HSLuv -> LuvLCh -> CIELUV -> CIEXYZ -> Linear RGB -> sRGB
	l, u, v := LuvLChToLuv(HSLuvToLuvLCh(h, s, l))
	return LinearRgb(XyzToLinearRgb(LuvToXyzWhiteRef(l, u, v, hSLuvD65))).Clamped()
}

// HPLuv creates a new Color from values in the HPLuv color space.
// Hue in [0..360], a Saturation [0..1], and a Luminance (lightness) in [0..1].
//
// The returned color values are clamped (using .Clamped), so this will never output
// an invalid color.
func HPLuv(h, s, l float64) Color {
	// HPLuv -> LuvLCh -> CIELUV -> CIEXYZ -> Linear RGB -> sRGB
	l, u, v := LuvLChToLuv(HPLuvToLuvLCh(h, s, l))
	return LinearRgb(XyzToLinearRgb(LuvToXyzWhiteRef(l, u, v, hSLuvD65))).Clamped()
}

// HSLuv returns the Hue, Saturation and Luminance of the color in the HSLuv
// color space. Hue in [0..360], a Saturation [0..1], and a Luminance
// (lightness) in [0..1].
func (col Color) HSLuv() (h, s, l float64) {
	// sRGB -> Linear RGB -> CIEXYZ -> CIELUV -> LuvLCh -> HSLuv
	return LuvLChToHSLuv(col.LuvLChWhiteRef(hSLuvD65))
}

// HPLuv returns the Hue, Saturation and Luminance of the color in the HSLuv
// color space. Hue in [0..360], a Saturation [0..1], and a Luminance
// (lightness) in [0..1].
//
// Note that HPLuv can only represent pastel colors, and so the Saturation
// value could be much larger than 1 for colors it can't represent.
func (col Color) HPLuv() (h, s, l float64) {
	return LuvLChToHPLuv(col.LuvLChWhiteRef(hSLuvD65))
}

// DistanceHSLuv calculates Euclidean distance in the HSLuv colorspace. No idea
// how useful this is.
//
// The Hue value is divided by 100 before the calculation, so that H, S, and L
// have the same relative ranges.
func (c1 Color) DistanceHSLuv(c2 Color) float64 {
	h1, s1, l1 := c1.HSLuv()
	h2, s2, l2 := c2.HSLuv()
	return math.Sqrt(sq((h1-h2)/100.0) + sq(s1-s2) + sq(l1-l2))
}

// DistanceHPLuv calculates Euclidean distance in the HPLuv colorspace. No idea
// how useful this is.
//
// The Hue value is divided by 100 before the calculation, so that H, S, and L
// have the same relative ranges.
func (c1 Color) DistanceHPLuv(c2 Color) float64 {
	h1, s1, l1 := c1.HPLuv()
	h2, s2, l2 := c2.HPLuv()
	return math.Sqrt(sq((h1-h2)/100.0) + sq(s1-s2) + sq(l1-l2))
}

var m = [3][3]float64{
	{3.2409699419045214, -1.5373831775700935, -0.49861076029300328},
	{-0.96924363628087983, 1.8759675015077207, 0.041555057407175613},
	{0.055630079696993609, -0.20397695888897657, 1.0569715142428786},
}

const kappa = 903.2962962962963
const epsilon = 0.0088564516790356308

func maxChromaForLH(l, h float64) float64 {
	hRad := h / 360.0 * math.Pi * 2.0
	minLength := math.MaxFloat64
	for _, line := range getBounds(l) {
		length := lengthOfRayUntilIntersect(hRad, line[0], line[1])
		if length > 0.0 && length < minLength {
			minLength = length
		}
	}
	return minLength
}

func getBounds(l float64) [6][2]float64 {
	var sub2 float64
	var ret [6][2]float64
	sub1 := math.Pow(l+16.0, 3.0) / 1560896.0
	if sub1 > epsilon {
		sub2 = sub1
	} else {
		sub2 = l / kappa
	}
	for i := range m {
		for k := 0; k < 2; k++ {
			top1 := (284517.0*m[i][0] - 94839.0*m[i][2]) * sub2
			top2 := (838422.0*m[i][2]+769860.0*m[i][1]+731718.0*m[i][0])*l*sub2 - 769860.0*float64(k)*l
			bottom := (632260.0*m[i][2]-126452.0*m[i][1])*sub2 + 126452.0*float64(k)
			ret[i*2+k][0] = top1 / bottom
			ret[i*2+k][1] = top2 / bottom
		}
	}
	return ret
}

func lengthOfRayUntilIntersect(theta, x, y float64) (length float64) {
	length = y / (math.Sin(theta) - x*math.Cos(theta))
	return
}

func maxSafeChromaForL(l float64) float64 {
	minLength := math.MaxFloat64
	for _, line := range getBounds(l) {
		m1 := line[0]
		b1 := line[1]
		x := intersectLineLine(m1, b1, -1.0/m1, 0.0)
		dist := distanceFromPole(x, b1+x*m1)
		if dist < minLength {
			minLength = dist
		}
	}
	return minLength
}

func intersectLineLine(x1, y1, x2, y2 float64) float64 {
	return (y1 - y2) / (x2 - x1)
}

func distanceFromPole(x, y float64) float64 {
	return math.Sqrt(math.Pow(x, 2.0) + math.Pow(y, 2.0))
}
