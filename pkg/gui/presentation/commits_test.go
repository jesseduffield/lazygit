package presentation

import "testing"

func TestGetInitials(t *testing.T) {
	for input, output := range map[string]string{
		"Jesse Duffield":     "JD",
		"Jesse Duffield Man": "JD",
		"JesseDuffield":      "Je",
		"J":                  "J",
		"":                   "",
	} {
		if output != getInitials(input) {
			t.Errorf("Expected %s to be %s", input, output)
		}
	}
}
