package generator

import (
	"errors"
	"go/types"
)

func (f *Fake) loadMethodForFunction() error {
	t, ok := f.Target.Type().(*types.Named)
	if !ok {
		return errors.New("target is not a named type")
	}
	sig, ok := t.Underlying().(*types.Signature)
	if !ok {
		return errors.New("target does not have an underlying function signature")
	}
	f.addTypesForMethod(sig)
	f.Function = methodForSignature(sig, f.TargetName, f.Imports)
	return nil
}
