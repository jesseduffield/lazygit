package syntax

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
)

func Escape(input string) string {
	b := &bytes.Buffer{}
	for _, r := range input {
		escape(b, r, false)
	}
	return b.String()
}

const meta = `\.+*?()|[]{}^$# `

func escape(b *bytes.Buffer, r rune, force bool) {
	if unicode.IsPrint(r) {
		if strings.IndexRune(meta, r) >= 0 || force {
			b.WriteRune('\\')
		}
		b.WriteRune(r)
		return
	}

	switch r {
	case '\a':
		b.WriteString(`\a`)
	case '\f':
		b.WriteString(`\f`)
	case '\n':
		b.WriteString(`\n`)
	case '\r':
		b.WriteString(`\r`)
	case '\t':
		b.WriteString(`\t`)
	case '\v':
		b.WriteString(`\v`)
	default:
		if r < 0x100 {
			b.WriteString(`\x`)
			s := strconv.FormatInt(int64(r), 16)
			if len(s) == 1 {
				b.WriteRune('0')
			}
			b.WriteString(s)
			break
		}
		b.WriteString(`\u`)
		b.WriteString(strconv.FormatInt(int64(r), 16))
	}
}

func Unescape(input string) (string, error) {
	idx := strings.IndexRune(input, '\\')
	// no slashes means no unescape needed
	if idx == -1 {
		return input, nil
	}

	buf := bytes.NewBufferString(input[:idx])
	// get the runes for the rest of the string -- we're going full parser scan on this

	p := parser{}
	p.setPattern(input[idx+1:])
	for {
		if p.rightMost() {
			return "", p.getErr(ErrIllegalEndEscape)
		}
		r, err := p.scanCharEscape()
		if err != nil {
			return "", err
		}
		buf.WriteRune(r)
		// are we done?
		if p.rightMost() {
			return buf.String(), nil
		}

		r = p.moveRightGetChar()
		for r != '\\' {
			buf.WriteRune(r)
			if p.rightMost() {
				// we're done, no more slashes
				return buf.String(), nil
			}
			// keep scanning until we get another slash
			r = p.moveRightGetChar()
		}
	}
}
