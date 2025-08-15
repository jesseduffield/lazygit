package logfmt

import (
	"errors"
	"fmt"
)

var ErrUnterminatedString = errors.New("logfmt: unterminated string")

func gotoScanner(data []byte, h Handler) (err error) {
	saveError := func(e error) {
		if err == nil {
			err = e
		}
	}

	var c byte
	var i int
	var m int
	var key []byte
	var val []byte
	var ok bool
	var esc bool

garbage:
	if i == len(data) {
		return
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		key, val = nil, nil
		m = i
		i++
		goto key
	default:
		i++
		goto garbage
	}

key:
	if i >= len(data) {
		if m >= 0 {
			key = data[m:i]
			saveError(h.HandleLogfmt(key, nil))
		}
		return
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		i++
		goto key
	case c == '=':
		key = data[m:i]
		i++
		goto equal
	default:
		key = data[m:i]
		i++
		saveError(h.HandleLogfmt(key, nil))
		goto garbage
	}

equal:
	if i >= len(data) {
		if m >= 0 {
			i--
			key = data[m:i]
			saveError(h.HandleLogfmt(key, nil))
		}
		return
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		m = i
		i++
		goto ivalue
	case c == '"':
		m = i
		i++
		esc = false
		goto qvalue
	default:
		if key != nil {
			saveError(h.HandleLogfmt(key, val))
		}
		i++
		goto garbage
	}

ivalue:
	if i >= len(data) {
		if m >= 0 {
			val = data[m:i]
			saveError(h.HandleLogfmt(key, val))
		}
		return
	}

	c = data[i]
	switch {
	case c > ' ' && c != '"' && c != '=':
		i++
		goto ivalue
	default:
		val = data[m:i]
		saveError(h.HandleLogfmt(key, val))
		i++
		goto garbage
	}

qvalue:
	if i >= len(data) {
		if m >= 0 {
			saveError(ErrUnterminatedString)
		}
		return
	}

	c = data[i]
	switch c {
	case '\\':
		i += 2
		esc = true
		goto qvalue
	case '"':
		i++
		val = data[m:i]
		if esc {
			val, ok = unquoteBytes(val)
			if !ok {
				saveError(fmt.Errorf("logfmt: error unquoting bytes %q", string(val)))
				goto garbage
			}
		} else {
			val = val[1 : len(val)-1]
		}
		saveError(h.HandleLogfmt(key, val))
		goto garbage
	default:
		i++
		goto qvalue
	}
}
