// +build codegen

package api

import (
	"fmt"
	"regexp"
	"strings"
)

// updateTopLevelShapeReferences moves resultWrapper, locationName, and
// xmlNamespace traits from toplevel shape references to the toplevel
// shapes for easier code generation
func (a *API) updateTopLevelShapeReferences() {
	for _, o := range a.Operations {
		// these are for REST-XML services
		if o.InputRef.LocationName != "" {
			o.InputRef.Shape.LocationName = o.InputRef.LocationName
		}
		if o.InputRef.Location != "" {
			o.InputRef.Shape.Location = o.InputRef.Location
		}
		if o.InputRef.Payload != "" {
			o.InputRef.Shape.Payload = o.InputRef.Payload
		}
		if o.InputRef.XMLNamespace.Prefix != "" {
			o.InputRef.Shape.XMLNamespace.Prefix = o.InputRef.XMLNamespace.Prefix
		}
		if o.InputRef.XMLNamespace.URI != "" {
			o.InputRef.Shape.XMLNamespace.URI = o.InputRef.XMLNamespace.URI
		}
	}

}

// writeShapeNames sets each shape's API and shape name values. Binding the
// shape to its parent API.
func (a *API) writeShapeNames() {
	for n, s := range a.Shapes {
		s.API = a
		s.ShapeName = n
	}
}

func (a *API) resolveReferences() {
	resolver := referenceResolver{API: a, visited: map[*ShapeRef]bool{}}

	for _, s := range a.Shapes {
		resolver.resolveShape(s)
	}

	for _, o := range a.Operations {
		o.API = a // resolve parent reference

		resolver.resolveReference(&o.InputRef)
		resolver.resolveReference(&o.OutputRef)

		// Resolve references for errors also
		for i := range o.ErrorRefs {
			resolver.resolveReference(&o.ErrorRefs[i])
			o.ErrorRefs[i].Shape.Exception = true
			o.ErrorRefs[i].Shape.ErrorInfo.Type = o.ErrorRefs[i].Shape.ShapeName
		}
	}
}

// A referenceResolver provides a way to resolve shape references to
// shape definitions.
type referenceResolver struct {
	*API
	visited map[*ShapeRef]bool
}

// resolveReference updates a shape reference to reference the API and
// its shape definition. All other nested references are also resolved.
func (r *referenceResolver) resolveReference(ref *ShapeRef) {
	if ref.ShapeName == "" {
		return
	}

	shape, ok := r.API.Shapes[ref.ShapeName]
	if !ok {
		panic(fmt.Sprintf("unable resolve reference, %s", ref.ShapeName))
	}

	if ref.JSONValue {
		ref.ShapeName = "JSONValue"
		if _, ok := r.API.Shapes[ref.ShapeName]; !ok {
			r.API.Shapes[ref.ShapeName] = &Shape{
				API:       r.API,
				ShapeName: "JSONValue",
				Type:      "jsonvalue",
				ValueRef: ShapeRef{
					JSONValue: true,
				},
			}
		}
	}

	ref.API = r.API   // resolve reference back to API
	ref.Shape = shape // resolve shape reference

	if r.visited[ref] {
		return
	}
	r.visited[ref] = true

	shape.refs = append(shape.refs, ref) // register the ref

	// resolve shape's references, if it has any
	r.resolveShape(shape)
}

// resolveShape resolves a shape's Member Key Value, and nested member
// shape references.
func (r *referenceResolver) resolveShape(shape *Shape) {
	r.resolveReference(&shape.MemberRef)
	r.resolveReference(&shape.KeyRef)
	r.resolveReference(&shape.ValueRef)
	for _, m := range shape.MemberRefs {
		r.resolveReference(m)
	}
}

// fixStutterNames fixes all name struttering based on Go naming conventions.
// "Stuttering" is when the prefix of a structure or function matches the
// package name (case insensitive).
func (a *API) fixStutterNames() {
	str, end := a.StructName(), ""
	if len(str) > 1 {
		l := len(str) - 1
		str, end = str[0:l], str[l:]
	}
	re := regexp.MustCompile(fmt.Sprintf(`\A(?i:%s)%s`, str, end))

	for name, op := range a.Operations {
		newName := re.ReplaceAllString(name, "")
		if newName != name && len(newName) > 0 {
			delete(a.Operations, name)
			a.Operations[newName] = op
		}
		op.ExportedName = newName
	}

	for k, s := range a.Shapes {
		newName := re.ReplaceAllString(k, "")
		if newName != s.ShapeName && len(newName) > 0 {
			s.Rename(newName)
		}
	}
}

// renameExportable renames all operation names to be exportable names.
// All nested Shape names are also updated to the exportable variant.
func (a *API) renameExportable() {
	for name, op := range a.Operations {
		newName := a.ExportableName(name)
		if newName != name {
			delete(a.Operations, name)
			a.Operations[newName] = op
		}
		op.ExportedName = newName
	}

	for k, s := range a.Shapes {
		// FIXME SNS has lower and uppercased shape names with the same name,
		// except the lowercased variant is used exclusively for string and
		// other primitive types. Renaming both would cause a collision.
		// We work around this by only renaming the structure shapes.
		if s.Type == "string" {
			continue
		}

		for mName, member := range s.MemberRefs {
			ref := s.MemberRefs[mName]
			ref.OrigShapeName = mName
			s.MemberRefs[mName] = ref

			newName := a.ExportableName(mName)
			if newName != mName {
				delete(s.MemberRefs, mName)
				s.MemberRefs[newName] = member

				// also apply locationName trait so we keep the old one
				// but only if there's no locationName trait on ref or shape
				if member.LocationName == "" && member.Shape.LocationName == "" {
					member.LocationName = mName
				}
			}

			if newName == "_" {
				panic("Shape " + s.ShapeName + " uses reserved member name '_'")
			}
		}

		newName := a.ExportableName(k)
		if newName != s.ShapeName {
			s.Rename(newName)
		}

		s.Payload = a.ExportableName(s.Payload)

		// fix required trait names
		for i, n := range s.Required {
			s.Required[i] = a.ExportableName(n)
		}
	}

	for _, s := range a.Shapes {
		// fix enum names
		if s.IsEnum() {
			s.EnumConsts = make([]string, len(s.Enum))
			for i := range s.Enum {
				shape := s.ShapeName
				shape = strings.ToUpper(shape[0:1]) + shape[1:]
				s.EnumConsts[i] = shape + s.EnumName(i)
			}
		}
	}
}

// renameCollidingFields will rename any fields that uses an SDK or Golang
// specific name.
func (a *API) renameCollidingFields() {
	for _, v := range a.Shapes {
		namesWithSet := map[string]struct{}{}
		for k, field := range v.MemberRefs {
			if _, ok := v.MemberRefs["Set"+k]; ok {
				namesWithSet["Set"+k] = struct{}{}
			}

			if collides(k) || (v.Exception && exceptionCollides(k)) {
				renameCollidingField(k, v, field)
			}
		}

		// checks if any field names collide with setters.
		for name := range namesWithSet {
			field := v.MemberRefs[name]
			renameCollidingField(name, v, field)
		}
	}
}

func renameCollidingField(name string, v *Shape, field *ShapeRef) {
	newName := name + "_"
	debugLogger.Logf("Shape %s's field %q renamed to %q", v.ShapeName, name, newName)
	delete(v.MemberRefs, name)
	v.MemberRefs[newName] = field
}

// collides will return true if it is a name used by the SDK or Golang.
func collides(name string) bool {
	switch name {
	case "String",
		"GoString",
		"Validate":
		return true
	}
	return false
}

func exceptionCollides(name string) bool {
	switch name {
	case "Code",
		"Message",
		"OrigErr":
		return true
	}
	return false
}

func (a *API) applyShapeNameAliases() {
	service, ok := shapeNameAliases[a.name]
	if !ok {
		return
	}

	// Generic Shape Aliases
	for name, s := range a.Shapes {
		if alias, ok := service[name]; ok {
			s.Rename(alias)
			s.AliasedShapeName = true
		}
	}
}

// createInputOutputShapes creates toplevel input/output shapes if they
// have not been defined in the API. This normalizes all APIs to always
// have an input and output structure in the signature.
func (a *API) createInputOutputShapes() {
	for _, op := range a.Operations {
		createAPIParamShape(a, op.Name, &op.InputRef, op.ExportedName+"Input",
			shamelist.Input,
		)
		createAPIParamShape(a, op.Name, &op.OutputRef, op.ExportedName+"Output",
			shamelist.Output,
		)
	}
}

func (a *API) renameAPIPayloadShapes() {
	for _, op := range a.Operations {
		op.InputRef.Payload = a.ExportableName(op.InputRef.Payload)
		op.OutputRef.Payload = a.ExportableName(op.OutputRef.Payload)
	}
}

func createAPIParamShape(a *API, opName string, ref *ShapeRef, shapeName string, shamelistLookup func(string, string) bool) {
	if len(ref.ShapeName) == 0 {
		setAsPlacholderShape(ref, shapeName, a)
		return
	}

	// nothing to do if already the correct name.
	if s := ref.Shape; s.AliasedShapeName || s.ShapeName == shapeName || shamelistLookup(a.name, opName) {
		return
	}

	if s, ok := a.Shapes[shapeName]; ok {
		panic(fmt.Sprintf(
			"attempting to create duplicate API parameter shape, %v, %v, %v, %v\n",
			shapeName, opName, ref.ShapeName, s.OrigShapeName,
		))
	}

	ref.Shape.removeRef(ref)
	ref.OrigShapeName = shapeName
	ref.ShapeName = shapeName
	ref.Shape = ref.Shape.Clone(shapeName)
	ref.Shape.refs = append(ref.Shape.refs, ref)
}

func setAsPlacholderShape(tgtShapeRef *ShapeRef, name string, a *API) {
	shape := a.makeIOShape(name)
	shape.Placeholder = true
	*tgtShapeRef = ShapeRef{API: a, ShapeName: shape.ShapeName, Shape: shape}
	shape.refs = append(shape.refs, tgtShapeRef)
}

// makeIOShape returns a pointer to a new Shape initialized by the name provided.
func (a *API) makeIOShape(name string) *Shape {
	shape := &Shape{
		API: a, ShapeName: name, Type: "structure",
		MemberRefs: map[string]*ShapeRef{},
	}
	a.Shapes[name] = shape
	return shape
}

// removeUnusedShapes removes shapes from the API which are not referenced by any
// other shape in the API.
func (a *API) removeUnusedShapes() {
	for _, s := range a.Shapes {
		if len(s.refs) == 0 {
			a.removeShape(s)
		}
	}
}

// Represents the service package name to EndpointsID mapping
var custEndpointsKey = map[string]string{
	"applicationautoscaling": "application-autoscaling",
}

// Sents the EndpointsID field of Metadata  with the value of the
// EndpointPrefix if EndpointsID is not set. Also adds
// customizations for services if EndpointPrefix is not a valid key.
func (a *API) setMetadataEndpointsKey() {
	if len(a.Metadata.EndpointsID) != 0 {
		return
	}

	if v, ok := custEndpointsKey[a.PackageName()]; ok {
		a.Metadata.EndpointsID = v
	} else {
		a.Metadata.EndpointsID = a.Metadata.EndpointPrefix
	}
}

func (a *API) findEndpointDiscoveryOp() {
	for _, op := range a.Operations {
		if op.IsEndpointDiscoveryOp {
			a.EndpointDiscoveryOp = op
			return
		}
	}
}
