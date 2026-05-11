package utils

import (
	"io"
	"sync"
)

// This wraps a writer and ensures that before we actually write anything we call a given function first

type OnceWriter struct {
	writer io.Writer
	once   sync.Once
	f      func()
}

var _ io.Writer = &OnceWriter{}

func NewOnceWriter(writer io.Writer, f func()) *OnceWriter {
	return &OnceWriter{
		writer: writer,
		f:      f,
	}
}

func (self *OnceWriter) Write(p []byte) (n int, err error) {
	self.once.Do(func() {
		self.f()
	})

	return self.writer.Write(p)
}
