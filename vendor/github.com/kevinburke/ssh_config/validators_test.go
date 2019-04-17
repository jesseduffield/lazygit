package ssh_config

import (
	"testing"
)

var validateTests = []struct {
	key string
	val string
	err string
}{
	{"IdentitiesOnly", "yes", ""},
	{"IdentitiesOnly", "Yes", `ssh_config: value for key "IdentitiesOnly" must be 'yes' or 'no', got "Yes"`},
	{"Port", "22", ``},
	{"Port", "yes", `ssh_config: strconv.ParseUint: parsing "yes": invalid syntax`},
}

func TestValidate(t *testing.T) {
	for _, tt := range validateTests {
		err := validate(tt.key, tt.val)
		if tt.err == "" && err != nil {
			t.Errorf("validate(%q, %q): got %v, want nil", tt.key, tt.val, err)
		}
		if tt.err != "" {
			if err == nil {
				t.Errorf("validate(%q, %q): got nil error, want %v", tt.key, tt.val, tt.err)
			} else if err.Error() != tt.err {
				t.Errorf("validate(%q, %q): got err %v, want %v", tt.key, tt.val, err, tt.err)
			}
		}
	}
}

func TestDefault(t *testing.T) {
	if v := Default("VisualHostKey"); v != "no" {
		t.Errorf("Default(%q): got %v, want 'no'", "VisualHostKey", v)
	}
	if v := Default("visualhostkey"); v != "no" {
		t.Errorf("Default(%q): got %v, want 'no'", "visualhostkey", v)
	}
	if v := Default("notfound"); v != "" {
		t.Errorf("Default(%q): got %v, want ''", "notfound", v)
	}
}
