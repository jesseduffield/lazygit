package eventstreamtest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/private/protocol/eventstream"
	"github.com/aws/aws-sdk-go/private/protocol/eventstream/eventstreamapi"
)

// ServeEventStream provides serving EventStream messages from a HTTP server to
// the client. The events are sent sequentially to the client without delay.
type ServeEventStream struct {
	T      *testing.T
	Events []eventstream.Message
}

func (s ServeEventStream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	encoder := eventstream.NewEncoder(flushWriter{w})

	for _, event := range s.Events {
		encoder.Encode(event)
	}
}

// SetupEventStreamSession creates a HTTP server SDK session for communicating
// with that server to be used for EventStream APIs. If HTTP/2 is enabled the
// server/client will only attempt to use HTTP/2.
func SetupEventStreamSession(
	t *testing.T, handler http.Handler, h2 bool,
) (sess *session.Session, cleanupFn func(), err error) {
	server := httptest.NewUnstartedServer(handler)

	client := setupServer(server, h2)

	cleanupFn = func() {
		server.Close()
	}

	sess, err = session.NewSession(unit.Session.Config, &aws.Config{
		Endpoint:               &server.URL,
		DisableParamValidation: aws.Bool(true),
		HTTPClient:             client,
		//		LogLevel:               aws.LogLevel(aws.LogDebugWithEventStreamBody),
	})
	if err != nil {
		return nil, nil, err
	}

	return sess, cleanupFn, nil
}

type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}

// MarshalEventPayload marshals a SDK API shape into its associated wire
// protocol payload.
func MarshalEventPayload(
	payloadMarshaler protocol.PayloadMarshaler,
	v interface{},
) []byte {
	var w bytes.Buffer
	err := payloadMarshaler.MarshalPayload(&w, v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal event %T, %v", v, v))
	}

	return w.Bytes()
}

// EventMessageTypeHeader is an event message type header for specifying an
// event is an message type.
var EventMessageTypeHeader = eventstream.Header{
	Name:  eventstreamapi.MessageTypeHeader,
	Value: eventstream.StringValue(eventstreamapi.EventMessageType),
}

// EventExceptionTypeHeader is an event exception type header for specifying an
// event is an exeption type.
var EventExceptionTypeHeader = eventstream.Header{
	Name:  eventstreamapi.MessageTypeHeader,
	Value: eventstream.StringValue(eventstreamapi.ExceptionMessageType),
}
