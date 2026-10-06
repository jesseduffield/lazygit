// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux netbsd openbsd solaris dragonfly

package clipboard

import (
	"errors"
	"os"
	"os/exec"
	"unicode/utf16"
)

const (
	xsel               = "xsel"
	xclip              = "xclip"
	powershellExe      = "powershell.exe"
	clipExe            = "clip.exe"
	wlcopy             = "wl-copy"
	wlpaste            = "wl-paste"
	termuxClipboardGet = "termux-clipboard-get"
	termuxClipboardSet = "termux-clipboard-set"
)

var (
	Primary   bool
	trimDos   bool
	copyUtf16 bool

	pasteCmdArgs []string
	copyCmdArgs  []string

	xselPasteArgs = []string{xsel, "--output", "--clipboard"}
	xselCopyArgs  = []string{xsel, "--input", "--clipboard"}

	xclipPasteArgs = []string{xclip, "-out", "-selection", "clipboard"}
	xclipCopyArgs  = []string{xclip, "-in", "-selection", "clipboard"}

	powershellExePasteArgs = []string{powershellExe, "-noprofile", "-command", "[Console]::OutputEncoding=[System.Text.Encoding]::UTF8; Get-Clipboard"}
	clipExeCopyArgs        = []string{clipExe}

	wlpasteArgs = []string{wlpaste, "--no-newline"}
	wlcopyArgs  = []string{wlcopy}

	termuxPasteArgs = []string{termuxClipboardGet}
	termuxCopyArgs  = []string{termuxClipboardSet}

	missingCommands = errors.New("No clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set.")
)

func init() {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		pasteCmdArgs = wlpasteArgs
		copyCmdArgs = wlcopyArgs

		if _, err := exec.LookPath(wlcopy); err == nil {
			if _, err := exec.LookPath(wlpaste); err == nil {
				return
			}
		}
	}

	pasteCmdArgs = xclipPasteArgs
	copyCmdArgs = xclipCopyArgs

	if _, err := exec.LookPath(xclip); err == nil {
		return
	}

	pasteCmdArgs = xselPasteArgs
	copyCmdArgs = xselCopyArgs

	if _, err := exec.LookPath(xsel); err == nil {
		return
	}

	pasteCmdArgs = termuxPasteArgs
	copyCmdArgs = termuxCopyArgs

	if _, err := exec.LookPath(termuxClipboardSet); err == nil {
		if _, err := exec.LookPath(termuxClipboardGet); err == nil {
			return
		}
	}

	pasteCmdArgs = powershellExePasteArgs
	copyCmdArgs = clipExeCopyArgs
	trimDos = true
	copyUtf16 = true

	if _, err := exec.LookPath(clipExe); err == nil {
		if _, err := exec.LookPath(powershellExe); err == nil {
			return
		}
	}

	Unsupported = true
}

func getPasteCommand() *exec.Cmd {
	if Primary {
		pasteCmdArgs = pasteCmdArgs[:1]
	}
	return exec.Command(pasteCmdArgs[0], pasteCmdArgs[1:]...)
}

func getCopyCommand() *exec.Cmd {
	if Primary {
		copyCmdArgs = copyCmdArgs[:1]
	}
	return exec.Command(copyCmdArgs[0], copyCmdArgs[1:]...)
}

// clip.exe interprets its input in the console's ANSI code page unless it
// starts with a UTF-16LE byte order mark, so encode the text that way to keep
// non-ASCII characters intact.
func utf16LeWithBom(text string) []byte {
	codes := utf16.Encode([]rune(text))
	data := make([]byte, 0, len(codes)*2+2)
	data = append(data, 0xff, 0xfe)
	for _, c := range codes {
		data = append(data, byte(c), byte(c>>8))
	}
	return data
}

func readAll() (string, error) {
	if Unsupported {
		return "", missingCommands
	}
	pasteCmd := getPasteCommand()
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	result := string(out)
	if trimDos && len(result) > 1 {
		result = result[:len(result)-2]
	}
	return result, nil
}

func writeAll(text string) error {
	if Unsupported {
		return missingCommands
	}
	copyCmd := getCopyCommand()
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	data := []byte(text)
	if copyUtf16 {
		data = utf16LeWithBom(text)
	}
	if _, err := in.Write(data); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
