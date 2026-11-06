package helpers

import "testing"

func TestSanitizedBranchName(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		whitespaceReplacement string
		expected             string
	}{
		{name: "default fallback", input: "feature new branch", whitespaceReplacement: "", expected: "feature-new-branch"},
		{name: "custom replacement", input: "feature new branch", whitespaceReplacement: "_", expected: "feature_new_branch"},
		{name: "multiple characters", input: "feature new branch", whitespaceReplacement: "--", expected: "feature--new--branch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizedBranchName(tt.input, tt.whitespaceReplacement)
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
