package humanlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
)

// JSONHandler can handle logs emmited by logrus.TextFormatter loggers.
type JSONHandler struct {
	buf     *bytes.Buffer
	out     *tabwriter.Writer
	truncKV int

	Opts *HandlerOptions

	Level   string
	Time    time.Time
	Message string
	Fields  map[string]string

	last map[string]string
}

func checkEachUntilFound(fieldList []string, found func(string) bool) bool {
	for _, field := range fieldList {
		if found(field) {
			return true
		}
	}
	return false
}

// supportedTimeFields enumerates supported timestamp field names
var supportedTimeFields = []string{"time", "ts", "@timestamp", "timestamp"}

// supportedMessageFields enumarates supported Message field names
var supportedMessageFields = []string{"message", "msg"}

// supportedLevelFields enumarates supported level field names
var supportedLevelFields = []string{"level", "lvl", "loglevel"}

func (h *JSONHandler) clear() {
	h.Level = ""
	h.Time = time.Time{}
	h.Message = ""
	h.last = h.Fields
	h.Fields = make(map[string]string)
	if h.buf != nil {
		h.buf.Reset()
	}
}

// TryHandle tells if this line was handled by this handler.
func (h *JSONHandler) TryHandle(d []byte) bool {
	if !h.UnmarshalJSON(d) {
		h.clear()
		return false
	}
	return true
}

// UnmarshalJSON sets the fields of the handler.
func (h *JSONHandler) UnmarshalJSON(data []byte) bool {
	raw := make(map[string]interface{})
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return false
	}

	checkEachUntilFound(supportedTimeFields, func(field string) bool {
		time, ok := tryParseTime(raw[field])
		if ok {
			h.Time = time
			delete(raw, field)
		}
		return ok
	})

	checkEachUntilFound(supportedMessageFields, func(field string) bool {
		msg, ok := raw[field].(string)
		if ok {
			h.Message = msg
			delete(raw, field)
		}
		return ok
	})

	checkEachUntilFound(supportedLevelFields, func(field string) bool {
		lvl, ok := raw[field]
		if !ok {
			return false
		}
		if strLvl, ok := lvl.(string); ok {
			h.Level = strLvl
		} else if flLvl, ok := lvl.(float64); ok {
			h.Level = convertBunyanLogLevel(flLvl)
		} else {
			h.Level = "???"
		}
		delete(raw, field)
		return true
	})

	if h.Fields == nil {
		h.Fields = make(map[string]string)
	}

	for key, val := range raw {
		switch v := val.(type) {
		case float64:
			if v-math.Floor(v) < 0.000001 && v < 1e9 {
				// looks like an integer that's not too large
				h.Fields[key] = fmt.Sprintf("%d", int(v))
			} else {
				h.Fields[key] = fmt.Sprintf("%g", v)
			}
		case string:
			h.Fields[key] = fmt.Sprintf("%q", v)
		default:
			h.Fields[key] = fmt.Sprintf("%v", v)
		}
	}

	return true
}
func (h *JSONHandler) setField(key, val []byte) {
	if h.Fields == nil {
		h.Fields = make(map[string]string)
	}
	h.Fields[string(key)] = string(val)
}

// Prettify the output in a logrus like fashion.
func (h *JSONHandler) Prettify(skipUnchanged bool) []byte {
	defer h.clear()
	if h.out == nil {
		if h.Opts == nil {
			h.Opts = DefaultOptions
		}
		h.buf = bytes.NewBuffer(nil)
		h.out = tabwriter.NewWriter(h.buf, 0, 1, 0, '\t', 0)
	}

	var (
		msgColor       *color.Color
		msgAbsentColor *color.Color
	)
	if h.Opts.LightBg {
		msgColor = h.Opts.MsgLightBgColor
		msgAbsentColor = h.Opts.MsgAbsentLightBgColor
	} else {
		msgColor = h.Opts.MsgDarkBgColor
		msgAbsentColor = h.Opts.MsgAbsentDarkBgColor
	}
	msgColor = color.New(color.FgHiWhite)
	msgAbsentColor = color.New(color.FgHiWhite)

	var msg string
	if h.Message == "" {
		msg = msgAbsentColor.Sprint("<no msg>")
	} else {
		msg = msgColor.Sprint(h.Message)
	}

	lvl := strings.ToUpper(h.Level)[:imin(4, len(h.Level))]
	var level string
	switch h.Level {
	case "debug":
		level = h.Opts.DebugLevelColor.Sprint(lvl)
	case "info":
		level = h.Opts.InfoLevelColor.Sprint(lvl)
	case "warn", "warning":
		level = h.Opts.WarnLevelColor.Sprint(lvl)
	case "error":
		level = h.Opts.ErrorLevelColor.Sprint(lvl)
	case "fatal", "panic":
		level = h.Opts.FatalLevelColor.Sprint(lvl)
	default:
		level = h.Opts.UnknownLevelColor.Sprint(lvl)
	}

	var timeColor *color.Color
	if h.Opts.LightBg {
		timeColor = h.Opts.TimeLightBgColor
	} else {
		timeColor = h.Opts.TimeDarkBgColor
	}
	_, _ = fmt.Fprintf(h.out, "%s |%s| %s\t %s",
		timeColor.Sprint(h.Time.Format(h.Opts.TimeFormat)),
		level,
		msg,
		strings.Join(h.joinKVs(skipUnchanged, "="), "\t "),
	)

	_ = h.out.Flush()

	return h.buf.Bytes()
}

func (h *JSONHandler) joinKVs(skipUnchanged bool, sep string) []string {

	kv := make([]string, 0, len(h.Fields))
	for k, v := range h.Fields {
		if !h.Opts.shouldShowKey(k) {
			continue
		}

		if skipUnchanged {
			if lastV, ok := h.last[k]; ok && lastV == v && !h.Opts.shouldShowUnchanged(k) {
				continue
			}
		}
		kstr := h.Opts.KeyColor.Sprint(k)

		var vstr string
		if h.Opts.Truncates && len(v) > h.Opts.TruncateLength {
			vstr = v[:h.Opts.TruncateLength] + "..."
		} else {
			vstr = v
		}
		vstr = h.Opts.ValColor.Sprint(vstr)
		kv = append(kv, kstr+sep+vstr)
	}

	sort.Strings(kv)

	if h.Opts.SortLongest {
		sort.Stable(byLongest(kv))
	}

	return kv
}

// convertBunyanLogLevel returns a human readable log level given a numerical bunyan level
// https://github.com/trentm/node-bunyan#levels
func convertBunyanLogLevel(level float64) string {
	switch level {
	case 10:
		return "trace"
	case 20:
		return "debug"
	case 30:
		return "info"
	case 40:
		return "warn"
	case 50:
		return "error"
	case 60:
		return "fatal"
	default:
		return "???"
	}
}
