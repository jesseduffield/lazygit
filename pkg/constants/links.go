package constants

type Docs struct {
	CustomPagers      string
	CustomCommands    string
	CustomKeybindings string
	Keybindings       string
	Undoing           string
	Config            string
	Tutorial          string
	CustomPatchDemo   string
}

var Links = struct {
	Docs        Docs
	Issues      string
	GitHub      string
	Discussions string
	RepoUrl     string
	Releases    string
}{
	RepoUrl:     "https://github.com/lobes/lazytask",
	Issues:      "https://github.com/lobes/lazytask/issues",
	GitHub:      "https://github.com/lobes/lazytask",
	Discussions: "https://github.com/lobes/lazytask/discussions",
	Releases:    "https://github.com/lobes/lazytask/releases",
	Docs: Docs{
		CustomPagers:      "https://github.com/lobes/lazytask/blob/master/docs/Custom_Pagers.md",
		CustomKeybindings: "https://github.com/lobes/lazytask/blob/master/docs/keybindings/Custom_Keybindings.md",
		CustomCommands:    "https://github.com/lobes/lazytask/wiki/Custom-Commands-Compendium",
		Keybindings:       "https://github.com/lobes/lazytask/blob/%s/docs/keybindings",
		Undoing:           "https://github.com/lobes/lazytask/blob/master/docs/Undoing.md",
		Config:            "https://github.com/lobes/lazytask/blob/%s/docs/Config.md",
		Tutorial:          "https://youtu.be/VDXvbHZYeKY",
		CustomPatchDemo:   "https://github.com/lobes/lazytask#rebase-magic-custom-patches",
	},
}
