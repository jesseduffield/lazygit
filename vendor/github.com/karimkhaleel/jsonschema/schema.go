package jsonschema

import (
	"encoding/json"

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
	Properties           *orderedmap.OrderedMap[string, *Schema] `json:"properties,omitempty"`           // section 10.3.2.1
	PatternProperties    map[string]*Schema                      `json:"patternProperties,omitempty"`    // section 10.3.2.2
	AdditionalProperties *Schema                                 `json:"additionalProperties,omitempty"` // section 10.3.2.3
	PropertyNames        *Schema                                 `json:"propertyNames,omitempty"`        // section 10.3.2.4
	// RFC draft-bhutton-json-schema-validation-00, section 6
	Type              string              `json:"type,omitempty"`              // section 6.1.1
	Enum              []any               `json:"enum,omitempty"`              // section 6.1.2
	Const             any                 `json:"const,omitempty"`             // section 6.1.3
	MultipleOf        json.Number         `json:"multipleOf,omitempty"`        // section 6.2.1
	Maximum           json.Number         `json:"maximum,omitempty"`           // section 6.2.2
	ExclusiveMaximum  json.Number         `json:"exclusiveMaximum,omitempty"`  // section 6.2.3
	Minimum           json.Number         `json:"minimum,omitempty"`           // section 6.2.4
	ExclusiveMinimum  json.Number         `json:"exclusiveMinimum,omitempty"`  // section 6.2.5
	MaxLength         *uint64             `json:"maxLength,omitempty"`         // section 6.3.1
	MinLength         *uint64             `json:"minLength,omitempty"`         // section 6.3.2
	Pattern           string              `json:"pattern,omitempty"`           // section 6.3.3
	MaxItems          *uint64             `json:"maxItems,omitempty"`          // section 6.4.1
	MinItems          *uint64             `json:"minItems,omitempty"`          // section 6.4.2
	UniqueItems       bool                `json:"uniqueItems,omitempty"`       // section 6.4.3
	MaxContains       *uint64             `json:"maxContains,omitempty"`       // section 6.4.4
	MinContains       *uint64             `json:"minContains,omitempty"`       // section 6.4.5
	MaxProperties     *uint64             `json:"maxProperties,omitempty"`     // section 6.5.1
	MinProperties     *uint64             `json:"minProperties,omitempty"`     // section 6.5.2
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

// Definitions hold schema definitions.
// http://json-schema.org/latest/json-schema-validation.html#rfc.section.5.26
// RFC draft-wright-json-schema-validation-00, section 5.26
type Definitions map[string]*Schema
