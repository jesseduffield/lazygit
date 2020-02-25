package commands

// CommitFile : A git commit file
type CommitFile struct {
	Sha           string
	Name          string
	DisplayString string
	Status        int // one of 'WHOLE' 'PART' 'NONE'
}

const (
	// UNSELECTED is for when the commit file has not been added to the patch in any way
	UNSELECTED = iota
	// WHOLE is for when you want to add the whole diff of a file to the patch,
	// including e.g. if it was deleted
	WHOLE = iota
	// PART is for when you're only talking about specific lines that have been modified
	PART
)
