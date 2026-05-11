package colorful

// Uses the HSV color space to generate colors with similar S,V but distributed
// evenly along their Hue. This is fast but not always pretty.
// If you've got time to spare, use Lab (the non-fast below).
func FastHappyPaletteWithRand(colorsCount int, rand RandInterface) (colors []Color) {
	colors = make([]Color, colorsCount)

	for i := 0; i < colorsCount; i++ {
		colors[i] = Hsv(float64(i)*(360.0/float64(colorsCount)), 0.8+rand.Float64()*0.2, 0.65+rand.Float64()*0.2)
	}
	return
}

func FastHappyPalette(colorsCount int) (colors []Color) {
	return FastHappyPaletteWithRand(colorsCount, getDefaultGlobalRand())
}

func HappyPaletteWithRand(colorsCount int, rand RandInterface) ([]Color, error) {
	pimpy := func(l, a, b float64) bool {
		_, c, _ := LabToHcl(l, a, b)
		return 0.3 <= c && 0.4 <= l && l <= 0.8
	}
	return SoftPaletteExWithRand(colorsCount, SoftPaletteSettings{pimpy, 50, true}, rand)
}

func HappyPalette(colorsCount int) ([]Color, error) {
	return HappyPaletteWithRand(colorsCount, getDefaultGlobalRand())
}
