package humanlog

import (
	"bufio"
	"bytes"
	"io"
)

var (
	eol = [...]byte{'\n'}
)

// Scanner reads logfmt'd lines from src and prettify them onto dst.
// If the lines aren't logfmt, it will simply write them out with no
// prettification.
func Scanner(src io.Reader, dst io.Writer, opts *HandlerOptions) error {
	in := bufio.NewScanner(src)
	in.Split(bufio.ScanLines)

	var line uint64

	var lastLogfmt bool
	var lastJSON bool

	logfmtEntry := LogfmtHandler{Opts: opts}
	jsonEntry := JSONHandler{Opts: opts}

	for in.Scan() {
		line++
		lineData := in.Bytes()

		// remove that pesky syslog crap
		lineData = bytes.TrimPrefix(lineData, []byte("@cee: "))

		switch {

		case jsonEntry.TryHandle(lineData):
			dst.Write(jsonEntry.Prettify(opts.SkipUnchanged && lastJSON))
			lastJSON = true

		case logfmtEntry.TryHandle(lineData):
			dst.Write(logfmtEntry.Prettify(opts.SkipUnchanged && lastLogfmt))
			lastLogfmt = true

		case tryDockerComposePrefix(lineData, &jsonEntry):
			dst.Write(jsonEntry.Prettify(opts.SkipUnchanged && lastJSON))
			lastJSON = true

		case tryDockerComposePrefix(lineData, &logfmtEntry):
			dst.Write(logfmtEntry.Prettify(opts.SkipUnchanged && lastLogfmt))
			lastLogfmt = true

		default:
			lastLogfmt = false
			lastJSON = false
			dst.Write(lineData)
		}
		dst.Write(eol[:])

	}

	switch err := in.Err(); err {
	case nil, io.EOF:
		return nil
	default:
		return err
	}
}
