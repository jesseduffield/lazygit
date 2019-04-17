package logrus_test

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	mystring string
)

type GlobalHook struct {
}

func (h *GlobalHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *GlobalHook) Fire(e *logrus.Entry) error {
	e.Data["mystring"] = mystring
	return nil
}

func ExampleGlobalVariableHook() {
	l := logrus.New()
	l.Out = os.Stdout
	l.Formatter = &logrus.TextFormatter{DisableTimestamp: true, DisableColors: true}
	l.AddHook(&GlobalHook{})
	mystring = "first value"
	l.Info("first log")
	mystring = "another value"
	l.Info("second log")
	// Output:
	// level=info msg="first log" mystring="first value"
	// level=info msg="second log" mystring="another value"
}
