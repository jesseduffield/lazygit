package models

type IFileChange interface {
	GetHasUnstagedChanges() bool
	GetHasStagedChanges() bool
	GetIsTracked() bool
	GetPath() string
}
