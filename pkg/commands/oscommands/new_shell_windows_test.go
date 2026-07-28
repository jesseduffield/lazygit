//go:build windows

package oscommands

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// These tests run only on Windows because they exercise real cmd.exe
// quote-parsing behaviour, which has only been a problem on Windows.

// makeWindowsShellBuilder returns a CmdObjBuilder configured for a real
// Windows cmd shell, bypassing the test "dummy" platform (which is darwin).
func makeWindowsShellBuilder() *CmdObjBuilder {
	log := utils.NewDummyLog()
	return &CmdObjBuilder{
		runner:   &cmdObjRunner{log: log, guiIO: NewNullGuiIO(log)},
		platform: &Platform{OS: "windows", Shell: "cmd", ShellArg: "/c"},
	}
}

// fakeEditorSrc is a minimal Go program that records the args it received,
// one per line, to marker.txt in its own directory. Using a real .exe (not a
// .bat) means args are parsed by Go's runtime via CommandLineToArgvW — the
// same algorithm used by ~all real Windows GUI editors. A .bat would parse
// args via cmd.exe's own rules, which can hide bugs that affect editors.
const fakeEditorSrc = `package main

import (
	"os"
	"path/filepath"
	"strings"
)

func main() {
	exe, err := os.Executable()
	if err != nil {
		os.Exit(2)
	}
	marker := filepath.Join(filepath.Dir(exe), "marker.txt")
	body := strings.Join(os.Args[1:], "\n")
	if err := os.WriteFile(marker, []byte(body), 0o644); err != nil {
		os.Exit(3)
	}
}
`

var (
	fakeEditorOnce  sync.Once
	fakeEditorBytes []byte
	fakeEditorErr   error
)

// loadFakeEditorBytes builds the fake editor exactly once per test process
// and returns its bytes. Tests then drop a copy at a path containing spaces.
func loadFakeEditorBytes(t *testing.T) []byte {
	t.Helper()
	fakeEditorOnce.Do(func() {
		buildDir, err := os.MkdirTemp("", "lazygit-fake-editor-build-*")
		if err != nil {
			fakeEditorErr = err
			return
		}
		defer os.RemoveAll(buildDir)

		srcPath := filepath.Join(buildDir, "main.go")
		binPath := filepath.Join(buildDir, "fake-editor.exe")
		if err := os.WriteFile(srcPath, []byte(fakeEditorSrc), 0o644); err != nil {
			fakeEditorErr = err
			return
		}
		cmd := exec.Command("go", "build", "-o", binPath, srcPath)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fakeEditorErr = err
			return
		}
		fakeEditorBytes, fakeEditorErr = os.ReadFile(binPath)
	})
	if fakeEditorErr != nil {
		t.Fatalf("failed to build fake editor helper: %v", fakeEditorErr)
	}
	return fakeEditorBytes
}

// placeFakeEditor builds the fake editor and places it at a path containing
// a space (mirroring `C:\Program Files\...`). The marker the editor writes
// lives next to the exe.
func placeFakeEditor(t *testing.T) (exe, markerFile string) {
	t.Helper()
	bin := loadFakeEditorBytes(t)
	exeDir := filepath.Join(t.TempDir(), "Program Files", "FakeEditor")
	if err := os.MkdirAll(exeDir, 0o755); err != nil {
		t.Fatalf("mkdir exeDir: %v", err)
	}
	exe = filepath.Join(exeDir, "fake-editor.exe")
	markerFile = filepath.Join(exeDir, "marker.txt")
	if err := os.WriteFile(exe, bin, 0o755); err != nil {
		t.Fatalf("write fake editor: %v", err)
	}
	return exe, markerFile
}

// placeTargetFile creates a file at <freshTempDir>/<dirName>/<basename> with
// a trivial body. Use a dirName containing a space (e.g. "my repo") to put
// the file at a path with spaces.
func placeTargetFile(t *testing.T, dirName, basename string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir target dir: %v", err)
	}
	target := filepath.Join(dir, basename)
	if err := os.WriteFile(target, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}
	return target
}

// setupFakeEditor is a convenience wrapper for the common case: editor at a
// spacey path AND target file at a spacey path — the conditions that
// trigger the cmd.exe quote-stripping bug.
func setupFakeEditor(t *testing.T) (fakeExe, targetFile, markerFile string) {
	t.Helper()
	fakeExe, markerFile = placeFakeEditor(t)
	targetFile = placeTargetFile(t, "my repo", "file.txt")
	return fakeExe, targetFile, markerFile
}

// resolveTemplate mirrors what pkg/commands/git_commands/file.go does: it
// substitutes {{filename}} with the Windows-quoted filename and {{line}}
// with a line number.
func resolveTemplate(builder *CmdObjBuilder, template, filename, line string) string {
	out := strings.ReplaceAll(template, "{{filename}}", builder.Quote(filename))
	out = strings.ReplaceAll(out, "{{line}}", line)
	return out
}

// readMarkerArgs reads the args the fake editor recorded. Each arg is on its
// own line (so an arg that itself contains a space stays one element). An
// empty file means the editor ran with zero args.
func readMarkerArgs(t *testing.T, markerFile string) []string {
	t.Helper()
	data, err := os.ReadFile(markerFile)
	if err != nil {
		t.Fatalf("marker file was not written; the fake editor never ran: %v", err)
	}
	s := strings.TrimRight(string(data), "\r\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func TestNewShell_QuotedExePath_FilenameWithSpaces_EditAtLineAndWait(t *testing.T) {
	builder := makeWindowsShellBuilder()
	fakeExe, targetFile, markerFile := setupFakeEditor(t)

	template := `"` + fakeExe + `" -multiInst -nosession -noPlugin -n{{line}} {{filename}}`
	cmdStr := resolveTemplate(builder, template, targetFile, "42")

	out, err := builder.NewShell(cmdStr, "").GetCmd().CombinedOutput()
	if err != nil {
		t.Fatalf("shell command failed: %v\ncmd.exe output:\n%s", err, string(out))
	}

	assert.Equal(t,
		[]string{"-multiInst", "-nosession", "-noPlugin", "-n42", targetFile},
		readMarkerArgs(t, markerFile),
	)
}

func TestNewShell_QuotedExePath_FilenameWithSpaces_EditAtLine(t *testing.T) {
	builder := makeWindowsShellBuilder()
	fakeExe, targetFile, markerFile := setupFakeEditor(t)

	template := `"` + fakeExe + `" -n{{line}} {{filename}}`
	cmdStr := resolveTemplate(builder, template, targetFile, "42")

	out, err := builder.NewShell(cmdStr, "").GetCmd().CombinedOutput()
	if err != nil {
		t.Fatalf("shell command failed: %v\ncmd.exe output:\n%s", err, string(out))
	}

	assert.Equal(t,
		[]string{"-n42", targetFile},
		readMarkerArgs(t, markerFile),
	)
}

func TestNewShell_QuotedExePath_FilenameWithSpaces_Edit(t *testing.T) {
	builder := makeWindowsShellBuilder()
	fakeExe, targetFile, markerFile := setupFakeEditor(t)

	template := `"` + fakeExe + `" {{filename}}`
	cmdStr := resolveTemplate(builder, template, targetFile, "")

	out, err := builder.NewShell(cmdStr, "").GetCmd().CombinedOutput()
	if err != nil {
		t.Fatalf("shell command failed: %v\ncmd.exe output:\n%s", err, string(out))
	}

	assert.Equal(t,
		[]string{targetFile},
		readMarkerArgs(t, markerFile),
	)
}

// Sanity check: for a filename WITHOUT spaces the same templates already work,
// because the resulting cmd.exe line has exactly two quote characters and
// cmd /c keeps them. This pins the difference down to filename quoting and
// guards against a regression where the no-spaces case starts failing too.
func TestNewShell_QuotedExePath_FilenameWithoutSpaces_StillWorks(t *testing.T) {
	builder := makeWindowsShellBuilder()
	fakeExe, markerFile := placeFakeEditor(t)
	plainTarget := placeTargetFile(t, "repo", "plain.txt") // no-space dir

	template := `"` + fakeExe + `" -multiInst -nosession -noPlugin -n{{line}} {{filename}}`
	cmdStr := resolveTemplate(builder, template, plainTarget, "42")

	out, err := builder.NewShell(cmdStr, "").GetCmd().CombinedOutput()
	if err != nil {
		t.Fatalf("shell command failed: %v\ncmd.exe output:\n%s", err, string(out))
	}

	assert.Equal(t,
		[]string{"-multiInst", "-nosession", "-noPlugin", "-n42", plainTarget},
		readMarkerArgs(t, markerFile),
	)
}

// TestNewShell_VarietyOfEditorTemplates exercises NewShell with a range of
// realistic editor templates, all with the trigger conditions of the bug
// (quoted exe at a spacey path + filename at a spacey path). Each subtest
// asserts the editor receives the exact args lazygit intended.
//
// Args in `wantArgs` may use the literal "<file>" placeholder; it gets
// substituted with the resolved target file path before comparison.
func TestNewShell_VarietyOfEditorTemplates(t *testing.T) {
	const filePlaceholder = "<file>"

	cases := []struct {
		name     string
		template string // <exe> stands for the fake editor's full path
		line     string
		wantArgs []string
	}{
		{
			name:     "vim/nvim style: +line filename",
			template: `"<exe>" +{{line}} {{filename}}`,
			line:     "42",
			wantArgs: []string{"+42", filePlaceholder},
		},
		{
			name:     "emacs-like with explicit +N",
			template: `"<exe>" +{{line}} -nw {{filename}}`,
			line:     "7",
			wantArgs: []string{"+7", "-nw", filePlaceholder},
		},
		{
			name:     "long flag with =value",
			template: `"<exe>" --line={{line}} --tab-size=4 {{filename}}`,
			line:     "42",
			wantArgs: []string{"--line=42", "--tab-size=4", filePlaceholder},
		},
		{
			name:     "many short and long flags before filename",
			template: `"<exe>" -a -b -c --foo --bar -n{{line}} {{filename}}`,
			line:     "42",
			wantArgs: []string{"-a", "-b", "-c", "--foo", "--bar", "-n42", filePlaceholder},
		},
		{
			name:     "flag after filename",
			template: `"<exe>" {{filename}} --readonly`,
			line:     "",
			wantArgs: []string{filePlaceholder, "--readonly"},
		},
		{
			name:     "single short flag attached to value",
			template: `"<exe>" -n{{line}} {{filename}}`,
			line:     "1",
			wantArgs: []string{"-n1", filePlaceholder},
		},
		{
			name:     "no flags, just filename",
			template: `"<exe>" {{filename}}`,
			line:     "",
			wantArgs: []string{filePlaceholder},
		},
		{
			name:     "flag with separate value (space-separated)",
			template: `"<exe>" --goto {{line}} {{filename}}`,
			line:     "42",
			wantArgs: []string{"--goto", "42", filePlaceholder},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			builder := makeWindowsShellBuilder()
			fakeExe, targetFile, markerFile := setupFakeEditor(t)

			template := strings.ReplaceAll(tc.template, "<exe>", fakeExe)
			cmdStr := resolveTemplate(builder, template, targetFile, tc.line)

			out, err := builder.NewShell(cmdStr, "").GetCmd().CombinedOutput()
			if err != nil {
				t.Fatalf("shell command failed: %v\ncmd.exe output:\n%s", err, string(out))
			}

			want := make([]string, len(tc.wantArgs))
			for i, a := range tc.wantArgs {
				want[i] = strings.ReplaceAll(a, filePlaceholder, targetFile)
			}
			assert.Equal(t, want, readMarkerArgs(t, markerFile))
		})
	}
}

// TestNewShell_FilenameSpecialCharacters varies the basename of the target
// file across characters that are legal in Windows filenames but might
// interact badly with cmd.exe / Quote(): parentheses, brackets, single
// quote, comma, semicolon, equals, etc. The exe is at a spacey path and
// the target dir has spaces, so the bug-trigger conditions are still met.
func TestNewShell_FilenameSpecialCharacters(t *testing.T) {
	cases := []struct {
		name     string
		basename string
	}{
		{"parens", "file (1).txt"},
		{"brackets", "file[v2].txt"},
		{"single quote", "it's a file.txt"},
		{"comma", "a,b,c.txt"},
		{"semicolon", "a;b.txt"},
		{"equals", "key=value.txt"},
		{"plus", "a+b.txt"},
		{"hash", "issue#42.txt"},
		{"at sign", "user@host.txt"},
		{"tilde", "~backup.txt"},
		{"dot leading", ".gitignore.txt"},
		{"multiple dots", "v1.2.3.txt"},
		{"dash leading", "-flag-looking.txt"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			builder := makeWindowsShellBuilder()
			fakeExe, markerFile := placeFakeEditor(t)
			targetFile := placeTargetFile(t, "my repo", tc.basename)

			template := `"` + fakeExe + `" -n{{line}} {{filename}}`
			cmdStr := resolveTemplate(builder, template, targetFile, "42")

			out, err := builder.NewShell(cmdStr, "").GetCmd().CombinedOutput()
			if err != nil {
				t.Fatalf("shell command failed: %v\ncmd.exe output:\n%s", err, string(out))
			}

			assert.Equal(t,
				[]string{"-n42", targetFile},
				readMarkerArgs(t, markerFile),
			)
		})
	}
}

// Command chaining with && must work: cmd /s /c runs the assembled line verbatim,
// so cmd treats && as a separator and runs both commands. The two echoes
// therefore produce two separate output lines.
func TestNewShell_CommandChaining(t *testing.T) {
	builder := makeWindowsShellBuilder()

	out, err := builder.NewShell("echo first&&echo second", "").GetCmd().CombinedOutput()
	if err != nil {
		t.Fatalf("shell command failed: %v\ncmd.exe output:\n%s", err, string(out))
	}

	normalized := strings.ReplaceAll(string(out), "\r\n", "\n")
	lines := strings.Split(strings.TrimSpace(normalized), "\n")
	assert.Equal(t, []string{"first", "second"}, lines)
}
