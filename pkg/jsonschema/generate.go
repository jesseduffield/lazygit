//go:generate go run generator.go

package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/karimkhaleel/jsonschema"
	"github.com/samber/lo"
)

func GetSchemaDir() string {
	return utils.GetLazyRootDirectory() + "/schema-master"
}

func GenerateSchema() *jsonschema.Schema {
	schema := customReflect(&config.UserConfig{})
	obj, _ := json.MarshalIndent(schema, "", "  ")
	obj = append(obj, '\n')

	if err := os.WriteFile(GetSchemaDir()+"/config.json", obj, 0o644); err != nil {
		fmt.Println("Error writing to file:", err)
		return nil
	}
	return schema
}

func getSubSchema(rootSchema, parentSchema *jsonschema.Schema, key string) *jsonschema.Schema {
	subSchema, found := parentSchema.Properties.Get(key)
	if !found {
		panic(fmt.Sprintf("Failed to find subSchema at %s on parent", key))
	}

	// This means the schema is defined on the rootSchema's Definitions
	if subSchema.Ref != "" {
		key, _ = strings.CutPrefix(subSchema.Ref, "#/$defs/")
		refSchema, ok := rootSchema.Definitions[key]
		if !ok {
			panic(fmt.Sprintf("Failed to find #/$defs/%s", key))
		}
		refSchema.Description = subSchema.Description
		return refSchema
	}

	return subSchema
}

func customReflect(v *config.UserConfig) *jsonschema.Schema {
	r := &jsonschema.Reflector{FieldNameTag: "yaml", RequiredFromJSONSchemaTags: true}
	if err := r.AddGoComments("github.com/jesseduffield/lazygit/pkg/config", "../config"); err != nil {
		panic(err)
	}
	filterOutDevComments(r)
	schema := r.Reflect(v)
	defaultConfig := config.GetDefaultConfig()
	userConfigSchema := schema.Definitions["UserConfig"]

	defaultValue := reflect.ValueOf(defaultConfig).Elem()

	yamlToFieldNames := lo.Invert(userConfigSchema.OriginalPropertiesMapping)

	for pair := userConfigSchema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		yamlName := pair.Key
		fieldName := yamlToFieldNames[yamlName]

		subSchema := getSubSchema(schema, userConfigSchema, yamlName)

		setDefaultVals(schema, subSchema, defaultValue.FieldByName(fieldName).Interface())
	}

	return schema
}

func filterOutDevComments(r *jsonschema.Reflector) {
	for k, v := range r.CommentMap {
		commentLines := strings.Split(v, "\n")
		filteredCommentLines := lo.Filter(commentLines, func(line string, _ int) bool {
			return !strings.Contains(line, "[dev]")
		})
		r.CommentMap[k] = strings.Join(filteredCommentLines, "\n")
	}
}

func setDefaultVals(rootSchema, schema *jsonschema.Schema, defaults any) {
	t := reflect.TypeOf(defaults)
	v := reflect.ValueOf(defaults)

	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
		v = v.Elem()
	}

	k := t.Kind()
	_ = k

	switch t.Kind() {
	case reflect.Bool:
		schema.Default = v.Bool()
	case reflect.Int:
		schema.Default = v.Int()
	case reflect.String:
		schema.Default = v.String()
	default:
		// Do nothing
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := range t.NumField() {
		value := v.Field(i).Interface()
		parentKey := t.Field(i).Name

		key, ok := schema.OriginalPropertiesMapping[parentKey]
		if !ok {
			continue
		}

		subSchema := getSubSchema(rootSchema, schema, key)

		if isStruct(value) {
			setDefaultVals(rootSchema, subSchema, value)
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
		for i := range rv.NumField() {
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
