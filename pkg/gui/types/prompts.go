package types

type AskOpts struct {
	Title               string
	Prompt              string
	HandleConfirm       func() error
	HandleClose         func() error
	HandlersManageFocus bool
	FindSuggestionsFunc func(string) []*Suggestion
}

type PromptOpts struct {
	Title               string
	InitialContent      string
	HandleConfirm       func(string) error
	FindSuggestionsFunc func(string) []*Suggestion
}

type Suggestion struct {
	// value is the thing that we're matching on and the thing that will be submitted if you select the suggestion
	Value string
	// label is what is actually displayed so it can e.g. contain color
	Label string
}
