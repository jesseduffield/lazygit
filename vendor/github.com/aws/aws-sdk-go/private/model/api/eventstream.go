// +build codegen

package api

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

// EventStreamAPI provides details about the event stream async API and
// associated EventStream shapes.
type EventStreamAPI struct {
	API       *API
	Name      string
	Operation *Operation
	Shape     *Shape
	Inbound   *EventStream
	Outbound  *EventStream
}

// EventStream represents a single eventstream group (input/output) and the
// modeled events that are known for the stream.
type EventStream struct {
	Name       string
	Shape      *Shape
	Events     []*Event
	Exceptions []*Event
}

// Event is a single EventStream event that can be sent or received in an
// EventStream.
type Event struct {
	Name  string
	Shape *Shape
	For   *EventStream
}

// ShapeDoc returns the docstring for the EventStream API.
func (esAPI *EventStreamAPI) ShapeDoc() string {
	tmpl := template.Must(template.New("eventStreamShapeDoc").Parse(`
{{- $.Name }} provides handling of EventStreams for
the {{ $.Operation.ExportedName }} API.
{{- if $.Inbound }}

Use this type to receive {{ $.Inbound.Name }} events. The events
can be read from the Events channel member.

The events that can be received are:
{{ range $_, $event := $.Inbound.Events }}
    * {{ $event.Shape.ShapeName }}
{{- end }}

{{- end }}

{{- if $.Outbound }}

Use this type to send {{ $.Outbound.Name }} events. The events
can be sent with the Send method.

The events that can be sent are:
{{ range $_, $event := $.Outbound.Events -}}
    * {{ $event.Shape.ShapeName }}
{{- end }}

{{- end }}`))

	var w bytes.Buffer
	if err := tmpl.Execute(&w, esAPI); err != nil {
		panic(fmt.Sprintf("failed to generate eventstream shape template for %v, %v", esAPI.Name, err))
	}

	return commentify(w.String())
}

func hasEventStream(topShape *Shape) bool {
	for _, ref := range topShape.MemberRefs {
		if ref.Shape.IsEventStream {
			return true
		}
	}

	return false
}

func eventStreamAPIShapeRefDoc(refName string) string {
	return commentify(fmt.Sprintf("Use %s to use the API's stream.", refName))
}

func (a *API) setupEventStreams() {
	const eventStreamMemberName = "EventStream"

	for _, op := range a.Operations {
		outbound := setupEventStream(op.InputRef.Shape)
		inbound := setupEventStream(op.OutputRef.Shape)

		if outbound == nil && inbound == nil {
			continue
		}

		if outbound != nil {
			panic(fmt.Sprintf("Outbound stream support not implemented, %s, %s",
				outbound.Name, outbound.Shape.ShapeName))
		}

		switch a.Metadata.Protocol {
		case `rest-json`, `rest-xml`, `json`:
		default:
			panic(fmt.Sprintf("EventStream not supported for protocol %v",
				a.Metadata.Protocol))
		}

		op.EventStreamAPI = &EventStreamAPI{
			API:       a,
			Name:      op.ExportedName + eventStreamMemberName,
			Operation: op,
			Outbound:  outbound,
			Inbound:   inbound,
		}

		streamShape := &Shape{
			API:            a,
			ShapeName:      op.EventStreamAPI.Name,
			Documentation:  op.EventStreamAPI.ShapeDoc(),
			Type:           "structure",
			EventStreamAPI: op.EventStreamAPI,
			IsEventStream:  true,
			MemberRefs: map[string]*ShapeRef{
				"Inbound": &ShapeRef{
					ShapeName: inbound.Shape.ShapeName,
				},
			},
		}
		inbound.Shape.refs = append(inbound.Shape.refs, streamShape.MemberRefs["Inbound"])
		streamShapeRef := &ShapeRef{
			API:           a,
			ShapeName:     streamShape.ShapeName,
			Shape:         streamShape,
			Documentation: eventStreamAPIShapeRefDoc(eventStreamMemberName),
		}
		streamShape.refs = []*ShapeRef{streamShapeRef}
		op.EventStreamAPI.Shape = streamShape

		if _, ok := op.OutputRef.Shape.MemberRefs[eventStreamMemberName]; ok {
			panic(fmt.Sprintf("shape ref already exists, %s.%s",
				op.OutputRef.Shape.ShapeName, eventStreamMemberName))
		}
		op.OutputRef.Shape.MemberRefs[eventStreamMemberName] = streamShapeRef
		op.OutputRef.Shape.EventStreamsMemberName = eventStreamMemberName
		if _, ok := a.Shapes[streamShape.ShapeName]; ok {
			panic("shape already exists, " + streamShape.ShapeName)
		}
		a.Shapes[streamShape.ShapeName] = streamShape

		a.HasEventStream = true
	}
}

func setupEventStream(topShape *Shape) *EventStream {
	var eventStream *EventStream
	for refName, ref := range topShape.MemberRefs {
		if !ref.Shape.IsEventStream {
			continue
		}
		if eventStream != nil {
			panic(fmt.Sprintf("multiple shape ref eventstreams, %s, prev: %s",
				refName, eventStream.Name))
		}

		eventStream = &EventStream{
			Name:  ref.Shape.ShapeName,
			Shape: ref.Shape,
		}

		if topShape.API.Metadata.Protocol == "json" {
			topShape.EventFor = append(topShape.EventFor, eventStream)
		}

		for _, eventRefName := range ref.Shape.MemberNames() {
			eventRef := ref.Shape.MemberRefs[eventRefName]
			if !(eventRef.Shape.IsEvent || eventRef.Shape.Exception) {
				panic(fmt.Sprintf("unexpected non-event member reference %s.%s",
					ref.Shape.ShapeName, eventRefName))
			}

			updateEventPayloadRef(eventRef.Shape)

			eventRef.Shape.EventFor = append(eventRef.Shape.EventFor, eventStream)

			// Exceptions and events are two different lists to allow the SDK
			// to easly generate code with the two handled differently.
			event := &Event{
				Name:  eventRefName,
				Shape: eventRef.Shape,
				For:   eventStream,
			}
			if eventRef.Shape.Exception {
				eventStream.Exceptions = append(eventStream.Exceptions, event)
			} else {
				eventStream.Events = append(eventStream.Events, event)
			}
		}

		// Remove the eventstream references as they will be added elsewhere.
		ref.Shape.removeRef(ref)
		delete(topShape.MemberRefs, refName)
		delete(topShape.API.Shapes, ref.Shape.ShapeName)
	}

	return eventStream
}

func updateEventPayloadRef(parent *Shape) {
	refName := parent.PayloadRefName()
	if len(refName) == 0 {
		return
	}

	payloadRef := parent.MemberRefs[refName]

	if payloadRef.Shape.Type == "blob" {
		return
	}

	if len(payloadRef.LocationName) != 0 {
		return
	}

	payloadRef.LocationName = refName
}

func renderEventStreamAPIShape(w io.Writer, s *Shape) error {
	// Imports needed by the EventStream APIs.
	s.API.AddImport("fmt")
	s.API.AddImport("bytes")
	s.API.AddImport("io")
	s.API.AddImport("sync")
	s.API.AddImport("sync/atomic")
	s.API.AddSDKImport("aws")
	s.API.AddSDKImport("aws/awserr")
	s.API.AddSDKImport("private/protocol/eventstream")
	s.API.AddSDKImport("private/protocol/eventstream/eventstreamapi")

	return eventStreamAPIShapeTmpl.Execute(w, s)
}

// Template for an EventStream API Shape that will provide read/writing events
// across the EventStream. This is a special shape that's only public members
// are the Events channel and a Close and Err method.
//
// Executed in the context of a Shape.
var eventStreamAPIShapeTmpl = func() *template.Template {
	t := template.Must(
		template.New("eventStreamAPIShapeTmpl").
			Funcs(template.FuncMap{}).
			Parse(eventStreamAPITmplDef),
	)

	template.Must(
		t.AddParseTree(
			"eventStreamAPIReaderTmpl", eventStreamAPIReaderTmpl.Tree),
	)

	return t
}()

const eventStreamAPITmplDef = `
{{ $.Documentation }}
type {{ $.ShapeName }} struct {
	{{- if $.EventStreamAPI.Inbound }}
		// Reader is the EventStream reader for the {{ $.EventStreamAPI.Inbound.Name }}
		// events. This value is automatically set by the SDK when the API call is made
		// Use this member when unit testing your code with the SDK to mock out the
		// EventStream Reader.
		//
		// Must not be nil.
		Reader {{ $.ShapeName }}Reader

	{{ end -}}

	{{- if $.EventStreamAPI.Outbound }}
		// Writer is the EventStream reader for the {{ $.EventStreamAPI.Inbound.Name }}
		// events. This value is automatically set by the SDK when the API call is made
		// Use this member when unit testing your code with the SDK to mock out the
		// EventStream Writer.
		//
		// Must not be nil.
		Writer *{{ $.ShapeName }}Writer

	{{ end -}}

	// StreamCloser is the io.Closer for the EventStream connection. For HTTP
	// EventStream this is the response Body. The stream will be closed when
	// the Close method of the EventStream is called.
	StreamCloser io.Closer
}

// Close closes the EventStream. This will also cause the Events channel to be
// closed. You can use the closing of the Events channel to terminate your
// application's read from the API's EventStream.
{{- if $.EventStreamAPI.Inbound }}
//
// Will close the underlying EventStream reader. For EventStream over HTTP
// connection this will also close the HTTP connection.
{{ end -}}
//
// Close must be called when done using the EventStream API. Not calling Close
// may result in resource leaks.
func (es *{{ $.ShapeName }}) Close() (err error) {
	{{- if $.EventStreamAPI.Inbound }}
		es.Reader.Close()
	{{ end -}}
	{{- if $.EventStreamAPI.Outbound }}
		es.Writer.Close()
	{{ end -}}

	return es.Err()
}

// Err returns any error that occurred while reading EventStream Events from
// the service API's response. Returns nil if there were no errors.
func (es *{{ $.ShapeName }}) Err() error {
	{{- if $.EventStreamAPI.Outbound }}
		if err := es.Writer.Err(); err != nil {
			return err
		}
	{{ end -}}

	{{- if $.EventStreamAPI.Inbound }}
		if err := es.Reader.Err(); err != nil {
			return err
		}
	{{ end -}}

	es.StreamCloser.Close()

	return nil
}

{{ if $.EventStreamAPI.Inbound }}
	// Events returns a channel to read EventStream Events from the
	// {{ $.EventStreamAPI.Operation.ExportedName }} API.
	//
	// These events are:
	// {{ range $_, $event := $.EventStreamAPI.Inbound.Events }}
	//     * {{ $event.Shape.ShapeName }}
	{{- end }}
	func (es *{{ $.ShapeName }}) Events() <-chan {{ $.EventStreamAPI.Inbound.Name }}Event {
		return es.Reader.Events()
	}

	{{ template "eventStreamAPIReaderTmpl" $ }}
{{ end }}

{{ if $.EventStreamAPI.Outbound }}
	// TODO writer helper method.
{{ end }}

`

var eventStreamAPIReaderTmpl = template.Must(template.New("eventStreamAPIReaderTmpl").
	Funcs(template.FuncMap{}).
	Parse(`
// {{ $.EventStreamAPI.Inbound.Name }}Event groups together all EventStream
// events read from the {{ $.EventStreamAPI.Operation.ExportedName }} API.
//
// These events are:
// {{ range $_, $event := $.EventStreamAPI.Inbound.Events }}
//     * {{ $event.Shape.ShapeName }}
{{- end }}
type {{ $.EventStreamAPI.Inbound.Name }}Event interface {
	event{{ $.EventStreamAPI.Inbound.Name }}()
}

// {{ $.ShapeName }}Reader provides the interface for reading EventStream
// Events from the {{ $.EventStreamAPI.Operation.ExportedName }} API. The
// default implementation for this interface will be {{ $.ShapeName }}.
//
// The reader's Close method must allow multiple concurrent calls.
//
// These events are:
// {{ range $_, $event := $.EventStreamAPI.Inbound.Events }}
//     * {{ $event.Shape.ShapeName }}
{{- end }}
type {{ $.ShapeName }}Reader interface {
	// Returns a channel of events as they are read from the event stream.
	Events() <-chan {{ $.EventStreamAPI.Inbound.Name }}Event

	// Close will close the underlying event stream reader. For event stream over
	// HTTP this will also close the HTTP connection.
	Close() error

	// Returns any error that has occurred while reading from the event stream.
	Err() error
}

type read{{ $.ShapeName }} struct {
	eventReader *eventstreamapi.EventReader
	stream chan {{ $.EventStreamAPI.Inbound.Name }}Event
	errVal atomic.Value

	done      chan struct{}
	closeOnce sync.Once

	{{ if eq $.API.Metadata.Protocol "json" -}}
		initResp eventstreamapi.Unmarshaler
	{{ end -}}
}

func newRead{{ $.ShapeName }}(
	reader io.ReadCloser,
	unmarshalers request.HandlerList,
	logger aws.Logger,
	logLevel aws.LogLevelType,
	{{ if eq $.API.Metadata.Protocol "json" -}}
		initResp eventstreamapi.Unmarshaler,
	{{ end -}}
) *read{{ $.ShapeName }} {
	r := &read{{ $.ShapeName }}{
		stream: make(chan {{ $.EventStreamAPI.Inbound.Name }}Event),
		done: make(chan struct{}),
		{{ if eq $.API.Metadata.Protocol "json" -}}
			initResp: initResp,
		{{ end -}}
	}

	r.eventReader = eventstreamapi.NewEventReader(
		reader,
		protocol.HandlerPayloadUnmarshal{
			Unmarshalers: unmarshalers,
		},
		r.unmarshalerForEventType,
	)
	r.eventReader.UseLogger(logger, logLevel)

	return r
}

// Close will close the underlying event stream reader. For EventStream over
// HTTP this will also close the HTTP connection.
func (r *read{{ $.ShapeName }}) Close() error {
	r.closeOnce.Do(r.safeClose)

	return r.Err()
}

func (r *read{{ $.ShapeName }}) safeClose() {
	close(r.done)
	err := r.eventReader.Close()
	if err != nil {
		r.errVal.Store(err)
	}
}

func (r *read{{ $.ShapeName }}) Err() error {
	if v := r.errVal.Load(); v != nil {
		return v.(error)
	}

	return nil
}

func (r *read{{ $.ShapeName }}) Events() <-chan {{ $.EventStreamAPI.Inbound.Name }}Event {
	return r.stream
}

func (r *read{{ $.ShapeName }}) readEventStream() {
	defer close(r.stream)

	for {
		event, err := r.eventReader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				return
			}
			select {
			case <-r.done:
				// If closed already ignore the error
				return
			default:
			}
			r.errVal.Store(err)
			return
		}

		select {
		case r.stream <- event.({{ $.EventStreamAPI.Inbound.Name }}Event):
		case <-r.done:
			return
		}
	}
}

func (r *read{{ $.ShapeName }}) unmarshalerForEventType(
	eventType string,
) (eventstreamapi.Unmarshaler, error) {
	switch eventType {
		{{- if eq $.API.Metadata.Protocol "json" }}
			case "initial-response":
				return r.initResp, nil
		{{ end -}}
		{{- range $_, $event := $.EventStreamAPI.Inbound.Events }}
			case {{ printf "%q" $event.Name }}:
				return &{{ $event.Shape.ShapeName }}{}, nil
		{{ end -}}
		{{- range $_, $event := $.EventStreamAPI.Inbound.Exceptions }}
			case {{ printf "%q" $event.Name }}:
				return &{{ $event.Shape.ShapeName }}{}, nil
		{{ end -}}
	default:
		return nil, awserr.New(
			request.ErrCodeSerialization,
			fmt.Sprintf("unknown event type name, %s, for {{ $.ShapeName }}", eventType),
			nil,
		)
	}
}
`))

// Template for the EventStream API Output shape that contains the EventStream
// member.
//
// Executed in the context of a Shape.
var eventStreamAPILoopMethodTmpl = template.Must(
	template.New("eventStreamAPILoopMethodTmpl").Parse(`
func (s *{{ $.ShapeName }}) runEventStreamLoop(r *request.Request) {
	if r.Error != nil {
		return
	}

	{{- $esMemberRef := index $.MemberRefs $.EventStreamsMemberName }}
	{{- if $esMemberRef.Shape.EventStreamAPI.Inbound }}
		reader := newRead{{ $esMemberRef.ShapeName }}(
			r.HTTPResponse.Body,
			r.Handlers.UnmarshalStream,
			r.Config.Logger,
			r.Config.LogLevel.Value(),
			{{ if eq $.API.Metadata.Protocol "json" -}}
				s,
			{{ end -}}
		)
		go reader.readEventStream()

		eventStream := &{{ $esMemberRef.ShapeName }} {
			StreamCloser: r.HTTPResponse.Body,
			Reader: reader,
		}
	{{ end -}}

	s.{{ $.EventStreamsMemberName }} = eventStream
}

{{ if eq $.API.Metadata.Protocol "json" -}}
	func (s *{{ $.ShapeName }}) unmarshalInitialResponse(r *request.Request) {
		// Wait for the initial response event, which must be the first event to be
		// received from the API.
		select {
		case event, ok := <-s.EventStream.Events():
			if !ok {
				return
			}
			es := s.EventStream
			v, ok := event.(*{{ $.ShapeName }})
			if !ok || v == nil {
				r.Error = awserr.New(
					request.ErrCodeSerialization,
					fmt.Sprintf("invalid event, %T, expect *SubscribeToShardOutput, %v", event, v),
					nil,
				)
				return
			}
			*s = *v
			s.EventStream = es
		}
	}
{{ end -}}
`))

// EventStreamHeaderTypeMap provides the mapping of a EventStream Header's
// Value type to the shape reference's member type.
type EventStreamHeaderTypeMap struct {
	Header string
	Member string
}

var eventStreamEventShapeTmplFuncs = template.FuncMap{
	"EventStreamHeaderTypeMap": func(ref *ShapeRef) EventStreamHeaderTypeMap {
		switch ref.Shape.Type {
		case "boolean":
			return EventStreamHeaderTypeMap{Header: "bool", Member: "bool"}
		case "byte":
			return EventStreamHeaderTypeMap{Header: "int8", Member: "int64"}
		case "short":
			return EventStreamHeaderTypeMap{Header: "int16", Member: "int64"}
		case "integer":
			return EventStreamHeaderTypeMap{Header: "int32", Member: "int64"}
		case "long":
			return EventStreamHeaderTypeMap{Header: "int64", Member: "int64"}
		case "timestamp":
			return EventStreamHeaderTypeMap{Header: "time.Time", Member: "time.Time"}
		case "blob":
			return EventStreamHeaderTypeMap{Header: "[]byte", Member: "[]byte"}
		case "string":
			return EventStreamHeaderTypeMap{Header: "string", Member: "string"}
		// TODO case "uuid"  what is modeled type
		default:
			panic("unsupported EventStream header type, " + ref.Shape.Type)
		}
	},
	"HasNonBlobPayloadMembers": eventHasNonBlobPayloadMembers,
}

// Returns if the event has any members which are not the event's blob payload,
// nor a header.
func eventHasNonBlobPayloadMembers(s *Shape) bool {
	num := len(s.MemberRefs)
	for _, ref := range s.MemberRefs {
		if ref.IsEventHeader || (ref.IsEventPayload && (ref.Shape.Type == "blob" || ref.Shape.Type == "string")) {
			num--
		}
	}
	return num > 0
}

// Template for an EventStream Event shape. This is a normal API shape that is
// decorated as an EventStream Event.
//
// Executed in the context of a Shape.
var eventStreamEventShapeTmpl = template.Must(template.New("eventStreamEventShapeTmpl").
	Funcs(eventStreamEventShapeTmplFuncs).Parse(`
{{ range $_, $eventstream := $.EventFor }}
	// The {{ $.ShapeName }} is and event in the {{ $eventstream.Name }} group of events.
	func (s *{{ $.ShapeName }}) event{{ $eventstream.Name }}() {}
{{ end }}

// UnmarshalEvent unmarshals the EventStream Message into the {{ $.ShapeName }} value.
// This method is only used internally within the SDK's EventStream handling.
func (s *{{ $.ShapeName }}) UnmarshalEvent(
	payloadUnmarshaler protocol.PayloadUnmarshaler,
	msg eventstream.Message,
) error {
	{{- range $memName, $memRef := $.MemberRefs }}
		{{- if $memRef.IsEventHeader }}
			if hv := msg.Headers.Get("{{ $memName }}"); hv != nil {
				{{ $types := EventStreamHeaderTypeMap $memRef -}}
				v := hv.Get().({{ $types.Header }})
				{{- if ne $types.Header $types.Member }}
					m := {{ $types.Member }}(v)
					s.{{ $memName }} = {{ if $memRef.UseIndirection }}&{{ end }}m
				{{- else }}
					s.{{ $memName }} = {{ if $memRef.UseIndirection }}&{{ end }}v
				{{- end }}
			}
		{{- else if (and ($memRef.IsEventPayload) (eq $memRef.Shape.Type "blob")) }}
			s.{{ $memName }} = make([]byte, len(msg.Payload))
			copy(s.{{ $memName }}, msg.Payload)
		{{- else if (and ($memRef.IsEventPayload) (eq $memRef.Shape.Type "string")) }}
			s.{{ $memName }} = aws.String(string(msg.Payload))
		{{- end }}
	{{- end }}
	{{- if HasNonBlobPayloadMembers $ }}
		if err := payloadUnmarshaler.UnmarshalPayload(
			bytes.NewReader(msg.Payload), s,
		); err != nil {
			return err
		}
	{{- end }}
	return nil
}
`))

var eventStreamExceptionEventShapeTmpl = template.Must(
	template.New("eventStreamExceptionEventShapeTmpl").Parse(`
// Code returns the exception type name.
func (s {{ $.ShapeName }}) Code() string {
	{{- if $.ErrorInfo.Code }}
		return "{{ $.ErrorInfo.Code }}"
	{{- else }}
		return "{{ $.ShapeName }}"
	{{ end -}}
}

// Message returns the exception's message.
func (s {{ $.ShapeName }}) Message() string {
	{{- if index $.MemberRefs "Message_" }}
		return *s.Message_
	{{- else }}
		return ""
	{{ end -}}
}

// OrigErr always returns nil, satisfies awserr.Error interface.
func (s {{ $.ShapeName }}) OrigErr() error {
	return nil
}

func (s {{ $.ShapeName }}) Error() string {
	return fmt.Sprintf("%s: %s", s.Code(), s.Message())
}
`))

// APIEventStreamTestGoCode generates Go code for EventStream operation tests.
func (a *API) APIEventStreamTestGoCode() string {
	var buf bytes.Buffer

	a.resetImports()
	a.AddImport("bytes")
	a.AddImport("io/ioutil")
	a.AddImport("net/http")
	a.AddImport("reflect")
	a.AddImport("testing")
	a.AddImport("time")
	a.AddSDKImport("aws")
	a.AddSDKImport("aws/corehandlers")
	a.AddSDKImport("aws/request")
	a.AddSDKImport("aws/awserr")
	a.AddSDKImport("awstesting/unit")
	a.AddSDKImport("private/protocol")
	a.AddSDKImport("private/protocol/", a.ProtocolPackage())
	a.AddSDKImport("private/protocol/eventstream")
	a.AddSDKImport("private/protocol/eventstream/eventstreamapi")
	a.AddSDKImport("private/protocol/eventstream/eventstreamtest")

	unused := `
	var _ time.Time
	var _ awserr.Error
	`

	if err := eventStreamTestTmpl.Execute(&buf, a); err != nil {
		panic(err)
	}

	return a.importsGoCode() + unused + strings.TrimSpace(buf.String())
}

func valueForType(s *Shape, visited []string) string {
	for _, v := range visited {
		if v == s.ShapeName {
			return "nil"
		}
	}

	visited = append(visited, s.ShapeName)

	switch s.Type {
	case "blob":
		return `[]byte("blob value goes here")`
	case "string":
		return `aws.String("string value goes here")`
	case "boolean":
		return `aws.Bool(true)`
	case "byte":
		return `aws.Int64(1)`
	case "short":
		return `aws.Int64(12)`
	case "integer":
		return `aws.Int64(123)`
	case "long":
		return `aws.Int64(1234)`
	case "float":
		return `aws.Float64(123.4)`
	case "double":
		return `aws.Float64(123.45)`
	case "timestamp":
		return `aws.Time(time.Unix(1396594860, 0).UTC())`
	case "structure":
		w := bytes.NewBuffer(nil)
		fmt.Fprintf(w, "&%s{\n", s.ShapeName)
		for _, refName := range s.MemberNames() {
			fmt.Fprintf(w, "%s: %s,\n", refName, valueForType(s.MemberRefs[refName].Shape, visited))
		}
		fmt.Fprintf(w, "}")
		return w.String()
	case "list":
		w := bytes.NewBuffer(nil)
		fmt.Fprintf(w, "%s{\n", s.GoType())
		for i := 0; i < 3; i++ {
			fmt.Fprintf(w, "%s,\n", valueForType(s.MemberRef.Shape, visited))
		}
		fmt.Fprintf(w, "}")
		return w.String()

	case "map":
		w := bytes.NewBuffer(nil)
		fmt.Fprintf(w, "%s{\n", s.GoType())
		for _, k := range []string{"a", "b", "c"} {
			fmt.Fprintf(w, "%q: %s,\n", k, valueForType(s.ValueRef.Shape, visited))
		}
		fmt.Fprintf(w, "}")
		return w.String()

	default:
		panic(fmt.Sprintf("valueForType does not support %s, %s", s.ShapeName, s.Type))
	}
}

func setEventHeaderValueForType(s *Shape, memVar string) string {
	switch s.Type {
	case "blob":
		return fmt.Sprintf("eventstream.BytesValue(%s)", memVar)
	case "string":
		return fmt.Sprintf("eventstream.StringValue(*%s)", memVar)
	case "boolean":
		return fmt.Sprintf("eventstream.BoolValue(*%s)", memVar)
	case "byte":
		return fmt.Sprintf("eventstream.Int8Value(int8(*%s))", memVar)
	case "short":
		return fmt.Sprintf("eventstream.Int16Value(int16(*%s))", memVar)
	case "integer":
		return fmt.Sprintf("eventstream.Int32Value(int32(*%s))", memVar)
	case "long":
		return fmt.Sprintf("eventstream.Int64Value(*%s)", memVar)
	case "float":
		return fmt.Sprintf("eventstream.Float32Value(float32(*%s))", memVar)
	case "double":
		return fmt.Sprintf("eventstream.Float64Value(*%s)", memVar)
	case "timestamp":
		return fmt.Sprintf("eventstream.TimestampValue(*%s)", memVar)
	default:
		panic(fmt.Sprintf("value type %s not supported for event headers, %s", s.Type, s.ShapeName))
	}
}

func templateMap(args ...interface{}) map[string]interface{} {
	if len(args)%2 != 0 {
		panic(fmt.Sprintf("invalid map call, non-even args %v", args))
	}

	m := map[string]interface{}{}
	for i := 0; i < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			panic(fmt.Sprintf("invalid map call, arg is not string, %T, %v", args[i], args[i]))
		}
		m[k] = args[i+1]
	}

	return m
}

var eventStreamTestTmpl = template.Must(
	template.New("eventStreamTestTmpl").Funcs(template.FuncMap{
		"ValueForType":               valueForType,
		"HasNonBlobPayloadMembers":   eventHasNonBlobPayloadMembers,
		"SetEventHeaderValueForType": setEventHeaderValueForType,
		"Map":                        templateMap,
		"OptionalAddInt": func(do bool, a, b int) int {
			if !do {
				return a
			}
			return a + b
		},
		"HasNonEventStreamMember": func(s *Shape) bool {
			for _, ref := range s.MemberRefs {
				if !ref.Shape.IsEventStream {
					return true
				}
			}
			return false
		},
	}).Parse(`
{{ range $opName, $op := $.Operations }}
	{{ if $op.EventStreamAPI }}
		{{ if $op.EventStreamAPI.Inbound }}
			{{ template "event stream inbound tests" $op.EventStreamAPI }}
		{{ end }}
	{{ end }}
{{ end }}

type loopReader struct {
	source *bytes.Reader
}

func (c *loopReader) Read(p []byte) (int, error) {
	if c.source.Len() == 0 {
		c.source.Seek(0, 0)
	}

	return c.source.Read(p)
}

{{ define "event stream inbound tests" }}
	func Test{{ $.Operation.ExportedName }}_Read(t *testing.T) {
		expectEvents, eventMsgs := mock{{ $.Operation.ExportedName }}ReadEvents()
		sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
			eventstreamtest.ServeEventStream{
				T:      t,
				Events: eventMsgs,
			},
			true,
		)
		if err != nil {
			t.Fatalf("expect no error, %v", err)
		}
		defer cleanupFn()

		svc := New(sess)
		resp, err := svc.{{ $.Operation.ExportedName }}(nil)
		if err != nil {
			t.Fatalf("expect no error got, %v", err)
		}
		defer resp.EventStream.Close()

		{{- if eq $.Operation.API.Metadata.Protocol "json" }}
			{{- if HasNonEventStreamMember $.Operation.OutputRef.Shape }}
				expectResp := expectEvents[0].(*{{ $.Operation.OutputRef.Shape.ShapeName }})
				{{- range $name, $ref := $.Operation.OutputRef.Shape.MemberRefs }}
					{{- if not $ref.Shape.IsEventStream }}
						if e, a := expectResp.{{ $name }}, resp.{{ $name }}; !reflect.DeepEqual(e,a) {
							t.Errorf("expect %v, got %v", e, a)
						}
					{{- end }}
				{{- end }}
			{{- end }}
			// Trim off response output type pseudo event so only event messages remain.
			expectEvents = expectEvents[1:]
		{{ end }}

		var i int
		for event := range resp.EventStream.Events() {
			if event == nil {
				t.Errorf("%d, expect event, got nil", i)
			}
			if e, a := expectEvents[i], event; !reflect.DeepEqual(e, a) {
				t.Errorf("%d, expect %T %v, got %T %v", i, e, e, a, a)
			}
			i++
		}

		if err := resp.EventStream.Err(); err != nil {
			t.Errorf("expect no error, %v", err)
		}
	}

	func Test{{ $.Operation.ExportedName }}_ReadClose(t *testing.T) {
		_, eventMsgs := mock{{ $.Operation.ExportedName }}ReadEvents()
		sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
			eventstreamtest.ServeEventStream{
				T:      t,
				Events: eventMsgs,
			},
			true,
		)
		if err != nil {
			t.Fatalf("expect no error, %v", err)
		}
		defer cleanupFn()

		svc := New(sess)
		resp, err := svc.{{ $.Operation.ExportedName }}(nil)
		if err != nil {
			t.Fatalf("expect no error got, %v", err)
		}

		resp.EventStream.Close()
		<-resp.EventStream.Events()

		if err := resp.EventStream.Err(); err != nil {
			t.Errorf("expect no error, %v", err)
		}
	}

	func Benchmark{{ $.Operation.ExportedName }}_Read(b *testing.B) {
		_, eventMsgs := mock{{ $.Operation.ExportedName }}ReadEvents()
		var buf bytes.Buffer
		encoder := eventstream.NewEncoder(&buf)
		for _, msg := range eventMsgs {
			if err := encoder.Encode(msg); err != nil {
				b.Fatalf("failed to encode message, %v", err)
			}
		}
		stream := &loopReader{source: bytes.NewReader(buf.Bytes())}

		sess := unit.Session
		svc := New(sess, &aws.Config{
			Endpoint:               aws.String("https://example.com"),
			DisableParamValidation: aws.Bool(true),
		})
		svc.Handlers.Send.Swap(corehandlers.SendHandler.Name,
			request.NamedHandler{Name: "mockSend",
				Fn: func(r *request.Request) {
					r.HTTPResponse = &http.Response{
						Status:     "200 OK",
						StatusCode: 200,
						Header:     http.Header{},
						Body:       ioutil.NopCloser(stream),
					}
				},
			},
		)

		resp, err := svc.{{ $.Operation.ExportedName }}(nil)
		if err != nil {
			b.Fatalf("failed to create request, %v", err)
		}
		defer resp.EventStream.Close()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			if err = resp.EventStream.Err(); err != nil {
				b.Fatalf("expect no error, got %v", err)
			}
			event := <-resp.EventStream.Events()
			if event == nil {
				b.Fatalf("expect event, got nil, %v, %d", resp.EventStream.Err(), i)
			}
		}
	}

	func mock{{ $.Operation.ExportedName }}ReadEvents() (
		[]{{ $.Inbound.Name }}Event,
		[]eventstream.Message,
	) {
		expectEvents := []{{ $.Inbound.Name }}Event {
			{{- if eq $.Operation.API.Metadata.Protocol "json" }}
				{{- template "set event type" $.Operation.OutputRef.Shape }}
			{{- end }}
			{{- range $_, $event := $.Inbound.Events }}
				{{- template "set event type" $event.Shape }}
			{{- end }}
		}

		var marshalers request.HandlerList
		marshalers.PushBackNamed({{ $.API.ProtocolPackage }}.BuildHandler)
		payloadMarshaler := protocol.HandlerPayloadMarshal{
			Marshalers: marshalers,
		}
		_ = payloadMarshaler

		eventMsgs := []eventstream.Message{
			{{- if eq $.Operation.API.Metadata.Protocol "json" }}
				{{- template "set event message" Map "idx" 0 "parentShape" $.Operation.OutputRef.Shape "eventName" "initial-response" }}
			{{- end }}
			{{- range $idx, $event := $.Inbound.Events }}
				{{- $offsetIdx := OptionalAddInt (eq $.Operation.API.Metadata.Protocol "json") $idx 1 }}
				{{- template "set event message" Map "idx" $offsetIdx "parentShape" $event.Shape "eventName" $event.Name }}
			{{- end }}
		}

		return expectEvents, eventMsgs
	}

	{{- if $.Inbound.Exceptions }}
		func Test{{ $.Operation.ExportedName }}_ReadException(t *testing.T) {
			expectEvents := []{{ $.Inbound.Name }}Event {
				{{- if eq $.Operation.API.Metadata.Protocol "json" }}
					{{- template "set event type" $.Operation.OutputRef.Shape }}
				{{- end }}

				{{- $exception := index $.Inbound.Exceptions 0 }}
				{{- template "set event type" $exception.Shape }}
			}

			var marshalers request.HandlerList
			marshalers.PushBackNamed({{ $.API.ProtocolPackage }}.BuildHandler)
			payloadMarshaler := protocol.HandlerPayloadMarshal{
				Marshalers: marshalers,
			}

			eventMsgs := []eventstream.Message{
				{{- if eq $.Operation.API.Metadata.Protocol "json" }}
					{{- template "set event message" Map "idx" 0 "parentShape" $.Operation.OutputRef.Shape "eventName" "initial-response" }}
				{{- end }}

				{{- $offsetIdx := OptionalAddInt (eq $.Operation.API.Metadata.Protocol "json") 0 1 }}
				{{- $exception := index $.Inbound.Exceptions 0 }}
				{{- template "set event message" Map "idx" $offsetIdx "parentShape" $exception.Shape "eventName" $exception.Name }}
			}

			sess, cleanupFn, err := eventstreamtest.SetupEventStreamSession(t,
				eventstreamtest.ServeEventStream{
					T:      t,
					Events: eventMsgs,
				},
				true,
			)
			if err != nil {
				t.Fatalf("expect no error, %v", err)
			}
			defer cleanupFn()

			svc := New(sess)
			resp, err := svc.{{ $.Operation.ExportedName }}(nil)
			if err != nil {
				t.Fatalf("expect no error got, %v", err)
			}

			defer resp.EventStream.Close()

			<-resp.EventStream.Events()

			err = resp.EventStream.Err()
			if err == nil {
				t.Fatalf("expect err, got none")
			}

			expectErr := {{ ValueForType $exception.Shape nil }}
			aerr, ok := err.(awserr.Error)
			if !ok {
				t.Errorf("expect exception, got %T, %#v", err, err)
			}
			if e, a := expectErr.Code(), aerr.Code(); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}
			if e, a := expectErr.Message(), aerr.Message(); e != a {
				t.Errorf("expect %v, got %v", e, a)
			}

			if e, a := expectErr, aerr; !reflect.DeepEqual(e, a) {
				t.Errorf("expect %#v, got %#v", e, a)
			}
		}

		{{- range $_, $exception := $.Inbound.Exceptions }}
			var _ awserr.Error = (*{{ $exception.Shape.ShapeName }})(nil)
		{{- end }}

	{{ end }}
{{ end }}

{{/* Params: *Shape */}}
{{ define "set event type" }}
	&{{ $.ShapeName }}{
		{{- range $memName, $memRef := $.MemberRefs }}
			{{- if not $memRef.Shape.IsEventStream }}
				{{ $memName }}: {{ ValueForType $memRef.Shape nil }},
			{{- end }}
		{{- end }}
	},
{{- end }}

{{/* Params: idx:int, parentShape:*Shape, eventName:string */}}
{{ define "set event message" }}
	{
		Headers: eventstream.Headers{
			{{- if $.parentShape.Exception }}
				eventstreamtest.EventExceptionTypeHeader,
				{
					Name:  eventstreamapi.ExceptionTypeHeader,
					Value: eventstream.StringValue("{{ $.eventName }}"),
				},
			{{- else }}
				eventstreamtest.EventMessageTypeHeader,
				{
					Name:  eventstreamapi.EventTypeHeader,
					Value: eventstream.StringValue("{{ $.eventName }}"),
				},
			{{- end }}
			{{- range $memName, $memRef := $.parentShape.MemberRefs }}
				{{- template "set event message header" Map "idx" $.idx "parentShape" $.parentShape "memName" $memName "memRef" $memRef }}
			{{- end }}
		},
		{{- template "set event message payload" Map "idx" $.idx "parentShape" $.parentShape }}
	},
{{- end }}

{{/* Params: idx:int, parentShape:*Shape, memName:string, memRef:*ShapeRef */}}
{{ define "set event message header" }}
	{{- if $.memRef.IsEventHeader }}
		{
			Name: "{{ $.memName }}",
			{{- $shapeValueVar := printf "expectEvents[%d].(%s).%s" $.idx $.parentShape.GoType $.memName }}
			Value: {{ SetEventHeaderValueForType $.memRef.Shape $shapeValueVar }},
		},
	{{- end }}
{{- end }}

{{/* Params: idx:int, parentShape:*Shape, memName:string, memRef:*ShapeRef */}}
{{ define "set event message payload" }}
	{{- $payloadMemName := $.parentShape.PayloadRefName }}
	{{- if HasNonBlobPayloadMembers $.parentShape }}
		Payload: eventstreamtest.MarshalEventPayload(payloadMarshaler, expectEvents[{{ $.idx }}]),
	{{- else if $payloadMemName }}
		{{- $shapeType := (index $.parentShape.MemberRefs $payloadMemName).Shape.Type }}
		{{- if eq $shapeType "blob" }}
			Payload: expectEvents[{{ $.idx }}].({{ $.parentShape.GoType }}).{{ $payloadMemName }},
		{{- else if eq $shapeType "string" }}
			Payload: []byte(*expectEvents[{{ $.idx }}].({{ $.parentShape.GoType }}).{{ $payloadMemName }}),
		{{- end }}
	{{- end }}
{{- end }}
`))
