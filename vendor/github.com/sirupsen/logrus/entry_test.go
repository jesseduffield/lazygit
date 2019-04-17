package logrus

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntryWithError(t *testing.T) {

	assert := assert.New(t)

	defer func() {
		ErrorKey = "error"
	}()

	err := fmt.Errorf("kaboom at layer %d", 4711)

	assert.Equal(err, WithError(err).Data["error"])

	logger := New()
	logger.Out = &bytes.Buffer{}
	entry := NewEntry(logger)

	assert.Equal(err, entry.WithError(err).Data["error"])

	ErrorKey = "err"

	assert.Equal(err, entry.WithError(err).Data["err"])

}

func TestEntryWithContext(t *testing.T) {
	assert := assert.New(t)
	ctx := context.WithValue(context.Background(), "foo", "bar")

	assert.Equal(ctx, WithContext(ctx).Context)

	logger := New()
	logger.Out = &bytes.Buffer{}
	entry := NewEntry(logger)

	assert.Equal(ctx, entry.WithContext(ctx).Context)
}

func TestEntryPanicln(t *testing.T) {
	errBoom := fmt.Errorf("boom time")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *Entry:
			assert.Equal(t, "kaboom", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := New()
	logger.Out = &bytes.Buffer{}
	entry := NewEntry(logger)
	entry.WithField("err", errBoom).Panicln("kaboom")
}

func TestEntryPanicf(t *testing.T) {
	errBoom := fmt.Errorf("boom again")

	defer func() {
		p := recover()
		assert.NotNil(t, p)

		switch pVal := p.(type) {
		case *Entry:
			assert.Equal(t, "kaboom true", pVal.Message)
			assert.Equal(t, errBoom, pVal.Data["err"])
		default:
			t.Fatalf("want type *Entry, got %T: %#v", pVal, pVal)
		}
	}()

	logger := New()
	logger.Out = &bytes.Buffer{}
	entry := NewEntry(logger)
	entry.WithField("err", errBoom).Panicf("kaboom %v", true)
}

const (
	badMessage   = "this is going to panic"
	panicMessage = "this is broken"
)

type panickyHook struct{}

func (p *panickyHook) Levels() []Level {
	return []Level{InfoLevel}
}

func (p *panickyHook) Fire(entry *Entry) error {
	if entry.Message == badMessage {
		panic(panicMessage)
	}

	return nil
}

func TestEntryHooksPanic(t *testing.T) {
	logger := New()
	logger.Out = &bytes.Buffer{}
	logger.Level = InfoLevel
	logger.Hooks.Add(&panickyHook{})

	defer func() {
		p := recover()
		assert.NotNil(t, p)
		assert.Equal(t, panicMessage, p)

		entry := NewEntry(logger)
		entry.Info("another message")
	}()

	entry := NewEntry(logger)
	entry.Info(badMessage)
}

func TestEntryWithIncorrectField(t *testing.T) {
	assert := assert.New(t)

	fn := func() {}

	e := Entry{}
	eWithFunc := e.WithFields(Fields{"func": fn})
	eWithFuncPtr := e.WithFields(Fields{"funcPtr": &fn})

	assert.Equal(eWithFunc.err, `can not add field "func"`)
	assert.Equal(eWithFuncPtr.err, `can not add field "funcPtr"`)

	eWithFunc = eWithFunc.WithField("not_a_func", "it is a string")
	eWithFuncPtr = eWithFuncPtr.WithField("not_a_func", "it is a string")

	assert.Equal(eWithFunc.err, `can not add field "func"`)
	assert.Equal(eWithFuncPtr.err, `can not add field "funcPtr"`)

	eWithFunc = eWithFunc.WithTime(time.Now())
	eWithFuncPtr = eWithFuncPtr.WithTime(time.Now())

	assert.Equal(eWithFunc.err, `can not add field "func"`)
	assert.Equal(eWithFuncPtr.err, `can not add field "funcPtr"`)
}

func TestEntryLogfLevel(t *testing.T) {
	logger := New()
	buffer := &bytes.Buffer{}
	logger.Out = buffer
	logger.SetLevel(InfoLevel)
	entry := NewEntry(logger)

	entry.Logf(DebugLevel, "%s", "debug")
	assert.NotContains(t, buffer.String(), "debug", )

	entry.Logf(WarnLevel, "%s", "warn")
	assert.Contains(t, buffer.String(), "warn", )
}