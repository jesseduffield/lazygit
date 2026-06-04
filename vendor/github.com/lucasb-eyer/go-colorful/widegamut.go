package colorful

import "math"

// Wide-gamut RGB color spaces from CSS Color Level 4.
// https://www.w3.org/TR/css-color-4/#color-conversion-code

/// Bradford ///
////////////////
// Bradford chromatic adaptation between D50 and D65 illuminants.

func D50ToD65(x, y, z float64) (xo, yo, zo float64) {
	xo = 0.9555766*x - 0.0230393*y + 0.0631636*z
	yo = -0.0282895*x + 1.0099416*y + 0.0210077*z
	zo = 0.0122982*x - 0.0204830*y + 1.3299098*z
	return
}

func D65ToD50(x, y, z float64) (xo, yo, zo float64) {
	xo = 1.0479298208405488*x + 0.022946793341019088*y - 0.05019222954313557*z
	yo = 0.029627815688159344*x + 0.990434484573249*y - 0.01707382502938514*z
	zo = -0.009243058152591178*x + 0.015055144896577895*y + 0.7518742899580008*z
	return
}

/// XYZ D50 ///
///////////////

func XyzD50(x, y, z float64) Color {
	return Xyz(D50ToD65(x, y, z))
}

func (col Color) XyzD50() (x, y, z float64) {
	return D65ToD50(col.Xyz())
}

/// Display P3 ///
//////////////////
// Uses the sRGB transfer function with DCI-P3 primaries.

func DisplayP3ToLinearRgb(r, g, b float64) (rl, gl, bl float64) {
	rl = linearize(r)
	gl = linearize(g)
	bl = linearize(b)
	return
}

func LinearDisplayP3ToXyz(r, g, b float64) (x, y, z float64) {
	x = 0.4865709486482162*r + 0.26566769316909306*g + 0.1982172852343625*b
	y = 0.2289745640697488*r + 0.6917385218365064*g + 0.079286914093745*b
	z = 0.04511338185890264*g + 1.043944368900976*b
	return
}

func XyzToLinearDisplayP3(x, y, z float64) (r, g, b float64) {
	r = 2.493496911941425*x - 0.9313836179191239*y - 0.40271078445071684*z
	g = -0.8294889695615747*x + 1.7626640603183463*y + 0.023624685841943577*z
	b = 0.035845830243784335*x - 0.07617238926804182*y + 0.9568845240076872*z
	return
}

func DisplayP3(r, g, b float64) Color {
	rl, gl, bl := DisplayP3ToLinearRgb(r, g, b)
	x, y, z := LinearDisplayP3ToXyz(rl, gl, bl)
	return Xyz(x, y, z)
}

func (col Color) DisplayP3() (r, g, b float64) {
	x, y, z := col.Xyz()
	rl, gl, bl := XyzToLinearDisplayP3(x, y, z)
	r = delinearize(rl)
	g = delinearize(gl)
	b = delinearize(bl)
	return
}

// BlendDisplayP3 blends two colors in the Display P3 color-space.
// t == 0 results in c1, t == 1 results in c2
func (c1 Color) BlendDisplayP3(c2 Color, t float64) Color {
	r1, g1, b1 := c1.DisplayP3()
	r2, g2, b2 := c2.DisplayP3()
	return DisplayP3(
		r1+t*(r2-r1),
		g1+t*(g2-g1),
		b1+t*(b2-b1))
}

/// A98 RGB ///
///////////////
// Adobe RGB (1998) color space.

func linearizeA98(v float64) float64 {
	sign := 1.0
	if v < 0 {
		sign = -1.0
		v = -v
	}
	return sign * math.Pow(v, 563.0/256.0)
}

func delinearizeA98(v float64) float64 {
	sign := 1.0
	if v < 0 {
		sign = -1.0
		v = -v
	}
	return sign * math.Pow(v, 256.0/563.0)
}

func A98RgbToLinearRgb(r, g, b float64) (rl, gl, bl float64) {
	rl = linearizeA98(r)
	gl = linearizeA98(g)
	bl = linearizeA98(b)
	return
}

func LinearA98RgbToXyz(r, g, b float64) (x, y, z float64) {
	x = 0.5766690429101305*r + 0.1855582379065463*g + 0.1882286462349947*b
	y = 0.29734497525053605*r + 0.6273635662554661*g + 0.07529145849399788*b
	z = 0.02703136138641234*r + 0.07068885253582723*g + 0.9913375368376388*b
	return
}

func XyzToLinearA98Rgb(x, y, z float64) (r, g, b float64) {
	r = 2.0415879038107327*x - 0.5650069742788597*y - 0.34473135077832956*z
	g = -0.9692436362808795*x + 1.8759675015077202*y + 0.04155505740717559*z
	b = 0.013444280632031142*x - 0.11836239223101838*y + 1.0151749943912054*z
	return
}

func A98Rgb(r, g, b float64) Color {
	rl, gl, bl := A98RgbToLinearRgb(r, g, b)
	x, y, z := LinearA98RgbToXyz(rl, gl, bl)
	return Xyz(x, y, z)
}

func (col Color) A98Rgb() (r, g, b float64) {
	x, y, z := col.Xyz()
	rl, gl, bl := XyzToLinearA98Rgb(x, y, z)
	r = delinearizeA98(rl)
	g = delinearizeA98(gl)
	b = delinearizeA98(bl)
	return
}

// BlendA98Rgb blends two colors in the A98 RGB color-space.
// t == 0 results in c1, t == 1 results in c2
func (c1 Color) BlendA98Rgb(c2 Color, t float64) Color {
	r1, g1, b1 := c1.A98Rgb()
	r2, g2, b2 := c2.A98Rgb()
	return A98Rgb(
		r1+t*(r2-r1),
		g1+t*(g2-g1),
		b1+t*(b2-b1))
}

/// ProPhoto RGB ///
////////////////////
// ProPhoto RGB (ROMM RGB) uses D50 illuminant.

func linearizeProPhoto(v float64) float64 {
	if v <= 16.0/512.0 {
		return v / 16.0
	}
	return math.Pow(v, 1.8)
}

func delinearizeProPhoto(v float64) float64 {
	if v < 1.0/512.0 {
		return 16.0 * v
	}
	return math.Pow(v, 1.0/1.8)
}

func ProPhotoRgbToLinearRgb(r, g, b float64) (rl, gl, bl float64) {
	rl = linearizeProPhoto(r)
	gl = linearizeProPhoto(g)
	bl = linearizeProPhoto(b)
	return
}

func LinearProPhotoRgbToXyzD50(r, g, b float64) (x, y, z float64) {
	x = 0.7977604896723027*r + 0.13518583717574031*g + 0.0313493495815248*b
	y = 0.2880711282292934*r + 0.7118432178101014*g + 0.00008565396060525902*b
	z = 0.8251046025104602 * b
	return
}

func XyzD50ToLinearProPhotoRgb(x, y, z float64) (r, g, b float64) {
	r = 1.3457989731028281*x - 0.25558010007997534*y - 0.05110628506753401*z
	g = -0.5446224939028347*x + 1.5082327413132781*y + 0.02053603239147973*z
	b = 1.2119675456389454 * z
	return
}

func ProPhotoRgb(r, g, b float64) Color {
	rl, gl, bl := ProPhotoRgbToLinearRgb(r, g, b)
	x, y, z := LinearProPhotoRgbToXyzD50(rl, gl, bl)
	return XyzD50(x, y, z)
}

func (col Color) ProPhotoRgb() (r, g, b float64) {
	x, y, z := col.XyzD50()
	rl, gl, bl := XyzD50ToLinearProPhotoRgb(x, y, z)
	r = delinearizeProPhoto(rl)
	g = delinearizeProPhoto(gl)
	b = delinearizeProPhoto(bl)
	return
}

// BlendProPhotoRgb blends two colors in the ProPhoto RGB color-space.
// t == 0 results in c1, t == 1 results in c2
func (c1 Color) BlendProPhotoRgb(c2 Color, t float64) Color {
	r1, g1, b1 := c1.ProPhotoRgb()
	r2, g2, b2 := c2.ProPhotoRgb()
	return ProPhotoRgb(
		r1+t*(r2-r1),
		g1+t*(g2-g1),
		b1+t*(b2-b1))
}

/// Rec. 2020 ///
/////////////////
// ITU-R BT.2020 color space.

const (
	rec2020Alpha = 1.09929682680944
	rec2020Beta  = 0.018053968510807
)

func linearizeRec2020(v float64) float64 {
	if v < rec2020Beta*4.5 {
		return v / 4.5
	}
	return math.Pow((v+rec2020Alpha-1)/rec2020Alpha, 1.0/0.45)
}

func delinearizeRec2020(v float64) float64 {
	if v < rec2020Beta {
		return 4.5 * v
	}
	return rec2020Alpha*math.Pow(v, 0.45) - (rec2020Alpha - 1)
}

func Rec2020ToLinearRgb(r, g, b float64) (rl, gl, bl float64) {
	rl = linearizeRec2020(r)
	gl = linearizeRec2020(g)
	bl = linearizeRec2020(b)
	return
}

func LinearRec2020ToXyz(r, g, b float64) (x, y, z float64) {
	x = 0.6369580483012914*r + 0.14461690358620832*g + 0.1688809751641721*b
	y = 0.2627002120112671*r + 0.6779980715188708*g + 0.05930171646986196*b
	z = 0.028072693049087428*g + 1.0609850577107909*b
	return
}

func XyzToLinearRec2020(x, y, z float64) (r, g, b float64) {
	r = 1.7166511879712674*x - 0.35567078377639233*y - 0.25336628137365974*z
	g = -0.666684351832489*x + 1.616481236634939*y + 0.0157685458139402*z
	b = 0.017639857445310783*x - 0.042770613257808524*y + 0.9421031212354738*z
	return
}

func Rec2020(r, g, b float64) Color {
	rl, gl, bl := Rec2020ToLinearRgb(r, g, b)
	x, y, z := LinearRec2020ToXyz(rl, gl, bl)
	return Xyz(x, y, z)
}

func (col Color) Rec2020() (r, g, b float64) {
	x, y, z := col.Xyz()
	rl, gl, bl := XyzToLinearRec2020(x, y, z)
	r = delinearizeRec2020(rl)
	g = delinearizeRec2020(gl)
	b = delinearizeRec2020(bl)
	return
}

// BlendRec2020 blends two colors in the Rec. 2020 color-space.
// t == 0 results in c1, t == 1 results in c2
func (c1 Color) BlendRec2020(c2 Color, t float64) Color {
	r1, g1, b1 := c1.Rec2020()
	r2, g2, b2 := c2.Rec2020()
	return Rec2020(
		r1+t*(r2-r1),
		g1+t*(g2-g1),
		b1+t*(b2-b1))
}
