// Package i18n provides support for looking up messages
// according to a set of locale preferences.
//
// Create a Bundle to use for the lifetime of your application.
//     bundle := i18n.NewBundle(language.English)
//
// Load translations into your bundle during initialization.
//     bundle.LoadMessageFile("en-US.yaml")
//
// Create a Localizer to use for a set of language preferences.
//     func(w http.ResponseWriter, r *http.Request) {
//         lang := r.FormValue("lang")
//         accept := r.Header.Get("Accept-Language")
//         localizer := i18n.NewLocalizer(bundle, lang, accept)
//     }
//
// Use the Localizer to lookup messages.
//     localizer.MustLocalize(&i18n.LocalizeConfig{
// 	        DefaultMessage: &i18n.Message{
// 	            ID: "HelloWorld",
// 	            Other: "Hello World!",
// 	        },
//     })
package i18n
