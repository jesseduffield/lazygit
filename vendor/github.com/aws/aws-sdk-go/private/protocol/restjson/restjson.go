// Package restjson provides RESTful JSON serialization of AWS
// requests and responses.
package restjson

//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/input/rest-json.json build_test.go
//go:generate go run -tags codegen ../../../models/protocol_tests/generate.go ../../../models/protocol_tests/output/rest-json.json unmarshal_test.go

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
)

// BuildHandler is a named request handler for building restjson protocol requests
var BuildHandler = request.NamedHandler{Name: "awssdk.restjson.Build", Fn: Build}

// UnmarshalHandler is a named request handler for unmarshaling restjson protocol requests
var UnmarshalHandler = request.NamedHandler{Name: "awssdk.restjson.Unmarshal", Fn: Unmarshal}

// UnmarshalMetaHandler is a named request handler for unmarshaling restjson protocol request metadata
var UnmarshalMetaHandler = request.NamedHandler{Name: "awssdk.restjson.UnmarshalMeta", Fn: UnmarshalMeta}

// UnmarshalErrorHandler is a named request handler for unmarshaling restjson protocol request errors
var UnmarshalErrorHandler = request.NamedHandler{Name: "awssdk.restjson.UnmarshalError", Fn: UnmarshalError}

// Build builds a request for the REST JSON protocol.
func Build(r *request.Request) {
	rest.Build(r)

	if t := rest.PayloadType(r.Params); t == "structure" || t == "" {
		if v := r.HTTPRequest.Header.Get("Content-Type"); len(v) == 0 {
			r.HTTPRequest.Header.Set("Content-Type", "application/json")
		}
		jsonrpc.Build(r)
	}
}

// Unmarshal unmarshals a response body for the REST JSON protocol.
func Unmarshal(r *request.Request) {
	if t := rest.PayloadType(r.Data); t == "structure" || t == "" {
		jsonrpc.Unmarshal(r)
	} else {
		rest.Unmarshal(r)
	}
}

// UnmarshalMeta unmarshals response headers for the REST JSON protocol.
func UnmarshalMeta(r *request.Request) {
	rest.UnmarshalMeta(r)
}

// UnmarshalError unmarshals a response error for the REST JSON protocol.
func UnmarshalError(r *request.Request) {
	defer r.HTTPResponse.Body.Close()

	var jsonErr jsonErrorResponse
	err := json.NewDecoder(r.HTTPResponse.Body).Decode(&jsonErr)
	if err == io.EOF {
		r.Error = awserr.NewRequestFailure(
			awserr.New("SerializationError", r.HTTPResponse.Status, nil),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	} else if err != nil {
		r.Error = awserr.NewRequestFailure(
			awserr.New("SerializationError", "failed decoding REST JSON error response", err),
			r.HTTPResponse.StatusCode,
			r.RequestID,
		)
		return
	}

	code := r.HTTPResponse.Header.Get("X-Amzn-Errortype")
	if code == "" {
		code = jsonErr.Code
	}

	code = strings.SplitN(code, ":", 2)[0]
	r.Error = awserr.NewRequestFailure(
		awserr.New(code, jsonErr.Message, nil),
		r.HTTPResponse.StatusCode,
		r.RequestID,
	)
}

type jsonErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
