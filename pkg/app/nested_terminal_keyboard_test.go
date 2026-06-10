package app

import "testing"

func TestShouldForceWin32KeyboardProtocol(t *testing.T) {
	tests := []struct {
		name string
		goos string
		env  map[string]string
		want bool
	}{
		{
			name: "linux ignores nested editor terminals",
			goos: "linux",
			env:  map[string]string{"VIM_TERMINAL": "1"},
			want: false,
		},
		{
			name: "windows vim terminal",
			goos: "windows",
			env:  map[string]string{"VIM_TERMINAL": "1"},
			want: true,
		},
		{
			name: "windows neovim terminal",
			goos: "windows",
			env:  map[string]string{"NVIM": "/tmp/nvim.sock"},
			want: true,
		},
		{
			name: "windows standalone terminal",
			goos: "windows",
			env:  map[string]string{},
			want: false,
		},
		{
			name: "respects explicit override",
			goos: "windows",
			env: map[string]string{
				"VIM_TERMINAL":            "1",
				"TCELL_KEYBOARD_PROTOCOL": "legacy",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string {
				return tt.env[key]
			}
			if got := shouldForceWin32KeyboardProtocol(tt.goos, getenv); got != tt.want {
				t.Fatalf("shouldForceWin32KeyboardProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}