package i18n

import (
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

func TestLocalizer_Localize(t *testing.T) {
	testCases := []struct {
		name              string
		defaultLanguage   language.Tag
		messages          map[language.Tag][]*Message
		acceptLangs       []string
		conf              *LocalizeConfig
		expectedErr       error
		expectedLocalized string
	}{
		{
			name:            "message id mismatch",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloWorld",
				DefaultMessage: &Message{
					ID: "DefaultHelloWorld",
				},
			},
			expectedErr: &messageIDMismatchErr{messageID: "HelloWorld", defaultMessageID: "DefaultHelloWorld"},
		},
		{
			name:            "message id not mismatched",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "Hello!"}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloWorld",
				DefaultMessage: &Message{
					ID: "HelloWorld",
				},
			},
			expectedLocalized: "Hello!",
		},
		{
			name:              "missing translation from default language",
			defaultLanguage:   language.English,
			acceptLangs:       []string{"en"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:            "empty translation without fallback",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld"}},
			},
			acceptLangs: []string{"es"},
			conf:        &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr: &MessageNotFoundErr{messageID: "HelloWorld"},
		},
		{
			name:            "empty translation with fallback",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "Hello World!"}},
				language.Spanish: {{ID: "HelloWorld"}},
			},
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedLocalized: "Hello World!",
		},
		{
			name:            "missing translation from default language with other translation",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "other"}},
			},
			acceptLangs:       []string{"en"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:              "missing translation from not default language",
			defaultLanguage:   language.English,
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:            "missing translation not default language with other translation",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.French: {{ID: "HelloWorld", Other: "other"}},
			},
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:            "accept default language, message in bundle",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "other"}},
			},
			acceptLangs:       []string{"en"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedLocalized: "other",
		},
		{
			name:            "accept default language, message in bundle, default message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "accept not default language, message in bundle",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "other"}},
			},
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedLocalized: "other",
		},
		{
			name:            "accept not default language, other message in bundle, default message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "accept not default language, message in bundle, default message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "accept default language, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "default other",
		},
		{
			name:            "accept not default language, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "default other",
		},
		{
			name:            "fallback to non-default less specific language",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es-ES"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "fallback to non-default more specific language",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.EuropeanSpanish: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "plural count one, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID:   "Cats",
				PluralCount: 1,
			},
			expectedLocalized: "I have 1 cat",
		},
		{
			name:            "plural count other, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID:   "Cats",
				PluralCount: 2,
			},
			expectedLocalized: "I have 2 cats",
		},
		{
			name:            "plural count float, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID:   "Cats",
				PluralCount: "2.5",
			},
			expectedLocalized: "I have 2.5 cats",
		},
		{
			name:            "plural count one, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: 1,
				DefaultMessage: &Message{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				},
			},
			expectedLocalized: "I have 1 cat",
		},
		{
			name:            "plural count other, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: 2,
				DefaultMessage: &Message{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				},
			},
			expectedLocalized: "I have 2 cats",
		},
		{
			name:            "plural count float, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: "2.5",
				DefaultMessage: &Message{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				},
			},
			expectedLocalized: "I have 2.5 cats",
		},
		{
			name:            "template data, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "HelloPerson",
					Other: "Hello {{.Person}}",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloPerson",
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "HelloPerson",
					Other: "Hello {{.Person}}",
				},
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, custom delims, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:         "HelloPerson",
					Other:      "Hello <<.Person>>",
					LeftDelim:  "<<",
					RightDelim: ">>",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloPerson",
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, custom delims, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:         "HelloPerson",
					Other:      "Hello <<.Person>>",
					LeftDelim:  "<<",
					RightDelim: ">>",
				},
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, plural count one, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "PersonCats",
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  1,
				},
				PluralCount: 1,
			},
			expectedLocalized: "Nick has 1 cat",
		},
		{
			name:            "template data, plural count other, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "PersonCats",
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  2,
				},
				PluralCount: 2,
			},
			expectedLocalized: "Nick has 2 cats",
		},
		{
			name:            "template data, plural count float, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "PersonCats",
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  "2.5",
				},
				PluralCount: "2.5",
			},
			expectedLocalized: "Nick has 2.5 cats",
		},
		{
			name:            "template data, plural count one, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				},
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  1,
				},
				PluralCount: 1,
			},
			expectedLocalized: "Nick has 1 cat",
		},
		{
			name:            "template data, plural count other, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				},
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  2,
				},
				PluralCount: 2,
			},
			expectedLocalized: "Nick has 2 cats",
		},
		{
			name:            "template data, plural count float, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				},
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  "2.5",
				},
				PluralCount: "2.5",
			},
			expectedLocalized: "Nick has 2.5 cats",
		},
		{
			name:            "test slow path",
			defaultLanguage: language.Spanish,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Hello",
					Other: "Hello!",
				}},
				language.AmericanEnglish: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
			},
			acceptLangs: []string{"en-US"},
			conf: &LocalizeConfig{
				MessageID: "Hello",
			},
			expectedLocalized: "Hello!",
		},
		{
			name:            "test slow path default message",
			defaultLanguage: language.Spanish,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
				language.AmericanEnglish: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
			},
			acceptLangs: []string{"en-US"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "Hello",
					Other: "Hola!",
				},
			},
			expectedLocalized: "Hola!",
		},
		{
			name:            "test slow path no message",
			defaultLanguage: language.Spanish,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
				language.AmericanEnglish: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
			},
			acceptLangs: []string{"en-US"},
			conf: &LocalizeConfig{
				MessageID: "Hello",
			},
			expectedErr: &MessageNotFoundErr{messageID: "Hello"},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bundle := &Bundle{DefaultLanguage: testCase.defaultLanguage}
			for tag, messages := range testCase.messages {
				bundle.AddMessages(tag, messages...)
			}
			localizer := NewLocalizer(bundle, testCase.acceptLangs...)
			localized, err := localizer.Localize(testCase.conf)
			if !reflect.DeepEqual(err, testCase.expectedErr) {
				t.Errorf("expected error %#v; got %#v", testCase.expectedErr, err)
			}
			if localized != testCase.expectedLocalized {
				t.Errorf("expected localized string %q; got %q", testCase.expectedLocalized, localized)
			}
		})
	}
}
