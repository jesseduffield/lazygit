package helpers

import "testing"

func TestMatchesPromptTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		prompt   string
		template string
		want     bool
	}{
		{
			name:     "matches template with placeholder",
			prompt:   "All merge conflicts resolved. Continue the rebase?",
			template: "All merge conflicts resolved. Continue the %s?",
			want:     true,
		},
		{
			name:     "matches translated-like template with placeholder",
			prompt:   "Tous les conflits sont resolus. Continuer le rebase ?",
			template: "Tous les conflits sont resolus. Continuer le %s ?",
			want:     true,
		},
		{
			name:     "rejects unrelated prompt",
			prompt:   "Do you really want to quit?",
			template: "All merge conflicts resolved. Continue the %s?",
			want:     false,
		},
		{
			name:     "rejects empty template",
			prompt:   "All merge conflicts resolved. Continue the rebase?",
			template: "",
			want:     false,
		},
		{
			name:     "exact match when no placeholder",
			prompt:   "Continue",
			template: "Continue",
			want:     true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := matchesPromptTemplate(tt.prompt, tt.template); got != tt.want {
				t.Fatalf("matchesPromptTemplate(%q, %q) = %v, want %v", tt.prompt, tt.template, got, tt.want)
			}
		})
	}
}
