// Various ways to generate single random colors

package colorful

// Creates a random dark, "warm" color through a restricted HSV space.
func FastWarmColorWithRand(rand RandInterface) Color {
	return Hsv(
		rand.Float64()*360.0,
		0.5+rand.Float64()*0.3,
		0.3+rand.Float64()*0.3)
}

func FastWarmColor() Color {
	return FastWarmColorWithRand(getDefaultGlobalRand())
}

// Creates a random dark, "warm" color through restricted HCL space.
// This is slower than FastWarmColor but will likely give you colors which have
// the same "warmness" if you run it many times.
func WarmColorWithRand(rand RandInterface) (c Color) {
	for c = randomWarmWithRand(rand); !c.IsValid(); c = randomWarmWithRand(rand) {
	}
	return
}

func WarmColor() (c Color) {
	return WarmColorWithRand(getDefaultGlobalRand())
}

func randomWarmWithRand(rand RandInterface) Color {
	return Hcl(
		rand.Float64()*360.0,
		0.1+rand.Float64()*0.3,
		0.2+rand.Float64()*0.3)
}

// Creates a random bright, "pimpy" color through a restricted HSV space.
func FastHappyColorWithRand(rand RandInterface) Color {
	return Hsv(
		rand.Float64()*360.0,
		0.7+rand.Float64()*0.3,
		0.6+rand.Float64()*0.3)
}

func FastHappyColor() Color {
	return FastHappyColorWithRand(getDefaultGlobalRand())
}

// Creates a random bright, "pimpy" color through restricted HCL space.
// This is slower than FastHappyColor but will likely give you colors which
// have the same "brightness" if you run it many times.
func HappyColorWithRand(rand RandInterface) (c Color) {
	for c = randomPimpWithRand(rand); !c.IsValid(); c = randomPimpWithRand(rand) {
	}
	return
}

func HappyColor() (c Color) {
	return HappyColorWithRand(getDefaultGlobalRand())
}

func randomPimpWithRand(rand RandInterface) Color {
	return Hcl(
		rand.Float64()*360.0,
		0.5+rand.Float64()*0.3,
		0.5+rand.Float64()*0.3)
}
