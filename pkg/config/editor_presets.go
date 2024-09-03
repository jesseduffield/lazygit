package config

import "os"

func GetEditTemplate(osConfig *OSConfig, guessDefaultEditor func() string) (string, bool) {
	preset := getPreset(osConfig, guessDefaultEditor)
	template := osConfig.Edit
	if template == "" {
		template = preset.editTemplate
	}

	return template, getEditInTerminal(osConfig, preset)
}

func GetEditAtLineTemplate(osConfig *OSConfig, guessDefaultEditor func() string) (string, bool) {
	preset := getPreset(osConfig, guessDefaultEditor)
	template := osConfig.EditAtLine
	if template == "" {
		template = preset.editAtLineTemplate
	}
	return template, getEditInTerminal(osConfig, preset)
}

func GetEditAtLineAndWaitTemplate(osConfig *OSConfig, guessDefaultEditor func() string) string {
	preset := getPreset(osConfig, guessDefaultEditor)
	template := osConfig.EditAtLineAndWait
	if template == "" {
		template = preset.editAtLineAndWaitTemplate
	}
	return template
}

func GetOpenDirInEditorTemplate(osConfig *OSConfig, guessDefaultEditor func() string) (string, bool) {
	preset := getPreset(osConfig, guessDefaultEditor)
	template := osConfig.OpenDirInEditor
	if template == "" {
		template = preset.openDirInEditorTemplate
	}
	return template, getEditInTerminal(osConfig, preset)
}

type editPreset struct {
	editTemplate              string
	editAtLineTemplate        string
	editAtLineAndWaitTemplate string
	openDirInEditorTemplate   string
	suspend                   func() bool
}

func returnBool(a bool) func() bool { return (func() bool { return a }) }

// IF YOU ADD A PRESET TO THIS FUNCTION YOU MUST UPDATE THE `Supported presets` SECTION OF docs/Config.md
func getPreset(osConfig *OSConfig, guessDefaultEditor func() string) *editPreset {
	presets := map[string]*editPreset{
		"vi":   standardTerminalEditorPreset("vi"),
		"vim":  standardTerminalEditorPreset("vim"),
		"nvim": standardTerminalEditorPreset("nvim"),
		"nvim-remote": {
			editTemplate:       `[ -z "$NVIM" ] && (nvim -- {{filename}}) || (nvim --server "$NVIM" --remote-send "q" && nvim --server "$NVIM" --remote-tab {{filename}})`,
			editAtLineTemplate: `[ -z "$NVIM" ] && (nvim +{{line}} -- {{filename}}) || (nvim --server "$NVIM" --remote-send "q" &&  nvim --server "$NVIM" --remote-tab {{filename}} && nvim --server "$NVIM" --remote-send ":{{line}}<CR>")`,
			// No remote-wait support yet. See https://github.com/neovim/neovim/pull/17856
			editAtLineAndWaitTemplate: `nvim +{{line}} {{filename}}`,
			openDirInEditorTemplate:   `[ -z "$NVIM" ] && (nvim -- {{dir}}) || (nvim --server "$NVIM" --remote-send "q" && nvim --server "$NVIM" --remote-tab {{dir}})`,
			suspend: func() bool {
				_, ok := os.LookupEnv("NVIM")
				return !ok
			},
		},
		"lvim":    standardTerminalEditorPreset("lvim"),
		"emacs":   standardTerminalEditorPreset("emacs"),
		"micro":   standardTerminalEditorPreset("micro"),
		"nano":    standardTerminalEditorPreset("nano"),
		"kakoune": standardTerminalEditorPreset("kak"),
		"helix": {
			editTemplate:              "helix -- {{filename}}",
			editAtLineTemplate:        "helix -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "helix -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "helix -- {{dir}}",
			suspend:                   returnBool(true),
		},
		"helix (hx)": {
			editTemplate:              "hx -- {{filename}}",
			editAtLineTemplate:        "hx -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "hx -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "hx -- {{dir}}",
			suspend:                   returnBool(true),
		},
		"vscode": {
			editTemplate:              "code --reuse-window -- {{filename}}",
			editAtLineTemplate:        "code --reuse-window --goto -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "code --reuse-window --goto --wait -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "code -- {{dir}}",
			suspend:                   returnBool(false),
		},
		"sublime": {
			editTemplate:              "subl -- {{filename}}",
			editAtLineTemplate:        "subl -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "subl --wait -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "subl -- {{dir}}",
			suspend:                   returnBool(false),
		},
		"bbedit": {
			editTemplate:              "bbedit -- {{filename}}",
			editAtLineTemplate:        "bbedit +{{line}} -- {{filename}}",
			editAtLineAndWaitTemplate: "bbedit +{{line}} --wait -- {{filename}}",
			openDirInEditorTemplate:   "bbedit -- {{dir}}",
			suspend:                   returnBool(false),
		},
		"xcode": {
			editTemplate:              "xed -- {{filename}}",
			editAtLineTemplate:        "xed --line {{line}} -- {{filename}}",
			editAtLineAndWaitTemplate: "xed --line {{line}} --wait -- {{filename}}",
			openDirInEditorTemplate:   "xed -- {{dir}}",
			suspend:                   returnBool(false),
		},
		"zed": {
			editTemplate:              "zed -- {{filename}}",
			editAtLineTemplate:        "zed -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "zed --wait -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "zed -- {{dir}}",
			suspend:                   returnBool(false),
		},
	}

	// Some of our presets have a different name than the editor they are using.
	editorToPreset := map[string]string{
		"kak":   "kakoune",
		"helix": "helix",
		"hx":    "helix (hx)",
		"code":  "vscode",
		"subl":  "sublime",
		"xed":   "xcode",
	}

	presetName := osConfig.EditPreset
	if presetName == "" {
		defaultEditor := guessDefaultEditor()
		if presets[defaultEditor] != nil {
			presetName = defaultEditor
		} else if p := editorToPreset[defaultEditor]; p != "" {
			presetName = p
		}
	}

	if presetName == "" || presets[presetName] == nil {
		presetName = "vim"
	}

	return presets[presetName]
}

func standardTerminalEditorPreset(editor string) *editPreset {
	return &editPreset{
		editTemplate:              editor + " -- {{filename}}",
		editAtLineTemplate:        editor + " +{{line}} -- {{filename}}",
		editAtLineAndWaitTemplate: editor + " +{{line}} -- {{filename}}",
		openDirInEditorTemplate:   editor + " -- {{dir}}",
		suspend:                   returnBool(true),
	}
}

func getEditInTerminal(osConfig *OSConfig, preset *editPreset) bool {
	if osConfig.SuspendOnEdit != nil {
		return *osConfig.SuspendOnEdit
	}
	return preset.suspend()
}
