package authors

import "testing"

func TestGetInitials(t *testing.T) {
	for input, expectedOutput := range map[string]string{
		"Jesse Duffield":     "JD",
		"Jesse Duffield Man": "JD",
		"JesseDuffield":      "Je",
		"J":                  "J",
		"六书六書":               "六",
		"書":                  "書",
		"":                   "",
	} {
		output := getInitials(input)
		if output != expectedOutput {
			t.Errorf("Expected %s to be %s", output, expectedOutput)
		}
	}
}

func TestSetUnspecifiedAuthorColors(t *testing.T) {
	customAuthorColors := []string{"#FF0000", "red", "cyan"}

	SetUnspecifiedAuthorColors(customAuthorColors)

	if len(unspecifiedAuthorColors) != len(customAuthorColors) {
		t.Errorf("Expected %d unspecified author colors, but got %d", len(customAuthorColors), len(unspecifiedAuthorColors))
	}
}
