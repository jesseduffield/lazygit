package constants

type Docs struct {
	CustomPagers      string
	CustomCommands    string
	CustomKeybindings string
	Keybindings       string
	Undoing           string
	Config            string
	Tutorial          string
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
		CustomPagers:      "https://github.com/jesseduffield/lazygit/blob/master/docs/Custom_Pagers.md",
		CustomKeybindings: "https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md",
		CustomCommands:    "https://github.com/jesseduffield/lazygit/wiki/Custom-Commands-Compendium",
		Keybindings:       "https://github.com/jesseduffield/lazygit/blob/master/docs/keybindings",
		Undoing:           "https://github.com/jesseduffield/lazygit/blob/master/docs/Undoing.md",
		Config:            "https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md",
		Tutorial:          "https://youtu.be/VDXvbHZYeKY",
	},
}
