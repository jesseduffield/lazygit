// The goi18n command formats and merges translation files.
//
//     go get -u github.com/nicksnyder/go-i18n/goi18n
//     goi18n -help
//
// Help documentation:
//
//     goi18n manages translation files.
//
//     Usage:
//
//         goi18n merge     Merge translation files
//         goi18n constants Generate constant file from translation file
//
//     For more details execute:
//
//         goi18n [command] -help
//
//     Merge translation files.
//
//     Usage:
//
//         goi18n merge [options] [files...]
//
//     Translation files:
//
//         A translation file contains the strings and translations for a single language.
//
//         Translation file names must have a suffix of a supported format (e.g. .json) and
//         contain a valid language tag as defined by RFC 5646 (e.g. en-us, fr, zh-hant, etc.).
//
//         For each language represented by at least one input translation file, goi18n will produce 2 output files:
//
//             xx-yy.all.format
//                 This file contains all strings for the language (translated and untranslated).
//                 Use this file when loading strings at runtime.
//
//             xx-yy.untranslated.format
//                 This file contains the strings that have not been translated for this language.
//                 The translations for the strings in this file will be extracted from the source language.
//                 After they are translated, merge them back into xx-yy.all.format using goi18n.
//
//     Merging:
//
//         goi18n will merge multiple translation files for the same language.
//         Duplicate translations will be merged into the existing translation.
//         Non-empty fields in the duplicate translation will overwrite those fields in the existing translation.
//         Empty fields in the duplicate translation are ignored.
//
//     Adding a new language:
//
//         To produce translation files for a new language, create an empty translation file with the
//         appropriate name and pass it in to goi18n.
//
//     Options:
//
//         -sourceLanguage tag
//             goi18n uses the strings from this language to seed the translations for other languages.
//             Default: en-us
//
//         -outdir directory
//             goi18n writes the output translation files to this directory.
//             Default: .
//
//         -format format
//             goi18n encodes the output translation files in this format.
//             Supported formats: json, yaml
//             Default: json
//
//     Generate constant file from translation file.
//
//     Usage:
//
//         goi18n constants [options] [file]
//
//     Translation files:
//
//         A translation file contains the strings and translations for a single language.
//
//         Translation file names must have a suffix of a supported format (e.g. .json) and
//         contain a valid language tag as defined by RFC 5646 (e.g. en-us, fr, zh-hant, etc.).
//
//     Options:
//
//         -package name
//             goi18n generates the constant file under the package name.
//             Default: R
//
//         -outdir directory
//             goi18n writes the constant file to this directory.
//             Default: .
//
package main
