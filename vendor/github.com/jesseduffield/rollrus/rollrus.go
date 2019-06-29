// Package rollrus combines github.com/jesseduffield/roll with github.com/sirupsen/logrus
// via logrus.Hook mechanism, so that whenever logrus' logger.Error/f(),
// logger.Fatal/f() or logger.Panic/f() are used the messages are
// intercepted and sent to rollbar.
//
// Using SetupLogging should suffice for basic use cases that use the logrus
// singleton logger.
//
// More custom uses are supported by creating a new Hook with NewHook and
// registering that hook with the logrus Logger of choice.
//
// The levels can be customized with the WithLevels OptionFunc.
//
// Specific errors can be ignored with the WithIgnoredErrors OptionFunc. This is
// useful for ignoring errors such as context.Canceled.
//
// See the Examples in the tests for more usage.
package rollrus

import (
	"fmt"
	"os"
	"time"

	"github.com/jesseduffield/roll"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var defaultTriggerLevels = []logrus.Level{
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

// Hook is a wrapper for the rollbar Client and is usable as a logrus.Hook.
type Hook struct {
	roll.Client
	triggers        []logrus.Level
	ignoredErrors   []error
	ignoreErrorFunc func(error) bool
	ignoreFunc      func(error, map[string]string) bool

	// only used for tests to verify whether or not a report happened.
	reported bool
}

// OptionFunc that can be passed to NewHook.
type OptionFunc func(*Hook)

// wellKnownErrorFields are the names of the fields to be checked for values of
// type `error`, in priority order.
var wellKnownErrorFields = []string{
	logrus.ErrorKey, "err",
}

// WithLevels is an OptionFunc that customizes the log.Levels the hook will
// report on.
func WithLevels(levels ...logrus.Level) OptionFunc {
	return func(h *Hook) {
		h.triggers = levels
	}
}

// WithMinLevel is an OptionFunc that customizes the log.Levels the hook will
// report on by selecting all levels more severe than the one provided.
func WithMinLevel(level logrus.Level) OptionFunc {
	var levels []logrus.Level
	for _, l := range logrus.AllLevels {
		if l <= level {
			levels = append(levels, l)
		}
	}

	return func(h *Hook) {
		h.triggers = levels
	}
}

// WithIgnoredErrors is an OptionFunc that whitelists certain errors to prevent
// them from firing. See https://golang.org/ref/spec#Comparison_operators
func WithIgnoredErrors(errors ...error) OptionFunc {
	return func(h *Hook) {
		h.ignoredErrors = append(h.ignoredErrors, errors...)
	}
}

// WithIgnoreErrorFunc is an OptionFunc that receives the error that is about
// to be logged and returns true/false if it wants to fire a rollbar alert for.
func WithIgnoreErrorFunc(fn func(error) bool) OptionFunc {
	return func(h *Hook) {
		h.ignoreErrorFunc = fn
	}
}

// WithIgnoreFunc is an OptionFunc that receives the error and custom fields that are about
// to be logged and returns true/false if it wants to fire a rollbar alert for.
func WithIgnoreFunc(fn func(err error, fields map[string]string) bool) OptionFunc {
	return func(h *Hook) {
		h.ignoreFunc = fn
	}
}

// NewHook creates a hook that is intended for use with your own logrus.Logger
// instance. Uses the defualt report levels defined in wellKnownErrorFields.
func NewHook(token string, env string, opts ...OptionFunc) *Hook {
	h := NewHookForLevels(token, env, defaultTriggerLevels)

	for _, o := range opts {
		o(h)
	}

	return h
}

// NewHookForLevels provided by the caller. Otherwise works like NewHook.
func NewHookForLevels(token string, env string, levels []logrus.Level) *Hook {
	return &Hook{
		Client:          roll.New(token, env),
		triggers:        levels,
		ignoredErrors:   make([]error, 0),
		ignoreErrorFunc: func(error) bool { return false },
		ignoreFunc:      func(error, map[string]string) bool { return false },
	}
}

// SetupLogging for use on Heroku. If token is not an empty string a rollbar
// hook is added with the environment set to env. The log formatter is set to a
// TextFormatter with timestamps disabled.
func SetupLogging(token, env string) {
	setupLogging(token, env, defaultTriggerLevels)
}

// SetupLoggingForLevels works like SetupLogging, but allows you to
// set the levels on which to trigger this hook.
func SetupLoggingForLevels(token, env string, levels []logrus.Level) {
	setupLogging(token, env, levels)
}

func setupLogging(token, env string, levels []logrus.Level) {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})

	if token != "" {
		logrus.AddHook(NewHookForLevels(token, env, levels))
	}
}

// ReportPanic attempts to report the panic to rollbar using the provided
// client and then re-panic. If it can't report the panic it will print an
// error to stderr.
func (r *Hook) ReportPanic() {
	if p := recover(); p != nil {
		if _, err := r.Client.Critical(fmt.Errorf("panic: %q", p), nil); err != nil {
			fmt.Fprintf(os.Stderr, "reporting_panic=false err=%q\n", err)
		}
		panic(p)
	}
}

// ReportPanic attempts to report the panic to rollbar if the token is set
func ReportPanic(token, env string) {
	if token != "" {
		h := &Hook{Client: roll.New(token, env)}
		h.ReportPanic()
	}
}

// Levels returns the logrus log.Levels that this hook handles
func (r *Hook) Levels() []logrus.Level {
	if r.triggers == nil {
		return defaultTriggerLevels
	}
	return r.triggers
}

// Fire the hook. This is called by Logrus for entries that match the levels
// returned by Levels().
func (r *Hook) Fire(entry *logrus.Entry) error {
	trace, cause := extractError(entry)
	for _, ie := range r.ignoredErrors {
		if ie == cause {
			return nil
		}
	}

	if r.ignoreErrorFunc(cause) {
		return nil
	}

	m := convertFields(entry.Data)
	if _, exists := m["time"]; !exists {
		m["time"] = entry.Time.Format(time.RFC3339)
	}

	if r.ignoreFunc(cause, m) {
		return nil
	}

	return r.report(entry, cause, m, trace)
}

func (r *Hook) report(entry *logrus.Entry, cause error, m map[string]string, trace []uintptr) (err error) {
	hasTrace := len(trace) > 0
	level := entry.Level

	r.reported = true

	switch {
	case hasTrace && level == logrus.FatalLevel:
		_, err = r.Client.CriticalStack(cause, trace, m)
	case hasTrace && level == logrus.PanicLevel:
		_, err = r.Client.CriticalStack(cause, trace, m)
	case hasTrace && level == logrus.ErrorLevel:
		_, err = r.Client.ErrorStack(cause, trace, m)
	case hasTrace && level == logrus.WarnLevel:
		_, err = r.Client.WarningStack(cause, trace, m)
	case level == logrus.FatalLevel || level == logrus.PanicLevel:
		_, err = r.Client.Critical(cause, m)
	case level == logrus.ErrorLevel:
		_, err = r.Client.Error(cause, m)
	case level == logrus.WarnLevel:
		_, err = r.Client.Warning(cause, m)
	case level == logrus.InfoLevel:
		_, err = r.Client.Info(entry.Message, m)
	case level == logrus.DebugLevel:
		_, err = r.Client.Debug(entry.Message, m)
	}
	return err
}

// convertFields converts from log.Fields to map[string]string so that we can
// report extra fields to Rollbar
func convertFields(fields logrus.Fields) map[string]string {
	m := make(map[string]string)
	for k, v := range fields {
		switch t := v.(type) {
		case time.Time:
			m[k] = t.Format(time.RFC3339)
		default:
			if s, ok := v.(fmt.Stringer); ok {
				m[k] = s.String()
			} else {
				m[k] = fmt.Sprintf("%+v", t)
			}
		}
	}

	return m
}

// extractError attempts to extract an error from a well known field, err or error
func extractError(entry *logrus.Entry) ([]uintptr, error) {
	var trace []uintptr
	fields := entry.Data

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	for _, f := range wellKnownErrorFields {
		e, ok := fields[f]
		if !ok {
			continue
		}
		err, ok := e.(error)
		if !ok {
			continue
		}

		cause := errors.Cause(err)
		tracer, ok := err.(stackTracer)
		if ok {
			return copyStackTrace(tracer.StackTrace()), cause
		}
		return trace, cause
	}

	// when no error found, default to the logged message.
	return trace, fmt.Errorf(entry.Message)
}

func copyStackTrace(trace errors.StackTrace) (out []uintptr) {
	for _, frame := range trace {
		out = append(out, uintptr(frame))
	}
	return
}
