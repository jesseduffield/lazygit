package types

type Ref interface {
	RefName() string
	ParentRefName() string
	Description() string
}
