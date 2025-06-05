package types

type Ref interface {
	FullRefName() string
	RefName() string
	ShortRefName() string
	ParentRefName() string
	Description() string
}

type RefRange struct {
	From Ref
	To   Ref
}
