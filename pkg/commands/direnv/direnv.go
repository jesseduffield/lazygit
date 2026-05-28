package direnv

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// Load runs `direnv export json` for the current working directory and applies
// the resulting env-var delta to the current process. If direnv isn't on PATH,
// it's a no-op — users who don't use direnv pay nothing, and users who do need
// no config to opt in.
//
// direnv prints diagnostics to stderr ("direnv: loading .envrc", "direnv:
// error /path/.envrc is blocked", etc.); whatever it printed is returned in
// message so callers can surface it in their command log.
func Load(cmd oscommands.ICmdObjBuilder) (message string, err error) {
	if _, lookupErr := exec.LookPath("direnv"); lookupErr != nil {
		return "", nil
	}

	stdout, stderr, runErr := cmd.New([]string{
		"direnv", "export", "json",
	}).DontLog().RunWithOutputs()
	message = strings.TrimRight(stderr, "\n")

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
		return message, runErr
	}
	return message, parseErr
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
