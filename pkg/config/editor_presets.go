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

type editPreset struct {
	editTemplate              string
	editAtLineTemplate        string
	editAtLineAndWaitTemplate string
	editInTerminal            bool
}

func getPreset(osConfig *OSConfig, guessDefaultEditor func() string) *editPreset {
	presets := map[string]*editPreset{
		"vi":      standardTerminalEditorPreset("vi"),
		"vim":     standardTerminalEditorPreset("vim"),
		"nvim":    standardTerminalEditorPreset("nvim"),
		"emacs":   standardTerminalEditorPreset("emacs"),
		"nano":    standardTerminalEditorPreset("nano"),
		"kakoune": standardTerminalEditorPreset("kakoune"),
		"vscode": {
			editTemplate:              "code --reuse-window -- {{filename}}",
			editAtLineTemplate:        "code --reuse-window --goto -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "code --reuse-window --goto --wait -- {{filename}}:{{line}}",
			editInTerminal:            false,
		},
		"sublime": {
			editTemplate:              "subl -- {{filename}}",
			editAtLineTemplate:        "subl -- {{filename}}:{{line}}",
			editAtLineAndWaitTemplate: "subl --wait -- {{filename}}:{{line}}",
			editInTerminal:            false,
		},
		"bbedit": {
			editTemplate:              "bbedit -- {{filename}}",
			editAtLineTemplate:        "bbedit +{{line}} -- {{filename}}",
			editAtLineAndWaitTemplate: "bbedit +{{line}} --wait -- {{filename}}",
			editInTerminal:            false,
		},
		"xcode": {
			editTemplate:              "xed -- {{filename}}",
			editAtLineTemplate:        "xed --line {{line}} -- {{filename}}",
			editAtLineAndWaitTemplate: "xed --line {{line}} --wait -- {{filename}}",
			editInTerminal:            false,
		},
	}

	// Some of our presets have a different name than the editor they are using.
	editorToPreset := map[string]string{
		"kak":  "kakoune",
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
		editInTerminal:            true,
	}
}

func getEditInTerminal(osConfig *OSConfig, preset *editPreset) bool {
	if osConfig.EditInTerminal != nil {
		return *osConfig.EditInTerminal
	}
	return preset.editInTerminal
}
