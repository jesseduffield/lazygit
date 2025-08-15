// Package jsonschema uses reflection to generate JSON Schemas from Go types [1].
//
// If json tags are present on struct fields, they will be used to infer
// property names and if a property is required (omitempty is present).
//
// [1] http://json-schema.org/latest/json-schema-validation.html
package jsonschema

import (
	"bytes"
	"encoding/json"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// Version is the JSON Schema version.
var Version = "https://json-schema.org/draft/2020-12/schema"

// Schema represents a JSON Schema object type.
// RFC draft-bhutton-json-schema-00 section 4.3
type Schema struct {
	// RFC draft-bhutton-json-schema-00
	Version     string      `json:"$schema,omitempty"`     // section 8.1.1
	ID          ID          `json:"$id,omitempty"`         // section 8.2.1
	Anchor      string      `json:"$anchor,omitempty"`     // section 8.2.2
	Ref         string      `json:"$ref,omitempty"`        // section 8.2.3.1
	DynamicRef  string      `json:"$dynamicRef,omitempty"` // section 8.2.3.2
	Definitions Definitions `json:"$defs,omitempty"`       // section 8.2.4
	Comments    string      `json:"$comment,omitempty"`    // section 8.3
	// RFC draft-bhutton-json-schema-00 section 10.2.1 (Sub-schemas with logic)
	AllOf []*Schema `json:"allOf,omitempty"` // section 10.2.1.1
	AnyOf []*Schema `json:"anyOf,omitempty"` // section 10.2.1.2
	OneOf []*Schema `json:"oneOf,omitempty"` // section 10.2.1.3
	Not   *Schema   `json:"not,omitempty"`   // section 10.2.1.4
	// RFC draft-bhutton-json-schema-00 section 10.2.2 (Apply sub-schemas conditionally)
	If               *Schema            `json:"if,omitempty"`               // section 10.2.2.1
	Then             *Schema            `json:"then,omitempty"`             // section 10.2.2.2
	Else             *Schema            `json:"else,omitempty"`             // section 10.2.2.3
	DependentSchemas map[string]*Schema `json:"dependentSchemas,omitempty"` // section 10.2.2.4
	// RFC draft-bhutton-json-schema-00 section 10.3.1 (arrays)
	PrefixItems []*Schema `json:"prefixItems,omitempty"` // section 10.3.1.1
	Items       *Schema   `json:"items,omitempty"`       // section 10.3.1.2  (replaces additionalItems)
	Contains    *Schema   `json:"contains,omitempty"`    // section 10.3.1.3
	// RFC draft-bhutton-json-schema-00 section 10.3.2 (sub-schemas)
	Properties                *orderedmap.OrderedMap[string, *Schema] `json:"properties,omitempty"` // section 10.3.2.1
	OriginalPropertiesMapping map[string]string                       `json:"-"`
	PatternProperties         map[string]*Schema                      `json:"patternProperties,omitempty"`    // section 10.3.2.2
	AdditionalProperties      *Schema                                 `json:"additionalProperties,omitempty"` // section 10.3.2.3
	PropertyNames             *Schema                                 `json:"propertyNames,omitempty"`        // section 10.3.2.4
	// RFC draft-bhutton-json-schema-validation-00, section 6
	Type              string              `json:"type,omitempty"`              // section 6.1.1
	Enum              []any               `json:"enum,omitempty"`              // section 6.1.2
	Const             any                 `json:"const,omitempty"`             // section 6.1.3
	MultipleOf        json.Number         `json:"multipleOf,omitempty"`        // section 6.2.1
	Maximum           json.Number         `json:"maximum,omitempty"`           // section 6.2.2
	ExclusiveMaximum  json.Number         `json:"exclusiveMaximum,omitempty"`  // section 6.2.3
	Minimum           json.Number         `json:"minimum,omitempty"`           // section 6.2.4
	ExclusiveMinimum  json.Number         `json:"exclusiveMinimum,omitempty"`  // section 6.2.5
	MaxLength         int                 `json:"maxLength,omitempty"`         // section 6.3.1
	MinLength         int                 `json:"minLength,omitempty"`         // section 6.3.2
	Pattern           string              `json:"pattern,omitempty"`           // section 6.3.3
	MaxItems          int                 `json:"maxItems,omitempty"`          // section 6.4.1
	MinItems          int                 `json:"minItems,omitempty"`          // section 6.4.2
	UniqueItems       bool                `json:"uniqueItems,omitempty"`       // section 6.4.3
	MaxContains       uint                `json:"maxContains,omitempty"`       // section 6.4.4
	MinContains       uint                `json:"minContains,omitempty"`       // section 6.4.5
	MaxProperties     int                 `json:"maxProperties,omitempty"`     // section 6.5.1
	MinProperties     int                 `json:"minProperties,omitempty"`     // section 6.5.2
	Required          []string            `json:"required,omitempty"`          // section 6.5.3
	DependentRequired map[string][]string `json:"dependentRequired,omitempty"` // section 6.5.4
	// RFC draft-bhutton-json-schema-validation-00, section 7
	Format string `json:"format,omitempty"`
	// RFC draft-bhutton-json-schema-validation-00, section 8
	ContentEncoding  string  `json:"contentEncoding,omitempty"`  // section 8.3
	ContentMediaType string  `json:"contentMediaType,omitempty"` // section 8.4
	ContentSchema    *Schema `json:"contentSchema,omitempty"`    // section 8.5
	// RFC draft-bhutton-json-schema-validation-00, section 9
	Title       string `json:"title,omitempty"`       // section 9.1
	Description string `json:"description,omitempty"` // section 9.1
	Default     any    `json:"default,omitempty"`     // section 9.2
	Deprecated  bool   `json:"deprecated,omitempty"`  // section 9.3
	ReadOnly    bool   `json:"readOnly,omitempty"`    // section 9.4
	WriteOnly   bool   `json:"writeOnly,omitempty"`   // section 9.4
	Examples    []any  `json:"examples,omitempty"`    // section 9.5

	Extras map[string]any `json:"-"`

	// Special boolean representation of the Schema - section 4.3.2
	boolean *bool
}

var (
	// TrueSchema defines a schema with a true value
	TrueSchema = &Schema{boolean: &[]bool{true}[0]}
	// FalseSchema defines a schema with a false value
	FalseSchema = &Schema{boolean: &[]bool{false}[0]}
)

// customSchemaImpl is used to detect if the type provides it's own
// custom Schema Type definition to use instead. Very useful for situations
// where there are custom JSON Marshal and Unmarshal methods.
type customSchemaImpl interface {
	JSONSchema() *Schema
}

// Function to be run after the schema has been generated.
// this will let you modify a schema afterwards
type extendSchemaImpl interface {
	JSONSchemaExtend(*Schema)
}

var customType = reflect.TypeOf((*customSchemaImpl)(nil)).Elem()
var extendType = reflect.TypeOf((*extendSchemaImpl)(nil)).Elem()

// customSchemaGetFieldDocString
type customSchemaGetFieldDocString interface {
	GetFieldDocString(fieldName string) string
}

type customGetFieldDocString func(fieldName string) string

var customStructGetFieldDocString = reflect.TypeOf((*customSchemaGetFieldDocString)(nil)).Elem()

// Reflect reflects to Schema from a value using the default Reflector
func Reflect(v any) *Schema {
	return ReflectFromType(reflect.TypeOf(v))
}

// ReflectFromType generates root schema using the default Reflector
func ReflectFromType(t reflect.Type) *Schema {
	r := &Reflector{}
	return r.ReflectFromType(t)
}

// A Reflector reflects values into a Schema.
type Reflector struct {
	// BaseSchemaID defines the URI that will be used as a base to determine Schema
	// IDs for models. For example, a base Schema ID of `https://invopop.com/schemas`
	// when defined with a struct called `User{}`, will result in a schema with an
	// ID set to `https://invopop.com/schemas/user`.
	//
	// If no `BaseSchemaID` is provided, we'll take the type's complete package path
	// and use that as a base instead. Set `Anonymous` to try if you do not want to
	// include a schema ID.
	BaseSchemaID ID

	// Anonymous when true will hide the auto-generated Schema ID and provide what is
	// known as an "anonymous schema". As a rule, this is not recommended.
	Anonymous bool

	// AssignAnchor when true will use the original struct's name as an anchor inside
	// every definition, including the root schema. These can be useful for having a
	// reference to the original struct's name in CamelCase instead of the snake-case used
	// by default for URI compatibility.
	//
	// Anchors do not appear to be widely used out in the wild, so at this time the
	// anchors themselves will not be used inside generated schema.
	AssignAnchor bool

	// AllowAdditionalProperties will cause the Reflector to generate a schema
	// without additionalProperties set to 'false' for all struct types. This means
	// the presence of additional keys in JSON objects will not cause validation
	// to fail. Note said additional keys will simply be dropped when the
	// validated JSON is unmarshaled.
	AllowAdditionalProperties bool

	// RequiredFromJSONSchemaTags will cause the Reflector to generate a schema
	// that requires any key tagged with `jsonschema:required`, overriding the
	// default of requiring any key *not* tagged with `json:,omitempty`.
	RequiredFromJSONSchemaTags bool

	// Do not reference definitions. This will remove the top-level $defs map and
	// instead cause the entire structure of types to be output in one tree. The
	// list of type definitions (`$defs`) will not be included.
	DoNotReference bool

	// ExpandedStruct when true will include the reflected type's definition in the
	// root as opposed to a definition with a reference.
	ExpandedStruct bool

	// FieldNameTag will change the tag used to get field names. json tags are used by default.
	FieldNameTag string

	// IgnoredTypes defines a slice of types that should be ignored in the schema,
	// switching to just allowing additional properties instead.
	IgnoredTypes []any

	// Lookup allows a function to be defined that will provide a custom mapping of
	// types to Schema IDs. This allows existing schema documents to be referenced
	// by their ID instead of being embedded into the current schema definitions.
	// Reflected types will never be pointers, only underlying elements.
	Lookup func(reflect.Type) ID

	// Mapper is a function that can be used to map custom Go types to jsonschema schemas.
	Mapper func(reflect.Type) *Schema

	// Namer allows customizing of type names. The default is to use the type's name
	// provided by the reflect package.
	Namer func(reflect.Type) string

	// KeyNamer allows customizing of key names.
	// The default is to use the key's name as is, or the json tag if present.
	// If a json tag is present, KeyNamer will receive the tag's name as an argument, not the original key name.
	KeyNamer func(string) string

	// AdditionalFields allows adding structfields for a given type
	AdditionalFields func(reflect.Type) []reflect.StructField

	// CommentMap is a dictionary of fully qualified go types and fields to comment
	// strings that will be used if a description has not already been provided in
	// the tags. Types and fields are added to the package path using "." as a
	// separator.
	//
	// Type descriptions should be defined like:
	//
	//   map[string]string{"github.com/invopop/jsonschema.Reflector": "A Reflector reflects values into a Schema."}
	//
	// And Fields defined as:
	//
	//   map[string]string{"github.com/invopop/jsonschema.Reflector.DoNotReference": "Do not reference definitions."}
	//
	// See also: AddGoComments
	CommentMap map[string]string
}

// Reflect reflects to Schema from a value.
func (r *Reflector) Reflect(v any) *Schema {
	return r.ReflectFromType(reflect.TypeOf(v))
}

// ReflectFromType generates root schema
func (r *Reflector) ReflectFromType(t reflect.Type) *Schema {
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // re-assign from pointer
	}

	name := r.typeName(t)

	s := new(Schema)
	definitions := Definitions{}
	s.Definitions = definitions
	bs := r.reflectTypeToSchemaWithID(definitions, t)
	if r.ExpandedStruct {
		*s = *definitions[name]
		delete(definitions, name)
	} else {
		*s = *bs
	}

	// Attempt to set the schema ID
	if !r.Anonymous && s.ID == EmptyID {
		baseSchemaID := r.BaseSchemaID
		if baseSchemaID == EmptyID {
			id := ID("https://" + t.PkgPath())
			if err := id.Validate(); err == nil {
				// it's okay to silently ignore URL errors
				baseSchemaID = id
			}
		}
		if baseSchemaID != EmptyID {
			s.ID = baseSchemaID.Add(ToSnakeCase(name))
		}
	}

	s.Version = Version
	if !r.DoNotReference {
		s.Definitions = definitions
	}

	return s
}

// Definitions hold schema definitions.
// http://json-schema.org/latest/json-schema-validation.html#rfc.section.5.26
// RFC draft-wright-json-schema-validation-00, section 5.26
type Definitions map[string]*Schema

// Available Go defined types for JSON Schema Validation.
// RFC draft-wright-json-schema-validation-00, section 7.3
var (
	timeType = reflect.TypeOf(time.Time{}) // date-time RFC section 7.3.1
	ipType   = reflect.TypeOf(net.IP{})    // ipv4 and ipv6 RFC section 7.3.4, 7.3.5
	uriType  = reflect.TypeOf(url.URL{})   // uri RFC section 7.3.6
)

// Byte slices will be encoded as base64
var byteSliceType = reflect.TypeOf([]byte(nil))

// Except for json.RawMessage
var rawMessageType = reflect.TypeOf(json.RawMessage{})

// Go code generated from protobuf enum types should fulfil this interface.
type protoEnum interface {
	EnumDescriptor() ([]byte, []int)
}

var protoEnumType = reflect.TypeOf((*protoEnum)(nil)).Elem()

// SetBaseSchemaID is a helper use to be able to set the reflectors base
// schema ID from a string as opposed to then ID instance.
func (r *Reflector) SetBaseSchemaID(id string) {
	r.BaseSchemaID = ID(id)
}

func (r *Reflector) refOrReflectTypeToSchema(definitions Definitions, t reflect.Type) *Schema {
	id := r.lookupID(t)
	if id != EmptyID {
		return &Schema{
			Ref: id.String(),
		}
	}

	// Already added to definitions?
	if def := r.refDefinition(definitions, t); def != nil {
		return def
	}

	return r.reflectTypeToSchemaWithID(definitions, t)
}

func (r *Reflector) reflectTypeToSchemaWithID(defs Definitions, t reflect.Type) *Schema {
	s := r.reflectTypeToSchema(defs, t)
	if s != nil {
		if r.Lookup != nil {
			id := r.Lookup(t)
			if id != EmptyID {
				s.ID = id
			}
		}
	}
	return s
}

func (r *Reflector) reflectTypeToSchema(definitions Definitions, t reflect.Type) *Schema {
	// only try to reflect non-pointers
	if t.Kind() == reflect.Ptr {
		return r.refOrReflectTypeToSchema(definitions, t.Elem())
	}

	// Do any pre-definitions exist?
	if r.Mapper != nil {
		if t := r.Mapper(t); t != nil {
			return t
		}
	}
	if rt := r.reflectCustomSchema(definitions, t); rt != nil {
		return rt
	}

	// Prepare a base to which details can be added
	st := new(Schema)

	// jsonpb will marshal protobuf enum options as either strings or integers.
	// It will unmarshal either.
	if t.Implements(protoEnumType) {
		st.OneOf = []*Schema{
			{Type: "string"},
			{Type: "integer"},
		}
		return st
	}

	// Defined format types for JSON Schema Validation
	// RFC draft-wright-json-schema-validation-00, section 7.3
	// TODO email RFC section 7.3.2, hostname RFC section 7.3.3, uriref RFC section 7.3.7
	if t == ipType {
		// TODO differentiate ipv4 and ipv6 RFC section 7.3.4, 7.3.5
		st.Type = "string"
		st.Format = "ipv4"
		return st
	}

	switch t.Kind() {
	case reflect.Struct:
		r.reflectStruct(definitions, t, st)

	case reflect.Slice, reflect.Array:
		r.reflectSliceOrArray(definitions, t, st)

	case reflect.Map:
		r.reflectMap(definitions, t, st)

	case reflect.Interface:
		// empty

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		st.Type = "integer"

	case reflect.Float32, reflect.Float64:
		st.Type = "number"

	case reflect.Bool:
		st.Type = "boolean"

	case reflect.String:
		st.Type = "string"

	default:
		panic("unsupported type " + t.String())
	}

	r.reflectSchemaExtend(definitions, t, st)

	// Always try to reference the definition which may have just been created
	if def := r.refDefinition(definitions, t); def != nil {
		return def
	}

	return st
}

func (r *Reflector) reflectCustomSchema(definitions Definitions, t reflect.Type) *Schema {
	if t.Kind() == reflect.Ptr {
		return r.reflectCustomSchema(definitions, t.Elem())
	}

	if t.Implements(customType) {
		v := reflect.New(t)
		o := v.Interface().(customSchemaImpl)
		st := o.JSONSchema()
		r.addDefinition(definitions, t, st)
		if ref := r.refDefinition(definitions, t); ref != nil {
			return ref
		}
		return st
	}

	return nil
}

func (r *Reflector) reflectSchemaExtend(definitions Definitions, t reflect.Type, s *Schema) *Schema {
	if t.Implements(extendType) {
		v := reflect.New(t)
		o := v.Interface().(extendSchemaImpl)
		o.JSONSchemaExtend(s)
		if ref := r.refDefinition(definitions, t); ref != nil {
			return ref
		}
	}

	return s
}

func (r *Reflector) reflectSliceOrArray(definitions Definitions, t reflect.Type, st *Schema) {
	if t == rawMessageType {
		return
	}

	r.addDefinition(definitions, t, st)

	if st.Description == "" {
		st.Description = r.lookupComment(t, "")
	}

	if t.Kind() == reflect.Array {
		st.MinItems = t.Len()
		st.MaxItems = st.MinItems
	}
	if t.Kind() == reflect.Slice && t.Elem() == byteSliceType.Elem() {
		st.Type = "string"
		// NOTE: ContentMediaType is not set here
		st.ContentEncoding = "base64"
	} else {
		st.Type = "array"
		st.Items = r.refOrReflectTypeToSchema(definitions, t.Elem())
	}
}

func (r *Reflector) reflectMap(definitions Definitions, t reflect.Type, st *Schema) {
	r.addDefinition(definitions, t, st)

	st.Type = "object"
	if st.Description == "" {
		st.Description = r.lookupComment(t, "")
	}

	switch t.Key().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		st.PatternProperties = map[string]*Schema{
			"^[0-9]+$": r.refOrReflectTypeToSchema(definitions, t.Elem()),
		}
		st.AdditionalProperties = FalseSchema
		return
	}
	if t.Elem().Kind() != reflect.Interface {
		st.AdditionalProperties = r.refOrReflectTypeToSchema(definitions, t.Elem())
	}
}

// Reflects a struct to a JSON Schema type.
func (r *Reflector) reflectStruct(definitions Definitions, t reflect.Type, s *Schema) {
	// Handle special types
	switch t {
	case timeType: // date-time RFC section 7.3.1
		s.Type = "string"
		s.Format = "date-time"
		return
	case uriType: // uri RFC section 7.3.6
		s.Type = "string"
		s.Format = "uri"
		return
	}

	r.addDefinition(definitions, t, s)
	s.Type = "object"
	s.Properties = NewProperties()
	s.OriginalPropertiesMapping = make(map[string]string)
	s.Description = r.lookupComment(t, "")
	if r.AssignAnchor {
		s.Anchor = t.Name()
	}
	if !r.AllowAdditionalProperties && s.AdditionalProperties == nil {
		s.AdditionalProperties = FalseSchema
	}

	ignored := false
	for _, it := range r.IgnoredTypes {
		if reflect.TypeOf(it) == t {
			ignored = true
			break
		}
	}
	if !ignored {
		r.reflectStructFields(s, definitions, t)
	}
}

func (r *Reflector) reflectStructFields(st *Schema, definitions Definitions, t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	var getFieldDocString customGetFieldDocString
	if t.Implements(customStructGetFieldDocString) {
		v := reflect.New(t)
		o := v.Interface().(customSchemaGetFieldDocString)
		getFieldDocString = o.GetFieldDocString
	}

	handleField := func(f reflect.StructField) {
		name, originalName, shouldEmbed, required, nullable := r.reflectFieldName(f)
		// if anonymous and exported type should be processed recursively
		// current type should inherit properties of anonymous one
		if name == "" {
			if shouldEmbed {
				r.reflectStructFields(st, definitions, f.Type)
			}
			return
		}

		property := r.refOrReflectTypeToSchema(definitions, f.Type)
		property.structKeywordsFromTags(f, st, name)
		if property.Description == "" {
			property.Description = r.lookupComment(t, f.Name)
		}
		if getFieldDocString != nil {
			property.Description = getFieldDocString(f.Name)
		}

		if nullable {
			property = &Schema{
				OneOf: []*Schema{
					property,
					{
						Type: "null",
					},
				},
			}
		}

		st.Properties.Set(name, property)
		st.OriginalPropertiesMapping[originalName] = name
		if required {
			st.Required = appendUniqueString(st.Required, name)
		}
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		handleField(f)
	}
	if r.AdditionalFields != nil {
		if af := r.AdditionalFields(t); af != nil {
			for _, sf := range af {
				handleField(sf)
			}
		}
	}
}

func appendUniqueString(base []string, value string) []string {
	for _, v := range base {
		if v == value {
			return base
		}
	}
	return append(base, value)
}

func (r *Reflector) lookupComment(t reflect.Type, name string) string {
	if r.CommentMap == nil {
		return ""
	}

	n := fullyQualifiedTypeName(t)
	if name != "" {
		n = n + "." + name
	}

	return r.CommentMap[n]
}

// addDefinition will append the provided schema. If needed, an ID and anchor will also be added.
func (r *Reflector) addDefinition(definitions Definitions, t reflect.Type, s *Schema) {
	name := r.typeName(t)
	if name == "" {
		return
	}
	definitions[name] = s
}

// refDefinition will provide a schema with a reference to an existing definition.
func (r *Reflector) refDefinition(definitions Definitions, t reflect.Type) *Schema {
	if r.DoNotReference {
		return nil
	}
	name := r.typeName(t)
	if name == "" {
		return nil
	}
	if _, ok := definitions[name]; !ok {
		return nil
	}
	return &Schema{
		Ref: "#/$defs/" + name,
	}
}

func (r *Reflector) lookupID(t reflect.Type) ID {
	if r.Lookup != nil {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		return r.Lookup(t)

	}
	return EmptyID
}

func (t *Schema) structKeywordsFromTags(f reflect.StructField, parent *Schema, propertyName string) {
	t.Description = f.Tag.Get("jsonschema_description")

	tags := splitOnUnescapedCommas(f.Tag.Get("jsonschema"))
	tags = t.genericKeywords(tags, parent, propertyName)

	switch t.Type {
	case "string":
		t.stringKeywords(tags)
	case "number":
		t.numericalKeywords(tags)
	case "integer":
		t.numericalKeywords(tags)
	case "array":
		t.arrayKeywords(tags)
	case "boolean":
		t.booleanKeywords(tags)
	}
	extras := strings.Split(f.Tag.Get("jsonschema_extras"), ",")
	t.extraKeywords(extras)
}

// read struct tags for generic keywords
func (t *Schema) genericKeywords(tags []string, parent *Schema, propertyName string) []string { //nolint:gocyclo
	unprocessed := make([]string, 0, len(tags))
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "title":
				t.Title = val
			case "description":
				t.Description = val
			case "type":
				t.Type = val
			case "anchor":
				t.Anchor = val
			case "oneof_required":
				var typeFound *Schema
				for i := range parent.OneOf {
					if parent.OneOf[i].Title == nameValue[1] {
						typeFound = parent.OneOf[i]
					}
				}
				if typeFound == nil {
					typeFound = &Schema{
						Title:    nameValue[1],
						Required: []string{},
					}
					parent.OneOf = append(parent.OneOf, typeFound)
				}
				typeFound.Required = append(typeFound.Required, propertyName)
			case "anyof_required":
				var typeFound *Schema
				for i := range parent.AnyOf {
					if parent.AnyOf[i].Title == nameValue[1] {
						typeFound = parent.AnyOf[i]
					}
				}
				if typeFound == nil {
					typeFound = &Schema{
						Title:    nameValue[1],
						Required: []string{},
					}
					parent.AnyOf = append(parent.AnyOf, typeFound)
				}
				typeFound.Required = append(typeFound.Required, propertyName)
			case "oneof_ref":
				subSchema := t
				if t.Items != nil {
					subSchema = t.Items
				}
				if subSchema.OneOf == nil {
					subSchema.OneOf = make([]*Schema, 0, 1)
				}
				subSchema.Ref = ""
				refs := strings.Split(nameValue[1], ";")
				for _, r := range refs {
					subSchema.OneOf = append(subSchema.OneOf, &Schema{
						Ref: r,
					})
				}
			case "oneof_type":
				if t.OneOf == nil {
					t.OneOf = make([]*Schema, 0, 1)
				}
				t.Type = ""
				types := strings.Split(nameValue[1], ";")
				for _, ty := range types {
					t.OneOf = append(t.OneOf, &Schema{
						Type: ty,
					})
				}
			case "anyof_ref":
				subSchema := t
				if t.Items != nil {
					subSchema = t.Items
				}
				if subSchema.AnyOf == nil {
					subSchema.AnyOf = make([]*Schema, 0, 1)
				}
				subSchema.Ref = ""
				refs := strings.Split(nameValue[1], ";")
				for _, r := range refs {
					subSchema.AnyOf = append(subSchema.AnyOf, &Schema{
						Ref: r,
					})
				}
			case "anyof_type":
				if t.AnyOf == nil {
					t.AnyOf = make([]*Schema, 0, 1)
				}
				t.Type = ""
				types := strings.Split(nameValue[1], ";")
				for _, ty := range types {
					t.AnyOf = append(t.AnyOf, &Schema{
						Type: ty,
					})
				}
			default:
				unprocessed = append(unprocessed, tag)
			}
		}
	}
	return unprocessed
}

// read struct tags for boolean type keywords
func (t *Schema) booleanKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) != 2 {
			continue
		}
		name, val := nameValue[0], nameValue[1]
		if name == "default" {
			if val == "true" {
				t.Default = true
			} else if val == "false" {
				t.Default = false
			}
		}
	}
}

// read struct tags for string type keywords
func (t *Schema) stringKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "minLength":
				i, _ := strconv.Atoi(val)
				t.MinLength = i
			case "maxLength":
				i, _ := strconv.Atoi(val)
				t.MaxLength = i
			case "pattern":
				t.Pattern = val
			case "format":
				switch val {
				case "date-time", "email", "hostname", "ipv4", "ipv6", "uri", "uuid":
					t.Format = val
				}
			case "readOnly":
				i, _ := strconv.ParseBool(val)
				t.ReadOnly = i
			case "writeOnly":
				i, _ := strconv.ParseBool(val)
				t.WriteOnly = i
			case "default":
				t.Default = val
			case "example":
				t.Examples = append(t.Examples, val)
			case "enum":
				t.Enum = append(t.Enum, val)
			}
		}
	}
}

// read struct tags for numerical type keywords
func (t *Schema) numericalKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "multipleOf":
				t.MultipleOf, _ = toJSONNumber(val)
			case "minimum":
				t.Minimum, _ = toJSONNumber(val)
			case "maximum":
				t.Maximum, _ = toJSONNumber(val)
			case "exclusiveMaximum":
				t.ExclusiveMaximum, _ = toJSONNumber(val)
			case "exclusiveMinimum":
				t.ExclusiveMinimum, _ = toJSONNumber(val)
			case "default":
				if num, ok := toJSONNumber(val); ok {
					t.Default = num
				}
			case "example":
				if num, ok := toJSONNumber(val); ok {
					t.Examples = append(t.Examples, num)
				}
			case "enum":
				if num, ok := toJSONNumber(val); ok {
					t.Enum = append(t.Enum, num)
				}
			}
		}
	}
}

// read struct tags for object type keywords
// func (t *Type) objectKeywords(tags []string) {
//     for _, tag := range tags{
//         nameValue := strings.Split(tag, "=")
//         name, val := nameValue[0], nameValue[1]
//         switch name{
//             case "dependencies":
//                 t.Dependencies = val
//                 break;
//             case "patternProperties":
//                 t.PatternProperties = val
//                 break;
//         }
//     }
// }

// read struct tags for array type keywords
func (t *Schema) arrayKeywords(tags []string) {
	var defaultValues []any

	unprocessed := make([]string, 0, len(tags))
	for _, tag := range tags {
		nameValue := strings.Split(tag, "=")
		if len(nameValue) == 2 {
			name, val := nameValue[0], nameValue[1]
			switch name {
			case "minItems":
				i, _ := strconv.Atoi(val)
				t.MinItems = i
			case "maxItems":
				i, _ := strconv.Atoi(val)
				t.MaxItems = i
			case "uniqueItems":
				t.UniqueItems = true
			case "default":
				defaultValues = append(defaultValues, val)
			case "format":
				t.Items.Format = val
			case "pattern":
				t.Items.Pattern = val
			default:
				unprocessed = append(unprocessed, tag) // left for further processing by underlying type
			}
		}
	}
	if len(defaultValues) > 0 {
		t.Default = defaultValues
	}

	if len(unprocessed) == 0 {
		// we don't have anything else to process
		return
	}

	switch t.Items.Type {
	case "string":
		t.Items.stringKeywords(unprocessed)
	case "number":
		t.Items.numericalKeywords(unprocessed)
	case "integer":
		t.Items.numericalKeywords(unprocessed)
	case "array":
		// explicitly don't support traversal for the [][]..., as it's unclear where the array tags belong
	case "boolean":
		t.Items.booleanKeywords(unprocessed)
	}
}

func (t *Schema) extraKeywords(tags []string) {
	for _, tag := range tags {
		nameValue := strings.SplitN(tag, "=", 2)
		if len(nameValue) == 2 {
			t.setExtra(nameValue[0], nameValue[1])
		}
	}
}

func (t *Schema) setExtra(key, val string) {
	if t.Extras == nil {
		t.Extras = map[string]any{}
	}
	if existingVal, ok := t.Extras[key]; ok {
		switch existingVal := existingVal.(type) {
		case string:
			t.Extras[key] = []string{existingVal, val}
		case []string:
			t.Extras[key] = append(existingVal, val)
		case int:
			t.Extras[key], _ = strconv.Atoi(val)
		case bool:
			t.Extras[key] = (val == "true" || val == "t")
		}
	} else {
		switch key {
		case "minimum":
			t.Extras[key], _ = strconv.Atoi(val)
		default:
			var x any
			if val == "true" {
				x = true
			} else if val == "false" {
				x = false
			} else {
				x = val
			}
			t.Extras[key] = x
		}
	}
}

func requiredFromJSONTags(tags []string, val *bool) {
	if ignoredByJSONTags(tags) {
		return
	}

	for _, tag := range tags[1:] {
		if tag == "omitempty" {
			*val = false
			return
		}
	}
	*val = true
}

func requiredFromJSONSchemaTags(tags []string, val *bool) {
	if ignoredByJSONSchemaTags(tags) {
		return
	}
	for _, tag := range tags {
		if tag == "required" {
			*val = true
		}
	}
}

func nullableFromJSONSchemaTags(tags []string) bool {
	if ignoredByJSONSchemaTags(tags) {
		return false
	}
	for _, tag := range tags {
		if tag == "nullable" {
			return true
		}
	}
	return false
}

func ignoredByJSONTags(tags []string) bool {
	return tags[0] == "-"
}

func ignoredByJSONSchemaTags(tags []string) bool {
	return tags[0] == "-"
}

// toJSONNumber converts string to *json.Number.
// It'll aso return whether the number is valid.
func toJSONNumber(s string) (json.Number, bool) {
	num := json.Number(s)
	if _, err := num.Int64(); err == nil {
		return num, true
	}
	if _, err := num.Float64(); err == nil {
		return num, true
	}
	return json.Number(""), false
}

func (r *Reflector) fieldNameTag() string {
	if r.FieldNameTag != "" {
		return r.FieldNameTag
	}
	return "json"
}

func (r *Reflector) reflectFieldName(f reflect.StructField) (string, string, bool, bool, bool) {
	jsonTagString := f.Tag.Get(r.fieldNameTag())
	jsonTags := strings.Split(jsonTagString, ",")

	if ignoredByJSONTags(jsonTags) {
		return "", "", false, false, false
	}

	schemaTags := strings.Split(f.Tag.Get("jsonschema"), ",")
	if ignoredByJSONSchemaTags(schemaTags) {
		return "", "", false, false, false
	}

	var required bool
	if !r.RequiredFromJSONSchemaTags {
		requiredFromJSONTags(jsonTags, &required)
	}
	requiredFromJSONSchemaTags(schemaTags, &required)

	nullable := nullableFromJSONSchemaTags(schemaTags)

	if f.Anonymous && jsonTags[0] == "" {
		// As per JSON Marshal rules, anonymous structs are inherited
		if f.Type.Kind() == reflect.Struct {
			return "", "", true, false, false
		}

		// As per JSON Marshal rules, anonymous pointer to structs are inherited
		if f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct {
			return "", "", true, false, false
		}
	}

	// Try to determine the name from the different combos
	name := f.Name
	originalName := f.Name
	if jsonTags[0] != "" {
		name = jsonTags[0]
	}
	if !f.Anonymous && f.PkgPath != "" {
		// field not anonymous and not export has no export name
		name = ""
	} else if r.KeyNamer != nil {
		name = r.KeyNamer(name)
	}

	return name, originalName, false, required, nullable
}

// UnmarshalJSON is used to parse a schema object or boolean.
func (t *Schema) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("true")) {
		*t = *TrueSchema
		return nil
	} else if bytes.Equal(data, []byte("false")) {
		*t = *FalseSchema
		return nil
	}
	type SchemaAlt Schema
	aux := &struct {
		*SchemaAlt
	}{
		SchemaAlt: (*SchemaAlt)(t),
	}
	return json.Unmarshal(data, aux)
}

// MarshalJSON is used to serialize a schema object or boolean.
func (t *Schema) MarshalJSON() ([]byte, error) {
	if t.boolean != nil {
		if *t.boolean {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	}
	if reflect.DeepEqual(&Schema{}, t) {
		// Don't bother returning empty schemas
		return []byte("true"), nil
	}
	type SchemaAlt Schema
	b, err := json.Marshal((*SchemaAlt)(t))
	if err != nil {
		return nil, err
	}
	if t.Extras == nil || len(t.Extras) == 0 {
		return b, nil
	}
	m, err := json.Marshal(t.Extras)
	if err != nil {
		return nil, err
	}
	if len(b) == 2 {
		return m, nil
	}
	b[len(b)-1] = ','
	return append(b, m[1:]...), nil
}

func (r *Reflector) typeName(t reflect.Type) string {
	if r.Namer != nil {
		if name := r.Namer(t); name != "" {
			return name
		}
	}
	return t.Name()
}

// Split on commas that are not preceded by `\`.
// This way, we prevent splitting regexes
func splitOnUnescapedCommas(tagString string) []string {
	ret := make([]string, 0)
	separated := strings.Split(tagString, ",")
	ret = append(ret, separated[0])
	i := 0
	for _, nextTag := range separated[1:] {
		if len(ret[i]) == 0 {
			ret = append(ret, nextTag)
			i++
			continue
		}

		if ret[i][len(ret[i])-1] == '\\' {
			ret[i] = ret[i][:len(ret[i])-1] + "," + nextTag
		} else {
			ret = append(ret, nextTag)
			i++
		}
	}

	return ret
}

func fullyQualifiedTypeName(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}

// AddGoComments will update the reflectors comment map with all the comments
// found in the provided source directories. See the #ExtractGoComments method
// for more details.
func (r *Reflector) AddGoComments(base, path string) error {
	if r.CommentMap == nil {
		r.CommentMap = make(map[string]string)
	}
	return ExtractGoComments(base, path, r.CommentMap)
}
