package internal

import (
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

func TestParseMessageFileBytes(t *testing.T) {
	testCases := []struct {
		file           string
		path           string
		unmarshalFuncs map[string]UnmarshalFunc
		messageFile    *MessageFile
		err            error
	}{
		{
			file: `{"hello": "world"}`,
			path: "en.json",
			messageFile: &MessageFile{
				Path:   "en.json",
				Tag:    language.English,
				Format: "json",
				Messages: []*Message{{
					ID:    "hello",
					Other: "world",
				}},
			},
		},
	}
	for _, testCase := range testCases {
		actual, err := ParseMessageFileBytes([]byte(testCase.file), testCase.path, testCase.unmarshalFuncs)
		if err != testCase.err {
			t.Fatalf("expected error %#v; got %#v", testCase.err, err)
		}
		if actual.Path != testCase.messageFile.Path {
			t.Fatalf("expected path %q; got %q", testCase.messageFile.Path, actual.Path)
		}
		if actual.Tag != testCase.messageFile.Tag {
			t.Fatalf("expected tag %q; got %q", testCase.messageFile.Tag, actual.Tag)
		}
		if actual.Format != testCase.messageFile.Format {
			t.Fatalf("expected format %q; got %q", testCase.messageFile.Format, actual.Format)
		}
		if !reflect.DeepEqual(actual.Messages, testCase.messageFile.Messages) {
			t.Fatalf("expected %#v; got %#v", testCase.messageFile.Messages, actual.Messages)
		}
	}
}
