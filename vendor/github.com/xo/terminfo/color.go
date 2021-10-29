package terminfo

import (
	"os"
	"strconv"
	"strings"
)

// ColorLevel is the color level supported by a terminal.
type ColorLevel uint

// ColorLevel values.
const (
	ColorLevelNone ColorLevel = iota
	ColorLevelBasic
	ColorLevelHundreds
	ColorLevelMillions
)

// String satisfies the Stringer interface.
func (c ColorLevel) String() string {
	switch c {
	case ColorLevelBasic:
		return "basic"
	case ColorLevelHundreds:
		return "hundreds"
	case ColorLevelMillions:
		return "millions"
	}
	return "none"
}

// ChromaFormatterName returns the github.com/alecthomas/chroma compatible
// formatter name for the color level.
func (c ColorLevel) ChromaFormatterName() string {
	switch c {
	case ColorLevelBasic:
		return "terminal"
	case ColorLevelHundreds:
		return "terminal256"
	case ColorLevelMillions:
		return "terminal16m"
	}
	return "noop"
}

// ColorLevelFromEnv returns the color level COLORTERM, FORCE_COLOR,
// TERM_PROGRAM, or determined from the TERM environment variable.
func ColorLevelFromEnv() (ColorLevel, error) {
	// check for overriding environment variables
	colorTerm, termProg, forceColor := os.Getenv("COLORTERM"), os.Getenv("TERM_PROGRAM"), os.Getenv("FORCE_COLOR")
	switch {
	case strings.Contains(colorTerm, "truecolor") || strings.Contains(colorTerm, "24bit") || termProg == "Hyper":
		return ColorLevelMillions, nil
	case colorTerm != "" || forceColor != "":
		return ColorLevelBasic, nil
	case termProg == "Apple_Terminal":
		return ColorLevelHundreds, nil
	case termProg == "iTerm.app":
		ver := os.Getenv("TERM_PROGRAM_VERSION")
		if ver == "" {
			return ColorLevelHundreds, nil
		}
		i, err := strconv.Atoi(strings.Split(ver, ".")[0])
		if err != nil {
			return ColorLevelNone, ErrInvalidTermProgramVersion
		}
		if i == 3 {
			return ColorLevelMillions, nil
		}
		return ColorLevelHundreds, nil
	}

	// otherwise determine from TERM's max_colors capability
	if term := os.Getenv("TERM"); term != "" {
		ti, err := Load(term)
		if err != nil {
			return ColorLevelNone, err
		}

		v, ok := ti.Nums[MaxColors]
		switch {
		case !ok || v <= 16:
			return ColorLevelNone, nil
		case ok && v >= 256:
			return ColorLevelHundreds, nil
		}
	}

	return ColorLevelBasic, nil
}
