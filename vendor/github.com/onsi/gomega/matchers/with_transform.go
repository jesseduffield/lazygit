package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/internal/oraclematcher"
	"github.com/onsi/gomega/types"
)

type WithTransformMatcher struct {
	// input
	Transform interface{} // must be a function of one parameter that returns one value
	Matcher   types.GomegaMatcher

	// cached value
	transformArgType reflect.Type

	// state
	transformedValue interface{}
}

func NewWithTransformMatcher(transform interface{}, matcher types.GomegaMatcher) *WithTransformMatcher {
	if transform == nil {
		panic("transform function cannot be nil")
	}
	txType := reflect.TypeOf(transform)
	if txType.NumIn() != 1 {
		panic("transform function must have 1 argument")
	}
	if txType.NumOut() != 1 {
		panic("transform function must have 1 return value")
	}

	return &WithTransformMatcher{
		Transform:        transform,
		Matcher:          matcher,
		transformArgType: reflect.TypeOf(transform).In(0),
	}
}

func (m *WithTransformMatcher) Match(actual interface{}) (bool, error) {
	// prepare a parameter to pass to the Transform function
	var param reflect.Value
	if actual != nil && reflect.TypeOf(actual).AssignableTo(m.transformArgType) {
		// The dynamic type of actual is compatible with the transform argument.
		param = reflect.ValueOf(actual)

	} else if actual == nil && m.transformArgType.Kind() == reflect.Interface {
		// The dynamic type of actual is unknown, so there's no way to make its
		// reflect.Value. Create a nil of the transform argument, which is known.
		param = reflect.Zero(m.transformArgType)

	} else {
		return false, fmt.Errorf("Transform function expects '%s' but we have '%T'", m.transformArgType, actual)
	}

	// call the Transform function with `actual`
	fn := reflect.ValueOf(m.Transform)
	result := fn.Call([]reflect.Value{param})
	m.transformedValue = result[0].Interface() // expect exactly one value

	return m.Matcher.Match(m.transformedValue)
}

func (m *WithTransformMatcher) FailureMessage(_ interface{}) (message string) {
	return m.Matcher.FailureMessage(m.transformedValue)
}

func (m *WithTransformMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return m.Matcher.NegatedFailureMessage(m.transformedValue)
}

func (m *WithTransformMatcher) MatchMayChangeInTheFuture(_ interface{}) bool {
	// TODO: Maybe this should always just return true? (Only an issue for non-deterministic transformers.)
	//
	// Querying the next matcher is fine if the transformer always will return the same value.
	// But if the transformer is non-deterministic and returns a different value each time, then there
	// is no point in querying the next matcher, since it can only comment on the last transformed value.
	return oraclematcher.MatchMayChangeInTheFuture(m.Matcher, m.transformedValue)
}
