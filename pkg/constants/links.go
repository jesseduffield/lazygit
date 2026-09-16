package constants

type Docs struct {
	CustomDiffRenderers string
	CustomCommands      string
	CustomKeybindings   string
	Keybindings         string
	Undoing             string
	Config              string
	Tutorial            string
	CustomPatchDemo     string
}

var Links = struct {
	Docs        Docs
	Issues      string
	Donate      string
	Discussions string
	RepoUrl     string
	Releases    string
}{
	RepoUrl:     "https://github.com/jesseduffield/lazygit",
	Issues:      "https://github.com/jesseduffield/lazygit/issues",
	Donate:      "https://github.com/sponsors/jesseduffield",
	Discussions: "https://github.com/jesseduffield/lazygit/discussions",
	Releases:    "https://github.com/jesseduffield/lazygit/releases",
	Docs: Docs{
		CustomDiffRenderers: "https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_DiffRenderers.md",
		CustomKeybindings:   "https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md",
		CustomCommands:      "https://github.com/jesseduffield/lazygit/wiki/Custom-Commands-Compendium",
		Keybindings:         "https://github.com/jesseduffield/lazygit/blob/%s/docs/keybindings",
		Undoing:             "https://github.com/jesseduffield/lazygit/blob/master/docs/Undoing.md",
		Config:              "https://github.com/jesseduffield/lazygit/blob/%s/docs/Config.md",
		Tutorial:            "https://youtu.be/VDXvbHZYeKY",
		CustomPatchDemo:     "https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches",
	},
}
