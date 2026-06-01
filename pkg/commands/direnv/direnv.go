package direnv

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// LoadResult bundles everything callers might want to know about a direnv
// invocation. The env-var delta has already been applied to the process by
// the time Load returns.
type LoadResult struct {
	// Message is whatever direnv printed to stderr — useful to log
	// (success: "direnv: loading .envrc"; error: the error text).
	Message string

	// Err is non-nil when direnv exited non-zero or its stdout could
	// not be parsed.
	Err error

	// Blocked is true when the target .envrc exists but hasn't been
	// approved with `direnv allow` yet. EnvrcPath then holds the path
	// direnv said was blocked, suitable for passing to Allow.
	Blocked   bool
	EnvrcPath string
}

// Load runs `direnv export json` for the current working directory and applies
// the resulting env-var delta to the current process. If direnv isn't on PATH,
// it's a no-op — users who don't use direnv pay nothing, and users who do need
// no config to opt in.
func Load(cmd oscommands.ICmdObjBuilder) LoadResult {
	if _, lookupErr := exec.LookPath("direnv"); lookupErr != nil {
		return LoadResult{}
	}

	stdout, stderr, runErr := cmd.New([]string{
		"direnv", "export", "json",
	}).DontLog().RunWithOutputs()

	result := LoadResult{Message: strings.TrimRight(stderr, "\n")}

	// Apply whatever delta direnv produced even if it exited non-zero.
	// When the new dir's .envrc is blocked, direnv still emits a valid
	// JSON delta on stdout that unloads vars from the previous dir;
	// without applying it the old env would leak into the new repo.
	delta, parseErr := parseDirenvExport([]byte(stdout))
	for k, v := range delta {
		if v == nil {
			_ = os.Unsetenv(k)
		} else {
			_ = os.Setenv(k, *v)
		}
	}

	// Prefer the runtime error (whose Error() text is direnv's stderr)
	// over a parse error, since it's the more actionable signal.
	if runErr != nil {
		result.Err = runErr
		if envrcPath := queryBlockedEnvrc(cmd); envrcPath != "" {
			result.Blocked = true
			result.EnvrcPath = envrcPath
		}
	} else {
		result.Err = parseErr
	}
	return result
}

// Allow runs `direnv allow <envrcPath>` to approve a .envrc file so the next
// Load can read it.
func Allow(cmd oscommands.ICmdObjBuilder, envrcPath string) error {
	return cmd.New([]string{"direnv", "allow", envrcPath}).DontLog().Run()
}

func parseDirenvExport(stdout []byte) (map[string]*string, error) {
	trimmed := bytes.TrimSpace(stdout)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	var delta map[string]*string
	if err := json.Unmarshal(trimmed, &delta); err != nil {
		return nil, err
	}
	return delta, nil
}

// queryBlockedEnvrc asks direnv (via `status --json`) whether the current
// directory has a found-but-not-yet-allowed .envrc, and returns its path
// if so. We use direnv's structured output rather than parsing the
// human-readable "is blocked" line because the status output is more
// stable across versions and locales.
func queryBlockedEnvrc(cmd oscommands.ICmdObjBuilder) string {
	stdout, _, err := cmd.New([]string{
		"direnv", "status", "--json",
	}).DontLog().RunWithOutputs()
	if err != nil {
		return ""
	}
	return parseDirenvStatus([]byte(stdout))
}

func parseDirenvStatus(stdout []byte) string {
	var status struct {
		State struct {
			FoundRC *struct {
				Allowed int    `json:"allowed"`
				Path    string `json:"path"`
			} `json:"foundRC"`
		} `json:"state"`
	}
	if err := json.Unmarshal(stdout, &status); err != nil {
		return ""
	}
	if status.State.FoundRC == nil {
		return ""
	}
	// direnv's AllowStatus enum (`internal/cmd/rc.go`): 0=Allowed,
	// 1=NotAllowed, 2=Denied. Only NotAllowed is something the user
	// can approve; Denied means they already said no.
	const notAllowed = 1
	if status.State.FoundRC.Allowed != notAllowed {
		return ""
	}
	return status.State.FoundRC.Path
}
