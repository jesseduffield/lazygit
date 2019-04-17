package i18n_test

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func ExampleLocalizer_MustLocalize() {
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	localizer := i18n.NewLocalizer(bundle, "en")
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HelloWorld",
			Other: "Hello World!",
		},
	}))
	// Output:
	// Hello World!
}

func ExampleLocalizer_MustLocalize_noDefaultMessage() {
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
HelloWorld = "Hello World!"
`), "en.toml")
	bundle.MustParseMessageFileBytes([]byte(`
HelloWorld = "Hola Mundo!"
`), "es.toml")

	{
		localizer := i18n.NewLocalizer(bundle, "en-US")
		fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "HelloWorld"}))
	}
	{
		localizer := i18n.NewLocalizer(bundle, "es-ES")
		fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "HelloWorld"}))
	}
	// Output:
	// Hello World!
	// Hola Mundo!
}

func ExampleLocalizer_MustLocalize_plural() {
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	localizer := i18n.NewLocalizer(bundle, "en")
	catsMessage := &i18n.Message{
		ID:    "Cats",
		One:   "I have {{.PluralCount}} cat.",
		Other: "I have {{.PluralCount}} cats.",
	}
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: catsMessage,
		PluralCount:    1,
	}))
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: catsMessage,
		PluralCount:    2,
	}))
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: catsMessage,
		PluralCount:    "2.5",
	}))
	// Output:
	// I have 1 cat.
	// I have 2 cats.
	// I have 2.5 cats.
}

func ExampleLocalizer_MustLocalize_template() {
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	localizer := i18n.NewLocalizer(bundle, "en")
	helloPersonMessage := &i18n.Message{
		ID:    "HelloPerson",
		Other: "Hello {{.Name}}!",
	}
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: helloPersonMessage,
		TemplateData:   map[string]string{"Name": "Nick"},
	}))
	// Output:
	// Hello Nick!
}

func ExampleLocalizer_MustLocalize_plural_template() {
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	localizer := i18n.NewLocalizer(bundle, "en")
	personCatsMessage := &i18n.Message{
		ID:    "PersonCats",
		One:   "{{.Name}} has {{.Count}} cat.",
		Other: "{{.Name}} has {{.Count}} cats.",
	}
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: personCatsMessage,
		PluralCount:    1,
		TemplateData: map[string]interface{}{
			"Name":  "Nick",
			"Count": 1,
		},
	}))
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: personCatsMessage,
		PluralCount:    2,
		TemplateData: map[string]interface{}{
			"Name":  "Nick",
			"Count": 2,
		},
	}))
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: personCatsMessage,
		PluralCount:    "2.5",
		TemplateData: map[string]interface{}{
			"Name":  "Nick",
			"Count": "2.5",
		},
	}))
	// Output:
	// Nick has 1 cat.
	// Nick has 2 cats.
	// Nick has 2.5 cats.
}

func ExampleLocalizer_MustLocalize_customTemplateDelims() {
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	localizer := i18n.NewLocalizer(bundle, "en")
	helloPersonMessage := &i18n.Message{
		ID:         "HelloPerson",
		Other:      "Hello <<.Name>>!",
		LeftDelim:  "<<",
		RightDelim: ">>",
	}
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: helloPersonMessage,
		TemplateData:   map[string]string{"Name": "Nick"},
	}))
	// Output:
	// Hello Nick!
}
