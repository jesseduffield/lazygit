package commands

// Conflict : A git conflict with a start middle and end corresponding to line
// numbers in the file where the conflict bars appear
type Conflict struct {
	Start  int
	Middle int
	End    int
}
