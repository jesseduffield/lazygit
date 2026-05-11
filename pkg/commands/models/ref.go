package models

type Ref interface {
	FullRefName() string
	RefName() string
	ShortRefName() string
	ParentRefName() string
	Description() string
}
