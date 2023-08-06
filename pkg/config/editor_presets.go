package config

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

func GetOpenDirInEditorTemplate(osConfig *OSConfig, guessDefaultEditor func() string) string {
	preset := getPreset(osConfig, guessDefaultEditor)
	template := osConfig.OpenDirInEditor
	if template == "" {
		template = preset.openDirInEditorTemplate
	}
	return template
}

type editPreset struct {
	editTemplate              string
	editAtLineTemplate        string
	editAtLineAndWaitTemplate string
	openDirInEditorTemplate   string
	editInTerminal            bool
}

// IF YOU ADD A PRESET TO THIS FUNCTION YOU MUST UPDATE THE `Supported presets` SECTION OF docs/Config.md
func getPreset(osConfig *OSConfig, guessDefaultEditor func() string) *editPreset {
	presets := map[string]*editPreset{
		"vi":      standardTerminalEditorPreset("vi"),
		"vim":     standardTerminalEditorPreset("vim"),
		"nvim":    standardTerminalEditorPreset("nvim"),
		"emacs":   standardTerminalEditorPreset("emacs"),
		"nano":    standardTerminalEditorPreset("nano"),
		"kakoune": standardTerminalEditorPreset("kak"),
		"helix": {
			editTemplate:              "hx -- {{filename}}",
			editAtLineTemplate:        "hx -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "hx -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "hx -- {{dir}}",
			editInTerminal:            true,
		},
		"emacsclient": {
			editTemplate:              "emacsclient --create-frame --alternate-editor="" {{filename}}",
			editAtLineTemplate:        "emacsclient --create-frame --alternate-editor="" +{{line}} {{filename}}",
			editAtLineAndWaitTemplate: "emacsclient --create-frame --alternate-editor="" +{{line}} {{filename}}",
			editInTerminal:            true,
		},
		"emacsclient_tty": {
			editTemplate:              "emacsclient --create-frame --tty --alternate-editor="" {{filename}}",
			editAtLineTemplate:        "emacsclient --create-frame --tty --alternate-editor="" +{{line}} {{filename}}",
			editAtLineAndWaitTemplate: "emacsclient --create-frame --tty --alternate-editor="" +{{line}} {{filename}}",
			editInTerminal:            true,
		},
		"vscode": {
			editTemplate:              "code --reuse-window -- {{filename}}",
			editAtLineTemplate:        "code --reuse-window --goto -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "code --reuse-window --goto --wait -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "code -- {{dir}}",
			editInTerminal:            false,
		},
		"sublime": {
			editTemplate:              "subl -- {{filename}}",
			editAtLineTemplate:        "subl -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "subl --wait -- {{filename}}:{{line}}",
			openDirInEditorTemplate:   "subl -- {{dir}}",
			editInTerminal:            false,
		},
		"bbedit": {
			editTemplate:              "bbedit -- {{filename}}",
			editAtLineTemplate:        "bbedit +{{line}} -- {{filename}}",
			editAtLineAndWaitTemplate: "bbedit +{{line}} --wait -- {{filename}}",
			openDirInEditorTemplate:   "bbedit -- {{dir}}",
			editInTerminal:            false,
		},
		"xcode": {
			editTemplate:              "xed -- {{filename}}",
			editAtLineTemplate:        "xed --line {{line}} -- {{filename}}",
			editAtLineAndWaitTemplate: "xed --line {{line}} --wait -- {{filename}}",
			openDirInEditorTemplate:   "xed -- {{dir}}",
			editInTerminal:            false,
		},
	}

	// Some of our presets have a different name than the editor they are using.
	editorToPreset := map[string]string{
		"kak":  "kakoune",
		"hx":   "helix",
		"code": "vscode",
		"subl": "sublime",
		"xed":  "xcode",
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
		editInTerminal:            true,
	}
}

func getEditInTerminal(osConfig *OSConfig, preset *editPreset) bool {
	if osConfig.EditInTerminal != nil {
		return *osConfig.EditInTerminal
	}
	return preset.editInTerminal
}
