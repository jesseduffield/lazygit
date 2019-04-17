package main

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	logrus.Info("succeeded")
	logrus.Warn("not correct")
	logrus.Error("something error")
	logrus.Fatal("panic")
}
