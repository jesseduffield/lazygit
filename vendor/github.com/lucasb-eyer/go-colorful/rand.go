package colorful

import "math/rand"

type RandInterface interface {
	Float64() float64
	Intn(n int) int
}

type defaultGlobalRand struct{}

func (df defaultGlobalRand) Float64() float64 {
	return rand.Float64()
}

func (df defaultGlobalRand) Intn(n int) int {
	return rand.Intn(n)
}

func getDefaultGlobalRand() RandInterface {
	return defaultGlobalRand{}
}
