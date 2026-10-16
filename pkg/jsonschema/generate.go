//go:generate go run generator.go

package jsonschema

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/karimkhaleel/jsonschema"
	"github.com/samber/lo"
)

func GetSchemaDir() string {
	return utils.MustFindLazygitRootDirectory() + "/schema-master"
}

// GenerateSchema writes the schema of the user config, and returns it together
// with the properties that repeat the fields of a struct that the schema has at
// an earlier path.
func GenerateSchema() (*jsonschema.Schema, []*jsonschema.Schema) {
	schema, repeatedStructs := customReflect(&config.UserConfig{})
	obj, _ := json.MarshalIndent(schema, "", "  ")
	obj = append(obj, '\n')

	if err := os.WriteFile(GetSchemaDir()+"/config.json", obj, 0o644); err != nil {
		fmt.Println("Error writing to file:", err)
		return nil, nil
	}
	return schema, repeatedStructs
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

func customReflect(v *config.UserConfig) (*jsonschema.Schema, []*jsonschema.Schema) {
	r := &jsonschema.Reflector{FieldNameTag: "yaml", RequiredFromJSONSchemaTags: true}
	if err := r.AddGoComments("github.com/jesseduffield/lazygit/pkg/config", "../config"); err != nil {
		panic(err)
	}
	filterOutDevComments(r)
	schema := r.Reflect(v)
	inlineKeybindingRefs(schema)
	repeatedStructs := inlineSharedStructDefinitions(schema, "UserConfig")
	defaultConfig := config.GetDefaultConfig()
	defaultConfig.Keybinding.MergeLegacyAltKeybindings()
	userConfigSchema := schema.Definitions["UserConfig"]

	defaultValue := reflect.ValueOf(defaultConfig).Elem()

	yamlToFieldNames := lo.Invert(userConfigSchema.OriginalPropertiesMapping)

	for pair := userConfigSchema.Properties.Oldest(); pair != nil; pair = pair.Next() {
		yamlName := pair.Key
		fieldName := yamlToFieldNames[yamlName]

		subSchema := getSubSchema(schema, userConfigSchema, yamlName)

		setDefaultVals(schema, subSchema, defaultValue.FieldByName(fieldName).Interface())
	}

	return schema, repeatedStructs
}

// inlineSharedStructDefinitions gives each property whose struct definition is
// shared with other properties a copy of that definition of its own. The
// defaults of a struct's fields depend on its path, and setDefaultVals writes
// them into the definition; properties that share a definition would get each
// other's defaults.
//
// It returns the copies at all but the first path of each struct, in the order
// in which a depth-first walk from the root definition meets them.
func inlineSharedStructDefinitions(schema *jsonschema.Schema, rootDefinition string) []*jsonschema.Schema {
	var repeatedStructs []*jsonschema.Schema
	inlinedDefinitions := map[string]bool{}

	// A property inside a shared definition is met once for each path of that
	// definition. It only gets copies of its own once its parent has them, so
	// this repeats until no definition is shared.
	for {
		var properties []*jsonschema.Schema
		var collect func(s *jsonschema.Schema)
		collect = func(s *jsonschema.Schema) {
			if s.Properties == nil {
				return
			}
			for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
				property := pair.Value
				if def := structDefinition(schema, property); def != nil {
					properties = append(properties, property)
					collect(def)
				} else {
					collect(property)
				}
			}
		}
		collect(schema.Definitions[rootDefinition])
		properties = lo.Uniq(properties)

		refCounts := lo.CountValuesBy(properties, func(s *jsonschema.Schema) string { return s.Ref })
		shared := lo.Filter(properties, func(s *jsonschema.Schema, _ int) bool { return refCounts[s.Ref] > 1 })
		if len(shared) == 0 {
			break
		}

		seenRefs := map[string]bool{}
		for _, property := range shared {
			ref := property.Ref
			description := property.Description
			*property = *copySchema(structDefinition(schema, property))
			property.Description = description

			if seenRefs[ref] {
				repeatedStructs = append(repeatedStructs, property)
			}
			seenRefs[ref] = true
			inlinedDefinitions[strings.TrimPrefix(ref, "#/$defs/")] = true
		}
	}

	for name := range inlinedDefinitions {
		if !isReferenced(schema, "#/$defs/"+name) {
			delete(schema.Definitions, name)
		}
	}

	return repeatedStructs
}

func isReferenced(schema *jsonschema.Schema, ref string) bool {
	var refers func(s *jsonschema.Schema) bool
	refers = func(s *jsonschema.Schema) bool {
		if s == nil {
			return false
		}
		if s.Ref == ref || refers(s.Items) || refers(s.AdditionalProperties) {
			return true
		}
		if s.Properties != nil {
			for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
				if refers(pair.Value) {
					return true
				}
			}
		}
		return false
	}
	return lo.SomeBy(lo.Values(schema.Definitions), refers)
}

// structDefinition returns the definition that s refers to if it is the
// definition of a struct, and nil otherwise.
func structDefinition(schema, s *jsonschema.Schema) *jsonschema.Schema {
	if s.Ref == "" {
		return nil
	}
	def := schema.Definitions[strings.TrimPrefix(s.Ref, "#/$defs/")]
	if def == nil || def.Properties == nil {
		return nil
	}
	return def
}

// copySchema copies s and the schemas that it has for its properties, items and
// additional properties.
func copySchema(s *jsonschema.Schema) *jsonschema.Schema {
	c := *s
	if s.Properties != nil {
		c.Properties = jsonschema.NewProperties()
		for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
			c.Properties.Set(pair.Key, copySchema(pair.Value))
		}
	}
	if s.Items != nil {
		c.Items = copySchema(s.Items)
	}
	if s.AdditionalProperties != nil {
		c.AdditionalProperties = copySchema(s.AdditionalProperties)
	}
	return &c
}

// inlineKeybindingRefs replaces every `$ref: #/$defs/Keybinding` in the
// schema with the inlined oneOf union, then drops the Keybinding definition.
//
// The schema generator stores types that implement JSONSchema() as shared
// definitions and uses $ref to point at them. That works for most types
// (where every reference logically points at the same data), but for
// Keybinding fields each property carries its own description and default,
// and writing those onto the shared definition would clobber siblings.
// Inlining sidesteps the issue.
func inlineKeybindingRefs(schema *jsonschema.Schema) {
	const ref = "#/$defs/Keybinding"
	keybindingDef, ok := schema.Definitions["Keybinding"]
	if !ok {
		return
	}
	inline := func(s *jsonschema.Schema) {
		desc := s.Description
		*s = *keybindingDef
		s.Description = desc
	}
	var visit func(s *jsonschema.Schema)
	visit = func(s *jsonschema.Schema) {
		if s == nil {
			return
		}
		if s.Properties != nil {
			for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
				if pair.Value.Ref == ref {
					inline(pair.Value)
				} else {
					visit(pair.Value)
				}
			}
		}
		if s.Items != nil {
			if s.Items.Ref == ref {
				inline(s.Items)
			} else {
				visit(s.Items)
			}
		}
		if s.AdditionalProperties != nil {
			visit(s.AdditionalProperties)
		}
	}
	for _, def := range schema.Definitions {
		visit(def)
	}
	delete(schema.Definitions, "Keybinding")
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

	if t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface {
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
	case reflect.Pointer, reflect.Interface:
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
