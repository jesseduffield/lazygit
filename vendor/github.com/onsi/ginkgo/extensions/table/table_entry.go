package table

import (
	"fmt"
	"reflect"

	"github.com/onsi/ginkgo/internal/codelocation"
	"github.com/onsi/ginkgo/internal/global"
	"github.com/onsi/ginkgo/types"
)

/*
TableEntry represents an entry in a table test.  You generally use the `Entry` constructor.
*/
type TableEntry struct {
	Description  interface{}
	Parameters   []interface{}
	Pending      bool
	Focused      bool
	codeLocation types.CodeLocation
}

func (t TableEntry) generateIt(itBody reflect.Value) {
	var description string
	descriptionValue := reflect.ValueOf(t.Description)
	switch descriptionValue.Kind() {
	case reflect.String:
		description = descriptionValue.String()
	case reflect.Func:
		values := castParameters(descriptionValue, t.Parameters)
		res := descriptionValue.Call(values)
		if len(res) != 1 {
			panic(fmt.Sprintf("The describe function should return only a value, returned %d", len(res)))
		}
		if res[0].Kind() != reflect.String {
			panic(fmt.Sprintf("The describe function should return a string, returned %#v", res[0]))
		}
		description = res[0].String()
	default:
		panic(fmt.Sprintf("Description can either be a string or a function, got %#v", descriptionValue))
	}

	if t.Pending {
		global.Suite.PushItNode(description, func() {}, types.FlagTypePending, t.codeLocation, 0)
		return
	}

	values := castParameters(itBody, t.Parameters)
	body := func() {
		itBody.Call(values)
	}

	if t.Focused {
		global.Suite.PushItNode(description, body, types.FlagTypeFocused, t.codeLocation, global.DefaultTimeout)
	} else {
		global.Suite.PushItNode(description, body, types.FlagTypeNone, t.codeLocation, global.DefaultTimeout)
	}
}

func castParameters(function reflect.Value, parameters []interface{}) []reflect.Value {
	res := make([]reflect.Value, len(parameters))
	funcType := function.Type()
	for i, param := range parameters {
		if param == nil {
			inType := funcType.In(i)
			res[i] = reflect.Zero(inType)
		} else {
			res[i] = reflect.ValueOf(param)
		}
	}
	return res
}

/*
Entry constructs a TableEntry.

The first argument is a required description (this becomes the content of the generated Ginkgo `It`).
Subsequent parameters are saved off and sent to the callback passed in to `DescribeTable`.

Each Entry ends up generating an individual Ginkgo It.
*/
func Entry(description interface{}, parameters ...interface{}) TableEntry {
	return TableEntry{
		Description:  description,
		Parameters:   parameters,
		Pending:      false,
		Focused:      false,
		codeLocation: codelocation.New(1),
	}
}

/*
You can focus a particular entry with FEntry.  This is equivalent to FIt.
*/
func FEntry(description interface{}, parameters ...interface{}) TableEntry {
	return TableEntry{
		Description:  description,
		Parameters:   parameters,
		Pending:      false,
		Focused:      true,
		codeLocation: codelocation.New(1),
	}
}

/*
You can mark a particular entry as pending with PEntry.  This is equivalent to PIt.
*/
func PEntry(description interface{}, parameters ...interface{}) TableEntry {
	return TableEntry{
		Description:  description,
		Parameters:   parameters,
		Pending:      true,
		Focused:      false,
		codeLocation: codelocation.New(1),
	}
}

/*
You can mark a particular entry as pending with XEntry.  This is equivalent to XIt.
*/
func XEntry(description interface{}, parameters ...interface{}) TableEntry {
	return TableEntry{
		Description:  description,
		Parameters:   parameters,
		Pending:      true,
		Focused:      false,
		codeLocation: codelocation.New(1),
	}
}
