// +build example

package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/private/protocol/jsonrpc"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const TxAwareErrorUnmarshallerName = "awssdk.jsonrpc.TxAwareErrorUnmarshaller"

// New creates a new instance of the DynamoDB client with a session.
// The client's behaviour is same as what is returned by dynamodb.New(), except for richer error reasons.
func NewTxErrorAwareDynamoDBClient(p client.ConfigProvider, cfgs ...*aws.Config) *dynamodb.DynamoDB {
	c := dynamodb.New(p, cfgs...)
	// NOTE: Ignore if swap failed. Returning nil might fail app startup which is worse than inadequate error details.
	c.Handlers.UnmarshalError.Swap(jsonrpc.UnmarshalErrorHandler.Name, request.NamedHandler{
		Name: TxAwareErrorUnmarshallerName,
		Fn:   TxAwareUnmarshalError,
	})
	return c
}

// A RequestFailure is an interface to extract request failure information from an Error.
type TxRequestFailure interface {
	awserr.RequestFailure
	CancellationReasons() []dynamodb.CancellationReason
}

// TxAwareUnmarshalError unmarshals an error response for a JSON RPC service.
// This is exactly same as jsonrpc.UnmarshalError, except for attempt to parse CancellationReasons
func TxAwareUnmarshalError(req *request.Request) {
	defer req.HTTPResponse.Body.Close()

	var jsonErr jsonTxErrorResponse
	err := json.NewDecoder(req.HTTPResponse.Body).Decode(&jsonErr)
	if err == io.EOF {
		req.Error = awserr.NewRequestFailure(
			awserr.New("SerializationError", req.HTTPResponse.Status, nil),
			req.HTTPResponse.StatusCode,
			req.RequestID,
		)
		return
	} else if err != nil {
		req.Error = awserr.NewRequestFailure(
			awserr.New("SerializationError", "failed decoding JSON RPC error response", err),
			req.HTTPResponse.StatusCode,
			req.RequestID,
		)
		return
	}

	codes := strings.SplitN(jsonErr.Code, "#", 2)
	req.Error = newTxRequestError(
		awserr.New(codes[len(codes)-1], jsonErr.Message, nil),
		req.HTTPResponse.StatusCode,
		req.RequestID,
		jsonErr.CancellationReasons,
	)
}

type jsonTxErrorResponse struct {
	Code                string                        `json:"__type"`
	Message             string                        `json:"message"`
	CancellationReasons []dynamodb.CancellationReason `json:"CancellationReasons"`
}

// So that the Error interface type can be included as an anonymous field
// in the requestError struct and not conflict with the error.Error() method.
type awsError awserr.Error

// A TxRequestError wraps a request or service error.
// TxRequestError is awserr.requestError with additional cancellationReasons field
type txRequestError struct {
	awsError
	statusCode          int
	requestID           string
	cancellationReasons []dynamodb.CancellationReason
}

func newTxRequestError(err awserr.Error, statusCode int, requestID string, cancellationReasons []dynamodb.CancellationReason) TxRequestFailure {
	return &txRequestError{
		awsError:            err,
		statusCode:          statusCode,
		requestID:           requestID,
		cancellationReasons: cancellationReasons,
	}
}

func (r txRequestError) Error() string {
	extra := fmt.Sprintf("status code: %d, request id: %s",
		r.statusCode, r.requestID)
	return awserr.SprintError(r.Code(), r.Message(), extra, r.OrigErr())
}

func (r txRequestError) String() string {
	return r.Error()
}

func (r txRequestError) StatusCode() int {
	return r.statusCode
}

func (r txRequestError) RequestID() string {
	return r.requestID
}

func (r txRequestError) CancellationReasons() []dynamodb.CancellationReason {
	return r.cancellationReasons
}

