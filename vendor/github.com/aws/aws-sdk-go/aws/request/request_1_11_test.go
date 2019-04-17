// +build go1.11

package request_test

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/client/metadata"
	"github.com/aws/aws-sdk-go/aws/corehandlers"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/awstesting/unit"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
)

func TestSerializationErrConnectionReset_read(t *testing.T) {
	count := 0
	handlers := request.Handlers{}
	handlers.Send.PushBack(func(r *request.Request) {
		count++
		r.HTTPResponse = &http.Response{}
		r.HTTPResponse.Body = &connResetCloser{
			Err: errReadConnectionResetStub,
		}
	})

	handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	handlers.Build.PushBackNamed(jsonrpc.BuildHandler)
	handlers.Unmarshal.PushBackNamed(jsonrpc.UnmarshalHandler)
	handlers.UnmarshalMeta.PushBackNamed(jsonrpc.UnmarshalMetaHandler)
	handlers.UnmarshalError.PushBackNamed(jsonrpc.UnmarshalErrorHandler)
	handlers.AfterRetry.PushBackNamed(corehandlers.AfterRetryHandler)

	op := &request.Operation{
		Name:       "op",
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	meta := metadata.ClientInfo{
		ServiceName:   "fooService",
		SigningName:   "foo",
		SigningRegion: "foo",
		Endpoint:      "localhost",
		APIVersion:    "2001-01-01",
		JSONVersion:   "1.1",
		TargetPrefix:  "Foo",
	}
	cfg := unit.Session.Config.Copy()
	cfg.MaxRetries = aws.Int(5)

	req := request.New(
		*cfg,
		meta,
		handlers,
		client.DefaultRetryer{NumMaxRetries: 5},
		op,
		&struct {
		}{},
		&struct {
		}{},
	)

	osErr := errReadConnectionResetStub
	req.ApplyOptions(request.WithResponseReadTimeout(time.Second))
	err := req.Send()
	if err == nil {
		t.Error("Expected rror 'SerializationError', but received nil")
	}
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "SerializationError" {
		t.Errorf("Expected 'SerializationError', but received %q", aerr.Code())
	} else if !ok {
		t.Errorf("Expected 'awserr.Error', but received %v", reflect.TypeOf(err))
	} else if aerr.OrigErr().Error() != osErr.Error() {
		t.Errorf("Expected %q, but received %q", osErr.Error(), aerr.OrigErr().Error())
	}

	if count != 1 {
		t.Errorf("Expected '1', but received %d", count)
	}
}
