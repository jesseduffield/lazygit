package colorful

// Uses the HSV color space to generate colors with similar S,V but distributed
// evenly along their Hue. This is fast but not always pretty.
// If you've got time to spare, use Lab (the non-fast below).
func FastWarmPaletteWithRand(colorsCount int, rand RandInterface) (colors []Color) {
	colors = make([]Color, colorsCount)

	for i := 0; i < colorsCount; i++ {
		colors[i] = Hsv(float64(i)*(360.0/float64(colorsCount)), 0.55+rand.Float64()*0.2, 0.35+rand.Float64()*0.2)
	}
	return
}

func FastWarmPalette(colorsCount int) (colors []Color) {
	return FastWarmPaletteWithRand(colorsCount, getDefaultGlobalRand())
}

func WarmPaletteWithRand(colorsCount int, rand RandInterface) ([]Color, error) {
	warmy := func(l, a, b float64) bool {
		_, c, _ := LabToHcl(l, a, b)
		return 0.1 <= c && c <= 0.4 && 0.2 <= l && l <= 0.5
	}
	return SoftPaletteExWithRand(colorsCount, SoftPaletteSettings{warmy, 50, true}, rand)
}

func WarmPalette(colorsCount int) ([]Color, error) {
	return WarmPaletteWithRand(colorsCount, getDefaultGlobalRand())
}
