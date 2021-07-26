package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHexColorValues(t *testing.T) {
	scenarios := []struct {
		name     string
		hexColor string
		rgb      []int32
		valid    bool
	}{
		{
			name:     "valid uppercase hex color",
			hexColor: "#00FF00",
			rgb:      []int32{0, 255, 0},
			valid:    true,
		},
		{
			name:     "valid lowercase hex color",
			hexColor: "#00ff00",
			rgb:      []int32{0, 255, 0},
			valid:    true,
		},
		{
			name:     "valid short hex color",
			hexColor: "#0bf",
			rgb:      []int32{0, 187, 255},
			valid:    true,
		},
		{
			name:     "invalid hex value",
			hexColor: "#zz00ff",
			valid:    false,
		},
		{
			name:     "invalid length hex color",
			hexColor: "#",
			valid:    false,
		},
		{
			name:     "invalid length hex color",
			hexColor: "#aaaaaaaaaaa",
			valid:    false,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			r, g, b, valid := getHexColorValues(s.hexColor)
			assert.EqualValues(t, s.valid, valid, s.hexColor)
			if valid {
				assert.EqualValues(t, s.rgb[0], r, s.hexColor)
				assert.EqualValues(t, s.rgb[1], g, s.hexColor)
				assert.EqualValues(t, s.rgb[2], b, s.hexColor)
			}
		})
	}
}
