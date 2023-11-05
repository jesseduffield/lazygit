# Go JSON Schema Reflection

[![Lint](https://github.com/invopop/jsonschema/actions/workflows/lint.yaml/badge.svg)](https://github.com/invopop/jsonschema/actions/workflows/lint.yaml)
[![Test Go](https://github.com/invopop/jsonschema/actions/workflows/test.yaml/badge.svg)](https://github.com/invopop/jsonschema/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/jsonschema)](https://goreportcard.com/report/github.com/invopop/jsonschema)
[![GoDoc](https://godoc.org/github.com/invopop/jsonschema?status.svg)](https://godoc.org/github.com/invopop/jsonschema)
![Latest Tag](https://img.shields.io/github/v/tag/invopop/jsonschema)

This package can be used to generate [JSON Schemas](http://json-schema.org/latest/json-schema-validation.html) from Go types through reflection.

- Supports arbitrarily complex types, including `interface{}`, maps, slices, etc.
- Supports json-schema features such as minLength, maxLength, pattern, format, etc.
- Supports simple string and numeric enums.
- Supports custom property fields via the `jsonschema_extras` struct tag.

This repository is a fork of the original [jsonschema](https://github.com/alecthomas/jsonschema) by [@alecthomas](https://github.com/alecthomas). At [Invopop](https://invopop.com) we use jsonschema as a cornerstone in our [GOBL library](https://github.com/invopop/gobl), and wanted to be able to continue building and adding features without taking up Alec's time. There have been a few significant changes that probably mean this version is a not compatible with with Alec's:

- The original was stuck on the draft-04 version of JSON Schema, we've now moved to the latest JSON Schema Draft 2020-12.
- Schema IDs are added automatically from the current Go package's URL in order to be unique, and can be disabled with the `Anonymous` option.
- Support for the `FullyQualifyTypeName` option has been removed. If you have conflicts, you should use multiple schema files with different IDs, set the `DoNotReference` option to true to hide definitions completely, or add your own naming strategy using the `Namer` property.
- Support for `yaml` tags and related options has been dropped for the sake of simplification. There were a [few inconsistencies](https://github.com/invopop/jsonschema/pull/21) around this that have now been fixed.

## Versions

This project is still under v0 scheme, as per Go convention, breaking changes are likely. Please pin go modules to version tags or branches, and reach out if you think something can be improved.

Go version >= 1.18 is required as generics are now being used.

## Example

The following Go type:

```go
type TestUser struct {
  ID            int                    `json:"id"`
  Name          string                 `json:"name" jsonschema:"title=the name,description=The name of a friend,example=joe,example=lucy,default=alex"`
  Friends       []int                  `json:"friends,omitempty" jsonschema_description:"The list of IDs, omitted when empty"`
  Tags          map[string]interface{} `json:"tags,omitempty" jsonschema_extras:"a=b,foo=bar,foo=bar1"`
  BirthDate     time.Time              `json:"birth_date,omitempty" jsonschema:"oneof_required=date"`
  YearOfBirth   string                 `json:"year_of_birth,omitempty" jsonschema:"oneof_required=year"`
  Metadata      interface{}            `json:"metadata,omitempty" jsonschema:"oneof_type=string;array"`
  FavColor      string                 `json:"fav_color,omitempty" jsonschema:"enum=red,enum=green,enum=blue"`
}
```

Results in following JSON Schema:

```go
jsonschema.Reflect(&TestUser{})
```

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://github.com/invopop/jsonschema_test/sample-user",
  "$ref": "#/$defs/SampleUser",
  "$defs": {
    "SampleUser": {
      "oneOf": [
        {
          "required": ["birth_date"],
          "title": "date"
        },
        {
          "required": ["year_of_birth"],
          "title": "year"
        }
      ],
      "properties": {
        "id": {
          "type": "integer"
        },
        "name": {
          "type": "string",
          "title": "the name",
          "description": "The name of a friend",
          "default": "alex",
          "examples": ["joe", "lucy"]
        },
        "friends": {
          "items": {
            "type": "integer"
          },
          "type": "array",
          "description": "The list of IDs, omitted when empty"
        },
        "tags": {
          "type": "object",
          "a": "b",
          "foo": ["bar", "bar1"]
        },
        "birth_date": {
          "type": "string",
          "format": "date-time"
        },
        "year_of_birth": {
          "type": "string"
        },
        "metadata": {
          "oneOf": [
            {
              "type": "string"
            },
            {
              "type": "array"
            }
          ]
        },
        "fav_color": {
          "type": "string",
          "enum": ["red", "green", "blue"]
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": ["id", "name"]
    }
  }
}
```

## YAML

Support for `yaml` tags has now been removed. If you feel very strongly about this, we've opened a discussion to hear your comments: https://github.com/invopop/jsonschema/discussions/28

The recommended approach if you need to deal with YAML data is to first convert to JSON. The [invopop/yaml](https://github.com/invopop/yaml) library will make this trivial.

## Configurable behaviour

The behaviour of the schema generator can be altered with parameters when a `jsonschema.Reflector`
instance is created.

### ExpandedStruct

If set to `true`, makes the top level struct not to reference itself in the definitions. But type passed should be a struct type.

eg.

```go
type GrandfatherType struct {
	FamilyName string `json:"family_name" jsonschema:"required"`
}

type SomeBaseType struct {
	SomeBaseProperty int `json:"some_base_property"`
	// The jsonschema required tag is nonsensical for private and ignored properties.
	// Their presence here tests that the fields *will not* be required in the output
	// schema, even if they are tagged required.
	somePrivateBaseProperty            string `json:"i_am_private" jsonschema:"required"`
	SomeIgnoredBaseProperty            string `json:"-" jsonschema:"required"`
	SomeSchemaIgnoredProperty          string `jsonschema:"-,required"`
	SomeUntaggedBaseProperty           bool   `jsonschema:"required"`
	someUnexportedUntaggedBaseProperty bool
	Grandfather                        GrandfatherType `json:"grand"`
}
```

will output:

```json
{
  "$schema": "http://json-schema.org/draft/2020-12/schema",
  "required": ["some_base_property", "grand", "SomeUntaggedBaseProperty"],
  "properties": {
    "SomeUntaggedBaseProperty": {
      "type": "boolean"
    },
    "grand": {
      "$schema": "http://json-schema.org/draft/2020-12/schema",
      "$ref": "#/definitions/GrandfatherType"
    },
    "some_base_property": {
      "type": "integer"
    }
  },
  "type": "object",
  "$defs": {
    "GrandfatherType": {
      "required": ["family_name"],
      "properties": {
        "family_name": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  }
}
```

### Using Go Comments

Writing a good schema with descriptions inside tags can become cumbersome and tedious, especially if you already have some Go comments around your types and field definitions. If you'd like to take advantage of these existing comments, you can use the `AddGoComments(base, path string)` method that forms part of the reflector to parse your go files and automatically generate a dictionary of Go import paths, types, and fields, to individual comments. These will then be used automatically as description fields, and can be overridden with a manual definition if needed.

Take a simplified example of a User struct which for the sake of simplicity we assume is defined inside this package:

```go
package main

// User is used as a base to provide tests for comments.
type User struct {
	// Unique sequential identifier.
	ID int `json:"id" jsonschema:"required"`
	// Name of the user
	Name string `json:"name"`
}
```

To get the comments provided into your JSON schema, use a regular `Reflector` and add the go code using an import module URL and path. Fully qualified go module paths cannot be determined reliably by the `go/parser` library, so we need to introduce this manually:

```go
r := new(Reflector)
if err := r.AddGoComments("github.com/invopop/jsonschema", "./"); err != nil {
  // deal with error
}
s := r.Reflect(&User{})
// output
```

Expect the results to be similar to:

```json
{
  "$schema": "http://json-schema.org/draft/2020-12/schema",
  "$ref": "#/$defs/User",
  "$defs": {
    "User": {
      "required": ["id"],
      "properties": {
        "id": {
          "type": "integer",
          "description": "Unique sequential identifier."
        },
        "name": {
          "type": "string",
          "description": "Name of the user"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "description": "User is used as a base to provide tests for comments."
    }
  }
}
```

### Custom Key Naming

In some situations, the keys actually used to write files are different from Go structs'.

This is often the case when writing a configuration file to YAML or JSON from a Go struct, or when returning a JSON response for a Web API: APIs typically use snake_case, while Go uses PascalCase.

You can pass a `func(string) string` function to `Reflector`'s `KeyNamer` option to map Go field names to JSON key names and reflect the aforementioned transformations, without having to specify `json:"..."` on every struct field.

For example, consider the following struct

```go
type User struct {
  GivenName       string
  PasswordSalted  []byte `json:"salted_password"`
}
```

We can transform field names to snake_case in the generated JSON schema:

```go
r := new(jsonschema.Reflector)
r.KeyNamer = strcase.SnakeCase // from package github.com/stoewer/go-strcase

r.Reflect(&User{})
```

Will yield

```diff
  {
    "$schema": "http://json-schema.org/draft/2020-12/schema",
    "$ref": "#/$defs/User",
    "$defs": {
      "User": {
        "properties": {
-         "GivenName": {
+         "given_name": {
            "type": "string"
          },
          "salted_password": {
            "type": "string",
            "contentEncoding": "base64"
          }
        },
        "additionalProperties": false,
        "type": "object",
-       "required": ["GivenName", "salted_password"]
+       "required": ["given_name", "salted_password"]
      }
    }
  }
```

As you can see, if a field name has a `json:""` tag set, the `key` argument to `KeyNamer` will have the value of that tag.

### Custom Type Definitions

Sometimes it can be useful to have custom JSON Marshal and Unmarshal methods in your structs that automatically convert for example a string into an object.

To override auto-generating an object type for your type, implement the `JSONSchema() *Schema` method and whatever is defined will be provided in the schema definitions.

You also have the option of defining a `JSONSchemaExtend(schema *jsonschema.Schema)` method for your types that will be called _after_ the schema has been generated, allowing you to add or manipulate the fields easily.

Take the following simplified example of a `CompactDate` that only includes the Year and Month:

```go
type CompactDate struct {
	Year  int
	Month int
}

func (d *CompactDate) UnmarshalJSON(data []byte) error {
  if len(data) != 9 {
    return errors.New("invalid compact date length")
  }
  var err error
  d.Year, err = strconv.Atoi(string(data[1:5]))
  if err != nil {
    return err
  }
  d.Month, err = strconv.Atoi(string(data[7:8]))
  if err != nil {
    return err
  }
  return nil
}

func (d *CompactDate) MarshalJSON() ([]byte, error) {
  buf := new(bytes.Buffer)
  buf.WriteByte('"')
  buf.WriteString(fmt.Sprintf("%d-%02d", d.Year, d.Month))
  buf.WriteByte('"')
  return buf.Bytes(), nil
}

func (CompactDate) JSONSchema() *Schema {
	return &Schema{
		Type:        "string",
		Title:       "Compact Date",
		Description: "Short date that only includes year and month",
		Pattern:     "^[0-9]{4}-[0-1][0-9]$",
	}
}
```

The resulting schema generated for this struct would look like:

```json
{
  "$schema": "http://json-schema.org/draft/2020-12/schema",
  "$ref": "#/$defs/CompactDate",
  "$defs": {
    "CompactDate": {
      "pattern": "^[0-9]{4}-[0-1][0-9]$",
      "type": "string",
      "title": "Compact Date",
      "description": "Short date that only includes year and month"
    }
  }
}
```
