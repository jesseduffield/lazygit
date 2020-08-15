package commands

// CommitFile : A git commit file
type CommitFile struct {
	Sha           string
	Name          string
	DisplayString string
	Status        int // one of 'WHOLE' 'PART' 'NONE'
}
