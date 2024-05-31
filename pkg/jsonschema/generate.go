//go:generate go run generator.go

package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/karimkhaleel/jsonschema"
)

func GetSchemaDir() string {
	return utils.GetLazyRootDirectory() + "/schema"
}

func GenerateSchema() {
	schema := customReflect(&config.UserConfig{})
	obj, _ := json.MarshalIndent(schema, "", "  ")

	if err := os.WriteFile(GetSchemaDir()+"/config.json", obj, 0o644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func customReflect(v *config.UserConfig) *jsonschema.Schema {
	defaultConfig := config.GetDefaultConfig()
	r := &jsonschema.Reflector{FieldNameTag: "yaml", RequiredFromJSONSchemaTags: true, DoNotReference: true}
	if err := r.AddGoComments("github.com/jesseduffield/lazygit/pkg/config", "../config"); err != nil {
		panic(err)
	}
	schema := r.Reflect(v)

	setDefaultVals(defaultConfig, schema)

	return schema
}

func setDefaultVals(defaults any, schema *jsonschema.Schema) {
	t := reflect.TypeOf(defaults)
	v := reflect.ValueOf(defaults)

	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i).Interface()
		parentKey := t.Field(i).Name

		key, ok := schema.OriginalPropertiesMapping[parentKey]
		if !ok {
			continue
		}

		subSchema, ok := schema.Properties.Get(key)
		if !ok {
			continue
		}

		if isStruct(value) {
			setDefaultVals(value, subSchema)
		} else if !isZeroValue(value) {
			subSchema.Default = value
		}
	}
}

func isZeroValue(v any) bool {
	switch v := v.(type) {
	case int, int32, int64, float32, float64:
		return v == 0
	case string:
		return v == ""
	case bool:
		return false
	case nil:
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Map:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			if !isZeroValue(rv.Field(i).Interface()) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func isStruct(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}
