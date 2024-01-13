package osc52

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

// output is the default output for Copy which uses os.Stdout and os.Environ.
var output = NewOutput(os.Stdout, os.Environ())

// envs is a map of environment variables.
type envs map[string]string

// Get returns the value of the environment variable named by the key.
func (e envs) Get(key string) string {
	v, ok := e[key]
	if !ok {
		return ""
	}
	return v
}

// Output is where the OSC52 string should be written.
type Output struct {
	out  io.Writer
	envs envs
}

// NewOutput returns a new Output.
func NewOutput(out io.Writer, envs []string) *Output {
	e := make(map[string]string, 0)
	for _, env := range envs {
		s := strings.Split(env, "=")
		k := s[0]
		v := strings.Join(s[1:], "=")
		e[k] = v
	}
	o := &Output{
		out:  out,
		envs: e,
	}
	return o
}

// Copy copies the OSC52 string to the output. This is the default copy function.
func Copy(str string) {
	output.Copy(str)
}

// Copy copies the OSC52 string to the output.
func (o *Output) Copy(str string) {
	mode := "default"
	term := o.envs.Get("TERM")
	switch {
	case o.envs.Get("TMUX") != "", strings.HasPrefix(term, "tmux"):
		mode = "tmux"
	case strings.HasPrefix(term, "screen"):
		mode = "screen"
	case strings.Contains(term, "kitty"):
		mode = "kitty"
	}

	switch mode {
	case "default":
		o.copyDefault(str)
	case "tmux":
		o.copyTmux(str)
	case "screen":
		o.copyDCS(str)
	case "kitty":
		o.copyKitty(str)
	}
}

// copyDefault copies the OSC52 string to the output.
func (o *Output) copyDefault(str string) {
	b64 := base64.StdEncoding.EncodeToString([]byte(str))
	o.out.Write([]byte("\x1b]52;c;" + b64 + "\x07"))
}

// copyTmux copies the OSC52 string escaped for Tmux.
func (o *Output) copyTmux(str string) {
	b64 := base64.StdEncoding.EncodeToString([]byte(str))
	o.out.Write([]byte("\x1bPtmux;\x1b\x1b]52;c;" + b64 + "\x07\x1b\\"))
}

// copyDCS copies the OSC52 string wrapped in a DCS sequence which is
// appropriate when using screen.
//
// Screen doesn't support OSC52 but will pass the contents of a DCS sequence to
// the outer terminal unchanged.
func (o *Output) copyDCS(str string) {
	// Here, we split the encoded string into 76 bytes chunks and then join the
	// chunks with <end-dsc><start-dsc> sequences. Finally, wrap the whole thing in
	// <start-dsc><start-osc52><joined-chunks><end-osc52><end-dsc>.
	b64 := base64.StdEncoding.EncodeToString([]byte(str))
	s := strings.SplitN(b64, "", 76)
	q := fmt.Sprintf("\x1bP\x1b]52;c;%s\x07\x1b\x5c", strings.Join(s, "\x1b\\\x1bP"))
	o.out.Write([]byte(q))
}

// copyKitty copies the OSC52 string to Kitty. First, it flushes the keyboard
// before copying, this is required for Kitty < 0.22.0.
func (o *Output) copyKitty(str string) {
	o.out.Write([]byte("\x1b]52;c;!\x07"))
	o.copyDefault(str)
}
