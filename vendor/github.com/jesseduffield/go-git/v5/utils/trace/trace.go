package trace

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
)

var (
	// logger is the logger to use for tracing.
	logger = newLogger()

	// current is the targets that are enabled for tracing.
	current atomic.Int32
)

func newLogger() *log.Logger {
	return log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

// Target is a tracing target.
type Target int32

const (
	// General traces general operations.
	General Target = 1 << iota

	// Packet traces git packets.
	Packet
)

// SetTarget sets the tracing targets.
func SetTarget(target Target) {
	current.Store(int32(target))
}

// SetLogger sets the logger to use for tracing.
func SetLogger(l *log.Logger) {
	logger = l
}

// Print prints the given message only if the target is enabled.
func (t Target) Print(args ...interface{}) {
	if int32(t)&current.Load() != 0 {
		logger.Output(2, fmt.Sprint(args...)) // nolint: errcheck
	}
}

// Printf prints the given message only if the target is enabled.
func (t Target) Printf(format string, args ...interface{}) {
	if int32(t)&current.Load() != 0 {
		logger.Output(2, fmt.Sprintf(format, args...)) // nolint: errcheck
	}
}
