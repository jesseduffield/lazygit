package config

import (
	"fmt"
	"io"
	"strings"
)

// An Encoder writes config files to an output stream.
type Encoder struct {
	w io.Writer
}

var (
	subsectionReplacer = strings.NewReplacer(`"`, `\"`, `\`, `\\`)
	valueReplacer = strings.NewReplacer(`"`, `\"`, `\`, `\\`, "\n", `\n`, "\t", `\t`, "\b", `\b`)
)
// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encode writes the config in git config format to the stream of the encoder.
func (e *Encoder) Encode(cfg *Config) error {
	for _, s := range cfg.Sections {
		if err := e.encodeSection(s); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeSection(s *Section) error {
	if len(s.Options) > 0 {
		if err := e.printf("[%s]\n", s.Name); err != nil {
			return err
		}

		if err := e.encodeOptions(s.Options); err != nil {
			return err
		}
	}

	for _, ss := range s.Subsections {
		if err := e.encodeSubsection(s.Name, ss); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeSubsection(sectionName string, s *Subsection) error {
	if err := e.printf("[%s \"%s\"]\n", sectionName, subsectionReplacer.Replace(s.Name)); err != nil {
		return err
	}

	return e.encodeOptions(s.Options)
}

func (e *Encoder) encodeOptions(opts Options) error {
	for _, o := range opts {
		var value string
		if strings.ContainsAny(o.Value, "#;\"\t\n\\") || strings.HasPrefix(o.Value, " ") || strings.HasSuffix(o.Value, " ") {
			value = `"`+valueReplacer.Replace(o.Value)+`"`
		} else {
			value = o.Value
		}

		if err := e.printf("\t%s = %s\n", o.Key, value); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) printf(msg string, args ...interface{}) error {
	_, err := fmt.Fprintf(e.w, msg, args...)
	return err
}
