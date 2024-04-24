package todo

type TodoCommand int

const (
	Pick TodoCommand = iota + 1
	Revert
	Edit
	Reword
	Fixup
	Squash

	Exec
	Break
	Label
	Reset
	Merge

	NoOp
	Drop
	UpdateRef

	Comment
)

type Todo struct {
	Command     TodoCommand
	Commit      string
	Flag        string
	Comment     string
	ExecCommand string
	Label       string
	Msg         string
	Ref         string
}

func (t TodoCommand) String() string {
	return commandToString[t]
}

var commandToString = map[TodoCommand]string{
	Pick:      "pick",
	Revert:    "revert",
	Edit:      "edit",
	Reword:    "reword",
	Fixup:     "fixup",
	Squash:    "squash",
	Exec:      "exec",
	Break:     "break",
	Label:     "label",
	Reset:     "reset",
	Merge:     "merge",
	NoOp:      "noop",
	Drop:      "drop",
	UpdateRef: "update-ref",
	Comment:   "comment",
}

var todoCommandInfo = [15]struct {
	nickname string
	cmd      string
}{
	{"", ""}, // dummy value since we're using 1-based indexing
	{"p", "pick"},
	{"", "revert"},
	{"e", "edit"},
	{"r", "reword"},
	{"f", "fixup"},
	{"s", "squash"},
	{"x", "exec"},
	{"b", "break"},
	{"l", "label"},
	{"t", "reset"},
	{"m", "merge"},
	{"", "noop"},
	{"d", "drop"},
	{"u", "update-ref"},
}
