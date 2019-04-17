package logrus_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/sirupsen/logrus"
	. "github.com/sirupsen/logrus/internal/testutils"
)

// TestReportCaller verifies that when ReportCaller is set, the 'func' field
// is added, and when it is unset it is not set or modified
// Verify that functions within the Logrus package aren't considered when
// discovering the caller.
func TestReportCallerWhenConfigured(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.ReportCaller = false
		log.Print("testNoCaller")
	}, func(fields Fields) {
		assert.Equal(t, "testNoCaller", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Equal(t, nil, fields["func"])
	})

	LogAndAssertJSON(t, func(log *Logger) {
		log.ReportCaller = true
		log.Print("testWithCaller")
	}, func(fields Fields) {
		assert.Equal(t, "testWithCaller", fields["msg"])
		assert.Equal(t, "info", fields["level"])
		assert.Equal(t,
			"github.com/sirupsen/logrus_test.TestReportCallerWhenConfigured.func3", fields[FieldKeyFunc])
	})

	LogAndAssertJSON(t, func(log *Logger) {
		log.ReportCaller = true
		log.Formatter.(*JSONFormatter).CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			return "somekindoffunc", "thisisafilename"
		}
		log.Print("testWithCallerPrettyfier")
	}, func(fields Fields) {
		assert.Equal(t, "somekindoffunc", fields[FieldKeyFunc])
		assert.Equal(t, "thisisafilename", fields[FieldKeyFile])
	})

	LogAndAssertText(t, func(log *Logger) {
		log.ReportCaller = true
		log.Formatter.(*TextFormatter).CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			return "somekindoffunc", "thisisafilename"
		}
		log.Print("testWithCallerPrettyfier")
	}, func(fields map[string]string) {
		assert.Equal(t, "somekindoffunc", fields[FieldKeyFunc])
		assert.Equal(t, "thisisafilename", fields[FieldKeyFile])
	})
}

func logSomething(t *testing.T, message string) Fields {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)
	logger.ReportCaller = true

	entry := logger.WithFields(Fields{
		"foo": "bar",
	})

	entry.Info(message)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	return fields
}

// TestReportCallerHelperDirect - verify reference when logging from a regular function
func TestReportCallerHelperDirect(t *testing.T) {
	fields := logSomething(t, "direct")

	assert.Equal(t, "direct", fields["msg"])
	assert.Equal(t, "info", fields["level"])
	assert.Regexp(t, "github.com/.*/logrus_test.logSomething", fields["func"])
}

// TestReportCallerHelperDirect - verify reference when logging from a function called via pointer
func TestReportCallerHelperViaPointer(t *testing.T) {
	fptr := logSomething
	fields := fptr(t, "via pointer")

	assert.Equal(t, "via pointer", fields["msg"])
	assert.Equal(t, "info", fields["level"])
	assert.Regexp(t, "github.com/.*/logrus_test.logSomething", fields["func"])
}

func TestPrint(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Print("test")
	}, func(fields Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})
}

func TestInfo(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test")
	}, func(fields Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "info", fields["level"])
	})
}

func TestWarn(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Warn("test")
	}, func(fields Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "warning", fields["level"])
	})
}

func TestLog(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Log(WarnLevel, "test")
	}, func(fields Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "warning", fields["level"])
	})
}

func TestInfolnShouldAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", "test")
	}, func(fields Fields) {
		assert.Equal(t, "test test", fields["msg"])
	})
}

func TestInfolnShouldAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", 10)
	}, func(fields Fields) {
		assert.Equal(t, "test 10", fields["msg"])
	})
}

func TestInfolnShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(fields Fields) {
		assert.Equal(t, "10 10", fields["msg"])
	})
}

func TestInfoShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(fields Fields) {
		assert.Equal(t, "10 10", fields["msg"])
	})
}

func TestInfoShouldNotAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", 10)
	}, func(fields Fields) {
		assert.Equal(t, "test10", fields["msg"])
	})
}

func TestInfoShouldNotAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", "test")
	}, func(fields Fields) {
		assert.Equal(t, "testtest", fields["msg"])
	})
}

func TestWithFieldsShouldAllowAssignments(t *testing.T) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	localLog := logger.WithFields(Fields{
		"key1": "value1",
	})

	localLog.WithField("key2", "value2").Info("test")
	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	assert.Equal(t, "value2", fields["key2"])
	assert.Equal(t, "value1", fields["key1"])

	buffer = bytes.Buffer{}
	fields = Fields{}
	localLog.Info("test")
	err = json.Unmarshal(buffer.Bytes(), &fields)
	assert.Nil(t, err)

	_, ok := fields["key2"]
	assert.Equal(t, false, ok)
	assert.Equal(t, "value1", fields["key1"])
}

func TestUserSuppliedFieldDoesNotOverwriteDefaults(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("msg", "hello").Info("test")
	}, func(fields Fields) {
		assert.Equal(t, "test", fields["msg"])
	})
}

func TestUserSuppliedMsgFieldHasPrefix(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("msg", "hello").Info("test")
	}, func(fields Fields) {
		assert.Equal(t, "test", fields["msg"])
		assert.Equal(t, "hello", fields["fields.msg"])
	})
}

func TestUserSuppliedTimeFieldHasPrefix(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("time", "hello").Info("test")
	}, func(fields Fields) {
		assert.Equal(t, "hello", fields["fields.time"])
	})
}

func TestUserSuppliedLevelFieldHasPrefix(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("level", 1).Info("test")
	}, func(fields Fields) {
		assert.Equal(t, "info", fields["level"])
		assert.Equal(t, 1.0, fields["fields.level"]) // JSON has floats only
	})
}

func TestDefaultFieldsAreNotPrefixed(t *testing.T) {
	LogAndAssertText(t, func(log *Logger) {
		ll := log.WithField("herp", "derp")
		ll.Info("hello")
		ll.Info("bye")
	}, func(fields map[string]string) {
		for _, fieldName := range []string{"fields.level", "fields.time", "fields.msg"} {
			if _, ok := fields[fieldName]; ok {
				t.Fatalf("should not have prefixed %q: %v", fieldName, fields)
			}
		}
	})
}

func TestWithTimeShouldOverrideTime(t *testing.T) {
	now := time.Now().Add(24 * time.Hour)

	LogAndAssertJSON(t, func(log *Logger) {
		log.WithTime(now).Info("foobar")
	}, func(fields Fields) {
		assert.Equal(t, fields["time"], now.Format(time.RFC3339))
	})
}

func TestWithTimeShouldNotOverrideFields(t *testing.T) {
	now := time.Now().Add(24 * time.Hour)

	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("herp", "derp").WithTime(now).Info("blah")
	}, func(fields Fields) {
		assert.Equal(t, fields["time"], now.Format(time.RFC3339))
		assert.Equal(t, fields["herp"], "derp")
	})
}

func TestWithFieldShouldNotOverrideTime(t *testing.T) {
	now := time.Now().Add(24 * time.Hour)

	LogAndAssertJSON(t, func(log *Logger) {
		log.WithTime(now).WithField("herp", "derp").Info("blah")
	}, func(fields Fields) {
		assert.Equal(t, fields["time"], now.Format(time.RFC3339))
		assert.Equal(t, fields["herp"], "derp")
	})
}

func TestTimeOverrideMultipleLogs(t *testing.T) {
	var buffer bytes.Buffer
	var firstFields, secondFields Fields

	logger := New()
	logger.Out = &buffer
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = time.StampMilli
	logger.Formatter = formatter

	llog := logger.WithField("herp", "derp")
	llog.Info("foo")

	err := json.Unmarshal(buffer.Bytes(), &firstFields)
	assert.NoError(t, err, "should have decoded first message")

	buffer.Reset()

	time.Sleep(10 * time.Millisecond)
	llog.Info("bar")

	err = json.Unmarshal(buffer.Bytes(), &secondFields)
	assert.NoError(t, err, "should have decoded second message")

	assert.NotEqual(t, firstFields["time"], secondFields["time"], "timestamps should not be equal")
}

func TestDoubleLoggingDoesntPrefixPreviousFields(t *testing.T) {

	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	llog := logger.WithField("context", "eating raw fish")

	llog.Info("looks delicious")

	err := json.Unmarshal(buffer.Bytes(), &fields)
	assert.NoError(t, err, "should have decoded first message")
	assert.Equal(t, len(fields), 4, "should only have msg/time/level/context fields")
	assert.Equal(t, fields["msg"], "looks delicious")
	assert.Equal(t, fields["context"], "eating raw fish")

	buffer.Reset()

	llog.Warn("omg it is!")

	err = json.Unmarshal(buffer.Bytes(), &fields)
	assert.NoError(t, err, "should have decoded second message")
	assert.Equal(t, len(fields), 4, "should only have msg/time/level/context fields")
	assert.Equal(t, "omg it is!", fields["msg"])
	assert.Equal(t, "eating raw fish", fields["context"])
	assert.Nil(t, fields["fields.msg"], "should not have prefixed previous `msg` entry")

}

func TestNestedLoggingReportsCorrectCaller(t *testing.T) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)
	logger.ReportCaller = true

	llog := logger.WithField("context", "eating raw fish")

	llog.Info("looks delicious")
	_, _, line, _ := runtime.Caller(0)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	require.NoError(t, err, "should have decoded first message")
	assert.Equal(t, 6, len(fields), "should have msg/time/level/func/context fields")
	assert.Equal(t, "looks delicious", fields["msg"])
	assert.Equal(t, "eating raw fish", fields["context"])
	assert.Equal(t,
		"github.com/sirupsen/logrus_test.TestNestedLoggingReportsCorrectCaller", fields["func"])
	cwd, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, filepath.ToSlash(fmt.Sprintf("%s/logrus_test.go:%d", cwd, line-1)), filepath.ToSlash(fields["file"].(string)))

	buffer.Reset()

	logger.WithFields(Fields{
		"Clyde": "Stubblefield",
	}).WithFields(Fields{
		"Jab'o": "Starks",
	}).WithFields(Fields{
		"uri": "https://www.youtube.com/watch?v=V5DTznu-9v0",
	}).WithFields(Fields{
		"func": "y drummer",
	}).WithFields(Fields{
		"James": "Brown",
	}).Print("The hardest workin' man in show business")
	_, _, line, _ = runtime.Caller(0)

	err = json.Unmarshal(buffer.Bytes(), &fields)
	assert.NoError(t, err, "should have decoded second message")
	assert.Equal(t, 11, len(fields), "should have all builtin fields plus foo,bar,baz,...")
	assert.Equal(t, "Stubblefield", fields["Clyde"])
	assert.Equal(t, "Starks", fields["Jab'o"])
	assert.Equal(t, "https://www.youtube.com/watch?v=V5DTznu-9v0", fields["uri"])
	assert.Equal(t, "y drummer", fields["fields.func"])
	assert.Equal(t, "Brown", fields["James"])
	assert.Equal(t, "The hardest workin' man in show business", fields["msg"])
	assert.Nil(t, fields["fields.msg"], "should not have prefixed previous `msg` entry")
	assert.Equal(t,
		"github.com/sirupsen/logrus_test.TestNestedLoggingReportsCorrectCaller", fields["func"])
	require.NoError(t, err)
	assert.Equal(t, filepath.ToSlash(fmt.Sprintf("%s/logrus_test.go:%d", cwd, line-1)), filepath.ToSlash(fields["file"].(string)))

	logger.ReportCaller = false // return to default value
}

func logLoop(iterations int, reportCaller bool) {
	var buffer bytes.Buffer

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)
	logger.ReportCaller = reportCaller

	for i := 0; i < iterations; i++ {
		logger.Infof("round %d of %d", i, iterations)
	}
}

// Assertions for upper bounds to reporting overhead
func TestCallerReportingOverhead(t *testing.T) {
	iterations := 5000
	before := time.Now()
	logLoop(iterations, false)
	during := time.Now()
	logLoop(iterations, true)
	after := time.Now()

	elapsedNotReporting := during.Sub(before).Nanoseconds()
	elapsedReporting := after.Sub(during).Nanoseconds()

	maxDelta := 1 * time.Second
	assert.WithinDuration(t, during, before, maxDelta,
		"%d log calls without caller name lookup takes less than %d second(s) (was %d nanoseconds)",
		iterations, maxDelta.Seconds(), elapsedNotReporting)
	assert.WithinDuration(t, after, during, maxDelta,
		"%d log calls without caller name lookup takes less than %d second(s) (was %d nanoseconds)",
		iterations, maxDelta.Seconds(), elapsedReporting)
}

// benchmarks for both with and without caller-function reporting
func BenchmarkWithoutCallerTracing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logLoop(1000, false)
	}
}

func BenchmarkWithCallerTracing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logLoop(1000, true)
	}
}

func TestConvertLevelToString(t *testing.T) {
	assert.Equal(t, "trace", TraceLevel.String())
	assert.Equal(t, "debug", DebugLevel.String())
	assert.Equal(t, "info", InfoLevel.String())
	assert.Equal(t, "warning", WarnLevel.String())
	assert.Equal(t, "error", ErrorLevel.String())
	assert.Equal(t, "fatal", FatalLevel.String())
	assert.Equal(t, "panic", PanicLevel.String())
}

func TestParseLevel(t *testing.T) {
	l, err := ParseLevel("panic")
	assert.Nil(t, err)
	assert.Equal(t, PanicLevel, l)

	l, err = ParseLevel("PANIC")
	assert.Nil(t, err)
	assert.Equal(t, PanicLevel, l)

	l, err = ParseLevel("fatal")
	assert.Nil(t, err)
	assert.Equal(t, FatalLevel, l)

	l, err = ParseLevel("FATAL")
	assert.Nil(t, err)
	assert.Equal(t, FatalLevel, l)

	l, err = ParseLevel("error")
	assert.Nil(t, err)
	assert.Equal(t, ErrorLevel, l)

	l, err = ParseLevel("ERROR")
	assert.Nil(t, err)
	assert.Equal(t, ErrorLevel, l)

	l, err = ParseLevel("warn")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("WARN")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("warning")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("WARNING")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("info")
	assert.Nil(t, err)
	assert.Equal(t, InfoLevel, l)

	l, err = ParseLevel("INFO")
	assert.Nil(t, err)
	assert.Equal(t, InfoLevel, l)

	l, err = ParseLevel("debug")
	assert.Nil(t, err)
	assert.Equal(t, DebugLevel, l)

	l, err = ParseLevel("DEBUG")
	assert.Nil(t, err)
	assert.Equal(t, DebugLevel, l)

	l, err = ParseLevel("trace")
	assert.Nil(t, err)
	assert.Equal(t, TraceLevel, l)

	l, err = ParseLevel("TRACE")
	assert.Nil(t, err)
	assert.Equal(t, TraceLevel, l)

	l, err = ParseLevel("invalid")
	assert.Equal(t, "not a valid logrus Level: \"invalid\"", err.Error())
}

func TestLevelString(t *testing.T) {
	var loggerlevel Level
	loggerlevel = 32000

	_ = loggerlevel.String()
}

func TestGetSetLevelRace(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				SetLevel(InfoLevel)
			} else {
				GetLevel()
			}
		}(i)

	}
	wg.Wait()
}

func TestLoggingRace(t *testing.T) {
	logger := New()

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			logger.Info("info")
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestLoggingRaceWithHooksOnEntry(t *testing.T) {
	logger := New()
	hook := new(ModifyHook)
	logger.AddHook(hook)
	entry := logger.WithField("context", "clue")

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			entry.Info("info")
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestReplaceHooks(t *testing.T) {
	old, cur := &TestHook{}, &TestHook{}

	logger := New()
	logger.SetOutput(ioutil.Discard)
	logger.AddHook(old)

	hooks := make(LevelHooks)
	hooks.Add(cur)
	replaced := logger.ReplaceHooks(hooks)

	logger.Info("test")

	assert.Equal(t, old.Fired, false)
	assert.Equal(t, cur.Fired, true)

	logger.ReplaceHooks(replaced)
	logger.Info("test")
	assert.Equal(t, old.Fired, true)
}

// Compile test
func TestLogrusInterfaces(t *testing.T) {
	var buffer bytes.Buffer
	// This verifies FieldLogger and Ext1FieldLogger work as designed.
	// Please don't use them. Use Logger and Entry directly.
	fn := func(xl Ext1FieldLogger) {
		var l FieldLogger = xl
		b := l.WithField("key", "value")
		b.Debug("Test")
	}
	// test logger
	logger := New()
	logger.Out = &buffer
	fn(logger)

	// test Entry
	e := logger.WithField("another", "value")
	fn(e)
}

// Implements io.Writer using channels for synchronization, so we can wait on
// the Entry.Writer goroutine to write in a non-racey way. This does assume that
// there is a single call to Logger.Out for each message.
type channelWriter chan []byte

func (cw channelWriter) Write(p []byte) (int, error) {
	cw <- p
	return len(p), nil
}

func TestEntryWriter(t *testing.T) {
	cw := channelWriter(make(chan []byte, 1))
	log := New()
	log.Out = cw
	log.Formatter = new(JSONFormatter)
	log.WithField("foo", "bar").WriterLevel(WarnLevel).Write([]byte("hello\n"))

	bs := <-cw
	var fields Fields
	err := json.Unmarshal(bs, &fields)
	assert.Nil(t, err)
	assert.Equal(t, fields["foo"], "bar")
	assert.Equal(t, fields["level"], "warning")
}

func TestLogLevelEnabled(t *testing.T) {
	log := New()
	log.SetLevel(PanicLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, false, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, false, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, false, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))
	assert.Equal(t, false, log.IsLevelEnabled(TraceLevel))

	log.SetLevel(FatalLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, false, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, false, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))
	assert.Equal(t, false, log.IsLevelEnabled(TraceLevel))

	log.SetLevel(ErrorLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, false, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))
	assert.Equal(t, false, log.IsLevelEnabled(TraceLevel))

	log.SetLevel(WarnLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, true, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))
	assert.Equal(t, false, log.IsLevelEnabled(TraceLevel))

	log.SetLevel(InfoLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, true, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, true, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))
	assert.Equal(t, false, log.IsLevelEnabled(TraceLevel))

	log.SetLevel(DebugLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, true, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, true, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, true, log.IsLevelEnabled(DebugLevel))
	assert.Equal(t, false, log.IsLevelEnabled(TraceLevel))

	log.SetLevel(TraceLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, true, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, true, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, true, log.IsLevelEnabled(DebugLevel))
	assert.Equal(t, true, log.IsLevelEnabled(TraceLevel))
}

func TestReportCallerOnTextFormatter(t *testing.T) {
	l := New()

	l.Formatter.(*TextFormatter).ForceColors = true
	l.Formatter.(*TextFormatter).DisableColors = false
	l.WithFields(Fields{"func": "func", "file": "file"}).Info("test")

	l.Formatter.(*TextFormatter).ForceColors = false
	l.Formatter.(*TextFormatter).DisableColors = true
	l.WithFields(Fields{"func": "func", "file": "file"}).Info("test")
}

func TestSetReportCallerRace(t *testing.T) {
	l := New()
	l.Out = ioutil.Discard
	l.SetReportCaller(true)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			l.Error("Some Error")
			wg.Done()
		}()
	}
	wg.Wait()
}
