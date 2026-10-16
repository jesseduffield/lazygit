package shared

import "github.com/jesseduffield/lazygit/pkg/config"

// RecordsDiffRenderer stands in for a diff renderer that states each row it renders in
// an OSC 1717 record, as delta and diff-so-fancy do. It passes the diff through as it
// is, but since it states its rows, lazygit places them by the records alone and reads
// nothing out of the text, the diff's header included. It reads the diff's +/- column,
// so it wants its input uncoloured.
var RecordsDiffRenderer = config.DiffRendererConfig{
	Command: `printf '\033]1717;1\007'; awk '
		/^diff --git / { file = $0; sub(/^.* b\//, "", file); inHunk = 0 }
		/^@@ / {
			split($2, o, ","); split($3, n, ","); old = -o[1]; new = n[1] + 0; inHunk = 1
			printf "\033]1717;1;h;%d;;%s\007%s\n", new, file, $0; next
		}
		!inHunk { printf "\033]1717;1;f;;;%s\007%s\n", file, $0; next }
		/^\+/ { printf "\033]1717;1;a;%d;;%s\007%s\n", new, file, $0; new++; next }
		/^-/ { printf "\033]1717;1;d;%d;%d;%s\007%s\n", new, old, file, $0; old++; next }
		/^ / { printf "\033]1717;1;c;%d;;%s\007%s\n", new, file, $0; new++; old++; next }
		{ print }'`,
	ColorArg: "never",
}
