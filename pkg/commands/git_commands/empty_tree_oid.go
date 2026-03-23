package git_commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// GetEmptyTreeOID returns the empty tree object ID for the repository's object format
// (SHA-1 vs SHA-256). The SHA-1 empty tree constant is wrong for sha256 repositories.
func GetEmptyTreeOID(cmd oscommands.ICmdObjBuilder) (string, error) {
	stdout, _, err := cmd.New(
		NewGitCmd("hash-object").Arg("-t", "tree", "--stdin").ToArgv(),
	).SetStdin("").DontLog().RunWithOutputs()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout), nil
}
