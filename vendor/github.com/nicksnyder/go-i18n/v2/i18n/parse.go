package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/text/language"
)

// MessageFile represents a parsed message file.
type MessageFile struct {
	Path     string
	Tag      language.Tag
	Format   string
	Messages []*Message
}

// ParseMessageFileBytes returns the messages parsed from file.
func ParseMessageFileBytes(buf []byte, path string, unmarshalFuncs map[string]UnmarshalFunc) (*MessageFile, error) {
	lang, format := parsePath(path)
	tag := language.Make(lang)
	messageFile := &MessageFile{
		Path:   path,
		Tag:    tag,
		Format: format,
	}
	if len(buf) == 0 {
		return messageFile, nil
	}
	unmarshalFunc := unmarshalFuncs[messageFile.Format]
	if unmarshalFunc == nil {
		if messageFile.Format == "json" {
			unmarshalFunc = json.Unmarshal
		} else {
			return nil, fmt.Errorf("no unmarshaler registered for %s", messageFile.Format)
		}
	}
	var err error
	var raw interface{}
	if err = unmarshalFunc(buf, &raw); err != nil {
		return nil, err
	}

	if messageFile.Messages, err = recGetMessages(raw, isMessage(raw), true); err != nil {
		return nil, err
	}

	return messageFile, nil
}

const nestedSeparator = "."

var errInvalidTranslationFile = errors.New("invalid translation file, expected key-values, got a single value")

// recGetMessages looks for translation messages inside "raw" parameter,
// scanning nested maps using recursion.
func recGetMessages(raw interface{}, isMapMessage, isInitialCall bool) ([]*Message, error) {
	var messages []*Message
	var err error

	switch data := raw.(type) {
	case string:
		if isInitialCall {
			return nil, errInvalidTranslationFile
		}
		m, err := NewMessage(data)
		return []*Message{m}, err

	case map[string]interface{}:
		if isMapMessage {
			m, err := NewMessage(data)
			return []*Message{m}, err
		}
		messages = make([]*Message, 0, len(data))
		for id, data := range data {
			// recursively scan map items
			messages, err = addChildMessages(id, data, messages)
			if err != nil {
				return nil, err
			}
		}

	case map[interface{}]interface{}:
		if isMapMessage {
			m, err := NewMessage(data)
			return []*Message{m}, err
		}
		messages = make([]*Message, 0, len(data))
		for id, data := range data {
			strid, ok := id.(string)
			if !ok {
				return nil, fmt.Errorf("expected key to be string but got %#v", id)
			}
			// recursively scan map items
			messages, err = addChildMessages(strid, data, messages)
			if err != nil {
				return nil, err
			}
		}

	case []interface{}:
		// Backward compatibility for v1 file format.
		messages = make([]*Message, 0, len(data))
		for _, data := range data {
			// recursively scan slice items
			childMessages, err := recGetMessages(data, isMessage(data), false)
			if err != nil {
				return nil, err
			}
			messages = append(messages, childMessages...)
		}

	default:
		return nil, fmt.Errorf("unsupported file format %T", raw)
	}

	return messages, nil
}

func addChildMessages(id string, data interface{}, messages []*Message) ([]*Message, error) {
	isChildMessage := isMessage(data)
	childMessages, err := recGetMessages(data, isChildMessage, false)
	if err != nil {
		return nil, err
	}
	for _, m := range childMessages {
		if isChildMessage {
			if m.ID == "" {
				m.ID = id // start with innermost key
			}
		} else {
			m.ID = id + nestedSeparator + m.ID // update ID with each nested key on the way
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func parsePath(path string) (langTag, format string) {
	formatStartIdx := -1
	for i := len(path) - 1; i >= 0; i-- {
		c := path[i]
		if os.IsPathSeparator(c) {
			if formatStartIdx != -1 {
				langTag = path[i+1 : formatStartIdx]
			}
			return
		}
		if path[i] == '.' {
			if formatStartIdx != -1 {
				langTag = path[i+1 : formatStartIdx]
				return
			}
			if formatStartIdx == -1 {
				format = path[i+1:]
				formatStartIdx = i
			}
		}
	}
	if formatStartIdx != -1 {
		langTag = path[:formatStartIdx]
	}
	return
}
