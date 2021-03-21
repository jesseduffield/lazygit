package models

// sometimes we need to deal with either a node (which contains a file) or an actual file
type IFileChange interface {
	GetHasUnstagedChanges() bool
	GetHasStagedChanges() bool
	GetIsTracked() bool
	GetPath() string
}
