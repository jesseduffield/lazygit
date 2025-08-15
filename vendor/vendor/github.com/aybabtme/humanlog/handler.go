package humanlog

import (
	"time"

	"github.com/fatih/color"
	"github.com/kr/logfmt"
)

// Handler can recognize it's log lines, parse them and prettify them.
type Handler interface {
	CanHandle(line []byte) bool
	Prettify(skipUnchanged bool) []byte
	logfmt.Handler
}

var DefaultOptions = &HandlerOptions{
	SortLongest:    true,
	SkipUnchanged:  true,
	Truncates:      true,
	LightBg:        false,
	TruncateLength: 15,
	TimeFormat:     time.Stamp,

	KeyColor:              color.New(color.FgGreen),
	ValColor:              color.New(color.FgHiWhite),
	TimeLightBgColor:      color.New(color.FgBlack),
	TimeDarkBgColor:       color.New(color.FgWhite),
	MsgLightBgColor:       color.New(color.FgBlack),
	MsgAbsentLightBgColor: color.New(color.FgHiBlack),
	MsgDarkBgColor:        color.New(color.FgHiWhite),
	MsgAbsentDarkBgColor:  color.New(color.FgWhite),
	DebugLevelColor:       color.New(color.FgMagenta),
	InfoLevelColor:        color.New(color.FgCyan),
	WarnLevelColor:        color.New(color.FgYellow),
	ErrorLevelColor:       color.New(color.FgRed),
	PanicLevelColor:       color.New(color.BgRed),
	FatalLevelColor:       color.New(color.BgHiRed, color.FgHiWhite),
	UnknownLevelColor:     color.New(color.FgMagenta),
}

type HandlerOptions struct {
	Skip           map[string]struct{}
	Keep           map[string]struct{}
	SortLongest    bool
	SkipUnchanged  bool
	Truncates      bool
	LightBg        bool
	TruncateLength int
	TimeFormat     string

	KeyColor              *color.Color
	ValColor              *color.Color
	TimeLightBgColor      *color.Color
	TimeDarkBgColor       *color.Color
	MsgLightBgColor       *color.Color
	MsgAbsentLightBgColor *color.Color
	MsgDarkBgColor        *color.Color
	MsgAbsentDarkBgColor  *color.Color
	DebugLevelColor       *color.Color
	InfoLevelColor        *color.Color
	WarnLevelColor        *color.Color
	ErrorLevelColor       *color.Color
	PanicLevelColor       *color.Color
	FatalLevelColor       *color.Color
	UnknownLevelColor     *color.Color
}

func (h *HandlerOptions) shouldShowKey(key string) bool {
	if len(h.Keep) != 0 {
		if _, keep := h.Keep[key]; keep {
			return true
		}
	}
	if len(h.Skip) != 0 {
		if _, skip := h.Skip[key]; skip {
			return false
		}
	}
	return true
}

func (h *HandlerOptions) shouldShowUnchanged(key string) bool {
	if len(h.Keep) != 0 {
		if _, keep := h.Keep[key]; keep {
			return true
		}
	}
	return false
}

func (h *HandlerOptions) SetSkip(skip []string) {
	if h.Skip == nil {
		h.Skip = make(map[string]struct{})
	}
	for _, key := range skip {
		h.Skip[key] = struct{}{}
	}
}

func (h *HandlerOptions) SetKeep(keep []string) {
	if h.Keep == nil {
		h.Keep = make(map[string]struct{})
	}
	for _, key := range keep {
		h.Keep[key] = struct{}{}
	}
}
