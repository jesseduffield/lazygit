package types

type Ref interface {
	FullRefName() string
	RefName() string
	ParentRefName() string
	Description() string
}
