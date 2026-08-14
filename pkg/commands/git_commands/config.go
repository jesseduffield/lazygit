package git_commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

// BranchConfig holds the tracking configuration for a branch.
type BranchConfig struct {
	Remote string
	Merge  string // short ref name of upstream branch
}

type ConfigCommands struct {
	*common.Common

	gitConfig git_config.IGitConfig

	// usesTerminalPinentryOnce lazily determines, once per instance, whether
	// the user's gpg-agent is configured to use a terminal-based pinentry
	// program (pinentry-tty/pinentry-curses).
	usesTerminalPinentryOnce func() bool

	// canUseGpgLoopbackOnce lazily determines, once per instance, whether we
	// can sign using `gpg --pinentry-mode=loopback`, which lets us prompt for
	// the passphrase in our own popup instead of handing off the terminal.
	canUseGpgLoopbackOnce func() bool
}

func NewConfigCommands(
	common *common.Common,
	gitConfig git_config.IGitConfig,
) *ConfigCommands {
	self := &ConfigCommands{
		Common:    common,
		gitConfig: gitConfig,
	}
	self.usesTerminalPinentryOnce = sync.OnceValue(usesTerminalPinentry)
	self.canUseGpgLoopbackOnce = sync.OnceValue(func() bool { return canUseGpgLoopback(self.GetGpgProgram()) })
	return self
}

type GpgConfigKey string

const (
	CommitGpgSign GpgConfigKey = "commit.gpgSign"
	TagGpgSign    GpgConfigKey = "tag.gpgSign"
)

// NeedsGpgSubprocess tells us whether the user has gpg enabled for the specified action type
// and needs a subprocess because they have a process where they manually
// enter their password every time a GPG action is taken
func (self *ConfigCommands) NeedsGpgSubprocess(key GpgConfigKey) bool {
	overrideGpg := self.UserConfig().Git.OverrideGpg
	// A terminal-based pinentry (pinentry-tty/pinentry-curses) draws directly
	// on the TTY, which corrupts Lazygit's own screen buffer unless we hand
	// off the terminal via a subprocess. In that case we must ignore
	// overrideGpg, since honoring it would break the UI.
	if overrideGpg && !self.usesTerminalPinentryOnce() {
		return false
	}

	return self.gitConfig.GetBool(string(key))
}

func (self *ConfigCommands) NeedsGpgSubprocessForCommit() bool {
	return self.NeedsGpgSubprocess(CommitGpgSign)
}

// IsGpgSignEnabled tells us whether the user has gpg signing enabled for the
// specified action type (commit.gpgSign or tag.gpgSign), independent of
// overrideGpg/subprocess considerations.
func (self *ConfigCommands) IsGpgSignEnabled(key GpgConfigKey) bool {
	return self.gitConfig.GetBool(string(key))
}

func (self *ConfigCommands) GetGpgTagSign() bool {
	return self.gitConfig.GetBool(string(TagGpgSign))
}

// GetGpgProgram returns the gpg binary git will invoke to sign things,
// respecting the user's gpg.program override and falling back to git's own
// default of "gpg".
func (self *ConfigCommands) GetGpgProgram() string {
	if program := self.gitConfig.Get("gpg.program"); program != "" {
		return program
	}

	return "gpg"
}

// CanUseGpgLoopback tells us whether we can sign using
// `gpg --pinentry-mode=loopback`, which makes gpg print a plain textual
// passphrase prompt on its own stdio instead of invoking a pinentry program.
// This lets us detect the prompt and answer it from our own popup, so
// signing never has to hand off the terminal at all. It requires GnuPG 2.1+,
// and is unavailable if the agent has been hardened with
// `no-allow-loopback-pinentry`.
func (self *ConfigCommands) CanUseGpgLoopback() bool {
	return self.canUseGpgLoopbackOnce()
}

// AddGpgLoopbackEnvVars arranges for cmdObj (expected to be a `git commit` or
// `git tag` invocation that may need to sign) to invoke gpg via
// `--pinentry-mode=loopback`, by overriding gpg.program to re-invoke lazygit
// itself as a thin wrapper (see daemon.NewGpgWrapperInstruction).
func (self *ConfigCommands) AddGpgLoopbackEnvVars(cmdObj *oscommands.CmdObj) {
	cmdObj.AddEnvVars(daemon.ToEnvVars(daemon.NewGpgWrapperInstruction(self.GetGpgProgram()))...)
	// gpg.program is not passed through a shell by git (unlike e.g.
	// GIT_EDITOR), so we must use the unquoted executable path here.
	cmdObj.AddEnvVars(
		"GIT_CONFIG_COUNT=1",
		"GIT_CONFIG_KEY_0=gpg.program",
		"GIT_CONFIG_VALUE_0="+oscommands.GetLazygitExecutablePath(),
	)
}

func (self *ConfigCommands) GetCoreEditor() string {
	return self.gitConfig.Get("core.editor")
}

// GetRemoteURL returns current repo remote url
func (self *ConfigCommands) GetRemoteURL() string {
	return self.gitConfig.Get("remote.origin.url")
}

func (self *ConfigCommands) GetShowUntrackedFiles() string {
	return self.gitConfig.Get("status.showUntrackedFiles")
}

// this determines whether the user has configured to push to the remote branch of the same name as the current or not
func (self *ConfigCommands) GetPushToCurrent() bool {
	return self.gitConfig.Get("push.default") == "current"
}

// returns the repo's branches as specified in the git config
func (self *ConfigCommands) Branches(cmd oscommands.ICmdObjBuilder) map[string]*BranchConfig {
	cmdArgs := NewGitCmd("config").
		Arg("--local", "--get-regexp", `^branch\.`).ToArgv()
	output, err := cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		// exit code 1 means no matching keys (no branches with config)
		return nil
	}

	result := make(map[string]*BranchConfig)
	for _, line := range strings.Split(output, "\n") {
		key, value, found := strings.Cut(strings.TrimSpace(line), " ")
		if !found {
			continue
		}
		// key is like "branch.<name>.remote" or "branch.<name>.merge"
		lastDot := strings.LastIndex(key, ".")
		// ignore key like branch.autosetuprebase
		if lastDot < len("branch.") {
			continue
		}
		configKey := key[lastDot+1:]
		branchName := key[len("branch."):lastDot]
		if _, ok := result[branchName]; !ok {
			result[branchName] = &BranchConfig{}
		}
		switch configKey {
		case "remote":
			result[branchName].Remote = value
		case "merge":
			result[branchName].Merge = strings.TrimPrefix(value, "refs/heads/")
		}
	}
	return result
}

// git-flow config key patterns: legacy uses gitflow.prefix.<type>, git-flow-next uses gitflow.branch.<type>.prefix
const (
	gitFlowLegacyConfigArgs = "--local --get-regexp gitflow.prefix"
	gitFlowNextConfigArgs   = "--local --get-regexp gitflow\\.branch\\..*\\.prefix"
)

func (self *ConfigCommands) getGitFlowPrefixes() string {
	return self.gitConfig.GetGeneral(gitFlowLegacyConfigArgs)
}

func (self *ConfigCommands) getGitFlowNextPrefixes() string {
	return self.gitConfig.GetGeneral(gitFlowNextConfigArgs)
}

// parseGitFlowLines parses lines matching re (submatch 1 = branch type, 2 = prefix) into prefixToType.
// When overwrite is false, existing keys are left unchanged so legacy entries win over next.
func parseGitFlowLines(output string, re *regexp.Regexp, prefixToType map[string]string, overwrite bool) {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if m := re.FindStringSubmatch(line); len(m) == 3 {
			prefix := normalizeGitFlowPrefix(m[2])
			if prefix == "" {
				continue
			}
			if overwrite || prefixToType[prefix] == "" {
				prefixToType[prefix] = m[1]
			}
		}
	}
}

// parseGitFlowPrefixMap parses legacy and git-flow-next config output into a unified prefix → branchType map.
// Legacy line format: "gitflow.prefix.<type> <prefix>"
// Next line format: "gitflow.branch.<type>.prefix <prefix>"
// Prefixes are normalized to end in "/". Legacy entries win on duplicate prefix.
func parseGitFlowPrefixMap(legacyOutput, nextOutput string) map[string]string {
	legacyRegexp := regexp.MustCompile(`gitflow\.prefix\.(\S+)\s+(.*)`)
	nextRegexp := regexp.MustCompile(`gitflow\.branch\.([^.]+)\.prefix\s+(.*)`)
	prefixToType := make(map[string]string)
	parseGitFlowLines(legacyOutput, legacyRegexp, prefixToType, true)
	parseGitFlowLines(nextOutput, nextRegexp, prefixToType, false)
	return prefixToType
}

func normalizeGitFlowPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	if !strings.HasSuffix(prefix, "/") {
		return prefix + "/"
	}
	return prefix
}

func (self *ConfigCommands) GetGitFlowPrefixMap() map[string]string {
	return parseGitFlowPrefixMap(self.getGitFlowPrefixes(), self.getGitFlowNextPrefixes())
}

func (self *ConfigCommands) GetCoreCommentChar() byte {
	if commentCharStr := self.gitConfig.Get("core.commentChar"); len(commentCharStr) == 1 {
		return commentCharStr[0]
	}

	return '#'
}

func (self *ConfigCommands) GetRebaseUpdateRefs() bool {
	return self.gitConfig.GetBool("rebase.updateRefs")
}

func (self *ConfigCommands) GetMergeFF() string {
	return self.gitConfig.Get("merge.ff")
}

func (self *ConfigCommands) DropConfigCache() {
	self.gitConfig.DropCache()
}

// usesTerminalPinentry determines whether the configured pinentry program is
// terminal-based (pinentry-tty/pinentry-curses), by asking gpgconf, falling
// back to reading gpg-agent.conf directly if gpgconf isn't available.
func usesTerminalPinentry() bool {
	program := pinentryProgramFromGpgConf()
	if program == "" {
		program = pinentryProgramFromGpgAgentConfFile(gpgAgentConfPath())
	}

	return isTerminalPinentryProgram(program)
}

// isTerminalPinentryProgram returns true if the given pinentry binary name
// looks like a terminal-based pinentry program.
func isTerminalPinentryProgram(program string) bool {
	if program == "" {
		return false
	}

	name := strings.ToLower(filepath.Base(program))
	return strings.Contains(name, "tty") || strings.Contains(name, "curses")
}

// pinentryProgramFromGpgConf asks gpgconf for the configured pinentry-program
// option of gpg-agent, returning "" if it can't be determined.
func pinentryProgramFromGpgConf() string {
	output, err := exec.Command("gpgconf", "--list-options", "gpg-agent").Output()
	if err != nil {
		return ""
	}

	return parseGpgConfPinentryProgram(string(output))
}

// parseGpgConfPinentryProgram parses the output of
// `gpgconf --list-options gpg-agent`, which consists of colon-separated
// lines of the form:
//
//	name:flags:level:description:type:alt-type:argname:default:argdef:value
//
// where "value" is percent-encoded if present, and holds the pinentry program
// path when the user has overridden it.
func parseGpgConfPinentryProgram(output string) string {
	for line := range strings.Lines(output) {
		line = strings.TrimSuffix(line, "\n")
		fields := strings.Split(line, ":")
		if len(fields) < 10 || fields[0] != "pinentry-program" {
			continue
		}

		return gpgConfUnescape(fields[9])
	}

	return ""
}

// gpgConfUnescape decodes the %XX percent-encoding used by gpgconf's
// colon-separated output format.
func gpgConfUnescape(s string) string {
	var builder strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			if b, err := strconv.ParseUint(s[i+1:i+3], 16, 8); err == nil {
				builder.WriteByte(byte(b))
				i += 2
				continue
			}
		}
		builder.WriteByte(s[i])
	}

	return builder.String()
}

// canUseGpgLoopback determines whether we can sign using
// `gpg --pinentry-mode=loopback` with the given gpg program: this requires
// GnuPG 2.1+, and requires that the agent hasn't been hardened with
// `no-allow-loopback-pinentry`.
func canUseGpgLoopback(program string) bool {
	versionOutput, err := exec.Command(program, "--version").Output()
	if err != nil {
		return false
	}

	major, minor, ok := parseGpgVersion(string(versionOutput))
	if !ok || major < 2 || (major == 2 && minor < 1) {
		return false
	}

	gpgConfOutput, err := exec.Command("gpgconf", "--list-options", "gpg-agent").Output()
	if err != nil {
		// if gpgconf isn't available we can't check for
		// no-allow-loopback-pinentry, but that option is opt-in and rare, so
		// assume loopback is fine.
		return true
	}

	return !parseGpgConfNoAllowLoopbackPinentry(string(gpgConfOutput))
}

var gpgVersionRe = regexp.MustCompile(`(?m)^gpg \(GnuPG(?:/MacGPG2)?\)\s+(\d+)\.(\d+)`)

// parseGpgVersion extracts the major and minor version numbers from the
// output of `gpg --version`, whose first line looks like
// "gpg (GnuPG) 2.2.27" or "gpg (GnuPG/MacGPG2) 2.2.27".
func parseGpgVersion(output string) (major int, minor int, ok bool) {
	match := gpgVersionRe.FindStringSubmatch(output)
	if match == nil {
		return 0, 0, false
	}

	major, errMajor := strconv.Atoi(match[1])
	minor, errMinor := strconv.Atoi(match[2])
	if errMajor != nil || errMinor != nil {
		return 0, 0, false
	}

	return major, minor, true
}

// parseGpgConfNoAllowLoopbackPinentry parses the output of
// `gpgconf --list-options gpg-agent` (see parseGpgConfPinentryProgram for the
// column format) to determine whether the agent has been hardened with
// `no-allow-loopback-pinentry`, which would make `--pinentry-mode=loopback`
// fail.
func parseGpgConfNoAllowLoopbackPinentry(output string) bool {
	for line := range strings.Lines(output) {
		line = strings.TrimSuffix(line, "\n")
		fields := strings.Split(line, ":")
		if len(fields) < 10 || fields[0] != "no-allow-loopback-pinentry" {
			continue
		}

		return gpgConfUnescape(fields[9]) == "1"
	}

	return false
}

// gpgAgentConfPath returns the path to the user's gpg-agent.conf, or "" if
// the home directory can't be determined.
func gpgAgentConfPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".gnupg", "gpg-agent.conf")
}

// pinentryProgramFromGpgAgentConfFile reads the pinentry-program setting out
// of a gpg-agent.conf file, returning "" if the file doesn't exist or doesn't
// set it.
func pinentryProgramFromGpgAgentConfFile(path string) string {
	if path == "" {
		return ""
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return parseGpgAgentConfPinentryProgram(string(content))
}

var pinentryProgramLineRe = regexp.MustCompile(`(?m)^\s*pinentry-program\s+(\S+)\s*$`)

// parseGpgAgentConfPinentryProgram extracts the value of the pinentry-program
// option from the contents of a gpg-agent.conf file.
func parseGpgAgentConfPinentryProgram(content string) string {
	match := pinentryProgramLineRe.FindStringSubmatch(content)
	if match == nil {
		return ""
	}

	return match[1]
}
