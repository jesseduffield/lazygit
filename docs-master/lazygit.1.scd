lazygit(1) "Lazygit Manual" "lazygit"

# NAME

lazygit - simple terminal UI for git commands

# SYNOPSIS

*lazygit* [_options_] [_git-arg_]

# DESCRIPTION

*lazygit* is a terminal user interface for interacting with git repositories.
Run it inside a git working tree and it shows status, branches, commits,
stash, reflog and other views side by side, letting you stage hunks and
individual lines, resolve conflicts, rebase interactively, cherry-pick,
bisect, manage worktrees and submodules, and invoke custom shell commands
without leaving the terminal.

*lazygit* must normally be started from inside a git repository. Use the *-p*,
*-w* and *-g* options (or the *GIT_WORK_TREE* and *GIT_DIR* environment
variables) to point it at a repository elsewhere.

Press *?* inside *lazygit* for the full keybinding menu. Keybindings are
configurable; see the configuration documentation for details.

# POSITIONAL ARGUMENTS

*git-arg*
	Panel to focus upon opening *lazygit*. Accepted values (based on git
	terminology): *status*, *branch*, *log*, *stash*. Ignored if *--filter* is
	passed.

# OPTIONS

*-p*, *--path*=_path_
	Path of git repository. Equivalent to *--work-tree*=_path_
	*--git-dir*=_path_*/.git*. Incompatible with *--work-tree* and *--git-dir*.

*-f*, *--filter*=_path_
	Path to filter on in `git log -- <path>`. When in filter mode, the commits,
	reflog and stash views are filtered to the given path, and some operations
	are restricted.

*-v*, *--version*
	Print the current version, including commit hash, build date, build source,
	operating system, architecture and git version, then exit.

*-d*, *--debug*
	Run in debug mode with logging. See *--logs* below. Use the *LOG_LEVEL*
	environment variable to set the log level (*debug*/*info*/*warn*/*error*).
	*DEBUG=TRUE* in the environment has the same effect.

*-l*, *--logs*
	Tail the *lazygit* log file and exit. Intended to be run in a separate
	terminal tab while `lazygit --debug` runs in another.

*--profile*
	Start the profiler and serve it on HTTP port 6060. See CONTRIBUTING.md for
	more information. Intended for developers.

*-c*, *--config*
	Print the default configuration as YAML, then exit.

*-cd*, *--print-config-dir*
	Print the configuration directory, then exit.

*-ucd*, *--use-config-dir*=_directory_
	Override the default configuration directory with the given directory.
	Equivalent to setting *CONFIG_DIR*.

*-w*, *--work-tree*=_path_
	Equivalent of the *--work-tree* git argument. Defaults to the value of
	*GIT_WORK_TREE*.

*-g*, *--git-dir*=_path_
	Equivalent of the *--git-dir* git argument. Defaults to the value of
	*GIT_DIR*.

*-ucf*, *--use-config-file*=_file_[,_file_...]
	Comma-separated list of custom configuration file(s) to use instead of the
	default location. Equivalent to setting *LG_CONFIG_FILE*.

*-sm*, *--screen-mode*=_mode_
	Initial screen mode, which determines the size of the focused panel. Valid
	options: *normal* (default), *half*, *full*.

*-h*, *--help*
	Show the help text and exit.

# ENVIRONMENT

*GIT_WORK_TREE*
	Default for *--work-tree*.

*GIT_DIR*
	Default for *--git-dir*.

*DEBUG*
	If set to *TRUE*, enables debug mode like *--debug*.

*LOG_LEVEL*
	Log level when running with *--debug*: *debug*, *info*, *warn* or *error*.
	Defaults to *debug* when unset or unparsable.

*LG_CONFIG_FILE*
	Comma-separated list of custom configuration file(s). Set implicitly by
	*--use-config-file*.

*CONFIG_DIR*
	Override the configuration (and log) directory. Set implicitly by
	*--use-config-dir*.

*LAZYGIT_NEW_DIR_FILE*
	If set, *lazygit* writes its current directory to this file on exit so that
	the calling shell can change to it (see EXAMPLES). Exiting with *Shift+Q*
	records the original directory instead, so the shell stays where it started.

*LAZYGIT_LOG_PATH*
	Override the log file path read by *--logs* and written by *--debug*.
	Defaults to *development.log* in the state directory.

*XDG_CONFIG_HOME*, *XDG_CONFIG_DIRS*
	Searched for the *lazygit/config.yml* configuration file (and the legacy
	*jesseduffield/lazygit/config.yml* location).

*XDG_STATE_HOME*
	Location of persistent state (*state.yml*) and the default log file when no
	existing state file is found.

# FILES

*~/.config/lazygit/config.yml* (Linux)
	Global configuration file. Only settings you want to change need to be
	included; defaults are used for everything else. Use `lazygit --config` to
	view all available options with their defaults. A JSON schema is published
	for editor completion; see the configuration docs.

*~/Library/Application Support/lazygit/config.yml* (macOS)
	Global configuration file (see above).

*%LOCALAPPDATA%\\lazygit\\config.yml* (Windows)
	Global configuration file (see above).

Legacy global locations (still honoured if present):
	*~/.config/jesseduffield/lazygit/config.yml* (Linux),
	*~/Library/Application Support/jesseduffield/lazygit/config.yml* (macOS),
	*%APPDATA%\\jesseduffield\\lazygit\\config.yml* (Windows).

_<repo>/.git/lazygit.yml_
	Repository-specific configuration file. Settings here override the global
	configuration file.

_.lazygit.yml_ in any parent directory of a repository
	Also loaded, with settings nearer the repository taking precedence. Useful
	for applying settings to a group of repositories.

# EXIT STATUS

*0*
	Successful execution.

Non-zero
	Fatal error, for example an invalid repository path, an invalid _git-arg_
	value, an unreadable configuration file, or an unwritable temporary
	directory. Quitting the interface normally (for example with *q*) exits *0*.

# EXAMPLES

Open *lazygit* in the current repository:

	$ lazygit

Open focused on the branches view:

	$ lazygit branch

Open a repository by path:

	$ lazygit -p /path/to/repo

Show only commits touching a path:

	$ lazygit -f path/to/file

# SEE ALSO

*git*(1)

Configuration reference: configuration docs (*docs/Config.md* in a release,
*docs-master/Config.md* in the source tree).

Keybindings: *docs/keybindings* (release) or *docs-master/keybindings* (source),
and the *?* menu inside *lazygit*.

Project home: https://github.com/jesseduffield/lazygit

# BUGS

Report issues at https://github.com/jesseduffield/lazygit/issues. Contributor
discussion takes place on Discord; see the contributing guide.

# AUTHORS

Jesse Duffield and contributors. See the contributor graph at
https://github.com/jesseduffield/lazygit/graphs/contributors.

# COPYRIGHT

MIT License. See the LICENSE file in the distribution.
