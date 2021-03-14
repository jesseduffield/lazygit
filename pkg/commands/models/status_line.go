package models

type IStatusLine interface {
	GetHasUnstagedChanges() bool
	GetHasStagedChanges() bool
	GetIsTracked() bool
	GetPath() string
}
