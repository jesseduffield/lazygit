package plural

// Form represents a language pluralization form as defined here:
// http://cldr.unicode.org/index/cldr-spec/plural-rules
type Form string

// All defined plural forms.
const (
	Invalid Form = ""
	Zero    Form = "zero"
	One     Form = "one"
	Two     Form = "two"
	Few     Form = "few"
	Many    Form = "many"
	Other   Form = "other"
)
